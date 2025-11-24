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
		<div class="flex items-center justify-between mb-12">
			<h2 class="text-4xl font-bold">{t('projects.title')}</h2>
			<a
				href="/projects"
				class="text-blue-400 hover:text-blue-300 transition-colors duration-200 flex items-center gap-2"
			>
				{t('projects.viewAll')}
				<ExternalLink class="w-5 h-5" />
			</a>
		</div>
		<ProjectsShowcase projects={projects} caseStudyMap={caseStudyMap} loading={loading} />
	</div>
</section>

