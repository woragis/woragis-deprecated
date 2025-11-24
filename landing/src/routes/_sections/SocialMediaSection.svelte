<script lang="ts">
	import SocialMediaPostsSection from '../_components/SocialMediaPostsSection.svelte';
	import type { SocialMediaPost } from '$lib/types/social-media-post';

	interface Props {
		socialMediaPosts: SocialMediaPost[];
		loading: boolean;
		translations: any;
	}

	let { socialMediaPosts, loading, translations: t }: Props = $props();
</script>

<section id="social-media" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-pink-400 to-purple-400 bg-clip-text text-transparent">
				{t('socialMedia.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('socialMedia.subtitle')}
			</p>
		</div>
		<SocialMediaPostsSection posts={socialMediaPosts} loading={loading} />
	</div>
</section>

