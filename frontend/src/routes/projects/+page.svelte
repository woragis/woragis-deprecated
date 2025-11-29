<script lang="ts">
	import { createProjectsListLogic } from './projects-list.logic';
	import AuthNotice from './_components/AuthNotice.svelte';
	import ErrorBanner from './_components/ErrorBanner.svelte';
	import CreateProjectForm from './_components/CreateProjectForm.svelte';
	import ProjectsHero from './_components/ProjectsHero.svelte';
	import ProjectsStatsBar from './_components/ProjectsStatsBar.svelte';
	import ProjectsGrid from './_components/ProjectsGrid.svelte';
	import CollapsibleSection from '$lib/components/CollapsibleSection.svelte';

	const {
		isAuthenticated,
		error,
		projects,
		projectForm,
		loadProjects,
		handleCreateProject,
		updateProjectFormField,
		statusOptions
	} = createProjectsListLogic();

	let showCreateForm = false;

	function toggleCreateForm() {
		showCreateForm = !showCreateForm;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<ProjectsHero
		isAuthenticated={$isAuthenticated}
		{showCreateForm}
		onToggleCreateForm={toggleCreateForm}
		onRefresh={loadProjects}
	/>

	<!-- Main Content -->
	<div class="mx-auto max-w-7xl px-6 py-8 lg:px-8">
		{#if $error}
			<div class="mb-6">
				<ErrorBanner message={$error} />
			</div>
		{/if}

		{#if !$isAuthenticated}
			<div class="mx-auto max-w-2xl">
				<AuthNotice />
			</div>
		{:else}
			<CollapsibleSection open={showCreateForm}>
				<CreateProjectForm
					formState={$projectForm}
					statusOptions={statusOptions}
					onFieldChange={updateProjectFormField}
					onSubmit={handleCreateProject}
				/>
			</CollapsibleSection>

			<!-- Projects Section -->
			<div class="space-y-6">
				<ProjectsStatsBar projects={$projects} />
				<ProjectsGrid projects={$projects} />
			</div>
		{/if}
	</div>
</div>

