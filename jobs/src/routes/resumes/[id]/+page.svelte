<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getResume,
		updateResume,
		deleteResume,
		downloadResumeById,
		markAsMain,
		unmarkAsMain,
		markAsFeatured,
		unmarkAsFeatured,
		recalculateResumeMetrics,
		type Resume,
		type UpdateResumeInput
	} from '$lib/api/resumes';
	import { listJobApplications, type JobApplication } from '$lib/api/jobapplications';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import ConfirmationModal from '$lib/components/ui/ConfirmationModal.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import { toastSuccess, toastError, getApiErrorMessage } from '$lib/utils/toast';
	import ResumeInfoSection from './_sections/ResumeInfoSection.svelte';
	import ResumePreviewSection from './_sections/ResumePreviewSection.svelte';
	import ResumeMetricsSection from './_sections/ResumeMetricsSection.svelte';
	import ResumeActionsSection from './_sections/ResumeActionsSection.svelte';
	import TagInput from '$lib/components/ui/TagInput.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let resume: Resume | null = $state(null);
	let associatedApplications: JobApplication[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showDeleteConfirm = $state(false);
	let editing = $state(false);
	let formTitle = $state('');
	let formTags = $state<string[]>([]);
	let saving = $state(false);

	onMount(async () => {
		await fetchResume();
		await fetchAssociatedApplications();
	});

	async function fetchResume() {
		loading = true;
		error = null;
		const resumeId = $page.params.id;
		try {
			const data = await getResume(resumeId);
			if (data) {
				resume = data;
				formTitle = data.title;
				formTags = data.tags || [];
			} else {
				error = 'Resume not found';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load resume';
			console.error('Error fetching resume:', err);
		} finally {
			loading = false;
		}
	}

	async function fetchAssociatedApplications() {
		try {
			const applications = await listJobApplications();
			if (resume) {
				associatedApplications = applications.filter(app => app.resumeId === resume.id);
			}
		} catch (err) {
			console.error('Error fetching applications:', err);
		}
	}

	async function handleUpdate() {
		if (!resume) return;
		saving = true;
		try {
			const input: UpdateResumeInput = {
				title: formTitle,
				tags: formTags
			};
			resume = await updateResume(resume.id, input);
			editing = false;
			toastSuccess('Resume updated successfully');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to update resume');
			toastError(message);
			console.error('Error updating resume:', err);
		} finally {
			saving = false;
		}
	}

	function handleDelete() {
		showDeleteConfirm = true;
	}

	async function confirmDelete() {
		if (!resume) return;
		try {
			await deleteResume(resume.id);
			toastSuccess('Resume deleted successfully');
			goto('/resumes');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to delete resume');
			toastError(message);
			console.error('Error deleting resume:', err);
		} finally {
			showDeleteConfirm = false;
		}
	}

	async function handleDownload() {
		if (!resume) return;
		try {
			await downloadResumeById(resume.id, resume.fileName);
		} catch (err) {
			toastError(getApiErrorMessage(err, 'Failed to download resume'));
			console.error('Error downloading resume:', err);
		}
	}

	async function handleMarkAsMain() {
		if (!resume) return;
		try {
			await markAsMain(resume.id);
			await fetchResume();
			toastSuccess('Resume marked as main');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to mark resume as main');
			toastError(message);
			console.error('Error marking resume as main:', err);
		}
	}

	async function handleUnmarkAsMain() {
		if (!resume) return;
		try {
			await unmarkAsMain(resume.id);
			await fetchResume();
			toastSuccess('Resume unmarked as main');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to unmark resume as main');
			toastError(message);
			console.error('Error unmarking resume as main:', err);
		}
	}

	async function handleMarkAsFeatured() {
		if (!resume) return;
		try {
			await markAsFeatured(resume.id);
			await fetchResume();
			toastSuccess('Resume marked as featured');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to mark resume as featured');
			toastError(message);
			console.error('Error marking resume as featured:', err);
		}
	}

	async function handleUnmarkAsFeatured() {
		if (!resume) return;
		try {
			await unmarkAsFeatured(resume.id);
			await fetchResume();
			toastSuccess('Resume unmarked as featured');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to unmark resume as featured');
			toastError(message);
			console.error('Error unmarking resume as featured:', err);
		}
	}

	async function handleRecalculateMetrics() {
		if (!resume) return;
		try {
			resume = await recalculateResumeMetrics(resume.id);
			toastSuccess('Metrics recalculated successfully');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to recalculate metrics');
			toastError(message);
			console.error('Error recalculating metrics:', err);
		}
	}

	function cancelEdit() {
		if (resume) {
			formTitle = resume.title;
			formTags = resume.tags || [];
		}
		editing = false;
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="page-header">
		<button class="back-button" onclick={() => goto('/resumes')}>
			← Back to Resumes
		</button>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchResume} />
	{:else if loading}
		<LoadingState message="Loading resume..." />
	{:else if resume}
		<div class="resume-detail">
			<div class="detail-header">
				<div class="header-content">
					<div class="title-section">
						{#if editing}
							<Input
								bind:value={formTitle}
								placeholder="Resume title"
								class="title-input"
							/>
						{:else}
							<h1 class="resume-title">{resume.title}</h1>
						{/if}
						<div class="status-badges">
							{#if resume.isMain}
								<StatusBadge status="main" label="Main" />
							{/if}
							{#if resume.isFeatured}
								<StatusBadge status="featured" label="Featured" />
							{/if}
						</div>
					</div>
					<div class="file-info">
						<span class="filename">{resume.fileName}</span>
						<span class="file-size">{(resume.fileSize / 1024).toFixed(2)} KB</span>
					</div>
				</div>
				<div class="header-actions">
					{#if editing}
						<Button onclick={handleUpdate} variant="primary" disabled={saving}>
							{saving ? 'Saving...' : 'Save'}
						</Button>
						<Button onclick={cancelEdit} variant="secondary" disabled={saving}>
							Cancel
						</Button>
					{:else}
						<Button onclick={() => editing = true} variant="secondary">
							Edit
						</Button>
						<Button onclick={handleDownload} variant="primary">
							Download
						</Button>
					{/if}
				</div>
			</div>

			<div class="detail-content">
				<div class="main-content">
					<ResumePreviewSection resume={resume} />
					
					<ResumeInfoSection 
						resume={resume}
						editing={editing}
						bind:formTitle
						bind:formTags
					/>

					{#if associatedApplications.length > 0}
						<div class="section">
							<h2 class="section-title">Associated Job Applications</h2>
							<div class="applications-list">
								{#each associatedApplications as app}
									<div class="application-item" onclick={() => goto(`/applications/${app.id}`)}>
										<div class="app-info">
											<span class="app-company">{app.companyName}</span>
											<span class="app-title">{app.jobTitle}</span>
										</div>
										<StatusBadge status={app.status} label={app.status} />
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<div class="sidebar">
					<ResumeMetricsSection resume={resume} />
					<ResumeActionsSection
						resume={resume}
						onMarkAsMain={handleMarkAsMain}
						onUnmarkAsMain={handleUnmarkAsMain}
						onMarkAsFeatured={handleMarkAsFeatured}
						onUnmarkAsFeatured={handleUnmarkAsFeatured}
						onRecalculateMetrics={handleRecalculateMetrics}
						onDelete={handleDelete}
					/>
				</div>
			</div>
		</div>
	{/if}
</div>

<ConfirmationModal
	bind:open={showDeleteConfirm}
	title="Delete Resume"
	message="Are you sure you want to delete this resume? This action cannot be undone."
	onConfirm={confirmDelete}
/>

<ToastContainer />

<style>
	.page-header {
		margin-bottom: 1.5rem;
	}

	.back-button {
		padding: 0.5rem 1rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		background-color: white;
		color: #1f2937;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s;
	}

	.back-button:hover {
		background-color: #f9fafb;
	}

	.resume-detail {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.detail-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 1.5rem;
		background-color: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
	}

	.header-content {
		flex: 1;
	}

	.title-section {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 0.75rem;
	}

	.resume-title {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 600;
		color: #1f2937;
	}

	.title-input {
		flex: 1;
		max-width: 400px;
	}

	.status-badges {
		display: flex;
		gap: 0.5rem;
	}

	.file-info {
		display: flex;
		gap: 1rem;
		align-items: center;
		font-size: 0.875rem;
		color: #6b7280;
	}

	.filename {
		font-weight: 500;
	}

	.file-size {
		color: #9ca3af;
	}

	.header-actions {
		display: flex;
		gap: 0.75rem;
	}

	.detail-content {
		display: grid;
		grid-template-columns: 1fr 300px;
		gap: 1.5rem;
	}

	.main-content {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.sidebar {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.section {
		background-color: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 1.5rem;
	}

	.section-title {
		margin: 0 0 1rem 0;
		font-size: 1.125rem;
		font-weight: 600;
		color: #1f2937;
	}

	.applications-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.application-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.application-item:hover {
		background-color: #f9fafb;
	}

	.app-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.app-company {
		font-weight: 500;
		color: #1f2937;
	}

	.app-title {
		font-size: 0.875rem;
		color: #6b7280;
	}

	@media (max-width: 1024px) {
		.detail-content {
			grid-template-columns: 1fr;
		}
	}
</style>
