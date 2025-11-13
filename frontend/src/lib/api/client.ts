import { browser } from '$app/environment';
import axios from 'axios';
import { get } from 'svelte/store';

import { authStore } from '$lib';

const baseURL = (import.meta.env.PUBLIC_API_BASE_URL ?? 'http://localhost:8080').replace(/\/+$/, '');
const apiBaseURL = `${baseURL}/api`;

const isAuthBypass = (url?: string) => {
	if (!url) {
		return false;
	}

	const endpoint = url.startsWith('http') ? new URL(url).pathname : url;
	return endpoint.startsWith('/auth/login') || endpoint.startsWith('/auth/register');
};

export const apiClient = axios.create({
	baseURL: apiBaseURL,
	headers: {
		'Content-Type': 'application/json'
	}
});

apiClient.interceptors.request.use((config) => {
	if (!browser || isAuthBypass(config.url ?? config.baseURL)) {
		return config;
	}

	const { token } = get(authStore);
	if (token) {
		config.headers = config.headers ?? {};
		config.headers['Authorization'] = `Bearer ${token}`;
	}

	return config;
});

apiClient.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response) {
			console.error('API error', {
				url: error.config?.url,
				status: error.response.status,
				data: error.response.data
			});
		}
		return Promise.reject(error);
	}
);
