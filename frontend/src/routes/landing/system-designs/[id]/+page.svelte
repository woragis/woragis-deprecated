<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Code, Image } from 'lucide-svelte';
	import { getSystemDesign, type SystemDesign } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let design: SystemDesign | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const designId = $page.params.id;
		if (!designId) {
			error = 'System design ID is required';
			loading = false;
			return;
		}

		try {
			design = await getSystemDesign(designId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load system design';
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
		title={design?.title || 'System Design Details'}
		description={design?.description}
		gradientFrom="from-cyan-950/30"
		gradientVia="via-blue-950/30"
		gradientTo="to-cyan-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-cyan-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/system-designs')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to System Designs
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-cyan-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading system design...</p>
			</div>
		{:else if error || !design}
			<EmptyState
				title={error || 'System design not found'}
				description="The system design you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-cyan-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-cyan-700"
					onclick={() => goto('/landing/system-designs')}
				>
					Return to System Designs
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Created"
						value={formatDate(design.created_at)}
						accentColor="cyan"
					/>
					<StatCard
						label="Featured"
						value={design.featured ? 'Yes' : 'No'}
						accentColor={design.featured ? 'emerald' : 'slate'}
					/>
					{#if design.technologies && design.technologies.length > 0}
						<StatCard
							label="Technologies"
							value={design.technologies.length.toString()}
							accentColor="blue"
						/>
					{/if}
				</div>

				<!-- Description Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Description</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{design.description}</div>
					</div>
				</div>

				<!-- Content Card -->
				{#if design.content}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Content</h2>
						<div class="prose prose-invert max-w-none">
							<div class="text-slate-300 whitespace-pre-wrap">{design.content}</div>
						</div>
					</div>
				{/if}

				<!-- Diagram Card -->
				{#if design.diagram_url}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Diagram</h2>
						<a
							href={design.diagram_url}
							target="_blank"
							class="flex items-center gap-2 text-cyan-400 hover:text-cyan-300"
						>
							<Image class="h-4 w-4" />
							<span>View Diagram</span>
						</a>
					</div>
				{/if}

				<!-- Technologies Card -->
				{#if design.technologies && design.technologies.length > 0}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Technologies</h2>
						<div class="flex flex-wrap gap-2">
							{#each design.technologies as tech}
								<span
									class="inline-flex items-center rounded-full bg-cyan-500/20 px-3 py-1 text-xs font-medium text-cyan-300"
								>
									{tech}
								</span>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#if design.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{design.slug}</p>
							</div>
						{/if}
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(design.updated_at)}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

