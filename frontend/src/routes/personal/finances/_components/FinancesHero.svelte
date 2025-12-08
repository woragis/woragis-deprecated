<script lang="ts">
	import PageHero from '$lib/components/PageHero.svelte';

	export let isAuthenticated = false;
	export let onRefresh: () => void;
</script>

<PageHero
	title="Finance Workspace"
	description="Track your income, expenses, and financial goals"
	gradientFrom="from-emerald-950/30"
	gradientVia="via-teal-950/30"
	gradientTo="to-emerald-950/30"
>
	<svelte:fragment slot="actions">
		{#if isAuthenticated}
			<button
				class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2.5 text-sm font-medium text-slate-200 backdrop-blur-sm transition-all hover:border-emerald-500/50 hover:bg-slate-800/80"
				on:click={onRefresh}
			>
				<svg class="mr-2 inline h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
					></path>
				</svg>
				Refresh
			</button>
		{/if}
	</svelte:fragment>
</PageHero>

