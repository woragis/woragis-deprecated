<script lang="ts">
	import { Calendar, Clock, ExternalLink, Tag as TagIcon, FolderOpen, Star } from 'lucide-svelte';
	import { calculateReadingTime, getPostSkills, getPostCategories, getPostTags } from '$lib/api/posts';
	import { usePostsQuery } from '$lib/queries/posts';
	import type { Post } from '$lib/types/post';

	// Fetch posts using TanStack Query
	const featuredPostsQuery = usePostsQuery({
		status: 'published',
		featured: true,
		limit: 3,
		orderBy: 'publishedAt',
		order: 'desc'
	});

	const latestPostsQuery = usePostsQuery({
		status: 'published',
		limit: 6,
		orderBy: 'publishedAt',
		order: 'desc'
	});

	let featuredPosts: Post[] = $state([]);
	let latestPosts: Post[] = $state([]);
	let posts: Post[] = $derived([...featuredPosts, ...latestPosts.filter(p => !featuredPosts.some(fp => fp.id === p.id))]);
	let loading = $derived(featuredPostsQuery.isPending || latestPostsQuery.isPending);

	// Enrich posts when data is available
	$effect(async () => {
		if (featuredPostsQuery.data) {
			featuredPosts = await enrichPosts(featuredPostsQuery.data);
		}
	});

	$effect(async () => {
		if (latestPostsQuery.data) {
			latestPosts = await enrichPosts(latestPostsQuery.data.slice(0, 6));
		}
	});

	async function enrichPosts(postsToEnrich: Post[]): Promise<Post[]> {
		return Promise.all(
			postsToEnrich.map(async (post) => {
				try {
					const [skills, categories, tags] = await Promise.all([
						getPostSkills(post.id),
						getPostCategories(post.id),
						getPostTags(post.id)
					]);
					return {
						...post,
						skills,
						categories,
						tags
					};
				} catch (error) {
					console.error(`Error enriching post ${post.id}:`, error);
					return post;
				}
			})
		);
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getPostUrl(post: Post): string {
		return `/blog/${post.slug}`;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if posts.length === 0}
		<div class="text-center py-20">
			<p class="text-gray-400 text-lg mb-2">No blog posts available</p>
			<p class="text-gray-500 text-sm">Check back later for new posts</p>
		</div>
	{:else}
		<!-- Featured Posts -->
		{#if featuredPosts.length > 0}
			<div class="mb-12">
				<div class="flex items-center gap-2 mb-6">
					<Star class="w-5 h-5 text-yellow-400 fill-yellow-400" />
					<h3 class="text-2xl font-bold text-white">Featured Posts</h3>
				</div>
				<div class="grid md:grid-cols-3 gap-6">
					{#each featuredPosts as post}
						<a
							href={getPostUrl(post)}
							class="group bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl overflow-hidden border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
						>
							{#if post.featuredImage}
								<div class="relative h-48 overflow-hidden">
									<img
										src={post.featuredImage}
										alt={post.title}
										class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
									/>
									<div
										class="absolute top-2 right-2 px-2 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded flex items-center gap-1"
									>
										<Star class="w-3 h-3 fill-current" />
										Featured
									</div>
								</div>
							{/if}
							<div class="p-6">
								<h4
									class="text-xl font-bold text-white mb-2 group-hover:text-blue-400 transition-colors line-clamp-2"
								>
									{post.title}
								</h4>
								{#if post.excerpt}
									<p class="text-gray-300 text-sm mb-4 line-clamp-3">{post.excerpt}</p>
								{/if}
								<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
									{#if post.publishedAt}
										<div class="flex items-center gap-1">
											<Calendar class="w-3 h-3" />
											<span>{formatDate(post.publishedAt)}</span>
										</div>
									{/if}
									<div class="flex items-center gap-1">
										<Clock class="w-3 h-3" />
										<span>{calculateReadingTime(post.content)} min read</span>
									</div>
								</div>
								{#if post.categories && post.categories.length > 0}
									<div class="flex flex-wrap gap-2 mb-3">
										{#each post.categories.slice(0, 2) as category}
											<span
												class="px-2 py-1 bg-blue-600/20 text-blue-300 text-xs rounded border border-blue-600/30 flex items-center gap-1"
											>
												<FolderOpen class="w-3 h-3" />
												{category.name}
											</span>
										{/each}
									</div>
								{/if}
								<div class="flex items-center gap-2 text-blue-400 text-sm font-medium group-hover:gap-3 transition-all">
									<span>Read More</span>
									<ExternalLink class="w-4 h-4" />
								</div>
							</div>
						</a>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Latest Posts -->
		<div>
			<h3 class="text-2xl font-bold text-white mb-6">Latest Posts</h3>
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each latestPosts.slice(0, 6) as post}
					<a
						href={getPostUrl(post)}
						class="group bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
					>
						<h4
							class="text-lg font-bold text-white mb-2 group-hover:text-blue-400 transition-colors line-clamp-2"
						>
							{post.title}
						</h4>
						{#if post.excerpt}
							<p class="text-gray-300 text-sm mb-4 line-clamp-2">{post.excerpt}</p>
						{/if}
						<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
							{#if post.publishedAt}
								<div class="flex items-center gap-1">
									<Calendar class="w-3 h-3" />
									<span>{formatDate(post.publishedAt)}</span>
								</div>
							{/if}
							<div class="flex items-center gap-1">
								<Clock class="w-3 h-3" />
								<span>{calculateReadingTime(post.content)} min read</span>
							</div>
						</div>
						{#if post.tags && post.tags.length > 0}
							<div class="flex flex-wrap gap-2 mb-3">
								{#each post.tags.slice(0, 3) as tag}
									<span
										class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded border border-purple-600/30 flex items-center gap-1"
									>
										<TagIcon class="w-3 h-3" />
										{tag.name}
									</span>
								{/each}
							</div>
						{/if}
						<div class="flex items-center gap-2 text-blue-400 text-sm font-medium group-hover:gap-3 transition-all">
							<span>Read More</span>
							<ExternalLink class="w-4 h-4" />
						</div>
					</a>
				{/each}
			</div>
		</div>
	{/if}
</div>

