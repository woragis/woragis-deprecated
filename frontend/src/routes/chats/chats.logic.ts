import { browser } from '$app/environment';
import { onDestroy } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import type { AppendMessageInput } from '$lib/api/chats';
import { CHATS_STREAM_BASE } from '$lib/api/chats';
import type { ChatConversation, Idea, Project, UUID } from '$lib/api/types';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';
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
		enabled: true
	});

	const conversationsQuery = useConversationsQuery(conversationsFilters);
	const selectedConversationId = derived(
		selectedConversation,
		($conversation) => $conversation?.id ?? null
	);
	const messagesQuery = useConversationMessagesQuery(selectedConversationId);
	const transcriptsQuery = useConversationTranscriptsQuery(selectedConversationId);
	const assignmentsQuery = useConversationAssignmentsQuery(selectedConversationId);

	const projectsQuery = useProjectsListQuery({ enabled: true });
	const ideasQuery = useIdeasReferenceQuery(true);

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

	let websocket: WebSocket | null = null;

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

	const connectToStream = (conversationId: UUID) => {
		if (!browser) return;
		disconnectStream();
		const base = CHATS_STREAM_BASE.replace(/\/+$/, '');
		const url = `${base}/chats/conversations/${conversationId}/stream`;
		try {
			websocket = new WebSocket(url);
			websocket.onmessage = (event) => {
				try {
					const payload = JSON.parse(event.data);
					queryClient.setQueryData(
						['conversation', conversationId, 'messages'],
						(previous: any) => {
							if (!payload?.conversation_id || payload.conversation_id !== conversationId) {
								return previous ?? [];
							}
							const baseMessages = previous ?? [];
							const nextMessages = [...baseMessages, payload];
							return nextMessages.sort(
								(a: any, b: any) =>
									new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
							);
						}
					);
				} catch (error) {
					console.warn('Unable to parse stream payload', error);
				}
			};
			websocket.onclose = () => {
				websocket = null;
			};
		} catch (error) {
			console.error('Unable to open chat stream', error);
		}
	};

	const disconnectStream = () => {
		if (websocket) {
			websocket.close();
			websocket = null;
		}
	};

	const setSearchInput = (value: string) => searchInput.set(value);
	const setIncludeArchived = (value: boolean) => includeArchived.set(value);

	const applyFilters = async () => {
		const nextSearch = get(searchInput).trim();
		const include = get(includeArchived);
		conversationsFilters.set({
			search: nextSearch,
			includeArchived: include,
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
		if (!conversation || !content) return;
		const payload: AppendMessageInput = {
			role,
			content,
			generate_reply: shouldReply
		};
		const newMessages = await get(appendMessageMutation).mutateAsync({
			conversationId: conversation.id,
			input: payload
		});
		queryClient.setQueryData(
			['conversation', conversation.id, 'messages'],
			(previous: any) => {
				const baseMessages = previous ?? [];
				const nextMessages = [...baseMessages, ...newMessages];
				return nextMessages.sort(
					(a: any, b: any) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
				);
			}
		);
		queryClient.invalidateQueries({ queryKey: ['conversations'] });
		composeContent.set('');
		generateReply.set(false);
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
			transcriptStatus.set(`Share code ${transcript.share_code}`);
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
			selectedConversation.set(newConversation);
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
		createConversationMutation,
		appendMessageMutation,
		shareTranscriptMutation
	};
}

