<script lang="ts">
	import { derived } from 'svelte/store';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	import WorkspaceHeader from '../_components/WorkspaceHeader.svelte';
	import AuthNotice from '../_components/AuthNotice.svelte';
	import ErrorBanner from '../_components/ErrorBanner.svelte';
	import DuplicateProjectForm from '../_components/DuplicateProjectForm.svelte';
	import ActiveProjectPanel from '../_components/ActiveProjectPanel.svelte';
	import KanbanBoardSection from '../_components/KanbanBoardSection.svelte';
	import MilestonesSection from '../_components/MilestonesSection.svelte';
	import DependenciesSection from '../_components/DependenciesSection.svelte';
	import type { UUID } from '$lib/api/types';

	import { createProjectDetailLogic } from './project-detail.logic';

	const slugStore = derived(page, ($page) => $page.params.slug ?? null);

	const {
		isAuthenticated,
		error,
	projectQuery,
		activeProject,
		board,
		milestones,
		dependencies,
		columnForm,
		cardForm,
		milestoneForm,
		dependencyForm,
		duplicateForm,
		loadProjects,
		updateColumnFormField,
		updateCardFormField,
		updateMilestoneFormField,
		updateDependencyFormField,
		updateDuplicateFormField,
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
	} = createProjectDetailLogic(slugStore);

	const handleDuplicateAndNavigate = async () => {
		if (!$activeProject) return;
		const duplicate = await handleDuplicateProject(($activeProject.id ?? '') as UUID);
	if (duplicate?.slug) {
		await goto(`/projects/${duplicate.slug}`);
		}
	};
</script>

<section class="space-y-6">
	<WorkspaceHeader isAuthenticated={$isAuthenticated} onRefresh={loadProjects} />

	{#if $error}
		<ErrorBanner message={$error} />
	{/if}

{#if !$isAuthenticated}
		<AuthNotice />
{:else if $projectQuery.isLoading}
	<div class="rounded border border-slate-800 bg-slate-950/60 p-6 text-center text-sm text-slate-300">
		Loading project details...
	</div>
{:else if !$activeProject}
		<div class="rounded border border-slate-800 bg-slate-950/60 p-6 text-center text-sm text-slate-300">
			Project not found. Return to the <a class="text-indigo-400 underline" href="/projects">projects list</a>.
		</div>
	{:else}
		<div class="space-y-6">
			<div class="grid gap-6 lg:grid-cols-[1.05fr_2fr]">
				<div class="space-y-6">
					<DuplicateProjectForm
						formState={$duplicateForm}
						statusOptions={statusOptions}
						activeProjectName={$activeProject?.name ?? ''}
						onFieldChange={updateDuplicateFormField}
						onSubmit={handleDuplicateAndNavigate}
					/>
				</div>

				<div class="space-y-6">
					<ActiveProjectPanel
						project={$activeProject}
						statusOptions={statusOptions}
						onStatusChange={(status) => updateActiveProjectField('status', status)}
						onSaveStatus={saveProjectStatus}
						onMetricChange={(field, value) => updateActiveProjectField(field, value)}
						onSaveMetrics={saveProjectMetrics}
					/>
				</div>
			</div>

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
</section>


