<script lang="ts">
	import type { Project, ProjectDependency, UUID } from '$lib/api/types';
	import type { DependencyFormState } from '../projects.logic';

	export let dependencies: ProjectDependency[] = [];
	export let formState: DependencyFormState;
	export let availableProjects: Project[] = [];
	export let onFieldChange: <K extends keyof DependencyFormState>(
		field: K,
		value: DependencyFormState[K]
	) => void;
	export let onAdd: () => void;
	export let onDelete: (dependencyId: UUID) => void;
</script>

<div class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300">
	<header class="flex items-center justify-between">
		<h4 class="text-sm font-semibold text-slate-100">Dependencies ({dependencies.length})</h4>
	</header>
	<ul class="mt-3 space-y-2">
		{#each dependencies as dependency (dependency.id)}
			<li class="flex items-center justify-between rounded border border-slate-800 bg-slate-950/60 px-3 py-2">
				<span>{dependency.type} → {dependency.depends_on_project_id}</span>
				<button class="rounded bg-rose-500/70 px-2 py-1 text-white" on:click={() => onDelete(dependency.id)}>
					Remove
				</button>
			</li>
		{/each}
	</ul>
	<form class="mt-4 space-y-2" on:submit|preventDefault={onAdd}>
		<label class="flex flex-col">
			Depends On
			<select
				class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
				value={formState.dependsOnProjectId}
				on:change={(event) =>
					onFieldChange('dependsOnProjectId', (event.target as HTMLSelectElement).value as UUID | '')}
			>
				<option value="">Select project</option>
				{#each availableProjects as project (project.id)}
					<option value={project.id}>{project.name}</option>
				{/each}
			</select>
		</label>
		<label class="flex flex-col">
			Type
			<select
				class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
				value={formState.type}
				on:change={(event) =>
					onFieldChange(
						'type',
						(event.target as HTMLSelectElement).value as DependencyFormState['type']
					)}
			>
				<option value="blocks">blocks</option>
				<option value="relates">relates</option>
				<option value="supports">supports</option>
			</select>
		</label>
		<button class="w-full rounded bg-indigo-500 px-2 py-2 text-xs font-semibold text-white">
			Add Dependency
		</button>
	</form>
</div>

