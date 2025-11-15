import { derived, readable, type Readable } from 'svelte/store';
import { createMutation, createQuery } from '@tanstack/svelte-query';

import { ideasApi } from '$lib/api/ideas';
import type {
	CreateIdeaInput,
	CreateLinkInput,
	MoveIdeaInput,
	UpdateIdeaInput
} from '$lib/api/ideas';
import type { Idea, IdeaLink, IdeaVersion, UUID } from '$lib/api/types';

type MaybeReadable<T> = T | Readable<T>;

const isReadable = (value: unknown): value is Readable<unknown> =>
	typeof value === 'object' && value !== null && 'subscribe' in value;

const toReadable = <T>(value: MaybeReadable<T>): Readable<T> =>
	(isReadable(value) ? (value as Readable<T>) : readable(value)) as Readable<T>;

const toOptionsReadable = <T extends object>(value: MaybeReadable<T> | undefined, fallback: T): Readable<T> =>
	value === undefined ? readable(fallback) : toReadable(value);

export const useIdeasCanvasQuery = (enabled = true) =>
	createQuery<Idea[]>({
		queryKey: ['ideas', 'canvas'],
		queryFn: () => ideasApi.fetchIdeas(),
		enabled
	});

export const useIdeaBySlugQuery = (
	slug: MaybeReadable<string | null>,
	options?: MaybeReadable<{ enabled?: boolean }>
) => {
	const slugStore = toReadable(slug);
	const optionsStore = toOptionsReadable(options, {});

	return createQuery<Idea>(
		derived([slugStore, optionsStore], ([$slug, $options]) => ({
			queryKey: ['ideas', 'slug', $slug],
			queryFn: () => {
				if (!$slug) {
					throw new Error('Idea slug is required to load idea details');
				}
				return ideasApi.fetchIdeaBySlug($slug);
			},
			enabled: Boolean($slug) && ($options.enabled ?? true),
			retry: false
		}))
	);
};

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

export const useIdeasReferenceQuery = (enabled = true) =>
	createQuery<Idea[]>({
		queryKey: ['ideas', 'reference'],
		queryFn: () => ideasApi.fetchIdeas(),
		enabled,
		staleTime: Infinity,
		placeholderData: (previous) => previous ?? []
	});

