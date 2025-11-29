<script lang="ts">
	export let title = 'No items yet';
	export let description: string | undefined = undefined;
	export let icon: string | undefined = undefined;
</script>

<div class="rounded-xl border border-dashed border-slate-700/50 bg-slate-900/20 p-12 text-center">
	{#if icon}
		<div class="mx-auto mb-4 h-12 w-12 text-slate-600">{@html icon}</div>
	{/if}
	<p class="text-sm font-medium text-slate-400">{title}</p>
	{#if description}
		<p class="mt-2 text-xs text-slate-500">{description}</p>
	{/if}
	<slot />
</div>

