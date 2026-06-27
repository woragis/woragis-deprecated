<script lang="ts">
	import type { SocialMediaPost } from '$lib/api/socialmediaposts';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { goto } from '$app/navigation';

	let {
		posts = [],
		onDelete,
		onView
	}: {
		posts?: SocialMediaPost[];
		onDelete?: (id: string) => void;
		onView?: (post: SocialMediaPost) => void;
	} = $props();

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function getStatusColor(status: string): string {
		const colors: Record<string, string> = {
			draft: 'var(--color-text-secondary)',
			ready: 'var(--color-primary)',
			scheduled: 'var(--color-warning)',
			posted: 'var(--color-success)',
			analyzed: 'var(--color-primary)',
			archived: 'var(--color-text-tertiary)'
		};
		return colors[status] || 'var(--color-text-secondary)';
	}
</script>

{#if posts.length === 0}
	<EmptyState
		title="No posts found"
		description="Create your first social media post to get started"
	/>
{:else}
	<div class="table-wrapper">
		<table class="table">
			<thead>
				<tr>
					<th>Platform</th>
					<th>Title</th>
					<th>Format</th>
					<th>Status</th>
					<th>Scheduled</th>
					<th>Posted</th>
					<th>Engagement</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each posts as post}
					<tr class="table-row" onclick={() => onView?.(post) || goto(`/posts/${post.id}`)}>
						<td>
							<span class="platform-badge">{post.platform}</span>
						</td>
						<td>
							<div class="post-title">{post.title}</div>
						</td>
						<td>
							<span class="format-badge">{post.format}</span>
						</td>
						<td>
							<span class="status-badge" style="color: {getStatusColor(post.status)}">
								{post.status}
							</span>
						</td>
						<td>
							{#if post.scheduledAt}
								{formatDate(post.scheduledAt)}
							{:else}
								<span class="text-muted">—</span>
							{/if}
						</td>
						<td>
							{#if post.postedAt}
								{formatDate(post.postedAt)}
							{:else}
								<span class="text-muted">—</span>
							{/if}
						</td>
						<td>
							<div class="engagement">
								{#if post.likes}
									<span>👍 {post.likes}</span>
								{/if}
								{#if post.views}
									<span>👁️ {post.views}</span>
								{/if}
							</div>
						</td>
						<td>{formatDate(post.createdAt)}</td>
						<td>
							<div class="actions" onclick={(e) => e.stopPropagation()}>
								<button
									class="action-button"
									onclick={() => onView?.(post) || goto(`/posts/${post.id}`)}
								>
									View
								</button>
								{#if onDelete}
									<button class="action-button danger" onclick={() => onDelete(post.id)}>
										Delete
									</button>
								{/if}
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

	.platform-badge,
	.format-badge {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		text-transform: capitalize;
	}

	.status-badge {
		font-weight: var(--font-weight-medium);
		text-transform: capitalize;
	}

	.post-title {
		font-weight: var(--font-weight-medium);
		max-width: 200px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.engagement {
		display: flex;
		gap: var(--spacing-sm);
		font-size: var(--font-size-xs);
	}

	.text-muted {
		color: var(--color-text-tertiary);
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

	.action-button.danger:hover {
		background-color: var(--color-danger);
		border-color: var(--color-danger);
	}
</style>
