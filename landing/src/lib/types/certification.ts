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
	skills?: Array<{ id: string; name: string }>;
	createdAt: string;
	updatedAt: string;
}

export interface ListCertificationsParams {
	status?: CertificationStatus;
	category?: CertificationCategory;
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

