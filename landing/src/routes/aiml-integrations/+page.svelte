<script lang="ts">
	import { Search, Filter, X, ChevronDown, ExternalLink, Calendar, Brain, Star, Code2, Globe, Github } from 'lucide-svelte';
	import type { AIMLIntegration, IntegrationType, Framework } from '$lib/types/aiml-integration';
	import { useAIMLIntegrationsQuery } from '$lib/queries/aiml-integrations';

	// Filters
	let searchQuery = $state('');
	let typeFilter: IntegrationType | 'all' = $state('all');
	let frameworkFilter: Framework | 'all' = $state('all');
	let featuredFilter: boolean | 'all' = $state('all');
	let sortBy: 'updatedAt' | 'title' | 'createdAt' = $state('updatedAt');
	let sortOrder: 'asc' | 'desc' = $state('desc');
	let showFilters = $state(false);

	// Fetch all AI/ML integrations
	const integrationsQuery = useAIMLIntegrationsQuery();
	let integrations = $derived(integrationsQuery.data || []);
	let filteredIntegrations: AIMLIntegration[] = $state([]);
	let loading = $derived(integrationsQuery.isPending);
	let error = $derived(integrationsQuery.error ? (integrationsQuery.error instanceof Error ? integrationsQuery.error.message : 'Failed to fetch AI/ML integrations') : null);

	const typeOptions: Array<{ value: IntegrationType | 'all'; label: string }> = [
		{ value: 'all', label: 'All Types' },
		{ value: 'rag', label: 'RAG' },
		{ value: 'llm', label: 'LLM' },
		{ value: 'ml_model', label: 'ML Model' },
		{ value: 'computer_vision', label: 'Computer Vision' },
		{ value: 'nlp', label: 'NLP' },
		{ value: 'recommendation', label: 'Recommendation' },
		{ value: 'chatbot', label: 'Chatbot' },
		{ value: 'anomaly_detection', label: 'Anomaly Detection' },
		{ value: 'predictive_analytics', label: 'Predictive Analytics' },
		{ value: 'generative_ai', label: 'Generative AI' },
		{ value: 'other', label: 'Other' }
	];

	const frameworkOptions: Array<{ value: Framework | 'all'; label: string }> = [
		{ value: 'all', label: 'All Frameworks' },
		{ value: 'openai', label: 'OpenAI' },
		{ value: 'anthropic', label: 'Anthropic' },
		{ value: 'huggingface', label: 'Hugging Face' },
		{ value: 'tensorflow', label: 'TensorFlow' },
		{ value: 'pytorch', label: 'PyTorch' },
		{ value: 'langchain', label: 'LangChain' },
		{ value: 'llamaindex', label: 'LlamaIndex' },
		{ value: 'cohere', label: 'Cohere' },
		{ value: 'google_ai', label: 'Google AI' },
		{ value: 'azure_ai', label: 'Azure AI' },
		{ value: 'aws_bedrock', label: 'AWS Bedrock' },
		{ value: 'custom', label: 'Custom' },
		{ value: 'other', label: 'Other' }
	];

	const sortOptions = [
		{ value: 'updatedAt', label: 'Last Updated' },
		{ value: 'createdAt', label: 'Date Created' },
		{ value: 'title', label: 'Title' }
	];

	// Apply filters when integrations or filter values change
	$effect(() => {
		applyFilters();
	});

	function applyFilters() {
		let filtered = [...integrations];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(integration: AIMLIntegration) =>
					integration.title.toLowerCase().includes(query) ||
					integration.description?.toLowerCase().includes(query) ||
					integration.useCase?.toLowerCase().includes(query) ||
					integration.modelName?.toLowerCase().includes(query) ||
					integration.technologies?.some((tech) => tech.toLowerCase().includes(query))
			);
		}

		// Apply type filter
		if (typeFilter !== 'all') {
			filtered = filtered.filter((integration: AIMLIntegration) => integration.type === typeFilter);
		}

		// Apply framework filter
		if (frameworkFilter !== 'all') {
			filtered = filtered.filter((integration: AIMLIntegration) => integration.framework === frameworkFilter);
		}

		// Apply featured filter
		if (featuredFilter !== 'all') {
			filtered = filtered.filter((integration: AIMLIntegration) => integration.featured === featuredFilter);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aVal: string | number | undefined;
			let bVal: string | number | undefined;

			switch (sortBy) {
				case 'updatedAt':
					aVal = a.updatedAt ? new Date(a.updatedAt).getTime() : 0;
					bVal = b.updatedAt ? new Date(b.updatedAt).getTime() : 0;
					break;
				case 'createdAt':
					aVal = a.createdAt ? new Date(a.createdAt).getTime() : 0;
					bVal = b.createdAt ? new Date(b.createdAt).getTime() : 0;
					break;
				case 'title':
					aVal = a.title.toLowerCase();
					bVal = b.title.toLowerCase();
					break;
				default:
					return 0;
			}

			if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
			if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
			return 0;
		});

		filteredIntegrations = filtered;
	}

	function clearFilters() {
		searchQuery = '';
		typeFilter = 'all';
		frameworkFilter = 'all';
		featuredFilter = 'all';
		applyFilters();
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getIntegrationUrl(integration: AIMLIntegration): string {
		return `/aiml-integrations/${integration.id}`;
	}

	function getTypeLabel(type: IntegrationType): string {
		return type.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
	}

	function getFrameworkLabel(framework: Framework): string {
		return framework.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Header -->
	<section class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
					AI/ML Integrations
				</h1>
				<p class="text-xl text-gray-300 max-w-2xl mx-auto">
					Artificial Intelligence and Machine Learning integrations showcasing innovative solutions and implementations
				</p>
			</div>

			<!-- Search and Filter Bar -->
			<div class="mb-8">
				<div class="flex flex-col md:flex-row gap-4 mb-4">
					<!-- Search Input -->
					<div class="flex-1 relative">
						<Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
						<input
							type="text"
							placeholder="Search integrations by title, description, model, use case..."
							bind:value={searchQuery}
							class="w-full pl-10 pr-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
						/>
					</div>
					<!-- Filter Toggle -->
					<div class="flex items-center justify-between gap-4">
						<button
							onclick={() => (showFilters = !showFilters)}
							class="flex items-center gap-2 px-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white hover:bg-gray-700/50 transition-colors"
						>
							<Filter class="w-5 h-5" />
							Filters
							{#if showFilters}
								<ChevronDown class="w-4 h-4 rotate-180 transition-transform" />
							{:else}
								<ChevronDown class="w-4 h-4 transition-transform" />
							{/if}
						</button>

						<div class="flex items-center gap-4">
							<!-- Sort By -->
							<select
								bind:value={sortBy}
								onchange={applyFilters}
								class="px-4 py-2 bg-gray-800/50 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white text-sm"
							>
								{#each sortOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
							</select>

							<!-- Sort Order -->
							<button
								onclick={() => {
									sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
									applyFilters();
								}}
								class="px-4 py-2 bg-gray-800/50 hover:bg-gray-700/50 border border-gray-700 rounded-lg transition-colors text-sm"
								title="Toggle sort order"
							>
								{sortOrder === 'asc' ? '↑' : '↓'}
							</button>

							{#if searchQuery || typeFilter !== 'all' || frameworkFilter !== 'all' || featuredFilter !== 'all'}
								<button
									onclick={clearFilters}
									class="flex items-center gap-2 px-4 py-2 bg-red-600/20 hover:bg-red-600/30 border border-red-700/30 rounded-lg transition-colors text-sm"
								>
									<X class="w-4 h-4" />
									Clear
								</button>
							{/if}
						</div>
					</div>
				</div>

				<!-- Filters Panel -->
				{#if showFilters}
					<div class="p-4 bg-gray-800/50 border border-gray-700 rounded-lg space-y-4">
						<!-- Type Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Type</div>
							<div class="flex flex-wrap gap-2">
								{#each typeOptions as option}
									<button
										onclick={() => {
											typeFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {typeFilter === option.value
											? 'bg-blue-600 border-blue-500 text-white'
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										{option.label}
									</button>
								{/each}
							</div>
						</div>

						<!-- Framework Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Framework</div>
							<div class="flex flex-wrap gap-2">
								{#each frameworkOptions as option}
									<button
										onclick={() => {
											frameworkFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {frameworkFilter === option.value
											? 'bg-purple-600 border-purple-500 text-white'
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										{option.label}
									</button>
								{/each}
							</div>
						</div>

						<!-- Featured Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Featured</div>
							<div class="flex flex-wrap gap-2">
								<button
									onclick={() => {
										featuredFilter = 'all';
										applyFilters();
									}}
									class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {featuredFilter === 'all'
										? 'bg-blue-600 border-blue-500 text-white'
										: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
								>
									All
								</button>
								<button
									onclick={() => {
										featuredFilter = true;
										applyFilters();
									}}
									class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {featuredFilter === true
										? 'bg-yellow-600 border-yellow-500 text-white'
										: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
								>
									Featured Only
								</button>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</section>

	<!-- Integrations Grid -->
	<section class="container mx-auto px-6 py-12">
		<div class="max-w-7xl mx-auto">
			{#if loading}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-6 text-center">
					<p class="text-red-400 mb-2">Error loading AI/ML integrations</p>
					<p class="text-gray-400 text-sm">{error}</p>
					<button
						onclick={() => integrationsQuery.refetch()}
						class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{:else if filteredIntegrations.length === 0}
				<div class="text-center py-20">
					<Brain class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No AI/ML integrations found</p>
					<p class="text-gray-500 text-sm">Try adjusting your filters or search query</p>
				</div>
			{:else}
				<div class="mb-4 text-gray-400 text-sm">
					Showing {filteredIntegrations.length} of {integrations.length} AI/ML integrations
				</div>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredIntegrations as integration, index}
						{@const animationDelay = index * 0.1}
						<a
							href={getIntegrationUrl(integration)}
							class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl overflow-hidden border border-gray-700 hover:border-cyan-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-cyan-500/20 hover:scale-[1.02] relative animate-fadeInUp"
							style="animation-delay: {animationDelay}s"
						>
							<!-- Decorative gradient overlay -->
							<div class="absolute inset-0 bg-gradient-to-br from-cyan-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-cyan-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
							<div class="relative z-10">
								<div class="p-6">
									<div class="flex items-start justify-between mb-3">
										<div
											class="w-12 h-12 bg-gradient-to-br from-cyan-600 to-purple-600 rounded-lg flex items-center justify-center flex-shrink-0"
										>
											<Brain class="w-6 h-6 text-white" />
										</div>
										{#if integration.featured}
											<div
												class="px-2 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded flex items-center gap-1"
											>
												<Star class="w-3 h-3 fill-current" />
												Featured
											</div>
										{/if}
									</div>

									<h3
										class="text-xl font-bold text-white mb-2 group-hover:text-cyan-400 transition-colors line-clamp-2"
									>
										{integration.title}
									</h3>

									{#if integration.description}
										<p class="text-gray-300 text-sm mb-4 line-clamp-3">{integration.description}</p>
									{/if}

									<div class="flex items-center gap-2 mb-4 flex-wrap">
										<span
											class="px-2 py-1 bg-cyan-600/20 text-cyan-300 text-xs rounded-lg border border-cyan-500/30"
										>
											{getTypeLabel(integration.type)}
										</span>
										<span
											class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded-lg border border-purple-500/30"
										>
											{getFrameworkLabel(integration.framework)}
										</span>
									</div>

									{#if integration.modelName}
										<div class="mb-4">
											<p class="text-xs text-gray-400 mb-1">Model</p>
											<p class="text-sm text-gray-300 font-medium">{integration.modelName}
												{#if integration.modelVersion}
													<span class="text-gray-500"> v{integration.modelVersion}</span>
												{/if}
											</p>
										</div>
									{/if}

									<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
										{#if integration.updatedAt}
											<div class="flex items-center gap-1">
												<Calendar class="w-3 h-3" />
												<span>Updated {formatDate(integration.updatedAt)}</span>
											</div>
										{/if}
									</div>

									<!-- Links -->
									<div class="flex items-center gap-2 mb-4 flex-wrap">
										{#if integration.demoUrl}
											<a
												href={integration.demoUrl}
												target="_blank"
												rel="noopener noreferrer"
												onclick={(e) => e.stopPropagation()}
												class="flex items-center gap-1 px-2 py-1 bg-green-600/20 text-green-300 text-xs rounded-lg border border-green-500/30 hover:bg-green-600/30 transition-colors"
											>
												<Globe class="w-3 h-3" />
												Demo
											</a>
										{/if}
										{#if integration.githubUrl}
											<a
												href={integration.githubUrl}
												target="_blank"
												rel="noopener noreferrer"
												onclick={(e) => e.stopPropagation()}
												class="flex items-center gap-1 px-2 py-1 bg-gray-700/50 text-gray-300 text-xs rounded-lg border border-gray-600 hover:bg-gray-600/50 transition-colors"
											>
												<Github class="w-3 h-3" />
												Code
											</a>
										{/if}
									</div>

									<div class="flex items-center justify-between pt-4 border-t border-gray-700">
										<div class="flex items-center gap-2 text-cyan-400 text-sm font-medium group-hover:gap-3 transition-all">
											<span>View Integration</span>
											<ExternalLink class="w-4 h-4" />
										</div>
									</div>
								</div>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	</section>
</div>

<style>
	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(30px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.animate-fadeInUp) {
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}
</style>
