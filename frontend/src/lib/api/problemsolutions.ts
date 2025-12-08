import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface MetricsData {
	before?: string;
	after?: string;
	improvement?: string;
}

export interface ProblemSolution {
	id: string;
	userId: string;
	problem: string;
	context: string;
	solution: string;
	technologies: string[];
	impact?: string;
	metrics?: MetricsData;
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateProblemSolutionInput {
	problem: string;
	context: string;
	solution: string;
	technologies?: string[];
	impact?: string;
	metrics?: MetricsData;
	featured?: boolean;
}

export interface UpdateProblemSolutionInput {
	problem?: string;
	context?: string;
	solution?: string;
	technologies?: string[];
	impact?: string;
	metrics?: MetricsData;
	featured?: boolean;
}

export async function listProblemSolutions(): Promise<ProblemSolution[]> {
	const response = await apiClient.get<ApiResponse<ProblemSolution[]>>('/problem-solutions');
	return response.data.data ?? [];
}

export async function getProblemSolution(id: string): Promise<ProblemSolution> {
	const response = await apiClient.get<ApiResponse<ProblemSolution>>(`/problem-solutions/${id}`);
	return response.data.data;
}

export async function createProblemSolution(input: CreateProblemSolutionInput): Promise<ProblemSolution> {
	const response = await apiClient.post<ApiResponse<ProblemSolution>>('/problem-solutions', input);
	return response.data.data;
}

export async function updateProblemSolution(
	id: string,
	input: UpdateProblemSolutionInput
): Promise<ProblemSolution> {
	const response = await apiClient.patch<ApiResponse<ProblemSolution>>(`/problem-solutions/${id}`, input);
	return response.data.data;
}

export async function deleteProblemSolution(id: string): Promise<void> {
	await apiClient.delete(`/problem-solutions/${id}`);
}

