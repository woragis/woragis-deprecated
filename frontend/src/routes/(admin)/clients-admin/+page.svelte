<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listClients,
		createClient,
		updateClient,
		toggleArchiveClient,
		deleteClient,
		type Client,
		type CreateClientInput,
		type UpdateClientInput
	} from '$lib/api/clients';

	let clients: Client[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingClient: Client | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formName = $state('');
	let formEmail = $state('');
	let formPhoneNumber = $state('');
	let formCompany = $state('');
	let formNotes = $state('');

	onMount(async () => {
		await fetchClients();
	});

	async function fetchClients() {
		loading = true;
		error = null;
		try {
			clients = await listClients();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load clients';
			console.error('Error fetching clients:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(client: Client) {
		editingClient = client;
		formName = client.name;
		formEmail = client.email || '';
		formPhoneNumber = client.phoneNumber;
		formCompany = client.company || '';
		formNotes = client.notes || '';
		showEditModal = true;
	}

	function resetForm() {
		formName = '';
		formEmail = '';
		formPhoneNumber = '';
		formCompany = '';
		formNotes = '';
		editingClient = null;
	}

	async function handleCreate() {
		if (!formName.trim() || !formPhoneNumber.trim()) {
			alert('Name and phone number are required');
			return;
		}

		try {
			const input: CreateClientInput = {
				name: formName.trim(),
				email: formEmail.trim() || undefined,
				phoneNumber: formPhoneNumber.trim(),
				company: formCompany.trim() || undefined,
				notes: formNotes.trim() || undefined
			};

			await createClient(input);
			showCreateModal = false;
			resetForm();
			await fetchClients();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create client');
			console.error('Error creating client:', err);
		}
	}

	async function handleUpdate() {
		if (!editingClient) {
			alert('No client selected');
			return;
		}

		try {
			const input: UpdateClientInput = {
				name: formName.trim(),
				email: formEmail.trim() || undefined,
				phoneNumber: formPhoneNumber.trim(),
				company: formCompany.trim() || undefined,
				notes: formNotes.trim() || undefined
			};

			await updateClient(editingClient.id, input);
			showEditModal = false;
			resetForm();
			await fetchClients();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update client');
			console.error('Error updating client:', err);
		}
	}

	async function handleToggleArchive(id: string) {
		try {
			await toggleArchiveClient(id);
			await fetchClients();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to toggle archive');
			console.error('Error toggling archive:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this client?')) return;

		try {
			await deleteClient(id);
			await fetchClients();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete client');
			console.error('Error deleting client:', err);
		}
	}

	function filteredClients() {
		if (!searchQuery.trim()) return clients;
		const query = searchQuery.toLowerCase();
		return clients.filter(
			(c) =>
				c.name.toLowerCase().includes(query) ||
				c.email?.toLowerCase().includes(query) ||
				c.phoneNumber.toLowerCase().includes(query) ||
				c.company?.toLowerCase().includes(query)
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
			<h1>Clients Management</h1>
			<p>Manage clients and contacts</p>
		</div>
		<button onclick={openCreateModal}>Create Client</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search clients..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredClients().length === 0}
		<div class="empty">No clients found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Email</th>
					<th>Phone</th>
					<th>Company</th>
					<th>Archived</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredClients() as client}
					<tr>
						<td><strong>{client.name}</strong></td>
						<td>{client.email || '—'}</td>
						<td>{client.phoneNumber}</td>
						<td>{client.company || '—'}</td>
						<td>{client.isArchived ? 'Yes' : 'No'}</td>
						<td>{formatDate(client.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(client)}>Edit</button>
							<button
								onclick={() => handleToggleArchive(client.id)}
								class="archive-btn"
							>
								{client.isArchived ? 'Unarchive' : 'Archive'}
							</button>
							<button onclick={() => handleDelete(client.id)} class="delete-btn">Delete</button>
						</td>
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
			<h2>Create Client</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Phone Number *</label>
					<input type="tel" bind:value={formPhoneNumber} />
				</div>
				<div class="form-group">
					<label>Email</label>
					<input type="email" bind:value={formEmail} />
				</div>
				<div class="form-group">
					<label>Company</label>
					<input type="text" bind:value={formCompany} />
				</div>
				<div class="form-group">
					<label>Notes</label>
					<textarea bind:value={formNotes} rows="4"></textarea>
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>Create</button>
					<button onclick={() => (showCreateModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingClient}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Client</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Phone Number *</label>
					<input type="tel" bind:value={formPhoneNumber} />
				</div>
				<div class="form-group">
					<label>Email</label>
					<input type="email" bind:value={formEmail} />
				</div>
				<div class="form-group">
					<label>Company</label>
					<input type="text" bind:value={formCompany} />
				</div>
				<div class="form-group">
					<label>Notes</label>
					<textarea bind:value={formNotes} rows="4"></textarea>
				</div>
				<div class="form-actions">
					<button onclick={handleUpdate}>Update</button>
					<button onclick={() => (showEditModal = false)}>Cancel</button>
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

	.table button {
		padding: 0.25rem 0.75rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		margin-right: 0.5rem;
	}

	.table button:hover {
		background: #218838;
	}

	.archive-btn {
		background: #17a2b8 !important;
	}

	.archive-btn:hover {
		background: #138496 !important;
	}

	.delete-btn {
		background: #dc3545 !important;
	}

	.delete-btn:hover {
		background: #c82333 !important;
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

