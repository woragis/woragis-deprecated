import { createMutation, createQuery } from '@tanstack/svelte-query';

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
import type { ReportDefinition, ReportDefinitionDetail, ReportRun, UUID } from '$lib/api/types';

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
			listReportDefinitions({
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
				getReportDefinition(definitionId!),
				listReportRuns(definitionId!)
			]);
			return {
				detail: detailResponse,
				runs: runsResponse
			};
		},
		enabled: Boolean(definitionId) && (options.enabled ?? true),
		placeholderData: () => ({
			detail: { definition: null, schedules: [], deliveries: [] },
			runs: []
		})
	});

export const useCreateReportDefinitionMutation = () =>
	createMutation({
		mutationFn: createReportDefinition
	});

export const useUpdateReportDefinitionMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof updateReportDefinition>[1];
		}) => updateReportDefinition(definitionId, payload)
	});

export const useArchiveReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: archiveReportDefinitions
	});

export const useRestoreReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: restoreReportDefinitions
	});

export const useDeleteReportDefinitionsMutation = () =>
	createMutation({
		mutationFn: deleteReportDefinitions
	});

export const useToggleReportFavoriteMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			favorite
		}: {
			definitionId: UUID;
			favorite: boolean;
		}) => toggleReportFavorite(definitionId, favorite)
	});

export const useQueueReportRunsMutation = () =>
	createMutation({
		mutationFn: ({
			definitionIds,
			metadata
		}: {
			definitionIds: UUID[];
			metadata?: Record<string, unknown>;
		}) => queueReportRuns(definitionIds, metadata)
	});

export const useCreateReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof createReportSchedule>[1];
		}) => createReportSchedule(definitionId, payload)
	});

export const useUpdateReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			scheduleId,
			payload
		}: {
			scheduleId: UUID;
			payload: Parameters<typeof updateReportSchedule>[1];
		}) => updateReportSchedule(scheduleId, payload)
	});

export const useToggleReportScheduleMutation = () =>
	createMutation({
		mutationFn: ({
			scheduleId,
			enabled
		}: {
			scheduleId: UUID;
			enabled: boolean;
		}) => toggleReportSchedule(scheduleId, enabled)
	});

export const useDeleteReportScheduleMutation = () =>
	createMutation({
		mutationFn: (scheduleId: UUID) => deleteReportSchedule(scheduleId)
	});

export const useCreateReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			definitionId,
			payload
		}: {
			definitionId: UUID;
			payload: Parameters<typeof createReportDelivery>[1];
		}) => createReportDelivery(definitionId, payload)
	});

export const useUpdateReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			deliveryId,
			payload
		}: {
			deliveryId: UUID;
			payload: Parameters<typeof updateReportDelivery>[1];
		}) => updateReportDelivery(deliveryId, payload)
	});

export const useToggleReportDeliveryMutation = () =>
	createMutation({
		mutationFn: ({
			deliveryId,
			enabled
		}: {
			deliveryId: UUID;
			enabled: boolean;
		}) => toggleReportDelivery(deliveryId, enabled)
	});

export const useDeleteReportDeliveryMutation = () =>
	createMutation({
		mutationFn: (deliveryId: UUID) => deleteReportDelivery(deliveryId)
	});

