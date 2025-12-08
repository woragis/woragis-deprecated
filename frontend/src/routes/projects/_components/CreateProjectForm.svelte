<script lang="ts">
	import type { ProjectFormState } from '../projects-list.logic';
	import type { ProjectStatus } from '$lib/api/types';

	export let formState: ProjectFormState;
	export let statusOptions: ProjectStatus[] = [];
	export let onFieldChange: <K extends keyof ProjectFormState>(
		field: K,
		value: ProjectFormState[K]
	) => void;
	export let onSubmit: () => void;
</script>

<div class="space-y-3 rounded border border-slate-800 bg-slate-900/60 p-4">
	<h3 class="text-sm font-semibold text-slate-100">Create Project</h3>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Name
		<input
			class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
			value={formState.name}
			on:input={(event) => onFieldChange('name', (event.target as HTMLInputElement).value)}
		/>
	</label>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Description
		<textarea
			class="min-h-[80px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
			value={formState.description}
			on:input={(event) => onFieldChange('description', (event.target as HTMLTextAreaElement).value)}
		></textarea>
	</label>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Status
		<select
			class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
			value={formState.status}
			on:change={(event) => onFieldChange('status', (event.target as HTMLSelectElement).value as ProjectStatus)}
		>
			{#each statusOptions as status (status)}
				<option value={status}>{status}</option>
			{/each}
		</select>
	</label>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Health Score
		<input
			class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
			type="number"
			min="0"
			max="100"
			value={formState.healthScore}
			on:input={(event) =>
				onFieldChange('healthScore', Number((event.target as HTMLInputElement).value))}
		/>
	</label>
	<button
		class="w-full rounded bg-emerald-500 px-3 py-2 text-xs font-semibold text-slate-900"
		on:click={onSubmit}
	>
		Create Project
	</button>
</div>

