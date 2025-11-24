<script lang="ts">
	import { Brain, ExternalLink, Code, Zap, Target, Layers, Sparkles, Cpu, MessageSquare, TrendingUp, Eye, FileText, Github } from 'lucide-svelte';
	import type { AIMLIntegration, IntegrationType, Framework } from '$lib/types/aiml-integration';

	interface Props {
		integrations: AIMLIntegration[];
		loading?: boolean;
	}

	let { integrations = [], loading = false }: Props = $props();

	// Group integrations by type
	let groupedIntegrations = $derived.by(() => {
		const grouped: Record<IntegrationType, AIMLIntegration[]> = {
			rag: [],
			llm: [],
			ml_model: [],
			computer_vision: [],
			nlp: [],
			recommendation: [],
			chatbot: [],
			anomaly_detection: [],
			predictive_analytics: [],
			generative_ai: [],
			other: []
		};

		integrations.forEach((integration) => {
			if (grouped[integration.type]) {
				grouped[integration.type].push(integration);
			} else {
				grouped.other.push(integration);
			}
		});

		// Sort each group by display order
		Object.keys(grouped).forEach((key) => {
			grouped[key as IntegrationType].sort((a, b) => a.displayOrder - b.displayOrder);
		});

		// Filter out empty groups
		return Object.entries(grouped).filter(([_, items]) => items.length > 0);
	});

	function getTypeIcon(type: IntegrationType) {
		const icons: Record<IntegrationType, typeof Brain> = {
			rag: Layers,
			llm: Brain,
			ml_model: Cpu,
			computer_vision: Eye,
			nlp: MessageSquare,
			recommendation: TrendingUp,
			chatbot: MessageSquare,
			anomaly_detection: Zap,
			predictive_analytics: Target,
			generative_ai: Sparkles,
			other: Code
		};
		return icons[type] || Code;
	}

	function getTypeColor(type: IntegrationType): string {
		const colors: Record<IntegrationType, string> = {
			rag: 'from-purple-500 to-pink-500',
			llm: 'from-blue-500 to-cyan-500',
			ml_model: 'from-green-500 to-emerald-500',
			computer_vision: 'from-orange-500 to-red-500',
			nlp: 'from-indigo-500 to-purple-500',
			recommendation: 'from-yellow-500 to-amber-500',
			chatbot: 'from-teal-500 to-cyan-500',
			anomaly_detection: 'from-red-500 to-pink-500',
			predictive_analytics: 'from-violet-500 to-purple-500',
			generative_ai: 'from-pink-500 via-purple-500 to-blue-500',
			other: 'from-gray-500 to-slate-500'
		};
		return colors[type] || colors.other;
	}

	function getTypeLabel(type: IntegrationType): string {
		const labels: Record<IntegrationType, string> = {
			rag: 'RAG Systems',
			llm: 'Large Language Models',
			ml_model: 'ML Models',
			computer_vision: 'Computer Vision',
			nlp: 'Natural Language Processing',
			recommendation: 'Recommendation Systems',
			chatbot: 'Chatbots & Assistants',
			anomaly_detection: 'Anomaly Detection',
			predictive_analytics: 'Predictive Analytics',
			generative_ai: 'Generative AI',
			other: 'Other'
		};
		return labels[type] || 'Other';
	}

	function getFrameworkLabel(framework: Framework): string {
		const labels: Record<Framework, string> = {
			openai: 'OpenAI',
			anthropic: 'Anthropic',
			huggingface: 'Hugging Face',
			tensorflow: 'TensorFlow',
			pytorch: 'PyTorch',
			langchain: 'LangChain',
			llamaindex: 'LlamaIndex',
			cohere: 'Cohere',
			google_ai: 'Google AI',
			azure_ai: 'Azure AI',
			aws_bedrock: 'AWS Bedrock',
			custom: 'Custom',
			other: 'Other'
		};
		return labels[framework] || framework;
	}

	function truncateText(text: string, maxLength: number): string {
		if (text.length <= maxLength) return text;
		return text.slice(0, maxLength).trim() + '...';
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if integrations.length === 0}
		<div class="text-center py-20">
			<Brain class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No AI/ML integrations available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="space-y-8">
			{#each groupedIntegrations as [type, items]}
				{@const TypeIcon = getTypeIcon(type as IntegrationType)}
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<!-- Type Header -->
					<div class="flex items-center gap-3 mb-6">
						<div
							class="w-12 h-12 bg-gradient-to-br {getTypeColor(type as IntegrationType)} rounded-lg flex items-center justify-center"
						>
							<TypeIcon class="w-6 h-6 text-white" />
						</div>
						<div>
							<h3 class="text-2xl font-bold text-white">{getTypeLabel(type as IntegrationType)}</h3>
							<p class="text-sm text-gray-400">{items.length} {items.length === 1 ? 'integration' : 'integrations'}</p>
						</div>
					</div>

					<!-- Integrations Grid -->
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each items as integration}
							<div
								class="group bg-gray-800/50 rounded-lg p-5 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 flex flex-col"
							>
								<!-- Header -->
								<div class="mb-3">
									<h4 class="text-lg font-bold text-white mb-2 group-hover:text-blue-400 transition-colors">
										{integration.title}
									</h4>
									<div class="flex items-center gap-2 mb-2">
										<span
											class="px-2 py-1 text-xs rounded bg-blue-600/20 text-blue-300 border border-blue-500/30"
										>
											{getFrameworkLabel(integration.framework)}
										</span>
										{#if integration.modelName}
											<span
												class="px-2 py-1 text-xs rounded bg-purple-600/20 text-purple-300 border border-purple-500/30"
											>
												{integration.modelName}
												{#if integration.modelVersion}
													<span class="text-purple-400/70"> v{integration.modelVersion}</span>
												{/if}
											</span>
										{/if}
									</div>
								</div>

								<!-- Description -->
								<p class="text-sm text-gray-300 mb-4 line-clamp-3 flex-1">
									{integration.description}
								</p>

								<!-- Use Case -->
								{#if integration.useCase}
									<div class="mb-4">
										<div class="flex items-start gap-2 mb-2">
											<Target class="w-4 h-4 text-cyan-400 mt-0.5 flex-shrink-0" />
											<div>
												<p class="text-xs font-semibold text-cyan-400 mb-1">Use Case</p>
												<p class="text-xs text-gray-300">{truncateText(integration.useCase, 120)}</p>
											</div>
										</div>
									</div>
								{/if}

								<!-- Impact -->
								{#if integration.impact}
									<div class="mb-4">
										<div class="flex items-start gap-2 mb-2">
											<Zap class="w-4 h-4 text-yellow-400 mt-0.5 flex-shrink-0" />
											<div>
												<p class="text-xs font-semibold text-yellow-400 mb-1">Impact</p>
												<p class="text-xs text-gray-300">{truncateText(integration.impact, 120)}</p>
											</div>
										</div>
									</div>
								{/if}

								<!-- Technologies -->
								{#if integration.technologies && integration.technologies.length > 0}
									<div class="mb-4">
										<div class="flex flex-wrap gap-2">
											{#each integration.technologies.slice(0, 4) as tech}
												<span
													class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
												>
													{tech}
												</span>
											{/each}
											{#if integration.technologies.length > 4}
												<span
													class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
												>
													+{integration.technologies.length - 4}
												</span>
											{/if}
										</div>
									</div>
								{/if}

								<!-- Links -->
								<div class="flex items-center gap-3 pt-4 border-t border-gray-700 mt-auto">
									{#if integration.demoUrl}
										<a
											href={integration.demoUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center gap-1 text-sm text-blue-400 hover:text-blue-300 transition-colors"
											on:click|stopPropagation
										>
											<Eye class="w-4 h-4" />
											<span>Demo</span>
										</a>
									{/if}
									{#if integration.documentationUrl}
										<a
											href={integration.documentationUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center gap-1 text-sm text-gray-400 hover:text-gray-300 transition-colors"
											on:click|stopPropagation
										>
											<FileText class="w-4 h-4" />
											<span>Docs</span>
										</a>
									{/if}
									{#if integration.githubUrl}
										<a
											href={integration.githubUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center gap-1 text-sm text-gray-400 hover:text-gray-300 transition-colors"
											on:click|stopPropagation
										>
											<Github class="w-4 h-4" />
											<span>Code</span>
										</a>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

