<script lang="ts">
	import ResumeTableRow from '../_components/ResumeTableRow.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resumes = [],
		onDelete,
		onMarkAsMain,
		onUnmarkAsMain,
		onMarkAsFeatured,
		onUnmarkAsFeatured,
		selectedResumes = $bindable(new Set<string>()),
		onToggleSelection,
		onToggleSelectAll,
		onView
	}: {
		resumes?: Resume[];
		onDelete?: (id: string) => void;
		onMarkAsMain?: (id: string) => void;
		onUnmarkAsMain?: (id: string) => void;
		onMarkAsFeatured?: (id: string) => void;
		onUnmarkAsFeatured?: (id: string) => void;
		selectedResumes?: Set<string>;
		onToggleSelection?: (id: string) => void;
		onToggleSelectAll?: () => void;
		onView?: (resume: Resume) => void;
	} = $props();
</script>

{#if resumes.length === 0}
	<EmptyState
		title="No resumes found"
		description="Upload your first CV or generate one from a job application"
	/>
{:else}
	<div class="table-wrapper">
		<table class="table">
			<thead>
				<tr>
					<th>
						<input
							type="checkbox"
							checked={resumes.length > 0 && resumes.every(r => selectedResumes.has(r.id))}
							onchange={() => onToggleSelectAll?.()}
							class="select-all-checkbox"
						/>
					</th>
					<th>Title</th>
					<th>Status</th>
					<th>Tags</th>
					<th>Usage</th>
					<th>Interview Rate</th>
					<th>Offer Rate</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each resumes as resume}
					<ResumeTableRow 
						resume={resume} 
						{onDelete}
						{onMarkAsMain}
						{onUnmarkAsMain}
						{onMarkAsFeatured}
						{onUnmarkAsFeatured}
						selected={selectedResumes.has(resume.id)}
						onToggleSelection={onToggleSelection}
						onView={onView}
					/>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<style>
	.table-wrapper {
		background-color: var(--color-bg-primary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}
	
	.table {
		width: 100%;
		border-collapse: collapse;
		background-color: var(--color-bg-primary);
	}
	
	.table th {
		padding: var(--spacing-md);
		text-align: left;
		background-color: var(--color-bg-secondary);
		font-weight: var(--font-weight-semibold);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		border-bottom: 2px solid var(--color-border);
	}
	
	.table tbody tr {
		border-bottom: 1px solid var(--color-border);
	}
	
	.table tbody tr:last-child {
		border-bottom: none;
	}

	.select-all-checkbox {
		cursor: pointer;
	}
</style>
