<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listTechnicalWritings,
		createTechnicalWriting,
		updateTechnicalWriting,
		deleteTechnicalWriting,
		type TechnicalWriting,
		type CreateTechnicalWritingInput,
		type UpdateTechnicalWritingInput,
		type WritingType,
		type PublicationPlatform
	} from '$lib/api/technicalwritings';

	let writings: TechnicalWriting[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingWriting: TechnicalWriting | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formTitle = $state('');
	let formDescription = $state('');
	let formType = $state<WritingType>('article');
	let formPlatform = $state<PublicationPlatform>('medium');
	let formUrl = $state('');
	let formCanonicalUrl = $state('');
	let formContent = $state('');
	let formPublishedAt = $state('');
	let formReadingTime = $state<number | ''>('');
	let formTopics = $state('');
	let formTechnologies = $state('');
	let formViews = $state<number | ''>('');
	let formLikes = $state<number | ''>('');
	let formShares = $state<number | ''>('');
	let formComments = $state<number | ''>('');
	let formProjectId = $state('');
	let formCaseStudyId = $state('');
	let formFeatured = $state(false);
	let formDisplayOrder = $state<number | ''>(0);
	let formExcerpt = $state('');
	let formCoverImageUrl = $state('');

	const types: WritingType[] = [
		'article',
		'documentation',
		'tutorial',
		'guide',
		'blog_post',
		'case_study',
		'other'
	];
	const platforms: PublicationPlatform[] = [
		'medium',
		'dev_to',
		'hashnode',
		'personal_blog',
		'github',
		'company_blog',
		'substack',
		'linkedin',
		'other'
	];

	onMount(async () => {
		await fetchWritings();
	});

	async function fetchWritings() {
		loading = true;
		error = null;
		try {
			writings = await listTechnicalWritings();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load technical writings';
			console.error('Error fetching technical writings:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(writing: TechnicalWriting) {
		editingWriting = writing;
		formTitle = writing.title;
		formDescription = writing.description;
		formType = writing.type;
		formPlatform = writing.platform;
		formUrl = writing.url;
		formCanonicalUrl = writing.canonicalUrl || '';
		formContent = writing.content || '';
		formPublishedAt = writing.publishedAt ? writing.publishedAt.split('T')[0] : '';
		formReadingTime = writing.readingTime || '';
		formTopics = writing.topics?.join(', ') || '';
		formTechnologies = writing.technologies?.join(', ') || '';
		formViews = writing.views || '';
		formLikes = writing.likes || '';
		formShares = writing.shares || '';
		formComments = writing.comments || '';
		formProjectId = writing.projectId || '';
		formCaseStudyId = writing.caseStudyId || '';
		formFeatured = writing.featured;
		formDisplayOrder = writing.displayOrder;
		formExcerpt = writing.excerpt || '';
		formCoverImageUrl = writing.coverImageUrl || '';
		showEditModal = true;
	}

	function resetForm() {
		formTitle = '';
		formDescription = '';
		formType = 'article';
		formPlatform = 'medium';
		formUrl = '';
		formCanonicalUrl = '';
		formContent = '';
		formPublishedAt = '';
		formReadingTime = '';
		formTopics = '';
		formTechnologies = '';
		formViews = '';
		formLikes = '';
		formShares = '';
		formComments = '';
		formProjectId = '';
		formCaseStudyId = '';
		formFeatured = false;
		formDisplayOrder = 0;
		formExcerpt = '';
		formCoverImageUrl = '';
		editingWriting = null;
	}

	async function handleCreate() {
		if (!formTitle.trim() || !formDescription.trim() || !formUrl.trim()) {
			alert('Title, description, and URL are required');
			return;
		}

		try {
			const input: CreateTechnicalWritingInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				type: formType,
				platform: formPlatform,
				url: formUrl.trim(),
				canonicalUrl: formCanonicalUrl.trim() || undefined,
				content: formContent.trim() || undefined,
				publishedAt: formPublishedAt || undefined,
				readingTime: formReadingTime ? Number(formReadingTime) : undefined,
				topics: formTopics.trim() ? formTopics.split(',').map((s) => s.trim()).filter((s) => s) : undefined,
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: undefined,
				views: formViews ? Number(formViews) : undefined,
				likes: formLikes ? Number(formLikes) : undefined,
				shares: formShares ? Number(formShares) : undefined,
				comments: formComments ? Number(formComments) : undefined,
				projectId: formProjectId.trim() || undefined,
				caseStudyId: formCaseStudyId.trim() || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				excerpt: formExcerpt.trim() || undefined,
				coverImageUrl: formCoverImageUrl.trim() || undefined
			};

			await createTechnicalWriting(input);
			showCreateModal = false;
			resetForm();
			await fetchWritings();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create technical writing');
			console.error('Error creating technical writing:', err);
		}
	}

	async function handleUpdate() {
		if (!editingWriting || !formTitle.trim() || !formDescription.trim() || !formUrl.trim()) {
			alert('Title, description, and URL are required');
			return;
		}

		try {
			const input: UpdateTechnicalWritingInput = {
				title: formTitle.trim(),
				description: formDescription.trim(),
				type: formType,
				platform: formPlatform,
				url: formUrl.trim(),
				canonicalUrl: formCanonicalUrl.trim() || undefined,
				content: formContent.trim() || undefined,
				publishedAt: formPublishedAt || undefined,
				readingTime: formReadingTime ? Number(formReadingTime) : undefined,
				topics: formTopics.trim() ? formTopics.split(',').map((s) => s.trim()).filter((s) => s) : undefined,
				technologies: formTechnologies.trim()
					? formTechnologies.split(',').map((s) => s.trim()).filter((s) => s)
					: undefined,
				views: formViews ? Number(formViews) : undefined,
				likes: formLikes ? Number(formLikes) : undefined,
				shares: formShares ? Number(formShares) : undefined,
				comments: formComments ? Number(formComments) : undefined,
				projectId: formProjectId.trim() || undefined,
				caseStudyId: formCaseStudyId.trim() || undefined,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0,
				excerpt: formExcerpt.trim() || undefined,
				coverImageUrl: formCoverImageUrl.trim() || undefined
			};

			await updateTechnicalWriting(editingWriting.id, input);
			showEditModal = false;
			resetForm();
			await fetchWritings();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update technical writing');
			console.error('Error updating technical writing:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this technical writing?')) return;

		try {
			await deleteTechnicalWriting(id);
			await fetchWritings();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete technical writing');
			console.error('Error deleting technical writing:', err);
		}
	}

	function filteredWritings() {
		if (!searchQuery.trim()) return writings;
		const query = searchQuery.toLowerCase();
		return writings.filter(
			(w) =>
				w.title.toLowerCase().includes(query) ||
				w.description.toLowerCase().includes(query) ||
				w.platform.toLowerCase().includes(query)
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
			<h1>Technical Writings Management</h1>
			<p>Manage technical writing portfolio</p>
		</div>
		<button onclick={openCreateModal}>Create Writing</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search writings..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredWritings().length === 0}
		<div class="empty">No technical writings found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Type</th>
					<th>Platform</th>
					<th>Views</th>
					<th>Featured</th>
					<th>Published</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredWritings() as writing}
					<tr>
						<td>
							<strong>{writing.title}</strong>
							<br />
							<small>{writing.description.substring(0, 80)}...</small>
						</td>
						<td>{writing.type}</td>
						<td>{writing.platform}</td>
						<td>{writing.views || 0}</td>
						<td>{writing.featured ? 'Yes' : 'No'}</td>
						<td>{formatDate(writing.publishedAt)}</td>
						<td>
							<button onclick={() => openEditModal(writing)}>Edit</button>
							<button onclick={() => handleDelete(writing.id)} class="delete-btn">Delete</button>
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
			<h2>Create Technical Writing</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
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
					<label>URL *</label>
					<input type="url" bind:value={formUrl} />
				</div>
				<div class="form-group">
					<label>Canonical URL</label>
					<input type="url" bind:value={formCanonicalUrl} />
				</div>
				<div class="form-group">
					<label>Published Date</label>
					<input type="date" bind:value={formPublishedAt} />
				</div>
				<div class="form-group">
					<label>Reading Time (minutes)</label>
					<input type="number" bind:value={formReadingTime} />
				</div>
				<div class="form-group">
					<label>Topics (comma separated)</label>
					<input type="text" bind:value={formTopics} />
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Views</label>
					<input type="number" bind:value={formViews} />
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
					<label>Project ID</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-group">
					<label>Case Study ID</label>
					<input type="text" bind:value={formCaseStudyId} />
				</div>
				<div class="form-group">
					<label>Excerpt</label>
					<textarea bind:value={formExcerpt} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Cover Image URL</label>
					<input type="url" bind:value={formCoverImageUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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
{#if showEditModal && editingWriting}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Technical Writing</h2>
			<div class="form">
				<div class="form-group">
					<label>Title *</label>
					<input type="text" bind:value={formTitle} />
				</div>
				<div class="form-group">
					<label>Description *</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Type</label>
					<select bind:value={formType}>
						{#each types as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
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
					<label>URL *</label>
					<input type="url" bind:value={formUrl} />
				</div>
				<div class="form-group">
					<label>Canonical URL</label>
					<input type="url" bind:value={formCanonicalUrl} />
				</div>
				<div class="form-group">
					<label>Published Date</label>
					<input type="date" bind:value={formPublishedAt} />
				</div>
				<div class="form-group">
					<label>Reading Time (minutes)</label>
					<input type="number" bind:value={formReadingTime} />
				</div>
				<div class="form-group">
					<label>Topics (comma separated)</label>
					<input type="text" bind:value={formTopics} />
				</div>
				<div class="form-group">
					<label>Technologies (comma separated)</label>
					<input type="text" bind:value={formTechnologies} />
				</div>
				<div class="form-group">
					<label>Views</label>
					<input type="number" bind:value={formViews} />
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
					<label>Project ID</label>
					<input type="text" bind:value={formProjectId} />
				</div>
				<div class="form-group">
					<label>Case Study ID</label>
					<input type="text" bind:value={formCaseStudyId} />
				</div>
				<div class="form-group">
					<label>Excerpt</label>
					<textarea bind:value={formExcerpt} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Cover Image URL</label>
					<input type="url" bind:value={formCoverImageUrl} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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

