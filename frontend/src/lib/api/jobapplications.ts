import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type ApplicationStatus =
	| 'pending'
	| 'processing'
	| 'applied'
	| 'contacted'
	| 'rejected'
	| 'accepted'
	| 'failed';

export interface JobApplication {
	id: string;
	userId: string;
	companyName: string;
	location?: string;
	jobTitle: string;
	jobUrl: string;
	website: string;
	appliedAt?: string;
	coverLetter?: string;
	linkedInContact: boolean;
	status: ApplicationStatus;
	errorMessage?: string;
	resumeId?: string;
	salaryMin?: number;
	salaryMax?: number;
	salaryCurrency?: string;
	jobDescription?: string;
	deadline?: string;
	interestLevel?: string;
	notes?: string;
	tags?: string[];
	followUpDate?: string;
	responseReceivedAt?: string;
	rejectionReason?: string;
	interviewCount: number;
	nextInterviewDate?: string;
	source?: string;
	applicationMethod?: string;
	language?: string; // ISO 639-1 language code (2 characters)
	createdAt: string;
	updatedAt: string;
}

export interface CreateJobApplicationInput {
	companyName: string;
	location?: string;
	jobTitle: string;
	jobUrl: string;
	website: string;
	coverLetter?: string;
	linkedInContact?: boolean;
	status?: ApplicationStatus;
}

export interface UpdateJobApplicationStatusInput {
	status: ApplicationStatus;
}

export interface UpdateJobApplicationInput {
	resumeId?: string;
	salaryMin?: number;
	salaryMax?: number;
	salaryCurrency?: string;
	jobDescription?: string;
	deadline?: string;
	interestLevel?: string;
	notes?: string;
	tags?: string[];
	followUpDate?: string;
	responseReceivedAt?: string;
	rejectionReason?: string;
	nextInterviewDate?: string;
	source?: string;
	applicationMethod?: string;
	language?: string; // ISO 639-1 language code (2 characters)
}

export async function listJobApplications(): Promise<JobApplication[]> {
	const response = await apiClient.get<ApiResponse<JobApplication[]>>('/job-applications');
	return response.data.data ?? [];
}

export async function getJobApplication(id: string): Promise<JobApplication> {
	const response = await apiClient.get<ApiResponse<JobApplication>>(`/job-applications/${id}`);
	return response.data.data;
}

export async function createJobApplication(input: CreateJobApplicationInput): Promise<JobApplication> {
	const response = await apiClient.post<ApiResponse<JobApplication>>('/job-applications', input);
	return response.data.data;
}

export async function updateJobApplicationStatus(
	id: string,
	input: UpdateJobApplicationStatusInput
): Promise<JobApplication> {
	const response = await apiClient.patch<ApiResponse<JobApplication>>(
		`/job-applications/${id}/status`,
		input
	);
	return response.data.data;
}

export async function updateJobApplication(
	id: string,
	input: UpdateJobApplicationInput
): Promise<JobApplication> {
	const response = await apiClient.patch<ApiResponse<JobApplication>>(
		`/job-applications/${id}`,
		input
	);
	return response.data.data;
}

export async function deleteJobApplication(id: string): Promise<void> {
	await apiClient.delete(`/job-applications/${id}`);
}

