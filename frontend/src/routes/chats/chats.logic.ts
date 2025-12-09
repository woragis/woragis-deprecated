import { browser } from '$app/environment';
import { onDestroy } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import type { AppendMessageInput } from '$lib/api/chats';
// WebSocket implementation replaces CHATS_STREAM_BASE
import { apiClient } from '@clients/apiClient';
import type { ChatConversation, Idea, Project, UUID } from '$lib/api/types';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';
import { authStore } from '$lib';
import {
	useAppendMessageMutation,
	useArchiveConversationsMutation,
	useConversationAssignmentsQuery,
	useConversationMessagesQuery,
	useConversationsQuery,
	useConversationTranscriptsQuery,
	useCreateConversationMutation,
	useDeleteConversationsMutation,
	useRestoreConversationsMutation,
	useShareTranscriptMutation,
	type ConversationsQueryOptions
} from '@hooks/chats';
import { useProjectsListQuery } from '@hooks/projects';
import { useIdeasReferenceQuery } from '@hooks/ideas';

export interface CreateConversationForm {
	title: string;
	description: string;
	ideaId: string;
	projectId: string;
}

const defaultCreateForm = (): CreateConversationForm => ({
	title: '',
	description: '',
	ideaId: '',
	projectId: ''
});

const normalize = (value: string) => value.toLowerCase().trim();

export function createChatsLogic() {
	const queryClient = useQueryClient();

	const selectedConversation = writable<ChatConversation | null>(null);
	const selectedConversationIds = writable<Set<UUID>>(new Set());

	const conversationsFilters = writable<ConversationsQueryOptions>({
		search: '',
		includeArchived: false,
		jobApplicationId: undefined,
		enabled: browser
	});

	const conversationsQuery = useConversationsQuery(conversationsFilters);
	const selectedConversationId = derived(
		selectedConversation,
		($conversation) => $conversation?.id ?? null
	);
	const messagesQuery = useConversationMessagesQuery(selectedConversationId);
	const transcriptsQuery = useConversationTranscriptsQuery(selectedConversationId);
	const assignmentsQuery = useConversationAssignmentsQuery(selectedConversationId);

	const projectsQuery = useProjectsListQuery({ enabled: browser });
	const ideasQuery = useIdeasReferenceQuery(browser);

	const createConversationMutation = useCreateConversationMutation();
	const archiveMutation = useArchiveConversationsMutation();
	const restoreMutation = useRestoreConversationsMutation();
	const deleteMutation = useDeleteConversationsMutation();
	const appendMessageMutation = useAppendMessageMutation();
	const shareTranscriptMutation = useShareTranscriptMutation();

	const createError = writable('');
	const ideaQuery = writable('');
	const projectQuery = writable('');
	const ideaDropdownOpen = writable(false);
	const projectDropdownOpen = writable(false);
	const createForm = writable<CreateConversationForm>(defaultCreateForm());
	const updateCreateFormField = <K extends keyof CreateConversationForm>(
		field: K,
		value: CreateConversationForm[K]
	) => {
		createForm.update((current) => ({ ...current, [field]: value }));
	};


	const searchInput = writable('');
	const includeArchived = writable(false);
	const bulkStatus = writable('');
	const composeContent = writable('');
	const composeRole = writable<'user' | 'assistant'>('user');
	const generateReply = writable(false);
	const transcriptStatus = writable('');
	const showCreateModal = writable(false);

	// WebSocket connection state
	let ws: WebSocket | null = null;
	let streamingMessageId: string | null = null;
	let streamingContent: string = '';
	let reconnectAttempts = 0;
	const maxReconnectAttempts = 5;
	const reconnectDelay = 1000;

	// Per-message AI controls
	const provider = writable<string>('openai');
	const model = writable<string>('');
	const agent = writable<string>('startup');
	const autoAgent = writable<boolean>(false);
	const isSending = writable<boolean>(false);

	const conversationsUnsubscribe = conversationsQuery.subscribe((result) => {
		const data = result.data ?? [];
		const current = get(selectedConversation);
		if (current) {
			const match = data.find((item) => item.id === current.id);
			if (match) {
				selectedConversation.set(match);
			} else {
				selectedConversation.set(null);
				disconnectStream();
			}
		}
		selectedConversationIds.update(
			(set) => new Set([...set].filter((id) => data.some((item) => item.id === id)))
		);
	});

	const getWebSocketURL = (conversationId: UUID): string | null => {
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

	const connectToStream = (conversationId: UUID) => {
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
				reconnectAttempts = 0;
			};

			ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data);
					
					// Handle delta events (streaming AI response)
					if (data.type === 'delta' && data.delta) {
						handleStreamingDelta(conversationId, data.delta);
					}
					// Handle full message events
					else if (data.id && data.conversation_id && data.role && data.content) {
						handleFullMessage(conversationId, data);
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
				
				// Attempt reconnection if conversation is still selected
				const current = get(selectedConversation);
				if (current && current.id === conversationId && reconnectAttempts < maxReconnectAttempts) {
					reconnectAttempts++;
					setTimeout(() => {
						console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...`);
						connectToStream(conversationId);
					}, reconnectDelay * reconnectAttempts);
				}
			};
		} catch (error) {
			console.error('Failed to create WebSocket connection:', error);
		}
	};

	const handleStreamingDelta = (conversationId: UUID, delta: string) => {
		// Create or update a streaming message in the cache
		queryClient.setQueryData(
			['chats', 'conversation', conversationId, 'messages'],
			(previous: any) => {
				const messages: any[] = Array.isArray(previous) ? previous : [];
				
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
						conversation_id: conversationId,
						role: 'assistant',
						content: delta,
						created_at: new Date().toISOString(),
						_streaming: true
					};
					return [...messages, streamingMsg];
				} else {
					// Update existing streaming message
					streamingMessageId = streamingMsg.id;
					streamingContent += delta;
					streamingMsg.content = streamingContent;
					streamingMsg._streaming = true;
					return messages.map((m) => (m.id === streamingMsg.id ? streamingMsg : m));
				}
			}
		);
	};

	const handleFullMessage = (conversationId: UUID, message: any) => {
		// Clear streaming state if this is the final message
		if (streamingMessageId && message.role === 'assistant') {
			streamingMessageId = null;
			streamingContent = '';
		}

		queryClient.setQueryData(
			['chats', 'conversation', conversationId, 'messages'],
			(previous: any) => {
				const messages: any[] = Array.isArray(previous) ? previous : [];
				
				// Remove any streaming placeholder
				const withoutStreaming = messages.filter((m) => m.id !== streamingMessageId && !m._streaming);
				
				// Check if message already exists
				const exists = withoutStreaming.some((m) => m.id === message.id);
				if (exists) {
					// Update existing message
					return withoutStreaming.map((m) => (m.id === message.id ? { ...message, _streaming: false } : m));
				} else {
					// Add new message
					return [...withoutStreaming, { ...message, _streaming: false }].sort(
						(a: any, b: any) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
					);
				}
			}
		);
	};

	const disconnectStream = () => {
		if (ws) {
			ws.close();
			ws = null;
		}
		streamingMessageId = null;
		streamingContent = '';
		reconnectAttempts = 0;
	};

	// Polling disabled - using WebSocket streaming instead

	const setSearchInput = (value: string) => searchInput.set(value);
	const setIncludeArchived = (value: boolean) => includeArchived.set(value);

	const applyFilters = async (jobApplicationId?: string) => {
		const nextSearch = get(searchInput).trim();
		const include = get(includeArchived);
		conversationsFilters.set({
			search: nextSearch,
			includeArchived: include,
			jobApplicationId: jobApplicationId,
			enabled: true
		});
		await queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
	};

	const refreshConversations = () => {
		void queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
	};

	const toggleConversationSelection = (id: UUID) => {
		selectedConversationIds.update((set) => {
			const next = new Set(set);
			if (next.has(id)) {
				next.delete(id);
			} else {
				next.add(id);
			}
			return next;
		});
	};

	const clearSelections = () => {
		selectedConversationIds.set(new Set());
	};

	const selectConversation = (conversation: ChatConversation) => {
		if (get(selectedConversation)?.id === conversation.id) {
			return;
		}
		selectedConversation.set(conversation);
		connectToStream(conversation.id);
		void queryClient.invalidateQueries({ queryKey: ['chats', 'conversation', conversation.id, 'messages'] });
		void queryClient.invalidateQueries({ queryKey: ['chats', 'conversation', conversation.id, 'transcripts'] });
		void queryClient.invalidateQueries({ queryKey: ['chats', 'conversation', conversation.id, 'assignments'] });
	};

	const handleSendMessage = async (event: SubmitEvent) => {
		event.preventDefault();
		const conversation = get(selectedConversation);
		const content = get(composeContent).trim();
		const role = get(composeRole);
		const shouldReply = get(generateReply);
		const currentProvider = get(provider);
		const currentModel = get(model);
		const currentAgent = get(agent);
		const useAutoAgent = get(autoAgent);
		if (!conversation || !content) return;
		const payload: AppendMessageInput = {
			role,
			content,
			generate_reply: shouldReply,
			provider: currentProvider,
			model: currentModel
		};
		// Direct POST without TanStack mutation to avoid caching layers
		const body: any = {
			...payload
		};
		if (!useAutoAgent) {
			body.agent = currentAgent;
		}
		// Optimistically add the user message
		const nowIso = new Date().toISOString();
		const optimisticId = `__local_${Math.random().toString(36).slice(2)}`;
		isSending.set(true);
		queryClient.setQueryData(
			['chats', 'conversation', conversation.id, 'messages'],
			(previous: any) => {
				const base: any[] = Array.isArray(previous) ? previous : [];
				const optimistic = {
					id: optimisticId,
					conversation_id: conversation.id,
					role,
					content,
					created_at: nowIso
				};
				return [...base, optimistic];
			}
		);

		let newMessages: any[] = [];
		try {
			const response = await apiClient.post<{ success: boolean; data: any[] }>(
				`/chats/conversations/${conversation.id}/messages`,
				body
			);
			if (response.data?.success === false) {
				throw new Error('Request failed');
			}
			newMessages = (response.data?.data as any[]) ?? [];
		} catch (err) {
			toastError('Failed to send message.');
			// Rollback optimistic user message
			queryClient.setQueryData(
				['chats', 'conversation', conversation.id, 'messages'],
				(previous: any) => {
					const base: any[] = Array.isArray(previous) ? previous : [];
					return base.filter((m) => m?.id !== optimisticId);
				}
			);
			isSending.set(false);
			return;
		}
		queryClient.setQueryData(
			['chats', 'conversation', conversation.id, 'messages'],
			(previous: any) => {
				const baseMessages: any[] = Array.isArray(previous) ? previous : [];
				// Remove optimistic user message if backend returned it
				const withoutOptimistic = baseMessages.filter((m) => m?.id !== optimisticId);
				const nextMessages = [...withoutOptimistic, ...newMessages];
				return nextMessages.sort(
					(a: any, b: any) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
				);
			}
		);
		// Refresh conversations list using the correct key
		queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
		composeContent.set('');
		generateReply.set(false);
		isSending.set(false);

		// AI reply will arrive via WebSocket stream - no polling needed
	};

	const handleShareTranscript = async () => {
		const conversation = get(selectedConversation);
		if (!conversation) {
			toastError('Select a conversation first.');
			return;
		}
		transcriptStatus.set('Generating transcript...');
		toastInfo('Generating transcript...');
		try {
			const transcript = await get(shareTranscriptMutation).mutateAsync(conversation.id);
			queryClient.setQueryData(
				['chats', 'conversation', conversation.id, 'transcripts'],
				(previous: any) => [transcript, ...(previous ?? [])]
			);
			transcriptStatus.set(`Share code ${(transcript as { share_code: string }).share_code}`);
			setTimeout(() => transcriptStatus.set(''), 4000);
			toastSuccess('Transcript generated.');
		} catch (error) {
			const message = getApiErrorMessage(error, 'Failed to generate transcript.');
			transcriptStatus.set(message);
			setTimeout(() => transcriptStatus.set(''), 4000);
			toastError(message);
		}
	};

	const applyBulkAction = async (action: 'archive' | 'restore' | 'delete') => {
		const ids = Array.from(get(selectedConversationIds));
		if (ids.length === 0) {
			toastError('Select at least one conversation.');
			return;
		}
		bulkStatus.set('Working...');
		try {
			if (action === 'archive') {
				await get(archiveMutation).mutateAsync(ids);
				bulkStatus.set('Archived conversations');
				toastSuccess('Archived conversations.');
			} else if (action === 'restore') {
				await get(restoreMutation).mutateAsync(ids);
				bulkStatus.set('Restored conversations');
				toastSuccess('Restored conversations.');
			} else {
			await get(deleteMutation).mutateAsync(ids);
				bulkStatus.set('Deleted conversations');
				toastSuccess('Deleted conversations.');
			}
			clearSelections();
		await queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
		} catch {
			bulkStatus.set('Something went wrong.');
			toastError('Unable to complete bulk action.');
		} finally {
			setTimeout(() => bulkStatus.set(''), 3000);
		}
	};

	const openCreateModal = async () => {
		createForm.set(defaultCreateForm());
		createError.set('');
		showCreateModal.set(true);
		ideaQuery.set('');
		projectQuery.set('');
		ideaDropdownOpen.set(false);
		projectDropdownOpen.set(false);
	};

	const setIdeaQuery = (value: string) => ideaQuery.set(value);
	const setProjectQuery = (value: string) => projectQuery.set(value);
	const setIdeaDropdownOpen = (value: boolean) => ideaDropdownOpen.set(value);
	const setProjectDropdownOpen = (value: boolean) => projectDropdownOpen.set(value);

	const closeCreateModal = () => {
		showCreateModal.set(false);
		createError.set('');
	};

	const handleCreateConversation = async (event: SubmitEvent) => {
		event.preventDefault();
		const form = get(createForm);
		if (!form.title.trim()) {
			const message = 'Title is required.';
			createError.set(message);
			toastError(message);
			return;
		}
		try {
			const payload = {
				title: form.title.trim(),
				description: form.description.trim(),
				ideaId: form.ideaId.trim() || undefined,
				projectId: form.projectId.trim() || undefined
			};
			const newConversation = await get(createConversationMutation).mutateAsync(payload);
			showCreateModal.set(false);
			searchInput.set('');
			conversationsFilters.set({
				search: '',
				includeArchived: get(includeArchived),
				enabled: true
			});
			// Convert Conversation to ChatConversation
			const chatConversation: ChatConversation = {
				id: newConversation.id,
				user_id: newConversation.userId,
				title: newConversation.title,
				description: newConversation.description,
				idea_id: newConversation.ideaId,
				project_id: newConversation.projectId,
				assigned_agent_id: newConversation.assignedAgentId,
				shared_transcript: newConversation.sharedTranscript,
				archived_at: newConversation.archivedAt,
				deleted_at: newConversation.deletedAt,
				created_at: newConversation.createdAt,
				updated_at: newConversation.updatedAt
			};
			selectedConversation.set(chatConversation);
			connectToStream(newConversation.id);
			await queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
		} catch {
			// handled via mutation callbacks
		}
	};

	const selectProject = (project: Project) => {
		createForm.update((form) => ({ ...form, projectId: project.id }));
		projectQuery.set(project.name);
		projectDropdownOpen.set(false);
	};

	const clearProjectSelection = () => {
		createForm.update((form) => ({ ...form, projectId: '' }));
		projectQuery.set('');
		projectDropdownOpen.set(false);
	};

	const selectIdea = (idea: Idea) => {
		createForm.update((form) => ({ ...form, ideaId: idea.id }));
		ideaQuery.set(idea.title);
		ideaDropdownOpen.set(false);
	};

	const clearIdeaSelection = () => {
		createForm.update((form) => ({ ...form, ideaId: '' }));
		ideaQuery.set('');
		ideaDropdownOpen.set(false);
	};

	const getFilteredProjects = () => {
		const items = (get(projectsQuery).data ?? []) as Project[];
		const query = normalize(get(projectQuery));
		if (!query) return items.slice(0, 10);
		return items.filter((project) => normalize(project.name).includes(query)).slice(0, 10);
	};

	const getFilteredIdeas = () => {
		const items = (get(ideasQuery).data ?? []) as Idea[];
		const query = normalize(get(ideaQuery));
		if (!query) return items.slice(0, 10);
		return items.filter((idea) => normalize(idea.title).includes(query)).slice(0, 10);
	};

	onDestroy(() => {
		conversationsUnsubscribe();
		disconnectStream();
	});

	return {
		isSending,
		provider,
		model,
		agent,
		autoAgent,
		conversationsQuery,
		messagesQuery,
		transcriptsQuery,
		assignmentsQuery,
		projectsQuery,
		ideasQuery,
		selectedConversation,
		selectedConversationIds,
		searchInput,
		includeArchived,
		bulkStatus,
		composeContent,
		composeRole,
		generateReply,
		transcriptStatus,
		showCreateModal,
		createError,
		createForm,
		ideaQuery,
		projectQuery,
		ideaDropdownOpen,
		projectDropdownOpen,
		refreshConversations,
		applyFilters,
		setSearchInput,
		setIncludeArchived,
		selectConversation,
		toggleConversationSelection,
		clearSelections,
		handleSendMessage,
		handleShareTranscript,
		applyBulkAction,
		openCreateModal,
		closeCreateModal,
		handleCreateConversation,
		updateCreateFormField,
		selectProject,
		clearProjectSelection,
		selectIdea,
		clearIdeaSelection,
		getFilteredProjects,
		getFilteredIdeas,
		setIdeaQuery,
		setProjectQuery,
		setIdeaDropdownOpen,
		setProjectDropdownOpen,
		setComposeContent: (value: string) => composeContent.set(value),
		setComposeRole: (value: 'user' | 'assistant') => composeRole.set(value),
		setGenerateReply: (value: boolean) => generateReply.set(value),
		setProvider: (v: string) => provider.set(v),
		setModel: (v: string) => model.set(v),
		setAgent: (v: string) => agent.set(v),
		setAutoAgent: (v: boolean) => autoAgent.set(v),
		createConversationMutation,
		appendMessageMutation,
		shareTranscriptMutation
	};
}

