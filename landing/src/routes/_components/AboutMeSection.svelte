<script lang="ts">
	import { downloadResume } from '$lib/api/resumes';
	import { userId } from '$lib/constants';
	import { language } from '$lib/i18n';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();

	// Fallback image SVG as data URL
	const fallbackImage = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Crect fill='%23334155' width='200' height='200'/%3E%3Ctext fill='%239CA3AF' font-family='sans-serif' font-size='18' dy='10.5' font-weight='bold' x='50%25' y='50%25' text-anchor='middle'%3EProfile%3C/text%3E%3C/svg%3E";

	function handleImageError(event: Event) {
		const img = event.target as HTMLImageElement;
		if (img) {
			img.src = fallbackImage;
		}
	}

	async function handleDownloadResume() {
		try {
			// Map language to backend format (en, pt-BR -> en, pt)
			const backendLang = $language === 'pt-BR' ? 'pt' : $language === 'en' ? 'en' : 'en';
			await downloadResume(userId, backendLang);
		} catch (error) {
			console.error('Failed to download resume:', error);
			// You could show a toast notification here
		}
	}
</script>

<section id="about-me" class="container mx-auto px-6 py-20 md:py-32">
	<div class="max-w-6xl mx-auto">
		<div class="flex flex-col md:flex-row items-center gap-12 md:gap-16">
			<!-- Foto -->
			<div class="flex-shrink-0">
				<div class="relative">
					<div
						class="w-64 h-64 md:w-80 md:h-80 rounded-full overflow-hidden border-4 border-blue-500/30 shadow-2xl shadow-blue-500/20 bg-gradient-to-br from-blue-600/20 to-purple-600/20"
					>
						<img
							src="/profile-photo.jpg"
							alt={t('aboutMe.photoAlt')}
							class="w-full h-full object-cover"
							onerror={handleImageError}
						/>
					</div>
					<!-- Decorative elements -->
					<div
						class="absolute -top-4 -right-4 w-24 h-24 bg-blue-500/20 rounded-full blur-2xl -z-10"
					></div>
					<div
						class="absolute -bottom-4 -left-4 w-32 h-32 bg-purple-500/20 rounded-full blur-2xl -z-10"
					></div>
				</div>
			</div>

			<!-- Conteúdo -->
			<div class="flex-1 text-center md:text-left">
				<div
					class="inline-block mb-4 px-4 py-2 bg-blue-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium"
				>
					{t('aboutMe.badge')}
				</div>
				<h1
					class="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent"
				>
					{t('aboutMe.title')}
				</h1>
				<div class="space-y-4 text-lg md:text-xl text-gray-300 leading-relaxed">
					<p>{t('aboutMe.description1')}</p>
					<p>{t('aboutMe.description2')}</p>
					<p class="text-blue-300 font-medium">{t('aboutMe.description3')}</p>
				</div>
				<div class="mt-8 flex flex-wrap gap-4 justify-center md:justify-start">
					<a
						href="#projects"
						class="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
					>
						{t('aboutMe.viewProjects')}
					</a>
					<button
						onclick={handleDownloadResume}
						class="px-6 py-3 bg-green-600 hover:bg-green-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-green-500/50 flex items-center gap-2"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="w-5 h-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
						{t('aboutMe.downloadResume')}
					</button>
					<a
						href="#contact"
						class="px-6 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
					>
						{t('aboutMe.contactMe')}
					</a>
				</div>
			</div>
		</div>
	</div>
</section>

<style>
	@keyframes float {
		0%,
		100% {
			transform: translateY(0px);
		}
		50% {
			transform: translateY(-10px);
		}
	}

	img {
		animation: float 6s ease-in-out infinite;
	}
</style>

