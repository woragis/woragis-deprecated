<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getProject,
		updateProjectStatus,
		updateProjectMetrics,
		deleteProject,
		type Project,
		type ProjectStatus
	} from '$lib/api/projects';

	let project: Project | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showEditModal = $state(false);

	// Form state
	let formStatus = $state<ProjectStatus>('idea');
	let formHealthScore = $state<number | ''>(0);
	let formMrr = $state<number | ''>(0);
	let formCac = $state<number | ''>(0);
	let formLtv = $state<number | ''>(0);
	let formChurnRate = $state<number | ''>(0);

	const projectId = $derived($page.params.id);

	const statuses: ProjectStatus[] = ['idea', 'planning', 'executing', 'monitoring', 'completed'];

	onMount(async () => {
		if (projectId) {
			await loadProject();
		}
	});

	async function loadProject() {
		if (!projectId) return;
		loading = true;
		error = null;
		try {
			project = await getProject(projectId);
			formStatus = project.status;
			formHealthScore = project.health_score;
			formMrr = project.mrr;
			formCac = project.cac;
			formLtv = project.ltv;
			formChurnRate = project.churn_rate;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load project';
			console.error('Error loading project:', err);
		} finally {
			loading = false;
		}
	}

	function openEditModal() {
		if (!project) return;
		formStatus = project.status;
		formHealthScore = project.health_score;
		formMrr = project.mrr;
		formCac = project.cac;
		formLtv = project.ltv;
		formChurnRate = project.churn_rate;
		showEditModal = true;
	}

	async function handleUpdate() {
		if (!project) return;

		try {
			await updateProjectStatus(project.id, formStatus);
			await updateProjectMetrics(project.id, {
				healthScore: formHealthScore ? Number(formHealthScore) : undefined,
				mrr: formMrr ? Number(formMrr) : undefined,
				cac: formCac ? Number(formCac) : undefined,
				ltv: formLtv ? Number(formLtv) : undefined,
				churnRate: formChurnRate ? Number(formChurnRate) : undefined
			});
			showEditModal = false;
			await loadProject();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update project');
			console.error('Error updating project:', err);
		}
	}

	async function handleDelete() {
		if (!project || !confirm('Are you sure you want to delete this project? This will also delete all related milestones, kanban boards, dependencies, documentation, technologies, file structures, and architecture diagrams.')) return;

		try {
			await deleteProject(project.id);
			await goto('/projects-admin');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete project');
			console.error('Error deleting project:', err);
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="header">
		<a href="/projects-admin" class="back-link">← Back to Projects</a>
		<div class="header-actions">
			{#if project}
				<button onclick={openEditModal}>Edit Project</button>
				<button onclick={handleDelete} class="delete-btn">Delete</button>
			{/if}
		</div>
	</div>

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if error}
		<div class="error">{error}</div>
	{:else if project}
		<div class="details-container">
			<div class="details-header">
				<h2>{project.name}</h2>
			</div>

			<div class="details-grid">
				<div class="detail-section">
					<h3>Project Information</h3>
					<div class="detail-item">
						<strong>Status:</strong> <span class="status status-{project.status}">{project.status}</span>
					</div>
					<div class="detail-item">
						<strong>Description:</strong> {project.description || '—'}
					</div>
					<div class="detail-item">
						<strong>Slug:</strong> {project.slug || '—'}
					</div>
					<div class="detail-item">
						<strong>Created:</strong> {formatDate(project.created_at)}
					</div>
					<div class="detail-item">
						<strong>Updated:</strong> {formatDate(project.updated_at)}
					</div>
				</div>

				<div class="detail-section">
					<h3>Metrics</h3>
					<div class="detail-item">
						<strong>Health Score:</strong> {project.health_score}
					</div>
					<div class="detail-item">
						<strong>MRR (Monthly Recurring Revenue):</strong> ${project.mrr.toFixed(2)}
					</div>
					<div class="detail-item">
						<strong>CAC (Customer Acquisition Cost):</strong> ${project.cac.toFixed(2)}
					</div>
					<div class="detail-item">
						<strong>LTV (Lifetime Value):</strong> ${project.ltv.toFixed(2)}
					</div>
					<div class="detail-item">
						<strong>Churn Rate:</strong> {project.churn_rate.toFixed(2)}%
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Edit Modal -->
{#if showEditModal && project}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Project</h2>
			<div class="form">
				<div class="form-group">
					<label>Name</label>
					<input type="text" value={project.name} disabled />
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Health Score</label>
					<input type="number" bind:value={formHealthScore} />
				</div>
				<div class="form-group">
					<label>MRR (Monthly Recurring Revenue)</label>
					<input type="number" step="0.01" bind:value={formMrr} />
				</div>
				<div class="form-group">
					<label>CAC (Customer Acquisition Cost)</label>
					<input type="number" step="0.01" bind:value={formCac} />
				</div>
				<div class="form-group">
					<label>LTV (Lifetime Value)</label>
					<input type="number" step="0.01" bind:value={formLtv} />
				</div>
				<div class="form-group">
					<label>Churn Rate (%)</label>
					<input type="number" step="0.01" bind:value={formChurnRate} />
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
		max-width: 1400px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.back-link {
		color: #007bff;
		text-decoration: none;
		font-size: 0.9rem;
		padding: 0.5rem 0;
	}

	.back-link:hover {
		color: #0056b3;
		text-decoration: underline;
	}

	.header-actions {
		display: flex;
		gap: 0.5rem;
	}

	.header-actions button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.header-actions button:hover {
		background: #0056b3;
	}

	.delete-btn {
		background: #dc3545 !important;
	}

	.delete-btn:hover {
		background: #c82333 !important;
	}

	.loading,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #666;
	}

	.error {
		padding: 0.75rem;
		background: #fee;
		color: #c33;
		border: 1px solid #fcc;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.details-container {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		color: #333;
	}

	.details-header {
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 2px solid #e0e0e0;
	}

	.details-header h2 {
		margin: 0;
		font-size: 1.5rem;
		color: #333;
	}

	.details-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}

	.detail-section {
		background: #f9f9f9;
		padding: 1rem;
		border-radius: 4px;
	}

	.detail-section h3 {
		margin: 0 0 1rem 0;
		font-size: 1.1rem;
		color: #333;
	}

	.detail-item {
		margin-bottom: 0.5rem;
		font-size: 0.9rem;
		color: #333;
	}

	.detail-item strong {
		color: #333;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-idea {
		background: #e2e3e5;
		color: #383d41;
	}

	.status-planning {
		background: #fff3cd;
		color: #856404;
	}

	.status-executing {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-monitoring {
		background: #d4edda;
		color: #155724;
	}

	.status-completed {
		background: #d1ecf1;
		color: #0c5460;
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
	.form-group select {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
	}

	.form-group input:disabled {
		background: #f5f5f5;
		cursor: not-allowed;
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

