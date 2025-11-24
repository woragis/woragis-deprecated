<script lang="ts">
	import ProblemSolutionMatrix from '../_components/ProblemSolutionMatrix.svelte';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section id="problem-solution-matrix" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-yellow-400 via-orange-400 to-red-400 bg-clip-text text-transparent">
				{t('problemSolutions.matrix.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('problemSolutions.matrix.subtitle')}
			</p>
		</div>
		<ProblemSolutionMatrix />
	</div>
</section>

