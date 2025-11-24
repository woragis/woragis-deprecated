export type WritingType =
	| 'article'
	| 'documentation'
	| 'tutorial'
	| 'guide'
	| 'blog_post'
	| 'case_study'
	| 'other';

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

export interface ListTechnicalWritingsParams {
	type?: WritingType;
	platform?: PublicationPlatform;
	projectId?: string;
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

