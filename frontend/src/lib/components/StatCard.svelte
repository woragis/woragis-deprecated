<script lang="ts">
	export let label: string;
	export let value: string | number;
	export let accentColor = 'indigo';
	export let icon: string | undefined = undefined;

	const colorClasses: Record<string, string> = {
		indigo: 'hover:border-indigo-500/30 hover:shadow-indigo-500/10 bg-indigo-500/10',
		emerald: 'hover:border-emerald-500/30 hover:shadow-emerald-500/10 bg-emerald-500/10',
		purple: 'hover:border-purple-500/30 hover:shadow-purple-500/10 bg-purple-500/10',
		blue: 'hover:border-blue-500/30 hover:shadow-blue-500/10 bg-blue-500/10',
		amber: 'hover:border-amber-500/30 hover:shadow-amber-500/10 bg-amber-500/10',
		red: 'hover:border-red-500/30 hover:shadow-red-500/10 bg-red-500/10'
	};
</script>

<div
	class="group relative overflow-hidden rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-5 backdrop-blur-sm transition-all {colorClasses[accentColor]} hover:shadow-lg"
>
	<div class="relative z-10">
		<p class="text-xs font-medium uppercase tracking-wider text-slate-400">{label}</p>
		<p class="mt-2 text-3xl font-bold text-white">{value}</p>
	</div>
	<div class="absolute -right-4 -top-4 h-24 w-24 rounded-full {colorClasses[accentColor]} blur-2xl"></div>
	{#if icon}
		<div class="absolute bottom-4 right-4 opacity-10">
			{@html icon}
		</div>
	{/if}
</div>

