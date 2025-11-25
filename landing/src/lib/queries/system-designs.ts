import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listSystemDesigns,
	getFeaturedSystemDesigns,
	getSystemDesign
} from '$lib/api/system-designs';
import type { SystemDesign } from '$lib/types/system-design';

// Query keys factory - includes language for proper cache separation
export const systemDesignKeys = {
	all: ['system-designs'] as const,
	lists: () => [...systemDesignKeys.all, 'list'] as const,
	list: (params?: { featured?: boolean }, lang?: string) => [...systemDesignKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...systemDesignKeys.all, 'featured', lang] as const,
	details: () => [...systemDesignKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...systemDesignKeys.details(), id, lang] as const
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

// Hook for featured system designs - reactive to language changes
export function useFeaturedSystemDesignsQuery(lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: systemDesignKeys.featured(currentLang),
			queryFn: () => getFeaturedSystemDesigns()
		};
	});
}

// Hook for getting a system design by ID
export function useSystemDesignQuery(id: string) {
	return createQuery(() => ({
		queryKey: systemDesignKeys.detail(id),
		queryFn: () => getSystemDesign(id),
		enabled: !!id
	}));
}

