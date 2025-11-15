import { browser } from '$app/environment';
import type { Edge, Node } from '@xyflow/svelte';
import { MarkerType, Position, type Connection } from '@xyflow/system';
import { useQueryClient } from '@tanstack/svelte-query';
import { onDestroy, onMount } from 'svelte';
import { get, writable } from 'svelte/store';

import type { ChatConversation, Idea, IdeaLink, IdeaVersion } from '$lib/api/types';
import {
	useCreateIdeaMutation,
	useCreateLinkMutation,
	useIdeaLinksQuery,
	useIdeaVersionsQuery,
	useIdeasCanvasQuery,
	useMoveIdeaMutation,
	useUpdateIdeaMutation
} from '@hooks/ideas';
import { useConversationsQuery } from '@hooks/chats';
import { useCreateConversationMutation } from '@hooks/chats';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export type IdeaNodeData = {
	idea: Idea;
};

export interface IdeaFormState {
	title: string;
	description: string;
	color: string;
}

const defaultIdeaForm = (): IdeaFormState => ({
	title: '',
	description: '',
	color: '#2563eb'
});

const toNode = (idea: Idea): Node<IdeaNodeData> => ({
	id: idea.id,
	type: 'idea',
	position: {
		x: idea.pos_x ?? Math.random() * 320,
		y: idea.pos_y ?? Math.random() * 320
	},
	data: { idea },
	draggable: true,
	selectable: true,
	connectable: true,
	sourcePosition: Position.Right,
	targetPosition: Position.Left
});

const toEdge = (link: IdeaLink): Edge => ({
	id: link.id,
	source: link.source_idea_id,
	target: link.target_idea_id,
	type: 'smoothstep',
	label: link.relation,
	animated: Boolean(link.bidirectional),
	markerEnd: {
		type: MarkerType.ArrowClosed,
		width: 18,
		height: 18,
		color: '#94a3b8'
	}
});

export function createIdeasLogic() {
	const nodes = writable<Node<IdeaNodeData>[]>([]);
	const edges = writable<Edge[]>([]);
	const versions = writable<IdeaVersion[]>([]);
	const ideaChats = writable<ChatConversation[]>([]);
	const isSaving = writable(false);
	const isCreatingChat = writable(false);
	const isLoading = writable(false);
	const uiError = writable('');
	const selectedIdea = writable<Idea | null>(null);
	const editForm = writable<IdeaFormState | null>(null);
	const showCreateModal = writable(false);
	const newIdea = writable<IdeaFormState>(defaultIdeaForm());
	const ideasQueryError = writable<unknown>(null);
	const linksQueryError = writable<unknown>(null);

	let ideas: Idea[] = [];
	let links: IdeaLink[] = [];
	let conversations: ChatConversation[] = [];

	const queryClient = useQueryClient();

	const ideasQuery = useIdeasCanvasQuery(browser);
	const linksQuery = useIdeaLinksQuery(browser);
	const conversationsQuery = useConversationsQuery({ search: '', includeArchived: false, enabled: true });
	let versionsQuery = useIdeaVersionsQuery(null, { enabled: false, limit: 15 });
	let versionsUnsubscribe: (() => void) | null = null;
	let currentVersionsIdeaId: string | null = null;

	const subscribeToVersionsQuery = () => {
		versionsUnsubscribe?.();
		versionsUnsubscribe = versionsQuery.subscribe((state) => {
			versions.set(state.data ?? []);
		});
	};

	subscribeToVersionsQuery();

	const setVersionsQueryTarget = (ideaId: string | null) => {
		if (currentVersionsIdeaId === ideaId) {
			if (!ideaId) {
				versions.set([]);
			}
			return;
		}

		currentVersionsIdeaId = ideaId;
		versionsQuery = useIdeaVersionsQuery(ideaId, {
			enabled: Boolean(ideaId),
			limit: 15
		});
		subscribeToVersionsQuery();
		if (!ideaId) {
			versions.set([]);
		}
	};

	const updateRelatedChats = () => {
		const ideaId = get(selectedIdea)?.id;
		if (!ideaId) {
			ideaChats.set([]);
			return;
		}
		const filtered = conversations
			.filter((chat) => chat.idea_id === ideaId)
			.sort(
				(a, b) =>
					new Date(b.updated_at ?? b.created_at ?? '').getTime() -
					new Date(a.updated_at ?? a.created_at ?? '').getTime()
			);
		ideaChats.set(filtered);
	};

	const applySelectedIdea = (idea: Idea | null) => {
		selectedIdea.set(idea);
		if (idea) {
			editForm.set({
				title: idea.title,
				description: idea.description ?? '',
				color: idea.color ?? '#2563eb'
			});
		} else {
			editForm.set(null);
		}
		setVersionsQueryTarget(idea?.id ?? null);
		updateRelatedChats();
	};

	const syncSelectedIdea = () => {
		const current = get(selectedIdea);
		if (!current) {
			applySelectedIdea(null);
			return;
		}
		const refreshed = ideas.find((idea) => idea.id === current.id) ?? null;
		applySelectedIdea(refreshed);
	};

const sanitizeSlug = (value: string) =>
	value
		.toLowerCase()
		.trim()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '') || 'idea';

const buildIdeaSlug = (title: string, id: string) => `${sanitizeSlug(title)}--${id}`;

const normalizeIdea = (raw: Idea | Record<string, any>): Idea | null => {
	const id = raw.id ?? raw.ID;
	if (!id) return null;

	const userId = raw.user_id ?? raw.userId ?? '';
	const posX = raw.pos_x ?? raw.posX ?? 0;
	const posY = raw.pos_y ?? raw.posY ?? 0;
	const createdAt = raw.created_at ?? raw.createdAt ?? new Date().toISOString();
	const updatedAt = raw.updated_at ?? raw.updatedAt ?? createdAt;
	const title = raw.title ?? raw.Title ?? '';
	const slug = raw.slug ?? raw.Slug ?? buildIdeaSlug(title, String(id));

	return {
		...(raw as Idea),
		id: String(id),
		user_id: userId ? String(userId) : undefined,
		slug,
		pos_x: Number.isFinite(posX) ? Number(posX) : 0,
		pos_y: Number.isFinite(posY) ? Number(posY) : 0,
		color: raw.color ?? '#2563eb',
		version: raw.version ?? 1,
		created_at: createdAt,
		updated_at: updatedAt
	};
};

const sanitizeIdeas = (items: Idea[]) => {
	const seen = new Set<string>();
	const result: Idea[] = [];
	items.forEach((idea) => {
		const normalized = normalizeIdea(idea);
		if (!normalized) {
			console.warn('Skipping idea without id', idea);
			return;
		}
		const id = normalized.id;
		if (seen.has(id)) {
			console.warn('Skipping duplicate idea id', id);
			return;
		}
		seen.add(id);
		result.push(normalized);
	});
	return result;
};

const sanitizeLinks = (items: IdeaLink[]) => {
	const seen = new Set<string>();
	const result: IdeaLink[] = [];
	items.forEach((link) => {
		const id = link.id ? String(link.id) : '';
		if (!id) {
			console.warn('Skipping link without id', link);
		 return;
		}
		if (seen.has(id)) {
			return;
		}
		if (!link.source_idea_id || !link.target_idea_id) {
			return;
		}
		seen.add(id);
		result.push({ ...link, id });
	});
	return result;
};

const ideasUnsubscribe = ideasQuery.subscribe((state) => {
	ideasQueryError.set(state.error);
	if (!state.data) return;
	ideas = sanitizeIdeas(state.data);
	nodes.set(ideas.map(toNode));
	syncSelectedIdea();
});

const linksUnsubscribe = linksQuery.subscribe((state) => {
	linksQueryError.set(state.error);
	if (!state.data) return;
	links = sanitizeLinks(state.data);
	edges.set(links.map(toEdge));
});

	const chatsUnsubscribe = conversationsQuery.subscribe((state) => {
		if (!state.data) return;
		conversations = state.data;
		updateRelatedChats();
	});

	onDestroy(() => {
		ideasUnsubscribe();
		linksUnsubscribe();
		versionsUnsubscribe?.();
	chatsUnsubscribe();
	});

	const refreshCanvas = async () => {
		if (!browser) return;
		isLoading.set(true);
		uiError.set('');
		try {
			await Promise.all([
				queryClient.invalidateQueries({ queryKey: ['ideas', 'canvas'] }),
				queryClient.invalidateQueries({ queryKey: ['ideas', 'links'] })
			]);
			toastInfo('Ideas refreshed.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to load ideas. Please try again.');
			uiError.set(message);
			toastError(message);
		} finally {
			isLoading.set(false);
		}
	};

	onMount(async () => {
		await refreshCanvas();
	});

	const createIdeaMutation = useCreateIdeaMutation();
	const updateIdeaMutation = useUpdateIdeaMutation();
	const moveIdeaMutation = useMoveIdeaMutation();
	const createLinkMutation = useCreateLinkMutation();
	const createConversationMutation = useCreateConversationMutation();

	const updateEdgesWithNewLink = (link: IdeaLink) => {
		const newEdge = toEdge(link);
		edges.update((current) => [...current, newEdge]);
	};

	const handleConnect = async (connection: Connection) => {
		if (!connection.source || !connection.target) {
			return;
		}

		try {
			const link = await get(createLinkMutation).mutateAsync({
				source_idea_id: connection.source,
				target_idea_id: connection.target,
				relation: 'relates',
				bidirectional: false
			});
			links = [...links, link];
			updateEdgesWithNewLink(link);
			toastSuccess('Ideas linked.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to create relation.');
			uiError.set(message);
			toastError(message);
		}
	};

	const handleNodeDragStop = async ({
		targetNode
	}: {
		targetNode: Node<IdeaNodeData> | null;
		event: MouseEvent | TouchEvent;
		nodes: Node<IdeaNodeData>[];
	}) => {
		const node = targetNode;
		if (!node) return;

		try {
			await get(moveIdeaMutation).mutateAsync({
				ideaId: node.id,
				input: {
					pos_x: node.position.x,
					pos_y: node.position.y
				}
			});

			ideas = ideas.map((idea) =>
				idea.id === node.id
					? {
							...idea,
							pos_x: node.position.x,
							pos_y: node.position.y,
							version: idea.version + 1,
							updated_at: new Date().toISOString()
						}
					: idea
			);
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to persist position.');
			uiError.set(message);
			toastError(message);
			await refreshCanvas();
			return;
		}
		nodes.set(ideas.map(toNode));
		toastInfo('Idea position saved.');
	};

	const handleSelectIdea = (idea: Idea) => {
		applySelectedIdea(idea);
	};

	const handleNodeClick = ({
		node
	}: {
		node: Node<IdeaNodeData>;
		event: MouseEvent | TouchEvent;
	}) => {
		const idea = node?.data?.idea;
		if (idea) {
			handleSelectIdea(idea);
		}
	};

	const handleIdeaSave = async () => {
		const idea = get(selectedIdea);
		const form = get(editForm);
		if (!idea || !form) return;

		isSaving.set(true);
		try {
			const updated = await get(updateIdeaMutation).mutateAsync({
				ideaId: idea.id,
				input: {
					title: form.title,
					description: form.description,
					color: form.color
				}
			});

			ideas = ideas.map((item) => (item.id === updated.id ? updated : item));
			nodes.set(ideas.map(toNode));
			applySelectedIdea(updated);
			toastSuccess('Idea updated.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to save idea changes.');
			uiError.set(message);
			toastError(message);
		} finally {
			isSaving.set(false);
		}
	};

	const resetNewIdea = () => {
		newIdea.set(defaultIdeaForm());
	};

	const handleCreateIdea = async (event?: SubmitEvent) => {
		event?.preventDefault();
		const form = get(newIdea);
		if (!form.title.trim()) {
			const message = 'Title is required to create a new idea.';
			uiError.set(message);
			toastError(message);
			return;
		}

		isSaving.set(true);
		try {
			const idea = await get(createIdeaMutation).mutateAsync({
				title: form.title,
				description: form.description,
				color: form.color,
				pos_x: Math.random() * 420,
				pos_y: Math.random() * 300
			});
			ideas = [...ideas, idea];
			nodes.set(ideas.map(toNode));
			await refreshCanvas();
			resetNewIdea();
			showCreateModal.set(false);
			toastSuccess('Idea created.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to create idea.');
			uiError.set(message);
			toastError(message);
		} finally {
			isSaving.set(false);
		}
	};

	const updateEditFormField = <K extends keyof IdeaFormState>(field: K, value: IdeaFormState[K]) => {
		editForm.update((current) => (current ? { ...current, [field]: value } : current));
	};

	const updateNewIdeaField = <K extends keyof IdeaFormState>(field: K, value: IdeaFormState[K]) => {
		newIdea.update((current) => ({ ...current, [field]: value }));
	};

	const createIdeaChat = async () => {
		const idea = get(selectedIdea);
		if (!idea) {
			toastError('Select an idea before creating a chat.');
			return null;
		}

		isCreatingChat.set(true);
		try {
			const conversation = await get(createConversationMutation).mutateAsync({
				title: idea.title || 'New idea conversation',
				description: idea.description ?? '',
				ideaId: idea.id,
				projectId: idea.project_id ?? undefined
			});
			toastSuccess('Chat created for idea.');
			await queryClient.invalidateQueries({ queryKey: ['chats', 'conversations'] });
			return conversation;
		} catch (error) {
			const message = getApiErrorMessage(error, 'Unable to create chat for idea.');
			uiError.set(message);
			toastError(message);
			return null;
		} finally {
			isCreatingChat.set(false);
		}
	};

	return {
		conversationsQuery,
		nodes,
		edges,
		versions,
		ideaChats,
		isSaving,
		isCreatingChat,
		isLoading,
		uiError,
		selectedIdea,
		editForm,
		showCreateModal,
		newIdea,
		ideasQueryError,
		linksQueryError,
		refreshCanvas,
		handleConnect,
		handleNodeDragStop,
		handleNodeClick,
		handleIdeaSave,
		handleCreateIdea,
		createIdeaChat,
		updateEditFormField,
		updateNewIdeaField
	};
}


