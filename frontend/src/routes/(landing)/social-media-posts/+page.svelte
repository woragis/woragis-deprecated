<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Share2, AlertCircle, Edit, Trash2, Globe, CheckSquare, Square } from 'lucide-svelte';
	import {
		listSocialMediaPosts,
		createSocialMediaPost,
		updateSocialMediaPost,
		deleteSocialMediaPost,
		type SocialMediaPost,
		type CreateSocialMediaPostInput
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let posts: SocialMediaPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let showBulkActions = $state(false);
	let editingPost: SocialMediaPost | null = $state(null);
	let translatingPostId: string | null = $state(null);
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formContent = $state('');
	let formPlatform = $state('twitter');
	let formUrl = $state('');
	let formPublishedAt = $state('');
	let formFeatured = $state(false);

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
			posts = await listSocialMediaPosts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch social media posts';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		editingPost = null;
		formContent = '';
		formPlatform = 'twitter';
		formUrl = '';
		formPublishedAt = '';
		formFeatured = false;
		showCreateModal = true;
	}

	function startEdit(post: SocialMediaPost) {
		editingPost = post;
		formContent = post.content;
		formPlatform = post.platform;
		formUrl = post.url || '';
		formPublishedAt = post.published_at ? post.published_at.split('T')[0] : '';
		formFeatured = post.featured;
		showEditModal = true;
	}

	function cancelEdit() {
		showCreateModal = false;
		showEditModal = false;
		editingPost = null;
	}

	async function handleSave() {
		if (!formContent.trim()) {
			toastError('Content is required');
			return;
		}

		try {
			const payload: CreateSocialMediaPostInput = {
				content: formContent.trim(),
				platform: formPlatform,
				url: formUrl.trim() || undefined,
				published_at: formPublishedAt || undefined,
				featured: formFeatured
			};

			if (editingPost) {
				await updateSocialMediaPost(editingPost.id, payload);
				toastSuccess('Social media post updated successfully');
			} else {
				await createSocialMediaPost(payload);
				toastSuccess('Social media post created successfully');
			}

			cancelEdit();
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to save social media post');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this social media post?')) return;
		try {
			await deleteSocialMediaPost(id);
			toastSuccess('Social media post deleted successfully');
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete social media post');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} post(s)?`)) return;

		try {
			await Promise.all(Array.from(selectedIds).map((id) => deleteSocialMediaPost(id)));
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
			await Promise.all(
				Array.from(selectedIds).map((id) => updateSocialMediaPost(id, { featured }))
			);
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
					entityType: 'social_media_post',
					entityId: translatingPostId,
					targetLanguages: [selectedLanguage]
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				// Translate to each language separately
				const results = await Promise.all(
					languages.map((lang) =>
						translateEntity({
							entityType: 'social_media_post',
							entityId: translatingPostId!,
							language: lang,
							fields: { content: '' } // Will be filled by the API
						})
					)
				);
				const result = { queuedCount: results.length };
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
			<h1 class="page-title">Social Media Posts</h1>
			<p class="page-description">Manage social media posts and content</p>
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
			<Share2 class="empty-icon" />
			<p class="empty-title">No social media posts found</p>
			<p class="empty-description">Create your first social media post to get started</p>
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
						<th>Platform</th>
						<th>Content</th>
						<th>URL</th>
						<th>Published</th>
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
								<span class="badge badge-{post.platform.toLowerCase()}">{post.platform}</span>
							</td>
							<td>
								<div class="content-preview">{post.content.substring(0, 80)}...</div>
							</td>
							<td>
								{#if post.url}
									<a href={post.url} target="_blank" class="link">View →</a>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td class="text-muted">
								{#if post.published_at}
									{formatDate(post.published_at)}
								{:else}
									—
								{/if}
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
			<h2 class="modal-title">{editingPost ? 'Edit Social Media Post' : 'Create Social Media Post'}</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Platform *</label>
					<select bind:value={formPlatform} class="input">
						<option value="twitter">Twitter/X</option>
						<option value="linkedin">LinkedIn</option>
						<option value="facebook">Facebook</option>
						<option value="instagram">Instagram</option>
						<option value="youtube">YouTube</option>
						<option value="github">GitHub</option>
						<option value="medium">Medium</option>
						<option value="dev.to">Dev.to</option>
						<option value="other">Other</option>
					</select>
				</div>
				<div class="form-group">
					<label class="form-label">Content *</label>
					<textarea
						bind:value={formContent}
						class="input textarea"
						rows="6"
						placeholder="Post content"
					></textarea>
				</div>
				<div class="form-group">
					<label class="form-label">URL</label>
					<input type="url" bind:value={formUrl} class="input" placeholder="https://..." />
				</div>
				<div class="form-row">
					<div class="form-group">
						<label class="form-label">Published Date</label>
						<input type="date" bind:value={formPublishedAt} class="input" />
					</div>
					<div class="form-group">
						<label class="checkbox-label">
							<input type="checkbox" bind:checked={formFeatured} />
							Featured
						</label>
					</div>
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
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.12);
		color: #d4d4d4;
	}

	.btn-primary:hover {
		background: rgba(255, 255, 255, 0.12);
		border-color: rgba(255, 255, 255, 0.2);
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
		border-color: rgba(255, 255, 255, 0.08);
		color: #cbd5e1;
	}

	.btn-secondary:hover {
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.12);
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
		border: 2px solid rgba(255, 255, 255, 0.06);
		border-top-color: #737373;
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
		color: rgba(255, 255, 255, 0.12);
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
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0.5rem;
	}

	.bulk-count {
		font-weight: 500;
		color: #d4d4d4;
	}

	.bulk-buttons {
		display: flex;
		gap: 0.5rem;
	}

	.table-container {
		background: rgba(15, 15, 15, 0.4);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.75rem;
		overflow: hidden;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
	}

	.table thead {
		background: rgba(15, 15, 15, 0.6);
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
		border-top: 1px solid rgba(255, 255, 255, 0.06);
	}

	.table tbody tr:hover {
		background: rgba(255, 255, 255, 0.03);
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
		color: #d4d4d4;
	}

	.text-right {
		text-align: right;
	}

	.text-muted {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.content-preview {
		max-width: 300px;
		color: rgba(203, 213, 225, 0.9);
		font-size: 0.875rem;
	}

	.badge {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 500;
		text-transform: capitalize;
	}

	.badge-twitter {
		background: rgba(29, 161, 242, 0.2);
		color: #d4d4d4;
	}

	.badge-linkedin {
		background: rgba(0, 119, 181, 0.2);
		color: #a3a3a3;
	}

	.badge-facebook {
		background: rgba(24, 119, 242, 0.2);
		color: #d4d4d4;
	}

	.badge-instagram {
		background: rgba(225, 48, 108, 0.2);
		color: #f9a8d4;
	}

	.badge-youtube {
		background: rgba(255, 0, 0, 0.2);
		color: #fca5a5;
	}

	.badge-github {
		background: rgba(36, 41, 47, 0.2);
		color: #cbd5e1;
	}

	.badge-medium {
		background: rgba(0, 0, 0, 0.2);
		color: #cbd5e1;
	}

	.badge-dev.to {
		background: rgba(10, 10, 10, 0.2);
		color: #cbd5e1;
	}

	.badge-other {
		background: rgba(255, 255, 255, 0.05);
		color: #cbd5e1;
	}

	.link {
		color: #d4d4d4;
		text-decoration: none;
		font-size: 0.875rem;
	}

	.link:hover {
		text-decoration: underline;
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

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.75);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal {
		background: rgba(15, 15, 15, 0.98);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.75rem;
		padding: 1.5rem;
		width: 100%;
		max-width: 28rem;
		box-shadow: 0 20px 45px rgba(0, 0, 0, 0.8);
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
		background: rgba(15, 15, 15, 0.6);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.5rem;
		color: #f8fafc;
		font-size: 0.875rem;
		font-family: inherit;
	}

	.input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.2);
		box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.05);
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
		background: rgba(15, 15, 15, 0.3);
		border-radius: 0.5rem;
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
</style>
