<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getJobApplication, type JobApplication } from '$lib/api/jobapplications';
	import InlineChat from '$lib/components/InlineChat.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';

	let application: JobApplication | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const applicationId = $derived($page.params.id);

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

<div class="page-container">
	<div class="header">
		<a href="/job-applications-admin/{applicationId}" class="back-link">← Back to Application</a>
		<div class="header-actions">
			{#if application}
				<Button onclick={() => goto(`/job-applications-admin/${applicationId}`)} variant="secondary">
					View Details
				</Button>
			{/if}
		</div>
	</div>

	{#if loading}
		<LoadingState message="Loading..." />
	{:else if error}
		<ErrorState message={error} onRetry={loadApplication} />
	{:else if application}
		<div class="chat-page-container">
			<div class="chat-header">
				<div class="chat-header-content">
					<h1 class="chat-title">{application.companyName} - {application.jobTitle}</h1>
					{#if application.jobDescription}
						<p class="chat-description">{application.jobDescription}</p>
					{/if}
				</div>
			</div>
			
			<div class="chat-wrapper">
				<InlineChat
					jobApplicationId={application.id}
					title={`${application.jobTitle} at ${application.companyName}`}
					description={application.jobDescription || `Chat about the ${application.jobTitle} position at ${application.companyName}`}
				/>
			</div>
		</div>
	{/if}
</div>

<style>
	.page-container {
		padding: var(--spacing-md);
		max-width: 1400px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}

	.back-link {
		color: var(--color-primary);
		text-decoration: none;
		font-size: var(--font-size-sm);
		transition: color var(--transition-base);
	}

	.back-link:hover {
		color: var(--color-primary-hover);
		text-decoration: underline;
	}

	.header-actions {
		display: flex;
		gap: var(--spacing-sm);
	}

	.chat-page-container {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 200px);
		min-height: 600px;
	}

	.chat-header {
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		border: 1px solid var(--color-border);
		padding: var(--spacing-lg);
		margin-bottom: var(--spacing-lg);
	}

	.chat-header-content {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.chat-title {
		margin: 0;
		font-size: var(--font-size-2xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.chat-description {
		margin: 0;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		line-height: 1.5;
	}

	.chat-wrapper {
		flex: 1;
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		border: 1px solid var(--color-border);
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	/* Ensure InlineChat fills the container */
	.chat-wrapper :global(> *) {
		flex: 1;
		min-height: 0;
		display: flex;
		flex-direction: column;
	}
</style>
