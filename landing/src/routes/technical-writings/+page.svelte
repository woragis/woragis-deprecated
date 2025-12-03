<script lang="ts">
	import { Search, Filter, X, ChevronDown, ExternalLink, Calendar, Clock, Eye, Heart, Share2, MessageCircle, FileText, GraduationCap, Book, BookOpen, FileCode, Layers } from 'lucide-svelte';
	import type { TechnicalWriting, WritingType, PublicationPlatform } from '$lib/types/technical-writing';
	import { useTechnicalWritingsQuery } from '$lib/queries/technical-writings';

	// Filters
	let searchQuery = $state('');
	let typeFilter: WritingType | 'all' = $state('all');
	let platformFilter: PublicationPlatform | 'all' = $state('all');
	let sortBy: 'publishedAt' | 'title' | 'views' | 'likes' = $state('publishedAt');
	let sortOrder: 'asc' | 'desc' = $state('desc');
	let showFilters = $state(false);

	// Fetch all technical writings
	const writingsQuery = useTechnicalWritingsQuery();
	let writings = $derived(writingsQuery.data || []);
	let filteredWritings: TechnicalWriting[] = $state([]);
	let loading = $derived(writingsQuery.isPending);
	let error = $derived(writingsQuery.error ? (writingsQuery.error instanceof Error ? writingsQuery.error.message : 'Failed to fetch technical writings') : null);

	const typeOptions: Array<{ value: WritingType | 'all'; label: string }> = [
		{ value: 'all', label: 'All Types' },
		{ value: 'article', label: 'Articles' },
		{ value: 'documentation', label: 'Documentation' },
		{ value: 'tutorial', label: 'Tutorials' },
		{ value: 'guide', label: 'Guides' },
		{ value: 'blog_post', label: 'Blog Posts' },
		{ value: 'case_study', label: 'Case Studies' },
		{ value: 'other', label: 'Other' }
	];

	const platformOptions: Array<{ value: PublicationPlatform | 'all'; label: string }> = [
		{ value: 'all', label: 'All Platforms' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'dev_to', label: 'Dev.to' },
		{ value: 'hashnode', label: 'Hashnode' },
		{ value: 'personal_blog', label: 'Personal Blog' },
		{ value: 'github', label: 'GitHub' },
		{ value: 'company_blog', label: 'Company Blog' },
		{ value: 'substack', label: 'Substack' },
		{ value: 'linkedin', label: 'LinkedIn' },
		{ value: 'other', label: 'Other' }
	];

	const sortOptions = [
		{ value: 'publishedAt', label: 'Published Date' },
		{ value: 'title', label: 'Title' },
		{ value: 'views', label: 'Views' },
		{ value: 'likes', label: 'Likes' }
	];

	// Apply filters when writings or filter values change
	$effect(() => {
		applyFilters();
	});

	function applyFilters() {
		let filtered = [...writings];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(writing: TechnicalWriting) =>
					writing.title.toLowerCase().includes(query) ||
					writing.description?.toLowerCase().includes(query) ||
					writing.excerpt?.toLowerCase().includes(query) ||
					writing.topics?.some((topic) => topic.toLowerCase().includes(query)) ||
					writing.technologies?.some((tech) => tech.toLowerCase().includes(query))
			);
		}

		// Apply type filter
		if (typeFilter !== 'all') {
			filtered = filtered.filter((writing: TechnicalWriting) => writing.type === typeFilter);
		}

		// Apply platform filter
		if (platformFilter !== 'all') {
			filtered = filtered.filter((writing: TechnicalWriting) => writing.platform === platformFilter);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aVal: string | number | undefined;
			let bVal: string | number | undefined;

			switch (sortBy) {
				case 'publishedAt':
					aVal = a.publishedAt ? new Date(a.publishedAt).getTime() : 0;
					bVal = b.publishedAt ? new Date(b.publishedAt).getTime() : 0;
					break;
				case 'title':
					aVal = a.title.toLowerCase();
					bVal = b.title.toLowerCase();
					break;
				case 'views':
					aVal = a.views || 0;
					bVal = b.views || 0;
					break;
				case 'likes':
					aVal = a.likes || 0;
					bVal = b.likes || 0;
					break;
				default:
					return 0;
			}

			if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
			if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
			return 0;
		});

		filteredWritings = filtered;
	}

	function clearFilters() {
		searchQuery = '';
		typeFilter = 'all';
		platformFilter = 'all';
		applyFilters();
	}

	function getTypeIcon(type: WritingType) {
		const icons: Record<WritingType, typeof FileText> = {
			article: FileText,
			documentation: Book,
			tutorial: GraduationCap,
			guide: BookOpen,
			blog_post: FileText,
			case_study: Layers,
			other: FileCode
		};
		return icons[type] || FileText;
	}

	function getTypeColor(type: WritingType): string {
		const colors: Record<WritingType, string> = {
			article: 'bg-blue-600/20 text-blue-300 border-blue-500/30',
			documentation: 'bg-purple-600/20 text-purple-300 border-purple-500/30',
			tutorial: 'bg-green-600/20 text-green-300 border-green-500/30',
			guide: 'bg-yellow-600/20 text-yellow-300 border-yellow-500/30',
			blog_post: 'bg-indigo-600/20 text-indigo-300 border-indigo-500/30',
			case_study: 'bg-red-600/20 text-red-300 border-red-500/30',
			other: 'bg-gray-600/20 text-gray-300 border-gray-500/30'
		};
		return colors[type] || colors.other;
	}

	function getPlatformLabel(platform: PublicationPlatform): string {
		const labels: Record<PublicationPlatform, string> = {
			medium: 'Medium',
			dev_to: 'Dev.to',
			hashnode: 'Hashnode',
			personal_blog: 'Personal Blog',
			github: 'GitHub',
			company_blog: 'Company Blog',
			substack: 'Substack',
			linkedin: 'LinkedIn',
			other: 'Other'
		};
		return labels[platform] || 'Other';
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function formatNumber(num?: number): string {
		if (!num) return '0';
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'k';
		}
		return num.toString();
	}

	function getWritingUrl(writing: TechnicalWriting): string {
		// Use ID as slug since TechnicalWriting doesn't have a slug field
		return `/technical-writings/${writing.id}`;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Header -->
	<section class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
					Technical Writings
				</h1>
				<p class="text-xl text-gray-300 max-w-2xl mx-auto">
					Articles, tutorials, and documentation showcasing technical communication skills
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
							placeholder="Search writings by title, description, topics..."
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

							{#if searchQuery || typeFilter !== 'all' || platformFilter !== 'all'}
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

						<!-- Platform Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Platform</div>
							<div class="flex flex-wrap gap-2">
								{#each platformOptions as option}
									<button
										onclick={() => {
											platformFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {platformFilter === option.value
											? 'bg-purple-600 border-purple-500 text-white'
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										{option.label}
									</button>
								{/each}
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</section>

	<!-- Writings Grid -->
	<section class="container mx-auto px-6 py-12">
		<div class="max-w-7xl mx-auto">
			{#if loading}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-6 text-center">
					<p class="text-red-400 mb-2">Error loading technical writings</p>
					<p class="text-gray-400 text-sm">{error}</p>
					<button
						onclick={() => writingsQuery.refetch()}
						class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{:else if filteredWritings.length === 0}
				<div class="text-center py-20">
					<FileText class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No technical writings found</p>
					<p class="text-gray-500 text-sm">Try adjusting your filters or search query</p>
				</div>
			{:else}
				<div class="mb-4 text-gray-400 text-sm">
					Showing {filteredWritings.length} of {writings.length} technical writings
				</div>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredWritings as writing, index}
						{@const TypeIcon = getTypeIcon(writing.type)}
						{@const animationDelay = index * 0.1}
						<a
							href={getWritingUrl(writing)}
							class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl overflow-hidden border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/20 hover:scale-[1.02] relative animate-fadeInUp"
							style="animation-delay: {animationDelay}s"
						>
							<!-- Decorative gradient overlay -->
							<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
							<div class="relative z-10">
								<!-- Cover Image or Icon -->
								{#if writing.coverImageUrl}
									<div class="relative h-48 overflow-hidden">
										<img
											src={writing.coverImageUrl}
											alt={writing.title}
											class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
										/>
										<div class="absolute top-2 right-2">
											<span
												class="px-2 py-1 rounded-lg text-xs font-medium border {getTypeColor(writing.type)}"
											>
												{getPlatformLabel(writing.platform)}
											</span>
										</div>
									</div>
								{:else}
									<div
										class="h-48 bg-gradient-to-br from-gray-700/30 to-gray-800/30 flex items-center justify-center"
									>
										<TypeIcon class="w-16 h-16 text-white opacity-50" />
									</div>
								{/if}

								<div class="p-6">
									<!-- Type Badge -->
									<div class="flex items-center gap-2 mb-3">
										<span
											class="inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium border {getTypeColor(writing.type)}"
										>
											<TypeIcon class="w-3 h-3" />
											{writing.type.replace('_', ' ')}
										</span>
									</div>

									<!-- Title -->
									<h3 class="text-xl font-bold text-white mb-2 group-hover:text-blue-400 transition-colors line-clamp-2">
										{writing.title}
									</h3>

									<!-- Description/Excerpt -->
									{#if writing.excerpt}
										<p class="text-gray-300 text-sm mb-4 line-clamp-2">{writing.excerpt}</p>
									{:else if writing.description}
										<p class="text-gray-300 text-sm mb-4 line-clamp-2">{writing.description}</p>
									{/if}

									<!-- Metadata -->
									<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
										{#if writing.publishedAt}
											<div class="flex items-center gap-1">
												<Calendar class="w-3 h-3" />
												<span>{formatDate(writing.publishedAt)}</span>
											</div>
										{/if}
										{#if writing.readingTime}
											<div class="flex items-center gap-1">
												<Clock class="w-3 h-3" />
												<span>{writing.readingTime} min</span>
											</div>
										{/if}
									</div>

									<!-- Engagement Metrics -->
									{#if writing.views || writing.likes || writing.shares || writing.comments}
										<div class="flex items-center gap-4 pt-4 border-t border-gray-700 mb-4">
											{#if writing.views}
												<div class="flex items-center gap-1 text-xs text-gray-400">
													<Eye class="w-3 h-3" />
													<span>{formatNumber(writing.views)}</span>
												</div>
											{/if}
											{#if writing.likes}
												<div class="flex items-center gap-1 text-xs text-gray-400">
													<Heart class="w-3 h-3" />
													<span>{formatNumber(writing.likes)}</span>
												</div>
											{/if}
											{#if writing.shares}
												<div class="flex items-center gap-1 text-xs text-gray-400">
													<Share2 class="w-3 h-3" />
													<span>{formatNumber(writing.shares)}</span>
												</div>
											{/if}
											{#if writing.comments}
												<div class="flex items-center gap-1 text-xs text-gray-400">
													<MessageCircle class="w-3 h-3" />
													<span>{formatNumber(writing.comments)}</span>
												</div>
											{/if}
										</div>
									{/if}

									<!-- Topics/Technologies -->
									{#if writing.topics && writing.topics.length > 0}
										<div class="mb-4">
											<div class="flex flex-wrap gap-2">
												{#each writing.topics.slice(0, 3) as topic}
													<span
														class="px-2 py-1 bg-blue-600/20 text-blue-300 text-xs rounded-lg border border-blue-500/30"
													>
														{topic}
													</span>
												{/each}
											</div>
										</div>
									{/if}

									<!-- External Link Indicator -->
									<div class="flex items-center justify-between pt-4 border-t border-gray-700">
										<span class="text-xs text-gray-400">Read article</span>
										<ExternalLink class="w-4 h-4 text-gray-400 group-hover:text-blue-400 transition-colors" />
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
