<script lang="ts">
	import { Award, ExternalLink, Calendar, CheckCircle2, Clock, XCircle, Building2 } from 'lucide-svelte';
	import type { Certification, CertificationCategory } from '$lib/types/certification';

	interface Props {
		certifications: Certification[];
		loading?: boolean;
	}

	let { certifications = [], loading = false }: Props = $props();

	// Group certifications by category
	let groupedCertifications = $derived.by(() => {
		const grouped: Record<CertificationCategory, Certification[]> = {
			cloud: [],
			security: [],
			programming: [],
			database: [],
			devops: [],
			architecture: [],
			other: []
		};

		certifications.forEach((cert) => {
			if (grouped[cert.category]) {
				grouped[cert.category].push(cert);
			} else {
				grouped.other.push(cert);
			}
		});

		// Sort each group by display order
		Object.keys(grouped).forEach((key) => {
			grouped[key as CertificationCategory].sort((a, b) => a.displayOrder - b.displayOrder);
		});

		// Filter out empty categories
		return Object.entries(grouped).filter(([_, certs]) => certs.length > 0);
	});

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
	}

	function isExpired(expiryDate?: string): boolean {
		if (!expiryDate) return false;
		return new Date(expiryDate) < new Date();
	}

	function daysUntilExpiry(expiryDate?: string): number | null {
		if (!expiryDate) return null;
		const days = Math.ceil((new Date(expiryDate).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));
		return days;
	}

	function getCategoryColor(category: CertificationCategory): string {
		const colors: Record<CertificationCategory, string> = {
			cloud: 'from-blue-500 to-cyan-500',
			security: 'from-red-500 to-orange-500',
			programming: 'from-purple-500 to-pink-500',
			database: 'from-green-500 to-emerald-500',
			devops: 'from-yellow-500 to-amber-500',
			architecture: 'from-indigo-500 to-purple-500',
			other: 'from-gray-500 to-slate-500'
		};
		return colors[category] || colors.other;
	}

	function getCategoryLabel(category: CertificationCategory): string {
		const labels: Record<CertificationCategory, string> = {
			cloud: 'Cloud',
			security: 'Security',
			programming: 'Programming',
			database: 'Database',
			devops: 'DevOps',
			architecture: 'Architecture',
			other: 'Other'
		};
		return labels[category] || 'Other';
	}

	function getStatusIcon(status: string, expiryDate?: string) {
		if (expiryDate && isExpired(expiryDate)) {
			return XCircle;
		}
		if (status === 'active') {
			return CheckCircle2;
		}
		return Clock;
	}

	function getStatusColor(status: string, expiryDate?: string): string {
		if (expiryDate && isExpired(expiryDate)) {
			return 'text-red-400';
		}
		if (status === 'active') {
			return 'text-green-400';
		}
		return 'text-yellow-400';
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if certifications.length === 0}
		<div class="text-center py-20">
			<Award class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No certifications available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="space-y-8">
			{#each groupedCertifications as [category, certs]}
				<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
					<!-- Category Header -->
					<div class="flex items-center gap-3 mb-6">
						<div
							class="w-12 h-12 bg-gradient-to-br {getCategoryColor(category as CertificationCategory)} rounded-lg flex items-center justify-center"
						>
							<Award class="w-6 h-6 text-white" />
						</div>
						<div>
							<h3 class="text-2xl font-bold text-white">{getCategoryLabel(category as CertificationCategory)}</h3>
							<p class="text-sm text-gray-400">{certs.length} certification{certs.length !== 1 ? 's' : ''}</p>
						</div>
					</div>

					<!-- Certifications Grid -->
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each certs as cert}
							<div
								class="bg-gray-800/50 rounded-lg p-5 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
							>
								<!-- Header -->
								<div class="flex items-start justify-between mb-3">
									<div class="flex-1">
										<h4 class="text-lg font-bold text-white mb-1">{cert.name}</h4>
										<div class="flex items-center gap-2 text-sm text-gray-400">
											<Building2 class="w-4 h-4" />
											<span>{cert.issuer}</span>
										</div>
									</div>
									{#if cert.status || cert.expiryDate}
										{@const StatusIcon = getStatusIcon(cert.status, cert.expiryDate)}
										<StatusIcon class="w-5 h-5 {getStatusColor(cert.status, cert.expiryDate)}" />
									{/if}
								</div>

								<!-- Description -->
								{#if cert.description}
									<p class="text-sm text-gray-300 mb-4 line-clamp-2">{cert.description}</p>
								{/if}

								<!-- Dates -->
								<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
									<div class="flex items-center gap-1">
										<Calendar class="w-3 h-3" />
										<span>Issued: {formatDate(cert.issueDate)}</span>
									</div>
									{#if cert.expiryDate}
										<div class="flex items-center gap-1 {isExpired(cert.expiryDate) ? 'text-red-400' : ''}">
											<Clock class="w-3 h-3" />
											<span>
												{#if isExpired(cert.expiryDate)}
													Expired: {formatDate(cert.expiryDate)}
												{:else}
													{@const days = daysUntilExpiry(cert.expiryDate)}
													{#if days !== null}
														{#if days < 30}
															Expires in {days} day{days !== 1 ? 's' : ''}
														{:else}
															Expires: {formatDate(cert.expiryDate)}
														{/if}
													{/if}
												{/if}
											</span>
										</div>
									{/if}
								</div>

								<!-- Skills -->
								{#if cert.skills && cert.skills.length > 0}
									<div class="mb-4">
										<div class="flex flex-wrap gap-2">
											{#each cert.skills as skill}
												<span
													class="px-2 py-1 text-xs rounded bg-blue-600/20 text-blue-300 border border-blue-500/30"
												>
													{skill.name}
												</span>
											{/each}
										</div>
									</div>
								{/if}

								<!-- Actions -->
								<div class="flex items-center gap-3 pt-4 border-t border-gray-700">
									{#if cert.verificationUrl}
										<a
											href={cert.verificationUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center gap-2 text-sm text-blue-400 hover:text-blue-300 transition-colors"
										>
											<ExternalLink class="w-4 h-4" />
											<span>Verify</span>
										</a>
									{/if}
									{#if cert.certificateUrl}
										<a
											href={cert.certificateUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="flex items-center gap-2 text-sm text-gray-400 hover:text-gray-300 transition-colors"
										>
											<ExternalLink class="w-4 h-4" />
											<span>View Certificate</span>
										</a>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

