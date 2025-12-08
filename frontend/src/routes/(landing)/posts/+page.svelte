<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listPosts,
		createPost,
		updatePost,
		deletePost,
		type Post,
		type CreatePostInput,
		type UpdatePostInput,
		type PostStatus
	} from '$lib/api/posts';

	let posts: Post[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingPost: Post | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formTitle = $state('');
	let formContent = $state('');
	let formExcerpt = $state('');
	let formStatus = $state<PostStatus>('draft');
	let formFeaturedImage = $state('');
	let formMetaTitle = $state('');
	let formMetaDescription = $state('');
	let formMetaKeywords = $state('');
	let formOgTitle = $state('');
	let formOgDescription = $state('');
	let formOgImage = $state('');
	let formFeatured = $state(false);

	const statuses: PostStatus[] = ['draft', 'published', 'archived'];

	onMount(async () => {
		await fetchPosts();
	});

	async function fetchPosts() {
		loading = true;
		error = null;
		try {
			posts = await listPosts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load posts';
			console.error('Error fetching posts:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(post: Post) {
		editingPost = post;
		formTitle = post.title;
		formContent = post.content;
		formExcerpt = post.excerpt || '';
		formStatus = post.status;
		formFeaturedImage = post.featuredImage || '';
		formMetaTitle = post.metaTitle || '';
		formMetaDescription = post.metaDescription || '';
		formMetaKeywords = post.metaKeywords || '';
		formOgTitle = post.ogTitle || '';
		formOgDescription = post.ogDescription || '';
		formOgImage = post.ogImage || '';
		formFeatured = post.featured;
		showEditModal = true;
	}

	function resetForm() {
		formTitle = '';
		formContent = '';
		formExcerpt = '';
		formStatus = 'draft';
		formFeaturedImage = '';
		formMetaTitle = '';
		formMetaDescription = '';
		formMetaKeywords = '';
		formOgTitle = '';
		formOgDescription = '';
		formOgImage = '';
		formFeatured = false;
		editingPost = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formContent.trim()) {
			alert('Title and content are required');
			return;
		}

		try {
			const input: CreatePostInput = {
				title: formTitle.trim(),
				content: formContent.trim(),
				excerpt: formExcerpt.trim() || undefined,
				status: formStatus,
				featuredImage: formFeaturedImage.trim() || undefined,
				metaTitle: formMetaTitle.trim() || undefined,
				metaDescription: formMetaDescription.trim() || undefined,
				metaKeywords: formMetaKeywords.trim() || undefined,
				ogTitle: formOgTitle.trim() || undefined,
				ogDescription: formOgDescription.trim() || undefined,
				ogImage: formOgImage.trim() || undefined,
				featured: formFeatured
			};

			await createPost(input);
			showCreateModal = false;
			resetForm();
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create post');
			console.error('Error creating post:', err);
		}
	}

	async function handleUpdate() {
		if (!editingPost || !formTitle.trim() || !formContent.trim()) {
			alert('Title and content are required');
			return;
		}

		try {
			const input: UpdatePostInput = {
				title: formTitle.trim(),
				content: formContent.trim(),
				excerpt: formExcerpt.trim() || undefined,
				status: formStatus,
				featuredImage: formFeaturedImage.trim() || undefined,
				metaTitle: formMetaTitle.trim() || undefined,
				metaDescription: formMetaDescription.trim() || undefined,
				metaKeywords: formMetaKeywords.trim() || undefined,
				ogTitle: formOgTitle.trim() || undefined,
				ogDescription: formOgDescription.trim() || undefined,
				ogImage: formOgImage.trim() || undefined,
				featured: formFeatured
			};

			await updatePost(editingPost.id, input);
			showEditModal = false;
			resetForm();
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update post');
			console.error('Error updating post:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this post?')) return;

		try {
			await deletePost(id);
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete post');
			console.error('Error deleting post:', err);
		}
	}

	function filteredPosts() {
		if (!searchQuery.trim()) return posts;
		const query = searchQuery.toLowerCase();
		return posts.filter(
			(p) =>
				p.title.toLowerCase().includes(query) ||
				p.excerpt?.toLowerCase().includes(query) ||
				p.content.toLowerCase().includes(query)
		);
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="header">
		<div>
			<h1>Posts Management</h1>
			<p>Manage blog posts and content</p>
		</div>
		<button onclick={openCreateModal}>Create Post</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search posts..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredPosts().length === 0}
		<div class="empty">No posts found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Status</th>
					<th>Featured</th>
					<th>Views</th>
					<th>Published</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredPosts() as post}
					<tr>
						<td>
							<strong>{post.title}</strong>
							{#if post.excerpt}
								<br />
								<small>{post.excerpt}</small>
							{/if}
						</td>
						<td>
							<span class="status status-{post.status}">{post.status}</span>
						</td>
						<td>{post.featured ? 'Yes' : 'No'}</td>
						<td>{post.viewsCount}</td>
						<td>{formatDate(post.publishedAt)}</td>
						<td>{formatDate(post.createdAt)}</td>
						<td>
							<button onclick={() => openEditModal(post)}>Edit</button>
							<button onclick={() => handleDelete(post.id)} class="delete-btn">Delete</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="modal-overlay" onclick={() => (showCreateModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Create Post</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Content *</label>
					<textarea bind:value={formContent} rows="10"></textarea>
				</div>
				<div class="form-group">
					<label>Excerpt</label>
					<textarea bind:value={formExcerpt} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Featured Image URL</label>
					<input type="url" bind:value={formFeaturedImage} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Meta Title</label>
					<input type="text" bind:value={formMetaTitle} />
				</div>
				<div class="form-group">
					<label>Meta Description</label>
					<textarea bind:value={formMetaDescription} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Meta Keywords</label>
					<input type="text" bind:value={formMetaKeywords} />
				</div>
				<div class="form-group">
					<label>OG Title</label>
					<input type="text" bind:value={formOgTitle} />
				</div>
				<div class="form-group">
					<label>OG Description</label>
					<textarea bind:value={formOgDescription} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>OG Image URL</label>
					<input type="url" bind:value={formOgImage} />
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>Create</button>
					<button onclick={() => (showCreateModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingPost}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Post</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Content *</label>
					<textarea bind:value={formContent} rows="10"></textarea>
				</div>
				<div class="form-group">
					<label>Excerpt</label>
					<textarea bind:value={formExcerpt} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Featured Image URL</label>
					<input type="url" bind:value={formFeaturedImage} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Meta Title</label>
					<input type="text" bind:value={formMetaTitle} />
				</div>
				<div class="form-group">
					<label>Meta Description</label>
					<textarea bind:value={formMetaDescription} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Meta Keywords</label>
					<input type="text" bind:value={formMetaKeywords} />
				</div>
				<div class="form-group">
					<label>OG Title</label>
					<input type="text" bind:value={formOgTitle} />
				</div>
				<div class="form-group">
					<label>OG Description</label>
					<textarea bind:value={formOgDescription} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>OG Image URL</label>
					<input type="url" bind:value={formOgImage} />
				</div>
				<div class="form-actions">
					<button onclick={handleUpdate}>Update</button>
					<button onclick={() => (showEditModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		padding: 1rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.header h1 {
		margin: 0 0 0.25rem 0;
		font-size: 1.5rem;
	}

	.header p {
		margin: 0;
		color: #666;
		font-size: 0.9rem;
	}

	.header button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.header button:hover {
		background: #0056b3;
	}

	.search-bar {
		margin-bottom: 1rem;
	}

	.search-input {
		width: 100%;
		max-width: 400px;
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	.error {
		padding: 0.75rem;
		background: #fee;
		color: #c33;
		border: 1px solid #fcc;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.loading,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #666;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
		background: white;
	}

	.table th,
	.table td {
		padding: 0.75rem;
		text-align: left;
		border-bottom: 1px solid #ddd;
	}

	.table th {
		background: #f5f5f5;
		font-weight: 600;
	}

	.table tbody tr:hover {
		background: #f9f9f9;
	}

	.table button {
		padding: 0.25rem 0.75rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		margin-right: 0.5rem;
	}

	.table button:hover {
		background: #218838;
	}

	.delete-btn {
		background: #dc3545 !important;
	}

	.delete-btn:hover {
		background: #c82333 !important;
	}

	.table small {
		color: #666;
		font-size: 0.875rem;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-draft {
		background: #fff3cd;
		color: #856404;
	}

	.status-published {
		background: #d4edda;
		color: #155724;
	}

	.status-archived {
		background: #e2e3e5;
		color: #383d41;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		max-width: 800px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-large {
		max-width: 900px;
	}

	.modal h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.form-group label {
		font-weight: 500;
		font-size: 0.875rem;
	}

	.form-group input,
	.form-group textarea,
	.form-group select {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
	}

	.form-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.form-actions button {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.form-actions button:first-child {
		background: #007bff;
		color: white;
	}

	.form-actions button:first-child:hover {
		background: #0056b3;
	}

	.form-actions button:last-child {
		background: #6c757d;
		color: white;
	}

	.form-actions button:last-child:hover {
		background: #5a6268;
	}
</style>

