import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listSocialMediaPosts,
	getSocialMediaPost,
	getSocialMediaPostByURL
} from '$lib/api/social-media-posts';
import type { SocialMediaPost, Platform, PostStatus } from '$lib/types/social-media-post';

// Query keys factory - includes language for proper cache separation
export const socialMediaPostKeys = {
	all: ['social-media-posts'] as const,
	lists: () => [...socialMediaPostKeys.all, 'list'] as const,
	list: (params?: { platform?: Platform; status?: PostStatus }, lang?: string) =>
		[...socialMediaPostKeys.lists(), params, lang] as const,
	details: () => [...socialMediaPostKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...socialMediaPostKeys.details(), id, lang] as const,
	byUrl: (url: string, lang?: string) => [...socialMediaPostKeys.all, 'url', url, lang] as const
};

// Query options for listing social media posts
export function getSocialMediaPostsQueryOptions(params?: {
	platform?: Platform;
	status?: PostStatus;
}) {
	return queryOptions({
		queryKey: socialMediaPostKeys.list(params),
		queryFn: () => listSocialMediaPosts(params)
	});
}

// Query options for getting a social media post by ID
export function getSocialMediaPostQueryOptions(id: string) {
	return queryOptions({
		queryKey: socialMediaPostKeys.detail(id),
		queryFn: () => getSocialMediaPost(id),
		enabled: !!id
	});
}

// Query options for getting a social media post by URL
export function getSocialMediaPostByURLQueryOptions(url: string) {
	return queryOptions({
		queryKey: socialMediaPostKeys.byUrl(url),
		queryFn: () => getSocialMediaPostByURL(url),
		enabled: !!url
	});
}

// Hook for listing social media posts - reactive to language changes
export function useSocialMediaPostsQuery(params?: { platform?: Platform; status?: PostStatus }, lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: socialMediaPostKeys.list(params, currentLang),
			queryFn: () => listSocialMediaPosts(params)
		};
	});
}

// Hook for getting a social media post by ID
export function useSocialMediaPostQuery(id: string) {
	return createQuery(() => ({
		queryKey: socialMediaPostKeys.detail(id),
		queryFn: () => getSocialMediaPost(id),
		enabled: !!id
	}));
}

// Hook for getting a social media post by URL
export function useSocialMediaPostByURLQuery(url: string) {
	return createQuery(() => ({
		queryKey: socialMediaPostKeys.byUrl(url),
		queryFn: () => getSocialMediaPostByURL(url),
		enabled: !!url
	}));
}

