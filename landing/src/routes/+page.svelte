<script lang="ts">
	import AboutMeSection from './_components/AboutMeSection.svelte';
	import ExperienceSection from './_sections/ExperienceSection.svelte';
	import VolunteerWorkSection from './_sections/VolunteerWorkSection.svelte';
	import SkillsSection from './_components/SkillsSection.svelte';
	import ProjectsSection from './_sections/ProjectsSection.svelte';
	import BlogSection from './_sections/BlogSection.svelte';
	import TechnicalWritingsSection from './_sections/TechnicalWritingsSection.svelte';
	import CaseStudiesSection from './_sections/CaseStudiesSection.svelte';
	import SkillsTimelineSection from './_sections/SkillsTimelineSection.svelte';
	import SystemDesignsSection from './_sections/SystemDesignsSection.svelte';
	import ProblemSolutionsSection from './_sections/ProblemSolutionsSection.svelte';
	import ProblemSolutionMatrixSection from './_sections/ProblemSolutionMatrixSection.svelte';
	import SystemThinkingSection from './_sections/SystemThinkingSection.svelte';
	import InterestsSection from './_sections/InterestsSection.svelte';
	import TestimonialsSection from './_sections/TestimonialsSection.svelte';
	import CertificationsSection from './_sections/CertificationsSection.svelte';
	import SocialMediaSection from './_sections/SocialMediaSection.svelte';
	import AIMLIntegrationsSection from './_sections/AIMLIntegrationsSection.svelte';
	import ImpactMetricsSection from './_sections/ImpactMetricsSection.svelte';
	import ContactSection from './_sections/ContactSection.svelte';
	import FooterSection from './_sections/FooterSection.svelte';
	import LanguagesSection from './_sections/LanguagesSection.svelte';
	import { translationsStore, language } from '$lib/i18n';
	import { useProjectsQuery, projectKeys } from '$lib/queries/projects';
	import { useSkillsWithCountsQuery, skillKeys } from '$lib/queries/skills';
	import { usePostsQuery, postKeys } from '$lib/queries/posts';
	import { useCaseStudiesQuery, caseStudyKeys } from '$lib/queries/case-studies';
	import { useFeaturedCertificationsQuery, certificationKeys } from '$lib/queries/certifications';
	import { useFeaturedTechnicalWritingsQuery, technicalWritingKeys } from '$lib/queries/technical-writings';
	import { useSocialMediaPostsQuery, socialMediaPostKeys } from '$lib/queries/social-media-posts';
	import { useFeaturedAIMLIntegrationsQuery, aimlIntegrationKeys } from '$lib/queries/aiml-integrations';
	import { useFeaturedImpactMetricsQuery, impactMetricKeys } from '$lib/queries/impact-metrics';
	import { useFeaturedSystemDesignsQuery, systemDesignKeys } from '$lib/queries/system-designs';
	import { useFeaturedProblemSolutionsQuery, problemSolutionKeys } from '$lib/queries/problem-solutions';
	import { interestKeys } from '$lib/queries/interests';
	import { listTestimonials } from '$lib/api/testimonials';
	import { listCaseStudies } from '$lib/api/case-studies';
	import { onMount } from 'svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import type { Testimonial } from '$lib/types/testimonial';
	import type { CaseStudy } from '$lib/types/case-study';
	import type { SkillWithCount } from '$lib/api/skills';

	// Reactive translation helper
	let t = $derived($translationsStore);

	// Get query client for invalidating queries
	const queryClient = useQueryClient();

	// Featured projects - using TanStack Query - reactive to language changes
	const featuredProjectsQuery = useProjectsQuery({ limit: 50, sortBy: 'updatedAt', sortOrder: 'desc' }, $language);

	// Invalidate all language-dependent queries when language changes
	$effect(() => {
		// Track language changes
		const currentLang = $language;
		// Invalidate all queries that depend on language
		queryClient.invalidateQueries({ queryKey: projectKeys.all });
		queryClient.invalidateQueries({ queryKey: postKeys.all });
		queryClient.invalidateQueries({ queryKey: caseStudyKeys.all });
		queryClient.invalidateQueries({ queryKey: certificationKeys.all });
		queryClient.invalidateQueries({ queryKey: technicalWritingKeys.all });
		queryClient.invalidateQueries({ queryKey: socialMediaPostKeys.all });
		queryClient.invalidateQueries({ queryKey: aimlIntegrationKeys.all });
		queryClient.invalidateQueries({ queryKey: impactMetricKeys.all });
		queryClient.invalidateQueries({ queryKey: systemDesignKeys.all });
		queryClient.invalidateQueries({ queryKey: problemSolutionKeys.all });
		queryClient.invalidateQueries({ queryKey: skillKeys.all });
		queryClient.invalidateQueries({ queryKey: interestKeys.all });
	});

	let projects = $derived(featuredProjectsQuery.data || []);
	let loadingProjects = $derived(featuredProjectsQuery.isPending);

	// Case studies for projects
	let caseStudyMap: Map<string, CaseStudy> = $state(new Map());
	let loadingCaseStudies = $state(false);

	// Blog posts queries - reactive to language changes
	const featuredPostsQuery = usePostsQuery({
		status: 'published',
		featured: true,
		limit: 3,
		orderBy: 'publishedAt',
		order: 'desc'
	}, $language);

	const latestPostsQuery = usePostsQuery({
		status: 'published',
		limit: 6,
		orderBy: 'publishedAt',
		order: 'desc'
	}, $language);

	let featuredPosts = $derived(featuredPostsQuery.data || []);
	let latestPosts = $derived(latestPostsQuery.data?.slice(0, 6) || []);
	let loadingPosts = $derived(featuredPostsQuery.isPending || latestPostsQuery.isPending);

	// Case studies queries - reactive to language changes
	const featuredCaseStudiesQuery = useCaseStudiesQuery({
		featured: true,
		limit: 3,
		orderBy: 'updatedAt',
		order: 'desc'
	}, $language);

	const latestCaseStudiesQuery = useCaseStudiesQuery({
		limit: 6,
		orderBy: 'updatedAt',
		order: 'desc'
	}, $language);

	let featuredCaseStudies = $derived(featuredCaseStudiesQuery.data || []);
	let latestCaseStudies = $derived(
		latestCaseStudiesQuery.data?.filter(
			(cs) => !featuredCaseStudies.some((fcs) => fcs.id === cs.id)
		) || []
	);
	let loadingCaseStudiesSection = $derived(featuredCaseStudiesQuery.isPending || latestCaseStudiesQuery.isPending);

	// Certifications query - reactive to language changes
	const certificationsQuery = useFeaturedCertificationsQuery($language);
	let certifications = $derived(certificationsQuery.data || []);
	let loadingCertifications = $derived(certificationsQuery.isPending);

	// Technical writings query - reactive to language changes
	const technicalWritingsQuery = useFeaturedTechnicalWritingsQuery($language);
	let technicalWritings = $derived(technicalWritingsQuery.data || []);
	let loadingTechnicalWritings = $derived(technicalWritingsQuery.isPending);

	// Social media posts query - reactive to language changes
	const socialMediaPostsQuery = useSocialMediaPostsQuery({ status: 'active' }, $language);
	let socialMediaPosts = $derived(socialMediaPostsQuery.data || []);
	let loadingSocialMediaPosts = $derived(socialMediaPostsQuery.isPending);

	// AI/ML integrations query - reactive to language changes
	const aimlIntegrationsQuery = useFeaturedAIMLIntegrationsQuery($language);
	let aimlIntegrations = $derived(aimlIntegrationsQuery.data || []);
	let loadingAIMLIntegrations = $derived(aimlIntegrationsQuery.isPending);

	// Impact metrics query - reactive to language changes
	const impactMetricsQuery = useFeaturedImpactMetricsQuery($language);
	let impactMetrics = $derived(impactMetricsQuery.data || []);
	let loadingImpactMetrics = $derived(impactMetricsQuery.isPending);

	// System designs query - reactive to language changes
	const systemDesignsQuery = useFeaturedSystemDesignsQuery($language);
	let systemDesigns = $derived(systemDesignsQuery.data || []);
	let loadingSystemDesigns = $derived(systemDesignsQuery.isPending);

	// Problem solutions query - reactive to language changes
	const problemSolutionsQuery = useFeaturedProblemSolutionsQuery($language);
	let problemSolutions = $derived(problemSolutionsQuery.data || []);
	let loadingProblemSolutions = $derived(problemSolutionsQuery.isPending);

	// Testimonials
	let testimonials: Testimonial[] = $state([]);
	let loadingTestimonials = $state(false);

	// Skills - using TanStack Query - reactive to language changes
	const skillsQuery = useSkillsWithCountsQuery($language);
	let popularSkills = $derived.by(() => {
		const allSkills = skillsQuery.data || [];
		return allSkills
			.filter((skill: SkillWithCount) => skill.projectCount > 0)
			.sort((a: SkillWithCount, b: SkillWithCount) => b.projectCount - a.projectCount)
			.slice(0, 6);
	});
	let loadingSkills = $derived(skillsQuery.isPending);

	// Fetch testimonials and case studies on mount
	onMount(async () => {
		// Fetch testimonials
		loadingTestimonials = true;
		try {
			const data = await listTestimonials({
				status: 'approved',
				orderBy: 'displayOrder',
				order: 'asc'
			});
			testimonials = data;
		} catch (error) {
			console.error('Error fetching testimonials:', error);
		} finally {
			loadingTestimonials = false;
		}

	// 	// Fetch case studies for project mapping
		loadingCaseStudies = true;
		try {
			const caseStudies = await listCaseStudies();
			const newMap = new Map<string, CaseStudy>();
			caseStudies.forEach((caseStudy) => {
				if (caseStudy.projectSlug) {
					newMap.set(caseStudy.projectSlug, caseStudy);
				}
			});
			caseStudyMap = newMap;
		} catch (error) {
			console.error('Error fetching case studies:', error);
		} finally {
			loadingCaseStudies = false;
		}
	});

</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- 1. About Me -->
	<AboutMeSection translations={t} />

	<!-- 2. Experience -->
	<ExperienceSection translations={t} />

	<!-- 3. Volunteer Work -->
	<VolunteerWorkSection translations={t} />

	<!-- 4. Skills -->
	<SkillsSection skills={popularSkills} loading={loadingSkills} />

	<!-- 5. Projects -->
	<ProjectsSection projects={projects} caseStudyMap={caseStudyMap} loading={loadingProjects} translations={t} />

	<!-- 6. Posts & Writings -->
	<BlogSection featuredPosts={featuredPosts} latestPosts={latestPosts} loading={loadingPosts} translations={t} />
	<TechnicalWritingsSection technicalWritings={technicalWritings} loading={loadingTechnicalWritings} translations={t} />

	<!-- 7. Languages -->
	<LanguagesSection translations={t} />

	<!-- Additional sections -->
	<CaseStudiesSection featuredCaseStudies={featuredCaseStudies} latestCaseStudies={latestCaseStudies} loading={loadingCaseStudiesSection} translations={t} />

	<SkillsTimelineSection translations={t} />

	<SystemDesignsSection systemDesigns={systemDesigns} loading={loadingSystemDesigns} translations={t} />

	<ProblemSolutionsSection problemSolutions={problemSolutions} loading={loadingProblemSolutions} translations={t} />

	<ProblemSolutionMatrixSection translations={t} />

	<SystemThinkingSection />

	<InterestsSection translations={t} />

	<TestimonialsSection testimonials={testimonials} loading={loadingTestimonials} translations={t} />

	<CertificationsSection certifications={certifications} loading={loadingCertifications} translations={t} />

	<SocialMediaSection socialMediaPosts={socialMediaPosts} loading={loadingSocialMediaPosts} translations={t} />

	<AIMLIntegrationsSection aimlIntegrations={aimlIntegrations} loading={loadingAIMLIntegrations} translations={t} />

	<ImpactMetricsSection impactMetrics={impactMetrics} loading={loadingImpactMetrics} translations={t} />

	<ContactSection translations={t} />

	<FooterSection translations={t} />
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	section {
		animation: fadeIn 0.6s ease-out;
	}

	:global(html) {
		scroll-behavior: smooth;
	}

</style>
