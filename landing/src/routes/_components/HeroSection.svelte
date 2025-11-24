<script lang="ts">
	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section class="container mx-auto px-6 py-20 md:py-32">
	<div class="max-w-4xl mx-auto text-center">
		<div
			class="inline-block mb-6 px-4 py-2 bg-blue-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium"
		>
			{t('hero.badge')}
		</div>
		<h1
			class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent animate-pulse"
		>
			{t('hero.title')}
		</h1>
		<p class="text-xl md:text-2xl text-gray-300 mb-8 leading-relaxed">
			{t('hero.subtitle')}
		</p>
		<div class="flex flex-wrap gap-4 justify-center">
			<a
				href="#projects"
				class="px-8 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
			>
				{t('hero.viewProjects')}
			</a>
			<a
				href="#technical-depth"
				class="px-8 py-3 bg-purple-600 hover:bg-purple-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-purple-500/50"
			>
				{t('hero.technicalDepth')}
			</a>
			<a
				href="/skills"
				class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
			>
				{t('hero.viewSkills')}
			</a>
			<a
				href="#testimonials"
				class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
			>
				{t('testimonials.title')}
			</a>
			<a
				href="#contact"
				class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
			>
				{t('hero.contactMe')}
			</a>
		</div>
	</div>
</section>

