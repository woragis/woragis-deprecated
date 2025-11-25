import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listTechnicalWritings,
	getFeaturedTechnicalWritings,
	getTechnicalWriting,
	searchTechnicalWritings,
	getWritingsByType,
	getWritingsByPlatform
} from '$lib/api/technical-writings';
import type { TechnicalWriting, WritingType, PublicationPlatform } from '$lib/types/technical-writing';

// Query keys factory - includes language for proper cache separation
export const technicalWritingKeys = {
	all: ['technical-writings'] as const,
	lists: () => [...technicalWritingKeys.all, 'list'] as const,
	list: (params?: {
		type?: WritingType;
		platform?: PublicationPlatform;
		featured?: boolean;
	}, lang?: string) => [...technicalWritingKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...technicalWritingKeys.all, 'featured', lang] as const,
	details: () => [...technicalWritingKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...technicalWritingKeys.details(), id, lang] as const,
	search: (query: string, lang?: string) => [...technicalWritingKeys.all, 'search', query, lang] as const,
	byType: (type: WritingType, lang?: string) => [...technicalWritingKeys.all, 'type', type, lang] as const,
	byPlatform: (platform: PublicationPlatform, lang?: string) =>
		[...technicalWritingKeys.all, 'platform', platform, lang] as const
};

// Query options for listing technical writings
export function getTechnicalWritingsQueryOptions(params?: {
	type?: WritingType;
	platform?: PublicationPlatform;
	featured?: boolean;
}) {
	return queryOptions({
		queryKey: technicalWritingKeys.list(params),
		queryFn: () => listTechnicalWritings(params)
	});
}

// Query options for featured technical writings
export function getFeaturedTechnicalWritingsQueryOptions() {
	return queryOptions({
		queryKey: technicalWritingKeys.featured(),
		queryFn: () => getFeaturedTechnicalWritings()
	});
}

// Query options for getting a technical writing by ID
export function getTechnicalWritingQueryOptions(id: string) {
	return queryOptions({
		queryKey: technicalWritingKeys.detail(id),
		queryFn: () => getTechnicalWriting(id),
		enabled: !!id
	});
}

// Query options for searching technical writings
export function getSearchTechnicalWritingsQueryOptions(query: string) {
	return queryOptions({
		queryKey: technicalWritingKeys.search(query),
		queryFn: () => searchTechnicalWritings(query),
		enabled: !!query && query.length > 0
	});
}

// Query options for getting writings by type
export function getWritingsByTypeQueryOptions(type: WritingType) {
	return queryOptions({
		queryKey: technicalWritingKeys.byType(type),
		queryFn: () => getWritingsByType(type)
	});
}

// Query options for getting writings by platform
export function getWritingsByPlatformQueryOptions(platform: PublicationPlatform) {
	return queryOptions({
		queryKey: technicalWritingKeys.byPlatform(platform),
		queryFn: () => getWritingsByPlatform(platform)
	});
}

// Hook for listing technical writings
export function useTechnicalWritingsQuery(params?: {
	type?: WritingType;
	platform?: PublicationPlatform;
	featured?: boolean;
}) {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.list(params),
		queryFn: () => listTechnicalWritings(params)
	}));
}

// Hook for featured technical writings - reactive to language changes
export function useFeaturedTechnicalWritingsQuery(lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: technicalWritingKeys.featured(currentLang),
			queryFn: () => getFeaturedTechnicalWritings()
		};
	});
}

// Hook for getting a technical writing by ID
export function useTechnicalWritingQuery(id: string) {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.detail(id),
		queryFn: () => getTechnicalWriting(id),
		enabled: !!id
	}));
}

// Hook for searching technical writings
export function useSearchTechnicalWritingsQuery(query: string) {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.search(query),
		queryFn: () => searchTechnicalWritings(query),
		enabled: !!query && query.length > 0
	}));
}

// Hook for getting writings by type
export function useWritingsByTypeQuery(type: WritingType) {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.byType(type),
		queryFn: () => getWritingsByType(type)
	}));
}

// Hook for getting writings by platform
export function useWritingsByPlatformQuery(platform: PublicationPlatform) {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.byPlatform(platform),
		queryFn: () => getWritingsByPlatform(platform)
	}));
}

