<script lang="ts">
	import type { JobApplication } from '$lib/api/jobapplications';

	let { applications = [] }: { applications?: JobApplication[] } = $props();

	const metrics = $derived(() => {
		const total = applications.length;
		const byStatus = {
			pending: applications.filter(a => a.status === 'pending').length,
			processing: applications.filter(a => a.status === 'processing').length,
			applied: applications.filter(a => a.status === 'applied').length,
			contacted: applications.filter(a => a.status === 'contacted').length,
			rejected: applications.filter(a => a.status === 'rejected').length,
			accepted: applications.filter(a => a.status === 'accepted').length,
			failed: applications.filter(a => a.status === 'failed').length
		};

		const byWebsite = applications.reduce((acc, app) => {
			const website = app.website.toLowerCase();
			acc[website] = (acc[website] || 0) + 1;
			return acc;
		}, {} as Record<string, number>);

		const withResponse = applications.filter(a => a.responseReceivedAt).length;
		const responseRate = total > 0 ? (withResponse / total) * 100 : 0;

		const highInterest = applications.filter(a => 
			a.interestLevel === 'high' || a.interestLevel === 'very-high'
		).length;

		const needsFollowUp = applications.filter(a => {
			if (!a.followUpDate) return false;
			const followUp = new Date(a.followUpDate);
			const today = new Date();
			today.setHours(0, 0, 0, 0);
			followUp.setHours(0, 0, 0, 0);
			return followUp <= today && a.status !== 'rejected' && a.status !== 'accepted';
		}).length;

		return {
			total,
			byStatus,
			byWebsite,
			withResponse,
			responseRate: Math.round(responseRate * 10) / 10,
			highInterest,
			needsFollowUp
		};
	});
</script>

<div class="metrics-container">
	<div class="metrics-grid">
		<div class="metric-card">
			<div class="metric-value">{metrics().total}</div>
			<div class="metric-label">Total Applications</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().applied}</div>
			<div class="metric-label">Applied</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().contacted}</div>
			<div class="metric-label">Contacted</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().accepted}</div>
			<div class="metric-label">Accepted</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().responseRate}%</div>
			<div class="metric-label">Response Rate</div>
		</div>
		<div class="metric-card">
			<div class="metric-value">{metrics().highInterest}</div>
			<div class="metric-label">High Interest</div>
		</div>
		<div class="metric-card highlight">
			<div class="metric-value">{metrics().needsFollowUp}</div>
			<div class="metric-label">Needs Follow-up</div>
		</div>
	</div>
	{#if Object.keys(metrics().byWebsite).length > 0}
		<div class="website-breakdown">
			<h4>By Website</h4>
			<div class="website-list">
				{#each Object.entries(metrics().byWebsite) as [website, count]}
					<div class="website-item">
						<span class="website-name">{website}</span>
						<span class="website-count">{count}</span>
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

	.website-breakdown {
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
	}

	.website-breakdown h4 {
		margin: 0 0 var(--spacing-sm) 0;
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.website-list {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-sm);
	}

	.website-item {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
	}

	.website-name {
		text-transform: capitalize;
		color: var(--color-text-primary);
	}

	.website-count {
		font-weight: var(--font-weight-semibold);
		color: var(--color-primary);
	}
</style>

