import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type WritingType = 'article' | 'documentation' | 'tutorial' | 'guide' | 'blog_post' | 'case_study' | 'other';

export type PublicationPlatform =
	| 'medium'
	| 'dev_to'
	| 'hashnode'
	| 'personal_blog'
	| 'github'
	| 'company_blog'
	| 'substack'
	| 'linkedin'
	| 'other';

export interface TechnicalWriting {
	id: string;
	userId: string;
	title: string;
	description: string;
	type: WritingType;
	platform: PublicationPlatform;
	content?: string;
	url: string;
	canonicalUrl?: string;
	publishedAt?: string;
	readingTime?: number;
	topics?: string[];
	technologies?: string[];
	views?: number;
	likes?: number;
	shares?: number;
	comments?: number;
	projectId?: string;
	caseStudyId?: string;
	featured: boolean;
	displayOrder: number;
	excerpt?: string;
	coverImageUrl?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateTechnicalWritingInput {
	title: string;
	description: string;
	type: WritingType;
	platform: PublicationPlatform;
	url: string;
	canonicalUrl?: string;
	content?: string;
	publishedAt?: string;
	readingTime?: number;
	topics?: string[];
	technologies?: string[];
	views?: number;
	likes?: number;
	shares?: number;
	comments?: number;
	projectId?: string;
	caseStudyId?: string;
	featured?: boolean;
	displayOrder?: number;
	excerpt?: string;
	coverImageUrl?: string;
}

export interface UpdateTechnicalWritingInput {
	title?: string;
	description?: string;
	type?: WritingType;
	platform?: PublicationPlatform;
	url?: string;
	canonicalUrl?: string;
	content?: string;
	publishedAt?: string;
	readingTime?: number;
	topics?: string[];
	technologies?: string[];
	views?: number;
	likes?: number;
	shares?: number;
	comments?: number;
	projectId?: string;
	caseStudyId?: string;
	featured?: boolean;
	displayOrder?: number;
	excerpt?: string;
	coverImageUrl?: string;
}

export async function listTechnicalWritings(): Promise<TechnicalWriting[]> {
	const response = await apiClient.get<ApiResponse<TechnicalWriting[]>>('/technical-writings');
	return response.data.data ?? [];
}

export async function getTechnicalWriting(id: string): Promise<TechnicalWriting> {
	const response = await apiClient.get<ApiResponse<TechnicalWriting>>(`/technical-writings/${id}`);
	return response.data.data;
}

export async function createTechnicalWriting(input: CreateTechnicalWritingInput): Promise<TechnicalWriting> {
	const response = await apiClient.post<ApiResponse<TechnicalWriting>>('/technical-writings', input);
	return response.data.data;
}

export async function updateTechnicalWriting(
	id: string,
	input: UpdateTechnicalWritingInput
): Promise<TechnicalWriting> {
	const response = await apiClient.patch<ApiResponse<TechnicalWriting>>(`/technical-writings/${id}`, input);
	return response.data.data;
}

export async function deleteTechnicalWriting(id: string): Promise<void> {
	await apiClient.delete(`/technical-writings/${id}`);
}

