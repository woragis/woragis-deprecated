<script lang="ts">
	import Select from '$lib/components/ui/Select.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { ApplicationStatus } from '$lib/api/jobapplications';

	let {
		filters = $bindable({
			status: '',
			website: '',
			interestLevel: ''
		}),
		onClear
	}: {
		filters?: {
			status: string;
			website: string;
			interestLevel: string;
		};
		onClear?: () => void;
	} = $props();

	const statuses: (ApplicationStatus | '')[] = [
		'',
		'pending',
		'processing',
		'applied',
		'contacted',
		'rejected',
		'accepted',
		'failed'
	];

	const interestLevels = [
		{ value: '', label: 'All' },
		{ value: 'low', label: 'Low' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'high', label: 'High' },
		{ value: 'very-high', label: 'Very High' }
	];

	const websites = [
		{ value: '', label: 'All' },
		{ value: 'linkedin', label: 'LinkedIn' },
		{ value: 'glassdoor', label: 'Glassdoor' },
		{ value: 'indeed', label: 'Indeed' },
		{ value: 'monster', label: 'Monster' },
		{ value: 'ziprecruiter', label: 'ZipRecruiter' },
		{ value: 'careerbuilder', label: 'CareerBuilder' }
	];

	function handleClear() {
		filters = {
			status: '',
			website: '',
			interestLevel: ''
		};
		if (onClear) {
			onClear();
		}
	}

	const hasActiveFilters = filters.status || filters.website || filters.interestLevel;
</script>

<div class="filters-container">
	<div class="filters-header">
		<h3>Filters</h3>
		{#if hasActiveFilters}
			<Button variant="secondary" size="sm" onclick={handleClear}>Clear All</Button>
		{/if}
	</div>
	<div class="filters-grid">
		<Select label="Status" bind:value={filters.status}>
			{#each statuses as status}
				<option value={status}>{status || 'All'}</option>
			{/each}
		</Select>
		<Select label="Website" bind:value={filters.website}>
			{#each websites as site}
				<option value={site.value}>{site.label}</option>
			{/each}
		</Select>
		<Select label="Interest Level" bind:value={filters.interestLevel}>
			{#each interestLevels as level}
				<option value={level.value}>{level.label}</option>
			{/each}
		</Select>
	</div>
</div>

<style>
	.filters-container {
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		padding: var(--spacing-md);
		margin-bottom: var(--spacing-md);
	}

	.filters-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-md);
	}

	.filters-header h3 {
		margin: 0;
		font-size: var(--font-size-md);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.filters-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--spacing-md);
	}
</style>

