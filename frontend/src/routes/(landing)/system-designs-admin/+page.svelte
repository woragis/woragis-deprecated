<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listSystemDesigns,
		createSystemDesign,
		updateSystemDesign,
		deleteSystemDesign,
		type SystemDesign,
		type CreateSystemDesignInput,
		type UpdateSystemDesignInput
	} from '$lib/api/systemdesigns';

	let designs: SystemDesign[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingDesign: SystemDesign | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formTitle = $state('');
	let formDescription = $state('');
	let formDataFlow = $state('');
	let formScalability = $state('');
	let formReliability = $state('');
	let formDiagram = $state('');
	let formFeatured = $state(false);

	onMount(async () => {
		await fetchDesigns();
	});

	async function fetchDesigns() {
		loading = true;
		error = null;
		try {
			designs = await listSystemDesigns();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load system designs';
			console.error('Error fetching system designs:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(design: SystemDesign) {
		editingDesign = design;
		formTitle = design.title;
		formDescription = design.description;
		formDataFlow = design.dataFlow || '';
		formScalability = design.scalability || '';
		formReliability = design.reliability || '';
		formDiagram = design.diagram || '';
		formFeatured = design.featured;
		showEditModal = true;
	}

	function resetForm() {
		formTitle = '';
		formDescription = '';
		formDataFlow = '';
		formScalability = '';
		formReliability = '';
		formDiagram = '';
		formFeatured = false;
		editingDesign = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: CreateSystemDesignInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				dataFlow: formDataFlow.trim() || undefined,
				scalability: formScalability.trim() || undefined,
				reliability: formReliability.trim() || undefined,
				diagram: formDiagram.trim() || undefined,
				featured: formFeatured
			};

			await createSystemDesign(input);
			showCreateModal = false;
			resetForm();
			await fetchDesigns();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create system design');
			console.error('Error creating system design:', err);
		}
	}

	async function handleUpdate() {
		if (!editingDesign || !formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: UpdateSystemDesignInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				dataFlow: formDataFlow.trim() || undefined,
				scalability: formScalability.trim() || undefined,
				reliability: formReliability.trim() || undefined,
				diagram: formDiagram.trim() || undefined,
				featured: formFeatured
			};

			await updateSystemDesign(editingDesign.id, input);
			showEditModal = false;
			resetForm();
			await fetchDesigns();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update system design');
			console.error('Error updating system design:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this system design?')) return;

		try {
			await deleteSystemDesign(id);
			await fetchDesigns();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete system design');
			console.error('Error deleting system design:', err);
		}
	}

	function filteredDesigns() {
		if (!searchQuery.trim()) return designs;
		const query = searchQuery.toLowerCase();
		return designs.filter(
			(d) =>
				d.title.toLowerCase().includes(query) ||
				d.description.toLowerCase().includes(query)
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
			<h1>System Designs Management</h1>
			<p>Manage system design documents</p>
		</div>
		<button onclick={openCreateModal}>Create System Design</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search system designs..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredDesigns().length === 0}
		<div class="empty">No system designs found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Description</th>
					<th>Featured</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredDesigns() as design}
					<tr>
						<td>
							<strong>{design.title}</strong>
						</td>
						<td>
							<div class="content-preview">{design.description.substring(0, 100)}...</div>
						</td>
						<td>{design.featured ? 'Yes' : 'No'}</td>
						<td>{formatDate(design.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(design)}>Edit</button>
							<button onclick={() => handleDelete(design.id)} class="delete-btn">Delete</button>
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
			<h2>Create System Design</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Data Flow</label>
					<textarea bind:value={formDataFlow} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Scalability</label>
					<textarea bind:value={formScalability} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Reliability</label>
					<textarea bind:value={formReliability} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Diagram URL</label>
					<input type="url" bind:value={formDiagram} />
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
{#if showEditModal && editingDesign}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit System Design</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Data Flow</label>
					<textarea bind:value={formDataFlow} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Scalability</label>
					<textarea bind:value={formScalability} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Reliability</label>
					<textarea bind:value={formReliability} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Diagram URL</label>
					<input type="url" bind:value={formDiagram} />
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

