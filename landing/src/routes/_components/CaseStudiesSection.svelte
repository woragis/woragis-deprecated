<script lang="ts">
	import { Calendar, ExternalLink, Star, Code2, Target, TrendingUp } from 'lucide-svelte';
	import type { CaseStudy } from '$lib/types/case-study';

	interface Props {
		featuredCaseStudies?: CaseStudy[];
		latestCaseStudies?: CaseStudy[];
		loading?: boolean;
	}

	let { featuredCaseStudies = [], latestCaseStudies = [], loading = false }: Props = $props();

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getCaseStudyUrl(caseStudy: CaseStudy): string {
		return `/case-studies/${caseStudy.id}`;
	}

	function getProjectUrl(caseStudy: CaseStudy): string {
		return `/projects/${caseStudy.projectSlug}`;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if featuredCaseStudies.length === 0 && latestCaseStudies.length === 0}
		<div class="text-center py-20">
			<Code2 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No case studies available</p>
			<p class="text-gray-500 text-sm">Check back later for new case studies</p>
		</div>
	{:else}
		<!-- Featured Case Studies -->
		{#if featuredCaseStudies.length > 0}
			<div class="mb-12">
				<div class="flex items-center gap-2 mb-6">
					<Star class="w-5 h-5 text-yellow-400 fill-yellow-400" />
					<h3 class="text-2xl font-bold text-white">Featured Case Studies</h3>
				</div>
				<div class="grid md:grid-cols-3 gap-6">
					{#each featuredCaseStudies as caseStudy}
						<a
							href={getCaseStudyUrl(caseStudy)}
							class="group bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl overflow-hidden border border-gray-700 hover:border-purple-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-purple-500/20"
						>
							<div class="p-6">
								<div class="flex items-start justify-between mb-3">
									<div
										class="w-12 h-12 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center flex-shrink-0"
									>
										<Target class="w-6 h-6 text-white" />
									</div>
									<div
										class="px-2 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded flex items-center gap-1"
									>
										<Star class="w-3 h-3 fill-current" />
										Featured
									</div>
								</div>
								<h4
									class="text-xl font-bold text-white mb-2 group-hover:text-purple-400 transition-colors line-clamp-2"
								>
									{caseStudy.title}
								</h4>
								{#if caseStudy.context}
									<p class="text-gray-300 text-sm mb-4 line-clamp-3">{caseStudy.context}</p>
								{:else if caseStudy.problem}
									<p class="text-gray-300 text-sm mb-4 line-clamp-3">{caseStudy.problem}</p>
								{/if}
								<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
									{#if caseStudy.updatedAt}
										<div class="flex items-center gap-1">
											<Calendar class="w-3 h-3" />
											<span>Updated {formatDate(caseStudy.updatedAt)}</span>
										</div>
									{/if}
								</div>
								{#if caseStudy.technologies && caseStudy.technologies.length > 0}
									<div class="flex flex-wrap gap-2 mb-3">
										{#each caseStudy.technologies.slice(0, 3) as tech}
											<span
												class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded border border-purple-600/30"
											>
												{tech}
											</span>
										{/each}
									</div>
								{/if}
								<div class="flex items-center justify-between pt-4 border-t border-gray-700">
									<div class="flex items-center gap-2 text-purple-400 text-sm font-medium group-hover:gap-3 transition-all">
										<span>View Case Study</span>
										<ExternalLink class="w-4 h-4" />
									</div>
									{#if caseStudy.projectSlug}
										<a
											href={getProjectUrl(caseStudy)}
											onclick={(e) => e.stopPropagation()}
											class="text-xs text-gray-400 hover:text-blue-400 transition-colors"
										>
											View Project
										</a>
									{/if}
								</div>
							</div>
						</a>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Latest Case Studies -->
		{#if latestCaseStudies.length > 0}
			<div>
				<h3 class="text-2xl font-bold text-white mb-6">Latest Case Studies</h3>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each latestCaseStudies.slice(0, 6) as caseStudy}
						<a
							href={getCaseStudyUrl(caseStudy)}
							class="group bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-purple-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-purple-500/20"
						>
							<div class="flex items-start gap-3 mb-3">
								<div
									class="w-10 h-10 bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center flex-shrink-0"
								>
									<Code2 class="w-5 h-5 text-white" />
								</div>
								<div class="flex-1">
									<h4
										class="text-lg font-bold text-white mb-2 group-hover:text-purple-400 transition-colors line-clamp-2"
									>
										{caseStudy.title}
									</h4>
								</div>
							</div>
							{#if caseStudy.context}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{caseStudy.context}</p>
							{:else if caseStudy.problem}
								<p class="text-gray-300 text-sm mb-4 line-clamp-2">{caseStudy.problem}</p>
							{/if}
							<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
								{#if caseStudy.updatedAt}
									<div class="flex items-center gap-1">
										<Calendar class="w-3 h-3" />
										<span>{formatDate(caseStudy.updatedAt)}</span>
									</div>
								{/if}
							</div>
							{#if caseStudy.technologies && caseStudy.technologies.length > 0}
								<div class="flex flex-wrap gap-2 mb-3">
									{#each caseStudy.technologies.slice(0, 3) as tech}
										<span
											class="px-2 py-1 bg-purple-600/20 text-purple-300 text-xs rounded border border-purple-600/30"
										>
											{tech}
										</span>
									{/each}
								</div>
							{/if}
							<div class="flex items-center justify-between pt-4 border-t border-gray-700">
								<div class="flex items-center gap-2 text-purple-400 text-sm font-medium group-hover:gap-3 transition-all">
									<span>Read More</span>
									<ExternalLink class="w-4 h-4" />
								</div>
								{#if caseStudy.projectSlug}
									<a
										href={getProjectUrl(caseStudy)}
										onclick={(e) => e.stopPropagation()}
										class="text-xs text-gray-400 hover:text-blue-400 transition-colors"
									>
										Project
									</a>
								{/if}
							</div>
						</a>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>

