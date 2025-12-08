<script lang="ts">
	import { derived } from 'svelte/store';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	import ProjectDetailHero from '../_components/ProjectDetailHero.svelte';
	import AuthNotice from '../_components/AuthNotice.svelte';
	import ErrorBanner from '../_components/ErrorBanner.svelte';
	import DuplicateProjectForm from '../_components/DuplicateProjectForm.svelte';
	import ActiveProjectPanel from '../_components/ActiveProjectPanel.svelte';
	import KanbanBoardSection from '../_components/KanbanBoardSection.svelte';
	import MilestonesSection from '../_components/MilestonesSection.svelte';
	import DependenciesSection from '../_components/DependenciesSection.svelte';
	import DocumentationSection from '../_components/DocumentationSection.svelte';
	import CollapsibleSection from '$lib/components/CollapsibleSection.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
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
		documentation,
		documentationSections,
		technologies,
		fileStructures,
		architectureDiagrams,
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
		statusOptions,
		handleCreateDocumentation,
		handleUpdateDocumentationVisibility,
		handleCreateDocumentationSection,
		handleDeleteDocumentationSection,
		handleCreateTechnology,
		handleDeleteTechnology,
		handleCreateFileStructure,
		handleDeleteFileStructure,
		handleCreateArchitectureDiagram,
		handleDeleteArchitectureDiagram
	} = createProjectDetailLogic(slugStore);

	const handleDuplicateAndNavigate = async () => {
		if (!$activeProject) return;
		const duplicate = await handleDuplicateProject(($activeProject.id ?? '') as UUID);
	if (duplicate?.slug) {
		await goto(`/projects/${duplicate.slug}`);
		}
	};
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<ProjectDetailHero
		project={$activeProject}
		isAuthenticated={$isAuthenticated}
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
		{:else if $projectQuery.isLoading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-indigo-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading project details...</p>
			</div>
		{:else if !$activeProject}
			<EmptyState
				title="Project not found"
				description="The project you're looking for doesn't exist or has been removed."
			>
				<a
					href="/projects"
					class="mt-4 inline-block rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-indigo-700"
				>
					Return to Projects
				</a>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<div class="grid gap-6 lg:grid-cols-[1.05fr_2fr]">
					<div class="space-y-6">
						<div
							class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
						>
							<DuplicateProjectForm
								formState={$duplicateForm}
								statusOptions={statusOptions}
								activeProjectName={$activeProject?.name ?? ''}
								onFieldChange={updateDuplicateFormField}
								onSubmit={handleDuplicateAndNavigate}
							/>
						</div>
					</div>

					<div class="space-y-6">
						<div
							class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
						>
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
				</div>

				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
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
				</div>

				<section class="grid gap-6 md:grid-cols-2">
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<MilestonesSection
							milestones={$milestones}
							formState={$milestoneForm}
							onFieldChange={updateMilestoneFormField}
							onAdd={handleAddMilestone}
							onToggleMilestone={handleToggleMilestone}
						/>
					</div>

					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<DependenciesSection
							dependencies={$dependencies}
							formState={$dependencyForm}
							availableProjects={getOtherProjects($activeProject?.id)}
							onFieldChange={updateDependencyFormField}
							onAdd={handleCreateDependency}
							onDelete={handleDeleteDependency}
						/>
					</div>
				</section>

				{#if $activeProject}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<DocumentationSection
							projectId={$activeProject.id}
							documentation={$documentation}
							sections={$documentationSections}
							technologies={$technologies}
							fileStructures={$fileStructures}
							diagrams={$architectureDiagrams}
							isLoading={$projectQuery.isLoading}
							onCreateDocumentation={handleCreateDocumentation}
							onUpdateVisibility={handleUpdateDocumentationVisibility}
							onCreateSection={handleCreateDocumentationSection}
							onDeleteSection={handleDeleteDocumentationSection}
							onCreateTechnology={handleCreateTechnology}
							onDeleteTechnology={handleDeleteTechnology}
							onCreateFileStructure={handleCreateFileStructure}
							onDeleteFileStructure={handleDeleteFileStructure}
							onCreateDiagram={handleCreateArchitectureDiagram}
							onDeleteDiagram={handleDeleteArchitectureDiagram}
						/>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>


