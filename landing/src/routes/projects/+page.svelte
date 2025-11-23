<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiGo, SiDocker, SiKubernetes, SiRedis, SiPython, SiGithub } from 'svelte-icons-pack/si';
	import { Search, Filter, X, ChevronDown, ExternalLink, Calendar, TrendingUp } from 'lucide-svelte';
	import type { Project, ProjectFilters, ProjectStatus, ProjectTechnology } from '$lib/types/project';
	import { useProjectsQuery } from '$lib/queries/projects';

	// Filters
	let searchQuery = $state('');
	let statusFilter: ProjectStatus | 'all' = $state('all');
	let technologyFilter = $state('');
	let sortBy: 'name' | 'createdAt' | 'updatedAt' | 'status' | 'healthScore' = $state('updatedAt');
	let sortOrder: 'asc' | 'desc' = $state('desc');
	let showFilters = $state(false);

	// Fetch all projects using TanStack Query (we'll filter client-side)
	const projectsQuery = useProjectsQuery();

	let projects = $derived(projectsQuery.data || []);
	let filteredProjects: Project[] = $state([]);
	let loading = $derived(projectsQuery.isPending);
	let error = $derived(projectsQuery.error ? (projectsQuery.error instanceof Error ? projectsQuery.error.message : 'Failed to fetch projects') : null);

	// Get unique technologies from all projects for filter
	let availableTechnologies = $derived.by(() => {
		const techSet = new Set<string>();
		projects.forEach((project: Project) => {
			project.technologies?.forEach((tech: ProjectTechnology) => {
				techSet.add(tech.name.toLowerCase());
			});
		});
		return Array.from(techSet).sort();
	});

	const statusOptions: Array<{ value: ProjectStatus | 'all'; label: string; color: string }> = [
		{ value: 'all', label: 'All Status', color: 'gray' },
		{ value: 'idea', label: 'Idea', color: 'purple' },
		{ value: 'planning', label: 'Planning', color: 'blue' },
		{ value: 'executing', label: 'Executing', color: 'yellow' },
		{ value: 'monitoring', label: 'Monitoring', color: 'green' },
		{ value: 'completed', label: 'Completed', color: 'cyan' }
	];

	const statusColors: Record<ProjectStatus | 'all', string> = {
		all: 'bg-gray-600',
		idea: 'bg-purple-600',
		planning: 'bg-blue-600',
		executing: 'bg-yellow-600',
		monitoring: 'bg-green-600',
		completed: 'bg-cyan-600'
	};

	const sortOptions = [
		{ value: 'name', label: 'Name' },
		{ value: 'createdAt', label: 'Date Created' },
		{ value: 'updatedAt', label: 'Last Updated' },
		{ value: 'status', label: 'Status' },
		{ value: 'healthScore', label: 'Health Score' }
	];

	// Apply filters when projects or filter values change
	$effect(() => {
		applyFilters();
	});

	function applyFilters() {
		let filtered = [...projects];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(project: Project) =>
					project.name.toLowerCase().includes(query) ||
					project.description?.toLowerCase().includes(query) ||
					project.slug.toLowerCase().includes(query)
			);
		}

		// Apply status filter
		if (statusFilter !== 'all') {
			filtered = filtered.filter((project: Project) => project.status === statusFilter);
		}

		// Apply technology filter
		if (technologyFilter) {
			filtered = filtered.filter((project: Project) =>
				project.technologies?.some((tech: ProjectTechnology) =>
					tech.name.toLowerCase().includes(technologyFilter.toLowerCase())
				)
			);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aVal: string | number;
			let bVal: string | number;

			switch (sortBy) {
				case 'name':
					aVal = a.name.toLowerCase();
					bVal = b.name.toLowerCase();
					break;
				case 'createdAt':
					aVal = new Date(a.createdAt).getTime();
					bVal = new Date(b.createdAt).getTime();
					break;
				case 'updatedAt':
					aVal = new Date(a.updatedAt).getTime();
					bVal = new Date(b.updatedAt).getTime();
					break;
				case 'status':
					aVal = a.status;
					bVal = b.status;
					break;
				case 'healthScore':
					aVal = a.healthScore;
					bVal = b.healthScore;
					break;
				default:
					return 0;
			}

			if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
			if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
			return 0;
		});

		filteredProjects = filtered;
	}

	// Watch for changes and apply filters
	$effect(() => {
		// Access all reactive values to track them
		const _ = [searchQuery, statusFilter, technologyFilter, sortBy, sortOrder, projects.length];
		applyFilters();
	});

	function clearFilters() {
		searchQuery = '';
		statusFilter = 'all';
		technologyFilter = '';
		sortBy = 'updatedAt';
		sortOrder = 'desc';
	}

	function getStatusColor(status: ProjectStatus): string {
		return statusColors[status] || statusColors.all;
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
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Header -->
	<section class="container mx-auto px-6 py-12 border-b border-gray-700">
		<div class="max-w-7xl mx-auto">
			<div class="flex items-center justify-between mb-6">
				<div>
					<h1 class="text-4xl md:text-5xl font-bold mb-2 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
						My Projects
					</h1>
					<p class="text-gray-400">Explore my portfolio of backend systems and applications</p>
				</div>
				<a
					href="/"
					class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors duration-200 text-sm font-medium"
				>
					← Back Home
				</a>
			</div>

			<!-- Search and Filter Bar -->
			<div class="space-y-4">
				<!-- Search Input -->
				<div class="relative">
					<Search class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
					<input
						type="text"
						placeholder="Search projects by name, description, or slug..."
						bind:value={searchQuery}
						oninput={applyFilters}
						class="w-full pl-12 pr-4 py-3 bg-gray-800 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-500"
					/>
				</div>

				<!-- Filter Toggle -->
				<div class="flex items-center justify-between gap-4">
					<button
						onclick={() => (showFilters = !showFilters)}
						class="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg transition-colors duration-200"
					>
						<Filter class="w-5 h-5" />
						<span>Filters</span>
						{#if showFilters}
							<ChevronDown class="w-4 h-4 rotate-180 transition-transform" />
						{:else}
							<ChevronDown class="w-4 h-4 transition-transform" />
						{/if}
					</button>

					<div class="flex items-center gap-4">
						<!-- Sort By -->
						<select
							bind:value={sortBy}
							onchange={applyFilters}
							class="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white text-sm"
						>
							{#each sortOptions as option}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>

						<!-- Sort Order -->
						<button
							onclick={() => {
								sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
								applyFilters();
							}}
							class="px-4 py-2 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg transition-colors duration-200 text-sm"
							title="Toggle sort order"
						>
							{sortOrder === 'asc' ? '↑' : '↓'}
						</button>

						{#if searchQuery || statusFilter !== 'all' || technologyFilter}
							<button
								onclick={clearFilters}
								class="flex items-center gap-2 px-4 py-2 bg-red-600/20 hover:bg-red-600/30 border border-red-700/30 rounded-lg transition-colors duration-200 text-sm"
							>
								<X class="w-4 h-4" />
								Clear
							</button>
						{/if}
					</div>
				</div>

				<!-- Filters Panel -->
				{#if showFilters}
					<div class="p-4 bg-gray-800/50 border border-gray-700 rounded-lg space-y-4">
						<!-- Status Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Status</div>
							<div class="flex flex-wrap gap-2">
								{#each statusOptions as option}
									<button
										onclick={() => {
											statusFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {statusFilter === option.value
											? `bg-${option.color}-600 border-${option.color}-500 text-white`
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										{option.label}
									</button>
								{/each}
							</div>
						</div>

						<!-- Technology Filter -->
						{#if availableTechnologies.length > 0}
							<div>
								<label for="technology-filter" class="block text-sm font-medium text-gray-300 mb-2">Technology</label>
								<select
									id="technology-filter"
									bind:value={technologyFilter}
									onchange={applyFilters}
									class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
								>
									<option value="">All Technologies</option>
									{#each availableTechnologies as tech}
										<option value={tech}>{tech.charAt(0).toUpperCase() + tech.slice(1)}</option>
									{/each}
								</select>
							</div>
						{/if}
					</div>
				{/if}
			</div>
		</div>
	</section>

	<!-- Projects Grid -->
	<section class="container mx-auto px-6 py-12">
		<div class="max-w-7xl mx-auto">
			{#if loading}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-6 text-center">
					<p class="text-red-400 mb-2">Error loading projects</p>
					<p class="text-gray-400 text-sm">{error}</p>
					<button
						onclick={() => projectsQuery.refetch()}
						class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{:else if filteredProjects.length === 0}
				<div class="text-center py-20">
					<p class="text-gray-400 text-lg mb-2">No projects found</p>
					<p class="text-gray-500 text-sm">Try adjusting your filters or search query</p>
				</div>
			{:else}
				<div class="mb-4 text-gray-400 text-sm">
					Showing {filteredProjects.length} of {projects.length} projects
				</div>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredProjects as project}
						<div
							class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20 group"
						>
							<!-- Header -->
							<div class="flex items-start justify-between mb-4">
								<div class="flex-1">
									<h3 class="text-xl font-bold text-white mb-1 group-hover:text-blue-400 transition-colors">
										{project.name}
									</h3>
									<div class="flex items-center gap-2 mb-2">
										<span
											class="px-2 py-1 rounded text-xs font-medium {getStatusColor(project.status)}"
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
								<a
									href="/projects/{project.slug}"
									class="text-gray-400 hover:text-blue-400 transition-colors"
									aria-label="View project details"
								>
									<ExternalLink class="w-5 h-5" />
								</a>
							</div>

							<!-- Description -->
							{#if project.description}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{project.description}</p>
							{/if}

							<!-- Skills -->
							{#if project.skills && project.skills.length > 0}
								<div class="mb-4">
									<p class="text-xs text-gray-400 mb-2">Skills</p>
									<div class="flex flex-wrap gap-2">
										{#each project.skills.slice(0, 6) as skill}
											<a
												href="/skills"
												class="flex items-center gap-1 px-2 py-1 bg-blue-600/20 hover:bg-blue-600/30 rounded text-xs text-blue-300 border border-blue-600/30 transition-colors"
												title={skill.name}
											>
												{#if skill.icon}
													<span class="text-xs">{skill.icon}</span>
												{/if}
												<span>{skill.name}</span>
											</a>
										{/each}
										{#if project.skills.length > 6}
											<span class="px-2 py-1 text-xs text-gray-400"
												>+{project.skills.length - 6} more</span
											>
										{/if}
									</div>
								</div>
							{/if}

							<!-- Technologies -->
							{#if project.technologies && project.technologies.length > 0}
								<div class="mb-4">
									<p class="text-xs text-gray-400 mb-2">Technologies</p>
									<div class="flex flex-wrap gap-2">
										{#each project.technologies.slice(0, 5) as tech}
											{@const TechIcon = getTechIcon(tech.name)}
											<div
												class="flex items-center gap-1 px-2 py-1 bg-gray-700/50 rounded text-xs text-gray-300 border border-gray-600"
											>
												{#if TechIcon}
													<Icon src={TechIcon} size="0.875rem" />
												{/if}
												<span>{tech.name}</span>
												{#if tech.version}
													<span class="text-gray-500">v{tech.version}</span>
												{/if}
											</div>
										{/each}
										{#if project.technologies.length > 5}
											<span class="px-2 py-1 text-xs text-gray-400"
												>+{project.technologies.length - 5} more</span
											>
										{/if}
									</div>
								</div>
							{/if}

							<!-- Footer -->
							<div class="flex items-center justify-between pt-4 border-t border-gray-700">
								<div class="flex items-center gap-2 text-xs text-gray-400">
									<Calendar class="w-3 h-3" />
									<span>Updated {formatDate(project.updatedAt)}</span>
								</div>
								<a
									href="/projects/{project.slug}"
									class="text-xs text-blue-400 hover:text-blue-300 transition-colors font-medium"
								>
									View Details →
								</a>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</section>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>

