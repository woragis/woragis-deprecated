<script lang="ts">
	import { onMount } from 'svelte';
	import { Plus, Code, AlertCircle, Edit, Trash2, Globe, CheckSquare, Square } from 'lucide-svelte';
	import {
		listSkills,
		createSkill,
		updateSkill,
		deleteSkill,
		type Skill,
		type CreateSkillInput
	} from '$lib/api/landing';
	import { requestTranslation, translateEntity, SUPPORTED_LANGUAGES, type Language } from '$lib/api/translations';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let skills: Skill[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let showTranslationModal = $state(false);
	let showBulkActions = $state(false);
	let editingSkill: Skill | null = $state(null);
	let translatingSkillId: string | null = $state(null);
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formName = $state('');
	let formDescription = $state('');
	let formCategory = $state('other');
	let formProficiencyLevel = $state('');
	let formYearsOfExperience = $state<number | null>(null);
	let formFeatured = $state(false);

	// Translation state
	let selectedLanguage: Language = $state('pt-BR');
	let selectedLanguages: Language[] = $state([]);
	let translationMode: 'single' | 'multiple' = $state('single');

	onMount(async () => {
		await fetchSkills();
	});

	async function fetchSkills() {
		loading = true;
		error = null;
		try {
			skills = await listSkills();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch skills';
			toastError(error);
		} finally {
			loading = false;
		}
	}

	function startCreate() {
		editingSkill = null;
		formName = '';
		formDescription = '';
		formCategory = 'other';
		formProficiencyLevel = '';
		formYearsOfExperience = null;
		formFeatured = false;
		showCreateModal = true;
	}

	function startEdit(skill: Skill) {
		editingSkill = skill;
		formName = skill.name;
		formDescription = skill.description || '';
		formCategory = skill.category || 'other';
		formProficiencyLevel = skill.proficiency_level || '';
		formYearsOfExperience = skill.years_of_experience || null;
		formFeatured = skill.featured;
		showEditModal = true;
	}

	function cancelEdit() {
		showCreateModal = false;
		showEditModal = false;
		editingSkill = null;
	}

	async function handleSave() {
		if (!formName.trim()) {
			toastError('Name is required');
			return;
		}

		try {
			const payload: CreateSkillInput = {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				category: formCategory,
				proficiency_level: formProficiencyLevel || undefined,
				years_of_experience: formYearsOfExperience || undefined,
				featured: formFeatured
			};

			if (editingSkill) {
				await updateSkill(editingSkill.id, payload);
				toastSuccess('Skill updated successfully');
			} else {
				await createSkill(payload);
				toastSuccess('Skill created successfully');
			}

			cancelEdit();
			await fetchSkills();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to save skill');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this skill?')) return;
		try {
			await deleteSkill(id);
			toastSuccess('Skill deleted successfully');
			await fetchSkills();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete skill');
		}
	}

	async function handleBulkDelete() {
		if (selectedIds.size === 0) return;
		if (!confirm(`Are you sure you want to delete ${selectedIds.size} skill(s)?`)) return;

		try {
			await Promise.all(Array.from(selectedIds).map((id) => deleteSkill(id)));
			toastSuccess(`Deleted ${selectedIds.size} skill(s) successfully`);
			selectedIds.clear();
			showBulkActions = false;
			await fetchSkills();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to delete skills');
		}
	}

	async function handleBulkUpdate(featured: boolean) {
		if (selectedIds.size === 0) return;

		try {
			await Promise.all(Array.from(selectedIds).map((id) => updateSkill(id, { featured })));
			toastSuccess(`Updated ${selectedIds.size} skill(s) successfully`);
			selectedIds.clear();
			showBulkActions = false;
			await fetchSkills();
		} catch (err) {
			toastError(err instanceof Error ? err.message : 'Failed to update skills');
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
		if (selectedIds.size === skills.length) {
			selectedIds.clear();
		} else {
			selectedIds = new Set(skills.map((s) => s.id));
		}
		showBulkActions = selectedIds.size > 0;
	}

	function startTranslation(skillId: string) {
		translatingSkillId = skillId;
		selectedLanguage = 'pt-BR';
		selectedLanguages = [];
		translationMode = 'single';
		showTranslationModal = true;
	}

	function cancelTranslation() {
		showTranslationModal = false;
		translatingSkillId = null;
	}

	async function handleRequestTranslation() {
		if (!translatingSkillId) return;

		try {
			if (translationMode === 'single') {
				await requestTranslation({
					entityType: 'skill',
					entityId: translatingSkillId,
					targetLanguages: [selectedLanguage]
				});
				toastSuccess(`Translation to ${SUPPORTED_LANGUAGES.find((l) => l.value === selectedLanguage)?.label} queued`);
			} else {
				const languages = selectedLanguages.length > 0 ? selectedLanguages : SUPPORTED_LANGUAGES.map((l) => l.value);
				// Translate to each language separately
				const results = await Promise.all(
					languages.map((lang) =>
						translateEntity({
							entityType: 'skill',
							entityId: translatingSkillId!,
							language: lang,
							fields: { name: '', description: '' } // Will be filled by the API
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
			<h1 class="page-title">Skills</h1>
			<p class="page-description">Manage technical skills and proficiencies</p>
		</div>
		<button type="button" class="btn btn-primary" onclick={startCreate}>
			<Plus class="icon" />
			Create Skill
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
	{:else if skills.length === 0}
		<div class="empty-state">
			<Code class="empty-icon" />
			<p class="empty-title">No skills found</p>
			<p class="empty-description">Create your first skill to get started</p>
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
								{#if selectedIds.size === skills.length}
									<CheckSquare class="icon-sm" />
								{:else}
									<Square class="icon-sm" />
								{/if}
							</button>
						</th>
						<th>Name</th>
						<th>Category</th>
						<th>Proficiency</th>
						<th>Experience</th>
						<th>Featured</th>
						<th>Created</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each skills as skill}
						<tr>
							<td class="checkbox-column">
								<button
									type="button"
									class="checkbox-btn"
									onclick={() => toggleSelect(skill.id)}
								>
									{#if selectedIds.has(skill.id)}
										<CheckSquare class="icon-sm" />
									{:else}
										<Square class="icon-sm" />
									{/if}
								</button>
							</td>
							<td>
								<a
									href="/landing/skills/{skill.id}"
									class="font-medium hover:text-indigo-400 transition-colors"
								>
									{skill.name}
								</a>
							</td>
							<td>
								<span class="badge">{skill.category || 'other'}</span>
							</td>
							<td>
								{#if skill.proficiency_level}
									<span class="badge badge-proficiency">{skill.proficiency_level}</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td>
								{#if skill.years_of_experience}
									<span class="text-muted">{skill.years_of_experience} years</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td>
								{#if skill.featured}
									<span class="status-badge status-active">Featured</span>
								{:else}
									<span class="text-muted">—</span>
								{/if}
							</td>
							<td class="text-muted">{formatDate(skill.created_at)}</td>
							<td class="text-right">
								<div class="actions">
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => startTranslation(skill.id)}
										title="Request Translation"
									>
										<Globe class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-primary"
										onclick={() => startEdit(skill)}
										title="Edit"
									>
										<Edit class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-danger"
										onclick={() => handleDelete(skill.id)}
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
			<h2 class="modal-title">{editingSkill ? 'Edit Skill' : 'Create Skill'}</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Name *</label>
					<input type="text" bind:value={formName} class="input" placeholder="Skill name" />
				</div>
				<div class="form-group">
					<label class="form-label">Description</label>
					<textarea
						bind:value={formDescription}
						class="input textarea"
						rows="4"
						placeholder="Skill description"
					></textarea>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label class="form-label">Category</label>
						<select bind:value={formCategory} class="input">
							<option value="backend">Backend</option>
							<option value="frontend">Frontend</option>
							<option value="database">Database</option>
							<option value="infrastructure">Infrastructure</option>
							<option value="devops">DevOps</option>
							<option value="language">Language</option>
							<option value="framework">Framework</option>
							<option value="tool">Tool</option>
							<option value="service">Service</option>
							<option value="library">Library</option>
							<option value="other">Other</option>
						</select>
					</div>
					<div class="form-group">
						<label class="form-label">Proficiency Level</label>
						<select bind:value={formProficiencyLevel} class="input">
							<option value="">Not specified</option>
							<option value="expert">Expert</option>
							<option value="advanced">Advanced</option>
							<option value="proficient">Proficient</option>
							<option value="learning">Learning</option>
						</select>
					</div>
				</div>
				<div class="form-group">
					<label class="form-label">Years of Experience</label>
					<input
						type="number"
						bind:value={formYearsOfExperience}
						class="input"
						placeholder="0"
						min="0"
						step="0.5"
					/>
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
{#if showTranslationModal && translatingSkillId}
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

	.font-medium {
		font-weight: 500;
	}

	.badge {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 0.375rem;
		font-size: 0.75rem;
		color: #cbd5e1;
		text-transform: capitalize;
	}

	.badge-proficiency {
		background: rgba(255, 255, 255, 0.08);
		color: #d4d4d4;
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
