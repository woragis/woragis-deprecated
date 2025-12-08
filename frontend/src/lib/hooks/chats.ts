import { get, type Readable } from 'svelte/store';
import { createMutation, createQuery } from '@tanstack/svelte-query';

import type { AppendMessageInput } from '$lib/api/chats';
import {
	appendMessage,
	archiveConversations,
	createConversation as createConversationRequest,
	deleteConversations,
	fetchConversations,
	fetchMessages,
	listAssignments,
	listTranscripts,
	restoreConversations,
	searchConversations,
	shareTranscript as shareTranscriptRequest
} from '$lib/api/chats';
import type {
	ChatAssignment,
	ChatConversation,
	ChatMessage,
	ChatTranscript,
	UUID
} from '$lib/api/types';

type MaybeReadable<T> = T | Readable<T>;

const isReadable = <T>(value: MaybeReadable<T>): value is Readable<T> =>
	typeof value === 'object' && value !== null && 'subscribe' in value;

const resolve = <T>(value: MaybeReadable<T>): T => (isReadable(value) ? get(value) : value);

export interface ConversationsQueryOptions {
	search?: string;
	includeArchived?: boolean;
	enabled?: boolean;
}

// Helper to convert Conversation to ChatConversation
const mapConversation = (conv: any): ChatConversation => ({
	id: conv.id,
	user_id: conv.userId || conv.user_id,
	title: conv.title,
	description: conv.description,
	idea_id: conv.ideaId || conv.idea_id,
	project_id: conv.projectId || conv.project_id,
	assigned_agent_id: conv.assignedAgentId || conv.assigned_agent_id,
	shared_transcript: conv.sharedTranscript || conv.shared_transcript,
	archived_at: conv.archivedAt || conv.archived_at,
	deleted_at: conv.deletedAt || conv.deleted_at,
	created_at: conv.createdAt || conv.created_at,
	updated_at: conv.updatedAt || conv.updated_at
});

export const useConversationsQuery = (options: MaybeReadable<ConversationsQueryOptions>) =>
	createQuery<ChatConversation[]>({
		queryKey: ['chats', 'conversations', options],
		queryFn: async () => {
			const resolved = resolve(options);
			const search = resolved.search?.trim() ?? '';
			const includeArchived = resolved.includeArchived ?? false;
			const conversations = search
				? await searchConversations(search, includeArchived)
				: await fetchConversations();
			return conversations.map(mapConversation);
		},
		placeholderData: (previous) => previous ?? [],
		enabled: resolve(options).enabled ?? true
	});

// Helper to convert Message to ChatMessage
const mapMessage = (msg: any): ChatMessage => ({
	id: msg.id,
	conversation_id: msg.conversationId || msg.conversation_id,
	role: msg.role,
	content: msg.content,
	created_at: msg.createdAt || msg.created_at
});

export const useConversationMessagesQuery = (
	conversationId: MaybeReadable<UUID | null>,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ChatMessage[]>({
		queryKey: ['chats', 'conversation', conversationId, 'messages'],
		queryFn: async () => {
			const messages = await fetchMessages(resolve(conversationId)!);
			return messages.map(mapMessage);
		},
		enabled: Boolean(resolve(conversationId)) && (options.enabled ?? true),
		select: (data: ChatMessage[]) =>
			[...data].sort(
				(a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
			),
		placeholderData: () => []
	});

export const useConversationTranscriptsQuery = (
	conversationId: MaybeReadable<UUID | null>,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ChatTranscript[]>({
		queryKey: ['chats', 'conversation', conversationId, 'transcripts'],
		queryFn: () => listTranscripts(resolve(conversationId)!),
		enabled: Boolean(resolve(conversationId)) && (options.enabled ?? true),
		placeholderData: () => []
	});

export const useConversationAssignmentsQuery = (
	conversationId: MaybeReadable<UUID | null>,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ChatAssignment[]>({
		queryKey: ['chats', 'conversation', conversationId, 'assignments'],
		queryFn: () => listAssignments(resolve(conversationId)!),
		enabled: Boolean(resolve(conversationId)) && (options.enabled ?? true),
		placeholderData: () => []
	});

export const useCreateConversationMutation = () =>
	createMutation({
		mutationFn: createConversationRequest
	});

export const useArchiveConversationsMutation = () =>
	createMutation({
		mutationFn: (ids: UUID[]) => archiveConversations(ids)
	});

export const useRestoreConversationsMutation = () =>
	createMutation({
		mutationFn: (ids: UUID[]) => restoreConversations(ids)
	});

export const useDeleteConversationsMutation = () =>
	createMutation({
		mutationFn: (ids: UUID[]) => deleteConversations(ids)
	});

export const useAppendMessageMutation = () =>
	createMutation({
		mutationFn: ({
			conversationId,
			input
		}: {
			conversationId: UUID;
			input: AppendMessageInput;
		}) => appendMessage(conversationId, input)
	});

export const useShareTranscriptMutation = () =>
	createMutation({
		mutationFn: (conversationId: UUID) => shareTranscriptRequest(conversationId)
	});

