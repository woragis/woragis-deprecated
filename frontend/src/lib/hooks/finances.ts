import { createMutation, createQuery } from '@tanstack/svelte-query';

import { financesApi } from '$lib/api/finances';
import type { Transaction, TransactionSummary, UUID } from '$lib/api/types';

export interface FinanceSummaryOptions {
	enabled?: boolean;
}

export const useFinanceSummaryQuery = (
	userId: string | null,
	options: FinanceSummaryOptions = {}
) =>
	createQuery<TransactionSummary>({
		queryKey: ['finances', 'summary', userId],
		queryFn: () => financesApi.fetchSummary(),
		enabled: Boolean(userId) && (options.enabled ?? true)
	});

export interface FinanceTransactionsOptions {
	search?: string;
	includeArchived?: boolean;
	enabled?: boolean;
}

export const useFinanceTransactionsQuery = (
	userId: string | null,
	options: FinanceTransactionsOptions = {}
) => {
	const filters = {
		search: options.search ?? '',
		includeArchived: options.includeArchived ?? false
	};

	return createQuery<Transaction[]>({
		queryKey: ['finances', 'transactions', { userId, ...filters }],
		queryFn: () => financesApi.fetchTransactions(filters),
		enabled: Boolean(userId) && (options.enabled ?? true),
		placeholderData: (previous) => previous ?? []
	});
};

export const useCreateTransactionMutation = () =>
	createMutation({
		mutationFn: financesApi.createTransaction
	});

export const useDeleteTransactionsMutation = () =>
	createMutation({
		mutationFn: financesApi.deleteTransactions
	});

interface TogglePayload {
	id: UUID;
	value: boolean;
}

export const useToggleArchivedMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => financesApi.toggleArchived(id, value)
	});

export const useToggleRecurringMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => financesApi.toggleRecurring(id, value)
	});

export const useToggleEssentialMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => financesApi.toggleEssential(id, value)
	});

