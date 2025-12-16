<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		listResumes,
		deleteResume,
		markAsMain,
		unmarkAsMain,
		markAsFeatured,
		unmarkAsFeatured,
		listResumeTags,
		type Resume
	} from '$lib/api/resumes';
	import PageHeader from './_sections/PageHeader.svelte';
	import SearchBar from './_sections/SearchBar.svelte';
	import ResumeFilters from './_sections/ResumeFilters.svelte';
	import ResumesTable from './_sections/ResumesTable.svelte';
	import ResumeMetrics from './_sections/ResumeMetrics.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import ConfirmationModal from '$lib/components/ui/ConfirmationModal.svelte';
	import { toastSuccess, toastError, getApiErrorMessage } from '$lib/utils/toast';

	let resumes: Resume[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');
	let filters = $state({
		tags: [] as string[],
		main: false,
		featured: false
	});
	let availableTags: string[] = $state([]);
	let sortBy = $state<'date' | 'title' | 'usage' | 'interviewRate' | 'offerRate'>('date');
	let sortOrder = $state<'asc' | 'desc'>('desc');
	let selectedResumes = $state<Set<string>>(new Set());
	let showBatchActions = $state(false);
	let showDeleteConfirm = $state(false);
	let deleteTargetId: string | null = $state(null);
	let showBatchDeleteConfirm = $state(false);

	onMount(async () => {
		await Promise.all([fetchResumes(), fetchTags()]);
	});

	async function fetchResumes() {
		loading = true;
		error = null;
		try {
			// Always fetch all resumes, filtering happens client-side
			resumes = await listResumes();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load resumes';
			console.error('Error fetching resumes:', err);
		} finally {
			loading = false;
		}
	}

	async function fetchTags() {
		try {
			availableTags = await listResumeTags();
		} catch (err) {
			console.error('Error fetching tags:', err);
		}
	}

	function openUploadPage() {
		goto('/resumes/upload');
	}

	function openGeneratePage() {
		goto('/resumes/generate');
	}

	function handleView(resume: Resume) {
		goto(`/resumes/${resume.id}`);
	}

	function handleDelete(id: string) {
		deleteTargetId = id;
		showDeleteConfirm = true;
	}

	async function confirmDelete() {
		if (!deleteTargetId) return;
		try {
			await deleteResume(deleteTargetId);
			await fetchResumes();
			toastSuccess('Resume deleted successfully');
			selectedResumes.delete(deleteTargetId);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to delete resume');
			toastError(message);
			console.error('Error deleting resume:', err);
		} finally {
			deleteTargetId = null;
			showDeleteConfirm = false;
		}
	}

	async function handleMarkAsMain(id: string) {
		try {
			await markAsMain(id);
			await fetchResumes();
			toastSuccess('Resume marked as main');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to mark resume as main');
			toastError(message);
			console.error('Error marking resume as main:', err);
		}
	}

	async function handleUnmarkAsMain(id: string) {
		try {
			await unmarkAsMain(id);
			await fetchResumes();
			toastSuccess('Resume unmarked as main');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to unmark resume as main');
			toastError(message);
			console.error('Error unmarking resume as main:', err);
		}
	}

	async function handleMarkAsFeatured(id: string) {
		try {
			await markAsFeatured(id);
			await fetchResumes();
			toastSuccess('Resume marked as featured');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to mark resume as featured');
			toastError(message);
			console.error('Error marking resume as featured:', err);
		}
	}

	async function handleUnmarkAsFeatured(id: string) {
		try {
			await unmarkAsFeatured(id);
			await fetchResumes();
			toastSuccess('Resume unmarked as featured');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to unmark resume as featured');
			toastError(message);
			console.error('Error unmarking resume as featured:', err);
		}
	}

	function handleBatchDelete() {
		if (selectedResumes.size === 0) return;
		showBatchDeleteConfirm = true;
	}

	async function confirmBatchDelete() {
		if (selectedResumes.size === 0) return;
		try {
			const deletePromises = Array.from(selectedResumes).map(id => deleteResume(id));
			await Promise.all(deletePromises);
			await fetchResumes();
			toastSuccess(`Deleted ${selectedResumes.size} resume(s)`);
			selectedResumes.clear();
			showBatchActions = false;
		} catch (err) {
			toastError(getApiErrorMessage(err, 'Error deleting resumes'));
			console.error('Error batch deleting:', err);
		} finally {
			showBatchDeleteConfirm = false;
		}
	}

	function toggleSelection(id: string) {
		if (selectedResumes.has(id)) {
			selectedResumes.delete(id);
		} else {
			selectedResumes.add(id);
		}
		selectedResumes = new Set(selectedResumes);
		showBatchActions = selectedResumes.size > 0;
	}

	function toggleSelectAll() {
		const filtered = filteredResumes();
		if (selectedResumes.size === filtered.length) {
			selectedResumes.clear();
		} else {
			selectedResumes = new Set(filtered.map(r => r.id));
		}
		showBatchActions = selectedResumes.size > 0;
	}

	function filteredResumes() {
		let filtered = resumes;

		// Search filter
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			filtered = filtered.filter(
				(r) =>
					r.title.toLowerCase().includes(query) ||
					r.fileName.toLowerCase().includes(query) ||
					r.tags?.some(tag => tag.toLowerCase().includes(query))
			);
		}

		// Tag filter
		if (filters.tags.length > 0) {
			filtered = filtered.filter(r => 
				r.tags && filters.tags.some(tag => r.tags!.includes(tag))
			);
		}

		// Main filter
		if (filters.main) {
			filtered = filtered.filter(r => r.isMain);
		}

		// Featured filter
		if (filters.featured) {
			filtered = filtered.filter(r => r.isFeatured);
		}

		// Sort
		filtered = [...filtered].sort((a, b) => {
			let comparison = 0;
			
			switch (sortBy) {
				case 'date':
					const dateA = new Date(a.createdAt).getTime();
					const dateB = new Date(b.createdAt).getTime();
					comparison = dateA - dateB;
					break;
				case 'title':
					comparison = a.title.localeCompare(b.title);
					break;
				case 'usage':
					comparison = a.applicationsUsed - b.applicationsUsed;
					break;
				case 'interviewRate':
					comparison = a.interviewRate - b.interviewRate;
					break;
				case 'offerRate':
					comparison = a.offerRate - b.offerRate;
					break;
			}
			
			return sortOrder === 'asc' ? comparison : -comparison;
		});

		return filtered;
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<PageHeader 
		onUploadClick={openUploadPage} 
		onGenerateClick={openGeneratePage} 
	/>

	<ResumeMetrics resumes={resumes} />

	<SearchBar bind:searchQuery />
	
	<ResumeFilters 
		bind:filters 
		availableTags={availableTags}
	/>

	{#if showBatchActions}
		<div class="batch-actions">
			<span class="batch-count">{selectedResumes.size} selected</span>
			<div class="batch-buttons">
				<button class="batch-btn danger" onclick={handleBatchDelete}>
					Delete Selected
				</button>
				<button class="batch-btn" onclick={() => { selectedResumes.clear(); showBatchActions = false; }}>
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
				<option value="title">Title</option>
				<option value="usage">Usage</option>
				<option value="interviewRate">Interview Rate</option>
				<option value="offerRate">Offer Rate</option>
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
		<ErrorState message={error} onRetry={fetchResumes} />
	{:else if loading}
		<LoadingState message="Loading resumes..." />
	{:else}
		<ResumesTable 
			resumes={filteredResumes()} 
			onDelete={handleDelete}
			onMarkAsMain={handleMarkAsMain}
			onUnmarkAsMain={handleUnmarkAsMain}
			onMarkAsFeatured={handleMarkAsFeatured}
			onUnmarkAsFeatured={handleUnmarkAsFeatured}
			bind:selectedResumes
			onToggleSelection={toggleSelection}
			onToggleSelectAll={toggleSelectAll}
			onView={handleView}
		/>
	{/if}
</div>

<ConfirmationModal
	bind:open={showDeleteConfirm}
	title="Delete Resume"
	message="Are you sure you want to delete this resume? This action cannot be undone."
	onConfirm={confirmDelete}
/>

<ConfirmationModal
	bind:open={showBatchDeleteConfirm}
	title="Delete Selected Resumes"
	message="Are you sure you want to delete {selectedResumes.size} resume(s)? This action cannot be undone."
	onConfirm={confirmBatchDelete}
/>

<ToastContainer />

<style>
	.batch-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background-color: #f3f4f6;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		margin-bottom: 1rem;
	}

	.batch-count {
		font-weight: 600;
		color: #1f2937;
	}

	.batch-buttons {
		display: flex;
		gap: 0.5rem;
	}

	.batch-btn {
		padding: 0.5rem 1rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		background-color: white;
		color: #1f2937;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s;
	}

	.batch-btn:hover {
		background-color: #f9fafb;
	}

	.batch-btn.danger {
		background-color: #ef4444;
		color: white;
		border-color: #ef4444;
	}

	.batch-btn.danger:hover {
		background-color: #dc2626;
	}

	.sort-controls {
		display: flex;
		gap: 1rem;
		margin-bottom: 1rem;
		align-items: center;
	}

	.sort-controls label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		color: #1f2937;
	}

	.sort-controls select {
		padding: 0.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		background-color: white;
		color: #1f2937;
		font-size: 0.875rem;
	}
</style>
