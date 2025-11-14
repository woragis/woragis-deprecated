import { createMutation, createQuery } from '@tanstack/svelte-query';

import {
	listOAuthAccounts,
	listOAuthProviders,
	listSessions,
	logout,
	revokeOtherSessions,
	unlinkOAuthAccount
} from '$lib/api/auth';
import type { OAuthAccountPayload, SessionPayload } from '$lib/api/types';

export interface AuthQueryOptions {
	enabled?: boolean;
}

export const useSessionsQuery = (options: AuthQueryOptions = {}) =>
	createQuery<SessionPayload[]>({
		queryKey: ['auth', 'sessions'],
		queryFn: async () => {
			const response = await listSessions();
			return response.data?.data?.sessions ?? [];
		},
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});

export const useOAuthProvidersQuery = (options: AuthQueryOptions = {}) =>
	createQuery({
		queryKey: ['auth', 'providers'],
		queryFn: async () => {
			const response = await listOAuthProviders();
			return response.data?.data?.providers ?? [];
		},
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});

export const useOAuthAccountsQuery = (options: AuthQueryOptions = {}) =>
	createQuery<OAuthAccountPayload[]>({
		queryKey: ['auth', 'accounts'],
		queryFn: async () => {
			const response = await listOAuthAccounts();
			return response.data?.data?.accounts ?? [];
		},
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});

export const useLogoutSessionMutation = () =>
	createMutation({
		mutationFn: (sessionId: string) => logout(sessionId)
	});

export const useRevokeOtherSessionsMutation = () =>
	createMutation({
		mutationFn: (keepSessionId?: string) => revokeOtherSessions(keepSessionId)
	});

export const useUnlinkOAuthAccountMutation = () =>
	createMutation({
		mutationFn: (provider: string) => unlinkOAuthAccount(provider)
	});

