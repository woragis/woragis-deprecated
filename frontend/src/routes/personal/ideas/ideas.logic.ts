import { browser } from '$app/environment';
import type { Edge, Node } from '@xyflow/svelte';
import { MarkerType, Position, type Connection } from '@xyflow/system';
import { useQueryClient } from '@tanstack/svelte-query';
import { onDestroy, onMount } from 'svelte';
import { get, writable } from 'svelte/store';

import type {
	ChatConversation,
	ConnectionDirection,
	Idea,
	IdeaDocument,
	IdeaLink,
	IdeaNode,
	IdeaNodeConnection,
	IdeaVersion,
	UUID
} from '$lib/api/types';
import {
	useCreateIdeaMutation,
	useCreateIdeaDocumentMutation,
	useCreateIdeaNodeMutation,
	useCreateIdeaNodeConnectionMutation,
	useCreateLinkMutation,
	useDeleteIdeaDocumentMutation,
	useDeleteIdeaNodeMutation,
	useDeleteIdeaNodeConnectionMutation,
	useIdeaLinksQuery,
	useIdeaVersionsQuery,
	useIdeasCanvasQuery,
	useIdeaDocumentsQuery,
	useIdeaNodeConnectionsQuery,
	useIdeaNodesQuery,
	useMoveIdeaNodeMutation,
	useUpdateIdeaDocumentMutation,
	useUpdateIdeaMutation,
	useUpdateIdeaNodeMutation
} from '@hooks/ideas';
import { useConversationsQuery } from '@hooks/chats';
import { useCreateConversationMutation } from '@hooks/chats';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export type CanvasNodeData = {
	node: IdeaNode;
};

export interface IdeaFormState {
	title: string;
	description: string;
	color: string;
}

export interface IdeaNodeFormState {
	title: string;
	description: string;
	color: string;
	type: string;
}

const defaultIdeaForm = (): IdeaFormState => ({
	title: '',
	description: '',
	color: '#2563eb'
});

const defaultNodeForm = (): IdeaNodeFormState => ({
	title: '',
	description: '',
	color: '#2563eb',
	type: 'default'
});

// Convert IdeaNode to SvelteFlow Node
const toNode = (node: IdeaNode): Node<CanvasNodeData> => {
	return {
		id: node.id,
		type: 'idea-node',
		position: {
			x: node.pos_x ?? Math.random() * 320,
			y: node.pos_y ?? Math.random() * 320
		},
		style: {
			width: node.width ?? 200,
			height: node.height ?? 100,
			borderColor: node.color || '#2563eb'
		} as any,
		data: { node },
		draggable: true,
		selectable: true,
		connectable: true
		// Note: Handle positions are defined in the IdeaNode component for 4-directional support
	};
};

// Convert IdeaNodeConnection to SvelteFlow Edge
const toEdge = (conn: IdeaNodeConnection, nodes: IdeaNode[]): Edge => {
	// Map direction to handle positions
	let sourceHandle: string = 'right';
	let targetHandle: string = 'left';

	switch (conn.direction) {
		case 'north':
			sourceHandle = 'top';
			targetHandle = 'bottom';
			break;
		case 'south':
			sourceHandle = 'bottom';
			targetHandle = 'top';
			break;
		case 'east':
			sourceHandle = 'right';
			targetHandle = 'left';
			break;
		case 'west':
			sourceHandle = 'left';
			targetHandle = 'right';
			break;
	}

	return {
		id: conn.id,
		source: conn.source_node_id,
		target: conn.target_node_id,
		type: 'smoothstep',
		label: conn.label,
		sourceHandle,
		targetHandle,
		markerEnd: {
			type: MarkerType.ArrowClosed,
			width: 18,
			height: 18,
			color: '#94a3b8'
		},
		style: {
			stroke: '#64748b',
			strokeWidth: 2
		} as any
	};
};

// Determine connection direction based on connection handle
function determineDirection(connection: Connection, nodes: IdeaNode[]): ConnectionDirection {
	const { sourceHandle, targetHandle } = connection;

	// Map handle positions to directions
	if (sourceHandle === 'top' || sourceHandle === 'Top') return 'north';
	if (sourceHandle === 'bottom' || sourceHandle === 'Bottom') return 'south';
	if (sourceHandle === 'right' || sourceHandle === 'Right') return 'east';
	if (sourceHandle === 'left' || sourceHandle === 'Left') return 'west';

	// Fallback: determine by node positions if handles aren't available
	const sourceNode = nodes.find((n: IdeaNode) => n.id === connection.source);
	const targetNode = nodes.find((n: IdeaNode) => n.id === connection.target);

	if (sourceNode && targetNode) {
		const dx = targetNode.pos_x - sourceNode.pos_x;
		const dy = targetNode.pos_y - sourceNode.pos_y;
		const absDx = Math.abs(dx);
		const absDy = Math.abs(dy);

		if (absDy > absDx) {
			return dy > 0 ? 'south' : 'north';
		} else {
			return dx > 0 ? 'east' : 'west';
		}
	}

	// Default to east
	return 'east';
}

export function createIdeasLogic() {
	const nodes = writable<Node<CanvasNodeData>[]>([]);
	const edges = writable<Edge[]>([]);
	const versions = writable<IdeaVersion[]>([]);
	const ideaChats = writable<ChatConversation[]>([]);
	const documents = writable<IdeaDocument[]>([]);
	const isSaving = writable(false);
	const isCreatingChat = writable(false);
	const isLoading = writable(false);
	const uiError = writable('');
	const selectedIdea = writable<Idea | null>(null);
	const selectedNode = writable<IdeaNode | null>(null);
	const editForm = writable<IdeaFormState | null>(null);
	const nodeEditForm = writable<IdeaNodeFormState | null>(null);
	const showCreateModal = writable(false);
	const showCreateNodeModal = writable(false);
	const showDocumentsPanel = writable(false);
	const newIdea = writable<IdeaFormState>(defaultIdeaForm());
	const newNode = writable<IdeaNodeFormState>(defaultNodeForm());
	const ideasQueryError = writable<unknown>(null);

	let ideas: Idea[] = [];
	let canvasNodes: IdeaNode[] = [];
	let canvasConnections: IdeaNodeConnection[] = [];
	let conversations: ChatConversation[] = [];

	const queryClient = useQueryClient();

	const ideasQuery = useIdeasCanvasQuery(browser);
	let nodesQuery = useIdeaNodesQuery(null, false);
	let connectionsQuery = useIdeaNodeConnectionsQuery(null, false);
	let documentsQuery = useIdeaDocumentsQuery(null, undefined, false);
	let versionsQuery = useIdeaVersionsQuery(null, { enabled: false, limit: 15 });
	const conversationsQuery = useConversationsQuery({ search: '', includeArchived: false, enabled: true });

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

	const updateCanvasForIdea = (ideaId: string | null) => {
		if (!ideaId) {
			nodes.set([]);
			edges.set([]);
			canvasNodes = [];
			canvasConnections = [];
			return;
		}

		// Subscribe to nodes and connections for the selected idea
		nodesQuery = useIdeaNodesQuery(ideaId, true);
		connectionsQuery = useIdeaNodeConnectionsQuery(ideaId, true);
		documentsQuery = useIdeaDocumentsQuery(ideaId, undefined, true);

		nodesQuery.subscribe((state) => {
			if (state.data) {
				canvasNodes = state.data;
				nodes.set(state.data.map(toNode));
				// Update edges when nodes change
				if (canvasConnections.length > 0) {
					edges.set(canvasConnections.map((conn) => toEdge(conn, canvasNodes)));
				}
			}
		});

		connectionsQuery.subscribe((state) => {
			if (state.data) {
				canvasConnections = state.data;
				edges.set(state.data.map((conn) => toEdge(conn, canvasNodes)));
			}
		});

		documentsQuery.subscribe((state) => {
			if (state.data) {
				documents.set(state.data);
			}
		});
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
		selectedNode.set(null);
		if (idea) {
			editForm.set({
				title: idea.title,
				description: idea.description ?? '',
				color: idea.color ?? '#2563eb'
			});
			updateCanvasForIdea(idea.id);
		} else {
			editForm.set(null);
			updateCanvasForIdea(null);
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
		const rawAny = raw as Record<string, any>;
		const id = rawAny.id ?? rawAny.ID;
		if (!id) return null;

		const userId = rawAny.user_id ?? rawAny.userId ?? '';
		const posX = rawAny.pos_x ?? rawAny.posX ?? 0;
		const posY = rawAny.pos_y ?? rawAny.posY ?? 0;
		const createdAt = rawAny.created_at ?? rawAny.createdAt ?? new Date().toISOString();
		const updatedAt = rawAny.updated_at ?? rawAny.updatedAt ?? createdAt;
		const title = rawAny.title ?? rawAny.Title ?? '';
		const slug = rawAny.slug ?? rawAny.Slug ?? buildIdeaSlug(title, String(id));

		return {
			...(raw as Idea),
			id: String(id),
			user_id: userId ? String(userId) : '',
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

	const ideasUnsubscribe = ideasQuery.subscribe((state) => {
		ideasQueryError.set(state.error);
		if (!state.data) return;
		ideas = sanitizeIdeas(state.data);
		syncSelectedIdea();
	});

	const chatsUnsubscribe = conversationsQuery.subscribe((state) => {
		if (!state.data) return;
		conversations = state.data;
		updateRelatedChats();
	});

	onDestroy(() => {
		ideasUnsubscribe();
		chatsUnsubscribe();
		versionsUnsubscribe?.();
	});

	const refreshCanvas = async () => {
		if (!browser) return;
		const currentIdea = get(selectedIdea);
		isLoading.set(true);
		uiError.set('');
		try {
			await Promise.all([
				queryClient.invalidateQueries({ queryKey: ['ideas', 'canvas'] }),
				currentIdea ? queryClient.invalidateQueries({ queryKey: ['ideas', currentIdea.id, 'nodes'] }) : Promise.resolve(),
				currentIdea ? queryClient.invalidateQueries({ queryKey: ['ideas', currentIdea.id, 'node-connections'] }) : Promise.resolve(),
				currentIdea ? queryClient.invalidateQueries({ queryKey: ['ideas', currentIdea.id, 'documents'] }) : Promise.resolve()
			]);
			toastInfo('Canvas refreshed.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to load canvas. Please try again.');
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
	const createNodeMutation = useCreateIdeaNodeMutation();
	const updateNodeMutation = useUpdateIdeaNodeMutation();
	const moveNodeMutation = useMoveIdeaNodeMutation();
	const createConnectionMutation = useCreateIdeaNodeConnectionMutation();
	const deleteNodeMutation = useDeleteIdeaNodeMutation();
	const deleteConnectionMutation = useDeleteIdeaNodeConnectionMutation();
	const createDocumentMutation = useCreateIdeaDocumentMutation();
	const updateDocumentMutation = useUpdateIdeaDocumentMutation();
	const deleteDocumentMutation = useDeleteIdeaDocumentMutation();
	const createConversationMutation = useCreateConversationMutation();

	const handleConnect = async (connection: Connection) => {
		const currentIdea = get(selectedIdea);
		if (!connection.source || !connection.target || !currentIdea) {
			return;
		}

		try {
			// Check if connection already exists
			const exists = canvasConnections.some(
				(c) => c.source_node_id === connection.source && c.target_node_id === connection.target
			);

			if (exists) {
				toastInfo('Connection already exists.');
				return;
			}

			const ideaNodes = get(nodesQuery).data ?? [];
			const direction = determineDirection(connection, ideaNodes);

			const conn = await get(createConnectionMutation).mutateAsync({
				idea_id: currentIdea.id,
				source_node_id: connection.source,
				target_node_id: connection.target,
				direction,
				label: ''
			});

			canvasConnections = [...canvasConnections, conn];
			edges.set(canvasConnections.map((c) => toEdge(c, canvasNodes)));
			toastSuccess('Nodes connected.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to create connection.');
			uiError.set(message);
			toastError(message);
		}
	};

	const handleNodeDragStop = async ({
		targetNode
	}: {
		targetNode: Node<CanvasNodeData> | null;
		event: MouseEvent | TouchEvent;
		nodes: Node<CanvasNodeData>[];
	}) => {
		const node = targetNode;
		if (!node) return;

		try {
			await get(moveNodeMutation).mutateAsync({
				nodeId: node.id,
				input: {
					pos_x: node.position.x,
					pos_y: node.position.y
				}
			});

			canvasNodes = canvasNodes.map((n) =>
				n.id === node.id
					? {
							...n,
							pos_x: node.position.x,
							pos_y: node.position.y,
							version: n.version + 1,
							updated_at: new Date().toISOString()
						}
					: n
			);
			nodes.set(canvasNodes.map(toNode));
			toastInfo('Node position saved.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to persist position.');
			uiError.set(message);
			toastError(message);
			await refreshCanvas();
			return;
		}
	};

	const handleSelectIdea = (idea: Idea) => {
		applySelectedIdea(idea);
	};

	const handleNodeClick = ({
		node
	}: {
		node: Node<CanvasNodeData>;
		event: MouseEvent | TouchEvent;
	}) => {
		const canvasNode = node?.data?.node;
		if (canvasNode) {
			selectedNode.set(canvasNode);
			nodeEditForm.set({
				title: canvasNode.title,
				description: canvasNode.description ?? '',
				color: canvasNode.color ?? '#2563eb',
				type: canvasNode.type ?? 'default'
			});
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

	const handleNodeSave = async () => {
		const node = get(selectedNode);
		const form = get(nodeEditForm);
		if (!node || !form) return;

		isSaving.set(true);
		try {
			const updated = await get(updateNodeMutation).mutateAsync({
				nodeId: node.id,
				input: {
					title: form.title,
					description: form.description,
					color: form.color,
					type: form.type
				}
			});

			canvasNodes = canvasNodes.map((n) => (n.id === updated.id ? updated : n));
			nodes.set(canvasNodes.map(toNode));
			selectedNode.set(updated);
			toastSuccess('Node updated.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to save node changes.');
			uiError.set(message);
			toastError(message);
		} finally {
			isSaving.set(false);
		}
	};

	const resetNewIdea = () => {
		newIdea.set(defaultIdeaForm());
	};

	const resetNewNode = () => {
		newNode.set(defaultNodeForm());
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
				pos_x: 0,
				pos_y: 0
			});
			ideas = [...ideas, idea];
			await refreshCanvas();
			resetNewIdea();
			showCreateModal.set(false);
			applySelectedIdea(idea);
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

	const handleCreateNode = async (event?: SubmitEvent) => {
		event?.preventDefault();
		const form = get(newNode);
		const currentIdea = get(selectedIdea);
		if (!form.title.trim()) {
			const message = 'Title is required to create a new node.';
			uiError.set(message);
			toastError(message);
			return;
		}
		if (!currentIdea) {
			const message = 'Please select an idea first.';
			uiError.set(message);
			toastError(message);
			return;
		}

		isSaving.set(true);
		try {
			const node = await get(createNodeMutation).mutateAsync({
				idea_id: currentIdea.id,
				title: form.title,
				description: form.description,
				color: form.color,
				type: form.type,
				pos_x: Math.random() * 400 + 100,
				pos_y: Math.random() * 300 + 100,
				width: 200,
				height: 100
			});
			canvasNodes = [...canvasNodes, node];
			nodes.set(canvasNodes.map(toNode));
			await refreshCanvas();
			resetNewNode();
			showCreateNodeModal.set(false);
			toastSuccess('Node created.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to create node.');
			uiError.set(message);
			toastError(message);
		} finally {
			isSaving.set(false);
		}
	};

	const handleDeleteNode = async (nodeId: UUID) => {
		try {
			await get(deleteNodeMutation).mutateAsync(nodeId);
			canvasNodes = canvasNodes.filter((n) => n.id !== nodeId);
			nodes.set(canvasNodes.map(toNode));
			// Remove connections involving this node
			canvasConnections = canvasConnections.filter(
				(c) => c.source_node_id !== nodeId && c.target_node_id !== nodeId
			);
			edges.set(canvasConnections.map((c) => toEdge(c, canvasNodes)));
			if (get(selectedNode)?.id === nodeId) {
				selectedNode.set(null);
				nodeEditForm.set(null);
			}
			toastSuccess('Node deleted.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to delete node.');
			uiError.set(message);
			toastError(message);
		}
	};

	const handleDeleteConnection = async (connectionId: UUID) => {
		try {
			await get(deleteConnectionMutation).mutateAsync(connectionId);
			canvasConnections = canvasConnections.filter((c) => c.id !== connectionId);
			edges.set(canvasConnections.map((c) => toEdge(c, canvasNodes)));
			toastSuccess('Connection deleted.');
		} catch (error) {
			console.error(error);
			const message = getApiErrorMessage(error, 'Unable to delete connection.');
			uiError.set(message);
			toastError(message);
		}
	};

	const updateEditFormField = <K extends keyof IdeaFormState>(field: K, value: IdeaFormState[K]) => {
		editForm.update((current) => (current ? { ...current, [field]: value } : current));
	};

	const updateNewIdeaField = <K extends keyof IdeaFormState>(field: K, value: IdeaFormState[K]) => {
		newIdea.update((current) => ({ ...current, [field]: value }));
	};

	const updateNodeEditFormField = <K extends keyof IdeaNodeFormState>(
		field: K,
		value: IdeaNodeFormState[K]
	) => {
		nodeEditForm.update((current) => (current ? { ...current, [field]: value } : current));
	};

	const updateNewNodeField = <K extends keyof IdeaNodeFormState>(field: K, value: IdeaNodeFormState[K]) => {
		newNode.update((current) => ({ ...current, [field]: value }));
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
		ideasQuery,
		conversationsQuery,
		documents,
		nodes,
		edges,
		versions,
		ideaChats,
		isSaving,
		isCreatingChat,
		isLoading,
		uiError,
		selectedIdea,
		selectedNode,
		ideas,
		editForm,
		nodeEditForm,
		showCreateModal,
		showCreateNodeModal,
		showDocumentsPanel,
		newIdea,
		newNode,
		ideasQueryError,
		refreshCanvas,
		handleConnect,
		handleNodeDragStop,
		handleNodeClick,
		handleIdeaSave,
		handleNodeSave,
		handleCreateIdea,
		handleCreateNode,
		handleSelectIdea,
		handleDeleteNode,
		handleDeleteConnection,
		createIdeaChat,
		updateEditFormField,
		updateNewIdeaField,
		updateNodeEditFormField,
		updateNewNodeField
	};
}
