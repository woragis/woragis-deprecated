import { derived, get, writable } from 'svelte/store';
import { onDestroy } from 'svelte';
import { useQueryClient } from '@tanstack/svelte-query';
import { goto } from '$app/navigation';

import type { ReportDefinition, ReportDefinitionDetail, ReportSchedule, UUID } from '$lib/api/types';
import {
	useDeleteReportScheduleMutation,
	useReportDefinitionsQuery,
	useReportDetailQuery,
	useToggleReportScheduleMutation,
	useUpdateReportScheduleMutation
} from '@hooks/reports';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export interface ScheduleFormState {
	cron: string;
	frequency: string;
	timezone: string;
	nextRun: string;
	enabled: boolean;
	metaText: string;
}

export interface ScheduleWithReport {
	schedule: ReportSchedule;
	report: ReportDefinition;
}

function defaultScheduleForm(): ScheduleFormState {
	return {
		cron: '0 8 * * *',
		frequency: 'daily',
		timezone: Intl.DateTimeFormat().resolvedOptions().timeZone ?? 'UTC',
		nextRun: '',
		enabled: true,
		metaText: ''
	};
}

const parseJsonField = (value: string, fieldName: string) => {
	if (!value.trim()) return {};
	try {
		const parsed = JSON.parse(value);
		if (typeof parsed !== 'object' || parsed === null) {
			throw new Error(`${fieldName} must be a JSON object`);
		}
		return parsed as Record<string, unknown>;
	} catch (error) {
		throw new Error(`${fieldName} is not valid JSON`);
	}
};

export function createSchedulesLogic() {
	const errorMessage = writable<string | null>(null);
	const showScheduleForm = writable(false);
	const scheduleFormMode = writable<'create' | 'edit'>('edit');
	const scheduleForm = writable<ScheduleFormState>(defaultScheduleForm());
	const selectedSchedule = writable<ScheduleWithReport | null>(null);
	const selectedReportId = writable<UUID | null>(null);

	const queryClient = useQueryClient();
	const schedulesData = writable<ScheduleWithReport[]>([]);
	const isLoadingData = writable(false);

	// Fetch all report definitions
	const definitionsQuery = useReportDefinitionsQuery(
		writable({
			search: '',
			includeArchived: false,
			favoritesOnly: false,
			channel: '',
			limit: 1000, // Get all reports
			offset: 0
		})
	);

	const definitions = derived(
		definitionsQuery,
		($query) => ($query.data as ReportDefinition[] | undefined) ?? []
	);

	const isLoading = derived(
		[definitionsQuery, isLoadingData],
		([$defsQuery, $isLoadingData]) => $defsQuery.isFetching || $defsQuery.isLoading || $isLoadingData
	);

	const schedules = derived(schedulesData, ($data) => $data);

	// Fetch all schedules by loading details for each definition
	const fetchAllSchedules = async () => {
		const defs = get(definitions);
		if (defs.length === 0) {
			schedulesData.set([]);
			return;
		}

		isLoadingData.set(true);
		try {
			const { listReportSchedules } = await import('$lib/api/reports');
			const allSchedules: ScheduleWithReport[] = [];

			// Fetch schedules for all definitions in parallel
			const schedulePromises = defs.map(async (def) => {
				try {
					const schedules = await listReportSchedules(def.id);
					return schedules.map((schedule: ReportSchedule) => ({
						schedule,
						report: def
					}));
				} catch {
					return [];
				}
			});
			const scheduleArrays = await Promise.all(schedulePromises);
			allSchedules.push(...scheduleArrays.flat());

			// Sort by next run time, then by report name
			allSchedules.sort((a, b) => {
				if (a.schedule.next_run && b.schedule.next_run) {
					return new Date(a.schedule.next_run).getTime() - new Date(b.schedule.next_run).getTime();
				}
				if (a.schedule.next_run) return -1;
				if (b.schedule.next_run) return 1;
				return a.report.name.localeCompare(b.report.name);
			});

			schedulesData.set(allSchedules);
		} catch (error) {
			console.error('Failed to fetch schedules:', error);
			schedulesData.set([]);
		} finally {
			isLoadingData.set(false);
		}
	};

	// Watch definitions and fetch schedules when they change
	let unsubscribeDefinitions: (() => void) | null = null;
	const setupDefinitionsWatcher = () => {
		if (unsubscribeDefinitions) {
			unsubscribeDefinitions();
		}
		unsubscribeDefinitions = definitions.subscribe(() => {
			fetchAllSchedules();
		});
	};

	setupDefinitionsWatcher();

	const toggleScheduleMutation = useToggleReportScheduleMutation();
	const updateScheduleMutation = useUpdateReportScheduleMutation();
	const deleteScheduleMutation = useDeleteReportScheduleMutation();

	const setError = (message: string | null) => {
		errorMessage.set(message);
		if (message) {
			toastError(message);
		}
	};

	const invalidateSchedules = async () => {
		await queryClient.invalidateQueries({ queryKey: ['reports', 'definitions'] });
		// Invalidate all detail queries
		const defs = get(definitions);
		for (const def of defs) {
			await queryClient.invalidateQueries({ queryKey: ['reports', def.id, 'detail'] });
		}
		// Refetch schedules
		await fetchAllSchedules();
	};

	const loadSchedules = async () => {
		await invalidateSchedules();
	};

	const updateScheduleFormField = <K extends keyof ScheduleFormState>(
		field: K,
		value: ScheduleFormState[K]
	) => {
		scheduleForm.update((current) => ({ ...current, [field]: value }));
	};

	const openEditSchedule = (scheduleWithReport: ScheduleWithReport) => {
		selectedSchedule.set(scheduleWithReport);
		selectedReportId.set(scheduleWithReport.report.id);
		scheduleFormMode.set('edit');
		scheduleForm.set({
			cron: scheduleWithReport.schedule.cron,
			frequency: scheduleWithReport.schedule.frequency,
			timezone: scheduleWithReport.schedule.timezone,
			nextRun: scheduleWithReport.schedule.next_run ?? '',
			enabled: scheduleWithReport.schedule.enabled,
			metaText: scheduleWithReport.schedule.meta ? JSON.stringify(scheduleWithReport.schedule.meta, null, 2) : ''
		});
		showScheduleForm.set(true);
	};

	const closeScheduleForm = () => {
		showScheduleForm.set(false);
		scheduleFormMode.set('edit');
		scheduleForm.set(defaultScheduleForm());
		selectedSchedule.set(null);
		selectedReportId.set(null);
	};

	const submitScheduleForm = async () => {
		const schedule = get(selectedSchedule);
		if (!schedule) return;

		const form = get(scheduleForm);
		const meta = form.metaText.trim() ? parseJsonField(form.metaText, 'Meta') : undefined;
		try {
			errorMessage.set(null);
			await get(updateScheduleMutation).mutateAsync({
				scheduleId: schedule.schedule.id,
				payload: {
					cron: form.cron,
					frequency: form.frequency,
					timezone: form.timezone,
					nextRun: form.nextRun || undefined,
					enabled: form.enabled,
					meta
				}
			});
			toastSuccess('Schedule updated.');
			closeScheduleForm();
			await invalidateSchedules();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save schedule.'));
		}
	};

	const toggleScheduleEnabled = async (scheduleWithReport: ScheduleWithReport) => {
		try {
			await get(toggleScheduleMutation).mutateAsync({
				scheduleId: scheduleWithReport.schedule.id,
				enabled: !scheduleWithReport.schedule.enabled
			});
			toastSuccess('Schedule updated.');
			await invalidateSchedules();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to toggle schedule.'));
		}
	};

	const deleteSchedule = async (scheduleWithReport: ScheduleWithReport) => {
		if (!confirm(`Are you sure you want to delete this schedule for "${scheduleWithReport.report.name}"?`)) {
			return;
		}
		try {
			await get(deleteScheduleMutation).mutateAsync(scheduleWithReport.schedule.id);
			toastInfo('Schedule deleted.');
			await invalidateSchedules();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to delete schedule.'));
		}
	};

	onDestroy(() => {
		if (unsubscribeDefinitions) {
			unsubscribeDefinitions();
		}
	});

	return {
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
	};
}

