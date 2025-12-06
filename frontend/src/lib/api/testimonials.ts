import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type TestimonialStatus = 'pending' | 'approved' | 'rejected' | 'hidden';

export type TestimonialType = 'general' | 'project_specific' | 'skill_specific';

export interface Testimonial {
	id: string;
	userId: string;
	authorName: string;
	authorRole?: string;
	authorCompany?: string;
	authorPhoto?: string;
	content: string;
	context?: string;
	videoUrl?: string;
	type: TestimonialType;
	rating?: number;
	linkedinUrl?: string;
	status: TestimonialStatus;
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
}

export interface CreateTestimonialInput {
	authorName: string;
	authorRole?: string;
	authorCompany?: string;
	authorPhoto?: string;
	content: string;
	context?: string;
	videoUrl?: string;
	type?: TestimonialType;
	rating?: number;
	linkedinUrl?: string;
	status?: TestimonialStatus;
	displayOrder?: number;
}

export interface UpdateTestimonialInput {
	authorName?: string;
	authorRole?: string;
	authorCompany?: string;
	authorPhoto?: string;
	content?: string;
	context?: string;
	videoUrl?: string;
	type?: TestimonialType;
	rating?: number;
	linkedinUrl?: string;
	status?: TestimonialStatus;
	displayOrder?: number;
}

export async function listTestimonials(): Promise<Testimonial[]> {
	const response = await apiClient.get<ApiResponse<Testimonial[]>>('/testimonials');
	return response.data.data ?? [];
}

export async function getTestimonial(id: string): Promise<Testimonial> {
	const response = await apiClient.get<ApiResponse<Testimonial>>(`/testimonials/${id}`);
	return response.data.data;
}

export async function createTestimonial(input: CreateTestimonialInput): Promise<Testimonial> {
	const response = await apiClient.post<ApiResponse<Testimonial>>('/testimonials', input);
	return response.data.data;
}

export async function updateTestimonial(
	id: string,
	input: UpdateTestimonialInput
): Promise<Testimonial> {
	const response = await apiClient.patch<ApiResponse<Testimonial>>(`/testimonials/${id}`, input);
	return response.data.data;
}

export async function deleteTestimonial(id: string): Promise<void> {
	await apiClient.delete(`/testimonials/${id}`);
}

export async function approveTestimonial(id: string): Promise<void> {
	await apiClient.post(`/testimonials/${id}/approve`);
}

export async function rejectTestimonial(id: string): Promise<void> {
	await apiClient.post(`/testimonials/${id}/reject`);
}

export async function hideTestimonial(id: string): Promise<void> {
	await apiClient.post(`/testimonials/${id}/hide`);
}

