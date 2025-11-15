import { onDestroy, onMount } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

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
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';
import {
	useAddMilestoneMutation,
	useBulkUpdateMilestonesMutation,
	useCreateDependencyMutation,
	useCreateKanbanCardMutation,
	useCreateKanbanColumnMutation,
	useCreateProjectMutation,
	useDeleteDependencyMutation,
	useDeleteKanbanCardMutation,
	useDeleteKanbanColumnMutation,
	useDuplicateProjectMutation,
	useMoveKanbanCardMutation,
	useProjectBoardQuery,
	useProjectDependenciesQuery,
	useProjectMilestonesQuery,
	useProjectsListQuery,
	useUpdateProjectMetricsMutation,
	useUpdateProjectStatusMutation
} from '@hooks/projects';

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

interface AuthState {
	isAuthenticated: boolean;
	userId: UUID | null;
}

export function createProjectsLogic() {
	const queryClient = useQueryClient();

	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);
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

	const activeProjectId = derived(activeProject, ($project) => $project?.id ?? null);

	const projectsQueryOptions = derived(authStateStore, ($auth) => ({ enabled: $auth.isAuthenticated }));
	const detailQueryOptions = derived(authStateStore, ($auth) => ({ enabled: $auth.isAuthenticated }));

	const projectsQuery = useProjectsListQuery(projectsQueryOptions);
	const boardQuery = useProjectBoardQuery(activeProjectId, detailQueryOptions);
	const milestonesQuery = useProjectMilestonesQuery(activeProjectId, detailQueryOptions);
	const dependenciesQuery = useProjectDependenciesQuery(activeProjectId, detailQueryOptions);

	const invalidateProjectsList = () => queryClient.invalidateQueries({ queryKey: ['projects', 'list'] });
	const invalidateBoard = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProject)?.id ?? null;
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'board', targetId] });
	};
	const invalidateMilestones = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProject)?.id ?? null;
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'milestones', targetId] });
	};
	const invalidateDependencies = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProject)?.id ?? null;
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'dependencies', targetId] });
	};

	const createProjectMutation = useCreateProjectMutation();
	const createColumnMutation = useCreateKanbanColumnMutation();
	const createCardMutation = useCreateKanbanCardMutation();
	const moveCardMutation = useMoveKanbanCardMutation();
	const deleteCardMutation = useDeleteKanbanCardMutation();
	const deleteColumnMutation = useDeleteKanbanColumnMutation();
	const addMilestoneMutation = useAddMilestoneMutation();
	const bulkUpdateMilestonesMutation = useBulkUpdateMilestonesMutation();
	const createDependencyMutation = useCreateDependencyMutation();
	const deleteDependencyMutation = useDeleteDependencyMutation();
	const duplicateProjectMutation = useDuplicateProjectMutation();
	const updateStatusMutation = useUpdateProjectStatusMutation();
	const updateMetricsMutation = useUpdateProjectMetricsMutation();

	const resetStateForUnauthenticated = () => {
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
		error.set(null);
		queryClient.removeQueries({ queryKey: ['projects'] });
	};

	const refetchProjectDetails = async () => {
		const auth = get(authStateStore);
		const current = get(activeProject);
		if (!auth.isAuthenticated || !current) return;
		await Promise.all([invalidateBoard(current.id), invalidateMilestones(current.id), invalidateDependencies(current.id)]);
	};

	const projectsQueryUnsubscribe = projectsQuery.subscribe((result) => {
		if (result.data) {
			const currentActive = get(activeProject);
			const data = result.data;
			const nextActive =
				currentActive && data.length
					? data.find((project) => project.id === currentActive.id) ?? data[0]
					: data[0] ?? null;
			projects.set(data);
			activeProject.set(nextActive ?? null);
			if (nextActive) {
				void refetchProjectDetails();
			} else {
				board.set(null);
				milestones.set([]);
				dependencies.set([]);
			}
			error.set(null);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load projects.');
		}
	});

	const boardQueryUnsubscribe = boardQuery.subscribe((result) => {
		if (result.data) {
			board.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load kanban board.');
		}
	});

	const milestonesQueryUnsubscribe = milestonesQuery.subscribe((result) => {
		if (result.data) {
			milestones.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load milestones.');
		}
	});

	const dependenciesQueryUnsubscribe = dependenciesQuery.subscribe((result) => {
		if (result.data) {
			dependencies.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load dependencies.');
		}
	});

	const loadProjects = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			const message = 'You must be signed in to load projects.';
			error.set(message);
			toastError(message);
			return;
		}

		try {
			await invalidateProjectsList();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to load projects.');
			error.set(message);
			toastError(message);
		}
	};

	const selectProject = async (project: Project) => {
		activeProject.set(project);
		await refetchProjectDetails();
	};

	onMount(() => {
		const unsubscribeAuth = authStore.subscribe((state) => {
			const next: AuthState = {
				isAuthenticated: state.isAuthenticated,
				userId: state.user?.id ?? null
			};
			authStateStore.set(next);

			if (!next.isAuthenticated) {
				resetStateForUnauthenticated();
				return;
			}

			void loadProjects();
		});

		return () => {
			unsubscribeAuth();
		};
	});

	onDestroy(() => {
		projectsQueryUnsubscribe();
		boardQueryUnsubscribe();
		milestonesQueryUnsubscribe();
		dependenciesQueryUnsubscribe();
	});

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

	const handleCreateProject = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
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

		try {
			const newProject = await get(createProjectMutation).mutateAsync({
				name: form.name.trim(),
				description: form.description,
				status: form.status,
				healthScore: form.healthScore
			});
			projectForm.set(defaultProjectForm());
			await invalidateProjectsList();
			await selectProject(newProject);
			toastSuccess('Project created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create project.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateColumn = async () => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding columns.');
			return;
		}

		const form = get(columnForm);
		if (!form.name.trim()) {
			const message = 'Column name is required.';
			error.set(message);
			toastError(message);
			return;
		}

		try {
			await get(createColumnMutation).mutateAsync({
				projectId: project.id,
				payload: {
					name: form.name,
					position: form.position,
					wipLimit: form.wipLimit
				}
			});
			columnForm.set(defaultColumnForm());
			await invalidateBoard();
			toastSuccess('Column created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create column.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateCard = async () => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding cards.');
			return;
		}

		const form = get(cardForm);
		if (!form.columnId || !form.title.trim()) {
			const message = 'Column and title are required.';
			error.set(message);
			toastError(message);
			return;
		}

		try {
			await get(createCardMutation).mutateAsync({
				projectId: project.id,
				payload: {
					columnId: form.columnId,
					title: form.title,
					description: form.description,
					dueDate: form.dueDate ? new Date(form.dueDate).toISOString() : undefined,
					milestoneId: form.milestoneId || undefined,
					position: form.position
				}
			});
			cardForm.set(defaultCardForm());
			await invalidateBoard();
			toastSuccess('Card created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create card.');
			error.set(message);
			toastError(message);
		}
	};

	const handleMoveCard = async (card: KanbanCard, targetColumnId: UUID) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before moving cards.');
			return;
		}

		try {
			await get(moveCardMutation).mutateAsync({
				projectId: project.id,
				cardId: card.id,
				targetColumnId,
				targetPosition: 0
			});
			await invalidateBoard();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to move card.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteCard = async (card: KanbanCard) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before deleting cards.');
			return;
		}

		try {
			await get(deleteCardMutation).mutateAsync({ projectId: project.id, cardId: card.id });
			await invalidateBoard();
			toastInfo('Card deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete card.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteColumn = async (columnId: UUID) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before deleting columns.');
			return;
		}

		try {
			await get(deleteColumnMutation).mutateAsync({ projectId: project.id, columnId });
			await invalidateBoard();
			toastInfo('Column deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete column.');
			error.set(message);
			toastError(message);
		}
	};

	const handleAddMilestone = async () => {
		const project = get(activeProject);
		if (!project) {
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

		try {
			await get(addMilestoneMutation).mutateAsync({
				projectId: project.id,
				payload: {
					title: form.title,
					description: form.description,
					dueDate: new Date(form.dueDate).toISOString()
				}
			});
			milestoneForm.set(defaultMilestoneForm());
			await invalidateMilestones();
			toastSuccess('Milestone added.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to add milestone.');
			error.set(message);
			toastError(message);
		}
	};

	const handleToggleMilestone = async (milestone: Milestone) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Sign in to update milestones.');
			return;
		}

		try {
			await get(bulkUpdateMilestonesMutation).mutateAsync({
				projectId: project.id,
				updates: [{ milestoneId: milestone.id, completed: !milestone.completed }]
			});
			await invalidateMilestones();
			toastInfo('Milestone updated.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update milestone.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateDependency = async () => {
		const project = get(activeProject);
		if (!project) {
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

		try {
			await get(createDependencyMutation).mutateAsync({
				projectId: project.id,
				payload: {
					dependsOnProjectId: form.dependsOnProjectId,
					type: form.type
				}
			});
			dependencyForm.set(defaultDependencyForm());
			await invalidateDependencies();
			toastSuccess('Dependency created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create dependency.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteDependency = async (dependencyId: UUID) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before deleting dependencies.');
			return;
		}

		try {
			await get(deleteDependencyMutation).mutateAsync({ projectId: project.id, dependencyId });
			await invalidateDependencies();
			toastInfo('Dependency removed.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete dependency.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDuplicateProject = async (templateId: UUID) => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
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

		try {
			const duplicate = await get(duplicateProjectMutation).mutateAsync({
				templateProjectId: templateId,
				payload: {
					name: form.name,
					description: form.description,
					status: form.status,
					copyBoard: form.copyBoard,
					copyMilestones: form.copyMilestones,
					copyDependencies: form.copyDependencies
				}
			});
			duplicateForm.set(defaultDuplicateForm());
			await invalidateProjectsList();
			await selectProject(duplicate);
			toastSuccess('Project duplicated.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to duplicate project.');
			error.set(message);
			toastError(message);
		}
	};

	const updateActiveProjectField = <K extends keyof Project>(field: K, value: Project[K]) => {
		activeProject.update((project) => (project ? { ...project, [field]: value } : project));
	};

	const saveProjectStatus = async () => {
		const project = get(activeProject);
		if (!project) return;

		try {
			await get(updateStatusMutation).mutateAsync({ projectId: project.id, status: project.status });
			await invalidateProjectsList();
			toastSuccess('Status saved.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to save status.');
			error.set(message);
			toastError(message);
		}
	};

	const saveProjectMetrics = async () => {
		const project = get(activeProject);
		if (!project) return;

		try {
			await get(updateMetricsMutation).mutateAsync({
				projectId: project.id,
				payload: {
					healthScore: project.health_score,
					mrr: project.mrr,
					cac: project.cac,
					ltv: project.ltv,
					churnRate: project.churn_rate
				}
			});
			await invalidateProjectsList();
			toastSuccess('Metrics saved.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to save metrics.');
			error.set(message);
			toastError(message);
		}
	};

	const getOtherProjects = (projectId?: UUID | null) => get(projects).filter((project) => project.id !== projectId);

	return {
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
	};
}

