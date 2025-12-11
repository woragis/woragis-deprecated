<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import type { ApplicationStatus } from '$lib/api/jobapplications';

	const TEMPLATES_KEY = 'jobApplicationTemplates';

	export interface ApplicationTemplate {
		name: string;
		location?: string;
		website?: string;
		interestLevel?: string;
		tags?: string[];
		notes?: string;
		coverLetter?: string;
		status?: ApplicationStatus;
	}

	let {
		open = $bindable(false),
		onLoadTemplate
	}: {
		open?: boolean;
		onLoadTemplate?: (template: ApplicationTemplate) => void;
	} = $props();

	let templates = $state<ApplicationTemplate[]>([]);
	let showCreateModal = $state(false);
	let templateName = $state('');
	let newTemplate: ApplicationTemplate = $state({
		name: '',
		location: '',
		website: '',
		interestLevel: '',
		tags: [],
		notes: '',
		coverLetter: '',
		status: 'pending'
	});

	$effect(() => {
		if (typeof window !== 'undefined') {
			loadTemplates();
		}
	});

	function loadTemplates() {
		if (typeof window !== 'undefined') {
			const stored = localStorage.getItem(TEMPLATES_KEY);
			if (stored) {
				try {
					templates = JSON.parse(stored);
				} catch (e) {
					console.warn('Failed to load templates:', e);
				}
			}
		}
	}

	function saveTemplates() {
		if (typeof window !== 'undefined') {
			localStorage.setItem(TEMPLATES_KEY, JSON.stringify(templates));
		}
	}

	function handleCreateTemplate() {
		if (!templateName.trim()) return;

		const template: ApplicationTemplate = {
			name: templateName.trim(),
			...newTemplate
		};

		templates = [...templates, template];
		saveTemplates();
		
		templateName = '';
		newTemplate = {
			name: '',
			location: '',
			website: '',
			interestLevel: '',
			tags: [],
			notes: '',
			coverLetter: '',
			status: 'pending'
		};
		showCreateModal = false;
	}

	function handleLoadTemplate(template: ApplicationTemplate) {
		if (onLoadTemplate) {
			onLoadTemplate(template);
		}
		open = false;
	}

	function handleDeleteTemplate(index: number) {
		if (confirm('Delete this template?')) {
			templates = templates.filter((_, i) => i !== index);
			saveTemplates();
		}
	}
</script>

<Modal bind:open size="md" title="Application Templates">
	<div class="templates-container">
		<div class="templates-header">
			<Button variant="primary" size="sm" onclick={() => showCreateModal = true}>
				Create Template
			</Button>
		</div>
		
		{#if templates.length === 0}
			<p class="empty-message">No templates saved. Create one to save common fields.</p>
		{:else}
			<div class="templates-list">
				{#each templates as template, index}
					<div class="template-item">
						<div class="template-info">
							<h4>{template.name}</h4>
							<div class="template-details">
								{#if template.location}
									<span class="detail-badge">Location: {template.location}</span>
								{/if}
								{#if template.website}
									<span class="detail-badge">Website: {template.website}</span>
								{/if}
								{#if template.interestLevel}
									<span class="detail-badge">Interest: {template.interestLevel}</span>
								{/if}
								{#if template.tags && template.tags.length > 0}
									<span class="detail-badge">Tags: {template.tags.join(', ')}</span>
								{/if}
							</div>
						</div>
						<div class="template-actions">
							<Button variant="primary" size="sm" onclick={() => handleLoadTemplate(template)}>
								Load
							</Button>
							<Button variant="danger" size="sm" onclick={() => handleDeleteTemplate(index)}>
								Delete
							</Button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</Modal>

{#if showCreateModal}
	<Modal bind:open={showCreateModal} size="lg" title="Create Template">
		<div class="create-template-form">
			<Input
				label="Template Name"
				bind:value={templateName}
				placeholder="e.g., Remote LinkedIn Jobs"
				required
			/>
			<div class="form-row">
				<Input
					label="Default Location"
					bind:value={newTemplate.location}
					placeholder="e.g., remote"
				/>
				<Input
					label="Default Website"
					bind:value={newTemplate.website}
					placeholder="e.g., linkedin"
				/>
			</div>
			<div class="form-row">
				<select bind:value={newTemplate.interestLevel} class="form-select">
					<option value="">No default</option>
					<option value="low">Low</option>
					<option value="medium">Medium</option>
					<option value="high">High</option>
					<option value="very-high">Very High</option>
				</select>
			</div>
			<div class="form-actions">
				<Button variant="primary" onclick={handleCreateTemplate} disabled={!templateName.trim()}>
					Save Template
				</Button>
				<Button variant="secondary" onclick={() => showCreateModal = false}>
					Cancel
				</Button>
			</div>
		</div>
	</Modal>
{/if}

<style>
	.templates-container {
		padding: var(--spacing-md);
	}

	.templates-header {
		margin-bottom: var(--spacing-md);
	}

	.templates-list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
	}

	.template-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--spacing-md);
		background-color: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
	}

	.template-info h4 {
		margin: 0 0 var(--spacing-xs) 0;
		font-size: var(--font-size-md);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.template-details {
		display: flex;
		flex-wrap: wrap;
		gap: var(--spacing-xs);
		margin-top: var(--spacing-xs);
	}

	.detail-badge {
		font-size: var(--font-size-xs);
		color: var(--color-text-secondary);
	}

	.template-actions {
		display: flex;
		gap: var(--spacing-xs);
	}

	.empty-message {
		text-align: center;
		color: var(--color-text-secondary);
		padding: var(--spacing-lg);
	}

	.create-template-form {
		padding: var(--spacing-md);
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--spacing-md);
	}

	.form-select {
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background-color: var(--color-bg-primary);
		color: var(--color-text-primary);
	}

	.form-actions {
		display: flex;
		gap: var(--spacing-sm);
		justify-content: flex-end;
		margin-top: var(--spacing-md);
	}
</style>

