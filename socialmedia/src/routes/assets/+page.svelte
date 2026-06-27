<script lang="ts">
	import { onMount } from 'svelte';
	import { listAssets, deleteAsset, type ContentAsset } from '$lib/api/assets';
	import PageHeader from '../posts/_sections/PageHeader.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';

	let assets: ContentAsset[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchAssets();
	});

	async function fetchAssets() {
		loading = true;
		error = null;
		try {
			assets = await listAssets();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load assets';
			console.error('Error fetching assets:', err);
		} finally {
			loading = false;
		}
	}

	function handleDelete(id: string) {
		if (confirm('Are you sure you want to delete this asset?')) {
			deleteAsset(id)
				.then(() => {
					assets = assets.filter((a) => a.id !== id);
				})
				.catch((err) => {
					console.error('Error deleting asset:', err);
					alert('Failed to delete asset');
				});
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<PageHeader onCreateClick={() => {}} />
	
	{#if error}
		<ErrorState message={error} onRetry={fetchAssets} />
	{:else if loading}
		<LoadingState message="Loading assets..." />
	{:else if assets.length === 0}
		<EmptyState
			title="No assets found"
			description="Upload your first asset to get started"
		/>
	{:else}
		<div class="assets-grid">
			{#each assets as asset}
				<div class="asset-card">
					<div class="asset-header">
						<span class="asset-type">{asset.assetType}</span>
						<button class="delete-button" onclick={() => handleDelete(asset.id)}>×</button>
					</div>
					{#if asset.fileUrl}
						<div class="asset-preview">
							{#if asset.assetType === 'image'}
								<img src={asset.fileUrl} alt={asset.altText || 'Asset'} class="preview-image" />
							{:else}
								<div class="preview-placeholder">{asset.assetType}</div>
							{/if}
						</div>
					{/if}
					<div class="asset-info">
						<div class="asset-path">{asset.filePath}</div>
						{#if asset.altText}
							<div class="asset-alt">{asset.altText}</div>
						{/if}
						<div class="asset-date">{formatDate(asset.createdAt)}</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.assets-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: var(--spacing-lg);
	}

	.asset-card {
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
		transition: transform var(--transition-base), box-shadow var(--transition-base);
	}

	.asset-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.asset-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border-bottom: 1px solid var(--color-border);
	}

	.asset-type {
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		text-transform: capitalize;
		color: var(--color-text-secondary);
	}

	.delete-button {
		background: none;
		border: none;
		font-size: 1.5rem;
		color: var(--color-text-secondary);
		cursor: pointer;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.delete-button:hover {
		color: var(--color-danger);
	}

	.asset-preview {
		width: 100%;
		height: 200px;
		background-color: var(--color-bg-secondary);
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
	}

	.preview-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.preview-placeholder {
		color: var(--color-text-tertiary);
		font-size: var(--font-size-sm);
		text-transform: capitalize;
	}

	.asset-info {
		padding: var(--spacing-md);
	}

	.asset-path {
		font-size: var(--font-size-xs);
		color: var(--color-text-secondary);
		margin-bottom: var(--spacing-xs);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.asset-alt {
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		margin-bottom: var(--spacing-xs);
	}

	.asset-date {
		font-size: var(--font-size-xs);
		color: var(--color-text-tertiary);
	}
</style>
