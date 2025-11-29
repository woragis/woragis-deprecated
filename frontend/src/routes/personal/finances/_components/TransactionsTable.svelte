<script lang="ts">
	import type { Transaction, UUID } from '$lib/api/types';
	import type { ToggleAction } from '../finances.logic';

	export let transactions: Transaction[] = [];
	export let isFetching = false;
	export let error: unknown = null;
	export let selection: Set<UUID>;
	export let isMutating = false;
	export let numberFormatter: Intl.NumberFormat;
	export let onToggleSelection: (id: UUID, checked: boolean) => void;
	export let onToggle: (transaction: Transaction, action: ToggleAction, value: boolean) => void;
	export let onBulkDelete: () => void;
</script>

<div class="space-y-3 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
	<div class="flex items-center justify-between">
		<h3 class="text-sm font-semibold text-slate-100">Transactions ({transactions.length})</h3>
		<button
			class="rounded bg-rose-500/80 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
			disabled={selection.size === 0 || isMutating}
			on:click={onBulkDelete}
		>
			Delete Selected
		</button>
	</div>

	{#if isFetching}
		<p class="text-sm text-slate-400">Loading…</p>
	{:else if error}
		<p class="rounded border border-red-500/40 bg-red-500/10 px-3 py-2 text-sm text-red-200">
			Unable to load transactions. Please try again.
		</p>
	{:else if transactions.length === 0}
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
					{#each transactions as transaction (transaction.id)}
						<tr class="rounded border border-slate-800 bg-slate-950/40">
							<td class="px-2">
								<input
									type="checkbox"
									checked={selection.has(transaction.id)}
									on:change={(event) =>
										onToggleSelection(transaction.id, (event.target as HTMLInputElement).checked)}
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
									<span class="mr-2 rounded bg-indigo-500/20 px-2 py-1 text-indigo-200">Recurring</span>
								{/if}
								{#if transaction.is_essential}
									<span class="rounded bg-amber-500/20 px-2 py-1 text-amber-200">Essential</span>
								{/if}
								{#if transaction.is_archived}
									<span class="ml-2 rounded bg-slate-700/40 px-2 py-1 text-slate-300">Archived</span>
								{/if}
							</td>
							<td class="flex flex-wrap gap-2 px-2 py-2">
								<button
									class="rounded bg-slate-800 px-2 py-1"
									on:click={() =>
										onToggle(transaction, 'archive', !transaction.is_archived)}
								>
									{transaction.is_archived ? 'Unarchive' : 'Archive'}
								</button>
								<button
									class="rounded bg-slate-800 px-2 py-1"
									on:click={() =>
										onToggle(transaction, 'recurring', !transaction.is_recurring)}
								>
									Toggle Recurring
								</button>
								<button
									class="rounded bg-slate-800 px-2 py-1"
									on:click={() =>
										onToggle(transaction, 'essential', !transaction.is_essential)}
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
</div>

