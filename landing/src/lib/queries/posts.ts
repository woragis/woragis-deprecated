import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import { listPosts, getPostBySlug, getPost, listCategories, listTags } from '$lib/api/posts';
import type { Post, Category, Tag, ListPostsParams } from '$lib/types/post';

// Query keys factory - includes language for proper cache separation
export const postKeys = {
	all: ['posts'] as const,
	lists: () => [...postKeys.all, 'list'] as const,
	list: (params?: ListPostsParams, lang?: string) => [...postKeys.lists(), params, lang] as const,
	details: () => [...postKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...postKeys.details(), id, lang] as const,
	bySlug: (slug: string, lang?: string) => [...postKeys.details(), 'slug', slug, lang] as const,
	categories: () => [...postKeys.all, 'categories'] as const,
	tags: () => [...postKeys.all, 'tags'] as const
};

// Query options for listing posts
export function getPostsQueryOptions(params?: ListPostsParams) {
	return queryOptions({
		queryKey: postKeys.list(params),
		queryFn: () => listPosts(params)
	});
}

// Query options for getting a post by ID
export function getPostQueryOptions(id: string) {
	return queryOptions({
		queryKey: postKeys.detail(id),
		queryFn: () => getPost(id),
		enabled: !!id
	});
}

// Query options for getting a post by slug
export function getPostBySlugQueryOptions(slug: string) {
	return queryOptions({
		queryKey: postKeys.bySlug(slug),
		queryFn: () => getPostBySlug(slug),
		enabled: !!slug
	});
}

// Query options for listing categories
export function getCategoriesQueryOptions() {
	return queryOptions({
		queryKey: postKeys.categories(),
		queryFn: () => listCategories()
	});
}

// Query options for listing tags
export function getTagsQueryOptions() {
	return queryOptions({
		queryKey: postKeys.tags(),
		queryFn: () => listTags()
	});
}

// Hook for listing posts - reactive to language changes
export function usePostsQuery(params?: ListPostsParams, lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: postKeys.list(params, currentLang),
			queryFn: () => listPosts(params)
		};
	});
}

// Hook for getting a post by ID - reactive to language changes
export function usePostQuery(id: string, lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: postKeys.detail(id, currentLang),
			queryFn: () => getPost(id),
			enabled: !!id
		};
	});
}

// Hook for getting a post by slug - reactive to language changes
export function usePostBySlugQuery(slug: string, lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: postKeys.bySlug(slug, currentLang),
			queryFn: () => getPostBySlug(slug),
			enabled: !!slug
		};
	});
}

// Hook for listing categories
export function useCategoriesQuery() {
	return createQuery(() => ({
		queryKey: postKeys.categories(),
		queryFn: () => listCategories()
	}));
}

// Hook for listing tags
export function useTagsQuery() {
	return createQuery(() => ({
		queryKey: postKeys.tags(),
		queryFn: () => listTags()
	}));
}

