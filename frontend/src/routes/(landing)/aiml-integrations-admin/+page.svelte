<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listAIMLIntegrations,
		createAIMLIntegration,
		updateAIMLIntegration,
		deleteAIMLIntegration,
		type AIMLIntegration,
		type CreateAIMLIntegrationInput,
		type UpdateAIMLIntegrationInput,
		type IntegrationType,
		type Framework
	} from '$lib/api/aimlintegrations';

	let integrations: AIMLIntegration[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingIntegration: AIMLIntegration | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formTitle = $state('');
	let formDescription = $state('');
	let formType = $state<IntegrationType>('llm');
	let formFramework = $state<Framework>('openai');
	let formModelName = $state('');
	let formModelVersion = $state('');
	let formUseCase = $state('');
	let formImpact = $state('');
	let formTechnologies = $state('');
	let formArchitecture = $state('');
	let formMetrics = $state('');
	let formProjectId = $state('');
	let formCaseStudyId = $state('');
	let formFeatured = $state(false);
	let formDisplayOrder = $state<number | ''>(0);
	let formDemoUrl = $state('');
	let formDocumentationUrl = $state('');
	let formGithubUrl = $state('');

	const types: IntegrationType[] = [
		'rag',
		'llm',
		'ml_model',
		'computer_vision',
		'nlp',
		'recommendation',
		'chatbot',
		'anomaly_detection',
		'predictive_analytics',
		'generative_ai',
		'other'
	];
	const frameworks: Framework[] = [
		'openai',
		'anthropic',
		'huggingface',
		'tensorflow',
		'pytorch',
		'langchain',
		'llamaindex',
		'cohere',
		'google_ai',
		'azure_ai',
		'aws_bedrock',
		'custom',
		'other'
	];

	onMount(async () => {
		await fetchIntegrations();
	});

	async function fetchIntegrations() {
		loading = true;
		error = null;
		try {
			integrations = await listAIMLIntegrations();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load AI/ML integrations';
			console.error('Error fetching AI/ML integrations:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(integration: AIMLIntegration) {
		editingIntegration = integration;
		formTitle = integration.title;
		formDescription = integration.description;
		formType = integration.type;
		formFramework = integration.framework;
		formModelName = integration.modelName || '';
		formModelVersion = integration.modelVersion || '';
		formUseCase = integration.useCase || '';
		formImpact = integration.impact || '';
		formTechnologies = integration.technologies?.join(', ') || '';
		formArchitecture = integration.architecture || '';
		formMetrics = integration.metrics || '';
		formProjectId = integration.projectId || '';
		formCaseStudyId = integration.caseStudyId || '';
		formFeatured = integration.featured;
		formDisplayOrder = integration.displayOrder;
		formDemoUrl = integration.demoUrl || '';
		formDocumentationUrl = integration.documentationUrl || '';
		formGithubUrl = integration.githubUrl || '';
		showEditModal = true;
	}

	function resetForm() {
		formTitle = '';
		formDescription = '';
		formType = 'llm';
		formFramework = 'openai';
		formModelName = '';
		formModelVersion = '';
		formUseCase = '';
		formImpact = '';
		formTechnologies = '';
		formArchitecture = '';
		formMetrics = '';
		formProjectId = '';
		formCaseStudyId = '';
		formFeatured = false;
		formDisplayOrder = 0;
		formDemoUrl = '';
		formDocumentationUrl = '';
		formGithubUrl = '';
		editingIntegration = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: CreateAIMLIntegrationInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				type: formType,
				framework: formFramework,
				modelName: formModelName.trim() || undefined,
				modelVersion: formModelVersion.trim() || undefined,
				useCase: formUseCase.trim() || undefined,
				impact: formImpact.trim() || undefined,
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: undefined,
				architecture: formArchitecture.trim() || undefined,
				metrics: formMetrics.trim() || undefined,
				projectId: formProjectId.trim() || undefined,
				caseStudyId: formCaseStudyId.trim() || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				demoUrl: formDemoUrl.trim() || undefined,
				documentationUrl: formDocumentationUrl.trim() || undefined,
				githubUrl: formGithubUrl.trim() || undefined
			};

			await createAIMLIntegration(input);
			showCreateModal = false;
			resetForm();
			await fetchIntegrations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create AI/ML integration');
			console.error('Error creating AI/ML integration:', err);
		}
	}

	async function handleUpdate() {
		if (!editingIntegration || !formTitle.trim() || !formDescription.trim()) {
			alert('Title and description are required');
			return;
		}

		try {
			const input: UpdateAIMLIntegrationInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				type: formType,
				framework: formFramework,
				modelName: formModelName.trim() || undefined,
				modelVersion: formModelVersion.trim() || undefined,
				useCase: formUseCase.trim() || undefined,
				impact: formImpact.trim() || undefined,
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: undefined,
				architecture: formArchitecture.trim() || undefined,
				metrics: formMetrics.trim() || undefined,
				projectId: formProjectId.trim() || undefined,
				caseStudyId: formCaseStudyId.trim() || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				demoUrl: formDemoUrl.trim() || undefined,
				documentationUrl: formDocumentationUrl.trim() || undefined,
				githubUrl: formGithubUrl.trim() || undefined
			};

			await updateAIMLIntegration(editingIntegration.id, input);
			showEditModal = false;
			resetForm();
			await fetchIntegrations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update AI/ML integration');
			console.error('Error updating AI/ML integration:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this AI/ML integration?')) return;

		try {
			await deleteAIMLIntegration(id);
			await fetchIntegrations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete AI/ML integration');
			console.error('Error deleting AI/ML integration:', err);
		}
	}

	function filteredIntegrations() {
		if (!searchQuery.trim()) return integrations;
		const query = searchQuery.toLowerCase();
		return integrations.filter(
			(i) =>
				i.title.toLowerCase().includes(query) ||
				i.description.toLowerCase().includes(query) ||
				i.type.toLowerCase().includes(query) ||
				i.framework.toLowerCase().includes(query)
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
			<h1>AI/ML Integrations Management</h1>
			<p>Manage AI/ML integration showcases</p>
		</div>
		<button onclick={openCreateModal}>Create Integration</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search integrations..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredIntegrations().length === 0}
		<div class="empty">No AI/ML integrations found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Type</th>
					<th>Framework</th>
					<th>Model</th>
					<th>Featured</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredIntegrations() as integration}
					<tr>
						<td>
							<strong>{integration.title}</strong>
							<br />
							<small>{integration.description.substring(0, 60)}...</small>
						</td>
						<td>{integration.type}</td>
						<td>{integration.framework}</td>
						<td>{integration.modelName || '—'}</td>
						<td>{integration.featured ? 'Yes' : 'No'}</td>
						<td>{formatDate(integration.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(integration)}>Edit</button>
							<button onclick={() => handleDelete(integration.id)} class="delete-btn">Delete</button>
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
			<h2>Create AI/ML Integration</h2>
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
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Framework</label>
					<select bind:value={formFramework}>
						{#each frameworks as framework}
							<option value={framework}>{framework}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Model Name</label>
					<input type="text" bind:value={formModelName} />
				</div>
				<div class="form-group">
					<label>Model Version</label>
					<input type="text" bind:value={formModelVersion} />
				</div>
				<div class="form-group">
					<label>Use Case</label>
					<textarea bind:value={formUseCase} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Impact</label>
					<textarea bind:value={formImpact} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Architecture</label>
					<textarea bind:value={formArchitecture} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Metrics</label>
					<textarea bind:value={formMetrics} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Project ID</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-group">
					<label>Case Study ID</label>
					<input type="text" bind:value={formCaseStudyId} />
				</div>
				<div class="form-group">
					<label>Demo URL</label>
					<input type="url" bind:value={formDemoUrl} />
				</div>
				<div class="form-group">
					<label>Documentation URL</label>
					<input type="url" bind:value={formDocumentationUrl} />
				</div>
				<div class="form-group">
					<label>GitHub URL</label>
					<input type="url" bind:value={formGithubUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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
{#if showEditModal && editingIntegration}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit AI/ML Integration</h2>
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
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Framework</label>
					<select bind:value={formFramework}>
						{#each frameworks as framework}
							<option value={framework}>{framework}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Model Name</label>
					<input type="text" bind:value={formModelName} />
				</div>
				<div class="form-group">
					<label>Model Version</label>
					<input type="text" bind:value={formModelVersion} />
				</div>
				<div class="form-group">
					<label>Use Case</label>
					<textarea bind:value={formUseCase} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Impact</label>
					<textarea bind:value={formImpact} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Architecture</label>
					<textarea bind:value={formArchitecture} rows="4"></textarea>
				</div>
				<div class="form-group">
					<label>Metrics</label>
					<textarea bind:value={formMetrics} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Project ID</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-group">
					<label>Case Study ID</label>
					<input type="text" bind:value={formCaseStudyId} />
				</div>
				<div class="form-group">
					<label>Demo URL</label>
					<input type="url" bind:value={formDemoUrl} />
				</div>
				<div class="form-group">
					<label>Documentation URL</label>
					<input type="url" bind:value={formDocumentationUrl} />
				</div>
				<div class="form-group">
					<label>GitHub URL</label>
					<input type="url" bind:value={formGithubUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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

