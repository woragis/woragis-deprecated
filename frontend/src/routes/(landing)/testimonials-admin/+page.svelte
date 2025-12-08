<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listTestimonials,
		createTestimonial,
		updateTestimonial,
		deleteTestimonial,
		approveTestimonial,
		rejectTestimonial,
		hideTestimonial,
		type Testimonial,
		type CreateTestimonialInput,
		type UpdateTestimonialInput,
		type TestimonialStatus,
		type TestimonialType
	} from '$lib/api/testimonials';

	let testimonials: Testimonial[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingTestimonial: Testimonial | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formAuthorName = $state('');
	let formAuthorRole = $state('');
	let formAuthorCompany = $state('');
	let formAuthorPhoto = $state('');
	let formContent = $state('');
	let formContext = $state('');
	let formVideoUrl = $state('');
	let formType = $state<TestimonialType>('general');
	let formRating = $state<number | ''>('');
	let formLinkedinUrl = $state('');
	let formStatus = $state<TestimonialStatus>('pending');
	let formDisplayOrder = $state<number | ''>(0);

	const statuses: TestimonialStatus[] = ['pending', 'approved', 'rejected', 'hidden'];
	const types: TestimonialType[] = ['general', 'project_specific', 'skill_specific'];

	onMount(async () => {
		await fetchTestimonials();
	});

	async function fetchTestimonials() {
		loading = true;
		error = null;
		try {
			testimonials = await listTestimonials();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load testimonials';
			console.error('Error fetching testimonials:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(testimonial: Testimonial) {
		editingTestimonial = testimonial;
		formAuthorName = testimonial.authorName;
		formAuthorRole = testimonial.authorRole || '';
		formAuthorCompany = testimonial.authorCompany || '';
		formAuthorPhoto = testimonial.authorPhoto || '';
		formContent = testimonial.content;
		formContext = testimonial.context || '';
		formVideoUrl = testimonial.videoUrl || '';
		formType = testimonial.type;
		formRating = testimonial.rating || '';
		formLinkedinUrl = testimonial.linkedinUrl || '';
		formStatus = testimonial.status;
		formDisplayOrder = testimonial.displayOrder;
		showEditModal = true;
	}

	function resetForm() {
		formAuthorName = '';
		formAuthorRole = '';
		formAuthorCompany = '';
		formAuthorPhoto = '';
		formContent = '';
		formContext = '';
		formVideoUrl = '';
		formType = 'general';
		formRating = '';
		formLinkedinUrl = '';
		formStatus = 'pending';
		formDisplayOrder = 0;
		editingTestimonial = null;
	}

	async function handleCreate() {
		if (!formAuthorName.trim() || !formContent.trim()) {
			alert('Author name and content are required');
			return;
		}

		try {
			const input: CreateTestimonialInput = {
				authorName: formAuthorName.trim(),
				authorRole: formAuthorRole.trim() || undefined,
				authorCompany: formAuthorCompany.trim() || undefined,
				authorPhoto: formAuthorPhoto.trim() || undefined,
				content: formContent.trim(),
				context: formContext.trim() || undefined,
				videoUrl: formVideoUrl.trim() || undefined,
				type: formType,
				rating: formRating ? Number(formRating) : undefined,
				linkedinUrl: formLinkedinUrl.trim() || undefined,
				status: formStatus,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await createTestimonial(input);
			showCreateModal = false;
			resetForm();
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create testimonial');
			console.error('Error creating testimonial:', err);
		}
	}

	async function handleUpdate() {
		if (!editingTestimonial || !formAuthorName.trim() || !formContent.trim()) {
			alert('Author name and content are required');
			return;
		}

		try {
			const input: UpdateTestimonialInput = {
				authorName: formAuthorName.trim(),
				authorRole: formAuthorRole.trim() || undefined,
				authorCompany: formAuthorCompany.trim() || undefined,
				authorPhoto: formAuthorPhoto.trim() || undefined,
				content: formContent.trim(),
				context: formContext.trim() || undefined,
				videoUrl: formVideoUrl.trim() || undefined,
				type: formType,
				rating: formRating ? Number(formRating) : undefined,
				linkedinUrl: formLinkedinUrl.trim() || undefined,
				status: formStatus,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await updateTestimonial(editingTestimonial.id, input);
			showEditModal = false;
			resetForm();
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update testimonial');
			console.error('Error updating testimonial:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this testimonial?')) return;

		try {
			await deleteTestimonial(id);
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete testimonial');
			console.error('Error deleting testimonial:', err);
		}
	}

	async function handleApprove(id: string) {
		try {
			await approveTestimonial(id);
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to approve testimonial');
		}
	}

	async function handleReject(id: string) {
		try {
			await rejectTestimonial(id);
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to reject testimonial');
		}
	}

	async function handleHide(id: string) {
		try {
			await hideTestimonial(id);
			await fetchTestimonials();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to hide testimonial');
		}
	}

	function filteredTestimonials() {
		if (!searchQuery.trim()) return testimonials;
		const query = searchQuery.toLowerCase();
		return testimonials.filter(
			(t) =>
				t.authorName.toLowerCase().includes(query) ||
				t.content.toLowerCase().includes(query) ||
				t.authorCompany?.toLowerCase().includes(query)
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
			<h1>Testimonials Management</h1>
			<p>Manage testimonials and reviews</p>
		</div>
		<button onclick={openCreateModal}>Create Testimonial</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search testimonials..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredTestimonials().length === 0}
		<div class="empty">No testimonials found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Author</th>
					<th>Content</th>
					<th>Type</th>
					<th>Status</th>
					<th>Rating</th>
					<th>Order</th>
					<th>Created</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredTestimonials() as testimonial}
					<tr>
						<td>
							<strong>{testimonial.authorName}</strong>
							{#if testimonial.authorRole}
								<br />
								<small>{testimonial.authorRole}</small>
							{/if}
							{#if testimonial.authorCompany}
								<br />
								<small>{testimonial.authorCompany}</small>
							{/if}
						</td>
						<td>
							<div class="content-preview">{testimonial.content.substring(0, 100)}...</div>
						</td>
						<td>{testimonial.type}</td>
						<td>
							<span class="status status-{testimonial.status}">{testimonial.status}</span>
						</td>
						<td>{testimonial.rating || '—'}</td>
						<td>{testimonial.displayOrder}</td>
						<td>{formatDate(testimonial.createdAt)}</td>
						<td>
							<div class="actions">
								<button onclick={() => openEditModal(testimonial)}>Edit</button>
								{#if testimonial.status === 'pending'}
									<button onclick={() => handleApprove(testimonial.id)} class="approve-btn">
										Approve
									</button>
									<button onclick={() => handleReject(testimonial.id)} class="reject-btn">
										Reject
									</button>
								{/if}
								{#if testimonial.status === 'approved'}
									<button onclick={() => handleHide(testimonial.id)} class="hide-btn">Hide</button>
								{/if}
								<button onclick={() => handleDelete(testimonial.id)} class="delete-btn">
									Delete
								</button>
							</div>
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
			<h2>Create Testimonial</h2>
			<div class="form">
				<div class="form-group">
					<label>Author Name *</label>
					<input type="text" bind:value={formAuthorName} />
				</div>
				<div class="form-group">
					<label>Author Role</label>
					<input type="text" bind:value={formAuthorRole} />
				</div>
				<div class="form-group">
					<label>Author Company</label>
					<input type="text" bind:value={formAuthorCompany} />
				</div>
				<div class="form-group">
					<label>Author Photo URL</label>
					<input type="url" bind:value={formAuthorPhoto} />
				</div>
				<div class="form-group">
					<label>Content *</label>
					<textarea bind:value={formContent} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Context</label>
					<textarea bind:value={formContext} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Video URL</label>
					<input type="url" bind:value={formVideoUrl} />
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
					<label>Rating (1-5)</label>
					<input type="number" min="1" max="5" bind:value={formRating} />
				</div>
				<div class="form-group">
					<label>LinkedIn URL</label>
					<input type="url" bind:value={formLinkedinUrl} />
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
{#if showEditModal && editingTestimonial}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Testimonial</h2>
			<div class="form">
				<div class="form-group">
					<label>Author Name *</label>
					<input type="text" bind:value={formAuthorName} />
				</div>
				<div class="form-group">
					<label>Author Role</label>
					<input type="text" bind:value={formAuthorRole} />
				</div>
				<div class="form-group">
					<label>Author Company</label>
					<input type="text" bind:value={formAuthorCompany} />
				</div>
				<div class="form-group">
					<label>Author Photo URL</label>
					<input type="url" bind:value={formAuthorPhoto} />
				</div>
				<div class="form-group">
					<label>Content *</label>
					<textarea bind:value={formContent} rows="5"></textarea>
				</div>
				<div class="form-group">
					<label>Context</label>
					<textarea bind:value={formContext} rows="2"></textarea>
				</div>
				<div class="form-group">
					<label>Video URL</label>
					<input type="url" bind:value={formVideoUrl} />
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
					<label>Rating (1-5)</label>
					<input type="number" min="1" max="5" bind:value={formRating} />
				</div>
				<div class="form-group">
					<label>LinkedIn URL</label>
					<input type="url" bind:value={formLinkedinUrl} />
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

	.content-preview {
		max-width: 300px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.actions {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.actions button {
		padding: 0.25rem 0.5rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.75rem;
	}

	.actions button:first-child {
		background: #28a745;
		color: white;
	}

	.actions button:first-child:hover {
		background: #218838;
	}

	.approve-btn {
		background: #17a2b8 !important;
		color: white;
	}

	.approve-btn:hover {
		background: #138496 !important;
	}

	.reject-btn {
		background: #ffc107 !important;
		color: #000;
	}

	.reject-btn:hover {
		background: #e0a800 !important;
	}

	.hide-btn {
		background: #6c757d !important;
		color: white;
	}

	.hide-btn:hover {
		background: #5a6268 !important;
	}

	.delete-btn {
		background: #dc3545 !important;
		color: white;
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

	.status-pending {
		background: #fff3cd;
		color: #856404;
	}

	.status-approved {
		background: #d4edda;
		color: #155724;
	}

	.status-rejected {
		background: #f8d7da;
		color: #721c24;
	}

	.status-hidden {
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

