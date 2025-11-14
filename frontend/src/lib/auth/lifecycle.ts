import { browser } from '$app/environment';
import { get } from 'svelte/store';

import { logout as logoutRequest, refreshSession as refreshSessionRequest } from '$lib/api/auth';
import { authStore } from '$lib';

const REFRESH_INTERVAL_MS = 9 * 60 * 1000; // 9 minutes
const INITIAL_REFRESH_DELAY_MS = 60 * 1000; // 1 minute

let refreshInterval: ReturnType<typeof setInterval> | null = null;
let initialRefreshTimeout: ReturnType<typeof setTimeout> | null = null;
let isRefreshing = false;
let unsubscribe: (() => void) | null = null;

function startRefreshTimer() {
	if (refreshInterval) {
		return;
	}
	refreshInterval = setInterval(triggerRefresh, REFRESH_INTERVAL_MS);
	initialRefreshTimeout = setTimeout(triggerRefresh, INITIAL_REFRESH_DELAY_MS);
}

function stopRefreshTimer() {
	if (refreshInterval) {
		clearInterval(refreshInterval);
		refreshInterval = null;
	}
	if (initialRefreshTimeout) {
		clearTimeout(initialRefreshTimeout);
		initialRefreshTimeout = null;
	}
}

async function triggerRefresh() {
	if (isRefreshing) {
		return;
	}

	const state = get(authStore);
	if (!state.isAuthenticated || !state.refreshToken) {
		stopRefreshTimer();
		return;
	}

	isRefreshing = true;
	try {
		const response = await refreshSessionRequest(state.refreshToken, globalThis?.navigator?.userAgent);
		const payload = response.data?.data;
		if (payload?.access_token && payload?.refresh_token && payload?.session_id) {
			authStore.updateTokens({
				accessToken: payload.access_token,
				refreshToken: payload.refresh_token,
				sessionId: payload.session_id
			});
		}
	} catch (error) {
		// Errors are already logged by the API client.
		authStore.clear();
		stopRefreshTimer();
	} finally {
		isRefreshing = false;
	}
}

export function initAuthLifecycle() {
	if (!browser || unsubscribe) {
		return;
	}

	authStore.refreshFromCookies();

	unsubscribe = authStore.subscribe((state) => {
		if (state.isAuthenticated && state.refreshToken) {
			startRefreshTimer();
		} else {
			stopRefreshTimer();
		}
	});
}

export async function performLogout() {
	const state = get(authStore);
	try {
		if (state.sessionId) {
			await logoutRequest(state.sessionId);
		}
	} catch (error) {
		// Allow logout to continue even if the server rejects the request.
		console.error('Logout failed', error);
	} finally {
		authStore.clear();
		stopRefreshTimer();
	}
}

export async function refreshSessionNow() {
	await triggerRefresh();
}

export function teardownAuthLifecycle() {
	if (unsubscribe) {
		unsubscribe();
		unsubscribe = null;
	}
	stopRefreshTimer();
}

