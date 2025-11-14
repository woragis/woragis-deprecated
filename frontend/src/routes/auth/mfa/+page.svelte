<script lang="ts">
	import { onMount } from 'svelte';

import { enableMFA, verifyMFA, disableMFA } from '$lib/api/auth';
import { authStore } from '$lib';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

	let isAuthenticated = false;
	let mfaEnabled = false;
	let email = '';

	let issuer = 'Woragis';
	let label = '';
	let verificationCode = '';

	let provisioningURI = '';
	let secret = '';
	let backupCodes: string[] = [];

	let info = '';
	let error = '';
	let loadingSetup = false;
	let loadingVerify = false;
	let loadingDisable = false;

	const resetSetup = () => {
		provisioningURI = '';
		secret = '';
		backupCodes = [];
		verificationCode = '';
	};

	const handleStartSetup = async () => {
		if (!label.trim()) {
			error = 'Enter a label to display in your authenticator app.';
			toastError(error);
			return;
		}

		loadingSetup = true;
		error = '';
		info = '';
		try {
			const response = await enableMFA({
				issuer: issuer.trim() || 'Woragis',
				label: label.trim()
			});
			const payload = response.data?.data;
			if (payload) {
				provisioningURI = payload.provisioning_uri;
				secret = payload.secret;
				backupCodes = payload.backup_codes ?? [];
				info = 'Scan the QR or enter the secret, then enter a code below to verify.';
				toastInfo('MFA secret generated. Complete verification to finish setup.');
			}
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to start MFA setup.');
			toastError(error);
			console.error(err);
		} finally {
			loadingSetup = false;
		}
	};

	const handleVerify = async () => {
		if (!verificationCode.trim()) {
			error = 'Enter the code from your authenticator app.';
			toastError(error);
			return;
		}

		loadingVerify = true;
		error = '';
		info = '';
		try {
			await verifyMFA(verificationCode.trim());
			authStore.updateUser({ mfa_enabled: true });
			mfaEnabled = true;
			info = 'Multi-factor authentication enabled.';
			toastSuccess(info);
			resetSetup();
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to verify MFA code.');
			toastError(error);
			console.error(err);
		} finally {
			loadingVerify = false;
		}
	};

	const handleDisable = async () => {
		loadingDisable = true;
		error = '';
		info = '';
		try {
			await disableMFA();
			authStore.updateUser({ mfa_enabled: false });
			mfaEnabled = false;
			resetSetup();
			info = 'Multi-factor authentication disabled.';
			toastInfo(info);
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to disable MFA.');
			toastError(error);
			console.error(err);
		} finally {
			loadingDisable = false;
		}
	};

	onMount(() => {
		const unsubscribe = authStore.subscribe((state) => {
			isAuthenticated = state.isAuthenticated;
			mfaEnabled = state.user?.mfa_enabled ?? false;
			email = state.user?.email ?? '';
			if (!label && email) {
				label = `${email}`;
			}
		});
		return () => {
			unsubscribe();
		};
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Multi-factor authentication</h2>
		<p class="text-sm text-slate-400">
			Enhance account security by requiring a time-based one-time code in addition to your password.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if info}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">{info}</p>
	{/if}

	{#if !isAuthenticated}
		<div class="card">
			<p class="text-sm text-slate-300">Sign in to manage multi-factor authentication.</p>
			<a class="btn-outline" href="/auth/login">Go to sign in</a>
		</div>
	{:else}
		{#if mfaEnabled}
			<div class="card">
				<h3 class="card__title">MFA is enabled</h3>
				<p class="text-sm text-slate-300">
					You'll be asked for a six-digit code from your authenticator app each time you sign in.
				</p>
				<button class="btn-danger" type="button" on:click={handleDisable} disabled={loadingDisable}>
					{#if loadingDisable}
						Disabling...
					{:else}
						Disable MFA
					{/if}
				</button>
			</div>
		{:else}
			<div class="card">
				<h3 class="card__title">Enable MFA</h3>
				<p class="text-sm text-slate-300">
					Provide a label to identify this account in your authenticator app. The issuer is shown alongside the label.
				</p>
				<form class="space-y-4" on:submit|preventDefault={handleStartSetup}>
					<label class="flex flex-col gap-1 text-sm text-slate-200">
						Issuer
						<input class="input" type="text" bind:value={issuer} placeholder="Woragis" />
					</label>
					<label class="flex flex-col gap-1 text-sm text-slate-200">
						Label
						<input class="input" type="text" bind:value={label} placeholder="you@example.com" required />
					</label>
					<button class="btn-primary" type="submit" disabled={loadingSetup}>
						{#if loadingSetup}
							Generating secret...
						{:else}
							Generate secret
						{/if}
					</button>
				</form>
			</div>

			{#if secret}
				<div class="card space-y-4">
					<h3 class="card__title">Verify and activate</h3>
					<p class="text-sm text-slate-300">
						Scan the QR code or manually enter the secret below. Then enter a code from your authenticator app to
						activate MFA.
					</p>

					{#if provisioningURI}
						<div class="qr-placeholder">
							<span>Scan with your authenticator app</span>
							<a class="text-primary hover:underline text-sm" href={provisioningURI} rel="noopener">
								Open in authenticator
							</a>
							<code class="uri">{provisioningURI}</code>
						</div>
					{/if}

					<div class="secret-block">
						<span class="secret">{secret}</span>
						<span class="hint">Keep this secret safe.</span>
					</div>

					{#if backupCodes.length}
						<div>
							<h4 class="text-sm font-semibold text-slate-200">Backup codes</h4>
							<p class="text-xs text-slate-400">Store these one-time codes in a secure place.</p>
							<ul class="backup-codes">
								{#each backupCodes as code}
									<li>{code}</li>
								{/each}
							</ul>
						</div>
					{/if}

					<form class="space-y-4" on:submit|preventDefault={handleVerify}>
						<label class="flex flex-col gap-1 text-sm text-slate-200">
							Authenticator code
							<input
								class="input"
								type="text"
								inputmode="numeric"
								minlength="6"
								maxlength="6"
								placeholder="123456"
								bind:value={verificationCode}
								required
							/>
						</label>
						<button class="btn-primary" type="submit" disabled={loadingVerify}>
							{#if loadingVerify}
								Verifying...
							{:else}
								Verify and enable
							{/if}
						</button>
					</form>
				</div>
			{/if}
		{/if}
	{/if}
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
	.btn-outline,
	.btn-danger {
		font-size: 0.8rem;
		text-transform: uppercase;
		font-weight: 600;
		letter-spacing: 0.08em;
		padding: 0.45rem 1rem;
		border-radius: 0.65rem;
		transition: background 120ms ease, color 120ms ease, border-color 120ms ease;
		text-align: center;
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
		display: inline-flex;
		justify-content: center;
		align-items: center;
	}

	.btn-outline:hover,
	.btn-outline:focus-visible {
		background: rgba(148, 163, 184, 0.12);
		color: #f8fafc;
		border-color: rgba(148, 163, 184, 0.6);
		outline: none;
	}

	.btn-danger {
		background: rgba(239, 68, 68, 0.18);
		color: #f87171;
		border: 1px solid rgba(239, 68, 68, 0.45);
	}

	.btn-danger:hover,
	.btn-danger:focus-visible {
		background: rgba(239, 68, 68, 0.28);
		color: #fee2e2;
		border-color: rgba(239, 68, 68, 0.65);
		outline: none;
	}

	.secret-block {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		padding: 0.85rem 1rem;
		border-radius: 0.75rem;
		background: rgba(15, 118, 110, 0.14);
		border: 1px dashed rgba(45, 212, 191, 0.35);
	}

	.secret {
		font-family: 'Fira Code', 'SFMono-Regular', Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
		letter-spacing: 0.12em;
		font-size: 1rem;
		color: #5eead4;
	}

	.hint {
		font-size: 0.7rem;
		color: rgba(148, 163, 184, 0.9);
		text-transform: uppercase;
		letter-spacing: 0.12em;
	}

	.backup-codes {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(8rem, 1fr));
		gap: 0.5rem;
		padding: 0;
		margin: 0;
		list-style: none;
	}

	.backup-codes li {
		font-family: 'Fira Code', monospace;
		font-size: 0.85rem;
		background: rgba(30, 64, 175, 0.18);
		color: #bfdbfe;
		border-radius: 0.55rem;
		padding: 0.4rem 0.6rem;
		text-align: center;
	}

	.qr-placeholder {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 0.35rem;
		font-size: 0.8rem;
		color: rgba(148, 163, 184, 0.9);
	}

	.uri {
		font-family: 'Fira Code', monospace;
		font-size: 0.65rem;
		color: rgba(148, 163, 184, 0.8);
		word-break: break-all;
	}
</style>

