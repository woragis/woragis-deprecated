export interface Testimonial {
	id: string;
	userId: string;
	authorName: string;
	authorRole?: string;
	authorCompany?: string;
	authorPhoto?: string;
	content: string;
	rating?: number; // 1-5 stars
	linkedinUrl?: string;
	status: 'pending' | 'approved' | 'rejected' | 'hidden';
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
}

export interface ListTestimonialsParams {
	status?: 'pending' | 'approved' | 'rejected' | 'hidden';
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

