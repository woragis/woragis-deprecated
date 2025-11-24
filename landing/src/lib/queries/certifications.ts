import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listCertifications,
	getFeaturedCertifications,
	getCertification,
	getCertificationsBySkill,
	type Certification,
	type CertificationStatus,
	type CertificationCategory
} from '$lib/api/certifications';

// Query keys factory
export const certificationKeys = {
	all: ['certifications'] as const,
	lists: () => [...certificationKeys.all, 'list'] as const,
	list: (params?: { status?: CertificationStatus; category?: CertificationCategory; featured?: boolean }) =>
		[...certificationKeys.lists(), params] as const,
	featured: () => [...certificationKeys.all, 'featured'] as const,
	details: () => [...certificationKeys.all, 'detail'] as const,
	detail: (id: string) => [...certificationKeys.details(), id] as const,
	bySkill: (skillId: string) => [...certificationKeys.all, 'skill', skillId] as const
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

// Hook for featured certifications
export function useFeaturedCertificationsQuery() {
	return createQuery(() => ({
		queryKey: certificationKeys.featured(),
		queryFn: () => getFeaturedCertifications()
	}));
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

