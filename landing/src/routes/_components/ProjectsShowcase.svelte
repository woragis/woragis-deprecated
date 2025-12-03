<script lang="ts">
	import { Search, Filter, X, ExternalLink, Calendar, TrendingUp, Code2, Globe, Github } from 'lucide-svelte';
	import { Icon } from 'svelte-icons-pack';
	import { SiGo, SiDocker, SiKubernetes, SiRedis, SiPython } from 'svelte-icons-pack/si';
	import type { Project, ProjectStatus } from '$lib/types/project';
	import type { CaseStudy } from '$lib/types/case-study';

	interface Props {
		projects: Project[];
		caseStudyMap?: Map<string, CaseStudy>;
		loading?: boolean;
	}

	let { projects = [], caseStudyMap = new Map(), loading = false }: Props = $props();

	let filteredProjects: Project[] = $state([]);
	let searchQuery = $state('');
	let statusFilter: ProjectStatus | 'all' = $state('all');
	let technologyFilter = $state('');
	let showFilters = $state(false);

	// Get unique technologies from all projects
	let availableTechnologies = $derived(
		Array.from(
			new Set(
				projects
					.flatMap((p) => p.technologies || [])
					.map((t) => t.name.toLowerCase())
			)
		).sort()
	);

	// Watch for projects data changes and apply filters
	$effect(() => {
		if (projects.length > 0) {
			applyFilters();
		}
	});

	function applyFilters() {
		let filtered = [...projects];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(project) =>
					project.name.toLowerCase().includes(query) ||
					project.description?.toLowerCase().includes(query) ||
					project.slug.toLowerCase().includes(query) ||
					project.technologies?.some((tech) =>
						tech.name.toLowerCase().includes(query)
					) ||
					project.skills?.some((skill) =>
						skill.name.toLowerCase().includes(query)
					)
			);
		}

		// Apply status filter
		if (statusFilter !== 'all') {
			filtered = filtered.filter((project) => project.status === statusFilter);
		}

		// Apply technology filter
		if (technologyFilter) {
			filtered = filtered.filter((project) =>
				project.technologies?.some((tech) =>
					tech.name.toLowerCase().includes(technologyFilter.toLowerCase())
				)
			);
		}

		filteredProjects = filtered;
	}

	function clearFilters() {
		searchQuery = '';
		statusFilter = 'all';
		technologyFilter = '';
		applyFilters();
	}

	// Watch for filter changes
	$effect(() => {
		applyFilters();
	});

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

	function getTechIcon(techName: string) {
		const name = techName.toLowerCase();
		if (name.includes('go') || name === 'golang') return SiGo;
		if (name.includes('docker')) return SiDocker;
		if (name.includes('kubernetes') || name.includes('k8s')) return SiKubernetes;
		if (name.includes('redis')) return SiRedis;
		if (name.includes('python')) return SiPython;
		return null;
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}
</script>

<div class="w-full">
	<!-- Search and Filter Bar -->
	<div class="mb-8">
		<div class="flex flex-col md:flex-row gap-4 mb-4">
			<!-- Search Input -->
			<div class="flex-1 relative">
				<Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
				<input
					type="text"
					placeholder="Search projects by name, description, tech stack..."
					bind:value={searchQuery}
					class="w-full pl-10 pr-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
				/>
			</div>
			<!-- Filter Toggle -->
			<button
				onclick={() => (showFilters = !showFilters)}
				class="px-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white hover:bg-gray-700/50 transition-colors flex items-center gap-2"
			>
				<Filter class="w-5 h-5" />
				Filters
			</button>
		</div>

		<!-- Filter Panel -->
		{#if showFilters}
			<div class="bg-gray-800/50 border border-gray-700 rounded-lg p-4 mb-4">
				<div class="flex flex-wrap gap-4">
					<!-- Status Filter -->
					<div class="flex-1 min-w-[200px]">
						<label class="block text-sm font-medium text-gray-300 mb-2">Status</label>
						<select
							bind:value={statusFilter}
							class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500"
						>
							<option value="all">All Statuses</option>
							<option value="idea">Idea</option>
							<option value="planning">Planning</option>
							<option value="executing">Executing</option>
							<option value="monitoring">Monitoring</option>
							<option value="completed">Completed</option>
						</select>
					</div>

					<!-- Technology Filter -->
					<div class="flex-1 min-w-[200px]">
						<label class="block text-sm font-medium text-gray-300 mb-2">Technology</label>
						<select
							bind:value={technologyFilter}
							class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500"
						>
							<option value="">All Technologies</option>
							{#each availableTechnologies as tech}
								<option value={tech}>{tech}</option>
							{/each}
						</select>
					</div>

					<!-- Clear Filters -->
					<div class="flex items-end">
						<button
							onclick={clearFilters}
							class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-white text-sm transition-colors flex items-center gap-2"
						>
							<X class="w-4 h-4" />
							Clear
						</button>
					</div>
				</div>
			</div>
		{/if}

		<!-- Active Filters Display -->
		{#if searchQuery || statusFilter !== 'all' || technologyFilter}
			<div class="flex flex-wrap gap-2 items-center">
				<span class="text-sm text-gray-400">Active filters:</span>
				{#if searchQuery}
					<span
						class="px-3 py-1 bg-blue-600/20 text-blue-300 text-xs rounded-full border border-blue-600/30 flex items-center gap-2"
					>
						Search: "{searchQuery}"
						<button onclick={() => (searchQuery = '')} class="hover:text-blue-200">
							<X class="w-3 h-3" />
						</button>
					</span>
				{/if}
				{#if statusFilter !== 'all'}
					<span
						class="px-3 py-1 bg-purple-600/20 text-purple-300 text-xs rounded-full border border-purple-600/30 flex items-center gap-2"
					>
						Status: {statusFilter}
						<button onclick={() => (statusFilter = 'all')} class="hover:text-purple-200">
							<X class="w-3 h-3" />
						</button>
					</span>
				{/if}
				{#if technologyFilter}
					<span
						class="px-3 py-1 bg-green-600/20 text-green-300 text-xs rounded-full border border-green-600/30 flex items-center gap-2"
					>
						Tech: {technologyFilter}
						<button onclick={() => (technologyFilter = '')} class="hover:text-green-200">
							<X class="w-3 h-3" />
						</button>
					</span>
				{/if}
			</div>
		{/if}
	</div>

	<!-- Projects Grid -->
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if filteredProjects.length === 0}
		<div class="text-center py-20">
			<Code2 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No projects found</p>
			<p class="text-gray-500 text-sm">
				{#if searchQuery || statusFilter !== 'all' || technologyFilter}
					Try adjusting your filters
				{:else}
					Check back later for new projects
				{/if}
			</p>
		</div>
	{:else}
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each filteredProjects as project, index}
				{@const animationDelay = index * 0.1}
				<div
					class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/20 hover:scale-[1.02] cursor-pointer relative overflow-hidden animate-fadeInUp"
					style="animation-delay: {animationDelay}s"
					onclick={() => (window.location.href = `/projects/${project.slug}`)}
					role="link"
					tabindex="0"
					onkeydown={(e) => e.key === 'Enter' && (window.location.href = `/projects/${project.slug}`)}
				>
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<div class="flex items-start justify-between mb-4">
							<div class="flex-1">
								<h3
									class="text-xl font-bold text-white mb-2 group-hover:text-blue-400 transition-colors"
								>
									{project.name}
								</h3>
								<div class="flex items-center gap-2 flex-wrap">
									<span
										class="px-2 py-1 rounded-lg text-xs font-medium capitalize {getStatusColor(project.status)}"
									>
										{project.status}
									</span>
									{#if project.healthScore >= 0}
										<div
											class="flex items-center gap-1 px-2 py-1 rounded-lg text-xs bg-green-600/20 text-green-400 border border-green-600/30"
										>
											<TrendingUp class="w-3 h-3" />
											{project.healthScore}%
										</div>
									{/if}
								</div>
							</div>
							<ExternalLink class="w-5 h-5 text-gray-400 group-hover:text-blue-400 transition-colors flex-shrink-0" />
						</div>

						{#if project.description}
							<p class="text-gray-300 text-sm mb-4 line-clamp-2">{project.description}</p>
						{/if}

						<!-- Technologies with Icons -->
						{#if project.technologies && project.technologies.length > 0}
							<div class="mb-4">
								<div class="flex flex-wrap gap-2">
									{#each project.technologies.slice(0, 5) as tech}
										{@const TechIcon = getTechIcon(tech.name)}
										<div
											class="flex items-center gap-1 px-2 py-1 bg-gray-700/50 rounded-lg text-xs text-gray-300 border border-gray-600 hover:border-blue-500/50 transition-colors"
											title={tech.name}
										>
											{#if TechIcon}
												<Icon src={TechIcon} size="0.875rem" color="currentColor" />
											{/if}
											<span>{tech.name}</span>
										</div>
									{/each}
									{#if project.technologies.length > 5}
										<span class="px-2 py-1 text-xs text-gray-400"
											>+{project.technologies.length - 5}</span
										>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Skills -->
						{#if project.skills && project.skills.length > 0}
							<div class="mb-4">
								<div class="flex flex-wrap gap-2">
									{#each project.skills.slice(0, 3) as skill}
										<div
											class="flex items-center gap-1 px-2 py-1 bg-blue-600/20 rounded-lg text-xs text-blue-300 border border-blue-600/30"
											title={skill.name}
										>
											{#if skill.icon}
												<span class="text-xs">{skill.icon}</span>
											{/if}
											<span>{skill.name}</span>
										</div>
									{/each}
									{#if project.skills.length > 3}
										<span class="px-2 py-1 text-xs text-gray-400"
											>+{project.skills.length - 3}</span
										>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Action Buttons -->
						<div class="flex items-center gap-2 mb-4 flex-wrap">
							<!-- Case Study Link -->
							{#if caseStudyMap.has(project.slug)}
								{@const caseStudy = caseStudyMap.get(project.slug)}
								{#if caseStudy}
									<a
										href="/case-studies/{caseStudy.id}"
										onclick={(e) => e.stopPropagation()}
										class="flex items-center gap-1 px-3 py-1.5 bg-purple-600/20 hover:bg-purple-600/30 text-purple-300 text-xs rounded-lg border border-purple-600/30 transition-colors"
									>
										<Code2 class="w-3 h-3" />
										Case Study
									</a>
								{/if}
							{/if}
							<!-- View Live Button (if project has a URL) -->
							{#if project.technologies?.some((t) => t.link)}
								<a
									href={project.technologies.find((t) => t.link)?.link || '#'}
									target="_blank"
									rel="noopener noreferrer"
									onclick={(e) => e.stopPropagation()}
									class="flex items-center gap-1 px-3 py-1.5 bg-green-600/20 hover:bg-green-600/30 text-green-300 text-xs rounded-lg border border-green-600/30 transition-colors"
								>
									<Globe class="w-3 h-3" />
									View Live
								</a>
							{/if}
							<!-- View Code Button (if GitHub link exists) -->
							{#if project.technologies?.some((t) => t.link?.includes('github'))}
								<a
									href={project.technologies.find((t) => t.link?.includes('github'))?.link || '#'}
									target="_blank"
									rel="noopener noreferrer"
									onclick={(e) => e.stopPropagation()}
									class="flex items-center gap-1 px-3 py-1.5 bg-gray-700/50 hover:bg-gray-600/50 text-gray-300 text-xs rounded-lg border border-gray-600 transition-colors"
								>
									<Github class="w-3 h-3" />
									View Code
								</a>
							{/if}
						</div>

						<div class="flex items-center justify-between pt-4 border-t border-gray-700">
							<div class="flex items-center gap-2 text-xs text-gray-400">
								<Calendar class="w-3 h-3" />
								<span>Updated {formatDate(project.updatedAt)}</span>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

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

