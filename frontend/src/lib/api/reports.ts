import { apiClient } from '@clients/apiClient';
import type { ReportDefinition } from './types';

// Re-export for convenience
export type { ReportDefinition };

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

// Helper function to map API response to ReportDefinition type
const mapReportDefinition = (apiReport: any): ReportDefinition => ({
	id: apiReport.id,
	name: apiReport.name,
	description: apiReport.description || '',
	sections: apiReport.sections || {},
	filters: apiReport.filters || {},
	is_favorite: apiReport.isFavorite ?? apiReport.is_favorite ?? false,
	archived_at: apiReport.archivedAt || apiReport.archived_at || null,
	created_at: apiReport.createdAt || apiReport.created_at,
	updated_at: apiReport.updatedAt || apiReport.updated_at
});

export interface CreateReportDefinitionInput {
	name: string;
	description?: string;
	sections?: Record<string, any>;
	filters?: Record<string, any>;
	isFavorite?: boolean;
}

export interface UpdateReportDefinitionInput {
	name?: string;
	description?: string;
	sections?: Record<string, any>;
	filters?: Record<string, any>;
	isFavorite?: boolean;
}

export interface ListReportDefinitionsOptions {
	search?: string;
	includeArchived?: boolean;
	favorites?: boolean;
	channel?: string;
	limit?: number;
	offset?: number;
}

export async function listReportDefinitions(
	options?: ListReportDefinitionsOptions
): Promise<ReportDefinition[]> {
	const params = new URLSearchParams();
	if (options?.search) params.append('search', options.search);
	if (options?.includeArchived) params.append('includeArchived', 'true');
	if (options?.favorites) params.append('favorites', 'true');
	if (options?.channel) params.append('channel', options.channel);
	if (options?.limit) params.append('limit', options.limit.toString());
	if (options?.offset) params.append('offset', options.offset.toString());

	const queryString = params.toString();
	const url = queryString ? `/reports?${queryString}` : '/reports';
	const response = await apiClient.get<ApiResponse<any[]>>(url);
	return (response.data.data ?? []).map(mapReportDefinition);
}

export async function getReportDefinition(id: string): Promise<ReportDefinition> {
	const response = await apiClient.get<ApiResponse<any>>(`/reports/${id}`);
	return mapReportDefinition(response.data.data);
}

export async function createReportDefinition(
	input: CreateReportDefinitionInput
): Promise<ReportDefinition> {
	const response = await apiClient.post<ApiResponse<any>>('/reports', input);
	return mapReportDefinition(response.data.data);
}

export async function updateReportDefinition(
	id: string,
	input: UpdateReportDefinitionInput
): Promise<ReportDefinition> {
	const response = await apiClient.put<ApiResponse<any>>(`/reports/${id}`, input);
	return mapReportDefinition(response.data.data);
}

export async function archiveReportDefinitions(ids: string[]): Promise<void> {
	await apiClient.post('/reports/archive', { ids });
}

export async function restoreReportDefinitions(ids: string[]): Promise<void> {
	await apiClient.post('/reports/restore', { ids });
}

export async function deleteReportDefinitions(ids: string[]): Promise<void> {
	await apiClient.post('/reports/delete', { ids });
}

export async function toggleFavorite(id: string, favorite: boolean): Promise<void> {
	await apiClient.post('/reports/favorite', { id, favorite });
}

// Import types
import type { ReportDefinitionDetail, ReportRun, ReportSchedule, ReportDelivery, UUID } from './types';

// Report Runs
export async function listReportRuns(definitionId: string): Promise<ReportRun[]> {
	const response = await apiClient.get<ApiResponse<ReportRun[]>>(`/reports/${definitionId}/runs`);
	return response.data.data ?? [];
}

// Report Schedules
export async function listReportSchedules(definitionId: string): Promise<ReportSchedule[]> {
	const response = await apiClient.get<ApiResponse<ReportSchedule[]>>(
		`/reports/${definitionId}/schedules`
	);
	return response.data.data ?? [];
}

export interface CreateReportScheduleInput {
	cron: string;
	frequency: string;
	timezone: string;
	nextRun?: string;
	enabled?: boolean;
	meta?: Record<string, unknown>;
}

export async function createReportSchedule(
	definitionId: string,
	input: CreateReportScheduleInput
): Promise<ReportSchedule> {
	const response = await apiClient.post<ApiResponse<ReportSchedule>>(
		`/reports/${definitionId}/schedules`,
		input
	);
	return response.data.data;
}

export interface UpdateReportScheduleInput {
	cron?: string;
	frequency?: string;
	timezone?: string;
	nextRun?: string;
	enabled?: boolean;
	meta?: Record<string, unknown>;
}

export async function updateReportSchedule(
	scheduleId: string,
	input: UpdateReportScheduleInput
): Promise<ReportSchedule> {
	const response = await apiClient.patch<ApiResponse<ReportSchedule>>(
		`/reports/schedules/${scheduleId}`,
		input
	);
	return response.data.data;
}

export async function toggleReportSchedule(scheduleId: string, enabled: boolean): Promise<void> {
	await apiClient.patch(`/reports/schedules/${scheduleId}`, { enabled });
}

export async function deleteReportSchedule(scheduleId: string): Promise<void> {
	await apiClient.delete(`/reports/schedules/${scheduleId}`);
}

// Report Deliveries
export interface CreateReportDeliveryInput {
	channel: string;
	target: string;
	enabled?: boolean;
	template?: Record<string, unknown>;
}

export async function createReportDelivery(
	definitionId: string,
	input: CreateReportDeliveryInput
): Promise<ReportDelivery> {
	const response = await apiClient.post<ApiResponse<ReportDelivery>>(
		`/reports/${definitionId}/deliveries`,
		input
	);
	return response.data.data;
}

export interface UpdateReportDeliveryInput {
	channel?: string;
	target?: string;
	enabled?: boolean;
	template?: Record<string, unknown>;
}

export async function updateReportDelivery(
	deliveryId: string,
	input: UpdateReportDeliveryInput
): Promise<ReportDelivery> {
	const response = await apiClient.patch<ApiResponse<ReportDelivery>>(
		`/reports/deliveries/${deliveryId}`,
		input
	);
	return response.data.data;
}

export async function toggleReportDelivery(deliveryId: string, enabled: boolean): Promise<void> {
	await apiClient.patch(`/reports/deliveries/${deliveryId}`, { enabled });
}

export async function deleteReportDelivery(deliveryId: string): Promise<void> {
	await apiClient.delete(`/reports/deliveries/${deliveryId}`);
}

// Queue Report Runs
export async function queueReportRuns(
	definitionIds: string[],
	metadata?: Record<string, unknown>
): Promise<void> {
	await apiClient.post('/reports/queue-runs', { definitionIds, metadata });
}

// API Object
export const reportsApi = {
	listReportDefinitions,
	getReportDefinition,
	createReportDefinition,
	updateReportDefinition,
	archiveReportDefinitions,
	restoreReportDefinitions,
	deleteReportDefinitions,
	toggleReportFavorite: toggleFavorite,
	listReportRuns,
	queueReportRuns,
	listReportSchedules,
	createReportSchedule,
	updateReportSchedule,
	toggleReportSchedule,
	deleteReportSchedule,
	createReportDelivery,
	updateReportDelivery,
	toggleReportDelivery,
	deleteReportDelivery
};