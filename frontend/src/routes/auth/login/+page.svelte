<!-- eslint-disable svelte/no-navigation-without-resolve -->
<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { onDestroy, onMount } from 'svelte';
	import { API_BASE_URL } from '@clients/apiClient';
	import { login as loginRequest, startOAuth } from '$lib/api/auth';
	import type { LoginResponsePayload } from '$lib/api/auth';
	import { authStore } from '$lib';
	import { getDeviceFingerprint, getDeviceName, getUserAgent } from '$lib/device';
	import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';

	interface LoginResponseEnvelope {
		data: LoginResponsePayload | null;
		success: boolean;
	}

	const oauthProviders = [
		{ id: 'google', label: 'Continue with Google' },
		{ id: 'github', label: 'Continue with GitHub' },
		{ id: 'microsoft', label: 'Continue with Microsoft' }
	] as const;

	const API_ORIGIN = browser ? new URL(API_BASE_URL).origin : '';

	let email = '';
	let password = '';
	let loading = false;
	let error = '';
	let oauthError = '';
	let oauthLoading: Record<string, boolean> = {};
	let oauthWindow: Window | null = null;

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		error = '';
		loading = true;

		try {
			const response = await loginRequest({
				email,
				password,
				device_fingerprint: getDeviceFingerprint(),
				device_name: getDeviceName(),
				user_agent: getUserAgent()
			});
			const payload = response.data?.data;

			if (payload?.user && payload?.access_token) {
				authStore.setSession(payload.user, payload.access_token, {
					refreshToken: payload.refresh_token ?? null,
					sessionId: payload.session_id ?? null
				});
				toastSuccess('Signed in successfully.');
				await goto('/');
				return;
			}

			error = 'Unable to sign in. Please try again.';
			toastError(error);
		} catch (err) {
			error = getApiErrorMessage(err, 'Unable to sign in. Check your credentials and try again.');
			toastError(error);
			console.error(err);
		} finally {
			loading = false;
		}
	};

	const handleOAuthMessage = (event: MessageEvent) => {
		if (event.origin !== API_ORIGIN) {
			return;
		}

		const payload = event.data;
		if (!payload || payload.type !== 'oauth:result' || payload.mode !== 'login') {
			return;
		}

		if (oauthWindow && !oauthWindow.closed) {
			oauthWindow.close();
		}
		oauthWindow = null;

		if (payload.success && payload.payload) {
			const data = payload.payload as LoginResponsePayload;
			if (data?.user && data?.access_token) {
				authStore.setSession(data.user, data.access_token, {
					refreshToken: data.refresh_token ?? null,
					sessionId: data.session_id ?? null
				});
				toastSuccess('Signed in successfully.');
				goto('/');
				return;
			}
			oauthError = 'Unable to complete sign in.';
			toastError(oauthError);
			return;
		}

		oauthError = payload.message ?? 'OAuth sign in failed.';
		toastError(oauthError);
	};

	const handleOAuth = async (provider: string) => {
		if (!browser || typeof window === 'undefined') {
			return;
		}

		oauthError = '';
		oauthLoading = { ...oauthLoading, [provider]: true };

		try {
			const response = await startOAuth(
				provider,
				'login',
				window.location.origin,
				getDeviceFingerprint(),
				getDeviceName(),
				getUserAgent()
			);
			const payload = response.data?.data;

			if (!payload?.authorization_url) {
				throw new Error('Authorization URL missing from response');
			}

			if (oauthWindow && !oauthWindow.closed) {
				oauthWindow.close();
			}

			oauthWindow = window.open(
				payload.authorization_url,
				`oauth_${provider}`,
				'width=500,height=640,menubar=no,toolbar=no'
			);
		} catch (err) {
			oauthError = getApiErrorMessage(err, 'Unable to start OAuth flow. Please try again.');
			toastError(oauthError);
			console.error(err);
		} finally {
			oauthLoading = { ...oauthLoading, [provider]: false };
		}
	};

	onMount(() => {
		if (!browser || typeof window === 'undefined') {
			return;
		}

		window.addEventListener('message', handleOAuthMessage);
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

<div class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-lg font-semibold">Sign in</h2>
		<p class="text-sm text-slate-400">Choose a provider or sign in with your email.</p>
	</div>

	{#if oauthError}
		<p class="rounded border border-orange-500/30 bg-orange-500/10 px-3 py-2 text-sm text-orange-100">
			{oauthError}
		</p>
	{/if}

	<div class="space-y-3">
		{#each oauthProviders as provider}
			<button
				class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm font-semibold text-slate-100 transition hover:bg-slate-800 disabled:opacity-50"
				type="button"
				on:click={() => handleOAuth(provider.id)}
				disabled={oauthLoading[provider.id]}
			>
				{#if oauthLoading[provider.id]}
					Connecting...
				{:else}
					{provider.label}
				{/if}
			</button>
		{/each}
	</div>

	<div class="flex items-center gap-3 text-xs text-slate-500">
		<div class="h-px flex-1 bg-slate-800" />
		<span>or continue with email</span>
		<div class="h-px flex-1 bg-slate-800" />
	</div>

	<form class="space-y-6" on:submit={handleSubmit}>
		{#if error}
			<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
				{error}
			</p>
		{/if}

		<div class="space-y-4">
			<label class="flex flex-col gap-1 text-sm text-slate-200">
				Email
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					type="email"
					placeholder="you@example.com"
					bind:value={email}
					required
				/>
			</label>
			<label class="flex flex-col gap-1 text-sm text-slate-200">
				Password
				<input
					class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
					type="password"
					placeholder="********"
					bind:value={password}
					required
				/>
			</label>
		</div>

		<button
			class="w-full rounded bg-emerald-500 px-3 py-2 text-sm font-semibold text-slate-900 transition-opacity disabled:opacity-50"
			type="submit"
			disabled={loading}
		>
			{#if loading}
				Signing in...
			{:else}
				Sign in
			{/if}
		</button>

		<div class="text-center text-xs text-slate-400 space-y-1">
			<p>
				Don't have an account?
				<a class="text-primary hover:underline" href="/auth/register">Create one</a>
			</p>
			<p>
				Forgot your password?
				<a class="text-primary hover:underline" href="/auth/forgot">Reset it</a>
			</p>
			<p>
				Need to confirm your email?
				<a class="text-primary hover:underline" href="/auth/confirm-email">Resend confirmation</a>
			</p>
		</div>
	</form>
</div>
