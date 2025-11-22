import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios';
import { api } from '$lib/constants';

class ApiClient {
	private client: AxiosInstance;

	constructor() {
		this.client = axios.create({
			baseURL: api.baseURL,
			timeout: api.timeout,
			headers: {
				'Content-Type': 'application/json'
			}
		});

		// Request interceptor for adding API key for GET requests
		this.client.interceptors.request.use(
			(config) => {
				// Add API key for GET requests (read-only access)
				if (config.method?.toLowerCase() === 'get') {
					const apiKey = this.getAPIKey();
					if (apiKey) {
						config.headers['X-API-Key'] = apiKey;
					}
				}
				return config;
			},
			(error) => {
				return Promise.reject(error);
			}
		);

		// Response interceptor for error handling
		this.client.interceptors.response.use(
			(response) => response,
			(error) => {
				// Handle common errors
				if (error.response) {
					// Server responded with error
					console.error('API Error:', error.response.data);
				} else if (error.request) {
					// Request was made but no response
					console.error('Network Error:', error.request);
				} else {
					// Something else happened
					console.error('Error:', error.message);
				}
				return Promise.reject(error);
			}
		);
	}

	async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
		const response = await this.client.get<T>(url, config);
		return response.data;
	}

	async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
		const response = await this.client.post<T>(url, data, config);
		return response.data;
	}

	async patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
		const response = await this.client.patch<T>(url, data, config);
		return response.data;
	}

	async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
		const response = await this.client.delete<T>(url, config);
		return response.data;
	}

	// API Key management
	private getAPIKey(): string | null {
		if (typeof window === 'undefined') {
			return null;
		}
		// Try to get from localStorage first
		const stored = localStorage.getItem('woragis_api_key');
		if (stored) {
			return stored;
		}
		// Fallback to environment variable
		return import.meta.env.PUBLIC_API_KEY || null;
	}

	setAPIKey(apiKey: string | null): void {
		if (typeof window === 'undefined') {
			return;
		}
		if (apiKey) {
			localStorage.setItem('woragis_api_key', apiKey);
		} else {
			localStorage.removeItem('woragis_api_key');
		}
	}
}

export const apiClient = new ApiClient();

