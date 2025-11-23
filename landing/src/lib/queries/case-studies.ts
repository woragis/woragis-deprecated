import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listCaseStudies,
	getCaseStudy,
	getCaseStudyByProjectSlug,
	type ListCaseStudiesParams
} from '$lib/api/case-studies';
import type { CaseStudy } from '$lib/types/case-study';

// Query keys factory
export const caseStudyKeys = {
	all: ['case-studies'] as const,
	lists: () => [...caseStudyKeys.all, 'list'] as const,
	list: (params?: ListCaseStudiesParams) => [...caseStudyKeys.lists(), params] as const,
	details: () => [...caseStudyKeys.all, 'detail'] as const,
	detail: (id: string) => [...caseStudyKeys.details(), id] as const,
	byProjectSlug: (slug: string) => [...caseStudyKeys.details(), 'project-slug', slug] as const
};

// Query options for listing case studies
export function getCaseStudiesQueryOptions(params?: ListCaseStudiesParams) {
	return queryOptions({
		queryKey: caseStudyKeys.list(params),
		queryFn: () => listCaseStudies(params)
	});
}

// Query options for getting a case study by ID
export function getCaseStudyQueryOptions(id: string) {
	return queryOptions({
		queryKey: caseStudyKeys.detail(id),
		queryFn: () => getCaseStudy(id),
		enabled: !!id
	});
}

// Query options for getting a case study by project slug
export function getCaseStudyByProjectSlugQueryOptions(projectSlug: string) {
	return queryOptions({
		queryKey: caseStudyKeys.byProjectSlug(projectSlug),
		queryFn: () => getCaseStudyByProjectSlug(projectSlug),
		enabled: !!projectSlug
	});
}

// Hook for listing case studies
export function useCaseStudiesQuery(params?: ListCaseStudiesParams) {
	return createQuery(() => ({
		queryKey: caseStudyKeys.list(params),
		queryFn: () => listCaseStudies(params)
	}));
}

// Hook for getting a case study by ID
export function useCaseStudyQuery(id: string) {
	return createQuery(() => ({
		queryKey: caseStudyKeys.detail(id),
		queryFn: () => getCaseStudy(id),
		enabled: !!id
	}));
}

// Hook for getting a case study by project slug
export function useCaseStudyByProjectSlugQuery(projectSlug: string) {
	return createQuery(() => ({
		queryKey: caseStudyKeys.byProjectSlug(projectSlug),
		queryFn: () => getCaseStudyByProjectSlug(projectSlug),
		enabled: !!projectSlug
	}));
}

