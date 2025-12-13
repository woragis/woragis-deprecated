<script lang="ts">
	let {
		title = 'No items yet',
		description,
		icon
	}: {
		title?: string;
		description?: string;
		icon?: string;
	} = $props();
</script>

<div class="empty-state">
	{#if icon}
		<div class="empty-icon">{@html icon}</div>
	{/if}
	<p class="empty-title">{title}</p>
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
		border: 1px dashed var(--color-border);
		border-radius: var(--radius-lg);
		background-color: var(--color-bg-secondary);
	}
	
	.empty-icon {
		font-size: var(--font-size-3xl);
		margin-bottom: var(--spacing-md);
		opacity: 0.6;
	}
	
	.empty-title {
		font-size: var(--font-size-base);
		font-weight: var(--font-weight-medium);
		color: var(--color-text-secondary);
		margin: 0;
	}
	
	.empty-description {
		font-size: var(--font-size-sm);
		color: var(--color-text-tertiary);
		margin: var(--spacing-sm) 0 0 0;
	}
</style>
