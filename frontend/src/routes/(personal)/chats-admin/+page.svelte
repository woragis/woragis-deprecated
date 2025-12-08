<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listConversations,
		createConversation,
		archiveConversations,
		deleteConversations,
		restoreConversations,
		type Conversation,
		type CreateConversationInput
	} from '$lib/api/chats';

	let conversations: Conversation[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let searchQuery = $state('');
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formTitle = $state('');
	let formDescription = $state('');
	let formIdeaId = $state('');
	let formProjectId = $state('');

	onMount(async () => {
		await fetchConversations();
	});

	async function fetchConversations() {
		loading = true;
		error = null;
		try {
			conversations = await listConversations();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load conversations';
			console.error('Error fetching conversations:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function resetForm() {
		formTitle = '';
		formDescription = '';
		formIdeaId = '';
		formProjectId = '';
	}

	async function handleCreate() {
		if (!formTitle.trim()) {
			alert('Title is required');
			return;
		}

		try {
			const input: CreateConversationInput = {
				title: formTitle.trim(),
				description: formDescription.trim() || undefined,
				ideaId: formIdeaId.trim() || undefined,
				projectId: formProjectId.trim() || undefined
			};

			await createConversation(input);
			showCreateModal = false;
			resetForm();
			await fetchConversations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create conversation');
			console.error('Error creating conversation:', err);
		}
	}

	function toggleSelect(id: string) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = new Set(selectedIds);
	}

	async function handleBulkArchive() {
		if (selectedIds.size === 0) return;
		try {
			await archiveConversations(Array.from(selectedIds));
			selectedIds.clear();
			await fetchConversations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to archive conversations');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} conversation(s)?`)) return;
		try {
			await deleteConversations(Array.from(selectedIds));
			selectedIds.clear();
			await fetchConversations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete conversations');
		}
	}

	async function handleBulkRestore() {
		if (selectedIds.size === 0) return;
		try {
			await restoreConversations(Array.from(selectedIds));
			selectedIds.clear();
			await fetchConversations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to restore conversations');
		}
	}

	function filteredConversations() {
		if (!searchQuery.trim()) return conversations;
		const query = searchQuery.toLowerCase();
		return conversations.filter(
			(c) =>
				c.title.toLowerCase().includes(query) ||
				c.description?.toLowerCase().includes(query)
		);
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="header">
		<div>
			<h1>Chats Management</h1>
			<p>Manage conversations</p>
		</div>
		<button onclick={openCreateModal}>Create Conversation</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search conversations..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if selectedIds.size > 0}
		<div class="bulk-actions">
			<span>{selectedIds.size} selected</span>
			<button onclick={handleBulkArchive}>Archive</button>
			<button onclick={handleBulkDelete} class="delete-btn">Delete</button>
			<button onclick={handleBulkRestore}>Restore</button>
			<button onclick={() => (selectedIds.clear())}>Clear</button>
		</div>
	{/if}

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredConversations().length === 0}
		<div class="empty">No conversations found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>
						<input
							type="checkbox"
							checked={selectedIds.size === filteredConversations().length && filteredConversations().length > 0}
							onchange={() => {
								if (selectedIds.size === filteredConversations().length) {
									selectedIds.clear();
								} else {
									selectedIds = new Set(filteredConversations().map((c) => c.id));
								}
							}}
						/>
					</th>
					<th>Title</th>
					<th>Description</th>
					<th>Idea ID</th>
					<th>Project ID</th>
					<th>Status</th>
					<th>Created</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredConversations() as conversation}
					<tr>
						<td>
							<input
								type="checkbox"
								checked={selectedIds.has(conversation.id)}
								onchange={() => toggleSelect(conversation.id)}
							/>
						</td>
						<td><strong>{conversation.title}</strong></td>
						<td>{conversation.description || '—'}</td>
						<td>{conversation.ideaId ? conversation.ideaId.substring(0, 8) + '...' : '—'}</td>
						<td>{conversation.projectId ? conversation.projectId.substring(0, 8) + '...' : '—'}</td>
						<td>
							{#if conversation.deletedAt}
								<span class="status status-deleted">Deleted</span>
							{:else if conversation.archivedAt}
								<span class="status status-archived">Archived</span>
							{:else}
								<span class="status status-active">Active</span>
							{/if}
						</td>
						<td>{formatDate(conversation.createdAt)}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="modal-overlay" onclick={() => (showCreateModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Create Conversation</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Idea ID</label>
					<input type="text" bind:value={formIdeaId} />
				</div>
				<div class="form-group">
					<label>Project ID</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>Create</button>
					<button onclick={() => (showCreateModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		padding: 1rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.header h1 {
		margin: 0 0 0.25rem 0;
		font-size: 1.5rem;
	}

	.header p {
		margin: 0;
		color: #666;
		font-size: 0.9rem;
	}

	.header button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.header button:hover {
		background: #0056b3;
	}

	.search-bar {
		margin-bottom: 1rem;
	}

	.search-input {
		width: 100%;
		max-width: 400px;
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	.bulk-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem;
		background: #f5f5f5;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.bulk-actions button {
		padding: 0.375rem 0.75rem;
		background: #6c757d;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.bulk-actions button:hover {
		background: #5a6268;
	}

	.bulk-actions .delete-btn {
		background: #dc3545 !important;
	}

	.bulk-actions .delete-btn:hover {
		background: #c82333 !important;
	}

	.error {
		padding: 0.75rem;
		background: #fee;
		color: #c33;
		border: 1px solid #fcc;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.loading,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #666;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
		background: white;
	}

	.table th,
	.table td {
		padding: 0.75rem;
		text-align: left;
		border-bottom: 1px solid #ddd;
	}

	.table th {
		background: #f5f5f5;
		font-weight: 600;
	}

	.table tbody tr:hover {
		background: #f9f9f9;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-active {
		background: #d4edda;
		color: #155724;
	}

	.status-archived {
		background: #fff3cd;
		color: #856404;
	}

	.status-deleted {
		background: #f8d7da;
		color: #721c24;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.form-group label {
		font-weight: 500;
		font-size: 0.875rem;
	}

	.form-group input,
	.form-group textarea {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
	}

	.form-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.form-actions button {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.form-actions button:first-child {
		background: #007bff;
		color: white;
	}

	.form-actions button:first-child:hover {
		background: #0056b3;
	}

	.form-actions button:last-child {
		background: #6c757d;
		color: white;
	}

	.form-actions button:last-child:hover {
		background: #5a6268;
	}
</style>

