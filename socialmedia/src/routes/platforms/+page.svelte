<script lang="ts">
	import { onMount } from 'svelte';
	import { listPlatforms, type PlatformConfig } from '$lib/api/platforms';
	import PageHeader from '../posts/_sections/PageHeader.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Card from '$lib/components/ui/Card.svelte';

	let platforms: PlatformConfig[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchPlatforms();
	});

	async function fetchPlatforms() {
		loading = true;
		error = null;
		try {
			platforms = await listPlatforms();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load platforms';
			console.error('Error fetching platforms:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="page-header">
		<div class="header-content">
			<h1 class="page-title">Platforms</h1>
			<p class="page-subtitle">Configure and manage your social media platforms</p>
		</div>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchPlatforms} />
	{:else if loading}
		<LoadingState message="Loading platforms..." />
	{:else}
		<div class="platforms-grid">
			{#each platforms as platform}
				<Card title={platform.displayName} subtitle={platform.name}>
					<div class="platform-info">
						<div class="info-item">
							<span class="info-label">Status:</span>
							<span class="info-value {platform.isActive ? 'active' : 'inactive'}">
								{platform.isActive ? 'Active' : 'Inactive'}
							</span>
						</div>
						{#if platform.postingFrequency}
							<div class="info-item">
								<span class="info-label">Posting Frequency:</span>
								<span class="info-value">{platform.postingFrequency} posts/week</span>
							</div>
						{/if}
						{#if platform.supportedFormats && platform.supportedFormats.length > 0}
							<div class="info-item">
								<span class="info-label">Supported Formats:</span>
								<div class="formats-list">
									{#each platform.supportedFormats as format}
										<span class="format-tag">{format}</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</Card>
			{/each}
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

	.platforms-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: var(--spacing-lg);
	}

	.platform-info {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}

	.info-item {
		display: flex;
		align-items: flex-start;
		gap: var(--spacing-sm);
		font-size: var(--font-size-sm);
	}

	.info-label {
		font-weight: var(--font-weight-medium);
		color: var(--color-text-secondary);
		min-width: 120px;
	}

	.info-value {
		color: var(--color-text-primary);
	}

	.info-value.active {
		color: var(--color-success);
		font-weight: var(--font-weight-semibold);
	}

	.info-value.inactive {
		color: var(--color-text-tertiary);
	}

	.formats-list {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
	}

	.format-tag {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-tertiary);
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		text-transform: capitalize;
	}
</style>
