<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { ApplicationStatus, CreateJobApplicationInput } from '$lib/api/jobapplications';
	
	let {
		open = $bindable(false),
		onSubmit
	}: {
		open?: boolean;
		onSubmit?: (input: CreateJobApplicationInput) => Promise<void>;
	} = $props();
	
	const tFn = useTranslation();
	
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
	
	function handleClose() {
		open = false;
		resetForm();
	}
	
	async function handleSubmit() {
		if (!formCompanyName.trim() || !formJobTitle.trim() || !formJobUrl.trim() || !formWebsite.trim()) {
			alert(tFn('jobApplications.modal.required') + ' ' + tFn('jobApplications.modal.companyName') + ', ' + tFn('jobApplications.modal.jobTitle') + ', ' + tFn('jobApplications.modal.jobUrl') + ', ' + tFn('jobApplications.modal.website'));
			return;
		}
		
		if (onSubmit) {
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
