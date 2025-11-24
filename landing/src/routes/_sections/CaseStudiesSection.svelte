<script lang="ts">
	import CaseStudiesSectionComponent from '../_components/CaseStudiesSection.svelte';
	import type { CaseStudy } from '$lib/types/case-study';

	interface Props {
		featuredCaseStudies: CaseStudy[];
		latestCaseStudies: CaseStudy[];
		loading: boolean;
		translations: any;
	}

	let { featuredCaseStudies, latestCaseStudies, loading, translations: t }: Props = $props();
</script>

<section id="case-studies" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
				Case Studies
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				Real-world projects with detailed technical implementations, challenges, and solutions
			</p>
		</div>
		<CaseStudiesSectionComponent featuredCaseStudies={featuredCaseStudies} latestCaseStudies={latestCaseStudies} loading={loading} />
	</div>
</section>

