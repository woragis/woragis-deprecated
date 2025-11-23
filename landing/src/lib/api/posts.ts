import { apiClient } from './client';
import type { Post, ListPostsParams, Category, Tag } from '$lib/types/post';

// Calculate reading time from markdown content (average 200 words per minute)
export function calculateReadingTime(content: string): number {
	const words = content.split(/\s+/).length;
	return Math.ceil(words / 200); // minutes
}

// List all posts
export async function listPosts(params?: ListPostsParams): Promise<Post[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.status) {
			queryParams.append('status', params.status);
		}
		if (params?.featured !== undefined) {
			queryParams.append('featured', params.featured.toString());
		}
		if (params?.categoryId) {
			queryParams.append('categoryId', params.categoryId);
		}
		if (params?.tagId) {
			queryParams.append('tagId', params.tagId);
		}
		if (params?.skillId) {
			queryParams.append('skillId', params.skillId);
		}
		if (params?.search) {
			queryParams.append('search', params.search);
		}
		if (params?.limit) {
			queryParams.append('limit', params.limit.toString());
		}
		if (params?.offset) {
			queryParams.append('offset', params.offset.toString());
		}
		if (params?.orderBy) {
			queryParams.append('orderBy', params.orderBy);
		}
		if (params?.order) {
			queryParams.append('order', params.order);
		}

		const queryString = queryParams.toString();
		const url = queryString ? `/posts?${queryString}` : '/posts';

		const response = await apiClient.get<{ success: boolean; data: Post[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching posts:', error);
		throw error;
	}
}

// Get post by slug
export async function getPostBySlug(slug: string): Promise<Post | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Post }>(
			`/posts/slug/${slug}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching post by slug ${slug}:`, error);
		return null;
	}
}

// Get post by ID
export async function getPost(id: string): Promise<Post | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Post }>(`/posts/${id}`);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching post ${id}:`, error);
		return null;
	}
}

// List categories
export async function listCategories(): Promise<Category[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Category[] }>(
			'/posts/categories'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching categories:', error);
		return [];
	}
}

// List tags
export async function listTags(): Promise<Tag[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Tag[] }>('/posts/tags');
		return response.data || [];
	} catch (error) {
		console.error('Error fetching tags:', error);
		return [];
	}
}

// Get post skills
export async function getPostSkills(postId: string): Promise<Array<{ id: string; name: string; slug: string }>> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Array<{ id: string; name: string; slug: string }> }>(
			`/posts/${postId}/skills`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching skills for post ${postId}:`, error);
		return [];
	}
}

// Get post categories
export async function getPostCategories(postId: string): Promise<Category[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Category[] }>(
			`/posts/${postId}/categories`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching categories for post ${postId}:`, error);
		return [];
	}
}

// Get post tags
export async function getPostTags(postId: string): Promise<Tag[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Tag[] }>(
			`/posts/${postId}/tags`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching tags for post ${postId}:`, error);
		return [];
	}
}

