<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Code2, ExternalLink, Calendar, Target, TrendingUp, ArrowLeft, CheckCircle, Lightbulb } from 'lucide-svelte';
	import { getCaseStudy } from '$lib/api/case-studies';
	import type { CaseStudy } from '$lib/types/case-study';

	let caseStudy: CaseStudy | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const caseStudyId = $derived($page.params.slug);

	onMount(async () => {
		if (caseStudyId) {
			await fetchCaseStudy(caseStudyId);
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

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
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
				<p class="text-gray-400 mb-8">{error || 'The case study you are looking for does not exist.'}</p>
				<a
					href="/case-studies"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to Case Studies
				</a>
			</div>
		</div>
	{:else}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/case-studies"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to Case Studies
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-purple-500/0 via-pink-500/0 to-orange-500/0 hover:from-purple-500/5 hover:via-pink-500/5 hover:to-orange-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<div class="flex items-center gap-3 mb-4">
								<div
									class="w-12 h-12 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center"
								>
									<Target class="w-6 h-6 text-white" />
								</div>
								{#if caseStudy.featured}
									<span class="px-3 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded">
										⭐ Featured
									</span>
								{/if}
							</div>

							<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
								{caseStudy.title}
							</h1>

							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if caseStudy.updatedAt}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Updated {formatDate(caseStudy.updatedAt)}</span>
									</div>
								{/if}
								{#if caseStudy.projectSlug}
									<a
										href="/projects/{caseStudy.projectSlug}"
										class="flex items-center gap-2 text-purple-400 hover:text-purple-300 transition-colors"
									>
										<Code2 class="w-4 h-4" />
										View Project
									</a>
								{/if}
							</div>
						</div>

						<!-- Context -->
						{#if caseStudy.context}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Context</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed">{caseStudy.context}</p>
								</div>
							</div>
						{/if}

						<!-- Problem -->
						{#if caseStudy.problem}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Problem</h2>
								<div class="bg-red-900/20 rounded-lg p-6 border border-red-700/30">
									<p class="text-gray-300 leading-relaxed">{caseStudy.problem}</p>
								</div>
							</div>
						{/if}

						<!-- Solution -->
						{#if caseStudy.solution}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Solution</h2>
								<div class="bg-green-900/20 rounded-lg p-6 border border-green-700/30">
									<p class="text-gray-300 leading-relaxed">{caseStudy.solution}</p>
								</div>
							</div>
						{/if}

						<!-- Approach -->
						{#if caseStudy.approach && caseStudy.approach.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Approach</h2>
								<div class="space-y-3">
									{#each caseStudy.approach as step, index}
										<div class="flex items-start gap-3 bg-gray-800/50 rounded-lg p-4 border border-gray-700">
											<div
												class="flex-shrink-0 w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold"
											>
												{index + 1}
											</div>
											<p class="text-gray-300 leading-relaxed flex-1">{step}</p>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Architecture -->
						{#if caseStudy.architecture}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Architecture</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									{#if caseStudy.architecture.description}
										<p class="text-gray-300 leading-relaxed mb-4">{caseStudy.architecture.description}</p>
									{/if}

									{#if caseStudy.architecture.components && caseStudy.architecture.components.length > 0}
										<div class="grid md:grid-cols-2 gap-4 mt-6">
											{#each caseStudy.architecture.components as component}
												<div class="bg-gray-700/30 rounded-lg p-4 border border-gray-600">
													<h3 class="text-lg font-bold text-white mb-2">{component.name}</h3>
													<p class="text-gray-300 text-sm mb-3">{component.description}</p>
													{#if component.technologies && component.technologies.length > 0}
														<div class="flex flex-wrap gap-2">
															{#each component.technologies as tech}
																<span
																	class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded-lg border border-purple-500/30"
																>
																	{tech}
																</span>
															{/each}
														</div>
													{/if}
												</div>
											{/each}
										</div>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Metrics -->
						{#if caseStudy.metrics}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Metrics & Impact</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									{#if caseStudy.metrics.before && caseStudy.metrics.before.length > 0}
										<div class="mb-6">
											<h3 class="text-lg font-bold text-red-300 mb-3">Before</h3>
											<div class="grid md:grid-cols-2 gap-4">
												{#each caseStudy.metrics.before as metric}
													<div class="bg-red-900/20 rounded-lg p-4 border border-red-700/30">
														<p class="text-sm text-gray-400 mb-1">{metric.label}</p>
														<p class="text-xl font-bold text-white">{metric.value}</p>
													</div>
												{/each}
											</div>
										</div>
									{/if}

									{#if caseStudy.metrics.after && caseStudy.metrics.after.length > 0}
										<div class="mb-6">
											<h3 class="text-lg font-bold text-green-300 mb-3">After</h3>
											<div class="grid md:grid-cols-2 gap-4">
												{#each caseStudy.metrics.after as metric}
													<div class="bg-green-900/20 rounded-lg p-4 border border-green-700/30">
														<p class="text-sm text-gray-400 mb-1">{metric.label}</p>
														<p class="text-xl font-bold text-white">{metric.value}</p>
													</div>
												{/each}
											</div>
										</div>
									{/if}

									{#if caseStudy.metrics.impact}
										<div class="bg-blue-900/20 rounded-lg p-4 border border-blue-700/30">
											<div class="flex items-start gap-3">
												<TrendingUp class="w-5 h-5 text-blue-400 flex-shrink-0 mt-0.5" />
												<div>
													<p class="text-sm text-gray-400 mb-1">Impact</p>
													<p class="text-gray-300 leading-relaxed">{caseStudy.metrics.impact}</p>
												</div>
											</div>
										</div>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Lessons Learned -->
						{#if caseStudy.lessonsLearned && caseStudy.lessonsLearned.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Lessons Learned</h2>
								<div class="space-y-3">
									{#each caseStudy.lessonsLearned as lesson}
										<div class="flex items-start gap-3 bg-yellow-900/20 rounded-lg p-4 border border-yellow-700/30">
											<Lightbulb class="w-5 h-5 text-yellow-400 flex-shrink-0 mt-0.5" />
											<p class="text-gray-300 leading-relaxed">{lesson}</p>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Technologies -->
						{#if caseStudy.technologies && caseStudy.technologies.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Technologies Used</h2>
								<div class="flex flex-wrap gap-2">
									{#each caseStudy.technologies as tech}
										<span
											class="px-3 py-1 bg-purple-600/20 text-purple-300 text-sm rounded-lg border border-purple-500/30"
										>
											{tech}
										</span>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Project Link -->
						{#if caseStudy.projectSlug}
							<div class="mt-8 pt-8 border-t border-gray-700">
								<a
									href="/projects/{caseStudy.projectSlug}"
									class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 rounded-lg transition-all duration-200 font-medium shadow-lg shadow-purple-500/50"
								>
									<Code2 class="w-5 h-5" />
									View Related Project
									<ExternalLink class="w-5 h-5" />
								</a>
							</div>
						{/if}
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
