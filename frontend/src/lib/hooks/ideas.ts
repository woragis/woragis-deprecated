import { createMutation, createQuery } from '@tanstack/svelte-query';

import {
	bulkMoveIdeas,
	bulkUpdateIdeas,
	createIdea,
	createLink,
	deleteIdeas,
	fetchIdeas,
	fetchLinks,
	fetchVersions,
	moveIdea,
	restoreIdeas,
	updateIdea
} from '$lib/api/ideas';
import type {
	CreateIdeaInput,
	CreateLinkInput,
	MoveIdeaInput,
	UpdateIdeaInput
} from '$lib/api/ideas';
import type { Idea, IdeaLink, IdeaVersion, UUID } from '$lib/api/types';

export const useIdeasCanvasQuery = (enabled = true) =>
	createQuery<Idea[]>({
		queryKey: ['ideas', 'canvas'],
		queryFn: () => fetchIdeas(),
		enabled
	});

export const useIdeaLinksQuery = (enabled = true) =>
	createQuery<IdeaLink[]>({
		queryKey: ['ideas', 'links'],
		queryFn: () => fetchLinks(),
		enabled
	});

export const useIdeaVersionsQuery = (
	ideaId: string | null,
	options?: { enabled?: boolean; limit?: number }
) =>
	createQuery<IdeaVersion[]>({
		queryKey: ['ideas', ideaId, 'versions', options?.limit ?? 15],
		queryFn: () => fetchVersions(ideaId!, options?.limit ?? 15),
		enabled: Boolean(ideaId) && (options?.enabled ?? true),
		placeholderData: () => []
	});

export const useCreateIdeaMutation = () =>
	createMutation({
		mutationFn: (input: CreateIdeaInput) => createIdea(input)
	});

interface UpdateIdeaVariables {
	ideaId: UUID;
	input: UpdateIdeaInput;
}

export const useUpdateIdeaMutation = () =>
	createMutation({
		mutationFn: ({ ideaId, input }: UpdateIdeaVariables) => updateIdea(ideaId, input)
	});

interface MoveIdeaVariables {
	ideaId: UUID;
	input: MoveIdeaInput;
}

export const useMoveIdeaMutation = () =>
	createMutation({
		mutationFn: ({ ideaId, input }: MoveIdeaVariables) => moveIdea(ideaId, input)
	});

export const useBulkMoveIdeasMutation = () =>
	createMutation({
		mutationFn: bulkMoveIdeas
	});

export const useBulkUpdateIdeasMutation = () =>
	createMutation({
		mutationFn: bulkUpdateIdeas
	});

export const useDeleteIdeasMutation = () =>
	createMutation({
		mutationFn: deleteIdeas
	});

export const useRestoreIdeasMutation = () =>
	createMutation({
		mutationFn: restoreIdeas
	});

export const useCreateLinkMutation = () =>
	createMutation({
		mutationFn: (input: CreateLinkInput) => createLink(input)
	});

