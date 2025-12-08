import { onDestroy, onMount } from 'svelte';
import { derived, get, writable, type Readable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type { Transaction, UUID } from '$lib/api/types';
import {
	useFinanceTransactionsQuery,
	useToggleArchivedMutation,
	useToggleEssentialMutation,
	useToggleRecurringMutation
} from '@hooks/finances';
import { getApiErrorMessage, toastError, toastInfo } from '$lib/utils/toast';
import type { ToggleAction } from '../finances.logic';

interface AuthState {
	isAuthenticated: boolean;
	userId: UUID | null;
}

export const numberFormatter = new Intl.NumberFormat('en-US', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2
});

export function createTransactionDetailLogic(transactionIdStore: Readable<string | null>) {
	const queryClient = useQueryClient();

	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);
	const error = writable<string | null>(null);

	const transaction = writable<Transaction | null>(null);

	const transactionsQueryOptions = derived(authStateStore, ($auth) => ({
		userId: $auth.userId,
		enabled: Boolean($auth.isAuthenticated && $auth.userId)
	}));

	const transactionsQuery = useFinanceTransactionsQuery(transactionsQueryOptions);

	const toggleArchivedMutation = useToggleArchivedMutation();
	const toggleRecurringMutation = useToggleRecurringMutation();
	const toggleEssentialMutation = useToggleEssentialMutation();

	const invalidateTransactions = () => {
		const { userId } = get(authStateStore);
		if (!userId) return Promise.resolve();
		return queryClient.invalidateQueries({ queryKey: ['finances', 'transactions'] });
	};

	const transactionsQueryUnsubscribe = transactionsQuery.subscribe((result) => {
		if (result.data) {
			const transactionId = get(transactionIdStore);
			if (transactionId) {
				const found = result.data.find((t) => t.id === transactionId);
				if (found) {
					transaction.set(found);
					error.set(null);
				} else {
					transaction.set(null);
					error.set('Transaction not found.');
				}
			}
		}
		if (result.error) {
			error.set(result.error.message ?? 'Unable to load transaction.');
		}
	});

	const refresh = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated || !auth.userId) {
			toastError('Sign in to refresh your transaction.');
			return;
		}

		try {
			await invalidateTransactions();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to refresh transaction.');
			error.set(message);
			toastError(message);
		}
	};

	const handleToggle = async (action: ToggleAction, value: boolean) => {
		const currentTransaction = get(transaction);
		if (!currentTransaction) {
			toastError('Transaction not loaded.');
			return;
		}

		try {
			if (action === 'archive') {
				await get(toggleArchivedMutation).mutateAsync({ id: currentTransaction.id, value });
			} else if (action === 'recurring') {
				await get(toggleRecurringMutation).mutateAsync({ id: currentTransaction.id, value });
			} else {
				await get(toggleEssentialMutation).mutateAsync({ id: currentTransaction.id, value });
			}

			await invalidateTransactions();

			const actionLabels: Record<ToggleAction, { on: string; off: string }> = {
				archive: { on: 'Transaction archived.', off: 'Transaction restored.' },
				recurring: { on: 'Recurring flag enabled.', off: 'Recurring flag disabled.' },
				essential: { on: 'Essential flag enabled.', off: 'Essential flag disabled.' }
			};
			const label = actionLabels[action][value ? 'on' : 'off'];
			toastInfo(label);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update transaction.');
			error.set(message);
			toastError(message);
		}
	};

	onMount(() => {
		const unsubscribeAuth = authStore.subscribe((state) => {
			const next: AuthState = {
				isAuthenticated: state.isAuthenticated,
				userId: state.user?.id ?? null
			};
			authStateStore.set(next);

			if (!next.isAuthenticated) {
				transaction.set(null);
				error.set(null);
			}
		});

		return () => {
			unsubscribeAuth();
		};
	});

	onDestroy(() => {
		transactionsQueryUnsubscribe();
	});

	return {
		isAuthenticated,
		error,
		transactionQuery: transactionsQuery,
		transaction,
		refresh,
		handleToggle,
		numberFormatter
	};
}

