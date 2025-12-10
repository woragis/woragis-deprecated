<script lang="ts">
	let {
		tabs = [],
		activeTab = $bindable(''),
		onTabChange
	}: {
		tabs?: Array<{ id: string; label: string; icon?: any }>;
		activeTab?: string;
		onTabChange?: (tabId: string) => void;
	} = $props();
	
	function handleTabClick(tabId: string) {
		activeTab = tabId;
		onTabChange?.(tabId);
	}
</script>

<div class="tabs">
	{#each tabs as tab}
		<button
			class="tab"
			class:active={activeTab === tab.id}
			onclick={() => handleTabClick(tab.id)}
			type="button"
		>
			{#if tab.icon}
				<svelte:component this={tab.icon} class="tab-icon" />
			{/if}
			{tab.label}
		</button>
	{/each}
</div>

<style>
	.tabs {
		display: flex;
		gap: var(--spacing-sm);
		margin-bottom: var(--spacing-lg);
		border-bottom: 2px solid var(--color-border);
	}
	
	.tab {
		padding: var(--spacing-md) var(--spacing-lg);
		background: none;
		border: none;
		border-bottom: 2px solid transparent;
		cursor: pointer;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		transition: color var(--transition-base), border-color var(--transition-base);
		margin-bottom: -2px;
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		font-family: inherit;
	}
	
	.tab:hover {
		color: var(--color-text-primary);
	}
	
	.tab.active {
		color: var(--color-primary);
		border-bottom-color: var(--color-primary);
		font-weight: var(--font-weight-semibold);
	}
	
	.tab-icon {
		width: 1rem;
		height: 1rem;
	}
</style>
