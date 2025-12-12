<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Trash2, Copy, Check, Eye, EyeOff, Calendar, Key } from 'lucide-svelte';
	import {
		listAPIKeys,
		createAPIKey,
		updateAPIKey,
		deleteAPIKey,
		type APIKey,
		type APIKeyWithToken
	} from '$lib/api/apikeys';

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
			console.error('Error fetching API keys:', err);
		} finally {
			loading = false;
		}
	}

	async function handleCreate() {
		if (!formName.trim()) {
			error = 'Name is required';
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
			await fetchAPIKeys();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create API key';
			console.error('Error creating API key:', err);
		}
	}

	async function handleUpdate(id: string) {
		if (!editingName.trim()) {
			error = 'Name is required';
			return;
		}

		try {
			await updateAPIKey(id, { name: editingName.trim() });
			editingId = null;
			editingName = '';
			await fetchAPIKeys();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update API key';
			console.error('Error updating API key:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this API key? This action cannot be undone.')) {
			return;
		}

		try {
			await deleteAPIKey(id);
			await fetchAPIKeys();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete API key';
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
			setTimeout(() => {
				copiedKeyId = null;
			}, 2000);
		} catch (err) {
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

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white">
	<div class="container mx-auto px-4 py-8">
		<!-- Header -->
		<div class="mb-8 flex items-center justify-between">
			<div>
				<h1 class="text-4xl font-bold mb-2">API Keys</h1>
				<p class="text-gray-400">Manage API keys for public read-only access to your data</p>
			</div>
			<button
				onclick={() => {
					showCreateModal = true;
					formName = '';
					formExpiresAt = '';
				}}
				class="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
			>
				<Plus class="w-5 h-5" />
				Create API Key
			</button>
		</div>

		<!-- Error State -->
		{#if error}
			<div class="bg-red-900/50 border border-red-700 rounded-lg p-4 mb-6">
				<p class="text-red-200">{error}</p>
			</div>
		{/if}

		<!-- Loading State -->
		{#if loading}
			<div class="flex justify-center items-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
			</div>
		{:else if apiKeys.length === 0}
			<!-- Empty State -->
			<div class="text-center py-20">
				<Key class="w-16 h-16 mx-auto text-gray-600 mb-4" />
				<p class="text-gray-400 text-lg mb-2">No API keys found</p>
				<p class="text-gray-500 text-sm mb-6">Create your first API key to enable public read access</p>
				<button
					onclick={() => (showCreateModal = true)}
					class="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
				>
					Create API Key
				</button>
			</div>
		{:else}
			<!-- API Keys Table -->
			<div class="bg-gray-800/50 backdrop-blur-sm rounded-xl border border-gray-700 overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-gray-700/50">
							<tr>
								<th class="px-6 py-4 text-left text-sm font-semibold text-gray-300">Name</th>
								<th class="px-6 py-4 text-left text-sm font-semibold text-gray-300">Prefix</th>
								<th class="px-6 py-4 text-left text-sm font-semibold text-gray-300">Created</th>
								<th class="px-6 py-4 text-left text-sm font-semibold text-gray-300">Last Used</th>
								<th class="px-6 py-4 text-left text-sm font-semibold text-gray-300">Expires</th>
								<th class="px-6 py-4 text-right text-sm font-semibold text-gray-300">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-700">
							{#each apiKeys as apiKey}
								<tr class="hover:bg-gray-700/30 transition-colors">
									<td class="px-6 py-4">
										{#if editingId === apiKey.id}
											<input
												type="text"
												bind:value={editingName}
												class="w-full px-3 py-1 bg-gray-700 border border-gray-600 rounded text-white"
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
									<td class="px-6 py-4">
										<code class="px-2 py-1 bg-gray-700 rounded text-sm">{apiKey.prefix}...</code>
									</td>
									<td class="px-6 py-4 text-sm text-gray-400">
										{formatDate(apiKey.createdAt)}
									</td>
									<td class="px-6 py-4 text-sm text-gray-400">
										{apiKey.lastUsedAt ? formatDate(apiKey.lastUsedAt) : 'Never'}
									</td>
									<td class="px-6 py-4 text-sm">
										{#if apiKey.expiresAt}
											<span class={isExpired(apiKey.expiresAt) ? 'text-red-400' : 'text-gray-400'}>
												{formatDate(apiKey.expiresAt)}
											</span>
										{:else}
											<span class="text-gray-500">Never</span>
										{/if}
									</td>
									<td class="px-6 py-4 text-right">
										<div class="flex items-center justify-end gap-2">
											{#if editingId === apiKey.id}
												<button
													onclick={() => handleUpdate(apiKey.id)}
													class="px-3 py-1 bg-green-600 hover:bg-green-700 rounded text-sm transition-colors"
												>
													Save
												</button>
												<button
													onclick={cancelEdit}
													class="px-3 py-1 bg-gray-600 hover:bg-gray-700 rounded text-sm transition-colors"
												>
													Cancel
												</button>
											{:else}
												<button
													onclick={() => startEdit(apiKey)}
													class="px-3 py-1 bg-blue-600 hover:bg-blue-700 rounded text-sm transition-colors"
												>
													Edit
												</button>
												<button
													onclick={() => handleDelete(apiKey.id)}
													class="px-3 py-1 bg-red-600 hover:bg-red-700 rounded text-sm transition-colors"
												>
													<Trash2 class="w-4 h-4" />
												</button>
											{/if}
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}
	</div>
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
		onclick={() => (showCreateModal = false)}
	>
		<div
			class="bg-gray-800 rounded-xl p-6 w-full max-w-md border border-gray-700"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold mb-4">Create API Key</h2>
			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-300 mb-2">Name</label>
					<input
						type="text"
						bind:value={formName}
						placeholder="e.g., Landing Page Key"
						class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
					/>
				</div>
				<div>
					<label class="block text-sm font-medium text-gray-300 mb-2">Expires At (Optional)</label>
					<input
						type="datetime-local"
						bind:value={formExpiresAt}
						class="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
					/>
				</div>
				<div class="flex gap-3 pt-4">
					<button
						onclick={handleCreate}
						class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
					>
						Create
					</button>
					<button
						onclick={() => (showCreateModal = false)}
						class="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
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
		class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50"
		onclick={() => {
			showTokenModal = false;
			newToken = null;
		}}
	>
		<div
			class="bg-gray-800 rounded-xl p-6 w-full max-w-md border border-gray-700"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-2xl font-bold mb-2">API Key Created</h2>
			<p class="text-gray-400 text-sm mb-4">
				Copy this key now. You won't be able to see it again!
			</p>
			<div class="bg-gray-900 rounded-lg p-4 mb-4 relative">
				<code class="text-sm break-all text-gray-300">{newToken}</code>
				<button
					onclick={() => copyToClipboard(newToken!, 'new')}
					class="absolute top-2 right-2 p-2 hover:bg-gray-700 rounded transition-colors"
					title="Copy to clipboard"
				>
					{#if copiedKeyId === 'new'}
						<Check class="w-4 h-4 text-green-400" />
					{:else}
						<Copy class="w-4 h-4 text-gray-400" />
					{/if}
				</button>
			</div>
			<button
				onclick={() => {
					showTokenModal = false;
					newToken = null;
				}}
				class="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
			>
				I've copied the key
			</button>
		</div>
	</div>
{/if}

