import { apiClient } from './client';
import type { AIMLIntegration, ListAIMLIntegrationsParams } from '$lib/types/aiml-integration';

// List all AI/ML integrations
export async function listAIMLIntegrations(
	params?: ListAIMLIntegrationsParams
): Promise<AIMLIntegration[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.type) {
			queryParams.append('type', params.type);
		}
		if (params?.framework) {
			queryParams.append('framework', params.framework);
		}
		if (params?.projectId) {
			queryParams.append('projectId', params.projectId);
		}
		if (params?.featured !== undefined) {
			queryParams.append('featured', params.featured.toString());
		}
		if (params?.limit) {
			queryParams.append('limit', params.limit.toString());
		}
		if (params?.offset) {
			queryParams.append('offset', params.offset.toString());
		}
		if (params?.orderBy) {
			queryParams.append('orderBy', params.orderBy);
		}
		if (params?.order) {
			queryParams.append('order', params.order);
		}

		const queryString = queryParams.toString();
		const url = queryString ? `/aiml-integrations?${queryString}` : '/aiml-integrations';

		const response = await apiClient.get<{ success: boolean; data: AIMLIntegration[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching AI/ML integrations:', error);
		throw error;
	}
}

// Get featured AI/ML integrations (public endpoint)
export async function getFeaturedAIMLIntegrations(): Promise<AIMLIntegration[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: AIMLIntegration[] }>(
			'/aiml-integrations/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured AI/ML integrations:', error);
		throw error;
	}
}

// Get AI/ML integration by ID
export async function getAIMLIntegration(id: string): Promise<AIMLIntegration | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: AIMLIntegration }>(
			`/aiml-integrations/${id}/public`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching AI/ML integration ${id}:`, error);
		return null;
	}
}

// Get integrations by type
export async function getIntegrationsByType(type: string): Promise<AIMLIntegration[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: AIMLIntegration[] }>(
			`/aiml-integrations/type/${type}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching AI/ML integrations by type ${type}:`, error);
		return [];
	}
}

// Get integrations by framework
export async function getIntegrationsByFramework(framework: string): Promise<AIMLIntegration[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: AIMLIntegration[] }>(
			`/aiml-integrations/framework/${framework}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching AI/ML integrations by framework ${framework}:`, error);
		return [];
	}
}

