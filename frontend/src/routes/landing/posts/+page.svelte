<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Plus,
		FileText,
		AlertCircle,
		Edit,
		Trash2,
		Globe,
		X,
		CheckSquare,
		Square,
		Link,
		Calendar,
		Send,
		Clock,
		Tag as TagIcon,
		Folder,
		Code
	} from 'lucide-svelte';
	import {
		listPosts,
		getPost,
		createPost,
		updatePost,
		deletePost,
		getPostSkills,
		attachSkillToPost,
		detachSkillFromPost,
		listCategories,
		createCategory,
		getPostCategories,
		attachCategoryToPost,
		detachCategoryFromPost,
		listTags,
		getPostTags,
		attachTagToPost,
		detachTagFromPost,
		listSkills,
		type Post,
		type CreatePostInput,
		type Category,
		type Tag,
		type Skill
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let posts: Post[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let showRelationshipsModal = $state(false);
	let showPublishModal = $state(false);
	let showBulkActions = $state(false);
	let editingPost: Post | null = $state(null);
	let translatingPostId: string | null = $state(null);
	let managingRelationshipsPostId: string | null = $state(null);
	let publishingPostId: string | null = $state(null);
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
	let formScheduledPublish = $state(false);
	let formPublishDate = $state('');
	let formPublishTime = $state('');

	// Relationship state
	let availableSkills: Skill[] = $state([]);
	let availableCategories: Category[] = $state([]);
	let availableTags: Tag[] = $state([]);
	let postSkills: string[] = $state([]);
	let postCategories: Category[] = $state([]);
	let postTags: Tag[] = $state([]);
	let newCategoryName = $state('');
	let newTagName = $state('');
	let relationshipLoading = $state(false);

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

	async function startEdit(post: Post) {
		editingPost = post;
		formTitle = post.title;
		formContent = post.content;
		formExcerpt = post.excerpt || '';
		formStatus = post.status;
		formFeatured = post.featured;
		formFeaturedImage = post.featured_image || '';
		formMetaTitle = post.meta_title || '';
		formMetaDescription = post.meta_description || '';
		formScheduledPublish = false;
		formPublishDate = '';
		formPublishTime = '';
		showEditModal = true;
	}

	async function startManageRelationships(postId: string) {
		managingRelationshipsPostId = postId;
		relationshipLoading = true;
		try {
			// Load all available options
			[availableSkills, availableCategories, availableTags] = await Promise.all([
				listSkills(),
				listCategories(),
				listTags()
			]);

			// Load current relationships
			[postSkills, postCategories, postTags] = await Promise.all([
				getPostSkills(postId),
				getPostCategories(postId),
				getPostTags(postId)
			]);
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to load relationships');
		} finally {
			relationshipLoading = false;
		}
		showRelationshipsModal = true;
	}

	function cancelRelationships() {
		showRelationshipsModal = false;
		managingRelationshipsPostId = null;
		postSkills = [];
		postCategories = [];
		postTags = [];
	}

	async function toggleSkill(skillId: string) {
		if (!managingRelationshipsPostId) return;

		const isAttached = postSkills.includes(skillId);
		try {
			if (isAttached) {
				await detachSkillFromPost(managingRelationshipsPostId, skillId);
				postSkills = postSkills.filter((id) => id !== skillId);
			} else {
				await attachSkillToPost(managingRelationshipsPostId, skillId);
				postSkills = [...postSkills, skillId];
			}
			toastSuccess(isAttached ? 'Skill removed' : 'Skill attached');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update skill');
		}
	}

	async function toggleCategory(categoryId: string) {
		if (!managingRelationshipsPostId) return;

		const isAttached = postCategories.some((c) => c.id === categoryId);
		try {
			if (isAttached) {
				await detachCategoryFromPost(managingRelationshipsPostId, categoryId);
				postCategories = postCategories.filter((c) => c.id !== categoryId);
			} else {
				await attachCategoryToPost(managingRelationshipsPostId, categoryId);
				const category = availableCategories.find((c) => c.id === categoryId);
				if (category) {
					postCategories = [...postCategories, category];
				}
			}
			toastSuccess(isAttached ? 'Category removed' : 'Category attached');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update category');
		}
	}

	async function toggleTag(tagId: string) {
		if (!managingRelationshipsPostId) return;

		const isAttached = postTags.some((t) => t.id === tagId);
		try {
			if (isAttached) {
				await detachTagFromPost(managingRelationshipsPostId, tagId);
				postTags = postTags.filter((t) => t.id !== tagId);
			} else {
				await attachTagToPost(managingRelationshipsPostId, tagId);
				const tag = availableTags.find((t) => t.id === tagId);
				if (tag) {
					postTags = [...postTags, tag];
				}
			}
			toastSuccess(isAttached ? 'Tag removed' : 'Tag attached');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update tag');
		}
	}

	async function createNewCategory() {
		if (!newCategoryName.trim() || !managingRelationshipsPostId) return;

		try {
			const category = await createCategory(newCategoryName.trim());
			availableCategories = [...availableCategories, category];
			await attachCategoryToPost(managingRelationshipsPostId, category.id);
			postCategories = [...postCategories, category];
			newCategoryName = '';
			toastSuccess('Category created and attached');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to create category');
		}
	}

	async function createNewTag() {
		if (!newTagName.trim() || !managingRelationshipsPostId) return;

		try {
			// Tags are created on-the-fly when attaching, so we need to check if it exists first
			// For now, we'll just try to attach and let the backend handle creation
			// This is a simplified approach - in reality, you'd need a createTag endpoint
			toastError('Tag creation not yet implemented - please use existing tags');
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to create tag');
		}
	}

	function startPublish(post: Post) {
		publishingPostId = post.id;
		formStatus = post.status;
		
		// If post has a published_at date, pre-fill the scheduled publish fields
		if (post.published_at) {
			const pubDate = new Date(post.published_at);
			formScheduledPublish = true;
			formPublishDate = pubDate.toISOString().split('T')[0];
			formPublishTime = pubDate.toTimeString().slice(0, 5);
		} else {
			formScheduledPublish = false;
			formPublishDate = '';
			formPublishTime = '';
		}
		showPublishModal = true;
	}

	function cancelPublish() {
		showPublishModal = false;
		publishingPostId = null;
	}

	async function handlePublish() {
		if (!publishingPostId) return;

		try {
			const updates: Partial<CreatePostInput> = {
				status: formStatus
			};

			// Handle scheduled publishing
			if (formScheduledPublish && formPublishDate && formPublishTime) {
				const publishDateTime = new Date(`${formPublishDate}T${formPublishTime}`);
				// Format as ISO string for backend
				// Note: Backend should handle this in the updatePost endpoint
				// For now, we'll include it in the payload
				updates.published_at = publishDateTime.toISOString();
			} else if (formStatus === 'published' && !formScheduledPublish) {
				// If publishing now, set published_at to current time
				updates.published_at = new Date().toISOString();
			}

			await updatePost(publishingPostId, updates);
			toastSuccess(
				formScheduledPublish
					? 'Post scheduled for publishing'
					: `Post ${formStatus === 'published' ? 'published' : 'status updated'} successfully`
			);
			cancelPublish();
			await fetchPosts();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update post status');
		}
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
					targetLanguages: [selectedLanguage]
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				// Translate to each language separately
				const results = await Promise.all(
					languages.map((lang) =>
						translateEntity({
							entityType: 'post',
							entityId: translatingPostId!,
							language: lang,
							fields: { title: '', content: '', excerpt: '' } // Will be filled by the API
						})
					)
				);
				toastSuccess(`Queued ${results.length} translation(s)`);
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
									<a
										href="/landing/posts/{post.id}"
										class="font-medium hover:text-blue-400 transition-colors"
									>
										{post.title}
									</a>
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
										onclick={() => startManageRelationships(post.id)}
										title="Manage Relationships"
									>
										<Link class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => startPublish(post)}
										title="Publish/Status"
									>
										<Send class="icon-sm" />
									</button>
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

<!-- Relationships Management Modal -->
{#if showRelationshipsModal && managingRelationshipsPostId}
	<div class="modal-overlay" onclick={cancelRelationships}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Manage Relationships</h2>
			<div class="modal-content">
				{#if relationshipLoading}
					<div class="loading-container">
						<div class="spinner"></div>
					</div>
				{:else}
					<!-- Skills Section -->
					<div class="relationship-section">
						<div class="section-header">
							<Code class="icon" />
							<h3 class="section-title">Skills</h3>
						</div>
						<div class="relationship-list">
							{#each availableSkills as skill}
								<label class="relationship-item">
									<input
										type="checkbox"
										checked={postSkills.includes(skill.id)}
										onchange={() => toggleSkill(skill.id)}
									/>
									<span>{skill.name}</span>
									{#if skill.category}
										<span class="relationship-meta">{skill.category}</span>
									{/if}
								</label>
							{/each}
						</div>
					</div>

					<!-- Categories Section -->
					<div class="relationship-section">
						<div class="section-header">
							<Folder class="icon" />
							<h3 class="section-title">Categories</h3>
						</div>
						<div class="create-new-item">
							<input
								type="text"
								bind:value={newCategoryName}
								class="input input-sm"
								placeholder="New category name"
								onkeydown={(e) => {
									if (e.key === 'Enter') {
										createNewCategory();
									}
								}}
							/>
							<button type="button" class="btn btn-sm btn-primary" onclick={createNewCategory}>
								<Plus class="icon-sm" />
								Add
							</button>
						</div>
						<div class="relationship-list">
							{#each availableCategories as category}
								<label class="relationship-item">
									<input
										type="checkbox"
										checked={postCategories.some((c) => c.id === category.id)}
										onchange={() => toggleCategory(category.id)}
									/>
									<span>{category.name}</span>
									{#if category.description}
										<span class="relationship-meta">{category.description}</span>
									{/if}
								</label>
							{/each}
						</div>
					</div>

					<!-- Tags Section -->
					<div class="relationship-section">
						<div class="section-header">
							<TagIcon class="icon" />
							<h3 class="section-title">Tags</h3>
						</div>
						<div class="relationship-list">
							{#each availableTags as tag}
								<label class="relationship-item">
									<input
										type="checkbox"
										checked={postTags.some((t) => t.id === tag.id)}
										onchange={() => toggleTag(tag.id)}
									/>
									<span>{tag.name}</span>
								</label>
							{/each}
						</div>
					</div>

					<div class="modal-actions">
						<button type="button" class="btn btn-secondary" onclick={cancelRelationships}>Done</button>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<!-- Publishing Workflow Modal -->
{#if showPublishModal && publishingPostId}
	<div class="modal-overlay" onclick={cancelPublish}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Publishing Workflow</h2>
			<div class="modal-content">
				<div class="workflow-steps">
					<div class="workflow-step {formStatus === 'draft' ? 'active' : formStatus === 'published' || formStatus === 'archived' ? 'completed' : ''}">
						<div class="step-indicator"></div>
						<div class="step-content">
							<h4 class="step-title">Draft</h4>
							<p class="step-description">Work in progress</p>
						</div>
					</div>
					<div class="workflow-arrow">→</div>
					<div class="workflow-step {formStatus === 'published' ? 'active' : formStatus === 'archived' ? 'completed' : ''}">
						<div class="step-indicator"></div>
						<div class="step-content">
							<h4 class="step-title">Published</h4>
							<p class="step-description">Live on site</p>
						</div>
					</div>
					<div class="workflow-arrow">→</div>
					<div class="workflow-step {formStatus === 'archived' ? 'active' : ''}">
						<div class="step-indicator"></div>
						<div class="step-content">
							<h4 class="step-title">Archived</h4>
							<p class="step-description">No longer active</p>
						</div>
					</div>
				</div>

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
						<input type="checkbox" bind:checked={formScheduledPublish} />
						Schedule Publishing
					</label>
				</div>

				{#if formScheduledPublish}
					<div class="form-row">
						<div class="form-group">
							<label class="form-label">Publish Date</label>
							<input type="date" bind:value={formPublishDate} class="input" />
						</div>
						<div class="form-group">
							<label class="form-label">Publish Time</label>
							<input type="time" bind:value={formPublishTime} class="input" />
						</div>
					</div>
				{/if}

				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handlePublish}>
						<Send class="icon-sm" />
						Update Status
					</button>
					<button type="button" class="btn btn-secondary" onclick={cancelPublish}>Cancel</button>
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
		color: #e5e5e5;
	}

	.btn-primary:hover {
		background: rgba(255, 255, 255, 0.12);
		border-color: rgba(255, 255, 255, 0.18);
		color: #f5f5f5;
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
		border: 1px solid rgba(255, 255, 255, 0.1);
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
		border: 1px solid rgba(255, 255, 255, 0.1);
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

	.relationship-section {
		margin-bottom: 1.5rem;
		padding-bottom: 1.5rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.06);
	}

	.relationship-section:last-child {
		border-bottom: none;
		margin-bottom: 0;
		padding-bottom: 0;
	}

	.section-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.section-title {
		font-size: 1rem;
		font-weight: 600;
		color: #f5f5f5;
		margin: 0;
	}

	.relationship-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 200px;
		overflow-y: auto;
		padding: 0.5rem;
		background: rgba(15, 15, 15, 0.3);
		border-radius: 0.5rem;
	}

	.relationship-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: background 120ms ease;
	}

	.relationship-item:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.relationship-item input[type="checkbox"] {
		cursor: pointer;
	}

	.relationship-item span:first-of-type {
		flex: 1;
		color: #e5e5e5;
		font-size: 0.875rem;
	}

	.relationship-meta {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.4);
		padding: 0.125rem 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 0.25rem;
	}

	.create-new-item {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 0.75rem;
	}

	.input-sm {
		flex: 1;
		padding: 0.375rem 0.5rem;
		font-size: 0.875rem;
	}

	.workflow-steps {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1.5rem;
		padding: 1rem;
		background: rgba(15, 15, 15, 0.3);
		border-radius: 0.5rem;
	}

	.workflow-step {
		flex: 1;
		display: flex;
		align-items: center;
		gap: 0.75rem;
		opacity: 0.5;
		transition: opacity 200ms ease;
	}

	.workflow-step.active {
		opacity: 1;
	}

	.workflow-step.completed {
		opacity: 0.7;
	}

	.step-indicator {
		width: 2rem;
		height: 2rem;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.1);
		border: 2px solid rgba(255, 255, 255, 0.2);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.workflow-step.active .step-indicator {
		background: rgba(255, 255, 255, 0.15);
		border-color: rgba(255, 255, 255, 0.3);
	}

	.workflow-step.completed .step-indicator {
		background: rgba(34, 197, 94, 0.2);
		border-color: rgba(34, 197, 94, 0.4);
	}

	.workflow-step.completed .step-indicator::after {
		content: '✓';
		color: #86efac;
		font-size: 0.875rem;
	}

	.step-content {
		flex: 1;
	}

	.step-title {
		font-size: 0.875rem;
		font-weight: 600;
		color: #f5f5f5;
		margin: 0 0 0.25rem 0;
	}

	.step-description {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.5);
		margin: 0;
	}

	.workflow-arrow {
		color: rgba(255, 255, 255, 0.3);
		font-size: 1.25rem;
		margin: 0 0.5rem;
	}
</style>

