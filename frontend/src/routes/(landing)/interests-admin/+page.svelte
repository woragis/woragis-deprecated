<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listInterests,
		createInterest,
		updateInterest,
		deleteInterest,
		type Interest,
		type CreateInterestInput,
		type UpdateInterestInput
	} from '$lib/api/interests';

	let interests: Interest[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingInterest: Interest | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formTitle = $state('');
	let formDescription = $state('');
	let formIcon = $state('');
	let formColor = $state('');
	let formBgGradient = $state('');
	let formBorderColor = $state('');
	let formHoverBorderColor = $state('');
	let formShadowColor = $state('');
	let formFullWidth = $state(false);
	let formFeatured = $state(false);

	onMount(async () => {
		await fetchInterests();
	});

	async function fetchInterests() {
		loading = true;
		error = null;
		try {
			interests = await listInterests();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load interests';
			console.error('Error fetching interests:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(interest: Interest) {
		editingInterest = interest;
		formTitle = interest.title;
		formDescription = interest.description;
		formIcon = interest.icon || '';
		formColor = interest.color || '';
		formBgGradient = interest.bgGradient || '';
		formBorderColor = interest.borderColor || '';
		formHoverBorderColor = interest.hoverBorderColor || '';
		formShadowColor = interest.shadowColor || '';
		formFullWidth = interest.fullWidth;
		formFeatured = interest.featured;
		showEditModal = true;
	}

	function resetForm() {
		formTitle = '';
		formDescription = '';
		formIcon = '';
		formColor = '';
		formBgGradient = '';
		formBorderColor = '';
		formHoverBorderColor = '';
		formShadowColor = '';
		formFullWidth = false;
		formFeatured = false;
		editingInterest = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: CreateInterestInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				icon: formIcon.trim() || undefined,
				color: formColor.trim() || undefined,
				bgGradient: formBgGradient.trim() || undefined,
				borderColor: formBorderColor.trim() || undefined,
				hoverBorderColor: formHoverBorderColor.trim() || undefined,
				shadowColor: formShadowColor.trim() || undefined,
				fullWidth: formFullWidth,
				featured: formFeatured
			};

			await createInterest(input);
			showCreateModal = false;
			resetForm();
			await fetchInterests();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create interest');
			console.error('Error creating interest:', err);
		}
	}

	async function handleUpdate() {
		if (!editingInterest || !formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: UpdateInterestInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				icon: formIcon.trim() || undefined,
				color: formColor.trim() || undefined,
				bgGradient: formBgGradient.trim() || undefined,
				borderColor: formBorderColor.trim() || undefined,
				hoverBorderColor: formHoverBorderColor.trim() || undefined,
				shadowColor: formShadowColor.trim() || undefined,
				fullWidth: formFullWidth,
				featured: formFeatured
			};

			await updateInterest(editingInterest.id, input);
			showEditModal = false;
			resetForm();
			await fetchInterests();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update interest');
			console.error('Error updating interest:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this interest?')) return;

		try {
			await deleteInterest(id);
			await fetchInterests();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete interest');
			console.error('Error deleting interest:', err);
		}
	}

	function filteredInterests() {
		if (!searchQuery.trim()) return interests;
		const query = searchQuery.toLowerCase();
		return interests.filter(
			(i) =>
				i.title.toLowerCase().includes(query) ||
				i.description.toLowerCase().includes(query)
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
			<h1>Interests Management</h1>
			<p>Manage personal interests</p>
		</div>
		<button onclick={openCreateModal}>Create Interest</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search interests..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredInterests().length === 0}
		<div class="empty">No interests found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Description</th>
					<th>Icon</th>
					<th>Featured</th>
					<th>Full Width</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredInterests() as interest}
					<tr>
						<td><strong>{interest.title}</strong></td>
						<td>
							<div class="content-preview">{interest.description.substring(0, 80)}...</div>
						</td>
						<td>{interest.icon || '—'}</td>
						<td>{interest.featured ? 'Yes' : 'No'}</td>
						<td>{interest.fullWidth ? 'Yes' : 'No'}</td>
						<td>{formatDate(interest.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(interest)}>Edit</button>
							<button onclick={() => handleDelete(interest.id)} class="delete-btn">Delete</button>
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
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Create Interest</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Icon</label>
					<input type="text" bind:value={formIcon} placeholder="Icon identifier" />
				</div>
				<div class="form-group">
					<label>Color</label>
					<input type="text" bind:value={formColor} placeholder="Color name" />
				</div>
				<div class="form-group">
					<label>Background Gradient</label>
					<input type="text" bind:value={formBgGradient} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Border Color</label>
					<input type="text" bind:value={formBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Hover Border Color</label>
					<input type="text" bind:value={formHoverBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Shadow Color</label>
					<input type="text" bind:value={formShadowColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFullWidth} />
						Full Width
					</label>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
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
{#if showEditModal && editingInterest}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Interest</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Icon</label>
					<input type="text" bind:value={formIcon} placeholder="Icon identifier" />
				</div>
				<div class="form-group">
					<label>Color</label>
					<input type="text" bind:value={formColor} placeholder="Color name" />
				</div>
				<div class="form-group">
					<label>Background Gradient</label>
					<input type="text" bind:value={formBgGradient} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Border Color</label>
					<input type="text" bind:value={formBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Hover Border Color</label>
					<input type="text" bind:value={formHoverBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Shadow Color</label>
					<input type="text" bind:value={formShadowColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFullWidth} />
						Full Width
					</label>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
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

	.content-preview {
		max-width: 300px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
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
		max-width: 800px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-large {
		max-width: 900px;
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

