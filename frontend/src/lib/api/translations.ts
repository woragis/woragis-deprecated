import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type EntityType =
	| 'post'
	| 'skill'
	| 'project'
	| 'case_study'
	| 'certification'
	| 'testimonial'
	| 'technical_writing'
	| 'problem_solution'
	| 'social_media_post'
	| 'system_design';

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

export type TranslationStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface Translation {
	id: string;
	entity_type: EntityType;
	entity_id: string;
	language: Language;
	fields: Record<string, string>;
	status: TranslationStatus;
	error_message?: string;
	created_at: string;
	updated_at: string;
}

export interface RequestTranslationPayload {
	entityType: EntityType;
	entityId: string;
	language: Language;
	fields?: string[];
	sourceText?: Record<string, string>;
}

export interface TranslateEntityPayload {
	entityType: EntityType;
	entityId: string;
	languages?: Language[]; // If empty, translates to all languages
}

export interface TranslationFilters {
	entityType?: EntityType;
	entityId?: string;
	language?: Language;
	status?: TranslationStatus;
}

// Request a translation for a specific entity and language
export async function requestTranslation(payload: RequestTranslationPayload): Promise<{ message: string }> {
	const response = await apiClient.post<ApiResponse<{ message: string }>>('/translations/request', payload);
	return response.data.data;
}

// Translate an entity to all or specified languages
export async function translateEntity(payload: TranslateEntityPayload): Promise<{ message: string; queuedCount: number }> {
	const response = await apiClient.post<ApiResponse<{ message: string; queuedCount: number }>>(
		'/translations/translate-entity',
		payload
	);
	return response.data.data;
}

// Get translation for a specific entity and language
export async function getTranslation(
	entityType: EntityType,
	entityId: string,
	language: Language
): Promise<Translation> {
	const response = await apiClient.get<ApiResponse<Translation>>(
		`/translations/get?entityType=${entityType}&entityId=${entityId}&language=${language}`
	);
	return response.data.data;
}

// List translations with optional filters
export async function listTranslations(filters?: TranslationFilters): Promise<Translation[]> {
	const params = new URLSearchParams();
	if (filters?.entityType) params.append('entityType', filters.entityType);
	if (filters?.entityId) params.append('entityId', filters.entityId);
	if (filters?.language) params.append('language', filters.language);
	if (filters?.status) params.append('status', filters.status);

	const queryString = params.toString();
	const url = `/translations${queryString ? `?${queryString}` : ''}`;

	const response = await apiClient.get<ApiResponse<Translation[]>>(url);
	return response.data.data ?? [];
}

// Supported languages for translation
export const SUPPORTED_LANGUAGES: { value: Language; label: string }[] = [
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

