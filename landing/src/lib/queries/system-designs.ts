import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listSystemDesigns,
	getFeaturedSystemDesigns,
	getSystemDesign
} from '$lib/api/system-designs';
import type { SystemDesign } from '$lib/types/system-design';

// Query keys factory
export const systemDesignKeys = {
	all: ['system-designs'] as const,
	lists: () => [...systemDesignKeys.all, 'list'] as const,
	list: (params?: { featured?: boolean }) => [...systemDesignKeys.lists(), params] as const,
	featured: () => [...systemDesignKeys.all, 'featured'] as const,
	details: () => [...systemDesignKeys.all, 'detail'] as const,
	detail: (id: string) => [...systemDesignKeys.details(), id] as const
};

// Query options for listing system designs
export function getSystemDesignsQueryOptions(params?: { featured?: boolean }) {
	return queryOptions({
		queryKey: systemDesignKeys.list(params),
		queryFn: () => listSystemDesigns(params)
	});
}

// Query options for featured system designs
export function getFeaturedSystemDesignsQueryOptions() {
	return queryOptions({
		queryKey: systemDesignKeys.featured(),
		queryFn: () => getFeaturedSystemDesigns()
	});
}

// Query options for getting a system design by ID
export function getSystemDesignQueryOptions(id: string) {
	return queryOptions({
		queryKey: systemDesignKeys.detail(id),
		queryFn: () => getSystemDesign(id),
		enabled: !!id
	});
}

// Hook for listing system designs
export function useSystemDesignsQuery(params?: { featured?: boolean }) {
	return createQuery(() => ({
		queryKey: systemDesignKeys.list(params),
		queryFn: () => listSystemDesigns(params)
	}));
}

// Hook for featured system designs
export function useFeaturedSystemDesignsQuery() {
	return createQuery(() => ({
		queryKey: systemDesignKeys.featured(),
		queryFn: () => getFeaturedSystemDesigns()
	}));
}

// Hook for getting a system design by ID
export function useSystemDesignQuery(id: string) {
	return createQuery(() => ({
		queryKey: systemDesignKeys.detail(id),
		queryFn: () => getSystemDesign(id),
		enabled: !!id
	}));
}

