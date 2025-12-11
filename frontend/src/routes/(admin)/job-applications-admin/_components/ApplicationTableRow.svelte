<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { JobApplication } from '$lib/api/jobapplications';
	
	let {
		application,
		onDelete
	}: {
		application: JobApplication;
		onDelete?: (id: string) => void;
	} = $props();
	
	const tFn = useTranslation();
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}

	function formatDateTime(dateString?: string): string {
		if (!dateString) return '—';
		const date = new Date(dateString);
		return date.toLocaleString();
	}

	function getStatusInfo(application: JobApplication): { text: string; timestamp?: string } {
		switch (application.status) {
			case 'processing':
				return {
					text: 'Processing...',
					timestamp: application.updatedAt ? formatDateTime(application.updatedAt) : undefined
				};
			case 'applied':
				return {
					text: 'Applied',
					timestamp: application.appliedAt ? formatDateTime(application.appliedAt) : undefined
				};
			case 'failed':
				return {
					text: 'Failed',
					timestamp: application.updatedAt ? formatDateTime(application.updatedAt) : undefined
				};
			case 'contacted':
				return {
					text: 'Contacted',
					timestamp: application.responseReceivedAt ? formatDateTime(application.responseReceivedAt) : undefined
				};
			case 'rejected':
				return {
					text: 'Rejected',
					timestamp: application.updatedAt ? formatDateTime(application.updatedAt) : undefined
				};
			case 'accepted':
				return {
					text: 'Accepted',
					timestamp: application.updatedAt ? formatDateTime(application.updatedAt) : undefined
				};
			default:
				return {
					text: 'Pending',
					timestamp: application.createdAt ? formatDateTime(application.createdAt) : undefined
				};
		}
	}

	const statusInfo = $derived(getStatusInfo(application));
</script>

<tr class="table-row">
	<td><strong>{application.companyName}</strong></td>
	<td>{application.jobTitle}</td>
	<td>{application.location || '—'}</td>
	<td>{application.website}</td>
	<td>{application.language || '—'}</td>
	<td>
		<div class="status-cell">
			<StatusBadge status={application.status} type="status">
				{statusInfo.text}
			</StatusBadge>
			{#if statusInfo.timestamp}
				<span class="status-timestamp" title={statusInfo.timestamp}>
					{statusInfo.timestamp}
				</span>
			{/if}
		</div>
	</td>
	<td>{application.interestLevel || '—'}</td>
	<td>{formatDate(application.appliedAt)}</td>
	<td class="actions-cell">
		<a href="/job-applications-admin/{application.id}" class="view-link">
			{tFn('jobApplications.table.view')}
		</a>
		<Button variant="danger" size="sm" onclick={() => onDelete?.(application.id)}>
			{tFn('jobApplications.table.delete')}
		</Button>
	</td>
</tr>

<style>
	.table-row {
		transition: background-color var(--transition-base);
	}
	
	.table-row:hover {
		background-color: var(--color-bg-hover);
	}
	
	.table-row td {
		padding: var(--spacing-md);
		border-bottom: 1px solid var(--color-border);
		color: var(--color-text-primary);
	}
	
	.actions-cell {
		display: flex;
		gap: var(--spacing-sm);
		align-items: center;
	}
	
	.view-link {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-success);
		color: white;
		border-radius: var(--radius-md);
		text-decoration: none;
		font-size: var(--font-size-sm);
		transition: background-color var(--transition-base);
		display: inline-block;
	}
	
	.view-link:hover {
		background-color: var(--color-success-hover);
	}

	.status-cell {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.status-timestamp {
		font-size: var(--font-size-xs);
		color: var(--color-text-secondary);
		margin-top: var(--spacing-xs);
	}
</style>
