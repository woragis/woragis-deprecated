<script lang="ts">
	import { createFinancesLogic, numberFormatter } from './finances.logic';
	import type { Transaction, TransactionSummary } from '$lib/api/types';

	import FinancesHero from './_components/FinancesHero.svelte';
	import ActionError from './_components/ActionError.svelte';
	import AuthNotice from './_components/AuthNotice.svelte';
	import FinanceFilters from './_components/FinanceFilters.svelte';
	import FinanceSummary from './_components/FinanceSummary.svelte';
	import TransactionForm from './_components/TransactionForm.svelte';
	import TransactionsTable from './_components/TransactionsTable.svelte';
	import CollapsibleSection from '$lib/components/CollapsibleSection.svelte';

	let showTransactionForm = false;

	function toggleTransactionForm() {
		showTransactionForm = !showTransactionForm;
	}

	const {
		search,
		includeArchived,
		formState,
		actionError,
		isAuthenticated,
		summaryQuery,
		transactionsQuery,
		selection,
		isMutating,
		refresh,
		updateFormField,
		handleCreateTransaction,
		handleToggle,
		handleBulkDelete,
		toggleSelection
	} = createFinancesLogic();
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<FinancesHero isAuthenticated={$isAuthenticated} onRefresh={refresh} />

	<!-- Main Content -->
	<div class="mx-auto max-w-7xl px-6 py-8 lg:px-8">
		{#if $actionError}
			<div class="mb-6">
				<ActionError message={$actionError} />
			</div>
		{/if}

		{#if !$isAuthenticated}
			<div class="mx-auto max-w-2xl">
				<AuthNotice />
			</div>
		{:else}
			<div class="space-y-6">
				<!-- Filters -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<FinanceFilters
						search={$search}
						includeArchived={$includeArchived}
						isAuthenticated={$isAuthenticated}
						onSearchChange={(value) => search.set(value)}
						onIncludeArchivedChange={(value) => includeArchived.set(value)}
						onRefresh={refresh}
					/>
				</div>

				<!-- Summary Stats -->
				<FinanceSummary
					summary={(($summaryQuery.data as TransactionSummary | undefined) ?? null)}
					isLoading={$summaryQuery.isLoading}
					error={$summaryQuery.error}
					numberFormatter={numberFormatter}
				/>

				<!-- Action Button -->
				<div class="flex justify-end">
					<button
						class="group relative overflow-hidden rounded-lg bg-gradient-to-r from-emerald-600 to-teal-600 px-6 py-2.5 text-sm font-semibold text-white shadow-lg shadow-emerald-500/25 transition-all hover:scale-105 hover:shadow-xl hover:shadow-emerald-500/40"
						on:click={toggleTransactionForm}
					>
						<span class="relative z-10 flex items-center gap-2">
							<svg
								class="h-4 w-4 transition-transform group-hover:rotate-90"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 4v16m8-8H4"
								></path>
							</svg>
							{showTransactionForm ? 'Cancel' : 'Record Transaction'}
						</span>
					</button>
				</div>

				<!-- Transaction Form (Collapsible) -->
				<CollapsibleSection open={showTransactionForm}>
					<TransactionForm
						formState={$formState}
						isMutating={$isMutating}
						onFieldChange={updateFormField}
						onSubmit={handleCreateTransaction}
					/>
				</CollapsibleSection>

				<!-- Transactions Table -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<TransactionsTable
						transactions={(($transactionsQuery.data as Transaction[] | undefined) ?? [])}
						isFetching={$transactionsQuery.isFetching}
						error={$transactionsQuery.error}
						selection={$selection}
						isMutating={$isMutating}
						numberFormatter={numberFormatter}
						onToggleSelection={toggleSelection}
						onToggle={handleToggle}
						onBulkDelete={handleBulkDelete}
					/>
				</div>
			</div>
		{/if}
	</div>
</div>
