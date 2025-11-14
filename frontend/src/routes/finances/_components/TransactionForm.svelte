<script lang="ts">
	import type { TransactionFormState } from '../finances.logic

	export let formState: TransactionFormState;
	export let isMutating = false;
	export let onSubmit: () => void;
	export let onFieldChange: <K extends keyof TransactionFormState>(
		field: K,
		value: TransactionFormState[K]
	) => void;
</script>

<div class="space-y-4 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
	<h3 class="text-sm font-semibold text-slate-100">Record Transaction</h3>
	<form class="space-y-3" on:submit|preventDefault={onSubmit}>
		<div class="grid gap-3 md:grid-cols-2">
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Type
				<select
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.type}
					on:change={(event) =>
						onFieldChange('type', (event.target as HTMLSelectElement).value as 'income' | 'expense')}
				>
					<option value="income">Income</option>
					<option value="expense">Expense</option>
				</select>
			</label>
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Category
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.category}
					required
					on:input={(event) => onFieldChange('category', (event.target as HTMLInputElement).value)}
				/>
			</label>
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Amount
				<input
					type="number"
					step="0.01"
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.amount}
					required
					on:input={(event) =>
						onFieldChange('amount', Number((event.target as HTMLInputElement).value))}
				/>
			</label>
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Currency
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.currency}
					maxlength="3"
					on:input={(event) => onFieldChange('currency', (event.target as HTMLInputElement).value)}
				/>
			</label>
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Base Currency
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.baseCurrency}
					maxlength="3"
					on:input={(event) => onFieldChange('baseCurrency', (event.target as HTMLInputElement).value)}
				/>
			</label>
			<label class="flex flex-col gap-1 text-xs text-slate-400">
				Occurred At
				<input
					type="datetime-local"
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					value={formState.occurredAt}
					on:input={(event) => onFieldChange('occurredAt', (event.target as HTMLInputElement).value)}
				/>
			</label>
		</div>
		<label class="flex flex-col gap-1 text-xs text-slate-400">
			Description
			<textarea
				class="min-h-[80px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
				value={formState.description}
				on:input={(event) => onFieldChange('description', (event.target as HTMLTextAreaElement).value)}
			></textarea>
		</label>
		<label class="flex flex-col gap-1 text-xs text-slate-400">
			Tags (comma separated)
			<input
				class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
				value={formState.tags}
				on:input={(event) => onFieldChange('tags', (event.target as HTMLInputElement).value)}
			/>
		</label>
		<div class="flex flex-wrap items-center gap-4 text-xs text-slate-300">
			<label class="flex items-center gap-2">
				<input
					type="checkbox"
					checked={formState.isRecurring}
					on:change={(event) => onFieldChange('isRecurring', (event.target as HTMLInputElement).checked)}
				/>
				<span>Recurring</span>
			</label>
			<label class="flex items-center gap-2">
				<input
					type="checkbox"
					checked={formState.isEssential}
					on:change={(event) => onFieldChange('isEssential', (event.target as HTMLInputElement).checked)}
				/>
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

