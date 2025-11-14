<script lang="ts">
	import type { TransactionSummary } from '$lib/api/types';

	export let summary: TransactionSummary | null = null;
	export let isLoading = false;
	export let error: unknown = null;
	export let numberFormatter: Intl.NumberFormat;
</script>

{#if isLoading}
	<p class="text-sm text-slate-400">Loading summary…</p>
{:else if error}
	<p class="rounded border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
		Unable to load summary. Please try again.
	</p>
{:else if summary}
	<section class="grid gap-4 md:grid-cols-3">
		<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
			<h3 class="text-xs tracking-wide text-slate-400 uppercase">Income</h3>
			<p class="mt-2 text-2xl font-semibold text-emerald-400">
				{numberFormatter.format(summary.income_total)}
				{summary.base_currency}
			</p>
		</div>
		<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
			<h3 class="text-xs tracking-wide text-slate-400 uppercase">Expenses</h3>
			<p class="mt-2 text-2xl font-semibold text-rose-400">
				{numberFormatter.format(summary.expense_total)}
				{summary.base_currency}
			</p>
		</div>
		<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
			<h3 class="text-xs tracking-wide text-slate-400 uppercase">Savings Allocation</h3>
			<p class="mt-2 text-2xl font-semibold text-sky-400">
				{numberFormatter.format(summary.savings_allocation)}
				{summary.base_currency}
			</p>
		</div>
	</section>
{/if}

