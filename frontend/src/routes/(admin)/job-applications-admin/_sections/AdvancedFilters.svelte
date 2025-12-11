<script lang="ts">
	import Select from '$lib/components/ui/Select.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { ApplicationStatus } from '$lib/api/jobapplications';

	const FILTER_PRESETS_KEY = 'jobApplicationFilterPresets';

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

	let showPresets = $state(false);
	let presetName = $state('');
	let savedPresets = $state<Array<{ name: string; filters: typeof filters }>>([]);

	$effect(() => {
		if (typeof window !== 'undefined') {
			const stored = localStorage.getItem(FILTER_PRESETS_KEY);
			if (stored) {
				try {
					savedPresets = JSON.parse(stored);
				} catch (e) {
					console.warn('Failed to load filter presets:', e);
				}
			}
		}
	});

	function savePreset() {
		if (!presetName.trim()) return;
		
		const preset = {
			name: presetName.trim(),
			filters: { ...filters }
		};
		
		savedPresets = [...savedPresets, preset];
		
		if (typeof window !== 'undefined') {
			localStorage.setItem(FILTER_PRESETS_KEY, JSON.stringify(savedPresets));
		}
		
		presetName = '';
		showPresets = false;
	}

	function loadPreset(preset: typeof savedPresets[0]) {
		filters = { ...preset.filters };
	}

	function deletePreset(index: number) {
		savedPresets = savedPresets.filter((_, i) => i !== index);
		if (typeof window !== 'undefined') {
			localStorage.setItem(FILTER_PRESETS_KEY, JSON.stringify(savedPresets));
		}
	}

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
		<div class="filters-actions">
			<Button variant="secondary" size="sm" onclick={() => showPresets = !showPresets}>
				Presets
			</Button>
			{#if hasActiveFilters}
				<Button variant="secondary" size="sm" onclick={handleClear}>Clear All</Button>
			{/if}
		</div>
	</div>
	{#if showPresets}
		<div class="presets-section">
			<div class="presets-header">
				<h4>Saved Presets</h4>
				<Button variant="primary" size="sm" onclick={savePreset} disabled={!presetName.trim() || !hasActiveFilters}>
					Save Current
				</Button>
			</div>
			<div class="presets-input">
				<Input
					placeholder="Preset name..."
					bind:value={presetName}
				/>
			</div>
			{#if savedPresets.length > 0}
				<div class="presets-list">
					{#each savedPresets as preset, index}
						<div class="preset-item">
							<button class="preset-button" onclick={() => loadPreset(preset)}>
								{preset.name}
							</button>
							<button class="preset-delete" onclick={() => deletePreset(index)}>×</button>
						</div>
					{/each}
				</div>
			{:else}
				<p class="presets-empty">No saved presets. Save your current filters as a preset.</p>
			{/if}
		</div>
	{/if}
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

	.filters-actions {
		display: flex;
		gap: var(--spacing-sm);
	}

	.filters-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: var(--spacing-md);
	}

	.presets-section {
		margin-top: var(--spacing-md);
		padding-top: var(--spacing-md);
		border-top: 1px solid var(--color-border);
	}

	.presets-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-sm);
	}

	.presets-header h4 {
		margin: 0;
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-semibold);
	}

	.presets-input {
		margin-bottom: var(--spacing-sm);
	}

	.presets-list {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-sm);
	}

	.preset-item {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
	}

	.preset-button {
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		cursor: pointer;
		font-size: var(--font-size-sm);
		transition: background-color var(--transition-base);
	}

	.preset-button:hover {
		background-color: var(--color-bg-hover);
	}

	.preset-delete {
		background: none;
		border: none;
		color: var(--color-text-secondary);
		cursor: pointer;
		font-size: 1.2rem;
		padding: 0;
		width: 20px;
		height: 20px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.preset-delete:hover {
		color: var(--color-danger);
	}

	.presets-empty {
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
		margin: var(--spacing-sm) 0 0 0;
	}
</style>

