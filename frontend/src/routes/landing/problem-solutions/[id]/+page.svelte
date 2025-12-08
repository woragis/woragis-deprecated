<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Lightbulb, Code } from 'lucide-svelte';
	import { getProblemSolution, type ProblemSolution } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let solution: ProblemSolution | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const solutionId = $page.params.id;
		if (!solutionId) {
			error = 'Problem solution ID is required';
			loading = false;
			return;
		}

		try {
			solution = await getProblemSolution(solutionId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load problem solution';
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
		title="Problem Solution"
		description={solution?.problem}
		gradientFrom="from-green-950/30"
		gradientVia="via-emerald-950/30"
		gradientTo="to-green-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-green-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/problem-solutions')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Problem Solutions
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-green-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading problem solution...</p>
			</div>
		{:else if error || !solution}
			<EmptyState
				title={error || 'Problem solution not found'}
				description="The problem solution you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-green-700"
					onclick={() => goto('/landing/problem-solutions')}
				>
					Return to Problem Solutions
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Created"
						value={formatDate(solution.created_at)}
						accentColor="green"
					/>
					<StatCard
						label="Featured"
						value={solution.featured ? 'Yes' : 'No'}
						accentColor={solution.featured ? 'emerald' : 'slate'}
					/>
					{#if solution.technologies && solution.technologies.length > 0}
						<StatCard
							label="Technologies"
							value={solution.technologies.length.toString()}
							accentColor="blue"
						/>
					{/if}
				</div>

				<!-- Problem Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Problem</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{solution.problem}</div>
					</div>
				</div>

				<!-- Context Card -->
				{#if solution.context}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Context</h2>
						<div class="prose prose-invert max-w-none">
							<div class="text-slate-300 whitespace-pre-wrap">{solution.context}</div>
						</div>
					</div>
				{/if}

				<!-- Solution Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Solution</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{solution.solution}</div>
					</div>
				</div>

				<!-- Technologies & Impact -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Details</h2>
					<div class="space-y-4">
						{#if solution.technologies && solution.technologies.length > 0}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-2">Technologies</p>
								<div class="flex flex-wrap gap-2">
									{#each solution.technologies as tech}
										<span
											class="inline-flex items-center rounded-full bg-blue-500/20 px-3 py-1 text-xs font-medium text-blue-300"
										>
											{tech}
										</span>
									{/each}
								</div>
							</div>
						{/if}
						{#if solution.impact}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-2">Impact</p>
								<div class="prose prose-invert max-w-none">
									<div class="text-slate-300 whitespace-pre-wrap">{solution.impact}</div>
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

