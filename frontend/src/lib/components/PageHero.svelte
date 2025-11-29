<script lang="ts">
	export let title: string;
	export let description: string | undefined = undefined;
	export let gradientFrom = 'from-indigo-950/30';
	export let gradientVia = 'via-purple-950/30';
	export let gradientTo = 'to-indigo-950/30';
</script>

<div
	class="relative overflow-hidden border-b border-slate-800/50 bg-gradient-to-r {gradientFrom} {gradientVia} {gradientTo} backdrop-blur-sm"
>
	<div class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(99,102,241,0.1),transparent_50%)]"></div>
	<div class="relative px-6 py-8 lg:px-8">
		<div class="mx-auto max-w-7xl">
			<div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
				<div class="space-y-2">
					<h1 class="text-3xl font-bold tracking-tight text-white lg:text-4xl">{title}</h1>
					{#if description}
						<p class="text-sm text-slate-400 lg:text-base">{description}</p>
					{/if}
				</div>
				<div class="flex items-center gap-3">
					<slot name="actions" />
				</div>
			</div>
		</div>
	</div>
</div>

