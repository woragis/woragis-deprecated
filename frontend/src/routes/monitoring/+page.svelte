<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { monitoringStore, type MetricSeries } from '$lib';
	import { PUBLIC_GRAFANA_PANELS } from '$env/static/public';

	let status = 'disconnected';
	let error: string | null = null;
	let lastUpdated: number | null = null;
	let series: MetricSeries[] = [];

	const unsubscribe = monitoringStore.subscribe((state) => {
		status = state.status;
		error = state.error;
		lastUpdated = state.lastUpdated;
		series = state.series;
	});

	onMount(() => {
		monitoringStore.connect();
	});

	onDestroy(() => {
		unsubscribe();
		monitoringStore.disconnect();
	});

	function formatTimestamp(value: number | null) {
		if (!value) return '—';
		return new Date(value).toLocaleTimeString();
	}

	const highlightMetrics = new Set([
		'woragis_http_requests_total',
		'woragis_http_request_duration_seconds_bucket',
		'woragis_user_registrations_total'
	]);

	const formatLabels = (labels: Record<string, string>) => {
		const entries = Object.entries(labels);
		if (entries.length === 0) return '—';
		return entries.map(([key, value]) => `${key}=${value}`).join(', ');
	};

	const toPath = (points: MetricSeries['points']) => {
		if (points.length === 0) return '';
		const values = points.map((point) => point.value);
		const min = Math.min(...values);
		const max = Math.max(...values);
		const range = max - min || 1;
		const lastIndex = Math.max(points.length - 1, 1);

		return points
			.map((point, index) => {
				const x = (index / lastIndex) * 100;
				const y = 40 - ((point.value - min) / range) * 40;
				return `${index === 0 ? 'M' : 'L'}${x},${y}`;
			})
			.join(' ');
	};

	const latestValue = (metric: MetricSeries) =>
		metric.points.length ? metric.points[metric.points.length - 1].value : 0;

	const sortedSeries = () => [...series].sort((a, b) => a.name.localeCompare(b.name));

	type GrafanaPanel = {
		key: string;
		title: string;
		url: string;
	};

	const grafanaPanels: GrafanaPanel[] = (() => {
		const raw = (PUBLIC_GRAFANA_PANELS ?? '').split(',').map((value) => value.trim());
		return raw
			.filter((value) => value.length > 0)
			.map((entry, index) => {
				const [title, url] = entry.split('|');
				const resolvedUrl = (url ?? title ?? '').trim();
				return {
					key: `${index}-${resolvedUrl}`,
					title: title?.trim() || `Grafana panel #${index + 1}`,
					url: resolvedUrl
				};
			})
			.filter((panel) => panel.url.length > 0);
	})();

	$: console.log('Monitoring metrics series', series);
</script>

<section class="space-y-6">
	<header class="rounded border border-slate-800 bg-slate-900/60 p-4">
		<h2 class="text-lg font-semibold text-slate-100">Runtime Metrics Stream</h2>
		<div class="mt-3 flex flex-wrap items-center gap-4 text-xs">
			<span
				class={`inline-flex items-center gap-2 rounded px-3 py-1 ${
					status === 'connected'
						? 'bg-emerald-500/20 text-emerald-200'
						: status === 'polling'
							? 'bg-amber-500/20 text-amber-200'
							: status === 'connecting'
								? 'bg-sky-500/20 text-sky-200'
								: 'bg-slate-800 text-slate-300'
				}`}
			>
				<span class="h-2 w-2 rounded-full bg-current"></span>
				{status}
			</span>
			<button
				class="rounded bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-100 hover:bg-slate-700"
				type="button"
				on:click={() => monitoringStore.refresh()}
			>
				Refresh now
			</button>
			<span class="text-slate-400">Last update: {formatTimestamp(lastUpdated)}</span>
		</div>
		{#if error}
			<p
				class="mt-2 rounded border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-200"
			>
				{error}
			</p>
		{/if}
		<p class="mt-3 text-xs leading-relaxed text-slate-300">
			The monitoring service streams Prometheus metrics in real time. When the websocket is not
			available, the client falls back to polling <code>/metrics</code> every few seconds. Pin this panel
			while observing traffic or running automated tests.
		</p>
	</header>

	{#if grafanaPanels.length > 0}
		<section class="space-y-3 rounded border border-slate-800 bg-slate-900/60 p-4">
			<header class="flex items-center justify-between">
				<h3 class="text-sm font-semibold text-slate-100">Grafana Dashboards</h3>
				<span class="text-xs text-slate-400">
					{grafanaPanels.length} panel{grafanaPanels.length === 1 ? '' : 's'}
				</span>
			</header>
			<div class="grid gap-4 lg:grid-cols-2">
				{#each grafanaPanels as panel (panel.key)}
					<div class="space-y-2">
						<h4 class="text-xs font-semibold tracking-wide text-slate-400 uppercase">
							{panel.title}
						</h4>
						<iframe
							src={panel.url}
							class="h-[320px] w-full overflow-hidden rounded border border-slate-800 bg-slate-950"
							frameborder="0"
							loading="lazy"
							allow="fullscreen"
						/>
					</div>
				{/each}
			</div>
		</section>
	{:else}
		<section
			class="rounded border border-dashed border-slate-700 bg-slate-900/40 p-4 text-xs text-slate-300"
		>
			<p class="font-semibold text-slate-100">Grafana panels</p>
			<p class="mt-2">
				Set <code>PUBLIC_GRAFANA_PANELS</code> in <code>frontend/.env</code> with one or more
				Grafana embed URLs in the form
				<code>Title|http://localhost:3000/d-solo/uid/dashboard?orgId=1&panelId=2&refresh=5s</code>.
				Panels will appear here automatically.
			</p>
		</section>
	{/if}

	<section class="space-y-4">
		<h3 class="text-sm font-semibold text-slate-100">
			Live Series ({sortedSeries().length})
		</h3>
		{#if sortedSeries().length === 0}
			<p class="rounded border border-slate-800 bg-slate-900/60 p-4 text-sm text-slate-400">
				Waiting for metric samples. Exercise the API or use the refresh button to load a snapshot.
			</p>
		{:else}
			<div class="grid gap-4 lg:grid-cols-2">
				{#each sortedSeries() as metric (metric.key)}
					<div
						class={`space-y-3 rounded border p-4 ${
							highlightMetrics.has(metric.name)
								? 'border-emerald-500/40 bg-emerald-500/10'
								: 'border-slate-800 bg-slate-900/60'
						}`}
					>
						<div class="flex items-center justify-between">
							<div>
								<h4 class="text-sm font-semibold text-slate-100">{metric.name}</h4>
								<p class="text-[11px] text-slate-400">{formatLabels(metric.labels)}</p>
							</div>
							<div class="text-right text-xs text-slate-300">
								<div class="text-slate-100">{latestValue(metric).toLocaleString()}</div>
								<div>Latest value</div>
							</div>
						</div>
						<div class="relative h-32 overflow-hidden rounded bg-slate-950/60">
							{#if metric.points.length > 1}
								<svg class="h-full w-full" viewBox="0 0 100 40" preserveAspectRatio="none">
									<path
										class="text-primary stroke-current"
										d={toPath(metric.points)}
										fill="none"
										stroke-width="1.5"
									/>
								</svg>
							{:else}
								<div class="flex h-full items-center justify-center text-xs text-slate-500">
									Insufficient data for sparkline
								</div>
							{/if}
						</div>
						<div class="flex items-center justify-between text-[11px] text-slate-400">
							<span>Points: {metric.points.length}</span>
							<span>
								Range:
								{metric.points.length
									? `${Math.min(...metric.points.map((p) => p.value)).toLocaleString()} → ${Math.max(...metric.points.map((p) => p.value)).toLocaleString()}`
									: '—'}
							</span>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</section>
