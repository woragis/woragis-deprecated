<script lang="ts">
	import { Tag, ExternalLink } from 'lucide-svelte';
	import type { SkillWithCount } from '$lib/api/skills';

	interface Props {
		skills: SkillWithCount[];
		loading: boolean;
	}

	let { skills = [], loading = false }: Props = $props();
</script>

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
				{#each skills as skill (skill.id || skill.name)}
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

