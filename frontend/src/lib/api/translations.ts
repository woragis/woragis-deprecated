import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type EntityType =
	| 'testimonial'
	| 'post'
	| 'project'
	| 'case_study'
	| 'project_case_study'
	| 'system_design'
	| 'problem_solution'
	| 'certification'
	| 'aiml_integration'
	| 'impact_metric'
	| 'social_media_post'
	| 'technical_writing'
	| 'interest'
	| 'skill';

export type Language =
	| 'en'
	| 'pt-BR'
	| 'fr'
	| 'es'
	| 'de'
	| 'ru'
	| 'ja'
	| 'ko'
	| 'zh-CN'
	| 'el'
	| 'la';

export interface LanguageOption {
	value: Language;
	label: string;
}

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
	{ value: 'en', label: 'English' },
	{ value: 'pt-BR', label: 'Portuguese (Brazil)' },
	{ value: 'fr', label: 'French' },
	{ value: 'es', label: 'Spanish' },
	{ value: 'de', label: 'German' },
	{ value: 'ru', label: 'Russian' },
	{ value: 'ja', label: 'Japanese' },
	{ value: 'ko', label: 'Korean' },
	{ value: 'zh-CN', label: 'Chinese (Simplified)' },
	{ value: 'el', label: 'Greek' },
	{ value: 'la', label: 'Latin' }
];

export type TranslationStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface Translation {
	id: string;
	entityType: EntityType;
	entityId: string;
	language: Language;
	fields: Record<string, string>;
	status: TranslationStatus;
	errorMessage?: string;
	createdAt: string;
	updatedAt: string;
}

export interface RequestTranslationInput {
	entityType: EntityType;
	entityId: string;
	targetLanguages: Language[];
}

export interface TranslateEntityInput {
	entityType: EntityType;
	entityId: string;
	language: Language;
	fields: Record<string, string>;
}

export interface GetTranslationParams {
	entityType: EntityType;
	entityId: string;
	language: Language;
}

export async function listTranslations(): Promise<Translation[]> {
	const response = await apiClient.get<ApiResponse<Translation[]>>('/translations');
	return response.data.data ?? [];
}

export async function getTranslation(params: GetTranslationParams): Promise<Translation> {
	const queryParams = new URLSearchParams({
		entityType: params.entityType,
		entityId: params.entityId,
		language: params.language
	});
	const response = await apiClient.get<ApiResponse<Translation>>(
		`/translations/get?${queryParams.toString()}`
	);
	return response.data.data;
}

export async function requestTranslation(input: RequestTranslationInput): Promise<void> {
	await apiClient.post('/translations/request', input);
}

export async function translateEntity(input: TranslateEntityInput): Promise<Translation> {
	const response = await apiClient.post<ApiResponse<Translation>>('/translations/translate-entity', input);
	return response.data.data;
}
