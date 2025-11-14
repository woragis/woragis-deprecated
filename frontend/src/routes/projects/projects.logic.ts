import { onMount } from 'svelte';
import { get, writable } from 'svelte/store';
import { createQuery } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type {
	KanbanBoard,
	KanbanCard,
	Milestone,
	Project,
	ProjectDependency,
	ProjectStatus,
	UUID
} from '$lib/api/types';
import {
	addMilestone,
	bulkUpdateMilestones,
	createDependency,
	createKanbanCard,
	createKanbanColumn,
	createProject,
	deleteDependency,
	deleteKanbanCard,
	deleteKanbanColumn,
	duplicateProject,
	getKanbanBoard,
	listDependencies,
	listMilestones,
	listProjects,
	moveKanbanCard,
	updateProjectMetrics,
	updateProjectStatus
} from '$lib/api/projects';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export const statusOptions: ProjectStatus[] = ['idea', 'planning', 'executing', 'monitoring', 'completed'];

export interface ProjectFormState {
	name: string;
	description: string;
	status: ProjectStatus;
	healthScore: number;
}

const defaultProjectForm = (): ProjectFormState => ({
	name: '',
	description: '',
	status: 'planning',
	healthScore: 60
});

export interface ColumnFormState {
	name: string;
	position?: number;
	wipLimit?: number;
}

const defaultColumnForm = (): ColumnFormState => ({
	name: '',
	position: undefined,
	wipLimit: undefined
});

export interface CardFormState {
	columnId: UUID | '';
	title: string;
	description: string;
	dueDate: string;
	milestoneId: UUID | '';
	position?: number;
}

const defaultCardForm = (): CardFormState => ({
	columnId: '',
	title: '',
	description: '',
	dueDate: '',
	milestoneId: '',
	position: undefined
});

export interface MilestoneFormState {
	title: string;
	description: string;
	dueDate: string;
}

const defaultMilestoneForm = (): MilestoneFormState => ({
	title: '',
	description: '',
	dueDate: new Date().toISOString().slice(0, 10)
});

export interface DependencyFormState {
	dependsOnProjectId: UUID | '';
	type: 'blocks' | 'relates' | 'supports';
}

const defaultDependencyForm = (): DependencyFormState => ({
	dependsOnProjectId: '',
	type: 'relates'
});

export interface DuplicateFormState {
	name: string;
	description: string;
	status: ProjectStatus;
	copyBoard: boolean;
	copyMilestones: boolean;
	copyDependencies: boolean;
}

const defaultDuplicateForm = (): DuplicateFormState => ({
	name: '',
	description: '',
	status: 'planning',
	copyBoard: true,
	copyMilestones: true,
	copyDependencies: false
});

type FormUpdater<T> = <K extends keyof T>(field: K, value: T[K]) => void;

export function createProjectsLogic() {
	const isAuthenticated = writable(false);
	const loading = writable(false);
	const error = writable<string | null>(null);

	const projects = writable<Project[]>([]);
	const activeProject = writable<Project | null>(null);
	const board = writable<KanbanBoard | null>(null);
	const milestones = writable<Milestone[]>([]);
	const dependencies = writable<ProjectDependency[]>([]);

	const projectForm = writable<ProjectFormState>(defaultProjectForm());
	const columnForm = writable<ColumnFormState>(defaultColumnForm());
	const cardForm = writable<CardFormState>(defaultCardForm());
	const milestoneForm = writable<MilestoneFormState>(defaultMilestoneForm());
	const dependencyForm = writable<DependencyFormState>(defaultDependencyForm());
	const duplicateForm = writable<DuplicateFormState>(defaultDuplicateForm());

	const boardQuery = createQuery<KanbanBoard>(() => ({
		queryKey: ['projects', get(activeProject)?.id, 'board'],
		queryFn: () => getKanbanBoard(get(activeProject)!.id),
		enabled: false,
		onSuccess(data) {
			board.set(data);
			error.set(null);
		},
		onError(err: unknown) {
			error.set(err instanceof Error ? err.message : 'Unable to load kanban board');
		}
	}));

	const milestonesQuery = createQuery<Milestone[]>(() => ({
		queryKey: ['projects', get(activeProject)?.id, 'milestones'],
		queryFn: () => listMilestones(get(activeProject)!.id),
		enabled: false,
		onSuccess(data) {
			milestones.set(data);
			error.set(null);
		},
		onError(err: unknown) {
			error.set(err instanceof Error ? err.message : 'Unable to load milestones');
		}
	}));

	const dependenciesQuery = createQuery<ProjectDependency[]>(() => ({
		queryKey: ['projects', get(activeProject)?.id, 'dependencies'],
		queryFn: () => listDependencies(get(activeProject)!.id),
		enabled: false,
		onSuccess(data) {
			dependencies.set(data);
			error.set(null);
		},
		onError(err: unknown) {
			error.set(err instanceof Error ? err.message : 'Unable to load dependencies');
		}
	}));

	const projectsQuery = createQuery<Project[]>(() => ({
		queryKey: ['projects', 'list'],
		queryFn: listProjects,
		enabled: false,
		onSuccess(data) {
			projects.set(data);
			error.set(null);

			if (data.length === 0) {
				activeProject.set(null);
				board.set(null);
				milestones.set([]);
				dependencies.set([]);
				return;
			}

			const currentActive = get(activeProject);
			if (!currentActive) {
				activeProject.set(data[0]);
			} else {
				const existing = data.find((project) => project.id === currentActive.id);
				activeProject.set(existing ?? data[0]);
			}

			void refetchProjectDetails();
		},
		onError(err: unknown) {
			error.set(err instanceof Error ? err.message : 'Unable to load projects');
		}
	}));

	async function refetchProjectDetails() {
		if (!get(isAuthenticated) || !get(activeProject)) return;
		await Promise.all([boardQuery.refetch(), milestonesQuery.refetch(), dependenciesQuery.refetch()]);
	}

	async function loadProjects() {
		if (!get(isAuthenticated)) {
			const message = 'You must be signed in to load projects.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			await projectsQuery.refetch();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to load projects.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function selectProject(project: Project) {
		activeProject.set(project);
		await refetchProjectDetails();
	}

	onMount(() => {
		const unsubscribe = authStore.subscribe(async (state) => {
			isAuthenticated.set(state.isAuthenticated);

			if (!state.isAuthenticated) {
				resetStateForUnauthenticated();
				return;
			}

			await loadProjects();
		});

		return () => {
			unsubscribe();
		};
	});

	function resetStateForUnauthenticated() {
		projects.set([]);
		activeProject.set(null);
		board.set(null);
		milestones.set([]);
		dependencies.set([]);
		projectForm.set(defaultProjectForm());
		columnForm.set(defaultColumnForm());
		cardForm.set(defaultCardForm());
		milestoneForm.set(defaultMilestoneForm());
		dependencyForm.set(defaultDependencyForm());
		duplicateForm.set(defaultDuplicateForm());
	}

	const updateProjectFormField: FormUpdater<ProjectFormState> = (field, value) => {
		projectForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateColumnFormField: FormUpdater<ColumnFormState> = (field, value) => {
		columnForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateCardFormField: FormUpdater<CardFormState> = (field, value) => {
		cardForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateMilestoneFormField: FormUpdater<MilestoneFormState> = (field, value) => {
		milestoneForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateDependencyFormField: FormUpdater<DependencyFormState> = (field, value) => {
		dependencyForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateDuplicateFormField: FormUpdater<DuplicateFormState> = (field, value) => {
		duplicateForm.update((current) => ({ ...current, [field]: value }));
	};

	async function handleCreateProject() {
		if (!get(isAuthenticated)) {
			const message = 'You must be signed in to create projects.';
			error.set(message);
			toastError(message);
			return;
		}

		const form = get(projectForm);
		if (!form.name.trim()) {
			const message = 'Project name is required.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const created = await createProject({
				name: form.name,
				description: form.description,
				status: form.status,
				healthScore: form.healthScore
			});
			projectForm.set(defaultProjectForm());
			await projectsQuery.refetch();
			const currentProjects = get(projects);
			const nextActive = currentProjects.find((project) => project.id === created.id) ?? created;
			await selectProject(nextActive);
			toastSuccess('Project created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create project.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleCreateColumn() {
		if (!get(isAuthenticated) || !get(activeProject)) {
			const message = 'Select a project before adding columns.';
			error.set(message);
			toastError(message);
			return;
		}

		const form = get(columnForm);
		if (!form.name.trim()) {
			const message = 'Column name is required.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			const updatedBoard = await createKanbanColumn(projectId, {
				name: form.name,
				position: form.position,
				wipLimit: form.wipLimit
			});
			board.set(updatedBoard);
			columnForm.set(defaultColumnForm());
			toastSuccess('Column created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create column.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleCreateCard() {
		if (!get(isAuthenticated) || !get(activeProject)) {
			const message = 'Select a project before adding cards.';
			error.set(message);
			toastError(message);
			return;
		}

		const form = get(cardForm);
		if (!form.columnId || !form.title.trim()) {
			const message = 'Column and title are required.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			const updatedBoard = await createKanbanCard(projectId, {
				columnId: form.columnId,
				title: form.title,
				description: form.description,
				dueDate: form.dueDate ? new Date(form.dueDate).toISOString() : undefined,
				milestoneId: form.milestoneId || undefined,
				position: form.position
			});
			board.set(updatedBoard);
			cardForm.set(defaultCardForm());
			toastSuccess('Card created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create card.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleMoveCard(card: KanbanCard, targetColumnId: UUID) {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before moving cards.');
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			const updatedBoard = await moveKanbanCard(projectId, card.id, targetColumnId, 0);
			board.set(updatedBoard);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to move card.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleDeleteCard(card: KanbanCard) {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before deleting cards.');
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			const updatedBoard = await deleteKanbanCard(projectId, card.id);
			board.set(updatedBoard);
			toastInfo('Card deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete card.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleDeleteColumn(columnId: UUID) {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before deleting columns.');
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			const updatedBoard = await deleteKanbanColumn(projectId, columnId);
			board.set(updatedBoard);
			toastInfo('Column deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete column.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleAddMilestone() {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before adding milestones.');
			return;
		}

		const form = get(milestoneForm);
		if (!form.title.trim()) {
			const message = 'Milestone title is required.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			await addMilestone(projectId, {
				title: form.title,
				description: form.description,
				dueDate: new Date(form.dueDate).toISOString()
			});
			milestoneForm.set(defaultMilestoneForm());
			await milestonesQuery.refetch();
			toastSuccess('Milestone added.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to add milestone.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleToggleMilestone(milestone: Milestone) {
		if (!get(isAuthenticated)) {
			toastError('Sign in to update milestones.');
			return;
		}

		loading.set(true);
		try {
			await bulkUpdateMilestones(milestone.project_id, [
				{ milestoneId: milestone.id, completed: !milestone.completed }
			]);
			await milestonesQuery.refetch();
			toastInfo('Milestone updated.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update milestone.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleCreateDependency() {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before creating dependencies.');
			return;
		}

		const form = get(dependencyForm);
		if (!form.dependsOnProjectId) {
			const message = 'Select a project to depend on.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const projectId = get(activeProject)!.id;
			await createDependency(projectId, {
				dependsOnProjectId: form.dependsOnProjectId,
				type: form.type
			});
			dependencyForm.set(defaultDependencyForm());
			await dependenciesQuery.refetch();
			toastSuccess('Dependency created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create dependency.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	async function handleDeleteDependency(dependencyId: UUID) {
		if (!get(isAuthenticated) || !get(activeProject)) {
			toastError('Select a project before deleting dependencies.');
			return;
		}

		try {
			const projectId = get(activeProject)!.id;
			await deleteDependency(projectId, dependencyId);
			await dependenciesQuery.refetch();
			toastInfo('Dependency removed.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete dependency.');
			error.set(message);
			toastError(message);
		}
	}

	async function handleDuplicateProject(templateId: UUID) {
		if (!get(isAuthenticated)) {
			toastError('Sign in to duplicate projects.');
			return;
		}

		const form = get(duplicateForm);
		if (!templateId || !form.name.trim()) {
			const message = 'Provide a name for the duplicate project.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		try {
			const duplicate = await duplicateProject(templateId, {
				name: form.name,
				description: form.description,
				status: form.status,
				copyBoard: form.copyBoard,
				copyMilestones: form.copyMilestones,
				copyDependencies: form.copyDependencies
			});
			duplicateForm.set(defaultDuplicateForm());
			await loadProjects();
			await selectProject(duplicate);
			toastSuccess('Project duplicated.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to duplicate project.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	}

	const updateActiveProjectField = <K extends keyof Project>(field: K, value: Project[K]) => {
		activeProject.update((project) => (project ? { ...project, [field]: value } : project));
	};

	async function saveProjectStatus() {
		const project = get(activeProject);
		if (!project) return;

		try {
			await updateProjectStatus(project.id, project.status);
			await loadProjects();
			toastSuccess('Status saved.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to save status.');
			error.set(message);
			toastError(message);
		}
	}

	async function saveProjectMetrics() {
		const project = get(activeProject);
		if (!project) return;

		try {
			await updateProjectMetrics(project.id, {
				healthScore: project.health_score,
				mrr: project.mrr,
				cac: project.cac,
				ltv: project.ltv,
				churnRate: project.churn_rate
			});
			await loadProjects();
			toastSuccess('Metrics saved.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to save metrics.');
			error.set(message);
			toastError(message);
		}
	}

	function getOtherProjects(projectId?: UUID | null) {
		return get(projects).filter((project) => project.id !== projectId);
	}

	return {
		isAuthenticated,
		loading,
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
	};
}

