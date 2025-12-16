<script lang="ts">
	import { goto } from '$app/navigation';
	import { uploadResume, listResumeTags, type Resume } from '$lib/api/resumes';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import TagInput from '$lib/components/ui/TagInput.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import { toastSuccess, toastError, getApiErrorMessage } from '$lib/utils/toast';
	import { Upload, X } from 'lucide-svelte';

	let selectedFile: File | null = $state(null);
	let title = $state('');
	let tags = $state<string[]>([]);
	let availableTags = $state<string[]>([]);
	let uploading = $state(false);
	let dragActive = $state(false);

	onMount(async () => {
		try {
			availableTags = await listResumeTags();
		} catch (err) {
			console.error('Error fetching tags:', err);
		}
	});

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			const file = target.files[0];
			if (file.type === 'application/pdf') {
				selectedFile = file;
				if (!title) {
					title = file.name.replace('.pdf', '');
				}
			} else {
				toastError('Only PDF files are allowed');
			}
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		dragActive = true;
	}

	function handleDragLeave(event: DragEvent) {
		event.preventDefault();
		dragActive = false;
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		dragActive = false;
		
		if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
			const file = event.dataTransfer.files[0];
			if (file.type === 'application/pdf') {
				selectedFile = file;
				if (!title) {
					title = file.name.replace('.pdf', '');
				}
			} else {
				toastError('Only PDF files are allowed');
			}
		}
	}

	function removeFile() {
		selectedFile = null;
	}

	async function handleUpload() {
		if (!selectedFile || !title.trim()) {
			toastError('Please select a file and enter a title');
			return;
		}

		uploading = true;
		try {
			const resume = await uploadResume(selectedFile, title.trim());
			toastSuccess('Resume uploaded successfully!');
			goto(`/resumes/${resume.id}`);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Failed to upload resume');
			toastError(message);
			console.error('Error uploading resume:', err);
		} finally {
			uploading = false;
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-3xl">
	<div class="page-header">
		<button class="back-button" onclick={() => goto('/resumes')}>
			← Back to Resumes
		</button>
		<h1 class="page-title">Upload CV</h1>
		<p class="page-subtitle">Upload a PDF file to add it to your resume collection</p>
	</div>

	<div class="upload-form">
		<div
			class="drop-zone"
			class:active={dragActive}
			ondragover={handleDragOver}
			ondragleave={handleDragLeave}
			ondrop={handleDrop}
		>
			{#if selectedFile}
				<div class="file-selected">
					<div class="file-info">
						<span class="file-name">{selectedFile.name}</span>
						<span class="file-size">{formatFileSize(selectedFile.size)}</span>
					</div>
					<button class="remove-button" onclick={removeFile}>
						<X size={20} />
					</button>
				</div>
			{:else}
				<div class="drop-zone-content">
					<Upload size={48} class="upload-icon" />
					<p class="drop-zone-text">Drag and drop a PDF file here, or click to select</p>
					<input
						type="file"
						accept=".pdf,application/pdf"
						onchange={handleFileSelect}
						class="file-input"
						id="file-input"
					/>
					<label for="file-input" class="file-label">
						Select File
					</label>
				</div>
			{/if}
		</div>

		<div class="form-section">
			<label class="form-label">Title *</label>
			<Input
				bind:value={title}
				placeholder="Enter resume title"
				disabled={uploading}
			/>
		</div>

		<div class="form-section">
			<label class="form-label">Tags</label>
			<TagInput
				bind:tags
				{availableTags}
				placeholder="Add tags..."
			/>
		</div>

		<div class="form-actions">
			<Button onclick={() => goto('/resumes')} variant="secondary" disabled={uploading}>
				Cancel
			</Button>
			<Button onclick={handleUpload} variant="primary" disabled={uploading || !selectedFile || !title.trim()}>
				{uploading ? 'Uploading...' : 'Upload Resume'}
			</Button>
		</div>
	</div>
</div>

<ToastContainer />

<style>
	.page-header {
		margin-bottom: 2rem;
	}

	.back-button {
		padding: 0.5rem 1rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		background-color: white;
		color: #1f2937;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s;
		margin-bottom: 1rem;
	}

	.back-button:hover {
		background-color: #f9fafb;
	}

	.page-title {
		margin: 0 0 0.5rem 0;
		font-size: 1.875rem;
		font-weight: 600;
		color: #1f2937;
	}

	.page-subtitle {
		margin: 0;
		color: #6b7280;
		font-size: 0.875rem;
	}

	.upload-form {
		background-color: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 2rem;
	}

	.drop-zone {
		border: 2px dashed #e5e7eb;
		border-radius: 0.5rem;
		padding: 3rem;
		text-align: center;
		transition: all 0.2s;
		margin-bottom: 2rem;
	}

	.drop-zone.active {
		border-color: #3b82f6;
		background-color: #eff6ff;
	}

	.drop-zone-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.upload-icon {
		color: #9ca3af;
	}

	.drop-zone-text {
		color: #6b7280;
		font-size: 0.875rem;
	}

	.file-input {
		display: none;
	}

	.file-label {
		padding: 0.5rem 1rem;
		background-color: #3b82f6;
		color: white;
		border-radius: 0.375rem;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s;
	}

	.file-label:hover {
		background-color: #2563eb;
	}

	.file-selected {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background-color: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
	}

	.file-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.file-name {
		font-weight: 500;
		color: #1f2937;
	}

	.file-size {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.remove-button {
		padding: 0.5rem;
		border: none;
		background-color: transparent;
		color: #6b7280;
		cursor: pointer;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	.remove-button:hover {
		background-color: #fee2e2;
		color: #dc2626;
	}

	.form-section {
		margin-bottom: 1.5rem;
	}

	.form-label {
		display: block;
		margin-bottom: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #1f2937;
	}

	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		margin-top: 2rem;
	}
</style>
