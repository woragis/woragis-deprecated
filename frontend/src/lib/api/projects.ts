import { apiClient } from '@clients/apiClient';
import type {
	KanbanBoard,
	Milestone,
	Project,
	ProjectDependency,
	ProjectStatus,
	UUID
} from './types';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface CreateProjectInput {
	name: string;
	description: string;
	status?: ProjectStatus;
	healthScore?: number;
	mrr?: number;
	cac?: number;
	ltv?: number;
	churnRate?: number;
}

export async function listProjects(): Promise<Project[]> {
	const response = await apiClient.get<ApiResponse<Project[]>>('/projects');
	return response.data.data ?? [];
}

export async function createProject(input: CreateProjectInput): Promise<Project> {
	const response = await apiClient.post<ApiResponse<Project>>('/projects', {
		name: input.name,
		description: input.description,
		status: input.status ?? 'planning',
		health_score: input.healthScore ?? 50,
		mrr: input.mrr ?? 0,
		cac: input.cac ?? 0,
		ltv: input.ltv ?? 0,
		churn_rate: input.churnRate ?? 0
	});
	return response.data.data;
}

export async function updateProjectStatus(
	projectId: UUID,
	status: ProjectStatus
): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<Project>>(`/projects/${projectId}/status`, {
		status
	});
	return response.data.data;
}

export async function updateProjectMetrics(
	projectId: UUID,
	payload: Partial<{
		healthScore: number;
		mrr: number;
		cac: number;
		ltv: number;
		churnRate: number;
	}>
): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<Project>>(`/projects/${projectId}/metrics`, {
		health_score: payload.healthScore,
		mrr: payload.mrr,
		cac: payload.cac,
		ltv: payload.ltv,
		churn_rate: payload.churnRate
	});
	return response.data.data;
}

export async function addMilestone(
	projectId: UUID,
	payload: { title: string; description: string; dueDate: string }
): Promise<Milestone> {
	const response = await apiClient.post<ApiResponse<Milestone>>(
		`/projects/${projectId}/milestones`,
		{
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate
		}
	);
	return response.data.data;
}

export async function listMilestones(projectId: UUID): Promise<Milestone[]> {
	const response = await apiClient.get<ApiResponse<Milestone[]>>(
		`/projects/${projectId}/milestones`
	);
	return response.data.data ?? [];
}

export async function bulkUpdateMilestones(
	projectId: UUID,
	updates: Array<{
		milestoneId: UUID;
		title?: string;
		description?: string;
		dueDate?: string;
		completed?: boolean;
	}>
): Promise<Milestone[]> {
	const payload = updates.map((item) => ({
		milestone_id: item.milestoneId,
		title: item.title,
		description: item.description,
		due_date: item.dueDate,
		completed: item.completed
	}));

	const response = await apiClient.post<ApiResponse<Milestone[]>>(
		`/projects/${projectId}/milestones/bulk`,
		{
			updates: payload
		}
	);
	return response.data.data ?? [];
}

export async function toggleMilestone(milestoneId: UUID, completed: boolean): Promise<Milestone> {
	const response = await apiClient.patch<ApiResponse<Milestone>>(
		`/projects/milestones/${milestoneId}`,
		{
			completed
		}
	);
	return response.data.data;
}

export async function getKanbanBoard(projectId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.get<ApiResponse<KanbanBoard>>(`/projects/${projectId}/kanban`);
	return response.data.data;
}

export async function createKanbanColumn(
	projectId: UUID,
	payload: { name: string; position?: number; wipLimit?: number }
): Promise<KanbanBoard> {
	const response = await apiClient.post<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns`,
		{
			name: payload.name,
			position: payload.position,
			wip_limit: payload.wipLimit
		}
	);
	return response.data.data;
}

export async function updateKanbanColumn(
	projectId: UUID,
	columnId: UUID,
	payload: { name?: string; wipLimit?: number }
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/${columnId}`,
		{
			name: payload.name,
			wip_limit: payload.wipLimit
		}
	);
	return response.data.data;
}

export async function reorderKanbanColumns(
	projectId: UUID,
	columnOrder: UUID[]
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/reorder`,
		{
			column_order: columnOrder
		}
	);
	return response.data.data;
}

export async function deleteKanbanColumn(projectId: UUID, columnId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.delete<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/${columnId}`,
		{
			data: {}
		}
	);
	return response.data.data;
}

export async function createKanbanCard(
	projectId: UUID,
	payload: {
		columnId: UUID;
		title: string;
		description?: string;
		dueDate?: string;
		milestoneId?: UUID;
		position?: number;
	}
): Promise<KanbanBoard> {
	const response = await apiClient.post<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards`,
		{
			column_id: payload.columnId,
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate,
			milestone_id: payload.milestoneId,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function updateKanbanCard(
	projectId: UUID,
	cardId: UUID,
	payload: {
		title?: string;
		description?: string;
		dueDate?: string;
		milestoneId?: UUID;
	}
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}`,
		{
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate,
			milestone_id: payload.milestoneId
		}
	);
	return response.data.data;
}

export async function moveKanbanCard(
	projectId: UUID,
	cardId: UUID,
	targetColumnId: UUID,
	targetPosition: number
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}/move`,
		{
			target_column_id: targetColumnId,
			target_position: targetPosition
		}
	);
	return response.data.data;
}

export async function deleteKanbanCard(projectId: UUID, cardId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.delete<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}`,
		{
			data: {}
		}
	);
	return response.data.data;
}

export async function createDependency(
	projectId: UUID,
	payload: { dependsOnProjectId: UUID; type: 'blocks' | 'relates' | 'supports' }
): Promise<ProjectDependency> {
	const response = await apiClient.post<ApiResponse<ProjectDependency>>(
		`/projects/${projectId}/dependencies`,
		{
			depends_on_project_id: payload.dependsOnProjectId,
			type: payload.type
		}
	);
	return response.data.data;
}

export async function listDependencies(projectId: UUID): Promise<ProjectDependency[]> {
	const response = await apiClient.get<ApiResponse<ProjectDependency[]>>(
		`/projects/${projectId}/dependencies`
	);
	return response.data.data ?? [];
}

export async function deleteDependency(projectId: UUID, dependencyId: UUID): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/dependencies/${dependencyId}`, {
		data: {}
	});
}

export async function duplicateProject(
	templateProjectId: UUID,
	payload: {
		name: string;
		description?: string;
		status?: ProjectStatus;
		healthScore?: number;
		mrr?: number;
		cac?: number;
		ltv?: number;
		churnRate?: number;
		copyBoard?: boolean;
		copyMilestones?: boolean;
		copyDependencies?: boolean;
	}
): Promise<Project> {
	const response = await apiClient.post<ApiResponse<Project>>(
		`/projects/${templateProjectId}/duplicate`,
		{
			name: payload.name,
			description: payload.description,
			status: payload.status,
			health_score: payload.healthScore,
			mrr: payload.mrr,
			cac: payload.cac,
			ltv: payload.ltv,
			churn_rate: payload.churnRate,
			copy_board: payload.copyBoard ?? true,
			copy_milestones: payload.copyMilestones ?? true,
			copy_dependencies: payload.copyDependencies ?? false
		}
	);
	return response.data.data;
}
