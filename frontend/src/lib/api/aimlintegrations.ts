import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type IntegrationType =
	| 'rag'
	| 'llm'
	| 'ml_model'
	| 'computer_vision'
	| 'nlp'
	| 'recommendation'
	| 'chatbot'
	| 'anomaly_detection'
	| 'predictive_analytics'
	| 'generative_ai'
	| 'other';

export type Framework =
	| 'openai'
	| 'anthropic'
	| 'huggingface'
	| 'tensorflow'
	| 'pytorch'
	| 'langchain'
	| 'llamaindex'
	| 'cohere'
	| 'google_ai'
	| 'azure_ai'
	| 'aws_bedrock'
	| 'custom'
	| 'other';

export interface AIMLIntegration {
	id: string;
	userId: string;
	title: string;
	description: string;
	type: IntegrationType;
	framework: Framework;
	modelName?: string;
	modelVersion?: string;
	useCase?: string;
	impact?: string;
	technologies?: string[];
	architecture?: string;
	metrics?: string;
	projectId?: string;
	caseStudyId?: string;
	featured: boolean;
	displayOrder: number;
	demoUrl?: string;
	documentationUrl?: string;
	githubUrl?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateAIMLIntegrationInput {
	title: string;
	description: string;
	type: IntegrationType;
	framework: Framework;
	modelName?: string;
	modelVersion?: string;
	useCase?: string;
	impact?: string;
	technologies?: string[];
	architecture?: string;
	metrics?: string;
	projectId?: string;
	caseStudyId?: string;
	featured?: boolean;
	displayOrder?: number;
	demoUrl?: string;
	documentationUrl?: string;
	githubUrl?: string;
}

export interface UpdateAIMLIntegrationInput {
	title?: string;
	description?: string;
	type?: IntegrationType;
	framework?: Framework;
	modelName?: string;
	modelVersion?: string;
	useCase?: string;
	impact?: string;
	technologies?: string[];
	architecture?: string;
	metrics?: string;
	projectId?: string;
	caseStudyId?: string;
	featured?: boolean;
	displayOrder?: number;
	demoUrl?: string;
	documentationUrl?: string;
	githubUrl?: string;
}

export async function listAIMLIntegrations(): Promise<AIMLIntegration[]> {
	const response = await apiClient.get<ApiResponse<AIMLIntegration[]>>('/aiml-integrations');
	return response.data.data ?? [];
}

export async function getAIMLIntegration(id: string): Promise<AIMLIntegration> {
	const response = await apiClient.get<ApiResponse<AIMLIntegration>>(`/aiml-integrations/${id}`);
	return response.data.data;
}

export async function createAIMLIntegration(input: CreateAIMLIntegrationInput): Promise<AIMLIntegration> {
	const response = await apiClient.post<ApiResponse<AIMLIntegration>>('/aiml-integrations', input);
	return response.data.data;
}

export async function updateAIMLIntegration(
	id: string,
	input: UpdateAIMLIntegrationInput
): Promise<AIMLIntegration> {
	const response = await apiClient.patch<ApiResponse<AIMLIntegration>>(`/aiml-integrations/${id}`, input);
	return response.data.data;
}

export async function deleteAIMLIntegration(id: string): Promise<void> {
	await apiClient.delete(`/aiml-integrations/${id}`);
}

