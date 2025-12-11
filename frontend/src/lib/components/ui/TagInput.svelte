<script lang="ts">
	import { X } from 'lucide-svelte';
	
	let {
		tags = $bindable([]),
		availableTags: availableTagsProp = [],
		placeholder = 'Add tags...',
		maxTags = 10,
		onFetchTags: onFetchTagsProp = null
	}: {
		tags?: string[];
		availableTags?: string[];
		placeholder?: string;
		maxTags?: number;
		onFetchTags?: (() => Promise<string[]>) | null;
	} = $props();

	let availableTags = $state(availableTagsProp);
	let onFetchTags = $state<(() => Promise<string[]>) | null>(onFetchTagsProp);
	
	let inputValue = $state('');
	let showSuggestions = $state(false);
	let filteredSuggestions = $state<string[]>([]);
	let selectedIndex = $state(-1);
	
	$effect(() => {
		if (inputValue.trim()) {
			const query = inputValue.toLowerCase().trim();
			filteredSuggestions = availableTags
				.filter(tag => 
					tag.toLowerCase().includes(query) && 
					!tags.includes(tag)
				)
				.slice(0, 10);
			showSuggestions = filteredSuggestions.length > 0;
		} else {
			showSuggestions = false;
			filteredSuggestions = [];
		}
	});
	
	async function loadAvailableTags() {
		if (onFetchTags) {
			try {
				availableTags = await onFetchTags();
			} catch (err) {
				console.error('Error fetching tags:', err);
			}
		}
	}
	
	function addTag(tag: string) {
		const normalized = tag.toLowerCase().trim();
		if (normalized && !tags.includes(normalized) && tags.length < maxTags) {
			tags = [...tags, normalized];
			inputValue = '';
			showSuggestions = false;
		}
	}
	
	function removeTag(index: number) {
		tags = tags.filter((_, i) => i !== index);
	}
	
	function handleInput(event: KeyboardEvent) {
		const target = event.target as HTMLInputElement;
		inputValue = target.value;
		
		if (event.key === 'Enter' && inputValue.trim()) {
			event.preventDefault();
			addTag(inputValue);
		} else if (event.key === 'Backspace' && !inputValue && tags.length > 0) {
			removeTag(tags.length - 1);
		} else if (event.key === 'ArrowDown') {
			event.preventDefault();
			selectedIndex = Math.min(selectedIndex + 1, filteredSuggestions.length - 1);
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			selectedIndex = Math.max(selectedIndex - 1, -1);
		} else if (event.key === 'Escape') {
			showSuggestions = false;
			selectedIndex = -1;
		} else if (event.key === 'Enter' && selectedIndex >= 0) {
			event.preventDefault();
			addTag(filteredSuggestions[selectedIndex]);
			selectedIndex = -1;
		}
	}
	
	function handleSuggestionClick(tag: string) {
		addTag(tag);
	}
	
	function handleBlur() {
		// Delay to allow click events on suggestions
		setTimeout(() => {
			showSuggestions = false;
			selectedIndex = -1;
		}, 200);
	}
	
	function handleFocus() {
		if (inputValue.trim()) {
			showSuggestions = filteredSuggestions.length > 0;
		}
		if (availableTags.length === 0) {
			loadAvailableTags();
		}
	}
</script>

<div class="tag-input-container">
	<div class="tags-wrapper">
		{#each tags as tag, index}
			<span class="tag-chip">
				{tag}
				<button
					type="button"
					class="tag-remove"
					onclick={() => removeTag(index)}
					aria-label="Remove tag"
				>
					<X size={14} />
				</button>
			</span>
		{/each}
		<input
			type="text"
			class="tag-input"
			bind:value={inputValue}
			onkeydown={handleInput}
			onfocus={handleFocus}
			onblur={handleBlur}
			{placeholder}
			disabled={tags.length >= maxTags}
		/>
	</div>
	
	{#if showSuggestions}
		<div class="suggestions">
			{#each filteredSuggestions as suggestion, index}
				<button
					type="button"
					class="suggestion-item"
					class:selected={index === selectedIndex}
					onclick={() => handleSuggestionClick(suggestion)}
					onmousedown={(e) => e.preventDefault()}
				>
					{suggestion}
				</button>
			{/each}
		</div>
	{/if}
	
	{#if tags.length >= maxTags}
		<small class="max-tags-hint">Maximum {maxTags} tags allowed</small>
	{/if}
</div>

<style>
	.tag-input-container {
		position: relative;
		width: 100%;
	}
	
	.tags-wrapper {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
		align-items: center;
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		min-height: 42px;
	}
	
	.tags-wrapper:focus-within {
		border-color: var(--color-primary);
		outline: none;
	}
	
	.tag-chip {
		display: inline-flex;
		align-items: center;
		gap: var(--spacing-xs);
		padding: 4px 8px;
		background-color: var(--color-primary);
		color: white;
		border-radius: var(--radius-sm);
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
	}
	
	.tag-remove {
		display: flex;
		align-items: center;
		justify-content: center;
		background: none;
		border: none;
		color: white;
		cursor: pointer;
		padding: 0;
		margin: 0;
		opacity: 0.8;
		transition: opacity var(--transition-base);
	}
	
	.tag-remove:hover {
		opacity: 1;
	}
	
	.tag-input {
		flex: 1;
		min-width: 120px;
		border: none;
		outline: none;
		background: transparent;
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		padding: 4px;
	}
	
	.tag-input:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	
	.suggestions {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		margin-top: var(--spacing-xs);
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		max-height: 200px;
		overflow-y: auto;
		z-index: 1000;
	}
	
	.suggestion-item {
		width: 100%;
		padding: var(--spacing-sm) var(--spacing-md);
		text-align: left;
		background: none;
		border: none;
		cursor: pointer;
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		transition: background-color var(--transition-base);
	}
	
	.suggestion-item:hover,
	.suggestion-item.selected {
		background-color: var(--color-bg-secondary);
	}
	
	.suggestion-item:first-child {
		border-top-left-radius: var(--radius-md);
		border-top-right-radius: var(--radius-md);
	}
	
	.suggestion-item:last-child {
		border-bottom-left-radius: var(--radius-md);
		border-bottom-right-radius: var(--radius-md);
	}
	
	.max-tags-hint {
		display: block;
		margin-top: var(--spacing-xs);
		color: var(--color-text-secondary);
		font-size: var(--font-size-xs);
	}
</style>

