import { apiClient } from './client';
import type { Transaction, TransactionSummary, TransactionType, UUID } from './types';

const DEFAULT_BASE_CURRENCY = 'USD';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface TransactionFilters {
	search?: string;
	includeArchived?: boolean;
	types?: TransactionType[];
	categories?: string[];
}

export interface CreateTransactionInput {
	type: TransactionType;
	category: string;
	description?: string;
	amount: number;
	currency: string;
	baseCurrency?: string;
	occurredAt?: string;
	isRecurring?: boolean;
	isEssential?: boolean;
	tags?: string[];
}

export async function fetchTransactions(filters: TransactionFilters = {}): Promise<Transaction[]> {
	const { search, includeArchived, types, categories } = filters;
	const response = await apiClient.get<ApiResponse<Transaction[]>>('/finance/transactions', {
		params: {
			search,
			include_archived: includeArchived,
			types: types?.join(','),
			categories: categories?.join(',')
		}
	});
	return response.data.data ?? [];
}

export async function fetchSummary(from?: string, to?: string): Promise<TransactionSummary> {
	const response = await apiClient.get<ApiResponse<TransactionSummary>>('/finance/summary', {
		params: {
			from,
			to
		}
	});
	return response.data.data;
}

export async function createTransaction(input: CreateTransactionInput): Promise<Transaction> {
	const {
		type,
		category,
		description,
		amount,
		currency,
		baseCurrency = DEFAULT_BASE_CURRENCY,
		occurredAt,
		isRecurring,
		isEssential,
		tags
	} = input;

	const response = await apiClient.post<ApiResponse<Transaction>>('/finance/transactions', {
		type,
		category,
		description,
		amount: Number(amount),
		currency,
		base_currency: baseCurrency,
		occurred_at: occurredAt ?? new Date().toISOString(),
		is_recurring: Boolean(isRecurring),
		is_essential: Boolean(isEssential),
		tags
	});

	return response.data.data;
}

export async function deleteTransactions(ids: UUID[]): Promise<void> {
	await apiClient.delete('/finance/transactions/bulk', {
		data: {
			transaction_ids: ids
		}
	});
}

async function toggleTransactionFlag(
	transactionId: UUID,
	value: boolean,
	path: 'archive' | 'recurring' | 'essential'
): Promise<void> {
	await apiClient.patch(`/finance/transactions/${transactionId}/${path}`, {
		value
	});
}

export function toggleArchived(transactionId: UUID, value: boolean) {
	return toggleTransactionFlag(transactionId, value, 'archive');
}

export function toggleRecurring(transactionId: UUID, value: boolean) {
	return toggleTransactionFlag(transactionId, value, 'recurring');
}

export function toggleEssential(transactionId: UUID, value: boolean) {
	return toggleTransactionFlag(transactionId, value, 'essential');
}
