<script lang="ts">
	import type { DuplicateFormState } from '../projects.logic';
	import type { ProjectStatus } from '$lib/api/types';

	export let formState: DuplicateFormState;
	export let statusOptions: ProjectStatus[] = [];
	export let activeProjectName = '';
	export let onFieldChange: <K extends keyof DuplicateFormState>(
		field: K,
		value: DuplicateFormState[K]
	) => void;
	export let onSubmit: () => void;
</script>

<div class="space-y-3 rounded border border-slate-800 bg-slate-900/60 p-4">
	<h3 class="text-sm font-semibold text-slate-100">Duplicate Active Project</h3>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Name
		<input
			class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
			value={formState.name}
			placeholder="Copy of ..."
			on:input={(event) => onFieldChange('name', (event.target as HTMLInputElement).value)}
		/>
	</label>
	<label class="flex flex-col gap-1 text-xs text-slate-300">
		Description
		<textarea
			class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
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
	<div class="flex flex-col gap-2 text-xs text-slate-300">
		<label class="flex items-center gap-2">
			<input
				type="checkbox"
				checked={formState.copyBoard}
				on:change={(event) => onFieldChange('copyBoard', (event.target as HTMLInputElement).checked)}
			/>
			Copy board
		</label>
		<label class="flex items-center gap-2">
			<input
				type="checkbox"
				checked={formState.copyMilestones}
				on:change={(event) =>
					onFieldChange('copyMilestones', (event.target as HTMLInputElement).checked)}
			/>
			Copy milestones
		</label>
		<label class="flex items-center gap-2">
			<input
				type="checkbox"
				checked={formState.copyDependencies}
				on:change={(event) =>
					onFieldChange('copyDependencies', (event.target as HTMLInputElement).checked)}
			/>
			Copy dependencies
		</label>
	</div>
	<button
		class="w-full rounded bg-sky-500 px-3 py-2 text-xs font-semibold text-white"
		on:click={onSubmit}
	>
		Duplicate {activeProjectName}
	</button>
</div>

