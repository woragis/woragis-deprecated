<script lang="ts">
	import TagInput from '$lib/components/ui/TagInput.svelte';
	
	let {
		filters = $bindable({
			tags: [] as string[],
			main: false,
			featured: false
		}),
		availableTags = []
	}: {
		filters?: {
			tags: string[];
			main: boolean;
			featured: boolean;
		};
		availableTags?: string[];
	} = $props();
</script>

<div class="filters">
	<div class="filter-group">
		<label class="filter-label">Tags:</label>
		<TagInput
			bind:tags={filters.tags}
			{availableTags}
			placeholder="Filter by tags..."
		/>
	</div>
	
	<div class="filter-group">
		<label class="filter-checkbox">
			<input
				type="checkbox"
				bind:checked={filters.main}
				class="checkbox-input"
			/>
			<span>Main CV only</span>
		</label>
	</div>
	
	<div class="filter-group">
		<label class="filter-checkbox">
			<input
				type="checkbox"
				bind:checked={filters.featured}
				class="checkbox-input"
			/>
			<span>Featured only</span>
		</label>
	</div>
</div>

<style>
	.filters {
		display: flex;
		gap: 1.5rem;
		align-items: center;
		margin-bottom: 1rem;
		flex-wrap: wrap;
	}

	.filter-group {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.filter-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
		white-space: nowrap;
	}

	.filter-checkbox {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		font-size: 0.875rem;
		color: #1f2937;
	}

	.checkbox-input {
		width: 1rem;
		height: 1rem;
		cursor: pointer;
	}
</style>
