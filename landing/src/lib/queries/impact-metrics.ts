import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listImpactMetrics,
	getFeaturedImpactMetrics,
	getImpactMetric,
	getMetricsByEntity
} from '$lib/api/impact-metrics';
import type { ImpactMetric, MetricType, EntityType } from '$lib/types/impact-metric';

// Query keys factory
export const impactMetricKeys = {
	all: ['impact-metrics'] as const,
	lists: () => [...impactMetricKeys.all, 'list'] as const,
	list: (params?: { type?: MetricType; featured?: boolean }) =>
		[...impactMetricKeys.lists(), params] as const,
	featured: () => [...impactMetricKeys.all, 'featured'] as const,
	details: () => [...impactMetricKeys.all, 'detail'] as const,
	detail: (id: string) => [...impactMetricKeys.details(), id] as const,
	byEntity: (entityType: EntityType, entityId: string) =>
		[...impactMetricKeys.all, 'entity', entityType, entityId] as const
};

// Query options for listing impact metrics
export function getImpactMetricsQueryOptions(params?: {
	type?: MetricType;
	featured?: boolean;
}) {
	return queryOptions({
		queryKey: impactMetricKeys.list(params),
		queryFn: () => listImpactMetrics(params)
	});
}

// Query options for featured impact metrics
export function getFeaturedImpactMetricsQueryOptions() {
	return queryOptions({
		queryKey: impactMetricKeys.featured(),
		queryFn: () => getFeaturedImpactMetrics()
	});
}

// Query options for getting an impact metric by ID
export function getImpactMetricQueryOptions(id: string) {
	return queryOptions({
		queryKey: impactMetricKeys.detail(id),
		queryFn: () => getImpactMetric(id),
		enabled: !!id
	});
}

// Query options for getting metrics by entity
export function getMetricsByEntityQueryOptions(entityType: EntityType, entityId: string) {
	return queryOptions({
		queryKey: impactMetricKeys.byEntity(entityType, entityId),
		queryFn: () => getMetricsByEntity(entityType, entityId),
		enabled: !!entityType && !!entityId
	});
}

// Hook for listing impact metrics
export function useImpactMetricsQuery(params?: { type?: MetricType; featured?: boolean }) {
	return createQuery(() => ({
		queryKey: impactMetricKeys.list(params),
		queryFn: () => listImpactMetrics(params)
	}));
}

// Hook for featured impact metrics
export function useFeaturedImpactMetricsQuery() {
	return createQuery(() => ({
		queryKey: impactMetricKeys.featured(),
		queryFn: () => getFeaturedImpactMetrics()
	}));
}

// Hook for getting an impact metric by ID
export function useImpactMetricQuery(id: string) {
	return createQuery(() => ({
		queryKey: impactMetricKeys.detail(id),
		queryFn: () => getImpactMetric(id),
		enabled: !!id
	}));
}

// Hook for getting metrics by entity
export function useMetricsByEntityQuery(entityType: EntityType, entityId: string) {
	return createQuery(() => ({
		queryKey: impactMetricKeys.byEntity(entityType, entityId),
		queryFn: () => getMetricsByEntity(entityType, entityId),
		enabled: !!entityType && !!entityId
	}));
}

