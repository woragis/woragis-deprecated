<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getSkill,
		updateSkill,
		deleteSkill,
		type Skill,
		type UpdateSkillInput,
		type SkillCategory,
		type ProficiencyLevel
	} from '$lib/api/skills';

	let skill: Skill | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showEditModal = $state(false);

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

	const skillId = $derived($page.params.id);

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
		if (skillId) {
			await loadSkill();
		}
	});

	async function loadSkill() {
		if (!skillId) return;
		loading = true;
		error = null;
		try {
			skill = await getSkill(skillId);
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
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load skill';
			console.error('Error loading skill:', err);
		} finally {
			loading = false;
		}
	}

	function openEditModal() {
		if (!skill) return;
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

	async function handleUpdate() {
		if (!skill || !formName.trim()) {
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

			await updateSkill(skill.id, input);
			showEditModal = false;
			await loadSkill();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update skill');
			console.error('Error updating skill:', err);
		}
	}

	async function handleDelete() {
		if (!skill || !confirm('Are you sure you want to delete this skill?')) return;

		try {
			await deleteSkill(skill.id);
			await goto('/skills');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete skill');
			console.error('Error deleting skill:', err);
		}
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="header">
		<a href="/skills" class="back-link">← Back to Skills</a>
		<div class="header-actions">
			{#if skill}
				<button onclick={openEditModal}>Edit Skill</button>
				<button onclick={handleDelete} class="delete-btn">Delete</button>
			{/if}
		</div>
	</div>

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if error}
		<div class="error">{error}</div>
	{:else if skill}
		<div class="details-container">
			<div class="details-header">
				<h2>{skill.name}</h2>
			</div>

			<div class="details-grid">
				<div class="detail-section">
					<h3>Basic Information</h3>
					<div class="detail-item">
						<strong>Category:</strong> {skill.category}
					</div>
					<div class="detail-item">
						<strong>Description:</strong> {skill.description || '—'}
					</div>
					<div class="detail-item">
						<strong>Proficiency Level:</strong> {skill.proficiencyLevel || '—'}
					</div>
					<div class="detail-item">
						<strong>Years of Experience:</strong> {skill.yearsOfExperience || '—'}
					</div>
					<div class="detail-item">
						<strong>First Used:</strong> {formatDate(skill.firstUsedDate)}
					</div>
					<div class="detail-item">
						<strong>Last Used:</strong> {formatDate(skill.lastUsedDate)}
					</div>
				</div>

				<div class="detail-section">
					<h3>Styling</h3>
					<div class="detail-item">
						<strong>Icon:</strong> {skill.icon || '—'}
					</div>
					<div class="detail-item">
						<strong>Color:</strong> {skill.color || '—'}
					</div>
					<div class="detail-item">
						<strong>Background Gradient:</strong> {skill.bgGradient || '—'}
					</div>
					<div class="detail-item">
						<strong>Border Color:</strong> {skill.borderColor || '—'}
					</div>
					<div class="detail-item">
						<strong>Hover Border Color:</strong> {skill.hoverBorderColor || '—'}
					</div>
					<div class="detail-item">
						<strong>Shadow Color:</strong> {skill.shadowColor || '—'}
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Edit Modal -->
{#if showEditModal && skill}
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
		max-width: 1400px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.back-link {
		color: #007bff;
		text-decoration: none;
		font-size: 0.9rem;
		padding: 0.5rem 0;
	}

	.back-link:hover {
		color: #0056b3;
		text-decoration: underline;
	}

	.header-actions {
		display: flex;
		gap: 0.5rem;
	}

	.header-actions button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.header-actions button:hover {
		background: #0056b3;
	}

	.delete-btn {
		background: #dc3545 !important;
	}

	.delete-btn:hover {
		background: #c82333 !important;
	}

	.loading,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #666;
	}

	.error {
		padding: 0.75rem;
		background: #fee;
		color: #c33;
		border: 1px solid #fcc;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.details-container {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		color: #333;
	}

	.details-header {
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 2px solid #e0e0e0;
	}

	.details-header h2 {
		margin: 0;
		font-size: 1.5rem;
		color: #333;
	}

	.details-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}

	.detail-section {
		background: #f9f9f9;
		padding: 1rem;
		border-radius: 4px;
	}

	.detail-section h3 {
		margin: 0 0 1rem 0;
		font-size: 1.1rem;
		color: #333;
	}

	.detail-item {
		margin-bottom: 0.5rem;
		font-size: 0.9rem;
		color: #333;
	}

	.detail-item strong {
		color: #333;
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

