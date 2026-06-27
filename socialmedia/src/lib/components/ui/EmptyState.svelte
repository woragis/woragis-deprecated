<script lang="ts">
	let {
		title = 'No items found',
		description = ''
	}: {
		title?: string;
		description?: string;
	} = $props();
</script>

<div class="empty-state">
	<div class="empty-icon">📭</div>
	<h3 class="empty-title">{title}</h3>
	{#if description}
		<p class="empty-description">{description}</p>
	{/if}
	<slot />
</div>

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: var(--spacing-2xl);
		text-align: center;
		color: var(--color-text-secondary);
	}
	
	.empty-icon {
		font-size: 3rem;
		margin-bottom: var(--spacing-md);
		opacity: 0.5;
	}
	
	.empty-title {
		margin: 0 0 var(--spacing-sm) 0;
		font-size: var(--font-size-lg);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}
	
	.empty-description {
		margin: 0;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
	}
</style>
