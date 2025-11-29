<script lang="ts">
	import { onMount } from 'svelte';
	import { Users, Key, Activity, Settings, Shield, AlertCircle, ArrowRight } from 'lucide-svelte';
	import { listUsers } from '$lib/api/admin';
	import { getCurrentUser } from '$lib/api/auth';

	let currentUser: any = $state(null);
	let userCount = $state(0);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await Promise.all([loadCurrentUser(), loadUserStats()]);
	});

	async function loadCurrentUser() {
		try {
			currentUser = await getCurrentUser();
		} catch (err) {
			console.error('Error loading current user:', err);
		}
	}

	async function loadUserStats() {
		loading = true;
		error = null;
		try {
			const response = await listUsers({ limit: 1, offset: 0 });
			userCount = response.total;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load stats';
			console.error('Error loading user stats:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="page-container">
	<!-- Header -->
	<div class="page-header">
		<div>
			<h1 class="page-title">Admin Dashboard</h1>
			<p class="page-description">Manage your application and users</p>
		</div>
		{#if currentUser}
			<div class="user-info">
				<span class="user-email">{currentUser.email}</span>
				<span class="user-role">{currentUser.role}</span>
			</div>
		{/if}
	</div>

	<!-- Error State -->
	{#if error}
		<div class="alert alert-error">
			<AlertCircle class="icon" />
			<p>{error}</p>
		</div>
	{/if}

	<!-- Stats Cards -->
	<div class="stats-grid">
		<a href="/admin/users" class="stat-card">
			<div class="stat-icon stat-icon-users">
				<Users class="icon" />
			</div>
			<div class="stat-content">
				<div class="stat-value">
					{#if loading}
						<span class="loading-dots">...</span>
					{:else}
						{userCount.toLocaleString()}
					{/if}
				</div>
				<div class="stat-label">Total Users</div>
			</div>
			<ArrowRight class="stat-arrow" />
		</a>

		<a href="/api-keys" class="stat-card">
			<div class="stat-icon stat-icon-keys">
				<Key class="icon" />
			</div>
			<div class="stat-content">
				<div class="stat-value">API Keys</div>
				<div class="stat-label">Manage Access</div>
			</div>
			<ArrowRight class="stat-arrow" />
		</a>

		<a href="/monitoring" class="stat-card">
			<div class="stat-icon stat-icon-monitoring">
				<Activity class="icon" />
			</div>
			<div class="stat-content">
				<div class="stat-value">Metrics</div>
				<div class="stat-label">System Monitoring</div>
			</div>
			<ArrowRight class="stat-arrow" />
		</a>

		<div class="stat-card stat-card-disabled">
			<div class="stat-icon stat-icon-settings">
				<Settings class="icon" />
			</div>
			<div class="stat-content">
				<div class="stat-value">Settings</div>
				<div class="stat-label">Coming Soon</div>
			</div>
		</div>
	</div>

	<!-- Quick Actions -->
	<div class="section">
		<h2 class="section-title">Quick Actions</h2>
		<div class="actions-grid">
			<a href="/admin/users" class="action-card">
				<Users class="action-icon" />
				<div class="action-content">
					<h3 class="action-title">User Management</h3>
					<p class="action-description">View, edit, and manage user accounts</p>
				</div>
				<ArrowRight class="action-arrow" />
			</a>

			<a href="/api-keys" class="action-card">
				<Key class="action-icon" />
				<div class="action-content">
					<h3 class="action-title">API Keys</h3>
					<p class="action-description">Create and manage API keys for public access</p>
				</div>
				<ArrowRight class="action-arrow" />
			</a>

			<a href="/monitoring" class="action-card">
				<Activity class="action-icon" />
				<div class="action-content">
					<h3 class="action-title">System Monitoring</h3>
					<p class="action-description">View real-time metrics and system health</p>
				</div>
				<ArrowRight class="action-arrow" />
			</a>
		</div>
	</div>

	<!-- Admin Info -->
	<div class="section">
		<h2 class="section-title">Admin Information</h2>
		<div class="info-card">
			<div class="info-item">
				<Shield class="info-icon" />
				<div class="info-content">
					<div class="info-label">Admin Access</div>
					<div class="info-value">You have full administrative privileges</div>
				</div>
			</div>
			<div class="info-item">
				<Users class="info-icon" />
				<div class="info-content">
					<div class="info-label">User Management</div>
					<div class="info-value">Manage users, roles, and permissions</div>
				</div>
			</div>
			<div class="info-item">
				<Activity class="info-icon" />
				<div class="info-content">
					<div class="info-label">System Monitoring</div>
					<div class="info-value">Monitor system health and metrics</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.875rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 0.5rem;
	}

	.page-description {
		color: rgba(148, 163, 184, 0.9);
		font-size: 0.9rem;
	}

	.user-info {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.25rem;
	}

	.user-email {
		color: #f8fafc;
		font-weight: 500;
		font-size: 0.875rem;
	}

	.user-role {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.75rem;
		text-transform: uppercase;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1rem;
	}

	.stat-card {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1.5rem;
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		transition: all 200ms ease;
		cursor: pointer;
		text-decoration: none;
		color: inherit;
	}

	.stat-card:hover {
		background: rgba(15, 23, 42, 0.8);
		border-color: rgba(59, 130, 246, 0.6);
		transform: translateY(-2px);
	}

	.stat-card-disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.stat-card-disabled:hover {
		transform: none;
		border-color: rgba(71, 85, 105, 0.4);
	}

	.stat-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 3rem;
		height: 3rem;
		border-radius: 0.5rem;
	}

	.stat-icon-users {
		background: rgba(59, 130, 246, 0.2);
		color: #93c5fd;
	}

	.stat-icon-keys {
		background: rgba(34, 197, 94, 0.2);
		color: #86efac;
	}

	.stat-icon-monitoring {
		background: rgba(251, 191, 36, 0.2);
		color: #fde047;
	}

	.stat-icon-settings {
		background: rgba(71, 85, 105, 0.2);
		color: #cbd5e1;
	}

	.icon {
		width: 1.5rem;
		height: 1.5rem;
	}

	.stat-content {
		flex: 1;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 0.25rem;
	}

	.stat-label {
		font-size: 0.875rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.stat-arrow {
		width: 1.25rem;
		height: 1.25rem;
		color: rgba(148, 163, 184, 0.6);
	}

	.loading-dots {
		display: inline-block;
		width: 1.5rem;
		text-align: left;
	}

	.section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.section-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: #f8fafc;
	}

	.actions-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 1rem;
	}

	.action-card {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1.5rem;
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		transition: all 200ms ease;
		cursor: pointer;
		text-decoration: none;
		color: inherit;
	}

	.action-card:hover {
		background: rgba(15, 23, 42, 0.8);
		border-color: rgba(59, 130, 246, 0.6);
		transform: translateY(-2px);
	}

	.action-icon {
		width: 2rem;
		height: 2rem;
		color: #3b82f6;
		flex-shrink: 0;
	}

	.action-content {
		flex: 1;
	}

	.action-title {
		font-size: 1rem;
		font-weight: 600;
		color: #f8fafc;
		margin-bottom: 0.25rem;
	}

	.action-description {
		font-size: 0.875rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.action-arrow {
		width: 1.25rem;
		height: 1.25rem;
		color: rgba(148, 163, 184, 0.6);
		flex-shrink: 0;
	}

	.info-card {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1.5rem;
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
	}

	.info-item {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
	}

	.info-icon {
		width: 1.5rem;
		height: 1.5rem;
		color: #3b82f6;
		flex-shrink: 0;
		margin-top: 0.125rem;
	}

	.info-content {
		flex: 1;
	}

	.info-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #f8fafc;
		margin-bottom: 0.25rem;
	}

	.info-value {
		font-size: 0.875rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.alert {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 0.5rem;
		border: 1px solid;
	}

	.alert-error {
		background: rgba(239, 68, 68, 0.1);
		border-color: rgba(239, 68, 68, 0.3);
		color: #fca5a5;
	}
</style>

