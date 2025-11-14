<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

import { onMount } from 'svelte';
import { browser } from '$app/environment';
import { useQueryClient } from '@tanstack/svelte-query';
import {
	SvelteFlow,
	Background,
	Controls,
	MiniMap,
	Panel,
	type Edge,
	type Node,
	type NodeTypes
} from '@xyflow/svelte';
import { MarkerType, Position, type Connection } from '@xyflow/system';

	import IdeaNode from './IdeaNode.svelte';
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

type IdeaNodeData = {
	idea: Idea;
};

	let nodes: Node<IdeaNodeData>[] = [];
	let edges: Edge[] = [];

	let ideas: Idea[] = [];
	let links: IdeaLink[] = [];
	let versions: IdeaVersion[] = [];

let isSaving = false;
	let isLoading = false;
let uiError = '';
	let selectedIdea: Idea | null = null;
	let editForm: { title: string; description: string; color: string } | null = null;
	let showCreateModal = false;
	let newIdea = {
		title: '',
		description: '',
		color: '#2563eb'
	};

	const nodeTypes: NodeTypes = {
		idea: IdeaNode
	};

	const toNode = (idea: Idea): Node<IdeaNodeData> => ({
		id: idea.id,
		type: 'idea',
		position: {
			x: idea.pos_x ?? Math.random() * 320,
			y: idea.pos_y ?? Math.random() * 320
		},
		data: {
			idea
		},
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

const syncSelectedIdea = () => {
	if (!selectedIdea) {
		editForm = null;
		return;
	}
	const refreshed = ideas.find((idea) => idea.id === selectedIdea?.id) ?? null;
	selectedIdea = refreshed ?? null;
	if (selectedIdea) {
		editForm = {
			title: selectedIdea.title,
			description: selectedIdea.description ?? '',
			color: selectedIdea.color ?? '#2563eb'
		};
	} else {
		editForm = null;
		versions = [];
	}
};

const ideasQuery = useIdeasCanvasQuery(false);
const linksQuery = useIdeaLinksQuery(false);
let versionsQuery = useIdeaVersionsQuery(null, { enabled: false, limit: 15 });
const queryClient = useQueryClient();

const createIdeaMutation = useCreateIdeaMutation();
const updateIdeaMutation = useUpdateIdeaMutation();
const moveIdeaMutation = useMoveIdeaMutation();
const createLinkMutation = useCreateLinkMutation();

$: if ($ideasQuery.data) {
	ideas = $ideasQuery.data;
	nodes = ideas.map(toNode);
	syncSelectedIdea();
}

$: if ($linksQuery.data) {
	links = $linksQuery.data;
	edges = links.map(toEdge);
}

$: versions = $versionsQuery.data ?? [];

$: {
	const currentId = selectedIdea?.id ?? null;
	versionsQuery = useIdeaVersionsQuery(currentId, {
		enabled: Boolean(currentId),
		limit: 15
	});
	if (!currentId) {
		versions = [];
	}
}

const refreshCanvas = async () => {
	if (!browser) return;
	isLoading = true;
	uiError = '';
	try {
		await Promise.all([
			queryClient.invalidateQueries({ queryKey: ['ideas', 'canvas'] }),
			queryClient.invalidateQueries({ queryKey: ['ideas', 'links'] })
		]);
		toastInfo('Ideas refreshed.');
	} catch (error) {
		console.error(error);
		uiError = getApiErrorMessage(error, 'Unable to load ideas. Please try again.');
		toastError(uiError);
	} finally {
		isLoading = false;
	}
};

onMount(async () => {
	await refreshCanvas();
});

const handleConnect = async (connection: Connection) => {
		if (!connection.source || !connection.target) {
			return;
		}

		try {
			const link = await $createLinkMutation.mutateAsync({
				source_idea_id: connection.source,
				target_idea_id: connection.target,
				relation: 'relates',
				bidirectional: false
			});
			links = [...links, link];
		const newEdge: Edge = {
			id: link.id,
			source: connection.source,
			target: connection.target,
			type: 'smoothstep',
			label: link.relation,
			animated: Boolean(link.bidirectional),
			markerEnd: {
				type: MarkerType.ArrowClosed,
				width: 18,
				height: 18,
				color: '#94a3b8'
			}
		};
		edges = [...edges, newEdge];
			toastSuccess('Ideas linked.');
		} catch (error) {
			console.error(error);
			uiError = getApiErrorMessage(error, 'Unable to create relation.');
			toastError(uiError);
		}
	};

	async function handleNodeDragStop({
		targetNode
	}: {
		targetNode: Node<IdeaNodeData> | null;
		event: MouseEvent | TouchEvent;
		nodes: Node<IdeaNodeData>[];
	}) {
		const node = targetNode;
		if (!node) return;

		try {
			await $moveIdeaMutation.mutateAsync({
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
			uiError = getApiErrorMessage(error, 'Unable to persist position.');
			toastError(uiError);
			await refreshCanvas();
			return;
		}
		toastInfo('Idea position saved.');
	}

	function handleSelectIdea(idea: Idea) {
		selectedIdea = idea;
		editForm = {
			title: idea.title,
			description: idea.description ?? '',
			color: idea.color ?? '#2563eb'
		};
	if (selectedIdea?.id) {
		versionsQuery = useIdeaVersionsQuery(selectedIdea.id, {
			enabled: true,
			limit: 15
		});
	}
	}

	function handleNodeClick({
		node
	}: {
		node: Node<IdeaNodeData>;
		event: MouseEvent | TouchEvent;
	}) {
		const idea = node?.data?.idea;
		if (idea) {
			handleSelectIdea(idea);
		}
	}

	async function handleIdeaSave() {
		if (!selectedIdea || !editForm) return;

		isSaving = true;
		try {
			const updated = await $updateIdeaMutation.mutateAsync({
				ideaId: selectedIdea.id,
				input: {
					title: editForm.title,
					description: editForm.description,
					color: editForm.color
				}
			});

			ideas = ideas.map((idea) => (idea.id === updated.id ? updated : idea));
			nodes = nodes.map((node) =>
				node.id === updated.id
					? {
							...node,
							data: {
								...node.data,
								idea: updated
							}
						}
					: node
			);
			selectedIdea = updated;
			toastSuccess('Idea updated.');
		} catch (error) {
			console.error(error);
			uiError = getApiErrorMessage(error, 'Unable to save idea changes.');
			toastError(uiError);
		} finally {
			isSaving = false;
		}
	}

	async function handleCreateIdea(event: SubmitEvent) {
		event.preventDefault();
		if (!newIdea.title.trim()) {
			uiError = 'Title is required to create a new idea.';
			toastError(uiError);
			return;
		}

		isSaving = true;
		try {
			const idea = await $createIdeaMutation.mutateAsync({
				title: newIdea.title,
				description: newIdea.description,
				color: newIdea.color,
				pos_x: Math.random() * 420,
				pos_y: Math.random() * 300
			});
			ideas = [...ideas, idea];
			nodes = [...nodes, toNode(idea)];
			await refreshCanvas();
			newIdea = {
				title: '',
				description: '',
				color: '#2563eb'
			};
			showCreateModal = false;
			toastSuccess('Idea created.');
		} catch (error) {
			console.error(error);
			uiError = getApiErrorMessage(error, 'Unable to create idea.');
			toastError(uiError);
		} finally {
			isSaving = false;
		}
	}
</script>

<svelte:head>
	<title>Ideas Canvas · Woragis</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<header class="flex flex-wrap items-center justify-between gap-4">
		<div>
			<h1 class="text-2xl font-semibold text-slate-100">Ideas Canvas</h1>
			<p class="text-sm text-slate-400">
				Map your product ideas, connect concepts, and keep a history of every iteration.
			</p>
		</div>
		<div class="flex items-center gap-3">
			<button
				class="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
				type="button"
				on:click={refreshCanvas}
				disabled={isLoading}
			>
				{isLoading ? 'Refreshing…' : 'Refresh'}
			</button>
			<button
				class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90"
				type="button"
				on:click={() => (showCreateModal = true)}
			>
				New Idea
			</button>
		</div>
	</header>

	{#if uiError}
		<div class="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			{uiError}
		</div>
	{:else if $ideasQuery.error || $linksQuery.error}
		<div class="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			Unable to load ideas. Please try again.
		</div>
	{/if}

	<div class="grid gap-6 lg:grid-cols-[2fr_1fr]">
		<section class="relative min-h-[540px] overflow-hidden rounded-2xl border border-slate-800/80 bg-slate-950/60">
			{#if isLoading}
				<div class="absolute inset-0 z-10 flex items-center justify-center bg-slate-950/70">
					<div class="flex items-center gap-3 text-sm text-slate-300">
						<span class="h-3 w-3 animate-ping rounded-full bg-primary"></span>
						Loading canvas…
					</div>
				</div>
			{/if}

			<SvelteFlow
				bind:nodes
				bind:edges
				{nodeTypes}
				fitView
				class="h-[600px]"
				onconnect={handleConnect}
				onnodedragstop={handleNodeDragStop}
				onnodeclick={handleNodeClick}
			>
				<Background patternColor="#1e293b" />
				<Controls />
				<MiniMap
					nodeStrokeColor="#38bdf8"
					nodeColor="#0f1729"
					class="border border-slate-800/60 bg-slate-900/80"
				/>
				<Panel position="top-left" class="rounded-lg border border-slate-800/60 bg-slate-900/80 px-3 py-2 text-xs text-slate-300">
					Tip: drag ideas to reposition and connect them by drawing from one node handle to another.
				</Panel>
			</SvelteFlow>
		</section>

		<aside class="flex flex-col gap-6 rounded-2xl border border-slate-800/80 bg-slate-950/60 p-5">
			{#if selectedIdea && editForm}
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-slate-100">Idea Details</h2>
					<span class="rounded-full border border-slate-700/70 px-2 py-0.5 text-xs text-slate-400"
						>Version {selectedIdea.version}</span
					>
				</div>
				<form class="flex flex-col gap-4 text-sm" on:submit|preventDefault={handleIdeaSave}>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
						<input
							bind:value={editForm.title}
							class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Idea title"
							required
						/>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
						<textarea
							bind:value={editForm.description}
							class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Describe the problem or solution…"
						></textarea>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Color</span>
						<input
							type="color"
							bind:value={editForm.color}
							class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
						/>
					</label>
					<button
						class="mt-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-70"
						type="submit"
						disabled={isSaving}
					>
						{isSaving ? 'Saving…' : 'Save changes'}
					</button>
				</form>

				<div class="space-y-3">
					<h3 class="text-sm font-semibold text-slate-200">Recent versions</h3>
					{#if versions.length === 0}
						<p class="text-xs text-slate-500">Changes will show up here as you iterate.</p>
					{:else}
						<ul class="flex flex-col gap-2 text-xs text-slate-300">
							{#each versions as version}
								<li class="rounded-lg border border-slate-800/70 bg-slate-900/60 px-3 py-2">
									<div class="flex items-center justify-between">
										<span class="font-medium text-slate-200">v{version.version}</span>
										<span class="text-slate-500">{new Date(version.created_at).toLocaleString()}</span>
									</div>
									<p class="mt-1 text-slate-400 capitalize">{version.change_type}</p>
									{#if version.title && version.title !== selectedIdea.title}
										<p class="mt-1 text-slate-500">
											Title: {version.title}
										</p>
									{/if}
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			{:else}
				<div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
					<h2 class="text-lg font-semibold text-slate-200">Select an Idea</h2>
					<p class="max-w-xs text-sm text-slate-500">
						Choose a node on the canvas to review its details, update the description, or inspect its change history.
					</p>
				</div>
			{/if}
		</aside>
	</div>

	{#if showCreateModal}
		<div class="fixed inset-0 z-30 flex items-center justify-center bg-slate-950/80 backdrop-blur">
			<form
				class="w-full max-w-md space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
				on:submit={handleCreateIdea}
			>
				<header class="flex items-center justify-between">
					<div>
						<h3 class="text-lg font-semibold text-slate-100">Create new idea</h3>
						<p class="text-xs text-slate-400">Start with a clear name and optional description.</p>
					</div>
					<button
						class="rounded-full border border-slate-700/70 px-2 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
						type="button"
						on:click={() => (showCreateModal = false)}
					>
						Close
					</button>
				</header>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
					<input
						bind:value={newIdea.title}
						class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Idea title"
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
					<textarea
						bind:value={newIdea.description}
						class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Optional context"
					></textarea>
				</label>
				<label class="flex items-center justify-between text-sm text-slate-300">
					<span>Color</span>
					<input
						type="color"
						bind:value={newIdea.color}
						class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
					/>
				</label>
				<div class="flex items-center justify-end gap-3 pt-2">
					<button
						class="rounded-lg border border-slate-700/70 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
						type="button"
						on:click={() => (showCreateModal = false)}
					>
						Cancel
					</button>
					<button
						class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
						type="submit"
						disabled={isSaving}
					>
						{isSaving ? 'Creating…' : 'Create idea'}
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>

<style>
	:global(.svelte-flow__attribution) {
		display: none;
	}
</style>

