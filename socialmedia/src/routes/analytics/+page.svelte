<script lang="ts">
	import { onMount } from 'svelte';
	import { getAnalyticsSummary, getTopPosts, type AnalyticsSummary, type TopPost } from '$lib/api/analytics';
	import PageHeader from '../posts/_sections/PageHeader.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Card from '$lib/components/ui/Card.svelte';

	let analyticsSummary: AnalyticsSummary | null = $state(null);
	let topPosts: TopPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchAnalytics();
	});

	async function fetchAnalytics() {
		loading = true;
		error = null;
		try {
			const [summary, top] = await Promise.all([
				getAnalyticsSummary(),
				getTopPosts({ limit: 10, metric: 'likes' })
			]);
			analyticsSummary = summary;
			topPosts = top;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load analytics';
			console.error('Error fetching analytics:', err);
		} finally {
			loading = false;
		}
	}

	const engagementRate = $derived(() => {
		if (!analyticsSummary || analyticsSummary.totalImpressions === 0) return 0;
		const totalEngagement =
			analyticsSummary.totalLikes +
			analyticsSummary.totalComments +
			analyticsSummary.totalShares +
			(analyticsSummary.totalClicks || 0);
		return Math.round((totalEngagement / analyticsSummary.totalImpressions) * 100 * 10) / 10;
	});
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="page-header">
		<div class="header-content">
			<h1 class="page-title">Analytics</h1>
			<p class="page-subtitle">Track and analyze your social media post performance</p>
		</div>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchAnalytics} />
	{:else if loading}
		<LoadingState message="Loading analytics..." />
	{:else if analyticsSummary}
		<div class="analytics-grid">
			<Card title="Summary Metrics">
				<div class="metrics-grid">
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalLikes.toLocaleString()}</div>
						<div class="metric-label">Total Likes</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalComments.toLocaleString()}</div>
						<div class="metric-label">Total Comments</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalShares.toLocaleString()}</div>
						<div class="metric-label">Total Shares</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalViews.toLocaleString()}</div>
						<div class="metric-label">Total Views</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{engagementRate()}%</div>
						<div class="metric-label">Engagement Rate</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalReach.toLocaleString()}</div>
						<div class="metric-label">Total Reach</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.totalImpressions.toLocaleString()}</div>
						<div class="metric-label">Total Impressions</div>
					</div>
					<div class="metric-item">
						<div class="metric-value">{analyticsSummary.postCount}</div>
						<div class="metric-label">Posts Tracked</div>
					</div>
				</div>
			</Card>

			<Card title="Top Performing Posts" subtitle="By likes">
				{#if topPosts.length > 0}
					<div class="top-posts-list">
						{#each topPosts as post}
							<div class="top-post-item">
								<span class="post-id">{post.socialPostId.slice(0, 8)}...</span>
								<span class="post-metric">{post.metricValue.toLocaleString()} {post.metricName}</span>
							</div>
						{/each}
					</div>
				{:else}
					<p class="empty-text">No top posts data available</p>
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

	.analytics-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--spacing-lg);
	}

	@media (min-width: 1024px) {
		.analytics-grid {
			grid-template-columns: 2fr 1fr;
		}
	}

	.metrics-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: var(--spacing-md);
	}

	.metric-item {
		text-align: center;
		padding: var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-md);
	}

	.metric-value {
		font-size: 1.5rem;
		font-weight: var(--font-weight-bold);
		color: var(--color-text-primary);
		margin-bottom: var(--spacing-xs);
	}

	.metric-label {
		font-size: var(--font-size-xs);
		color: var(--color-text-secondary);
	}

	.top-posts-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}

	.top-post-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
	}

	.post-id {
		color: var(--color-text-secondary);
		font-family: monospace;
	}

	.post-metric {
		font-weight: var(--font-weight-semibold);
		color: var(--color-primary);
	}

	.empty-text {
		color: var(--color-text-secondary);
		text-align: center;
		padding: var(--spacing-md);
	}
</style>
