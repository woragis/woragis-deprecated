<script lang="ts">
	import type { ReportSchedule, ReportDefinition } from '$lib/api/types';
	import { goto } from '$app/navigation';

	export interface ScheduleWithReport {
		schedule: ReportSchedule;
		report: ReportDefinition;
	}

	export let schedules: ScheduleWithReport[] = [];
	export let isLoading = false;

	export let onEditSchedule: (schedule: ScheduleWithReport) => void;
	export let onToggleSchedule: (schedule: ScheduleWithReport) => void;
	export let onDeleteSchedule: (schedule: ScheduleWithReport) => void;

	const navigateToReport = (reportId: string) => {
		goto(`/reports`);
		// Could scroll to or highlight the report, but for now just navigate
	};
</script>

{#if isLoading}
	<div class="flex items-center justify-center py-12">
		<p class="text-sm text-slate-400">Loading schedules...</p>
	</div>
{:else if schedules.length === 0}
	<div class="flex flex-col items-center justify-center gap-3 rounded-2xl border border-dashed border-slate-800/80 bg-slate-950/60 p-10 text-center">
		<h2 class="text-xl font-semibold text-slate-100">No schedules found</h2>
		<p class="max-w-md text-sm text-slate-400">
			Create report definitions and add schedules to them in the{' '}
			<a href="/reports" class="text-primary hover:underline">Reports</a> page.
		</p>
	</div>
{:else}
	<div class="space-y-4">
		{#each schedules as scheduleWithReport (scheduleWithReport.schedule.id)}
			<div
				class="rounded-xl border border-slate-800/80 bg-slate-950/60 p-4 transition hover:border-slate-700"
			>
				<div class="flex items-start justify-between gap-4">
					<div class="flex-1 space-y-2">
						<div class="flex items-center gap-3">
							<h3 class="text-base font-semibold text-slate-100">
								{scheduleWithReport.report.name}
							</h3>
							<button
								class="text-xs text-primary hover:underline"
								type="button"
								on:click={() => navigateToReport(scheduleWithReport.report.id)}
							>
								View Report →
							</button>
						</div>
						{#if scheduleWithReport.report.description}
							<p class="text-sm text-slate-400">{scheduleWithReport.report.description}</p>
						{/if}
						<div class="flex flex-wrap items-center gap-3 text-xs uppercase tracking-wide text-slate-500">
							<span>{scheduleWithReport.schedule.frequency}</span>
							<span>•</span>
							<span>{scheduleWithReport.schedule.timezone}</span>
							{#if scheduleWithReport.schedule.next_run}
								<span>•</span>
								<span>
									Next: {new Date(scheduleWithReport.schedule.next_run).toLocaleString()}
								</span>
							{/if}
						</div>
						<div class="mt-2 rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
							<div class="text-xs font-medium uppercase tracking-wide text-slate-400">Cron Expression</div>
							<div class="mt-1 font-mono text-sm text-slate-100">{scheduleWithReport.schedule.cron}</div>
						</div>
						{#if scheduleWithReport.schedule.meta && Object.keys(scheduleWithReport.schedule.meta).length > 0}
							<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
								<div class="text-xs font-medium uppercase tracking-wide text-slate-400">Metadata</div>
								<pre class="mt-1 max-h-32 overflow-auto text-xs text-slate-300">
									{JSON.stringify(scheduleWithReport.schedule.meta, null, 2)}
								</pre>
							</div>
						{/if}
					</div>
					<div class="flex flex-col gap-2">
						<button
							class={`rounded-md border px-3 py-1.5 text-xs font-medium transition ${
								scheduleWithReport.schedule.enabled
									? 'border-emerald-500/60 bg-emerald-500/10 text-emerald-200 hover:bg-emerald-500/20'
									: 'border-slate-600 bg-slate-800/40 text-slate-400 hover:bg-slate-800/60'
							}`}
							type="button"
							on:click={() => onToggleSchedule(scheduleWithReport)}
						>
							{scheduleWithReport.schedule.enabled ? '✓ Enabled' : 'Disabled'}
						</button>
						<div class="flex gap-2">
							<button
								class="rounded-md border border-slate-700/60 px-3 py-1.5 text-xs text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
								type="button"
								on:click={() => onEditSchedule(scheduleWithReport)}
							>
								Edit
							</button>
							<button
								class="rounded-md border border-red-500/60 px-3 py-1.5 text-xs text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
								type="button"
								on:click={() => onDeleteSchedule(scheduleWithReport)}
							>
								Delete
							</button>
						</div>
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}

