<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Upload,
		Image,
		File,
		Trash2,
		Edit,
		Search,
		AlertCircle,
		X,
		CheckSquare,
		Square
	} from 'lucide-svelte';
	import {
		listMedia,
		uploadMedia,
		updateMedia,
		deleteMedia,
		bulkDeleteMedia,
		formatFileSize,
		isImage,
		isVideo,
		type MediaFile
	} from '$lib/api/media';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let files: MediaFile[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showUploadModal = $state(false);
	let showEditModal = $state(false);
	let showBulkActions = $state(false);
	let editingFile: MediaFile | null = $state(null);
	let selectedIds = $state<Set<string>>(new Set());
	let searchQuery = $state('');
	let currentPage = $state(1);
	let total = $state(0);
	let limit = $state(20);

	// Upload state
	let uploadFile: File | null = $state(null);
	let uploadAltText = $state('');
	let uploadCaption = $state('');
	let uploadFolder = $state('');
	let uploading = $state(false);

	// Edit state
	let editAltText = $state('');
	let editCaption = $state('');
	let editFilename = $state('');

	onMount(async () => {
		await fetchMedia();
	});

	async function fetchMedia() {
		loading = true;
		error = null;
		try {
			const response = await listMedia({
				page: currentPage,
				limit: limit,
				search: searchQuery || undefined
			});
			files = response.files;
			total = response.total;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch media files';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			uploadFile = target.files[0];
		}
	}

	async function handleUpload() {
		if (!uploadFile) {
			toastError('Please select a file to upload');
			return;
		}

		uploading = true;
		try {
			await uploadMedia({
				file: uploadFile,
				alt_text: uploadAltText.trim() || undefined,
				caption: uploadCaption.trim() || undefined,
				folder: uploadFolder.trim() || undefined
			});
			toastSuccess('File uploaded successfully');
			showUploadModal = false;
			uploadFile = null;
			uploadAltText = '';
			uploadCaption = '';
			uploadFolder = '';
			await fetchMedia();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to upload file');
		} finally {
			uploading = false;
		}
	}

	function startEdit(file: MediaFile) {
		editingFile = file;
		editAltText = file.alt_text || '';
		editCaption = file.caption || '';
		editFilename = file.filename;
		showEditModal = true;
	}

	function cancelEdit() {
		showEditModal = false;
		editingFile = null;
	}

	async function handleUpdate() {
		if (!editingFile) return;

		try {
			await updateMedia(editingFile.id, {
				alt_text: editAltText.trim() || undefined,
				caption: editCaption.trim() || undefined,
				filename: editFilename.trim() || undefined
			});
			toastSuccess('Media file updated successfully');
			cancelEdit();
			await fetchMedia();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update media file');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this media file?')) return;
		try {
			await deleteMedia(id);
			toastSuccess('Media file deleted successfully');
			await fetchMedia();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete media file');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} file(s)?`)) return;

		try {
			await bulkDeleteMedia(Array.from(selectedIds));
			toastSuccess(`Deleted ${selectedIds.size} file(s) successfully`);
			selectedIds.clear();
			showBulkActions = false;
			await fetchMedia();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete files');
		}
	}

	function toggleSelect(id: string) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = new Set(selectedIds);
		showBulkActions = selectedIds.size > 0;
	}

	function toggleSelectAll() {
		if (selectedIds.size === files.length) {
			selectedIds.clear();
		} else {
			selectedIds = new Set(files.map((f) => f.id));
		}
		showBulkActions = selectedIds.size > 0;
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}

	function copyUrl(url: string) {
		navigator.clipboard.writeText(url);
		toastSuccess('URL copied to clipboard');
	}
</script>

<div class="page-container">
	<div class="page-header">
		<div>
			<h1 class="page-title">Media Library</h1>
			<p class="page-description">Manage uploaded media files and images</p>
		</div>
		<button type="button" class="btn btn-primary" onclick={() => (showUploadModal = true)}>
			<Upload class="icon" />
			Upload Media
		</button>
	</div>

	<div class="search-bar">
		<div class="search-input-wrapper">
			<Search class="search-icon" />
			<input
				type="text"
				bind:value={searchQuery}
				class="search-input"
				placeholder="Search files..."
				onkeydown={(e) => {
					if (e.key === 'Enter') {
						currentPage = 1;
						fetchMedia();
					}
				}}
			/>
			{#if searchQuery}
				<button
					type="button"
					class="clear-search"
					onclick={() => {
						searchQuery = '';
						currentPage = 1;
						fetchMedia();
					}}
				>
					<X class="icon-sm" />
				</button>
			{/if}
		</div>
		<button type="button" class="btn btn-secondary" onclick={fetchMedia}>
			Search
		</button>
	</div>

	{#if error}
		<div class="alert alert-error">
			<AlertCircle class="icon" />
			<p>{error}</p>
		</div>
	{/if}

	{#if showBulkActions && selectedIds.size > 0}
		<div class="bulk-actions-bar">
			<span class="bulk-count">{selectedIds.size} selected</span>
			<div class="bulk-buttons">
				<button type="button" class="btn btn-sm btn-danger" onclick={handleBulkDelete}>
					Delete Selected
				</button>
				<button type="button" class="btn btn-sm btn-secondary" onclick={() => {
					selectedIds.clear();
					showBulkActions = false;
				}}>
					Cancel
				</button>
			</div>
		</div>
	{/if}

	{#if loading}
		<div class="loading-container">
			<div class="spinner"></div>
		</div>
	{:else if files.length === 0}
		<div class="empty-state">
			<Image class="empty-icon" />
			<p class="empty-title">No media files found</p>
			<p class="empty-description">Upload your first media file to get started</p>
		</div>
	{:else}
		<div class="media-grid">
			{#each files as file}
				<div class="media-card">
					<div class="media-checkbox">
						<button
							type="button"
							class="checkbox-btn"
							onclick={() => toggleSelect(file.id)}
						>
							{#if selectedIds.has(file.id)}
								<CheckSquare class="icon-sm" />
							{:else}
								<Square class="icon-sm" />
							{/if}
						</button>
					</div>
					<div class="media-preview">
						{#if isImage(file.mime_type)}
							<img src={file.thumbnail_url || file.url} alt={file.alt_text || file.filename} />
						{:else if isVideo(file.mime_type)}
							<video src={file.url} muted></video>
						{:else}
							<File class="file-icon" />
						{/if}
					</div>
					<div class="media-info">
						<div class="media-filename" title={file.original_filename}>
							{file.original_filename}
						</div>
						<div class="media-meta">
							<span class="media-size">{formatFileSize(file.size)}</span>
							{#if file.width && file.height}
								<span class="media-dimensions">{file.width}×{file.height}</span>
							{/if}
						</div>
						<div class="media-actions">
							<button
								type="button"
								class="btn btn-sm btn-secondary"
								onclick={() => copyUrl(file.url)}
								title="Copy URL"
							>
								Copy URL
							</button>
							<button
								type="button"
								class="btn btn-sm btn-primary"
								onclick={() => startEdit(file)}
								title="Edit"
							>
								<Edit class="icon-sm" />
							</button>
							<button
								type="button"
								class="btn btn-sm btn-danger"
								onclick={() => handleDelete(file.id)}
								title="Delete"
							>
								<Trash2 class="icon-sm" />
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>

		{#if total > limit}
			<div class="pagination">
				<button
					type="button"
					class="btn btn-sm btn-secondary"
					disabled={currentPage === 1}
					onclick={() => {
						currentPage--;
						fetchMedia();
					}}
				>
					Previous
				</button>
				<span class="page-info">
					Page {currentPage} of {Math.ceil(total / limit)}
				</span>
				<button
					type="button"
					class="btn btn-sm btn-secondary"
					disabled={currentPage >= Math.ceil(total / limit)}
					onclick={() => {
						currentPage++;
						fetchMedia();
					}}
				>
					Next
				</button>
			</div>
		{/if}
	{/if}
</div>

<!-- Upload Modal -->
{#if showUploadModal}
	<div class="modal-overlay" onclick={() => (showUploadModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Upload Media</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">File *</label>
					<input type="file" onchange={handleFileSelect} class="input" />
					{#if uploadFile}
						<div class="file-info">
							<span>{uploadFile.name}</span>
							<span class="text-muted">({formatFileSize(uploadFile.size)})</span>
						</div>
					{/if}
				</div>
				<div class="form-group">
					<label class="form-label">Alt Text</label>
					<input type="text" bind:value={uploadAltText} class="input" placeholder="Alt text for images" />
				</div>
				<div class="form-group">
					<label class="form-label">Caption</label>
					<input type="text" bind:value={uploadCaption} class="input" placeholder="Caption" />
				</div>
				<div class="form-group">
					<label class="form-label">Folder</label>
					<input type="text" bind:value={uploadFolder} class="input" placeholder="Optional folder path" />
				</div>
				<div class="modal-actions">
					<button
						type="button"
						class="btn btn-primary"
						onclick={handleUpload}
						disabled={!uploadFile || uploading}
					>
						{uploading ? 'Uploading...' : 'Upload'}
					</button>
					<button type="button" class="btn btn-secondary" onclick={() => (showUploadModal = false)}>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingFile}
	<div class="modal-overlay" onclick={cancelEdit}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Edit Media</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Filename</label>
					<input type="text" bind:value={editFilename} class="input" />
				</div>
				<div class="form-group">
					<label class="form-label">Alt Text</label>
					<input type="text" bind:value={editAltText} class="input" placeholder="Alt text for images" />
				</div>
				<div class="form-group">
					<label class="form-label">Caption</label>
					<input type="text" bind:value={editCaption} class="input" placeholder="Caption" />
				</div>
				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handleUpdate}>Save</button>
					<button type="button" class="btn btn-secondary" onclick={cancelEdit}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
	}

	.page-title {
		font-size: 1.875rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 0.5rem;
	}

	.page-description {
		color: rgba(148, 163, 184, 0.9);
		font-size: 0.9rem;
	}

	.search-bar {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	.search-input-wrapper {
		flex: 1;
		position: relative;
		display: flex;
		align-items: center;
	}

	.search-icon {
		position: absolute;
		left: 0.75rem;
		width: 1rem;
		height: 1rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.search-input {
		width: 100%;
		padding: 0.5rem 0.75rem 0.5rem 2.5rem;
		background: rgba(15, 23, 42, 0.8);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.5rem;
		color: #f8fafc;
		font-size: 0.875rem;
	}

	.search-input:focus {
		outline: none;
		border-color: rgba(59, 130, 246, 0.6);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.clear-search {
		position: absolute;
		right: 0.5rem;
		background: none;
		border: none;
		cursor: pointer;
		color: rgba(148, 163, 184, 0.8);
		padding: 0.25rem;
		display: flex;
		align-items: center;
	}

	.clear-search:hover {
		color: #f8fafc;
	}

	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		border: 1px solid;
		transition: all 120ms ease;
		cursor: pointer;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-primary {
		background: rgba(59, 130, 246, 0.15);
		border-color: rgba(59, 130, 246, 0.4);
		color: #93c5fd;
	}

	.btn-primary:hover:not(:disabled) {
		background: rgba(59, 130, 246, 0.25);
		border-color: rgba(59, 130, 246, 0.6);
	}

	.btn-sm {
		padding: 0.375rem 0.75rem;
		font-size: 0.8rem;
	}

	.btn-danger {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.4);
		color: #fca5a5;
	}

	.btn-danger:hover {
		background: rgba(239, 68, 68, 0.25);
		border-color: rgba(239, 68, 68, 0.6);
	}

	.btn-secondary {
		background: rgba(71, 85, 105, 0.15);
		border-color: rgba(71, 85, 105, 0.4);
		color: #cbd5e1;
	}

	.btn-secondary:hover {
		background: rgba(71, 85, 105, 0.25);
		border-color: rgba(71, 85, 105, 0.6);
	}

	.icon {
		width: 1rem;
		height: 1rem;
	}

	.icon-sm {
		width: 0.875rem;
		height: 0.875rem;
	}

	.alert {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 0.5rem;
		border: 1px solid;
	}

	.alert-error {
		background: rgba(239, 68, 68, 0.1);
		border-color: rgba(239, 68, 68, 0.3);
		color: #fca5a5;
	}

	.loading-container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem 0;
	}

	.spinner {
		width: 3rem;
		height: 3rem;
		border: 2px solid rgba(71, 85, 105, 0.3);
		border-top-color: #3b82f6;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
	}

	.empty-icon {
		width: 4rem;
		height: 4rem;
		color: rgba(71, 85, 105, 0.6);
		margin: 0 auto 1rem;
	}

	.empty-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: rgba(203, 213, 225, 0.9);
		margin-bottom: 0.5rem;
	}

	.empty-description {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.bulk-actions-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem;
		background: rgba(59, 130, 246, 0.1);
		border: 1px solid rgba(59, 130, 246, 0.3);
		border-radius: 0.5rem;
	}

	.bulk-count {
		font-weight: 500;
		color: #93c5fd;
	}

	.bulk-buttons {
		display: flex;
		gap: 0.5rem;
	}

	.media-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
	}

	.media-card {
		position: relative;
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		overflow: hidden;
		transition: all 120ms ease;
	}

	.media-card:hover {
		border-color: rgba(59, 130, 246, 0.6);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
	}

	.media-checkbox {
		position: absolute;
		top: 0.5rem;
		left: 0.5rem;
		z-index: 10;
	}

	.checkbox-btn {
		background: rgba(15, 23, 42, 0.9);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.375rem;
		cursor: pointer;
		color: rgba(148, 163, 184, 0.8);
		padding: 0.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.checkbox-btn:hover {
		color: #93c5fd;
		border-color: rgba(59, 130, 246, 0.6);
	}

	.media-preview {
		width: 100%;
		aspect-ratio: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(15, 23, 42, 0.8);
		overflow: hidden;
	}

	.media-preview img,
	.media-preview video {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.file-icon {
		width: 3rem;
		height: 3rem;
		color: rgba(148, 163, 184, 0.6);
	}

	.media-info {
		padding: 0.75rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.media-filename {
		font-size: 0.875rem;
		font-weight: 500;
		color: #f8fafc;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.media-meta {
		display: flex;
		gap: 0.5rem;
		font-size: 0.75rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.media-actions {
		display: flex;
		gap: 0.25rem;
		margin-top: 0.25rem;
	}

	.pagination {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		padding: 1rem;
	}

	.page-info {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.text-muted {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.file-info {
		margin-top: 0.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.875rem;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(2, 6, 23, 0.7);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal {
		background: rgba(15, 23, 42, 0.95);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		padding: 1.5rem;
		width: 100%;
		max-width: 28rem;
		box-shadow: 0 20px 45px rgba(2, 6, 23, 0.6);
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-large {
		max-width: 42rem;
	}

	.modal-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 1rem;
	}

	.modal-content {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: rgba(203, 213, 225, 0.9);
	}

	.input {
		width: 100%;
		padding: 0.5rem 0.75rem;
		background: rgba(15, 23, 42, 0.8);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.5rem;
		color: #f8fafc;
		font-size: 0.875rem;
		font-family: inherit;
	}

	.input:focus {
		outline: none;
		border-color: rgba(59, 130, 246, 0.6);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
</style>

