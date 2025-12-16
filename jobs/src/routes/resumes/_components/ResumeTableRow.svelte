<script lang="ts">
	import { Eye, Trash2, Star, StarOff, Award, AwardOff, Download } from 'lucide-svelte';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import { downloadResumeById } from '$lib/api/resumes';
	import { toastError, getApiErrorMessage } from '$lib/utils/toast';
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resume,
		onDelete,
		onMarkAsMain,
		onUnmarkAsMain,
		onMarkAsFeatured,
		onUnmarkAsFeatured,
		selected = false,
		onToggleSelection,
		onView
	}: {
		resume: Resume;
		onDelete?: (id: string) => void;
		onMarkAsMain?: (id: string) => void;
		onUnmarkAsFeatured?: (id: string) => void;
		onMarkAsFeatured?: (id: string) => void;
		onUnmarkAsMain?: (id: string) => void;
		selected?: boolean;
		onToggleSelection?: (id: string) => void;
		onView?: (resume: Resume) => void;
	} = $props();

	async function handleDownload() {
		try {
			await downloadResumeById(resume.id, resume.fileName);
		} catch (err) {
			toastError(getApiErrorMessage(err, 'Failed to download resume'));
			console.error('Error downloading resume:', err);
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}
</script>

<tr class="table-row" class:selected>
	<td>
		<input
			type="checkbox"
			checked={selected}
			onchange={() => onToggleSelection?.(resume.id)}
			class="row-checkbox"
		/>
	</td>
	<td>
		<div class="title-cell">
			<span class="title">{resume.title}</span>
			<span class="filename">{resume.fileName}</span>
		</div>
	</td>
	<td>
		<div class="status-badges">
			{#if resume.isMain}
				<StatusBadge status="main" label="Main" />
			{/if}
			{#if resume.isFeatured}
				<StatusBadge status="featured" label="Featured" />
			{/if}
		</div>
	</td>
	<td>
		<div class="tags">
			{#each (resume.tags || []).slice(0, 3) as tag}
				<span class="tag">{tag}</span>
			{/each}
			{#if (resume.tags || []).length > 3}
				<span class="tag-more">+{(resume.tags || []).length - 3}</span>
			{/if}
		</div>
	</td>
	<td class="metric-cell">{resume.applicationsUsed}</td>
	<td class="metric-cell">{resume.interviewRate.toFixed(1)}%</td>
	<td class="metric-cell">{resume.offerRate.toFixed(1)}%</td>
	<td class="date-cell">{formatDate(resume.createdAt)}</td>
	<td>
		<div class="actions">
			<button
				class="action-btn"
				onclick={() => onView?.(resume)}
				title="View"
			>
				<Eye size={16} />
			</button>
			<button
				class="action-btn"
				onclick={handleDownload}
				title="Download"
			>
				<Download size={16} />
			</button>
			{#if resume.isMain}
				<button
					class="action-btn"
					onclick={() => onUnmarkAsMain?.(resume.id)}
					title="Unmark as Main"
				>
					<StarOff size={16} />
				</button>
			{:else}
				<button
					class="action-btn"
					onclick={() => onMarkAsMain?.(resume.id)}
					title="Mark as Main"
				>
					<Star size={16} />
				</button>
			{/if}
			{#if resume.isFeatured}
				<button
					class="action-btn"
					onclick={() => onUnmarkAsFeatured?.(resume.id)}
					title="Unmark as Featured"
				>
					<AwardOff size={16} />
				</button>
			{:else}
				<button
					class="action-btn"
					onclick={() => onMarkAsFeatured?.(resume.id)}
					title="Mark as Featured"
				>
					<Award size={16} />
				</button>
			{/if}
			<button
				class="action-btn danger"
				onclick={() => onDelete?.(resume.id)}
				title="Delete"
			>
				<Trash2 size={16} />
			</button>
		</div>
	</td>
</tr>

<style>
	.table-row {
		transition: background-color 0.2s;
	}

	.table-row:hover {
		background-color: #f9fafb;
	}

	.table-row.selected {
		background-color: #eff6ff;
	}

	.row-checkbox {
		cursor: pointer;
	}

	.title-cell {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.title {
		font-weight: 500;
		color: #1f2937;
	}

	.filename {
		font-size: 0.75rem;
		color: #6b7280;
	}

	.status-badges {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.tags {
		display: flex;
		gap: 0.25rem;
		flex-wrap: wrap;
	}

	.tag {
		padding: 0.125rem 0.5rem;
		background-color: #f3f4f6;
		border: 1px solid #e5e7eb;
		border-radius: 0.25rem;
		font-size: 0.75rem;
		color: #1f2937;
	}

	.tag-more {
		padding: 0.125rem 0.5rem;
		font-size: 0.75rem;
		color: #6b7280;
	}

	.metric-cell {
		text-align: center;
		font-weight: 500;
	}

	.date-cell {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.375rem;
		border: none;
		background-color: transparent;
		color: #6b7280;
		cursor: pointer;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	.action-btn:hover {
		background-color: #f3f4f6;
		color: #1f2937;
	}

	.action-btn.danger:hover {
		background-color: #fee2e2;
		color: #dc2626;
	}
</style>
