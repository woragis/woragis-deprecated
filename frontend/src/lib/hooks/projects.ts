import { createMutation, createQuery } from '@tanstack/svelte-query';

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
	reorderKanbanColumns,
	updateKanbanCard,
	updateKanbanColumn,
	updateProjectMetrics,
	updateProjectStatus
} from '$lib/api/projects';
import type { CreateProjectInput } from '$lib/api/projects';
import type { KanbanBoard, Milestone, Project, ProjectDependency, ProjectStatus, UUID } from '$lib/api/types';

export interface ProjectsListOptions {
	enabled?: boolean;
}

export const useProjectsListQuery = (options: ProjectsListOptions = {}) =>
	createQuery<Project[]>({
		queryKey: ['projects', 'list'],
		queryFn: () => listProjects(),
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});

export const useProjectBoardQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<KanbanBoard>({
		queryKey: ['projects', projectId, 'board'],
		queryFn: () => getKanbanBoard(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true)
	});

export const useProjectMilestonesQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<Milestone[]>({
		queryKey: ['projects', projectId, 'milestones'],
		queryFn: () => listMilestones(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true),
		placeholderData: () => []
	});

export const useProjectDependenciesQuery = (
	projectId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ProjectDependency[]>({
		queryKey: ['projects', projectId, 'dependencies'],
		queryFn: () => listDependencies(projectId!),
		enabled: Boolean(projectId) && (options.enabled ?? true),
		placeholderData: () => []
	});

export const useCreateProjectMutation = () =>
	createMutation({
		mutationFn: (input: CreateProjectInput) => createProject(input)
	});

export const useUpdateProjectStatusMutation = () =>
	createMutation({
		mutationFn: ({ projectId, status }: { projectId: UUID; status: ProjectStatus }) =>
			updateProjectStatus(projectId, status)
	});

export const useUpdateProjectMetricsMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof updateProjectMetrics>[1];
		}) => updateProjectMetrics(projectId, payload)
	});

export const useAddMilestoneMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof addMilestone>[1];
		}) => addMilestone(projectId, payload)
	});

export const useBulkUpdateMilestonesMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			updates
		}: {
			projectId: UUID;
			updates: Parameters<typeof bulkUpdateMilestones>[1];
		}) => bulkUpdateMilestones(projectId, updates)
	});

export const useCreateKanbanColumnMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof createKanbanColumn>[1];
		}) => createKanbanColumn(projectId, payload)
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
			payload: Parameters<typeof updateKanbanColumn>[2];
		}) => updateKanbanColumn(projectId, columnId, payload)
	});

export const useReorderKanbanColumnsMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			columnOrder
		}: {
			projectId: UUID;
			columnOrder: UUID[];
		}) => reorderKanbanColumns(projectId, columnOrder)
	});

export const useDeleteKanbanColumnMutation = () =>
	createMutation({
		mutationFn: ({ projectId, columnId }: { projectId: UUID; columnId: UUID }) =>
			deleteKanbanColumn(projectId, columnId)
	});

export const useCreateKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof createKanbanCard>[1];
		}) => createKanbanCard(projectId, payload)
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
			payload: Parameters<typeof updateKanbanCard>[2];
		}) => updateKanbanCard(projectId, cardId, payload)
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
		}) => moveKanbanCard(projectId, cardId, targetColumnId, targetPosition)
	});

export const useDeleteKanbanCardMutation = () =>
	createMutation({
		mutationFn: ({ projectId, cardId }: { projectId: UUID; cardId: UUID }) =>
			deleteKanbanCard(projectId, cardId)
	});

export const useCreateDependencyMutation = () =>
	createMutation({
		mutationFn: ({
			projectId,
			payload
		}: {
			projectId: UUID;
			payload: Parameters<typeof createDependency>[1];
		}) => createDependency(projectId, payload)
	});

export const useDeleteDependencyMutation = () =>
	createMutation({
		mutationFn: ({ projectId, dependencyId }: { projectId: UUID; dependencyId: UUID }) =>
			deleteDependency(projectId, dependencyId)
	});

export const useDuplicateProjectMutation = () =>
	createMutation({
		mutationFn: ({
			templateProjectId,
			payload
		}: {
			templateProjectId: UUID;
			payload: Parameters<typeof duplicateProject>[1];
		}) => duplicateProject(templateProjectId, payload)
	});

