<script lang="ts">
	import { Search, Filter, X, ChevronDown, ExternalLink, Calendar, Code2, Target, TrendingUp, ArrowLeft } from 'lucide-svelte';
	import type { CaseStudy } from '$lib/types/case-study';
	import { useCaseStudiesQuery } from '$lib/queries/case-studies';

	// Filters
	let searchQuery = $state('');
	let featuredFilter: boolean | 'all' = $state('all');
	let sortBy: 'updatedAt' | 'title' | 'createdAt' = $state('updatedAt');
	let sortOrder: 'asc' | 'desc' = $state('desc');
	let showFilters = $state(false);

	// Fetch all case studies
	const caseStudiesQuery = useCaseStudiesQuery();
	let caseStudies = $derived(caseStudiesQuery.data || []);
	let filteredCaseStudies: CaseStudy[] = $state([]);
	let loading = $derived(caseStudiesQuery.isPending);
	let error = $derived(caseStudiesQuery.error ? (caseStudiesQuery.error instanceof Error ? caseStudiesQuery.error.message : 'Failed to fetch case studies') : null);

	const sortOptions = [
		{ value: 'updatedAt', label: 'Last Updated' },
		{ value: 'createdAt', label: 'Date Created' },
		{ value: 'title', label: 'Title' }
	];

	// Apply filters when case studies or filter values change
	$effect(() => {
		applyFilters();
	});

	function applyFilters() {
		let filtered = [...caseStudies];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(caseStudy: CaseStudy) =>
					caseStudy.title.toLowerCase().includes(query) ||
					caseStudy.problem?.toLowerCase().includes(query) ||
					caseStudy.context?.toLowerCase().includes(query) ||
					caseStudy.solution?.toLowerCase().includes(query) ||
					caseStudy.technologies?.some((tech) => tech.toLowerCase().includes(query))
			);
		}

		// Apply featured filter
		if (featuredFilter !== 'all') {
			filtered = filtered.filter((caseStudy: CaseStudy) => caseStudy.featured === featuredFilter);
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

		filteredCaseStudies = filtered;
	}

	function clearFilters() {
		searchQuery = '';
		featuredFilter = 'all';
		applyFilters();
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getCaseStudyUrl(caseStudy: CaseStudy): string {
		return `/case-studies/${caseStudy.id}`;
	}

	function getProjectUrl(caseStudy: CaseStudy): string {
		return `/projects/${caseStudy.projectSlug}`;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Header -->
	<section class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
					Case Studies
				</h1>
				<p class="text-xl text-gray-300 max-w-2xl mx-auto">
					In-depth case studies showcasing project challenges, solutions, metrics, and lessons learned
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
							placeholder="Search case studies by title, problem, solution, technologies..."
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

							{#if searchQuery || featuredFilter !== 'all'}
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

	<!-- Case Studies Grid -->
	<section class="container mx-auto px-6 py-12">
		<div class="max-w-7xl mx-auto">
			{#if loading}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-6 text-center">
					<p class="text-red-400 mb-2">Error loading case studies</p>
					<p class="text-gray-400 text-sm">{error}</p>
					<button
						onclick={() => caseStudiesQuery.refetch()}
						class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{:else if filteredCaseStudies.length === 0}
				<div class="text-center py-20">
					<Code2 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No case studies found</p>
					<p class="text-gray-500 text-sm">Try adjusting your filters or search query</p>
				</div>
			{:else}
				<div class="mb-4 text-gray-400 text-sm">
					Showing {filteredCaseStudies.length} of {caseStudies.length} case studies
				</div>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredCaseStudies as caseStudy, index}
						{@const animationDelay = index * 0.1}
						<a
							href={getCaseStudyUrl(caseStudy)}
							class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl overflow-hidden border border-gray-700 hover:border-purple-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-purple-500/20 hover:scale-[1.02] relative animate-fadeInUp"
							style="animation-delay: {animationDelay}s"
						>
							<!-- Decorative gradient overlay -->
							<div class="absolute inset-0 bg-gradient-to-br from-purple-500/0 via-pink-500/0 to-orange-500/0 group-hover:from-purple-500/5 group-hover:via-pink-500/5 group-hover:to-orange-500/5 transition-all duration-300 pointer-events-none"></div>
							<div class="relative z-10">
								<div class="p-6">
									<div class="flex items-start justify-between mb-3">
										<div
											class="w-12 h-12 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center flex-shrink-0"
										>
											<Target class="w-6 h-6 text-white" />
										</div>
										{#if caseStudy.featured}
											<div
												class="px-2 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded flex items-center gap-1"
											>
												⭐ Featured
											</div>
										{/if}
									</div>

									<h3
										class="text-xl font-bold text-white mb-2 group-hover:text-purple-400 transition-colors line-clamp-2"
									>
										{caseStudy.title}
									</h3>

									{#if caseStudy.context}
										<p class="text-gray-300 text-sm mb-4 line-clamp-3">{caseStudy.context}</p>
									{:else if caseStudy.problem}
										<p class="text-gray-300 text-sm mb-4 line-clamp-3">{caseStudy.problem}</p>
									{/if}

									<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
										{#if caseStudy.updatedAt}
											<div class="flex items-center gap-1">
												<Calendar class="w-3 h-3" />
												<span>Updated {formatDate(caseStudy.updatedAt)}</span>
											</div>
										{/if}
									</div>

									{#if caseStudy.technologies && caseStudy.technologies.length > 0}
										<div class="flex flex-wrap gap-2 mb-4">
											{#each caseStudy.technologies.slice(0, 3) as tech}
												<span
													class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded-lg border border-purple-600/30"
												>
													{tech}
												</span>
											{/each}
											{#if caseStudy.technologies.length > 3}
												<span class="px-2 py-1 text-xs text-gray-400">+{caseStudy.technologies.length - 3}</span>
											{/if}
										</div>
									{/if}

									<div class="flex items-center justify-between pt-4 border-t border-gray-700">
										<div class="flex items-center gap-2 text-purple-400 text-sm font-medium group-hover:gap-3 transition-all">
											<span>View Case Study</span>
											<ExternalLink class="w-4 h-4" />
										</div>
										{#if caseStudy.projectSlug}
											<a
												href={getProjectUrl(caseStudy)}
												onclick={(e) => e.stopPropagation()}
												class="text-xs text-gray-400 hover:text-blue-400 transition-colors"
											>
												View Project
											</a>
										{/if}
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
