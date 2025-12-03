<script lang="ts">
	import { ExternalLink } from 'lucide-svelte';
	import ProjectsShowcase from '../_components/ProjectsShowcase.svelte';
	import type { Project } from '$lib/types/project';
	import type { CaseStudy } from '$lib/types/case-study';

	interface Props {
		projects: Project[];
		caseStudyMap: Map<string, CaseStudy>;
		loading: boolean;
		translations: any;
	}

	let { projects, caseStudyMap, loading, translations: t }: Props = $props();
</script>

<section id="projects" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('projects.title')}
			</h2>
			<p class="text-xl text-gray-300 mb-6">{t('projects.showcase')}</p>
			<a
				href="/projects"
				class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors duration-200 shadow-lg shadow-blue-500/50"
			>
				{t('projects.viewAll')}
				<ExternalLink class="w-5 h-5" />
			</a>
		</div>
		<ProjectsShowcase projects={projects} caseStudyMap={caseStudyMap} loading={loading} />
	</div>
</section>

