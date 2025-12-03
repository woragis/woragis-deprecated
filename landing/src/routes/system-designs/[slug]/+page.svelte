<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Layers, ExternalLink, Calendar, ArrowLeft, Server, Zap, Star } from 'lucide-svelte';
	import { getSystemDesign } from '$lib/api/system-designs';
	import type { SystemDesign } from '$lib/types/system-design';

	let design: SystemDesign | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const designId = $derived($page.params.slug);

	onMount(async () => {
		if (designId) {
			await fetchDesign(designId);
		}
	});

	async function fetchDesign(id: string) {
		loading = true;
		error = null;
		try {
			design = await getSystemDesign(id);
			if (!design) {
				error = 'System design not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch system design';
			console.error('Error fetching system design:', err);
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !design}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">System Design Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The system design you are looking for does not exist.'}</p>
				<a
					href="/system-designs"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to System Designs
				</a>
			</div>
		</div>
	{:else}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/system-designs"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to System Designs
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 hover:from-blue-500/5 hover:via-purple-500/5 hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<div class="flex items-center gap-3 mb-4">
								<div
									class="w-12 h-12 bg-gradient-to-br from-blue-600 to-cyan-600 rounded-lg flex items-center justify-center"
								>
									<Layers class="w-6 h-6 text-white" />
								</div>
								{#if design.featured}
									<span class="px-3 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded">
										⭐ Featured
									</span>
								{/if}
							</div>

							<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
								{design.title}
							</h1>

							{#if design.description}
								<p class="text-xl text-gray-300 mb-6 leading-relaxed">{design.description}</p>
							{/if}

							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if design.updatedAt}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Updated {formatDate(design.updatedAt)}</span>
									</div>
								{/if}
							</div>
						</div>

						<!-- Components -->
						{#if design.components && design.components.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Components</h2>
								<div class="grid md:grid-cols-2 gap-4">
									{#each design.components as component}
										<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
											<div class="flex items-center gap-3 mb-3">
												<Server class="w-5 h-5 text-blue-400" />
												<h3 class="text-lg font-bold text-white">{component.name}</h3>
											</div>
											<p class="text-gray-300 text-sm mb-3">{component.description}</p>
											<div class="flex items-center gap-2">
												<span class="px-2 py-1 bg-blue-600/20 text-blue-300 text-xs rounded-lg border border-blue-500/30">
													{component.technology}
												</span>
											</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Data Flow -->
						{#if design.dataFlow}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Data Flow</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed whitespace-pre-wrap">{design.dataFlow}</p>
								</div>
							</div>
						{/if}

						<!-- Scalability & Reliability -->
						<div class="grid md:grid-cols-2 gap-6 mb-8">
							{#if design.scalability}
								<div class="bg-green-900/20 rounded-lg p-6 border border-green-700/30">
									<div class="flex items-center gap-3 mb-3">
										<Zap class="w-5 h-5 text-green-400" />
										<h3 class="text-lg font-bold text-white">Scalability</h3>
									</div>
									<p class="text-gray-300 leading-relaxed">{design.scalability}</p>
								</div>
							{/if}

							{#if design.reliability}
								<div class="bg-blue-900/20 rounded-lg p-6 border border-blue-700/30">
									<div class="flex items-center gap-3 mb-3">
										<Server class="w-5 h-5 text-blue-400" />
										<h3 class="text-lg font-bold text-white">Reliability</h3>
									</div>
									<p class="text-gray-300 leading-relaxed">{design.reliability}</p>
								</div>
							{/if}
						</div>

						<!-- Diagram -->
						{#if design.diagram}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Architecture Diagram</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<pre class="text-gray-300 text-sm overflow-x-auto whitespace-pre-wrap">{design.diagram}</pre>
								</div>
							</div>
						{/if}
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
