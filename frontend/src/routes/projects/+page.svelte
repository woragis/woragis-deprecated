<script lang="ts">
	import { createProjectsLogic } from './projects.logic';
	import WorkspaceHeader from './_components/WorkspaceHeader.svelte';
	import AuthNotice from './_components/AuthNotice.svelte';
	import ErrorBanner from './_components/ErrorBanner.svelte';
	import CreateProjectForm from './_components/CreateProjectForm.svelte';
	import DuplicateProjectForm from './_components/DuplicateProjectForm.svelte';
	import ProjectsTable from './_components/ProjectsTable.svelte';
	import ActiveProjectPanel from './_components/ActiveProjectPanel.svelte';
	import KanbanBoardSection from './_components/KanbanBoardSection.svelte';
	import MilestonesSection from './_components/MilestonesSection.svelte';
	import DependenciesSection from './_components/DependenciesSection.svelte';
	import type { UUID } from '$lib/api/types';

	const {
		isAuthenticated,
		error,
		projects,
		activeProject,
		board,
		milestones,
		dependencies,
		projectForm,
		columnForm,
		cardForm,
		milestoneForm,
		dependencyForm,
		duplicateForm,
		loadProjects,
		selectProject,
		updateProjectFormField,
		updateColumnFormField,
		updateCardFormField,
		updateMilestoneFormField,
		updateDependencyFormField,
		updateDuplicateFormField,
		handleCreateProject,
		handleCreateColumn,
		handleCreateCard,
		handleMoveCard,
		handleDeleteCard,
		handleDeleteColumn,
		handleAddMilestone,
		handleToggleMilestone,
		handleCreateDependency,
		handleDeleteDependency,
		handleDuplicateProject,
		saveProjectStatus,
		saveProjectMetrics,
		updateActiveProjectField,
		getOtherProjects,
		statusOptions
	} = createProjectsLogic();
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

				{#if $activeProject}
					<DuplicateProjectForm
						formState={$duplicateForm}
						statusOptions={statusOptions}
						activeProjectName={$activeProject?.name ?? ''}
						onFieldChange={updateDuplicateFormField}
						onSubmit={() => handleDuplicateProject(($activeProject?.id ?? '') as UUID)}
					/>
				{/if}
			</div>

			<div class="space-y-6">
				<ProjectsTable
					projects={$projects}
					activeProjectId={$activeProject?.id ?? null}
					onSelect={selectProject}
				/>

				{#if $activeProject}
					<div class="space-y-6">
						<ActiveProjectPanel
							project={$activeProject}
							statusOptions={statusOptions}
							onStatusChange={(status) => updateActiveProjectField('status', status)}
							onSaveStatus={saveProjectStatus}
							onMetricChange={(field, value) => updateActiveProjectField(field, value)}
							onSaveMetrics={saveProjectMetrics}
						/>

						<KanbanBoardSection
							board={$board}
							columnForm={$columnForm}
							cardForm={$cardForm}
							milestones={$milestones}
							onDeleteColumn={handleDeleteColumn}
							onMoveCard={handleMoveCard}
							onDeleteCard={handleDeleteCard}
							onColumnFieldChange={updateColumnFormField}
							onCardFieldChange={updateCardFormField}
							onAddColumn={handleCreateColumn}
							onAddCard={handleCreateCard}
						/>

						<section class="grid gap-4 md:grid-cols-2">
							<MilestonesSection
								milestones={$milestones}
								formState={$milestoneForm}
								onFieldChange={updateMilestoneFormField}
								onAdd={handleAddMilestone}
								onToggleMilestone={handleToggleMilestone}
							/>

							<DependenciesSection
								dependencies={$dependencies}
								formState={$dependencyForm}
								availableProjects={getOtherProjects($activeProject?.id)}
								onFieldChange={updateDependencyFormField}
								onAdd={handleCreateDependency}
								onDelete={handleDeleteDependency}
							/>
						</section>
					</div>
				{/if}
			</div>
		</section>
	{/if}
</section>

