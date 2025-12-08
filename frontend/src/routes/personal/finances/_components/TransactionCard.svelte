<script lang="ts">
	import type { Transaction, UUID } from '$lib/api/types';
	import type { ToggleAction } from '../finances.logic';

	export let transaction: Transaction;
	export let isSelected = false;
	export let numberFormatter: Intl.NumberFormat;
	export let onToggleSelection: (id: UUID, checked: boolean) => void;
	export let onToggle: (transaction: Transaction, action: ToggleAction, value: boolean) => void;
	export let onClick: ((transaction: Transaction) => void) | undefined = undefined;

	const isIncome = transaction.type === 'income';
	const typeIcon = isIncome ? '↑' : '↓';
	const typeColorClasses = isIncome
		? {
				border: 'hover:border-emerald-500/50',
				shadow: 'hover:shadow-emerald-500/20',
				text: 'text-emerald-400 group-hover:text-emerald-300',
				badge: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30',
				iconBg: 'bg-emerald-500/20 text-emerald-300',
				gradient: 'from-emerald-500/0 to-emerald-600/0 group-hover:from-emerald-500/10 group-hover:to-emerald-600/10',
				ring: 'ring-emerald-500',
				checkbox: 'text-emerald-600 focus:ring-emerald-500/20',
				button: 'hover:border-emerald-500/50',
				link: 'text-emerald-400 group-hover:text-emerald-300',
				blur: 'bg-emerald-500/5 group-hover:bg-emerald-500/10'
			}
		: {
				border: 'hover:border-rose-500/50',
				shadow: 'hover:shadow-rose-500/20',
				text: 'text-rose-400 group-hover:text-rose-300',
				badge: 'bg-rose-500/20 text-rose-300 border-rose-500/30',
				iconBg: 'bg-rose-500/20 text-rose-300',
				gradient: 'from-rose-500/0 to-rose-600/0 group-hover:from-rose-500/10 group-hover:to-rose-600/10',
				ring: 'ring-rose-500',
				checkbox: 'text-rose-600 focus:ring-rose-500/20',
				button: 'hover:border-rose-500/50',
				link: 'text-rose-400 group-hover:text-rose-300',
				blur: 'bg-rose-500/5 group-hover:bg-rose-500/10'
			};
</script>

<div
	class="group relative overflow-hidden rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-5 backdrop-blur-sm transition-all hover:scale-[1.01] {typeColorClasses.border} hover:shadow-xl {typeColorClasses.shadow}"
	class:ring-2={isSelected}
	class:{typeColorClasses.ring}={isSelected}
>
	<!-- Gradient overlay on hover -->
	<div class="absolute inset-0 bg-gradient-to-br {typeColorClasses.gradient} transition-all"></div>

	<div class="relative z-10 space-y-4">
		<!-- Header with checkbox and type -->
		<div class="flex items-start justify-between">
			<div class="flex items-start gap-3 flex-1">
				<input
					type="checkbox"
					class="mt-1 h-4 w-4 rounded border-slate-700 bg-slate-950 {typeColorClasses.checkbox} focus:ring-2"
					checked={isSelected}
					onchange={(e) => onToggleSelection(transaction.id, (e.target as HTMLInputElement).checked)}
					onclick={(e) => e.stopPropagation()}
				/>
				<div class="flex-1">
					<div class="flex items-center gap-2 mb-1">
						<span class="inline-flex h-6 w-6 items-center justify-center rounded-full {typeColorClasses.iconBg} text-xs font-bold">
							{typeIcon}
						</span>
						<span class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium capitalize {typeColorClasses.badge}">
							{transaction.type}
						</span>
					</div>
					<h3 class="text-base font-bold text-white {typeColorClasses.text} transition-colors">
						{transaction.category}
					</h3>
					{#if transaction.description}
						<p class="mt-1 line-clamp-2 text-xs text-slate-400">{transaction.description}</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Amount and Currency -->
		<div class="flex items-baseline justify-between border-t border-slate-800/50 pt-3">
			<div>
				<p class="text-xs text-slate-400">Amount</p>
				<p class="mt-1 text-2xl font-bold {typeColorClasses.text}">
					{numberFormatter.format(transaction.amount)} {transaction.currency}
				</p>
				{#if transaction.currency !== transaction.base_currency}
					<p class="mt-1 text-xs text-slate-500">
						≈ {numberFormatter.format(transaction.normalized_amount)} {transaction.base_currency}
					</p>
				{/if}
			</div>
		</div>

		<!-- Metadata -->
		<div class="grid grid-cols-2 gap-3 text-xs">
			<div>
				<p class="text-slate-400">Date</p>
				<p class="mt-1 font-medium text-slate-200">
					{new Date(transaction.occurred_at).toLocaleDateString()}
				</p>
				<p class="text-slate-500">{new Date(transaction.occurred_at).toLocaleTimeString()}</p>
			</div>
			<div>
				<p class="text-slate-400">Flags</p>
				<div class="mt-1 flex flex-wrap gap-1.5">
					{#if transaction.is_recurring}
						<span
							class="inline-flex items-center rounded-full bg-indigo-500/20 px-2 py-0.5 text-[10px] font-medium text-indigo-300"
						>
							Recurring
						</span>
					{/if}
					{#if transaction.is_essential}
						<span
							class="inline-flex items-center rounded-full bg-amber-500/20 px-2 py-0.5 text-[10px] font-medium text-amber-300"
						>
							Essential
						</span>
					{/if}
					{#if transaction.is_archived}
						<span
							class="inline-flex items-center rounded-full bg-slate-700/40 px-2 py-0.5 text-[10px] font-medium text-slate-300"
						>
							Archived
						</span>
					{/if}
					{#if transaction.tags && transaction.tags.length > 0}
						<span
							class="inline-flex items-center rounded-full bg-purple-500/20 px-2 py-0.5 text-[10px] font-medium text-purple-300"
						>
							{transaction.tags.length} tag{transaction.tags.length !== 1 ? 's' : ''}
						</span>
					{/if}
				</div>
			</div>
		</div>

		<!-- Actions -->
		<div class="flex items-center justify-between gap-2 pt-2 border-t border-slate-800/50">
			<div class="flex gap-2">
				<button
					class="rounded-lg border border-slate-700 bg-slate-800/50 px-3 py-1.5 text-xs font-medium text-slate-200 transition-all {typeColorClasses.button} hover:bg-slate-800/80"
					onclick={(e) => {
						e.stopPropagation();
						onToggle(transaction, 'archive', !transaction.is_archived);
					}}
				>
					{transaction.is_archived ? 'Unarchive' : 'Archive'}
				</button>
				<button
					class="rounded-lg border border-slate-700 bg-slate-800/50 px-3 py-1.5 text-xs font-medium text-slate-200 transition-all hover:border-indigo-500/50 hover:bg-slate-800/80"
					onclick={(e) => {
						e.stopPropagation();
						onToggle(transaction, 'recurring', !transaction.is_recurring);
					}}
				>
					{transaction.is_recurring ? 'Not Recurring' : 'Recurring'}
				</button>
				<button
					class="rounded-lg border border-slate-700 bg-slate-800/50 px-3 py-1.5 text-xs font-medium text-slate-200 transition-all hover:border-amber-500/50 hover:bg-slate-800/80"
					onclick={(e) => {
						e.stopPropagation();
						onToggle(transaction, 'essential', !transaction.is_essential);
					}}
				>
					{transaction.is_essential ? 'Not Essential' : 'Essential'}
				</button>
			</div>
			{#if onClick}
				<button
					class="text-xs font-medium {typeColorClasses.link} transition-colors"
					onclick={(e) => {
						e.stopPropagation();
						onClick?.(transaction);
					}}
				>
					View Details →
				</button>
			{/if}
		</div>
	</div>

	<!-- Decorative gradient -->
	<div class="absolute -right-8 -top-8 h-24 w-24 rounded-full {typeColorClasses.blur} blur-2xl transition-all"></div>
</div>

