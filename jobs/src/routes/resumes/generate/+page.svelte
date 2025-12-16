<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		requestCVGeneration,
		getResumeJobStatus,
		retryResumeJob,
		cancelResumeJob,
		type ResumeGenerationJobResponse,
		type ResumeJobStatus
	} from '$lib/api/resumes';
	import { listJobApplications, type JobApplication } from '$lib/api/jobapplications';
	import Button from '$lib/components/ui/Button.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import { toastSuccess, toastError, getApiErrorMessage } from '$lib/utils/toast';
	import { RefreshCw, X, CheckCircle, AlertCircle } from 'lucide-svelte';

	let applications: JobApplication[] = $state([]);
	let selectedApplicationId = $state('');
	let language = $state('en');
	let loading = $state(true);
	let error: string | null = $state(null);
	let generating = $state(false);
	let currentJob: ResumeJobStatus | null = $state(null);
	let pollingInterval: number | null = $state(null);

	const languages = [
		{ value: 'en', label: 'English' },
		{ value: 'pt', label: 'Portuguese' },
		{ value: 'es', label: 'Spanish' },
		{ value: 'fr', label: 'French' }
	];

	onMount(async () => {
		await fetchApplications();
		
		return () => {
			if (pollingInterval) {
				clearInterval(pollingInterval);
			}
		};
	});

	async function fetchApplications() {
		loading = true;
		error = null;
		try {
			applications = await listJobApplications();
			if (applications.length > 0 && !selectedApplicationId) {
				selectedApplicationId = applications[0].id;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load job applications';
			console.error('Error fetching applications:', err);
		} finally {
			loading = false;
		}
	}

	async function handleGenerate() {
		if (!selectedApplicationId) {
			toastError('Please select a job application');
			return;
		}

		generating = true;
		try {
			const response = await requestCVGeneration({
				jobApplicationId: selectedApplicationId,
				language: language || 'en'
			});
			
			toastSuccess('CV generation started! Checking status...');
			currentJob = {
				id: response.jobId,
				status: response.status as any,
				retryCount: 0,
				maxRetries: 3,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			};
			
			startPolling(response.jobId);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to start CV generation');
			toastError(message);
			console.error('Error generating CV:', err);
		} finally {
			generating = false;
		}
	}

	function startPolling(jobId: string) {
		if (pollingInterval) {
			clearInterval(pollingInterval);
		}

		pollingInterval = window.setInterval(async () => {
			try {
				const status = await getResumeJobStatus(jobId);
				currentJob = status;

				if (status.status === 'completed') {
					stopPolling();
					toastSuccess('CV generated successfully!');
					// Redirect to the generated resume if we have the result
					if (status.result) {
						// We'll need to get the resume ID from the result or redirect to list
						setTimeout(() => {
							goto('/resumes');
						}, 2000);
					}
				} else if (status.status === 'failed' || status.status === 'dead_letter') {
					stopPolling();
					toastError(status.error || 'CV generation failed');
				}
			} catch (err) {
				console.error('Error polling job status:', err);
				stopPolling();
			}
		}, 2000); // Poll every 2 seconds
	}

	function stopPolling() {
		if (pollingInterval) {
			clearInterval(pollingInterval);
			pollingInterval = null;
		}
	}

	async function handleRetry() {
		if (!currentJob) return;
		try {
			const response = await retryResumeJob(currentJob.id);
			currentJob = response;
			toastSuccess('Job retried successfully');
			startPolling(currentJob.id);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to retry job');
			toastError(message);
			console.error('Error retrying job:', err);
		}
	}

	async function handleCancel() {
		if (!currentJob) return;
		try {
			const response = await cancelResumeJob(currentJob.id);
			currentJob = response;
			stopPolling();
			toastSuccess('Job cancelled successfully');
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to cancel job');
			toastError(message);
			console.error('Error cancelling job:', err);
		}
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'completed':
				return 'accepted';
			case 'failed':
			case 'dead_letter':
				return 'rejected';
			case 'processing':
				return 'processing';
			case 'pending':
				return 'pending';
			default:
				return 'pending';
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-3xl">
	<div class="page-header">
		<button class="back-button" onclick={() => goto('/resumes')}>
			← Back to Resumes
		</button>
		<h1 class="page-title">Generate CV</h1>
		<p class="page-subtitle">Generate a tailored CV for a specific job application</p>
	</div>

	{#if error}
		<ErrorState message={error} onRetry={fetchApplications} />
	{:else if loading}
		<LoadingState message="Loading job applications..." />
	{:else}
		<div class="generate-form">
			<div class="form-section">
				<label class="form-label">Job Application *</label>
				<select
					bind:value={selectedApplicationId}
					class="form-select"
					disabled={generating || !!currentJob}
				>
					<option value="">Select a job application</option>
					{#each applications as app}
						<option value={app.id}>{app.companyName} - {app.jobTitle}</option>
					{/each}
				</select>
				{#if selectedApplicationId}
					{@const selectedApp = applications.find(a => a.id === selectedApplicationId)}
					{#if selectedApp}
						<div class="app-preview">
							<div class="app-info">
								<span class="app-company">{selectedApp.companyName}</span>
								<span class="app-title">{selectedApp.jobTitle}</span>
							</div>
							<StatusBadge status={selectedApp.status} label={selectedApp.status} />
						</div>
					{/if}
				{/if}
			</div>

			<div class="form-section">
				<label class="form-label">Language</label>
				<select
					bind:value={language}
					class="form-select"
					disabled={generating || !!currentJob}
				>
					{#each languages as lang}
						<option value={lang.value}>{lang.label}</option>
					{/each}
				</select>
			</div>

			{#if !currentJob}
				<div class="form-actions">
					<Button onclick={() => goto('/resumes')} variant="secondary" disabled={generating}>
						Cancel
					</Button>
					<Button onclick={handleGenerate} variant="primary" disabled={generating || !selectedApplicationId}>
						{generating ? 'Starting...' : 'Generate CV'}
					</Button>
				</div>
			{/if}

			{#if currentJob}
				<div class="job-status">
					<div class="status-header">
						<h2 class="status-title">Generation Status</h2>
						<StatusBadge status={getStatusColor(currentJob.status)} label={currentJob.status} />
					</div>

					<div class="status-details">
						<div class="status-item">
							<label>Job ID:</label>
							<span class="mono">{currentJob.id}</span>
						</div>
						<div class="status-item">
							<label>Retry Count:</label>
							<span>{currentJob.retryCount} / {currentJob.maxRetries}</span>
						</div>
						<div class="status-item">
							<label>Created:</label>
							<span>{formatDate(currentJob.createdAt)}</span>
						</div>
						<div class="status-item">
							<label>Last Updated:</label>
							<span>{formatDate(currentJob.updatedAt)}</span>
						</div>
						{#if currentJob.error}
							<div class="status-error">
								<AlertCircle size={20} />
								<div class="error-content">
									<strong>Error:</strong> {currentJob.error}
									{#if currentJob.errorType}
										<br />
										<small>Type: {currentJob.errorType}</small>
									{/if}
								</div>
							</div>
						{/if}
						{#if currentJob.status === 'completed' && currentJob.result}
							<div class="status-success">
								<CheckCircle size={20} />
								<div class="success-content">
									<strong>CV Generated Successfully!</strong>
									<br />
									<small>File: {currentJob.result.fileName}</small>
								</div>
							</div>
						{/if}
					</div>

					<div class="status-actions">
						{#if currentJob.status === 'pending' || currentJob.status === 'processing'}
							<Button onclick={handleCancel} variant="secondary">
								<X size={16} />
								Cancel Job
							</Button>
						{/if}
						{#if currentJob.status === 'failed' || currentJob.status === 'dead_letter'}
							<Button onclick={handleRetry} variant="primary">
								<RefreshCw size={16} />
								Retry
							</Button>
						{/if}
						{#if currentJob.status === 'completed'}
							<Button onclick={() => goto('/resumes')} variant="primary">
								View Resumes
							</Button>
						{/if}
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

<ToastContainer />

<style>
	.page-header {
		margin-bottom: 2rem;
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
		margin-bottom: 1rem;
	}

	.back-button:hover {
		background-color: #f9fafb;
	}

	.page-title {
		margin: 0 0 0.5rem 0;
		font-size: 1.875rem;
		font-weight: 600;
		color: #1f2937;
	}

	.page-subtitle {
		margin: 0;
		color: #6b7280;
		font-size: 0.875rem;
	}

	.generate-form {
		background-color: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 2rem;
	}

	.form-section {
		margin-bottom: 1.5rem;
	}

	.form-label {
		display: block;
		margin-bottom: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
	}

	.form-select {
		width: 100%;
		padding: 0.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		color: #1f2937;
		background-color: white;
		cursor: pointer;
		transition: border-color 0.2s;
	}

	.form-select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.form-select:disabled {
		background-color: #f3f4f6;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.app-preview {
		margin-top: 0.75rem;
		padding: 0.75rem;
		background-color: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
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

	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		margin-top: 2rem;
	}

	.job-status {
		margin-top: 2rem;
		padding: 1.5rem;
		background-color: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
	}

	.status-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.status-title {
		margin: 0;
		font-size: 1.125rem;
		font-weight: 600;
		color: #1f2937;
	}

	.status-details {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		margin-bottom: 1.5rem;
	}

	.status-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem 0;
		border-bottom: 1px solid #e5e7eb;
	}

	.status-item:last-child {
		border-bottom: none;
	}

	.status-item label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #6b7280;
	}

	.status-item span {
		font-size: 0.875rem;
		color: #1f2937;
	}

	.mono {
		font-family: monospace;
		font-size: 0.75rem;
	}

	.status-error {
		display: flex;
		gap: 0.75rem;
		padding: 1rem;
		background-color: #fee2e2;
		border: 1px solid #fecaca;
		border-radius: 0.375rem;
		color: #991b1b;
	}

	.status-success {
		display: flex;
		gap: 0.75rem;
		padding: 1rem;
		background-color: #d1fae5;
		border: 1px solid #a7f3d0;
		border-radius: 0.375rem;
		color: #065f46;
	}

	.error-content,
	.success-content {
		flex: 1;
		font-size: 0.875rem;
	}

	.status-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
	}
</style>
