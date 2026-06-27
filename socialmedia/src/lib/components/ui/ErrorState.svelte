<script lang="ts">
	let {
		message = 'An error occurred',
		onRetry
	}: {
		message?: string;
		onRetry?: () => void;
	} = $props();
</script>

<div class="error-state">
	<div class="error-icon">⚠️</div>
	<p class="error-message">{message}</p>
	{#if onRetry}
		<button class="retry-button" onclick={onRetry} type="button">Try Again</button>
	{/if}
</div>

<style>
	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: var(--spacing-2xl);
		background-color: var(--color-danger-light);
		border: 1px solid var(--color-danger);
		border-radius: var(--radius-lg);
		color: var(--color-danger);
	}
	
	.error-icon {
		font-size: var(--font-size-2xl);
		margin-bottom: var(--spacing-md);
	}
	
	.error-message {
		font-size: var(--font-size-sm);
		margin-bottom: var(--spacing-md);
		text-align: center;
	}
	
	.retry-button {
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: var(--color-danger);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		cursor: pointer;
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
		transition: background-color var(--transition-base);
	}
	
	.retry-button:hover {
		background-color: var(--color-danger-hover);
	}
</style>
