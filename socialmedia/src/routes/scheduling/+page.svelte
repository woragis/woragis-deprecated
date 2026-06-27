<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { listScheduledPosts, getUpcomingScheduledPosts, type ScheduledPost } from '$lib/api/scheduling';
	import PageHeader from '../posts/_sections/PageHeader.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Card from '$lib/components/ui/Card.svelte';

	let scheduledPosts: ScheduledPost[] = $state([]);
	let upcomingPosts: ScheduledPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchSchedules();
	});

	async function fetchSchedules() {
		loading = true;
		error = null;
		try {
			const today = new Date();
			const nextWeek = new Date(today);
			nextWeek.setDate(today.getDate() + 7);

			const [allSchedules, upcoming] = await Promise.all([
				listScheduledPosts({
					startDate: today.toISOString().split('T')[0],
					endDate: nextWeek.toISOString().split('T')[0]
				}),
				getUpcomingScheduledPosts(10)
			]);

			scheduledPosts = allSchedules;
			upcomingPosts = upcoming;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load scheduled posts';
			console.error('Error fetching scheduled posts:', err);
		} finally {
			loading = false;
		}
	}

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

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="page-header">
		<div class="header-content">
			<h1 class="page-title">Scheduling</h1>
			<p class="page-subtitle">Manage and view your scheduled social media posts</p>
		</div>
		<button class="create-button" onclick={() => goto('/scheduling/new')}>
			Schedule Post
		</button>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchSchedules} />
	{:else if loading}
		<LoadingState message="Loading scheduled posts..." />
	{:else}
		<div class="scheduling-grid">
			<Card title="Upcoming Posts" subtitle="Next 10 scheduled posts">
				{#if upcomingPosts.length > 0}
					<div class="posts-list">
						{#each upcomingPosts as post}
							<div class="post-item" onclick={() => goto(`/scheduling/${post.id}`)}>
								<div class="post-time">{formatDateTime(post.scheduledAt)}</div>
								<div class="post-status">{post.status}</div>
							</div>
						{/each}
					</div>
				{:else}
					<p class="empty-text">No upcoming scheduled posts</p>
				{/if}
			</Card>

			<Card title="This Week's Schedule" subtitle="All scheduled posts for the next 7 days">
				{#if scheduledPosts.length > 0}
					<div class="posts-list">
						{#each scheduledPosts as post}
							<div class="post-item" onclick={() => goto(`/scheduling/${post.id}`)}>
								<div class="post-time">{formatDateTime(post.scheduledAt)}</div>
								<div class="post-status">{post.status}</div>
							</div>
						{/each}
					</div>
				{:else}
					<p class="empty-text">No scheduled posts for this week</p>
				{/if}
			</Card>
		</div>
	{/if}
</div>

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}

	.header-content {
		flex: 1;
	}

	.page-title {
		margin: 0 0 var(--spacing-xs) 0;
		font-size: var(--font-size-2xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.page-subtitle {
		margin: 0;
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
	}

	.create-button {
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-primary);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
		cursor: pointer;
		transition: background-color var(--transition-base);
	}

	.create-button:hover {
		background-color: var(--color-primary-hover);
	}

	.scheduling-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--spacing-lg);
	}

	@media (min-width: 1024px) {
		.scheduling-grid {
			grid-template-columns: 1fr 1fr;
		}
	}

	.posts-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}

	.post-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-bg-secondary);
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

	.empty-text {
		color: var(--color-text-secondary);
		text-align: center;
		padding: var(--spacing-md);
	}
</style>
