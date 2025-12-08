<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, Award, ExternalLink, FileText } from 'lucide-svelte';
	import { getCertification, type Certification } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let certification: Certification | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const certId = $page.params.id;
		if (!certId) {
			error = 'Certification ID is required';
			loading = false;
			return;
		}

		try {
			certification = await getCertification(certId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load certification';
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
		title={certification?.name || 'Certification Details'}
		description={certification?.description}
		gradientFrom="from-yellow-950/30"
		gradientVia="via-amber-950/30"
		gradientTo="to-yellow-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-yellow-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/certifications')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Certifications
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-yellow-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading certification...</p>
			</div>
		{:else if error || !certification}
			<EmptyState
				title={error || 'Certification not found'}
				description="The certification you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-yellow-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-yellow-700"
					onclick={() => goto('/landing/certifications')}
				>
					Return to Certifications
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Status"
						value={certification.status}
						accentColor="yellow"
					/>
					<StatCard
						label="Issue Date"
						value={formatDate(certification.issue_date)}
						accentColor="amber"
					/>
					<StatCard
						label="Featured"
						value={certification.featured ? 'Yes' : 'No'}
						accentColor={certification.featured ? 'emerald' : 'slate'}
					/>
				</div>

				<!-- Details Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Details</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Issuer</p>
							<p class="text-sm font-medium text-slate-200">{certification.issuer}</p>
						</div>
						{#if certification.category}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Category</p>
								<p class="text-sm font-medium text-slate-200">{certification.category}</p>
							</div>
						{/if}
						{#if certification.credential_id}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Credential ID</p>
								<p class="text-sm font-medium text-slate-200">{certification.credential_id}</p>
							</div>
						{/if}
						{#if certification.expiry_date}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Expiry Date</p>
								<p class="text-sm font-medium text-slate-200">{formatDate(certification.expiry_date)}</p>
							</div>
						{/if}
					</div>
				</div>

				<!-- Description Card -->
				{#if certification.description}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Description</h2>
						<div class="prose prose-invert max-w-none">
							<div class="text-slate-300 whitespace-pre-wrap">{certification.description}</div>
						</div>
					</div>
				{/if}

				<!-- Links Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Links</h2>
					<div class="space-y-3">
						{#if certification.verification_url}
							<a
								href={certification.verification_url}
								target="_blank"
								class="flex items-center gap-2 text-yellow-400 hover:text-yellow-300"
							>
								<ExternalLink class="h-4 w-4" />
								<span>Verification URL</span>
							</a>
						{/if}
						{#if certification.certificate_url}
							<a
								href={certification.certificate_url}
								target="_blank"
								class="flex items-center gap-2 text-yellow-400 hover:text-yellow-300"
							>
								<FileText class="h-4 w-4" />
								<span>Certificate URL</span>
							</a>
						{/if}
					</div>
				</div>

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#if certification.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{certification.slug}</p>
							</div>
						{/if}
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Updated</p>
							<p class="text-sm font-medium text-slate-200">{formatDate(certification.updated_at)}</p>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

