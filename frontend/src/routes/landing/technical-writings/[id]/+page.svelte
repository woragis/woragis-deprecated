<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, BookOpen, Globe, ExternalLink } from 'lucide-svelte';
	import { getTechnicalWriting, type TechnicalWriting } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let writing: TechnicalWriting | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const writingId = $page.params.id;
		if (!writingId) {
			error = 'Technical writing ID is required';
			loading = false;
			return;
		}

		try {
			writing = await getTechnicalWriting(writingId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load technical writing';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<PageHero
		title={writing?.title || 'Technical Writing Details'}
		description={writing?.excerpt}
		gradientFrom="from-violet-950/30"
		gradientVia="via-purple-950/30"
		gradientTo="to-violet-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-violet-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/technical-writings')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Technical Writings
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-violet-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading technical writing...</p>
			</div>
		{:else if error || !writing}
			<EmptyState
				title={error || 'Technical writing not found'}
				description="The technical writing you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-violet-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-violet-700"
					onclick={() => goto('/landing/technical-writings')}
				>
					Return to Technical Writings
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Type"
						value={writing.type}
						accentColor="violet"
					/>
					<StatCard
						label="Created"
						value={formatDate(writing.created_at)}
						accentColor="purple"
					/>
					<StatCard
						label="Featured"
						value={writing.featured ? 'Yes' : 'No'}
						accentColor={writing.featured ? 'emerald' : 'slate'}
					/>
				</div>

				<!-- Content Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Content</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{writing.content}</div>
					</div>
				</div>

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Platform</p>
							<p class="text-sm font-medium text-slate-200 capitalize">{writing.platform}</p>
						</div>
						{#if writing.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{writing.slug}</p>
							</div>
						{/if}
						{#if writing.published_at}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Published</p>
								<p class="text-sm font-medium text-slate-200">{formatDate(writing.published_at)}</p>
							</div>
						{/if}
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(writing.updated_at)}</p>
						</div>
						{#if writing.url}
							<div class="md:col-span-2">
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">URL</p>
								<a
									href={writing.url}
									target="_blank"
									class="flex items-center gap-2 text-violet-400 hover:text-violet-300"
								>
									<ExternalLink class="h-4 w-4" />
									<span>{writing.url}</span>
								</a>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

