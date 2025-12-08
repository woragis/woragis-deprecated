import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface ArchitectureData {
	diagram?: string;
	diagramType?: string;
	description?: string;
	components?: Array<{ name: string; description: string; technologies: string[] }>;
}

export interface MetricsData {
	before?: Array<{ label: string; value: string }>;
	after?: Array<{ label: string; value: string }>;
	impact?: string;
}

export interface CaseStudy {
	id: string;
	userId: string;
	projectId: string;
	projectSlug: string;
	title: string;
	problem: string;
	context: string;
	solution: string;
	approach: string[];
	architecture?: ArchitectureData;
	metrics?: MetricsData;
	lessonsLearned: string[];
	technologies: string[];
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateCaseStudyInput {
	projectId: string;
	projectSlug: string;
	title: string;
	problem: string;
	context: string;
	solution: string;
	approach?: string[];
	architecture?: ArchitectureData;
	metrics?: MetricsData;
	lessonsLearned?: string[];
	technologies?: string[];
	featured?: boolean;
}

export interface UpdateCaseStudyInput {
	title?: string;
	problem?: string;
	context?: string;
	solution?: string;
	approach?: string[];
	architecture?: ArchitectureData;
	metrics?: MetricsData;
	lessonsLearned?: string[];
	technologies?: string[];
	featured?: boolean;
}

export async function listCaseStudies(): Promise<CaseStudy[]> {
	const response = await apiClient.get<ApiResponse<CaseStudy[]>>('/case-studies');
	return response.data.data ?? [];
}

export async function getCaseStudy(id: string): Promise<CaseStudy> {
	const response = await apiClient.get<ApiResponse<CaseStudy>>(`/case-studies/${id}`);
	return response.data.data;
}

export async function createCaseStudy(input: CreateCaseStudyInput): Promise<CaseStudy> {
	const response = await apiClient.post<ApiResponse<CaseStudy>>('/case-studies', input);
	return response.data.data;
}

export async function updateCaseStudy(id: string, input: UpdateCaseStudyInput): Promise<CaseStudy> {
	const response = await apiClient.patch<ApiResponse<CaseStudy>>(`/case-studies/${id}`, input);
	return response.data.data;
}

export async function deleteCaseStudy(id: string): Promise<void> {
	await apiClient.delete(`/case-studies/${id}`);
}

