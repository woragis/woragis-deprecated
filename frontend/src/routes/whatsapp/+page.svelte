<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchQRCode, fetchStatus, type WhatsAppQRResponse, type WhatsAppStatusResponse } from '$lib/api/whatsapp';

	let qrCode: string | null = null;
	let qrText: string | null = null;
	let connected = false;
	let loading = true;
	let error: string | null = null;
	let message = 'Loading...';
	let pollInterval: ReturnType<typeof setInterval> | null = null;

	const loadQRCode = async () => {
		try {
			loading = true;
			error = null;
			const response = await fetchQRCode();
			connected = response.connected;
			qrCode = response.qr_code;
			qrText = response.qr_text || null;
			message = response.message;

			// Stop polling if connected
			if (connected && pollInterval) {
				clearInterval(pollInterval);
				pollInterval = null;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load QR code';
			message = error;
			console.error('Failed to fetch QR code:', err);
		} finally {
			loading = false;
		}
	};

	const startPolling = () => {
		// Poll every 2 seconds if not connected
		if (!connected && !pollInterval) {
			pollInterval = setInterval(() => {
				if (!connected) {
					loadQRCode();
				}
			}, 2000);
		}
	};

	const stopPolling = () => {
		if (pollInterval) {
			clearInterval(pollInterval);
			pollInterval = null;
		}
	};

	onMount(() => {
		loadQRCode();
		startPolling();

		return () => {
			stopPolling();
		};
	});
</script>

<div class="container mx-auto px-4 py-8 max-w-2xl">
	<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
		<h1 class="text-3xl font-bold mb-6 text-gray-900 dark:text-white">WhatsApp Connection</h1>

		{#if loading && !qrCode}
			<div class="flex flex-col items-center justify-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
				<p class="text-gray-600 dark:text-gray-400">{message}</p>
			</div>
		{:else if error && !qrCode}
			<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-4">
				<p class="text-red-800 dark:text-red-200">{error}</p>
			</div>
		{:else if connected}
			<div class="flex flex-col items-center justify-center py-12">
				<div class="bg-green-100 dark:bg-green-900/20 rounded-full p-4 mb-4">
					<svg
						class="w-16 h-16 text-green-600 dark:text-green-400"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M5 13l4 4L19 7"
						></path>
					</svg>
				</div>
				<h2 class="text-2xl font-semibold text-gray-900 dark:text-white mb-2">Connected!</h2>
				<p class="text-gray-600 dark:text-gray-400 text-center">
					WhatsApp is successfully connected. You can now receive notifications via WhatsApp.
				</p>
			</div>
		{:else if qrCode}
			<div class="flex flex-col items-center">
				<div class="mb-6">
					<p class="text-gray-700 dark:text-gray-300 text-center mb-4">{message}</p>
					<div class="bg-white p-4 rounded-lg border-2 border-gray-200 dark:border-gray-700 inline-block">
						<img src={qrCode} alt="WhatsApp QR Code" class="w-64 h-64" />
					</div>
				</div>

				<div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 w-full max-w-md">
					<h3 class="font-semibold text-blue-900 dark:text-blue-200 mb-2">How to connect:</h3>
					<ol class="list-decimal list-inside space-y-2 text-sm text-blue-800 dark:text-blue-300">
						<li>Open WhatsApp on your phone</li>
						<li>Go to Settings → Linked Devices</li>
						<li>Tap "Link a Device"</li>
						<li>Scan this QR code</li>
					</ol>
				</div>

				{#if loading}
					<div class="mt-4 flex items-center gap-2 text-gray-600 dark:text-gray-400">
						<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
						<span class="text-sm">Refreshing...</span>
					</div>
				{/if}
			</div>
		{:else}
			<div class="flex flex-col items-center justify-center py-12">
				<p class="text-gray-600 dark:text-gray-400 text-center">{message}</p>
			</div>
		{/if}

		<div class="mt-6 flex justify-center">
			<button
				on:click={loadQRCode}
				disabled={loading}
				class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-lg transition-colors"
			>
				{loading ? 'Loading...' : 'Refresh'}
			</button>
		</div>
	</div>
</div>

<style>
	.container {
		min-height: calc(100vh - 4rem);
	}
</style>

