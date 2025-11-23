import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios';
import { api } from '$lib/constants';

class ApiClient {
	private client: AxiosInstance;

	constructor() {
		// Debug: Log environment variables in development
		if (import.meta.env.DEV) {
			console.log('API Client initialized');
			console.log('api.apiKey from constants:', api.apiKey ? `${api.apiKey.substring(0, 8)}...` : 'NOT FOUND');
			console.log('api.baseURL from constants:', api.baseURL || 'NOT FOUND');
		}
		
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
						// Use uppercase X-API-Key header
						if (config.headers) {
							config.headers['X-API-Key'] = apiKey;
						}
					} else {
						console.warn('API Key not found in environment variables. Make sure PUBLIC_API_KEY is set in .env and restart the dev server.');
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

	// API Key management - only from environment variable
	private getAPIKey(): string | null {
		// Use the API key from constants (which reads from import.meta.env.PUBLIC_API_KEY)
		return api.apiKey;
	}
}

export const apiClient = new ApiClient();

