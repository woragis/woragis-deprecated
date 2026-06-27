<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getDashboardData, type DashboardData } from '$lib/api/socialmediaposts';
	import DashboardMetrics from './_sections/DashboardMetrics.svelte';
	import UpcomingPosts from './_sections/UpcomingPosts.svelte';
	import ContentBacklog from './_sections/ContentBacklog.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Button from '$lib/components/ui/Button.svelte';

	let dashboardData: DashboardData | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchDashboardData();
	});

	async function fetchDashboardData() {
		loading = true;
		error = null;
		try {
			dashboardData = await getDashboardData();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load dashboard data';
			console.error('Error fetching dashboard data:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
		<div>
			<h1 class="text-3xl font-semibold text-gray-900 mb-2">Social Media Dashboard</h1>
			<p class="text-sm text-gray-600">Overview and analytics of your social media posts</p>
		</div>
		<div class="flex gap-3">
			<Button onclick={() => goto('/posts')} variant="primary">
				View All Posts
			</Button>
			<Button onclick={() => goto('/content/backlog')} variant="secondary">
				Content Backlog
			</Button>
		</div>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchDashboardData} />
	{:else if loading}
		<LoadingState message="Loading dashboard..." />
	{:else if dashboardData}
		<div class="dashboard-grid">
			<div class="dashboard-main">
				<DashboardMetrics {dashboardData} />
			</div>
			<div class="dashboard-sidebar">
				<UpcomingPosts upcomingPosts={dashboardData.upcomingPosts} />
				<ContentBacklog contentBacklog={dashboardData.contentBacklog} />
			</div>
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1280px;
	}

	.dashboard-grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--spacing-lg);
	}

	@media (min-width: 1024px) {
		.dashboard-grid {
			grid-template-columns: 2fr 1fr;
		}
	}

	.dashboard-main {
		min-width: 0;
	}

	.dashboard-sidebar {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-lg);
	}
</style>
