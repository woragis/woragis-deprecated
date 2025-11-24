import { apiClient } from './client';
import type { ImpactMetric, ListImpactMetricsParams } from '$lib/types/impact-metric';

// List all impact metrics
export async function listImpactMetrics(
	params?: ListImpactMetricsParams
): Promise<ImpactMetric[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.type) {
			queryParams.append('type', params.type);
		}
		if (params?.entityType) {
			queryParams.append('entityType', params.entityType);
		}
		if (params?.entityId) {
			queryParams.append('entityId', params.entityId);
		}
		if (params?.featured !== undefined) {
			queryParams.append('featured', params.featured.toString());
		}
		if (params?.limit) {
			queryParams.append('limit', params.limit.toString());
		}
		if (params?.offset) {
			queryParams.append('offset', params.offset.toString());
		}
		if (params?.orderBy) {
			queryParams.append('orderBy', params.orderBy);
		}
		if (params?.order) {
			queryParams.append('order', params.order);
		}

		const queryString = queryParams.toString();
		const url = queryString ? `/impact-metrics?${queryString}` : '/impact-metrics';

		const response = await apiClient.get<{ success: boolean; data: ImpactMetric[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching impact metrics:', error);
		throw error;
	}
}

// Get featured impact metrics (public endpoint)
export async function getFeaturedImpactMetrics(): Promise<ImpactMetric[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ImpactMetric[] }>(
			'/impact-metrics/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured impact metrics:', error);
		throw error;
	}
}

// Get impact metric by ID
export async function getImpactMetric(id: string): Promise<ImpactMetric | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ImpactMetric }>(
			`/impact-metrics/${id}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching impact metric ${id}:`, error);
		return null;
	}
}

// Get metrics by entity
export async function getMetricsByEntity(
	entityType: string,
	entityId: string
): Promise<ImpactMetric[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ImpactMetric[] }>(
			`/impact-metrics/entity/${entityType}/${entityId}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching metrics for entity ${entityType}/${entityId}:`, error);
		return [];
	}
}

