<script lang="ts">
	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section id="about" class="container mx-auto px-6 py-20">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('about.title')}
			</h2>
		</div>
		<div class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl hover:border-blue-500/50 transition-all duration-300 relative overflow-hidden">
			<!-- Decorative gradient overlay -->
			<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 hover:from-blue-500/5 hover:via-purple-500/5 hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
			<p class="text-lg md:text-xl text-gray-300 leading-relaxed relative z-10">
				{t('about.description')}
			</p>
		</div>
	</div>
</section>

