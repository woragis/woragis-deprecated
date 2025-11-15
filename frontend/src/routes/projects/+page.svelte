<script lang="ts">
	import { createProjectsListLogic } from './projects-list.logic';
	import WorkspaceHeader from './_components/WorkspaceHeader.svelte';
	import AuthNotice from './_components/AuthNotice.svelte';
	import ErrorBanner from './_components/ErrorBanner.svelte';
	import CreateProjectForm from './_components/CreateProjectForm.svelte';
	import ProjectsTable from './_components/ProjectsTable.svelte';

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
</script>

<section class="space-y-6">
	<WorkspaceHeader isAuthenticated={$isAuthenticated} onRefresh={loadProjects} />

	{#if $error}
		<ErrorBanner message={$error} />
	{/if}

	{#if !$isAuthenticated}
		<AuthNotice />
	{:else}
		<section class="grid gap-6 lg:grid-cols-[1.05fr_2fr]">
			<div class="space-y-6">
				<CreateProjectForm
					formState={$projectForm}
					statusOptions={statusOptions}
					onFieldChange={updateProjectFormField}
					onSubmit={handleCreateProject}
				/>
			</div>

			<div class="space-y-6">
				<ProjectsTable projects={$projects} />
			</div>
		</section>
	{/if}
</section>

