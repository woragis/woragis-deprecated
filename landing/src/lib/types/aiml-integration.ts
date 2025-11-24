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

export interface ListAIMLIntegrationsParams {
	type?: IntegrationType;
	framework?: Framework;
	projectId?: string;
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

