<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listJobApplications,
		createJobApplication,
		updateJobApplicationStatus,
		deleteJobApplication,
		type JobApplication,
		type CreateJobApplicationInput,
		type ApplicationStatus
	} from '$lib/api/jobapplications';
	import { listResumes, type Resume } from '$lib/api/resumes';
	import { locale, t } from '$lib/i18n';

	let applications: JobApplication[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');

	// Modal states
	let showCreateModal = $state(false);

	// Form states
	let formCompanyName = $state('');
	let formLocation = $state('');
	let formJobTitle = $state('');
	let formJobUrl = $state('');
	let formWebsite = $state('');
	let formCoverLetter = $state('');
	let formLinkedInContact = $state(false);
	let formStatus = $state<ApplicationStatus>('pending');

	const statuses: ApplicationStatus[] = [
		'pending',
		'processing',
		'applied',
		'contacted',
		'rejected',
		'accepted',
		'failed'
	];

	onMount(async () => {
		await fetchApplications();
	});

	async function fetchApplications() {
		loading = true;
		error = null;
		try {
			applications = await listJobApplications();
		} catch (err) {
			error = err instanceof Error ? err.message : $t('jobApplications.error');
			console.error('Error fetching job applications:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function resetForm() {
		formCompanyName = '';
		formLocation = '';
		formJobTitle = '';
		formJobUrl = '';
		formWebsite = '';
		formCoverLetter = '';
		formLinkedInContact = false;
		formStatus = 'pending';
	}

	async function handleCreate() {
		if (!formCompanyName.trim() || !formJobTitle.trim() || !formJobUrl.trim() || !formWebsite.trim()) {
			alert($t('jobApplications.modal.required') + ' ' + $t('jobApplications.modal.companyName') + ', ' + $t('jobApplications.modal.jobTitle') + ', ' + $t('jobApplications.modal.jobUrl') + ', ' + $t('jobApplications.modal.website'));
			return;
		}

		try {
			const input: CreateJobApplicationInput = {
				companyName: formCompanyName.trim(),
				location: formLocation.trim() || undefined,
				jobTitle: formJobTitle.trim(),
				jobUrl: formJobUrl.trim(),
				website: formWebsite.trim(),
				coverLetter: formCoverLetter.trim() || undefined,
				linkedInContact: formLinkedInContact,
				status: formStatus
			};

			await createJobApplication(input);
			showCreateModal = false;
			resetForm();
			await fetchApplications();
		} catch (err) {
			alert(err instanceof Error ? err.message : $t('jobApplications.createError'));
			console.error('Error creating job application:', err);
		}
	}


	async function handleDelete(id: string) {
		if (!confirm($t('jobApplications.deleteConfirm'))) return;

		try {
			await deleteJobApplication(id);
			await fetchApplications();
		} catch (err) {
			alert(err instanceof Error ? err.message : $t('jobApplications.deleteError'));
			console.error('Error deleting job application:', err);
		}
	}

	function filteredApplications() {
		if (!searchQuery.trim()) return applications;
		const query = searchQuery.toLowerCase();
		return applications.filter(
			(a) =>
				a.companyName.toLowerCase().includes(query) ||
				a.jobTitle.toLowerCase().includes(query) ||
				a.location?.toLowerCase().includes(query)
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
			<h1>{$t('jobApplications.title')}</h1>
			<p>{$t('jobApplications.subtitle')}</p>
		</div>
		<button onclick={openCreateModal}>{$t('jobApplications.createButton')}</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder={$t('jobApplications.searchPlaceholder')}
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">{$t('jobApplications.loading')}</div>
	{:else if filteredApplications().length === 0}
		<div class="empty">{$t('jobApplications.empty')}</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>{$t('jobApplications.table.company')}</th>
					<th>{$t('jobApplications.table.jobTitle')}</th>
					<th>{$t('jobApplications.table.location')}</th>
					<th>{$t('jobApplications.table.website')}</th>
					<th>{$t('jobApplications.table.language')}</th>
					<th>{$t('jobApplications.table.status')}</th>
					<th>{$t('jobApplications.table.interest')}</th>
					<th>{$t('jobApplications.table.appliedAt')}</th>
					<th>{$t('jobApplications.table.actions')}</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredApplications() as app}
					<tr>
						<td><strong>{app.companyName}</strong></td>
						<td>{app.jobTitle}</td>
						<td>{app.location || '—'}</td>
						<td>{app.website}</td>
						<td>{app.language || '—'}</td>
						<td>
							<span class="status status-{app.status}">{$t(`jobApplications.status.${app.status}` as any)}</span>
						</td>
						<td>{app.interestLevel || '—'}</td>
						<td>{formatDate(app.appliedAt)}</td>
						<td>
							<a href="/job-applications-admin/{app.id}" class="view-link">{$t('jobApplications.table.view')}</a>
							<button onclick={() => handleDelete(app.id)} class="delete-btn">{$t('jobApplications.table.delete')}</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- Create Application Modal -->
{#if showCreateModal}
	<div class="modal-overlay" onclick={() => {
		showCreateModal = false;
		resetForm();
	}}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>{$t('jobApplications.modal.createTitle')}</h2>
			<div class="form">
				<div class="form-row">
					<div class="form-group">
						<label>{$t('jobApplications.modal.companyName')} {$t('jobApplications.modal.required')}</label>
						<input type="text" bind:value={formCompanyName} />
					</div>
					<div class="form-group">
						<label>{$t('jobApplications.modal.jobTitle')} {$t('jobApplications.modal.required')}</label>
						<input type="text" bind:value={formJobTitle} />
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>{$t('jobApplications.modal.jobUrl')} {$t('jobApplications.modal.required')}</label>
						<input type="url" bind:value={formJobUrl} />
					</div>
					<div class="form-group">
						<label>{$t('jobApplications.modal.website')} {$t('jobApplications.modal.required')}</label>
						<input type="text" bind:value={formWebsite} placeholder="linkedin, glassdoor, etc." />
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>{$t('jobApplications.modal.location')}</label>
						<input type="text" bind:value={formLocation} />
					</div>
					<div class="form-group">
						<label>{$t('jobApplications.modal.status')}</label>
						<select bind:value={formStatus}>
							{#each statuses as status}
								<option value={status}>{$t(`jobApplications.status.${status}` as any)}</option>
							{/each}
						</select>
					</div>
				</div>
				<div class="form-group">
					<label>{$t('jobApplications.modal.coverLetter')}</label>
					<textarea bind:value={formCoverLetter} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formLinkedInContact} />
						{$t('jobApplications.modal.linkedInContact')}
					</label>
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>{$t('jobApplications.modal.create')}</button>
					<button onclick={() => {
						showCreateModal = false;
						resetForm();
					}}>{$t('jobApplications.modal.cancel')}</button>
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
		color: #333;
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
		color: #333;
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

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-pending {
		background: #fff3cd;
		color: #856404;
	}

	.status-processing {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-applied {
		background: #d4edda;
		color: #155724;
	}

	.status-contacted {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-rejected {
		background: #f8d7da;
		color: #721c24;
	}

	.status-accepted {
		background: #d4edda;
		color: #155724;
	}

	.status-failed {
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
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
		color: #333;
	}

	.modal-large {
		max-width: 900px;
	}

	.modal h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
		color: #333;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
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
		color: #333;
	}

	.form-group input,
	.form-group textarea,
	.form-group select {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
		color: #333;
		background: white;
	}

	.form-group input:disabled,
	.form-group select:disabled {
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
