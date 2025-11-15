<script lang="ts">
	import type { ScheduleFormState } from '../reports.logic';

	export let open = false;
	export let mode: 'create' | 'edit' = 'create';
	export let form: ScheduleFormState;
	export let onFieldChange: <K extends keyof ScheduleFormState>(
		field: K,
		value: ScheduleFormState[K]
	) => void;
	export let onClose: () => void;
	export let onSubmit: () => void;
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
		<div class="w-full max-w-lg rounded-2xl border border-slate-800 bg-slate-950 p-6 shadow-2xl">
			<header class="mb-4 flex items-center justify-between">
				<h2 class="text-lg font-semibold text-slate-100">
					{mode === 'create' ? 'Add Schedule' : 'Edit Schedule'}
				</h2>
				<button class="rounded border border-slate-700 px-2 py-1 text-sm text-slate-300" on:click={onClose}>
					Close
				</button>
			</header>
			<div class="space-y-3 text-sm text-slate-100">
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Cron
					<input
						class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 font-mono"
						value={form.cron}
						on:input={(event) => onFieldChange('cron', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<div class="grid gap-3 md:grid-cols-2">
					<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
						Frequency
						<input
							class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							value={form.frequency}
							on:input={(event) =>
								onFieldChange('frequency', (event.target as HTMLInputElement).value)}
						/>
					</label>
					<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
						Timezone
						<input
							class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							value={form.timezone}
							on:input={(event) =>
								onFieldChange('timezone', (event.target as HTMLInputElement).value)}
						/>
					</label>
				</div>
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Next run (optional)
					<input
						type="datetime-local"
						class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						value={form.nextRun}
						on:input={(event) => onFieldChange('nextRun', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<label class="flex items-center gap-2 text-sm text-slate-300">
					<input
						type="checkbox"
						checked={form.enabled}
						on:change={(event) => onFieldChange('enabled', (event.target as HTMLInputElement).checked)}
					/>
					Enabled
				</label>
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Meta (JSON)
					<textarea
						class="mt-1 min-h-[100px] w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 font-mono"
						value={form.metaText}
						on:input={(event) => onFieldChange('metaText', (event.target as HTMLTextAreaElement).value)}
					></textarea>
				</label>
			</div>
			<div class="mt-4 flex justify-end gap-3">
				<button
					class="rounded-lg border border-slate-700 px-4 py-2 text-sm text-slate-200 hover:border-slate-500"
					type="button"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white hover:bg-primary/80"
					type="button"
					on:click={onSubmit}
				>
					{mode === 'create' ? 'Create schedule' : 'Update schedule'}
				</button>
			</div>
		</div>
	</div>
{/if}

