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
	import { renderMarkdown } from '$lib/utils/markdown';
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
		// If we have a streaming message and this is the final assistant message, replace it
		if (streamingMessageId && message.role === 'assistant') {
			// Remove the streaming message and replace with the final one
			messages = messages
				.filter((m) => m.id !== streamingMessageId && !m.id?.startsWith('__streaming_'))
				.concat(message)
				.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());
			streamingMessageId = null;
			streamingContent = '';
			return;
		}

		// Check if message already exists
		const exists = messages.some((m) => m.id === message.id);
		if (exists) {
			// Update existing message
			messages = messages.map((m) => (m.id === message.id ? message : m));
		} else {
			// Add new message (but filter out any streaming messages first)
			messages = messages
				.filter((m) => !m.id?.startsWith('__streaming_'))
				.concat(message)
				.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());
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
		// But filter out any streaming messages that might have been created
		await loadMessages();
		// Clean up any streaming messages after reload
		messages = messages.filter((m) => !m.id?.startsWith('__streaming_'));
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

<div class="chat-container">
	<!-- Header -->
	<div class="chat-header">
		<div class="chat-header-content">
			<MessageSquare class="chat-icon" />
			<h3 class="chat-title">{title}</h3>
		</div>
		{#if onClose}
			<button onclick={onClose} class="chat-close-button">
				<X class="icon" />
			</button>
		{/if}
	</div>

	<!-- Messages -->
	<div class="messages-container">
		{#if loading}
			<div class="loading-state">
				<div class="spinner"></div>
				<p class="loading-text">Loading chat...</p>
			</div>
		{:else if error}
			<div class="error-state">
				<p class="error-text">{error}</p>
			</div>
		{:else if !conversation}
			<div class="empty-state">
				<MessageSquare class="empty-icon" />
				<p class="empty-text">No conversation yet</p>
				<button onclick={createNewConversation} class="start-chat-button">
					Start Chat
				</button>
			</div>
		{:else if messages.length === 0}
			<div class="empty-messages">
				<p class="empty-text">No messages yet. Start the conversation!</p>
			</div>
		{:else}
			{#each messages as message}
				<div class="message-wrapper {message.role === 'user' ? 'message-user' : 'message-assistant'}">
					<div class="message-bubble {message.role === 'user' ? 'bubble-user' : 'bubble-assistant'}">
						{#if message.role === 'assistant'}
							<div class="message-content markdown-content">
								{@html renderMarkdown(message.content)}
							</div>
						{:else}
							<p class="message-content">{message.content}</p>
						{/if}
						<div class="message-footer">
							<span class="message-time">
								{new Date(message.createdAt).toLocaleTimeString()}
							</span>
							{#if message.id?.startsWith('__streaming_') || message.id === streamingMessageId}
								<span class="streaming-indicator">●</span>
							{/if}
						</div>
					</div>
				</div>
			{/each}
			{#if sending}
				<div class="message-wrapper message-assistant">
					<div class="message-bubble bubble-assistant">
						<div class="typing-indicator">
							<span></span>
							<span></span>
							<span></span>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>

	<!-- Input -->
	{#if conversation}
		<div class="input-container">
			<div class="input-wrapper">
				<textarea
					bind:value={messageContent}
					onkeypress={handleKeyPress}
					placeholder="Type your message..."
					rows="2"
					class="message-input"
					disabled={sending}
				></textarea>
				<button
					onclick={handleSend}
					disabled={!messageContent.trim() || sending}
					class="send-button"
				>
					<Send class="send-icon" />
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	.chat-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	.chat-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--spacing-md);
		border-bottom: 1px solid var(--color-border);
		background-color: var(--color-bg-secondary);
	}

	.chat-header-content {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
	}

	.chat-icon {
		width: 1.25rem;
		height: 1.25rem;
		color: var(--color-primary);
	}

	.chat-title {
		margin: 0;
		font-size: var(--font-size-lg);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.chat-close-button {
		background: none;
		border: none;
		cursor: pointer;
		padding: var(--spacing-xs);
		color: var(--color-text-tertiary);
		transition: color var(--transition-base);
	}

	.chat-close-button:hover {
		color: var(--color-text-primary);
	}

	.chat-close-button .icon {
		width: 1.25rem;
		height: 1.25rem;
	}

	.messages-container {
		flex: 1;
		overflow-y: auto;
		padding: var(--spacing-md);
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	/* Message Styles */
	.message-wrapper {
		display: flex;
		width: 100%;
	}

	.message-user {
		justify-content: flex-end;
	}

	.message-assistant {
		justify-content: flex-start;
	}

	.message-bubble {
		max-width: 80%;
		padding: var(--spacing-sm) var(--spacing-md);
		border-radius: var(--radius-lg);
		word-wrap: break-word;
	}

	.bubble-user {
		background-color: var(--color-primary);
		color: white;
		border-bottom-right-radius: var(--radius-sm);
	}

	.bubble-assistant {
		background-color: var(--color-bg-secondary);
		color: var(--color-text-primary);
		border: 1px solid var(--color-border);
		border-bottom-left-radius: var(--radius-sm);
	}

	.message-content {
		margin: 0;
		font-size: var(--font-size-sm);
		line-height: 1.5;
		white-space: pre-wrap;
		word-wrap: break-word;
	}

	/* Markdown Content Styles */
	.markdown-content {
		white-space: normal;
	}

	.markdown-content :global(p) {
		margin: 0 0 var(--spacing-sm) 0;
		line-height: 1.6;
	}

	.markdown-content :global(p:last-child) {
		margin-bottom: 0;
	}

	.markdown-content :global(h1),
	.markdown-content :global(h2),
	.markdown-content :global(h3),
	.markdown-content :global(h4),
	.markdown-content :global(h5),
	.markdown-content :global(h6) {
		margin: var(--spacing-md) 0 var(--spacing-sm) 0;
		font-weight: var(--font-weight-semibold);
		line-height: 1.3;
	}

	.markdown-content :global(h1:first-child),
	.markdown-content :global(h2:first-child),
	.markdown-content :global(h3:first-child),
	.markdown-content :global(h4:first-child),
	.markdown-content :global(h5:first-child),
	.markdown-content :global(h6:first-child) {
		margin-top: 0;
	}

	.markdown-content :global(h1) {
		font-size: var(--font-size-xl);
	}

	.markdown-content :global(h2) {
		font-size: var(--font-size-lg);
	}

	.markdown-content :global(h3) {
		font-size: var(--font-size-base);
	}

	.markdown-content :global(h4),
	.markdown-content :global(h5),
	.markdown-content :global(h6) {
		font-size: var(--font-size-sm);
	}

	.markdown-content :global(strong),
	.markdown-content :global(b) {
		font-weight: var(--font-weight-semibold);
	}

	.markdown-content :global(em),
	.markdown-content :global(i) {
		font-style: italic;
	}

	.markdown-content :global(code) {
		background-color: var(--color-bg-tertiary);
		padding: 0.125rem 0.25rem;
		border-radius: var(--radius-sm);
		font-family: 'Courier New', monospace;
		font-size: 0.875em;
	}

	.markdown-content :global(pre) {
		background-color: var(--color-bg-tertiary);
		padding: var(--spacing-sm);
		border-radius: var(--radius-md);
		overflow-x: auto;
		margin: var(--spacing-sm) 0;
		border: 1px solid var(--color-border);
	}

	.markdown-content :global(pre code) {
		background-color: transparent;
		padding: 0;
		border-radius: 0;
	}

	.markdown-content :global(ul),
	.markdown-content :global(ol) {
		margin: var(--spacing-sm) 0;
		padding-left: var(--spacing-lg);
	}

	.markdown-content :global(li) {
		margin: var(--spacing-xs) 0;
		line-height: 1.6;
	}

	.markdown-content :global(ul) {
		list-style-type: disc;
	}

	.markdown-content :global(ol) {
		list-style-type: decimal;
	}

	.markdown-content :global(blockquote) {
		border-left: 3px solid var(--color-border);
		padding-left: var(--spacing-md);
		margin: var(--spacing-sm) 0;
		opacity: 0.8;
		font-style: italic;
	}

	.markdown-content :global(hr) {
		border: none;
		border-top: 1px solid var(--color-border);
		margin: var(--spacing-md) 0;
	}

	.markdown-content :global(a) {
		color: var(--color-primary);
		text-decoration: underline;
		transition: color var(--transition-base);
	}

	.markdown-content :global(a:hover) {
		color: var(--color-primary-hover);
	}

	.markdown-content :global(table) {
		width: 100%;
		border-collapse: collapse;
		margin: var(--spacing-sm) 0;
		font-size: var(--font-size-xs);
	}

	.markdown-content :global(th),
	.markdown-content :global(td) {
		padding: var(--spacing-xs) var(--spacing-sm);
		border: 1px solid var(--color-border);
		text-align: left;
	}

	.markdown-content :global(th) {
		background-color: var(--color-bg-tertiary);
		font-weight: var(--font-weight-semibold);
	}

	/* Adjust colors for user messages (white text) */
	.bubble-user .markdown-content :global(*) {
		color: inherit;
	}

	.bubble-user .markdown-content :global(code),
	.bubble-user .markdown-content :global(pre) {
		background-color: rgba(255, 255, 255, 0.2);
		color: inherit;
	}

	.bubble-user .markdown-content :global(a) {
		color: rgba(255, 255, 255, 0.9);
	}

	.bubble-user .markdown-content :global(a:hover) {
		color: white;
	}

	.message-footer {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		margin-top: var(--spacing-xs);
	}

	.message-time {
		font-size: var(--font-size-xs);
		opacity: 0.7;
	}

	.streaming-indicator {
		font-size: var(--font-size-xs);
		opacity: 0.5;
		animation: pulse 1.5s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 0.5; }
		50% { opacity: 1; }
	}

	.typing-indicator {
		display: flex;
		gap: var(--spacing-xs);
		padding: var(--spacing-xs) 0;
	}

	.typing-indicator span {
		width: 0.5rem;
		height: 0.5rem;
		background-color: var(--color-text-secondary);
		border-radius: 50%;
		animation: bounce 1.4s ease-in-out infinite;
	}

	.typing-indicator span:nth-child(2) {
		animation-delay: 0.2s;
	}

	.typing-indicator span:nth-child(3) {
		animation-delay: 0.4s;
	}

	@keyframes bounce {
		0%, 80%, 100% { transform: translateY(0); }
		40% { transform: translateY(-0.5rem); }
	}

	/* Loading and Empty States */
	.loading-state,
	.empty-state,
	.empty-messages {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		text-align: center;
	}

	.spinner {
		width: 2rem;
		height: 2rem;
		border: 2px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.loading-text,
	.empty-text {
		margin-top: var(--spacing-md);
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
	}

	.empty-icon {
		width: 3rem;
		height: 3rem;
		color: var(--color-text-tertiary);
		margin-bottom: var(--spacing-md);
	}

	.start-chat-button {
		margin-top: var(--spacing-md);
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-primary);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
		cursor: pointer;
		transition: background-color var(--transition-base);
	}

	.start-chat-button:hover {
		background-color: var(--color-primary-hover);
	}

	.error-state {
		background-color: var(--color-danger-bg);
		border: 1px solid var(--color-danger);
		border-radius: var(--radius-lg);
		padding: var(--spacing-md);
	}

	.error-text {
		margin: 0;
		font-size: var(--font-size-sm);
		color: var(--color-danger);
	}

	/* Input Styles */
	.input-container {
		padding: var(--spacing-md);
		border-top: 1px solid var(--color-border);
		background-color: var(--color-bg-secondary);
	}

	.input-wrapper {
		display: flex;
		gap: var(--spacing-sm);
		align-items: flex-end;
	}

	.message-input {
		flex: 1;
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		font-size: var(--font-size-sm);
		font-family: inherit;
		resize: none;
		transition: border-color var(--transition-base);
	}

	.message-input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.message-input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.send-button {
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-primary);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: background-color var(--transition-base);
	}

	.send-button:hover:not(:disabled) {
		background-color: var(--color-primary-hover);
	}

	.send-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.send-icon {
		width: 1rem;
		height: 1rem;
	}
</style>

