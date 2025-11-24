import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listTechnicalWritings,
	getFeaturedTechnicalWritings,
	getTechnicalWriting,
	searchTechnicalWritings,
	getWritingsByType,
	getWritingsByPlatform,
	type TechnicalWriting,
	type WritingType,
	type PublicationPlatform
} from '$lib/api/technical-writings';

// Query keys factory
export const technicalWritingKeys = {
	all: ['technical-writings'] as const,
	lists: () => [...technicalWritingKeys.all, 'list'] as const,
	list: (params?: {
		type?: WritingType;
		platform?: PublicationPlatform;
		featured?: boolean;
	}) => [...technicalWritingKeys.lists(), params] as const,
	featured: () => [...technicalWritingKeys.all, 'featured'] as const,
	details: () => [...technicalWritingKeys.all, 'detail'] as const,
	detail: (id: string) => [...technicalWritingKeys.details(), id] as const,
	search: (query: string) => [...technicalWritingKeys.all, 'search', query] as const,
	byType: (type: WritingType) => [...technicalWritingKeys.all, 'type', type] as const,
	byPlatform: (platform: PublicationPlatform) =>
		[...technicalWritingKeys.all, 'platform', platform] as const
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

// Hook for featured technical writings
export function useFeaturedTechnicalWritingsQuery() {
	return createQuery(() => ({
		queryKey: technicalWritingKeys.featured(),
		queryFn: () => getFeaturedTechnicalWritings()
	}));
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

