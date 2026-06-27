<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	
	let {
		title,
		subtitle,
		onCreateClick
	}: {
		title: string;
		subtitle?: string;
		onCreateClick?: () => void;
	} = $props();
</script>

<div class="page-header">
	<div class="header-content">
		<h1 class="page-title">{title}</h1>
		{#if subtitle}
			<p class="page-subtitle">{subtitle}</p>
		{/if}
	</div>
	{#if onCreateClick}
		<Button onclick={onCreateClick} variant="primary">
			Create New
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
