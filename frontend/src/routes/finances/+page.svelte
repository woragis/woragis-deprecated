<script lang="ts">
import { onDestroy, onMount } from 'svelte';
import { SvelteSet } from 'svelte/reactivity';
import { derived } from 'svelte/store';
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

interface TransactionFormState {
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

const defaultForm = (): TransactionFormState => ({
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

let currentUserId: UUID | null = null;
let isAuthenticated = false;
let isMutating = false;
let actionError: string | null = null;
let search = '';
let includeArchived = false;
let formState = defaultForm();
let selection = new SvelteSet<UUID>();
let appliedFilters = { search: '', includeArchived: false };

const queryClient = useQueryClient();

const summaryQuery = useFinanceSummaryQuery(() => ({
	userId: currentUserId,
	enabled: Boolean(isAuthenticated && currentUserId)
}));

const transactionsQuery = useFinanceTransactionsQuery(() => ({
	userId: currentUserId,
	search: appliedFilters.search,
	includeArchived: appliedFilters.includeArchived,
	enabled: Boolean(isAuthenticated && currentUserId)
}));

const summaryData = derived(summaryQuery, ($query) => ($query.data ?? null) as TransactionSummary | null);
const transactionsData = derived(transactionsQuery, ($query) => ($query.data ?? []) as Transaction[]);

const numberFormatter = new Intl.NumberFormat('en-US', {
	minimumFractionDigits: 2,
	maximumFractionDigits: 2
});

const createTransactionMutation = useCreateTransactionMutation();
const deleteTransactionsMutation = useDeleteTransactionsMutation();
const toggleArchivedMutation = useToggleArchivedMutation();
const toggleRecurringMutation = useToggleRecurringMutation();
const toggleEssentialMutation = useToggleEssentialMutation();

type ToggleAction = 'archive' | 'recurring' | 'essential';

function resetStateForUnauthenticated() {
	selection = new SvelteSet<UUID>();
	actionError = null;
}

const invalidateFinances = async () => {
	if (!currentUserId) return;

	await Promise.all([
		queryClient.invalidateQueries({ queryKey: ['finances', 'summary', currentUserId] }),
		queryClient.invalidateQueries({ queryKey: ['finances', 'transactions'] })
	]);
};

onMount(() => {
	const unsubscribeAuth = authStore.subscribe((state) => {
		isAuthenticated = state.isAuthenticated;
		currentUserId = state.user?.id ?? null;

		if (!isAuthenticated || !currentUserId) {
			resetStateForUnauthenticated();
		}
	});

	const unsubscribeTransactions = transactionsData.subscribe(() => {
		selection = new SvelteSet<UUID>();
	});

	return () => {
		unsubscribeAuth();
		unsubscribeTransactions();
	};
});

onDestroy(() => {
	selection = new SvelteSet<UUID>();
});

function refresh() {
	if (!isAuthenticated || !currentUserId) {
		toastError('Sign in to refresh your finances.');
		return;
	}

	appliedFilters = {
		search: search.trim(),
		includeArchived
	};
}

async function handleCreateTransaction() {
	if (!currentUserId) {
		actionError = 'You must be signed in to create transactions.';
		toastError(actionError);
		return;
	}

	if (!formState.category.trim()) {
		actionError = 'Category is required.';
		toastError(actionError);
		return;
	}

	isMutating = true;
	actionError = null;
	try {
		await $createTransactionMutation.mutateAsync({
			type: formState.type,
			category: formState.category,
			description: formState.description,
			amount: formState.amount,
			currency: formState.currency,
			baseCurrency: formState.baseCurrency,
			occurredAt: new Date(formState.occurredAt).toISOString(),
			isRecurring: formState.isRecurring,
			isEssential: formState.isEssential,
			tags: formState.tags ? formState.tags.split(',').map((tag) => tag.trim()) : []
		});
		formState = defaultForm();
		await invalidateFinances();
		toastSuccess('Transaction recorded.');
	} catch (err) {
		actionError = getApiErrorMessage(err, 'Unable to create transaction.');
		toastError(actionError);
	} finally {
		isMutating = false;
	}
}

async function handleToggle(transaction: Transaction, action: ToggleAction, value: boolean) {
	if (!currentUserId) {
		actionError = 'You must be signed in to update transactions.';
		toastError(actionError);
		return;
	}

	try {
		if (action === 'archive') {
			await $toggleArchivedMutation.mutateAsync({ id: transaction.id, value });
		} else if (action === 'recurring') {
			await $toggleRecurringMutation.mutateAsync({ id: transaction.id, value });
		} else {
			await $toggleEssentialMutation.mutateAsync({ id: transaction.id, value });
		}
		await invalidateFinances();
		const actionLabels: Record<ToggleAction, { on: string; off: string }> = {
			archive: { on: 'Transaction archived.', off: 'Transaction restored.' },
			recurring: { on: 'Recurring flag enabled.', off: 'Recurring flag disabled.' },
			essential: { on: 'Essential flag enabled.', off: 'Essential flag disabled.' }
		};
		const label = actionLabels[action][value ? 'on' : 'off'];
		toastInfo(label);
	} catch (err) {
		actionError = getApiErrorMessage(err, 'Unable to update transaction.');
		toastError(actionError);
	}
}

async function handleBulkDelete() {
	if (selection.size === 0) return;
	if (!currentUserId) {
		actionError = 'You must be signed in to delete transactions.';
		toastError(actionError);
		return;
	}

	isMutating = true;
	actionError = null;
	try {
		await $deleteTransactionsMutation.mutateAsync(Array.from(selection));
		await invalidateFinances();
		toastSuccess('Selected transactions deleted.');
	} catch (err) {
		actionError = getApiErrorMessage(err, 'Unable to delete transactions.');
		toastError(actionError);
	} finally {
		isMutating = false;
	}
}

function toggleSelection(id: UUID, checked: boolean) {
	const next = new SvelteSet<UUID>(selection);
	if (checked) {
		next.add(id);
	} else {
		next.delete(id);
	}
	selection = next;
}
</script>

<section class="space-y-6">
	<div class="flex flex-col gap-4 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
		<h2 class="text-lg font-semibold text-slate-100">Finance Workspace</h2>
		<div class="grid gap-3 md:grid-cols-2">
			<label class="flex flex-col gap-1 text-xs font-medium tracking-wide text-slate-400 uppercase">
				Search
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					type="search"
					bind:value={search}
					placeholder="category, description, tags"
				/>
			</label>
		</div>
		<div class="flex flex-wrap items-center gap-3 text-xs">
			<label class="flex items-center gap-2">
				<input type="checkbox" bind:checked={includeArchived} />
				<span>Include archived</span>
			</label>
			<button
				class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
				on:click={refresh}
				disabled={!isAuthenticated}
			>
				Refresh
			</button>
		</div>
	</div>

	{#if actionError}
		<p class="rounded border border-red-500/50 bg-red-500/10 px-3 py-2 text-sm text-red-200">
			{actionError}
		</p>
	{/if}

	{#if !isAuthenticated}
		<div
			class="rounded border border-slate-800 bg-slate-900/60 p-6 text-center text-sm text-slate-300"
		>
			<p class="mb-2 text-base font-semibold text-slate-100">Sign in to manage your finances</p>
			<p>You need to be authenticated before you can view or record transactions.</p>
		</div>
	{:else}
		{#if $summaryQuery.isLoading}
			<p class="text-sm text-slate-400">Loading summary…</p>
		{:else if $summaryQuery.error}
			<p class="rounded border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
				Unable to load summary. Please try again.
			</p>
		{:else if $summaryData}
			<section class="grid gap-4 md:grid-cols-3">
				<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
					<h3 class="text-xs tracking-wide text-slate-400 uppercase">Income</h3>
					<p class="mt-2 text-2xl font-semibold text-emerald-400">
						{numberFormatter.format($summaryData.income_total)}
						{$summaryData.base_currency}
					</p>
				</div>
				<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
					<h3 class="text-xs tracking-wide text-slate-400 uppercase">Expenses</h3>
					<p class="mt-2 text-2xl font-semibold text-rose-400">
						{numberFormatter.format($summaryData.expense_total)}
						{$summaryData.base_currency}
					</p>
				</div>
				<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
					<h3 class="text-xs tracking-wide text-slate-400 uppercase">Savings Allocation</h3>
					<p class="mt-2 text-2xl font-semibold text-sky-400">
						{numberFormatter.format($summaryData.savings_allocation)}
						{$summaryData.base_currency}
					</p>
				</div>
			</section>
		{/if}

		<section class="grid gap-6 lg:grid-cols-[1.1fr_2fr]">
			<div class="space-y-4 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
				<h3 class="text-sm font-semibold text-slate-100">Record Transaction</h3>
				<form class="space-y-3" on:submit|preventDefault={handleCreateTransaction}>
					<div class="grid gap-3 md:grid-cols-2">
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Type
							<select
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.type}
							>
								<option value="income">Income</option>
								<option value="expense">Expense</option>
							</select>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Category
							<input
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.category}
								required
							/>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Amount
							<input
								type="number"
								step="0.01"
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.amount}
								required
							/>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Currency
							<input
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.currency}
								maxlength="3"
							/>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Base Currency
							<input
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.baseCurrency}
								maxlength="3"
							/>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-400">
							Occurred At
							<input
								type="datetime-local"
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={formState.occurredAt}
							/>
						</label>
					</div>
					<label class="flex flex-col gap-1 text-xs text-slate-400">
						Description
						<textarea
							class="min-h-[80px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							bind:value={formState.description}
						></textarea>
					</label>
					<label class="flex flex-col gap-1 text-xs text-slate-400">
						Tags (comma separated)
						<input
							class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							bind:value={formState.tags}
						/>
					</label>
					<div class="flex flex-wrap items-center gap-4 text-xs text-slate-300">
						<label class="flex items-center gap-2">
							<input type="checkbox" bind:checked={formState.isRecurring} />
							<span>Recurring</span>
						</label>
						<label class="flex items-center gap-2">
							<input type="checkbox" bind:checked={formState.isEssential} />
							<span>Essential</span>
						</label>
					</div>
					<button
						type="submit"
						class="w-full rounded bg-emerald-500 px-3 py-2 text-xs font-semibold text-slate-950 disabled:opacity-50"
						disabled={isMutating}
					>
						Create Transaction
					</button>
				</form>
			</div>

			<div class="space-y-3 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
				<div class="flex items-center justify-between">
					<h3 class="text-sm font-semibold text-slate-100">Transactions ({$transactionsData.length})</h3>
					<button
						class="rounded bg-rose-500/80 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
						disabled={selection.size === 0 || isMutating}
						on:click={handleBulkDelete}
					>
						Delete Selected
					</button>
				</div>

				{#if $transactionsQuery.isFetching}
					<p class="text-sm text-slate-400">Loading…</p>
				{:else if $transactionsQuery.error}
					<p class="rounded border border-red-500/40 bg-red-500/10 px-3 py-2 text-sm text-red-200">
						Unable to load transactions. Please try again.
					</p>
				{:else}
					{#if $transactionsData.length === 0}
						<p class="text-sm text-slate-400">
							No transactions found. Try adjusting filters or create your first entry.
						</p>
					{:else}
						<div class="overflow-x-auto text-xs">
							<table class="min-w-full border-separate border-spacing-y-2">
								<thead class="text-[10px] tracking-wide text-slate-400 uppercase">
									<tr>
										<th></th>
										<th class="text-left">Type</th>
										<th class="text-left">Category</th>
										<th class="text-left">Amount</th>
										<th class="text-left">Currency</th>
										<th class="text-left">Occurred</th>
										<th class="text-left">Flags</th>
										<th class="text-left">Actions</th>
									</tr>
								</thead>
								<tbody>
									{#each $transactionsData as transaction (transaction.id)}
										<tr class="rounded border border-slate-800 bg-slate-950/40">
											<td class="px-2">
												<input
													type="checkbox"
													checked={selection.has(transaction.id)}
													on:change={(event) =>
														toggleSelection(
															transaction.id,
															(event.target as HTMLInputElement).checked
														)}
												/>
											</td>
											<td
												class="px-2 py-2 font-semibold {transaction.type === 'income'
													? 'text-emerald-400'
													: 'text-rose-400'}"
											>
												{transaction.type}
											</td>
											<td class="px-2 py-2 text-slate-200">{transaction.category}</td>
											<td class="px-2 py-2">{numberFormatter.format(transaction.amount)}</td>
											<td class="px-2 py-2">{transaction.currency}</td>
											<td class="px-2 py-2 text-slate-300"
												>{new Date(transaction.occurred_at).toLocaleString()}</td
											>
											<td class="px-2 py-2 text-[11px] text-slate-300">
												{#if transaction.is_recurring}
													<span class="mr-2 rounded bg-indigo-500/20 px-2 py-1 text-indigo-200"
														>Recurring</span
													>
												{/if}
												{#if transaction.is_essential}
													<span class="rounded bg-amber-500/20 px-2 py-1 text-amber-200"
														>Essential</span
													>
												{/if}
												{#if transaction.is_archived}
													<span class="ml-2 rounded bg-slate-700/40 px-2 py-1 text-slate-300"
														>Archived</span
													>
												{/if}
											</td>
											<td class="flex flex-wrap gap-2 px-2 py-2">
												<button
													class="rounded bg-slate-800 px-2 py-1"
													on:click={() =>
														handleToggle(transaction, 'archive', !transaction.is_archived)}
												>
													{transaction.is_archived ? 'Unarchive' : 'Archive'}
												</button>
												<button
													class="rounded bg-slate-800 px-2 py-1"
													on:click={() =>
														handleToggle(transaction, 'recurring', !transaction.is_recurring)}
												>
													Toggle Recurring
												</button>
												<button
													class="rounded bg-slate-800 px-2 py-1"
													on:click={() =>
														handleToggle(transaction, 'essential', !transaction.is_essential)}
												>
													Toggle Essential
												</button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				{/if}
			</div>
		</section>
	{/if}
</section>
