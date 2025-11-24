<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { ArrowLeft, ExternalLink, TrendingUp, Lightbulb, Target, CheckCircle2 } from 'lucide-svelte';
	import { getCaseStudy } from '$lib/api/case-studies';
	import MermaidDiagram from '$lib/components/MermaidDiagram.svelte';
	import type { CaseStudy } from '$lib/types/case-study';

	let caseStudy: CaseStudy | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const id = $page.params.id;
		if (id) {
			await fetchCaseStudy(id);
		} else {
			loading = false;
		}
	});

	async function fetchCaseStudy(id: string) {
		loading = true;
		error = null;
		try {
			caseStudy = await getCaseStudy(id);
			if (!caseStudy) {
				error = 'Case study not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch case study';
			console.error('Error fetching case study:', err);
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
	{:else if error || !caseStudy}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">Case Study Not Found</h1>
				<p class="text-gray-400 mb-8">The case study you are looking for does not exist.</p>
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
			<div class="max-w-5xl mx-auto">
				<!-- Back Button -->
				<a
					href="/projects/{caseStudy.projectSlug}"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to Project
				</a>

				<!-- Header -->
				<header class="mb-12">
					<h1 class="text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
						{caseStudy.title}
					</h1>
					{#if caseStudy.featured}
						<span
							class="inline-block px-3 py-1 bg-yellow-500/20 text-yellow-300 text-sm font-medium rounded-full border border-yellow-500/30 mb-4"
						>
							Featured Case Study
						</span>
					{/if}
					<div class="flex flex-wrap gap-2 mt-4">
						{#each caseStudy.technologies as tech}
							<span
								class="px-3 py-1 bg-blue-600/20 text-blue-300 text-sm rounded border border-blue-600/30"
							>
								{tech}
							</span>
						{/each}
					</div>
				</header>

				<!-- Problem Statement -->
				<section class="mb-12 bg-gradient-to-br from-red-900/20 to-orange-900/20 rounded-xl p-8 border border-red-700/30">
					<div class="flex items-start gap-4 mb-4">
						<Target class="w-8 h-8 text-red-400 flex-shrink-0 mt-1" />
						<div>
							<h2 class="text-3xl font-bold mb-4 text-red-300">Problem Statement</h2>
							<p class="text-gray-200 text-lg leading-relaxed">{caseStudy.problem}</p>
						</div>
					</div>
				</section>

				<!-- Context -->
				<section class="mb-12">
					<h2 class="text-3xl font-bold mb-6">Context</h2>
					<p class="text-gray-300 text-lg leading-relaxed">{caseStudy.context}</p>
				</section>

				<!-- Solution -->
				<section class="mb-12 bg-gradient-to-br from-green-900/20 to-emerald-900/20 rounded-xl p-8 border border-green-700/30">
					<div class="flex items-start gap-4 mb-4">
						<CheckCircle2 class="w-8 h-8 text-green-400 flex-shrink-0 mt-1" />
						<div>
							<h2 class="text-3xl font-bold mb-4 text-green-300">Solution</h2>
							<p class="text-gray-200 text-lg leading-relaxed mb-6">{caseStudy.solution}</p>

							<!-- Approach -->
							{#if caseStudy.approach && caseStudy.approach.length > 0}
								<div class="mt-6">
									<h3 class="text-xl font-semibold mb-4 text-green-200">Approach</h3>
									<ul class="space-y-3">
										{#each caseStudy.approach as step, index}
											<li class="flex items-start gap-3 text-gray-200">
												<span
													class="flex-shrink-0 w-6 h-6 rounded-full bg-green-600/30 border border-green-500/50 flex items-center justify-center text-green-300 text-sm font-medium mt-0.5"
												>
													{index + 1}
												</span>
												<span>{step}</span>
											</li>
										{/each}
									</ul>
								</div>
							{/if}
						</div>
					</div>
				</section>

				<!-- Architecture -->
				{#if caseStudy.architecture}
					<section class="mb-12">
						<h2 class="text-3xl font-bold mb-6">Architecture</h2>
						{#if caseStudy.architecture.description}
							<p class="text-gray-300 text-lg leading-relaxed mb-6">
								{caseStudy.architecture.description}
							</p>
						{/if}

						{#if caseStudy.architecture.diagram && caseStudy.architecture.diagramType === 'mermaid'}
							<div class="mb-8">
								<MermaidDiagram diagram={caseStudy.architecture.diagram} />
							</div>
						{/if}

						{#if caseStudy.architecture.components && caseStudy.architecture.components.length > 0}
							<div class="grid md:grid-cols-2 gap-6">
								{#each caseStudy.architecture.components as component}
									<div
										class="bg-gray-800/50 rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-colors"
									>
										<h4 class="text-xl font-bold mb-2 text-blue-400">{component.name}</h4>
										<p class="text-gray-300 mb-4">{component.description}</p>
										<div class="flex flex-wrap gap-2">
											{#each component.technologies as tech}
												<span
													class="px-2 py-1 bg-gray-700/50 text-gray-300 text-xs rounded border border-gray-600"
												>
													{tech}
												</span>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</section>
				{/if}

				<!-- Metrics & Impact -->
				{#if caseStudy.metrics}
					<section class="mb-12">
						<h2 class="text-3xl font-bold mb-6 flex items-center gap-3">
							<TrendingUp class="w-8 h-8 text-blue-400" />
							Metrics & Impact
						</h2>

						{#if caseStudy.metrics.before && caseStudy.metrics.before.length > 0}
							<div class="mb-8">
								<h3 class="text-xl font-semibold mb-4 text-gray-300">Before</h3>
								<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
									{#each caseStudy.metrics.before as metric}
										<div
											class="bg-gray-800/50 rounded-lg p-4 border border-gray-700"
										>
											<div class="text-sm text-gray-400 mb-1">{metric.label}</div>
											<div class="text-xl font-bold text-gray-200">{metric.value}</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						{#if caseStudy.metrics.after && caseStudy.metrics.after.length > 0}
							<div class="mb-8">
								<h3 class="text-xl font-semibold mb-4 text-green-300">After</h3>
								<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
									{#each caseStudy.metrics.after as metric}
										<div
											class="bg-green-900/20 rounded-lg p-4 border border-green-700/30"
										>
											<div class="text-sm text-green-400 mb-1">{metric.label}</div>
											<div class="text-xl font-bold text-green-300">{metric.value}</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						{#if caseStudy.metrics.impact}
							<div
								class="bg-gradient-to-br from-blue-900/20 to-purple-900/20 rounded-xl p-6 border border-blue-700/30"
							>
								<h3 class="text-xl font-semibold mb-3 text-blue-300">Impact</h3>
								<p class="text-gray-200 leading-relaxed">{caseStudy.metrics.impact}</p>
							</div>
						{/if}
					</section>
				{/if}

				<!-- Lessons Learned -->
				{#if caseStudy.lessonsLearned && caseStudy.lessonsLearned.length > 0}
					<section class="mb-12">
						<h2 class="text-3xl font-bold mb-6 flex items-center gap-3">
							<Lightbulb class="w-8 h-8 text-yellow-400" />
							Lessons Learned
						</h2>
						<div class="grid md:grid-cols-2 gap-4">
							{#each caseStudy.lessonsLearned as lesson}
								<div
									class="bg-yellow-900/20 rounded-lg p-4 border border-yellow-700/30 flex items-start gap-3"
								>
									<Lightbulb class="w-5 h-5 text-yellow-400 flex-shrink-0 mt-0.5" />
									<p class="text-gray-200">{lesson}</p>
								</div>
							{/each}
						</div>
					</section>
				{/if}

				<!-- Project Link -->
				<div class="mt-12 pt-8 border-t border-gray-700">
					<a
						href="/projects/{caseStudy.projectSlug}"
						class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						<ExternalLink class="w-5 h-5" />
						View Full Project Details
					</a>
				</div>
			</div>
		</div>
	{/if}
</div>

