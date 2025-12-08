<script lang="ts">
	import type { TransactionSummary } from '$lib/api/types';
	import StatCard from '$lib/components/StatCard.svelte';

	export let summary: TransactionSummary | null = null;
	export let isLoading = false;
	export let error: unknown = null;
	export let numberFormatter: Intl.NumberFormat;
</script>

{#if isLoading}
	<div class="grid gap-4 md:grid-cols-3">
		<div
			class="rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-5 backdrop-blur-sm"
		>
			<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Loading…</p>
		</div>
		<div
			class="rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-5 backdrop-blur-sm"
		>
			<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Loading…</p>
		</div>
		<div
			class="rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-5 backdrop-blur-sm"
		>
			<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Loading…</p>
		</div>
	</div>
{:else if error}
	<div
		class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-200"
	>
		Unable to load summary. Please try again.
	</div>
{:else if summary}
	<div class="grid gap-4 md:grid-cols-3">
		<StatCard
			label="Income"
			value={`${numberFormatter.format(summary.income_total)} ${summary.base_currency}`}
			accentColor="emerald"
		/>
		<StatCard
			label="Expenses"
			value={`${numberFormatter.format(summary.expense_total)} ${summary.base_currency}`}
			accentColor="red"
		/>
		<StatCard
			label="Savings Allocation"
			value={`${numberFormatter.format(summary.savings_allocation)} ${summary.base_currency}`}
			accentColor="blue"
		/>
	</div>
{/if}

