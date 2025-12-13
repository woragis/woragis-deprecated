<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { JobApplication } from '$lib/api/jobapplications';
	import { listResumes, downloadResumeById, type Resume } from '$lib/api/resumes';
	import { Download, Eye, FileText } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { toastError, toastSuccess } from '$lib/utils/toast';
	
	let {
		application
	}: {
		application: JobApplication;
	} = $props();
	
	const tFn = useTranslation();
	let resumes: Resume[] = $state([]);
	let associatedResume: Resume | null = $state(null);
	
	// Use resume from application if available (from detail view), otherwise fetch from list
	const resumeFromApp = $derived(application.resume ? {
		id: application.resume.id,
		userId: application.resume.userId,
		title: application.resume.title,
		isMain: application.resume.isMain,
		isFeatured: application.resume.isFeatured,
		filePath: application.resume.filePath,
		fileName: application.resume.fileName,
		fileSize: application.resume.fileSize,
		tags: application.resume.tags,
		applicationsUsed: 0,
		interviewRate: 0,
		offerRate: 0,
		createdAt: application.resume.createdAt,
		updatedAt: application.resume.updatedAt
	} as Resume : null);
	
	onMount(async () => {
		// If resume is already in application object, use it
		if (resumeFromApp) {
			associatedResume = resumeFromApp;
			return;
		}
		
		// Otherwise, fetch from list
		try {
			resumes = await listResumes();
			if (application.resumeId) {
				associatedResume = resumes.find(r => r.id === application.resumeId) || null;
			}
		} catch (err) {
			console.error('Error fetching resumes:', err);
		}
	});
	
	// Update when application.resume changes
	$effect(() => {
		if (resumeFromApp) {
			associatedResume = resumeFromApp;
		}
	});
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}

	async function handleDownloadResume() {
		if (!application.resumeId || !associatedResume) {
			toastError('No resume available to download');
			return;
		}
		try {
			await downloadResumeById(application.resumeId, associatedResume.fileName);
			toastSuccess('Resume downloaded successfully');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to download resume');
		}
	}

	function handlePreviewResume() {
		if (!application.resumeId || !associatedResume) {
			toastError('No resume available to preview');
			return;
		}
		// Open in new tab for preview
		const baseURL = import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080';
		const url = `${baseURL}/api/resumes/${application.resumeId}/download`;
		window.open(url, '_blank');
	}
</script>

<Card>
	<div class="info-grid">
		<div class="info-item">
			<strong>Status:</strong>
			<StatusBadge status={application.status} type="status">
				{tFn(`jobApplications.status.${application.status}` as any)}
			</StatusBadge>
		</div>
		<div class="info-item"><strong>Location:</strong> {application.location || '—'}</div>
		<div class="info-item"><strong>Website:</strong> {application.website}</div>
		<div class="info-item">
			<strong>Job URL:</strong>
			<a href={application.jobUrl} target="_blank" rel="noopener noreferrer" class="link">
				{application.jobUrl}
			</a>
		</div>
		<div class="info-item"><strong>Interest Level:</strong> {application.interestLevel || '—'}</div>
		<div class="info-item"><strong>Source:</strong> {application.source || '—'}</div>
		<div class="info-item"><strong>Language:</strong> {application.language || '—'}</div>
		<div class="info-item"><strong>Applied At:</strong> {formatDate(application.appliedAt)}</div>
		{#if associatedResume}
			<div class="info-item info-item-full">
				<strong>Associated CV:</strong>
				<div class="cv-actions">
					<span class="cv-link">{associatedResume.title}</span>
					{#if associatedResume.isMain}
						<span class="cv-badge">Main</span>
					{/if}
					<div class="cv-buttons">
						<Button variant="secondary" size="sm" onclick={handleDownloadResume}>
							<Download size={14} />
							Download
						</Button>
						<Button variant="secondary" size="sm" onclick={handlePreviewResume}>
							<Eye size={14} />
							Preview
						</Button>
					</div>
				</div>
			</div>
		{:else if application.resumeId}
			<div class="info-item"><strong>Associated CV:</strong> <span class="cv-missing">CV ID: {application.resumeId} (not found)</span></div>
		{:else}
			<div class="info-item"><strong>Associated CV:</strong> <span class="cv-missing">No CV associated</span></div>
		{/if}
		{#if application.salaryMin || application.salaryMax}
			<div class="info-item">
				<strong>Salary:</strong> {application.salaryMin || '—'} - {application.salaryMax || '—'} {application.salaryCurrency || ''}
			</div>
		{/if}
		{#if application.deadline}
			<div class="info-item"><strong>Deadline:</strong> {formatDate(application.deadline)}</div>
		{/if}
		{#if application.followUpDate}
			<div class="info-item"><strong>Follow-up Date:</strong> {formatDate(application.followUpDate)}</div>
		{/if}
		{#if application.coverLetter}
			<div class="info-item info-item-full">
				<strong>Cover Letter:</strong>
				<div class="text-content">{application.coverLetter}</div>
			</div>
		{/if}
		{#if application.jobDescription}
			<div class="info-item info-item-full">
				<strong>Job Description:</strong>
				<div class="text-content">{application.jobDescription}</div>
			</div>
		{/if}
		{#if application.notes}
			<div class="info-item"><strong>Notes:</strong> {application.notes}</div>
		{/if}
		{#if application.tags && application.tags.length > 0}
			<div class="info-item">
				<strong>Tags:</strong> {application.tags.join(', ')}
			</div>
		{/if}
	</div>
</Card>

<style>
	.info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: var(--spacing-md);
	}
	
	.info-item {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
	}
	
	.info-item-full {
		grid-column: 1 / -1;
	}
	
	.info-item strong {
		color: var(--color-text-secondary);
		font-weight: var(--font-weight-medium);
	}
	
	.link {
		color: var(--color-primary);
		text-decoration: none;
		word-break: break-all;
	}
	
	.link:hover {
		text-decoration: underline;
	}
	
	.text-content {
		margin-top: var(--spacing-sm);
		padding: var(--spacing-md);
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		white-space: pre-wrap;
		font-size: var(--font-size-sm);
		line-height: 1.6;
		color: var(--color-text-primary);
	}
	
	.cv-link {
		color: var(--color-primary);
		font-weight: var(--font-weight-medium);
	}
	
	.cv-badge {
		display: inline-block;
		margin-left: var(--spacing-xs);
		padding: 2px 6px;
		background-color: var(--color-primary);
		color: white;
		border-radius: var(--radius-sm);
		font-size: var(--font-size-xs);
		font-weight: var(--font-weight-medium);
	}
	
	.cv-missing {
		color: var(--color-text-secondary);
		font-style: italic;
	}

	.cv-actions {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
		margin-top: var(--spacing-xs);
	}

	.cv-buttons {
		display: flex;
		gap: var(--spacing-xs);
		flex-wrap: wrap;
	}
</style>
