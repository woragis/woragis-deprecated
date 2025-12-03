<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { BookOpen, ExternalLink, Calendar, Clock, Eye, Heart, Share2, MessageCircle, FileText, GraduationCap, Book, FileCode, Layers, ArrowLeft } from 'lucide-svelte';
	import { getTechnicalWriting } from '$lib/api/technical-writings';
	import type { TechnicalWriting, WritingType, PublicationPlatform } from '$lib/types/technical-writing';

	let writing: TechnicalWriting | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const writingId = $derived($page.params.slug);

	onMount(async () => {
		if (writingId) {
			await fetchWriting(writingId);
		}
	});

	async function fetchWriting(id: string) {
		loading = true;
		error = null;
		try {
			writing = await getTechnicalWriting(id);
			if (!writing) {
				error = 'Technical writing not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch technical writing';
			console.error('Error fetching technical writing:', err);
		} finally {
			loading = false;
		}
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
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function formatNumber(num?: number): string {
		if (!num) return '0';
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'k';
		}
		return num.toString();
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !writing}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">Technical Writing Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The technical writing you are looking for does not exist.'}</p>
				<a
					href="/technical-writings"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to Technical Writings
				</a>
			</div>
		</div>
	{:else}
		{@const TypeIcon = getTypeIcon(writing.type)}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/technical-writings"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to Technical Writings
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 hover:from-blue-500/5 hover:via-purple-500/5 hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<!-- Type Badge -->
							<div class="flex items-center gap-3 mb-4">
								<span
									class="inline-flex items-center gap-2 px-3 py-1 rounded-lg text-sm font-medium border {getTypeColor(writing.type)}"
								>
									<TypeIcon class="w-4 h-4" />
									{writing.type.replace('_', ' ')}
								</span>
								<span class="px-3 py-1 rounded-lg text-sm font-medium bg-gray-700/50 text-gray-300 border border-gray-600">
									{getPlatformLabel(writing.platform)}
								</span>
							</div>

							<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
								{writing.title}
							</h1>

							{#if writing.excerpt || writing.description}
								<p class="text-xl text-gray-300 mb-6 leading-relaxed">
									{writing.excerpt || writing.description}
								</p>
							{/if}

							<!-- Metadata -->
							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if writing.publishedAt}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Published {formatDate(writing.publishedAt)}</span>
									</div>
								{/if}
								{#if writing.readingTime}
									<div class="flex items-center gap-2">
										<Clock class="w-4 h-4" />
										<span>{writing.readingTime} min read</span>
									</div>
								{/if}
							</div>

							<!-- Engagement Metrics -->
							{#if writing.views || writing.likes || writing.shares || writing.comments}
								<div class="flex flex-wrap items-center gap-6 pt-6 border-t border-gray-700 mb-6">
									{#if writing.views}
										<div class="flex items-center gap-2 text-gray-300">
											<Eye class="w-5 h-5 text-blue-400" />
											<span class="font-medium">{formatNumber(writing.views)}</span>
											<span class="text-gray-400 text-sm">views</span>
										</div>
									{/if}
									{#if writing.likes}
										<div class="flex items-center gap-2 text-gray-300">
											<Heart class="w-5 h-5 text-red-400" />
											<span class="font-medium">{formatNumber(writing.likes)}</span>
											<span class="text-gray-400 text-sm">likes</span>
										</div>
									{/if}
									{#if writing.shares}
										<div class="flex items-center gap-2 text-gray-300">
											<Share2 class="w-5 h-5 text-green-400" />
											<span class="font-medium">{formatNumber(writing.shares)}</span>
											<span class="text-gray-400 text-sm">shares</span>
										</div>
									{/if}
									{#if writing.comments}
										<div class="flex items-center gap-2 text-gray-300">
											<MessageCircle class="w-5 h-5 text-purple-400" />
											<span class="font-medium">{formatNumber(writing.comments)}</span>
											<span class="text-gray-400 text-sm">comments</span>
										</div>
									{/if}
								</div>
							{/if}
						</div>

						<!-- Content -->
						{#if writing.content}
							<div class="prose prose-invert prose-lg max-w-none mb-8">
								<div class="text-gray-300 leading-relaxed whitespace-pre-wrap">
									{writing.content}
								</div>
							</div>
						{:else}
							<div class="bg-blue-900/20 border border-blue-700/30 rounded-lg p-6 mb-8 text-center">
								<p class="text-blue-300 mb-4">This article is published externally</p>
								<a
									href={writing.url}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200 font-medium"
								>
									Read on {getPlatformLabel(writing.platform)}
									<ExternalLink class="w-5 h-5" />
								</a>
							</div>
						{/if}

						<!-- Topics and Technologies -->
						{#if (writing.topics && writing.topics.length > 0) || (writing.technologies && writing.technologies.length > 0)}
							<div class="pt-8 border-t border-gray-700 space-y-4">
								{#if writing.topics && writing.topics.length > 0}
									<div>
										<h3 class="text-lg font-bold text-white mb-3">Topics</h3>
										<div class="flex flex-wrap gap-2">
											{#each writing.topics as topic}
												<span
													class="px-3 py-1 bg-blue-600/20 text-blue-300 text-sm rounded-lg border border-blue-500/30"
												>
													{topic}
												</span>
											{/each}
										</div>
									</div>
								{/if}

								{#if writing.technologies && writing.technologies.length > 0}
									<div>
										<h3 class="text-lg font-bold text-white mb-3">Technologies</h3>
										<div class="flex flex-wrap gap-2">
											{#each writing.technologies as tech}
												<span
													class="px-3 py-1 bg-purple-600/20 text-purple-300 text-sm rounded-lg border border-purple-500/30"
												>
													{tech}
												</span>
											{/each}
										</div>
									</div>
								{/if}
							</div>
						{/if}

						<!-- External Link -->
						<div class="mt-8 pt-8 border-t border-gray-700">
							<a
								href={writing.url}
								target="_blank"
								rel="noopener noreferrer"
								class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 rounded-lg transition-all duration-200 font-medium shadow-lg shadow-blue-500/50"
							>
								Read Full Article
								<ExternalLink class="w-5 h-5" />
							</a>
						</div>
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
