<script lang="ts">
	import BlogPostsSection from '../_components/BlogPostsSection.svelte';
	import type { Post } from '$lib/types/post';

	interface Props {
		featuredPosts: Post[];
		latestPosts: Post[];
		loading: boolean;
		translations: any;
	}

	let { featuredPosts, latestPosts, loading, translations: t }: Props = $props();
</script>

<section id="blog" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('blog.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto">
				{t('blog.subtitle')}
			</p>
		</div>
		<BlogPostsSection featuredPosts={featuredPosts} latestPosts={latestPosts} loading={loading} />
	</div>
</section>

