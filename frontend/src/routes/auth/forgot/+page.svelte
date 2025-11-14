<script lang="ts">
import { requestPasswordReset } from '$lib/api/auth';
import { getApiErrorMessage, toastError, toastInfo } from '$lib/utils/toast';

	let email = '';
	let loading = false;
	let message = '';
	let error = '';

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		error = '';
		message = '';

		loading = true;
		try {
			await requestPasswordReset(email.trim());
			message = 'If an account exists for this email, a reset link has been sent.';
			toastInfo(message);
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to send password reset email.');
			toastError(error);
			console.error(err);
		} finally {
			loading = false;
		}
	};
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Forgot password</h2>
		<p class="text-sm text-slate-400">
			Enter your email address and we'll send a link to reset your password. If you no longer have access to your email,
			contact an administrator for assistance.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if message}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">{message}</p>
	{/if}

	<form class="card" on:submit={handleSubmit}>
		<label class="flex flex-col gap-2 text-sm text-slate-200">
			Email
			<input
				class="input"
				type="email"
				placeholder="you@example.com"
				bind:value={email}
				required
			/>
		</label>

		<button class="btn-primary" type="submit" disabled={loading}>
			{#if loading}
				Sending reset link...
			{:else}
				Send reset link
			{/if}
		</button>

		<p class="text-xs text-slate-400">
			Remembered your password?
			<a class="text-primary hover:underline" href="/auth/login">Sign in</a>
		</p>
	</form>
</section>

<style>
	.card {
		border: 1px solid rgba(71, 85, 105, 0.45);
		border-radius: 0.85rem;
		background: rgba(15, 23, 42, 0.8);
		padding: 1.75rem;
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		max-width: 28rem;
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

	.btn-primary {
		font-size: 0.8rem;
		text-transform: uppercase;
		font-weight: 600;
		letter-spacing: 0.08em;
		padding: 0.45rem 1rem;
		border-radius: 0.65rem;
		transition: background 120ms ease, color 120ms ease, border-color 120ms ease;
		background: rgba(59, 130, 246, 0.2);
		color: #60a5fa;
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

