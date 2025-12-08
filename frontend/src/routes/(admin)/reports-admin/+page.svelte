<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listReportDefinitions,
		createReportDefinition,
		updateReportDefinition,
		archiveReportDefinitions,
		restoreReportDefinitions,
		deleteReportDefinitions,
		toggleFavorite,
		type ReportDefinition,
		type CreateReportDefinitionInput
	} from '$lib/api/reports';

	let reports: ReportDefinition[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingReport: ReportDefinition | null = $state(null);
	let searchQuery = $state('');
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formName = $state('');
	let formDescription = $state('');
	let formSections = $state('{}');
	let formFilters = $state('{}');
	let formIsFavorite = $state(false);

	onMount(async () => {
		await fetchReports();
	});

	async function fetchReports() {
		loading = true;
		error = null;
		try {
			reports = await listReportDefinitions();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load reports';
			console.error('Error fetching reports:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(report: ReportDefinition) {
		editingReport = report;
		formName = report.name;
		formDescription = report.description || '';
		formSections = JSON.stringify(report.sections || {}, null, 2);
		formFilters = JSON.stringify(report.filters || {}, null, 2);
		formIsFavorite = report.is_favorite;
		showEditModal = true;
	}

	function resetForm() {
		formName = '';
		formDescription = '';
		formSections = '{}';
		formFilters = '{}';
		formIsFavorite = false;
		editingReport = null;
	}

	async function handleCreate() {
		if (!formName.trim()) {
			alert('Name is required');
			return;
		}

		try {
			let sectionsObj: Record<string, any> = {};
			let filtersObj: Record<string, any> = {};

			try {
				sectionsObj = JSON.parse(formSections);
			} catch {
				alert('Sections must be valid JSON');
				return;
			}

			try {
				filtersObj = JSON.parse(formFilters);
			} catch {
				alert('Filters must be valid JSON');
				return;
			}

			const input: CreateReportDefinitionInput = {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				sections: sectionsObj,
				filters: filtersObj,
				isFavorite: formIsFavorite
			};

			await createReportDefinition(input);
			showCreateModal = false;
			resetForm();
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create report');
			console.error('Error creating report:', err);
		}
	}

	async function handleUpdate() {
		if (!editingReport) {
			alert('No report selected');
			return;
		}

		try {
			let sectionsObj: Record<string, any> = {};
			let filtersObj: Record<string, any> = {};

			try {
				sectionsObj = JSON.parse(formSections);
			} catch {
				alert('Sections must be valid JSON');
				return;
			}

			try {
				filtersObj = JSON.parse(formFilters);
			} catch {
				alert('Filters must be valid JSON');
				return;
			}

			await updateReportDefinition(editingReport.id, {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				sections: sectionsObj,
				filters: filtersObj,
				isFavorite: formIsFavorite
			});
			showEditModal = false;
			resetForm();
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update report');
			console.error('Error updating report:', err);
		}
	}

	function toggleSelect(id: string) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = new Set(selectedIds);
	}

	async function handleBulkArchive() {
		if (selectedIds.size === 0) return;
		try {
			await archiveReportDefinitions(Array.from(selectedIds));
			selectedIds.clear();
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to archive reports');
		}
	}

	async function handleBulkRestore() {
		if (selectedIds.size === 0) return;
		try {
			await restoreReportDefinitions(Array.from(selectedIds));
			selectedIds.clear();
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to restore reports');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} report(s)?`)) return;
		try {
			await deleteReportDefinitions(Array.from(selectedIds));
			selectedIds.clear();
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete reports');
		}
	}

	async function handleToggleFavorite(id: string) {
		try {
			const report = reports.find((r) => r.id === id);
			await toggleFavorite(id, report ? !report.is_favorite : false);
			await fetchReports();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to toggle favorite');
		}
	}

	function filteredReports() {
		if (!searchQuery.trim()) return reports;
		const query = searchQuery.toLowerCase();
		return reports.filter(
			(r) =>
				r.name.toLowerCase().includes(query) ||
				r.description?.toLowerCase().includes(query)
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
			<h1>Reports Management</h1>
			<p>Manage report definitions</p>
		</div>
		<button onclick={openCreateModal}>Create Report</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search reports..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if selectedIds.size > 0}
		<div class="bulk-actions">
			<span>{selectedIds.size} selected</span>
			<button onclick={handleBulkArchive}>Archive</button>
			<button onclick={handleBulkRestore}>Restore</button>
			<button onclick={handleBulkDelete} class="delete-btn">Delete</button>
			<button onclick={() => (selectedIds.clear())}>Clear</button>
		</div>
	{/if}

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredReports().length === 0}
		<div class="empty">No reports found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>
						<input
							type="checkbox"
							checked={selectedIds.size === filteredReports().length && filteredReports().length > 0}
							onchange={() => {
								if (selectedIds.size === filteredReports().length) {
									selectedIds.clear();
								} else {
									selectedIds = new Set(filteredReports().map((r) => r.id));
								}
							}}
						/>
					</th>
					<th>Name</th>
					<th>Description</th>
					<th>Favorite</th>
					<th>Status</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredReports() as report}
					<tr>
						<td>
							<input
								type="checkbox"
								checked={selectedIds.has(report.id)}
								onchange={() => toggleSelect(report.id)}
							/>
						</td>
						<td><strong>{report.name}</strong></td>
						<td>{report.description || '—'}</td>
						<td>
							<button
								onclick={() => handleToggleFavorite(report.id)}
								class="favorite-btn {report.is_favorite ? 'active' : ''}"
							>
								{report.is_favorite ? '★' : '☆'}
							</button>
						</td>
						<td>
							{#if report.deleted_at}
								<span class="status status-deleted">Deleted</span>
							{:else if report.archived_at}
								<span class="status status-archived">Archived</span>
							{:else}
								<span class="status status-active">Active</span>
							{/if}
						</td>
						<td>{formatDate(report.created_at)}</td>
						<td>
							<button onclick={() => openEditModal(report)}>Edit</button>
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
			<h2>Create Report</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Sections (JSON)</label>
					<textarea bind:value={formSections} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>Filters (JSON)</label>
					<textarea bind:value={formFilters} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formIsFavorite} />
						Favorite
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
{#if showEditModal && editingReport}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Report</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Sections (JSON)</label>
					<textarea bind:value={formSections} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>Filters (JSON)</label>
					<textarea bind:value={formFilters} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formIsFavorite} />
						Favorite
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

	.bulk-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem;
		background: #f5f5f5;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.bulk-actions button {
		padding: 0.375rem 0.75rem;
		background: #6c757d;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.bulk-actions button:hover {
		background: #5a6268;
	}

	.bulk-actions .delete-btn {
		background: #dc3545 !important;
	}

	.bulk-actions .delete-btn:hover {
		background: #c82333 !important;
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

	.favorite-btn {
		background: transparent !important;
		color: #ffc107;
		font-size: 1.25rem;
		padding: 0;
		border: none;
		cursor: pointer;
	}

	.favorite-btn.active {
		color: #ffc107;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-active {
		background: #d4edda;
		color: #155724;
	}

	.status-archived {
		background: #fff3cd;
		color: #856404;
	}

	.status-deleted {
		background: #f8d7da;
		color: #721c24;
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
	.form-group textarea {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
		font-family: monospace;
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

