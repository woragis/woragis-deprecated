import { apiClient } from '$lib/clients/apiClient';

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

export interface JobApplicationResume {
	id: string;
	userId: string;
	title: string;
	isMain: boolean;
	isFeatured: boolean;
	filePath: string;
	fileName: string;
	fileSize: number;
	tags?: string[];
	createdAt: string;
	updatedAt: string;
}

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
	resume?: JobApplicationResume; // Full resume object when available from detail view
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
	interestLevel?: string;
	tags?: string[];
	followUpDate?: string;
	notes?: string;
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
	const response = await apiClient.get<ApiResponse<{ count: number; applications: JobApplication[] }>>('/job-applications');
	return response.data.data?.applications ?? [];
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

export interface GenerateCoverLetterInput {
	messageId?: string; // Optional: message ID from chat to use as additional context
}

export async function generateCoverLetter(
	id: string,
	input?: GenerateCoverLetterInput
): Promise<JobApplication> {
	const response = await apiClient.post<ApiResponse<JobApplication>>(
		`/job-applications/${id}/generate-cover-letter`,
		input || {}
	);
	return response.data.data;
}

