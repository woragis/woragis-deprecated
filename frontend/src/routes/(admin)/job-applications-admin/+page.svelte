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
		} catch (err) {
			alert(err instanceof Error ? err.message : tFn('jobApplications.createError'));
			console.error('Error creating job application:', err);
			throw err;
		}
	}

	async function handleDelete(id: string) {
		if (!confirm(tFn('jobApplications.deleteConfirm'))) return;

		try {
			await deleteJobApplication(id);
			await fetchApplications();
		} catch (err) {
			alert(err instanceof Error ? err.message : tFn('jobApplications.deleteError'));
			console.error('Error deleting job application:', err);
		}
	}

	function filteredApplications() {
		let filtered = applications;

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

		return filtered;
	}
</script>

<div class="page-container">
	<PageHeader onCreateClick={openCreateModal} />

	<SearchBar bind:searchQuery />
	
	<AdvancedFilters bind:filters />

	{#if error}
		<ErrorState message={error} onRetry={fetchApplications} />
	{:else if loading}
		<LoadingState message={tFn('jobApplications.loading')} />
	{:else}
		<ApplicationsTable applications={filteredApplications()} onDelete={handleDelete} />
	{/if}
</div>

<CreateApplicationModal 
	bind:open={showCreateModal} 
	onSubmit={handleCreate}
	existingApplications={applications}
/>

<style>
	.page-container {
		padding: var(--spacing-md);
		max-width: 1400px;
		margin: 0 auto;
	}
</style>
