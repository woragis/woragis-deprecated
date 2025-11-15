import { SvelteSet } from 'svelte/reactivity';
import { derived, get, writable } from 'svelte/store';
import { onDestroy, onMount } from 'svelte';
import { useQueryClient } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type { Transaction, TransactionSummary, UUID } from '$lib/api/types';
import {
	useCreateTransactionMutation,
	useDeleteTransactionsMutation,
	useFinanceSummaryQuery,
	useFinanceTransactionsQuery,
	useToggleArchivedMutation,
	useToggleEssentialMutation,
	useToggleRecurringMutation
} from '@hooks/finances';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export type ToggleAction = 'archive' | 'recurring' | 'essential';

export interface TransactionFormState {
	type: 'income' | 'expense';
	category: string;
	description: string;
	amount: number;
	currency: string;
	baseCurrency: string;
	occurredAt: string;
	isRecurring: boolean;
	isEssential: boolean;
	tags: string;
}

const defaultFormState = (): TransactionFormState => ({
	type: 'income',
	category: '',
	description: '',
	amount: 0,
	currency: 'USD',
	baseCurrency: 'USD',
	occurredAt: new Date().toISOString().slice(0, 16),
	isRecurring: false,
	isEssential: false,
	tags: ''
});

interface AuthState {
	isAuthenticated: boolean;
	userId: UUID | null;
}

export const numberFormatter = new Intl.NumberFormat('en-US', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2
});

export function createFinancesLogic() {
	const search = writable('');
	const includeArchived = writable(false);
	const formState = writable<TransactionFormState>(defaultFormState());
	const actionError = writable<string | null>(null);
	const isMutating = writable(false);
	const selection = writable<Set<UUID>>(new SvelteSet<UUID>());
	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const filterStore = writable({ search: '', includeArchived: false });

	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);

	const queryClient = useQueryClient();

	const summaryOptionsStore = derived(authStateStore, ($auth) => ({
		userId: $auth.userId,
		enabled: Boolean($auth.isAuthenticated && $auth.userId)
	}));

	const transactionsOptionsStore = derived(
		[authStateStore, filterStore],
		([$auth, $filters]) => ({
			userId: $auth.userId,
			search: $filters.search,
			includeArchived: $filters.includeArchived,
			enabled: Boolean($auth.isAuthenticated && $auth.userId)
		})
	);

	const summaryQuery = useFinanceSummaryQuery(summaryOptionsStore);
	const transactionsQuery = useFinanceTransactionsQuery(transactionsOptionsStore);

	const createTransactionMutation = useCreateTransactionMutation();
	const deleteTransactionsMutation = useDeleteTransactionsMutation();
	const toggleArchivedMutation = useToggleArchivedMutation();
	const toggleRecurringMutation = useToggleRecurringMutation();
	const toggleEssentialMutation = useToggleEssentialMutation();

	const resetSelection = () => {
		selection.set(new SvelteSet<UUID>());
	};

	const transactionsUnsubscribe = transactionsQuery.subscribe(() => {
		resetSelection();
	});

	const resetStateForUnauthenticated = () => {
		resetSelection();
		actionError.set(null);
		formState.set(defaultFormState());
		filterStore.set({ search: '', includeArchived: false });
		search.set('');
		includeArchived.set(false);
	};

	onMount(() => {
		const unsubscribeAuth = authStore.subscribe((state) => {
			const next: AuthState = {
				isAuthenticated: state.isAuthenticated,
				userId: state.user?.id ?? null
			};
			authStateStore.set(next);

			if (!next.isAuthenticated || !next.userId) {
				resetStateForUnauthenticated();
			}
		});

		return () => {
			unsubscribeAuth();
		};
	});

	onDestroy(() => {
		transactionsUnsubscribe();
	});

	const invalidateFinances = async () => {
		const { userId } = get(authStateStore);
		if (!userId) return;

		await Promise.all([
			queryClient.invalidateQueries({ queryKey: ['finances', 'summary', userId] }),
			queryClient.invalidateQueries({ queryKey: ['finances', 'transactions'] })
		]);
	};

	const refresh = () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated || !auth.userId) {
			toastError('Sign in to refresh your finances.');
			return;
		}

		filterStore.set({
			search: get(search).trim(),
			includeArchived: get(includeArchived)
		});
	};

	const updateFormField = <K extends keyof TransactionFormState>(
		field: K,
		value: TransactionFormState[K]
	) => {
		formState.update((current) => ({
			...current,
			[field]: value
		}));
	};

	const handleCreateTransaction = async () => {
		const auth = get(authStateStore);
		if (!auth.userId) {
			const message = 'You must be signed in to create transactions.';
			actionError.set(message);
			toastError(message);
			return;
		}

		const currentForm = get(formState);

		if (!currentForm.category.trim()) {
			const message = 'Category is required.';
			actionError.set(message);
			toastError(message);
			return;
		}

		isMutating.set(true);
		actionError.set(null);

		try {
			await get(createTransactionMutation).mutateAsync({
				type: currentForm.type,
				category: currentForm.category,
				description: currentForm.description,
				amount: currentForm.amount,
				currency: currentForm.currency,
				baseCurrency: currentForm.baseCurrency,
				occurredAt: new Date(currentForm.occurredAt).toISOString(),
				isRecurring: currentForm.isRecurring,
				isEssential: currentForm.isEssential,
				tags: currentForm.tags ? currentForm.tags.split(',').map((tag) => tag.trim()) : []
			});

			formState.set(defaultFormState());
			await invalidateFinances();
			toastSuccess('Transaction recorded.');
		} catch (error) {
			const message = getApiErrorMessage(error, 'Unable to create transaction.');
			actionError.set(message);
			toastError(message);
		} finally {
			isMutating.set(false);
		}
	};

	const handleToggle = async (transaction: Transaction, action: ToggleAction, value: boolean) => {
		const auth = get(authStateStore);
		if (!auth.userId) {
			const message = 'You must be signed in to update transactions.';
			actionError.set(message);
			toastError(message);
			return;
		}

		try {
			if (action === 'archive') {
				await get(toggleArchivedMutation).mutateAsync({ id: transaction.id, value });
			} else if (action === 'recurring') {
				await get(toggleRecurringMutation).mutateAsync({ id: transaction.id, value });
			} else {
				await get(toggleEssentialMutation).mutateAsync({ id: transaction.id, value });
			}

			await invalidateFinances();

			const actionLabels: Record<ToggleAction, { on: string; off: string }> = {
				archive: { on: 'Transaction archived.', off: 'Transaction restored.' },
				recurring: { on: 'Recurring flag enabled.', off: 'Recurring flag disabled.' },
				essential: { on: 'Essential flag enabled.', off: 'Essential flag disabled.' }
			};
			const label = actionLabels[action][value ? 'on' : 'off'];
			toastInfo(label);
		} catch (error) {
			const message = getApiErrorMessage(error, 'Unable to update transaction.');
			actionError.set(message);
			toastError(message);
		}
	};

	const handleBulkDelete = async () => {
		const auth = get(authStateStore);
		const currentSelection = Array.from(get(selection));

		if (currentSelection.length === 0) return;

		if (!auth.userId) {
			const message = 'You must be signed in to delete transactions.';
			actionError.set(message);
			toastError(message);
			return;
		}

		isMutating.set(true);
		actionError.set(null);

		try {
			await get(deleteTransactionsMutation).mutateAsync(currentSelection);
			await invalidateFinances();
			toastSuccess('Selected transactions deleted.');
		} catch (error) {
			const message = getApiErrorMessage(error, 'Unable to delete transactions.');
			actionError.set(message);
			toastError(message);
		} finally {
			isMutating.set(false);
		}
	};

	const toggleSelection = (id: UUID, checked: boolean) => {
		selection.update((current) => {
			const next = new SvelteSet(current);
			if (checked) {
				next.add(id);
			} else {
				next.delete(id);
			}
			return next;
		});
	};

	return {
		search,
		includeArchived,
		formState,
		actionError,
		isMutating,
		selection,
		isAuthenticated,
		summaryQuery,
		transactionsQuery,
		refresh,
		updateFormField,
		handleCreateTransaction,
		handleToggle,
		handleBulkDelete,
		toggleSelection
	};
}

