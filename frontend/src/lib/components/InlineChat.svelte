<script lang="ts">
	import { onMount } from 'svelte';
	import {
		createConversation,
		listMessages,
		appendMessage,
		type Conversation,
		type Message,
		type CreateConversationInput
	} from '$lib/api/chats';
	import { toastError, toastSuccess } from '$lib/utils/toast';
	import { MessageSquare, Send, X } from 'lucide-svelte';

	interface Props {
		jobApplicationId?: string;
		title?: string;
		description?: string;
		onClose?: () => void;
	}

	let { jobApplicationId, title = 'Chat', description = '', onClose }: Props = $props();

	let conversation: Conversation | null = $state(null);
	let messages: Message[] = $state([]);
	let messageContent = $state('');
	let loading = $state(false);
	let sending = $state(false);
	let error: string | null = $state(null);

	onMount(async () => {
		await initializeChat();
	});

	async function initializeChat() {
		if (!jobApplicationId) return;

		loading = true;
		error = null;
		try {
			// Try to find existing conversation for this job application
			const { searchConversations } = await import('$lib/api/chats');
			const conversations = await searchConversations(undefined, false, jobApplicationId);
			
			if (conversations.length > 0) {
				conversation = conversations[0];
				await loadMessages();
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
		} catch (err) {
			console.error('Error loading messages:', err);
			toastError('Failed to load messages');
		}
	}

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

			// Reload messages to get the AI response
			await new Promise((resolve) => setTimeout(resolve, 1000)); // Wait a bit for AI response
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
						<p class="text-xs mt-1 opacity-70">
							{new Date(message.createdAt).toLocaleTimeString()}
						</p>
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

