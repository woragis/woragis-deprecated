<script lang="ts">
	import TestimonialsCarousel from '../_components/TestimonialsCarousel.svelte';
	import type { Testimonial } from '$lib/types/testimonial';

	interface Props {
		testimonials: Testimonial[];
		loading: boolean;
		translations: any;
	}

	let { testimonials, loading, translations: t }: Props = $props();
</script>

<section id="testimonials" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('testimonials.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto">
				{t('testimonials.subtitle')}
			</p>
		</div>
		<TestimonialsCarousel testimonials={testimonials} loading={loading} />
	</div>
</section>

