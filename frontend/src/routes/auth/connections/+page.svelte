<script lang="ts">
import { browser } from '$app/environment';
import { onDestroy, onMount } from 'svelte';

import { API_BASE_URL } from '@clients/apiClient';
import { listOAuthAccounts, listOAuthProviders, startOAuth, unlinkOAuthAccount } from '$lib/api/auth';
import { authStore } from '$lib';
import { getDeviceFingerprint, getDeviceName, getUserAgent } from '$lib/device';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

	interface OAuthProviderMeta {
		id: string;
		name: string;
	}

	interface OAuthAccount {
		provider: string;
		linked_at: string;
		updated_at: string;
		scopes: string[];
	}

const API_ORIGIN = browser ? new URL(API_BASE_URL).origin : '';

let providers: OAuthProviderMeta[] = [];
let accounts: OAuthAccount[] = [];
let loadingProviders = true;
let loadingAccounts = false;
let error = '';
let info = '';
let isAuthenticated = false;
	let oauthWindow: Window | null = null;
	let oauthLoading: Record<string, boolean> = {};

	const findAccount = (provider: string) => accounts.find((account) => account.provider === provider);

const refreshAccounts = async (showSpinner = true) => {
	if (!isAuthenticated) {
		accounts = [];
		return;
	}

	if (showSpinner) {
		loadingAccounts = true;
	}
	try {
		const response = await listOAuthAccounts();
		accounts = response.data?.data?.accounts ?? [];
		error = '';
	} catch (err) {
		error = 'Unable to load linked accounts.';
		toastError(error);
		console.error(err);
	} finally {
		if (showSpinner) {
			loadingAccounts = false;
		}
	}
};

const fetchProviders = async () => {
	try {
		const response = await listOAuthProviders();
		providers = response.data?.data?.providers ?? [];
		error = '';
	} catch (err) {
		error = 'Unable to load provider metadata.';
		toastError(error);
		console.error(err);
	} finally {
		loadingProviders = false;
	}
};

	const handleOAuthMessage = (event: MessageEvent) => {
		if (event.origin !== API_ORIGIN) {
			return;
		}

		const payload = event.data;
		if (!payload || payload.type !== 'oauth:result') {
			return;
		}

		if (oauthWindow && !oauthWindow.closed) {
			oauthWindow.close();
		}
		oauthWindow = null;

		if (payload.mode !== 'link') {
			return;
		}

		if (payload.success) {
			info = 'Provider linked successfully.';
			error = '';
			toastSuccess(info);
			void refreshAccounts();
		} else {
			error = payload.message ?? 'Unable to link provider.';
			toastError(error);
		}
	};

	const linkProvider = async (provider: string) => {
		if (!browser || !isAuthenticated || typeof window === 'undefined') {
			return;
		}

		error = '';
		info = '';
		oauthLoading = { ...oauthLoading, [provider]: true };

		try {
			const response = await startOAuth(
				provider,
				'link',
				window.location.origin,
				getDeviceFingerprint(),
				getDeviceName(),
				getUserAgent()
			);
			const payload = response.data?.data;
			if (!payload?.authorization_url) {
				throw new Error('Missing authorization URL');
			}

			if (oauthWindow && !oauthWindow.closed) {
				oauthWindow.close();
			}

			oauthWindow = window.open(
				payload.authorization_url,
				`oauth_link_${provider}`,
				'width=500,height=640,menubar=no,toolbar=no'
			);
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to start OAuth linking flow.');
			toastError(error);
			console.error(err);
		} finally {
			oauthLoading = { ...oauthLoading, [provider]: false };
		}
	};

	const unlinkProvider = async (provider: string) => {
		error = '';
		info = '';
		oauthLoading = { ...oauthLoading, [provider]: true };

		try {
			await unlinkOAuthAccount(provider);
			info = 'Provider disconnected.';
			toastInfo(info);
			await refreshAccounts();
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to disconnect provider.');
			toastError(error);
			console.error(err);
		} finally {
			oauthLoading = { ...oauthLoading, [provider]: false };
		}
	};

	onMount(() => {
		if (!browser || typeof window === 'undefined') {
			return;
		}

		fetchProviders();
		const unsubscribe = authStore.subscribe((state) => {
			isAuthenticated = state.isAuthenticated;
			if (isAuthenticated) {
				refreshAccounts();
			} else {
				accounts = [];
			}
		});

		window.addEventListener('message', handleOAuthMessage);

		return () => {
			unsubscribe();
		};
	});

	onDestroy(() => {
		if (!browser || typeof window === 'undefined') {
			return;
		}

		window.removeEventListener('message', handleOAuthMessage);
		if (oauthWindow && !oauthWindow.closed) {
			oauthWindow.close();
		}
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Connected providers</h2>
		<p class="text-sm text-slate-400">
			Link your Woragis account with external identity providers to enable passwordless sign-in.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">
			{error}
		</p>
	{/if}

	{#if info}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">
			{info}
		</p>
	{/if}

	{#if !isAuthenticated}
		<div class="rounded border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-300">
			<p>You need to sign in before managing linked providers.</p>
			<p class="mt-2">
				<a class="text-primary hover:underline" href="/auth/login">Go to sign in</a>
			</p>
		</div>
	{:else if loadingProviders}
		<p class="text-sm text-slate-400">Loading providers...</p>
	{:else if providers.length === 0}
		<p class="text-sm text-slate-400">No OAuth providers are configured for this environment.</p>
	{:else}
		<div class="space-y-4">
			{#if loadingAccounts}
				<p class="text-sm text-slate-400">Checking linked accounts...</p>
			{/if}

			<ul class="space-y-3">
				{#each providers as provider}
					{@const connected = findAccount(provider.id)}
					<li class="flex items-center justify-between rounded border border-slate-800 bg-slate-900/60 px-4 py-3">
						<div>
							<p class="text-sm font-semibold text-slate-100">{provider.name}</p>
							{#if connected}
								<p class="text-xs text-slate-400">
									Linked {new Date(connected.linked_at).toLocaleString()}
								</p>
							{:else}
								<p class="text-xs text-slate-400">Not linked</p>
							{/if}
						</div>
						<div class="flex items-center gap-2">
							{#if connected}
								<button
									class="rounded border border-slate-700 bg-slate-950 px-3 py-1 text-xs font-semibold text-slate-100 transition hover:bg-slate-800 disabled:opacity-50"
									type="button"
									on:click={() => unlinkProvider(provider.id)}
									disabled={oauthLoading[provider.id]}
								>
									{#if oauthLoading[provider.id]}
										Removing...
									{:else}
										Disconnect
									{/if}
								</button>
							{:else}
								<button
									class="rounded bg-emerald-500 px-3 py-1 text-xs font-semibold text-slate-900 transition disabled:opacity-50"
									type="button"
									on:click={() => linkProvider(provider.id)}
									disabled={oauthLoading[provider.id]}
								>
									{#if oauthLoading[provider.id]}
										Connecting...
									{:else}
										Connect
									{/if}
								</button>
							{/if}
						</div>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</section>

