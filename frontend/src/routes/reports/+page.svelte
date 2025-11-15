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

