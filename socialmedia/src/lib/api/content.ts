import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface ContentPost {
	id: string;
	postId: string;
	contentType?: string;
	project?: string;
	priority: 'low' | 'medium' | 'high';
	status: 'pending' | 'in_progress' | 'completed' | 'archived';
	createdAt: string;
	updatedAt: string;
}

export interface ContentPostWithSocialPosts {
	contentPost: ContentPost;
	socialPosts: Array<{
		id: string;
		platform: string;
		format: string;
		title: string;
		status: string;
	}>;
}

export interface CreateContentPostInput {
	postId: string;
	contentType?: string;
	project?: string;
	priority?: 'low' | 'medium' | 'high';
}

export interface RepurposeRequest {
	platforms: Array<{
		platform: string;
		format: string;
		title: string;
		content: string;
	}>;
}

export async function listContentPosts(filters?: {
	status?: string;
	priority?: string;
}): Promise<ContentPost[]> {
	const params = new URLSearchParams();
	if (filters?.status) params.append('status', filters.status);
	if (filters?.priority) params.append('priority', filters.priority);

	const response = await apiClient.get<ApiResponse<ContentPost[]>>(
		`/social-media-posts/content/posts?${params.toString()}`
	);
	return response.data.data ?? [];
}

export async function getContentPost(id: string): Promise<ContentPostWithSocialPosts> {
	const response = await apiClient.get<ApiResponse<ContentPostWithSocialPosts>>(
		`/social-media-posts/content/posts/${id}`
	);
	return response.data.data;
}

export async function createContentPost(input: CreateContentPostInput): Promise<ContentPost> {
	const response = await apiClient.post<ApiResponse<ContentPost>>(
		'/social-media-posts/content/posts',
		input
	);
	return response.data.data;
}

export async function updateContentPostPriority(
	id: string,
	priority: 'low' | 'medium' | 'high'
): Promise<ContentPost> {
	const response = await apiClient.patch<ApiResponse<ContentPost>>(
		`/social-media-posts/content/posts/${id}/priority`,
		{ priority }
	);
	return response.data.data;
}

export async function repurposeToPlatforms(
	id: string,
	request: RepurposeRequest
): Promise<Array<{ id: string; platform: string; format: string; title: string; status: string }>> {
	const response = await apiClient.post<ApiResponse<Array<{ id: string; platform: string; format: string; title: string; status: string }>>>(
		`/social-media-posts/content/posts/${id}/repurpose`,
		request
	);
	return response.data.data ?? [];
}

export async function getContentBacklog(): Promise<ContentPost[]> {
	const response = await apiClient.get<ApiResponse<ContentPost[]>>(
		'/social-media-posts/content/backlog'
	);
	return response.data.data ?? [];
}
