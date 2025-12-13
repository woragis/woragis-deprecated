<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { InterviewStage } from '$lib/api/jobapplicationstages';
	
	let {
		stage,
		onEdit,
		onComplete,
		onDelete
	}: {
		stage: InterviewStage;
		onEdit?: (stage: InterviewStage) => void;
		onComplete?: (stage: InterviewStage) => void;
		onDelete?: (id: string) => void;
	} = $props();
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="list-item">
	<div class="list-item-header">
		<span class="stage-type">{stage.stageType}</span>
		<StatusBadge status={stage.outcome} type="outcome">
			{stage.outcome}
		</StatusBadge>
	</div>
	{#if stage.scheduledDate}
		<p class="list-item-text"><strong>Scheduled:</strong> {formatDate(stage.scheduledDate)}</p>
	{/if}
	{#if stage.interviewerName}
		<p class="list-item-text">
			<strong>Interviewer:</strong> {stage.interviewerName} {stage.interviewerEmail ? `(${stage.interviewerEmail})` : ''}
		</p>
	{/if}
	{#if stage.notes}
		<p class="list-item-text">{stage.notes}</p>
	{/if}
	{#if stage.feedback}
		<p class="list-item-text"><strong>Feedback:</strong> {stage.feedback}</p>
	{/if}
	<div class="list-item-actions">
		{#if onEdit}
			<Button variant="secondary" size="sm" onclick={() => onEdit(stage)}>Edit</Button>
		{/if}
		{#if stage.outcome === 'pending' && onComplete}
			<Button variant="success" size="sm" onclick={() => onComplete(stage)}>Complete</Button>
		{/if}
		{#if onDelete}
			<Button variant="danger" size="sm" onclick={() => onDelete(stage.id)}>Delete</Button>
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
	
	.stage-type {
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
		color: var(--color-text-primary);
		text-transform: capitalize;
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
