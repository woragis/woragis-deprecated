import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const AUTH_COOKIE_KEY = 'woragis_auth';
const COOKIE_MAX_AGE_SECONDS = 60 * 60 * 24 * 7; // 7 days

export interface AuthUser {
	id: string;
	display_name?: string;
	email: string;
	created_at: string;
	email_confirmed?: boolean;
	mfa_enabled?: boolean;
	preferred_locale?: string;
	role?: string;
}

export interface AuthState {
	user: AuthUser | null;
	token: string | null;
	refreshToken: string | null;
	sessionId: string | null;
	isAuthenticated: boolean;
}

const defaultState: AuthState = {
	user: null,
	token: null,
	refreshToken: null,
	sessionId: null,
	isAuthenticated: false
};

const readCookie = (name: string): string | null => {
	if (!browser) {
		return null;
	}

	const cookies = document.cookie ? document.cookie.split('; ') : [];
	for (const cookie of cookies) {
		const [key, ...rest] = cookie.split('=');
		if (key === name) {
			return rest.join('=');
		}
	}
	return null;
};

const writeCookie = (name: string, value: string, maxAgeSeconds: number) => {
	if (!browser) {
		return;
	}

	document.cookie = `${name}=${value}; Path=/; Max-Age=${maxAgeSeconds}; SameSite=Lax`;
};

const deleteCookie = (name: string) => {
	if (!browser) {
		return;
	}

	document.cookie = `${name}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax`;
};

const normalizeUser = (user: Partial<AuthUser> & Record<string, any>): AuthUser => {
	if (!user) {
		return {
			id: '',
			email: '',
			created_at: ''
		};
	}

	return {
		id: user.id ?? '',
		display_name: user.display_name ?? user.displayName ?? undefined,
		email: user.email ?? '',
		created_at: (user.created_at ?? user.createdAt ?? '') as string,
		email_confirmed: user.email_confirmed ?? user.emailConfirmed ?? undefined,
		mfa_enabled: user.mfa_enabled ?? user.mfaEnabled ?? undefined,
		preferred_locale: user.preferred_locale ?? user.preferredLocale ?? undefined,
		role: user.role ?? undefined
	};
};

const loadInitialState = (): AuthState => {
	if (!browser) {
		return defaultState;
	}

	const cookieValue = readCookie(AUTH_COOKIE_KEY);
	if (!cookieValue) {
		return defaultState;
	}

	try {
		const parsed = JSON.parse(decodeURIComponent(cookieValue)) as {
			user: AuthUser;
			token: string;
			refreshToken?: string | null;
			sessionId?: string | null;
		};

		if (!parsed?.user || !parsed?.token) {
			throw new Error('Auth cookie missing properties');
		}

		const normalizedUser = normalizeUser(parsed.user);

		if (!normalizedUser.id || !normalizedUser.email) {
			throw new Error('Auth cookie missing properties');
		}

		return {
			user: normalizedUser,
			token: parsed.token,
			refreshToken: parsed.refreshToken ?? null,
			sessionId: parsed.sessionId ?? null,
			isAuthenticated: true
		};
	} catch (error) {
		console.warn('Unable to parse auth cookie, clearing it.', error);
		deleteCookie(AUTH_COOKIE_KEY);
		return defaultState;
	}
};

const persistState = (state: AuthState) => {
	if (!browser) {
		return;
	}

	if (!state.isAuthenticated || !state.user || !state.token) {
		deleteCookie(AUTH_COOKIE_KEY);
		return;
	}

	const encoded = encodeURIComponent(
		JSON.stringify({
			user: state.user,
			token: state.token,
			refreshToken: state.refreshToken,
			sessionId: state.sessionId
		})
	);
	writeCookie(AUTH_COOKIE_KEY, encoded, COOKIE_MAX_AGE_SECONDS);
};

const createAuthStore = () => {
	const { subscribe, set, update } = writable<AuthState>(loadInitialState());

	return {
		subscribe,
		setSession: (
			user: AuthUser | Record<string, any>,
			token: string,
			extras?: { refreshToken?: string | null; sessionId?: string | null }
		) => {
			const normalizedUser = normalizeUser(user as AuthUser);
			const nextState: AuthState = {
				user: normalizedUser,
				token,
				refreshToken: extras?.refreshToken ?? null,
				sessionId: extras?.sessionId ?? null,
				isAuthenticated: true
			};
			set(nextState);
			persistState(nextState);
		},
		updateTokens: (tokens: { accessToken: string; refreshToken: string; sessionId: string }) => {
			update((state) => {
				if (!state.isAuthenticated || !state.user) {
					return state;
				}

				const nextState: AuthState = {
					...state,
					token: tokens.accessToken ?? state.token,
					refreshToken: tokens.refreshToken ?? state.refreshToken,
					sessionId: tokens.sessionId ?? state.sessionId
				};

				persistState(nextState);
				return nextState;
			});
		},
		updateUser: (updates: Partial<AuthUser>) => {
			update((state) => {
				if (!state.isAuthenticated || !state.user) {
					return state;
				}

				const nextUser = {
					...state.user,
					...updates
				};

				const nextState: AuthState = {
					...state,
					user: nextUser
				};

				persistState(nextState);
				return nextState;
			});
		},
		clear: () => {
			set(defaultState);
			persistState(defaultState);
		},
		refreshFromCookies: () => {
			set(loadInitialState());
		}
	};
};

export const authStore = createAuthStore();
