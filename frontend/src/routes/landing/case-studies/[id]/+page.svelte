<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Briefcase } from 'lucide-svelte';
	import { getCaseStudy, type CaseStudy } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let study: CaseStudy | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const studyId = $page.params.id;
		if (!studyId) {
			error = 'Case study ID is required';
			loading = false;
			return;
		}

		try {
			study = await getCaseStudy(studyId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load case study';
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
		title={study?.title || 'Case Study Details'}
		description={study?.excerpt}
		gradientFrom="from-amber-950/30"
		gradientVia="via-orange-950/30"
		gradientTo="to-amber-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-amber-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/case-studies')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Case Studies
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-amber-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading case study...</p>
			</div>
		{:else if error || !study}
			<EmptyState
				title={error || 'Case study not found'}
				description="The case study you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-amber-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-amber-700"
					onclick={() => goto('/landing/case-studies')}
				>
					Return to Case Studies
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Created"
						value={formatDate(study.created_at)}
						accentColor="amber"
					/>
					<StatCard
						label="Featured"
						value={study.featured ? 'Yes' : 'No'}
						accentColor={study.featured ? 'emerald' : 'slate'}
					/>
					{#if study.project_id}
						<StatCard
							label="Project ID"
							value={study.project_id}
							accentColor="blue"
						/>
					{/if}
				</div>

				<!-- Content Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Content</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{study.content}</div>
					</div>
				</div>

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#if study.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{study.slug}</p>
							</div>
						{/if}
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(study.updated_at)}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

