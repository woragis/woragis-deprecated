<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listSkills,
		createSkill,
		updateSkill,
		deleteSkill,
		getSkill,
		type Skill,
		type CreateSkillInput,
		type UpdateSkillInput,
		type SkillCategory,
		type ProficiencyLevel
	} from '$lib/api/skills';

	let skills: Skill[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingSkill: Skill | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formName = $state('');
	let formDescription = $state('');
	let formCategory = $state<SkillCategory>('other');
	let formIcon = $state('');
	let formColor = $state('');
	let formBgGradient = $state('');
	let formBorderColor = $state('');
	let formHoverBorderColor = $state('');
	let formShadowColor = $state('');
	let formProficiencyLevel = $state<ProficiencyLevel | ''>('');
	let formYearsOfExperience = $state<number | ''>('');
	let formFirstUsedDate = $state('');
	let formLastUsedDate = $state('');

	const categories: SkillCategory[] = [
		'backend',
		'frontend',
		'database',
		'infrastructure',
		'devops',
		'language',
		'framework',
		'tool',
		'service',
		'library',
		'other'
	];

	const proficiencyLevels: ProficiencyLevel[] = ['expert', 'advanced', 'proficient', 'learning'];

	onMount(async () => {
		await fetchSkills();
	});

	async function fetchSkills() {
		loading = true;
		error = null;
		try {
			skills = await listSkills();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load skills';
			console.error('Error fetching skills:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(skill: Skill) {
		editingSkill = skill;
		formName = skill.name;
		formDescription = skill.description || '';
		formCategory = skill.category;
		formIcon = skill.icon || '';
		formColor = skill.color || '';
		formBgGradient = skill.bgGradient || '';
		formBorderColor = skill.borderColor || '';
		formHoverBorderColor = skill.hoverBorderColor || '';
		formShadowColor = skill.shadowColor || '';
		formProficiencyLevel = skill.proficiencyLevel || '';
		formYearsOfExperience = skill.yearsOfExperience || '';
		formFirstUsedDate = skill.firstUsedDate || '';
		formLastUsedDate = skill.lastUsedDate || '';
		showEditModal = true;
	}

	function resetForm() {
		formName = '';
		formDescription = '';
		formCategory = 'other';
		formIcon = '';
		formColor = '';
		formBgGradient = '';
		formBorderColor = '';
		formHoverBorderColor = '';
		formShadowColor = '';
		formProficiencyLevel = '';
		formYearsOfExperience = '';
		formFirstUsedDate = '';
		formLastUsedDate = '';
		editingSkill = null;
	}

	async function handleCreate() {
		if (!formName.trim()) {
			alert('Name is required');
			return;
		}

		try {
			const input: CreateSkillInput = {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				category: formCategory,
				icon: formIcon.trim() || undefined,
				color: formColor.trim() || undefined,
				bgGradient: formBgGradient.trim() || undefined,
				borderColor: formBorderColor.trim() || undefined,
				hoverBorderColor: formHoverBorderColor.trim() || undefined,
				shadowColor: formShadowColor.trim() || undefined,
				proficiencyLevel: formProficiencyLevel || undefined,
				yearsOfExperience: formYearsOfExperience ? Number(formYearsOfExperience) : undefined,
				firstUsedDate: formFirstUsedDate || undefined,
				lastUsedDate: formLastUsedDate || undefined
			};

			await createSkill(input);
			showCreateModal = false;
			resetForm();
			await fetchSkills();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create skill');
			console.error('Error creating skill:', err);
		}
	}

	async function handleUpdate() {
		if (!editingSkill || !formName.trim()) {
			alert('Name is required');
			return;
		}

		try {
			const input: UpdateSkillInput = {
				name: formName.trim(),
				description: formDescription.trim() || undefined,
				category: formCategory,
				icon: formIcon.trim() || undefined,
				color: formColor.trim() || undefined,
				bgGradient: formBgGradient.trim() || undefined,
				borderColor: formBorderColor.trim() || undefined,
				hoverBorderColor: formHoverBorderColor.trim() || undefined,
				shadowColor: formShadowColor.trim() || undefined,
				proficiencyLevel: formProficiencyLevel || undefined,
				yearsOfExperience: formYearsOfExperience ? Number(formYearsOfExperience) : undefined,
				firstUsedDate: formFirstUsedDate || undefined,
				lastUsedDate: formLastUsedDate || undefined
			};

			await updateSkill(editingSkill.id, input);
			showEditModal = false;
			resetForm();
			await fetchSkills();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update skill');
			console.error('Error updating skill:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this skill?')) return;

		try {
			await deleteSkill(id);
			await fetchSkills();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete skill');
			console.error('Error deleting skill:', err);
		}
	}

	function filteredSkills() {
		if (!searchQuery.trim()) return skills;
		const query = searchQuery.toLowerCase();
		return skills.filter(
			(s) =>
				s.name.toLowerCase().includes(query) ||
				s.description?.toLowerCase().includes(query) ||
				s.category.toLowerCase().includes(query)
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
			<h1>Skills Management</h1>
			<p>Manage skills and their details</p>
		</div>
		<button onclick={openCreateModal}>Create Skill</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search skills..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredSkills().length === 0}
		<div class="empty">No skills found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Category</th>
					<th>Proficiency</th>
					<th>Years</th>
					<th>First Used</th>
					<th>Last Used</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredSkills() as skill}
					<tr>
						<td>
							<strong>{skill.name}</strong>
							{#if skill.description}
								<br />
								<small>{skill.description}</small>
							{/if}
						</td>
						<td>{skill.category}</td>
						<td>{skill.proficiencyLevel || '—'}</td>
						<td>{skill.yearsOfExperience || '—'}</td>
						<td>{formatDate(skill.firstUsedDate)}</td>
						<td>{formatDate(skill.lastUsedDate)}</td>
						<td>
							<button onclick={() => openEditModal(skill)}>Edit</button>
							<button onclick={() => handleDelete(skill.id)} class="delete-btn">Delete</button>
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
			<h2>Create Skill</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Category *</label>
					<select bind:value={formCategory}>
						{#each categories as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Proficiency Level</label>
					<select bind:value={formProficiencyLevel}>
						<option value="">—</option>
						{#each proficiencyLevels as level}
							<option value={level}>{level}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Years of Experience</label>
					<input type="number" bind:value={formYearsOfExperience} />
				</div>
				<div class="form-group">
					<label>First Used Date</label>
					<input type="date" bind:value={formFirstUsedDate} />
				</div>
				<div class="form-group">
					<label>Last Used Date</label>
					<input type="date" bind:value={formLastUsedDate} />
				</div>
				<div class="form-group">
					<label>Icon</label>
					<input type="text" bind:value={formIcon} placeholder="Icon identifier" />
				</div>
				<div class="form-group">
					<label>Color</label>
					<input type="text" bind:value={formColor} placeholder="Color name" />
				</div>
				<div class="form-group">
					<label>Background Gradient</label>
					<input type="text" bind:value={formBgGradient} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Border Color</label>
					<input type="text" bind:value={formBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Hover Border Color</label>
					<input type="text" bind:value={formHoverBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Shadow Color</label>
					<input type="text" bind:value={formShadowColor} placeholder="Tailwind classes" />
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
{#if showEditModal && editingSkill}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Skill</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>Category *</label>
					<select bind:value={formCategory}>
						{#each categories as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Proficiency Level</label>
					<select bind:value={formProficiencyLevel}>
						<option value="">—</option>
						{#each proficiencyLevels as level}
							<option value={level}>{level}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Years of Experience</label>
					<input type="number" bind:value={formYearsOfExperience} />
				</div>
				<div class="form-group">
					<label>First Used Date</label>
					<input type="date" bind:value={formFirstUsedDate} />
				</div>
				<div class="form-group">
					<label>Last Used Date</label>
					<input type="date" bind:value={formLastUsedDate} />
				</div>
				<div class="form-group">
					<label>Icon</label>
					<input type="text" bind:value={formIcon} placeholder="Icon identifier" />
				</div>
				<div class="form-group">
					<label>Color</label>
					<input type="text" bind:value={formColor} placeholder="Color name" />
				</div>
				<div class="form-group">
					<label>Background Gradient</label>
					<input type="text" bind:value={formBgGradient} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Border Color</label>
					<input type="text" bind:value={formBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Hover Border Color</label>
					<input type="text" bind:value={formHoverBorderColor} placeholder="Tailwind classes" />
				</div>
				<div class="form-group">
					<label>Shadow Color</label>
					<input type="text" bind:value={formShadowColor} placeholder="Tailwind classes" />
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

