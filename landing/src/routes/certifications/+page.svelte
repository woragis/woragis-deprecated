<script lang="ts">
	import { Search, Filter, X, ChevronDown, ExternalLink, Calendar, Award, Star, CheckCircle, XCircle, Clock } from 'lucide-svelte';
	import type { Certification, CertificationStatus, CertificationCategory } from '$lib/types/certification';
	import { useCertificationsQuery } from '$lib/queries/certifications';

	// Filters
	let searchQuery = $state('');
	let statusFilter: CertificationStatus | 'all' = $state('all');
	let categoryFilter: CertificationCategory | 'all' = $state('all');
	let featuredFilter: boolean | 'all' = $state('all');
	let sortBy: 'issueDate' | 'name' | 'expiryDate' = $state('issueDate');
	let sortOrder: 'asc' | 'desc' = $state('desc');
	let showFilters = $state(false);

	// Fetch all certifications
	const certificationsQuery = useCertificationsQuery();
	let certifications = $derived(certificationsQuery.data || []);
	let filteredCertifications: Certification[] = $state([]);
	let loading = $derived(certificationsQuery.isPending);
	let error = $derived(certificationsQuery.error ? (certificationsQuery.error instanceof Error ? certificationsQuery.error.message : 'Failed to fetch certifications') : null);

	const statusOptions: Array<{ value: CertificationStatus | 'all'; label: string; icon: typeof CheckCircle }> = [
		{ value: 'all', label: 'All Status', icon: CheckCircle },
		{ value: 'active', label: 'Active', icon: CheckCircle },
		{ value: 'expired', label: 'Expired', icon: XCircle },
		{ value: 'revoked', label: 'Revoked', icon: XCircle },
		{ value: 'renewed', label: 'Renewed', icon: CheckCircle }
	];

	const categoryOptions: Array<{ value: CertificationCategory | 'all'; label: string }> = [
		{ value: 'all', label: 'All Categories' },
		{ value: 'cloud', label: 'Cloud' },
		{ value: 'security', label: 'Security' },
		{ value: 'programming', label: 'Programming' },
		{ value: 'database', label: 'Database' },
		{ value: 'devops', label: 'DevOps' },
		{ value: 'architecture', label: 'Architecture' },
		{ value: 'other', label: 'Other' }
	];

	const sortOptions = [
		{ value: 'issueDate', label: 'Issue Date' },
		{ value: 'expiryDate', label: 'Expiry Date' },
		{ value: 'name', label: 'Name' }
	];

	// Apply filters when certifications or filter values change
	$effect(() => {
		applyFilters();
	});

	function applyFilters() {
		let filtered = [...certifications];

		// Apply search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			filtered = filtered.filter(
				(cert: Certification) =>
					cert.name.toLowerCase().includes(query) ||
					cert.issuer?.toLowerCase().includes(query) ||
					cert.description?.toLowerCase().includes(query) ||
					cert.credentialId?.toLowerCase().includes(query)
			);
		}

		// Apply status filter
		if (statusFilter !== 'all') {
			filtered = filtered.filter((cert: Certification) => cert.status === statusFilter);
		}

		// Apply category filter
		if (categoryFilter !== 'all') {
			filtered = filtered.filter((cert: Certification) => cert.category === categoryFilter);
		}

		// Apply featured filter
		if (featuredFilter !== 'all') {
			filtered = filtered.filter((cert: Certification) => cert.featured === featuredFilter);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aVal: string | number | undefined;
			let bVal: string | number | undefined;

			switch (sortBy) {
				case 'issueDate':
					aVal = a.issueDate ? new Date(a.issueDate).getTime() : 0;
					bVal = b.issueDate ? new Date(b.issueDate).getTime() : 0;
					break;
				case 'expiryDate':
					aVal = a.expiryDate ? new Date(a.expiryDate).getTime() : 0;
					bVal = b.expiryDate ? new Date(b.expiryDate).getTime() : 0;
					break;
				case 'name':
					aVal = a.name.toLowerCase();
					bVal = b.name.toLowerCase();
					break;
				default:
					return 0;
			}

			if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
			if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
			return 0;
		});

		filteredCertifications = filtered;
	}

	function clearFilters() {
		searchQuery = '';
		statusFilter = 'all';
		categoryFilter = 'all';
		featuredFilter = 'all';
		applyFilters();
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
	}

	function getCertificationUrl(cert: Certification): string {
		return `/certifications/${cert.id}`;
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

	function isExpiringSoon(expiryDate?: string): boolean {
		if (!expiryDate) return false;
		const expiry = new Date(expiryDate);
		const now = new Date();
		const daysUntilExpiry = Math.ceil((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
		return daysUntilExpiry > 0 && daysUntilExpiry <= 90;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<!-- Header -->
	<section class="container mx-auto px-6 py-20">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-12">
				<h1 class="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
					Certifications
				</h1>
				<p class="text-xl text-gray-300 max-w-2xl mx-auto">
					Professional certifications and credentials demonstrating expertise and continuous learning
				</p>
			</div>

			<!-- Search and Filter Bar -->
			<div class="mb-8">
				<div class="flex flex-col md:flex-row gap-4 mb-4">
					<!-- Search Input -->
					<div class="flex-1 relative">
						<Search class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
						<input
							type="text"
							placeholder="Search certifications by name, issuer, credential ID..."
							bind:value={searchQuery}
							class="w-full pl-10 pr-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20"
						/>
					</div>
					<!-- Filter Toggle -->
					<div class="flex items-center justify-between gap-4">
						<button
							onclick={() => (showFilters = !showFilters)}
							class="flex items-center gap-2 px-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-white hover:bg-gray-700/50 transition-colors"
						>
							<Filter class="w-5 h-5" />
							Filters
							{#if showFilters}
								<ChevronDown class="w-4 h-4 rotate-180 transition-transform" />
							{:else}
								<ChevronDown class="w-4 h-4 transition-transform" />
							{/if}
						</button>

						<div class="flex items-center gap-4">
							<!-- Sort By -->
							<select
								bind:value={sortBy}
								onchange={applyFilters}
								class="px-4 py-2 bg-gray-800/50 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white text-sm"
							>
								{#each sortOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
							</select>

							<!-- Sort Order -->
							<button
								onclick={() => {
									sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
									applyFilters();
								}}
								class="px-4 py-2 bg-gray-800/50 hover:bg-gray-700/50 border border-gray-700 rounded-lg transition-colors text-sm"
								title="Toggle sort order"
							>
								{sortOrder === 'asc' ? '↑' : '↓'}
							</button>

							{#if searchQuery || statusFilter !== 'all' || categoryFilter !== 'all' || featuredFilter !== 'all'}
								<button
									onclick={clearFilters}
									class="flex items-center gap-2 px-4 py-2 bg-red-600/20 hover:bg-red-600/30 border border-red-700/30 rounded-lg transition-colors text-sm"
								>
									<X class="w-4 h-4" />
									Clear
								</button>
							{/if}
						</div>
					</div>
				</div>

				<!-- Filters Panel -->
				{#if showFilters}
					<div class="p-4 bg-gray-800/50 border border-gray-700 rounded-lg space-y-4">
						<!-- Status Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Status</div>
							<div class="flex flex-wrap gap-2">
								{#each statusOptions as option}
									<button
										onclick={() => {
											statusFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium flex items-center gap-2 {statusFilter === option.value
											? 'bg-blue-600 border-blue-500 text-white'
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										<svelte:component this={option.icon} class="w-4 h-4" />
										{option.label}
									</button>
								{/each}
							</div>
						</div>

						<!-- Category Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Category</div>
							<div class="flex flex-wrap gap-2">
								{#each categoryOptions as option}
									<button
										onclick={() => {
											categoryFilter = option.value;
											applyFilters();
										}}
										class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {categoryFilter === option.value
											? 'bg-purple-600 border-purple-500 text-white'
											: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
									>
										{option.label}
									</button>
								{/each}
							</div>
						</div>

						<!-- Featured Filter -->
						<div>
							<div class="block text-sm font-medium text-gray-300 mb-2">Featured</div>
							<div class="flex flex-wrap gap-2">
								<button
									onclick={() => {
										featuredFilter = 'all';
										applyFilters();
									}}
									class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {featuredFilter === 'all'
										? 'bg-blue-600 border-blue-500 text-white'
										: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
								>
									All
								</button>
								<button
									onclick={() => {
										featuredFilter = true;
										applyFilters();
									}}
									class="px-4 py-2 rounded-lg border transition-colors duration-200 text-sm font-medium {featuredFilter === true
										? 'bg-yellow-600 border-yellow-500 text-white'
										: 'bg-gray-700/50 border-gray-600 text-gray-300 hover:bg-gray-700'}"
								>
									Featured Only
								</button>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</section>

	<!-- Certifications Grid -->
	<section class="container mx-auto px-6 py-12">
		<div class="max-w-7xl mx-auto">
			{#if loading}
				<div class="flex items-center justify-center py-20">
					<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="bg-red-900/20 border border-red-700/30 rounded-lg p-6 text-center">
					<p class="text-red-400 mb-2">Error loading certifications</p>
					<p class="text-gray-400 text-sm">{error}</p>
					<button
						onclick={() => certificationsQuery.refetch()}
						class="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
					>
						Retry
					</button>
				</div>
			{:else if filteredCertifications.length === 0}
				<div class="text-center py-20">
					<Award class="w-16 h-16 mx-auto mb-4 text-gray-600" />
					<p class="text-gray-400 text-lg mb-2">No certifications found</p>
					<p class="text-gray-500 text-sm">Try adjusting your filters or search query</p>
				</div>
			{:else}
				<div class="mb-4 text-gray-400 text-sm">
					Showing {filteredCertifications.length} of {certifications.length} certifications
				</div>
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each filteredCertifications as cert, index}
						{@const animationDelay = index * 0.1}
						<a
							href={getCertificationUrl(cert)}
							class="group bg-gradient-to-br from-gray-800/50 via-gray-800/30 to-gray-900/50 backdrop-blur-sm rounded-2xl overflow-hidden border border-gray-700 hover:border-green-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-green-500/20 hover:scale-[1.02] relative animate-fadeInUp"
							style="animation-delay: {animationDelay}s"
						>
							<!-- Decorative gradient overlay -->
							<div class="absolute inset-0 bg-gradient-to-br from-green-500/0 via-emerald-500/0 to-teal-500/0 group-hover:from-green-500/5 group-hover:via-emerald-500/5 group-hover:to-teal-500/5 transition-all duration-300 pointer-events-none"></div>
							<div class="relative z-10">
								<div class="p-6">
									<div class="flex items-start justify-between mb-3">
										<div
											class="w-12 h-12 bg-gradient-to-br from-green-600 to-emerald-600 rounded-lg flex items-center justify-center flex-shrink-0"
										>
											<Award class="w-6 h-6 text-white" />
										</div>
										<div class="flex items-center gap-2">
											{#if cert.featured}
												<div
													class="px-2 py-1 bg-yellow-500/90 text-yellow-900 text-xs font-bold rounded flex items-center gap-1"
												>
													<Star class="w-3 h-3 fill-current" />
													Featured
												</div>
											{/if}
											{#if isExpiringSoon(cert.expiryDate)}
												<div
													class="px-2 py-1 bg-orange-500/90 text-orange-900 text-xs font-bold rounded flex items-center gap-1"
												>
													<Clock class="w-3 h-3" />
													Expiring
												</div>
											{/if}
										</div>
									</div>

									<h3
										class="text-xl font-bold text-white mb-2 group-hover:text-green-400 transition-colors line-clamp-2"
									>
										{cert.name}
									</h3>

									{#if cert.issuer}
										<p class="text-gray-300 text-sm mb-4">{cert.issuer}</p>
									{/if}

									<div class="flex items-center gap-2 mb-4 flex-wrap">
										<span
											class="px-2 py-1 rounded-lg text-xs font-medium border {getStatusColor(cert.status)}"
										>
											{cert.status}
										</span>
										<span
											class="px-2 py-1 rounded-lg text-xs font-medium border {getCategoryColor(cert.category)}"
										>
											{cert.category}
										</span>
									</div>

									{#if cert.description}
										<p class="text-gray-300 text-sm mb-4 line-clamp-2">{cert.description}</p>
									{/if}

									<div class="flex items-center gap-4 text-xs text-gray-400 mb-4">
										{#if cert.issueDate}
											<div class="flex items-center gap-1">
												<Calendar class="w-3 h-3" />
												<span>Issued {formatDate(cert.issueDate)}</span>
											</div>
										{/if}
										{#if cert.expiryDate}
											<div class="flex items-center gap-1">
												<Clock class="w-3 h-3" />
												<span>Expires {formatDate(cert.expiryDate)}</span>
											</div>
										{/if}
									</div>

									<div class="flex items-center justify-between pt-4 border-t border-gray-700">
										<div class="flex items-center gap-2 text-green-400 text-sm font-medium group-hover:gap-3 transition-all">
											<span>View Details</span>
											<ExternalLink class="w-4 h-4" />
										</div>
									</div>
								</div>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	</section>
</div>

<style>
	@keyframes fadeInUp {
		from {
			opacity: 0;
			transform: translateY(30px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.animate-fadeInUp) {
		animation: fadeInUp 0.6s ease-out;
		animation-fill-mode: both;
	}
</style>
