import { apiClient } from '@clients/apiClient';

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
}

export interface AppendMessageInput {
	role: string;
	content: string;
	generate_reply?: boolean;
	provider?: string;
	model?: string;
	agent?: string;
}

export async function listConversations(): Promise<Conversation[]> {
	const response = await apiClient.get<ApiResponse<Conversation[]>>('/chats/conversations');
	return response.data.data ?? [];
}

export async function getConversation(id: string): Promise<Conversation> {
	const response = await apiClient.get<ApiResponse<Conversation>>(`/chats/conversations/${id}`);
	return response.data.data;
}

export async function createConversation(input: CreateConversationInput): Promise<Conversation> {
	const response = await apiClient.post<ApiResponse<Conversation>>('/chats/conversations', input);
	return response.data.data;
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

export async function searchConversations(query: string, includeArchived: boolean = false): Promise<Conversation[]> {
	const response = await apiClient.get<ApiResponse<Conversation[]>>(
		`/chats/conversations/search?q=${encodeURIComponent(query)}&includeArchived=${includeArchived}`
	);
	return response.data.data ?? [];
}

export async function shareTranscript(conversationId: string): Promise<{ share_code: string }> {
	const response = await apiClient.post<ApiResponse<{ share_code: string }>>(
		`/chats/conversations/${conversationId}/transcripts/share`
	);
	return response.data.data;
}

export const CHATS_STREAM_BASE = '/chats/stream';