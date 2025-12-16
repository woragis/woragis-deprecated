<script lang="ts">
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resumes = []
	}: {
		resumes?: Resume[];
	} = $props();

	let totalResumes = $derived(resumes.length);
	let mainResumes = $derived(resumes.filter(r => r.isMain).length);
	let featuredResumes = $derived(resumes.filter(r => r.isFeatured).length);
	let totalUsage = $derived(resumes.reduce((sum, r) => sum + r.applicationsUsed, 0));
	let avgInterviewRate = $derived(
		resumes.length > 0
			? resumes.reduce((sum, r) => sum + r.interviewRate, 0) / resumes.length
			: 0
	);
	let avgOfferRate = $derived(
		resumes.length > 0
			? resumes.reduce((sum, r) => sum + r.offerRate, 0) / resumes.length
			: 0
	);
</script>

<div class="metrics">
	<div class="metric-card">
		<div class="metric-value">{totalResumes}</div>
		<div class="metric-label">Total CVs</div>
	</div>
	<div class="metric-card">
		<div class="metric-value">{mainResumes}</div>
		<div class="metric-label">Main CV</div>
	</div>
	<div class="metric-card">
		<div class="metric-value">{featuredResumes}</div>
		<div class="metric-label">Featured</div>
	</div>
	<div class="metric-card">
		<div class="metric-value">{totalUsage}</div>
		<div class="metric-label">Total Usage</div>
	</div>
	<div class="metric-card">
		<div class="metric-value">{avgInterviewRate.toFixed(1)}%</div>
		<div class="metric-label">Avg Interview Rate</div>
	</div>
	<div class="metric-card">
		<div class="metric-value">{avgOfferRate.toFixed(1)}%</div>
		<div class="metric-label">Avg Offer Rate</div>
	</div>
</div>

<style>
	.metrics {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 1rem;
		margin-bottom: 2rem;
	}

	.metric-card {
		background-color: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 1.25rem;
		text-align: center;
	}

	.metric-value {
		font-size: 1.875rem;
		font-weight: 600;
		color: #1f2937;
		margin-bottom: 0.25rem;
	}

	.metric-label {
		font-size: 0.875rem;
		color: #6b7280;
	}
</style>
