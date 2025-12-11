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

		// Website breakdown with success rates
		const byWebsite = applications.reduce((acc, app) => {
			const website = app.website.toLowerCase();
			if (!acc[website]) {
				acc[website] = { total: 0, contacted: 0, accepted: 0, rejected: 0, withResponse: 0 };
			}
			acc[website].total++;
			if (app.status === 'contacted') acc[website].contacted++;
			if (app.status === 'accepted') acc[website].accepted++;
			if (app.status === 'rejected') acc[website].rejected++;
			if (app.responseReceivedAt) acc[website].withResponse++;
			return acc;
		}, {} as Record<string, { total: number; contacted: number; accepted: number; rejected: number; withResponse: number }>);

		// Calculate success rates by website
		const websiteSuccessRates = Object.entries(byWebsite).map(([website, data]) => {
			const responseRate = data.total > 0 ? (data.withResponse / data.total) * 100 : 0;
			const contactRate = data.total > 0 ? (data.contacted / data.total) * 100 : 0;
			const acceptanceRate = data.contacted > 0 ? (data.accepted / data.contacted) * 100 : 0;
			return {
				website,
				...data,
				responseRate: Math.round(responseRate * 10) / 10,
				contactRate: Math.round(contactRate * 10) / 10,
				acceptanceRate: Math.round(acceptanceRate * 10) / 10
			};
		}).sort((a, b) => b.responseRate - a.responseRate);

		// Time to response calculation
		const responseTimes = applications
			.filter(a => a.appliedAt && a.responseReceivedAt)
			.map(a => {
				const applied = new Date(a.appliedAt).getTime();
				const responded = new Date(a.responseReceivedAt).getTime();
				return (responded - applied) / (1000 * 60 * 60 * 24); // days
			});
		const avgTimeToResponse = responseTimes.length > 0
			? Math.round((responseTimes.reduce((a, b) => a + b, 0) / responseTimes.length) * 10) / 10
			: 0;

		// Conversion funnel
		const applied = byStatus.applied + byStatus.contacted + byStatus.accepted + byStatus.rejected;
		const contacted = byStatus.contacted + byStatus.accepted;
		const accepted = byStatus.accepted;
		
		const funnel = {
			applied: {
				count: applied,
				percentage: 100
			},
			contacted: {
				count: contacted,
				percentage: applied > 0 ? Math.round((contacted / applied) * 100 * 10) / 10 : 0
			},
			accepted: {
				count: accepted,
				percentage: contacted > 0 ? Math.round((accepted / contacted) * 100 * 10) / 10 : 0
			}
		};

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
			byWebsite: websiteSuccessRates,
			withResponse,
			responseRate: Math.round(responseRate * 10) / 10,
			highInterest,
			needsFollowUp,
			avgTimeToResponse,
			funnel
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
		{#if metrics().avgTimeToResponse > 0}
			<div class="metric-card">
				<div class="metric-value">{metrics().avgTimeToResponse}</div>
				<div class="metric-label">Avg Days to Response</div>
			</div>
		{/if}
	</div>
	
	<!-- Conversion Funnel -->
	{#if metrics().funnel.applied.count > 0}
		<div class="funnel-section">
			<h4>Conversion Funnel</h4>
			<div class="funnel-container">
				<div class="funnel-stage">
					<div class="funnel-label">Applied</div>
					<div class="funnel-bar">
						<div class="funnel-fill" style="width: {metrics().funnel.applied.percentage}%"></div>
					</div>
					<div class="funnel-stats">
						<span class="funnel-count">{metrics().funnel.applied.count}</span>
						<span class="funnel-percentage">{metrics().funnel.applied.percentage}%</span>
					</div>
				</div>
				<div class="funnel-arrow">↓</div>
				<div class="funnel-stage">
					<div class="funnel-label">Contacted</div>
					<div class="funnel-bar">
						<div class="funnel-fill" style="width: {metrics().funnel.contacted.percentage}%"></div>
					</div>
					<div class="funnel-stats">
						<span class="funnel-count">{metrics().funnel.contacted.count}</span>
						<span class="funnel-percentage">{metrics().funnel.contacted.percentage}%</span>
					</div>
				</div>
				<div class="funnel-arrow">↓</div>
				<div class="funnel-stage">
					<div class="funnel-label">Accepted</div>
					<div class="funnel-bar">
						<div class="funnel-fill accepted" style="width: {metrics().funnel.accepted.percentage}%"></div>
					</div>
					<div class="funnel-stats">
						<span class="funnel-count">{metrics().funnel.accepted.count}</span>
						<span class="funnel-percentage">{metrics().funnel.accepted.percentage}%</span>
					</div>
				</div>
			</div>
		</div>
	{/if}
	
	<!-- Website Success Rates -->
	{#if metrics().byWebsite.length > 0}
		<div class="website-breakdown">
			<h4>Success Rate by Website</h4>
			<div class="website-table">
				<div class="website-header">
					<span>Website</span>
					<span>Total</span>
					<span>Response Rate</span>
					<span>Contact Rate</span>
					<span>Acceptance Rate</span>
				</div>
				{#each metrics().byWebsite as site}
					<div class="website-row">
						<span class="website-name">{site.website}</span>
						<span class="website-stat">{site.total}</span>
						<span class="website-stat">{site.responseRate}%</span>
						<span class="website-stat">{site.contactRate}%</span>
						<span class="website-stat">{site.acceptanceRate}%</span>
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

	.funnel-section {
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
		margin-top: var(--spacing-md);
	}

	.funnel-section h4 {
		margin: 0 0 var(--spacing-md) 0;
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.funnel-container {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}

	.funnel-stage {
		display: flex;
		align-items: center;
		gap: var(--spacing-md);
	}

	.funnel-label {
		min-width: 100px;
		font-weight: var(--font-weight-medium);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
	}

	.funnel-bar {
		flex: 1;
		height: 30px;
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-md);
		overflow: hidden;
		position: relative;
	}

	.funnel-fill {
		height: 100%;
		background-color: var(--color-primary);
		transition: width 0.3s ease;
	}

	.funnel-fill.accepted {
		background-color: var(--color-success, #10b981);
	}

	.funnel-stats {
		display: flex;
		gap: var(--spacing-sm);
		min-width: 120px;
		justify-content: flex-end;
		font-size: var(--font-size-sm);
	}

	.funnel-count {
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.funnel-percentage {
		color: var(--color-text-secondary);
	}

	.funnel-arrow {
		text-align: center;
		color: var(--color-text-secondary);
		font-size: 1.5rem;
		margin: var(--spacing-xs) 0;
	}

	.website-table {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.website-header {
		display: grid;
		grid-template-columns: 1fr 80px 120px 120px 120px;
		gap: var(--spacing-md);
		padding: var(--spacing-sm);
		font-weight: var(--font-weight-semibold);
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		border-bottom: 2px solid var(--color-border);
	}

	.website-row {
		display: grid;
		grid-template-columns: 1fr 80px 120px 120px 120px;
		gap: var(--spacing-md);
		padding: var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
	}

	.website-stat {
		text-align: right;
		color: var(--color-text-primary);
	}
</style>

