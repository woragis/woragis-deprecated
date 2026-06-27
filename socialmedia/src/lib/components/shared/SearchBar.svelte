<script lang="ts">
	let {
		searchQuery = $bindable(''),
		placeholder = 'Search...'
	}: {
		searchQuery?: string;
		placeholder?: string;
	} = $props();
</script>

<div class="search-bar">
	<input
		type="text"
		class="search-input"
		{placeholder}
		bind:value={searchQuery}
	/>
</div>

<style>
	.search-bar {
		margin-bottom: var(--spacing-md);
	}
	
	.search-input {
		width: 100%;
		max-width: 400px;
		padding: var(--spacing-sm) var(--spacing-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		background-color: var(--color-bg-primary);
		transition: border-color var(--transition-base);
	}
	
	.search-input:focus {
		outline: none;
		border-color: var(--color-primary);
	}
</style>
