<script lang="ts">
	import type { Project } from '$lib/api/types';

	export let project: Project;

	function getStatusColor(status: string): string {
		const colors: Record<string, string> = {
			planning: 'bg-blue-500/20 text-blue-300 border-blue-500/30',
			active: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30',
			on_hold: 'bg-amber-500/20 text-amber-300 border-amber-500/30',
			completed: 'bg-purple-500/20 text-purple-300 border-purple-500/30',
			archived: 'bg-slate-500/20 text-slate-300 border-slate-500/30'
		};
		return colors[status] || colors.planning;
	}

	function getHealthColor(score: number): string {
		if (score >= 80) return 'text-emerald-400';
		if (score >= 60) return 'text-amber-400';
		return 'text-red-400';
	}
</script>

<a
	href={`/projects/${project.slug ?? project.id}`}
	class="group relative overflow-hidden rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-6 backdrop-blur-sm transition-all hover:scale-[1.02] hover:border-indigo-500/50 hover:shadow-xl hover:shadow-indigo-500/20"
>
	<!-- Gradient overlay on hover -->
	<div
		class="absolute inset-0 bg-gradient-to-br from-indigo-500/0 to-purple-500/0 transition-all group-hover:from-indigo-500/10 group-hover:to-purple-500/10"
	></div>

	<div class="relative z-10 space-y-4">
		<!-- Header -->
		<div class="flex items-start justify-between">
			<div class="flex-1">
				<h3 class="text-lg font-bold text-white transition-colors group-hover:text-indigo-300">
					{project.name}
				</h3>
				{#if project.description}
					<p class="mt-1 line-clamp-2 text-xs text-slate-400">{project.description}</p>
				{/if}
			</div>
		</div>

		<!-- Status Badge -->
		<div class="flex items-center gap-2">
			<span
				class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium {getStatusColor(
					project.status
				)}"
			>
				{project.status.replace('_', ' ')}
			</span>
		</div>

		<!-- Metrics -->
		<div class="grid grid-cols-2 gap-3 border-t border-slate-800/50 pt-4">
			<div>
				<p class="text-xs text-slate-400">Health Score</p>
				<p class="mt-1 text-lg font-semibold {getHealthColor(project.health_score)}">
					{project.health_score}
				</p>
			</div>
			<div>
				<p class="text-xs text-slate-400">MRR</p>
				<p class="mt-1 text-lg font-semibold text-slate-200">${project.mrr.toFixed(2)}</p>
			</div>
		</div>

		<!-- View Button -->
		<div class="flex items-center justify-end pt-2">
			<span
				class="text-xs font-medium text-indigo-400 transition-colors group-hover:text-indigo-300"
			>
				View Details →
			</span>
		</div>
	</div>

	<!-- Decorative gradient -->
	<div
		class="absolute -right-8 -top-8 h-24 w-24 rounded-full bg-indigo-500/5 blur-2xl transition-all group-hover:bg-indigo-500/10"
	></div>
</a>

