<script lang="ts">
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resume
	}: {
		resume: Resume;
	} = $props();
</script>

<div class="section">
	<h2 class="section-title">Metrics</h2>
	<div class="metrics-grid">
		<div class="metric-item">
			<div class="metric-value">{resume.applicationsUsed}</div>
			<div class="metric-label">Applications Used</div>
		</div>
		<div class="metric-item">
			<div class="metric-value">{resume.interviewRate.toFixed(1)}%</div>
			<div class="metric-label">Interview Rate</div>
		</div>
		<div class="metric-item">
			<div class="metric-value">{resume.offerRate.toFixed(1)}%</div>
			<div class="metric-label">Offer Rate</div>
		</div>
	</div>
</div>

<style>
	.metrics-grid {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.metric-item {
		text-align: center;
		padding: 1rem;
		background-color: #f9fafb;
		border-radius: 0.375rem;
	}

	.metric-value {
		font-size: 1.5rem;
		font-weight: 600;
		color: #1f2937;
		margin-bottom: 0.25rem;
	}

	.metric-label {
		font-size: 0.875rem;
		color: #6b7280;
	}
</style>
