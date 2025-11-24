export type Platform = 'linkedin' | 'twitter' | 'instagram';
export type PostStatus = 'active' | 'deleted' | 'unavailable';

export interface SocialMediaPost {
	id: string;
	url: string;
	platform: Platform;
	title?: string;
	contentPreview?: string;
	publishedDate?: string;
	likes?: number;
	shares?: number;
	comments?: number;
	views?: number;
	status: PostStatus;
	createdAt: string;
	updatedAt: string;
}

export interface ListSocialMediaPostsParams {
	platform?: Platform;
	status?: PostStatus;
	limit?: number;
	offset?: number;
}

