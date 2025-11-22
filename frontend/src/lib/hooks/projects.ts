import { derived, readable, type Readable } from 'svelte/store';
import { createMutation, createQuery } from '@tanstack/svelte-query';

import { projectsApi } from '$lib/api/projects';
import type { CreateProjectInput } from '$lib/api/projects';
import type {
	ArchitectureDiagramType,
	DocumentationSection,
	DocumentationSectionType,
	DocumentationVisibility,
	KanbanBoard,
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

type MaybeReadable<T> = T | Readable<T>;

const isReadable = (value: unknown): value is Readable<unknown> =>
	typeof value === 'object' && value !== null && 'subscribe' in value;

const toReadable = <T>(value: MaybeReadable<T>): Readable<T> =>
	(isReadable(value) ? (value as Readable<T>) : readable(value)) as Readable<T>;

const toOptionsReadable = <T extends object>(value: MaybeReadable<T> | undefined, fallback: T): Readable<T> =>
	value === undefined ? readable(fallback) : toReadable(value);

export interface ProjectsListOptions {
	enabled?: boolean;
}

export const useProjectsListQuery = (options?: MaybeReadable<ProjectsListOptions>) => {
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<Project[]>(
		derived(optionsStore, ($options) => ({
			queryKey: ['projects', 'list'],
			queryFn: () => projectsApi.listProjects(),
			enabled: $options.enabled ?? true,
			placeholderData: () => []
		}))
	);
};

export const useProjectBySlugQuery = (
	slug: MaybeReadable<string | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const slugStore = toReadable(slug);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<Project>(
		derived([slugStore, optionsStore], ([$slug, $options]) => ({
			queryKey: ['projects', 'slug', $slug],
			queryFn: () => {
				if (!$slug) {
					throw new Error('Project slug is required to load project details');
				}
				return projectsApi.getProjectBySlug($slug);
			},
			enabled: Boolean($slug) && ($options.enabled ?? true),
			retry: false
		}))
	);
};

export const useProjectBoardQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<KanbanBoard>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'board', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load board');
				}
				return projectsApi.getKanbanBoard($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true)
		}))
	);
};

export const useProjectMilestonesQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<Milestone[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'milestones', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load milestones');
				}
				return projectsApi.listMilestones($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useProjectDependenciesQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<ProjectDependency[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'dependencies', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load dependencies');
				}
				return projectsApi.listDependencies($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useCreateProjectMutation = () =>
	createMutation({
		mutationFn: (input: CreateProjectInput) => projectsApi.createProject(input)
	});

export const useUpdateProjectStatusMutation = () =>
	createMutation({
		mutationFn: ({ projectId, status }: { projectId: UUID; status: ProjectStatus }) =>
			projectsApi.updateProjectStatus(projectId, status)
	});

export const useUpdateProjectMetricsMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.updateProjectMetrics>[1];
		}) => projectsApi.updateProjectMetrics(projectId, payload)
	});

export const useAddMilestoneMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.addMilestone>[1];
		}) => projectsApi.addMilestone(projectId, payload)
	});

export const useBulkUpdateMilestonesMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			updates
		}: {
			projectId: UUID;
			updates: Parameters<typeof projectsApi.bulkUpdateMilestones>[1];
		}) => projectsApi.bulkUpdateMilestones(projectId, updates)
	});

export const useCreateKanbanColumnMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createKanbanColumn>[1];
		}) => projectsApi.createKanbanColumn(projectId, payload)
	});

export const useUpdateKanbanColumnMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			columnId,
			payload
		}: {
			projectId: UUID;
			columnId: UUID;
			payload: Parameters<typeof projectsApi.updateKanbanColumn>[2];
		}) => projectsApi.updateKanbanColumn(projectId, columnId, payload)
	});

export const useReorderKanbanColumnsMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			columnOrder
		}: {
			projectId: UUID;
			columnOrder: UUID[];
		}) => projectsApi.reorderKanbanColumns(projectId, columnOrder)
	});

export const useDeleteKanbanColumnMutation = () =>
	createMutation({
		mutationFn: ({ projectId, columnId }: { projectId: UUID; columnId: UUID }) =>
			projectsApi.deleteKanbanColumn(projectId, columnId)
	});

export const useCreateKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createKanbanCard>[1];
		}) => projectsApi.createKanbanCard(projectId, payload)
	});

export const useUpdateKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			cardId,
			payload
		}: {
			projectId: UUID;
			cardId: UUID;
			payload: Parameters<typeof projectsApi.updateKanbanCard>[2];
		}) => projectsApi.updateKanbanCard(projectId, cardId, payload)
	});

export const useMoveKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			cardId,
			targetColumnId,
			targetPosition
		}: {
			projectId: UUID;
			cardId: UUID;
			targetColumnId: UUID;
			targetPosition: number;
		}) => projectsApi.moveKanbanCard(projectId, cardId, targetColumnId, targetPosition)
	});

export const useDeleteKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({ projectId, cardId }: { projectId: UUID; cardId: UUID }) =>
			projectsApi.deleteKanbanCard(projectId, cardId)
	});

export const useCreateDependencyMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createDependency>[1];
		}) => projectsApi.createDependency(projectId, payload)
	});

export const useDeleteDependencyMutation = () =>
	createMutation({
		mutationFn: ({ projectId, dependencyId }: { projectId: UUID; dependencyId: UUID }) =>
			projectsApi.deleteDependency(projectId, dependencyId)
	});

export const useDuplicateProjectMutation = () =>
	createMutation({
		mutationFn: ({
			templateProjectId,
			payload
		}: {
			templateProjectId: UUID;
			payload: Parameters<typeof projectsApi.duplicateProject>[1];
		}) => projectsApi.duplicateProject(templateProjectId, payload)
	});

// Documentation hooks

export const useProjectDocumentationQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<ProjectDocumentation>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'documentation', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load documentation');
				}
				return projectsApi.getDocumentation($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			retry: false
		}))
	);
};

export const useProjectDocumentationSectionsQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<DocumentationSection[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'documentation-sections', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load documentation sections');
				}
				return projectsApi.listDocumentationSections($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useProjectTechnologiesQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<ProjectTechnology[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'technologies', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load technologies');
				}
				return projectsApi.listTechnologies($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useProjectFileStructuresQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<ProjectFileStructure[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'file-structures', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load file structures');
				}
				return projectsApi.listFileStructures($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useProjectArchitectureDiagramsQuery = (
	projectId: MaybeReadable<UUID | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const projectIdStore = toReadable(projectId);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<ProjectArchitectureDiagram[]>(
		derived([projectIdStore, optionsStore], ([$projectId, $options]) => ({
			queryKey: ['projects', 'architecture-diagrams', $projectId],
			queryFn: () => {
				if (!$projectId) {
					throw new Error('Project id is required to load architecture diagrams');
				}
				return projectsApi.listArchitectureDiagrams($projectId);
			},
			enabled: Boolean($projectId) && ($options.enabled ?? true),
			placeholderData: () => []
		}))
	);
};

export const useCreateDocumentationMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createDocumentation>[1];
		}) => projectsApi.createDocumentation(projectId, payload)
	});

export const useUpdateDocumentationVisibilityMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			visibility
		}: {
			projectId: UUID;
			visibility: DocumentationVisibility;
		}) => projectsApi.updateDocumentationVisibility(projectId, visibility)
	});

export const useCreateDocumentationSectionMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createDocumentationSection>[1];
		}) => projectsApi.createDocumentationSection(projectId, payload)
	});

export const useUpdateDocumentationSectionMutation = () =>
	createMutation({
		mutationFn: ({
			sectionId,
			payload
		}: {
			sectionId: UUID;
			payload: Parameters<typeof projectsApi.updateDocumentationSection>[1];
		}) => projectsApi.updateDocumentationSection(sectionId, payload)
	});

export const useDeleteDocumentationSectionMutation = () =>
	createMutation({
		mutationFn: ({ sectionId }: { sectionId: UUID }) =>
			projectsApi.deleteDocumentationSection(sectionId)
	});

export const useReorderDocumentationSectionsMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			sectionOrder
		}: {
			projectId: UUID;
			sectionOrder: UUID[];
		}) => projectsApi.reorderDocumentationSections(projectId, sectionOrder)
	});

export const useCreateTechnologyMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createTechnology>[1];
		}) => projectsApi.createTechnology(projectId, payload)
	});

export const useUpdateTechnologyMutation = () =>
	createMutation({
		mutationFn: ({
			techId,
			payload
		}: {
			techId: UUID;
			payload: Parameters<typeof projectsApi.updateTechnology>[1];
		}) => projectsApi.updateTechnology(techId, payload)
	});

export const useDeleteTechnologyMutation = () =>
	createMutation({
		mutationFn: ({ techId }: { techId: UUID }) => projectsApi.deleteTechnology(techId)
	});

export const useBulkCreateTechnologiesMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			technologies
		}: {
			projectId: UUID;
			technologies: Parameters<typeof projectsApi.bulkCreateTechnologies>[1];
		}) => projectsApi.bulkCreateTechnologies(projectId, technologies)
	});

export const useBulkUpdateTechnologiesMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			technologies
		}: {
			projectId: UUID;
			technologies: Parameters<typeof projectsApi.bulkUpdateTechnologies>[1];
		}) => projectsApi.bulkUpdateTechnologies(projectId, technologies)
	});

export const useCreateFileStructureMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createFileStructure>[1];
		}) => projectsApi.createFileStructure(projectId, payload)
	});

export const useUpdateFileStructureMutation = () =>
	createMutation({
		mutationFn: ({
			fileStructureId,
			payload
		}: {
			fileStructureId: UUID;
			payload: Parameters<typeof projectsApi.updateFileStructure>[1];
		}) => projectsApi.updateFileStructure(fileStructureId, payload)
	});

export const useDeleteFileStructureMutation = () =>
	createMutation({
		mutationFn: ({ fileStructureId }: { fileStructureId: UUID }) =>
			projectsApi.deleteFileStructure(fileStructureId)
	});

export const useBulkCreateFileStructuresMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			structures
		}: {
			projectId: UUID;
			structures: Parameters<typeof projectsApi.bulkCreateFileStructures>[1];
		}) => projectsApi.bulkCreateFileStructures(projectId, structures)
	});

export const useBulkUpdateFileStructuresMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			structures
		}: {
			projectId: UUID;
			structures: Parameters<typeof projectsApi.bulkUpdateFileStructures>[1];
		}) => projectsApi.bulkUpdateFileStructures(projectId, structures)
	});

export const useCreateArchitectureDiagramMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof projectsApi.createArchitectureDiagram>[1];
		}) => projectsApi.createArchitectureDiagram(projectId, payload)
	});

export const useUpdateArchitectureDiagramMutation = () =>
	createMutation({
		mutationFn: ({
			diagramId,
			payload
		}: {
			diagramId: UUID;
			payload: Parameters<typeof projectsApi.updateArchitectureDiagram>[1];
		}) => projectsApi.updateArchitectureDiagram(diagramId, payload)
	});

export const useDeleteArchitectureDiagramMutation = () =>
	createMutation({
		mutationFn: ({ diagramId }: { diagramId: UUID }) =>
			projectsApi.deleteArchitectureDiagram(diagramId)
	});

