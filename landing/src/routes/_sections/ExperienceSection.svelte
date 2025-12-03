<script lang="ts">
	import {
		Building2,
		Briefcase,
		Code,
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

	// Helper function to get icon for work type
	function getWorkTypeIcon(type: string) {
		switch (type) {
			case 'full-time':
				return Briefcase;
			case 'freelance':
				return Code;
			case 'contract':
				return Building2;
			default:
				return Briefcase;
		}
	}

	// Helper function to get badge color classes for work type
	function getWorkTypeBadgeClasses(type: string) {
		switch (type) {
			case 'full-time':
				return 'bg-blue-600/20 border-blue-500/30 text-blue-300';
			case 'freelance':
				return 'bg-purple-600/20 border-purple-500/30 text-purple-300';
			case 'contract':
				return 'bg-green-600/20 border-green-500/30 text-green-300';
			default:
				return 'bg-gray-600/20 border-gray-500/30 text-gray-300';
		}
	}

	// Hard-coded experience data
	const experiences = [
		{
			company: 'Tech Solutions Inc.',
			position: 'Senior Backend Developer',
			period: '2022 - Presente',
			location: 'São Paulo, Brasil',
			description:
				'Desenvolvimento de sistemas backend escaláveis usando Go, Docker e Kubernetes. Liderança técnica em projetos de microsserviços, implementação de arquiteturas cloud-native e otimização de performance.',
			technologies: ['Go', 'Docker', 'Kubernetes', 'PostgreSQL', 'Redis', 'AWS'],
			type: 'full-time',
			companyUrl: 'https://example.com',
			linkedinUrl: 'https://linkedin.com/company/tech-solutions',
			projects: [
				{ name: 'Sistema de Microsserviços', url: 'https://example.com/project1' },
				{ name: 'Plataforma Cloud-Native', url: 'https://example.com/project2' }
			],
			achievements: [
				{ metric: '40%', description: 'Aumento de performance', icon: TrendingUp },
				{ metric: '5+', description: 'Projetos liderados', icon: Trophy },
				{ metric: '99.9%', description: 'Uptime alcançado', icon: Zap }
			]
		},
		{
			company: 'StartupXYZ',
			position: 'Backend Developer',
			period: '2020 - 2022',
			location: 'Remoto',
			description:
				'Desenvolvimento de APIs RESTful e GraphQL, implementação de sistemas de autenticação e autorização, integração com serviços de terceiros e otimização de queries de banco de dados.',
			technologies: ['Node.js', 'TypeScript', 'PostgreSQL', 'MongoDB', 'GraphQL', 'Docker'],
			type: 'full-time',
			companyUrl: 'https://startupxyz.com',
			linkedinUrl: 'https://linkedin.com/company/startupxyz',
			projects: [
				{ name: 'API Gateway', url: 'https://startupxyz.com/api-gateway' },
				{ name: 'Sistema de Autenticação', url: 'https://startupxyz.com/auth' }
			],
			achievements: [
				{ metric: '3x', description: 'Redução de latência', icon: Zap },
				{ metric: '50K+', description: 'Usuários atendidos', icon: Users },
				{ metric: '2x', description: 'Escalabilidade melhorada', icon: TrendingUp }
			]
		},
		{
			company: 'Freelancer',
			position: 'Desenvolvedor Full Stack',
			period: '2018 - 2020',
			location: 'Remoto',
			description:
				'Desenvolvimento de aplicações web completas para diversos clientes, desde landing pages até sistemas complexos. Trabalho com múltiplas tecnologias e frameworks.',
			technologies: ['JavaScript', 'React', 'Node.js', 'Python', 'MySQL', 'AWS'],
			type: 'freelance',
			companyUrl: null,
			linkedinUrl: null,
			projects: [
				{ name: 'E-commerce Platform', url: 'https://example.com/ecommerce' },
				{ name: 'SaaS Dashboard', url: 'https://example.com/saas' }
			],
			achievements: [
				{ metric: '20+', description: 'Projetos entregues', icon: Trophy },
				{ metric: '100%', description: 'Satisfação do cliente', icon: Users }
			]
		}
	];
</script>

<section id="experience" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4">{t('experience.title')}</h2>
			<p class="text-xl text-gray-300">{t('experience.subtitle')}</p>
		</div>

		<div class="relative">
			<!-- Timeline line -->
			<div
				class="absolute left-8 md:left-1/2 top-0 bottom-0 w-0.5 bg-gradient-to-b from-blue-500 via-purple-500 to-pink-500 transform md:-translate-x-1/2"
			></div>

			<div class="space-y-12">
				{#each experiences as experience, index}
					{@const WorkIcon = getWorkTypeIcon(experience.type)}
					<div class="relative flex flex-col md:flex-row items-start gap-8">
						<!-- Timeline dot with icon -->
						<div
							class="absolute left-6 md:left-1/2 w-10 h-10 bg-gray-800 border-4 border-gray-900 rounded-full transform md:-translate-x-1/2 z-10 flex items-center justify-center shadow-lg"
						>
							<WorkIcon class="w-5 h-5 text-blue-400" />
						</div>

						<!-- Content card -->
						<div
							class="flex-1 bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 md:p-8 border border-gray-700 shadow-xl hover:border-blue-500/50 transition-all duration-300 {index % 2 === 0
								? 'md:mr-auto md:pr-12'
								: 'md:ml-auto md:pl-12'}"
						>
							<div class="flex flex-col md:flex-row md:items-center md:justify-between mb-4 gap-2">
								<div class="flex items-start gap-3">
									<div
										class="flex-shrink-0 w-12 h-12 rounded-lg {getWorkTypeBadgeClasses(experience.type)} flex items-center justify-center"
									>
										<WorkIcon class="w-6 h-6" />
									</div>
									<div>
										<h3 class="text-2xl font-bold text-white mb-1">{experience.position}</h3>
										<div class="flex items-center gap-2 flex-wrap">
											<p class="text-lg text-blue-400 font-medium flex items-center gap-2">
												<Building2 class="w-4 h-4" />
												{#if experience.companyUrl}
													<a
														href={experience.companyUrl}
														target="_blank"
														rel="noopener noreferrer"
														class="hover:text-blue-300 transition-colors duration-200 flex items-center gap-1"
													>
														{experience.company}
														<ExternalLink class="w-3 h-3" />
													</a>
												{:else}
													<span>{experience.company}</span>
												{/if}
											</p>
											{#if experience.linkedinUrl}
												<a
													href={experience.linkedinUrl}
													target="_blank"
													rel="noopener noreferrer"
													class="text-blue-400 hover:text-blue-300 transition-colors duration-200"
													aria-label="LinkedIn da empresa"
												>
													<Linkedin class="w-4 h-4" />
												</a>
											{/if}
										</div>
									</div>
								</div>
								<div class="flex flex-col md:items-end gap-1">
									<span
										class="inline-flex items-center gap-2 px-3 py-1 border rounded-full text-sm font-medium {getWorkTypeBadgeClasses(experience.type)}"
									>
										<WorkIcon class="w-3 h-3" />
										{experience.type === 'full-time'
											? t('experience.fullTime')
											: experience.type === 'freelance'
												? t('experience.freelance')
												: t('experience.contract')}
									</span>
									<p class="text-sm text-gray-400">{experience.period}</p>
								</div>
							</div>

							<div class="flex items-center gap-2 text-gray-400 mb-4">
								<MapPin class="w-4 h-4" />
								<span class="text-sm">{experience.location}</span>
							</div>

							<p class="text-gray-300 leading-relaxed mb-6">{experience.description}</p>

							{#if experience.achievements && experience.achievements.length > 0}
								<div class="mb-6 p-4 bg-gradient-to-r from-blue-600/10 via-purple-600/10 to-pink-600/10 rounded-xl border border-blue-500/20">
									<p class="text-sm text-gray-400 mb-3 flex items-center gap-2">
										<Trophy class="w-4 h-4 text-yellow-400" />
										{t('experience.achievements') || 'Conquistas e Resultados'}:
									</p>
									<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
										{#each experience.achievements as achievement}
											{@const AchievementIcon = achievement.icon}
											<div
												class="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg border border-gray-700/50"
											>
												<div
													class="flex-shrink-0 w-10 h-10 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center"
												>
													<AchievementIcon class="w-5 h-5 text-blue-400" />
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
								{#each experience.technologies as tech}
									<span
										class="px-3 py-1 bg-gray-700/50 border border-gray-600 rounded-lg text-sm text-gray-300"
									>
										{tech}
									</span>
								{/each}
							</div>

							{#if experience.projects && experience.projects.length > 0}
								<div class="pt-4 border-t border-gray-700">
									<p class="text-sm text-gray-400 mb-3 flex items-center gap-2">
										<Globe class="w-4 h-4" />
										{t('experience.projects') || 'Projetos Relacionados'}:
									</p>
									<div class="flex flex-wrap gap-2">
										{#each experience.projects as project}
											<a
												href={project.url}
												target="_blank"
												rel="noopener noreferrer"
												class="px-3 py-1 bg-blue-600/10 border border-blue-500/30 rounded-lg text-sm text-blue-300 hover:bg-blue-600/20 hover:border-blue-500/50 transition-colors duration-200 flex items-center gap-1"
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

	:global(#experience > div > div > div) {
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}

	:global(#experience > div > div > div:nth-child(1)) {
		animation-delay: 0.1s;
	}

	:global(#experience > div > div > div:nth-child(2)) {
		animation-delay: 0.2s;
	}

	:global(#experience > div > div > div:nth-child(3)) {
		animation-delay: 0.3s;
	}
</style>

