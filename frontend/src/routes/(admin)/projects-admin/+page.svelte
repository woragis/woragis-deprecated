<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listProjects,
		createProject,
		updateProjectStatus,
		updateProjectMetrics,
		deleteProject,
		type Project,
		type CreateProjectInput,
		type ProjectStatus
	} from '$lib/api/projects';
	import { locale, t } from '$lib/i18n';

	let projects: Project[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let searchQuery = $state('');

	// Form state
	let formName = $state('');
	let formDescription = $state('');
	let formStatus = $state<ProjectStatus>('idea');
	let formHealthScore = $state<number | ''>(0);
	let formMrr = $state<number | ''>(0);
	let formCac = $state<number | ''>(0);
	let formLtv = $state<number | ''>(0);
	let formChurnRate = $state<number | ''>(0);

	const statuses: ProjectStatus[] = ['idea', 'planning', 'executing', 'monitoring', 'completed'];

	onMount(async () => {
		await fetchProjects();
	});

	async function fetchProjects() {
		loading = true;
		error = null;
		try {
			projects = await listProjects();
		} catch (err) {
			error = err instanceof Error ? err.message : $t('projects.error');
			console.error('Error fetching projects:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function resetForm() {
		formName = '';
		formDescription = '';
		formStatus = 'idea';
		formHealthScore = 0;
		formMrr = 0;
		formCac = 0;
		formLtv = 0;
		formChurnRate = 0;
	}

	async function handleCreate() {
		if (!formName.trim()) {
			alert($t('projects.modal.name') + ' ' + $t('projects.modal.required'));
			return;
		}

		try {
			const input: CreateProjectInput = {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				status: formStatus,
				healthScore: formHealthScore ? Number(formHealthScore) : undefined,
				mrr: formMrr ? Number(formMrr) : undefined,
				cac: formCac ? Number(formCac) : undefined,
				ltv: formLtv ? Number(formLtv) : undefined,
				churnRate: formChurnRate ? Number(formChurnRate) : undefined
			};

			await createProject(input);
			showCreateModal = false;
			resetForm();
			await fetchProjects();
		} catch (err) {
			alert(err instanceof Error ? err.message : $t('projects.createError'));
			console.error('Error creating project:', err);
		}
	}


	async function handleDelete(id: string) {
		if (!confirm($t('projects.deleteConfirm'))) return;

		try {
			await deleteProject(id);
			await fetchProjects();
		} catch (err) {
			alert(err instanceof Error ? err.message : $t('projects.deleteError'));
			console.error('Error deleting project:', err);
		}
	}

	function filteredProjects() {
		if (!searchQuery.trim()) return projects;
		const query = searchQuery.toLowerCase();
		return projects.filter(
			(p) =>
				p.name.toLowerCase().includes(query) ||
				p.description?.toLowerCase().includes(query) ||
				p.status.toLowerCase().includes(query)
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
			<h1>{$t('projects.title')}</h1>
			<p>{$t('projects.subtitle')}</p>
		</div>
		<button onclick={openCreateModal}>{$t('projects.createButton')}</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder={$t('projects.searchPlaceholder')}
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">{$t('projects.loading')}</div>
	{:else if filteredProjects().length === 0}
		<div class="empty">{$t('projects.empty')}</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>{$t('projects.table.name')}</th>
					<th>{$t('projects.table.status')}</th>
					<th>{$t('projects.table.healthScore')}</th>
					<th>{$t('projects.table.mrr')}</th>
					<th>{$t('projects.table.cac')}</th>
					<th>{$t('projects.table.ltv')}</th>
					<th>{$t('projects.table.churnRate')}</th>
					<th>{$t('projects.table.created')}</th>
					<th>{$t('projects.table.actions')}</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredProjects() as project}
					<tr>
						<td>
							<strong>{project.name}</strong>
							<br />
							<small>{project.description || '—'}</small>
						</td>
						<td>
							<span class="status status-{project.status}">{$t(`projects.status.${project.status}` as any)}</span>
						</td>
						<td>{project.health_score}</td>
						<td>${project.mrr.toFixed(2)}</td>
						<td>${project.cac.toFixed(2)}</td>
						<td>${project.ltv.toFixed(2)}</td>
						<td>{project.churn_rate.toFixed(2)}%</td>
						<td>{formatDate(project.created_at)}</td>
						<td>
							<a href="/projects-admin/{project.id}" class="view-link">{$t('projects.table.view')}</a>
							<button onclick={() => handleDelete(project.id)} class="delete-btn">{$t('projects.table.delete')}</button>
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
			<h2>{$t('projects.modal.createTitle')}</h2>
			<div class="form">
				<div class="form-group">
					<label>{$t('projects.modal.name')} {$t('projects.modal.required')}</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.description')}</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.status')}</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{$t(`projects.status.${status}` as any)}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.healthScore')}</label>
					<input type="number" bind:value={formHealthScore} />
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.mrr')}</label>
					<input type="number" step="0.01" bind:value={formMrr} />
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.cac')}</label>
					<input type="number" step="0.01" bind:value={formCac} />
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.ltv')}</label>
					<input type="number" step="0.01" bind:value={formLtv} />
				</div>
				<div class="form-group">
					<label>{$t('projects.modal.churnRate')}</label>
					<input type="number" step="0.01" bind:value={formChurnRate} />
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>{$t('projects.modal.create')}</button>
					<button onclick={() => (showCreateModal = false)}>{$t('projects.modal.cancel')}</button>
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

	.table small {
		color: #666;
		font-size: 0.875rem;
	}

	.table button {
		padding: 0.25rem 0.75rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
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

	.view-link {
		padding: 0.25rem 0.75rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		margin-right: 0.5rem;
		text-decoration: none;
		display: inline-block;
	}

	.view-link:hover {
		background: #218838;
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
	.form-group textarea,
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

