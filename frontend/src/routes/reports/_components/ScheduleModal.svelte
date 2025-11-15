<script lang="ts">
	import type { ScheduleFormState } from '../reports.logic';
	import {
		getCronPresets,
		WEEKDAYS,
		dailyCron,
		weeklyCron,
		monthlyCron,
		every7DaysCron,
		every14DaysCron,
		every30DaysCron,
		parseCron
	} from '../cron-utils';

	export let open = false;
	export let mode: 'create' | 'edit' = 'create';
	export let form: ScheduleFormState;
	export let onFieldChange: <K extends keyof ScheduleFormState>(
		field: K,
		value: ScheduleFormState[K]
	) => void;
	export let onClose: () => void;
	export let onSubmit: () => void;

	let scheduleType: 'preset' | 'custom' | 'weekly' | 'monthly' = 'preset';
	let selectedWeekday: string = '1'; // Monday
	let selectedDayOfMonth: string = '1';
	let selectedHour: number = 8;
	let selectedMinute: number = 0;

	$: {
		// Parse current cron to determine schedule type
		if (form.cron) {
			const parsed = parseCron(form.cron);
			if (parsed) {
				// Try to determine type
				if (parsed.dayOfWeek !== '*' && parsed.dayOfMonth === '*') {
					scheduleType = 'weekly';
					selectedWeekday = parsed.dayOfWeek;
				} else if (parsed.dayOfMonth !== '*' && parsed.dayOfWeek === '*') {
					scheduleType = 'monthly';
					selectedDayOfMonth = parsed.dayOfMonth;
				} else if (parsed.dayOfMonth === '*' && parsed.dayOfWeek === '*') {
					scheduleType = 'preset';
				} else {
					scheduleType = 'custom';
				}
				selectedHour = parseInt(parsed.hour) || 8;
				selectedMinute = parseInt(parsed.minute) || 0;
			}
		}
	}

	function applyPreset(preset: ReturnType<typeof getCronPresets>[0]) {
		onFieldChange('cron', preset.cron);
		onFieldChange('frequency', preset.frequency);
		scheduleType = 'preset';
	}

	function applyWeeklySchedule() {
		const cron = weeklyCron(parseInt(selectedWeekday), selectedHour, selectedMinute);
		onFieldChange('cron', cron);
		onFieldChange('frequency', 'weekly');
		scheduleType = 'weekly';
	}

	function applyMonthlySchedule() {
		const day = parseInt(selectedDayOfMonth) || 1;
		const cron = monthlyCron(day, selectedHour, selectedMinute);
		onFieldChange('cron', cron);
		onFieldChange('frequency', 'monthly');
		scheduleType = 'monthly';
	}

	function applyCustomCron() {
		scheduleType = 'custom';
	}

	$: presets = getCronPresets(selectedHour, selectedMinute);
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

