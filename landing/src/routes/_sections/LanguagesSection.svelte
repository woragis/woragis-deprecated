<script lang="ts">
	import { Languages, Globe, CheckCircle } from 'lucide-svelte';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();

	// Languages data
	const languages = [
		{
			name: 'Portuguese',
			proficiency: 'native',
			level: 100,
			flag: '🇧🇷'
		},
		{
			name: 'English',
			proficiency: 'fluent',
			level: 90,
			flag: '🇺🇸'
		},
		{
			name: 'French',
			proficiency: 'fluent',
			level: 90,
			flag: '🇫🇷'
		},
		{
			name: 'German',
			proficiency: 'intermediate',
			level: 60,
			flag: '🇩🇪'
		},
		{
			name: 'Japanese',
			proficiency: 'intermediate',
			level: 60,
			flag: '🇯🇵'
		},
		{
			name: 'Russian',
			proficiency: 'beginner',
			level: 30,
			flag: '🇷🇺'
		}
	];

	function getProficiencyColor(proficiency: string): string {
		switch (proficiency) {
			case 'native':
				return 'bg-purple-600/20 text-purple-300 border-purple-500/30';
			case 'fluent':
				return 'bg-blue-600/20 text-blue-300 border-blue-500/30';
			case 'intermediate':
				return 'bg-green-600/20 text-green-300 border-green-500/30';
			case 'beginner':
				return 'bg-yellow-600/20 text-yellow-300 border-yellow-500/30';
			default:
				return 'bg-gray-600/20 text-gray-300 border-gray-500/30';
		}
	}

	function getProficiencyBarColor(proficiency: string): string {
		switch (proficiency) {
			case 'native':
				return 'bg-gradient-to-r from-purple-500 to-purple-600';
			case 'fluent':
				return 'bg-gradient-to-r from-blue-500 to-blue-600';
			case 'intermediate':
				return 'bg-gradient-to-r from-green-500 to-green-600';
			case 'beginner':
				return 'bg-gradient-to-r from-yellow-500 to-yellow-600';
			default:
				return 'bg-gradient-to-r from-gray-500 to-gray-600';
		}
	}
</script>

<section id="languages" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('languages.title') || 'Languages'}
			</h2>
			<p class="text-xl text-gray-300 mb-6">
				{t('languages.subtitle') || 'Languages I speak and their proficiency levels'}
			</p>
		</div>

		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each languages as language, index}
				{@const animationDelay = index * 0.1}
				<div
					class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/20 hover:scale-[1.02] relative overflow-hidden animate-fadeInUp"
					style="animation-delay: {animationDelay}s"
				>
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
					
					<div class="relative z-10">
						<div class="flex items-center gap-4 mb-4">
							<div class="text-4xl">{language.flag}</div>
							<div class="flex-1">
								<h3 class="text-xl font-bold text-white group-hover:text-blue-400 transition-colors mb-1">
									{language.name}
								</h3>
								<span
									class="inline-block px-3 py-1 text-xs rounded-full border capitalize {getProficiencyColor(language.proficiency)}"
								>
									{language.proficiency}
								</span>
							</div>
							{#if language.proficiency === 'native'}
								<CheckCircle class="w-6 h-6 text-purple-400" />
							{/if}
						</div>

						<!-- Proficiency bar -->
						<div class="mb-2">
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm text-gray-400">Proficiency</span>
								<span class="text-sm font-medium text-gray-300">{language.level}%</span>
							</div>
							<div class="w-full h-2 bg-gray-700/50 rounded-full overflow-hidden">
								<div
									class="h-full {getProficiencyBarColor(language.proficiency)} transition-all duration-1000 ease-out"
									style="width: {language.level}%"
								></div>
							</div>
						</div>

						<!-- Language icon -->
						<div class="flex items-center gap-2 pt-4 border-t border-gray-700">
							<Globe class="w-4 h-4 text-gray-400" />
							<span class="text-sm text-gray-400">
								{language.proficiency === 'native'
									? 'Native speaker'
									: language.proficiency === 'fluent'
										? 'Fluent communication'
										: language.proficiency === 'intermediate'
											? 'Conversational level'
											: 'Learning'}
							</span>
						</div>
					</div>
				</div>
			{/each}
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
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}
</style>

