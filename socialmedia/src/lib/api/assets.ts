import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface ContentAsset {
	id: string;
	contentPostId?: string;
	socialPostId?: string;
	assetType: 'image' | 'video' | 'document' | 'other';
	filePath: string;
	fileUrl?: string;
	altText?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateAssetRequest {
	contentPostId?: string;
	socialPostId?: string;
	assetType: 'image' | 'video' | 'document' | 'other';
	filePath: string;
	fileUrl?: string;
	altText?: string;
}

export interface UpdateAssetRequest {
	fileUrl?: string;
	altText?: string;
}

export async function listAssets(params?: {
	contentPostId?: string;
	socialPostId?: string;
}): Promise<ContentAsset[]> {
	const queryParams = new URLSearchParams();
	if (params?.contentPostId) queryParams.append('contentPostId', params.contentPostId);
	if (params?.socialPostId) queryParams.append('socialPostId', params.socialPostId);

	const response = await apiClient.get<ApiResponse<ContentAsset[]>>(
		`/social-media-posts/assets?${queryParams.toString()}`
	);
	return response.data.data ?? [];
}

export async function getAsset(id: string): Promise<ContentAsset> {
	const response = await apiClient.get<ApiResponse<ContentAsset>>(`/social-media-posts/assets/${id}`);
	return response.data.data;
}

export async function getAssetsByContentPost(contentPostId: string): Promise<ContentAsset[]> {
	const response = await apiClient.get<ApiResponse<ContentAsset[]>>(
		`/social-media-posts/assets/content-posts/${contentPostId}`
	);
	return response.data.data ?? [];
}

export async function getAssetsBySocialPost(socialPostId: string): Promise<ContentAsset[]> {
	const response = await apiClient.get<ApiResponse<ContentAsset[]>>(
		`/social-media-posts/assets/social-posts/${socialPostId}`
	);
	return response.data.data ?? [];
}

export async function createAsset(input: CreateAssetRequest): Promise<ContentAsset> {
	const response = await apiClient.post<ApiResponse<ContentAsset>>(
		'/social-media-posts/assets',
		input
	);
	return response.data.data;
}

export async function updateAsset(id: string, input: UpdateAssetRequest): Promise<ContentAsset> {
	const response = await apiClient.patch<ApiResponse<ContentAsset>>(
		`/social-media-posts/assets/${id}`,
		input
	);
	return response.data.data;
}

export async function deleteAsset(id: string): Promise<void> {
	await apiClient.delete(`/social-media-posts/assets/${id}`);
}
