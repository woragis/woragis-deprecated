<script lang="ts">
	import TechnicalWritingsSectionComponent from '../_components/TechnicalWritingsSection.svelte';
	import type { TechnicalWriting } from '$lib/types/technical-writing';

	interface Props {
		technicalWritings: TechnicalWriting[];
		loading: boolean;
		translations: any;
	}

	let { technicalWritings, loading, translations: t }: Props = $props();
</script>

<section id="technical-writings" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
				{t('technicalWritings.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('technicalWritings.subtitle')}
			</p>
		</div>
		<TechnicalWritingsSectionComponent writings={technicalWritings} loading={loading} />
	</div>
</section>

