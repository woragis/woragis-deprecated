import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface MediaFile {
	id: string;
	filename: string;
	original_filename: string;
	mime_type: string;
	size: number;
	url: string;
	thumbnail_url?: string;
	width?: number;
	height?: number;
	alt_text?: string;
	caption?: string;
	created_at: string;
	updated_at: string;
}

export interface UploadMediaResponse {
	file: MediaFile;
	message: string;
}

export interface MediaListResponse {
	files: MediaFile[];
	total: number;
	page: number;
	limit: number;
}

export interface UploadMediaInput {
	file: File;
	alt_text?: string;
	caption?: string;
	folder?: string;
}

// Upload a media file
export async function uploadMedia(input: UploadMediaInput): Promise<MediaFile> {
	const formData = new FormData();
	formData.append('file', input.file);
	if (input.alt_text) formData.append('alt_text', input.alt_text);
	if (input.caption) formData.append('caption', input.caption);
	if (input.folder) formData.append('folder', input.folder);

	const response = await apiClient.post<ApiResponse<UploadMediaResponse>>('/media/upload', formData, {
		headers: {
			'Content-Type': 'multipart/form-data'
		}
	});

	return response.data.data.file;
}

// List media files
export async function listMedia(params?: {
	page?: number;
	limit?: number;
	search?: string;
	mime_type?: string;
	folder?: string;
}): Promise<MediaListResponse> {
	const response = await apiClient.get<ApiResponse<MediaListResponse>>('/media', { params });
	return response.data.data;
}

// Get a single media file
export async function getMedia(id: string): Promise<MediaFile> {
	const response = await apiClient.get<ApiResponse<MediaFile>>(`/media/${id}`);
	return response.data.data;
}

// Update media metadata
export async function updateMedia(
	id: string,
	updates: {
		alt_text?: string;
		caption?: string;
		filename?: string;
	}
): Promise<MediaFile> {
	const response = await apiClient.patch<ApiResponse<MediaFile>>(`/media/${id}`, updates);
	return response.data.data;
}

// Delete a media file
export async function deleteMedia(id: string): Promise<void> {
	await apiClient.delete(`/media/${id}`);
}

// Bulk delete media files
export async function bulkDeleteMedia(ids: string[]): Promise<void> {
	await apiClient.post('/media/bulk-delete', { ids });
}

// Get media URL (helper for constructing URLs)
export function getMediaUrl(id: string): string {
	return `/media/${id}`;
}

// Format file size
export function formatFileSize(bytes: number): string {
	if (bytes === 0) return '0 Bytes';
	const k = 1024;
	const sizes = ['Bytes', 'KB', 'MB', 'GB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

// Check if file is an image
export function isImage(mimeType: string): boolean {
	return mimeType.startsWith('image/');
}

// Check if file is a video
export function isVideo(mimeType: string): boolean {
	return mimeType.startsWith('video/');
}

