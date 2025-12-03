<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Lightbulb, ExternalLink, Calendar, ArrowLeft, CheckCircle, TrendingUp, Star } from 'lucide-svelte';
	import { getProblemSolution } from '$lib/api/problem-solutions';
	import type { ProblemSolution } from '$lib/types/problem-solution';

	let solution: ProblemSolution | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const solutionId = $derived($page.params.slug);

	onMount(async () => {
		if (solutionId) {
			await fetchSolution(solutionId);
		}
	});

	async function fetchSolution(id: string) {
		loading = true;
		error = null;
		try {
			solution = await getProblemSolution(id);
			if (!solution) {
				error = 'Problem solution not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch problem solution';
			console.error('Error fetching problem solution:', err);
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !solution}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">Problem Solution Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The problem solution you are looking for does not exist.'}</p>
				<a
					href="/problem-solutions"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to Problem Solutions
				</a>
			</div>
		</div>
	{:else}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/problem-solutions"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to Problem Solutions
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-yellow-500/0 via-orange-500/0 to-red-500/0 hover:from-yellow-500/5 hover:via-orange-500/5 hover:to-red-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<div class="flex items-center gap-3 mb-4">
								<div
									class="w-12 h-12 bg-gradient-to-br from-yellow-600 to-orange-600 rounded-lg flex items-center justify-center"
								>
									<Lightbulb class="w-6 h-6 text-white" />
								</div>
								{#if solution.featured}
									<span class="px-3 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded">
										⭐ Featured
									</span>
								{/if}
							</div>

							{#if solution.problem}
								<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
									{solution.problem}
								</h1>
							{/if}

							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if solution.updatedAt}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Updated {formatDate(solution.updatedAt)}</span>
									</div>
								{/if}
							</div>
						</div>

						<!-- Context -->
						{#if solution.context}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Context</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed">{solution.context}</p>
								</div>
							</div>
						{/if}

						<!-- Solution -->
						{#if solution.solution}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Solution</h2>
								<div class="bg-green-900/20 rounded-lg p-6 border border-green-700/30">
									<div class="flex items-start gap-3">
										<CheckCircle class="w-6 h-6 text-green-400 flex-shrink-0 mt-0.5" />
										<p class="text-gray-300 leading-relaxed">{solution.solution}</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Metrics -->
						{#if solution.metrics}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Metrics</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<div class="grid md:grid-cols-3 gap-4">
										<div class="bg-red-900/20 rounded-lg p-4 border border-red-700/30">
											<p class="text-sm text-gray-400 mb-1">Before</p>
											<p class="text-xl font-bold text-white">{solution.metrics.before}</p>
										</div>
										<div class="bg-green-900/20 rounded-lg p-4 border border-green-700/30">
											<p class="text-sm text-gray-400 mb-1">After</p>
											<p class="text-xl font-bold text-white">{solution.metrics.after}</p>
										</div>
										<div class="bg-blue-900/20 rounded-lg p-4 border border-blue-700/30">
											<p class="text-sm text-gray-400 mb-1">Improvement</p>
											<p class="text-xl font-bold text-white">{solution.metrics.improvement}</p>
										</div>
									</div>
								</div>
							</div>
						{/if}

						<!-- Impact -->
						{#if solution.impact}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Impact</h2>
								<div class="bg-blue-900/20 rounded-lg p-6 border border-blue-700/30">
									<div class="flex items-start gap-3">
										<TrendingUp class="w-6 h-6 text-blue-400 flex-shrink-0 mt-0.5" />
										<p class="text-gray-300 leading-relaxed">{solution.impact}</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Technologies -->
						{#if solution.technologies && solution.technologies.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Technologies Used</h2>
								<div class="flex flex-wrap gap-2">
									{#each solution.technologies as tech}
										<span
											class="px-3 py-1 bg-yellow-600/20 text-yellow-300 text-sm rounded-lg border border-yellow-500/30"
										>
											{tech}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
