import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listSkills,
	listSkillsWithCounts,
	getSkill,
	getSkillBySlug,
	searchSkills,
	listSkillsByCategory,
	getSkillsTimeline,
	type Skill,
	type SkillWithCount,
	type SkillCategory
} from '$lib/api/skills';

// Query keys factory - includes language for proper cache separation
export const skillKeys = {
	all: ['skills'] as const,
	lists: () => [...skillKeys.all, 'list'] as const,
	list: (lang?: string) => [...skillKeys.lists(), lang] as const,
	withCounts: (lang?: string) => [...skillKeys.all, 'with-counts', lang] as const,
	details: () => [...skillKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...skillKeys.details(), id, lang] as const,
	bySlug: (slug: string, lang?: string) => [...skillKeys.details(), 'slug', slug, lang] as const,
	search: (query: string, lang?: string) => [...skillKeys.all, 'search', query, lang] as const,
	byCategory: (category: SkillCategory, lang?: string) => [...skillKeys.all, 'category', category, lang] as const,
	timeline: (lang?: string) => [...skillKeys.all, 'timeline', lang] as const
};

// Query options for listing skills
export function getSkillsQueryOptions() {
	return queryOptions({
		queryKey: skillKeys.list(),
		queryFn: () => listSkills()
	});
}

// Query options for listing skills with counts
export function getSkillsWithCountsQueryOptions() {
	return queryOptions({
		queryKey: skillKeys.withCounts(),
		queryFn: () => listSkillsWithCounts()
	});
}

// Query options for getting a skill by ID
export function getSkillQueryOptions(id: string) {
	return queryOptions({
		queryKey: skillKeys.detail(id),
		queryFn: () => getSkill(id),
		enabled: !!id
	});
}

// Query options for getting a skill by slug
export function getSkillBySlugQueryOptions(slug: string) {
	return queryOptions({
		queryKey: skillKeys.bySlug(slug),
		queryFn: () => getSkillBySlug(slug),
		enabled: !!slug
	});
}

// Query options for searching skills
export function getSearchSkillsQueryOptions(query: string) {
	return queryOptions({
		queryKey: skillKeys.search(query),
		queryFn: () => searchSkills(query),
		enabled: !!query && query.length > 0
	});
}

// Query options for listing skills by category
export function getSkillsByCategoryQueryOptions(category: SkillCategory) {
	return queryOptions({
		queryKey: skillKeys.byCategory(category),
		queryFn: () => listSkillsByCategory(category)
	});
}

// Hook for listing skills
export function useSkillsQuery() {
	return createQuery(() => ({
		queryKey: skillKeys.list(),
		queryFn: () => listSkills()
	}));
}

// Hook for listing skills with counts
export function useSkillsWithCountsQuery() {
	return createQuery(() => ({
		queryKey: skillKeys.withCounts(),
		queryFn: () => listSkillsWithCounts()
	}));
}

// Hook for getting a skill by ID
export function useSkillQuery(id: string) {
	return createQuery(() => ({
		queryKey: skillKeys.detail(id),
		queryFn: () => getSkill(id),
		enabled: !!id
	}));
}

// Hook for getting a skill by slug
export function useSkillBySlugQuery(slug: string) {
	return createQuery(() => ({
		queryKey: skillKeys.bySlug(slug),
		queryFn: () => getSkillBySlug(slug),
		enabled: !!slug
	}));
}

// Hook for searching skills
export function useSearchSkillsQuery(query: string) {
	return createQuery(() => ({
		queryKey: skillKeys.search(query),
		queryFn: () => searchSkills(query),
		enabled: !!query && query.length > 0
	}));
}

// Hook for listing skills by category
export function useSkillsByCategoryQuery(category: SkillCategory) {
	return createQuery(() => ({
		queryKey: skillKeys.byCategory(category),
		queryFn: () => listSkillsByCategory(category)
	}));
}

// Query options for getting skills timeline
export function getSkillsTimelineQueryOptions(lang?: string) {
	return queryOptions({
		queryKey: skillKeys.timeline(lang),
		queryFn: () => getSkillsTimeline()
	});
}

// Hook for getting skills timeline - reactive to language changes
export function useSkillsTimelineQuery() {
	// Read language from store in the callback - TanStack Query will track this reactively
	// The query key includes the language, so it will refetch when language changes
	return createQuery(() => {
		const currentLang = get(language);
		return {
			queryKey: skillKeys.timeline(currentLang),
			queryFn: () => getSkillsTimeline()
		};
	});
}

