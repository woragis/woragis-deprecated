<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Code, Award, Clock } from 'lucide-svelte';
	import { getSkill, type Skill } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let skill: Skill | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const skillId = $page.params.id;
		if (!skillId) {
			error = 'Skill ID is required';
			loading = false;
			return;
		}

		try {
			skill = await getSkill(skillId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load skill';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<PageHero
		title={skill?.name || 'Skill Details'}
		description={skill?.description}
		gradientFrom="from-indigo-950/30"
		gradientVia="via-blue-950/30"
		gradientTo="to-indigo-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-indigo-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/skills')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Skills
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-indigo-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading skill...</p>
			</div>
		{:else if error || !skill}
			<EmptyState
				title={error || 'Skill not found'}
				description="The skill you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-indigo-700"
					onclick={() => goto('/landing/skills')}
				>
					Return to Skills
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Created"
						value={formatDate(skill.created_at)}
						accentColor="indigo"
					/>
					<StatCard
						label="Featured"
						value={skill.featured ? 'Yes' : 'No'}
						accentColor={skill.featured ? 'emerald' : 'slate'}
					/>
					{#if skill.years_of_experience}
						<StatCard
							label="Experience"
							value={`${skill.years_of_experience} years`}
							accentColor="blue"
						/>
					{/if}
				</div>

				<!-- Description Card -->
				{#if skill.description}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Description</h2>
						<div class="prose prose-invert max-w-none">
							<div class="text-slate-300 whitespace-pre-wrap">{skill.description}</div>
						</div>
					</div>
				{/if}

				<!-- Details Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Details</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#if skill.category}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Category</p>
								<p class="text-sm font-medium text-slate-200">{skill.category}</p>
							</div>
						{/if}
						{#if skill.proficiency_level}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Proficiency Level</p>
								<p class="text-sm font-medium text-slate-200 capitalize">{skill.proficiency_level}</p>
							</div>
						{/if}
						{#if skill.years_of_experience}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Years of Experience</p>
								<p class="text-sm font-medium text-slate-200">{skill.years_of_experience} years</p>
							</div>
						{/if}
						{#if skill.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{skill.slug}</p>
							</div>
						{/if}
					</div>
				</div>

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(skill.updated_at)}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

