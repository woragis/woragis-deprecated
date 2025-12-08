<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Transaction, UUID } from '$lib/api/types';
	import type { ToggleAction } from '../finances.logic';
	import TransactionCard from './TransactionCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	export let transactions: Transaction[] = [];
	export let isFetching = false;
	export let error: unknown = null;
	export let selection: Set<UUID>;
	export let isMutating = false;
	export let numberFormatter: Intl.NumberFormat;
	export let onToggleSelection: (id: UUID, checked: boolean) => void;
	export let onToggle: (transaction: Transaction, action: ToggleAction, value: boolean) => void;
	export let onBulkDelete: () => void;

	function handleCardClick(transaction: Transaction) {
		goto(`/personal/finances/${transaction.id}`);
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<div>
			<h3 class="text-xl font-semibold text-slate-100">Transactions</h3>
			<p class="mt-1 text-sm text-slate-400">{transactions.length} total</p>
		</div>
		{#if selection.size > 0}
			<button
				class="rounded-lg bg-rose-600 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-rose-500/25 transition-all hover:scale-105 hover:shadow-xl hover:shadow-rose-500/40 disabled:opacity-50"
				disabled={isMutating}
				onclick={onBulkDelete}
			>
				Delete Selected ({selection.size})
			</button>
		{/if}
	</div>

	{#if isFetching}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each Array(6) as _}
				<div
					class="h-48 animate-pulse rounded-xl border border-slate-800/50 bg-slate-900/40 backdrop-blur-sm"
				></div>
			{/each}
		</div>
	{:else if error}
		<div class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200">
			Unable to load transactions. Please try again.
		</div>
	{:else if transactions.length === 0}
		<EmptyState
			title="No transactions found"
			description="Try adjusting filters or create your first entry"
			icon='<svg class="mx-auto h-12 w-12 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>'
		/>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each transactions as transaction (transaction.id)}
				<TransactionCard
					{transaction}
					isSelected={selection.has(transaction.id)}
					{numberFormatter}
					{onToggleSelection}
					{onToggle}
					onClick={handleCardClick}
				/>
			{/each}
		</div>
	{/if}
</div>

