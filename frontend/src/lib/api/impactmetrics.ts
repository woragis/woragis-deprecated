import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type MetricType =
	| 'projects_delivered'
	| 'users_impacted'
	| 'performance_improvement'
	| 'cost_savings'
	| 'time_saved';

export type MetricUnit =
	| 'count'
	| 'percentage'
	| 'currency'
	| 'hours'
	| 'days'
	| 'months'
	| 'years'
	| 'milliseconds'
	| 'seconds'
	| 'minutes';

export type EntityType = 'project' | 'problem_solution' | 'case_study' | 'system_design';

export interface ImpactMetric {
	id: string;
	userId: string;
	type: MetricType;
	value: number;
	unit: MetricUnit;
	description?: string;
	entityType?: EntityType;
	entityId?: string;
	periodStart?: string;
	periodEnd?: string;
	featured: boolean;
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
}

export interface CreateImpactMetricInput {
	type: MetricType;
	value: number;
	unit: MetricUnit;
	description?: string;
	entityType?: EntityType;
	entityId?: string;
	periodStart?: string;
	periodEnd?: string;
	featured?: boolean;
	displayOrder?: number;
}

export interface UpdateImpactMetricInput {
	type?: MetricType;
	value?: number;
	unit?: MetricUnit;
	description?: string;
	entityType?: EntityType;
	entityId?: string;
	periodStart?: string;
	periodEnd?: string;
	featured?: boolean;
	displayOrder?: number;
}

export async function listImpactMetrics(): Promise<ImpactMetric[]> {
	const response = await apiClient.get<ApiResponse<ImpactMetric[]>>('/impact-metrics');
	return response.data.data ?? [];
}

export async function getImpactMetric(id: string): Promise<ImpactMetric> {
	const response = await apiClient.get<ApiResponse<ImpactMetric>>(`/impact-metrics/${id}`);
	return response.data.data;
}

export async function createImpactMetric(input: CreateImpactMetricInput): Promise<ImpactMetric> {
	const response = await apiClient.post<ApiResponse<ImpactMetric>>('/impact-metrics', input);
	return response.data.data;
}

export async function updateImpactMetric(
	id: string,
	input: UpdateImpactMetricInput
): Promise<ImpactMetric> {
	const response = await apiClient.patch<ApiResponse<ImpactMetric>>(`/impact-metrics/${id}`, input);
	return response.data.data;
}

export async function deleteImpactMetric(id: string): Promise<void> {
	await apiClient.delete(`/impact-metrics/${id}`);
}

