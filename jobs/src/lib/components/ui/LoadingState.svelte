<script lang="ts">
	let {
		message = 'Loading...'
	}: {
		message?: string;
	} = $props();
</script>

<div class="loading-state">
	<div class="spinner"></div>
	<p class="loading-message">{message}</p>
</div>

<style>
	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: var(--spacing-2xl);
		color: var(--color-text-secondary);
	}
	
	.spinner {
		width: 2rem;
		height: 2rem;
		border: 3px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
		margin-bottom: var(--spacing-md);
	}
	
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	
	.loading-message {
		font-size: var(--font-size-sm);
		color: var(--color-text-tertiary);
	}
</style>
