<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getProjectBySlug } from '$lib/api/projects';
	import type { Project } from '$lib/types/project';

	let project: Project | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const slug = $page.params.slug;
		if (slug) {
			await fetchProject(slug);
		}
	});

	async function fetchProject(slug: string) {
		loading = true;
		error = null;
		try {
			project = await getProjectBySlug(slug);
			if (!project) {
				error = 'Project not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch project';
			console.error('Error fetching project:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !project}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">Project Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The project you are looking for does not exist.'}</p>
				<a
					href="/projects"
					class="inline-block px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					← Back to Projects
				</a>
			</div>
		</div>
	{:else}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<a
					href="/projects"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					← Back to Projects
				</a>

				<article class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 border border-gray-700 shadow-xl">
					<h1 class="text-4xl font-bold mb-4">{project.name}</h1>
					{#if project.description}
						<p class="text-xl text-gray-300 mb-6">{project.description}</p>
					{/if}

					<div class="grid md:grid-cols-2 gap-6 mt-8">
						<div>
							<h2 class="text-2xl font-bold mb-4">Details</h2>
							<dl class="space-y-3">
								<div>
									<dt class="text-sm text-gray-400">Status</dt>
									<dd class="text-lg font-medium capitalize">{project.status}</dd>
								</div>
								<div>
									<dt class="text-sm text-gray-400">Health Score</dt>
									<dd class="text-lg font-medium">{project.healthScore}%</dd>
								</div>
								<div>
									<dt class="text-sm text-gray-400">Created</dt>
									<dd class="text-lg font-medium">
										{new Date(project.createdAt).toLocaleDateString()}
									</dd>
								</div>
								<div>
									<dt class="text-sm text-gray-400">Last Updated</dt>
									<dd class="text-lg font-medium">
										{new Date(project.updatedAt).toLocaleDateString()}
									</dd>
								</div>
							</dl>
						</div>

						{#if project.skills && project.skills.length > 0}
							<div>
								<h2 class="text-2xl font-bold mb-4">Skills</h2>
								<div class="flex flex-wrap gap-2">
									{#each project.skills as skill}
										<a
											href="/skills"
											class="flex items-center gap-2 px-4 py-2 bg-blue-600/20 hover:bg-blue-600/30 rounded-lg text-blue-300 border border-blue-600/30 transition-colors"
										>
											{#if skill.icon}
												<span class="text-lg">{skill.icon}</span>
											{/if}
											<span class="font-medium">{skill.name}</span>
											<span class="text-xs text-blue-400/70 capitalize">({skill.category})</span>
										</a>
									{/each}
								</div>
							</div>
						{/if}

						{#if project.technologies && project.technologies.length > 0}
							<div>
								<h2 class="text-2xl font-bold mb-4">Technologies</h2>
								<ul class="space-y-2">
									{#each project.technologies as tech}
										<li class="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg">
											<div>
												<span class="font-medium">{tech.name}</span>
												{#if tech.version}
													<span class="text-gray-400 text-sm ml-2">v{tech.version}</span>
												{/if}
												{#if tech.category}
													<span class="text-xs text-gray-500 ml-2 capitalize">{tech.category}</span>
												{/if}
											</div>
											{#if tech.link}
												<a
													href={tech.link}
													target="_blank"
													rel="noopener noreferrer"
													class="text-blue-400 hover:text-blue-300"
												>
													→
												</a>
											{/if}
										</li>
									{/each}
								</ul>
							</div>
						{/if}
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>

