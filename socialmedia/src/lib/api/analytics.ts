import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface PostAnalytics {
	id: string;
	socialPostId: string;
	metricDate: string;
	likes?: number;
	comments?: number;
	shares?: number;
	views?: number;
	clicks?: number;
	engagementRate?: number;
	reach?: number;
	impressions?: number;
	createdAt: string;
	updatedAt: string;
}

export interface AnalyticsSummary {
	totalLikes: number;
	totalComments: number;
	totalShares: number;
	totalViews: number;
	totalClicks: number;
	totalReach: number;
	totalImpressions: number;
	averageEngagementRate: number;
	postCount: number;
}

export interface TopPost {
	socialPostId: string;
	metricValue: number;
	metricName: string;
}

export interface RecordAnalyticsRequest {
	socialPostId: string;
	metricDate: string;
	likes?: number;
	comments?: number;
	shares?: number;
	views?: number;
	clicks?: number;
	reach?: number;
	impressions?: number;
}

export async function getPostAnalytics(
	postId: string,
	params?: { startDate?: string; endDate?: string }
): Promise<PostAnalytics[]> {
	const queryParams = new URLSearchParams();
	if (params?.startDate) queryParams.append('startDate', params.startDate);
	if (params?.endDate) queryParams.append('endDate', params.endDate);

	const response = await apiClient.get<ApiResponse<PostAnalytics[]>>(
		`/social-media-posts/analytics/posts/${postId}?${queryParams.toString()}`
	);
	return response.data.data ?? [];
}

export async function getAnalyticsSummary(params?: {
	socialPostId?: string;
	startDate?: string;
	endDate?: string;
}): Promise<AnalyticsSummary> {
	const queryParams = new URLSearchParams();
	if (params?.socialPostId) queryParams.append('socialPostId', params.socialPostId);
	if (params?.startDate) queryParams.append('startDate', params.startDate);
	if (params?.endDate) queryParams.append('endDate', params.endDate);

	const response = await apiClient.get<ApiResponse<AnalyticsSummary>>(
		`/social-media-posts/analytics/summary?${queryParams.toString()}`
	);
	return response.data.data;
}

export async function getTopPosts(params?: {
	limit?: number;
	metric?: string;
}): Promise<TopPost[]> {
	const queryParams = new URLSearchParams();
	if (params?.limit) queryParams.append('limit', params.limit.toString());
	if (params?.metric) queryParams.append('metric', params.metric);

	const response = await apiClient.get<ApiResponse<TopPost[]>>(
		`/social-media-posts/analytics/top-posts?${queryParams.toString()}`
	);
	return response.data.data ?? [];
}

export async function recordAnalytics(input: RecordAnalyticsRequest): Promise<PostAnalytics> {
	const response = await apiClient.post<ApiResponse<PostAnalytics>>(
		'/social-media-posts/analytics',
		input
	);
	return response.data.data;
}
