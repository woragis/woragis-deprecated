<script lang="ts">
	import { Target, CheckCircle2, XCircle, ArrowRight, ChevronDown, ChevronUp, Layers } from 'lucide-svelte';
	import { caseStudies } from '$lib/constants/technical';
	import type { TechnicalCaseStudy } from '$lib/types/technical';

	let expandedCaseStudy: string | null = $state(null);

	function toggleCaseStudy(id: string) {
		expandedCaseStudy = expandedCaseStudy === id ? null : id;
	}
</script>

<section id="system-thinking" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4">System Thinking</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				Architectural decisions, trade-offs, and lessons learned from building complex systems
			</p>
		</div>

		<div class="space-y-8">
			{#each caseStudies as study (study.id)}
				<div
					class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-8 border border-gray-700 hover:border-purple-500/50 transition-all duration-300"
				>
					<div class="flex items-start justify-between mb-6">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-3">
								<div
									class="w-12 h-12 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center"
								>
									<Target class="w-6 h-6" />
								</div>
								<div>
									<h3 class="text-2xl font-bold text-white">{study.title}</h3>
									<p class="text-sm text-gray-400">Case Study</p>
								</div>
							</div>
							<p class="text-gray-300 mb-4">{study.description}</p>
						</div>
						<button
							onclick={() => toggleCaseStudy(study.id)}
							class="text-gray-400 hover:text-white transition-colors ml-4"
						>
							{#if expandedCaseStudy === study.id}
								<ChevronUp class="w-6 h-6" />
							{:else}
								<ChevronDown class="w-6 h-6" />
							{/if}
						</button>
					</div>

					<div class="grid md:grid-cols-2 gap-4 mb-6">
						<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-4">
							<h4 class="text-sm font-semibold text-red-400 mb-2 flex items-center gap-2">
								<XCircle class="w-4 h-4" />
								Challenge
							</h4>
							<p class="text-sm text-gray-300">{study.challenge}</p>
						</div>
						<div class="bg-green-900/20 border border-green-700/30 rounded-lg p-4">
							<h4 class="text-sm font-semibold text-green-400 mb-2 flex items-center gap-2">
								<CheckCircle2 class="w-4 h-4" />
								Solution
							</h4>
							<p class="text-sm text-gray-300">{study.solution}</p>
						</div>
					</div>

					<div class="flex flex-wrap gap-2 mb-6">
						{#each study.technologies as tech (tech)}
							<span
								class="px-3 py-1 bg-purple-600/20 text-purple-300 rounded-full text-xs border border-purple-600/30"
							>
								{tech}
							</span>
						{/each}
					</div>

					{#if expandedCaseStudy === study.id}
						<div class="mt-6 pt-6 border-t border-gray-700 space-y-6">
							{#if study.architecture}
								<div>
									<h4 class="text-lg font-semibold text-blue-400 mb-3 flex items-center gap-2">
										<Layers class="w-5 h-5" />
										Architecture
									</h4>
									<p class="text-gray-300 leading-relaxed">{study.architecture}</p>
								</div>
							{/if}

							{#if study.metrics && study.metrics.length > 0}
								<div>
									<h4 class="text-lg font-semibold text-cyan-400 mb-3">Key Metrics</h4>
									<div class="grid md:grid-cols-3 gap-4">
										{#each study.metrics as metric (metric.label)}
											<div
												class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 rounded-lg p-4 border border-gray-700"
											>
												<p class="text-xs text-gray-400 mb-1">{metric.label}</p>
												<p class="text-2xl font-bold text-cyan-400">{metric.value}</p>
												{#if metric.improvement}
													<p class="text-xs text-gray-500 mt-1">{metric.improvement}</p>
												{/if}
											</div>
										{/each}
									</div>
								</div>
							{/if}

							{#if study.tradeoffs && study.tradeoffs.length > 0}
								<div>
									<h4 class="text-lg font-semibold text-yellow-400 mb-3">Trade-offs & Decisions</h4>
									<div class="space-y-4">
										{#each study.tradeoffs as tradeoff (tradeoff.decision)}
											<div
												class="bg-gray-800/30 rounded-lg p-4 border border-gray-700"
											>
												<h5 class="font-semibold text-white mb-3">{tradeoff.decision}</h5>
												<div class="grid md:grid-cols-2 gap-4">
													<div>
														<p class="text-xs font-semibold text-green-400 mb-2">Pros</p>
														<ul class="space-y-1">
															{#each tradeoff.pros as pro (pro)}
																<li class="text-sm text-gray-300 flex items-start gap-2">
																	<CheckCircle2 class="w-4 h-4 text-green-400 mt-0.5 flex-shrink-0" />
																	<span>{pro}</span>
																</li>
															{/each}
														</ul>
													</div>
													<div>
														<p class="text-xs font-semibold text-red-400 mb-2">Cons</p>
														<ul class="space-y-1">
															{#each tradeoff.cons as con (con)}
																<li class="text-sm text-gray-300 flex items-start gap-2">
																	<XCircle class="w-4 h-4 text-red-400 mt-0.5 flex-shrink-0" />
																	<span>{con}</span>
																</li>
															{/each}
														</ul>
													</div>
												</div>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							{#if study.lessonsLearned && study.lessonsLearned.length > 0}
								<div>
									<h4 class="text-lg font-semibold text-purple-400 mb-3">Lessons Learned</h4>
									<ul class="space-y-2">
										{#each study.lessonsLearned as lesson (lesson)}
											<li class="text-gray-300 flex items-start gap-3">
												<ArrowRight class="w-5 h-5 text-purple-400 mt-0.5 flex-shrink-0" />
												<span>{lesson}</span>
											</li>
										{/each}
									</ul>
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>

