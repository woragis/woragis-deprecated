<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiRedis } from 'svelte-icons-pack/si';
	import { Brain, GitBranch } from 'lucide-svelte';
	import { interests } from '$lib/constants';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section id="interests" class="container mx-auto px-6 py-20">
	<div class="max-w-6xl mx-auto">
		<h2 class="text-4xl font-bold mb-12 text-center">{t('interests.title')}</h2>
		<div class="grid md:grid-cols-2 gap-6">
			<!-- AI & RAG -->
			<div
				class="bg-gradient-to-br from-pink-900/30 to-purple-900/20 rounded-xl p-8 border border-pink-700/30 hover:border-pink-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-pink-500/20"
			>
				<div class="flex items-center mb-6">
					<div
						class="w-14 h-14 bg-gradient-to-br from-pink-600 to-purple-600 rounded-lg flex items-center justify-center mr-4"
					>
						<Brain class="w-8 h-8" stroke="white" />
					</div>
					<h3
						class="text-2xl font-bold bg-gradient-to-r from-pink-400 to-purple-400 bg-clip-text text-transparent"
					>
						{interests[0].title}
					</h3>
				</div>
				<p class="text-gray-300 leading-relaxed">{interests[0].description}</p>
			</div>

			<!-- Redis & Pub/Sub -->
			<div
				class="bg-gradient-to-br from-red-900/30 to-orange-900/20 rounded-xl p-8 border border-red-700/30 hover:border-red-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-red-500/20"
			>
				<div class="flex items-center mb-6">
					<div
						class="w-14 h-14 bg-gradient-to-br from-red-600 to-orange-600 rounded-lg flex items-center justify-center mr-4"
					>
						<Icon src={SiRedis} size="2rem" color="white" />
					</div>
					<h3
						class="text-2xl font-bold bg-gradient-to-r from-red-400 to-orange-400 bg-clip-text text-transparent"
					>
						{interests[1].title}
					</h3>
				</div>
				<p class="text-gray-300 leading-relaxed">{interests[1].description}</p>
			</div>

			<!-- Distributed Architecture -->
			<div
				class="bg-gradient-to-br from-green-900/30 to-emerald-900/20 rounded-xl p-8 border border-green-700/30 hover:border-green-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-green-500/20 md:col-span-2"
			>
				<div class="flex items-center mb-6">
					<div
						class="w-14 h-14 bg-gradient-to-br from-green-600 to-emerald-600 rounded-lg flex items-center justify-center mr-4"
					>
						<GitBranch class="w-8 h-8" stroke="white" />
					</div>
					<h3
						class="text-2xl font-bold bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent"
					>
						{interests[2].title}
					</h3>
				</div>
				<p class="text-gray-300 leading-relaxed">{interests[2].description}</p>
			</div>
		</div>
	</div>
</section>

