import { apiClient } from '$lib/clients/apiClient';

export interface WhatsAppQRResponse {
	connected: boolean;
	qr_code: string | null;
	qr_text?: string;
	message: string;
}

export interface WhatsAppStatusResponse {
	connected: boolean;
	has_qr_code: boolean;
	message: string;
}

export async function fetchQRCode(): Promise<WhatsAppQRResponse> {
	const response = await apiClient.get<{ data: WhatsAppQRResponse }>('/whatsapp/qr');
	return response.data.data;
}

export async function fetchStatus(): Promise<WhatsAppStatusResponse> {
	const response = await apiClient.get<{ data: WhatsAppStatusResponse }>('/whatsapp/status');
	return response.data.data;
}

export interface SendMessageInput {
	client_id: string;
	message?: string;
	use_ai?: boolean;
	template?: string;
	instructions?: string;
	client_context?: string;
}

export interface SendMessageResponse {
	message: string;
	sent_message: string;
}

export async function sendMessage(input: SendMessageInput): Promise<SendMessageResponse> {
	const response = await apiClient.post<{ data: SendMessageResponse }>('/whatsapp/send', {
		client_id: input.client_id,
		message: input.message,
		use_ai: input.use_ai ?? false,
		template: input.template,
		instructions: input.instructions,
		client_context: input.client_context
	});
	return response.data.data;
}

