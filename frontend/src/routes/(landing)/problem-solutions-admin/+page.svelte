<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listProblemSolutions,
		createProblemSolution,
		updateProblemSolution,
		deleteProblemSolution,
		type ProblemSolution,
		type CreateProblemSolutionInput,
		type UpdateProblemSolutionInput
	} from '$lib/api/problemsolutions';

	let solutions: ProblemSolution[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingSolution: ProblemSolution | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formProblem = $state('');
	let formContext = $state('');
	let formSolution = $state('');
	let formTechnologies = $state('');
	let formImpact = $state('');
	let formMetricsBefore = $state('');
	let formMetricsAfter = $state('');
	let formMetricsImprovement = $state('');
	let formFeatured = $state(false);

	onMount(async () => {
		await fetchSolutions();
	});

	async function fetchSolutions() {
		loading = true;
		error = null;
		try {
			solutions = await listProblemSolutions();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load problem solutions';
			console.error('Error fetching problem solutions:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(sol: ProblemSolution) {
		editingSolution = sol;
		formProblem = sol.problem;
		formContext = sol.context;
		formSolution = sol.solution;
		formTechnologies = sol.technologies.join(', ');
		formImpact = sol.impact || '';
		formMetricsBefore = sol.metrics?.before || '';
		formMetricsAfter = sol.metrics?.after || '';
		formMetricsImprovement = sol.metrics?.improvement || '';
		formFeatured = sol.featured;
		showEditModal = true;
	}

	function resetForm() {
		formProblem = '';
		formContext = '';
		formSolution = '';
		formTechnologies = '';
		formImpact = '';
		formMetricsBefore = '';
		formMetricsAfter = '';
		formMetricsImprovement = '';
		formFeatured = false;
		editingSolution = null;
	}

	async function handleCreate() {
		if (!formProblem.trim() || !formContext.trim() || !formSolution.trim()) {
			alert('Problem, context, and solution are required');
			return;
		}

		try {
			const input: CreateProblemSolutionInput = {
				problem: formProblem.trim(),
				context: formContext.trim(),
				solution: formSolution.trim(),
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: [],
				impact: formImpact.trim() || undefined,
				metrics:
					formMetricsBefore || formMetricsAfter || formMetricsImprovement
						? {
								before: formMetricsBefore.trim() || undefined,
								after: formMetricsAfter.trim() || undefined,
								improvement: formMetricsImprovement.trim() || undefined
							}
						: undefined,
				featured: formFeatured
			};

			await createProblemSolution(input);
			showCreateModal = false;
			resetForm();
			await fetchSolutions();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create problem solution');
			console.error('Error creating problem solution:', err);
		}
	}

	async function handleUpdate() {
		if (!editingSolution || !formProblem.trim() || !formContext.trim() || !formSolution.trim()) {
			alert('Problem, context, and solution are required');
			return;
		}

		try {
			const input: UpdateProblemSolutionInput = {
				problem: formProblem.trim(),
				context: formContext.trim(),
				solution: formSolution.trim(),
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: [],
				impact: formImpact.trim() || undefined,
				metrics:
					formMetricsBefore || formMetricsAfter || formMetricsImprovement
						? {
								before: formMetricsBefore.trim() || undefined,
								after: formMetricsAfter.trim() || undefined,
								improvement: formMetricsImprovement.trim() || undefined
							}
						: undefined,
				featured: formFeatured
			};

			await updateProblemSolution(editingSolution.id, input);
			showEditModal = false;
			resetForm();
			await fetchSolutions();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update problem solution');
			console.error('Error updating problem solution:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this problem solution?')) return;

		try {
			await deleteProblemSolution(id);
			await fetchSolutions();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete problem solution');
			console.error('Error deleting problem solution:', err);
		}
	}

	function filteredSolutions() {
		if (!searchQuery.trim()) return solutions;
		const query = searchQuery.toLowerCase();
		return solutions.filter(
			(s) =>
				s.problem.toLowerCase().includes(query) ||
				s.solution.toLowerCase().includes(query) ||
				s.context.toLowerCase().includes(query)
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
			<h1>Problem Solutions Management</h1>
			<p>Manage problem-solution documents</p>
		</div>
		<button onclick={openCreateModal}>Create Problem Solution</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search problem solutions..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredSolutions().length === 0}
		<div class="empty">No problem solutions found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Problem</th>
					<th>Technologies</th>
					<th>Featured</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredSolutions() as sol}
					<tr>
						<td>
							<strong>{sol.problem.substring(0, 80)}...</strong>
						</td>
						<td>{sol.technologies.slice(0, 3).join(', ')}</td>
						<td>{sol.featured ? 'Yes' : 'No'}</td>
						<td>{formatDate(sol.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(sol)}>Edit</button>
							<button onclick={() => handleDelete(sol.id)} class="delete-btn">Delete</button>
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
			<h2>Create Problem Solution</h2>
			<div class="form">
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
					<textarea bind:value={formSolution} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Impact</label>
					<textarea bind:value={formImpact} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Metrics - Before</label>
					<input type="text" bind:value={formMetricsBefore} />
				</div>
				<div class="form-group">
					<label>Metrics - After</label>
					<input type="text" bind:value={formMetricsAfter} />
				</div>
				<div class="form-group">
					<label>Metrics - Improvement</label>
					<input type="text" bind:value={formMetricsImprovement} />
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
{#if showEditModal && editingSolution}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Problem Solution</h2>
			<div class="form">
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
					<textarea bind:value={formSolution} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Impact</label>
					<textarea bind:value={formImpact} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Metrics - Before</label>
					<input type="text" bind:value={formMetricsBefore} />
				</div>
				<div class="form-group">
					<label>Metrics - After</label>
					<input type="text" bind:value={formMetricsAfter} />
				</div>
				<div class="form-group">
					<label>Metrics - Improvement</label>
					<input type="text" bind:value={formMetricsImprovement} />
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

