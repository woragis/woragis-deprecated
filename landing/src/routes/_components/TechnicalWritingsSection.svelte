<script lang="ts">
	import { BookOpen, ExternalLink, Calendar, Clock, Eye, Heart, Share2, MessageCircle, FileText, Code, GraduationCap, Book, FileCode, Layers } from 'lucide-svelte';
	import type { TechnicalWriting, WritingType, PublicationPlatform } from '$lib/types/technical-writing';

	interface Props {
		writings: TechnicalWriting[];
		loading?: boolean;
	}

	let { writings = [], loading = false }: Props = $props();

	// Group writings by type
	let groupedWritings = $derived.by(() => {
		const grouped: Record<WritingType, TechnicalWriting[]> = {
			article: [],
			documentation: [],
			tutorial: [],
			guide: [],
			blog_post: [],
			case_study: [],
			other: []
		};

		writings.forEach((writing) => {
			if (grouped[writing.type]) {
				grouped[writing.type].push(writing);
			} else {
				grouped.other.push(writing);
			}
		});

		// Sort each group by display order, then by published date
		Object.keys(grouped).forEach((key) => {
			grouped[key as WritingType].sort((a, b) => {
				if (a.displayOrder !== b.displayOrder) {
					return a.displayOrder - b.displayOrder;
				}
				if (a.publishedAt && b.publishedAt) {
					return new Date(b.publishedAt).getTime() - new Date(a.publishedAt).getTime();
				}
				return 0;
			});
		});

		// Filter out empty groups
		return Object.entries(grouped).filter(([_, items]) => items.length > 0);
	});

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
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
			article: 'from-blue-500 to-cyan-500',
			documentation: 'from-purple-500 to-pink-500',
			tutorial: 'from-green-500 to-emerald-500',
			guide: 'from-yellow-500 to-amber-500',
			blog_post: 'from-indigo-500 to-purple-500',
			case_study: 'from-red-500 to-orange-500',
			other: 'from-gray-500 to-slate-500'
		};
		return colors[type] || colors.other;
	}

	function getTypeLabel(type: WritingType): string {
		const labels: Record<WritingType, string> = {
			article: 'Articles',
			documentation: 'Documentation',
			tutorial: 'Tutorials',
			guide: 'Guides',
			blog_post: 'Blog Posts',
			case_study: 'Case Studies',
			other: 'Other'
		};
		return labels[type] || 'Other';
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

	function formatNumber(num?: number): string {
		if (!num) return '0';
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'k';
		}
		return num.toString();
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if writings.length === 0}
		<div class="text-center py-20">
			<BookOpen class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No technical writings available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="space-y-8">
			{#each groupedWritings as [type, items]}
				{@const TypeIcon = getTypeIcon(type as WritingType)}
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<!-- Type Header -->
					<div class="flex items-center gap-3 mb-6">
						<div
							class="w-12 h-12 bg-gradient-to-br {getTypeColor(type as WritingType)} rounded-lg flex items-center justify-center"
						>
							<TypeIcon class="w-6 h-6 text-white" />
						</div>
						<div>
							<h3 class="text-2xl font-bold text-white">{getTypeLabel(type as WritingType)}</h3>
							<p class="text-sm text-gray-400">{items.length} {items.length === 1 ? 'piece' : 'pieces'}</p>
						</div>
					</div>

					<!-- Writings Grid -->
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each items as writing}
							<a
								href={writing.url}
								target="_blank"
								rel="noopener noreferrer"
								class="group bg-gray-800/50 rounded-lg p-5 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 flex flex-col"
							>
								<!-- Cover Image or Icon -->
								{#if writing.coverImageUrl}
									<img
										src={writing.coverImageUrl}
										alt={writing.title}
										class="w-full h-40 object-cover rounded-lg mb-4"
									/>
								{:else}
									{@const WritingIcon = getTypeIcon(writing.type)}
									<div
										class="w-full h-40 bg-gradient-to-br {getTypeColor(writing.type)} rounded-lg mb-4 flex items-center justify-center"
									>
										<WritingIcon class="w-12 h-12 text-white opacity-50" />
									</div>
								{/if}

								<!-- Title -->
								<h4 class="text-lg font-bold text-white mb-2 group-hover:text-blue-400 transition-colors line-clamp-2">
									{writing.title}
								</h4>

								<!-- Description/Excerpt -->
								{#if writing.excerpt}
									<p class="text-sm text-gray-300 mb-4 line-clamp-3 flex-1">{writing.excerpt}</p>
								{:else if writing.description}
									<p class="text-sm text-gray-300 mb-4 line-clamp-3 flex-1">{writing.description}</p>
								{/if}

								<!-- Platform Badge -->
								<div class="mb-4">
									<span
										class="inline-block px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
									>
										{getPlatformLabel(writing.platform)}
									</span>
								</div>

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
									<div class="flex items-center gap-4 pt-4 border-t border-gray-700">
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
									<div class="mt-4 pt-4 border-t border-gray-700">
										<div class="flex flex-wrap gap-2">
											{#each writing.topics.slice(0, 3) as topic}
												<span
													class="px-2 py-1 text-xs rounded bg-blue-600/20 text-blue-300 border border-blue-500/30"
												>
													{topic}
												</span>
											{/each}
										</div>
									</div>
								{/if}

								<!-- External Link Indicator -->
								<div class="mt-4 pt-4 border-t border-gray-700 flex items-center justify-between">
									<span class="text-xs text-gray-400">Read article</span>
									<ExternalLink class="w-4 h-4 text-gray-400 group-hover:text-blue-400 transition-colors" />
								</div>
							</a>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

