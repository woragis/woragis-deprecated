import { createMutation, createQuery } from '@tanstack/svelte-query';

import { ideasApi } from '$lib/api/ideas';
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
		queryFn: () => ideasApi.fetchIdeas(),
		enabled
	});

export const useIdeaLinksQuery = (enabled = true) =>
	createQuery<IdeaLink[]>({
		queryKey: ['ideas', 'links'],
		queryFn: () => ideasApi.fetchLinks(),
		enabled
	});

export const useIdeaVersionsQuery = (
	ideaId: string | null,
	options?: { enabled?: boolean; limit?: number }
) =>
	createQuery<IdeaVersion[]>({
		queryKey: ['ideas', ideaId, 'versions', options?.limit ?? 15],
		queryFn: () => ideasApi.fetchVersions(ideaId!, options?.limit ?? 15),
		enabled: Boolean(ideaId) && (options?.enabled ?? true),
		placeholderData: () => []
	});

export const useCreateIdeaMutation = () =>
	createMutation({
		mutationFn: (input: CreateIdeaInput) => ideasApi.createIdea(input)
	});

interface UpdateIdeaVariables {
	ideaId: UUID;
	input: UpdateIdeaInput;
}

export const useUpdateIdeaMutation = () =>
	createMutation({
		mutationFn: ({ ideaId, input }: UpdateIdeaVariables) => ideasApi.updateIdea(ideaId, input)
	});

interface MoveIdeaVariables {
	ideaId: UUID;
	input: MoveIdeaInput;
}

export const useMoveIdeaMutation = () =>
	createMutation({
		mutationFn: ({ ideaId, input }: MoveIdeaVariables) => ideasApi.moveIdea(ideaId, input)
	});

export const useBulkMoveIdeasMutation = () =>
	createMutation({
		mutationFn: ideasApi.bulkMoveIdeas
	});

export const useBulkUpdateIdeasMutation = () =>
	createMutation({
		mutationFn: ideasApi.bulkUpdateIdeas
	});

export const useDeleteIdeasMutation = () =>
	createMutation({
		mutationFn: ideasApi.deleteIdeas
	});

export const useRestoreIdeasMutation = () =>
	createMutation({
		mutationFn: ideasApi.restoreIdeas
	});

export const useCreateLinkMutation = () =>
	createMutation({
		mutationFn: (input: CreateLinkInput) => ideasApi.createLink(input)
	});

