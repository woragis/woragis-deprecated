import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import { listInterests, getFeaturedInterests, type Interest } from '$lib/api/interests';

// Query keys factory - includes language for proper cache separation
export const interestKeys = {
	all: ['interests'] as const,
	lists: () => [...interestKeys.all, 'list'] as const,
	list: (lang: string) => [...interestKeys.lists(), lang] as const,
	featured: () => [...interestKeys.all, 'featured'] as const,
	featuredList: (lang: string) => [...interestKeys.featured(), lang] as const,
	details: () => [...interestKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...interestKeys.details(), id, lang] as const
};

// Query options for listing interests
export function getInterestsQueryOptions(lang: string) {
	return queryOptions({
		queryKey: interestKeys.list(lang),
		queryFn: () => listInterests(),
		enabled: true
	});
}

// Query options for getting featured interests
export function getFeaturedInterestsQueryOptions(lang: string) {
	return queryOptions({
		queryKey: interestKeys.featuredList(lang),
		queryFn: () => getFeaturedInterests(),
		enabled: true
	});
}

// Hook for listing interests - reactive to language changes
export function useInterestsQuery(lang: string) {
	return createQuery(() => {
		const currentLang = get(language);
		return {
			queryKey: interestKeys.list(currentLang),
			queryFn: () => listInterests()
		};
	});
}

// Hook for getting featured interests - reactive to language changes
export function useFeaturedInterestsQuery(lang: string) {
	return createQuery(() => {
		const currentLang = get(language);
		return {
			queryKey: interestKeys.featuredList(currentLang),
			queryFn: () => getFeaturedInterests()
		};
	});
}

