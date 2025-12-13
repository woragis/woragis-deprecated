import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Conversation {
	id: string;
	userId: string;
	title: string;
	description?: string;
	ideaId?: string;
	projectId?: string;
	jobApplicationId?: string;
	assignedAgentId?: string;
	sharedTranscript?: string;
	archivedAt?: string;
	deletedAt?: string;
	lastAssignedAt?: string;
	createdAt: string;
	updatedAt: string;
}

export interface Message {
	id: string;
	conversationId: string;
	role: string;
	content: string;
	createdAt: string;
}

export interface CreateConversationInput {
	title: string;
	description?: string;
	ideaId?: string;
	projectId?: string;
	jobApplicationId?: string;
}

export interface AppendMessageInput {
	role: string;
	content: string;
	generate_reply?: boolean;
	provider?: string;
	model?: string;
	agent?: string;
}

// Helper function to transform snake_case API response to camelCase
function transformConversation(apiConv: any): Conversation {
	return {
		id: apiConv.id,
		userId: apiConv.user_id,
		title: apiConv.title,
		description: apiConv.description,
		ideaId: apiConv.idea_id,
		projectId: apiConv.project_id,
		jobApplicationId: apiConv.job_application_id,
		assignedAgentId: apiConv.assigned_agent_id,
		sharedTranscript: apiConv.shared_transcript,
		archivedAt: apiConv.archived_at,
		deletedAt: apiConv.deleted_at,
		lastAssignedAt: apiConv.last_assigned_at,
		createdAt: apiConv.created_at,
		updatedAt: apiConv.updated_at
	};
}

export async function listConversations(): Promise<Conversation[]> {
	const response = await apiClient.get<ApiResponse<any[]>>('/chats/conversations');
	const conversations = response.data.data ?? [];
	return conversations.map(transformConversation);
}

export async function getConversation(id: string): Promise<Conversation> {
	const response = await apiClient.get<ApiResponse<any>>(`/chats/conversations/${id}`);
	return transformConversation(response.data.data);
}

export async function createConversation(input: CreateConversationInput): Promise<Conversation> {
	const response = await apiClient.post<ApiResponse<any>>('/chats/conversations', input);
	return transformConversation(response.data.data);
}

export async function listMessages(conversationId: string): Promise<Message[]> {
	const response = await apiClient.get<ApiResponse<Message[]>>(
		`/chats/conversations/${conversationId}/messages`
	);
	return response.data.data ?? [];
}

export async function appendMessage(
	conversationId: string,
	input: AppendMessageInput
): Promise<Message> {
	const response = await apiClient.post<ApiResponse<Message>>(
		`/chats/conversations/${conversationId}/messages`,
		input
	);
	return response.data.data;
}

export async function archiveConversations(ids: string[]): Promise<void> {
	await apiClient.post('/chats/conversations/archive', { ids });
}

export async function deleteConversations(ids: string[]): Promise<void> {
	await apiClient.post('/chats/conversations/delete', { ids });
}

export async function restoreConversations(ids: string[]): Promise<void> {
	await apiClient.post('/chats/conversations/restore', { ids });
}

// Aliases for hooks compatibility
export const fetchConversations = listConversations;
export const fetchMessages = listMessages;

export async function listAssignments(conversationId: string): Promise<any[]> {
	const response = await apiClient.get<ApiResponse<any[]>>(
		`/chats/conversations/${conversationId}/assignments`
	);
	return response.data.data ?? [];
}

export async function listTranscripts(conversationId: string): Promise<any[]> {
	const response = await apiClient.get<ApiResponse<any[]>>(
		`/chats/conversations/${conversationId}/transcripts`
	);
	return response.data.data ?? [];
}

export async function searchConversations(
	query?: string,
	includeArchived: boolean = false,
	jobApplicationId?: string
): Promise<Conversation[]> {
	const params = new URLSearchParams();
	if (query) params.append('q', query);
	params.append('include_archived', includeArchived.toString());
	if (jobApplicationId) params.append('job_application_id', jobApplicationId);

	const response = await apiClient.get<ApiResponse<any[]>>(
		`/chats/conversations/search?${params.toString()}`
	);
	const conversations = response.data.data ?? [];
	return conversations.map(transformConversation);
}

export interface ContextPreview {
	context: string;
	options: {
		includeJobApplication: boolean;
		includeResume: boolean;
		includeUserProfile: boolean;
		includeProjects: boolean;
		includeCaseStudies: boolean;
		includeTechnicalWritings: boolean;
		includePosts: boolean;
		includeProblemSolutions: boolean;
		includeSkills: boolean;
		includeExperiences: boolean;
	};
	message?: string;
}

export async function getContextPreview(conversationId: string): Promise<ContextPreview> {
	const response = await apiClient.get<ApiResponse<ContextPreview>>(`/chats/conversations/${conversationId}/context`);
	return response.data.data;
}

export async function shareTranscript(conversationId: string): Promise<{ share_code: string }> {
	const response = await apiClient.post<ApiResponse<{ share_code: string }>>(
		`/chats/conversations/${conversationId}/transcripts/share`
	);
	return response.data.data;
}

export const CHATS_STREAM_BASE = '/chats/stream';