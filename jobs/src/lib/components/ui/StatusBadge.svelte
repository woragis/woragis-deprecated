<script lang="ts">
	let {
		status,
		type = 'status' as 'status' | 'response' | 'outcome',
		label
	}: {
		status: string;
		type?: 'status' | 'response' | 'outcome';
		label?: string;
	} = $props();
	
	let className = $derived.by(() => {
		return `badge badge-${type} badge-${status}`;
	});
</script>

<span class={className}>
	{label || status}
</span>

<style>
	.badge {
		display: inline-block;
		padding: var(--spacing-xs) var(--spacing-sm);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		text-transform: capitalize;
	}
	
	/* Status Badges */
	.badge-status-pending {
		background-color: var(--status-pending-bg);
		color: var(--status-pending-text);
	}
	
	.badge-status-processing {
		background-color: var(--status-processing-bg);
		color: var(--status-processing-text);
	}
	
	.badge-status-applied {
		background-color: var(--status-applied-bg);
		color: var(--status-applied-text);
	}
	
	.badge-status-contacted {
		background-color: var(--status-contacted-bg);
		color: var(--status-contacted-text);
	}
	
	.badge-status-rejected {
		background-color: var(--status-rejected-bg);
		color: var(--status-rejected-text);
	}
	
	.badge-status-accepted {
		background-color: var(--status-accepted-bg);
		color: var(--status-accepted-text);
	}
	
	.badge-status-failed {
		background-color: var(--status-failed-bg);
		color: var(--status-failed-text);
	}
	
	/* Response Type Badges */
	.badge-response-rejection {
		background-color: var(--status-rejected-bg);
		color: var(--status-rejected-text);
	}
	
	.badge-response-interview {
		background-color: var(--status-processing-bg);
		color: var(--status-processing-text);
	}
	
	.badge-response-offer {
		background-color: var(--status-accepted-bg);
		color: var(--status-accepted-text);
	}
	
	.badge-response-no-response {
		background-color: var(--color-bg-tertiary);
		color: var(--color-text-secondary);
	}
	
	/* Outcome Badges */
	.badge-outcome-pending {
		background-color: var(--status-pending-bg);
		color: var(--status-pending-text);
	}
	
	.badge-outcome-passed {
		background-color: var(--status-accepted-bg);
		color: var(--status-accepted-text);
	}
	
	.badge-outcome-failed {
		background-color: var(--status-rejected-bg);
		color: var(--status-rejected-text);
	}
	
	.badge-outcome-cancelled {
		background-color: var(--color-bg-tertiary);
		color: var(--color-text-secondary);
	}

	/* Resume Status Badges */
	.badge-status-main {
		background-color: #fef3c7;
		color: #92400e;
	}

	.badge-status-featured {
		background-color: #dbeafe;
		color: #1e40af;
	}
</style>
