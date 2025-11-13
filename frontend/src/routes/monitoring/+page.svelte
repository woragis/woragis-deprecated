<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { monitoringStore, type MetricSample } from '$lib';

	let status = 'disconnected';
	let error: string | null = null;
	let lastUpdated: number | null = null;
	let samples: MetricSample[] = [];

	const unsubscribe = monitoringStore.subscribe((state) => {
		status = state.status;
		error = state.error;
		lastUpdated = state.lastUpdated;
		samples = state.samples;
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

	$: sortedSamples = [...samples].sort((a, b) => a.name.localeCompare(b.name));
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

	<section class="space-y-4">
		<h3 class="text-sm font-semibold text-slate-100">Live Samples ({sortedSamples.length})</h3>
		{#if sortedSamples.length === 0}
			<p class="rounded border border-slate-800 bg-slate-900/60 p-4 text-sm text-slate-400">
				Waiting for metric samples. Exercise the API or use the refresh button to load a snapshot.
			</p>
		{:else}
			<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
				<table class="min-w-full border-separate border-spacing-y-2 text-xs text-slate-200">
					<thead class="text-[11px] tracking-wide text-slate-500 uppercase">
						<tr>
							<th class="text-left">Metric</th>
							<th class="text-left">Value</th>
							<th class="text-left">Labels</th>
						</tr>
					</thead>
					<tbody>
						{#each sortedSamples as sample (sample.name + JSON.stringify(sample.labels))}
							<tr
								class={`rounded border ${
									highlightMetrics.has(sample.name)
										? 'border-emerald-500/40 bg-emerald-500/10'
										: 'border-slate-800 bg-slate-950/60'
								}`}
							>
								<td class="px-3 py-2 font-semibold">{sample.name}</td>
								<td class="px-3 py-2 text-slate-100">{sample.value}</td>
								<td class="px-3 py-2 text-[11px] text-slate-400">
									{#if Object.keys(sample.labels).length === 0}
										<span>—</span>
									{:else}
										<ul class="flex flex-wrap gap-2">
											{#each Object.entries(sample.labels) as [labelKey, labelValue] (`${sample.name}-${labelKey}-${labelValue}`)}
												<li class="rounded bg-slate-800/60 px-2 py-1">
													<strong>{labelKey}</strong>=<span>{labelValue}</span>
												</li>
											{/each}
										</ul>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</section>
</section>
