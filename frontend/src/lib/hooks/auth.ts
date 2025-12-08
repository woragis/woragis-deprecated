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
			const sessions = response.data?.data?.sessions ?? [];
			// Map API response to SessionPayload format
			return sessions.map((session: any): SessionPayload => ({
				id: session.id,
				user_id: session.user_id || '',
				device_fingerprint: session.device_id,
				device_name: session.device_name,
				user_agent: session.user_agent,
				ip_address: session.ip,
				created_at: session.created_at,
				last_activity_at: session.last_seen_at,
				expires_at: session.expires_at,
				is_current: !session.is_revoked
			}));
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

