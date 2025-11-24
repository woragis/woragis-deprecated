<script lang="ts">
	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section id="about" class="container mx-auto px-6 py-20">
	<div class="max-w-4xl mx-auto">
		<h2 class="text-4xl font-bold mb-8 text-center">{t('about.title')}</h2>
		<div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 border border-gray-700 shadow-xl">
			<p class="text-lg text-gray-300 leading-relaxed">
				{t('about.description')}
			</p>
		</div>
	</div>
</section>

