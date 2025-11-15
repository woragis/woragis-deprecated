<script lang="ts">
	import type {
		ReportDefinitionDetail,
		ReportDelivery,
		ReportRun,
		ReportSchedule,
		UUID
	} from '$lib/api/types';

	export let detail: ReportDefinitionDetail | null = null;
	export let runs: ReportRun[] = [];
	export let detailLoading = false;

	export let onToggleFavorite: (definitionId: UUID, isFavorite: boolean) => void;
	export let onEditDefinition: () => void;
	export let onArchive: (definitionId: UUID) => void;
	export let onRestore: (definitionId: UUID) => void;
	export let onAddSchedule: () => void;
	export let onEditSchedule: (schedule: ReportSchedule) => void;
	export let onToggleSchedule: (schedule: ReportSchedule) => void;
	export let onDeleteSchedule: (schedule: ReportSchedule) => void;
	export let onAddDelivery: () => void;
	export let onEditDelivery: (delivery: ReportDelivery) => void;
	export let onToggleDelivery: (delivery: ReportDelivery) => void;
	export let onDeleteDelivery: (delivery: ReportDelivery) => void;
	export let onRefreshRuns: () => void;
</script>

{#if !detail}
	<div class="flex flex-1 flex-col items-center justify-center gap-3 rounded-2xl border border-dashed border-slate-800/80 bg-slate-950/60 p-10 text-center">
		<h2 class="text-xl font-semibold text-slate-100">Select a report definition</h2>
		<p class="max-w-md text-sm text-slate-400">
			Pick a definition on the left to view its schedules, delivery settings, and run history.
		</p>
	</div>
{:else}
	<div class="flex flex-col gap-6 rounded-2xl border border-slate-800/80 bg-slate-950/60 p-6">
		<div class="flex flex-wrap items-start justify-between gap-4">
			<div class="space-y-2">
				<div class="flex items-center gap-3">
					<h2 class="text-xl font-semibold text-slate-100">{detail.definition.name}</h2>
					<button
						class={`text-sm ${
							detail.definition.is_favorite ? 'text-amber-300' : 'text-slate-500 hover:text-slate-200'
						}`}
						on:click={() => onToggleFavorite(detail.definition.id, detail.definition.is_favorite)}
					>
						{detail.definition.is_favorite ? '★' : '☆'}
					</button>
				</div>
				<p class="max-w-2xl text-sm text-slate-300">{detail.definition.description}</p>
				<div class="flex flex-wrap items-center gap-3 text-xs uppercase tracking-wide text-slate-500">
					<span>Updated {new Date(detail.definition.updated_at).toLocaleString()}</span>
					{#if detail.definition.archived_at}
						<span class="rounded-full border border-red-500/40 px-2 py-[1px] text-[10px] text-red-200">
							Archived at {new Date(detail.definition.archived_at).toLocaleDateString()}
						</span>
					{/if}
				</div>
			</div>
			<div class="flex items-center gap-2">
				<button
					class="rounded-lg border border-slate-700/60 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={onEditDefinition}
				>
					Edit definition
				</button>
				<button
					class="rounded-lg border border-slate-700/60 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={() => onArchive(detail.definition.id)}
				>
					Archive
				</button>
				<button
					class="rounded-lg border border-slate-700/60 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={() => onRestore(detail.definition.id)}
				>
					Restore
				</button>
			</div>
		</div>

		<div class="grid gap-4 md:grid-cols-2">
			<div class="rounded-xl border border-slate-800/70 bg-slate-900/40 p-4">
				<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">Sections</h3>
				<pre class="mt-3 max-h-56 overflow-auto rounded-lg bg-slate-950/80 p-3 text-xs text-slate-300">{JSON.stringify(detail.definition.sections ?? {}, null, 2)}</pre>
			</div>
			<div class="rounded-xl border border-slate-800/70 bg-slate-900/40 p-4">
				<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">Filters</h3>
				<pre class="mt-3 max-h-56 overflow-auto rounded-lg bg-slate-950/80 p-3 text-xs text-slate-300">{JSON.stringify(detail.definition.filters ?? {}, null, 2)}</pre>
			</div>
		</div>

		<div class="grid gap-6 lg:grid-cols-2">
			<div class="flex flex-col gap-3 rounded-xl border border-slate-800/80 bg-slate-900/50 p-4">
				<div class="flex items-center justify-between">
					<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">Schedules</h3>
					<button
						class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
						type="button"
						on:click={onAddSchedule}
					>
						Add schedule
					</button>
				</div>
				{#if detailLoading}
					<div class="flex items-center justify-center py-8 text-sm text-slate-400">Loading schedules…</div>
				{:else if detail.schedules.length === 0}
					<div class="rounded-lg border border-slate-800/60 bg-slate-950/70 px-3 py-6 text-center text-sm text-slate-400">
						No schedules configured. Create one to automate report generation.
					</div>
				{:else}
					<div class="space-y-3">
						{#each detail.schedules as schedule}
							<div class="rounded-lg border border-slate-800/60 bg-slate-950/70 p-3 text-sm text-slate-200">
								<div class="flex items-center justify-between gap-2">
									<div>
										<div class="flex items-center gap-2 text-xs uppercase tracking-wide text-slate-400">
											<span>{schedule.frequency}</span>
											<span>•</span>
											<span>{schedule.timezone}</span>
										</div>
										<div class="mt-1 font-mono text-slate-100">{schedule.cron}</div>
										<div class="mt-2 text-xs text-slate-400">
											Next run:{' '}
											{schedule.next_run ? new Date(schedule.next_run).toLocaleString() : 'Not scheduled'}
										</div>
									</div>
									<div class="flex flex-col gap-2">
										<button
											class={`rounded-md border px-2 py-1 text-xs ${
												schedule.enabled
													? 'border-emerald-500/60 text-emerald-200'
													: 'border-slate-600 text-slate-400'
											}`}
											type="button"
											on:click={() => onToggleSchedule(schedule)}
										>
											{schedule.enabled ? 'Enabled' : 'Disabled'}
										</button>
										<div class="flex gap-2">
											<button
												class="rounded-md border border-slate-700/60 px-2 py-1 text-xs text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
												type="button"
												on:click={() => onEditSchedule(schedule)}
											>
												Edit
											</button>
											<button
												class="rounded-md border border-red-500/60 px-2 py-1 text-xs text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
												type="button"
												on:click={() => onDeleteSchedule(schedule)}
											>
												Delete
											</button>
										</div>
									</div>
								</div>
								{#if schedule.meta && Object.keys(schedule.meta).length > 0}
									<pre class="mt-3 max-h-40 overflow-auto rounded-lg bg-slate-900/80 p-3 text-xs text-slate-300">{JSON.stringify(schedule.meta, null, 2)}</pre>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<div class="flex flex-col gap-3 rounded-xl border border-slate-800/80 bg-slate-900/50 p-4">
				<div class="flex items-center justify-between">
					<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">Deliveries</h3>
					<button
						class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
						type="button"
						on:click={onAddDelivery}
					>
						Add delivery
					</button>
				</div>
				{#if detailLoading}
					<div class="flex items-center justify-center py-8 text-sm text-slate-400">Loading deliveries…</div>
				{:else if detail.deliveries.length === 0}
					<div class="rounded-lg border border-slate-800/60 bg-slate-950/70 px-3 py-6 text-center text-sm text-slate-400">
						No delivery channels configured. Add one to distribute reports.
					</div>
				{:else}
					<div class="space-y-3">
						{#each detail.deliveries as delivery}
							<div class="rounded-lg border border-slate-800/60 bg-slate-950/70 p-3 text-sm text-slate-200">
								<div class="flex items-center justify-between gap-2">
									<div>
										<div class="text-xs uppercase tracking-wide text-slate-400">{delivery.channel}</div>
										<div class="mt-1 font-medium text-slate-100">{delivery.target}</div>
									</div>
									<div class="flex flex-col gap-2">
										<button
											class={`rounded-md border px-2 py-1 text-xs ${
												delivery.enabled
													? 'border-emerald-500/60 text-emerald-200'
													: 'border-slate-600 text-slate-400'
											}`}
											type="button"
											on:click={() => onToggleDelivery(delivery)}
										>
											{delivery.enabled ? 'Enabled' : 'Disabled'}
										</button>
										<div class="flex gap-2">
											<button
												class="rounded-md border border-slate-700/60 px-2 py-1 text-xs text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
												type="button"
												on:click={() => onEditDelivery(delivery)}
											>
												Edit
											</button>
											<button
												class="rounded-md border border-red-500/60 px-2 py-1 text-xs text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
												type="button"
												on:click={() => onDeleteDelivery(delivery)}
											>
												Delete
											</button>
										</div>
									</div>
								</div>
								{#if delivery.template && Object.keys(delivery.template).length > 0}
									<pre class="mt-3 max-h-40 overflow-auto rounded-lg bg-slate-900/80 p-3 text-xs text-slate-300">{JSON.stringify(delivery.template, null, 2)}</pre>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<div class="rounded-xl border border-slate-800/80 bg-slate-900/50 p-4">
			<div class="flex items-center justify-between">
				<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">Run history</h3>
				<button
					class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={onRefreshRuns}
				>
					Refresh
				</button>
			</div>
			{#if runs.length === 0}
				<div class="mt-3 rounded-lg border border-slate-800/60 bg-slate-950/70 px-3 py-6 text-center text-sm text-slate-400">
					No runs yet. Queue a run from the bulk panel to generate a report.
				</div>
			{:else}
				<div class="mt-3 space-y-3">
					{#each runs as run}
						<div class="rounded-lg border border-slate-800/60 bg-slate-950/70 p-3 text-sm text-slate-200">
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center gap-2 text-xs uppercase tracking-wide">
									<span
										class={`rounded-full px-2 py-[1px] ${
											run.status === 'completed'
												? 'bg-emerald-500/20 text-emerald-200'
												: run.status === 'failed'
												? 'bg-red-500/20 text-red-200'
												: run.status === 'running'
												? 'bg-sky-500/20 text-sky-200'
												: 'bg-slate-700/40 text-slate-300'
										}`}
									>
										{run.status}
									</span>
									<span class="text-slate-400">{new Date(run.created_at).toLocaleString()}</span>
								</div>
								{#if run.output_location}
									<a
										class="text-xs text-primary hover:underline"
										href={run.output_location}
										target="_blank"
										rel="noreferrer"
									>
										View output
									</a>
								{/if}
							</div>
							{#if run.error_message}
								<p class="mt-2 text-xs text-red-300">{run.error_message}</p>
							{/if}
							{#if run.metadata && Object.keys(run.metadata).length > 0}
								<pre class="mt-3 max-h-40 overflow-auto rounded-lg bg-slate-900/80 p-3 text-xs text-slate-300">{JSON.stringify(run.metadata, null, 2)}</pre>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}

