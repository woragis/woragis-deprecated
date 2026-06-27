<script lang="ts">
	import type { ContentPost } from '$lib/api/content';
	import { goto } from '$app/navigation';

	let { contentBacklog }: { contentBacklog: ContentPost[] } = $props();

	function getPriorityColor(priority: string): string {
		switch (priority) {
			case 'high':
				return 'var(--color-danger)';
			case 'medium':
				return 'var(--color-warning)';
			case 'low':
				return 'var(--color-success)';
			default:
				return 'var(--color-text-secondary)';
		}
	}
</script>

{#if contentBacklog.length > 0}
	<div class="content-backlog">
		<h3 class="section-title">Content Backlog</h3>
		<div class="backlog-list">
			{#each contentBacklog.slice(0, 5) as item}
				<div class="backlog-item" onclick={() => goto(`/content/posts/${item.id}`)}>
					<div class="backlog-info">
						<div class="backlog-priority" style="background-color: {getPriorityColor(item.priority)}">
							{item.priority}
						</div>
						<div class="backlog-details">
							<div class="backlog-status">{item.status}</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
		<button class="view-all-button" onclick={() => goto('/content/backlog')}>
			View All Backlog ({contentBacklog.length}) →
		</button>
	</div>
{:else}
	<div class="empty-state">
		<p>No content in backlog</p>
		<button class="view-all-button" onclick={() => goto('/content/posts')}>
			Create Content Post →
		</button>
	</div>
{/if}

<style>
	.content-backlog {
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		padding: var(--spacing-md);
	}

	.section-title {
		margin: 0 0 var(--spacing-md) 0;
		font-size: var(--font-size-lg);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.backlog-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
		margin-bottom: var(--spacing-md);
	}

	.backlog-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: background-color var(--transition-base);
	}

	.backlog-item:hover {
		background-color: var(--color-bg-hover);
	}

	.backlog-info {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
		flex: 1;
	}

	.backlog-priority {
		padding: var(--spacing-xs) var(--spacing-sm);
		border-radius: var(--radius-sm);
		color: white;
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		text-transform: capitalize;
	}

	.backlog-details {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.backlog-status {
		font-size: var(--font-size-xs);
		color: var(--color-text-secondary);
		text-transform: capitalize;
	}

	.view-all-button {
		width: 100%;
		padding: var(--spacing-sm);
		background-color: transparent;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		cursor: pointer;
		font-size: var(--font-size-sm);
		transition: all var(--transition-base);
	}

	.view-all-button:hover {
		background-color: var(--color-bg-hover);
		color: var(--color-text-primary);
	}

	.empty-state {
		text-align: center;
		padding: var(--spacing-xl);
		color: var(--color-text-secondary);
	}
</style>
