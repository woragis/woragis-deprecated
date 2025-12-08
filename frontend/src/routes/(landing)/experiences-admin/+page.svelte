<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listExperiences,
		createExperience,
		updateExperience,
		deleteExperience,
		type Experience,
		type CreateExperienceInput,
		type UpdateExperienceInput,
		type ExperienceType
	} from '$lib/api/experiences';

	let experiences: Experience[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingExperience: Experience | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formCompany = $state('');
	let formPosition = $state('');
	let formPeriodStart = $state('');
	let formPeriodEnd = $state('');
	let formPeriodText = $state('');
	let formLocation = $state('');
	let formDescription = $state('');
	let formType = $state<ExperienceType>('full-time');
	let formCompanyUrl = $state('');
	let formLinkedinUrl = $state('');
	let formDisplayOrder = $state<number | ''>(0);
	let formIsCurrent = $state(false);

	const types: ExperienceType[] = ['full-time', 'freelance', 'contract', 'internship'];

	onMount(async () => {
		await fetchExperiences();
	});

	async function fetchExperiences() {
		loading = true;
		error = null;
		try {
			experiences = await listExperiences();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load experiences';
			console.error('Error fetching experiences:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(exp: Experience) {
		editingExperience = exp;
		formCompany = exp.company;
		formPosition = exp.position;
		formPeriodStart = exp.periodStart ? exp.periodStart.split('T')[0] : '';
		formPeriodEnd = exp.periodEnd ? exp.periodEnd.split('T')[0] : '';
		formPeriodText = exp.periodText || '';
		formLocation = exp.location || '';
		formDescription = exp.description || '';
		formType = exp.type;
		formCompanyUrl = exp.companyUrl || '';
		formLinkedinUrl = exp.linkedinUrl || '';
		formDisplayOrder = exp.displayOrder;
		formIsCurrent = exp.isCurrent;
		showEditModal = true;
	}

	function resetForm() {
		formCompany = '';
		formPosition = '';
		formPeriodStart = '';
		formPeriodEnd = '';
		formPeriodText = '';
		formLocation = '';
		formDescription = '';
		formType = 'full-time';
		formCompanyUrl = '';
		formLinkedinUrl = '';
		formDisplayOrder = 0;
		formIsCurrent = false;
		editingExperience = null;
	}

	async function handleCreate() {
		if (!formCompany.trim() || !formPosition.trim()) {
			alert('Company and position are required');
			return;
		}

		try {
			const input: CreateExperienceInput = {
				company: formCompany.trim(),
				position: formPosition.trim(),
				periodStart: formPeriodStart || undefined,
				periodEnd: formPeriodEnd || undefined,
				periodText: formPeriodText.trim() || undefined,
				location: formLocation.trim() || undefined,
				description: formDescription.trim() || undefined,
				type: formType,
				companyUrl: formCompanyUrl.trim() || undefined,
				linkedinUrl: formLinkedinUrl.trim() || undefined,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				isCurrent: formIsCurrent
			};

			await createExperience(input);
			showCreateModal = false;
			resetForm();
			await fetchExperiences();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create experience');
			console.error('Error creating experience:', err);
		}
	}

	async function handleUpdate() {
		if (!editingExperience || !formCompany.trim() || !formPosition.trim()) {
			alert('Company and position are required');
			return;
		}

		try {
			const input: UpdateExperienceInput = {
				company: formCompany.trim(),
				position: formPosition.trim(),
				periodStart: formPeriodStart || undefined,
				periodEnd: formPeriodEnd || undefined,
				periodText: formPeriodText.trim() || undefined,
				location: formLocation.trim() || undefined,
				description: formDescription.trim() || undefined,
				type: formType,
				companyUrl: formCompanyUrl.trim() || undefined,
				linkedinUrl: formLinkedinUrl.trim() || undefined,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				isCurrent: formIsCurrent
			};

			await updateExperience(editingExperience.id, input);
			showEditModal = false;
			resetForm();
			await fetchExperiences();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update experience');
			console.error('Error updating experience:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this experience?')) return;

		try {
			await deleteExperience(id);
			await fetchExperiences();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete experience');
			console.error('Error deleting experience:', err);
		}
	}

	function filteredExperiences() {
		if (!searchQuery.trim()) return experiences;
		const query = searchQuery.toLowerCase();
		return experiences.filter(
			(e) =>
				e.company.toLowerCase().includes(query) ||
				e.position.toLowerCase().includes(query) ||
				e.location?.toLowerCase().includes(query)
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
			<h1>Experiences Management</h1>
			<p>Manage work experiences</p>
		</div>
		<button onclick={openCreateModal}>Create Experience</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search experiences..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredExperiences().length === 0}
		<div class="empty">No experiences found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Company</th>
					<th>Position</th>
					<th>Type</th>
					<th>Location</th>
					<th>Period</th>
					<th>Current</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredExperiences() as exp}
					<tr>
						<td><strong>{exp.company}</strong></td>
						<td>{exp.position}</td>
						<td>{exp.type}</td>
						<td>{exp.location || '—'}</td>
						<td>
							{exp.periodText || formatDate(exp.periodStart) + ' - ' + formatDate(exp.periodEnd)}
						</td>
						<td>{exp.isCurrent ? 'Yes' : 'No'}</td>
						<td>{formatDate(exp.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(exp)}>Edit</button>
							<button onclick={() => handleDelete(exp.id)} class="delete-btn">Delete</button>
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
			<h2>Create Experience</h2>
			<div class="form">
				<div class="form-group">
					<label>Company *</label>
					<input type="text" bind:value={formCompany} />
				</div>
				<div class="form-group">
					<label>Position *</label>
					<input type="text" bind:value={formPosition} />
				</div>
				<div class="form-group">
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Period Start</label>
					<input type="date" bind:value={formPeriodStart} />
				</div>
				<div class="form-group">
					<label>Period End</label>
					<input type="date" bind:value={formPeriodEnd} />
				</div>
				<div class="form-group">
					<label>Period Text (e.g., "2022 - Present")</label>
					<input type="text" bind:value={formPeriodText} />
				</div>
				<div class="form-group">
					<label>Location</label>
					<input type="text" bind:value={formLocation} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Company URL</label>
					<input type="url" bind:value={formCompanyUrl} />
				</div>
				<div class="form-group">
					<label>LinkedIn URL</label>
					<input type="url" bind:value={formLinkedinUrl} />
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formIsCurrent} />
						Is Current
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
{#if showEditModal && editingExperience}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Experience</h2>
			<div class="form">
				<div class="form-group">
					<label>Company *</label>
					<input type="text" bind:value={formCompany} />
				</div>
				<div class="form-group">
					<label>Position *</label>
					<input type="text" bind:value={formPosition} />
				</div>
				<div class="form-group">
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Period Start</label>
					<input type="date" bind:value={formPeriodStart} />
				</div>
				<div class="form-group">
					<label>Period End</label>
					<input type="date" bind:value={formPeriodEnd} />
				</div>
				<div class="form-group">
					<label>Period Text (e.g., "2022 - Present")</label>
					<input type="text" bind:value={formPeriodText} />
				</div>
				<div class="form-group">
					<label>Location</label>
					<input type="text" bind:value={formLocation} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Company URL</label>
					<input type="url" bind:value={formCompanyUrl} />
				</div>
				<div class="form-group">
					<label>LinkedIn URL</label>
					<input type="url" bind:value={formLinkedinUrl} />
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formIsCurrent} />
						Is Current
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

