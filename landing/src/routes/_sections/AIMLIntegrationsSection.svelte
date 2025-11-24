<script lang="ts">
	import AIMLIntegrationsSectionComponent from '../_components/AIMLIntegrationsSection.svelte';
	import type { AIMLIntegration } from '$lib/types/aiml-integration';

	interface Props {
		aimlIntegrations: AIMLIntegration[];
		loading: boolean;
		translations: any;
	}

	let { aimlIntegrations, loading, translations: t }: Props = $props();
</script>

<section id="aiml-integrations" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('aimlIntegrations.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('aimlIntegrations.subtitle')}
			</p>
		</div>
		<AIMLIntegrationsSectionComponent integrations={aimlIntegrations} loading={loading} />
	</div>
</section>

