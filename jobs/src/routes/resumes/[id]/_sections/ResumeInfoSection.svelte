<script lang="ts">
	import TagInput from '$lib/components/ui/TagInput.svelte';
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resume,
		editing = false,
		formTitle = $bindable(''),
		formTags = $bindable<string[]>([])
	}: {
		resume: Resume;
		editing?: boolean;
		formTitle?: string;
		formTags?: string[];
	} = $props();

	let availableTags = $state<string[]>([]);

	$effect(() => {
		if (resume.tags) {
			availableTags = [...new Set([...availableTags, ...resume.tags])];
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}
</script>

<div class="section">
	<h2 class="section-title">Information</h2>
	<div class="info-grid">
		<div class="info-item">
			<label class="info-label">Title</label>
			{#if editing}
				<input
					type="text"
					bind:value={formTitle}
					class="info-input"
					placeholder="Resume title"
				/>
			{:else}
				<span class="info-value">{resume.title}</span>
			{/if}
		</div>
		<div class="info-item">
			<label class="info-label">Filename</label>
			<span class="info-value">{resume.fileName}</span>
		</div>
		<div class="info-item">
			<label class="info-label">File Size</label>
			<span class="info-value">{(resume.fileSize / 1024).toFixed(2)} KB</span>
		</div>
		<div class="info-item">
			<label class="info-label">Created</label>
			<span class="info-value">{formatDate(resume.createdAt)}</span>
		</div>
		<div class="info-item">
			<label class="info-label">Updated</label>
			<span class="info-value">{formatDate(resume.updatedAt)}</span>
		</div>
		<div class="info-item full-width">
			<label class="info-label">Tags</label>
			{#if editing}
				<TagInput
					bind:tags={formTags}
					{availableTags}
					placeholder="Add tags..."
				/>
			{:else}
				<div class="tags">
					{#each (resume.tags || []) as tag}
						<span class="tag">{tag}</span>
					{:else}
						<span class="no-tags">No tags</span>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.info-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
	}

	.info-item {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.info-item.full-width {
		grid-column: 1 / -1;
	}

	.info-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #6b7280;
	}

	.info-value {
		font-size: 0.875rem;
		color: #1f2937;
	}

	.info-input {
		padding: 0.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		color: #1f2937;
	}

	.tags {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.tag {
		padding: 0.25rem 0.75rem;
		background-color: #f3f4f6;
		border: 1px solid #e5e7eb;
		border-radius: 0.25rem;
		font-size: 0.875rem;
		color: #1f2937;
	}

	.no-tags {
		font-size: 0.875rem;
		color: #9ca3af;
		font-style: italic;
	}
</style>
