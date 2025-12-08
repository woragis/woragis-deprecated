<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listJobWebsites,
		createJobWebsite,
		updateJobWebsite,
		deleteJobWebsite,
		resetJobWebsiteCounter,
		type JobWebsite,
		type CreateJobWebsiteInput,
		type UpdateJobWebsiteInput
	} from '$lib/api/jobwebsites';

	let websites: JobWebsite[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingWebsite: JobWebsite | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formName = $state('');
	let formDisplayName = $state('');
	let formDailyLimit = $state<number | ''>(50);
	let formBaseUrl = $state('');
	let formLoginUrl = $state('');
	let formEnabled = $state(true);

	onMount(async () => {
		await fetchWebsites();
	});

	async function fetchWebsites() {
		loading = true;
		error = null;
		try {
			websites = await listJobWebsites();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load job websites';
			console.error('Error fetching job websites:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(website: JobWebsite) {
		editingWebsite = website;
		formDisplayName = website.displayName;
		formDailyLimit = website.dailyLimit;
		formBaseUrl = website.baseUrl || '';
		formLoginUrl = website.loginUrl || '';
		formEnabled = website.enabled;
		showEditModal = true;
	}

	function resetForm() {
		formName = '';
		formDisplayName = '';
		formDailyLimit = 50;
		formBaseUrl = '';
		formLoginUrl = '';
		formEnabled = true;
		editingWebsite = null;
	}

	async function handleCreate() {
		if (!formName.trim() || !formDisplayName.trim()) {
			alert('Name and display name are required');
			return;
		}

		try {
			const input: CreateJobWebsiteInput = {
				name: formName.trim(),
				displayName: formDisplayName.trim(),
				dailyLimit: formDailyLimit ? Number(formDailyLimit) : 50,
				baseUrl: formBaseUrl.trim() || undefined,
				loginUrl: formLoginUrl.trim() || undefined,
				enabled: formEnabled
			};

			await createJobWebsite(input);
			showCreateModal = false;
			resetForm();
			await fetchWebsites();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create job website');
			console.error('Error creating job website:', err);
		}
	}

	async function handleUpdate() {
		if (!editingWebsite) {
			alert('No website selected');
			return;
		}

		try {
			const input: UpdateJobWebsiteInput = {
				displayName: formDisplayName.trim(),
				dailyLimit: formDailyLimit ? Number(formDailyLimit) : undefined,
				baseUrl: formBaseUrl.trim() || undefined,
				loginUrl: formLoginUrl.trim() || undefined,
				enabled: formEnabled
			};

			await updateJobWebsite(editingWebsite.id, input);
			showEditModal = false;
			resetForm();
			await fetchWebsites();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update job website');
			console.error('Error updating job website:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this job website?')) return;

		try {
			await deleteJobWebsite(id);
			await fetchWebsites();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete job website');
			console.error('Error deleting job website:', err);
		}
	}

	async function handleResetCounter(id: string) {
		try {
			await resetJobWebsiteCounter(id);
			await fetchWebsites();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to reset counter');
			console.error('Error resetting counter:', err);
		}
	}

	function filteredWebsites() {
		if (!searchQuery.trim()) return websites;
		const query = searchQuery.toLowerCase();
		return websites.filter(
			(w) =>
				w.name.toLowerCase().includes(query) ||
				w.displayName.toLowerCase().includes(query)
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
			<h1>Job Websites Management</h1>
			<p>Manage job website configurations</p>
		</div>
		<button onclick={openCreateModal}>Create Website</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search websites..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredWebsites().length === 0}
		<div class="empty">No job websites found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Display Name</th>
					<th>Daily Limit</th>
					<th>Current Count</th>
					<th>Enabled</th>
					<th>Last Reset</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredWebsites() as website}
					<tr>
						<td><strong>{website.name}</strong></td>
						<td>{website.displayName}</td>
						<td>{website.dailyLimit}</td>
						<td>
							<span class="count {website.currentCount >= website.dailyLimit ? 'limit-reached' : ''}">
								{website.currentCount} / {website.dailyLimit}
							</span>
						</td>
						<td>{website.enabled ? 'Yes' : 'No'}</td>
						<td>{formatDate(website.lastReset)}</td>
						<td>
							<button onclick={() => openEditModal(website)}>Edit</button>
							<button onclick={() => handleResetCounter(website.id)} class="reset-btn">
								Reset
							</button>
							<button onclick={() => handleDelete(website.id)} class="delete-btn">Delete</button>
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
			<h2>Create Job Website</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} placeholder="linkedin, glassdoor, etc." />
				</div>
				<div class="form-group">
					<label>Display Name *</label>
					<input type="text" bind:value={formDisplayName} />
				</div>
				<div class="form-group">
					<label>Daily Limit</label>
					<input type="number" bind:value={formDailyLimit} />
				</div>
				<div class="form-group">
					<label>Base URL</label>
					<input type="url" bind:value={formBaseUrl} />
				</div>
				<div class="form-group">
					<label>Login URL</label>
					<input type="url" bind:value={formLoginUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formEnabled} />
						Enabled
					</label>
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
{#if showEditModal && editingWebsite}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Job Website</h2>
			<div class="form">
				<div class="form-group">
					<label>Display Name *</label>
					<input type="text" bind:value={formDisplayName} />
				</div>
				<div class="form-group">
					<label>Daily Limit</label>
					<input type="number" bind:value={formDailyLimit} />
				</div>
				<div class="form-group">
					<label>Base URL</label>
					<input type="url" bind:value={formBaseUrl} />
				</div>
				<div class="form-group">
					<label>Login URL</label>
					<input type="url" bind:value={formLoginUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formEnabled} />
						Enabled
					</label>
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

	.count {
		font-weight: 500;
	}

	.count.limit-reached {
		color: #dc3545;
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

	.reset-btn {
		background: #17a2b8 !important;
	}

	.reset-btn:hover {
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
	.form-group textarea,
	.form-group select {
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

