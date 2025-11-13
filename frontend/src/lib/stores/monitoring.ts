import { browser } from '$app/environment';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { writable } from 'svelte/store';

export type MetricSample = {
	name: string;
	value: number;
	labels: Record<string, string>;
};

export type MetricSeriesPoint = {
	timestamp: number;
	value: number;
};

export type MetricSeries = {
	key: string;
	name: string;
	labels: Record<string, string>;
	points: MetricSeriesPoint[];
};

export type MonitoringState = {
	status: 'disconnected' | 'connecting' | 'connected' | 'polling';
	error: string | null;
	lastUpdated: number | null;
	samples: MetricSample[];
	series: MetricSeries[];
};

const INITIAL_STATE: MonitoringState = {
	status: 'disconnected',
	error: null,
	lastUpdated: null,
	samples: [],
	series: []
};

const METRICS_WS_PATH = '/metrics/stream';
const POLL_INTERVAL_MS = 5000;

const toWebSocketURL = (baseHttpUrl: string, path: string) => {
	try {
		const url = new URL(baseHttpUrl);
		url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
		url.pathname = path.startsWith('/') ? path : `/${path}`;
		url.search = '';
		url.hash = '';
		return url.toString();
	} catch {
		return path;
	}
};

const parsePrometheusText = (text: string): MetricSample[] => {
	const samples: MetricSample[] = [];
	const lines = text.split('\n');
	const metricRegex =
		/^([^{}\s]+)(?:\{([^}]*)\})?\s+(-?(?:\d*\.)?\d+(?:[eE][+-]?\d+)?)(?:\s+\d+)?$/;

	for (const line of lines) {
		const trimmed = line.trim();
		if (!trimmed || trimmed.startsWith('#')) {
			continue;
		}

		const match = trimmed.match(metricRegex);
		if (!match) {
			continue;
		}

		const [, name, rawLabels, rawValue] = match;
		const labels: Record<string, string> = {};

		if (rawLabels) {
			for (const label of rawLabels.split(',')) {
				const [key, value] = label.split('=');
				if (!key || value === undefined) continue;
				labels[key.trim()] = value.trim().replace(/^"|"$/g, '');
			}
		}

		const value = Number(rawValue);
		if (Number.isNaN(value)) {
			continue;
		}

		samples.push({ name, value, labels });
	}

	return samples;
};

const HISTORY_LIMIT = 200;

const createMonitoringStore = () => {
	let ws: WebSocket | null = null;
	let pollTimer: ReturnType<typeof setInterval> | null = null;
	let manualDisconnect = false;
	const seriesMap = new Map<string, MetricSeries>();
	const { subscribe, update, set } = writable<MonitoringState>(INITIAL_STATE);

	const applySamples = (samples: MetricSample[]) => {
		const timestamp = Date.now();

		for (const sample of samples) {
			const key = `${sample.name}|${JSON.stringify(sample.labels)}`;
			const existing = seriesMap.get(key) ?? {
				key,
				name: sample.name,
				labels: sample.labels,
				points: []
			};

			const nextPoints =
				existing.points.length >= HISTORY_LIMIT
					? existing.points.slice(existing.points.length - HISTORY_LIMIT + 1)
					: [...existing.points];

			nextPoints.push({ timestamp, value: sample.value });

			seriesMap.set(key, {
				key,
				name: sample.name,
				labels: sample.labels,
				points: nextPoints
			});
		}

		return Array.from(seriesMap.values()).map((series) => ({
			...series,
			points: [...series.points]
		}));
	};

	const stopPolling = () => {
		if (pollTimer) {
			clearInterval(pollTimer);
			pollTimer = null;
		}
	};

	const pollOnce = async () => {
		if (!browser) return;

		try {
			const response = await fetch(`${PUBLIC_API_BASE_URL ?? 'http://localhost:8080'}/metrics`, {
				headers: { Accept: 'text/plain' }
			});
			if (!response.ok) {
				throw new Error(`HTTP ${response.status}`);
			}
			const text = await response.text();
			const samples = parsePrometheusText(text);
			const series = applySamples(samples);
			set({
				status: 'polling',
				error: null,
				lastUpdated: Date.now(),
				samples,
				series
			});
		} catch (error) {
			update((state) => ({
				...state,
				error: error instanceof Error ? error.message : 'Unable to fetch metrics'
			}));
		}
	};

	const startPolling = () => {
		stopPolling();
		pollOnce();
		pollTimer = setInterval(pollOnce, POLL_INTERVAL_MS);
	};

	const disconnect = () => {
		manualDisconnect = true;
		stopPolling();
		if (ws) {
			ws.close();
			ws = null;
		}
		seriesMap.clear();
		set({ ...INITIAL_STATE, status: 'disconnected' });
		manualDisconnect = false;
	};

	const connect = () => {
		if (!browser) {
			return;
		}
		if (ws || pollTimer) {
			return;
		}

		const wsUrl = toWebSocketURL(PUBLIC_API_BASE_URL ?? 'http://localhost:8080', METRICS_WS_PATH);

		try {
			ws = new WebSocket(wsUrl);
			set({
				status: 'connecting',
				error: null,
				lastUpdated: null,
				samples: [],
				series: []
			});

			ws.onopen = () => {
				stopPolling();
				update((state) => ({ ...state, status: 'connected', error: null }));
			};

			ws.onmessage = (event) => {
				let payload = '';

				if (typeof event.data === 'string') {
					payload = event.data;
				} else if (event.data instanceof Blob) {
					event.data.text().then((text) => {
						const samples = parsePrometheusText(text);
						const series = applySamples(samples);
						set({
							status: 'connected',
							error: null,
							lastUpdated: Date.now(),
							samples,
							series
						});
					});
					return;
				}

				const samples = parsePrometheusText(payload);
				const series = applySamples(samples);
				set({
					status: 'connected',
					error: null,
					lastUpdated: Date.now(),
					samples,
					series
				});
			};

			ws.onerror = () => {
				update((state) => ({
					...state,
					error: 'WebSocket connection failed, using polling fallback.'
				}));
				ws?.close();
			};

			ws.onclose = () => {
				ws = null;
				// start polling fallback if we were not explicitly disconnected
				if (!manualDisconnect && !pollTimer) {
					startPolling();
				}
			};
		} catch (error) {
			update((state) => ({
				...state,
				error: error instanceof Error ? error.message : 'Unable to open metrics websocket'
			}));
			startPolling();
		}
	};

	return {
		subscribe,
		connect,
		disconnect,
		refresh: pollOnce
	};
};

export const monitoringStore = createMonitoringStore();
