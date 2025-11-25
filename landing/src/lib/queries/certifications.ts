import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listCertifications,
	getFeaturedCertifications,
	getCertification,
	getCertificationsBySkill
} from '$lib/api/certifications';
import type { Certification, CertificationStatus, CertificationCategory } from '$lib/types/certification';

// Query keys factory - includes language for proper cache separation
export const certificationKeys = {
	all: ['certifications'] as const,
	lists: () => [...certificationKeys.all, 'list'] as const,
	list: (params?: { status?: CertificationStatus; category?: CertificationCategory; featured?: boolean }, lang?: string) =>
		[...certificationKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...certificationKeys.all, 'featured', lang] as const,
	details: () => [...certificationKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...certificationKeys.details(), id, lang] as const,
	bySkill: (skillId: string, lang?: string) => [...certificationKeys.all, 'skill', skillId, lang] as const
};

// Query options for listing certifications
export function getCertificationsQueryOptions(params?: {
	status?: CertificationStatus;
	category?: CertificationCategory;
	featured?: boolean;
}) {
	return queryOptions({
		queryKey: certificationKeys.list(params),
		queryFn: () => listCertifications(params)
	});
}

// Query options for featured certifications
export function getFeaturedCertificationsQueryOptions() {
	return queryOptions({
		queryKey: certificationKeys.featured(),
		queryFn: () => getFeaturedCertifications()
	});
}

// Query options for getting a certification by ID
export function getCertificationQueryOptions(id: string) {
	return queryOptions({
		queryKey: certificationKeys.detail(id),
		queryFn: () => getCertification(id),
		enabled: !!id
	});
}

// Query options for getting certifications by skill
export function getCertificationsBySkillQueryOptions(skillId: string) {
	return queryOptions({
		queryKey: certificationKeys.bySkill(skillId),
		queryFn: () => getCertificationsBySkill(skillId),
		enabled: !!skillId
	});
}

// Hook for listing certifications
export function useCertificationsQuery(params?: {
	status?: CertificationStatus;
	category?: CertificationCategory;
	featured?: boolean;
}) {
	return createQuery(() => ({
		queryKey: certificationKeys.list(params),
		queryFn: () => listCertifications(params)
	}));
}

// Hook for featured certifications - reactive to language changes
export function useFeaturedCertificationsQuery(lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: certificationKeys.featured(currentLang),
			queryFn: () => getFeaturedCertifications()
		};
	});
}

// Hook for getting a certification by ID
export function useCertificationQuery(id: string) {
	return createQuery(() => ({
		queryKey: certificationKeys.detail(id),
		queryFn: () => getCertification(id),
		enabled: !!id
	}));
}

// Hook for getting certifications by skill
export function useCertificationsBySkillQuery(skillId: string) {
	return createQuery(() => ({
		queryKey: certificationKeys.bySkill(skillId),
		queryFn: () => getCertificationsBySkill(skillId),
		enabled: !!skillId
	}));
}

