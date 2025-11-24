<script lang="ts">
	import { page } from '$app/stores';

	// TODO: Implement Certification detail page
	// - Fetch certification by slug from backend API
	// - Display full certification details
	// - Show issuer information
	// - Display issue and expiry dates
	// - Show verification URL
	// - Display linked skills
	// - Show status (active, expired, expiring soon)
	// - Add related certifications section
	// - Add breadcrumb navigation
	// - Handle loading and error states

	const slug = $derived($page.params.slug);
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<div class="container mx-auto px-6 py-20">
		<div class="max-w-4xl mx-auto">
			<!-- Header -->
			<div class="text-center mb-12">
				<h1 class="text-5xl font-bold mb-4 bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent">
					Certification Detail
				</h1>
				<p class="text-gray-400 text-lg">
					Slug: {slug}
				</p>
			</div>

			<!-- Coming Soon Message -->
			<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-12 border border-gray-700 text-center">
				<div class="max-w-2xl mx-auto">
					<h2 class="text-3xl font-bold text-white mb-4">Coming Soon</h2>
					<p class="text-gray-300 text-lg mb-6">
						This page will display the full details of a certification, including issuer, dates, verification link, and linked skills.
					</p>
					<p class="text-gray-400">
						Features to be implemented:
					</p>
					<ul class="text-left text-gray-300 mt-4 space-y-2 max-w-md mx-auto">
						<li>• Full certification details</li>
						<li>• Issuer information</li>
						<li>• Issue and expiry dates</li>
						<li>• Verification URL</li>
						<li>• Linked skills</li>
						<li>• Status indicators</li>
						<li>• Related certifications section</li>
						<li>• Breadcrumb navigation</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</div>

