<script lang="ts">
	let {
		title = '',
		subtitle = ''
	}: {
		title?: string;
		subtitle?: string;
	} = $props();
</script>

<div class="card">
	{#if title || subtitle}
		<div class="card-header">
			{#if title}
				<h3 class="card-title">{title}</h3>
			{/if}
			{#if subtitle}
				<p class="card-subtitle">{subtitle}</p>
			{/if}
		</div>
	{/if}
	<div class="card-content">
		<slot />
	</div>
</div>

<style>
	.card {
		background-color: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		border: 1px solid var(--color-border);
		overflow: hidden;
	}
	
	.card-header {
		padding: var(--spacing-md);
		border-bottom: 1px solid var(--color-border);
	}
	
	.card-title {
		margin: 0;
		font-size: var(--font-size-lg);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}
	
	.card-subtitle {
		margin: var(--spacing-xs) 0 0 0;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
	}
	
	.card-content {
		padding: var(--spacing-md);
	}
</style>
