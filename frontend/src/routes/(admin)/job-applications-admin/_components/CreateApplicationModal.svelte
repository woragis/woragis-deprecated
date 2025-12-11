<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { ApplicationStatus, CreateJobApplicationInput } from '$lib/api/jobapplications';
	import { getUserPreferences } from '$lib/api/userpreferences';
	
	let {
		open = $bindable(false),
		onSubmit
	}: {
		open?: boolean;
		onSubmit?: (input: CreateJobApplicationInput) => Promise<void>;
	} = $props();
	
	const tFn = useTranslation();
	
	const LAST_WEBSITE_KEY = 'lastJobWebsite';
	
	let formCompanyName = $state('');
	let formLocation = $state('remote');
	let formJobTitle = $state('');
	let formJobUrl = $state('');
	let formWebsite = $state('');
	let formCoverLetter = $state('');
	let formLinkedInContact = $state(false);
	let formStatus = $state<ApplicationStatus>('pending');
	
	/**
	 * Normalizes a URL by adding https:// prefix if missing.
	 * Preserves existing http:// or https:// prefixes.
	 */
	function normalizeUrl(url: string): string {
		const trimmed = url.trim();
		if (!trimmed) return trimmed;
		
		// Check if URL already has a protocol
		if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
			return trimmed;
		}
		
		// Add https:// prefix
		return `https://${trimmed}`;
	}

	/**
	 * Parses job URL to extract company name, job title, and location.
	 * Only fills fields that are currently empty.
	 */
	function parseJobUrl(url: string) {
		if (!url || !url.trim()) return;

		try {
			const normalized = normalizeUrl(url);
			const urlObj = new URL(normalized);
			const hostname = urlObj.hostname.toLowerCase();
			const pathname = urlObj.pathname;

			// LinkedIn job URL pattern: /jobs/view/{id} or /jobs/collections/{collection}/jobs/{id}
			if (hostname.includes('linkedin.com')) {
				// Try to extract from URL path
				// LinkedIn URLs often have job title and company in the path or query params
				const pathParts = pathname.split('/').filter(p => p);
				
				// Check query params for title/company
				const titleParam = urlObj.searchParams.get('title');
				const companyParam = urlObj.searchParams.get('company');
				
				// Only fill if field is empty
				if (titleParam && !formJobTitle.trim()) {
					formJobTitle = decodeURIComponent(titleParam).replace(/-/g, ' ');
				}
				if (companyParam && !formCompanyName.trim()) {
					formCompanyName = decodeURIComponent(companyParam).replace(/-/g, ' ');
				}
			}
			// Glassdoor job URL pattern: /Job/{location}-{job-title}-{company}-{id}.htm
			else if (hostname.includes('glassdoor.com')) {
				const pathParts = pathname.split('/').filter(p => p);
				if (pathParts.length >= 2 && pathParts[0] === 'Job') {
					const jobPart = pathParts[1].replace('.htm', '').replace('.html', '');
					const segments = jobPart.split('-');
					
					// Glassdoor format: location-job-title-company-id
					// Try to extract location (first segment if it looks like a location)
					if (segments.length > 0 && !formLocation.trim() || formLocation === 'remote') {
						// First segment might be location
						const possibleLocation = segments[0];
						if (possibleLocation && possibleLocation.length > 2) {
							formLocation = possibleLocation.replace(/_/g, ' ');
						}
					}
					
					// Last segment before .htm is usually ID, second to last might be company
					// This is heuristic-based
					if (segments.length >= 3 && !formCompanyName.trim()) {
						// Try to extract company name (usually one of the middle segments)
						const companySegment = segments[segments.length - 2];
						if (companySegment && companySegment.length > 2) {
							formCompanyName = companySegment.replace(/_/g, ' ');
						}
					}
					
					if (segments.length >= 2 && !formJobTitle.trim()) {
						// Job title is usually in the middle segments
						const titleSegments = segments.slice(1, -2); // Skip location and company/id
						if (titleSegments.length > 0) {
							formJobTitle = titleSegments.join(' ').replace(/_/g, ' ');
						}
					}
				}
			}
			// Indeed job URL pattern: /viewjob?jk={id}
			else if (hostname.includes('indeed.com')) {
				const jk = urlObj.searchParams.get('jk');
				// Indeed URLs don't typically have title/company in URL, but we can try
				const pathParts = pathname.split('/').filter(p => p);
				// Indeed sometimes has job title in path
				if (pathParts.length > 0 && !formJobTitle.trim()) {
					const lastPart = pathParts[pathParts.length - 1];
					if (lastPart && lastPart !== 'viewjob') {
						formJobTitle = lastPart.replace(/-/g, ' ');
					}
				}
			}
			// Generic pattern: try to extract from URL path segments
			else {
				const pathParts = pathname.split('/').filter(p => p && p !== 'jobs' && p !== 'job');
				
				// Try to find job title in path (usually a descriptive segment)
				if (pathParts.length > 0 && !formJobTitle.trim()) {
					const possibleTitle = pathParts[pathParts.length - 1];
					if (possibleTitle && possibleTitle.length > 3 && !possibleTitle.match(/^\d+$/)) {
						formJobTitle = possibleTitle.replace(/-/g, ' ').replace(/_/g, ' ');
					}
				}
			}
		} catch (err) {
			// Silently fail - URL parsing is best effort
			console.debug('Failed to parse job URL:', err);
		}
	}

	
	const statuses: ApplicationStatus[] = [
		'pending',
		'processing',
		'applied',
		'contacted',
		'rejected',
		'accepted',
		'failed'
	];
	
	async function loadDefaults() {
		// Default location to "remote"
		formLocation = 'remote';
		
		// Try to load website from localStorage first
		if (typeof window !== 'undefined') {
			const lastWebsite = localStorage.getItem(LAST_WEBSITE_KEY);
			if (lastWebsite) {
				formWebsite = lastWebsite.toLowerCase();
				return;
			}
		}
		
		// If localStorage is empty, try to fetch from backend preferences
		try {
			const preferences = await getUserPreferences();
			if (preferences.defaultWebsite) {
				formWebsite = preferences.defaultWebsite.toLowerCase();
			}
		} catch (err) {
			// Silently fail - user can still enter website manually
			console.warn('Failed to load default website from preferences:', err);
		}
	}
	
	function resetForm() {
		formCompanyName = '';
		formLocation = 'remote';
		formJobTitle = '';
		formJobUrl = '';
		formWebsite = '';
		formCoverLetter = '';
		formLinkedInContact = false;
		formStatus = 'pending';
	}
	
	function handleClose() {
		open = false;
		resetForm();
	}
	
	// Load defaults when modal opens
	$effect(() => {
		if (open) {
			loadDefaults();
		}
	});

	// Parse URL when it changes (debounced to avoid parsing on every keystroke)
	let urlParseTimeout: ReturnType<typeof setTimeout> | null = null;
	$effect(() => {
		if (formJobUrl && formJobUrl.trim()) {
			// Clear previous timeout
			if (urlParseTimeout) {
				clearTimeout(urlParseTimeout);
			}
			
			// Debounce parsing - wait 500ms after user stops typing
			urlParseTimeout = setTimeout(() => {
				// Only parse if URL looks complete (has protocol or is a valid domain)
				if (formJobUrl.includes('://') || (formJobUrl.includes('.') && formJobUrl.length > 10)) {
					parseJobUrl(formJobUrl);
				}
			}, 500);
		}
		
		// Cleanup
		return () => {
			if (urlParseTimeout) {
				clearTimeout(urlParseTimeout);
			}
		};
	});
	
	async function handleSubmit() {
		if (!formCompanyName.trim() || !formJobTitle.trim() || !formJobUrl.trim() || !formWebsite.trim()) {
			alert(tFn('jobApplications.modal.required') + ' ' + tFn('jobApplications.modal.companyName') + ', ' + tFn('jobApplications.modal.jobTitle') + ', ' + tFn('jobApplications.modal.jobUrl') + ', ' + tFn('jobApplications.modal.website'));
			return;
		}
		
		// Normalize website to lowercase
		const normalizedWebsite = formWebsite.trim().toLowerCase();
		
		// Normalize URL (add https:// if missing)
		const normalizedJobUrl = normalizeUrl(formJobUrl);
		
		// Save website to localStorage for next time
		if (typeof window !== 'undefined' && normalizedWebsite) {
			localStorage.setItem(LAST_WEBSITE_KEY, normalizedWebsite);
		}
		
		if (onSubmit) {
			const input: CreateJobApplicationInput = {
				companyName: formCompanyName.trim(),
				location: formLocation.trim() || undefined,
				jobTitle: formJobTitle.trim(),
				jobUrl: normalizedJobUrl,
				website: normalizedWebsite,
				coverLetter: formCoverLetter.trim() || undefined,
				linkedInContact: formLinkedInContact,
				status: formStatus
			};
			
			await onSubmit(input);
			handleClose();
		}
	}
</script>

<Modal bind:open size="lg" title={tFn('jobApplications.modal.createTitle')}>
	<div class="modal-content">
		<div class="form">
			<div class="form-row">
				<Input
					label="{tFn('jobApplications.modal.companyName')} {tFn('jobApplications.modal.required')}"
					bind:value={formCompanyName}
					required
				/>
				<Input
					label="{tFn('jobApplications.modal.jobTitle')} {tFn('jobApplications.modal.required')}"
					bind:value={formJobTitle}
					required
				/>
			</div>
			<div class="form-row">
				<Input
					label="{tFn('jobApplications.modal.jobUrl')} {tFn('jobApplications.modal.required')}"
					type="url"
					bind:value={formJobUrl}
					required
				/>
				<Input
					label="{tFn('jobApplications.modal.website')} {tFn('jobApplications.modal.required')}"
					bind:value={formWebsite}
					placeholder="linkedin, glassdoor, etc."
					required
				/>
			</div>
			<div class="form-row">
				<Input
					label={tFn('jobApplications.modal.location')}
					bind:value={formLocation}
				/>
				<Select label={tFn('jobApplications.modal.status')} bind:value={formStatus}>
					{#each statuses as status}
						<option value={status}>{tFn(`jobApplications.status.${status}` as any)}</option>
					{/each}
				</Select>
			</div>
			<Textarea
				label={tFn('jobApplications.modal.coverLetter')}
				bind:value={formCoverLetter}
				rows={6}
			/>
			<div class="form-group">
				<label class="checkbox-label">
					<input type="checkbox" bind:checked={formLinkedInContact} />
					{tFn('jobApplications.modal.linkedInContact')}
				</label>
			</div>
			<div class="form-actions">
				<Button onclick={handleSubmit} variant="primary">{tFn('jobApplications.modal.create')}</Button>
				<Button onclick={handleClose} variant="secondary">{tFn('jobApplications.modal.cancel')}</Button>
			</div>
		</div>
	</div>
</Modal>

<style>
	.modal-content {
		padding: var(--spacing-lg);
	}
	
	.form {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}
	
	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--spacing-md);
	}
	
	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}
	
	.checkbox-label {
		display: flex;
		align-items: center;
		gap: var(--spacing-sm);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		cursor: pointer;
	}
	
	.checkbox-label input[type="checkbox"] {
		cursor: pointer;
	}
	
	.form-actions {
		display: flex;
		gap: var(--spacing-sm);
		margin-top: var(--spacing-sm);
		justify-content: flex-end;
	}
</style>
