<script lang="ts">
	import type { ContentPost } from '$lib/api/content';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { goto } from '$app/navigation';

	let {
		contentPosts = [],
		onView
	}: {
		contentPosts?: ContentPost[];
		onView?: (post: ContentPost) => void;
	} = $props();

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

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

{#if contentPosts.length === 0}
	<EmptyState
		title="No content posts found"
		description="Create your first content post to get started"
	/>
{:else}
	<div class="table-wrapper">
		<table class="table">
			<thead>
				<tr>
					<th>Priority</th>
					<th>Status</th>
					<th>Project</th>
					<th>Content Type</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each contentPosts as post}
					<tr class="table-row" onclick={() => onView?.(post) || goto(`/content/posts/${post.id}`)}>
						<td>
							<span class="priority-badge" style="background-color: {getPriorityColor(post.priority)}">
								{post.priority}
							</span>
						</td>
						<td>
							<span class="status-badge">{post.status}</span>
						</td>
						<td>{post.project || '—'}</td>
						<td>{post.contentType || '—'}</td>
						<td>{formatDate(post.createdAt)}</td>
						<td>
							<div class="actions" onclick={(e) => e.stopPropagation()}>
								<button
									class="action-button"
									onclick={() => onView?.(post) || goto(`/content/posts/${post.id}`)}
								>
									View
								</button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<style>
	.table-wrapper {
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
		background-color: var(--color-bg-primary);
	}

	.table th {
		padding: var(--spacing-md);
		text-align: left;
		background-color: var(--color-bg-secondary);
		font-weight: var(--font-weight-semibold);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		border-bottom: 2px solid var(--color-border);
	}

	.table-row {
		cursor: pointer;
		transition: background-color var(--transition-base);
	}

	.table-row:hover {
		background-color: var(--color-bg-hover);
	}

	.table td {
		padding: var(--spacing-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		border-bottom: 1px solid var(--color-border);
	}

	.priority-badge {
		padding: var(--spacing-xs) var(--spacing-sm);
		border-radius: var(--radius-sm);
		color: white;
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		text-transform: capitalize;
	}

	.status-badge {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		text-transform: capitalize;
	}

	.actions {
		display: flex;
		gap: var(--spacing-xs);
	}

	.action-button {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		cursor: pointer;
		transition: all var(--transition-base);
	}

	.action-button:hover {
		background-color: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}
</style>
