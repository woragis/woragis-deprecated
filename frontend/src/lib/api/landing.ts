import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

// ============================================================================
// POSTS
// ============================================================================

export interface Post {
	id: string;
	title: string;
	slug: string;
	content: string;
	excerpt?: string;
	status: string;
	featured_image?: string;
	meta_title?: string;
	meta_description?: string;
	meta_keywords?: string;
	og_title?: string;
	og_description?: string;
	og_image?: string;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreatePostInput {
	title: string;
	content: string;
	excerpt?: string;
	status?: string;
	featured_image?: string;
	meta_title?: string;
	meta_description?: string;
	meta_keywords?: string;
	og_title?: string;
	og_description?: string;
	og_image?: string;
	featured?: boolean;
	skill_ids?: string[];
	category_ids?: string[];
	tag_names?: string[];
}

export async function listPosts(): Promise<Post[]> {
	const response = await apiClient.get<ApiResponse<Post[]>>('/posts');
	return response.data.data ?? [];
}

export async function getPost(id: string): Promise<Post> {
	const response = await apiClient.get<ApiResponse<Post>>(`/posts/${id}`);
	return response.data.data;
}

export async function getPostBySlug(slug: string): Promise<Post> {
	const response = await apiClient.get<ApiResponse<Post>>(`/posts/slug/${slug}`);
	return response.data.data;
}

export async function createPost(input: CreatePostInput): Promise<Post> {
	const response = await apiClient.post<ApiResponse<Post>>('/posts', input);
	return response.data.data;
}

export async function updatePost(id: string, input: Partial<CreatePostInput>): Promise<Post> {
	const response = await apiClient.patch<ApiResponse<Post>>(`/posts/${id}`, input);
	return response.data.data;
}

export async function deletePost(id: string): Promise<void> {
	await apiClient.delete(`/posts/${id}`);
}

// ============================================================================
// TECHNICAL WRITINGS
// ============================================================================

export interface TechnicalWriting {
	id: string;
	title: string;
	slug: string;
	content: string;
	excerpt?: string;
	type: string;
	platform: string;
	url?: string;
	featured: boolean;
	published_at?: string;
	created_at: string;
	updated_at: string;
}

export interface CreateTechnicalWritingInput {
	title: string;
	content: string;
	excerpt?: string;
	type?: string;
	platform?: string;
	url?: string;
	featured?: boolean;
	published_at?: string;
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
	input: Partial<CreateTechnicalWritingInput>
): Promise<TechnicalWriting> {
	const response = await apiClient.patch<ApiResponse<TechnicalWriting>>(`/technical-writings/${id}`, input);
	return response.data.data;
}

export async function deleteTechnicalWriting(id: string): Promise<void> {
	await apiClient.delete(`/technical-writings/${id}`);
}

// ============================================================================
// CASE STUDIES
// ============================================================================

export interface CaseStudy {
	id: string;
	title: string;
	slug: string;
	content: string;
	excerpt?: string;
	project_id?: string;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateCaseStudyInput {
	title: string;
	content: string;
	excerpt?: string;
	project_id?: string;
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

export async function updateCaseStudy(id: string, input: Partial<CreateCaseStudyInput>): Promise<CaseStudy> {
	const response = await apiClient.patch<ApiResponse<CaseStudy>>(`/case-studies/${id}`, input);
	return response.data.data;
}

export async function deleteCaseStudy(id: string): Promise<void> {
	await apiClient.delete(`/case-studies/${id}`);
}

// ============================================================================
// PROBLEM SOLUTIONS
// ============================================================================

export interface ProblemSolution {
	id: string;
	problem: string;
	context?: string;
	solution: string;
	technologies?: string[];
	impact?: string;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateProblemSolutionInput {
	problem: string;
	context?: string;
	solution: string;
	technologies?: string[];
	impact?: string;
	featured?: boolean;
}

export async function listProblemSolutions(): Promise<ProblemSolution[]> {
	const response = await apiClient.get<ApiResponse<ProblemSolution[]>>('/problem-solutions');
	return response.data.data ?? [];
}

export async function getProblemSolution(id: string): Promise<ProblemSolution> {
	const response = await apiClient.get<ApiResponse<ProblemSolution>>(`/problem-solutions/${id}`);
	return response.data.data;
}

export async function createProblemSolution(input: CreateProblemSolutionInput): Promise<ProblemSolution> {
	const response = await apiClient.post<ApiResponse<ProblemSolution>>('/problem-solutions', input);
	return response.data.data;
}

export async function updateProblemSolution(
	id: string,
	input: Partial<CreateProblemSolutionInput>
): Promise<ProblemSolution> {
	const response = await apiClient.patch<ApiResponse<ProblemSolution>>(`/problem-solutions/${id}`, input);
	return response.data.data;
}

export async function deleteProblemSolution(id: string): Promise<void> {
	await apiClient.delete(`/problem-solutions/${id}`);
}

// ============================================================================
// SKILLS
// ============================================================================

export interface Skill {
	id: string;
	name: string;
	slug: string;
	description?: string;
	category?: string;
	proficiency_level?: string;
	years_of_experience?: number;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateSkillInput {
	name: string;
	description?: string;
	category?: string;
	proficiency_level?: string;
	years_of_experience?: number;
	featured?: boolean;
}

export async function listSkills(): Promise<Skill[]> {
	const response = await apiClient.get<ApiResponse<Skill[]>>('/skills');
	return response.data.data ?? [];
}

export async function getSkill(id: string): Promise<Skill> {
	const response = await apiClient.get<ApiResponse<Skill>>(`/skills/${id}`);
	return response.data.data;
}

export async function getSkillBySlug(slug: string): Promise<Skill> {
	const response = await apiClient.get<ApiResponse<Skill>>(`/skills/slug/${slug}`);
	return response.data.data;
}

export async function createSkill(input: CreateSkillInput): Promise<Skill> {
	const response = await apiClient.post<ApiResponse<Skill>>('/skills', input);
	return response.data.data;
}

export async function updateSkill(id: string, input: Partial<CreateSkillInput>): Promise<Skill> {
	const response = await apiClient.patch<ApiResponse<Skill>>(`/skills/${id}`, input);
	return response.data.data;
}

export async function deleteSkill(id: string): Promise<void> {
	await apiClient.delete(`/skills/${id}`);
}

// ============================================================================
// SOCIAL MEDIA POSTS
// ============================================================================

export interface SocialMediaPost {
	id: string;
	content: string;
	platform: string;
	url?: string;
	published_at?: string;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateSocialMediaPostInput {
	content: string;
	platform: string;
	url?: string;
	published_at?: string;
	featured?: boolean;
}

export async function listSocialMediaPosts(): Promise<SocialMediaPost[]> {
	const response = await apiClient.get<ApiResponse<SocialMediaPost[]>>('/social-media-posts');
	return response.data.data ?? [];
}

export async function getSocialMediaPost(id: string): Promise<SocialMediaPost> {
	const response = await apiClient.get<ApiResponse<SocialMediaPost>>(`/social-media-posts/${id}`);
	return response.data.data;
}

export async function createSocialMediaPost(input: CreateSocialMediaPostInput): Promise<SocialMediaPost> {
	const response = await apiClient.post<ApiResponse<SocialMediaPost>>('/social-media-posts', input);
	return response.data.data;
}

export async function updateSocialMediaPost(
	id: string,
	input: Partial<CreateSocialMediaPostInput>
): Promise<SocialMediaPost> {
	const response = await apiClient.patch<ApiResponse<SocialMediaPost>>(`/social-media-posts/${id}`, input);
	return response.data.data;
}

export async function deleteSocialMediaPost(id: string): Promise<void> {
	await apiClient.delete(`/social-media-posts/${id}`);
}

// ============================================================================
// CERTIFICATIONS
// ============================================================================

export interface Certification {
	id: string;
	name: string;
	slug: string;
	issuer: string;
	issue_date: string;
	expiry_date?: string;
	credential_id?: string;
	verification_url?: string;
	certificate_url?: string;
	description?: string;
	status: string;
	category?: string;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateCertificationInput {
	name: string;
	issuer: string;
	issue_date: string;
	expiry_date?: string;
	credential_id?: string;
	verification_url?: string;
	certificate_url?: string;
	description?: string;
	status?: string;
	category?: string;
	featured?: boolean;
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
	input: Partial<CreateCertificationInput>
): Promise<Certification> {
	const response = await apiClient.patch<ApiResponse<Certification>>(`/certifications/${id}`, input);
	return response.data.data;
}

export async function deleteCertification(id: string): Promise<void> {
	await apiClient.delete(`/certifications/${id}`);
}

// ============================================================================
// TESTIMONIALS
// ============================================================================

export interface Testimonial {
	id: string;
	content: string;
	author_name: string;
	author_title?: string;
	author_company?: string;
	author_image?: string;
	rating?: number;
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateTestimonialInput {
	content: string;
	author_name: string;
	author_title?: string;
	author_company?: string;
	author_image?: string;
	rating?: number;
	featured?: boolean;
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

export async function updateTestimonial(id: string, input: Partial<CreateTestimonialInput>): Promise<Testimonial> {
	const response = await apiClient.patch<ApiResponse<Testimonial>>(`/testimonials/${id}`, input);
	return response.data.data;
}

export async function deleteTestimonial(id: string): Promise<void> {
	await apiClient.delete(`/testimonials/${id}`);
}

// ============================================================================
// SYSTEM DESIGNS
// ============================================================================

export interface SystemDesign {
	id: string;
	title: string;
	slug: string;
	description: string;
	content?: string;
	diagram_url?: string;
	technologies?: string[];
	featured: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateSystemDesignInput {
	title: string;
	description: string;
	content?: string;
	diagram_url?: string;
	technologies?: string[];
	featured?: boolean;
}

export async function listSystemDesigns(): Promise<SystemDesign[]> {
	const response = await apiClient.get<ApiResponse<SystemDesign[]>>('/system-designs');
	return response.data.data ?? [];
}

export async function getSystemDesign(id: string): Promise<SystemDesign> {
	const response = await apiClient.get<ApiResponse<SystemDesign>>(`/system-designs/${id}`);
	return response.data.data;
}

export async function createSystemDesign(input: CreateSystemDesignInput): Promise<SystemDesign> {
	const response = await apiClient.post<ApiResponse<SystemDesign>>('/system-designs', input);
	return response.data.data;
}

export async function updateSystemDesign(id: string, input: Partial<CreateSystemDesignInput>): Promise<SystemDesign> {
	const response = await apiClient.patch<ApiResponse<SystemDesign>>(`/system-designs/${id}`, input);
	return response.data.data;
}

export async function deleteSystemDesign(id: string): Promise<void> {
	await apiClient.delete(`/system-designs/${id}`);
}

