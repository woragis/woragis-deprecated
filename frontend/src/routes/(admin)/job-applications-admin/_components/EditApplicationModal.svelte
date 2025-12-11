<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { ApplicationStatus, UpdateJobApplicationInput, JobApplication } from '$lib/api/jobapplications';
	
	let {
		open = $bindable(false),
		application,
		onSubmit
	}: {
		open?: boolean;
		application?: JobApplication | null;
		onSubmit?: (id: string, input: UpdateJobApplicationInput) => Promise<void>;
	} = $props();
	
	const tFn = useTranslation();
	
	let formStatus = $state<ApplicationStatus>('pending');
	let formInterestLevel = $state<string>('');
	let formTags = $state<string[]>([]);
	let formTagInput = $state('');
	let formFollowUpDate = $state('');
	let formNotes = $state('');
	
	const statuses: ApplicationStatus[] = [
		'pending',
		'processing',
		'applied',
		'contacted',
		'rejected',
		'accepted',
		'failed'
	];
	
	const interestLevels = [
		{ value: '', label: 'Not Set' },
		{ value: 'low', label: 'Low' },
		{ value: 'medium', label: 'Medium' },
		{ value: 'high', label: 'High' },
		{ value: 'very-high', label: 'Very High' }
	];
	
	$effect(() => {
		if (application && open) {
			formStatus = application.status;
			formInterestLevel = application.interestLevel || '';
			formTags = application.tags || [];
			formTagInput = '';
			formFollowUpDate = application.followUpDate 
				? new Date(application.followUpDate).toISOString().split('T')[0]
				: '';
			formNotes = application.notes || '';
		}
	});
	
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
	}
	
	async function handleSubmit() {
		if (!application) return;
		
		if (onSubmit) {
			// Update status separately if changed
			if (formStatus !== application.status) {
				try {
					const { updateJobApplicationStatus } = await import('$lib/api/jobapplications');
					await updateJobApplicationStatus(application.id, { status: formStatus });
				} catch (err) {
					console.error('Error updating status:', err);
				}
			}
			
			const input: UpdateJobApplicationInput = {
				interestLevel: formInterestLevel.trim() || undefined,
				tags: formTags.length > 0 ? formTags : undefined,
				followUpDate: formFollowUpDate ? new Date(formFollowUpDate).toISOString() : undefined,
				notes: formNotes.trim() || undefined
			};
			
			await onSubmit(application.id, input);
			handleClose();
		}
	}
</script>

{#if application}
	<Modal bind:open size="lg" title="Edit Application">
		<div class="modal-content">
			<div class="form">
				<div class="form-row">
					<Select label={tFn('jobApplications.modal.status')} bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{tFn(`jobApplications.status.${status}` as any)}</option>
						{/each}
					</Select>
					<Select label="Interest Level" bind:value={formInterestLevel}>
						{#each interestLevels as level}
							<option value={level.value}>{level.label}</option>
						{/each}
					</Select>
				</div>
				<div class="form-row">
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
					label="Notes"
					bind:value={formNotes}
					rows={6}
					placeholder="Add any notes about this application..."
				/>
				<div class="form-actions">
					<Button onclick={handleSubmit} variant="primary">Save Changes</Button>
					<Button onclick={handleClose} variant="secondary">Cancel</Button>
				</div>
			</div>
		</div>
	</Modal>
{/if}

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
</style>

