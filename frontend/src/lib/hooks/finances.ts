import { createMutation, createQuery } from '@tanstack/svelte-query';
import { derived, type Readable } from 'svelte/store';

import { financesApi } from '$lib/api/finances';
import type { Transaction, TransactionSummary, UUID } from '$lib/api/types';

type MaybeReadable<T> = T | Readable<T>;

const isReadable = <T>(value: MaybeReadable<T>): value is Readable<T> =>
	typeof value === 'object' && value !== null && typeof (value as Readable<T>).subscribe === 'function';

interface FinanceSummaryQueryArgs {
	userId: string | null;
	enabled?: boolean;
}

const createSummaryQueryOptions = (args: FinanceSummaryQueryArgs) => ({
	queryKey: ['finances', 'summary', args.userId],
	queryFn: () => financesApi.fetchSummary(),
	enabled: Boolean(args.userId) && (args.enabled ?? true)
});

export const useFinanceSummaryQuery = (source: MaybeReadable<FinanceSummaryQueryArgs>) => {
	if (isReadable(source)) {
		const optionsStore = derived(source, ($args) => createSummaryQueryOptions($args));
		return createQuery<TransactionSummary>(optionsStore);
	}

	return createQuery<TransactionSummary>(createSummaryQueryOptions(source));
};

interface FinanceTransactionsQueryArgs {
	userId: string | null;
	search?: string;
	includeArchived?: boolean;
	enabled?: boolean;
}

const createTransactionsQueryOptions = (args: FinanceTransactionsQueryArgs) => {
	const filters = {
		search: args.search ?? '',
		includeArchived: args.includeArchived ?? false
	};

	return {
		queryKey: [
			'finances',
			'transactions',
			args.userId,
			filters.search,
			filters.includeArchived
		],
		queryFn: () => financesApi.fetchTransactions(filters),
		enabled: Boolean(args.userId) && (args.enabled ?? true),
		placeholderData: (previous: Transaction[] | undefined) => previous ?? []
	};
};

export const useFinanceTransactionsQuery = (source: MaybeReadable<FinanceTransactionsQueryArgs>) => {
	if (isReadable(source)) {
		const optionsStore = derived(source, ($args) => createTransactionsQueryOptions($args));
		return createQuery<Transaction[]>(optionsStore);
	}

	return createQuery<Transaction[]>(createTransactionsQueryOptions(source));
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

