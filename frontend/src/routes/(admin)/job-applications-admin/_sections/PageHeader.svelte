<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	
	let {
		onCreateClick
	}: {
		onCreateClick?: () => void;
	} = $props();
	
	const tFn = useTranslation();
</script>

<div class="page-header">
	<div class="header-content">
		<h1 class="page-title">{tFn('jobApplications.title')}</h1>
		<p class="page-subtitle">{tFn('jobApplications.subtitle')}</p>
	</div>
	{#if onCreateClick}
		<Button onclick={onCreateClick} variant="primary">
			{tFn('jobApplications.createButton')}
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
