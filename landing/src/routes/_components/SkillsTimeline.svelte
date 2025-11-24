<script lang="ts">
	import { Calendar, Code2, TrendingUp, Star } from 'lucide-svelte';
	import { useSkillsTimelineQuery } from '$lib/queries/skills';
	import type { Skill, ProficiencyLevel } from '$lib/api/skills';

	const timelineQuery = useSkillsTimelineQuery();
	let skills = $derived(timelineQuery.data || []);
	let loading = $derived(timelineQuery.isPending);

	// Group skills by year
	let groupedByYear = $derived.by(() => {
		const grouped: Record<string, Skill[]> = {};

		skills.forEach((skill) => {
			if (skill.firstUsedDate) {
				const date = new Date(skill.firstUsedDate);
				const year = date.getFullYear().toString();
				if (!grouped[year]) {
					grouped[year] = [];
				}
				grouped[year].push(skill);
			}
		});

		// Sort years descending
		return Object.entries(grouped)
			.sort(([a], [b]) => parseInt(b) - parseInt(a))
			.map(([year, skills]) => ({
				year: parseInt(year),
				skills: skills.sort((a, b) => {
					const dateA = a.firstUsedDate ? new Date(a.firstUsedDate).getTime() : 0;
					const dateB = b.firstUsedDate ? new Date(b.firstUsedDate).getTime() : 0;
					return dateA - dateB;
				})
			}));
	});

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

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

	function getCategoryColor(category: string): string {
		const colors: Record<string, string> = {
			backend: 'from-blue-500 to-cyan-500',
			frontend: 'from-purple-500 to-pink-500',
			database: 'from-green-500 to-emerald-500',
			infrastructure: 'from-orange-500 to-red-500',
			devops: 'from-indigo-500 to-purple-500',
			language: 'from-yellow-500 to-amber-500',
			framework: 'from-teal-500 to-cyan-500',
			tool: 'from-gray-500 to-slate-500',
			service: 'from-pink-500 to-rose-500',
			library: 'from-violet-500 to-purple-500',
			other: 'from-gray-500 to-slate-500'
		};
		return colors[category] || colors.other;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if skills.length === 0}
		<div class="text-center py-20">
			<Calendar class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No timeline data available</p>
			<p class="text-gray-500 text-sm">Skills with first used dates will appear here</p>
		</div>
	{:else}
		<div class="relative">
			<!-- Timeline Line -->
			<div class="absolute left-8 top-0 bottom-0 w-0.5 bg-gradient-to-b from-blue-500 via-purple-500 to-pink-500"></div>

			<div class="space-y-12">
				{#each groupedByYear as { year, skills: yearSkills }}
					<div class="relative pl-20">
						<!-- Year Header -->
						<div class="flex items-center gap-4 mb-6">
							<div
								class="absolute left-0 w-16 h-16 bg-gradient-to-br from-blue-600 to-purple-600 rounded-full flex items-center justify-center border-4 border-gray-900 shadow-lg z-10"
							>
								<span class="text-white font-bold text-lg">{year}</span>
							</div>
							<div>
								<h3 class="text-3xl font-bold text-white">{year}</h3>
								<p class="text-gray-400">
									{yearSkills.length} {yearSkills.length === 1 ? 'skill' : 'skills'} learned
								</p>
							</div>
						</div>

						<!-- Skills for this year -->
						<div class="space-y-4">
							{#each yearSkills as skill}
								<div
									class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-5 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 ml-4"
								>
									<div class="flex items-start justify-between gap-4">
										<div class="flex-1">
											<div class="flex items-center gap-3 mb-2">
												<div
													class="w-3 h-3 rounded-full bg-gradient-to-br {getCategoryColor(skill.category)} border-2 border-gray-900"
												></div>
												<h4 class="text-lg font-bold text-white">{skill.name}</h4>
												{#if skill.proficiencyLevel}
													<span
														class="px-2 py-0.5 text-xs rounded-full border {getProficiencyColor(skill.proficiencyLevel)}"
													>
														{getProficiencyLabel(skill.proficiencyLevel)}
													</span>
												{/if}
											</div>

											<div class="flex items-center gap-4 text-sm text-gray-400 mb-2">
												<span
													class="px-2 py-1 rounded bg-gray-700/50 text-gray-300 capitalize"
												>
													{skill.category}
												</span>
												<div class="flex items-center gap-1">
													<Calendar class="w-4 h-4" />
													<span>{formatDate(skill.firstUsedDate)}</span>
												</div>
												{#if skill.yearsOfExperience}
													<div class="flex items-center gap-1">
														<TrendingUp class="w-4 h-4" />
														<span>{skill.yearsOfExperience} {skill.yearsOfExperience === 1 ? 'year' : 'years'}</span>
													</div>
												{/if}
											</div>

											{#if skill.description}
												<p class="text-sm text-gray-300 line-clamp-2">{skill.description}</p>
											{/if}
										</div>

										{#if skill.icon}
											<div class="w-12 h-12 bg-gray-700 rounded-lg flex items-center justify-center flex-shrink-0">
												<span class="text-xl">{skill.icon}</span>
											</div>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

