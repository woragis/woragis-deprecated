import { get, writable } from 'svelte/store';

import type {
	ReportDefinition,
	ReportDefinitionDetail,
	ReportDelivery,
	ReportRun,
	ReportSchedule,
	UUID
} from '$lib/api/types';
import {
	archiveReportDefinitions,
	createReportDefinition,
	createReportDelivery,
	createReportSchedule,
	deleteReportDefinitions,
	deleteReportDelivery,
	deleteReportSchedule,
	getReportDefinition,
	listReportDefinitions,
	listReportRuns,
	queueReportRuns,
	restoreReportDefinitions,
	toggleReportDelivery,
	toggleReportFavorite,
	toggleReportSchedule,
	updateReportDefinition,
	updateReportDelivery,
	updateReportSchedule
} from '$lib/api/reports';
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

export function createReportsLogic() {
	const isLoading = writable(false);
	const detailLoading = writable(false);
	const definitions = writable<ReportDefinition[]>([]);
	const detail = writable<ReportDefinitionDetail | null>(null);
	const runs = writable<ReportRun[]>([]);
	const filters = writable<ReportFilters>(defaultFilters());
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

	const loadDefinitions = async () => {
		isLoading.set(true);
		try {
			const currentFilters = get(filters);
			const data = await listReportDefinitions({
				search: currentFilters.search,
				includeArchived: currentFilters.includeArchived,
				favorites: currentFilters.favoritesOnly,
				channel: currentFilters.channel || undefined,
				limit: currentFilters.limit,
				offset: currentFilters.offset
			});
			definitions.set(data);
			selectedIds.update((set) => {
				const existingIds = new Set(data.map((item) => item.id));
				return new Set([...set].filter((id) => existingIds.has(id)));
			});
			const activeId = get(selectedDefinitionId);
			if (activeId && !data.find((definition) => definition.id === activeId)) {
				selectedDefinitionId.set(null);
				detail.set(null);
				runs.set([]);
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to load report definitions.'));
		} finally {
			isLoading.set(false);
		}
	};

	const loadDetail = async (definitionId: UUID) => {
		detailLoading.set(true);
		try {
			const [detailResponse, runsResponse] = await Promise.all([
				getReportDefinition(definitionId),
				listReportRuns(definitionId)
			]);
			detail.set(detailResponse);
			runs.set(runsResponse);
			const definition = detailResponse.definition;
			if (definition) {
				definitionForm.set({
					name: definition.name,
					description: definition.description ?? '',
					sectionsText: JSON.stringify(definition.sections ?? {}, null, 2),
					filtersText: JSON.stringify(definition.filters ?? {}, null, 2),
					favorite: definition.is_favorite
				});
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to load report detail.'));
		} finally {
			detailLoading.set(false);
		}
	};

	const selectDefinition = async (definitionId: UUID) => {
		if (get(selectedDefinitionId) === definitionId) return;
		selectedDefinitionId.set(definitionId);
		await loadDetail(definitionId);
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
		if (!get(detail)?.definition) return;
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
			if (get(definitionFormMode) === 'create') {
				const created = await createReportDefinition({
					name: form.name,
					description: form.description,
					sections,
					filters: reportFilters,
					favorite: form.favorite
				});
				toastSuccess('Report definition created.');
				closeDefinitionForm();
				await loadDefinitions();
				await selectDefinition(created.id);
			} else if (get(definitionFormMode) === 'edit' && get(selectedDefinitionId)) {
				const definitionId = get(selectedDefinitionId)!;
				await updateReportDefinition(definitionId, {
					name: form.name,
					description: form.description,
					sections,
					filters: reportFilters,
					favorite: form.favorite
				});
				toastSuccess('Report definition updated.');
				closeDefinitionForm();
				await loadDefinitions();
				await loadDetail(definitionId);
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
			if (action === 'archive') {
				await archiveReportDefinitions(ids);
				toastSuccess('Definitions archived.');
			} else if (action === 'restore') {
				await restoreReportDefinitions(ids);
				toastSuccess('Definitions restored.');
			} else {
				await deleteReportDefinitions(ids);
				toastSuccess('Definitions deleted.');
			}
			await loadDefinitions();
			const activeId = get(selectedDefinitionId);
			if (activeId && ids.includes(activeId)) {
				selectedDefinitionId.set(null);
				detail.set(null);
				runs.set([]);
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
			await queueReportRuns(ids, metadata);
			toastSuccess('Runs queued successfully.');
			if (get(selectedDefinitionId)) {
				await loadDetail(get(selectedDefinitionId)!);
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to queue runs.'));
		}
	};

	const toggleFavorite = async (definitionId: UUID, currentFavorite: boolean) => {
		try {
			await toggleReportFavorite(definitionId, !currentFavorite);
			toastInfo(!currentFavorite ? 'Marked as favorite.' : 'Removed from favorites.');
			await loadDefinitions();
			if (get(selectedDefinitionId) === definitionId) {
				await loadDetail(definitionId);
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
				await createReportSchedule(definitionId, {
					cron: form.cron,
					frequency: form.frequency,
					timezone: form.timezone,
					nextRun: form.nextRun || undefined,
					enabled: form.enabled,
					meta
				});
				toastSuccess('Schedule created.');
			} else if (get(scheduleFormMode) === 'edit' && get(editingScheduleId)) {
				await updateReportSchedule(get(editingScheduleId)!, {
					cron: form.cron,
					frequency: form.frequency,
					timezone: form.timezone,
					nextRun: form.nextRun || undefined,
					enabled: form.enabled,
					meta
				});
				toastSuccess('Schedule updated.');
			}
			closeScheduleForm();
			await loadDetail(definitionId);
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save schedule.'));
		}
	};

	const toggleScheduleEnabled = async (schedule: ReportSchedule) => {
		try {
			await toggleReportSchedule(schedule.id, !schedule.enabled);
			toastSuccess('Schedule updated.');
			if (get(selectedDefinitionId)) {
				await loadDetail(get(selectedDefinitionId)!);
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to toggle schedule.'));
		}
	};

	const deleteSchedule = async (schedule: ReportSchedule) => {
		try {
			await deleteReportSchedule(schedule.id);
			toastInfo('Schedule deleted.');
			if (get(selectedDefinitionId)) {
				await loadDetail(get(selectedDefinitionId)!);
			}
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
				await createReportDelivery(definitionId, {
					channel: form.channel,
					target: form.target,
					template,
					enabled: form.enabled
				});
				toastSuccess('Delivery created.');
			} else if (get(deliveryFormMode) === 'edit' && get(editingDeliveryId)) {
				await updateReportDelivery(get(editingDeliveryId)!, {
					channel: form.channel,
					target: form.target,
					template,
					enabled: form.enabled
				});
				toastSuccess('Delivery updated.');
			}
			closeDeliveryForm();
			await loadDetail(definitionId);
		} catch (error) {
			setError(getApiErrorMessage(error, 'Unable to save delivery.'));
		}
	};

	const toggleDeliveryEnabled = async (delivery: ReportDelivery) => {
		try {
			await toggleReportDelivery(delivery.id, !delivery.enabled);
			toastSuccess('Delivery updated.');
			if (get(selectedDefinitionId)) {
				await loadDetail(get(selectedDefinitionId)!);
			}
		} catch (error) {
			setError(getApiErrorMessage(error, 'Failed to toggle delivery.'));
		}
	};

	const deleteDelivery = async (delivery: ReportDelivery) => {
		try {
			await deleteReportDelivery(delivery.id);
			toastInfo('Delivery deleted.');
			if (get(selectedDefinitionId)) {
				await loadDetail(get(selectedDefinitionId)!);
			}
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

