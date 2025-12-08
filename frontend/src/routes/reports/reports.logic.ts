import { derived, get, writable } from 'svelte/store';
import { onDestroy } from 'svelte';
import { useQueryClient } from '@tanstack/svelte-query';

import type {
	ReportDefinition,
	ReportDefinitionDetail,
	ReportDelivery,
	ReportRun,
	ReportSchedule,
	UUID
} from '$lib/api/types';
import {
	useArchiveReportDefinitionsMutation,
	useCreateReportDefinitionMutation,
	useCreateReportDeliveryMutation,
	useCreateReportScheduleMutation,
	useDeleteReportDefinitionsMutation,
	useDeleteReportDeliveryMutation,
	useDeleteReportScheduleMutation,
	useQueueReportRunsMutation,
	useReportDefinitionsQuery,
	useReportDetailQuery,
	useRestoreReportDefinitionsMutation,
	useToggleReportDeliveryMutation,
	useToggleReportFavoriteMutation,
	useToggleReportScheduleMutation,
	useUpdateReportDefinitionMutation,
	useUpdateReportDeliveryMutation,
	useUpdateReportScheduleMutation
} from '@hooks/reports';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

export interface ReportFilters {
	search: string;
	includeArchived: boolean;
	favoritesOnly: boolean;
	channel: string;
	limit: number;
	offset: number;
}

export interface DefinitionFormState {
	name: string;
	description: string;
	sectionsText: string;
	filtersText: string;
	favorite: boolean;
}

export interface ScheduleFormState {
	cron: string;
	frequency: string;
	timezone: string;
	nextRun: string;
	enabled: boolean;
	metaText: string;
}

export interface DeliveryFormState {
	channel: string;
	target: string;
	templateText: string;
	enabled: boolean;
}

function defaultFilters(): ReportFilters {
	return {
		search: '',
		includeArchived: false,
		favoritesOnly: false,
		channel: '',
		limit: 25,
		offset: 0
	};
}

const DEFAULT_SECTIONS = '{\n  "overview": []\n}';
const DEFAULT_FILTERS = '{\n  "date_range": "last_30_days"\n}';

function defaultDefinitionForm(): DefinitionFormState {
	return {
		name: '',
		description: '',
		sectionsText: DEFAULT_SECTIONS,
		filtersText: DEFAULT_FILTERS,
		favorite: false
	};
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

function defaultDeliveryForm(): DeliveryFormState {
	return {
		channel: 'email',
		target: '',
		templateText: '',
		enabled: true
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

const areFiltersEqual = (a: ReportFilters, b: ReportFilters) =>
	a.search === b.search &&
	a.includeArchived === b.includeArchived &&
	a.favoritesOnly === b.favoritesOnly &&
	a.channel === b.channel &&
	a.limit === b.limit &&
	a.offset === b.offset;

export function createReportsLogic() {
	const filters = writable<ReportFilters>(defaultFilters());
	const appliedFilters = writable<ReportFilters>(defaultFilters());
	const selectedDefinitionId = writable<UUID | null>(null);
	const selectedIds = writable<Set<UUID>>(new Set());
	const queueMetadataText = writable('');

	const showDefinitionForm = writable(false);
	const definitionFormMode = writable<'create' | 'edit'>('create');
	const definitionForm = writable<DefinitionFormState>(defaultDefinitionForm());

	const showScheduleForm = writable(false);
	const scheduleFormMode = writable<'create' | 'edit'>('create');
	const editingScheduleId = writable<UUID | null>(null);
	const scheduleForm = writable<ScheduleFormState>(defaultScheduleForm());

	const showDeliveryForm = writable(false);
	const deliveryFormMode = writable<'create' | 'edit'>('create');
	const editingDeliveryId = writable<UUID | null>(null);
	const deliveryForm = writable<DeliveryFormState>(defaultDeliveryForm());

	const errorMessage = writable<string | null>(null);

	const queryClient = useQueryClient();

	const definitionsOptionsStore = derived(appliedFilters, ($filters) => ({
		search: $filters.search,
		includeArchived: $filters.includeArchived,
		favoritesOnly: $filters.favoritesOnly,
		channel: $filters.channel,
		limit: $filters.limit,
		offset: $filters.offset
	}));

	const detailOptionsStore = derived(selectedDefinitionId, ($id) => ({
		definitionId: $id,
		enabled: Boolean($id)
	}));

	const definitionsQuery = useReportDefinitionsQuery(definitionsOptionsStore);
	const detailQuery = useReportDetailQuery(detailOptionsStore);

	const definitions = derived(
		definitionsQuery,
		($query) => ($query.data as ReportDefinition[] | undefined) ?? []
	);
	const isLoading = derived(
		definitionsQuery,
		($query) => $query.isFetching || $query.isLoading
	);
	const detail = derived(
		[detailQuery, selectedDefinitionId],
		([$query, $selectedId]) => ($selectedId ? $query.data?.detail ?? null : null)
	);
	const runs = derived(
		[detailQuery, selectedDefinitionId],
		([$query, $selectedId]) => ($selectedId ? $query.data?.runs ?? [] : [])
	);
	const detailLoading = derived(
		detailQuery,
		($query) => $query.isFetching || $query.isLoading
	);

	const definitionsUnsubscribe = definitions.subscribe((items) => {
		selectedIds.update((current) => {
			const next = new Set<UUID>();
			for (const item of items) {
				if (current.has(item.id)) {
					next.add(item.id);
				}
			}
			return next;
		});

		const activeId = get(selectedDefinitionId);
		if (activeId && !items.find((definition) => definition.id === activeId)) {
			selectedDefinitionId.set(null);
		}
	});

	onDestroy(() => {
		definitionsUnsubscribe();
	});

	const createDefinitionMutation = useCreateReportDefinitionMutation();
	const updateDefinitionMutation = useUpdateReportDefinitionMutation();
	const archiveDefinitionsMutation = useArchiveReportDefinitionsMutation();
	const restoreDefinitionsMutation = useRestoreReportDefinitionsMutation();
	const deleteDefinitionsMutation = useDeleteReportDefinitionsMutation();
	const toggleFavoriteMutation = useToggleReportFavoriteMutation();
	const queueRunsMutation = useQueueReportRunsMutation();
	const createScheduleMutation = useCreateReportScheduleMutation();
	const updateScheduleMutation = useUpdateReportScheduleMutation();
	const toggleScheduleMutation = useToggleReportScheduleMutation();
	const deleteScheduleMutation = useDeleteReportScheduleMutation();
	const createDeliveryMutation = useCreateReportDeliveryMutation();
	const updateDeliveryMutation = useUpdateReportDeliveryMutation();
	const toggleDeliveryMutation = useToggleReportDeliveryMutation();
	const deleteDeliveryMutation = useDeleteReportDeliveryMutation();

	const updateDefinitionFormField = <K extends keyof DefinitionFormState>(
		field: K,
		value: DefinitionFormState[K]
	) => {
		definitionForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateScheduleFormField = <K extends keyof ScheduleFormState>(
		field: K,
		value: ScheduleFormState[K]
	) => {
		scheduleForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateDeliveryFormField = <K extends keyof DeliveryFormState>(
		field: K,
		value: DeliveryFormState[K]
	) => {
		deliveryForm.update((current) => ({ ...current, [field]: value }));
	};

	const updateQueueMetadata = (value: string) => {
		queueMetadataText.set(value);
	};

	const setError = (message: string | null) => {
		errorMessage.set(message);
		if (message) {
			toastError(message);
		}
	};

	const invalidateDefinitions = async () => {
		await queryClient.invalidateQueries({ queryKey: ['reports', 'definitions'] });
	};

	const invalidateDetail = async (definitionId: UUID | null = get(selectedDefinitionId)) => {
		if (!definitionId) return;
		await queryClient.invalidateQueries({ queryKey: ['reports', definitionId, 'detail'] });
	};

	const loadDefinitions = async () => {
		const next = { ...get(filters) };
		const currentApplied = get(appliedFilters);
		if (areFiltersEqual(next, currentApplied)) {
			return invalidateDefinitions();
		}
		appliedFilters.set(next);
	};

	const refreshDetail = async () => {
		await invalidateDetail();
	};

	const selectDefinition = (definitionId: UUID) => {
		if (get(selectedDefinitionId) === definitionId) return;
		selectedDefinitionId.set(definitionId);
	};

	const toggleSelection = (definitionId: UUID) => {
		selectedIds.update((set) => {
			const next = new Set(set);
			if (next.has(definitionId)) {
				next.delete(definitionId);
			} else {
				next.add(definitionId);
			}
			return next;
		});
	};

	const toggleSelectAll = () => {
		const currentDefinitions = get(definitions);
		const currentSelections = get(selectedIds);
		if (currentDefinitions.length === 0) {
			selectedIds.set(new Set());
			return;
		}
		if (currentSelections.size === currentDefinitions.length) {
			selectedIds.set(new Set());
		} else {
			selectedIds.set(new Set(currentDefinitions.map((item) => item.id)));
		}
	};

	const updateFilters = (partial: Partial<ReportFilters>) => {
		filters.update((current) => ({ ...current, ...partial }));
	};

	const openCreateDefinition = () => {
		definitionFormMode.set('create');
		definitionForm.set(defaultDefinitionForm());
		showDefinitionForm.set(true);
	};

	const openEditDefinition = () => {
		const currentDetail = get(detail);
		if (!currentDetail?.definition) return;
		definitionForm.set({
			name: currentDetail.definition.name,
			description: currentDetail.definition.description ?? '',
			sectionsText: JSON.stringify(currentDetail.definition.sections ?? {}, null, 2),
			filtersText: JSON.stringify(currentDetail.definition.filters ?? {}, null, 2),
			favorite: currentDetail.definition.is_favorite
		});
		definitionFormMode.set('edit');
		showDefinitionForm.set(true);
	};

	const closeDefinitionForm = () => {
		showDefinitionForm.set(false);
		definitionFormMode.set('create');
		definitionForm.set(defaultDefinitionForm());
	};

	const submitDefinitionForm = async () => {
		const form = get(definitionForm);
		const sections = parseJsonField(form.sectionsText, 'Sections');
		const reportFilters = parseJsonField(form.filtersText, 'Filters');
		try {
			errorMessage.set(null);
			if (get(definitionFormMode) === 'create') {
				const created = await get(createDefinitionMutation).mutateAsync({
					name: form.name,
					description: form.description,
					sections,
					filters: reportFilters,
					isFavorite: form.favorite
				});
				toastSuccess('Report definition created.');
				closeDefinitionForm();
				appliedFilters.set({ ...get(filters) });
				await loadDefinitions();
				await selectDefinition(created.id);
			} else if (get(definitionFormMode) === 'edit' && get(selectedDefinitionId)) {
				const definitionId = get(selectedDefinitionId)!;
				await get(updateDefinitionMutation).mutateAsync({
					definitionId,
					payload: {
						name: form.name,
						description: form.description,
						sections,
						filters: reportFilters,
						isFavorite: form.favorite
					}
				});
				toastSuccess('Report definition updated.');
				closeDefinitionForm();
				await invalidateDefinitions();
				await refreshDetail();
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save report definition.'));
		}
	};

	const getSelectedIds = () => new Set(get(selectedIds));

	const handleBulkAction = async (action: 'archive' | 'restore' | 'delete') => {
		const ids = Array.from(getSelectedIds());
		if (ids.length === 0) {
			toastInfo('Select at least one definition.');
			return;
		}
		try {
			errorMessage.set(null);
			if (action === 'archive') {
				await get(archiveDefinitionsMutation).mutateAsync(ids);
				toastSuccess('Definitions archived.');
			} else if (action === 'restore') {
				await get(restoreDefinitionsMutation).mutateAsync(ids);
				toastSuccess('Definitions restored.');
			} else {
				await get(deleteDefinitionsMutation).mutateAsync(ids);
				toastSuccess('Definitions deleted.');
			}
			await invalidateDefinitions();
			const activeId = get(selectedDefinitionId);
			if (activeId && ids.includes(activeId)) {
				selectedDefinitionId.set(null);
			} else if (activeId) {
				await refreshDetail();
			}
			selectedIds.set(new Set());
		} catch (error) {
			setError(getApiErrorMessage(error, 'Bulk action failed.'));
		}
	};

	const queueSelectedRuns = async () => {
		const ids = Array.from(getSelectedIds());
		if (ids.length === 0) {
			toastInfo('Select at least one definition.');
			return;
		}
		try {
			const metadataText = get(queueMetadataText);
			const metadata = metadataText.trim() ? parseJsonField(metadataText, 'Metadata') : undefined;
			await get(queueRunsMutation).mutateAsync({ definitionIds: ids, metadata });
			toastSuccess('Runs queued successfully.');
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to queue runs.'));
		}
	};

	const toggleFavorite = async (definitionId: UUID, currentFavorite: boolean) => {
		try {
			await get(toggleFavoriteMutation).mutateAsync({
				definitionId,
				favorite: !currentFavorite
			});
			toastInfo(!currentFavorite ? 'Marked as favorite.' : 'Removed from favorites.');
			await invalidateDefinitions();
			if (get(selectedDefinitionId) === definitionId) {
				await refreshDetail();
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to update favorite.'));
		}
	};
	const resetScheduleForm = () => {
		scheduleForm.set(defaultScheduleForm());
		scheduleFormMode.set('create');
		editingScheduleId.set(null);
	};

	const openCreateSchedule = () => {
		resetScheduleForm();
		showScheduleForm.set(true);
	};

	const openEditSchedule = (schedule: ReportSchedule) => {
		scheduleFormMode.set('edit');
		editingScheduleId.set(schedule.id);
		scheduleForm.set({
			cron: schedule.cron,
			frequency: schedule.frequency,
			timezone: schedule.timezone,
			nextRun: schedule.next_run ?? '',
			enabled: schedule.enabled,
			metaText: schedule.meta ? JSON.stringify(schedule.meta, null, 2) : ''
		});
		showScheduleForm.set(true);
	};

	const closeScheduleForm = () => {
		showScheduleForm.set(false);
		resetScheduleForm();
	};

	const submitScheduleForm = async () => {
		const definitionId = get(selectedDefinitionId);
		if (!definitionId) return;
		const form = get(scheduleForm);
		const meta = form.metaText.trim() ? parseJsonField(form.metaText, 'Meta') : undefined;
		try {
			if (get(scheduleFormMode) === 'create') {
				await get(createScheduleMutation).mutateAsync({
					definitionId,
					payload: {
						cron: form.cron,
						frequency: form.frequency,
						timezone: form.timezone,
						nextRun: form.nextRun || undefined,
						enabled: form.enabled,
						meta
					}
				});
				toastSuccess('Schedule created.');
			} else if (get(scheduleFormMode) === 'edit') {
				const scheduleId = get(editingScheduleId);
				if (!scheduleId) return;
				await get(updateScheduleMutation).mutateAsync({
					scheduleId,
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
			}
			closeScheduleForm();
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save schedule.'));
		}
	};

	const toggleScheduleEnabled = async (schedule: ReportSchedule) => {
		try {
			await get(toggleScheduleMutation).mutateAsync({
				scheduleId: schedule.id,
				enabled: !schedule.enabled
			});
			toastSuccess('Schedule updated.');
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to toggle schedule.'));
		}
	};

	const deleteSchedule = async (schedule: ReportSchedule) => {
		try {
			await get(deleteScheduleMutation).mutateAsync(schedule.id);
			toastInfo('Schedule deleted.');
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to delete schedule.'));
		}
	};

	const resetDeliveryForm = () => {
		deliveryForm.set(defaultDeliveryForm());
		deliveryFormMode.set('create');
		editingDeliveryId.set(null);
	};

	const openCreateDelivery = () => {
		resetDeliveryForm();
		showDeliveryForm.set(true);
	};

	const openEditDelivery = (delivery: ReportDelivery) => {
		deliveryFormMode.set('edit');
		editingDeliveryId.set(delivery.id);
		deliveryForm.set({
			channel: delivery.channel,
			target: delivery.target,
			templateText: delivery.template ? JSON.stringify(delivery.template, null, 2) : '',
			enabled: delivery.enabled
		});
		showDeliveryForm.set(true);
	};

	const closeDeliveryForm = () => {
		showDeliveryForm.set(false);
		resetDeliveryForm();
	};

	const submitDeliveryForm = async () => {
		const definitionId = get(selectedDefinitionId);
		if (!definitionId) return;
		const form = get(deliveryForm);
		const template = form.templateText.trim() ? parseJsonField(form.templateText, 'Template') : undefined;
		try {
			if (get(deliveryFormMode) === 'create') {
				await get(createDeliveryMutation).mutateAsync({
					definitionId,
					payload: {
						channel: form.channel,
						target: form.target,
						template,
						enabled: form.enabled
					}
				});
				toastSuccess('Delivery created.');
			} else if (get(deliveryFormMode) === 'edit') {
				const deliveryId = get(editingDeliveryId);
				if (!deliveryId) return;
				await get(updateDeliveryMutation).mutateAsync({
					deliveryId,
					payload: {
						channel: form.channel,
						target: form.target,
						template,
						enabled: form.enabled
					}
				});
				toastSuccess('Delivery updated.');
			}
			closeDeliveryForm();
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save delivery.'));
		}
	};

	const toggleDeliveryEnabled = async (delivery: ReportDelivery) => {
		try {
			await get(toggleDeliveryMutation).mutateAsync({
				deliveryId: delivery.id,
				enabled: !delivery.enabled
			});
			toastSuccess('Delivery updated.');
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to toggle delivery.'));
		}
	};

	const deleteDelivery = async (delivery: ReportDelivery) => {
		try {
			await get(deleteDeliveryMutation).mutateAsync(delivery.id);
			toastInfo('Delivery deleted.');
			await refreshDetail();
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to delete delivery.'));
		}
	};

	return {
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
		refreshDetail,
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
	};
}

