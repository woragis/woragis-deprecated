<script lang="ts">
import { onMount } from 'svelte';

import { listSessions, logout as logoutSession, revokeOtherSessions } from '$lib/api/auth';
import { authStore } from '$lib';
import { getApiErrorMessage, toastError, toastInfo, toastSuccess } from '$lib/utils/toast';

	interface SessionRow {
		id: string;
		device_id: string;
		user_agent: string;
		ip: string;
		created_at: string;
		expires_at: string;
		last_seen_at: string;
		is_revoked: boolean;
	}

let sessions: SessionRow[] = [];
	let loading = false;
	let error = '';
	let info = '';
	let isAuthenticated = false;
	let currentSessionId: string | null = null;

const fetchSessions = async (showSpinner = true) => {
	if (!isAuthenticated) {
		sessions = [];
		return;
	}

	if (showSpinner) {
		loading = true;
	}
	error = '';

	try {
		const response = await listSessions();
		sessions = response.data?.data?.sessions ?? [];
	} catch (err) {
		error = 'Unable to load sessions.';
		toastError(error);
		console.error(err);
	} finally {
		if (showSpinner) {
			loading = false;
		}
	}
};

	const formatDate = (value: string) => {
		if (!value) return '—';
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return value;
		return date.toLocaleString();
	};

	const handleRevokeOthers = async () => {
		if (!currentSessionId) {
			return;
		}

		error = '';
		info = '';
		loading = true;
		try {
			await revokeOtherSessions(currentSessionId);
			info = 'Other sessions revoked.';
			toastSuccess(info);
			await fetchSessions(false);
		} catch (err) {
			error = getApiErrorMessage(err, 'Unable to revoke other sessions.');
			toastError(error);
			console.error(err);
		} finally {
			loading = false;
		}
	};

	const handleLogoutSession = async (sessionId: string) => {
		error = '';
		info = '';
		loading = true;
		try {
			await logoutSession(sessionId);
			if (sessionId === currentSessionId) {
				authStore.clear();
			} else {
				info = 'Session revoked.';
				toastInfo(info);
			}
			await fetchSessions(false);
		} catch (err) {
			error = getApiErrorMessage(err, 'Unable to revoke session.');
			toastError(error);
			console.error(err);
		} finally {
			loading = false;
		}
	};

	onMount(() => {
		const unsubscribe = authStore.subscribe((state) => {
			isAuthenticated = state.isAuthenticated;
			currentSessionId = state.sessionId ?? null;
			if (isAuthenticated) {
				fetchSessions();
			} else {
				sessions = [];
			}
		});

		return () => {
			unsubscribe();
		};
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Active sessions</h2>
		<p class="text-sm text-slate-400">
			Track and manage devices that have access to your Woragis workspace. Revoke sessions individually or sign out of
			all other devices.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if info}
		<p class="rounded border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">{info}</p>
	{/if}

	{#if !isAuthenticated}
		<div class="rounded border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-300">
			<p>You need to sign in to view session details.</p>
			<p class="mt-2">
				<a class="text-primary hover:underline" href="/auth/login">Go to sign in</a>
			</p>
		</div>
	{:else if loading && sessions.length === 0}
		<p class="text-sm text-slate-400">Loading sessions...</p>
	{:else if sessions.length === 0}
		<p class="text-sm text-slate-400">No active sessions found.</p>
	{:else}
		<div class="space-y-4">
			<div class="flex flex-wrap items-center gap-3">
				<button class="btn-outline" type="button" on:click={() => fetchSessions()} disabled={loading}>
					Refresh
				</button>
				{#if currentSessionId}
					<button class="btn-primary" type="button" on:click={handleRevokeOthers} disabled={loading}>
						Sign out other devices
					</button>
				{/if}
			</div>

			<div class="session-table">
				<table>
					<thead>
						<tr>
							<th>Device</th>
							<th>User agent</th>
							<th>IP</th>
							<th>Last activity</th>
							<th>Status</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						{#each sessions as session (session.id)}
							{@const isCurrent = currentSessionId === session.id}
							<tr class:current={isCurrent}>
								<td>
									<div class="cell-primary">
										<span class="device-label">{session.device_id}</span>
										{#if isCurrent}
											<span class="badge">This device</span>
										{/if}
									</div>
									<div class="cell-secondary">
										Since {formatDate(session.created_at)}
									</div>
								</td>
								<td>
									<div class="cell-secondary ellipsis" title={session.user_agent}>
										{session.user_agent || 'Unknown'}
									</div>
								</td>
								<td>{session.ip || '—'}</td>
								<td>{formatDate(session.last_seen_at)}</td>
								<td>{session.is_revoked ? 'Revoked' : 'Active'}</td>
								<td>
									{#if !session.is_revoked}
										<button
											class="btn-link"
											type="button"
											on:click={() => handleLogoutSession(session.id)}
											disabled={loading}
										>
											{isCurrent ? 'Sign out' : 'Revoke'}
										</button>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</section>

<style>
	.session-table {
		overflow-x: auto;
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		background: rgba(15, 23, 42, 0.75);
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.85rem;
	}

	thead tr {
		background: rgba(15, 23, 42, 0.9);
	}

	th,
	td {
		padding: 0.85rem 1rem;
		text-align: left;
		vertical-align: top;
		border-bottom: 1px solid rgba(71, 85, 105, 0.25);
	}

	th {
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: rgba(148, 163, 184, 0.9);
	}

	tbody tr:last-child td {
		border-bottom: none;
	}

	tr.current {
		background: rgba(16, 185, 129, 0.08);
	}

	.cell-primary {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
	}

	.cell-secondary {
		font-size: 0.75rem;
		color: rgba(148, 163, 184, 0.9);
	}

	.device-label {
		font-weight: 600;
		color: #e2e8f0;
		word-break: break-all;
	}

	.badge {
		display: inline-flex;
		align-items: center;
		border-radius: 9999px;
		padding: 0.1rem 0.55rem;
		font-size: 0.65rem;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		background: rgba(94, 234, 212, 0.1);
		color: rgba(94, 234, 212, 0.95);
	}

	.btn-primary,
	.btn-outline,
	.btn-link {
		font-size: 0.75rem;
		text-transform: uppercase;
		font-weight: 600;
		letter-spacing: 0.08em;
		border-radius: 9999px;
		padding: 0.35rem 0.9rem;
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
		border-color: rgba(16, 185, 129, 0.8);
		color: #ecfdf5;
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

	.btn-link {
		background: none;
		border: none;
		color: rgba(248, 250, 252, 0.85);
		padding: 0.2rem 0.35rem;
		text-decoration: underline;
		cursor: pointer;
	}

	.btn-link:hover,
	.btn-link:focus-visible {
		color: #f8fafc;
		outline: none;
	}

	.ellipsis {
		max-width: 18rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	@media (max-width: 640px) {
		th,
		td {
			padding: 0.75rem 0.6rem;
		}

		.ellipsis {
			max-width: 12rem;
		}
	}
</style>

