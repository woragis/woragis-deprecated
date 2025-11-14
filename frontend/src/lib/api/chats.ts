import { apiClient, API_BASE_URL } from '@clients/apiClient';
import type {
	ChatAssignment,
	ChatConversation,
	ChatMessage,
	ChatTranscript,
	UUID
} from './types';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

const mapConversation = (item: ChatConversation): ChatConversation => ({
	...item,
	description: item.description ?? '',
	idea_id: item.idea_id ?? null,
	project_id: item.project_id ?? null,
	assigned_agent_id: item.assigned_agent_id ?? null,
	shared_transcript: item.shared_transcript ?? '',
	archived_at: item.archived_at ?? null,
	deleted_at: item.deleted_at ?? null
});

export async function fetchConversations(): Promise<ChatConversation[]> {
	const response = await apiClient.get<ApiResponse<ChatConversation[]>>('/chats/conversations');
	const conversations = response.data.data ?? [];
	return conversations.map(mapConversation);
}

export interface CreateConversationInput {
	title: string;
	description?: string;
	ideaId?: UUID | null;
	projectId?: UUID | null;
}

export async function createConversation(
	input: CreateConversationInput
): Promise<ChatConversation> {
	const response = await apiClient.post<ApiResponse<ChatConversation>>('/chats/conversations', {
		title: input.title,
		description: input.description ?? '',
		idea_id: input.ideaId ?? undefined,
		project_id: input.projectId ?? undefined
	});
	return mapConversation(response.data.data);
}

export async function searchConversations(
	query: string,
	includeArchived = false
): Promise<ChatConversation[]> {
	const response = await apiClient.get<ApiResponse<ChatConversation[]>>('/chats/conversations/search', {
		params: {
			q: query,
			include_archived: includeArchived
		}
	});
	const conversations = response.data.data ?? [];
	return conversations.map(mapConversation);
}

export async function fetchMessages(conversationId: UUID): Promise<ChatMessage[]> {
	const response = await apiClient.get<ApiResponse<ChatMessage[]>>(
		`/chats/conversations/${conversationId}/messages`
	);
	return response.data.data ?? [];
}

export interface AppendMessageInput {
	role: string;
	content: string;
	generate_reply?: boolean;
	provider?: string;
	model?: string;
	max_tokens?: number;
	temperature?: number;
}

export async function appendMessage(
	conversationId: UUID,
	payload: AppendMessageInput
): Promise<ChatMessage[]> {
	const response = await apiClient.post<ApiResponse<ChatMessage[]>>(
		`/chats/conversations/${conversationId}/messages`,
		{
			role: payload.role,
			content: payload.content,
			generate_reply: payload.generate_reply ?? false,
			provider: payload.provider ?? 'openai',
			model: payload.model ?? '',
			max_tokens: payload.max_tokens ?? 512,
			temperature: payload.temperature ?? 0.2
		}
	);
	return response.data.data ?? [];
}

async function postBulkAction(
	path: string,
	conversationIds: UUID[]
): Promise<{ status: string }> {
	const response = await apiClient.post<ApiResponse<{ status: string }>>(
		`/chats/conversations/${path}`,
		{
			conversation_ids: conversationIds
		}
	);
	return response.data.data ?? { status: 'ok' };
}

export const archiveConversations = (ids: UUID[]) => postBulkAction('archive', ids);
export const deleteConversations = (ids: UUID[]) => postBulkAction('delete', ids);
export const restoreConversations = (ids: UUID[]) => postBulkAction('restore', ids);

export async function shareTranscript(
	conversationId: UUID,
	expireAfter?: string
): Promise<ChatTranscript> {
	const response = await apiClient.post<ApiResponse<ChatTranscript>>(
		`/chats/conversations/${conversationId}/transcripts`,
		{
			expire_after: expireAfter
		}
	);
	return response.data.data!;
}

export async function listTranscripts(conversationId: UUID): Promise<ChatTranscript[]> {
	const response = await apiClient.get<ApiResponse<ChatTranscript[]>>(
		`/chats/conversations/${conversationId}/transcripts`
	);
	return response.data.data ?? [];
}

export async function listAssignments(conversationId: UUID): Promise<ChatAssignment[]> {
	const response = await apiClient.get<ApiResponse<ChatAssignment[]>>(
		`/chats/conversations/${conversationId}/assignments`
	);
	return response.data.data ?? [];
}

export const CHATS_STREAM_BASE = API_BASE_URL.replace(/^http/, 'ws');

