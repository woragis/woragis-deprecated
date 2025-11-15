import { browser } from '$app/environment';
import type { Edge, Node } from '@xyflow/svelte';
import { MarkerType, Position, type Connection } from '@xyflow/system';
import { useQueryClient } from '@tanstack/svelte-query';
import { onDestroy, onMount } from 'svelte';
import { get, writable } from 'svelte/store';

import type { Idea, IdeaLink, IdeaVersion } from '$lib/api/types';
import {
	useCreateIdeaMutation,
	useCreateLinkMutation,
	useIdeaLinksQuery,
	useIdeaVersionsQuery,
	useIdeasCanvasQuery,
	useMoveIdeaMutation,
	useUpdateIdeaMutation
} from '@hooks/ideas';
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
	const isSaving = writable(false);
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

	const queryClient = useQueryClient();

	const ideasQuery = useIdeasCanvasQuery(browser);
	const linksQuery = useIdeaLinksQuery(browser);
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

	const ideasUnsubscribe = ideasQuery.subscribe((state) => {
		ideasQueryError.set(state.error);
		if (!state.data) return;
		ideas = state.data;
		nodes.set(ideas.map(toNode));
		syncSelectedIdea();
	});

	const linksUnsubscribe = linksQuery.subscribe((state) => {
		linksQueryError.set(state.error);
		if (!state.data) return;
		links = state.data;
		edges.set(links.map(toEdge));
	});

	onDestroy(() => {
		ideasUnsubscribe();
		linksUnsubscribe();
		versionsUnsubscribe?.();
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

	return {
		nodes,
		edges,
		versions,
		isSaving,
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
		updateEditFormField,
		updateNewIdeaField
	};
}


