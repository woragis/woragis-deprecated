<script lang="ts">
	import type { Platform, PostStatus } from '$lib/api/socialmediaposts';

	let {
		platform = $bindable<Platform | undefined>(undefined),
		status = $bindable<PostStatus | undefined>(undefined)
	}: {
		platform?: Platform | undefined;
		status?: PostStatus | undefined;
	} = $props();

	const platforms: Platform[] = ['linkedin', 'twitter', 'instagram', 'medium', 'substack', 'valete', 'website'];
	const statuses: PostStatus[] = ['draft', 'ready', 'scheduled', 'posted', 'analyzed', 'archived'];
</script>

<div class="filters">
	<div class="filter-group">
		<label for="platform-filter">Platform</label>
		<select id="platform-filter" bind:value={platform} class="filter-select">
			<option value={undefined}>All Platforms</option>
			{#each platforms as p}
				<option value={p}>{p.charAt(0).toUpperCase() + p.slice(1)}</option>
			{/each}
		</select>
	</div>
	<div class="filter-group">
		<label for="status-filter">Status</label>
		<select id="status-filter" bind:value={status} class="filter-select">
			<option value={undefined}>All Statuses</option>
			{#each statuses as s}
				<option value={s}>{s.charAt(0).toUpperCase() + s.slice(1)}</option>
			{/each}
		</select>
	</div>
	{#if platform || status}
		<button class="clear-filters" onclick={() => { platform = undefined; status = undefined; }}>
			Clear Filters
		</button>
	{/if}
</div>

<style>
	.filters {
		display: flex;
		gap: var(--spacing-md);
		align-items: flex-end;
		margin-bottom: var(--spacing-md);
		flex-wrap: wrap;
	}

	.filter-group {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.filter-group label {
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
		color: var(--color-text-secondary);
	}

	.filter-select {
		padding: var(--spacing-sm) var(--spacing-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		min-width: 150px;
	}

	.filter-select:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.clear-filters {
		padding: var(--spacing-sm) var(--spacing-md);
		background-color: transparent;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-base);
	}

	.clear-filters:hover {
		background-color: var(--color-bg-hover);
		color: var(--color-text-primary);
	}
</style>
