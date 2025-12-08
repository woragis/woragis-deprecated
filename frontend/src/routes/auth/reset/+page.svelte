<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

import { confirmPasswordReset } from '$lib/api/auth';
import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';

	let token = '';
	let password = '';
	let confirmPassword = '';
	let loading = false;
	let message = '';
	let error = '';

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		error = '';
		message = '';

		if (password !== confirmPassword) {
			error = 'Passwords do not match.';
			toastError(error);
			return;
		}

		if (!token.trim()) {
			error = 'Reset token missing.';
			toastError(error);
			return;
		}

		loading = true;
		try {
			await confirmPasswordReset(token.trim(), password);
			message = 'Password updated successfully.';
			toastSuccess(message);
			setTimeout(() => goto('/auth/login'), 1500);
		} catch (err: unknown) {
			error = getApiErrorMessage(
				err,
				'Unable to update password. The token might be invalid or expired.'
			);
			toastError(error);
			console.error(err);
		} finally {
			loading = false;
		}
	};

	onMount(() => {
		const unsubscribe = page.subscribe((p) => {
			const tokenParam = p.url.searchParams.get('token');
			if (tokenParam) {
				token = tokenParam;
			}
		});
		return () => {
			unsubscribe();
		};
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Reset password</h2>
		<p class="text-sm text-slate-400">
			Create a new password for your Woragis account. Passwords must be at least eight characters long.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if message}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">{message}</p>
	{/if}

	<form class="card" on:submit={handleSubmit}>
		<label class="flex flex-col gap-1 text-sm text-slate-200">
			Token
			<input
				class="input"
				type="text"
				placeholder="Paste your reset token"
				bind:value={token}
				required
			/>
		</label>
		<label class="flex flex-col gap-1 text-sm text-slate-200">
			New password
			<input
				class="input"
				type="password"
				placeholder="********"
				minlength="8"
				bind:value={password}
				required
			/>
		</label>
		<label class="flex flex-col gap-1 text-sm text-slate-200">
			Confirm password
			<input
				class="input"
				type="password"
				placeholder="********"
				minlength="8"
				bind:value={confirmPassword}
				required
			/>
		</label>

		<button class="btn-primary" type="submit" disabled={loading}>
			{#if loading}
				Updating password...
			{:else}
				Update password
			{/if}
		</button>
	</form>
</section>

<style>
	.card {
		border: 1px solid rgba(71, 85, 105, 0.45);
		border-radius: 0.85rem;
		background: rgba(15, 15, 15, 0.6);
		padding: 1.75rem;
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		max-width: 32rem;
	}

	.input {
		border: 1px solid rgba(255, 255, 255, 0.12);
		border-radius: 0.65rem;
		background: rgba(2, 6, 23, 0.85);
		padding: 0.65rem 0.85rem;
		color: #f8fafc;
	}

	.input:focus-visible {
		outline: 2px solid rgba(255, 255, 255, 0.2);
		outline-offset: 2px;
	}

	.btn-primary {
		font-size: 0.8rem;
		text-transform: uppercase;
		font-weight: 600;
		letter-spacing: 0.08em;
		padding: 0.45rem 1rem;
		border-radius: 0.65rem;
		transition: background 120ms ease, color 120ms ease, border-color 120ms ease;
		background: rgba(255, 255, 255, 0.08);
		color: #a3a3a3;
		border: 1px solid rgba(59, 130, 246, 0.5);
	}

	.btn-primary:hover,
	.btn-primary:focus-visible {
		background: rgba(59, 130, 246, 0.35);
		color: #dbeafe;
		border-color: rgba(59, 130, 246, 0.7);
		outline: none;
	}
</style>

