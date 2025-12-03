<script lang="ts">
	import {
		Building2,
		Heart,
		MapPin,
		ExternalLink,
		Linkedin,
		Globe,
		Trophy,
		TrendingUp,
		Users,
		Zap
	} from 'lucide-svelte';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();

	// Hard-coded volunteer work data
	const volunteerWorks = [
		{
			organization: 'Tech for Good',
			role: 'Mentor de Desenvolvimento',
			period: '2021 - Presente',
			location: 'Remoto',
			description:
				'Mentoria de desenvolvedores iniciantes em projetos open source voltados para impacto social. Ajuda na arquitetura de sistemas, code reviews e orientação de carreira.',
			technologies: ['Go', 'JavaScript', 'React', 'Docker'],
			organizationUrl: 'https://techforgood.org',
			linkedinUrl: 'https://linkedin.com/company/techforgood',
			projects: [
				{ name: 'Sistema de Doações', url: 'https://techforgood.org/donations' },
				{ name: 'Plataforma de Voluntariado', url: 'https://techforgood.org/volunteer' }
			],
			achievements: [
				{ metric: '50+', description: 'Desenvolvedores mentorados', icon: Users },
				{ metric: '10+', description: 'Projetos apoiados', icon: Trophy },
				{ metric: '1000+', description: 'Horas voluntárias', icon: TrendingUp }
			]
		},
		{
			organization: 'Code Club Brasil',
			role: 'Instrutor de Programação',
			period: '2019 - 2021',
			location: 'São Paulo, Brasil',
			description:
				'Ensino de programação para crianças e adolescentes em comunidades carentes. Desenvolvimento de currículo e workshops práticos de introdução à programação.',
			technologies: ['Scratch', 'Python', 'JavaScript'],
			organizationUrl: 'https://codeclubbrasil.org',
			linkedinUrl: null,
			projects: [
				{ name: 'Curso de Scratch', url: 'https://codeclubbrasil.org/scratch' },
				{ name: 'Workshop Python', url: 'https://codeclubbrasil.org/python' }
			],
			achievements: [
				{ metric: '200+', description: 'Alunos atendidos', icon: Users },
				{ metric: '500+', description: 'Horas de ensino', icon: TrendingUp }
			]
		},
		{
			organization: 'Open Source Initiative',
			role: 'Contribuidor',
			period: '2020 - Presente',
			location: 'Remoto',
			description:
				'Contribuições regulares para projetos open source focados em educação e tecnologia. Manutenção de bibliotecas, correção de bugs e desenvolvimento de features.',
			technologies: ['Go', 'TypeScript', 'Docker', 'Kubernetes'],
			organizationUrl: null,
			linkedinUrl: null,
			projects: [
				{ name: 'EduTech Library', url: 'https://github.com/example/edutech' },
				{ name: 'DevTools Framework', url: 'https://github.com/example/devtools' }
			],
			achievements: [
				{ metric: '100+', description: 'PRs aprovados', icon: Trophy },
				{ metric: '50+', description: 'Issues resolvidas', icon: Zap }
			]
		}
	];
</script>

<section id="volunteer-work" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4">{t('volunteerWork.title')}</h2>
			<p class="text-xl text-gray-300">{t('volunteerWork.subtitle')}</p>
		</div>

		<div class="relative">
			<!-- Timeline line -->
			<div
				class="absolute left-8 md:left-1/2 top-0 bottom-0 w-0.5 bg-gradient-to-b from-pink-500 via-red-500 to-orange-500 transform md:-translate-x-1/2"
			></div>

			<div class="space-y-12">
				{#each volunteerWorks as work, index}
					<div class="relative flex flex-col md:flex-row items-start gap-8">
						<!-- Timeline dot with icon -->
						<div
							class="absolute left-6 md:left-1/2 w-10 h-10 bg-gray-800 border-4 border-gray-900 rounded-full transform md:-translate-x-1/2 z-10 flex items-center justify-center shadow-lg"
						>
							<Heart class="w-5 h-5 text-pink-400" />
						</div>

						<!-- Content card -->
						<div
							class="flex-1 bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 md:p-8 border border-gray-700 shadow-xl hover:border-pink-500/50 transition-all duration-300 {index % 2 === 0
								? 'md:mr-auto md:pr-12'
								: 'md:ml-auto md:pl-12'}"
						>
							<div class="flex flex-col md:flex-row md:items-center md:justify-between mb-4 gap-2">
								<div class="flex items-start gap-3">
									<div
										class="flex-shrink-0 w-12 h-12 rounded-lg bg-pink-600/20 border border-pink-500/30 flex items-center justify-center"
									>
										<Heart class="w-6 h-6 text-pink-300" />
									</div>
									<div>
										<h3 class="text-2xl font-bold text-white mb-1">{work.role}</h3>
										<div class="flex items-center gap-2 flex-wrap">
											<p class="text-lg text-pink-400 font-medium flex items-center gap-2">
												<Building2 class="w-4 h-4" />
												{#if work.organizationUrl}
													<a
														href={work.organizationUrl}
														target="_blank"
														rel="noopener noreferrer"
														class="hover:text-pink-300 transition-colors duration-200 flex items-center gap-1"
													>
														{work.organization}
														<ExternalLink class="w-3 h-3" />
													</a>
												{:else}
													<span>{work.organization}</span>
												{/if}
											</p>
											{#if work.linkedinUrl}
												<a
													href={work.linkedinUrl}
													target="_blank"
													rel="noopener noreferrer"
													class="text-pink-400 hover:text-pink-300 transition-colors duration-200"
													aria-label="LinkedIn da organização"
												>
													<Linkedin class="w-4 h-4" />
												</a>
											{/if}
										</div>
									</div>
								</div>
								<div class="flex flex-col md:items-end gap-1">
									<span
										class="inline-flex items-center gap-2 px-3 py-1 border rounded-full text-sm font-medium bg-pink-600/20 border-pink-500/30 text-pink-300"
									>
										<Heart class="w-3 h-3" />
										{t('volunteerWork.volunteer')}
									</span>
									<p class="text-sm text-gray-400">{work.period}</p>
								</div>
							</div>

							<div class="flex items-center gap-2 text-gray-400 mb-4">
								<MapPin class="w-4 h-4" />
								<span class="text-sm">{work.location}</span>
							</div>

							<p class="text-gray-300 leading-relaxed mb-6">{work.description}</p>

							{#if work.achievements && work.achievements.length > 0}
								<div class="mb-6 p-4 bg-gradient-to-r from-pink-600/10 via-red-600/10 to-orange-600/10 rounded-xl border border-pink-500/20">
									<p class="text-sm text-gray-400 mb-3 flex items-center gap-2">
										<Trophy class="w-4 h-4 text-yellow-400" />
										{t('volunteerWork.achievements') || 'Conquistas e Impacto'}:
									</p>
									<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
										{#each work.achievements as achievement}
											{@const AchievementIcon = achievement.icon}
											<div
												class="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg border border-gray-700/50"
											>
												<div
													class="flex-shrink-0 w-10 h-10 rounded-lg bg-pink-600/20 border border-pink-500/30 flex items-center justify-center"
												>
													<AchievementIcon class="w-5 h-5 text-pink-400" />
												</div>
												<div>
													<p class="text-2xl font-bold text-white mb-0.5">{achievement.metric}</p>
													<p class="text-sm text-gray-400">{achievement.description}</p>
												</div>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							<div class="flex flex-wrap gap-2 mb-6">
								{#each work.technologies as tech}
									<span
										class="px-3 py-1 bg-gray-700/50 border border-gray-600 rounded-lg text-sm text-gray-300"
									>
										{tech}
									</span>
								{/each}
							</div>

							{#if work.projects && work.projects.length > 0}
								<div class="pt-4 border-t border-gray-700">
									<p class="text-sm text-gray-400 mb-3 flex items-center gap-2">
										<Globe class="w-4 h-4" />
										{t('volunteerWork.projects') || 'Projetos Relacionados'}:
									</p>
									<div class="flex flex-wrap gap-2">
										{#each work.projects as project}
											<a
												href={project.url}
												target="_blank"
												rel="noopener noreferrer"
												class="px-3 py-1 bg-pink-600/10 border border-pink-500/30 rounded-lg text-sm text-pink-300 hover:bg-pink-600/20 hover:border-pink-500/50 transition-colors duration-200 flex items-center gap-1"
											>
												{project.name}
												<ExternalLink class="w-3 h-3" />
											</a>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
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

	:global(#volunteer-work > div > div > div) {
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}

	:global(#volunteer-work > div > div > div:nth-child(1)) {
		animation-delay: 0.1s;
	}

	:global(#volunteer-work > div > div > div:nth-child(2)) {
		animation-delay: 0.2s;
	}

	:global(#volunteer-work > div > div > div:nth-child(3)) {
		animation-delay: 0.3s;
	}
</style>

