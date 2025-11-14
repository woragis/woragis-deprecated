import { createMutation, createQuery } from '@tanstack/svelte-query';

import { projectsApi } from '$lib/api/projects';
import type { CreateProjectInput } from '$lib/api/projects';
import type { KanbanBoard, Milestone, Project, ProjectDependency, ProjectStatus, UUID } from '$lib/api/types';

export interface ProjectsListOptions {
	enabled?: boolean;
}

export const useProjectsListQuery = (options: ProjectsListOptions = {}) =>
	createQuery<Project[]>({
		queryKey: ['projects', 'list'],
		queryFn: () => projectsApi.listProjects(),
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});

export const useProjectBoardQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<KanbanBoard>({
		queryKey: ['projects', projectId, 'board'],
		queryFn: () => projectsApi.getKanbanBoard(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true)
	});

export const useProjectMilestonesQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<Milestone[]>({
		queryKey: ['projects', projectId, 'milestones'],
		queryFn: () => projectsApi.listMilestones(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true),
		placeholderData: () => []
	});

export const useProjectDependenciesQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ProjectDependency[]>({
		queryKey: ['projects', projectId, 'dependencies'],
		queryFn: () => projectsApi.listDependencies(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true),
		placeholderData: () => []
	});

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

