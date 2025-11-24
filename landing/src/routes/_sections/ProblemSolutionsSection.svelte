<script lang="ts">
	import ProblemSolutionsSectionComponent from '../_components/ProblemSolutionsSection.svelte';
	import type { ProblemSolution } from '$lib/types/problem-solution';

	interface Props {
		problemSolutions: ProblemSolution[];
		loading: boolean;
		translations: any;
	}

	let { problemSolutions, loading, translations: t }: Props = $props();
</script>

<section id="problem-solutions" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-yellow-400 to-orange-400 bg-clip-text text-transparent">
				{t('problemSolutions.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('problemSolutions.subtitle')}
			</p>
		</div>
		<ProblemSolutionsSectionComponent solutions={problemSolutions} loading={loading} />
	</div>
</section>

