<script lang="ts">
	import { Layers, ChevronDown, ChevronUp, Server, Database, Cpu, Network, Zap, Shield } from 'lucide-svelte';
	import type { SystemDesign, Component } from '$lib/types/system-design';

	interface Props {
		designs: SystemDesign[];
		loading: boolean;
	}

	let { designs = [], loading = false }: Props = $props();

	let expandedDesign = $state<string | null>(null);

	function toggleDesign(id: string) {
		expandedDesign = expandedDesign === id ? null : id;
	}

	function getComponentIcon(technology: string) {
		const tech = technology.toLowerCase();
		if (tech.includes('database') || tech.includes('db') || tech.includes('sql') || tech.includes('redis')) {
			return Database;
		}
		if (tech.includes('server') || tech.includes('api') || tech.includes('backend')) {
			return Server;
		}
		if (tech.includes('network') || tech.includes('cdn') || tech.includes('load')) {
			return Network;
		}
		if (tech.includes('cache') || tech.includes('queue') || tech.includes('message')) {
			return Zap;
		}
		if (tech.includes('security') || tech.includes('auth') || tech.includes('firewall')) {
			return Shield;
		}
		return Cpu;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if designs.length === 0}
		<div class="text-center py-20">
			<Layers class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No system designs available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div class="grid md:grid-cols-2 gap-6">
			{#each designs as design}
				<div
					class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300"
				>
					<!-- Header -->
					<div class="flex items-start justify-between mb-4">
						<div class="flex items-center gap-3 flex-1">
							<div
								class="w-12 h-12 bg-gradient-to-br from-blue-600 to-cyan-600 rounded-lg flex items-center justify-center flex-shrink-0"
							>
								<Layers class="w-6 h-6 text-white" />
							</div>
							<div class="flex-1 min-w-0">
								<h3 class="text-xl font-bold text-white mb-1">{design.title}</h3>
								<p class="text-sm text-gray-400">System Design</p>
							</div>
						</div>
						<button
							onclick={() => toggleDesign(design.id)}
							class="text-gray-400 hover:text-white transition-colors flex-shrink-0 ml-2"
							aria-label={expandedDesign === design.id ? 'Collapse' : 'Expand'}
						>
							{#if expandedDesign === design.id}
								<ChevronUp class="w-5 h-5" />
							{:else}
								<ChevronDown class="w-5 h-5" />
							{/if}
						</button>
					</div>

					<!-- Description -->
					<p class="text-gray-300 mb-4 line-clamp-3">{design.description}</p>

					<!-- Components Preview -->
					{#if design.components && design.components.length > 0}
						<div class="mb-4">
							<div class="flex flex-wrap gap-2">
								{#each design.components.slice(0, 3) as component}
									<span
										class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
									>
										{component.name}
									</span>
								{/each}
								{#if design.components.length > 3}
									<span
										class="px-2 py-1 text-xs rounded bg-gray-700/50 text-gray-300 border border-gray-600"
									>
										+{design.components.length - 3} more
									</span>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Diagram Preview -->
					{#if design.diagram}
						<div class="mb-4">
							<a
								href={design.diagram}
								target="_blank"
								rel="noopener noreferrer"
								class="text-sm text-blue-400 hover:text-blue-300 transition-colors flex items-center gap-1"
							>
								View Diagram
								<ChevronDown class="w-4 h-4 rotate-[-45deg]" />
							</a>
						</div>
					{/if}

					<!-- Expanded Details -->
					{#if expandedDesign === design.id}
						<div class="mt-4 space-y-4 pt-4 border-t border-gray-700">
							<!-- Components -->
							{#if design.components && design.components.length > 0}
								<div>
									<h4 class="text-sm font-semibold text-blue-400 mb-3">Components</h4>
									<div class="space-y-2">
										{#each design.components as component}
											{@const ComponentIcon = getComponentIcon(component.technology)}
											<div class="bg-gray-800/50 rounded-lg p-3 border border-gray-700">
												<div class="flex items-start gap-3">
													<div
														class="w-8 h-8 bg-gradient-to-br from-blue-600/20 to-cyan-600/20 rounded-lg flex items-center justify-center flex-shrink-0 border border-blue-500/30"
													>
														<ComponentIcon class="w-4 h-4 text-blue-400" />
													</div>
													<div class="flex-1 min-w-0">
														<div class="flex items-center justify-between mb-1">
															<span class="font-medium text-white text-sm">{component.name}</span>
															<span class="text-xs text-cyan-400 bg-cyan-600/20 px-2 py-0.5 rounded border border-cyan-500/30">
																{component.technology}
															</span>
														</div>
														<p class="text-xs text-gray-400">{component.description}</p>
													</div>
												</div>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Data Flow -->
							{#if design.dataFlow}
								<div>
									<h4 class="text-sm font-semibold text-blue-400 mb-2">Data Flow</h4>
									<p class="text-sm text-gray-300 leading-relaxed">{design.dataFlow}</p>
								</div>
							{/if}

							<!-- Scalability -->
							{#if design.scalability}
								<div>
									<h4 class="text-sm font-semibold text-green-400 mb-2 flex items-center gap-2">
										<Zap class="w-4 h-4" />
										Scalability
									</h4>
									<p class="text-sm text-gray-300 leading-relaxed">{design.scalability}</p>
								</div>
							{/if}

							<!-- Reliability -->
							{#if design.reliability}
								<div>
									<h4 class="text-sm font-semibold text-purple-400 mb-2 flex items-center gap-2">
										<Shield class="w-4 h-4" />
										Reliability
									</h4>
									<p class="text-sm text-gray-300 leading-relaxed">{design.reliability}</p>
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

