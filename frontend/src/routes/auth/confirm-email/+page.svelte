<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

import { confirmEmail, resendConfirmation } from '$lib/api/auth';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

	let token = '';
	let email = '';
	let loadingConfirm = false;
	let loadingResend = false;
	let confirmMessage = '';
	let resendMessage = '';
	let error = '';

	const loadPendingEmail = () => {
		if (!browser) return '';
		try {
			return localStorage.getItem('woragis_pending_confirmation') ?? '';
		} catch (err) {
			console.warn('Unable to read pending confirmation email.', err);
			return '';
		}
	};

	const clearPendingEmail = () => {
		if (!browser) return;
		try {
			localStorage.removeItem('woragis_pending_confirmation');
		} catch (err) {
			console.warn('Unable to clear pending confirmation email.', err);
		}
	};

	const handleConfirm = async () => {
		if (!token.trim()) {
			error = 'Enter the token from your confirmation email.';
			return;
		}

		loadingConfirm = true;
		confirmMessage = '';
		error = '';
		try {
			await confirmEmail(token.trim());
			clearPendingEmail();
			confirmMessage = 'Email confirmed. You can now sign in.';
			toastSuccess(confirmMessage);
		} catch (err: unknown) {
			error = getApiErrorMessage(
				err,
				'Unable to confirm email. The token might be invalid or expired.'
			);
			toastError(error);
			console.error(err);
		} finally {
			loadingConfirm = false;
		}
	};

	const handleResend = async () => {
		if (!email.trim()) {
			error = 'Enter the email address you registered with.';
			return;
		}

		loadingResend = true;
		resendMessage = '';
		error = '';
		try {
			await resendConfirmation(email.trim());
			storePendingEmail(email.trim());
			resendMessage = 'Confirmation email sent. Check your inbox.';
			toastInfo(resendMessage);
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to resend confirmation email.');
			toastError(error);
			console.error(err);
		} finally {
			loadingResend = false;
		}
	};

	const storePendingEmail = (value: string) => {
		if (!browser) return;
		try {
			localStorage.setItem('woragis_pending_confirmation', value);
		} catch (err) {
			console.warn('Unable to persist pending confirmation email.', err);
		}
	};

	onMount(() => {
		const params = page?.subscribe((p) => {
			const tokenParam = p.url.searchParams.get('token');
			const emailParam = p.url.searchParams.get('email');
			if (tokenParam) {
				token = tokenParam;
				handleConfirm();
			}
			if (emailParam) {
				email = emailParam;
			} else {
				const cached = loadPendingEmail();
				if (cached) {
					email = cached;
				}
			}
		});

		return () => {
			params?.();
		};
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Confirm your email</h2>
		<p class="text-sm text-slate-400">
			Complete registration by entering the token from your confirmation email. If you no longer have the email, request
			a new one below.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if confirmMessage}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">
			{confirmMessage}
			<a class="text-primary hover:underline" href="/auth/login"> Sign in</a>.
		</p>
	{/if}

	<div class="card">
		<h3 class="card__title">Enter confirmation token</h3>
		<form class="space-y-4" on:submit|preventDefault={handleConfirm}>
			<label class="flex flex-col gap-1 text-sm text-slate-200">
				Token
				<input
					class="input"
					type="text"
					placeholder="Paste your confirmation token"
					bind:value={token}
					required
				/>
			</label>
			<button class="btn-primary" type="submit" disabled={loadingConfirm}>
				{#if loadingConfirm}
					Confirming...
				{:else}
					Confirm email
				{/if}
			</button>
		</form>
	</div>

	<div class="card">
		<h3 class="card__title">Resend confirmation email</h3>
		{#if resendMessage}
			<p class="text-sm text-emerald-200">{resendMessage}</p>
		{/if}
		<form class="space-y-4" on:submit|preventDefault={handleResend}>
			<label class="flex flex-col gap-1 text-sm text-slate-200">
				Email
				<input
					class="input"
					type="email"
					placeholder="you@example.com"
					bind:value={email}
					required
				/>
			</label>
			<button class="btn-outline" type="submit" disabled={loadingResend}>
				{#if loadingResend}
					Sending...
				{:else}
					Send confirmation email
				{/if}
			</button>
		</form>
	</div>
</section>

<style>
	.card {
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		background: rgba(15, 23, 42, 0.75);
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.card__title {
		font-size: 1rem;
		font-weight: 600;
		color: #e2e8f0;
	}

	.input {
		border: 1px solid rgba(71, 85, 105, 0.6);
		border-radius: 0.65rem;
		background: rgba(2, 6, 23, 0.85);
		padding: 0.65rem 0.85rem;
		color: #f8fafc;
	}

	.input:focus-visible {
		outline: 2px solid rgba(94, 234, 212, 0.6);
		outline-offset: 2px;
	}

	.btn-primary,
	.btn-outline {
		font-size: 0.8rem;
		text-transform: uppercase;
		font-weight: 600;
		letter-spacing: 0.08em;
		padding: 0.45rem 1rem;
		border-radius: 0.65rem;
		transition: background 120ms ease, color 120ms ease, border-color 120ms ease;
	}

	.btn-primary {
		background: rgba(16, 185, 129, 0.2);
		color: #34d399;
		border: 1px solid rgba(16, 185, 129, 0.6);
	}

	.btn-primary:hover,
	.btn-primary:focus-visible {
		background: rgba(16, 185, 129, 0.35);
		color: #ecfdf5;
		border-color: rgba(16, 185, 129, 0.8);
		outline: none;
	}

	.btn-outline {
		background: transparent;
		color: rgba(148, 163, 184, 0.95);
		border: 1px solid rgba(148, 163, 184, 0.4);
	}

	.btn-outline:hover,
	.btn-outline:focus-visible {
		background: rgba(148, 163, 184, 0.12);
		color: #f8fafc;
		border-color: rgba(148, 163, 184, 0.6);
		outline: none;
	}
</style>

