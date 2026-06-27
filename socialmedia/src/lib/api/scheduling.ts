import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface ScheduledPost {
	id: string;
	socialPostId: string;
	platformId?: string;
	scheduledDate: string;
	scheduledTime: string;
	scheduledAt: string;
	status: 'pending' | 'scheduled' | 'posted' | 'cancelled';
	createdAt: string;
	updatedAt: string;
}

export interface SchedulePostRequest {
	socialPostId: string;
	scheduledAt: string;
	platformId?: string;
}

export interface UpdateScheduleRequest {
	scheduledAt?: string;
	status?: 'pending' | 'scheduled' | 'posted' | 'cancelled';
}

export async function listScheduledPosts(params?: {
	startDate?: string;
	endDate?: string;
}): Promise<ScheduledPost[]> {
	const queryParams = new URLSearchParams();
	if (params?.startDate) queryParams.append('startDate', params.startDate);
	if (params?.endDate) queryParams.append('endDate', params.endDate);

	const response = await apiClient.get<ApiResponse<ScheduledPost[]>>(
		`/social-media-posts/scheduling?${queryParams.toString()}`
	);
	return response.data.data ?? [];
}

export async function getUpcomingScheduledPosts(limit: number = 10): Promise<ScheduledPost[]> {
	const response = await apiClient.get<ApiResponse<ScheduledPost[]>>(
		`/social-media-posts/scheduling/upcoming?limit=${limit}`
	);
	return response.data.data ?? [];
}

export async function getScheduledPost(id: string): Promise<ScheduledPost> {
	const response = await apiClient.get<ApiResponse<ScheduledPost>>(
		`/social-media-posts/scheduling/${id}`
	);
	return response.data.data;
}

export async function schedulePost(input: SchedulePostRequest): Promise<ScheduledPost> {
	const response = await apiClient.post<ApiResponse<ScheduledPost>>(
		'/social-media-posts/scheduling',
		input
	);
	return response.data.data;
}

export async function updateSchedule(
	id: string,
	input: UpdateScheduleRequest
): Promise<ScheduledPost> {
	const response = await apiClient.patch<ApiResponse<ScheduledPost>>(
		`/social-media-posts/scheduling/${id}`,
		input
	);
	return response.data.data;
}

export async function cancelSchedule(id: string): Promise<void> {
	await apiClient.delete(`/social-media-posts/scheduling/${id}`);
}

export async function autoSchedule(
	socialPostId: string,
	platform: string
): Promise<ScheduledPost> {
	const response = await apiClient.post<ApiResponse<ScheduledPost>>(
		`/social-media-posts/scheduling/${socialPostId}/auto`,
		{ platform }
	);
	return response.data.data;
}
