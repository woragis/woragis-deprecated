<script lang="ts">
	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section class="container mx-auto px-6 py-20 md:py-32 relative overflow-hidden">
	<!-- Decorative background elements -->
	<div class="absolute inset-0 overflow-hidden pointer-events-none">
		<div class="absolute top-1/4 left-1/4 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl"></div>
		<div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl"></div>
	</div>
	
	<div class="max-w-4xl mx-auto text-center relative z-10">
		<div
			class="inline-block mb-6 px-4 py-2 bg-blue-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium animate-fadeInUp"
		>
			{t('hero.badge')}
		</div>
		<h1
			class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent animate-fadeInUp"
			style="animation-delay: 0.1s"
		>
			{t('hero.title')}
		</h1>
		<p class="text-xl md:text-2xl text-gray-300 mb-8 leading-relaxed animate-fadeInUp" style="animation-delay: 0.2s">
			{t('hero.subtitle')}
		</p>
		<div class="flex flex-wrap gap-4 justify-center animate-fadeInUp" style="animation-delay: 0.3s">
			<a
				href="#projects"
				class="px-8 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-all duration-200 shadow-lg shadow-blue-500/50 hover:shadow-xl hover:shadow-blue-500/70 hover:scale-105"
			>
				{t('hero.viewProjects')}
			</a>
			<a
				href="#technical-depth"
				class="px-8 py-3 bg-purple-600 hover:bg-purple-700 rounded-lg font-medium transition-all duration-200 shadow-lg shadow-purple-500/50 hover:shadow-xl hover:shadow-purple-500/70 hover:scale-105"
			>
				{t('hero.technicalDepth')}
			</a>
			<a
				href="/skills"
				class="px-8 py-3 bg-gray-700/50 hover:bg-gray-600/50 rounded-lg font-medium transition-all duration-200 border border-gray-600 hover:border-gray-500 hover:scale-105"
			>
				{t('hero.viewSkills')}
			</a>
			<a
				href="#testimonials"
				class="px-8 py-3 bg-gray-700/50 hover:bg-gray-600/50 rounded-lg font-medium transition-all duration-200 border border-gray-600 hover:border-gray-500 hover:scale-105"
			>
				{t('testimonials.title')}
			</a>
			<a
				href="#contact"
				class="px-8 py-3 bg-gray-700/50 hover:bg-gray-600/50 rounded-lg font-medium transition-all duration-200 border border-gray-600 hover:border-gray-500 hover:scale-105"
			>
				{t('hero.contactMe')}
			</a>
		</div>
	</div>
</section>

<style>
	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(30px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.animate-fadeInUp) {
		animation: fadeInUp 0.8s ease-out;
		animation-fill-mode: both;
	}
</style>

