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
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('problemSolutions.matrix.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto">
				{t('problemSolutions.matrix.subtitle')}
			</p>
		</div>
		<ProblemSolutionMatrix />
	</div>
</section>

