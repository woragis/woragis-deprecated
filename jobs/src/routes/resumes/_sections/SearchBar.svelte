<script lang="ts">
	import { Search } from 'lucide-svelte';
	
	let {
		searchQuery = $bindable('')
	}: {
		searchQuery?: string;
	} = $props();
</script>

<div class="search-bar">
	<div class="search-input-wrapper">
		<Search size={20} class="search-icon" />
		<input
			type="text"
			placeholder="Search by title, filename, or tags..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>
</div>

<style>
	.search-bar {
		margin-bottom: 1rem;
	}

	.search-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.search-icon {
		position: absolute;
		left: 0.75rem;
		color: #6b7280;
		pointer-events: none;
	}

	.search-input {
		width: 100%;
		padding: 0.75rem 0.75rem 0.75rem 2.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		background-color: white;
		color: #1f2937;
		transition: border-color 0.2s;
	}

	.search-input:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}
</style>
