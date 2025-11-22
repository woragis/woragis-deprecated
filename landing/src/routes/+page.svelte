<script lang="ts">
	import { onMount } from 'svelte';
	
	let mounted = false;
	
	onMount(() => {
		mounted = true;
	});
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
					href="#skills"
					class="px-8 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
				>
					View Skills
				</a>
				<a
					href="#interests"
					class="px-8 py-3 bg-gray-700 hover:bg-gray-600 rounded-lg font-medium transition-colors duration-200 border border-gray-600"
				>
					My Interests
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
							<svg class="w-7 h-7" fill="currentColor" viewBox="0 0 24 24">
								<path
									d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-1.94-1.4-2.72-2.55C9.29 15.98 8.02 14 6.88 12a3.92 3.92 0 0 0-2.43 1.67c-.3.4-.54.83-.67 1.3a4.13 4.13 0 0 0-.08 1c.01.4.1.78.26 1.13.17.3.4.56.68.78.57.43 1.31.66 2.2.7 1.3.06 2.79-.43 4.43-1.47 2.32-1.48 4.95-3.98 7.45-7.52 2.68-3.82 4.9-8.41 5.84-13.2.05-.3.14-.6.28-.85.18-.34.43-.62.75-.82.3-.18.62-.27.98-.3.37-.03.75.01 1.13.11.4.1.75.28 1.06.53.3.24.55.54.73.88.06.13.11.27.15.41.07.3.1.6.09.9-.02.58-.17 1.16-.46 1.72-.58 1.11-1.5 2.07-2.6 2.72-.96.57-2.02.9-3.1 1.03-.93.12-1.88.08-2.8-.06-1.77-.26-3.5-.87-5.1-1.81-.65-.38-1.27-.81-1.86-1.27-.3-.23-.58-.48-.85-.74-.27-.26-.52-.53-.76-.82-.27-.32-.5-.67-.7-1.04-.1-.18-.18-.37-.25-.56-.07-.2-.13-.4-.18-.6-.06-.23-.1-.46-.12-.7-.01-.24 0-.48.03-.72.03-.24.08-.48.15-.71.07-.23.16-.45.27-.66.1-.21.22-.41.36-.6.13-.18.28-.35.44-.51.16-.16.33-.31.51-.45.18-.14.37-.26.57-.37.2-.1.41-.19.63-.26.22-.07.45-.12.68-.16.23-.03.47-.05.71-.04.24.01.48.04.71.08.23.04.46.1.68.17.22.08.43.17.63.28.2.1.39.22.57.35.18.13.35.28.51.44.16.15.3.32.44.49.14.18.26.37.37.57.11.2.2.41.27.63.07.22.12.45.16.68.04.23.06.47.07.71.01.24 0 .48-.03.72-.03.24-.08.47-.15.7-.07.23-.16.45-.27.66-.1.21-.22.41-.36.6-.13.18-.28.35-.44.51-.16.16-.33.31-.51.45-.18.14-.37.26-.57.37-.2.1-.41.19-.63.26-.22.07-.45.12-.68.16-.23.03-.47.05-.71.04-.24-.01-.48-.04-.71-.08-.23-.04-.46-.1-.68-.17-.22-.08-.43-.17-.63-.28-.2-.1-.39-.22-.57-.35-.18-.13-.35-.28-.51-.44-.16-.15-.3-.32-.44-.49-.14-.18-.26-.37-.37-.57-.11-.2-.2-.41-.27-.63-.07-.22-.12-.45-.16-.68-.04-.23-.06-.47-.07-.71-.01-.24 0-.48.03-.72.03-.24.08-.47.15-.7.07-.23.16-.45.27-.66.1-.21.22-.41.36-.6.13-.18.28-.35.44-.51.16-.16.33-.31.51-.45.18-.14.37-.26.57-.37.2-.1.41-.19.63-.26.22-.07.45-.12.68-.16.23-.03.47-.05.71-.04.24.01.48.04.71.08.23.04.46.1.68.17.22.08.43.17.63.28.2.1.39.22.57.35.18.13.35.28.51.44.16.15.3.32.44.49.14.18.26.37.37.57.11.2.2.41.27.63.07.22.12.45.16.68.04.23.06.47.07.71.01.24 0 .48-.03.72-.03.24-.08.47-.15.7-.07.23-.16.45-.27.66-.1.21-.22.41-.36.6-.13.18-.28.35-.44.51-.16.16-.33.31-.51.45-.18.14-.37.26-.57.37-.2.1-.41.19-.63.26-.22.07-.45.12-.68.16-.23.03-.47.05-.71.04-.24-.01-.48-.04-.71-.08z"
								/>
							</svg>
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
							<svg class="w-7 h-7" fill="currentColor" viewBox="0 0 24 24">
								<path
									d="M13.983 0c-.035 0-.07.004-.105.004v.004H6.816c-.317 0-.625.106-.872.297-.248.191-.425.456-.497.75L4.2 6.7H1.116c-.304 0-.598.123-.813.344-.216.22-.337.52-.337.83v9.39c0 .31.12.61.337.83.215.22.51.344.813.344h16.768c.304 0 .598-.123.813-.344.216-.22.337-.52.337-.83V7.874c0-.31-.12-.61-.337-.83-.215-.22-.509-.344-.813-.344H19.8l-1.442-5.7c-.072-.294-.249-.559-.497-.75C17.613.106 17.305 0 16.988 0h-3.005zm.503 1.494h2.502c.09 0 .176.028.245.08.07.05.12.122.14.2l1.185 4.686c.01.042.01.086 0 .128v9.026c0 .12.05.236.138.321.088.086.206.134.33.134h.89c.124 0 .242-.048.33-.134.088-.085.138-.2.138-.32V7.888a.443.443 0 0 0-.003-.047l-1.18-4.693a.456.456 0 0 1 .138-.2c.07-.052.156-.08.246-.08h1.79c.09 0 .176.028.245.08.07.05.12.122.14.2l1.18 4.693c.003.016.003.031.003.047v9.626c0 .31-.12.61-.337.83-.216.22-.51.344-.813.344H1.116c-.304 0-.598-.123-.813-.344-.216-.22-.337-.52-.337-.83V7.874c0-.31.12-.61.337-.83.215-.22.51-.344.813-.344h.89c.124 0 .242.048.33.134.088.085.138.2.138.32v9.026c0 .042.01.086 0 .128l1.185 4.686c.02.078.07.15.14.2.069.052.155.08.245.08h11.974c.09 0 .176-.028.245-.08.07-.05.12-.122.14-.2l1.185-4.686a.443.443 0 0 0 0-.128V7.588l-1.442-5.7c-.02-.078-.07-.15-.14-.2-.069-.052-.155-.08-.245-.08h-2.502a.456.456 0 0 0-.246.08.456.456 0 0 1-.14.2l-.005.02v.004c-.01.042-.01.086 0 .128v.004z"
								/>
							</svg>
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
							<span class="text-white font-bold text-lg">K8s</span>
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
							<svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
								/>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
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
							<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
								/>
							</svg>
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
							<svg class="w-8 h-8" fill="currentColor" viewBox="0 0 24 24">
								<path
									d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"
								/>
							</svg>
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
							<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
								/>
							</svg>
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

	<!-- Footer -->
	<footer class="container mx-auto px-6 py-12 border-t border-gray-700">
		<div class="max-w-4xl mx-auto text-center text-gray-400">
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

	html {
		scroll-behavior: smooth;
	}
</style>
