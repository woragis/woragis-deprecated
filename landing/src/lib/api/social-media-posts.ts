import { apiClient } from './client';
import type { SocialMediaPost, ListSocialMediaPostsParams } from '$lib/types/social-media-post';

// List all social media posts
export async function listSocialMediaPosts(
	params?: ListSocialMediaPostsParams
): Promise<SocialMediaPost[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.platform) {
			queryParams.append('platform', params.platform);
		}
		if (params?.status) {
			queryParams.append('status', params.status);
		}
		if (params?.limit) {
			queryParams.append('limit', params.limit.toString());
		}
		if (params?.offset) {
			queryParams.append('offset', params.offset.toString());
		}

		const queryString = queryParams.toString();
		const url = queryString ? `/social-media-posts?${queryString}` : '/social-media-posts';

		const response = await apiClient.get<{ success: boolean; data: SocialMediaPost[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching social media posts:', error);
		throw error;
	}
}

// Get social media post by ID
export async function getSocialMediaPost(id: string): Promise<SocialMediaPost | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: SocialMediaPost }>(
			`/social-media-posts/${id}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching social media post ${id}:`, error);
		return null;
	}
}

// Get social media post by URL
export async function getSocialMediaPostByURL(url: string): Promise<SocialMediaPost | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: SocialMediaPost }>(
			`/social-media-posts/by-url?url=${encodeURIComponent(url)}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching social media post by URL:`, error);
		return null;
	}
}

