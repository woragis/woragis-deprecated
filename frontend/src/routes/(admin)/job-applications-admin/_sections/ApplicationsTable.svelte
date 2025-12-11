<script lang="ts">
	import ApplicationTableRow from '../_components/ApplicationTableRow.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { useTranslation } from '$lib/i18n';
	import type { JobApplication } from '$lib/api/jobapplications';
	
	let {
		applications = [],
		onDelete,
		selectedApplications = $bindable(new Set<string>()),
		onToggleSelection,
		onToggleSelectAll,
		onView,
		onEdit
	}: {
		applications?: JobApplication[];
		onDelete?: (id: string) => void;
		selectedApplications?: Set<string>;
		onToggleSelection?: (id: string) => void;
		onToggleSelectAll?: () => void;
		onView?: (application: JobApplication) => void;
		onEdit?: (application: JobApplication) => void;
	} = $props();
	
	const tFn = useTranslation();
</script>

{#if applications.length === 0}
	<EmptyState
		title={tFn('jobApplications.empty')}
		description="Create your first job application to get started"
	/>
{:else}
	<div class="table-wrapper">
		<table class="table">
			<thead>
				<tr>
					<th>
						<input
							type="checkbox"
							checked={applications.length > 0 && applications.every(a => selectedApplications.has(a.id))}
							onchange={() => onToggleSelectAll?.()}
							class="select-all-checkbox"
						/>
					</th>
					<th>{tFn('jobApplications.table.company')}</th>
					<th>{tFn('jobApplications.table.jobTitle')}</th>
					<th>{tFn('jobApplications.table.location')}</th>
					<th>{tFn('jobApplications.table.website')}</th>
					<th>{tFn('jobApplications.table.language')}</th>
					<th>{tFn('jobApplications.table.status')}</th>
					<th>{tFn('jobApplications.table.interest')}</th>
					<th>Tags</th>
					<th>Follow-up</th>
					<th>{tFn('jobApplications.table.appliedAt')}</th>
					<th>{tFn('jobApplications.table.actions')}</th>
				</tr>
			</thead>
			<tbody>
				{#each applications as app}
					<ApplicationTableRow 
						application={app} 
						{onDelete}
						selected={selectedApplications.has(app.id)}
						onToggleSelection={onToggleSelection}
						onView={onView}
						onEdit={onEdit}
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
