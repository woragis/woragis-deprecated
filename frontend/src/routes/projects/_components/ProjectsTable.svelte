<script lang="ts">
	import type { Project } from '$lib/api/types';

	export let projects: Project[] = [];
</script>

<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
	<h3 class="text-sm font-semibold text-slate-100">Projects ({projects.length})</h3>
	<div class="mt-3 overflow-x-auto text-xs text-slate-200">
		<table class="min-w-full border-separate border-spacing-y-2">
			<thead class="text-[11px] tracking-wide text-slate-400 uppercase">
				<tr>
					<th class="text-left">Name</th>
					<th class="text-left">Status</th>
					<th class="text-left">Health</th>
					<th class="text-left">MRR</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each projects as project (project.id)}
					<tr class="rounded border border-slate-800 bg-slate-950/40">
						<td class="px-3 py-2 font-semibold">{project.name}</td>
						<td class="px-3 py-2 text-slate-300 capitalize">{project.status}</td>
						<td class="px-3 py-2">{project.health_score}</td>
						<td class="px-3 py-2">{project.mrr.toFixed(2)}</td>
						<td class="px-3 py-2 text-right">
							<a
								class="rounded bg-slate-800 px-3 py-1 text-xs transition hover:bg-indigo-500"
								href={`/projects/${project.slug ?? project.id}`}
							>
								View
							</a>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

