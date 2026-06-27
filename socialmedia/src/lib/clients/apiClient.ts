import { browser } from '$app/environment';
import axios from 'axios';

const baseURL = (import.meta.env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '');
const apiBaseURL = `${baseURL}/api`;

export const API_BASE_URL = apiBaseURL;

const AUTH_BYPASS_PREFIXES = [
	'/auth/login',
	'/auth/register',
	'/auth/refresh',
	'/auth/confirm',
	'/auth/confirm/resend',
	'/auth/password/reset/request',
	'/auth/password/reset/confirm',
	'/auth/oauth/start',
	'/auth/oauth/providers'
];

const isAuthBypass = (url?: string) => {
	if (!url) {
		return false;
	}

	const endpoint = url.startsWith('http') ? new URL(url).pathname : url;
	return AUTH_BYPASS_PREFIXES.some((prefix) => endpoint.startsWith(prefix));
};

export const apiClient = axios.create({
	baseURL: apiBaseURL,
	headers: {
		'Content-Type': 'application/json'
	}
});

apiClient.interceptors.request.use((config: any) => {
	if (!browser || isAuthBypass(config.url ?? config.baseURL)) {
		return config;
	}

	// TODO: Add auth token when auth is implemented
	// const token = get(authStore);
	// if (token) {
	// 	config.headers = config.headers ?? {};
	// 	config.headers['Authorization'] = `Bearer ${token}`;
	// }

	return config;
});

apiClient.interceptors.response.use(
	(response: any) => response,
	(error: any) => {
		if (error.response) {
			console.error('API error', {
				url: error.config?.url,
				status: error.response.status,
				data: error.response.data
			});

			if (error.response.status === 401) {
				// TODO: Handle auth when implemented
				// authStore.clear();
			}
		}
		return Promise.reject(error);
	}
);
