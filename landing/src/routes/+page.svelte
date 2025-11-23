<script lang="ts">
	import { onMount } from 'svelte';
	import { Icon } from 'svelte-icons-pack';
	import {
		SiGo,
		SiDocker,
		SiKubernetes,
		SiRedis,
		SiGithub,
		SiInstagram,
		SiLinkedin,
		SiWhatsapp,
		SiPython
	} from 'svelte-icons-pack/si';
	import { Mail, Phone, MapPin, ExternalLink, Calendar, TrendingUp, Code2, Settings, Brain, GitBranch, Tag, Layers, Zap, Target, CheckCircle2, XCircle, ArrowRight, ChevronDown, ChevronUp } from 'lucide-svelte';
	import { contact, skills, interests } from '$lib/constants';
	import { caseStudies, systemDesigns, problemSolutions } from '$lib/constants/technical';
	import { listProjects } from '$lib/api/projects';
	import { listSkillsWithCounts, type SkillWithCount } from '$lib/api/skills';
	import type { Project, ProjectTechnology } from '$lib/types/project';
	import TestimonialsCarousel from '$lib/components/TestimonialsCarousel.svelte';
	import BlogPostsSection from '$lib/components/BlogPostsSection.svelte';
	import ProjectsShowcase from '$lib/components/ProjectsShowcase.svelte';
	import LanguageSwitcher from '$lib/components/LanguageSwitcher.svelte';
	import { language, translationsStore } from '$lib/i18n';

	// Reactive translation helper
	let t = $derived($translationsStore);

	// Generate WhatsApp URL with pre-filled message
	function getWhatsAppUrl(): string {
		if (!contact.whatsapp) return '#';
		const phoneNumber = contact.whatsapp.replace(/\D/g, '');
		const messages: Record<string, string> = {
			'en': 'Hello! I came across your website (www.woragis.me) and would like to get in touch.',
			'pt-BR': 'Olá! Encontrei seu site (www.woragis.me) e gostaria de entrar em contato.',
			'fr': 'Bonjour! J\'ai trouvé votre site (www.woragis.me) et j\'aimerais entrer en contact.',
			'es': '¡Hola! Encontré tu sitio web (www.woragis.me) y me gustaría ponerme en contacto.'
		};
		const message = encodeURIComponent(messages[$language] || messages['en']);
		return `https://wa.me/${phoneNumber}?text=${message}`;
	}

	// Featured projects
	let featuredProjects: Project[] = $state([]);
	let loadingProjects = $state(false);

	// Skills
	let popularSkills: SkillWithCount[] = $state([]);
	let loadingSkills = $state(false);

	onMount(async () => {
		await Promise.all([fetchFeaturedProjects(), fetchPopularSkills()]);
	});

	async function fetchFeaturedProjects() {
		loadingProjects = true;
		try {
			const projects = await listProjects({ limit: 6, sortBy: 'updatedAt', sortOrder: 'desc' });
			featuredProjects = projects.slice(0, 6);
		} catch (error) {
			console.error('Error fetching featured projects:', error);
		} finally {
			loadingProjects = false;
		}
	}

	async function fetchPopularSkills() {
		loadingSkills = true;
		try {
			const allSkills = await listSkillsWithCounts();
			// Get top 6 skills by project count
			popularSkills = allSkills
				.filter((skill) => skill.projectCount > 0)
				.sort((a, b) => b.projectCount - a.projectCount)
				.slice(0, 6);
		} catch (error) {
			console.error('Error fetching skills:', error);
		} finally {
			loadingSkills = false;
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getTechIcon(techName: string) {
		const name = techName.toLowerCase();
		if (name.includes('go') || name === 'golang') return SiGo;
		if (name.includes('docker')) return SiDocker;
		if (name.includes('kubernetes') || name.includes('k8s')) return SiKubernetes;
		if (name.includes('redis')) return SiRedis;
		if (name.includes('python')) return SiPython;
		return null;
	}

	function getStatusColor(status: string): string {
		const colors: Record<string, string> = {
			idea: 'bg-purple-600',
			planning: 'bg-blue-600',
			executing: 'bg-yellow-600',
			monitoring: 'bg-green-600',
			completed: 'bg-cyan-600'
		};
		return colors[status] || 'bg-gray-600';
	}

	// Technical sections state
	let expandedCaseStudy: string | null = $state(null);
	let expandedSystemDesign: string | null = $state(null);
	let expandedProblem: string | null = $state(null);

	function toggleCaseStudy(id: string) {
		expandedCaseStudy = expandedCaseStudy === id ? null : id;
	}

	function toggleSystemDesign(id: string) {
		expandedSystemDesign = expandedSystemDesign === id ? null : id;
	}

	function toggleProblem(id: string) {
		expandedProblem = expandedProblem === id ? null : id;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Language Switcher (Fixed Position) -->
	<div class="fixed top-4 right-4 z-50">
		<LanguageSwitcher />
	</div>

	<!-- Hero Section -->
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

	<!-- About Section -->
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

	<!-- Projects Section -->
	<section id="projects" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="flex items-center justify-between mb-12">
				<h2 class="text-4xl font-bold">{t('projects.title')}</h2>
				<a
					href="/projects"
					class="text-blue-400 hover:text-blue-300 transition-colors duration-200 flex items-center gap-2"
				>
					{t('projects.viewAll')}
					<ExternalLink class="w-5 h-5" />
				</a>
			</div>
			<ProjectsShowcase />
		</div>
	</section>

	<!-- Blog Posts Section -->
	<section id="blog" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
					{t('blog.title')}
				</h2>
				<p class="text-gray-400 text-lg max-w-2xl mx-auto">
					{t('blog.subtitle')}
				</p>
			</div>
			<BlogPostsSection />
		</div>
	</section>

	<!-- Skills Section -->
	<section id="skills" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center justify-between mb-12">
				<h2 class="text-4xl font-bold text-center flex-1">Popular Skills</h2>
				<a
					href="/skills"
					class="text-blue-400 hover:text-blue-300 transition-colors duration-200 flex items-center gap-2"
				>
					View All Skills
					<ExternalLink class="w-5 h-5" />
				</a>
			</div>

			{#if loadingSkills}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if popularSkills.length === 0}
				<div class="text-center py-20">
					<Tag class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No skills found</p>
					<p class="text-gray-500 text-sm">Skills will appear here once they're added to projects</p>
				</div>
			{:else}
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each popularSkills as skill}
						<a
							href="/skills"
							class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 group"
						>
							<div class="flex items-start justify-between mb-4">
								<div class="flex-1">
									<h3
										class="text-xl font-bold text-white mb-2 group-hover:text-blue-400 transition-colors"
									>
										{skill.name}
									</h3>
									<span
										class="inline-block px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 capitalize"
									>
										{skill.category}
									</span>
								</div>
								{#if skill.icon}
									<div class="w-10 h-10 bg-gray-700 rounded flex items-center justify-center">
										<span class="text-lg">{skill.icon}</span>
									</div>
								{/if}
							</div>

							{#if skill.description}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{skill.description}</p>
							{/if}

							<div class="flex items-center justify-between pt-4 border-t border-gray-700">
								<div class="flex items-center gap-2">
									<Tag class="w-4 h-4 text-gray-400" />
									<span class="text-sm text-gray-400">Projects</span>
								</div>
								<span class="text-lg font-bold text-blue-400">{skill.projectCount}</span>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	</section>

	<!-- Technical Depth Section -->
	<section id="technical-depth" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h2 class="text-4xl font-bold mb-4">Technical Depth</h2>
				<p class="text-gray-400 text-lg max-w-2xl mx-auto">
					Deep dive into system architectures, design decisions, and technical implementations
				</p>
			</div>

			<div class="grid md:grid-cols-2 gap-6 mb-12">
				{#each systemDesigns as design}
					<div
						class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300"
					>
						<div class="flex items-start justify-between mb-4">
							<div class="flex items-center gap-3">
								<div
									class="w-12 h-12 bg-gradient-to-br from-blue-600 to-cyan-600 rounded-lg flex items-center justify-center"
								>
									<Layers class="w-6 h-6" />
								</div>
								<div>
									<h3 class="text-xl font-bold text-white">{design.title}</h3>
									<p class="text-sm text-gray-400">System Design</p>
								</div>
							</div>
							<button
								onclick={() => toggleSystemDesign(design.id)}
								class="text-gray-400 hover:text-white transition-colors"
							>
								{#if expandedSystemDesign === design.id}
									<ChevronUp class="w-5 h-5" />
								{:else}
									<ChevronDown class="w-5 h-5" />
								{/if}
							</button>
						</div>

						<p class="text-gray-300 mb-4">{design.description}</p>

						{#if expandedSystemDesign === design.id}
							<div class="mt-4 space-y-4 pt-4 border-t border-gray-700">
								<div>
									<h4 class="text-sm font-semibold text-blue-400 mb-2">Components</h4>
									<div class="space-y-2">
										{#each design.components as component}
											<div class="bg-gray-800/50 rounded-lg p-3">
												<div class="flex items-center justify-between mb-1">
													<span class="font-medium text-white">{component.name}</span>
													<span class="text-xs text-cyan-400">{component.technology}</span>
												</div>
												<p class="text-sm text-gray-400">{component.description}</p>
											</div>
										{/each}
									</div>
								</div>

								{#if design.dataFlow}
									<div>
										<h4 class="text-sm font-semibold text-blue-400 mb-2">Data Flow</h4>
										<p class="text-sm text-gray-300">{design.dataFlow}</p>
									</div>
								{/if}

								{#if design.scalability}
									<div>
										<h4 class="text-sm font-semibold text-green-400 mb-2">Scalability</h4>
										<p class="text-sm text-gray-300">{design.scalability}</p>
									</div>
								{/if}

								{#if design.reliability}
									<div>
										<h4 class="text-sm font-semibold text-purple-400 mb-2">Reliability</h4>
										<p class="text-sm text-gray-300">{design.reliability}</p>
									</div>
								{/if}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Problem Solving Section -->
	<section id="problem-solving" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h2 class="text-4xl font-bold mb-4">Problem Solving & Communication</h2>
				<p class="text-gray-400 text-lg max-w-2xl mx-auto">
					Real challenges solved with clear communication of solutions and impact
				</p>
			</div>

			<div class="space-y-6">
				{#each problemSolutions as solution}
					<div
						class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-yellow-500/50 transition-all duration-300"
					>
						<div class="flex items-start justify-between mb-4">
							<div class="flex-1">
								<div class="flex items-center gap-3 mb-2">
									<div
										class="w-10 h-10 bg-gradient-to-br from-yellow-600 to-orange-600 rounded-lg flex items-center justify-center"
									>
										<Zap class="w-5 h-5" />
									</div>
									<h3 class="text-xl font-bold text-white">{solution.problem}</h3>
								</div>
								<p class="text-sm text-gray-400 mb-3">{solution.context}</p>
							</div>
							<button
								onclick={() => toggleProblem(solution.id)}
								class="text-gray-400 hover:text-white transition-colors ml-4"
							>
								{#if expandedProblem === solution.id}
									<ChevronUp class="w-5 h-5" />
								{:else}
									<ChevronDown class="w-5 h-5" />
								{/if}
							</button>
						</div>

						<div class="grid md:grid-cols-2 gap-4 mb-4">
							<div class="bg-gray-800/30 rounded-lg p-4">
								<h4 class="text-sm font-semibold text-red-400 mb-2">Problem</h4>
								<p class="text-sm text-gray-300">{solution.problem}</p>
							</div>
							<div class="bg-gray-800/30 rounded-lg p-4">
								<h4 class="text-sm font-semibold text-green-400 mb-2">Solution</h4>
								<p class="text-sm text-gray-300">{solution.solution}</p>
							</div>
						</div>

						{#if expandedProblem === solution.id}
							<div class="mt-4 pt-4 border-t border-gray-700 space-y-4">
								<div>
									<h4 class="text-sm font-semibold text-blue-400 mb-2">Technologies Used</h4>
									<div class="flex flex-wrap gap-2">
										{#each solution.technologies as tech}
											<span
												class="px-3 py-1 bg-blue-600/20 text-blue-300 rounded-full text-xs border border-blue-600/30"
											>
												{tech}
											</span>
										{/each}
									</div>
								</div>

								<div>
									<h4 class="text-sm font-semibold text-purple-400 mb-2">Impact</h4>
									<p class="text-sm text-gray-300">{solution.impact}</p>
								</div>

								{#if solution.metrics}
									<div class="bg-gradient-to-r from-gray-800/50 to-gray-900/50 rounded-lg p-4">
										<h4 class="text-sm font-semibold text-cyan-400 mb-3">Metrics</h4>
										<div class="grid grid-cols-3 gap-4">
											<div>
												<p class="text-xs text-gray-400 mb-1">Before</p>
												<p class="text-sm font-medium text-red-300">{solution.metrics.before}</p>
											</div>
											<div>
												<p class="text-xs text-gray-400 mb-1">After</p>
												<p class="text-sm font-medium text-green-300">{solution.metrics.after}</p>
											</div>
											<div>
												<p class="text-xs text-gray-400 mb-1">Improvement</p>
												<p class="text-sm font-medium text-cyan-300">
													{solution.metrics.improvement}
												</p>
											</div>
										</div>
									</div>
								{/if}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- System Thinking Section -->
	<section id="system-thinking" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h2 class="text-4xl font-bold mb-4">System Thinking</h2>
				<p class="text-gray-400 text-lg max-w-2xl mx-auto">
					Architectural decisions, trade-offs, and lessons learned from building complex systems
				</p>
			</div>

			<div class="space-y-8">
				{#each caseStudies as study}
					<div
						class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-8 border border-gray-700 hover:border-purple-500/50 transition-all duration-300"
					>
						<div class="flex items-start justify-between mb-6">
							<div class="flex-1">
								<div class="flex items-center gap-3 mb-3">
									<div
										class="w-12 h-12 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center"
									>
										<Target class="w-6 h-6" />
									</div>
									<div>
										<h3 class="text-2xl font-bold text-white">{study.title}</h3>
										<p class="text-sm text-gray-400">Case Study</p>
									</div>
								</div>
								<p class="text-gray-300 mb-4">{study.description}</p>
							</div>
							<button
								onclick={() => toggleCaseStudy(study.id)}
								class="text-gray-400 hover:text-white transition-colors ml-4"
							>
								{#if expandedCaseStudy === study.id}
									<ChevronUp class="w-6 h-6" />
								{:else}
									<ChevronDown class="w-6 h-6" />
								{/if}
							</button>
						</div>

						<div class="grid md:grid-cols-2 gap-4 mb-6">
							<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-4">
								<h4 class="text-sm font-semibold text-red-400 mb-2 flex items-center gap-2">
									<XCircle class="w-4 h-4" />
									Challenge
								</h4>
								<p class="text-sm text-gray-300">{study.challenge}</p>
							</div>
							<div class="bg-green-900/20 border border-green-700/30 rounded-lg p-4">
								<h4 class="text-sm font-semibold text-green-400 mb-2 flex items-center gap-2">
									<CheckCircle2 class="w-4 h-4" />
									Solution
								</h4>
								<p class="text-sm text-gray-300">{study.solution}</p>
							</div>
						</div>

						<div class="flex flex-wrap gap-2 mb-6">
							{#each study.technologies as tech}
								<span
									class="px-3 py-1 bg-purple-600/20 text-purple-300 rounded-full text-xs border border-purple-600/30"
								>
									{tech}
								</span>
							{/each}
						</div>

						{#if expandedCaseStudy === study.id}
							<div class="mt-6 pt-6 border-t border-gray-700 space-y-6">
								{#if study.architecture}
									<div>
										<h4 class="text-lg font-semibold text-blue-400 mb-3 flex items-center gap-2">
											<Layers class="w-5 h-5" />
											Architecture
										</h4>
										<p class="text-gray-300 leading-relaxed">{study.architecture}</p>
									</div>
								{/if}

								{#if study.metrics && study.metrics.length > 0}
									<div>
										<h4 class="text-lg font-semibold text-cyan-400 mb-3">Key Metrics</h4>
										<div class="grid md:grid-cols-3 gap-4">
											{#each study.metrics as metric}
												<div
													class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 rounded-lg p-4 border border-gray-700"
												>
													<p class="text-xs text-gray-400 mb-1">{metric.label}</p>
													<p class="text-2xl font-bold text-cyan-400">{metric.value}</p>
													{#if metric.improvement}
														<p class="text-xs text-gray-500 mt-1">{metric.improvement}</p>
													{/if}
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if study.tradeoffs && study.tradeoffs.length > 0}
									<div>
										<h4 class="text-lg font-semibold text-yellow-400 mb-3">Trade-offs & Decisions</h4>
										<div class="space-y-4">
											{#each study.tradeoffs as tradeoff}
												<div
													class="bg-gray-800/30 rounded-lg p-4 border border-gray-700"
												>
													<h5 class="font-semibold text-white mb-3">{tradeoff.decision}</h5>
													<div class="grid md:grid-cols-2 gap-4">
														<div>
															<p class="text-xs font-semibold text-green-400 mb-2">Pros</p>
															<ul class="space-y-1">
																{#each tradeoff.pros as pro}
																	<li class="text-sm text-gray-300 flex items-start gap-2">
																		<CheckCircle2 class="w-4 h-4 text-green-400 mt-0.5 flex-shrink-0" />
																		<span>{pro}</span>
																	</li>
																{/each}
															</ul>
														</div>
														<div>
															<p class="text-xs font-semibold text-red-400 mb-2">Cons</p>
															<ul class="space-y-1">
																{#each tradeoff.cons as con}
																	<li class="text-sm text-gray-300 flex items-start gap-2">
																		<XCircle class="w-4 h-4 text-red-400 mt-0.5 flex-shrink-0" />
																		<span>{con}</span>
																	</li>
																{/each}
															</ul>
														</div>
													</div>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if study.lessonsLearned && study.lessonsLearned.length > 0}
									<div>
										<h4 class="text-lg font-semibold text-purple-400 mb-3">Lessons Learned</h4>
										<ul class="space-y-2">
											{#each study.lessonsLearned as lesson}
												<li class="text-gray-300 flex items-start gap-3">
													<ArrowRight class="w-5 h-5 text-purple-400 mt-0.5 flex-shrink-0" />
													<span>{lesson}</span>
												</li>
											{/each}
										</ul>
									</div>
								{/if}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Interests & Technologies Section -->
	<section id="interests" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">{t('interests.title')}</h2>
			<div class="grid md:grid-cols-2 gap-6">
				<!-- AI & RAG -->
				<div
					class="bg-gradient-to-br from-pink-900/30 to-purple-900/20 rounded-xl p-8 border border-pink-700/30 hover:border-pink-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-pink-500/20"
				>
					<div class="flex items-center mb-6">
						<div
							class="w-14 h-14 bg-gradient-to-br from-pink-600 to-purple-600 rounded-lg flex items-center justify-center mr-4"
						>
							<Brain class="w-8 h-8" stroke="white" />
						</div>
						<h3
							class="text-2xl font-bold bg-gradient-to-r from-pink-400 to-purple-400 bg-clip-text text-transparent"
						>
							{interests[0].title}
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed">{interests[0].description}</p>
				</div>

				<!-- Redis & Pub/Sub -->
				<div
					class="bg-gradient-to-br from-red-900/30 to-orange-900/20 rounded-xl p-8 border border-red-700/30 hover:border-red-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-red-500/20"
				>
					<div class="flex items-center mb-6">
						<div
							class="w-14 h-14 bg-gradient-to-br from-red-600 to-orange-600 rounded-lg flex items-center justify-center mr-4"
						>
							<Icon src={SiRedis} size="2rem" color="white" />
						</div>
						<h3
							class="text-2xl font-bold bg-gradient-to-r from-red-400 to-orange-400 bg-clip-text text-transparent"
						>
							{interests[1].title}
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed">{interests[1].description}</p>
				</div>

				<!-- Distributed Architecture -->
				<div
					class="bg-gradient-to-br from-green-900/30 to-emerald-900/20 rounded-xl p-8 border border-green-700/30 hover:border-green-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-green-500/20 md:col-span-2"
				>
					<div class="flex items-center mb-6">
						<div
							class="w-14 h-14 bg-gradient-to-br from-green-600 to-emerald-600 rounded-lg flex items-center justify-center mr-4"
						>
							<GitBranch class="w-8 h-8" stroke="white" />
						</div>
						<h3
							class="text-2xl font-bold bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent"
						>
							{interests[2].title}
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed">{interests[2].description}</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Testimonials Section -->
	<section id="testimonials" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<div class="text-center mb-12">
				<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
					{t('testimonials.title')}
				</h2>
				<p class="text-gray-400 text-lg max-w-2xl mx-auto">
					{t('testimonials.subtitle')}
				</p>
			</div>
			<TestimonialsCarousel />
		</div>
	</section>

	<!-- Contact Section -->
	<section id="contact" class="container mx-auto px-6 py-20 relative overflow-hidden">
		<!-- Background decorative elements -->
		<div
			class="absolute inset-0 opacity-10 pointer-events-none"
			style="background-image: radial-gradient(circle at 20% 50%, rgba(59, 130, 246, 0.5) 0%, transparent 50%), radial-gradient(circle at 80% 80%, rgba(147, 51, 234, 0.5) 0%, transparent 50%);"
		></div>

		<div class="max-w-6xl mx-auto relative z-10">
			<!-- Header -->
			<div class="text-center mb-16">
				<div
					class="inline-block mb-4 px-4 py-2 bg-gradient-to-r from-blue-600/20 to-purple-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium"
				>
					{t('contact.title')}
				</div>
				<h2 class="text-5xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
					{t('contact.title')}
				</h2>
				<p class="text-xl text-gray-300 max-w-2xl mx-auto leading-relaxed">
					{t('contact.subtitle')}
				</p>
			</div>

			<!-- Contact Cards Grid -->
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
				<!-- Email Card -->
				<a
					href="mailto:{contact.email}"
					class="group relative bg-gradient-to-br from-blue-900/40 via-blue-800/30 to-blue-900/40 backdrop-blur-sm rounded-2xl p-6 border border-blue-700/30 hover:border-blue-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/30 hover:-translate-y-1 overflow-hidden"
				>
					<!-- Animated background gradient -->
					<div
						class="absolute inset-0 bg-gradient-to-br from-blue-600/0 via-blue-600/0 to-blue-600/0 group-hover:from-blue-600/10 group-hover:via-blue-600/5 group-hover:to-blue-600/10 transition-all duration-500"
					></div>

					<div class="relative z-10">
						<div
							class="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-blue-500/30"
						>
							<Mail class="w-8 h-8 text-white" />
						</div>
						<h3 class="text-xl font-bold text-white mb-2 group-hover:text-blue-300 transition-colors">
							{t('contact.email')}
						</h3>
						<p class="text-gray-300 text-sm break-all group-hover:text-white transition-colors">
							{contact.email}
						</p>
						<p class="text-blue-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
							{t('contact.clickToEmail')}
						</p>
					</div>
				</a>

				<!-- Phone Card -->
				{#if contact.phone}
					<a
						href="tel:{contact.phone}"
						class="group relative bg-gradient-to-br from-green-900/40 via-green-800/30 to-green-900/40 backdrop-blur-sm rounded-2xl p-6 border border-green-700/30 hover:border-green-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-green-500/30 hover:-translate-y-1 overflow-hidden"
					>
						<div
							class="absolute inset-0 bg-gradient-to-br from-green-600/0 via-green-600/0 to-green-600/0 group-hover:from-green-600/10 group-hover:via-green-600/5 group-hover:to-green-600/10 transition-all duration-500"
						></div>

						<div class="relative z-10">
							<div
								class="w-16 h-16 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-green-500/30"
							>
								<Phone class="w-8 h-8 text-white" />
							</div>
							<h3 class="text-xl font-bold text-white mb-2 group-hover:text-green-300 transition-colors">
								{t('contact.phone')}
							</h3>
							<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
								{contact.phone}
							</p>
							<p class="text-green-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
								{t('contact.clickToCall')}
							</p>
						</div>
					</a>
				{/if}

				<!-- Location Card -->
				{#if contact.location}
					<div
						class="group relative bg-gradient-to-br from-purple-900/40 via-purple-800/30 to-purple-900/40 backdrop-blur-sm rounded-2xl p-6 border border-purple-700/30 overflow-hidden"
					>
						<div class="relative z-10">
							<div
								class="w-16 h-16 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-purple-500/30"
							>
								<MapPin class="w-8 h-8 text-white" />
							</div>
							<h3 class="text-xl font-bold text-white mb-2">{t('contact.location')}</h3>
							<p class="text-gray-300 text-sm">{contact.location}</p>
						</div>
					</div>
				{/if}

				<!-- LinkedIn Card -->
				{#if contact.linkedin}
					<a
						href={contact.linkedin}
						target="_blank"
						rel="noopener noreferrer"
						class="group relative bg-gradient-to-br from-indigo-900/40 via-indigo-800/30 to-indigo-900/40 backdrop-blur-sm rounded-2xl p-6 border border-indigo-700/30 hover:border-indigo-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-indigo-500/30 hover:-translate-y-1 overflow-hidden"
					>
						<div
							class="absolute inset-0 bg-gradient-to-br from-indigo-600/0 via-indigo-600/0 to-indigo-600/0 group-hover:from-indigo-600/10 group-hover:via-indigo-600/5 group-hover:to-indigo-600/10 transition-all duration-500"
						></div>

						<div class="relative z-10">
							<div
								class="w-16 h-16 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-indigo-500/30"
							>
								<Icon src={SiLinkedin} size="2rem" color="white" />
							</div>
							<h3 class="text-xl font-bold text-white mb-2 group-hover:text-indigo-300 transition-colors">
								{t('contact.linkedin')}
							</h3>
							<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
								{t('contact.subtitle')}
							</p>
							<p class="text-indigo-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
								View profile →
							</p>
						</div>
					</a>
				{/if}

				<!-- WhatsApp Card -->
				{#if contact.whatsapp}
					<a
						href={getWhatsAppUrl()}
						target="_blank"
						rel="noopener noreferrer"
						class="group relative bg-gradient-to-br from-emerald-900/40 via-emerald-800/30 to-emerald-900/40 backdrop-blur-sm rounded-2xl p-6 border border-emerald-700/30 hover:border-emerald-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-emerald-500/30 hover:-translate-y-1 overflow-hidden"
					>
						<div
							class="absolute inset-0 bg-gradient-to-br from-emerald-600/0 via-emerald-600/0 to-emerald-600/0 group-hover:from-emerald-600/10 group-hover:via-emerald-600/5 group-hover:to-emerald-600/10 transition-all duration-500"
						></div>

						<div class="relative z-10">
							<div
								class="w-16 h-16 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-emerald-500/30"
							>
								<Icon src={SiWhatsapp} size="2rem" color="white" />
							</div>
							<h3 class="text-xl font-bold text-white mb-2 group-hover:text-emerald-300 transition-colors">
								{t('contact.whatsapp')}
							</h3>
							<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
								{t('contact.subtitle')}
							</p>
							<p class="text-emerald-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
								Start conversation →
							</p>
						</div>
					</a>
				{/if}
			</div>

			<!-- Social Links Section -->
			<div
				class="relative bg-gradient-to-br from-gray-800/60 via-gray-800/40 to-gray-800/60 backdrop-blur-sm rounded-3xl p-8 md:p-12 border border-gray-700/50 overflow-hidden"
			>
				<!-- Decorative gradient overlay -->
				<div
					class="absolute inset-0 bg-gradient-to-r from-blue-600/5 via-purple-600/5 to-pink-600/5 opacity-50"
				></div>

				<div class="relative z-10">
					<div class="text-center mb-8">
						<h3 class="text-2xl font-bold text-white mb-2">Follow My Journey</h3>
						<p class="text-gray-400">Connect with me on social media</p>
					</div>

					<div class="flex flex-wrap justify-center gap-8">
						<a
							href={contact.github}
							target="_blank"
							rel="noopener noreferrer"
							class="group relative flex flex-col items-center gap-3 p-6 bg-gray-700/30 hover:bg-gray-700/50 rounded-xl border border-gray-600/50 hover:border-gray-500 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-gray-500/20"
							aria-label="GitHub"
						>
							<div
								class="w-14 h-14 bg-gradient-to-br from-gray-700 to-gray-800 rounded-xl flex items-center justify-center group-hover:from-gray-600 group-hover:to-gray-700 transition-all duration-300 shadow-lg"
							>
								<Icon src={SiGithub} size="2rem" color="white" />
							</div>
							<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
								>{t('contact.github')}</span
							>
						</a>

						<a
							href={contact.instagram}
							target="_blank"
							rel="noopener noreferrer"
							class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-pink-900/30 to-purple-900/30 hover:from-pink-800/40 hover:to-purple-800/40 rounded-xl border border-pink-700/30 hover:border-pink-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-pink-500/20"
							aria-label="Instagram"
						>
							<div
								class="w-14 h-14 bg-gradient-to-br from-pink-500 to-purple-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-pink-500/30"
							>
								<Icon src={SiInstagram} size="2rem" color="white" />
							</div>
							<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
								>{t('contact.instagram')}</span
							>
						</a>

						{#if contact.linkedin}
							<a
								href={contact.linkedin}
								target="_blank"
								rel="noopener noreferrer"
								class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-blue-900/30 to-indigo-900/30 hover:from-blue-800/40 hover:to-indigo-800/40 rounded-xl border border-blue-700/30 hover:border-blue-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-blue-500/20"
								aria-label="LinkedIn"
							>
								<div
									class="w-14 h-14 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-blue-500/30"
								>
									<Icon src={SiLinkedin} size="2rem" color="white" />
								</div>
								<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
									>{t('contact.linkedin')}</span
								>
							</a>
						{/if}

						{#if contact.whatsapp}
							<a
								href={getWhatsAppUrl()}
								target="_blank"
								rel="noopener noreferrer"
								class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-green-900/30 to-emerald-900/30 hover:from-green-800/40 hover:to-emerald-800/40 rounded-xl border border-green-700/30 hover:border-green-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-green-500/20"
								aria-label="WhatsApp"
							>
								<div
									class="w-14 h-14 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-green-500/30"
								>
									<Icon src={SiWhatsapp} size="2rem" color="white" />
								</div>
								<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
									>{t('contact.whatsapp')}</span
								>
							</a>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="container mx-auto px-6 py-12 border-t border-gray-700">
		<div class="max-w-4xl mx-auto text-center text-gray-400">
			<div class="flex justify-center gap-6 mb-6">
				<a
					href={contact.github}
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-400 hover:text-white transition-colors duration-200"
					aria-label="GitHub"
				>
					<Icon src={SiGithub} size="1.5rem" />
				</a>
				<a
					href={contact.instagram}
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-400 hover:text-white transition-colors duration-200"
					aria-label="Instagram"
				>
					<Icon src={SiInstagram} size="1.5rem" />
				</a>
				{#if contact.linkedin}
					<a
						href={contact.linkedin}
						target="_blank"
						rel="noopener noreferrer"
						class="text-gray-400 hover:text-white transition-colors duration-200"
						aria-label="LinkedIn"
					>
						<Icon src={SiLinkedin} size="1.5rem" />
					</a>
				{/if}
			</div>
			<p class="mb-4">
				{t('footer.tagline')}
				<span class="text-gray-500">|</span> {t('footer.backendDeveloper')}
				<span class="text-gray-500">|</span> {t('footer.devopsEnthusiast')}
			</p>
			<p class="text-sm">© 2024</p>
		</div>
	</footer>
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	section {
		animation: fadeIn 0.6s ease-out;
	}

	:global(html) {
		scroll-behavior: smooth;
	}

	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
