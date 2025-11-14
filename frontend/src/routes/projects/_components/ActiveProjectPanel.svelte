<script lang="ts">
	import type { Project, ProjectStatus } from '$lib/api/types';

	export let project: Project | null = null;
	export let statusOptions: ProjectStatus[] = [];
	export let onStatusChange: (status: ProjectStatus) => void;
	export let onSaveStatus: () => void;
	export let onMetricChange: (field: keyof Project, value: number) => void;
	export let onSaveMetrics: () => void;
</script>

{#if project}
	<section class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300">
		<header class="flex flex-wrap items-center justify-between gap-2">
			<div>
				<h3 class="text-sm font-semibold text-slate-100">{project.name}</h3>
				<p class="text-[11px] text-slate-400">Status: {project.status}</p>
			</div>
			<div class="flex gap-2">
				<select
					value={project.status}
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					on:change={(event) => onStatusChange((event.target as HTMLSelectElement).value as ProjectStatus)}
				>
					{#each statusOptions as status (status)}
						<option value={status}>{status}</option>
					{/each}
				</select>
				<button
					class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white"
					on:click={onSaveStatus}
				>
					Save Status
				</button>
			</div>
		</header>
		<div class="mt-4 grid gap-3 md:grid-cols-4">
			<label class="flex flex-col gap-1"
				>MRR<input
					class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
					type="number"
					value={project.mrr}
					on:input={(event) =>
						onMetricChange('mrr', Number((event.target as HTMLInputElement).value))}
				/></label
			>
			<label class="flex flex-col gap-1"
				>CAC<input
					class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
					type="number"
					value={project.cac}
					on:input={(event) =>
						onMetricChange('cac', Number((event.target as HTMLInputElement).value))}
				/></label
			>
			<label class="flex flex-col gap-1"
				>LTV<input
					class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
					type="number"
					value={project.ltv}
					on:input={(event) =>
						onMetricChange('ltv', Number((event.target as HTMLInputElement).value))}
				/></label
			>
			<label class="flex flex-col gap-1"
				>Churn<input
					class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
					type="number"
					value={project.churn_rate}
					on:input={(event) =>
						onMetricChange('churn_rate', Number((event.target as HTMLInputElement).value))}
				/></label
			>
		</div>
		<button
			class="mt-3 rounded bg-emerald-500 px-3 py-2 text-xs font-semibold text-slate-900"
			on:click={onSaveMetrics}
		>
			Save Metrics
		</button>
	</section>
{/if}

