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
</script>

<tr class="table-row">
	<td><strong>{application.companyName}</strong></td>
	<td>{application.jobTitle}</td>
	<td>{application.location || '—'}</td>
	<td>{application.website}</td>
	<td>{application.language || '—'}</td>
	<td>
		<StatusBadge status={application.status} type="status">
			{tFn(`jobApplications.status.${application.status}` as any)}
		</StatusBadge>
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
</style>
