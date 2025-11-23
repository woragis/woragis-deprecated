export interface Post {
	id: string;
	userId: string;
	title: string;
	slug: string;
	content: string; // Markdown content
	excerpt?: string;
	status: 'draft' | 'published' | 'archived';
	publishedAt?: string;
	featuredImage?: string;
	metaTitle?: string;
	metaDescription?: string;
	metaKeywords?: string;
	ogTitle?: string;
	ogDescription?: string;
	ogImage?: string;
	featured: boolean;
	viewsCount: number;
	createdAt: string;
	updatedAt: string;
	skills?: Array<{ id: string; name: string; slug: string }>;
	categories?: Array<{ id: string; name: string; slug: string }>;
	tags?: Array<{ id: string; name: string; slug: string }>;
}

export interface ListPostsParams {
	status?: 'draft' | 'published' | 'archived';
	featured?: boolean;
	categoryId?: string;
	tagId?: string;
	skillId?: string;
	search?: string;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

export interface Category {
	id: string;
	name: string;
	slug: string;
	description?: string;
	createdAt: string;
	updatedAt: string;
}

export interface Tag {
	id: string;
	name: string;
	slug: string;
	createdAt: string;
	updatedAt: string;
}

