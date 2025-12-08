import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type ResponseType = 'rejection' | 'interview' | 'offer' | 'no-response';

export interface JobApplicationResponse {
	id: string;
	jobApplicationId: string;
	responseType: ResponseType;
	responseDate: string;
	message?: string;
	contactPerson?: string;
	contactEmail?: string;
	contactPhone?: string;
	responseChannel?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateResponseInput {
	jobApplicationId: string;
	responseType: ResponseType;
	responseDate?: string;
	message?: string;
	contactPerson?: string;
	contactEmail?: string;
	contactPhone?: string;
	responseChannel?: string;
}

export interface UpdateResponseInput {
	message?: string;
	contactPerson?: string;
	contactEmail?: string;
	contactPhone?: string;
	responseChannel?: string;
}

export async function listResponses(applicationId: string): Promise<JobApplicationResponse[]> {
	const response = await apiClient.get<ApiResponse<{ responses: JobApplicationResponse[]; count: number }>>(
		`/job-applications/${applicationId}/responses`
	);
	return response.data.data?.responses ?? [];
}

export async function getResponse(id: string, applicationId: string): Promise<JobApplicationResponse> {
	const response = await apiClient.get<ApiResponse<JobApplicationResponse>>(
		`/job-applications/${applicationId}/responses/${id}`
	);
	return response.data.data;
}

export async function createResponse(
	applicationId: string,
	input: Omit<CreateResponseInput, 'jobApplicationId'>
): Promise<JobApplicationResponse> {
	const payload: CreateResponseInput = {
		...input,
		jobApplicationId: applicationId
	};
	const response = await apiClient.post<ApiResponse<JobApplicationResponse>>(
		`/job-applications/${applicationId}/responses`,
		payload
	);
	return response.data.data;
}

export async function updateResponse(
	id: string,
	applicationId: string,
	input: UpdateResponseInput
): Promise<JobApplicationResponse> {
	const response = await apiClient.patch<ApiResponse<JobApplicationResponse>>(
		`/job-applications/${applicationId}/responses/${id}`,
		input
	);
	return response.data.data;
}

export async function deleteResponse(id: string, applicationId: string): Promise<void> {
	await apiClient.delete(`/job-applications/${applicationId}/responses/${id}`);
}

