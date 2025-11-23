<script lang="ts">
	import { Search, Filter, X, Tag, Code, Database, Server, Wrench } from 'lucide-svelte';
	import { useSkillsWithCountsQuery } from '$lib/queries/skills';
	import type { SkillWithCount, SkillCategory } from '$lib/api/skills';

	// Fetch skills using TanStack Query
	const skillsQuery = useSkillsWithCountsQuery();
	let skills = $derived(skillsQuery.data || []);
	let filteredSkills: SkillWithCount[] = $state([]);
	let loading = $derived(skillsQuery.isPending);
	let error = $derived(skillsQuery.error ? (skillsQuery.error instanceof Error ? skillsQuery.error.message : 'Failed to fetch skills') : null);

	// Filters
	let searchQuery = $state('');
	let categoryFilter: SkillCategory | 'all' = $state('all');
	let showFilters = $state(false);

	const categoryOptions: Array<{ value: SkillCategory | 'all'; label: string; icon: any }> = [
		{ value: 'all', label: 'All Categories', icon: Tag },
		{ value: 'backend', label: 'Backend', icon: Server },
		{ value: 'frontend', label: 'Frontend', icon: Code },
		{ value: 'database', label: 'Database', icon: Database },
		{ value: 'infrastructure', label: 'Infrastructure', icon: Server },
		{ value: 'devops', label: 'DevOps', icon: Wrench },
		{ value: 'language', label: 'Language', icon: Code },
		{ value: 'framework', label: 'Framework', icon: Code },
		{ value: 'tool', label: 'Tool', icon: Wrench },
		{ value: 'service', label: 'Service', icon: Server },
		{ value: 'library', label: 'Library', icon: Code },
		{ value: 'other', label: 'Other', icon: Tag }
	];

	const categoryColors: Record<SkillCategory | 'all', string> = {
		all: 'bg-gray-600',
		backend: 'bg-red-600',
		frontend: 'bg-blue-600',
		database: 'bg-green-600',
		infrastructure: 'bg-purple-600',
		devops: 'bg-orange-600',
		language: 'bg-yellow-600',
		framework: 'bg-indigo-600',
		tool: 'bg-pink-600',
		service: 'bg-cyan-600',
		library: 'bg-teal-600',
		other: 'bg-gray-600'
	};


	function applyFilters() {
		let filtered = [...skills];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(skill) =>
					skill.name.toLowerCase().includes(query) ||
					skill.description?.toLowerCase().includes(query) ||
					skill.slug.toLowerCase().includes(query)
			);
		}

		// Apply category filter
		if (categoryFilter !== 'all') {
			filtered = filtered.filter((skill) => skill.category === categoryFilter);
		}

		// Sort by project count (descending), then by name
		filtered.sort((a, b) => {
			if (b.projectCount !== a.projectCount) {
				return b.projectCount - a.projectCount;
			}
			return a.name.localeCompare(b.name);
		});

		filteredSkills = filtered;
	}

	$effect(() => {
		applyFilters();
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<div class="container mx-auto px-4 py-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-4xl font-bold mb-2">Skills</h1>
			<p class="text-gray-400">Explore all skills and see how many projects use each one</p>
		</div>

		<!-- Search and Filters -->
		<div class="mb-6 space-y-4">
			<!-- Search Bar -->
			<div class="relative">
				<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
				<input
					type="text"
					placeholder="Search skills..."
					bind:value={searchQuery}
					class="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white placeholder-gray-400"
				/>
			</div>

			<!-- Category Filter -->
			<div class="flex flex-wrap gap-2">
				{#each categoryOptions as option}
					<button
						onclick={() => {
							categoryFilter = option.value;
						}}
						class="px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2 {categoryFilter === option.value
							? categoryColors[option.value] + ' text-white'
							: 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
					>
						<svelte:component this={option.icon} class="w-4 h-4" />
						{option.label}
					</button>
				{/each}
			</div>
		</div>

		<!-- Error State -->
		{#if error}
			<div class="bg-red-900/50 border border-red-700 rounded-lg p-4 mb-6">
				<p class="text-red-200">{error}</p>
			</div>
		{/if}

		<!-- Loading State -->
		{#if loading}
			<div class="flex justify-center items-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
			</div>
		{:else if filteredSkills.length === 0}
			<!-- Empty State -->
			<div class="text-center py-20">
				<Tag class="w-16 h-16 mx-auto text-gray-600 mb-4" />
				<p class="text-gray-400 text-lg">No skills found</p>
				{#if searchQuery || categoryFilter !== 'all'}
					<p class="text-gray-500 text-sm mt-2">Try adjusting your filters</p>
				{/if}
			</div>
		{:else}
			<!-- Skills Grid -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
				{#each filteredSkills as skill}
					<div
						class="bg-gray-800 border border-gray-700 rounded-lg p-6 hover:border-blue-500 transition-colors cursor-pointer"
					>
						<div class="flex items-start justify-between mb-4">
							<div class="flex-1">
								<h3 class="text-xl font-semibold mb-1">{skill.name}</h3>
								<span
									class="inline-block px-2 py-1 text-xs rounded {categoryColors[skill.category]}"
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
							<p class="text-gray-400 text-sm mb-4 line-clamp-2">{skill.description}</p>
						{/if}

						<div class="flex items-center justify-between pt-4 border-t border-gray-700">
							<div class="flex items-center gap-2">
								<Tag class="w-4 h-4 text-gray-400" />
								<span class="text-sm text-gray-400">Projects</span>
							</div>
							<span class="text-lg font-bold text-blue-400">{skill.projectCount}</span>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>

