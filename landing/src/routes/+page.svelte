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
	import { Mail, Phone, MapPin, ExternalLink, Calendar, TrendingUp, Code2, Settings, Brain, GitBranch } from 'lucide-svelte';
	import { contact, skills, interests } from '$lib/constants';
	import { listProjects } from '$lib/api/projects';
	import type { Project, ProjectTechnology } from '$lib/types/project';

	// Generate WhatsApp URL with pre-filled message
	function getWhatsAppUrl(): string {
		if (!contact.whatsapp) return '#';
		const phoneNumber = contact.whatsapp.replace(/\D/g, '');
		const message = encodeURIComponent(
			'Hello! I came across your website (www.woragis.me) and would like to get in touch.'
		);
		return `https://wa.me/${phoneNumber}?text=${message}`;
	}

	// Featured projects
	let featuredProjects: Project[] = $state([]);
	let loadingProjects = $state(false);

	onMount(async () => {
		await fetchFeaturedProjects();
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
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Hero Section -->
	<section class="container mx-auto px-6 py-20 md:py-32">
		<div class="max-w-4xl mx-auto text-center">
			<div
				class="inline-block mb-6 px-4 py-2 bg-blue-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium"
			>
				Backend Developer
			</div>
			<h1
				class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent animate-pulse"
			>
				Hello, I'm a Developer
			</h1>
			<p class="text-xl md:text-2xl text-gray-300 mb-8 leading-relaxed">
				Passionate about building robust backend systems with <span class="text-cyan-400 font-semibold">Golang</span>,
				orchestrating with <span class="text-blue-400 font-semibold">Docker</span> &
				<span class="text-blue-300 font-semibold">Kubernetes</span>, and exploring the frontiers of
				<span class="text-purple-400 font-semibold">AI</span> & <span class="text-pink-400 font-semibold">RAG</span>.
			</p>
			<div class="flex flex-wrap gap-4 justify-center">
				<a
					href="#projects"
					class="px-8 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
				>
					View Projects
				</a>
				<a
					href="#skills"
					class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
				>
					View Skills
				</a>
				<a
					href="#contact"
					class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
				>
					Contact Me
				</a>
			</div>
		</div>
	</section>

	<!-- About Section -->
	<section id="about" class="container mx-auto px-6 py-20">
		<div class="max-w-4xl mx-auto">
			<h2 class="text-4xl font-bold mb-8 text-center">About Me</h2>
			<div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 border border-gray-700 shadow-xl">
				<p class="text-lg text-gray-300 leading-relaxed mb-4">
					I'm a developer with a passion for backend development, currently focusing on mastering
					<span class="text-cyan-400 font-semibold">Golang</span> to build high-performance, scalable server
					applications. My goal is to become an accomplished backend developer who can design and implement
					distributed systems that handle real-world challenges.
				</p>
				<p class="text-lg text-gray-300 leading-relaxed">
					Beyond coding, I'm deeply invested in <span class="text-blue-400 font-semibold">DevOps</span> practices,
					specifically containerization with <span class="text-blue-300 font-semibold">Docker</span> and orchestration
					with <span class="text-blue-300 font-semibold">Kubernetes</span>. I believe that understanding
					infrastructure and deployment is crucial for building modern, cloud-native applications.
				</p>
			</div>
		</div>
	</section>

	<!-- Projects Section -->
	<section id="projects" class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="flex items-center justify-between mb-12">
				<h2 class="text-4xl font-bold">Featured Projects</h2>
				<a
					href="/projects"
					class="text-blue-400 hover:text-blue-300 transition-colors duration-200 flex items-center gap-2"
				>
					View All Projects
					<ExternalLink class="w-5 h-5" />
				</a>
			</div>

			{#if loadingProjects}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if featuredProjects.length === 0}
				<div class="text-center py-20">
					<Code2 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No projects found</p>
					<p class="text-gray-500 text-sm">Check back later for featured projects</p>
				</div>
			{:else}
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each featuredProjects as project}
						<a
							href="/projects/{project.slug}"
							class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 group"
						>
							<div class="flex items-start justify-between mb-4">
								<div class="flex-1">
									<h3
										class="text-xl font-bold text-white mb-2 group-hover:text-blue-400 transition-colors"
									>
										{project.name}
									</h3>
									<div class="flex items-center gap-2 flex-wrap">
										<span
											class="px-2 py-1 rounded text-xs font-medium capitalize {getStatusColor(project.status)}"
										>
											{project.status}
										</span>
										{#if project.healthScore >= 0}
											<div
												class="flex items-center gap-1 px-2 py-1 rounded text-xs bg-green-600/20 text-green-400 border border-green-600/30"
											>
												<TrendingUp class="w-3 h-3" />
												{project.healthScore}%
											</div>
										{/if}
									</div>
								</div>
								<ExternalLink class="w-5 h-5 text-gray-400 group-hover:text-blue-400 transition-colors" />
							</div>

							{#if project.description}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{project.description}</p>
							{/if}

							{#if project.technologies && project.technologies.length > 0}
								<div class="mb-4">
									<div class="flex flex-wrap gap-2">
									{#each project.technologies.slice(0, 4) as tech}
										{@const TechIcon = getTechIcon(tech.name)}
										<div
											class="flex items-center gap-1 px-2 py-1 bg-gray-700/50 rounded text-xs text-gray-300 border border-gray-600"
										>
											{#if TechIcon}
												<Icon src={TechIcon} size="0.875rem" color="currentColor" />
											{/if}
											<span>{tech.name}</span>
										</div>
									{/each}
										{#if project.technologies.length > 4}
											<span class="px-2 py-1 text-xs text-gray-400"
												>+{project.technologies.length - 4}</span
											>
										{/if}
									</div>
								</div>
							{/if}

							<div class="flex items-center justify-between pt-4 border-t border-gray-700">
								<div class="flex items-center gap-2 text-xs text-gray-400">
									<Calendar class="w-3 h-3" />
									<span>Updated {formatDate(project.updatedAt)}</span>
								</div>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	</section>

	<!-- Skills Section -->
	<section id="skills" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">Core Skills & Focus Areas</h2>
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				<!-- Golang -->
				<div
					class="bg-gradient-to-br from-cyan-900/30 to-cyan-800/20 rounded-xl p-6 border border-cyan-700/30 hover:border-cyan-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-cyan-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-cyan-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiGo} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-cyan-400">{skills[0].name}</h3>
					</div>
					<p class="text-gray-300">{skills[0].description}</p>
				</div>

				<!-- Python -->
				<div
					class="bg-gradient-to-br from-yellow-900/30 to-yellow-800/20 rounded-xl p-6 border border-yellow-700/30 hover:border-yellow-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-yellow-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-yellow-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiPython} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-yellow-400">{skills[1].name}</h3>
					</div>
					<p class="text-gray-300">{skills[1].description}</p>
				</div>

				<!-- Docker -->
				<div
					class="bg-gradient-to-br from-blue-900/30 to-blue-800/20 rounded-xl p-6 border border-blue-700/30 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiDocker} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-blue-400">{skills[2].name}</h3>
					</div>
					<p class="text-gray-300">{skills[2].description}</p>
				</div>

				<!-- Kubernetes -->
				<div
					class="bg-gradient-to-br from-indigo-900/30 to-indigo-800/20 rounded-xl p-6 border border-indigo-700/30 hover:border-indigo-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-indigo-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-indigo-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiKubernetes} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-indigo-400">{skills[3].name}</h3>
					</div>
					<p class="text-gray-300">{skills[3].description}</p>
				</div>

				<!-- DevOps -->
				<div
					class="bg-gradient-to-br from-purple-900/30 to-purple-800/20 rounded-xl p-6 border border-purple-700/30 hover:border-purple-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-purple-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center mr-4">
							<Settings class="w-7 h-7" stroke="white" />
						</div>
						<h3 class="text-2xl font-bold text-purple-400">{skills[4].name}</h3>
					</div>
					<p class="text-gray-300">{skills[4].description}</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Interests & Technologies Section -->
	<section id="interests" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">Areas of Interest & Expertise</h2>
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

	<!-- Contact Section -->
	<section id="contact" class="container mx-auto px-6 py-20">
		<div class="max-w-4xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">Get In Touch</h2>
			<div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 border border-gray-700 shadow-xl">
				<p class="text-lg text-gray-300 text-center mb-8 leading-relaxed">
					I'm always open to discussing new projects, creative ideas, or opportunities to be part of your vision.
					Feel free to reach out!
				</p>
				<div class="grid md:grid-cols-2 gap-6">
					<!-- Email -->
					<a
						href="mailto:{contact.email}"
						class="flex items-center p-4 bg-gray-700/50 hover:bg-gray-700 rounded-lg border border-gray-600 hover:border-blue-500 transition-all duration-200 group"
					>
						<div
							class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-blue-500 transition-colors"
						>
							<Mail class="w-6 h-6 text-white" />
						</div>
						<div>
							<p class="text-sm text-gray-400 mb-1">Email</p>
							<p class="text-white font-medium">{contact.email}</p>
						</div>
					</a>

					<!-- Phone -->
					{#if contact.phone}
						<a
							href="tel:{contact.phone}"
							class="flex items-center p-4 bg-gray-700/50 hover:bg-gray-700 rounded-lg border border-gray-600 hover:border-green-500 transition-all duration-200 group"
						>
							<div
								class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-green-500 transition-colors"
							>
								<Phone class="w-6 h-6 text-white" />
							</div>
							<div>
								<p class="text-sm text-gray-400 mb-1">Phone</p>
								<p class="text-white font-medium">{contact.phone}</p>
							</div>
						</a>
					{/if}

					<!-- Location -->
					{#if contact.location}
						<div class="flex items-center p-4 bg-gray-700/50 rounded-lg border border-gray-600">
							<div
								class="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center mr-4"
							>
								<MapPin class="w-6 h-6 text-white" />
							</div>
							<div>
								<p class="text-sm text-gray-400 mb-1">Location</p>
								<p class="text-white font-medium">{contact.location}</p>
							</div>
						</div>
					{/if}

					<!-- LinkedIn -->
					{#if contact.linkedin}
						<a
							href={contact.linkedin}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-center p-4 bg-gray-700/50 hover:bg-gray-700 rounded-lg border border-gray-600 hover:border-blue-400 transition-all duration-200 group"
						>
							<div
								class="w-12 h-12 bg-blue-700 rounded-lg flex items-center justify-center mr-4 group-hover:bg-blue-600 transition-colors"
							>
								<Icon src={SiLinkedin} size="1.5rem" color="white" />
							</div>
							<div>
								<p class="text-sm text-gray-400 mb-1">LinkedIn</p>
								<p class="text-white font-medium">Connect with me</p>
							</div>
						</a>
					{/if}

					<!-- WhatsApp -->
					{#if contact.whatsapp}
						<a
							href={getWhatsAppUrl()}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-center p-4 bg-gray-700/50 hover:bg-gray-700 rounded-lg border border-gray-600 hover:border-green-500 transition-all duration-200 group"
						>
							<div
								class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-green-500 transition-colors"
							>
								<Icon src={SiWhatsapp} size="1.5rem" color="white" />
							</div>
							<div>
								<p class="text-sm text-gray-400 mb-1">WhatsApp</p>
								<p class="text-white font-medium">Message me</p>
							</div>
						</a>
					{/if}
				</div>

				<!-- Social Links -->
				<div class="mt-8 pt-8 border-t border-gray-700">
					<p class="text-center text-gray-400 mb-4">Or find me on</p>
					<div class="flex justify-center gap-6">
						<a
							href={contact.github}
							target="_blank"
							rel="noopener noreferrer"
							class="text-gray-400 hover:text-white transition-colors duration-200"
							aria-label="GitHub"
						>
							<Icon src={SiGithub} size="2rem" />
						</a>
						<a
							href={contact.instagram}
							target="_blank"
							rel="noopener noreferrer"
							class="text-gray-400 hover:text-white transition-colors duration-200"
							aria-label="Instagram"
						>
							<Icon src={SiInstagram} size="2rem" />
						</a>
						{#if contact.linkedin}
							<a
								href={contact.linkedin}
								target="_blank"
								rel="noopener noreferrer"
								class="text-gray-400 hover:text-white transition-colors duration-200"
								aria-label="LinkedIn"
							>
								<Icon src={SiLinkedin} size="2rem" />
							</a>
						{/if}
						{#if contact.whatsapp}
							<a
								href={getWhatsAppUrl()}
								target="_blank"
								rel="noopener noreferrer"
								class="text-gray-400 hover:text-white transition-colors duration-200"
								aria-label="WhatsApp"
							>
								<Icon src={SiWhatsapp} size="2rem" />
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
				Building the future, one service at a time.
				<span class="text-gray-500">|</span> Backend Developer
				<span class="text-gray-500">|</span> DevOps Enthusiast
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
