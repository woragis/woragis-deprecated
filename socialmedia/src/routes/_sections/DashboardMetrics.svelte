<script lang="ts">
	import type { DashboardData } from '$lib/api/socialmediaposts';

	let { dashboardData }: { dashboardData: DashboardData } = $props();

	const metrics = $derived(() => {
		const { analyticsSummary, upcomingPosts, contentBacklog, recentPosts } = dashboardData;

		// Status breakdown
		const postsByStatus = recentPosts.reduce(
			(acc, post) => {
				acc[post.status] = (acc[post.status] || 0) + 1;
				return acc;
			},
			{} as Record<string, number>
		);

		// Platform breakdown
		const postsByPlatform = recentPosts.reduce(
			(acc, post) => {
				acc[post.platform] = (acc[post.platform] || 0) + 1;
				return acc;
			},
			{} as Record<string, number>
		);

		// Engagement rate calculation
		const engagementRate =
			analyticsSummary.totalImpressions > 0
				? ((analyticsSummary.totalLikes +
						analyticsSummary.totalComments +
						analyticsSummary.totalShares +
						(analyticsSummary.totalClicks || 0)) /
						analyticsSummary.totalImpressions) *
				  100
				: 0;

		return {
			totalPosts: recentPosts.length,
			upcomingScheduled: upcomingPosts.length,
			contentBacklogCount: contentBacklog.length,
			totalEngagement: analyticsSummary.totalLikes + analyticsSummary.totalComments + analyticsSummary.totalShares,
			engagementRate: Math.round(engagementRate * 10) / 10,
			totalViews: analyticsSummary.totalViews,
			totalReach: analyticsSummary.totalReach,
			postsByStatus,
			postsByPlatform
		};
	});
</script>

<div class="metrics-container">
	<div class="metrics-grid">
		<div class="metric-card">
			<div class="metric-value">{metrics().totalPosts}</div>
			<div class="metric-label">Recent Posts</div>
		</div>
		<div class="metric-card highlight">
			<div class="metric-value">{metrics().upcomingScheduled}</div>
			<div class="metric-label">Upcoming Scheduled</div>
		</div>
		<div class="metric-card highlight">
			<div class="metric-value">{metrics().contentBacklogCount}</div>
			<div class="metric-label">Content Backlog</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().totalEngagement.toLocaleString()}</div>
			<div class="metric-label">Total Engagement</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().engagementRate}%</div>
			<div class="metric-label">Engagement Rate</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().totalViews.toLocaleString()}</div>
			<div class="metric-label">Total Views</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().totalReach.toLocaleString()}</div>
			<div class="metric-label">Total Reach</div>
		</div>
	</div>

	<!-- Status Breakdown -->
	{#if Object.keys(metrics().postsByStatus).length > 0}
		<div class="breakdown-section">
			<h4>Posts by Status</h4>
			<div class="breakdown-grid">
				{#each Object.entries(metrics().postsByStatus) as [status, count]}
					<div class="breakdown-item">
						<span class="breakdown-label">{status}</span>
						<span class="breakdown-value">{count}</span>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Platform Breakdown -->
	{#if Object.keys(metrics().postsByPlatform).length > 0}
		<div class="breakdown-section">
			<h4>Posts by Platform</h4>
			<div class="breakdown-grid">
				{#each Object.entries(metrics().postsByPlatform) as [platform, count]}
					<div class="breakdown-item">
						<span class="breakdown-label">{platform}</span>
						<span class="breakdown-value">{count}</span>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	.metrics-container {
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		padding: var(--spacing-md);
		margin-bottom: var(--spacing-md);
	}

	.metrics-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: var(--spacing-md);
		margin-bottom: var(--spacing-md);
	}

	.metric-card {
		text-align: center;
		padding: var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
	}

	.metric-card.highlight {
		background-color: var(--color-warning, #f59e0b);
		color: white;
		border-color: var(--color-warning, #f59e0b);
	}

	.metric-value {
		font-size: 2rem;
		font-weight: var(--font-weight-bold);
		color: var(--color-text-primary);
		margin-bottom: var(--spacing-xs);
	}

	.metric-card.highlight .metric-value {
		color: white;
	}

	.metric-label {
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
	}

	.metric-card.highlight .metric-label {
		color: rgba(255, 255, 255, 0.9);
	}

	.breakdown-section {
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
		margin-top: var(--spacing-md);
	}

	.breakdown-section h4 {
		margin: 0 0 var(--spacing-sm) 0;
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.breakdown-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: var(--spacing-sm);
	}

	.breakdown-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
	}

	.breakdown-label {
		text-transform: capitalize;
		color: var(--color-text-primary);
	}

	.breakdown-value {
		font-weight: var(--font-weight-semibold);
		color: var(--color-primary);
	}
</style>
