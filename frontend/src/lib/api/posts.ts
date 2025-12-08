import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type PostStatus = 'draft' | 'published' | 'archived';

export interface Post {
	id: string;
	userId: string;
	title: string;
	slug: string;
	content: string;
	excerpt?: string;
	status: PostStatus;
	publishedAt?: string;
	featuredImage?: string;
	metaTitle?: string;
	metaDescription?: string;
	metaKeywords?: string;
	ogTitle?: string;
	ogDescription?: string;
	ogImage?: string;
	featured: boolean;
	viewsCount: number;
	createdAt: string;
	updatedAt: string;
}

export interface Category {
	id: string;
	name: string;
	slug: string;
	description?: string;
	createdAt: string;
	updatedAt: string;
}

export interface Tag {
	id: string;
	name: string;
	slug: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreatePostInput {
	title: string;
	content: string;
	excerpt?: string;
	status?: PostStatus;
	featuredImage?: string;
	metaTitle?: string;
	metaDescription?: string;
	metaKeywords?: string;
	ogTitle?: string;
	ogDescription?: string;
	ogImage?: string;
	featured?: boolean;
	skillIds?: string[];
	categoryIds?: string[];
	tagNames?: string[];
}

export interface UpdatePostInput {
	title?: string;
	content?: string;
	excerpt?: string;
	status?: PostStatus;
	featuredImage?: string;
	metaTitle?: string;
	metaDescription?: string;
	metaKeywords?: string;
	ogTitle?: string;
	ogDescription?: string;
	ogImage?: string;
	featured?: boolean;
}

export async function listPosts(): Promise<Post[]> {
	const response = await apiClient.get<ApiResponse<Post[]>>('/posts');
	return response.data.data ?? [];
}

export async function getPost(id: string): Promise<Post> {
	const response = await apiClient.get<ApiResponse<Post>>(`/posts/${id}`);
	return response.data.data;
}

export async function getPostBySlug(slug: string): Promise<Post> {
	const response = await apiClient.get<ApiResponse<Post>>(`/posts/slug/${slug}`);
	return response.data.data;
}

export async function createPost(input: CreatePostInput): Promise<Post> {
	const response = await apiClient.post<ApiResponse<Post>>('/posts', input);
	return response.data.data;
}

export async function updatePost(id: string, input: UpdatePostInput): Promise<Post> {
	const response = await apiClient.patch<ApiResponse<Post>>(`/posts/${id}`, input);
	return response.data.data;
}

export async function deletePost(id: string): Promise<void> {
	await apiClient.delete(`/posts/${id}`);
}

export async function listCategories(): Promise<Category[]> {
	const response = await apiClient.get<ApiResponse<Category[]>>('/posts/categories');
	return response.data.data ?? [];
}

export async function listTags(): Promise<Tag[]> {
	const response = await apiClient.get<ApiResponse<Tag[]>>('/posts/tags');
	return response.data.data ?? [];
}

