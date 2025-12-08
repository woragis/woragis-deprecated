import { onDestroy, onMount } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type { Project, ProjectStatus, UUID } from '$lib/api/types';
import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';
import { useCreateProjectMutation, useProjectsListQuery } from '@hooks/projects';

import { statusOptions } from './status';

export interface ProjectFormState {
	name: string;
	description: string;
	status: ProjectStatus;
	healthScore: number;
}

const defaultProjectForm = (): ProjectFormState => ({
	name: '',
	description: '',
	status: 'planning',
	healthScore: 60
});

type FormUpdater<T> = {
	<K extends keyof T>(field: K, value: T[K]): void;
};

interface AuthState {
	isAuthenticated: boolean;
	userId: UUID | null;
}

export function createProjectsListLogic() {
	const queryClient = useQueryClient();

	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);
	const error = writable<string | null>(null);

	const projects = writable<Project[]>([]);
	const projectForm = writable<ProjectFormState>(defaultProjectForm());

	const projectsQueryOptions = derived(authStateStore, ($auth) => ({ enabled: $auth.isAuthenticated }));
	const projectsQuery = useProjectsListQuery(projectsQueryOptions);
	const createProjectMutation = useCreateProjectMutation();

	const invalidateProjectsList = () => queryClient.invalidateQueries({ queryKey: ['projects', 'list'] });

	const resetStateForUnauthenticated = () => {
		projects.set([]);
		projectForm.set(defaultProjectForm());
		error.set(null);
		queryClient.removeQueries({ queryKey: ['projects'] });
	};

	const projectsQueryUnsubscribe = projectsQuery.subscribe((result) => {
		if (result.data) {
			projects.set(result.data);
			error.set(null);
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load projects.');
		}
	});

	const loadProjects = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			const message = 'You must be signed in to load projects.';
			error.set(message);
			toastError(message);
			return;
		}

		try {
			await invalidateProjectsList();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to load projects.');
			error.set(message);
			toastError(message);
		}
	};

	const updateProjectFormField: FormUpdater<ProjectFormState> = (field, value) => {
		projectForm.update((current) => ({ ...current, [field]: value }));
	};

	const handleCreateProject = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			const message = 'You must be signed in to create projects.';
			error.set(message);
			toastError(message);
			return;
		}

		const form = get(projectForm);
		if (!form.name.trim()) {
			const message = 'Project name is required.';
			error.set(message);
			toastError(message);
			return;
		}

		try {
			await get(createProjectMutation).mutateAsync({
				name: form.name.trim(),
				description: form.description,
				status: form.status,
				healthScore: form.healthScore
			});
			projectForm.set(defaultProjectForm());
			await invalidateProjectsList();
			toastSuccess('Project created.');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create project.');
			error.set(message);
			toastError(message);
		}
	};

	onMount(() => {
		const unsubscribeAuth = authStore.subscribe((state) => {
			const next: AuthState = {
				isAuthenticated: state.isAuthenticated,
				userId: state.user?.id ?? null
			};
			authStateStore.set(next);

			if (!next.isAuthenticated) {
				resetStateForUnauthenticated();
				return;
			}

			void loadProjects();
		});

		return () => {
			unsubscribeAuth();
		};
	});

	onDestroy(() => {
		projectsQueryUnsubscribe();
	});

	return {
		isAuthenticated,
		error,
		projects,
		projectForm,
		loadProjects,
		updateProjectFormField,
		handleCreateProject,
		statusOptions
	};
}


