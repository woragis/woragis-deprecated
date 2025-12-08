<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Users, Star } from 'lucide-svelte';
	import { getTestimonial, type Testimonial } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let testimonial: Testimonial | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const testimonialId = $page.params.id;
		if (!testimonialId) {
			error = 'Testimonial ID is required';
			loading = false;
			return;
		}

		try {
			testimonial = await getTestimonial(testimonialId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load testimonial';
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
		title="Testimonial"
		description={testimonial?.author_name ? `From ${testimonial.author_name}` : undefined}
		gradientFrom="from-pink-950/30"
		gradientVia="via-rose-950/30"
		gradientTo="to-pink-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-pink-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/testimonials')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Testimonials
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-pink-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading testimonial...</p>
			</div>
		{:else if error || !testimonial}
			<EmptyState
				title={error || 'Testimonial not found'}
				description="The testimonial you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-pink-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-pink-700"
					onclick={() => goto('/landing/testimonials')}
				>
					Return to Testimonials
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Created"
						value={formatDate(testimonial.created_at)}
						accentColor="pink"
					/>
					<StatCard
						label="Featured"
						value={testimonial.featured ? 'Yes' : 'No'}
						accentColor={testimonial.featured ? 'emerald' : 'slate'}
					/>
					{#if testimonial.rating}
						<StatCard
							label="Rating"
							value={`${testimonial.rating}/5`}
							accentColor="amber"
						/>
					{/if}
				</div>

				<!-- Content Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Testimonial</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap text-lg italic">"{testimonial.content}"</div>
					</div>
				</div>

				<!-- Author Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Author</h2>
					<div class="flex items-start gap-4">
						{#if testimonial.author_image}
							<img
								src={testimonial.author_image}
								alt={testimonial.author_name}
								class="h-16 w-16 rounded-full object-cover"
							/>
						{/if}
						<div class="flex-1">
							<p class="text-lg font-semibold text-slate-200">{testimonial.author_name}</p>
							{#if testimonial.author_title}
								<p class="text-sm text-slate-400">{testimonial.author_title}</p>
							{/if}
							{#if testimonial.author_company}
								<p class="text-sm text-slate-400">{testimonial.author_company}</p>
							{/if}
						</div>
					</div>
				</div>

				<!-- Rating Card -->
				{#if testimonial.rating}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Rating</h2>
						<div class="flex items-center gap-2">
							{#each Array(testimonial.rating) as _}
								<Star class="h-5 w-5 fill-amber-400 text-amber-400" />
							{/each}
							{#each Array(5 - testimonial.rating) as _}
								<Star class="h-5 w-5 fill-slate-700 text-slate-700" />
							{/each}
							<span class="ml-2 text-sm text-slate-400">({testimonial.rating}/5)</span>
						</div>
					</div>
				{/if}

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(testimonial.updated_at)}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

