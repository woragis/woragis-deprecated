import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type StageType =
	| 'phone-screen'
	| 'technical'
	| 'behavioral'
	| 'system-design'
	| 'final'
	| 'hr'
	| 'manager'
	| 'panel'
	| 'other';

export type StageOutcome = 'pending' | 'passed' | 'failed' | 'cancelled';

export interface InterviewStage {
	id: string;
	jobApplicationId: string;
	stageType: StageType;
	scheduledDate?: string;
	completedDate?: string;
	interviewerName?: string;
	interviewerEmail?: string;
	location?: string;
	notes?: string;
	feedback?: string;
	outcome: StageOutcome;
	createdAt: string;
	updatedAt: string;
}

export interface CreateStageInput {
	jobApplicationId: string;
	stageType: StageType;
	scheduledDate?: string;
	interviewerName?: string;
	interviewerEmail?: string;
	location?: string;
	notes?: string;
}

export interface UpdateStageInput {
	scheduledDate?: string;
	interviewerName?: string;
	interviewerEmail?: string;
	location?: string;
	notes?: string;
	feedback?: string;
}

export interface ScheduleStageInput {
	scheduledDate: string;
}

export interface CompleteStageInput {
	completedDate: string;
	outcome: StageOutcome;
}

export async function listStages(applicationId: string): Promise<InterviewStage[]> {
	const response = await apiClient.get<ApiResponse<{ stages: InterviewStage[]; count: number }>>(
		`/job-applications/${applicationId}/interview-stages`
	);
	return response.data.data?.stages ?? [];
}

export async function getStage(id: string, applicationId: string): Promise<InterviewStage> {
	const response = await apiClient.get<ApiResponse<InterviewStage>>(
		`/job-applications/${applicationId}/interview-stages/${id}`
	);
	return response.data.data;
}

export async function createStage(
	applicationId: string,
	input: Omit<CreateStageInput, 'jobApplicationId'>
): Promise<InterviewStage> {
	const payload: CreateStageInput = {
		...input,
		jobApplicationId: applicationId
	};
	const response = await apiClient.post<ApiResponse<InterviewStage>>(
		`/job-applications/${applicationId}/interview-stages`,
		payload
	);
	return response.data.data;
}

export async function updateStage(
	id: string,
	applicationId: string,
	input: UpdateStageInput
): Promise<InterviewStage> {
	const response = await apiClient.patch<ApiResponse<InterviewStage>>(
		`/job-applications/${applicationId}/interview-stages/${id}`,
		input
	);
	return response.data.data;
}

export async function scheduleStage(
	id: string,
	applicationId: string,
	input: ScheduleStageInput
): Promise<InterviewStage> {
	const response = await apiClient.post<ApiResponse<InterviewStage>>(
		`/job-applications/${applicationId}/interview-stages/${id}/schedule`,
		input
	);
	return response.data.data;
}

export async function completeStage(
	id: string,
	applicationId: string,
	input: CompleteStageInput
): Promise<InterviewStage> {
	const response = await apiClient.post<ApiResponse<InterviewStage>>(
		`/job-applications/${applicationId}/interview-stages/${id}/complete`,
		input
	);
	return response.data.data;
}

export async function deleteStage(id: string, applicationId: string): Promise<void> {
	await apiClient.delete(`/job-applications/${applicationId}/interview-stages/${id}`);
}

