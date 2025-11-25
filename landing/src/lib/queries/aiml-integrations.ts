import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listAIMLIntegrations,
	getFeaturedAIMLIntegrations,
	getAIMLIntegration,
	getIntegrationsByType,
	getIntegrationsByFramework
} from '$lib/api/aiml-integrations';
import type { AIMLIntegration, IntegrationType, Framework } from '$lib/types/aiml-integration';

// Query keys factory - includes language for proper cache separation
export const aimlIntegrationKeys = {
	all: ['aiml-integrations'] as const,
	lists: () => [...aimlIntegrationKeys.all, 'list'] as const,
	list: (params?: {
		type?: IntegrationType;
		framework?: Framework;
		featured?: boolean;
	}, lang?: string) => [...aimlIntegrationKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...aimlIntegrationKeys.all, 'featured', lang] as const,
	details: () => [...aimlIntegrationKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...aimlIntegrationKeys.details(), id, lang] as const,
	byType: (type: IntegrationType, lang?: string) => [...aimlIntegrationKeys.all, 'type', type, lang] as const,
	byFramework: (framework: Framework, lang?: string) =>
		[...aimlIntegrationKeys.all, 'framework', framework, lang] as const
};

// Query options for listing AI/ML integrations
export function getAIMLIntegrationsQueryOptions(params?: {
	type?: IntegrationType;
	framework?: Framework;
	featured?: boolean;
}) {
	return queryOptions({
		queryKey: aimlIntegrationKeys.list(params),
		queryFn: () => listAIMLIntegrations(params)
	});
}

// Query options for featured AI/ML integrations
export function getFeaturedAIMLIntegrationsQueryOptions() {
	return queryOptions({
		queryKey: aimlIntegrationKeys.featured(),
		queryFn: () => getFeaturedAIMLIntegrations()
	});
}

// Query options for getting an AI/ML integration by ID
export function getAIMLIntegrationQueryOptions(id: string) {
	return queryOptions({
		queryKey: aimlIntegrationKeys.detail(id),
		queryFn: () => getAIMLIntegration(id),
		enabled: !!id
	});
}

// Query options for getting integrations by type
export function getIntegrationsByTypeQueryOptions(type: IntegrationType) {
	return queryOptions({
		queryKey: aimlIntegrationKeys.byType(type),
		queryFn: () => getIntegrationsByType(type)
	});
}

// Query options for getting integrations by framework
export function getIntegrationsByFrameworkQueryOptions(framework: Framework) {
	return queryOptions({
		queryKey: aimlIntegrationKeys.byFramework(framework),
		queryFn: () => getIntegrationsByFramework(framework)
	});
}

// Hook for listing AI/ML integrations
export function useAIMLIntegrationsQuery(params?: {
	type?: IntegrationType;
	framework?: Framework;
	featured?: boolean;
}) {
	return createQuery(() => ({
		queryKey: aimlIntegrationKeys.list(params),
		queryFn: () => listAIMLIntegrations(params)
	}));
}

// Hook for featured AI/ML integrations - reactive to language changes
export function useFeaturedAIMLIntegrationsQuery(lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: aimlIntegrationKeys.featured(currentLang),
			queryFn: () => getFeaturedAIMLIntegrations()
		};
	});
}

// Hook for getting an AI/ML integration by ID
export function useAIMLIntegrationQuery(id: string) {
	return createQuery(() => ({
		queryKey: aimlIntegrationKeys.detail(id),
		queryFn: () => getAIMLIntegration(id),
		enabled: !!id
	}));
}

// Hook for getting integrations by type
export function useIntegrationsByTypeQuery(type: IntegrationType) {
	return createQuery(() => ({
		queryKey: aimlIntegrationKeys.byType(type),
		queryFn: () => getIntegrationsByType(type)
	}));
}

// Hook for getting integrations by framework
export function useIntegrationsByFrameworkQuery(framework: Framework) {
	return createQuery(() => ({
		queryKey: aimlIntegrationKeys.byFramework(framework),
		queryFn: () => getIntegrationsByFramework(framework)
	}));
}

