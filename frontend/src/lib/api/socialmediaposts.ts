import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type Platform = 'linkedin' | 'twitter' | 'instagram';

export type PostStatus = 'active' | 'deleted' | 'unavailable';

export interface SocialMediaPost {
	id: string;
	url: string;
	platform: Platform;
	title?: string;
	contentPreview?: string;
	publishedDate?: string;
	likes?: number;
	shares?: number;
	comments?: number;
	views?: number;
	status: PostStatus;
	createdAt: string;
	updatedAt: string;
}

export interface CreateSocialMediaPostInput {
	url: string;
	platform: Platform;
	title?: string;
	contentPreview?: string;
	publishedDate?: string;
	likes?: number;
	shares?: number;
	comments?: number;
	views?: number;
	status?: PostStatus;
}

export interface UpdateSocialMediaPostInput {
	url?: string;
	platform?: Platform;
	title?: string;
	contentPreview?: string;
	publishedDate?: string;
	likes?: number;
	shares?: number;
	comments?: number;
	views?: number;
	status?: PostStatus;
}

export async function listSocialMediaPosts(): Promise<SocialMediaPost[]> {
	const response = await apiClient.get<ApiResponse<SocialMediaPost[]>>('/social-media-posts');
	return response.data.data ?? [];
}

export async function getSocialMediaPost(id: string): Promise<SocialMediaPost> {
	const response = await apiClient.get<ApiResponse<SocialMediaPost>>(`/social-media-posts/${id}`);
	return response.data.data;
}

export async function createSocialMediaPost(input: CreateSocialMediaPostInput): Promise<SocialMediaPost> {
	const response = await apiClient.post<ApiResponse<SocialMediaPost>>('/social-media-posts', input);
	return response.data.data;
}

export async function updateSocialMediaPost(
	id: string,
	input: UpdateSocialMediaPostInput
): Promise<SocialMediaPost> {
	const response = await apiClient.patch<ApiResponse<SocialMediaPost>>(`/social-media-posts/${id}`, input);
	return response.data.data;
}

export async function deleteSocialMediaPost(id: string): Promise<void> {
	await apiClient.delete(`/social-media-posts/${id}`);
}

