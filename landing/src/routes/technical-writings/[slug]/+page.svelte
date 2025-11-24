<script lang="ts">
	import { page } from '$app/stores';

	// TODO: Implement Technical Writing detail page
	// - Fetch technical writing by slug from backend API
	// - Display full content (if available) or link to external URL
	// - Show all metadata: type, platform, published date, reading time
	// - Display engagement metrics (views, likes, shares, comments)
	// - Show topics and technologies tags
	// - Add related writings section
	// - Add breadcrumb navigation
	// - Add social sharing buttons
	// - Handle loading and error states

	const slug = $derived($page.params.slug);
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<div class="container mx-auto px-6 py-20">
		<div class="max-w-4xl mx-auto">
			<!-- Header -->
			<div class="text-center mb-12">
				<h1 class="text-5xl font-bold mb-4 bg-gradient-to-r from-orange-400 to-red-400 bg-clip-text text-transparent">
					Technical Writing Detail
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
						This page will display the full details of a technical writing, including content, metadata, engagement metrics, and related writings.
					</p>
					<p class="text-gray-400">
						Features to be implemented:
					</p>
					<ul class="text-left text-gray-300 mt-4 space-y-2 max-w-md mx-auto">
						<li>• Full content display or external link</li>
						<li>• Metadata (type, platform, date, reading time)</li>
						<li>• Engagement metrics</li>
						<li>• Topics and technologies tags</li>
						<li>• Related writings section</li>
						<li>• Breadcrumb navigation</li>
						<li>• Social sharing buttons</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</div>

