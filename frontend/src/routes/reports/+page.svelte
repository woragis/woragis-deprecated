<script lang="ts">
	import { onMount } from 'svelte';

	import { createReportsLogic } from './reports.logic';
	import ReportsHeader from './_components/ReportsHeader.svelte';
	import ReportsFilters from './_components/ReportsFilters.svelte';
	import DefinitionsList from './_components/DefinitionsList.svelte';
	import BulkActionsPanel from './_components/BulkActionsPanel.svelte';
	import DefinitionDetailPanel from './_components/DefinitionDetailPanel.svelte';
	import DefinitionModal from './_components/DefinitionModal.svelte';
	import ScheduleModal from './_components/ScheduleModal.svelte';
	import DeliveryModal from './_components/DeliveryModal.svelte';

	const {
		isLoading,
		detailLoading,
		definitions,
		detail,
		runs,
		filters,
		selectedDefinitionId,
		selectedIds,
		queueMetadataText,
		showDefinitionForm,
		definitionFormMode,
		definitionForm,
		showScheduleForm,
		scheduleFormMode,
		scheduleForm,
		showDeliveryForm,
		deliveryFormMode,
		deliveryForm,
		errorMessage,
		updateDefinitionFormField,
		updateScheduleFormField,
		updateDeliveryFormField,
		updateQueueMetadata,
		loadDefinitions,
		selectDefinition,
		toggleSelection,
		toggleSelectAll,
		updateFilters,
		openCreateDefinition,
		openEditDefinition,
		closeDefinitionForm,
		submitDefinitionForm,
		handleBulkAction,
		queueSelectedRuns,
		toggleFavorite,
		openCreateSchedule,
		openEditSchedule,
		closeScheduleForm,
		submitScheduleForm,
		toggleScheduleEnabled,
		deleteSchedule,
		openCreateDelivery,
		openEditDelivery,
		closeDeliveryForm,
		submitDeliveryForm,
		toggleDeliveryEnabled,
		deleteDelivery
	} = createReportsLogic();

	let hasSelectedInitial = false;

	onMount(async () => {
		await loadDefinitions();
	});

	$: if ($definitions.length === 0) {
		hasSelectedInitial = false;
	}

	$: if (!hasSelectedInitial && $definitions.length > 0) {
		hasSelectedInitial = true;
		selectDefinition($definitions[0].id);
	}

	$: allSelected = $definitions.length > 0 && $selectedIds.size === $definitions.length;

	const archiveSelected = () => handleBulkAction('archive');
	const restoreSelected = () => handleBulkAction('restore');
	const deleteSelected = () => handleBulkAction('delete');
	const refreshDetail = () => {
		if ($selectedDefinitionId) {
			selectDefinition($selectedDefinitionId);
		}
	};
</script>

	const startCreateSchedule = () => {
		resetScheduleForm();
		showScheduleForm = true;
	};

	const startEditSchedule = (schedule: ReportSchedule) => {
		scheduleFormMode = 'edit';
		editingScheduleId = schedule.id;
		scheduleForm = {
			cron: schedule.cron,
			frequency: schedule.frequency,
			timezone: schedule.timezone,
			nextRun: schedule.next_run ?? '',
			enabled: schedule.enabled,
			metaText: schedule.meta ? JSON.stringify(schedule.meta, null, 2) : ''
		};
		showScheduleForm = true;
	};

	const submitScheduleForm = async () => {
		if (!selectedDefinitionId) return;
		const meta = scheduleForm.metaText.trim()
			? parseJsonField(scheduleForm.metaText, 'Meta')
			: undefined;

		try {
			if (scheduleFormMode === 'create') {
				await createReportSchedule(selectedDefinitionId, {
					cron: scheduleForm.cron,
					frequency: scheduleForm.frequency,
					timezone: scheduleForm.timezone,
					nextRun: scheduleForm.nextRun || undefined,
					enabled: scheduleForm.enabled,
					meta
				});
				toastSuccess('Schedule created.');
			} else if (scheduleFormMode === 'edit' && editingScheduleId) {
				await updateReportSchedule(editingScheduleId, {
					cron: scheduleForm.cron,
					frequency: scheduleForm.frequency,
					timezone: scheduleForm.timezone,
					nextRun: scheduleForm.nextRun || undefined,
					enabled: scheduleForm.enabled,
					meta
				});
				toastSuccess('Schedule updated.');
			}
			showScheduleForm = false;
			resetScheduleForm();
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Unable to save schedule.'));
		}
	};

	const toggleScheduleEnabled = async (schedule: ReportSchedule) => {
		try {
			await toggleReportSchedule(schedule.id, !schedule.enabled);
			toastSuccess('Schedule updated.');
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Failed to toggle schedule.'));
		}
	};

	const removeSchedule = async (schedule: ReportSchedule) => {
		try {
			await deleteReportSchedule(schedule.id);
			toastInfo('Schedule deleted.');
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Failed to delete schedule.'));
		}
	};

	const resetDeliveryForm = () => {
		deliveryForm = {
			channel: 'email',
			target: '',
			templateText: '',
			enabled: true
		};
		deliveryFormMode = 'create';
		editingDeliveryId = null;
	};

	const startCreateDelivery = () => {
		resetDeliveryForm();
		showDeliveryForm = true;
	};

	const startEditDelivery = (delivery: ReportDelivery) => {
		deliveryFormMode = 'edit';
		editingDeliveryId = delivery.id;
		deliveryForm = {
			channel: delivery.channel,
			target: delivery.target,
			templateText: delivery.template ? JSON.stringify(delivery.template, null, 2) : '',
			enabled: delivery.enabled
		};
		showDeliveryForm = true;
	};

	const submitDeliveryForm = async () => {
		if (!selectedDefinitionId) return;
		const template = deliveryForm.templateText.trim()
			? parseJsonField(deliveryForm.templateText, 'Template')
			: undefined;

		try {
			if (deliveryFormMode === 'create') {
				await createReportDelivery(selectedDefinitionId, {
					channel: deliveryForm.channel,
					target: deliveryForm.target,
					template,
					enabled: deliveryForm.enabled
				});
				toastSuccess('Delivery created.');
			} else if (deliveryFormMode === 'edit' && editingDeliveryId) {
				await updateReportDelivery(editingDeliveryId, {
					channel: deliveryForm.channel,
					target: deliveryForm.target,
					template,
					enabled: deliveryForm.enabled
				});
				toastSuccess('Delivery updated.');
			}
			showDeliveryForm = false;
			resetDeliveryForm();
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Unable to save delivery.'));
		}
	};

	const toggleDeliveryEnabled = async (delivery: ReportDelivery) => {
		try {
			await toggleReportDelivery(delivery.id, !delivery.enabled);
			toastSuccess('Delivery updated.');
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Failed to toggle delivery.'));
		}
	};

	const removeDelivery = async (delivery: ReportDelivery) => {
		try {
			await deleteReportDelivery(delivery.id);
			toastInfo('Delivery deleted.');
			if (selectedDefinitionId) {
				await loadDetail(selectedDefinitionId);
			}
		} catch (error) {
			console.error(error);
			toastError(getApiErrorMessage(error, 'Failed to delete delivery.'));
		}
	};

	onMount(async () => {
		await loadDefinitions();
		if (definitions.length > 0) {
			await selectDefinition(definitions[0].id);
		}
	});
</script>

<svelte:head>
	<title>Reports · Woragis</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<ReportsHeader
		isLoading={$isLoading}
		onRefresh={loadDefinitions}
		onCreate={openCreateDefinition}
	/>

	<section class="grid gap-6 lg:grid-cols-[350px_1fr]">
		<div class="flex flex-col gap-4 rounded-2xl border border-slate-800/80 bg-slate-950/60 p-4">
			<ReportsFilters filters={$filters} onChange={updateFilters} onApply={loadDefinitions} />
			<DefinitionsList
				definitions={$definitions}
				isLoading={$isLoading}
				selectedDefinitionId={$selectedDefinitionId}
				selectedIds={$selectedIds}
				allSelected={allSelected}
				onSelectDefinition={selectDefinition}
				onToggleSelection={toggleSelection}
				onToggleSelectAll={toggleSelectAll}
				onToggleFavorite={toggleFavorite}
			/>
			<BulkActionsPanel
				selectedCount={$selectedIds.size}
				queueMetadataText={$queueMetadataText}
				onArchive={archiveSelected}
				onRestore={restoreSelected}
				onDelete={deleteSelected}
				onQueueRuns={queueSelectedRuns}
				onMetadataChange={updateQueueMetadata}
			/>
		</div>

		<div class="flex flex-col gap-6">
			{#if $errorMessage}
				<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
					{$errorMessage}
				</p>
			{/if}

			<DefinitionDetailPanel
				detail={$detail}
				runs={$runs}
				detailLoading={$detailLoading}
				onToggleFavorite={toggleFavorite}
				onEditDefinition={openEditDefinition}
				onArchive={(id) => {
					selectedIds.set(new Set([id]));
					archiveSelected();
				}}
				onRestore={(id) => {
					selectedIds.set(new Set([id]));
					restoreSelected();
				}}
				onAddSchedule={openCreateSchedule}
				onEditSchedule={openEditSchedule}
				onToggleSchedule={toggleScheduleEnabled}
				onDeleteSchedule={deleteSchedule}
				onAddDelivery={openCreateDelivery}
				onEditDelivery={openEditDelivery}
				onToggleDelivery={toggleDeliveryEnabled}
				onDeleteDelivery={deleteDelivery}
				onRefreshRuns={refreshDetail}
			/>
		</div>
	</section>
</div>

<DefinitionModal
	open={$showDefinitionForm}
	mode={$definitionFormMode}
	form={$definitionForm}
	onFieldChange={updateDefinitionFormField}
	onClose={closeDefinitionForm}
	onSubmit={submitDefinitionForm}
/>

<ScheduleModal
	open={$showScheduleForm}
	mode={$scheduleFormMode}
	form={$scheduleForm}
	onFieldChange={updateScheduleFormField}
	onClose={closeScheduleForm}
	onSubmit={submitScheduleForm}
/>

<DeliveryModal
	open={$showDeliveryForm}
	mode={$deliveryFormMode}
	form={$deliveryForm}
	onFieldChange={updateDeliveryFormField}
	onClose={closeDeliveryForm}
	onSubmit={submitDeliveryForm}
/>
							<p class="max-w-2xl text-sm text-slate-300">
								{detail.definition.description}
							</p>
							<div class="flex flex-wrap items-center gap-3 text-xs uppercase tracking-wide text-slate-500">
								<span>
									Updated {new Date(detail.definition.updated_at).toLocaleString()}
								</span>
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
								on:click={openEditDefinition}
							>
								Edit definition
							</button>
							<button
								class="rounded-lg border border-slate-700/60 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
								type="button"
								on:click={() => handleBulkAction('archive', [detail.definition.id])}
							>
								Archive
							</button>
							<button
								class="rounded-lg border border-slate-700/60 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
								type="button"
								on:click={() => handleBulkAction('restore', [detail.definition.id])}
							>
								Restore
							</button>
						</div>
					</div>

					<div class="grid gap-4 md:grid-cols-2">
						<div class="rounded-xl border border-slate-800/70 bg-slate-900/40 p-4">
							<div class="flex items-center justify-between">
								<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">
									Sections
								</h3>
							</div>
							<pre class="mt-3 max-h-56 overflow-auto rounded-lg bg-slate-950/80 p-3 text-xs text-slate-300">{JSON.stringify(detail.definition.sections ?? {}, null, 2)}</pre>
						</div>
						<div class="rounded-xl border border-slate-800/70 bg-slate-900/40 p-4">
							<div class="flex items-center justify-between">
								<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">
									Filters
								</h3>
							</div>
							<pre class="mt-3 max-h-56 overflow-auto rounded-lg bg-slate-950/80 p-3 text-xs text-slate-300">{JSON.stringify(detail.definition.filters ?? {}, null, 2)}</pre>
						</div>
					</div>

					<div class="grid gap-6 lg:grid-cols-2">
						<div class="flex flex-col gap-3 rounded-xl border border-slate-800/80 bg-slate-900/50 p-4">
							<div class="flex items-center justify-between">
								<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">
									Schedules
								</h3>
								<button
									class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
									type="button"
									on:click={startCreateSchedule}
								>
									Add schedule
								</button>
							</div>
							{#if detailLoading}
								<div class="flex items-center justify-center py-8 text-sm text-slate-400">
									Loading schedules…
								</div>
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
														{schedule.next_run
															? new Date(schedule.next_run).toLocaleString()
															: 'Not scheduled'}
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
														on:click={() => toggleScheduleEnabled(schedule)}
													>
														{schedule.enabled ? 'Enabled' : 'Disabled'}
													</button>
													<div class="flex gap-2">
														<button
															class="rounded-md border border-slate-700/60 px-2 py-1 text-xs text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
															type="button"
															on:click={() => startEditSchedule(schedule)}
														>
															Edit
														</button>
														<button
															class="rounded-md border border-red-500/60 px-2 py-1 text-xs text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
															type="button"
															on:click={() => removeSchedule(schedule)}
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
								<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">
									Deliveries
								</h3>
								<button
									class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
									type="button"
									on:click={startCreateDelivery}
								>
									Add delivery
								</button>
							</div>
							{#if detailLoading}
								<div class="flex items-center justify-center py-8 text-sm text-slate-400">
									Loading deliveries…
								</div>
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
													<div class="text-xs uppercase tracking-wide text-slate-400">
														{delivery.channel}
													</div>
													<div class="mt-1 font-medium text-slate-100">
														{delivery.target}
													</div>
												</div>
												<div class="flex flex-col gap-2">
													<button
														class={`rounded-md border px-2 py-1 text-xs ${
															delivery.enabled
																? 'border-emerald-500/60 text-emerald-200'
																: 'border-slate-600 text-slate-400'
														}`}
														type="button"
														on:click={() => toggleDeliveryEnabled(delivery)}
													>
														{delivery.enabled ? 'Enabled' : 'Disabled'}
													</button>
													<div class="flex gap-2">
														<button
															class="rounded-md border border-slate-700/60 px-2 py-1 text-xs text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
															type="button"
															on:click={() => startEditDelivery(delivery)}
														>
															Edit
														</button>
														<button
															class="rounded-md border border-red-500/60 px-2 py-1 text-xs text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
															type="button"
															on:click={() => removeDelivery(delivery)}
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
							<h3 class="text-sm font-semibold uppercase tracking-wide text-slate-300">
								Run history
							</h3>
							<button
								class="rounded-lg border border-slate-700/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60"
								type="button"
								on:click={() => selectedDefinitionId && loadDetail(selectedDefinitionId)}
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
												<span class="text-slate-400">
													{new Date(run.created_at).toLocaleString()}
												</span>
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
		</div>
	</section>
</div>

{#if showDefinitionForm}
	<div class="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/80 backdrop-blur">
		<form
			class="w-full max-w-3xl space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
			on:submit|preventDefault={handleDefinitionSubmit}
		>
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-lg font-semibold text-slate-100">
						{definitionFormMode === 'create' ? 'Create report definition' : 'Edit report definition'}
					</h3>
					<p class="text-xs text-slate-400">
						Structure your report filters and sections using JSON for flexible dashboards.
					</p>
				</div>
				<button
					class="rounded-full border border-slate-700/70 px-3 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
					type="button"
					on:click={closeDefinitionForm}
				>
					Close
				</button>
			</div>
			<div class="grid gap-4 sm:grid-cols-2">
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Name</span>
					<input
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={definitionForm.name}
						placeholder="Quarterly executive summary"
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Favorite</span>
					<select
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={definitionForm.favorite}
					>
						<option value={true}>Yes</option>
						<option value={false}>No</option>
					</select>
				</label>
			</div>
			<label class="flex flex-col gap-2 text-sm text-slate-200">
				<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
				<textarea
					class="h-20 rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					bind:value={definitionForm.description}
					placeholder="Summarises KPIs and operational metrics for leadership."
				/>
			</label>
			<div class="grid gap-4 md:grid-cols-2">
				<label class="flex h-full flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Sections (JSON)</span>
					<textarea
						class="h-full min-h-[180px] flex-1 rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 font-mono text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={definitionForm.sectionsText}
						required
					/>
				</label>
				<label class="flex h-full flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Filters (JSON)</span>
					<textarea
						class="h-full min-h-[180px] flex-1 rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 font-mono text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={definitionForm.filtersText}
						required
					/>
				</label>
			</div>
			<div class="flex items-center justify-end gap-3">
				<button
					class="rounded-lg border border-slate-700/60 px-4 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={closeDefinitionForm}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white transition hover:bg-primary/80"
					type="submit"
				>
					{definitionFormMode === 'create' ? 'Create' : 'Save changes'}
				</button>
			</div>
		</form>
	</div>
{/if}

{#if showScheduleForm}
	<div class="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/80 backdrop-blur">
		<form
			class="w-full max-w-2xl space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
			on:submit|preventDefault={submitScheduleForm}
		>
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-lg font-semibold text-slate-100">
						{scheduleFormMode === 'create' ? 'Create schedule' : 'Edit schedule'}
					</h3>
					<p class="text-xs text-slate-400">
						Define cron expressions and metadata for automated report generation.
					</p>
				</div>
				<button
					class="rounded-full border border-slate-700/70 px-3 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
					type="button"
					on:click={() => {
						showScheduleForm = false;
						resetScheduleForm();
					}}
				>
					Close
				</button>
			</div>
			<div class="grid gap-4 md:grid-cols-2">
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Cron</span>
					<input
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 font-mono text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={scheduleForm.cron}
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Frequency</span>
					<select
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={scheduleForm.frequency}
					>
						<option value="daily">Daily</option>
						<option value="weekly">Weekly</option>
						<option value="monthly">Monthly</option>
						<option value="custom">Custom</option>
					</select>
				</label>
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Timezone</span>
					<input
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={scheduleForm.timezone}
						required
					/>
				</label>
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Next run (ISO)</span>
					<input
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={scheduleForm.nextRun}
						placeholder="2025-01-15T08:00:00Z"
					/>
				</label>
			</div>
			<label class="flex items-center justify-between gap-3 rounded-lg border border-slate-800/60 bg-slate-900/60 px-4 py-2 text-sm text-slate-300">
				<span>Enabled</span>
				<input type="checkbox" bind:checked={scheduleForm.enabled} />
			</label>
			<label class="flex flex-col gap-2 text-sm text-slate-200">
				<span class="text-xs uppercase tracking-wide text-slate-400">Metadata (JSON)</span>
				<textarea
					class="h-28 rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 font-mono text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					bind:value={scheduleForm.metaText}
				/>
			</label>
			<div class="flex items-center justify-end gap-3">
				<button
					class="rounded-lg border border-slate-700/60 px-4 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={() => {
						showScheduleForm = false;
						resetScheduleForm();
					}}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white transition hover:bg-primary/80"
					type="submit"
				>
					{scheduleFormMode === 'create' ? 'Create schedule' : 'Save schedule'}
				</button>
			</div>
		</form>
	</div>
{/if}

{#if showDeliveryForm}
	<div class="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/80 backdrop-blur">
		<form
			class="w-full max-w-2xl space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
			on:submit|preventDefault={submitDeliveryForm}
		>
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-lg font-semibold text-slate-100">
						{deliveryFormMode === 'create' ? 'Create delivery' : 'Edit delivery'}
					</h3>
					<p class="text-xs text-slate-400">
						Configure destinations for sending generated reports automatically.
					</p>
				</div>
				<button
					class="rounded-full border border-slate-700/70 px-3 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
					type="button"
					on:click={() => {
						showDeliveryForm = false;
						resetDeliveryForm();
					}}
				>
					Close
				</button>
			</div>
			<div class="grid gap-4 md:grid-cols-2">
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Channel</span>
					<select
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={deliveryForm.channel}
					>
						<option value="email">Email</option>
						<option value="slack">Slack</option>
						<option value="webhook">Webhook</option>
						<option value="whatsapp">WhatsApp</option>
					</select>
				</label>
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Target</span>
					<input
						class="rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						bind:value={deliveryForm.target}
						placeholder="ops@company.com"
						required
					/>
				</label>
			</div>
			<label class="flex items-center justify-between gap-3 rounded-lg border border-slate-800/60 bg-slate-900/60 px-4 py-2 text-sm text-slate-300">
				<span>Enabled</span>
				<input type="checkbox" bind:checked={deliveryForm.enabled} />
			</label>
			<label class="flex flex-col gap-2 text-sm text-slate-200">
				<span class="text-xs uppercase tracking-wide text-slate-400">Template metadata (JSON)</span>
				<textarea
					class="h-28 rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 font-mono text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					bind:value={deliveryForm.templateText}
				/>
			</label>
			<div class="flex items-center justify-end gap-3">
				<button
					class="rounded-lg border border-slate-700/60 px-4 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:bg-slate-800/60"
					type="button"
					on:click={() => {
						showDeliveryForm = false;
						resetDeliveryForm();
					}}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white transition hover:bg-primary/80"
					type="submit"
				>
					{deliveryFormMode === 'create' ? 'Create delivery' : 'Save delivery'}
				</button>
			</div>
		</form>
	</div>
{/if}

<style>
	button {
		cursor: pointer;
	}
</style>

