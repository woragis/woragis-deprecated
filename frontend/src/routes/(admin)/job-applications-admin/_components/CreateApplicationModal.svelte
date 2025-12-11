<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { ApplicationStatus, CreateJobApplicationInput, JobApplication } from '$lib/api/jobapplications';
	import { getUserPreferences } from '$lib/api/userpreferences';
	import { listResumes, type Resume } from '$lib/api/resumes';
	import ApplicationTemplates, { type ApplicationTemplate } from './ApplicationTemplates.svelte';
	
	let {
		open = $bindable(false),
		onSubmit,
		existingApplications = []
	}: {
		open?: boolean;
		onSubmit?: (input: CreateJobApplicationInput) => Promise<void>;
		existingApplications?: JobApplication[];
	} = $props();
	
	const tFn = useTranslation();
	
	const LAST_WEBSITE_KEY = 'lastJobWebsite';
	const LAST_RESUME_KEY_PREFIX = 'lastResume_';
	
	let formCompanyName = $state('');
	let formLocation = $state('remote');
	let formJobTitle = $state('');
	let formJobUrl = $state('');
	let formWebsite = $state('');
	let formCoverLetter = $state('');
	let formLinkedInContact = $state(false);
	let formStatus = $state<ApplicationStatus>('pending');
	let formResumeId = $state<string>('');
	let formInterestLevel = $state<string>('');
	let formTags = $state<string[]>([]);
	let formTagInput = $state('');
	let formFollowUpDate = $state('');
	let formNotes = $state('');
	
	let resumes: Resume[] = $state([]);
	let loadingResumes = $state(false);
	let showTemplatesModal = $state(false);
	
	const interestLevels = [
		{ value: '', label: 'Not set' },
		{ value: 'low', label: 'Low' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'high', label: 'High' },
		{ value: 'very-high', label: 'Very High' }
	];
	
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
	 * Detects website name from URL hostname.
	 * Returns website name like "linkedin", "glassdoor", "indeed", etc.
	 */
	function detectWebsiteFromUrl(url: string): string | null {
		if (!url || !url.trim()) return null;

		try {
			const normalized = normalizeUrl(url);
			const urlObj = new URL(normalized);
			const hostname = urlObj.hostname.toLowerCase();

			// Remove www. prefix if present
			const cleanHostname = hostname.replace(/^www\./, '');

			// Map common job board domains to website names
			const websiteMap: Record<string, string> = {
				'linkedin.com': 'linkedin',
				'glassdoor.com': 'glassdoor',
				'indeed.com': 'indeed',
				'monster.com': 'monster',
				'ziprecruiter.com': 'ziprecruiter',
				'careerbuilder.com': 'careerbuilder',
				'simplyhired.com': 'simplyhired',
				'angel.co': 'angel',
				'stackoverflow.com': 'stackoverflow',
				'github.com': 'github',
				'remoteok.io': 'remoteok',
				'weworkremotely.com': 'weworkremotely',
				'flexjobs.com': 'flexjobs',
				'upwork.com': 'upwork',
				'freelancer.com': 'freelancer'
			};

			// Check exact match first
			if (websiteMap[cleanHostname]) {
				return websiteMap[cleanHostname];
			}

			// Check if hostname contains any of the known domains
			for (const [domain, website] of Object.entries(websiteMap)) {
				if (cleanHostname.includes(domain)) {
					return website;
				}
			}

			// If no match, try to extract the main domain name
			// e.g., "jobs.example.com" -> "example"
			const parts = cleanHostname.split('.');
			if (parts.length >= 2) {
				// Return the main domain (second to last part, or last if only 2 parts)
				const mainDomain = parts.length > 2 ? parts[parts.length - 2] : parts[0];
				return mainDomain;
			}

			return null;
		} catch (err) {
			// Silently fail - URL parsing is best effort
			console.debug('Failed to detect website from URL:', err);
			return null;
		}
	}

	/**
	 * Parses job URL to extract company name, job title, location, and website.
	 * Only fills fields that are currently empty.
	 */
	function parseJobUrl(url: string) {
		if (!url || !url.trim()) return;

		try {
			const normalized = normalizeUrl(url);
			const urlObj = new URL(normalized);
			const hostname = urlObj.hostname.toLowerCase();
			const pathname = urlObj.pathname;

			// Auto-detect website from URL if website field is empty
			if (!formWebsite.trim()) {
				const detectedWebsite = detectWebsiteFromUrl(url);
				if (detectedWebsite) {
					formWebsite = detectedWebsite;
					// Load resume for the detected website
					loadResumeForWebsite(detectedWebsite);
				}
			}

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
	
	async function loadResumes() {
		if (loadingResumes || resumes.length > 0) return;
		
		loadingResumes = true;
		try {
			resumes = await listResumes();
		} catch (err) {
			console.warn('Failed to load resumes:', err);
			resumes = [];
		} finally {
			loadingResumes = false;
		}
	}

	function loadResumeForWebsite(website: string) {
		if (!website || typeof window === 'undefined') return;
		
		const key = `${LAST_RESUME_KEY_PREFIX}${website.toLowerCase()}`;
		const lastResumeId = localStorage.getItem(key);
		
		if (lastResumeId && resumes.some(r => r.id === lastResumeId)) {
			formResumeId = lastResumeId;
		} else {
			// If no saved resume for this website, try to use main resume
			const mainResume = resumes.find(r => r.isMain);
			if (mainResume) {
				formResumeId = mainResume.id;
			}
		}
	}

	async function loadDefaults() {
		// Default location to "remote"
		formLocation = 'remote';
		
		// Load resumes
		await loadResumes();
		
		// Try to load website from localStorage first
		if (typeof window !== 'undefined') {
			const lastWebsite = localStorage.getItem(LAST_WEBSITE_KEY);
			if (lastWebsite) {
				formWebsite = lastWebsite.toLowerCase();
				loadResumeForWebsite(lastWebsite);
				return;
			}
		}
		
		// If localStorage is empty, try to fetch from backend preferences
		try {
			const preferences = await getUserPreferences();
			if (preferences.defaultWebsite) {
				formWebsite = preferences.defaultWebsite.toLowerCase();
				loadResumeForWebsite(preferences.defaultWebsite);
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
		formResumeId = '';
		formInterestLevel = '';
		formTags = [];
		formTagInput = '';
		formFollowUpDate = '';
		formNotes = '';
	}

	function addTag() {
		const tag = formTagInput.trim().toLowerCase();
		if (tag && !formTags.includes(tag)) {
			formTags = [...formTags, tag];
			formTagInput = '';
		}
	}

	function removeTag(tag: string) {
		formTags = formTags.filter(t => t !== tag);
	}

	function handleTagInputKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ',') {
			e.preventDefault();
			addTag();
		}
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

	// Load resume when website changes
	$effect(() => {
		if (formWebsite && resumes.length > 0) {
			loadResumeForWebsite(formWebsite);
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
	
	function checkForDuplicate(jobUrl: string, website: string): JobApplication | null {
		const normalizedUrl = normalizeUrl(jobUrl);
		const normalizedWebsite = website.trim().toLowerCase();
		
		return existingApplications.find(
			app => app.jobUrl === normalizedUrl && app.website.toLowerCase() === normalizedWebsite
		) || null;
	}

	async function handleSubmit() {
		if (!formCompanyName.trim() || !formJobTitle.trim() || !formJobUrl.trim() || !formWebsite.trim()) {
			alert(tFn('jobApplications.modal.required') + ' ' + tFn('jobApplications.modal.companyName') + ', ' + tFn('jobApplications.modal.jobTitle') + ', ' + tFn('jobApplications.modal.jobUrl') + ', ' + tFn('jobApplications.modal.website'));
			return;
		}
		
		// Normalize website to lowercase
		const normalizedWebsite = formWebsite.trim().toLowerCase();
		
		// Normalize URL (add https:// if missing)
		const normalizedJobUrl = normalizeUrl(formJobUrl);
		
		// Check for duplicate
		const duplicate = checkForDuplicate(normalizedJobUrl, normalizedWebsite);
		if (duplicate) {
			const message = `You've already applied to this job!\n\n` +
				`Company: ${duplicate.companyName}\n` +
				`Title: ${duplicate.jobTitle}\n` +
				`Status: ${duplicate.status}\n` +
				`Applied: ${duplicate.appliedAt ? new Date(duplicate.appliedAt).toLocaleDateString() : 'Not yet'}\n\n` +
				`Do you still want to create a duplicate application?`;
			
			if (!confirm(message)) {
				return; // User cancelled
			}
		}
		
		// Save website to localStorage for next time
		if (typeof window !== 'undefined' && normalizedWebsite) {
			localStorage.setItem(LAST_WEBSITE_KEY, normalizedWebsite);
			
			// Save resume selection for this website
			if (formResumeId) {
				const key = `${LAST_RESUME_KEY_PREFIX}${normalizedWebsite}`;
				localStorage.setItem(key, formResumeId);
			}
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
				status: formStatus,
				interestLevel: formInterestLevel.trim() || undefined,
				tags: formTags.length > 0 ? formTags : undefined,
				followUpDate: formFollowUpDate || undefined,
				notes: formNotes.trim() || undefined
			};
			
			await onSubmit(input);
			handleClose();
		}
	}
</script>

<Modal bind:open size="lg" title={tFn('jobApplications.modal.createTitle')}>
	<div class="modal-header-actions">
		<Button variant="secondary" size="sm" onclick={() => showTemplatesModal = true}>
			Templates
		</Button>
	</div>
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
			<div class="form-row">
				<Select label="Interest Level" bind:value={formInterestLevel}>
					{#each interestLevels as level}
						<option value={level.value}>{level.label}</option>
					{/each}
				</Select>
				<Input
					label="Follow-up Date"
					type="date"
					bind:value={formFollowUpDate}
				/>
			</div>
			<div class="form-group">
				<label class="form-label">Tags</label>
				<div class="tags-input-container">
					<input
						type="text"
						class="tags-input"
						placeholder="Add tags (press Enter or comma)"
						bind:value={formTagInput}
						onkeydown={handleTagInputKeydown}
					/>
					<Button variant="secondary" size="sm" onclick={addTag}>Add</Button>
				</div>
				{#if formTags.length > 0}
					<div class="tags-list">
						{#each formTags as tag}
							<span class="tag">
								{tag}
								<button class="tag-remove" onclick={() => removeTag(tag)}>×</button>
							</span>
						{/each}
					</div>
				{/if}
			</div>
			<Textarea
				label={tFn('jobApplications.modal.coverLetter')}
				bind:value={formCoverLetter}
				rows={6}
			/>
			<Textarea
				label="Notes"
				bind:value={formNotes}
				rows={4}
				placeholder="Add any notes about this application..."
			/>
			{#if resumes.length > 0}
				<div class="form-group">
					<Select 
						label="Resume" 
						bind:value={formResumeId}
					>
						<option value="">No resume selected</option>
						{#each resumes as resume}
							<option value={resume.id}>
								{resume.title}{resume.isMain ? ' (Main)' : ''}
							</option>
						{/each}
					</Select>
				</div>
			{/if}
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

	.tags-input-container {
		display: flex;
		gap: var(--spacing-sm);
		align-items: center;
	}

	.tags-input {
		flex: 1;
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
	}

	.tags-list {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
		margin-top: var(--spacing-xs);
	}

	.tag {
		display: inline-flex;
		align-items: center;
		gap: var(--spacing-xs);
		padding: var(--spacing-xs) var(--spacing-sm);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-xs);
		color: var(--color-text-primary);
	}

	.tag-remove {
		background: none;
		border: none;
		color: var(--color-text-secondary);
		cursor: pointer;
		font-size: 1.2rem;
		padding: 0;
		line-height: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.tag-remove:hover {
		color: var(--color-danger);
	}

	.form-label {
		font-weight: var(--font-weight-medium);
		font-size: var(--font-size-sm);
		margin-bottom: var(--spacing-xs);
	}

	.modal-header-actions {
		position: absolute;
		top: var(--spacing-md);
		right: var(--spacing-md);
		z-index: 10;
	}
</style>

<ApplicationTemplates bind:open={showTemplatesModal} onLoadTemplate={loadTemplate} />
