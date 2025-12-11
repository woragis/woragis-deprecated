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
						<StatusBadge status={application.status} type="status">
							{tFn(`jobApplications.status.${application.status}` as any)}
						</StatusBadge>
					</div>
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
	
	.modal-actions {
		display: flex;
		gap: var(--spacing-sm);
		justify-content: flex-end;
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
	}
</style>

