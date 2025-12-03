<script lang="ts">
	import ImpactMetricsDashboard from '../_components/ImpactMetricsDashboard.svelte';
	import type { ImpactMetric } from '$lib/types/impact-metric';

	interface Props {
		impactMetrics: ImpactMetric[];
		loading: boolean;
		translations: any;
	}

	let { impactMetrics, loading, translations: t }: Props = $props();
</script>

<section id="impact-metrics" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('impactMetrics.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto">
				{t('impactMetrics.subtitle')}
			</p>
		</div>
		<ImpactMetricsDashboard metrics={impactMetrics} loading={loading} />
	</div>
</section>

