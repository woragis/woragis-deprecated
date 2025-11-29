<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, BookOpen, AlertCircle, Edit, Trash2, Globe, X } from 'lucide-svelte';
	import {
		listTechnicalWritings,
		createTechnicalWriting,
		updateTechnicalWriting,
		deleteTechnicalWriting,
		type TechnicalWriting,
		type CreateTechnicalWritingInput
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let writings: TechnicalWriting[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let editingWriting: TechnicalWriting | null = $state(null);
	let translatingWritingId: string | null = $state(null);

	// Form state
	let formTitle = $state('');
	let formContent = $state('');
	let formExcerpt = $state('');
	let formType = $state('article');
	let formPlatform = $state('medium');
	let formUrl = $state('');
	let formFeatured = $state(false);

	// Translation state
	let selectedLanguage: Language = $state('pt-BR');
	let selectedLanguages: Language[] = $state([]);
	let translationMode: 'single' | 'multiple' = $state('single');

	onMount(async () => {
		await fetchWritings();
	});

	async function fetchWritings() {
		loading = true;
		error = null;
		try {
			writings = await listTechnicalWritings();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch technical writings';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		editingWriting = null;
		formTitle = '';
		formContent = '';
		formExcerpt = '';
		formType = 'article';
		formPlatform = 'medium';
		formUrl = '';
		formFeatured = false;
		showCreateModal = true;
	}

	function startEdit(writing: TechnicalWriting) {
		editingWriting = writing;
		formTitle = writing.title;
		formContent = writing.content;
		formExcerpt = writing.excerpt || '';
		formType = writing.type;
		formPlatform = writing.platform;
		formUrl = writing.url || '';
		formFeatured = writing.featured;
		showEditModal = true;
	}

	function cancelEdit() {
		showCreateModal = false;
		showEditModal = false;
		editingWriting = null;
	}

	async function handleSave() {
		if (!formTitle.trim() || !formContent.trim()) {
			toastError('Title and content are required');
			return;
		}

		try {
			const payload: CreateTechnicalWritingInput = {
				title: formTitle.trim(),
				content: formContent.trim(),
				excerpt: formExcerpt.trim() || undefined,
				type: formType,
				platform: formPlatform,
				url: formUrl.trim() || undefined,
				featured: formFeatured
			};

			if (editingWriting) {
				await updateTechnicalWriting(editingWriting.id, payload);
				toastSuccess('Technical writing updated successfully');
			} else {
				await createTechnicalWriting(payload);
				toastSuccess('Technical writing created successfully');
			}

			cancelEdit();
			await fetchWritings();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to save technical writing');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this technical writing?')) return;
		try {
			await deleteTechnicalWriting(id);
			toastSuccess('Technical writing deleted successfully');
			await fetchWritings();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete technical writing');
		}
	}

	function startTranslation(writingId: string) {
		translatingWritingId = writingId;
		selectedLanguage = 'pt-BR';
		selectedLanguages = [];
		translationMode = 'single';
		showTranslationModal = true;
	}

	function cancelTranslation() {
		showTranslationModal = false;
		translatingWritingId = null;
	}

	async function handleRequestTranslation() {
		if (!translatingWritingId) return;

		try {
			if (translationMode === 'single') {
				await requestTranslation({
					entityType: 'technical_writing',
					entityId: translatingWritingId,
					language: selectedLanguage,
					fields: ['title', 'content', 'excerpt']
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				const result = await translateEntity({
					entityType: 'technical_writing',
					entityId: translatingWritingId,
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
			<h1 class="page-title">Technical Writings</h1>
			<p class="page-description">Manage technical articles and tutorials</p>
		</div>
		<button type="button" class="btn btn-primary" onclick={startCreate}>
			<Plus class="icon" />
			Create Writing
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
	{:else if writings.length === 0}
		<div class="empty-state">
			<BookOpen class="empty-icon" />
			<p class="empty-title">No technical writings found</p>
			<p class="empty-description">Create your first technical writing to get started</p>
		</div>
	{:else}
		<div class="table-container">
			<table class="table">
				<thead>
					<tr>
						<th>Title</th>
						<th>Type</th>
						<th>Platform</th>
						<th>Featured</th>
						<th>Created</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each writings as writing}
						<tr>
							<td>
								<div class="item-title">
									<span class="font-medium">{writing.title}</span>
									{#if writing.url}
										<a href={writing.url} target="_blank" class="text-muted text-sm">View →</a>
									{/if}
								</div>
							</td>
							<td><span class="badge">{writing.type}</span></td>
							<td><span class="badge">{writing.platform}</span></td>
							<td>
								{#if writing.featured}
									<span class="status-badge status-active">Featured</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td class="text-muted">{formatDate(writing.created_at)}</td>
							<td class="text-right">
								<div class="actions">
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => startTranslation(writing.id)}
										title="Request Translation"
									>
										<Globe class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-primary"
										onclick={() => startEdit(writing)}
										title="Edit"
									>
										<Edit class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-danger"
										onclick={() => handleDelete(writing.id)}
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
			<h2 class="modal-title">{editingWriting ? 'Edit Technical Writing' : 'Create Technical Writing'}</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Title *</label>
					<input type="text" bind:value={formTitle} class="input" placeholder="Article title" />
				</div>
				<div class="form-group">
					<label class="form-label">Content *</label>
					<textarea
						bind:value={formContent}
						class="input textarea"
						rows="10"
						placeholder="Article content"
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
						<label class="form-label">Type</label>
						<select bind:value={formType} class="input">
							<option value="article">Article</option>
							<option value="tutorial">Tutorial</option>
							<option value="guide">Guide</option>
							<option value="blog">Blog Post</option>
						</select>
					</div>
					<div class="form-group">
						<label class="form-label">Platform</label>
						<select bind:value={formPlatform} class="input">
							<option value="medium">Medium</option>
							<option value="dev.to">Dev.to</option>
							<option value="hashnode">Hashnode</option>
							<option value="personal">Personal Blog</option>
							<option value="other">Other</option>
						</select>
					</div>
				</div>
				<div class="form-group">
					<label class="form-label">URL</label>
					<input type="url" bind:value={formUrl} class="input" placeholder="https://..." />
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
{#if showTranslationModal && translatingWritingId}
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

	.text-sm {
		font-size: 0.75rem;
	}

	.item-title {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.font-medium {
		font-weight: 500;
	}

	.badge {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		background: rgba(71, 85, 105, 0.2);
		border-radius: 0.375rem;
		font-size: 0.75rem;
		color: #cbd5e1;
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
