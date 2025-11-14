import { createMutation, createQuery } from '@tanstack/svelte-query';

import { reportsApi } from '$lib/api/reports';
import type { ReportDefinition, ReportDefinitionDetail, ReportRun, UUID } from '$lib/api/types';

const createEmptyReportDetail = (): ReportDefinitionDetail => ({
	definition: {
		id: '',
		name: '',
		description: '',
		sections: {},
		filters: {},
		is_favorite: false,
		archived_at: null,
		created_at: '',
		updated_at: ''
	},
	schedules: [],
	deliveries: []
});

export interface ReportDefinitionsOptions {
	search?: string;
	includeArchived?: boolean;
	favoritesOnly?: boolean;
	channel?: string;
	limit?: number;
	offset?: number;
	enabled?: boolean;
}

export const useReportDefinitionsQuery = (options: ReportDefinitionsOptions = {}) => {
	const filters = {
		search: options.search ?? '',
		includeArchived: options.includeArchived ?? false,
		favoritesOnly: options.favoritesOnly ?? false,
		channel: options.channel ?? '',
		limit: options.limit ?? 25,
		offset: options.offset ?? 0
	};

	return createQuery<ReportDefinition[]>({
		queryKey: ['reports', 'definitions', filters],
		queryFn: () =>
			reportsApi.listReportDefinitions({
				search: filters.search,
				includeArchived: filters.includeArchived,
				favorites: filters.favoritesOnly,
				channel: filters.channel || undefined,
				limit: filters.limit,
				offset: filters.offset
			}),
		enabled: options.enabled ?? true,
		placeholderData: () => []
	});
};

export const useReportDetailQuery = (
	definitionId: string | null,
	options: { enabled?: boolean } = {}
) =>
	createQuery<{
		detail: ReportDefinitionDetail;
		runs: ReportRun[];
	}>({
		queryKey: ['reports', definitionId, 'detail'],
		queryFn: async () => {
			const [detailResponse, runsResponse] = await Promise.all([
				reportsApi.getReportDefinition(definitionId!),
				reportsApi.listReportRuns(definitionId!)
			]);
			return {
				detail: detailResponse,
				runs: runsResponse
			};
		},
		enabled: Boolean(definitionId) && (options.enabled ?? true),
		placeholderData: () => ({
			detail: createEmptyReportDetail(),
			runs: []
		})
	});

export const useCreateReportDefinitionMutation = () =>
	createMutation({
		mutationFn: reportsApi.createReportDefinition
	});

export const useUpdateReportDefinitionMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof reportsApi.updateReportDefinition>[1];
		}) => reportsApi.updateReportDefinition(definitionId, payload)
	});

export const useArchiveReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: reportsApi.archiveReportDefinitions
	});

export const useRestoreReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: reportsApi.restoreReportDefinitions
	});

export const useDeleteReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: reportsApi.deleteReportDefinitions
	});

export const useToggleReportFavoriteMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			favorite
		}: {
			definitionId: UUID;
			favorite: boolean;
		}) => reportsApi.toggleReportFavorite(definitionId, favorite)
	});

export const useQueueReportRunsMutation = () =>
	createMutation({
		mutationFn: ({
			definitionIds,
			metadata
		}: {
			definitionIds: UUID[];
			metadata?: Record<string, unknown>;
		}) => reportsApi.queueReportRuns(definitionIds, metadata)
	});

export const useCreateReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof reportsApi.createReportSchedule>[1];
		}) => reportsApi.createReportSchedule(definitionId, payload)
	});

export const useUpdateReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			scheduleId,
			payload
		}: {
			scheduleId: UUID;
			payload: Parameters<typeof reportsApi.updateReportSchedule>[1];
		}) => reportsApi.updateReportSchedule(scheduleId, payload)
	});

export const useToggleReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			scheduleId,
			enabled
		}: {
			scheduleId: UUID;
			enabled: boolean;
		}) => reportsApi.toggleReportSchedule(scheduleId, enabled)
	});

export const useDeleteReportScheduleMutation = () =>
	createMutation({
		mutationFn: (scheduleId: UUID) => reportsApi.deleteReportSchedule(scheduleId)
	});

export const useCreateReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof reportsApi.createReportDelivery>[1];
		}) => reportsApi.createReportDelivery(definitionId, payload)
	});

export const useUpdateReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			deliveryId,
			payload
		}: {
			deliveryId: UUID;
			payload: Parameters<typeof reportsApi.updateReportDelivery>[1];
		}) => reportsApi.updateReportDelivery(deliveryId, payload)
	});

export const useToggleReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			deliveryId,
			enabled
		}: {
			deliveryId: UUID;
			enabled: boolean;
		}) => reportsApi.toggleReportDelivery(deliveryId, enabled)
	});

export const useDeleteReportDeliveryMutation = () =>
	createMutation({
		mutationFn: (deliveryId: UUID) => reportsApi.deleteReportDelivery(deliveryId)
	});

