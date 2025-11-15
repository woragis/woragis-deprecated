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

	import IdeaNode from './IdeaNode.svelte';
import {
		createIdeasLogic,
		type IdeaNodeData,
		type IdeaFormState
	} from './ideas.logic';
	import { goto } from '$app/navigation';

	const nodeTypes: NodeTypes = {
		idea: IdeaNode
	};

	const {
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
	} = createIdeasLogic();
	const startIdeaChat = async () => {
		const conversation = await createIdeaChat();
		if (conversation) {
			goto(`/chats?conversation=${conversation.id}`);
		}
	};

let flowNodes: Node<IdeaNodeData>[] = [];
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
				disabled={$isLoading}
			>
				{$isLoading ? 'Refreshing…' : 'Refresh'}
			</button>
			<button
				class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90"
				type="button"
				on:click={() => showCreateModal.set(true)}
			>
				New Idea
			</button>
		</div>
	</header>

	{#if $uiError}
		<div class="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			{$uiError}
		</div>
	{:else if $ideasQueryError || $linksQueryError}
		<div class="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			Unable to load ideas. Please try again.
		</div>
	{/if}

	<div class="grid gap-6 lg:grid-cols-[2fr_1fr]">
		<section class="relative min-h-[540px] overflow-hidden rounded-2xl border border-slate-800/80 bg-slate-950/60">
			{#if $isLoading}
				<div class="absolute inset-0 z-10 flex items-center justify-center bg-slate-950/70">
					<div class="flex items-center gap-3 text-sm text-slate-300">
						<span class="h-3 w-3 animate-ping rounded-full bg-primary"></span>
						Loading canvas…
					</div>
				</div>
			{/if}

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
					Tip: drag ideas to reposition and connect them by drawing from one node handle to another.
				</Panel>
			</SvelteFlow>
		</section>

		<aside class="flex flex-col gap-6 rounded-2xl border border-slate-800/80 bg-slate-950/60 p-5">
			{#if $selectedIdea && $editForm}
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
							on:click={startIdeaChat}
							disabled={$isCreatingChat}
						>
							{$isCreatingChat ? 'Creating…' : 'Create chat'}
						</button>
					</div>
				</div>
				<form class="flex flex-col gap-4 text-sm" on:submit|preventDefault={handleIdeaSave}>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
						<input
							value={$editForm.title}
							on:input={handleInput((value) => updateEditFormField('title', value))}
							class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Idea title"
							required
						/>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
						<textarea
							value={$editForm.description}
							on:input={handleInput((value) => updateEditFormField('description', value))}
							class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Describe the problem or solution…"
						></textarea>
					</label>
					<label class="flex flex-col gap-2">
						<span class="text-xs uppercase tracking-wide text-slate-400">Color</span>
						<input
							type="color"
							value={$editForm.color}
							on:input={handleInput((value) => updateEditFormField('color', value))}
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

				<div class="space-y-3">
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
									{#if version.title && version.title !== $selectedIdea.title}
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

	<section class="rounded-2xl border border-slate-800/80 bg-slate-950/60 p-5">
		<div class="flex flex-wrap items-center justify-between gap-3">
			<div>
				<h3 class="text-lg font-semibold text-slate-100">Idea chats</h3>
				<p class="text-sm text-slate-400">Recent conversations linked to this idea.</p>
			</div>
			<a
				class="rounded-lg border border-slate-700/70 px-3 py-1.5 text-sm text-slate-200 transition hover:border-slate-500 hover:text-white"
				href="/chats"
				data-sveltekit-preload-data
			>
				Open chats
			</a>
		</div>

		{#if !$selectedIdea}
			<p class="mt-4 text-sm text-slate-500">Select an idea to see its related chats.</p>
		{:else if $conversationsQuery.isFetching}
			<p class="mt-4 text-sm text-slate-400">Loading chats…</p>
		{:else if chatsForIdea.length === 0}
			<p class="mt-4 text-sm text-slate-400">
				No chats linked to <span class="font-semibold text-slate-200">{$selectedIdea.title}</span> yet.
			</p>
		{:else}
			<div class="mt-5 grid gap-4 md:grid-cols-2">
				{#each chatsForIdea as chat}
					<a
						class="rounded-xl border border-slate-800/70 bg-slate-900/60 p-4 text-sm text-slate-200 shadow-inner transition hover:border-primary/60 hover:bg-slate-900"
						href={`/chats?conversation=${chat.id}`}
						data-sveltekit-preload-data
					>
						<div class="flex items-center justify-between gap-2">
							<h4 class="text-base font-semibold text-slate-100">{chat.title}</h4>
							<span class="text-[11px] uppercase text-slate-500">
								{new Date(chat.updated_at ?? chat.created_at ?? new Date()).toLocaleDateString()}
							</span>
						</div>
						<p class="mt-2 line-clamp-2 text-slate-400">
							{chat.description || 'No description provided yet.'}
						</p>
						<div class="mt-3 flex items-center justify-end text-xs uppercase tracking-wide text-primary/80">
							<span class="font-medium">Open chat →</span>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</section>

	{#if $showCreateModal}
		<div class="fixed inset-0 z-30 flex items-center justify-center bg-slate-950/80 backdrop-blur">
			<form
				class="w-full max-w-md space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
				on:submit|preventDefault={handleCreateIdea}
			>
				<header class="flex items-center justify-between">
					<div>
						<h3 class="text-lg font-semibold text-slate-100">Create new idea</h3>
						<p class="text-xs text-slate-400">Start with a clear name and optional description.</p>
					</div>
					<button
						class="rounded-full border border-slate-700/70 px-2 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
						type="button"
						on:click={() => showCreateModal.set(false)}
					>
						Close
					</button>
				</header>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
					<input
						value={$newIdea.title}
						on:input={handleInput((value) => updateNewIdeaField('title', value))}
						class="rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Idea title"
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm">
					<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
					<textarea
						value={$newIdea.description}
						on:input={handleInput((value) => updateNewIdeaField('description', value))}
						class="h-24 rounded-lg border border-slate-700/70 bg-slate-900/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						placeholder="Optional context"
					></textarea>
				</label>
				<label class="flex items-center justify-between text-sm text-slate-300">
					<span>Color</span>
					<input
						type="color"
						value={$newIdea.color}
						on:input={handleInput((value) => updateNewIdeaField('color', value))}
						class="h-10 w-20 cursor-pointer rounded border border-slate-700/70 bg-slate-900/80"
					/>
				</label>
				<div class="flex items-center justify-end gap-3 pt-2">
					<button
						class="rounded-lg border border-slate-700/70 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
						type="button"
						on:click={() => showCreateModal.set(false)}
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
</div>

<style>
	:global(.svelte-flow__attribution) {
		display: none;
	}
</style>

