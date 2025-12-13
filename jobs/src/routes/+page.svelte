<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { listJobApplications, type JobApplication } from '$lib/api/jobapplications';
	import ApplicationMetrics from './applications/_sections/ApplicationMetrics.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Button from '$lib/components/ui/Button.svelte';

	let applications: JobApplication[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchApplications();
	});

	async function fetchApplications() {
		loading = true;
		error = null;
		try {
			applications = await listJobApplications();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load applications';
			console.error('Error fetching job applications:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
		<div>
			<h1 class="text-3xl font-semibold text-gray-900 mb-2">Job Applications Dashboard</h1>
			<p class="text-sm text-gray-600">Overview and analytics of your job applications</p>
		</div>
		<div class="flex gap-3">
			<Button onclick={() => goto('/applications')} variant="primary">
				View All Applications
			</Button>
			<Button onclick={() => goto('/applications')} variant="secondary">
				📋 List
			</Button>
		</div>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchApplications} />
	{:else if loading}
		<LoadingState message="Loading analytics..." />
	{:else}
		<ApplicationMetrics {applications} />
	{/if}
</div>
