<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getJobApplication, type JobApplication } from '$lib/api/jobapplications';
	import InlineChat from '$lib/components/InlineChat.svelte';

	let application: JobApplication | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	$: applicationId = $page.params.id;

	onMount(async () => {
		if (applicationId) {
			await loadApplication();
		}
	});

	async function loadApplication() {
		if (!applicationId) return;
		loading = true;
		error = null;
		try {
			application = await getJobApplication(applicationId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load job application';
			console.error('Error loading application:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-4 py-8">
	{#if loading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600 dark:text-gray-400">Loading...</p>
		</div>
	{:else if error}
		<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
			<p class="text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if application}
		<div class="mb-4">
			<button
				onclick={() => goto('/job-applications-admin')}
				class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 mb-4"
			>
				← Back to Applications
			</button>
		</div>
		<div class="mb-6">
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white">
				Chat: {application.companyName} - {application.jobTitle}
			</h1>
		</div>
		<div style="height: 600px;">
			<InlineChat
				jobApplicationId={application.id}
				title={`${application.jobTitle} at ${application.companyName}`}
				description={application.jobDescription || `Chat about the ${application.jobTitle} position at ${application.companyName}`}
			/>
		</div>
	{/if}
</div>

