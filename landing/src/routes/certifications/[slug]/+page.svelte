<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Award, ExternalLink, Calendar, ArrowLeft, CheckCircle, XCircle, Clock, Star, Link as LinkIcon } from 'lucide-svelte';
	import { getCertification } from '$lib/api/certifications';
	import type { Certification, CertificationStatus, CertificationCategory } from '$lib/types/certification';

	let certification: Certification | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);

	const certificationId = $derived($page.params.slug);

	onMount(async () => {
		if (certificationId) {
			await fetchCertification(certificationId);
		}
	});

	async function fetchCertification(id: string) {
		loading = true;
		error = null;
		try {
			certification = await getCertification(id);
			if (!certification) {
				error = 'Certification not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch certification';
			console.error('Error fetching certification:', err);
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function getStatusColor(status: CertificationStatus): string {
		const colors: Record<CertificationStatus, string> = {
			active: 'bg-green-600/20 text-green-300 border-green-500/30',
			expired: 'bg-red-600/20 text-red-300 border-red-500/30',
			revoked: 'bg-gray-600/20 text-gray-300 border-gray-500/30',
			renewed: 'bg-blue-600/20 text-blue-300 border-blue-500/30'
		};
		return colors[status] || colors.active;
	}

	function getCategoryColor(category: CertificationCategory): string {
		const colors: Record<CertificationCategory, string> = {
			cloud: 'bg-blue-600/20 text-blue-300 border-blue-500/30',
			security: 'bg-red-600/20 text-red-300 border-red-500/30',
			programming: 'bg-purple-600/20 text-purple-300 border-purple-500/30',
			database: 'bg-yellow-600/20 text-yellow-300 border-yellow-500/30',
			devops: 'bg-green-600/20 text-green-300 border-green-500/30',
			architecture: 'bg-indigo-600/20 text-indigo-300 border-indigo-500/30',
			other: 'bg-gray-600/20 text-gray-300 border-gray-500/30'
		};
		return colors[category] || colors.other;
	}

	function getStatusIcon(status: CertificationStatus) {
		switch (status) {
			case 'active':
			case 'renewed':
				return CheckCircle;
			case 'expired':
			case 'revoked':
				return XCircle;
			default:
				return CheckCircle;
		}
	}

	function isExpiringSoon(expiryDate?: string): boolean {
		if (!expiryDate) return false;
		const expiry = new Date(expiryDate);
		const now = new Date();
		const daysUntilExpiry = Math.ceil((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
		return daysUntilExpiry > 0 && daysUntilExpiry <= 90;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	{#if loading}
		<div class="container mx-auto px-6 py-20">
			<div class="flex items-center justify-center">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else if error || !certification}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-2xl mx-auto text-center">
				<h1 class="text-4xl font-bold mb-4">Certification Not Found</h1>
				<p class="text-gray-400 mb-8">{error || 'The certification you are looking for does not exist.'}</p>
				<a
					href="/certifications"
					class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
				>
					<ArrowLeft class="w-5 h-5" />
					Back to Certifications
				</a>
			</div>
		</div>
	{:else}
		{@const StatusIcon = getStatusIcon(certification.status)}
		<div class="container mx-auto px-6 py-20">
			<div class="max-w-4xl mx-auto">
				<!-- Breadcrumb -->
				<a
					href="/certifications"
					class="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-8"
				>
					<ArrowLeft class="w-4 h-4" />
					Back to Certifications
				</a>

				<article class="bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl p-8 md:p-10 border border-gray-700 shadow-2xl relative overflow-hidden">
					<!-- Decorative gradient overlay -->
					<div class="absolute inset-0 bg-gradient-to-br from-green-500/0 via-emerald-500/0 to-teal-500/0 hover:from-green-500/5 hover:via-emerald-500/5 hover:to-teal-500/5 transition-all duration-300 pointer-events-none"></div>
					<div class="relative z-10">
						<!-- Header -->
						<div class="mb-8">
							<div class="flex items-center gap-3 mb-4">
								<div
									class="w-12 h-12 bg-gradient-to-br from-green-600 to-emerald-600 rounded-lg flex items-center justify-center"
								>
									<Award class="w-6 h-6 text-white" />
								</div>
								<div class="flex items-center gap-2 flex-wrap">
									<span
										class="px-3 py-1 rounded-lg text-sm font-medium border {getStatusColor(certification.status)} flex items-center gap-2"
									>
										<StatusIcon class="w-4 h-4" />
										{certification.status}
									</span>
									<span
										class="px-3 py-1 rounded-lg text-sm font-medium border {getCategoryColor(certification.category)}"
									>
										{certification.category}
									</span>
									{#if certification.featured}
										<span class="px-3 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded">
											⭐ Featured
										</span>
									{/if}
									{#if isExpiringSoon(certification.expiryDate)}
										<span class="px-3 py-1 bg-orange-500/90 text-orange-900 text-xs font-bold rounded flex items-center gap-1">
											<Clock class="w-3 h-3" />
											Expiring Soon
										</span>
									{/if}
								</div>
							</div>

							<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
								{certification.name}
							</h1>

							{#if certification.issuer}
								<p class="text-xl text-gray-300 mb-6">Issued by {certification.issuer}</p>
							{/if}

							<div class="flex flex-wrap items-center gap-6 text-sm text-gray-400 mb-6">
								{#if certification.issueDate}
									<div class="flex items-center gap-2">
										<Calendar class="w-4 h-4" />
										<span>Issued {formatDate(certification.issueDate)}</span>
									</div>
								{/if}
								{#if certification.expiryDate}
									<div class="flex items-center gap-2">
										<Clock class="w-4 h-4" />
										<span>Expires {formatDate(certification.expiryDate)}</span>
									</div>
								{/if}
							</div>
						</div>

						<!-- Description -->
						{#if certification.description}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Description</h2>
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-gray-300 leading-relaxed">{certification.description}</p>
								</div>
							</div>
						{/if}

						<!-- Credential Information -->
						<div class="grid md:grid-cols-2 gap-6 mb-8">
							{#if certification.credentialId}
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-sm text-gray-400 mb-1">Credential ID</p>
									<p class="text-lg font-bold text-white font-mono">{certification.credentialId}</p>
								</div>
							{/if}

							{#if certification.status}
								<div class="bg-gray-800/50 rounded-lg p-6 border border-gray-700">
									<p class="text-sm text-gray-400 mb-1">Status</p>
									<div class="flex items-center gap-2">
										<StatusIcon class="w-5 h-5 {certification.status === 'active' || certification.status === 'renewed' ? 'text-green-400' : 'text-red-400'}" />
										<p class="text-lg font-bold text-white capitalize">{certification.status}</p>
									</div>
								</div>
							{/if}
						</div>

						<!-- Skills -->
						{#if certification.skills && certification.skills.length > 0}
							<div class="mb-8">
								<h2 class="text-2xl font-bold text-white mb-4">Related Skills</h2>
								<div class="flex flex-wrap gap-2">
									{#each certification.skills as skill}
										<a
											href="/skills"
											class="px-3 py-1 bg-blue-600/20 text-blue-300 text-sm rounded-lg border border-blue-500/30 hover:bg-blue-600/30 transition-colors"
										>
											{skill.name}
										</a>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Links -->
						<div class="mt-8 pt-8 border-t border-gray-700 flex flex-wrap gap-4">
							{#if certification.verificationUrl}
								<a
									href={certification.verificationUrl}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200 font-medium"
								>
									<LinkIcon class="w-5 h-5" />
									Verify Credential
									<ExternalLink class="w-5 h-5" />
								</a>
							{/if}
							{#if certification.certificateUrl}
								<a
									href={certification.certificateUrl}
									target="_blank"
									rel="noopener noreferrer"
									class="inline-flex items-center gap-2 px-6 py-3 bg-green-600 hover:bg-green-700 rounded-lg transition-colors duration-200 font-medium"
								>
									<Award class="w-5 h-5" />
									View Certificate
									<ExternalLink class="w-5 h-5" />
								</a>
							{/if}
						</div>
					</div>
				</article>
			</div>
		</div>
	{/if}
</div>
