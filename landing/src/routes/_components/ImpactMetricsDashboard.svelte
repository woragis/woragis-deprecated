<script lang="ts">
	import { TrendingUp, Users, Rocket, DollarSign, Clock, BarChart3, Target, Zap } from 'lucide-svelte';
	import type { ImpactMetric, MetricType } from '$lib/types/impact-metric';

	interface Props {
		metrics: ImpactMetric[];
		loading: boolean;
	}

	let { metrics = [], loading = false }: Props = $props();

	// Aggregate metrics by type
	let aggregatedMetrics = $derived.by(() => {
		const grouped: Record<MetricType, ImpactMetric[]> = {
			projects_delivered: [],
			users_impacted: [],
			performance_improvement: [],
			cost_savings: [],
			time_saved: []
		};

		metrics.forEach((metric) => {
			if (grouped[metric.type]) {
				grouped[metric.type].push(metric);
			}
		});

		// Calculate totals for each type
		const aggregated: Array<{
			type: MetricType;
			total: number;
			count: number;
			unit: string;
			label: string;
			icon: typeof TrendingUp;
			color: string;
			metrics: ImpactMetric[];
		}> = [];

		Object.entries(grouped).forEach(([type, items]) => {
			if (items.length > 0) {
				const total = items.reduce((sum, m) => sum + m.value, 0);
				const unit = items[0].unit;
				const { label, icon, color } = getMetricTypeInfo(type as MetricType);
				aggregated.push({
					type: type as MetricType,
					total,
					count: items.length,
					unit,
					label,
					icon,
					color,
					metrics: items
				});
			}
		});

		// Sort by total value (descending)
		return aggregated.sort((a, b) => b.total - a.total);
	});

	function getMetricTypeInfo(type: MetricType) {
		const info: Record<
			MetricType,
			{ label: string; icon: typeof TrendingUp; color: string }
		> = {
			projects_delivered: {
				label: 'Projects Delivered',
				icon: Rocket,
				color: 'from-blue-500 to-cyan-500'
			},
			users_impacted: {
				label: 'Users Impacted',
				icon: Users,
				color: 'from-green-500 to-emerald-500'
			},
			performance_improvement: {
				label: 'Performance Improvement',
				icon: Zap,
				color: 'from-purple-500 to-pink-500'
			},
			cost_savings: {
				label: 'Cost Savings',
				icon: DollarSign,
				color: 'from-yellow-500 to-amber-500'
			},
			time_saved: {
				label: 'Time Saved',
				icon: Clock,
				color: 'from-indigo-500 to-purple-500'
			}
		};
		return info[type] || { label: type, icon: TrendingUp, color: 'from-gray-500 to-slate-500' };
	}

	function formatValue(value: number, unit: string): string {
		if (unit === 'currency') {
			return new Intl.NumberFormat('en-US', {
				style: 'currency',
				currency: 'USD',
				minimumFractionDigits: 0,
				maximumFractionDigits: 0
			}).format(value);
		}
		if (unit === 'percentage') {
			return `${value.toFixed(1)}%`;
		}
		if (value >= 1000000) {
			return `${(value / 1000000).toFixed(1)}M ${unit}`;
		}
		if (value >= 1000) {
			return `${(value / 1000).toFixed(1)}k ${unit}`;
		}
		return `${value.toFixed(0)} ${unit}`;
	}

	function getUnitLabel(unit: string): string {
		const labels: Record<string, string> = {
			count: '',
			percentage: '',
			currency: '',
			hours: 'hours',
			days: 'days',
			months: 'months',
			years: 'years',
			milliseconds: 'ms',
			seconds: 'seconds',
			minutes: 'minutes'
		};
		return labels[unit] || unit;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if metrics.length === 0}
		<div class="text-center py-20">
			<BarChart3 class="w-16 h-16 mx-auto mb-4 text-gray-600" />
			<p class="text-gray-400 text-lg mb-2">No impact metrics available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<!-- Summary Cards -->
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
			{#each aggregatedMetrics.slice(0, 6) as metric}
				{@const MetricIcon = metric.icon}
				<div
					class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700 hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/20"
				>
					<div class="flex items-start justify-between mb-4">
						<div class="flex items-center gap-3">
							<div
								class="w-12 h-12 bg-gradient-to-br {metric.color} rounded-lg flex items-center justify-center"
							>
								<MetricIcon class="w-6 h-6 text-white" />
							</div>
							<div>
								<h3 class="text-lg font-bold text-white">{metric.label}</h3>
								<p class="text-xs text-gray-400">{metric.count} {metric.count === 1 ? 'metric' : 'metrics'}</p>
							</div>
						</div>
					</div>

					<div class="mb-4">
						<div class="text-3xl font-bold text-white mb-1">
							{formatValue(metric.total, metric.unit)}
						</div>
						<p class="text-sm text-gray-400">Total {getUnitLabel(metric.unit)}</p>
					</div>

					{#if metric.metrics.length > 1}
						<div class="pt-4 border-t border-gray-700">
							<div class="flex items-center justify-between text-xs text-gray-400 mb-2">
								<span>Average</span>
								<span class="text-gray-300 font-medium">
									{formatValue(metric.total / metric.count, metric.unit)}
								</span>
							</div>
							<div class="flex items-center justify-between text-xs text-gray-400">
								<span>Range</span>
								<span class="text-gray-300 font-medium">
									{formatValue(Math.min(...metric.metrics.map((m) => m.value)), metric.unit)} - {formatValue(Math.max(...metric.metrics.map((m) => m.value)), metric.unit)}
								</span>
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<!-- Detailed Metrics by Type -->
		{#if aggregatedMetrics.length > 0}
			<div class="space-y-6">
				{#each aggregatedMetrics as metric}
					{@const MetricIcon = metric.icon}
					<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-6 border border-gray-700">
						<!-- Type Header -->
						<div class="flex items-center gap-3 mb-6">
							<div
								class="w-12 h-12 bg-gradient-to-br {metric.color} rounded-lg flex items-center justify-center"
							>
								<MetricIcon class="w-6 h-6 text-white" />
							</div>
							<div>
								<h3 class="text-2xl font-bold text-white">{metric.label}</h3>
								<p class="text-sm text-gray-400">
									{metric.count} {metric.count === 1 ? 'metric' : 'metrics'} • Total: {formatValue(metric.total, metric.unit)}
								</p>
							</div>
						</div>

						<!-- Individual Metrics -->
						<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
							{#each metric.metrics.slice(0, 6) as item}
								<div
									class="bg-gray-800/50 rounded-lg p-4 border border-gray-700 hover:border-blue-500/50 transition-all duration-300"
								>
									<div class="flex items-start justify-between mb-3">
										<div class="flex-1">
											<div class="text-2xl font-bold text-white mb-1">
												{formatValue(item.value, item.unit)}
											</div>
											{#if item.description}
												<p class="text-sm text-gray-300 line-clamp-2">{item.description}</p>
											{/if}
										</div>
									</div>

									{#if item.periodStart || item.periodEnd}
										<div class="flex items-center gap-2 text-xs text-gray-400 mt-3 pt-3 border-t border-gray-700">
											<Target class="w-3 h-3" />
											<span>
												{#if item.periodStart && item.periodEnd}
													{new Date(item.periodStart).toLocaleDateString('en-US', {
														month: 'short',
														year: 'numeric'
													})} - {new Date(item.periodEnd).toLocaleDateString('en-US', {
														month: 'short',
														year: 'numeric'
													})}
												{:else if item.periodStart}
													Since {new Date(item.periodStart).toLocaleDateString('en-US', {
														month: 'short',
														year: 'numeric'
													})}
												{:else if item.periodEnd}
													Until {new Date(item.periodEnd).toLocaleDateString('en-US', {
														month: 'short',
														year: 'numeric'
													})}
												{/if}
											</span>
										</div>
									{/if}
								</div>
							{/each}
						</div>

						{#if metric.metrics.length > 6}
							<div class="mt-4 text-center">
								<p class="text-sm text-gray-400">
									And {metric.metrics.length - 6} more {metric.metrics.length - 6 === 1 ? 'metric' : 'metrics'}...
								</p>
							</div>
						{/if}
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>

