import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { listPosts, getPostBySlug, getPost, listCategories, listTags } from '$lib/api/posts';
import type { Post, Category, Tag, ListPostsParams } from '$lib/types/post';

// Query keys factory
export const postKeys = {
	all: ['posts'] as const,
	lists: () => [...postKeys.all, 'list'] as const,
	list: (params?: ListPostsParams) => [...postKeys.lists(), params] as const,
	details: () => [...postKeys.all, 'detail'] as const,
	detail: (id: string) => [...postKeys.details(), id] as const,
	bySlug: (slug: string) => [...postKeys.details(), 'slug', slug] as const,
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

// Hook for listing posts
export function usePostsQuery(params?: ListPostsParams) {
	return createQuery(() => ({
		queryKey: postKeys.list(params),
		queryFn: () => listPosts(params)
	}));
}

// Hook for getting a post by ID
export function usePostQuery(id: string) {
	return createQuery(() => ({
		queryKey: postKeys.detail(id),
		queryFn: () => getPost(id),
		enabled: !!id
	}));
}

// Hook for getting a post by slug
export function usePostBySlugQuery(slug: string) {
	return createQuery(() => ({
		queryKey: postKeys.bySlug(slug),
		queryFn: () => getPostBySlug(slug),
		enabled: !!slug
	}));
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

