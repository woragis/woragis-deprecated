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

		// Request interceptor for adding auth token if needed
		this.client.interceptors.request.use(
			(config) => {
				// You can add auth token here when implementing authentication
				// const token = getAuthToken();
				// if (token) {
				//   config.headers.Authorization = `Bearer ${token}`;
				// }
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
}

export const apiClient = new ApiClient();

