<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listImpactMetrics,
		createImpactMetric,
		updateImpactMetric,
		deleteImpactMetric,
		type ImpactMetric,
		type CreateImpactMetricInput,
		type UpdateImpactMetricInput,
		type MetricType,
		type MetricUnit,
		type EntityType
	} from '$lib/api/impactmetrics';

	let metrics: ImpactMetric[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingMetric: ImpactMetric | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formType = $state<MetricType>('projects_delivered');
	let formValue = $state<number | ''>('');
	let formUnit = $state<MetricUnit>('count');
	let formDescription = $state('');
	let formEntityType = $state<EntityType | ''>('');
	let formEntityId = $state('');
	let formPeriodStart = $state('');
	let formPeriodEnd = $state('');
	let formFeatured = $state(false);
	let formDisplayOrder = $state<number | ''>(0);

	const types: MetricType[] = [
		'projects_delivered',
		'users_impacted',
		'performance_improvement',
		'cost_savings',
		'time_saved'
	];
	const units: MetricUnit[] = [
		'count',
		'percentage',
		'currency',
		'hours',
		'days',
		'months',
		'years',
		'milliseconds',
		'seconds',
		'minutes'
	];
	const entityTypes: EntityType[] = ['project', 'problem_solution', 'case_study', 'system_design'];

	onMount(async () => {
		await fetchMetrics();
	});

	async function fetchMetrics() {
		loading = true;
		error = null;
		try {
			metrics = await listImpactMetrics();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load impact metrics';
			console.error('Error fetching impact metrics:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(metric: ImpactMetric) {
		editingMetric = metric;
		formType = metric.type;
		formValue = metric.value;
		formUnit = metric.unit;
		formDescription = metric.description || '';
		formEntityType = metric.entityType || '';
		formEntityId = metric.entityId || '';
		formPeriodStart = metric.periodStart ? metric.periodStart.split('T')[0] : '';
		formPeriodEnd = metric.periodEnd ? metric.periodEnd.split('T')[0] : '';
		formFeatured = metric.featured;
		formDisplayOrder = metric.displayOrder;
		showEditModal = true;
	}

	function resetForm() {
		formType = 'projects_delivered';
		formValue = '';
		formUnit = 'count';
		formDescription = '';
		formEntityType = '';
		formEntityId = '';
		formPeriodStart = '';
		formPeriodEnd = '';
		formFeatured = false;
		formDisplayOrder = 0;
		editingMetric = null;
	}

	async function handleCreate() {
		if (formValue === '' || formValue === null) {
			alert('Value is required');
			return;
		}

		try {
			const input: CreateImpactMetricInput = {
				type: formType,
				value: Number(formValue),
				unit: formUnit,
				description: formDescription.trim() || undefined,
				entityType: formEntityType || undefined,
				entityId: formEntityId.trim() || undefined,
				periodStart: formPeriodStart || undefined,
				periodEnd: formPeriodEnd || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await createImpactMetric(input);
			showCreateModal = false;
			resetForm();
			await fetchMetrics();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create impact metric');
			console.error('Error creating impact metric:', err);
		}
	}

	async function handleUpdate() {
		if (!editingMetric || formValue === '' || formValue === null) {
			alert('Value is required');
			return;
		}

		try {
			const input: UpdateImpactMetricInput = {
				type: formType,
				value: Number(formValue),
				unit: formUnit,
				description: formDescription.trim() || undefined,
				entityType: formEntityType || undefined,
				entityId: formEntityId.trim() || undefined,
				periodStart: formPeriodStart || undefined,
				periodEnd: formPeriodEnd || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await updateImpactMetric(editingMetric.id, input);
			showEditModal = false;
			resetForm();
			await fetchMetrics();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update impact metric');
			console.error('Error updating impact metric:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this impact metric?')) return;

		try {
			await deleteImpactMetric(id);
			await fetchMetrics();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete impact metric');
			console.error('Error deleting impact metric:', err);
		}
	}

	function filteredMetrics() {
		if (!searchQuery.trim()) return metrics;
		const query = searchQuery.toLowerCase();
		return metrics.filter(
			(m) =>
				m.type.toLowerCase().includes(query) ||
				m.description?.toLowerCase().includes(query) ||
				m.unit.toLowerCase().includes(query)
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
			<h1>Impact Metrics Management</h1>
			<p>Manage impact metrics</p>
		</div>
		<button onclick={openCreateModal}>Create Metric</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search metrics..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredMetrics().length === 0}
		<div class="empty">No impact metrics found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Type</th>
					<th>Value</th>
					<th>Unit</th>
					<th>Description</th>
					<th>Entity</th>
					<th>Featured</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredMetrics() as metric}
					<tr>
						<td>{metric.type}</td>
						<td><strong>{metric.value}</strong></td>
						<td>{metric.unit}</td>
						<td>{metric.description || '—'}</td>
						<td>
							{#if metric.entityType && metric.entityId}
								{metric.entityType}: {metric.entityId.substring(0, 8)}...
							{:else}
								—
							{/if}
						</td>
						<td>{metric.featured ? 'Yes' : 'No'}</td>
						<td>{formatDate(metric.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(metric)}>Edit</button>
							<button onclick={() => handleDelete(metric.id)} class="delete-btn">Delete</button>
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
			<h2>Create Impact Metric</h2>
			<div class="form">
				<div class="form-group">
					<label>Type *</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Value *</label>
					<input type="number" step="0.01" bind:value={formValue} />
				</div>
				<div class="form-group">
					<label>Unit *</label>
					<select bind:value={formUnit}>
						{#each units as unit}
							<option value={unit}>{unit}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Entity Type</label>
					<select bind:value={formEntityType}>
						<option value="">—</option>
						{#each entityTypes as et}
							<option value={et}>{et}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Entity ID</label>
					<input type="text" bind:value={formEntityId} />
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
{#if showEditModal && editingMetric}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Impact Metric</h2>
			<div class="form">
				<div class="form-group">
					<label>Type *</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Value *</label>
					<input type="number" step="0.01" bind:value={formValue} />
				</div>
				<div class="form-group">
					<label>Unit *</label>
					<select bind:value={formUnit}>
						{#each units as unit}
							<option value={unit}>{unit}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Entity Type</label>
					<select bind:value={formEntityType}>
						<option value="">—</option>
						{#each entityTypes as et}
							<option value={et}>{et}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Entity ID</label>
					<input type="text" bind:value={formEntityId} />
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

