<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, FileText, AlertCircle, Edit, Trash2, Globe, X, CheckSquare, Square } from 'lucide-svelte';
	import {
		listPosts,
		getPost,
		createPost,
		updatePost,
		deletePost,
		type Post,
		type CreatePostInput
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let posts: Post[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let showBulkActions = $state(false);
	let editingPost: Post | null = $state(null);
	let translatingPostId: string | null = $state(null);
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formTitle = $state('');
	let formContent = $state('');
	let formExcerpt = $state('');
	let formStatus = $state('draft');
	let formFeatured = $state(false);
	let formFeaturedImage = $state('');
	let formMetaTitle = $state('');
	let formMetaDescription = $state('');

	// Translation state
	let selectedLanguage: Language = $state('pt-BR');
	let selectedLanguages: Language[] = $state([]);
	let translationMode: 'single' | 'multiple' = $state('single');

	onMount(async () => {
		await fetchPosts();
	});

	async function fetchPosts() {
		loading = true;
		error = null;
		try {
			posts = await listPosts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch posts';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		editingPost = null;
		formTitle = '';
		formContent = '';
		formExcerpt = '';
		formStatus = 'draft';
		formFeatured = false;
		formFeaturedImage = '';
		formMetaTitle = '';
		formMetaDescription = '';
		showCreateModal = true;
	}

	function startEdit(post: Post) {
		editingPost = post;
		formTitle = post.title;
		formContent = post.content;
		formExcerpt = post.excerpt || '';
		formStatus = post.status;
		formFeatured = post.featured;
		formFeaturedImage = post.featured_image || '';
		formMetaTitle = post.meta_title || '';
		formMetaDescription = post.meta_description || '';
		showEditModal = true;
	}

	function cancelEdit() {
		showCreateModal = false;
		showEditModal = false;
		editingPost = null;
	}

	async function handleSave() {
		if (!formTitle.trim() || !formContent.trim()) {
			toastError('Title and content are required');
			return;
		}

		try {
			const payload: CreatePostInput = {
				title: formTitle.trim(),
				content: formContent.trim(),
				excerpt: formExcerpt.trim() || undefined,
				status: formStatus,
				featured: formFeatured,
				featured_image: formFeaturedImage.trim() || undefined,
				meta_title: formMetaTitle.trim() || undefined,
				meta_description: formMetaDescription.trim() || undefined
			};

			if (editingPost) {
				await updatePost(editingPost.id, payload);
				toastSuccess('Post updated successfully');
			} else {
				await createPost(payload);
				toastSuccess('Post created successfully');
			}

			cancelEdit();
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to save post');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this post?')) return;
		try {
			await deletePost(id);
			toastSuccess('Post deleted successfully');
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete post');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} post(s)?`)) return;

		try {
			await Promise.all(Array.from(selectedIds).map((id) => deletePost(id)));
			toastSuccess(`Deleted ${selectedIds.size} post(s) successfully`);
			selectedIds.clear();
			showBulkActions = false;
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete posts');
		}
	}

	async function handleBulkUpdate(featured: boolean) {
		if (selectedIds.size === 0) return;

		try {
			await Promise.all(Array.from(selectedIds).map((id) => updatePost(id, { featured })));
			toastSuccess(`Updated ${selectedIds.size} post(s) successfully`);
			selectedIds.clear();
			showBulkActions = false;
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update posts');
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
		if (selectedIds.size === posts.length) {
			selectedIds.clear();
		} else {
			selectedIds = new Set(posts.map((p) => p.id));
		}
		showBulkActions = selectedIds.size > 0;
	}

	function startTranslation(postId: string) {
		translatingPostId = postId;
		selectedLanguage = 'pt-BR';
		selectedLanguages = [];
		translationMode = 'single';
		showTranslationModal = true;
	}

	function cancelTranslation() {
		showTranslationModal = false;
		translatingPostId = null;
	}

	async function handleRequestTranslation() {
		if (!translatingPostId) return;

		try {
			if (translationMode === 'single') {
				await requestTranslation({
					entityType: 'post',
					entityId: translatingPostId,
					language: selectedLanguage,
					fields: ['title', 'content', 'excerpt']
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				const result = await translateEntity({
					entityType: 'post',
					entityId: translatingPostId,
					languages: languages
				});
				toastSuccess(`Queued ${result.queuedCount} translation(s)`);
			}

			cancelTranslation();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to request translation');
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="page-header">
		<div>
			<h1 class="page-title">Posts</h1>
			<p class="page-description">Manage blog posts and articles</p>
		</div>
		<button type="button" class="btn btn-primary" onclick={startCreate}>
			<Plus class="icon" />
			Create Post
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
				<button type="button" class="btn btn-sm btn-secondary" onclick={() => handleBulkUpdate(true)}>
					Mark Featured
				</button>
				<button type="button" class="btn btn-sm btn-secondary" onclick={() => handleBulkUpdate(false)}>
					Unmark Featured
				</button>
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
	{:else if posts.length === 0}
		<div class="empty-state">
			<FileText class="empty-icon" />
			<p class="empty-title">No posts found</p>
			<p class="empty-description">Create your first post to get started</p>
		</div>
	{:else}
		<div class="table-container">
			<table class="table">
				<thead>
					<tr>
						<th class="checkbox-column">
							<button
								type="button"
								class="checkbox-btn"
								onclick={toggleSelectAll}
								title="Select all"
							>
								{#if selectedIds.size === posts.length}
									<CheckSquare class="icon-sm" />
								{:else}
									<Square class="icon-sm" />
								{/if}
							</button>
						</th>
						<th>Title</th>
						<th>Status</th>
						<th>Featured</th>
						<th>Created</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each posts as post}
						<tr>
							<td class="checkbox-column">
								<button
									type="button"
									class="checkbox-btn"
									onclick={() => toggleSelect(post.id)}
								>
									{#if selectedIds.has(post.id)}
										<CheckSquare class="icon-sm" />
									{:else}
										<Square class="icon-sm" />
									{/if}
								</button>
							</td>
							<td>
								<div class="post-title">
									<span class="font-medium">{post.title}</span>
									{#if post.slug}
										<span class="text-muted">/{post.slug}</span>
									{/if}
								</div>
							</td>
							<td>
								<span class="status-badge status-{post.status}">{post.status}</span>
							</td>
							<td>
								{#if post.featured}
									<span class="status-badge status-active">Featured</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td class="text-muted">{formatDate(post.created_at)}</td>
							<td class="text-right">
								<div class="actions">
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => startTranslation(post.id)}
										title="Request Translation"
									>
										<Globe class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-primary"
										onclick={() => startEdit(post)}
										title="Edit"
									>
										<Edit class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-danger"
										onclick={() => handleDelete(post.id)}
										title="Delete"
									>
										<Trash2 class="icon-sm" />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<!-- Create/Edit Modal -->
{#if showCreateModal || showEditModal}
	<div class="modal-overlay" onclick={cancelEdit}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">{editingPost ? 'Edit Post' : 'Create Post'}</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Title *</label>
					<input type="text" bind:value={formTitle} class="input" placeholder="Post title" />
				</div>
				<div class="form-group">
					<label class="form-label">Content *</label>
					<textarea
						bind:value={formContent}
						class="input textarea"
						rows="10"
						placeholder="Markdown content"
					></textarea>
				</div>
				<div class="form-group">
					<label class="form-label">Excerpt</label>
					<textarea
						bind:value={formExcerpt}
						class="input textarea"
						rows="3"
						placeholder="Short excerpt"
					></textarea>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label class="form-label">Status</label>
						<select bind:value={formStatus} class="input">
							<option value="draft">Draft</option>
							<option value="published">Published</option>
							<option value="archived">Archived</option>
						</select>
					</div>
					<div class="form-group">
						<label class="checkbox-label">
							<input type="checkbox" bind:checked={formFeatured} />
							Featured
						</label>
					</div>
				</div>
				<div class="form-group">
					<label class="form-label">Featured Image URL</label>
					<input type="url" bind:value={formFeaturedImage} class="input" placeholder="https://..." />
				</div>
				<div class="form-group">
					<label class="form-label">Meta Title</label>
					<input type="text" bind:value={formMetaTitle} class="input" placeholder="SEO title" />
				</div>
				<div class="form-group">
					<label class="form-label">Meta Description</label>
					<textarea
						bind:value={formMetaDescription}
						class="input textarea"
						rows="2"
						placeholder="SEO description"
					></textarea>
				</div>
				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handleSave}>Save</button>
					<button type="button" class="btn btn-secondary" onclick={cancelEdit}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Translation Modal -->
{#if showTranslationModal && translatingPostId}
	<div class="modal-overlay" onclick={cancelTranslation}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Request Translation</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Translation Mode</label>
					<div class="radio-group">
						<label class="radio-label">
							<input
								type="radio"
								bind:group={translationMode}
								value="single"
								onchange={() => {
									selectedLanguages = [];
								}}
							/>
							Single Language
						</label>
						<label class="radio-label">
							<input
								type="radio"
								bind:group={translationMode}
								value="multiple"
								onchange={() => {
									selectedLanguage = 'pt-BR';
								}}
							/>
							Multiple Languages
						</label>
					</div>
				</div>

				{#if translationMode === 'single'}
					<div class="form-group">
						<label class="form-label">Language</label>
						<select bind:value={selectedLanguage} class="input">
							{#each SUPPORTED_LANGUAGES as lang}
								<option value={lang.value}>{lang.label}</option>
							{/each}
						</select>
					</div>
				{:else}
					<div class="form-group">
						<label class="form-label">Languages (leave empty for all)</label>
						<div class="checkbox-list">
							{#each SUPPORTED_LANGUAGES as lang}
								<label class="checkbox-label">
									<input
										type="checkbox"
										checked={selectedLanguages.includes(lang.value)}
										onchange={(e) => {
											if (e.currentTarget.checked) {
												selectedLanguages = [...selectedLanguages, lang.value];
											} else {
												selectedLanguages = selectedLanguages.filter((l) => l !== lang.value);
											}
										}}
									/>
									{lang.label}
								</label>
							{/each}
						</div>
					</div>
				{/if}

				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handleRequestTranslation}>
						Request Translation
					</button>
					<button type="button" class="btn btn-secondary" onclick={cancelTranslation}>Cancel</button>
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

	.btn-primary {
		background: rgba(59, 130, 246, 0.15);
		border-color: rgba(59, 130, 246, 0.4);
		color: #93c5fd;
	}

	.btn-primary:hover {
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

	.table-container {
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		overflow: hidden;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
	}

	.table thead {
		background: rgba(15, 23, 42, 0.8);
	}

	.table th {
		padding: 1rem 1.5rem;
		text-align: left;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: rgba(148, 163, 184, 0.9);
	}

	.table td {
		padding: 1rem 1.5rem;
		border-top: 1px solid rgba(71, 85, 105, 0.3);
	}

	.table tbody tr:hover {
		background: rgba(51, 65, 85, 0.2);
	}

	.text-right {
		text-align: right;
	}

	.text-muted {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.post-title {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.font-medium {
		font-weight: 500;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-active {
		background: rgba(34, 197, 94, 0.2);
		color: #86efac;
	}

	.actions {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.5rem;
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

	.checkbox-column {
		width: 3rem;
		text-align: center;
	}

	.checkbox-btn {
		background: none;
		border: none;
		cursor: pointer;
		color: rgba(148, 163, 184, 0.8);
		padding: 0.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.checkbox-btn:hover {
		color: #93c5fd;
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

	.form-row {
		display: grid;
		grid-template-columns: 1fr auto;
		gap: 1rem;
		align-items: end;
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

	.textarea {
		resize: vertical;
		min-height: 100px;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #cbd5e1;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.radio-group {
		display: flex;
		gap: 1rem;
	}

	.radio-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #cbd5e1;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.checkbox-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 200px;
		overflow-y: auto;
		padding: 0.5rem;
		background: rgba(15, 23, 42, 0.4);
		border-radius: 0.5rem;
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
</style>

