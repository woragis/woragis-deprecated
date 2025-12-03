<script lang="ts">
	import { Tag, ExternalLink, Star, TrendingUp, Calendar } from 'lucide-svelte';
	import type { SkillWithCount, ProficiencyLevel } from '$lib/api/skills';

	interface Props {
		skills: SkillWithCount[];
		loading: boolean;
	}

	let { skills = [], loading = false }: Props = $props();

	function getProficiencyLabel(level?: ProficiencyLevel): string {
		if (!level) return '';
		const labels: Record<ProficiencyLevel, string> = {
			expert: 'Expert',
			advanced: 'Advanced',
			proficient: 'Proficient',
			learning: 'Learning'
		};
		return labels[level] || '';
	}

	function getProficiencyColor(level?: ProficiencyLevel): string {
		if (!level) return 'bg-gray-600/20 text-gray-300 border-gray-500/30';
		const colors: Record<ProficiencyLevel, string> = {
			expert: 'bg-purple-600/20 text-purple-300 border-purple-500/30',
			advanced: 'bg-blue-600/20 text-blue-300 border-blue-500/30',
			proficient: 'bg-green-600/20 text-green-300 border-green-500/30',
			learning: 'bg-yellow-600/20 text-yellow-300 border-yellow-500/30'
		};
		return colors[level] || colors.learning;
	}

	function getProficiencyStars(level?: ProficiencyLevel): number {
		if (!level) return 0;
		const stars: Record<ProficiencyLevel, number> = {
			expert: 4,
			advanced: 3,
			proficient: 2,
			learning: 1
		};
		return stars[level] || 0;
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
	}

	function formatYears(years?: number): string {
		if (!years) return '';
		if (years === 1) return '1 year';
		return `${years} years`;
	}
</script>

<section id="skills" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				Popular Skills
			</h2>
			<p class="text-xl text-gray-300 mb-6">Technologies and tools I work with</p>
			<a
				href="/skills"
				class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
			>
				View All Skills
				<ExternalLink class="w-5 h-5" />
			</a>
		</div>

		{#if loading}
			<div class="flex items-center justify-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		{:else if skills.length === 0}
			<div class="text-center py-20">
				<Tag class="w-16 h-16 mx-auto mb-4 text-gray-600" />
				<p class="text-gray-400 text-lg mb-2">No skills found</p>
				<p class="text-gray-500 text-sm">Skills will appear here once they're added to projects</p>
			</div>
		{:else}
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each skills as skill, index (skill.id || skill.name)}
					{@const animationDelay = index * 0.1}
					<a
						href="/skills"
						class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/20 hover:scale-[1.02] relative overflow-hidden animate-fadeInUp"
						style="animation-delay: {animationDelay}s"
					>
						<!-- Decorative gradient overlay -->
						<div class="absolute inset-0 bg-gradient-to-br from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-300 pointer-events-none"></div>
						<div class="relative z-10">
							<div class="flex items-start justify-between mb-4">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-2">
										<h3
											class="text-xl font-bold text-white group-hover:text-blue-400 transition-colors"
										>
											{skill.name}
										</h3>
										{#if skill.proficiencyLevel}
											<span
												class="px-2 py-0.5 text-xs rounded-full border {getProficiencyColor(skill.proficiencyLevel)}"
											>
												{getProficiencyLabel(skill.proficiencyLevel)}
											</span>
										{/if}
									</div>
									<div class="flex items-center gap-2 flex-wrap">
										<span
											class="inline-block px-2 py-1 text-xs rounded-lg bg-gray-700/50 text-gray-300 capitalize"
										>
											{skill.category}
										</span>
										{#if skill.proficiencyLevel}
											<div class="flex items-center gap-0.5">
												{#each Array(getProficiencyStars(skill.proficiencyLevel)) as _}
													<Star class="w-3 h-3 text-yellow-400 fill-yellow-400" />
												{/each}
												{#each Array(4 - getProficiencyStars(skill.proficiencyLevel)) as _}
													<Star class="w-3 h-3 text-gray-600" />
												{/each}
											</div>
										{/if}
									</div>
								</div>
								{#if skill.icon}
									<div class="w-12 h-12 bg-gray-700/50 rounded-lg flex items-center justify-center border border-gray-600">
										<span class="text-xl">{skill.icon}</span>
									</div>
								{/if}
							</div>

							{#if skill.description}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{skill.description}</p>
							{/if}

							<!-- Experience & Usage Info -->
							{#if skill.yearsOfExperience || skill.lastUsedDate}
								<div class="mb-4 space-y-2">
									{#if skill.yearsOfExperience}
										<div class="flex items-center gap-2 text-sm text-gray-400">
											<TrendingUp class="w-4 h-4" />
											<span>{formatYears(skill.yearsOfExperience)} of experience</span>
										</div>
									{/if}
									{#if skill.lastUsedDate}
										<div class="flex items-center gap-2 text-sm text-gray-400">
											<Calendar class="w-4 h-4" />
											<span>Last used: {formatDate(skill.lastUsedDate)}</span>
										</div>
									{/if}
								</div>
							{/if}

							<div class="flex items-center justify-between pt-4 border-t border-gray-700">
								<div class="flex items-center gap-2">
									<Tag class="w-4 h-4 text-gray-400" />
									<span class="text-sm text-gray-400">Projects</span>
								</div>
								<span class="text-lg font-bold text-blue-400">{skill.projectCount}</span>
							</div>
						</div>
					</a>
				{/each}
			</div>
		{/if}
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

	:global(.animate-fadeInUp) {
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}
</style>

