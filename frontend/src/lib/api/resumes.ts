import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Resume {
	id: string;
	userId: string;
	title: string;
	isMain: boolean;
	isFeatured: boolean;
	filePath: string;
	fileName: string;
	fileSize: number;
	tags?: string[];
	applicationsUsed: number;
	interviewRate: number; // Percentage (0-100)
	offerRate: number;     // Percentage (0-100)
	createdAt: string;
	updatedAt: string;
}

export interface CreateResumeInput {
	title: string;
	filePath: string;
	fileName: string;
	fileSize: number;
	tags?: string[];
}

export interface UpdateResumeInput {
	title?: string;
	tags?: string[];
}

// List all resumes for the authenticated user
export async function listResumes(): Promise<Resume[]> {
	const response = await apiClient.get<ApiResponse<Resume[]>>('/resumes');
	return response.data.data || [];
}

// Get resume by ID
export async function getResume(id: string): Promise<Resume | null> {
	try {
		const response = await apiClient.get<ApiResponse<Resume>>(`/resumes/${id}`);
		return response.data.data || null;
	} catch (error) {
		console.error('Error fetching resume:', error);
		return null;
	}
}

// Create resume
export async function createResume(input: CreateResumeInput): Promise<Resume> {
	const response = await apiClient.post<ApiResponse<Resume>>('/resumes', input);
	return response.data.data;
}

// Upload resume file
export async function uploadResume(file: File, title: string): Promise<Resume> {
	const formData = new FormData();
	formData.append('file', file);
	formData.append('title', title);

	const response = await apiClient.post<ApiResponse<Resume>>('/resumes/upload', formData, {
		headers: {
			'Content-Type': 'multipart/form-data'
		}
	});
	return response.data.data;
}

// Update resume
export async function updateResume(id: string, input: UpdateResumeInput): Promise<Resume> {
	const response = await apiClient.patch<ApiResponse<Resume>>(`/resumes/${id}`, input);
	return response.data.data;
}

// Delete resume
export async function deleteResume(id: string): Promise<void> {
	await apiClient.delete(`/resumes/${id}`);
}

// Mark resume as main
export async function markAsMain(id: string): Promise<Resume> {
	const response = await apiClient.patch<ApiResponse<Resume>>(`/resumes/${id}/main`);
	return response.data.data;
}

// Mark resume as featured
export async function markAsFeatured(id: string): Promise<Resume> {
	const response = await apiClient.patch<ApiResponse<Resume>>(`/resumes/${id}/featured`);
	return response.data.data;
}

// Unmark resume as main
export async function unmarkAsMain(id: string): Promise<Resume> {
	const response = await apiClient.delete<ApiResponse<Resume>>(`/resumes/${id}/main`);
	return response.data.data;
}

// Unmark resume as featured
export async function unmarkAsFeatured(id: string): Promise<Resume> {
	const response = await apiClient.delete<ApiResponse<Resume>>(`/resumes/${id}/featured`);
	return response.data.data;
}

// Download resume (public endpoint)
export function getResumeDownloadUrl(userId: string, language: string = 'en'): string {
	const params = new URLSearchParams({ userId, language });
	return `${apiClient.defaults.baseURL}/public/resume/download?${params.toString()}`;
}

// Preview resume (public endpoint)
export function getResumePreviewUrl(userId: string, language: string = 'en'): string {
	const params = new URLSearchParams({ userId, language });
	return `${apiClient.defaults.baseURL}/public/resume/preview?${params.toString()}`;
}

// Request CV generation for a job application
export interface RequestCVGenerationInput {
	jobApplicationId: string;
	language?: string;
}

export async function requestCVGeneration(input: RequestCVGenerationInput): Promise<Resume> {
	const response = await apiClient.post<ApiResponse<Resume>>('/resumes/generate', input);
	return response.data.data;
}

// List all unique tags from resumes (for autocomplete)
export async function listResumeTags(): Promise<string[]> {
	const response = await apiClient.get<ApiResponse<string[]>>('/resumes/tags');
	return response.data.data || [];
}

// List resumes with optional tag filter
export async function listResumesByTags(tags: string[]): Promise<Resume[]> {
	const tagsParam = tags.join(',');
	const response = await apiClient.get<ApiResponse<Resume[]>>(`/resumes?tags=${encodeURIComponent(tagsParam)}`);
	return response.data.data || [];
}

// Recalculate metrics for a resume
export async function recalculateResumeMetrics(id: string): Promise<Resume> {
	const response = await apiClient.post<ApiResponse<Resume>>(`/resumes/${id}/recalculate-metrics`);
	return response.data.data;
}

