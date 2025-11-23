import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listSkills,
	listSkillsWithCounts,
	getSkill,
	getSkillBySlug,
	searchSkills,
	listSkillsByCategory,
	type Skill,
	type SkillWithCount,
	type SkillCategory
} from '$lib/api/skills';

// Query keys factory
export const skillKeys = {
	all: ['skills'] as const,
	lists: () => [...skillKeys.all, 'list'] as const,
	list: () => [...skillKeys.lists()] as const,
	withCounts: () => [...skillKeys.all, 'with-counts'] as const,
	details: () => [...skillKeys.all, 'detail'] as const,
	detail: (id: string) => [...skillKeys.details(), id] as const,
	bySlug: (slug: string) => [...skillKeys.details(), 'slug', slug] as const,
	search: (query: string) => [...skillKeys.all, 'search', query] as const,
	byCategory: (category: SkillCategory) => [...skillKeys.all, 'category', category] as const
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

