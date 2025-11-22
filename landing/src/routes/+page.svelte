<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiGo, SiDocker, SiKubernetes, SiRedis, SiGithub, SiInstagram, SiLinkedin, SiWhatsapp } from 'svelte-icons-pack/si';
	import { Settings, Brain, GitBranch, Mail, Phone, MapPin } from 'lucide-svelte';
	import { contact } from '$lib/constants';

	// Generate WhatsApp URL with pre-filled message
	function getWhatsAppUrl(): string {
		if (!contact.whatsapp) return '#';
		// Remove all non-digit characters from phone number
		const phoneNumber = contact.whatsapp.replace(/\D/g, '');
		const message = encodeURIComponent('Hello! I came across your website (www.woragis.me) and would like to get in touch.');
		return `https://wa.me/${phoneNumber}?text=${message}`;
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
					href="/projects"
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
					<span class="text-cyan-400 font-semibold">Golang</span> to build high-performance, scalable server applications.
					My goal is to become an accomplished backend developer who can design and implement distributed systems that
					handle real-world challenges.
				</p>
				<p class="text-lg text-gray-300 leading-relaxed">
					Beyond coding, I'm deeply invested in <span class="text-blue-400 font-semibold">DevOps</span> practices,
					specifically containerization with <span class="text-blue-300 font-semibold">Docker</span> and orchestration
					with <span class="text-blue-300 font-semibold">Kubernetes</span>. I believe that understanding infrastructure
					and deployment is crucial for building modern, cloud-native applications.
				</p>
			</div>
		</div>
	</section>

	<!-- Skills Section -->
	<section id="skills" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">Core Skills & Focus Areas</h2>
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				<!-- Golang Card -->
				<div
					class="bg-gradient-to-br from-cyan-900/30 to-cyan-800/20 rounded-xl p-6 border border-cyan-700/30 hover:border-cyan-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-cyan-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-cyan-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiGo} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-cyan-400">Golang</h3>
					</div>
					<p class="text-gray-300">
						My primary language for building high-performance backend services. I develop main server applications
						with Go, leveraging its concurrency model and efficiency.
					</p>
				</div>

				<!-- Docker Card -->
				<div
					class="bg-gradient-to-br from-blue-900/30 to-blue-800/20 rounded-xl p-6 border border-blue-700/30 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiDocker} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-blue-400">Docker</h3>
					</div>
					<p class="text-gray-300">
						Expertise in containerization for consistent deployments. I containerize applications for development,
						testing, and production environments.
					</p>
				</div>

				<!-- Kubernetes Card -->
				<div
					class="bg-gradient-to-br from-indigo-900/30 to-indigo-800/20 rounded-xl p-6 border border-indigo-700/30 hover:border-indigo-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-indigo-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-indigo-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiKubernetes} size="1.75rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold text-indigo-400">Kubernetes</h3>
					</div>
					<p class="text-gray-300">
						Orchestrating containerized applications at scale. I design and manage K8s clusters for production-grade
						infrastructure and deployments.
					</p>
				</div>

				<!-- DevOps Card -->
				<div
					class="bg-gradient-to-br from-purple-900/30 to-purple-800/20 rounded-xl p-6 border border-purple-700/30 hover:border-purple-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-purple-500/20"
				>
					<div class="flex items-center mb-4">
						<div class="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center mr-4">
							<Settings class="w-7 h-7" stroke="white" />
						</div>
						<h3 class="text-2xl font-bold text-purple-400">DevOps</h3>
					</div>
					<p class="text-gray-300">
						Bridging development and operations. I focus on automation, CI/CD pipelines, infrastructure as code, and
						cloud-native practices to streamline deployments.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Interests & Technologies Section -->
	<section id="interests" class="container mx-auto px-6 py-20">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl font-bold mb-12 text-center">Areas of Interest & Expertise</h2>
			<div class="grid md:grid-cols-2 gap-6">
				<!-- AI & RAG Card -->
				<div
					class="bg-gradient-to-br from-pink-900/30 to-purple-900/20 rounded-xl p-8 border border-pink-700/30 hover:border-pink-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-pink-500/20"
				>
					<div class="flex items-center mb-6">
						<div class="w-14 h-14 bg-gradient-to-br from-pink-600 to-purple-600 rounded-lg flex items-center justify-center mr-4">
							<Brain class="w-8 h-8" stroke="white" />
						</div>
						<h3 class="text-2xl font-bold bg-gradient-to-r from-pink-400 to-purple-400 bg-clip-text text-transparent">
							AI & RAG
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed mb-4">
						Fascinated by <span class="text-pink-400 font-semibold">Artificial Intelligence</span> and
						<span class="text-purple-400 font-semibold">Retrieval-Augmented Generation (RAG)</span> systems. I'm
						exploring how to build intelligent applications that can retrieve and synthesize information effectively.
					</p>
					<p class="text-gray-300 leading-relaxed">
						I work with AI servers built in <span class="text-yellow-400 font-semibold">Python</span>, leveraging
						modern ML frameworks and libraries to create intelligent backend services.
					</p>
				</div>

				<!-- Redis & Distributed Systems Card -->
				<div
					class="bg-gradient-to-br from-red-900/30 to-orange-900/20 rounded-xl p-8 border border-red-700/30 hover:border-red-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-red-500/20"
				>
					<div class="flex items-center mb-6">
						<div class="w-14 h-14 bg-gradient-to-br from-red-600 to-orange-600 rounded-lg flex items-center justify-center mr-4">
							<Icon src={SiRedis} size="2rem" color="white" />
						</div>
						<h3 class="text-2xl font-bold bg-gradient-to-r from-red-400 to-orange-400 bg-clip-text text-transparent">
							Redis & Pub/Sub
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed mb-4">
						Deep interest in <span class="text-red-400 font-semibold">Redis Pub/Sub</span> patterns for real-time
						communication between distributed services. I design systems where multiple servers communicate
						seamlessly through Redis messaging.
					</p>
					<p class="text-gray-300 leading-relaxed">
						I implement <span class="text-orange-400 font-semibold">inter-service communication</span> architectures
						using Redis as the backbone, enabling scalable and responsive distributed applications.
					</p>
				</div>

				<!-- Architecture Card -->
				<div
					class="bg-gradient-to-br from-green-900/30 to-emerald-900/20 rounded-xl p-8 border border-green-700/30 hover:border-green-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-green-500/20 md:col-span-2"
				>
					<div class="flex items-center mb-6">
						<div class="w-14 h-14 bg-gradient-to-br from-green-600 to-emerald-600 rounded-lg flex items-center justify-center mr-4">
							<GitBranch class="w-8 h-8" stroke="white" />
						</div>
						<h3 class="text-2xl font-bold bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent">
							Distributed Architecture
						</h3>
					</div>
					<p class="text-gray-300 leading-relaxed mb-4">
						I specialize in building <span class="text-green-400 font-semibold">hybrid architectures</span> where
						the main server is built with <span class="text-cyan-400 font-semibold">Golang</span> for performance
						and reliability, while <span class="text-yellow-400 font-semibold">AI services</span> are implemented
						in <span class="text-yellow-400 font-semibold">Python</span> to leverage the rich ML ecosystem.
					</p>
					<p class="text-gray-300 leading-relaxed">
						These services communicate through <span class="text-red-400 font-semibold">Redis Pub/Sub</span>,
						enabling a microservices architecture that's both scalable and maintainable. The combination allows each
						service to use the best tools for its specific domain.
					</p>
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
						<div class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-blue-500 transition-colors">
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
							<div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-green-500 transition-colors">
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
						<div
							class="flex items-center p-4 bg-gray-700/50 rounded-lg border border-gray-600"
						>
							<div class="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center mr-4">
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
							<div class="w-12 h-12 bg-blue-700 rounded-lg flex items-center justify-center mr-4 group-hover:bg-blue-600 transition-colors">
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
							<div class="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center mr-4 group-hover:bg-green-500 transition-colors">
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
</style>
