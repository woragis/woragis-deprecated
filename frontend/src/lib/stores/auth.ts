import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const AUTH_COOKIE_KEY = 'woragis_auth';
const COOKIE_MAX_AGE_SECONDS = 60 * 60 * 24 * 7; // 7 days

export interface AuthUser {
	id: string;
	email: string;
	created_at: string;
}

export interface AuthState {
	user: AuthUser | null;
	token: string | null;
	isAuthenticated: boolean;
}

const defaultState: AuthState = {
	user: null,
	token: null,
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
		};

		if (!parsed?.user || !parsed?.token) {
			throw new Error('Auth cookie missing properties');
		}

		return {
			user: parsed.user,
			token: parsed.token,
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
			token: state.token
		})
	);
	writeCookie(AUTH_COOKIE_KEY, encoded, COOKIE_MAX_AGE_SECONDS);
};

const createAuthStore = () => {
	const { subscribe, set } = writable<AuthState>(loadInitialState());

	return {
		subscribe,
		setSession: (user: AuthUser, token: string) => {
			const nextState: AuthState = { user, token, isAuthenticated: true };
			set(nextState);
			persistState(nextState);
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

