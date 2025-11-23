<script lang="ts">
	import { onMount } from 'svelte';
	import mermaid from 'mermaid';

	interface Props {
		diagram: string;
		id?: string;
		theme?: 'default' | 'dark' | 'forest' | 'neutral';
	}

	let { diagram, id = `mermaid-${Math.random().toString(36).substr(2, 9)}`, theme = 'dark' }: Props = $props();

	let container: HTMLDivElement | null = $state(null);
	let error: string | null = $state(null);
	let loading = $state(true);

	onMount(async () => {
		if (!container || !diagram) return;

		try {
			// Initialize Mermaid with dark theme
			mermaid.initialize({
				startOnLoad: false,
				theme: theme === 'dark' ? 'dark' : theme,
				themeVariables: {
					primaryColor: '#3b82f6',
					primaryTextColor: '#ffffff',
					primaryBorderColor: '#1e40af',
					lineColor: '#60a5fa',
					secondaryColor: '#1e3a8a',
					tertiaryColor: '#1e40af',
					background: '#1f2937',
					mainBkgColor: '#1f2937',
					secondBkgColor: '#374151',
					textColor: '#ffffff',
					textInverseColor: '#000000',
					border1: '#4b5563',
					border2: '#6b7280',
					noteBkgColor: '#1e3a8a',
					noteTextColor: '#ffffff',
					noteBorderColor: '#3b82f6',
					actorBorder: '#3b82f6',
					actorBkg: '#1e3a8a',
					actorTextColor: '#ffffff',
					actorLineColor: '#60a5fa',
					signalColor: '#60a5fa',
					signalTextColor: '#ffffff',
					labelBoxBkgColor: '#1e3a8a',
					labelBoxBorderColor: '#3b82f6',
					labelTextColor: '#ffffff',
					loopTextColor: '#ffffff',
					activationBorderColor: '#3b82f6',
					activationBkgColor: '#1e40af',
					sequenceNumberColor: '#ffffff',
					sectionBkgColor: '#1e3a8a',
					altBkgColor: '#1e40af',
					optBkgColor: '#1e40af',
					sequenceLineColor: '#60a5fa',
					sectionBorderColor: '#3b82f6',
					altBorderColor: '#3b82f6',
					optBorderColor: '#3b82f6'
				},
				securityLevel: 'loose',
				flowchart: {
					useMaxWidth: true,
					htmlLabels: true,
					curve: 'basis'
				}
			});

			// Render the diagram
			const { svg } = await mermaid.render(id, diagram);
			if (container) {
				container.innerHTML = svg;
			}
			loading = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to render diagram';
			loading = false;
			console.error('Mermaid rendering error:', err);
		}
	});
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-8">
			<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if error}
		<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-4 text-red-300">
			<p class="text-sm">Failed to render diagram: {error}</p>
			<details class="mt-2">
				<summary class="cursor-pointer text-xs text-red-400">Show diagram code</summary>
				<pre class="mt-2 text-xs bg-gray-900 p-2 rounded overflow-auto"><code>{diagram}</code></pre>
			</details>
		</div>
	{:else}
		<div bind:this={container} class="w-full overflow-auto bg-gray-900/50 rounded-lg p-4 border border-gray-700"></div>
	{/if}
</div>

