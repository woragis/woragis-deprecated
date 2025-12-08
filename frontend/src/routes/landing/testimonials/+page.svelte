<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Users, AlertCircle, Edit, Trash2, Globe, Star } from 'lucide-svelte';
	import {
		listTestimonials,
		createTestimonial,
		updateTestimonial,
		deleteTestimonial,
		type Testimonial,
		type CreateTestimonialInput
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let testimonials: Testimonial[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let editingTestimonial: Testimonial | null = $state(null);
	let translatingTestimonialId: string | null = $state(null);

	let formContent = $state('');
	let formAuthorName = $state('');
	let formAuthorTitle = $state('');
	let formAuthorCompany = $state('');
	let formAuthorImage = $state('');
	let formRating = $state<number | null>(null);
	let formFeatured = $state(false);

	let selectedLanguage: Language = $state('pt-BR');
	let selectedLanguages: Language[] = $state([]);
	let translationMode: 'single' | 'multiple' = $state('single');

	onMount(async () => {
		await fetchTestimonials();
	});

	async function fetchTestimonials() {
		loading = true;
		error = null;
		try {
			testimonials = await listTestimonials();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch testimonials';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		editingTestimonial = null;
		formContent = '';
		formAuthorName = '';
		formAuthorTitle = '';
		formAuthorCompany = '';
		formAuthorImage = '';
		formRating = null;
		formFeatured = false;
		showCreateModal = true;
	}

	function startEdit(testimonial: Testimonial) {
		editingTestimonial = testimonial;
		formContent = testimonial.content;
		formAuthorName = testimonial.author_name;
		formAuthorTitle = testimonial.author_title || '';
		formAuthorCompany = testimonial.author_company || '';
		formAuthorImage = testimonial.author_image || '';
		formRating = testimonial.rating || null;
		formFeatured = testimonial.featured;
		showEditModal = true;
	}

	function cancelEdit() {
		showCreateModal = false;
		showEditModal = false;
		editingTestimonial = null;
	}

	async function handleSave() {
		if (!formContent.trim() || !formAuthorName.trim()) {
			toastError('Content and author name are required');
			return;
		}

		try {
			const payload: CreateTestimonialInput = {
				content: formContent.trim(),
				author_name: formAuthorName.trim(),
				author_title: formAuthorTitle.trim() || undefined,
				author_company: formAuthorCompany.trim() || undefined,
				author_image: formAuthorImage.trim() || undefined,
				rating: formRating || undefined,
				featured: formFeatured
			};

			if (editingTestimonial) {
				await updateTestimonial(editingTestimonial.id, payload);
				toastSuccess('Testimonial updated successfully');
			} else {
				await createTestimonial(payload);
				toastSuccess('Testimonial created successfully');
			}

			cancelEdit();
			await fetchTestimonials();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to save testimonial');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this testimonial?')) return;
		try {
			await deleteTestimonial(id);
			toastSuccess('Testimonial deleted successfully');
			await fetchTestimonials();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete testimonial');
		}
	}

	function startTranslation(testimonialId: string) {
		translatingTestimonialId = testimonialId;
		selectedLanguage = 'pt-BR';
		selectedLanguages = [];
		translationMode = 'single';
		showTranslationModal = true;
	}

	function cancelTranslation() {
		showTranslationModal = false;
		translatingTestimonialId = null;
	}

	async function handleRequestTranslation() {
		if (!translatingTestimonialId) return;

		try {
			if (translationMode === 'single') {
				await requestTranslation({
					entityType: 'testimonial',
					entityId: translatingTestimonialId,
					targetLanguages: [selectedLanguage]
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				// Translate to each language separately
				const results = await Promise.all(
					languages.map((lang) =>
						translateEntity({
							entityType: 'testimonial',
							entityId: translatingTestimonialId!,
							language: lang,
							fields: { content: '' } // Will be filled by the API
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

	function renderStars(rating?: number): string {
		if (!rating) return '—';
		return '★'.repeat(rating) + '☆'.repeat(5 - rating);
	}
</script>

<div class="page-container">
	<div class="page-header">
		<div>
			<h1 class="page-title">Testimonials</h1>
			<p class="page-description">Manage client and colleague testimonials</p>
		</div>
		<button type="button" class="btn btn-primary" onclick={startCreate}>
			<Plus class="icon" />
			Create Testimonial
		</button>
	</div>

	{#if error}
		<div class="alert alert-error">
			<AlertCircle class="icon" />
			<p>{error}</p>
		</div>
	{/if}

	{#if loading}
		<div class="loading-container">
			<div class="spinner"></div>
		</div>
	{:else if testimonials.length === 0}
		<div class="empty-state">
			<Users class="empty-icon" />
			<p class="empty-title">No testimonials found</p>
			<p class="empty-description">Create your first testimonial to get started</p>
		</div>
	{:else}
		<div class="table-container">
			<table class="table">
				<thead>
					<tr>
						<th>Author</th>
						<th>Content</th>
						<th>Rating</th>
						<th>Featured</th>
						<th>Created</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each testimonials as testimonial}
						<tr>
							<td>
								<div class="author-info">
									<a
										href="/landing/testimonials/{testimonial.id}"
										class="font-medium hover:text-pink-400 transition-colors"
									>
										{testimonial.author_name}
									</a>
									{#if testimonial.author_title || testimonial.author_company}
										<span class="text-muted text-sm">
											{testimonial.author_title}
											{#if testimonial.author_title && testimonial.author_company} at {/if}
											{testimonial.author_company}
										</span>
									{/if}
								</div>
							</td>
							<td>
								<div class="content-preview">{testimonial.content.substring(0, 80)}...</div>
							</td>
							<td>
								{#if testimonial.rating}
									<div class="rating">
										<Star class="icon-xs" />
										<span>{testimonial.rating}/5</span>
									</div>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td>
								{#if testimonial.featured}
									<span class="status-badge status-active">Featured</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td class="text-muted">{formatDate(testimonial.created_at)}</td>
							<td class="text-right">
								<div class="actions">
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => startTranslation(testimonial.id)}
										title="Request Translation"
									>
										<Globe class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-primary"
										onclick={() => startEdit(testimonial)}
										title="Edit"
									>
										<Edit class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-danger"
										onclick={() => handleDelete(testimonial.id)}
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
			<h2 class="modal-title">{editingTestimonial ? 'Edit Testimonial' : 'Create Testimonial'}</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Content *</label>
					<textarea
						bind:value={formContent}
						class="input textarea"
						rows="6"
						placeholder="Testimonial content"
					></textarea>
				</div>
				<div class="form-group">
					<label class="form-label">Author Name *</label>
					<input type="text" bind:value={formAuthorName} class="input" placeholder="Author name" />
				</div>
				<div class="form-row">
					<div class="form-group">
						<label class="form-label">Author Title</label>
						<input type="text" bind:value={formAuthorTitle} class="input" placeholder="Job title" />
					</div>
					<div class="form-group">
						<label class="form-label">Author Company</label>
						<input type="text" bind:value={formAuthorCompany} class="input" placeholder="Company name" />
					</div>
				</div>
				<div class="form-group">
					<label class="form-label">Author Image URL</label>
					<input type="url" bind:value={formAuthorImage} class="input" placeholder="https://..." />
				</div>
				<div class="form-group">
					<label class="form-label">Rating (1-5)</label>
					<select
						bind:value={formRating}
						class="input"
						onchange={(e) => {
							const val = parseInt(e.currentTarget.value);
							formRating = isNaN(val) ? null : val;
						}}
					>
						<option value="">No rating</option>
						<option value="1">1 Star</option>
						<option value="2">2 Stars</option>
						<option value="3">3 Stars</option>
						<option value="4">4 Stars</option>
						<option value="5">5 Stars</option>
					</select>
				</div>
				<div class="form-group">
					<label class="checkbox-label">
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
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
{#if showTranslationModal && translatingTestimonialId}
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

	.icon-xs {
		width: 0.75rem;
		height: 0.75rem;
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

	.text-sm {
		font-size: 0.75rem;
	}

	.author-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.font-medium {
		font-weight: 500;
	}

	.content-preview {
		max-width: 300px;
		color: rgba(203, 213, 225, 0.9);
		font-size: 0.875rem;
	}

	.rating {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		color: #fde047;
		font-size: 0.875rem;
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
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
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
