<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { JobApplication } from '$lib/api/jobapplications';
	import { goto } from '$app/navigation';
	
	let {
		open = $bindable(false),
		application,
		onEdit,
		onDelete
	}: {
		open?: boolean;
		application?: JobApplication | null;
		onEdit?: () => void;
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
	
	function handleViewFull() {
		if (application) {
			goto(`/admin/job-applications-admin/${application.id}`);
		}
	}
	
	function handleDelete() {
		if (application && onDelete) {
			if (confirm('Are you sure you want to delete this application?')) {
				onDelete(application.id);
				open = false;
			}
		}
	}
</script>

{#if application}
	<Modal bind:open size="lg" title="Application Details">
		<div class="modal-content">
			<div class="details-grid">
				<div class="detail-section">
					<h3>Company Information</h3>
					<div class="detail-item">
						<span class="label">Company:</span>
						<span class="value">{application.companyName}</span>
					</div>
					<div class="detail-item">
						<span class="label">Job Title:</span>
						<span class="value">{application.jobTitle}</span>
					</div>
					<div class="detail-item">
						<span class="label">Location:</span>
						<span class="value">{application.location || '—'}</span>
					</div>
					<div class="detail-item">
						<span class="label">Website:</span>
						<span class="value">{application.website}</span>
					</div>
					<div class="detail-item">
						<span class="label">Job URL:</span>
						<a href={application.jobUrl} target="_blank" rel="noopener noreferrer" class="link">
							View Job Posting
						</a>
					</div>
				</div>
				
				<div class="detail-section">
					<h3>Status & Interest</h3>
					<div class="detail-item">
						<span class="label">Status:</span>
						<div class="status-detail">
							<StatusBadge status={application.status} type="status">
								{tFn(`jobApplications.status.${application.status}` as any)}
							</StatusBadge>
							{#if application.status === 'processing'}
								<div class="processing-indicator" title="Application is being processed">
									<div class="spinner"></div>
									<span>Processing...</span>
								</div>
							{/if}
						</div>
					</div>
					{#if application.errorMessage}
						<div class="detail-item error-item">
							<span class="label">Error:</span>
							<div class="error-message">
								<span class="error-icon">⚠️</span>
								<span class="error-text">{application.errorMessage}</span>
							</div>
						</div>
					{/if}
					<div class="detail-item">
						<span class="label">Interest Level:</span>
						<span class="value">{application.interestLevel || '—'}</span>
					</div>
					<div class="detail-item">
						<span class="label">Language:</span>
						<span class="value">{application.language || '—'}</span>
					</div>
				</div>
				
				<div class="detail-section">
					<h3>Dates</h3>
					<div class="detail-item">
						<span class="label">Created:</span>
						<span class="value">{formatDateTime(application.createdAt)}</span>
					</div>
					<div class="detail-item">
						<span class="label">Applied:</span>
						<span class="value">{formatDate(application.appliedAt)}</span>
					</div>
					<div class="detail-item">
						<span class="label">Follow-up:</span>
						<span class="value">{formatDate(application.followUpDate)}</span>
					</div>
					<div class="detail-item">
						<span class="label">Response Received:</span>
						<span class="value">{formatDate(application.responseReceivedAt)}</span>
					</div>
				</div>
				
				{#if application.tags && application.tags.length > 0}
					<div class="detail-section">
						<h3>Tags</h3>
						<div class="tags-display">
							{#each application.tags as tag}
								<span class="tag-badge">{tag}</span>
							{/each}
						</div>
					</div>
				{/if}
				
				{#if application.notes}
					<div class="detail-section full-width">
						<h3>Notes</h3>
						<div class="notes-content">{application.notes}</div>
					</div>
				{/if}
			</div>
			
			<div class="modal-actions">
				<Button onclick={onEdit} variant="primary">Edit</Button>
				<Button onclick={handleViewFull} variant="secondary">View Full Details</Button>
				<Button onclick={handleDelete} variant="danger">Delete</Button>
			</div>
		</div>
	</Modal>
{/if}

<style>
	.modal-content {
		padding: var(--spacing-lg);
	}
	
	.details-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: var(--spacing-lg);
		margin-bottom: var(--spacing-lg);
	}
	
	.detail-section {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}
	
	.detail-section.full-width {
		grid-column: 1 / -1;
	}
	
	.detail-section h3 {
		margin: 0 0 var(--spacing-sm) 0;
		font-size: var(--font-size-md);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
		border-bottom: 2px solid var(--color-border);
		padding-bottom: var(--spacing-xs);
	}
	
	.detail-item {
		display: flex;
		align-items: flex-start;
		gap: var(--spacing-sm);
	}
	
	.detail-item .label {
		font-weight: var(--font-weight-medium);
		color: var(--color-text-secondary);
		min-width: 120px;
		font-size: var(--font-size-sm);
	}
	
	.detail-item .value {
		flex: 1;
		color: var(--color-text-primary);
		font-size: var(--font-size-sm);
	}
	
	.link {
		color: var(--color-primary);
		text-decoration: none;
	}
	
	.link:hover {
		text-decoration: underline;
	}
	
	.tags-display {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
	}
	
	.tag-badge {
		display: inline-block;
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-xs);
		color: var(--color-text-primary);
	}
	
	.notes-content {
		padding: var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		white-space: pre-wrap;
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		line-height: 1.6;
	}

	.status-detail {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
	}

	.processing-indicator {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
	}

	.spinner {
		width: 14px;
		height: 14px;
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

	.error-item {
		align-items: flex-start;
	}

	.error-message {
		flex: 1;
		display: flex;
		align-items: flex-start;
		gap: var(--spacing-xs);
		padding: var(--spacing-sm);
		background-color: rgba(239, 68, 68, 0.1);
		border: 1px solid var(--color-danger, #ef4444);
		border-radius: var(--radius-md);
	}

	.error-icon {
		font-size: 1.2rem;
		flex-shrink: 0;
	}

	.error-text {
		flex: 1;
		color: var(--color-danger, #ef4444);
		font-size: var(--font-size-sm);
		white-space: pre-wrap;
		word-break: break-word;
	}
	
	.modal-actions {
		display: flex;
		gap: var(--spacing-sm);
		justify-content: flex-end;
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
	}
</style>

