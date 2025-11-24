<script lang="ts">
	import { Grid, Zap, Code2, TrendingUp } from 'lucide-svelte';
	import { useProblemSolutionMatrixQuery } from '$lib/queries/problem-solutions';
	import type { ProblemSolutionMatrixEntry } from '$lib/types/problem-solution';

	const matrixQuery = useProblemSolutionMatrixQuery();
	let matrix = $derived(matrixQuery.data || []);
	let loading = $derived(matrixQuery.isPending);

	// Sort by count (descending) to show most used technologies first
	let sortedMatrix = $derived.by(() => {
		return [...matrix].sort((a, b) => b.count - a.count);
	});

	function getIntensityColor(count: number, maxCount: number): string {
		if (maxCount === 0) return 'bg-gray-700/30';
		const intensity = count / maxCount;
		if (intensity >= 0.8) return 'bg-green-600/80';
		if (intensity >= 0.6) return 'bg-green-600/60';
		if (intensity >= 0.4) return 'bg-yellow-600/60';
		if (intensity >= 0.2) return 'bg-yellow-600/40';
		return 'bg-yellow-600/20';
	}

	function truncateText(text: string, maxLength: number): string {
		if (text.length <= maxLength) return text;
		return text.slice(0, maxLength).trim() + '...';
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-yellow-500"></div>
		</div>
	{:else if matrix.length === 0}
		<div class="text-center py-20">
			<Grid class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No matrix data available</p>
			<p class="text-gray-500 text-sm">Problem solutions with technologies will appear here</p>
		</div>
	{:else}
		{@const maxCount = sortedMatrix.length > 0 ? Math.max(...sortedMatrix.map((m) => m.count)) : 0}
		<div class="space-y-6">
			<!-- Summary Stats -->
			<div class="grid md:grid-cols-3 gap-4 mb-8">
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<div class="flex items-center gap-3 mb-2">
						<Code2 class="w-6 h-6 text-blue-400" />
						<h3 class="text-lg font-bold text-white">Technologies</h3>
					</div>
					<p class="text-3xl font-bold text-blue-400">{sortedMatrix.length}</p>
					<p class="text-sm text-gray-400">Unique technologies used</p>
				</div>
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<div class="flex items-center gap-3 mb-2">
						<Zap class="w-6 h-6 text-yellow-400" />
						<h3 class="text-lg font-bold text-white">Problems Solved</h3>
					</div>
					<p class="text-3xl font-bold text-yellow-400">
						{sortedMatrix.reduce((sum, m) => sum + m.count, 0)}
					</p>
					<p class="text-sm text-gray-400">Total problem-solution pairs</p>
				</div>
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<div class="flex items-center gap-3 mb-2">
						<TrendingUp class="w-6 h-6 text-green-400" />
						<h3 class="text-lg font-bold text-white">Most Used</h3>
					</div>
					<p class="text-xl font-bold text-green-400">
						{#if sortedMatrix.length > 0}
							{sortedMatrix[0].technology}
						{:else}
							-
						{/if}
					</p>
					<p class="text-sm text-gray-400">
						{#if sortedMatrix.length > 0}
							{sortedMatrix[0].count} {sortedMatrix[0].count === 1 ? 'problem' : 'problems'}
						{:else}
							-
						{/if}
					</p>
				</div>
			</div>

			<!-- Matrix Grid -->
			<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
				<h3 class="text-2xl font-bold text-white mb-6 flex items-center gap-2">
					<Grid class="w-6 h-6" />
					Problem-Solution Matrix
				</h3>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each sortedMatrix as entry}
						<div
							class="bg-gray-800/50 rounded-lg p-5 border border-gray-700 hover:border-yellow-500/50 transition-all duration-300"
						>
							<!-- Technology Header -->
							<div class="flex items-center justify-between mb-3">
								<h4 class="text-lg font-bold text-white">{entry.technology}</h4>
								<span
									class="px-3 py-1 text-sm font-bold rounded-full bg-yellow-600/20 text-yellow-300 border border-yellow-500/30"
								>
									{entry.count}
								</span>
							</div>

							<!-- Problems List -->
							<div class="space-y-2">
								<p class="text-xs font-semibold text-gray-400 mb-2">
									Problems Solved ({entry.count}):
								</p>
								<ul class="space-y-2">
									{#each entry.problems.slice(0, 3) as problem}
										<li class="text-sm text-gray-300 flex items-start gap-2">
											<span class="text-yellow-400 mt-1">•</span>
											<span>{truncateText(problem, 80)}</span>
										</li>
									{/each}
									{#if entry.problems.length > 3}
										<li class="text-xs text-gray-400 italic">
											+{entry.problems.length - 3} more {entry.problems.length - 3 === 1 ? 'problem' : 'problems'}...
										</li>
									{/if}
								</ul>
							</div>

							<!-- Intensity Indicator -->
							<div class="mt-3 pt-3 border-t border-gray-700">
								<div class="flex items-center justify-between text-xs text-gray-400 mb-1">
									<span>Usage</span>
									<span>{entry.count} / {maxCount}</span>
								</div>
								<div class="w-full h-2 bg-gray-700 rounded-full overflow-hidden">
									<div
										class="h-full {getIntensityColor(entry.count, maxCount)} transition-all duration-300"
										style="width: {maxCount > 0 ? (entry.count / maxCount) * 100 : 0}%"
									></div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>

