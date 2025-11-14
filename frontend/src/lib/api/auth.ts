import { apiClient } from '@clients/apiClient';

export type OAuthMode = 'login' | 'link';

interface OAuthStartResponse {
	authorization_url: string;
	mode: OAuthMode;
	provider: string;
	state: string;
}

interface OAuthProvidersResponse {
	providers: Array<{ id: string; name: string }>;
}

interface OAuthAccountsResponse {
	accounts: Array<{
		provider: string;
		linked_at: string;
		updated_at: string;
		scopes: string[];
	}>;
}

export interface LoginRequestPayload {
	email: string;
	password: string;
	device_fingerprint?: string;
	device_name?: string;
	user_agent?: string;
}

export interface LoginResponsePayload {
	user: {
		id: string;
		email: string;
		created_at: string;
		email_confirmed?: boolean;
		mfa_enabled?: boolean;
		preferred_locale?: string;
		role?: string;
	};
	access_token: string;
	refresh_token?: string;
	session_id?: string;
}

export async function login(payload: LoginRequestPayload) {
	return apiClient.post<{ data: LoginResponsePayload }>('/auth/login', payload);
}

export async function refreshSession(refreshToken: string, userAgent?: string) {
	return apiClient.post<{ data: { access_token: string; refresh_token: string; session_id: string } }>(
		'/auth/refresh',
		{
			refresh_token: refreshToken,
			user_agent: userAgent ?? globalThis?.navigator?.userAgent ?? ''
		}
	);
}

export async function logout(sessionId: string) {
	return apiClient.post('/auth/logout', {
		session_id: sessionId
	});
}

export async function listSessions() {
	return apiClient.get<{ data: { sessions: Array<{ id: string; device_id: string; user_agent: string; ip: string; created_at: string; expires_at: string; last_seen_at: string; is_revoked: boolean }> } }>(
		'/auth/sessions'
	);
}

export async function revokeOtherSessions(keepSessionId?: string) {
	return apiClient.post('/auth/sessions/revoke', {
		keep_session_id: keepSessionId
	});
}

export async function startOAuth(
	provider: string,
	mode: OAuthMode,
	redirectOrigin: string,
	deviceFingerprint?: string,
	deviceName?: string,
	userAgent?: string
) {
	return apiClient.post<{ data: OAuthStartResponse }>('/auth/oauth/start', {
		provider,
		mode,
		redirect_origin: redirectOrigin,
		device_fingerprint: deviceFingerprint,
		device_name: deviceName,
		user_agent: userAgent
	});
}

export async function listOAuthProviders() {
	return apiClient.get<{ data: OAuthProvidersResponse }>('/auth/oauth/providers');
}

export async function listOAuthAccounts() {
	return apiClient.get<{ data: OAuthAccountsResponse }>('/auth/oauth/accounts');
}

export async function unlinkOAuthAccount(provider: string) {
	return apiClient.delete(`/auth/oauth/accounts/${provider}`);
}

export async function resendConfirmation(email: string) {
	return apiClient.post('/auth/confirm/resend', {
		email
	});
}

export async function confirmEmail(token: string) {
	return apiClient.post<{ data: { user: unknown } }>('/auth/confirm', {
		token
	});
}

export async function requestPasswordReset(email: string) {
	return apiClient.post('/auth/password/reset/request', {
		email
	});
}

export async function confirmPasswordReset(token: string, password: string) {
	return apiClient.post('/auth/password/reset/confirm', {
		token,
		password
	});
}

export async function enableMFA(payload: { issuer: string; label: string; code?: string }) {
	return apiClient.post<{ data: { secret: string; backup_codes: string[]; provisioning_uri: string } }>(
		'/auth/mfa/enable',
		payload
	);
}

export async function verifyMFA(code: string) {
	return apiClient.post('/auth/mfa/verify', {
		code
	});
}

export async function disableMFA() {
	return apiClient.post('/auth/mfa/disable', {});
}

