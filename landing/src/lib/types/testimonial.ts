export type TestimonialType = 'general' | 'project_specific' | 'skill_specific';

export type EntityType = 'project' | 'skill';

export interface TestimonialEntityLink {
	id: string;
	testimonialId: string;
	entityType: EntityType;
	entityId: string;
	createdAt: string;
	updatedAt: string;
}

export interface Testimonial {
	id: string;
	userId: string;
	authorName: string;
	authorRole?: string;
	authorCompany?: string;
	authorPhoto?: string;
	content: string;
	context?: string; // When/where/why the testimonial was given
	videoUrl?: string; // Optional video testimonial URL
	type: TestimonialType;
	rating?: number; // 1-5 stars
	linkedinUrl?: string;
	status: 'pending' | 'approved' | 'rejected' | 'hidden';
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
	entityLinks?: TestimonialEntityLink[]; // Links to projects or skills
}

export interface ListTestimonialsParams {
	status?: 'pending' | 'approved' | 'rejected' | 'hidden';
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

