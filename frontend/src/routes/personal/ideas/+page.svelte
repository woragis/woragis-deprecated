<script lang="ts">
	import '@xyflow/svelte/dist/style.css';

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

	import IdeaNodeComponent from './IdeaNode.svelte';
	import MarkdownRenderer from '$lib/components/MarkdownRenderer.svelte';
	import IdeasHero from './_components/IdeasHero.svelte';
	import {
		createIdeasLogic,
		type CanvasNodeData,
		type IdeaFormState,
		type IdeaNodeFormState
	} from './ideas.logic';
	import { goto } from '$app/navigation';

	const nodeTypes: NodeTypes = {
		'idea-node': IdeaNodeComponent
	};

	const {
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
	} = createIdeasLogic();

	const startIdeaChat = async () => {
		const conversation = await createIdeaChat();
		if (conversation) {
			goto(`/chats?conversation=${conversation.id}`);
		}
	};

	let flowNodes: Node<CanvasNodeData>[] = [];
	let flowEdges: Edge[] = [];
	let chatsForIdea = [];

	$: flowNodes = $nodes;
	$: flowEdges = $edges;
	$: chatsForIdea = $ideaChats;

	const handleInput =
		<T extends keyof IdeaFormState>(updater: (value: IdeaFormState[T]) => void) =>
		(event: Event) => {
			const target = event.currentTarget as HTMLInputElement | HTMLTextAreaElement;
			updater(target.value as IdeaFormState[T]);
		};

	const handleNodeInput =
		<T extends keyof IdeaNodeFormState>(updater: (value: IdeaNodeFormState[T]) => void) =>
		(event: Event) => {
			const target = event.currentTarget as HTMLInputElement | HTMLTextAreaElement;
			updater(target.value as IdeaNodeFormState[T]);
		};

	// Filter documents by selected node
	$: nodeDocuments = $selectedNode
		? $documents.filter((doc) => doc.node_id === $selectedNode.id)
		: $documents.filter((doc) => !doc.node_id);
</script>

<svelte:head>
	<title>Ideas Canvas · Woragis</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<IdeasHero
		onRefresh={refreshCanvas}
		onNewIdea={() => showCreateModal.set(true)}
		onNewNode={() => showCreateNodeModal.set(true)}
		showNewNode={!!$selectedIdea}
	/>

	<!-- Main Content -->
	<div class="mx-auto max-w-[1920px] px-6 py-8 lg:px-8">

		{#if $uiError}
			<div class="mb-6 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
				{$uiError}
			</div>
		{:else if $ideasQueryError}
			<div class="mb-6 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
				Unable to load ideas. Please try again.
			</div>
		{/if}

		<div class="grid gap-6 lg:grid-cols-[280px_1fr_360px]">
			<!-- Ideas List Sidebar -->
			<aside
				class="flex flex-col gap-4 rounded-xl border border-slate-800/50 bg-slate-900/40 p-4 shadow-xl backdrop-blur-sm"
			>
			<div class="flex items-center justify-between">
				<h2 class="text-lg font-semibold text-slate-100">Ideas</h2>
				<span class="rounded-full bg-slate-800 px-2 py-0.5 text-xs text-slate-400">
					{$ideasQuery.data?.length ?? 0}
				</span>
			</div>
			<div class="flex-1 space-y-2 overflow-y-auto">
				{#if $ideasQuery.isLoading}
					<p class="text-xs text-slate-500">Loading ideas…</p>
				{:else if $ideasQuery.data && $ideasQuery.data.length === 0}
					<p class="text-xs text-slate-500">No ideas yet. Create one to get started.</p>
				{:else}
					{#each $ideasQuery.data ?? [] as idea}
						<button
							class="group w-full rounded-lg border px-3 py-2 text-left text-sm transition-all hover:scale-[1.02] hover:border-violet-500/50 hover:bg-slate-900/60 hover:shadow-lg { $selectedIdea?.id === idea.id
								? 'border-violet-500/50 bg-slate-900/40 shadow-lg'
								: 'border-slate-800/50' }"
							type="button"
							onclick={() => handleSelectIdea(idea)}
						>
							<div class="flex items-center gap-2">
								<div
									class="h-3 w-3 rounded-full ring-2 ring-slate-700"
									style={`background-color: ${idea.color ?? '#8b5cf6'}`}
								></div>
								<span class="font-semibold text-slate-200 group-hover:text-violet-300 transition-colors">
									{idea.title}
								</span>
							</div>
							{#if idea.description}
								<p class="mt-1 line-clamp-2 text-xs text-slate-400">{idea.description}</p>
							{/if}
						</button>
					{/each}
				{/if}
			</div>
		</aside>

			<!-- Canvas -->
			<section
				class="relative min-h-[600px] overflow-hidden rounded-xl border border-slate-800/50 bg-slate-900/40 shadow-xl backdrop-blur-sm"
			>
			{#if !$selectedIdea}
				<div class="absolute inset-0 flex items-center justify-center">
					<div class="text-center">
						<h3 class="text-lg font-semibold text-slate-200">No Idea Selected</h3>
						<p class="mt-2 text-sm text-slate-500">
							Select an idea from the sidebar to view and edit its canvas
						</p>
					</div>
				</div>
			{:else if $isLoading}
				<div class="absolute inset-0 z-10 flex items-center justify-center bg-slate-950/70">
					<div class="flex items-center gap-3 text-sm text-slate-300">
						<span class="h-3 w-3 animate-ping rounded-full bg-primary"></span>
						Loading canvas…
					</div>
				</div>
			{:else}
				<SvelteFlow
					bind:nodes={flowNodes}
					bind:edges={flowEdges}
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
						Tip: Drag nodes to reposition. Connect them by drawing from one handle to another (4 directions supported).
					</Panel>
				</SvelteFlow>
			{/if}
		</section>

			<!-- Details Panel -->
			<aside
				class="flex flex-col gap-6 rounded-xl border border-slate-800/50 bg-slate-900/40 p-5 shadow-xl backdrop-blur-sm"
			>
			{#if $selectedNode && $nodeEditForm}
				<!-- Node Details -->
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-slate-100">Node Details</h2>
					<button
						class="rounded-lg border border-red-700/70 px-2.5 py-1 text-xs font-medium text-red-300 transition hover:border-red-500 hover:bg-red-500/10"
						type="button"
						onclick={() => handleDeleteNode($selectedNode.id)}
					>
						Delete
					</button>
				</div>
				<form class="flex flex-col gap-4 text-sm" onsubmit={(e) => { e.preventDefault(); handleNodeSave(); }}>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
						<input
							value={$nodeEditForm.title}
							oninput={handleNodeInput((value) => updateNodeEditFormField('title', value))}
							class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Node title"
							required
						/>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
						<textarea
							value={$nodeEditForm.description}
							oninput={handleNodeInput((value) => updateNodeEditFormField('description', value))}
							class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Node description…"
						></textarea>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Type</span>
						<input
							value={$nodeEditForm.type}
							oninput={handleNodeInput((value) => updateNodeEditFormField('type', value))}
							class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="default"
						/>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Color</span>
						<input
							type="color"
							value={$nodeEditForm.color}
							oninput={handleNodeInput((value) => updateNodeEditFormField('color', value))}
							class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
						/>
					</label>
					<button
						class="mt-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-70"
						type="submit"
						disabled={$isSaving}
					>
						{$isSaving ? 'Saving…' : 'Save changes'}
					</button>
				</form>

				{#if $showDocumentsPanel}
					<div class="mt-4 space-y-3 border-t border-slate-800 pt-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-semibold text-slate-200">Documents</h3>
						</div>
						{#if nodeDocuments.length === 0}
							<p class="text-xs text-slate-500">No documents for this node yet.</p>
						{:else}
							<div class="space-y-2">
								{#each nodeDocuments as doc}
									<div class="rounded-lg border border-slate-800/70 bg-slate-900/60 p-3 text-xs">
										<h4 class="font-medium text-slate-200">{doc.title}</h4>
										<div class="mt-2 max-h-32 overflow-y-auto">
											<MarkdownRenderer content={doc.content} className="text-xs" />
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/if}
			{:else if $selectedIdea && $editForm}
				<!-- Idea Details -->
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-slate-100">Idea Details</h2>
					<div class="flex items-center gap-2">
						<span class="rounded-full border border-slate-700/70 px-2 py-0.5 text-xs text-slate-400"
							>Version {$selectedIdea.version}</span
						>
						{#if $selectedIdea.slug}
							<a
								class="rounded-lg border border-slate-700/70 px-2.5 py-1 text-xs font-medium text-slate-300 transition hover:border-primary hover:text-primary"
								href={`/ideas/${$selectedIdea.slug}`}
							>
								View details
							</a>
						{/if}
						<button
							class="rounded-lg border border-slate-700/70 px-2.5 py-1 text-xs font-medium text-primary transition hover:border-primary hover:text-primary disabled:cursor-not-allowed disabled:opacity-60"
							type="button"
							onclick={startIdeaChat}
							disabled={$isCreatingChat}
						>
							{$isCreatingChat ? 'Creating…' : 'Create chat'}
						</button>
					</div>
				</div>
				<form class="flex flex-col gap-4 text-sm" onsubmit={(e) => { e.preventDefault(); handleIdeaSave(); }}>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
						<input
							value={$editForm.title}
							oninput={handleInput((value) => updateEditFormField('title', value))}
							class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Idea title"
							required
						/>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
						<textarea
							value={$editForm.description}
							oninput={handleInput((value) => updateEditFormField('description', value))}
							class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Describe the problem or solution…"
						></textarea>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Color</span>
						<input
							type="color"
							value={$editForm.color}
							oninput={handleInput((value) => updateEditFormField('color', value))}
							class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
						/>
					</label>
					<button
						class="mt-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-70"
						type="submit"
						disabled={$isSaving}
					>
						{$isSaving ? 'Saving…' : 'Save changes'}
					</button>
				</form>

				{#if $showDocumentsPanel}
					<div class="mt-4 space-y-3 border-t border-slate-800 pt-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-semibold text-slate-200">Idea Documents</h3>
						</div>
						{#if nodeDocuments.length === 0}
							<p class="text-xs text-slate-500">No documents for this idea yet.</p>
						{:else}
							<div class="space-y-2">
								{#each nodeDocuments as doc}
									<div class="rounded-lg border border-slate-800/70 bg-slate-900/60 p-3 text-xs">
										<h4 class="font-medium text-slate-200">{doc.title}</h4>
										<div class="mt-2 max-h-32 overflow-y-auto">
											<MarkdownRenderer content={doc.content} className="text-xs" />
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/if}

				<div class="space-y-3 border-t border-slate-800 pt-4">
					<h3 class="text-sm font-semibold text-slate-200">Recent versions</h3>
					{#if $versions.length === 0}
						<p class="text-xs text-slate-500">Changes will show up here as you iterate.</p>
					{:else}
						<ul class="flex flex-col gap-2 text-xs text-slate-300">
							{#each $versions as version}
								<li class="rounded-lg border border-slate-800/70 bg-slate-900/60 px-3 py-2">
									<div class="flex items-center justify-between">
										<span class="font-medium text-slate-200">v{version.version}</span>
										<span class="text-slate-500">{new Date(version.created_at).toLocaleString()}</span>
									</div>
									<p class="mt-1 text-slate-400 capitalize">{version.change_type}</p>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			{:else}
				<div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
					<h2 class="text-lg font-semibold text-slate-200">Select an Idea</h2>
					<p class="max-w-xs text-sm text-slate-500">
						Choose an idea from the sidebar to view its canvas and add nodes.
					</p>
				</div>
			{/if}
			</aside>
		</div>
	</div>

	<!-- Create Idea Modal -->
	{#if $showCreateModal}
		<div class="fixed inset-0 z-30 flex items-center justify-center bg-slate-950/80 backdrop-blur">
			<form
				class="w-full max-w-md space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
				onsubmit={(e) => { e.preventDefault(); handleCreateIdea(); }}
			>
				<header class="flex items-center justify-between">
					<div>
						<h3 class="text-lg font-semibold text-slate-100">Create new idea</h3>
						<p class="text-xs text-slate-400">Start with a clear name and optional description.</p>
					</div>
					<button
						class="rounded-full border border-slate-700/70 px-2 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
						type="button"
						onclick={() => showCreateModal.set(false)}
					>
						Close
					</button>
				</header>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
					<input
						value={$newIdea.title}
						oninput={handleInput((value) => updateNewIdeaField('title', value))}
						class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Idea title"
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
					<textarea
						value={$newIdea.description}
						oninput={handleInput((value) => updateNewIdeaField('description', value))}
						class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Optional context"
					></textarea>
				</label>
				<label class="flex items-center justify-between text-sm text-slate-300">
					<span>Color</span>
					<input
						type="color"
						value={$newIdea.color}
						oninput={handleInput((value) => updateNewIdeaField('color', value))}
						class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
					/>
				</label>
				<div class="flex items-center justify-end gap-3 pt-2">
					<button
						class="rounded-lg border border-slate-700/70 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
						type="button"
						onclick={() => showCreateModal.set(false)}
					>
						Cancel
					</button>
					<button
						class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
						type="submit"
						disabled={$isSaving}
					>
						{$isSaving ? 'Creating…' : 'Create idea'}
					</button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Create Node Modal -->
	{#if $showCreateNodeModal && $selectedIdea}
		<div class="fixed inset-0 z-30 flex items-center justify-center bg-slate-950/80 backdrop-blur">
			<form
				class="w-full max-w-md space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
				onsubmit={(e) => { e.preventDefault(); handleCreateNode(); }}
			>
				<header class="flex items-center justify-between">
					<div>
						<h3 class="text-lg font-semibold text-slate-100">Create new node</h3>
						<p class="text-xs text-slate-400">Add a node to the {$selectedIdea.title} canvas.</p>
					</div>
					<button
						class="rounded-full border border-slate-700/70 px-2 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
						type="button"
						onclick={() => showCreateNodeModal.set(false)}
					>
						Close
					</button>
				</header>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
					<input
						value={$newNode.title}
						oninput={handleNodeInput((value) => updateNewNodeField('title', value))}
						class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Node title"
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
					<textarea
						value={$newNode.description}
						oninput={handleNodeInput((value) => updateNewNodeField('description', value))}
						class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Optional description"
					></textarea>
				</label>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Type</span>
					<input
						value={$newNode.type}
						oninput={handleNodeInput((value) => updateNewNodeField('type', value))}
						class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="default"
					/>
				</label>
				<label class="flex items-center justify-between text-sm text-slate-300">
					<span>Color</span>
					<input
						type="color"
						value={$newNode.color}
						oninput={handleNodeInput((value) => updateNewNodeField('color', value))}
						class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
					/>
				</label>
				<div class="flex items-center justify-end gap-3 pt-2">
					<button
						class="rounded-lg border border-slate-700/70 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
						type="button"
						onclick={() => showCreateNodeModal.set(false)}
					>
						Cancel
					</button>
					<button
						class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
						type="submit"
						disabled={$isSaving}
					>
						{$isSaving ? 'Creating…' : 'Create node'}
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
