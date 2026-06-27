<script lang="ts">
	import type { ScheduledPost } from '$lib/api/scheduling';
	import { goto } from '$app/navigation';

	let { upcomingPosts }: { upcomingPosts: ScheduledPost[] } = $props();

	function formatDateTime(dateTime: string): string {
		const date = new Date(dateTime);
		return date.toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		});
	}
</script>

{#if upcomingPosts.length > 0}
	<div class="upcoming-posts">
		<h3 class="section-title">Upcoming Scheduled Posts</h3>
		<div class="posts-list">
			{#each upcomingPosts as post}
				<div class="post-item" onclick={() => goto(`/scheduling/${post.id}`)}>
					<div class="post-time">{formatDateTime(post.scheduledAt)}</div>
					<div class="post-status">{post.status}</div>
				</div>
			{/each}
		</div>
		<button class="view-all-button" onclick={() => goto('/scheduling')}>
			View All Scheduled Posts →
		</button>
	</div>
{:else}
	<div class="empty-state">
		<p>No upcoming scheduled posts</p>
		<button class="view-all-button" onclick={() => goto('/scheduling')}>
			Schedule a Post →
		</button>
	</div>
{/if}

<style>
	.upcoming-posts {
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

	.posts-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
		margin-bottom: var(--spacing-md);
	}

	.post-item {
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

	.post-item:hover {
		background-color: var(--color-bg-hover);
	}

	.post-time {
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		font-weight: var(--font-weight-medium);
	}

	.post-status {
		font-size: var(--font-size-xs);
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-primary);
		color: white;
		border-radius: var(--radius-sm);
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
