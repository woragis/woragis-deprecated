<script lang="ts">
	import SystemDesignsSectionComponent from '../_components/SystemDesignsSection.svelte';
	import type { SystemDesign } from '$lib/types/system-design';

	interface Props {
		systemDesigns: SystemDesign[];
		loading: boolean;
		translations: any;
	}

	let { systemDesigns, loading, translations: t }: Props = $props();
</script>

<section id="system-designs" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('systemDesigns.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto">
				{t('systemDesigns.subtitle')}
			</p>
		</div>
		<SystemDesignsSectionComponent designs={systemDesigns} loading={loading} />
	</div>
</section>

