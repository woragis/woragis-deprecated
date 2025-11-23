<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Trash2, Copy, Check, Calendar, Key, AlertCircle } from 'lucide-svelte';
	import {
		listAPIKeys,
		createAPIKey,
		updateAPIKey,
		deleteAPIKey,
		type APIKey,
		type APIKeyWithToken
	} from '$lib/api/apikeys';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let apiKeys: APIKey[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showTokenModal = $state(false);
	let newToken: string | null = $state(null);
	let copiedKeyId: string | null = $state(null);

	// Form state
	let formName = $state('');
	let formExpiresAt = $state('');
	let editingId: string | null = $state(null);
	let editingName = $state('');

	onMount(async () => {
		await fetchAPIKeys();
	});

	async function fetchAPIKeys() {
		loading = true;
		error = null;
		try {
			apiKeys = await listAPIKeys();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch API keys';
			toastError(error);
			console.error('Error fetching API keys:', err);
		} finally {
			loading = false;
		}
	}

	async function handleCreate() {
		if (!formName.trim()) {
			toastError('Name is required');
			return;
		}

		try {
			const payload: { name: string; expiresAt?: string } = {
				name: formName.trim()
			};

			if (formExpiresAt) {
				payload.expiresAt = new Date(formExpiresAt).toISOString();
			}

			const result: APIKeyWithToken = await createAPIKey(payload);
			newToken = result.token;
			showCreateModal = false;
			showTokenModal = true;
			formName = '';
			formExpiresAt = '';
			toastSuccess('API key created successfully');
			await fetchAPIKeys();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to create API key';
			toastError(errorMsg);
			console.error('Error creating API key:', err);
		}
	}

	async function handleUpdate(id: string) {
		if (!editingName.trim()) {
			toastError('Name is required');
			return;
		}

		try {
			await updateAPIKey(id, { name: editingName.trim() });
			editingId = null;
			editingName = '';
			toastSuccess('API key updated successfully');
			await fetchAPIKeys();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to update API key';
			toastError(errorMsg);
			console.error('Error updating API key:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this API key? This action cannot be undone.')) {
			return;
		}

		try {
			await deleteAPIKey(id);
			toastSuccess('API key deleted successfully');
			await fetchAPIKeys();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to delete API key';
			toastError(errorMsg);
			console.error('Error deleting API key:', err);
		}
	}

	function startEdit(apiKey: APIKey) {
		editingId = apiKey.id;
		editingName = apiKey.name;
	}

	function cancelEdit() {
		editingId = null;
		editingName = '';
	}

	async function copyToClipboard(text: string, keyId: string) {
		try {
			await navigator.clipboard.writeText(text);
			copiedKeyId = keyId;
			toastSuccess('Copied to clipboard');
			setTimeout(() => {
				copiedKeyId = null;
			}, 2000);
		} catch (err) {
			toastError('Failed to copy to clipboard');
			console.error('Failed to copy:', err);
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	function isExpired(expiresAt?: string): boolean {
		if (!expiresAt) return false;
		return new Date(expiresAt) < new Date();
	}
</script>

<div class="page-container">
	<!-- Header -->
	<div class="page-header">
		<div>
			<h1 class="page-title">API Keys</h1>
			<p class="page-description">Manage API keys for public read-only access to your data</p>
		</div>
		<button
			type="button"
			class="btn btn-primary"
			onclick={() => {
				showCreateModal = true;
				formName = '';
				formExpiresAt = '';
			}}
		>
			<Plus class="icon" />
			Create API Key
		</button>
	</div>

	<!-- Error State -->
	{#if error}
		<div class="alert alert-error">
			<AlertCircle class="icon" />
			<p>{error}</p>
		</div>
	{/if}

	<!-- Loading State -->
	{#if loading}
		<div class="loading-container">
			<div class="spinner"></div>
		</div>
	{:else if apiKeys.length === 0}
		<!-- Empty State -->
		<div class="empty-state">
			<Key class="empty-icon" />
			<p class="empty-title">No API keys found</p>
			<p class="empty-description">Create your first API key to enable public read access</p>
			<button type="button" class="btn btn-primary" onclick={() => (showCreateModal = true)}>
				Create API Key
			</button>
		</div>
	{:else}
		<!-- API Keys Table -->
		<div class="table-container">
			<table class="table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Prefix</th>
						<th>Created</th>
						<th>Last Used</th>
						<th>Expires</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each apiKeys as apiKey}
						<tr>
							<td>
								{#if editingId === apiKey.id}
									<input
										type="text"
										bind:value={editingName}
										class="input"
										onkeydown={(e) => {
											if (e.key === 'Enter') {
												handleUpdate(apiKey.id);
											} else if (e.key === 'Escape') {
												cancelEdit();
											}
										}}
									/>
								{:else}
									<span class="font-medium">{apiKey.name}</span>
								{/if}
							</td>
							<td>
								<code class="code-badge">{apiKey.prefix}...</code>
							</td>
							<td class="text-muted">
								{formatDate(apiKey.createdAt)}
							</td>
							<td class="text-muted">
								{apiKey.lastUsedAt ? formatDate(apiKey.lastUsedAt) : 'Never'}
							</td>
							<td>
								{#if apiKey.expiresAt}
									<span class={isExpired(apiKey.expiresAt) ? 'text-error' : 'text-muted'}>
										{formatDate(apiKey.expiresAt)}
									</span>
								{:else}
									<span class="text-muted">Never</span>
								{/if}
							</td>
							<td class="text-right">
								<div class="actions">
									{#if editingId === apiKey.id}
										<button
											type="button"
											class="btn btn-sm btn-success"
											onclick={() => handleUpdate(apiKey.id)}
										>
											Save
										</button>
										<button
											type="button"
											class="btn btn-sm btn-secondary"
											onclick={cancelEdit}
										>
											Cancel
										</button>
									{:else}
										<button
											type="button"
											class="btn btn-sm btn-primary"
											onclick={() => startEdit(apiKey)}
										>
											Edit
										</button>
										<button
											type="button"
											class="btn btn-sm btn-danger"
											onclick={() => handleDelete(apiKey.id)}
										>
											<Trash2 class="icon-sm" />
										</button>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="modal-overlay" onclick={() => (showCreateModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Create API Key</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Name</label>
					<input
						type="text"
						bind:value={formName}
						placeholder="e.g., Landing Page Key"
						class="input"
					/>
				</div>
				<div class="form-group">
					<label class="form-label">Expires At (Optional)</label>
					<input type="datetime-local" bind:value={formExpiresAt} class="input" />
				</div>
				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handleCreate}>
						Create
					</button>
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => (showCreateModal = false)}
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Token Display Modal -->
{#if showTokenModal && newToken}
	<div
		class="modal-overlay"
		onclick={() => {
			showTokenModal = false;
			newToken = null;
		}}
	>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">API Key Created</h2>
			<div class="modal-content">
				<div class="alert alert-warning">
					<AlertCircle class="icon" />
					<p>Copy this key now. You won't be able to see it again!</p>
				</div>
				<div class="token-display">
					<code class="token-code">{newToken}</code>
					<button
						type="button"
						class="btn-icon"
						onclick={() => copyToClipboard(newToken!, 'new')}
						title="Copy to clipboard"
					>
						{#if copiedKeyId === 'new'}
							<Check class="icon text-success" />
						{:else}
							<Copy class="icon" />
						{/if}
					</button>
				</div>
				<button
					type="button"
					class="btn btn-primary btn-full"
					onclick={() => {
						showTokenModal = false;
						newToken = null;
					}}
				>
					I've copied the key
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
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

	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		border: 1px solid;
		transition: all 120ms ease;
		cursor: pointer;
	}

	.btn-primary {
		background: rgba(59, 130, 246, 0.15);
		border-color: rgba(59, 130, 246, 0.4);
		color: #93c5fd;
	}

	.btn-primary:hover {
		background: rgba(59, 130, 246, 0.25);
		border-color: rgba(59, 130, 246, 0.6);
	}

	.btn-sm {
		padding: 0.375rem 0.75rem;
		font-size: 0.8rem;
	}

	.btn-success {
		background: rgba(34, 197, 94, 0.15);
		border-color: rgba(34, 197, 94, 0.4);
		color: #86efac;
	}

	.btn-success:hover {
		background: rgba(34, 197, 94, 0.25);
		border-color: rgba(34, 197, 94, 0.6);
	}

	.btn-secondary {
		background: rgba(71, 85, 105, 0.15);
		border-color: rgba(71, 85, 105, 0.4);
		color: #cbd5e1;
	}

	.btn-secondary:hover {
		background: rgba(71, 85, 105, 0.25);
		border-color: rgba(71, 85, 105, 0.6);
	}

	.btn-danger {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.4);
		color: #fca5a5;
	}

	.btn-danger:hover {
		background: rgba(239, 68, 68, 0.25);
		border-color: rgba(239, 68, 68, 0.6);
	}

	.btn-full {
		width: 100%;
		justify-content: center;
	}

	.btn-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 0.5rem;
		border-radius: 0.375rem;
		background: rgba(71, 85, 105, 0.15);
		border: 1px solid rgba(71, 85, 105, 0.4);
		color: #cbd5e1;
		cursor: pointer;
		transition: all 120ms ease;
	}

	.btn-icon:hover {
		background: rgba(71, 85, 105, 0.25);
		border-color: rgba(71, 85, 105, 0.6);
	}

	.icon {
		width: 1rem;
		height: 1rem;
	}

	.icon-sm {
		width: 0.875rem;
		height: 0.875rem;
	}

	.loading-container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem 0;
	}

	.spinner {
		width: 3rem;
		height: 3rem;
		border: 2px solid rgba(71, 85, 105, 0.3);
		border-top-color: #3b82f6;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
	}

	.empty-icon {
		width: 4rem;
		height: 4rem;
		color: rgba(71, 85, 105, 0.6);
		margin: 0 auto 1rem;
	}

	.empty-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: rgba(203, 213, 225, 0.9);
		margin-bottom: 0.5rem;
	}

	.empty-description {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
		margin-bottom: 1.5rem;
	}

	.table-container {
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		overflow: hidden;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
	}

	.table thead {
		background: rgba(15, 23, 42, 0.8);
	}

	.table th {
		padding: 1rem 1.5rem;
		text-align: left;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: rgba(148, 163, 184, 0.9);
	}

	.table td {
		padding: 1rem 1.5rem;
		border-top: 1px solid rgba(71, 85, 105, 0.3);
	}

	.table tbody tr:hover {
		background: rgba(51, 65, 85, 0.2);
	}

	.text-right {
		text-align: right;
	}

	.text-muted {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.text-error {
		color: #fca5a5;
	}

	.text-success {
		color: #86efac;
	}

	.code-badge {
		padding: 0.25rem 0.5rem;
		background: rgba(15, 23, 42, 0.8);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.375rem;
		font-size: 0.8rem;
		color: #cbd5e1;
	}

	.actions {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.5rem;
	}

	.input {
		width: 100%;
		padding: 0.5rem 0.75rem;
		background: rgba(15, 23, 42, 0.8);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.5rem;
		color: #f8fafc;
		font-size: 0.875rem;
	}

	.input:focus {
		outline: none;
		border-color: rgba(59, 130, 246, 0.6);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
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

	.alert-warning {
		background: rgba(251, 191, 36, 0.1);
		border-color: rgba(251, 191, 36, 0.3);
		color: #fde047;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(2, 6, 23, 0.7);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal {
		background: rgba(15, 23, 42, 0.95);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		padding: 1.5rem;
		width: 100%;
		max-width: 28rem;
		box-shadow: 0 20px 45px rgba(2, 6, 23, 0.6);
	}

	.modal-large {
		max-width: 32rem;
	}

	.modal-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 1rem;
	}

	.modal-content {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: rgba(203, 213, 225, 0.9);
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}

	.token-display {
		position: relative;
		background: rgba(2, 6, 23, 0.8);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.5rem;
		padding: 1rem;
	}

	.token-code {
		display: block;
		font-size: 0.8rem;
		word-break: break-all;
		color: rgba(203, 213, 225, 0.9);
		padding-right: 2.5rem;
	}

	.token-display .btn-icon {
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
	}

	.font-medium {
		font-weight: 500;
	}
</style>

