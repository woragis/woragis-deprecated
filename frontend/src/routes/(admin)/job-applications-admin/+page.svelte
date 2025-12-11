<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listJobApplications,
		createJobApplication,
		deleteJobApplication,
		type JobApplication,
		type CreateJobApplicationInput
	} from '$lib/api/jobapplications';
	import { useTranslation } from '$lib/i18n';
	import PageHeader from './_sections/PageHeader.svelte';
	import SearchBar from './_sections/SearchBar.svelte';
	import AdvancedFilters from './_sections/AdvancedFilters.svelte';
	import ApplicationsTable from './_sections/ApplicationsTable.svelte';
	import CreateApplicationModal from './_components/CreateApplicationModal.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import { showToast } from '$lib/utils/toast';
	import { updateJobApplicationStatus, type ApplicationStatus } from '$lib/api/jobapplications';

	const tFn = useTranslation();
	
	let applications: JobApplication[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');
	let showCreateModal = $state(false);
	let filters = $state({
		status: '',
		website: '',
		interestLevel: ''
	});
	let sortBy = $state<'date' | 'company' | 'status' | 'interest'>('date');
	let sortOrder = $state<'asc' | 'desc'>('desc');
	let selectedApplications = $state<Set<string>>(new Set());
	let showBatchActions = $state(false);
	let quickFilter = $state<string>('');

	onMount(async () => {
		await fetchApplications();
	});

	async function fetchApplications() {
		loading = true;
		error = null;
		try {
			applications = await listJobApplications();
		} catch (err) {
			error = err instanceof Error ? err.message : tFn('jobApplications.error');
			console.error('Error fetching job applications:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		showCreateModal = true;
	}

	async function handleCreate(input: CreateJobApplicationInput) {
		try {
			await createJobApplication(input);
			await fetchApplications();
			showToast('Job application created successfully!', 'success');
		} catch (err) {
			const message = err instanceof Error ? err.message : tFn('jobApplications.createError');
			showToast(message, 'error');
			console.error('Error creating job application:', err);
			throw err;
		}
	}

	async function handleDelete(id: string) {
		if (!confirm(tFn('jobApplications.deleteConfirm'))) return;

		try {
			await deleteJobApplication(id);
			await fetchApplications();
			showToast('Application deleted successfully', 'success');
			selectedApplications.delete(id);
		} catch (err) {
			const message = err instanceof Error ? err.message : tFn('jobApplications.deleteError');
			showToast(message, 'error');
			console.error('Error deleting job application:', err);
		}
	}

	async function handleBatchDelete() {
		if (selectedApplications.size === 0) return;
		if (!confirm(`Delete ${selectedApplications.size} application(s)?`)) return;

		try {
			const deletePromises = Array.from(selectedApplications).map(id => deleteJobApplication(id));
			await Promise.all(deletePromises);
			await fetchApplications();
			showToast(`Deleted ${selectedApplications.size} application(s)`, 'success');
			selectedApplications.clear();
			showBatchActions = false;
		} catch (err) {
			showToast('Error deleting applications', 'error');
			console.error('Error batch deleting:', err);
		}
	}

	async function handleBatchStatusUpdate(status: ApplicationStatus) {
		if (selectedApplications.size === 0) return;

		try {
			const updatePromises = Array.from(selectedApplications).map(id =>
				updateJobApplicationStatus(id, { status })
			);
			await Promise.all(updatePromises);
			await fetchApplications();
			showToast(`Updated ${selectedApplications.size} application(s) to ${status}`, 'success');
			selectedApplications.clear();
			showBatchActions = false;
		} catch (err) {
			showToast('Error updating applications', 'error');
			console.error('Error batch updating:', err);
		}
	}

	function toggleSelection(id: string) {
		if (selectedApplications.has(id)) {
			selectedApplications.delete(id);
		} else {
			selectedApplications.add(id);
		}
		selectedApplications = new Set(selectedApplications);
		showBatchActions = selectedApplications.size > 0;
	}

	function toggleSelectAll() {
		const filtered = filteredApplications();
		if (selectedApplications.size === filtered.length) {
			selectedApplications.clear();
		} else {
			selectedApplications = new Set(filtered.map(a => a.id));
		}
		showBatchActions = selectedApplications.size > 0;
	}

	function exportToCSV() {
		const filtered = filteredApplications();
		if (filtered.length === 0) {
			showToast('No applications to export', 'warning');
			return;
		}

		// CSV headers
		const headers = [
			'Company',
			'Job Title',
			'Location',
			'Website',
			'Status',
			'Interest Level',
			'Tags',
			'Follow-up Date',
			'Applied At',
			'Created At',
			'Job URL'
		];

		// CSV rows
		const rows = filtered.map(app => [
			app.companyName,
			app.jobTitle,
			app.location || '',
			app.website,
			app.status,
			app.interestLevel || '',
			app.tags?.join(', ') || '',
			app.followUpDate ? new Date(app.followUpDate).toLocaleDateString() : '',
			app.appliedAt ? new Date(app.appliedAt).toLocaleDateString() : '',
			app.createdAt ? new Date(app.createdAt).toLocaleDateString() : '',
			app.jobUrl
		]);

		// Escape CSV values
		function escapeCSV(value: string): string {
			if (value.includes(',') || value.includes('"') || value.includes('\n')) {
				return `"${value.replace(/"/g, '""')}"`;
			}
			return value;
		}

		// Build CSV content
		const csvContent = [
			headers.map(escapeCSV).join(','),
			...rows.map(row => row.map(cell => escapeCSV(String(cell || ''))).join(','))
		].join('\n');

		// Create download
		const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
		const link = document.createElement('a');
		const url = URL.createObjectURL(blob);
		link.setAttribute('href', url);
		link.setAttribute('download', `job-applications-${new Date().toISOString().split('T')[0]}.csv`);
		link.style.visibility = 'hidden';
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
		URL.revokeObjectURL(url);

		showToast(`Exported ${filtered.length} application(s) to CSV`, 'success');
	}

	function filteredApplications() {
		let filtered = applications;

		// Apply quick filters
		if (quickFilter === 'needs-follow-up') {
			const today = new Date();
			today.setHours(0, 0, 0, 0);
			filtered = filtered.filter(a => {
				if (!a.followUpDate) return false;
				const followUp = new Date(a.followUpDate);
				followUp.setHours(0, 0, 0, 0);
				return followUp <= today && a.status !== 'rejected' && a.status !== 'accepted';
			});
		} else if (quickFilter === 'no-response') {
			const twoWeeksAgo = new Date();
			twoWeeksAgo.setDate(twoWeeksAgo.getDate() - 14);
			filtered = filtered.filter(a => {
				if (!a.appliedAt) return false;
				const applied = new Date(a.appliedAt);
				return applied <= twoWeeksAgo && 
				       a.status === 'applied' && 
				       !a.responseReceivedAt;
			});
		} else if (quickFilter === 'high-interest') {
			filtered = filtered.filter(a => a.interestLevel === 'high' || a.interestLevel === 'very-high');
		}

		// Apply text search
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			filtered = filtered.filter(
				(a) =>
					a.companyName.toLowerCase().includes(query) ||
					a.jobTitle.toLowerCase().includes(query) ||
					a.location?.toLowerCase().includes(query)
			);
		}

		// Apply status filter
		if (filters.status) {
			filtered = filtered.filter(a => a.status === filters.status);
		}

		// Apply website filter
		if (filters.website) {
			filtered = filtered.filter(a => a.website.toLowerCase() === filters.website.toLowerCase());
		}

		// Apply interest level filter
		if (filters.interestLevel) {
			filtered = filtered.filter(a => a.interestLevel === filters.interestLevel);
		}

		// Apply sorting
		filtered = [...filtered].sort((a, b) => {
			let comparison = 0;
			
			switch (sortBy) {
				case 'date':
					const dateA = new Date(a.createdAt).getTime();
					const dateB = new Date(b.createdAt).getTime();
					comparison = dateA - dateB;
					break;
				case 'company':
					comparison = a.companyName.localeCompare(b.companyName);
					break;
				case 'status':
					comparison = a.status.localeCompare(b.status);
					break;
				case 'interest':
					const interestOrder = { 'very-high': 4, 'high': 3, 'medium': 2, 'low': 1, '': 0 };
					const aInterest = interestOrder[a.interestLevel as keyof typeof interestOrder] || 0;
					const bInterest = interestOrder[b.interestLevel as keyof typeof interestOrder] || 0;
					comparison = aInterest - bInterest;
					break;
			}
			
			return sortOrder === 'asc' ? comparison : -comparison;
		});

		return filtered;
	}
</script>

<div class="page-container">
	<PageHeader onCreateClick={openCreateModal} />

	<SearchBar bind:searchQuery />
	
	<div class="quick-filters">
		<button 
			class="quick-filter-btn" 
			class:active={quickFilter === 'needs-follow-up'}
			onclick={() => quickFilter = quickFilter === 'needs-follow-up' ? '' : 'needs-follow-up'}
		>
			Needs Follow-up
		</button>
		<button 
			class="quick-filter-btn"
			class:active={quickFilter === 'no-response'}
			onclick={() => quickFilter = quickFilter === 'no-response' ? '' : 'no-response'}
		>
			No Response (2+ weeks)
		</button>
		<button 
			class="quick-filter-btn"
			class:active={quickFilter === 'high-interest'}
			onclick={() => quickFilter = quickFilter === 'high-interest' ? '' : 'high-interest'}
		>
			High Interest
		</button>
		<button 
			class="quick-filter-btn export-btn"
			onclick={exportToCSV}
		>
			Export CSV
		</button>
	</div>
	
	<AdvancedFilters bind:filters />

	{#if showBatchActions}
		<div class="batch-actions">
			<span class="batch-count">{selectedApplications.size} selected</span>
			<div class="batch-buttons">
				<button class="batch-btn" onclick={() => handleBatchStatusUpdate('applied')}>
					Mark as Applied
				</button>
				<button class="batch-btn" onclick={() => handleBatchStatusUpdate('contacted')}>
					Mark as Contacted
				</button>
				<button class="batch-btn" onclick={() => handleBatchStatusUpdate('rejected')}>
					Mark as Rejected
				</button>
				<button class="batch-btn danger" onclick={handleBatchDelete}>
					Delete Selected
				</button>
				<button class="batch-btn" onclick={() => { selectedApplications.clear(); showBatchActions = false; }}>
					Cancel
				</button>
			</div>
		</div>
	{/if}

	<div class="sort-controls">
		<label>
			Sort by:
			<select bind:value={sortBy}>
				<option value="date">Date</option>
				<option value="company">Company</option>
				<option value="status">Status</option>
				<option value="interest">Interest Level</option>
			</select>
		</label>
		<label>
			Order:
			<select bind:value={sortOrder}>
				<option value="desc">Descending</option>
				<option value="asc">Ascending</option>
			</select>
		</label>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchApplications} />
	{:else if loading}
		<LoadingState message={tFn('jobApplications.loading')} />
	{:else}
		<ApplicationsTable 
			applications={filteredApplications()} 
			onDelete={handleDelete}
			bind:selectedApplications
			onToggleSelection={toggleSelection}
			onToggleSelectAll={toggleSelectAll}
		/>
	{/if}
</div>

<CreateApplicationModal 
	bind:open={showCreateModal} 
	onSubmit={handleCreate}
	existingApplications={applications}
/>

<ToastContainer />

<style>
	.page-container {
		padding: var(--spacing-md);
		max-width: 1400px;
		margin: 0 auto;
	}

	.batch-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		margin-bottom: var(--spacing-md);
	}

	.batch-count {
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.batch-buttons {
		display: flex;
		gap: var(--spacing-sm);
	}

	.batch-btn {
		padding: var(--spacing-xs) var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		cursor: pointer;
		font-size: var(--font-size-sm);
		transition: background-color var(--transition-base);
	}

	.batch-btn:hover {
		background-color: var(--color-bg-hover);
	}

	.batch-btn.danger {
		background-color: var(--color-danger, #ef4444);
		color: white;
		border-color: var(--color-danger, #ef4444);
	}

	.batch-btn.danger:hover {
		background-color: var(--color-danger-hover, #dc2626);
	}

	.sort-controls {
		display: flex;
		gap: var(--spacing-md);
		margin-bottom: var(--spacing-md);
		align-items: center;
	}

	.sort-controls label {
		display: flex;
		align-items: center;
		gap: var(--spacing-xs);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
	}

	.sort-controls select {
		padding: var(--spacing-xs) var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		font-size: var(--font-size-sm);
	}

	.quick-filters {
		display: flex;
		gap: var(--spacing-sm);
		margin-bottom: var(--spacing-md);
		flex-wrap: wrap;
	}

	.quick-filter-btn {
		padding: var(--spacing-xs) var(--spacing-md);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
		cursor: pointer;
		font-size: var(--font-size-sm);
		transition: all var(--transition-base);
	}

	.quick-filter-btn:hover {
		background-color: var(--color-bg-hover);
	}

	.quick-filter-btn.active {
		background-color: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}

	.quick-filter-btn.export-btn {
		background-color: var(--color-success, #10b981);
		color: white;
		border-color: var(--color-success, #10b981);
		margin-left: auto;
	}

	.quick-filter-btn.export-btn:hover {
		background-color: var(--color-success-hover, #059669);
	}
</style>
