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
	let isInitialized = false;

	function parseCurrentCron() {
		if (form.cron && !isInitialized) {
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
				isInitialized = true;
			}
		}
	}

	$: if (open && !isInitialized) {
		parseCurrentCron();
	}

	$: if (!open) {
		isInitialized = false;
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
			<div class="space-y-4 text-sm text-slate-100">
				<!-- Schedule Type Selection -->
				<div class="space-y-2">
					<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
						Schedule Type
					</label>
					<div class="grid grid-cols-2 gap-2">
						<button
							class={`rounded-lg border px-3 py-2 text-xs transition ${
								scheduleType === 'preset'
									? 'border-primary bg-primary/20 text-primary'
									: 'border-slate-700 text-slate-300 hover:border-slate-500'
							}`}
							type="button"
							on:click={() => (scheduleType = 'preset')}
						>
							Quick Presets
						</button>
						<button
							class={`rounded-lg border px-3 py-2 text-xs transition ${
								scheduleType === 'weekly'
									? 'border-primary bg-primary/20 text-primary'
									: 'border-slate-700 text-slate-300 hover:border-slate-500'
							}`}
							type="button"
							on:click={() => (scheduleType = 'weekly')}
						>
							Weekly
						</button>
						<button
							class={`rounded-lg border px-3 py-2 text-xs transition ${
								scheduleType === 'monthly'
									? 'border-primary bg-primary/20 text-primary'
									: 'border-slate-700 text-slate-300 hover:border-slate-500'
							}`}
							type="button"
							on:click={() => (scheduleType = 'monthly')}
						>
							Monthly
						</button>
						<button
							class={`rounded-lg border px-3 py-2 text-xs transition ${
								scheduleType === 'custom'
									? 'border-primary bg-primary/20 text-primary'
									: 'border-slate-700 text-slate-300 hover:border-slate-500'
							}`}
							type="button"
							on:click={() => (scheduleType = 'custom')}
						>
							Custom Cron
						</button>
					</div>
				</div>

				<!-- Quick Presets -->
				{#if scheduleType === 'preset'}
					<div class="space-y-2">
						<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
							Quick Presets
						</label>
						<div class="grid grid-cols-2 gap-2">
							{#each presets as preset}
								<button
									class="rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-left text-xs transition hover:border-primary hover:bg-primary/10"
									type="button"
									on:click={() => applyPreset(preset)}
								>
									<div class="font-medium text-slate-100">{preset.name}</div>
									<div class="mt-0.5 text-[10px] text-slate-400">{preset.description}</div>
								</button>
							{/each}
						</div>
						<div class="mt-2 grid grid-cols-2 gap-2">
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Hour (24h)
								<input
									type="number"
									min="0"
									max="23"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedHour}
									on:input={(e) => {
										selectedHour = parseInt((e.target as HTMLInputElement).value) || 8;
										const preset = presets[0];
										if (preset) applyPreset(preset);
									}}
								/>
							</label>
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Minute
								<input
									type="number"
									min="0"
									max="59"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedMinute}
									on:input={(e) => {
										selectedMinute = parseInt((e.target as HTMLInputElement).value) || 0;
										const preset = presets[0];
										if (preset) applyPreset(preset);
									}}
								/>
							</label>
						</div>
					</div>
				{/if}

				<!-- Weekly Schedule -->
				{#if scheduleType === 'weekly'}
					<div class="space-y-2">
						<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
							Select Weekday
						</label>
						<select
							class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							value={selectedWeekday}
							on:change={(e) => {
								selectedWeekday = (e.target as HTMLSelectElement).value;
								applyWeeklySchedule();
							}}
						>
							{#each WEEKDAYS as day}
								<option value={day.value}>{day.label}</option>
							{/each}
						</select>
						<div class="grid grid-cols-2 gap-2">
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Hour (24h)
								<input
									type="number"
									min="0"
									max="23"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedHour}
									on:input={(e) => {
										selectedHour = parseInt((e.target as HTMLInputElement).value) || 8;
										applyWeeklySchedule();
									}}
								/>
							</label>
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Minute
								<input
									type="number"
									min="0"
									max="59"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedMinute}
									on:input={(e) => {
										selectedMinute = parseInt((e.target as HTMLInputElement).value) || 0;
										applyWeeklySchedule();
									}}
								/>
							</label>
						</div>
					</div>
				{/if}

				<!-- Monthly Schedule -->
				{#if scheduleType === 'monthly'}
					<div class="space-y-2">
						<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
							Day of Month (1-31)
						</label>
						<input
							type="number"
							min="1"
							max="31"
							class="w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							value={selectedDayOfMonth}
							on:input={(e) => {
								selectedDayOfMonth = (e.target as HTMLInputElement).value;
								applyMonthlySchedule();
							}}
						/>
						<div class="grid grid-cols-2 gap-2">
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Hour (24h)
								<input
									type="number"
									min="0"
									max="23"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedHour}
									on:input={(e) => {
										selectedHour = parseInt((e.target as HTMLInputElement).value) || 8;
										applyMonthlySchedule();
									}}
								/>
							</label>
							<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
								Minute
								<input
									type="number"
									min="0"
									max="59"
									class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
									value={selectedMinute}
									on:input={(e) => {
										selectedMinute = parseInt((e.target as HTMLInputElement).value) || 0;
										applyMonthlySchedule();
									}}
								/>
							</label>
						</div>
					</div>
				{/if}

				<!-- Custom Cron -->
				{#if scheduleType === 'custom'}
					<div class="space-y-2">
						<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
							Cron Expression (minute hour day-of-month month day-of-week)
							<input
								class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 font-mono"
								value={form.cron}
								on:input={(event) => {
									onFieldChange('cron', (event.target as HTMLInputElement).value);
									applyCustomCron();
								}}
								placeholder="0 8 * * *"
							/>
						</label>
						<p class="text-xs text-slate-400">
							Format: minute hour day-of-month month day-of-week<br />
							Example: "0 8 * * *" = Daily at 8:00 AM
						</p>
					</div>
				{/if}

				<!-- Generated Cron Display -->
				{#if scheduleType !== 'custom'}
					<div class="rounded-lg border border-slate-700/60 bg-slate-900/40 p-3">
						<div class="text-xs font-medium uppercase tracking-wide text-slate-400">
							Generated Cron
						</div>
						<div class="mt-1 font-mono text-sm text-slate-100">{form.cron}</div>
					</div>
				{/if}

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

