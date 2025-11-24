import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listSocialMediaPosts,
	getSocialMediaPost,
	getSocialMediaPostByURL,
	type SocialMediaPost,
	type Platform,
	type PostStatus
} from '$lib/api/social-media-posts';

// Query keys factory
export const socialMediaPostKeys = {
	all: ['social-media-posts'] as const,
	lists: () => [...socialMediaPostKeys.all, 'list'] as const,
	list: (params?: { platform?: Platform; status?: PostStatus }) =>
		[...socialMediaPostKeys.lists(), params] as const,
	details: () => [...socialMediaPostKeys.all, 'detail'] as const,
	detail: (id: string) => [...socialMediaPostKeys.details(), id] as const,
	byUrl: (url: string) => [...socialMediaPostKeys.all, 'url', url] as const
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

// Hook for listing social media posts
export function useSocialMediaPostsQuery(params?: { platform?: Platform; status?: PostStatus }) {
	return createQuery(() => ({
		queryKey: socialMediaPostKeys.list(params),
		queryFn: () => listSocialMediaPosts(params)
	}));
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

