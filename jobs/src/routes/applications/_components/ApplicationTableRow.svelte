<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import { ExternalLink, Eye } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import type { JobApplication } from '$lib/api/jobapplications';
	
	let {
		application,
		onDelete,
		selected = false,
		onToggleSelection,
		onView,
		onEdit
	}: {
		application: JobApplication;
		onDelete?: (id: string) => void;
		selected?: boolean;
		onToggleSelection?: (id: string) => void;
		onView?: (application: JobApplication) => void;
		onEdit?: (application: JobApplication) => void;
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

	function getStatusInfo(application: JobApplication): { text: string; timestamp?: string; error?: string } {
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
					timestamp: application.updatedAt ? formatDateTime(application.updatedAt) : undefined,
					error: application.errorMessage
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
	let showFullError = $state(false);
	
	// Reset error display when application changes
	$effect(() => {
		if (application.id) {
			showFullError = false;
		}
	});

	function getErrorSummary(error: string): string {
		if (!error) return '';
		
		// Categorize errors
		if (error.toLowerCase().includes('rate limit')) {
			return 'Rate limit reached';
		}
		if (error.toLowerCase().includes('queue') || error.toLowerCase().includes('enqueue')) {
			return 'Queue error';
		}
		if (error.toLowerCase().includes('network') || error.toLowerCase().includes('timeout')) {
			return 'Network error';
		}
		if (error.toLowerCase().includes('authentication') || error.toLowerCase().includes('login')) {
			return 'Authentication error';
		}
		if (error.toLowerCase().includes('website') || error.toLowerCase().includes('page')) {
			return 'Website error';
		}
		
		// Return first 50 characters as summary
		return error.length > 50 ? error.substring(0, 50) + '...' : error;
	}
</script>

<tr class="table-row" class:selected data-status={application.status}>
	<td>
		<input
			type="checkbox"
			checked={selected}
			onchange={() => onToggleSelection?.(application.id)}
			class="row-checkbox"
		/>
	</td>
	<td><strong>{application.companyName}</strong></td>
	<td>{application.jobTitle}</td>
	<td>{application.location || '—'}</td>
	<td>{application.website}</td>
	<td>{application.language || '—'}</td>
	<td>
		<div class="status-cell">
			<div class="status-header">
				<StatusBadge status={application.status} type="status">
					{statusInfo.text}
				</StatusBadge>
				{#if application.status === 'processing'}
					<div class="processing-indicator" title="Application is being processed">
						<div class="spinner"></div>
					</div>
				{/if}
			</div>
			{#if statusInfo.timestamp}
				<span class="status-timestamp" title={statusInfo.timestamp}>
					{statusInfo.timestamp}
				</span>
			{/if}
			{#if statusInfo.error}
				<div class="error-details" title={statusInfo.error}>
					<span class="error-icon">⚠️</span>
					<span class="error-text">{getErrorSummary(statusInfo.error)}</span>
					<button class="error-toggle" onclick={() => showFullError = !showFullError}>
						{showFullError ? 'Hide' : 'Show'} details
					</button>
				</div>
				{#if showFullError}
					<div class="error-full">
						{statusInfo.error}
					</div>
				{/if}
			{/if}
		</div>
	</td>
	<td>{application.interestLevel || '—'}</td>
	<td>
		{#if application.tags && application.tags.length > 0}
			<div class="tags-display">
				{#each application.tags as tag}
					<span class="tag-badge">{tag}</span>
				{/each}
			</div>
		{:else}
			—
		{/if}
	</td>
	<td>
		{#if application.followUpDate}
			<span class="follow-up-date" title={formatDateTime(application.followUpDate)}>
				{formatDate(application.followUpDate)}
			</span>
		{:else}
			—
		{/if}
	</td>
	<td>{formatDate(application.appliedAt)}</td>
	<td class="actions-cell">
		<div class="actions">
			<button 
				class="action-btn preview-btn" 
				title="Preview"
				onclick={() => onView?.(application)}
			>
				<Eye size={16} />
			</button>
			<a 
				href="/applications/{application.id}" 
				class="action-btn external-link-btn" 
				title="Open in new page"
				onclick={(e) => {
					e.preventDefault();
					goto(`/applications/${application.id}`);
				}}
			>
				<ExternalLink size={16} />
			</a>
			<Button variant="secondary" size="sm" onclick={() => onView?.(application)}>
				View
			</Button>
			<Button variant="primary" size="sm" onclick={() => onEdit?.(application)}>
				Edit
			</Button>
			<Button variant="danger" size="sm" onclick={() => onDelete?.(application.id)}>
				Delete
			</Button>
		</div>
	</td>
</tr>

<style>
	.table-row {
		transition: background-color var(--transition-base);
	}
	
	.table-row:hover {
		background-color: var(--color-bg-hover);
	}

	.table-row.selected {
		background-color: var(--color-bg-secondary);
	}

	.table-row td:first-child {
		width: 40px;
	}

	/* Status-based row styling */
	.table-row[data-status="accepted"] {
		border-left: 3px solid var(--color-success, #10b981);
	}

	.table-row[data-status="rejected"] {
		border-left: 3px solid var(--color-danger, #ef4444);
	}

	.table-row[data-status="contacted"] {
		border-left: 3px solid var(--color-primary, #3b82f6);
	}

	.table-row[data-status="failed"] {
		opacity: 0.7;
	}

	.row-checkbox {
		cursor: pointer;
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

	.actions {
		display: flex;
		gap: var(--spacing-xs);
		align-items: center;
		flex-wrap: wrap;
	}

	.action-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: var(--spacing-xs);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		cursor: pointer;
		transition: all var(--transition-base);
		text-decoration: none;
	}

	.action-btn:hover {
		background-color: var(--color-bg-hover);
		border-color: var(--color-primary);
		color: var(--color-primary);
	}

	.preview-btn:hover {
		background-color: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}

	.external-link-btn:hover {
		background-color: var(--color-success);
		color: white;
		border-color: var(--color-success);
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

	.status-header {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
	}

	.processing-indicator {
		display: inline-flex;
		align-items: center;
	}

	.spinner {
		width: 12px;
		height: 12px;
		border: 2px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.error-details {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		margin-top: var(--spacing-xs);
		padding: var(--spacing-xs);
		background-color: var(--color-danger, #ef4444);
		background-color: rgba(239, 68, 68, 0.1);
		border: 1px solid var(--color-danger, #ef4444);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
	}

	.error-icon {
		font-size: 0.9rem;
	}

	.error-text {
		flex: 1;
		color: var(--color-danger, #ef4444);
		font-weight: var(--font-weight-medium);
	}

	.error-toggle {
		background: none;
		border: none;
		color: var(--color-danger, #ef4444);
		cursor: pointer;
		font-size: var(--font-size-xs);
		text-decoration: underline;
		padding: 0;
	}

	.error-toggle:hover {
		opacity: 0.8;
	}

	.error-full {
		margin-top: var(--spacing-xs);
		padding: var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		color: var(--color-text-primary);
		white-space: pre-wrap;
		word-break: break-word;
	}

	.tags-display {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
	}

	.tag-badge {
		display: inline-block;
		padding: 2px var(--spacing-xs);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		color: var(--color-text-primary);
	}

	.follow-up-date {
		color: var(--color-primary);
		font-weight: var(--font-weight-medium);
	}
</style>
