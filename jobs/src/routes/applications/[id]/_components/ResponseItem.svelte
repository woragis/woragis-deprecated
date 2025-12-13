<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { JobApplicationResponse } from '$lib/api/jobapplicationresponses';
	
	let {
		response,
		onEdit,
		onDelete
	}: {
		response: JobApplicationResponse;
		onEdit?: (response: JobApplicationResponse) => void;
		onDelete?: (id: string) => void;
	} = $props();
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="list-item">
	<div class="list-item-header">
		<StatusBadge status={response.responseType} type="response">
			{response.responseType}
		</StatusBadge>
		<span class="date">{formatDate(response.responseDate)}</span>
	</div>
	{#if response.message}
		<p class="list-item-text">{response.message}</p>
	{/if}
	{#if response.contactPerson}
		<p class="list-item-text">
			<strong>Contact:</strong> {response.contactPerson} {response.contactEmail ? `(${response.contactEmail})` : ''}
		</p>
	{/if}
	<div class="list-item-actions">
		{#if onEdit}
			<Button variant="secondary" size="sm" onclick={() => onEdit(response)}>Edit</Button>
		{/if}
		{#if onDelete}
			<Button variant="danger" size="sm" onclick={() => onDelete(response.id)}>Delete</Button>
		{/if}
	</div>
</div>

<style>
	.list-item {
		background-color: var(--color-bg-primary);
		padding: var(--spacing-md);
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
	}
	
	.list-item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-sm);
	}
	
	.date {
		font-size: var(--font-size-xs);
		color: var(--color-text-tertiary);
	}
	
	.list-item-text {
		margin: var(--spacing-xs) 0;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		line-height: 1.5;
	}
	
	.list-item-text strong {
		color: var(--color-text-primary);
	}
	
	.list-item-actions {
		display: flex;
		gap: var(--spacing-sm);
		margin-top: var(--spacing-sm);
	}
</style>
