import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listImpactMetrics,
	getFeaturedImpactMetrics,
	getImpactMetric,
	getMetricsByEntity
} from '$lib/api/impact-metrics';
import type { ImpactMetric, MetricType, EntityType } from '$lib/types/impact-metric';

// Query keys factory - includes language for proper cache separation
export const impactMetricKeys = {
	all: ['impact-metrics'] as const,
	lists: () => [...impactMetricKeys.all, 'list'] as const,
	list: (params?: { type?: MetricType; featured?: boolean }, lang?: string) =>
		[...impactMetricKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...impactMetricKeys.all, 'featured', lang] as const,
	details: () => [...impactMetricKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...impactMetricKeys.details(), id, lang] as const,
	byEntity: (entityType: EntityType, entityId: string, lang?: string) =>
		[...impactMetricKeys.all, 'entity', entityType, entityId, lang] as const
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
export function getFeaturedImpactMetricsQueryOptions(lang?: string) {
	return queryOptions({
		queryKey: impactMetricKeys.featured(lang),
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

// Hook for featured impact metrics - reactive to language changes
export function useFeaturedImpactMetricsQuery() {
	// Read language from store in the callback - TanStack Query will track this reactively
	// The query key includes the language, so it will refetch when language changes
	return createQuery(() => {
		const currentLang = get(language);
		return {
			queryKey: impactMetricKeys.featured(currentLang),
			queryFn: () => getFeaturedImpactMetrics()
		};
	});
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

