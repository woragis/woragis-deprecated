<script lang="ts">
	import { onMount } from 'svelte';
	import { createSchedulesLogic } from './schedules.logic';
	import SchedulesHeader from './_components/SchedulesHeader.svelte';
	import SchedulesList from './_components/SchedulesList.svelte';
	import ScheduleModal from '../reports/_components/ScheduleModal.svelte';

	const {
		isLoading,
		schedules,
		errorMessage,
		showScheduleForm,
		scheduleFormMode,
		scheduleForm,
		selectedSchedule,
		selectedReportId,
		updateScheduleFormField,
		loadSchedules,
		openEditSchedule,
		closeScheduleForm,
		submitScheduleForm,
		toggleScheduleEnabled,
		deleteSchedule
	} = createSchedulesLogic();

	onMount(async () => {
		await loadSchedules();
	});
</script>

<svelte:head>
	<title>Schedules · Woragis</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SchedulesHeader isLoading={$isLoading} onRefresh={loadSchedules} />

	{#if $errorMessage}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
			{$errorMessage}
		</p>
	{/if}

	<SchedulesList
		schedules={$schedules}
		isLoading={$isLoading}
		onEditSchedule={openEditSchedule}
		onToggleSchedule={toggleScheduleEnabled}
		onDeleteSchedule={deleteSchedule}
	/>
</div>

<ScheduleModal
	open={$showScheduleForm}
	mode={$scheduleFormMode}
	form={$scheduleForm}
	onFieldChange={updateScheduleFormField}
	onClose={closeScheduleForm}
	onSubmit={submitScheduleForm}
/>

