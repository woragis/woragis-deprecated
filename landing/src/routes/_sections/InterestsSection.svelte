<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiRedis } from 'svelte-icons-pack/si';
	import { Brain, GitBranch } from 'lucide-svelte';
	import { useFeaturedInterestsQuery } from '$lib/queries/interests';
	import { language } from '$lib/i18n';
	import type { Interest } from '$lib/types/interest';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();

	// Fetch featured interests from API
	const interestsQuery = useFeaturedInterestsQuery($language);
	let interests = $derived(interestsQuery.data || []);
	let loading = $derived(interestsQuery.isPending);

	// Helper function to get icon component based on icon name
	function getIconComponent(iconName?: string) {
		if (!iconName) return null;
		switch (iconName) {
			case 'Brain':
				return Brain;
			case 'SiRedis':
				return SiRedis;
			case 'GitBranch':
				return GitBranch;
			default:
				return null;
		}
	}

	// Helper function to get gradient classes from color
	function getGradientClasses(color?: string) {
		if (!color) return 'from-gray-900/30 to-gray-800/20';
		switch (color) {
			case 'pink-purple':
				return 'from-pink-900/30 to-purple-900/20';
			case 'red-orange':
				return 'from-red-900/30 to-orange-900/20';
			case 'green-emerald':
				return 'from-green-900/30 to-emerald-900/20';
			default:
				return 'from-gray-900/30 to-gray-800/20';
		}
	}

	// Helper function to get border classes
	function getBorderClasses(color?: string) {
		if (!color) return 'border-gray-700/30 hover:border-gray-500/50';
		switch (color) {
			case 'pink-purple':
				return 'border-pink-700/30 hover:border-pink-500/50';
			case 'red-orange':
				return 'border-red-700/30 hover:border-red-500/50';
			case 'green-emerald':
				return 'border-green-700/30 hover:border-green-500/50';
			default:
				return 'border-gray-700/30 hover:border-gray-500/50';
		}
	}

	// Helper function to get shadow classes
	function getShadowClasses(color?: string) {
		if (!color) return 'hover:shadow-gray-500/20';
		switch (color) {
			case 'pink-purple':
				return 'hover:shadow-pink-500/20';
			case 'red-orange':
				return 'hover:shadow-red-500/20';
			case 'green-emerald':
				return 'hover:shadow-green-500/20';
			default:
				return 'hover:shadow-gray-500/20';
		}
	}

	// Helper function to get text gradient classes
	function getTextGradientClasses(color?: string) {
		if (!color) return 'from-gray-400 to-gray-300';
		switch (color) {
			case 'pink-purple':
				return 'from-pink-400 to-purple-400';
			case 'red-orange':
				return 'from-red-400 to-orange-400';
			case 'green-emerald':
				return 'from-green-400 to-emerald-400';
			default:
				return 'from-gray-400 to-gray-300';
		}
	}

	// Helper function to get icon background gradient
	function getIconGradientClasses(color?: string) {
		if (!color) return 'from-gray-600 to-gray-500';
		switch (color) {
			case 'pink-purple':
				return 'from-pink-600 to-purple-600';
			case 'red-orange':
				return 'from-red-600 to-orange-600';
			case 'green-emerald':
				return 'from-green-600 to-emerald-600';
			default:
				return 'from-gray-600 to-gray-500';
		}
	}
</script>

<section id="interests" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<h2 class="text-4xl font-bold mb-12 text-center">{t('interests.title')}</h2>
		{#if loading}
			<div class="flex items-center justify-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500"></div>
			</div>
		{:else if interests.length === 0}
			<div class="text-center py-20">
				<p class="text-gray-400 text-lg">No interests available</p>
			</div>
		{:else}
			<div class="grid md:grid-cols-2 gap-6">
				{#each interests as interest (interest.id)}
					{@const IconComponent = getIconComponent(interest.icon)}
					{@const bgGradient = interest.bgGradient || getGradientClasses(interest.color)}
					{@const borderClasses = interest.borderColor || getBorderClasses(interest.color)}
					{@const shadowClasses = interest.shadowColor || getShadowClasses(interest.color)}
					{@const textGradient = getTextGradientClasses(interest.color)}
					{@const iconGradient = getIconGradientClasses(interest.color)}
					<div
						class="bg-gradient-to-br {bgGradient} rounded-xl p-8 border {borderClasses} transition-all duration-300 hover:shadow-lg {shadowClasses}"
						class:md:col-span-2={interest.fullWidth}
					>
						<div class="flex items-center mb-6">
							<div
								class="w-14 h-14 bg-gradient-to-br {iconGradient} rounded-lg flex items-center justify-center mr-4"
							>
								{#if IconComponent === Brain}
									<Brain class="w-8 h-8" stroke="white" />
								{:else if IconComponent === GitBranch}
									<GitBranch class="w-8 h-8" stroke="white" />
								{:else if IconComponent === SiRedis}
									<Icon src={SiRedis} size="2rem" color="white" />
								{/if}
							</div>
							<h3
								class="text-2xl font-bold bg-gradient-to-r {textGradient} bg-clip-text text-transparent"
							>
								{interest.title}
							</h3>
						</div>
						<p class="text-gray-300 leading-relaxed">{interest.description}</p>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</section>

