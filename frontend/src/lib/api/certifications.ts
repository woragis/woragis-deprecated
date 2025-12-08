import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type CertificationStatus = 'active' | 'expired' | 'revoked' | 'renewed';

export type CertificationCategory =
	| 'cloud'
	| 'security'
	| 'programming'
	| 'database'
	| 'devops'
	| 'architecture'
	| 'other';

export interface Certification {
	id: string;
	userId: string;
	name: string;
	issuer: string;
	issueDate: string;
	expiryDate?: string;
	credentialId?: string;
	verificationUrl?: string;
	certificateUrl?: string;
	description?: string;
	status: CertificationStatus;
	category: CertificationCategory;
	featured: boolean;
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
	skills?: Array<{ id: string; name: string }>;
}

export interface CreateCertificationInput {
	name: string;
	issuer: string;
	issueDate: string;
	expiryDate?: string;
	credentialId?: string;
	verificationUrl?: string;
	certificateUrl?: string;
	description?: string;
	status?: CertificationStatus;
	category?: CertificationCategory;
	featured?: boolean;
	displayOrder?: number;
}

export interface UpdateCertificationInput {
	name?: string;
	issuer?: string;
	issueDate?: string;
	expiryDate?: string;
	credentialId?: string;
	verificationUrl?: string;
	certificateUrl?: string;
	description?: string;
	status?: CertificationStatus;
	category?: CertificationCategory;
	featured?: boolean;
	displayOrder?: number;
}

export async function listCertifications(): Promise<Certification[]> {
	const response = await apiClient.get<ApiResponse<Certification[]>>('/certifications');
	return response.data.data ?? [];
}

export async function getCertification(id: string): Promise<Certification> {
	const response = await apiClient.get<ApiResponse<Certification>>(`/certifications/${id}`);
	return response.data.data;
}

export async function createCertification(input: CreateCertificationInput): Promise<Certification> {
	const response = await apiClient.post<ApiResponse<Certification>>('/certifications', input);
	return response.data.data;
}

export async function updateCertification(
	id: string,
	input: UpdateCertificationInput
): Promise<Certification> {
	const response = await apiClient.patch<ApiResponse<Certification>>(`/certifications/${id}`, input);
	return response.data.data;
}

export async function deleteCertification(id: string): Promise<void> {
	await apiClient.delete(`/certifications/${id}`);
}

export async function addCertificationSkill(certificationId: string, skillId: string): Promise<void> {
	await apiClient.post(`/certifications/${certificationId}/skills/${skillId}`);
}

export async function removeCertificationSkill(certificationId: string, skillId: string): Promise<void> {
	await apiClient.delete(`/certifications/${certificationId}/skills/${skillId}`);
}

export async function getCertificationSkills(certificationId: string): Promise<Array<{ id: string; name: string }>> {
	const response = await apiClient.get<ApiResponse<Array<{ id: string; name: string }>>>(
		`/certifications/${certificationId}/skills`
	);
	return response.data.data ?? [];
}

