<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	
	let {
		onCreateClick
	}: {
		onCreateClick?: () => void;
	} = $props();
</script>

<div class="page-header">
	<div class="header-content">
		<h1 class="page-title">Social Media Posts</h1>
		<p class="page-subtitle">Manage and track your social media content</p>
	</div>
	{#if onCreateClick}
		<Button onclick={onCreateClick} variant="primary">
			Create Post
		</Button>
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
</style>
