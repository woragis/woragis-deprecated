<script lang="ts">
	import type { Project } from '$lib/api/types';
	import ProjectCard from './ProjectCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	export let projects: Project[] = [];
	export let title = 'All Projects';
</script>

<div>
	<div class="mb-4 flex items-center justify-between">
		<h2 class="text-xl font-semibold text-slate-100">{title}</h2>
		<span class="text-sm text-slate-400">{projects.length} total</span>
	</div>

	{#if projects.length === 0}
		<EmptyState
			title="No projects yet"
			description="Create your first project to get started"
			icon='<svg class="mx-auto h-12 w-12 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>'
		/>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each projects as project (project.id)}
				<ProjectCard {project} />
			{/each}
		</div>
	{/if}
</div>

