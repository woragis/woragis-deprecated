import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type Platform = 'linkedin' | 'twitter' | 'instagram' | 'medium' | 'substack' | 'valete' | 'website';
export type ContentFormat = 'long-form' | 'thread' | 'carousel' | 'article' | 'newsletter' | 'post';
export type PostStatus = 'draft' | 'ready' | 'scheduled' | 'posted' | 'analyzed' | 'archived';

export interface SocialMediaPost {
	id: string;
	contentPostId?: string;
	platform: Platform;
	format: ContentFormat;
	status: PostStatus;
	title: string;
	content: string;
	wordCount: number;
	imageCount: number;
	scheduledAt?: string;
	postedAt?: string;
	analyzedAt?: string;
	url?: string;
	platformPostId?: string;
	likes?: number;
	shares?: number;
	comments?: number;
	views?: number;
	createdAt: string;
	updatedAt: string;
}

export interface PostFilters {
	platform?: Platform;
	status?: PostStatus;
}

// Import types from other modules to avoid duplication
import type { ScheduledPost } from './scheduling';
import type { ContentPost } from './content';
import type { AnalyticsSummary } from './analytics';

export interface DashboardData {
	upcomingPosts: ScheduledPost[];
	contentBacklog: ContentPost[];
	analyticsSummary: AnalyticsSummary;
	recentPosts: SocialMediaPost[];
}

export async function listSocialMediaPosts(filters?: PostFilters): Promise<SocialMediaPost[]> {
	const params = new URLSearchParams();
	if (filters?.platform) params.append('platform', filters.platform);
	if (filters?.status) params.append('status', filters.status);

	const response = await apiClient.get<ApiResponse<SocialMediaPost[]>>(
		`/social-media-posts?${params.toString()}`
	);
	return response.data.data ?? [];
}

export async function getSocialMediaPost(id: string): Promise<SocialMediaPost> {
	const response = await apiClient.get<ApiResponse<SocialMediaPost>>(`/social-media-posts/${id}`);
	return response.data.data;
}

export async function createSocialMediaPost(input: {
	platform: Platform;
	format: ContentFormat;
	title: string;
	content: string;
	contentPostId?: string;
}): Promise<SocialMediaPost> {
	const response = await apiClient.post<ApiResponse<SocialMediaPost>>('/social-media-posts', input);
	return response.data.data;
}

export async function updateSocialMediaPost(
	id: string,
	input: {
		title?: string;
		content?: string;
		status?: PostStatus;
	}
): Promise<SocialMediaPost> {
	const response = await apiClient.patch<ApiResponse<SocialMediaPost>>(
		`/social-media-posts/${id}`,
		input
	);
	return response.data.data;
}

export async function updateSocialMediaPostStatus(
	id: string,
	status: PostStatus
): Promise<SocialMediaPost> {
	const response = await apiClient.patch<ApiResponse<SocialMediaPost>>(
		`/social-media-posts/${id}/status`,
		{ status }
	);
	return response.data.data;
}

export async function deleteSocialMediaPost(id: string): Promise<void> {
	await apiClient.delete(`/social-media-posts/${id}`);
}

// Re-export functions from other modules for convenience
export { getUpcomingScheduledPosts } from './scheduling';
export { getContentBacklog } from './content';
export { getAnalyticsSummary } from './analytics';

import { getUpcomingScheduledPosts as fetchUpcomingScheduledPosts } from './scheduling';
import { getContentBacklog as fetchContentBacklog } from './content';
import { getAnalyticsSummary as fetchAnalyticsSummary } from './analytics';

export async function getDashboardData(): Promise<DashboardData> {
	// Fetch all dashboard data in parallel
	const [upcomingPosts, contentBacklog, analyticsSummary, recentPosts] = await Promise.all([
		fetchUpcomingScheduledPosts(5),
		fetchContentBacklog(),
		fetchAnalyticsSummary(),
		listSocialMediaPosts({ status: 'posted' }).then((posts) => posts.slice(0, 5))
	]);

	return {
		upcomingPosts,
		contentBacklog,
		analyticsSummary,
		recentPosts
	};
}
