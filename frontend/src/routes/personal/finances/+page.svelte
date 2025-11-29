<script lang="ts">
	import { createFinancesLogic, numberFormatter } from './finances.logic';
	import type { Transaction, TransactionSummary } from '$lib/api/types';

	import ActionError from './_components/ActionError.svelte';
	import AuthNotice from './_components/AuthNotice.svelte';
	import FinanceFilters from './_components/FinanceFilters.svelte';
	import FinanceSummary from './_components/FinanceSummary.svelte';
	import TransactionForm from './_components/TransactionForm.svelte';
	import TransactionsTable from './_components/TransactionsTable.svelte';

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

<section class="space-y-6">
	<FinanceFilters
		search={$search}
		includeArchived={$includeArchived}
		isAuthenticated={$isAuthenticated}
		onSearchChange={(value) => search.set(value)}
		onIncludeArchivedChange={(value) => includeArchived.set(value)}
		onRefresh={refresh}
	/>

	{#if $actionError}
		<ActionError message={$actionError} />
	{/if}

	{#if !$isAuthenticated}
		<AuthNotice />
	{:else}
		<FinanceSummary
			summary={(($summaryQuery.data as TransactionSummary | undefined) ?? null)}
			isLoading={$summaryQuery.isLoading}
			error={$summaryQuery.error}
			numberFormatter={numberFormatter}
		/>

		<section class="grid gap-6 lg:grid-cols-[1.1fr_2fr]">
			<TransactionForm
				formState={$formState}
				isMutating={$isMutating}
				onFieldChange={updateFormField}
				onSubmit={handleCreateTransaction}
			/>

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
		</section>
	{/if}
</section>
