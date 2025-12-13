<script lang="ts">
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { JobApplication } from '$lib/api/jobapplications';
	import { listResumes, type Resume } from '$lib/api/resumes';
	import { onMount } from 'svelte';
	
	let {
		application
	}: {
		application: JobApplication;
	} = $props();
	
	const tFn = useTranslation();
	let resumes: Resume[] = $state([]);
	let associatedResume: Resume | null = $state(null);
	
	onMount(async () => {
		try {
			resumes = await listResumes();
			if (application.resumeId) {
				associatedResume = resumes.find(r => r.id === application.resumeId) || null;
			}
		} catch (err) {
			console.error('Error fetching resumes:', err);
		}
	});
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
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
			<div class="info-item">
				<strong>Associated CV:</strong>
				<span class="cv-link">{associatedResume.title}</span>
				{#if associatedResume.isMain}
					<span class="cv-badge">Main</span>
				{/if}
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
</style>
