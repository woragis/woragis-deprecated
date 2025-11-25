import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listCaseStudies,
	getCaseStudy,
	getCaseStudyByProjectSlug,
	type ListCaseStudiesParams
} from '$lib/api/case-studies';
import type { CaseStudy } from '$lib/types/case-study';

// Query keys factory - includes language for proper cache separation
export const caseStudyKeys = {
	all: ['case-studies'] as const,
	lists: () => [...caseStudyKeys.all, 'list'] as const,
	list: (params?: ListCaseStudiesParams, lang?: string) => [...caseStudyKeys.lists(), params, lang] as const,
	details: () => [...caseStudyKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...caseStudyKeys.details(), id, lang] as const,
	byProjectSlug: (slug: string, lang?: string) => [...caseStudyKeys.details(), 'project-slug', slug, lang] as const
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

// Hook for listing case studies - reactive to language changes
export function useCaseStudiesQuery(params?: ListCaseStudiesParams, lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: caseStudyKeys.list(params, currentLang),
			queryFn: () => listCaseStudies(params)
		};
	});
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

