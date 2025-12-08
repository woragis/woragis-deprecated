import { onDestroy, onMount } from 'svelte';
import { derived, get, writable, type Readable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type {
	ArchitectureDiagramType,
	DocumentationSection,
	DocumentationSectionType,
	DocumentationVisibility,
	KanbanBoard,
	KanbanCard,
	Milestone,
	Project,
	ProjectArchitectureDiagram,
	ProjectDependency,
	ProjectDocumentation,
	ProjectFileStructure,
	ProjectStatus,
	ProjectTechnology,
	TechnologyCategory,
	UUID
} from '$lib/api/types';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';
import {
	useAddMilestoneMutation,
	useBulkUpdateMilestonesMutation,
	useCreateArchitectureDiagramMutation,
	useCreateDependencyMutation,
	useCreateDocumentationMutation,
	useCreateDocumentationSectionMutation,
	useCreateFileStructureMutation,
	useCreateKanbanCardMutation,
	useCreateKanbanColumnMutation,
	useCreateTechnologyMutation,
	useDeleteArchitectureDiagramMutation,
	useDeleteDependencyMutation,
	useDeleteDocumentationSectionMutation,
	useDeleteFileStructureMutation,
	useDeleteKanbanCardMutation,
	useDeleteKanbanColumnMutation,
	useDeleteTechnologyMutation,
	useDuplicateProjectMutation,
	useMoveKanbanCardMutation,
	useProjectArchitectureDiagramsQuery,
	useProjectBoardQuery,
	useProjectBySlugQuery,
	useProjectDependenciesQuery,
	useProjectDocumentationQuery,
	useProjectDocumentationSectionsQuery,
	useProjectFileStructuresQuery,
	useProjectMilestonesQuery,
	useProjectTechnologiesQuery,
	useProjectsListQuery,
	useUpdateDocumentationVisibilityMutation,
	useUpdateProjectMetricsMutation,
	useUpdateProjectStatusMutation
} from '@hooks/projects';
import { statusOptions } from '../status';

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

type FormUpdater<T> = {
	<K extends keyof T>(field: K, value: T[K]): void;
};

interface AuthState {
	isAuthenticated: boolean;
	userId: UUID | null;
}

export function createProjectDetailLogic(slugStore: Readable<string | null>) {
	const queryClient = useQueryClient();

	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);
	const error = writable<string | null>(null);

	const projects = writable<Project[]>([]);
	const activeProject = writable<Project | null>(null);
	const board = writable<KanbanBoard | null>(null);
	const milestones = writable<Milestone[]>([]);
	const dependencies = writable<ProjectDependency[]>([]);
	const documentation = writable<ProjectDocumentation | null>(null);
	const documentationSections = writable<DocumentationSection[]>([]);
	const technologies = writable<ProjectTechnology[]>([]);
	const fileStructures = writable<ProjectFileStructure[]>([]);
	const architectureDiagrams = writable<ProjectArchitectureDiagram[]>([]);

	const columnForm = writable<ColumnFormState>(defaultColumnForm());
	const cardForm = writable<CardFormState>(defaultCardForm());
	const milestoneForm = writable<MilestoneFormState>(defaultMilestoneForm());
	const dependencyForm = writable<DependencyFormState>(defaultDependencyForm());
	const duplicateForm = writable<DuplicateFormState>(defaultDuplicateForm());

	const activeProjectId = derived(activeProject, ($project) => $project?.id ?? null);
	const projectsQueryOptions = derived(authStateStore, ($auth) => ({ enabled: $auth.isAuthenticated }));
	const detailQueryOptions = derived(authStateStore, ($auth) => ({ enabled: $auth.isAuthenticated }));

	const projectsQuery = useProjectsListQuery(projectsQueryOptions);
	const projectQuery = useProjectBySlugQuery(slugStore, detailQueryOptions);
	const boardQuery = useProjectBoardQuery(activeProjectId, detailQueryOptions);
	const milestonesQuery = useProjectMilestonesQuery(activeProjectId, detailQueryOptions);
	const dependenciesQuery = useProjectDependenciesQuery(activeProjectId, detailQueryOptions);
	const documentationQuery = useProjectDocumentationQuery(activeProjectId, detailQueryOptions);
	const documentationSectionsQuery = useProjectDocumentationSectionsQuery(
		activeProjectId,
		detailQueryOptions
	);
	const technologiesQuery = useProjectTechnologiesQuery(activeProjectId, detailQueryOptions);
	const fileStructuresQuery = useProjectFileStructuresQuery(activeProjectId, detailQueryOptions);
	const architectureDiagramsQuery = useProjectArchitectureDiagramsQuery(
		activeProjectId,
		detailQueryOptions
	);

	const invalidateProjectsList = () => queryClient.invalidateQueries({ queryKey: ['projects', 'list'] });
	const invalidateProjectBySlug = () => {
		const slug = get(slugStore);
		if (!slug) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'slug', slug] });
	};
	const invalidateBoard = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'board', targetId] });
	};
	const invalidateMilestones = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'milestones', targetId] });
	};
	const invalidateDependencies = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'dependencies', targetId] });
	};
	const invalidateDocumentation = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'documentation', targetId] });
	};
	const invalidateDocumentationSections = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({
			queryKey: ['projects', 'documentation-sections', targetId]
		});
	};
	const invalidateTechnologies = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'technologies', targetId] });
	};
	const invalidateFileStructures = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['projects', 'file-structures', targetId] });
	};
	const invalidateArchitectureDiagrams = (projectId?: UUID | null) => {
		const targetId = projectId ?? get(activeProjectId);
		if (!targetId) return Promise.resolve();
		return queryClient.invalidateQueries({
			queryKey: ['projects', 'architecture-diagrams', targetId]
		});
	};

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
	const createDocumentationMutation = useCreateDocumentationMutation();
	const updateDocumentationVisibilityMutation = useUpdateDocumentationVisibilityMutation();
	const createDocumentationSectionMutation = useCreateDocumentationSectionMutation();
	const deleteDocumentationSectionMutation = useDeleteDocumentationSectionMutation();
	const createTechnologyMutation = useCreateTechnologyMutation();
	const deleteTechnologyMutation = useDeleteTechnologyMutation();
	const createFileStructureMutation = useCreateFileStructureMutation();
	const deleteFileStructureMutation = useDeleteFileStructureMutation();
	const createArchitectureDiagramMutation = useCreateArchitectureDiagramMutation();
	const deleteArchitectureDiagramMutation = useDeleteArchitectureDiagramMutation();

	const resetStateForUnauthenticated = () => {
		projects.set([]);
		activeProject.set(null);
		board.set(null);
		milestones.set([]);
		dependencies.set([]);
		documentation.set(null);
		documentationSections.set([]);
		technologies.set([]);
		fileStructures.set([]);
		architectureDiagrams.set([]);
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
		const currentProject = get(activeProject);
		if (!auth.isAuthenticated || !currentProject) return;
		await Promise.all([
			invalidateBoard(currentProject.id),
			invalidateMilestones(currentProject.id),
			invalidateDependencies(currentProject.id),
			invalidateDocumentation(currentProject.id),
			invalidateDocumentationSections(currentProject.id),
			invalidateTechnologies(currentProject.id),
			invalidateFileStructures(currentProject.id),
			invalidateArchitectureDiagrams(currentProject.id)
		]);
	};

	const projectsQueryUnsubscribe = projectsQuery.subscribe((result) => {
		if (result.data) {
			const data = result.data;
			projects.set(data);
			error.set(null);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load projects.');
		}
	});

	const projectQueryUnsubscribe = projectQuery.subscribe((result) => {
		if (result.data) {
			activeProject.set(result.data);
			void refetchProjectDetails();
			error.set(null);
		} else if (result.error) {
			activeProject.set(null);
			error.set(result.error.message ?? 'Unable to load project.');
		}

		if (!result.data && !result.isLoading) {
			board.set(null);
			milestones.set([]);
			dependencies.set([]);
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

	const documentationQueryUnsubscribe = documentationQuery.subscribe((result) => {
		if (result.data) {
			documentation.set(result.data);
		}
		if (result.error && result.error.message?.includes('not found') === false) {
			// Documentation might not exist yet, so only show error if it's not a 404
			error.set(result.error.message ?? 'Unable to load documentation.');
		}
	});

	const documentationSectionsQueryUnsubscribe = documentationSectionsQuery.subscribe((result) => {
		if (result.data) {
			documentationSections.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load documentation sections.');
		}
	});

	const technologiesQueryUnsubscribe = technologiesQuery.subscribe((result) => {
		if (result.data) {
			technologies.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load technologies.');
		}
	});

	const fileStructuresQueryUnsubscribe = fileStructuresQuery.subscribe((result) => {
		if (result.data) {
			fileStructures.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load file structures.');
		}
	});

	const architectureDiagramsQueryUnsubscribe = architectureDiagramsQuery.subscribe((result) => {
		if (result.data) {
			architectureDiagrams.set(result.data);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load architecture diagrams.');
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
			await Promise.all([invalidateProjectsList(), invalidateProjectBySlug()]);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to load projects.');
			error.set(message);
			toastError(message);
		}
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
		projectQueryUnsubscribe();
		boardQueryUnsubscribe();
		milestonesQueryUnsubscribe();
		dependenciesQueryUnsubscribe();
		documentationQueryUnsubscribe();
		documentationSectionsQueryUnsubscribe();
		technologiesQueryUnsubscribe();
		fileStructuresQueryUnsubscribe();
		architectureDiagramsQueryUnsubscribe();
	});

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
				updates: {
					updates: [{ id: milestone.id, completed: !milestone.completed }]
				}
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

	const handleDuplicateProject = async (templateId: UUID): Promise<Project | null> => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			toastError('Sign in to duplicate projects.');
			return null;
		}

		const form = get(duplicateForm);
		if (!templateId || !form.name.trim()) {
			const message = 'Provide a name for the duplicate project.';
			error.set(message);
			toastError(message);
			return null;
		}

		try {
			const duplicate = await get(duplicateProjectMutation).mutateAsync({
				templateProjectId: templateId,
				payload: {
					name: form.name,
					description: form.description,
					status: form.status,
					copyKanban: form.copyBoard,
					copyMilestones: form.copyMilestones,
					copyDependencies: form.copyDependencies
				}
			});
			duplicateForm.set(defaultDuplicateForm());
			await invalidateProjectsList();
			toastSuccess('Project duplicated.');
			return duplicate;
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to duplicate project.');
			error.set(message);
			toastError(message);
			return null;
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

	// Documentation handlers

	const handleCreateDocumentation = async () => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before creating documentation.');
			return;
		}

		try {
			await get(createDocumentationMutation).mutateAsync({
				projectId: project.id,
				payload: { visibility: 'collaborators' }
			});
			await invalidateDocumentation();
			toastSuccess('Documentation created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create documentation.');
			error.set(message);
			toastError(message);
		}
	};

	const handleUpdateDocumentationVisibility = async (visibility: DocumentationVisibility) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before updating documentation visibility.');
			return;
		}

		try {
			await get(updateDocumentationVisibilityMutation).mutateAsync({
				projectId: project.id,
				visibility
			});
			await invalidateDocumentation();
			toastSuccess('Visibility updated.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update visibility.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateDocumentationSection = async (payload: {
		type: DocumentationSectionType;
		title: string;
		content: string;
		position?: number;
	}) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding sections.');
			return;
		}

		try {
			await get(createDocumentationSectionMutation).mutateAsync({
				projectId: project.id,
				payload
			});
			await invalidateDocumentationSections();
			toastSuccess('Section created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create section.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteDocumentationSection = async (sectionId: UUID) => {
		try {
			await get(deleteDocumentationSectionMutation).mutateAsync({ sectionId });
			await invalidateDocumentationSections();
			toastInfo('Section deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete section.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateTechnology = async (payload: {
		name: string;
		version: string;
		category: TechnologyCategory;
		purpose: string;
		link?: string;
	}) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding technologies.');
			return;
		}

		try {
			await get(createTechnologyMutation).mutateAsync({
				projectId: project.id,
				payload
			});
			await invalidateTechnologies();
			toastSuccess('Technology added.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to add technology.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteTechnology = async (techId: UUID) => {
		try {
			await get(deleteTechnologyMutation).mutateAsync({ techId });
			await invalidateTechnologies();
			toastInfo('Technology deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete technology.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateFileStructure = async (payload: {
		path: string;
		name: string;
		is_directory: boolean;
		parent_id?: UUID;
		language?: string;
		line_count?: number;
		purpose?: string;
		position?: number;
	}) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding file structures.');
			return;
		}

		try {
			await get(createFileStructureMutation).mutateAsync({
				projectId: project.id,
				payload: {
					...payload,
					isDirectory: payload.is_directory
				}
			});
			await invalidateFileStructures();
			toastSuccess('File structure added.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to add file structure.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteFileStructure = async (fileStructureId: UUID) => {
		try {
			await get(deleteFileStructureMutation).mutateAsync({ fileStructureId });
			await invalidateFileStructures();
			toastInfo('File structure deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete file structure.');
			error.set(message);
			toastError(message);
		}
	};

	const handleCreateArchitectureDiagram = async (payload: {
		type: ArchitectureDiagramType;
		title: string;
		description: string;
		content: string;
		format?: string;
		image_url?: string;
	}) => {
		const project = get(activeProject);
		if (!project) {
			toastError('Select a project before adding diagrams.');
			return;
		}

		try {
			await get(createArchitectureDiagramMutation).mutateAsync({
				projectId: project.id,
				payload: {
					...payload,
					format: payload.format || 'mermaid'
				}
			});
			await invalidateArchitectureDiagrams();
			toastSuccess('Diagram created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create diagram.');
			error.set(message);
			toastError(message);
		}
	};

	const handleDeleteArchitectureDiagram = async (diagramId: UUID) => {
		try {
			await get(deleteArchitectureDiagramMutation).mutateAsync({ diagramId });
			await invalidateArchitectureDiagrams();
			toastInfo('Diagram deleted.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete diagram.');
			error.set(message);
			toastError(message);
		}
	};

	return {
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
	};
}

