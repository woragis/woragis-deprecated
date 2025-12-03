<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Brain, ExternalLink, Calendar, ArrowLeft, Code2, Globe, Github, Star, TrendingUp } from 'lucide-svelte';
	import { getAIMLIntegration } from '$lib/api/aiml-integrations';
	import type { AIMLIntegration, IntegrationType, Framework } from '$lib/types/aiml-integration';

	let integration: AIMLIntegration | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const integrationId = $derived($page.params.slug);

	onMount(async () => {
		if (integrationId) {
			await fetchIntegration(integrationId);
		}
	});

	async function fetchIntegration(id: string) {
		loading = true;
		error = null;
		try {
			integration = await getAIMLIntegration(id);
			if (!integration) {
				error = 'AI/ML integration not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch AI/ML integration';
			console.error('Error fetching AI/ML integration:', err);
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function getTypeLabel(type: IntegrationType): string {
		return type.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
	}

	function getFrameworkLabel(framework: Framework): string {
		return framework.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !integration}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">AI/ML Integration Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The AI/ML integration you are looking for does not exist.'}</p>
				<a
					href="/aiml-integrations"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to AI/ML Integrations
				</a>
			</div>
		</div>
	{:else}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/aiml-integrations"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to AI/ML Integrations
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-cyan-500/0 via-purple-500/0 to-pink-500/0 hover:from-cyan-500/5 hover:via-purple-500/5 hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<div class="flex items-center gap-3 mb-4">
								<div
									class="w-12 h-12 bg-gradient-to-br from-cyan-600 to-purple-600 rounded-lg flex items-center justify-center"
								>
									<Brain class="w-6 h-6 text-white" />
								</div>
								<div class="flex items-center gap-2">
									<span
										class="px-3 py-1 bg-cyan-600/20 text-cyan-300 text-sm rounded-lg border border-cyan-500/30"
									>
										{getTypeLabel(integration.type)}
									</span>
									<span
										class="px-3 py-1 bg-purple-600/20 text-purple-300 text-sm rounded-lg border border-purple-500/30"
									>
										{getFrameworkLabel(integration.framework)}
									</span>
									{#if integration.featured}
										<span class="px-3 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded">
											⭐ Featured
										</span>
									{/if}
								</div>
							</div>

							<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
								{integration.title}
							</h1>

							{#if integration.description}
								<p class="text-xl text-gray-300 mb-6 leading-relaxed">{integration.description}</p>
							{/if}

							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if integration.updatedAt}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Updated {formatDate(integration.updatedAt)}</span>
									</div>
								{/if}
							</div>
						</div>

						<!-- Model Info -->
						{#if integration.modelName}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Model Information</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<div class="grid md:grid-cols-2 gap-4">
										<div>
											<p class="text-sm text-gray-400 mb-1">Model Name</p>
											<p class="text-lg font-bold text-white">{integration.modelName}</p>
										</div>
										{#if integration.modelVersion}
											<div>
												<p class="text-sm text-gray-400 mb-1">Version</p>
												<p class="text-lg font-bold text-white">{integration.modelVersion}</p>
											</div>
										{/if}
									</div>
								</div>
							</div>
						{/if}

						<!-- Use Case -->
						{#if integration.useCase}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Use Case</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed">{integration.useCase}</p>
								</div>
							</div>
						{/if}

						<!-- Architecture -->
						{#if integration.architecture}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Architecture</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed whitespace-pre-wrap">{integration.architecture}</p>
								</div>
							</div>
						{/if}

						<!-- Impact -->
						{#if integration.impact}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Impact</h2>
								<div class="bg-blue-900/20 rounded-lg p-6 border border-blue-700/30">
									<div class="flex items-start gap-3">
										<TrendingUp class="w-6 h-6 text-blue-400 flex-shrink-0 mt-0.5" />
										<p class="text-gray-300 leading-relaxed">{integration.impact}</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Metrics -->
						{#if integration.metrics}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Metrics</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed whitespace-pre-wrap">{integration.metrics}</p>
								</div>
							</div>
						{/if}

						<!-- Technologies -->
						{#if integration.technologies && integration.technologies.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Technologies</h2>
								<div class="flex flex-wrap gap-2">
									{#each integration.technologies as tech}
										<span
											class="px-3 py-1 bg-cyan-600/20 text-cyan-300 text-sm rounded-lg border border-cyan-500/30"
										>
											{tech}
										</span>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Links -->
						<div class="mt-8 pt-8 border-t border-gray-700 flex flex-wrap gap-4">
							{#if integration.demoUrl}
								<a
									href={integration.demoUrl}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-green-600 hover:bg-green-700 rounded-lg transition-colors duration-200 font-medium"
								>
									<Globe class="w-5 h-5" />
									View Demo
									<ExternalLink class="w-5 h-5" />
								</a>
							{/if}
							{#if integration.githubUrl}
								<a
									href={integration.githubUrl}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors duration-200 font-medium"
								>
									<Github class="w-5 h-5" />
									View Code
									<ExternalLink class="w-5 h-5" />
								</a>
							{/if}
							{#if integration.documentationUrl}
								<a
									href={integration.documentationUrl}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200 font-medium"
								>
									<Code2 class="w-5 h-5" />
									Documentation
									<ExternalLink class="w-5 h-5" />
								</a>
							{/if}
						</div>
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
