<script lang="ts">
	import PageHero from '$lib/components/PageHero.svelte';

	export let onRefresh: () => void;
	export let onNewIdea: () => void;
	export let onNewNode: () => void;
	export let showNewNode = false;
</script>

<PageHero
	title="Ideas Canvas"
	description="Visualize and connect your ideas with nodes and connections"
	gradientFrom="from-violet-950/30"
	gradientVia="via-purple-950/30"
	gradientTo="to-violet-950/30"
>
	<svelte:fragment slot="actions">
		<button
			class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2.5 text-sm font-medium text-slate-200 backdrop-blur-sm transition-all hover:border-violet-500/50 hover:bg-slate-800/80"
			onclick={onRefresh}
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
		{#if showNewNode}
			<button
				class="rounded-lg border border-violet-500/50 bg-violet-600/20 px-4 py-2.5 text-sm font-medium text-violet-300 backdrop-blur-sm transition-all hover:border-violet-500 hover:bg-violet-600/30"
				onclick={onNewNode}
			>
				New Node
			</button>
		{/if}
		<button
			class="group relative overflow-hidden rounded-lg bg-gradient-to-r from-violet-600 to-purple-600 px-6 py-2.5 text-sm font-semibold text-white shadow-lg shadow-violet-500/25 transition-all hover:scale-105 hover:shadow-xl hover:shadow-violet-500/40"
			onclick={onNewIdea}
		>
			<span class="relative z-10 flex items-center gap-2">
				<svg
					class="h-4 w-4 transition-transform group-hover:rotate-90"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					></path>
				</svg>
				New Idea
			</span>
		</button>
	</svelte:fragment>
</PageHero>

