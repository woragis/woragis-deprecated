import { apiClient } from '@clients/apiClient';
import type {
	ReportDefinition,
	ReportDefinitionDetail,
	ReportDelivery,
	ReportRun,
	ReportSchedule,
	UUID
} from './types';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface ListDefinitionsParams {
	search?: string;
	includeArchived?: boolean;
	favorites?: boolean;
	channel?: string;
	limit?: number;
	offset?: number;
}

const buildQueryString = (params: Record<string, string | number | boolean | undefined>) => {
	const searchParams = new URLSearchParams();
	for (const [key, value] of Object.entries(params)) {
		if (value === undefined || value === null || value === '') continue;
		searchParams.set(key, String(value));
	}
	return searchParams.toString() ? `?${searchParams.toString()}` : '';
};

export async function listReportDefinitions(
	params: ListDefinitionsParams = {}
): Promise<ReportDefinition[]> {
	const query = buildQueryString({
		search: params.search,
		include_archived: params.includeArchived ? 'true' : undefined,
		favorites: params.favorites ? 'true' : undefined,
		channel: params.channel,
		limit: params.limit,
		offset: params.offset
	});
	const response = await apiClient.get<ApiResponse<ReportDefinition[]>>(`/reports${query}`);
	return response.data.data ?? [];
}

export interface UpsertDefinitionInput {
	name: string;
	description?: string;
	sections?: Record<string, unknown>;
	filters?: Record<string, unknown>;
	favorite?: boolean;
}

export async function createReportDefinition(
	input: UpsertDefinitionInput
): Promise<ReportDefinition> {
	const response = await apiClient.post<ApiResponse<ReportDefinition>>('/reports', {
		name: input.name,
		description: input.description,
		sections: input.sections ?? {},
		filters: input.filters ?? {},
		favorite: input.favorite ?? false
	});
	return response.data.data;
}

export async function updateReportDefinition(
	definitionId: UUID,
	input: UpsertDefinitionInput
): Promise<ReportDefinition> {
	const response = await apiClient.put<ApiResponse<ReportDefinition>>(
		`/reports/${definitionId}`,
		{
			name: input.name,
			description: input.description,
			sections: input.sections ?? {},
			filters: input.filters ?? {},
			favorite: input.favorite ?? false
		}
	);
	return response.data.data;
}

const postBulk = async (path: string, definitionIds: UUID[]): Promise<void> => {
	if (definitionIds.length === 0) return;
	await apiClient.post(path, {
		definition_ids: definitionIds
	});
};

export async function archiveReportDefinitions(definitionIds: UUID[]): Promise<void> {
	await postBulk('/reports/archive', definitionIds);
}

export async function restoreReportDefinitions(definitionIds: UUID[]): Promise<void> {
	await postBulk('/reports/restore', definitionIds);
}

export async function deleteReportDefinitions(definitionIds: UUID[]): Promise<void> {
	await postBulk('/reports/delete', definitionIds);
}

export async function toggleReportFavorite(
	definitionId: UUID,
	favorite: boolean
): Promise<void> {
	await apiClient.post('/reports/favorite', {
		definition_id: definitionId,
		favorite
	});
}

export async function getReportDefinition(
	definitionId: UUID
): Promise<ReportDefinitionDetail> {
	const response =
		await apiClient.get<ApiResponse<ReportDefinitionDetail>>(`/reports/${definitionId}`);
	return response.data.data;
}

export async function listReportRuns(definitionId: UUID): Promise<ReportRun[]> {
	const response =
		await apiClient.get<ApiResponse<ReportRun[]>>(`/reports/${definitionId}/runs`);
	return response.data.data ?? [];
}

export async function queueReportRuns(
	definitionIds: UUID[],
	metadata?: Record<string, unknown>
): Promise<void> {
	if (definitionIds.length === 0) return;
	await apiClient.post('/reports/runs/bulk', {
		definition_ids: definitionIds,
		metadata
	});
}

export interface ScheduleInput {
	cron: string;
	frequency: string;
	timezone: string;
	nextRun?: string | null;
	enabled?: boolean;
	meta?: Record<string, unknown>;
}

export async function createReportSchedule(
	definitionId: UUID,
	input: ScheduleInput
): Promise<ReportSchedule> {
	const response = await apiClient.post<ApiResponse<ReportSchedule>>(
		`/reports/${definitionId}/schedules`,
		{
			cron: input.cron,
			frequency: input.frequency,
			timezone: input.timezone,
			next_run: input.nextRun,
			enabled: input.enabled,
			meta: input.meta
		}
	);
	return response.data.data;
}

export async function updateReportSchedule(
	scheduleId: UUID,
	input: ScheduleInput
): Promise<ReportSchedule> {
	const response = await apiClient.put<ApiResponse<ReportSchedule>>(
		`/reports/schedules/${scheduleId}`,
		{
			cron: input.cron,
			frequency: input.frequency,
			timezone: input.timezone,
			next_run: input.nextRun,
			enabled: input.enabled,
			meta: input.meta
		}
	);
	return response.data.data;
}

export async function toggleReportSchedule(
	scheduleId: UUID,
	enabled: boolean
): Promise<void> {
	await apiClient.post(`/reports/schedules/${scheduleId}/toggle`, {
		enabled
	});
}

export async function deleteReportSchedule(scheduleId: UUID): Promise<void> {
	await apiClient.delete(`/reports/schedules/${scheduleId}`, {
		data: {}
	});
}

export interface DeliveryInput {
	channel: string;
	target: string;
	template?: Record<string, unknown>;
	enabled?: boolean;
}

export async function createReportDelivery(
	definitionId: UUID,
	input: DeliveryInput
): Promise<ReportDelivery> {
	const response = await apiClient.post<ApiResponse<ReportDelivery>>(
		`/reports/${definitionId}/deliveries`,
		{
			channel: input.channel,
			target: input.target,
			template: input.template,
			enabled: input.enabled
		}
	);
	return response.data.data;
}

export async function updateReportDelivery(
	deliveryId: UUID,
	input: DeliveryInput
): Promise<ReportDelivery> {
	const response = await apiClient.put<ApiResponse<ReportDelivery>>(
		`/reports/deliveries/${deliveryId}`,
		{
			channel: input.channel,
			target: input.target,
			template: input.template,
			enabled: input.enabled
		}
	);
	return response.data.data;
}

export async function toggleReportDelivery(deliveryId: UUID, enabled: boolean): Promise<void> {
	await apiClient.post(`/reports/deliveries/${deliveryId}/toggle`, {
		enabled
	});
}

export async function deleteReportDelivery(deliveryId: UUID): Promise<void> {
	await apiClient.delete(`/reports/deliveries/${deliveryId}`, {
		data: {}
	});
}

export const reportsApi = {
	listReportDefinitions,
	createReportDefinition,
	updateReportDefinition,
	archiveReportDefinitions,
	restoreReportDefinitions,
	deleteReportDefinitions,
	toggleReportFavorite,
	getReportDefinition,
	listReportRuns,
	queueReportRuns,
	createReportSchedule,
	updateReportSchedule,
	toggleReportSchedule,
	deleteReportSchedule,
	createReportDelivery,
	updateReportDelivery,
	toggleReportDelivery,
	deleteReportDelivery
};

