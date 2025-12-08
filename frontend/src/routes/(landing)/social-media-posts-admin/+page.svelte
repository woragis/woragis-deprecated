<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listSocialMediaPosts,
		createSocialMediaPost,
		updateSocialMediaPost,
		deleteSocialMediaPost,
		type SocialMediaPost,
		type CreateSocialMediaPostInput,
		type UpdateSocialMediaPostInput,
		type Platform,
		type PostStatus
	} from '$lib/api/socialmediaposts';

	let posts: SocialMediaPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingPost: SocialMediaPost | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formUrl = $state('');
	let formPlatform = $state<Platform>('linkedin');
	let formTitle = $state('');
	let formContentPreview = $state('');
	let formPublishedDate = $state('');
	let formLikes = $state<number | ''>('');
	let formShares = $state<number | ''>('');
	let formComments = $state<number | ''>('');
	let formViews = $state<number | ''>('');
	let formStatus = $state<PostStatus>('active');

	const platforms: Platform[] = ['linkedin', 'twitter', 'instagram'];
	const statuses: PostStatus[] = ['active', 'deleted', 'unavailable'];

	onMount(async () => {
		await fetchPosts();
	});

	async function fetchPosts() {
		loading = true;
		error = null;
		try {
			posts = await listSocialMediaPosts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load social media posts';
			console.error('Error fetching social media posts:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(post: SocialMediaPost) {
		editingPost = post;
		formUrl = post.url;
		formPlatform = post.platform;
		formTitle = post.title || '';
		formContentPreview = post.contentPreview || '';
		formPublishedDate = post.publishedDate ? post.publishedDate.split('T')[0] : '';
		formLikes = post.likes || '';
		formShares = post.shares || '';
		formComments = post.comments || '';
		formViews = post.views || '';
		formStatus = post.status;
		showEditModal = true;
	}

	function resetForm() {
		formUrl = '';
		formPlatform = 'linkedin';
		formTitle = '';
		formContentPreview = '';
		formPublishedDate = '';
		formLikes = '';
		formShares = '';
		formComments = '';
		formViews = '';
		formStatus = 'active';
		editingPost = null;
	}

	async function handleCreate() {
		if (!formUrl.trim()) {
			alert('URL is required');
			return;
		}

		try {
			const input: CreateSocialMediaPostInput = {
				url: formUrl.trim(),
				platform: formPlatform,
				title: formTitle.trim() || undefined,
				contentPreview: formContentPreview.trim() || undefined,
				publishedDate: formPublishedDate || undefined,
				likes: formLikes ? Number(formLikes) : undefined,
				shares: formShares ? Number(formShares) : undefined,
				comments: formComments ? Number(formComments) : undefined,
				views: formViews ? Number(formViews) : undefined,
				status: formStatus
			};

			await createSocialMediaPost(input);
			showCreateModal = false;
			resetForm();
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create social media post');
			console.error('Error creating social media post:', err);
		}
	}

	async function handleUpdate() {
		if (!editingPost || !formUrl.trim()) {
			alert('URL is required');
			return;
		}

		try {
			const input: UpdateSocialMediaPostInput = {
				url: formUrl.trim(),
				platform: formPlatform,
				title: formTitle.trim() || undefined,
				contentPreview: formContentPreview.trim() || undefined,
				publishedDate: formPublishedDate || undefined,
				likes: formLikes ? Number(formLikes) : undefined,
				shares: formShares ? Number(formShares) : undefined,
				comments: formComments ? Number(formComments) : undefined,
				views: formViews ? Number(formViews) : undefined,
				status: formStatus
			};

			await updateSocialMediaPost(editingPost.id, input);
			showEditModal = false;
			resetForm();
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update social media post');
			console.error('Error updating social media post:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this social media post?')) return;

		try {
			await deleteSocialMediaPost(id);
			await fetchPosts();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete social media post');
			console.error('Error deleting social media post:', err);
		}
	}

	function filteredPosts() {
		if (!searchQuery.trim()) return posts;
		const query = searchQuery.toLowerCase();
		return posts.filter(
			(p) =>
				p.url.toLowerCase().includes(query) ||
				p.title?.toLowerCase().includes(query) ||
				p.platform.toLowerCase().includes(query)
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
			<h1>Social Media Posts Management</h1>
			<p>Manage social media posts</p>
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
		<div class="empty">No social media posts found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Platform</th>
					<th>Title/Preview</th>
					<th>Likes</th>
					<th>Shares</th>
					<th>Status</th>
					<th>Published</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredPosts() as post}
					<tr>
						<td>
							<span class="platform platform-{post.platform}">{post.platform}</span>
						</td>
						<td>
							<strong>{post.title || 'No title'}</strong>
							{#if post.contentPreview}
								<br />
								<small>{post.contentPreview.substring(0, 60)}...</small>
							{/if}
						</td>
						<td>{post.likes || 0}</td>
						<td>{post.shares || 0}</td>
						<td>
							<span class="status status-{post.status}">{post.status}</span>
						</td>
						<td>{formatDate(post.publishedDate)}</td>
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
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Create Social Media Post</h2>
			<div class="form">
				<div class="form-group">
					<label>URL *</label>
					<input type="url" bind:value={formUrl} />
				</div>
				<div class="form-group">
					<label>Platform</label>
					<select bind:value={formPlatform}>
						{#each platforms as platform}
							<option value={platform}>{platform}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Title</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Content Preview</label>
					<textarea bind:value={formContentPreview} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Published Date</label>
					<input type="date" bind:value={formPublishedDate} />
				</div>
				<div class="form-group">
					<label>Likes</label>
					<input type="number" bind:value={formLikes} />
				</div>
				<div class="form-group">
					<label>Shares</label>
					<input type="number" bind:value={formShares} />
				</div>
				<div class="form-group">
					<label>Comments</label>
					<input type="number" bind:value={formComments} />
				</div>
				<div class="form-group">
					<label>Views</label>
					<input type="number" bind:value={formViews} />
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
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
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Social Media Post</h2>
			<div class="form">
				<div class="form-group">
					<label>URL *</label>
					<input type="url" bind:value={formUrl} />
				</div>
				<div class="form-group">
					<label>Platform</label>
					<select bind:value={formPlatform}>
						{#each platforms as platform}
							<option value={platform}>{platform}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Title</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Content Preview</label>
					<textarea bind:value={formContentPreview} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Published Date</label>
					<input type="date" bind:value={formPublishedDate} />
				</div>
				<div class="form-group">
					<label>Likes</label>
					<input type="number" bind:value={formLikes} />
				</div>
				<div class="form-group">
					<label>Shares</label>
					<input type="number" bind:value={formShares} />
				</div>
				<div class="form-group">
					<label>Comments</label>
					<input type="number" bind:value={formComments} />
				</div>
				<div class="form-group">
					<label>Views</label>
					<input type="number" bind:value={formViews} />
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
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

	.platform {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
		text-transform: capitalize;
	}

	.platform-linkedin {
		background: #0077b5;
		color: white;
	}

	.platform-twitter {
		background: #1da1f2;
		color: white;
	}

	.platform-instagram {
		background: #e4405f;
		color: white;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-active {
		background: #d4edda;
		color: #155724;
	}

	.status-deleted {
		background: #f8d7da;
		color: #721c24;
	}

	.status-unavailable {
		background: #fff3cd;
		color: #856404;
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
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
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

