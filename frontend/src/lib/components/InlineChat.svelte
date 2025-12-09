<script lang="ts">
	import { onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import { get } from 'svelte/store';
	import {
		createConversation,
		listMessages,
		appendMessage,
		type Conversation,
		type Message,
		type CreateConversationInput
	} from '$lib/api/chats';
	import { authStore } from '$lib';
	import { toastError, toastSuccess } from '$lib/utils/toast';
	import { MessageSquare, Send, X } from 'lucide-svelte';

	interface Props {
		jobApplicationId?: string;
		conversationId?: string;
		title?: string;
		description?: string;
		onClose?: () => void;
	}

	let { jobApplicationId, conversationId, title = 'Chat', description = '', onClose }: Props = $props();

	let conversation: Conversation | null = $state(null);
	let messages: Message[] = $state([]);
	let messageContent = $state('');
	let loading = $state(false);
	let sending = $state(false);
	let error: string | null = $state(null);
	let lastInitializedConversationId: string | undefined = $state(undefined);
	let lastInitializedJobApplicationId: string | undefined = $state(undefined);

	// WebSocket connection state
	let ws: WebSocket | null = $state(null);
	let streamingMessageId: string | null = $state(null);
	let streamingContent: string = $state('');

	onDestroy(() => {
		disconnectStream();
	});

	// Initialize chat when component mounts or when conversationId/jobApplicationId changes
	$effect(() => {
		// Reference both props so Svelte tracks them as dependencies
		const convId = conversationId;
		const appId = jobApplicationId;
		
		// Only initialize if:
		// 1. We have a jobApplicationId
		// 2. Either conversationId or jobApplicationId has changed since last initialization
		if (appId && (convId !== lastInitializedConversationId || appId !== lastInitializedJobApplicationId)) {
			lastInitializedConversationId = convId;
			lastInitializedJobApplicationId = appId;
			initializeChat();
		}
	});

	async function initializeChat() {
		if (!jobApplicationId) return;

		// Prevent multiple simultaneous initializations
		if (loading) return;

		loading = true;
		error = null;
		try {
			// If conversationId is provided, use it directly
			if (conversationId) {
				console.log('Loading conversation by ID:', conversationId);
				const { getConversation } = await import('$lib/api/chats');
				conversation = await getConversation(conversationId);
				await loadMessages();
			} else {
				// Try to find existing conversation for this job application
				console.log('Searching for conversations for jobApplicationId:', jobApplicationId);
				const { searchConversations } = await import('$lib/api/chats');
				
				// Try without archived first
				let conversations = await searchConversations(undefined, false, jobApplicationId);
				console.log('Found conversations (not archived):', conversations);
				
				// If none found, try including archived
				if (conversations.length === 0) {
					conversations = await searchConversations(undefined, true, jobApplicationId);
					console.log('Found conversations (including archived):', conversations);
				}
				
				if (conversations.length > 0) {
					// Use the most recent conversation
					conversation = conversations.sort((a, b) => 
						new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
					)[0];
					console.log('Selected conversation:', conversation.id);
					await loadMessages();
				} else {
					console.log('No conversations found for jobApplicationId:', jobApplicationId);
				}
			}
			// Don't auto-create - let user click "Start Chat" button
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to initialize chat';
			console.error('Error initializing chat:', err);
			toastError('Failed to initialize chat');
		} finally {
			loading = false;
		}
	}

	async function createNewConversation() {
		if (!jobApplicationId) return;

		loading = true;
		const input: CreateConversationInput = {
			title,
			description: description || `Chat about ${title}`,
			jobApplicationId
		};

		try {
			conversation = await createConversation(input);
			toastSuccess('Chat started');
			await loadMessages();
		} catch (err) {
			console.error('Error creating conversation:', err);
			toastError('Failed to create chat');
		} finally {
			loading = false;
		}
	}

	async function loadMessages() {
		if (!conversation) return;

		try {
			messages = await listMessages(conversation.id);
			// Connect to WebSocket after loading messages
			connectToStream(conversation.id);
		} catch (err) {
			console.error('Error loading messages:', err);
			toastError('Failed to load messages');
		}
	}

	const getWebSocketURL = (conversationId: string): string | null => {
		if (!browser) return null;
		
		const { token } = get(authStore);
		if (!token) {
			console.warn('No auth token available for WebSocket connection');
			return null;
		}

		// Get base URL from environment or default to localhost
		const baseURL = (import.meta.env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '');
		// Convert http/https to ws/wss
		const wsProtocol = baseURL.startsWith('https') ? 'wss' : 'ws';
		const wsBase = baseURL.replace(/^https?/, wsProtocol);
		
		return `${wsBase}/api/chats/conversations/${conversationId}/stream?token=${encodeURIComponent(token)}`;
	};

	const connectToStream = (conversationId: string) => {
		if (!browser) return;
		
		// Disconnect existing connection
		disconnectStream();

		const wsUrl = getWebSocketURL(conversationId);
		if (!wsUrl) {
			console.warn('Cannot connect to WebSocket: invalid URL or missing token');
			return;
		}

		try {
			ws = new WebSocket(wsUrl);

			ws.onopen = () => {
				console.log('WebSocket connected for conversation:', conversationId);
			};

			ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data);
					
					// Handle delta events (streaming AI response)
					if (data.type === 'delta' && data.delta) {
						handleStreamingDelta(data.delta);
					}
					// Handle full message events
					else if (data.id && data.conversation_id && data.role && data.content) {
						handleFullMessage(data);
					}
				} catch (err) {
					console.error('Error parsing WebSocket message:', err, event.data);
				}
			};

			ws.onerror = (error) => {
				console.error('WebSocket error:', error);
			};

			ws.onclose = () => {
				console.log('WebSocket disconnected');
				ws = null;
			};
		} catch (error) {
			console.error('Failed to create WebSocket connection:', error);
		}
	};

	const handleStreamingDelta = (delta: string) => {
		// Find or create the streaming assistant message
		let streamingMsg = messages.find(
			(m) => m.id === streamingMessageId || (m.id?.startsWith('__streaming_') && m.role === 'assistant')
		);

		if (!streamingMsg) {
			// Create a new streaming message
			streamingMessageId = `__streaming_${Date.now()}`;
			streamingContent = delta;
			streamingMsg = {
				id: streamingMessageId,
				conversationId: conversation!.id,
				role: 'assistant',
				content: delta,
				createdAt: new Date().toISOString()
			};
			messages = [...messages, streamingMsg];
		} else {
			// Update existing streaming message
			streamingMessageId = streamingMsg.id;
			streamingContent += delta;
			streamingMsg.content = streamingContent;
			messages = messages.map((m) => (m.id === streamingMsg.id ? streamingMsg : m));
		}
	};

	const handleFullMessage = (message: any) => {
		// Clear streaming state if this is the final message
		if (streamingMessageId && message.role === 'assistant') {
			streamingMessageId = null;
			streamingContent = '';
		}

		// Check if message already exists
		const exists = messages.some((m) => m.id === message.id);
		if (exists) {
			// Update existing message
			messages = messages.map((m) => (m.id === message.id ? message : m));
		} else {
			// Add new message
			messages = [...messages, message].sort(
				(a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
			);
		}
	};

	const disconnectStream = () => {
		if (ws) {
			ws.close();
			ws = null;
		}
		streamingMessageId = null;
		streamingContent = '';
	};

	async function handleSend() {
		if (!conversation || !messageContent.trim() || sending) return;

		const content = messageContent.trim();
		messageContent = '';
		sending = true;

		try {
		await appendMessage(conversation.id, {
			role: 'user',
			content,
			generate_reply: true
		});

		// Reload messages to get the user message (AI response will come via WebSocket)
		await loadMessages();
		} catch (err) {
			console.error('Error sending message:', err);
			toastError('Failed to send message');
			messageContent = content; // Restore message on error
		} finally {
			sending = false;
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			handleSend();
		}
	}
</script>

<div class="flex flex-col h-full bg-white dark:bg-gray-800 rounded-lg shadow-lg">
	<!-- Header -->
	<div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
		<div class="flex items-center gap-2">
			<MessageSquare class="w-5 h-5 text-blue-600 dark:text-blue-400" />
			<h3 class="text-lg font-semibold text-gray-900 dark:text-white">{title}</h3>
		</div>
		{#if onClose}
			<button
				onclick={onClose}
				class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
			>
				<X class="w-5 h-5" />
			</button>
		{/if}
	</div>

	<!-- Messages -->
	<div class="flex-1 overflow-y-auto p-4 space-y-4">
		{#if loading}
			<div class="flex items-center justify-center h-full">
				<div class="text-center">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
					<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">Loading chat...</p>
				</div>
			</div>
		{:else if error}
			<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
				<p class="text-sm text-red-800 dark:text-red-200">{error}</p>
			</div>
		{:else if !conversation}
			<div class="flex items-center justify-center h-full">
				<div class="text-center">
					<MessageSquare class="w-12 h-12 text-gray-400 mx-auto mb-2" />
					<p class="text-gray-600 dark:text-gray-400">No conversation yet</p>
					<button
						onclick={createNewConversation}
						class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
					>
						Start Chat
					</button>
				</div>
			</div>
		{:else if messages.length === 0}
			<div class="flex items-center justify-center h-full">
				<p class="text-gray-600 dark:text-gray-400">No messages yet. Start the conversation!</p>
			</div>
		{:else}
			{#each messages as message}
				<div
					class="flex {message.role === 'user' ? 'justify-end' : 'justify-start'}"
				>
					<div
						class="max-w-[80%] rounded-lg px-4 py-2 {message.role === 'user'
							? 'bg-blue-600 text-white'
							: 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white'}"
					>
						<p class="text-sm whitespace-pre-wrap">{message.content}</p>
						<div class="flex items-center gap-2 mt-1">
							<p class="text-xs opacity-70">
								{new Date(message.createdAt).toLocaleTimeString()}
							</p>
							{#if message.id?.startsWith('__streaming_') || message.id === streamingMessageId}
								<span class="text-xs opacity-50 animate-pulse">●</span>
							{/if}
						</div>
					</div>
				</div>
			{/each}
			{#if sending}
				<div class="flex justify-start">
					<div class="bg-gray-100 dark:bg-gray-700 rounded-lg px-4 py-2">
						<div class="flex items-center gap-2">
							<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
							<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
							<div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.4s"></div>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>

	<!-- Input -->
	{#if conversation}
		<div class="p-4 border-t border-gray-200 dark:border-gray-700">
			<div class="flex gap-2">
				<textarea
					bind:value={messageContent}
					onkeypress={handleKeyPress}
					placeholder="Type your message..."
					rows="2"
					class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white resize-none"
					disabled={sending}
				></textarea>
				<button
					onclick={handleSend}
					disabled={!messageContent.trim() || sending}
					class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
				>
					<Send class="w-4 h-4" />
				</button>
			</div>
		</div>
	{/if}
</div>

