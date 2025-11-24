<script lang="ts">
	import { Zap, ChevronDown, ChevronUp, AlertCircle, CheckCircle, TrendingUp, Code2 } from 'lucide-svelte';
	import type { ProblemSolution } from '$lib/types/problem-solution';

	interface Props {
		solutions: ProblemSolution[];
		loading: boolean;
	}

	let { solutions = [], loading = false }: Props = $props();

	let expandedSolution = $state<string | null>(null);

	function toggleSolution(id: string) {
		expandedSolution = expandedSolution === id ? null : id;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-yellow-500"></div>
		</div>
	{:else if solutions.length === 0}
		<div class="text-center py-20">
			<Zap class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No problem solutions available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="space-y-6">
			{#each solutions as solution}
				<div
					class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-yellow-500/50 transition-all duration-300"
				>
					<!-- Header -->
					<div class="flex items-start justify-between mb-4">
						<div class="flex items-center gap-3 flex-1">
							<div
								class="w-12 h-12 bg-gradient-to-br from-yellow-600 to-orange-600 rounded-lg flex items-center justify-center flex-shrink-0"
							>
								<Zap class="w-6 h-6 text-white" />
							</div>
							<div class="flex-1 min-w-0">
								<h3 class="text-xl font-bold text-white mb-1 line-clamp-2">{solution.problem}</h3>
								<p class="text-sm text-gray-400 line-clamp-1">{solution.context}</p>
							</div>
						</div>
						<button
							onclick={() => toggleSolution(solution.id)}
							class="text-gray-400 hover:text-white transition-colors flex-shrink-0 ml-2"
							aria-label={expandedSolution === solution.id ? 'Collapse' : 'Expand'}
						>
							{#if expandedSolution === solution.id}
								<ChevronUp class="w-5 h-5" />
							{:else}
								<ChevronDown class="w-5 h-5" />
							{/if}
						</button>
					</div>

					<!-- Problem & Solution Preview -->
					<div class="grid md:grid-cols-2 gap-4 mb-4">
						<div class="bg-red-600/10 rounded-lg p-4 border border-red-500/30">
							<div class="flex items-center gap-2 mb-2">
								<AlertCircle class="w-4 h-4 text-red-400" />
								<h4 class="text-sm font-semibold text-red-400">Problem</h4>
							</div>
							<p class="text-sm text-gray-300 line-clamp-3">{solution.problem}</p>
						</div>
						<div class="bg-green-600/10 rounded-lg p-4 border border-green-500/30">
							<div class="flex items-center gap-2 mb-2">
								<CheckCircle class="w-4 h-4 text-green-400" />
								<h4 class="text-sm font-semibold text-green-400">Solution</h4>
							</div>
							<p class="text-sm text-gray-300 line-clamp-3">{solution.solution}</p>
						</div>
					</div>

					<!-- Technologies Preview -->
					{#if solution.technologies && solution.technologies.length > 0}
						<div class="mb-4">
							<div class="flex flex-wrap gap-2">
								{#each solution.technologies.slice(0, 5) as tech}
									<span
										class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
									>
										{tech}
									</span>
								{/each}
								{#if solution.technologies.length > 5}
									<span
										class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
									>
										+{solution.technologies.length - 5} more
									</span>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Expanded Details -->
					{#if expandedSolution === solution.id}
						<div class="mt-4 space-y-4 pt-4 border-t border-gray-700">
							<!-- Full Problem -->
							<div>
								<h4 class="text-sm font-semibold text-red-400 mb-2 flex items-center gap-2">
									<AlertCircle class="w-4 h-4" />
									Problem
								</h4>
								<p class="text-sm text-gray-300 leading-relaxed">{solution.problem}</p>
							</div>

							<!-- Context -->
							<div>
								<h4 class="text-sm font-semibold text-gray-400 mb-2">Context</h4>
								<p class="text-sm text-gray-300 leading-relaxed">{solution.context}</p>
							</div>

							<!-- Full Solution -->
							<div>
								<h4 class="text-sm font-semibold text-green-400 mb-2 flex items-center gap-2">
									<CheckCircle class="w-4 h-4" />
									Solution
								</h4>
								<p class="text-sm text-gray-300 leading-relaxed">{solution.solution}</p>
							</div>

							<!-- Technologies -->
							{#if solution.technologies && solution.technologies.length > 0}
								<div>
									<h4 class="text-sm font-semibold text-blue-400 mb-2 flex items-center gap-2">
										<Code2 class="w-4 h-4" />
										Technologies Used
									</h4>
									<div class="flex flex-wrap gap-2">
										{#each solution.technologies as tech}
											<span
												class="px-3 py-1 text-sm rounded bg-blue-600/20 text-blue-300 border border-blue-500/30"
											>
												{tech}
											</span>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Impact -->
							{#if solution.impact}
								<div>
									<h4 class="text-sm font-semibold text-yellow-400 mb-2 flex items-center gap-2">
										<TrendingUp class="w-4 h-4" />
										Impact
									</h4>
									<p class="text-sm text-gray-300 leading-relaxed">{solution.impact}</p>
								</div>
							{/if}

							<!-- Metrics -->
							{#if solution.metrics}
								<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 rounded-lg p-4 border border-gray-700">
									<h4 class="text-sm font-semibold text-purple-400 mb-3 flex items-center gap-2">
										<TrendingUp class="w-4 h-4" />
										Metrics
									</h4>
									<div class="grid md:grid-cols-3 gap-4">
										<div>
											<p class="text-xs text-gray-400 mb-1">Before</p>
											<p class="text-sm text-red-300 font-medium">{solution.metrics.before}</p>
										</div>
										<div>
											<p class="text-xs text-gray-400 mb-1">After</p>
											<p class="text-sm text-green-300 font-medium">{solution.metrics.after}</p>
										</div>
										<div>
											<p class="text-xs text-gray-400 mb-1">Improvement</p>
											<p class="text-sm text-yellow-300 font-medium">
												{solution.metrics.improvement}
											</p>
										</div>
									</div>
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

