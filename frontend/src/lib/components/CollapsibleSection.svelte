<script lang="ts">
	import { fly } from 'svelte/transition';

	export let open = false;
</script>

{#if open}
	<div class="mb-8" transition:fly={{ y: -20, duration: 300 }}>
		<div class="rounded-2xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm">
			<slot />
		</div>
	</div>
{/if}

