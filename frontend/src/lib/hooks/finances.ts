import { createMutation, createQuery } from '@tanstack/svelte-query';

import {
	createTransaction,
	deleteTransactions,
	fetchSummary,
	fetchTransactions,
	toggleArchived,
	toggleEssential,
	toggleRecurring
} from '$lib/api/finances';
import type { Transaction, TransactionSummary, UUID } from '$lib/api/types';

export interface FinanceSummaryOptions {
	userId: string | null;
	enabled?: boolean;
}

export const useFinanceSummaryQuery = (getOptions: () => FinanceSummaryOptions) =>
	createQuery<TransactionSummary>(() => {
		const { userId, enabled = true } = getOptions();
		return {
			queryKey: ['finances', 'summary', userId],
			queryFn: () => fetchSummary(),
			enabled: Boolean(userId) && enabled
		};
	});

export interface FinanceTransactionsOptions {
	userId: string | null;
	search?: string;
	includeArchived?: boolean;
	enabled?: boolean;
}

export const useFinanceTransactionsQuery = (getOptions: () => FinanceTransactionsOptions) =>
	createQuery<Transaction[]>(() => {
		const {
			userId,
			search = '',
			includeArchived = false,
			enabled = true
		} = getOptions();

		const filters = {
			search,
			includeArchived
		};

		return {
			queryKey: ['finances', 'transactions', { userId, ...filters }],
			queryFn: () => fetchTransactions(filters),
			enabled: Boolean(userId) && enabled,
			placeholderData: (previous) => previous ?? []
		};
	});

export const useCreateTransactionMutation = () =>
	createMutation({
		mutationFn: createTransaction
	});

export const useDeleteTransactionsMutation = () =>
	createMutation({
		mutationFn: deleteTransactions
	});

interface TogglePayload {
	id: UUID;
	value: boolean;
}

export const useToggleArchivedMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => toggleArchived(id, value)
	});

export const useToggleRecurringMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => toggleRecurring(id, value)
	});

export const useToggleEssentialMutation = () =>
	createMutation({
		mutationFn: ({ id, value }: TogglePayload) => toggleEssential(id, value)
	});

