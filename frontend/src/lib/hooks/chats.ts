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

export const useConversationsQuery = (options: MaybeReadable<ConversationsQueryOptions>) =>
	createQuery<ChatConversation[]>({
		queryKey: ['chats', 'conversations', options],
		queryFn: () => {
			const resolved = resolve(options);
			const search = resolved.search?.trim() ?? '';
			const includeArchived = resolved.includeArchived ?? false;
			return search ? searchConversations(search, includeArchived) : fetchConversations();
		},
		placeholderData: (previous) => previous ?? [],
		enabled: resolve(options).enabled ?? true
	});

export const useConversationMessagesQuery = (
	conversationId: MaybeReadable<UUID | null>,
	options: { enabled?: boolean } = {}
) =>
	createQuery<ChatMessage[]>({
		queryKey: ['chats', 'conversation', conversationId, 'messages'],
		queryFn: () => fetchMessages(resolve(conversationId)!),
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

