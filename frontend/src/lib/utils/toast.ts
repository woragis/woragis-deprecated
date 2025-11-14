import { toast as sonnerToast } from 'svelte-sonner';
import type { ExternalToast } from 'svelte-sonner';

type ToastOptions = ExternalToast | undefined;

const DEFAULT_OPTIONS: Partial<ExternalToast> = {
	duration: 4500
};

const mergeOptions = (options?: ExternalToast): ExternalToast =>
	({
		...DEFAULT_OPTIONS,
		...(options ?? {})
	} satisfies Partial<ExternalToast>) as ExternalToast;

const extractAxiosErrorMessage = (error: unknown): string | null => {
	if (error && typeof error === 'object' && 'response' in error) {
		const response = (error as any).response;
		return (
			response?.data?.error?.details?.message ??
			response?.data?.error?.message ??
			response?.data?.message ??
			null
		);
	}
	return null;
};

export const toast = (message: string, options?: ToastOptions) => sonnerToast(message, mergeOptions(options));

export const toastSuccess = (message: string, options?: ToastOptions) =>
	sonnerToast.success(message, mergeOptions(options));

export const toastError = (message: string, options?: ToastOptions) =>
	sonnerToast.error(message, mergeOptions(options));

export const toastInfo = (message: string, options?: ToastOptions) =>
	sonnerToast.info(message, mergeOptions(options));

export const getApiErrorMessage = (error: unknown, fallback = 'Something went wrong') =>
	extractAxiosErrorMessage(error) ?? fallback;

export const toastApiError = (error: unknown, fallback = 'Something went wrong', options?: ToastOptions) => {
	const message = getApiErrorMessage(error, fallback);
	toastError(message, options);
	return message;
};

