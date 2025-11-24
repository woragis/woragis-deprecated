<script lang="ts">
	import { Share2, ExternalLink, Calendar, Heart, MessageCircle, Eye, TrendingUp, Linkedin, Twitter, Instagram } from 'lucide-svelte';
	import type { SocialMediaPost, Platform } from '$lib/types/social-media-post';

	interface Props {
		posts: SocialMediaPost[];
		loading?: boolean;
	}

	let { posts: allPosts = [], loading = false }: Props = $props();

	// Sort posts by engagement (total engagement score)
	let sortedPosts = $derived.by(() => {
		return [...allPosts].sort((a, b) => {
			const scoreA = getEngagementScore(a);
			const scoreB = getEngagementScore(b);
			return scoreB - scoreA;
		});
	});

	// Group posts by platform
	let groupedPosts = $derived.by(() => {
		const grouped: Record<Platform, SocialMediaPost[]> = {
			linkedin: [],
			twitter: [],
			instagram: []
		};

		sortedPosts.forEach((post) => {
			if (grouped[post.platform]) {
				grouped[post.platform].push(post);
			}
		});

		// Filter out empty platforms
		return Object.entries(grouped).filter(([_, posts]) => posts.length > 0);
	});

	function getEngagementScore(post: SocialMediaPost): number {
		const likes = post.likes || 0;
		const shares = post.shares || 0;
		const comments = post.comments || 0;
		const views = post.views || 0;
		// Weighted score: likes (3x), shares (5x), comments (4x), views (0.1x)
		return likes * 3 + shares * 5 + comments * 4 + views * 0.1;
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

		if (diffDays === 0) return 'Today';
		if (diffDays === 1) return 'Yesterday';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
		if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
	}

	function getPlatformIcon(platform: Platform) {
		const icons: Record<Platform, typeof Linkedin> = {
			linkedin: Linkedin,
			twitter: Twitter,
			instagram: Instagram
		};
		return icons[platform] || Linkedin;
	}

	function getPlatformColor(platform: Platform): string {
		const colors: Record<Platform, string> = {
			linkedin: 'from-blue-600 to-blue-700',
			twitter: 'from-sky-500 to-sky-600',
			instagram: 'from-pink-500 via-purple-500 to-orange-500'
		};
		return colors[platform] || 'from-gray-500 to-gray-600';
	}

	function getPlatformLabel(platform: Platform): string {
		const labels: Record<Platform, string> = {
			linkedin: 'LinkedIn',
			twitter: 'Twitter',
			instagram: 'Instagram'
		};
		return labels[platform] || platform;
	}

	function formatNumber(num?: number): string {
		if (!num) return '0';
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		}
		if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'k';
		}
		return num.toString();
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
	{:else if allPosts.length === 0}
		<div class="text-center py-20">
			<Share2 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No social media posts available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="space-y-8">
			{#each groupedPosts as [platform, posts]}
				{@const PlatformIcon = getPlatformIcon(platform as Platform)}
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<!-- Platform Header -->
					<div class="flex items-center gap-3 mb-6">
						<div
							class="w-12 h-12 bg-gradient-to-br {getPlatformColor(platform as Platform)} rounded-lg flex items-center justify-center"
						>
							<PlatformIcon class="w-6 h-6 text-white" />
						</div>
						<div>
							<h3 class="text-2xl font-bold text-white">{getPlatformLabel(platform as Platform)}</h3>
							<p class="text-sm text-gray-400">{posts.length} {posts.length === 1 ? 'post' : 'posts'}</p>
						</div>
					</div>

					<!-- Posts Grid -->
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each posts.slice(0, 6) as post}
							<a
								href={post.url}
								target="_blank"
								rel="noopener noreferrer"
								class="group bg-gray-800/50 rounded-lg p-5 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 flex flex-col"
							>
								<!-- Header -->
								<div class="flex items-start justify-between mb-3">
									<div class="flex-1">
										{#if post.title}
											<h4 class="text-lg font-bold text-white mb-2 group-hover:text-blue-400 transition-colors line-clamp-2">
												{post.title}
											</h4>
										{/if}
										{#if post.contentPreview}
											<p class="text-sm text-gray-300 mb-3 line-clamp-3">
												{truncateText(post.contentPreview, 150)}
											</p>
										{/if}
									</div>
								</div>

								<!-- Published Date -->
								{#if post.publishedDate}
									<div class="flex items-center gap-2 text-xs text-gray-400 mb-4">
										<Calendar class="w-3 h-3" />
										<span>{formatDate(post.publishedDate)}</span>
									</div>
								{/if}

								<!-- Engagement Metrics -->
								<div class="flex items-center gap-4 pt-4 border-t border-gray-700 mt-auto">
									{#if post.views}
										<div class="flex items-center gap-1 text-xs text-gray-400" title="Views">
											<Eye class="w-4 h-4" />
											<span>{formatNumber(post.views)}</span>
										</div>
									{/if}
									{#if post.likes}
										<div class="flex items-center gap-1 text-xs text-gray-400" title="Likes">
											<Heart class="w-4 h-4" />
											<span>{formatNumber(post.likes)}</span>
										</div>
									{/if}
									{#if post.comments}
										<div class="flex items-center gap-1 text-xs text-gray-400" title="Comments">
											<MessageCircle class="w-4 h-4" />
											<span>{formatNumber(post.comments)}</span>
										</div>
									{/if}
									{#if post.shares}
										<div class="flex items-center gap-1 text-xs text-gray-400" title="Shares">
											<Share2 class="w-4 h-4" />
											<span>{formatNumber(post.shares)}</span>
										</div>
									{/if}
									{#if getEngagementScore(post) > 0}
										<div
											class="ml-auto flex items-center gap-1 text-xs text-green-400"
											title="Engagement Score"
										>
											<TrendingUp class="w-4 h-4" />
											<span>{Math.round(getEngagementScore(post))}</span>
										</div>
									{/if}
								</div>

								<!-- External Link Indicator -->
								<div class="mt-4 pt-4 border-t border-gray-700 flex items-center justify-between">
									<span class="text-xs text-gray-400">View post</span>
									<ExternalLink class="w-4 h-4 text-gray-400 group-hover:text-blue-400 transition-colors" />
								</div>
							</a>
						{/each}
					</div>

					<!-- Show More Indicator -->
					{#if posts.length > 6}
						<div class="mt-4 text-center">
							<p class="text-sm text-gray-400">
								And {posts.length - 6} more {posts.length - 6 === 1 ? 'post' : 'posts'}...
							</p>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

