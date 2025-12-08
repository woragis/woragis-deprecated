<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listCaseStudies,
		createCaseStudy,
		updateCaseStudy,
		deleteCaseStudy,
		type CaseStudy,
		type CreateCaseStudyInput,
		type UpdateCaseStudyInput
	} from '$lib/api/casestudies';

	let caseStudies: CaseStudy[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingCaseStudy: CaseStudy | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formProjectId = $state('');
	let formProjectSlug = $state('');
	let formTitle = $state('');
	let formProblem = $state('');
	let formContext = $state('');
	let formSolution = $state('');
	let formApproach = $state('');
	let formLessonsLearned = $state('');
	let formTechnologies = $state('');
	let formFeatured = $state(false);

	onMount(async () => {
		await fetchCaseStudies();
	});

	async function fetchCaseStudies() {
		loading = true;
		error = null;
		try {
			caseStudies = await listCaseStudies();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load case studies';
			console.error('Error fetching case studies:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(cs: CaseStudy) {
		editingCaseStudy = cs;
		formProjectId = cs.projectId;
		formProjectSlug = cs.projectSlug;
		formTitle = cs.title;
		formProblem = cs.problem;
		formContext = cs.context;
		formSolution = cs.solution;
		formApproach = cs.approach.join('\n');
		formLessonsLearned = cs.lessonsLearned.join('\n');
		formTechnologies = cs.technologies.join(', ');
		formFeatured = cs.featured;
		showEditModal = true;
	}

	function resetForm() {
		formProjectId = '';
		formProjectSlug = '';
		formTitle = '';
		formProblem = '';
		formContext = '';
		formSolution = '';
		formApproach = '';
		formLessonsLearned = '';
		formTechnologies = '';
		formFeatured = false;
		editingCaseStudy = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formProblem.trim() || !formContext.trim() || !formSolution.trim()) {
			alert('Title, problem, context, and solution are required');
			return;
		}

		try {
			const input: CreateCaseStudyInput = {
				projectId: formProjectId.trim(),
				projectSlug: formProjectSlug.trim(),
				title: formTitle.trim(),
				problem: formProblem.trim(),
				context: formContext.trim(),
				solution: formSolution.trim(),
				approach: formApproach.trim() ? formApproach.trim().split('\n').filter((s) => s.trim()) : [],
				lessonsLearned: formLessonsLearned.trim()
					? formLessonsLearned.trim().split('\n').filter((s) => s.trim())
					: [],
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: [],
				featured: formFeatured
			};

			await createCaseStudy(input);
			showCreateModal = false;
			resetForm();
			await fetchCaseStudies();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create case study');
			console.error('Error creating case study:', err);
		}
	}

	async function handleUpdate() {
		if (!editingCaseStudy || !formTitle.trim() || !formProblem.trim() || !formContext.trim() || !formSolution.trim()) {
			alert('Title, problem, context, and solution are required');
			return;
		}

		try {
			const input: UpdateCaseStudyInput = {
				title: formTitle.trim(),
				problem: formProblem.trim(),
				context: formContext.trim(),
				solution: formSolution.trim(),
				approach: formApproach.trim() ? formApproach.trim().split('\n').filter((s) => s.trim()) : [],
				lessonsLearned: formLessonsLearned.trim()
					? formLessonsLearned.trim().split('\n').filter((s) => s.trim())
					: [],
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: [],
				featured: formFeatured
			};

			await updateCaseStudy(editingCaseStudy.id, input);
			showEditModal = false;
			resetForm();
			await fetchCaseStudies();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update case study');
			console.error('Error updating case study:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this case study?')) return;

		try {
			await deleteCaseStudy(id);
			await fetchCaseStudies();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete case study');
			console.error('Error deleting case study:', err);
		}
	}

	function filteredCaseStudies() {
		if (!searchQuery.trim()) return caseStudies;
		const query = searchQuery.toLowerCase();
		return caseStudies.filter(
			(cs) =>
				cs.title.toLowerCase().includes(query) ||
				cs.problem.toLowerCase().includes(query) ||
				cs.projectSlug.toLowerCase().includes(query)
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
			<h1>Case Studies Management</h1>
			<p>Manage project case studies</p>
		</div>
		<button onclick={openCreateModal}>Create Case Study</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search case studies..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredCaseStudies().length === 0}
		<div class="empty">No case studies found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Project</th>
					<th>Featured</th>
					<th>Technologies</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredCaseStudies() as cs}
					<tr>
						<td>
							<strong>{cs.title}</strong>
							<br />
							<small>{cs.problem.substring(0, 100)}...</small>
						</td>
						<td>{cs.projectSlug}</td>
						<td>{cs.featured ? 'Yes' : 'No'}</td>
						<td>{cs.technologies.slice(0, 3).join(', ')}</td>
						<td>{formatDate(cs.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(cs)}>Edit</button>
							<button onclick={() => handleDelete(cs.id)} class="delete-btn">Delete</button>
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
			<h2>Create Case Study</h2>
			<div class="form">
				<div class="form-group">
					<label>Project ID *</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-group">
					<label>Project Slug *</label>
					<input type="text" bind:value={formProjectSlug} />
				</div>
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Problem *</label>
					<textarea bind:value={formProblem} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Context *</label>
					<textarea bind:value={formContext} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Solution *</label>
					<textarea bind:value={formSolution} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Approach (one per line)</label>
					<textarea bind:value={formApproach} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Lessons Learned (one per line)</label>
					<textarea bind:value={formLessonsLearned} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
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
{#if showEditModal && editingCaseStudy}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Case Study</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Problem *</label>
					<textarea bind:value={formProblem} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Context *</label>
					<textarea bind:value={formContext} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Solution *</label>
					<textarea bind:value={formSolution} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Approach (one per line)</label>
					<textarea bind:value={formApproach} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Lessons Learned (one per line)</label>
					<textarea bind:value={formLessonsLearned} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
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

	.table small {
		color: #666;
		font-size: 0.875rem;
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

