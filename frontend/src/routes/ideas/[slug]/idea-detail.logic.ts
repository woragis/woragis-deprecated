import { derived, type Readable } from 'svelte/store';

import type { ChatConversation, Idea, Project } from '$lib/api/types';
import { useConversationsQuery } from '@hooks/chats';
import { useIdeaBySlugQuery } from '@hooks/ideas';
import { useProjectsListQuery } from '@hooks/projects';

export function createIdeaDetailLogic(slugStore: Readable<string | null>) {
	const ideaQuery = useIdeaBySlugQuery(slugStore);
	const idea = derived(ideaQuery, ($query) => $query.data ?? null);

	const conversationsQuery = useConversationsQuery({
		search: '',
		includeArchived: false,
		enabled: true
	});
	const projectsQuery = useProjectsListQuery({ enabled: true });

	const relatedChats = derived([idea, conversationsQuery], ([$idea, $conversations]) => {
		if (!$idea) {
			return [];
		}
		const chats = ($conversations.data ?? []) as ChatConversation[];
		return chats
			.filter((chat) => chat.idea_id === $idea.id)
			.sort(
				(a, b) =>
					new Date(b.updated_at ?? b.created_at ?? '').getTime() -
					new Date(a.updated_at ?? a.created_at ?? '').getTime()
			);
	});

	const relatedProjects = derived([idea, projectsQuery], ([$idea, $projects]) => {
		if (!$idea) return [];
		const projects = ($projects.data ?? []) as Project[];
		if (!$idea.project_id) return [];
		const project = projects.find((p) => p.id === $idea.project_id);
		return project ? [project] : [];
	});

	return {
		ideaQuery,
		idea,
		conversationsQuery,
		projectsQuery,
		relatedChats,
		relatedProjects
	};
}


