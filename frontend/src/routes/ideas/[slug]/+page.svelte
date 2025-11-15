<script lang="ts">
	import { derived } from 'svelte/store';
	import { page } from '$app/stores';

	import { createIdeaDetailLogic } from './idea-detail.logic';

	const slugStore = derived(page, ($page) => $page.params.slug ?? null);

	const {
		ideaQuery,
		idea,
		relatedChats,
		relatedProjects
	} = createIdeaDetailLogic(slugStore);

	const formatDate = (value?: string | null) => (value ? new Date(value).toLocaleString() : '—');
</script>

<svelte:head>
	<title>{#if $idea}{$idea.title} · Idea Details · Woragis{:else}Idea Details · Woragis{/if}</title>
</svelte:head>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-4">
		<div>
			<p class="text-sm text-slate-400">Idea Details</p>
			<h1 class="text-3xl font-semibold text-white">
				{$idea?.title ?? 'Loading idea…'}
			</h1>
		</div>

		<a
			class="rounded-lg border border-slate-700 px-4 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:text-white"
			href="/ideas"
		>
			Back to Canvas
		</a>
	</div>

	{#if $ideaQuery.isLoading}
		<div class="rounded-lg border border-slate-800 bg-slate-950/60 px-4 py-8 text-center text-sm text-slate-300">
			Loading idea details…
		</div>
	{:else if $ideaQuery.isError}
		<div class="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			Unable to load idea details right now. Please try again.
		</div>
	{:else if !$idea}
		<div class="rounded-lg border border-yellow-500/40 bg-yellow-500/10 px-4 py-3 text-sm text-yellow-100">
			This idea could not be found. Return to the <a class="underline" href="/ideas">ideas canvas</a>.
		</div>
	{:else}
		<div class="space-y-6">
			<div class="grid gap-6 lg:grid-cols-[1.4fr_1fr]">
				<div class="rounded-2xl border border-slate-800 bg-slate-950/70 p-6">
					<div class="flex flex-wrap items-center justify-between gap-4">
						<div>
							<h2 class="text-lg font-semibold text-white">Overview</h2>
							<p class="text-sm text-slate-400">Slug: {$idea.slug}</p>
						</div>
						<span
							class="rounded-full border border-white/20 px-3 py-1 text-xs uppercase tracking-wide text-slate-200"
							style={`background-color:${$idea.color ?? '#1e293b'}20;border-color:${$idea.color ?? '#1e293b'}40;color:${$idea.color ?? '#94a3b8'}`}
							>{new Date($idea.updated_at ?? $idea.created_at ?? '').toLocaleDateString()}</span
						>
					</div>
					<p class="mt-4 text-slate-200">
						{$idea.description ?? 'No description yet. Update this idea from the canvas to add more context.'}
					</p>
				</div>

				<div class="rounded-2xl border border-slate-800 bg-slate-950/70 p-6">
					<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-400">Metadata</h3>
					<dl class="mt-4 space-y-3 text-sm text-slate-200">
						<div class="flex justify-between gap-4">
							<dt class="text-slate-500">Version</dt>
							<dd class="font-medium text-white">v{$idea.version}</dd>
						</div>
						<div class="flex justify-between gap-4">
							<dt class="text-slate-500">Coordinates</dt>
							<dd>({$idea.pos_x.toFixed(1)}, {$idea.pos_y.toFixed(1)})</dd>
						</div>
						<div class="flex justify-between gap-4">
							<dt class="text-slate-500">Created</dt>
							<dd>{formatDate($idea.created_at)}</dd>
						</div>
						<div class="flex justify-between gap-4">
							<dt class="text-slate-500">Updated</dt>
							<dd>{formatDate($idea.updated_at)}</dd>
						</div>
					</dl>
				</div>
			</div>

			<div class="grid gap-6 lg:grid-cols-2">
				<section class="rounded-2xl border border-slate-800 bg-slate-950/70 p-6">
					<div class="flex items-center justify-between gap-4">
						<h3 class="text-lg font-semibold text-white">Associated Project</h3>
					</div>
					{#if $relatedProjects.length > 0}
						<ul class="mt-4 space-y-3 text-sm text-slate-200">
							{#each $relatedProjects as project}
								<li class="rounded-lg border border-slate-800/70 bg-slate-900/60 px-4 py-3">
									<div class="flex items-center justify-between gap-4">
										<div>
											<p class="font-semibold text-white">{project.name}</p>
											<p class="text-xs text-slate-400">Status: {project.status}</p>
										</div>
										<a
											class="text-xs font-medium text-primary underline"
											href={`/projects/${project.slug ?? project.id}`}
										>
											View project
										</a>
									</div>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="mt-4 text-sm text-slate-400">
							This idea is not linked to a project yet. Open the canvas to assign one.
						</p>
					{/if}
				</section>

				<section class="rounded-2xl border border-slate-800 bg-slate-950/70 p-6">
					<div class="flex items-center justify-between gap-4">
						<h3 class="text-lg font-semibold text-white">Idea Chats</h3>
						<span class="rounded-full bg-slate-900 px-2 py-0.5 text-xs text-slate-400">
							{$relatedChats.length}
						</span>
					</div>
					{#if $relatedChats.length > 0}
						<ul class="mt-4 space-y-3 text-sm text-slate-200">
							{#each $relatedChats as chat}
								<li class="rounded-lg border border-slate-800/70 bg-slate-900/60 px-4 py-3">
									<div class="flex items-center justify-between gap-4">
										<div>
											<p class="font-semibold text-white">{chat.title ?? 'Untitled Conversation'}</p>
											<p class="text-xs text-slate-400">
												Last activity {formatDate(chat.updated_at ?? chat.created_at)}
											</p>
										</div>
										<a
											class="text-xs font-medium text-primary underline"
											href={`/chats?conversation=${chat.id}`}
										>
											Open
										</a>
									</div>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="mt-4 text-sm text-slate-400">
							No chats reference this idea yet. Create one from the canvas sidebar.
						</p>
					{/if}
				</section>
			</div>
		</div>
	{/if}
</section>

<style>
	section {
		color: #e2e8f0;
	}
</style>


