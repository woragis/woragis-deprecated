<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listTranslations,
		requestTranslation,
		translateEntity,
		type Translation,
		type RequestTranslationInput,
		type TranslateEntityInput,
		type EntityType,
		type Language
	} from '$lib/api/translations';

	let translations: Translation[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showRequestModal = $state(false);
	let showTranslateModal = $state(false);
	let searchQuery = $state('');

	// Request form state
	let requestEntityType = $state<EntityType>('post');
	let requestEntityId = $state('');
	let requestTargetLanguages = $state<Language[]>([]);

	// Translate form state
	let translateEntityType = $state<EntityType>('post');
	let translateEntityId = $state('');
	let translateLanguage = $state<Language>('pt-BR');
	let translateFields = $state('');
	const translateFieldsPlaceholder = '{"title": "Título", "description": "Descrição"}';

	const entityTypes: EntityType[] = [
		'testimonial',
		'post',
		'project',
		'case_study',
		'project_case_study',
		'system_design',
		'problem_solution',
		'certification',
		'aiml_integration',
		'impact_metric',
		'social_media_post',
		'technical_writing',
		'interest',
		'skill'
	];

	const languages: Language[] = ['en', 'pt-BR', 'fr', 'es', 'de', 'ru', 'ja', 'ko', 'zh-CN', 'el', 'la'];

	onMount(async () => {
		await fetchTranslations();
	});

	async function fetchTranslations() {
		loading = true;
		error = null;
		try {
			translations = await listTranslations();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load translations';
			console.error('Error fetching translations:', err);
		} finally {
			loading = false;
		}
	}

	function openRequestModal() {
		requestEntityType = 'post';
		requestEntityId = '';
		requestTargetLanguages = [];
		showRequestModal = true;
	}

	function openTranslateModal() {
		translateEntityType = 'post';
		translateEntityId = '';
		translateLanguage = 'pt-BR';
		translateFields = '';
		showTranslateModal = true;
	}

	function toggleLanguage(lang: Language) {
		if (requestTargetLanguages.includes(lang)) {
			requestTargetLanguages = requestTargetLanguages.filter((l) => l !== lang);
		} else {
			requestTargetLanguages = [...requestTargetLanguages, lang];
		}
	}

	async function handleRequest() {
		if (!requestEntityId.trim() || requestTargetLanguages.length === 0) {
			alert('Entity ID and at least one target language are required');
			return;
		}

		try {
			const input: RequestTranslationInput = {
				entityType: requestEntityType,
				entityId: requestEntityId.trim(),
				targetLanguages: requestTargetLanguages
			};

			await requestTranslation(input);
			showRequestModal = false;
			await fetchTranslations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to request translation');
			console.error('Error requesting translation:', err);
		}
	}

	async function handleTranslate() {
		if (!translateEntityId.trim() || !translateFields.trim()) {
			alert('Entity ID and fields are required');
			return;
		}

		try {
			let fieldsObj: Record<string, string>;
			try {
				fieldsObj = JSON.parse(translateFields);
			} catch {
				alert('Fields must be valid JSON');
				return;
			}

			const input: TranslateEntityInput = {
				entityType: translateEntityType,
				entityId: translateEntityId.trim(),
				language: translateLanguage,
				fields: fieldsObj
			};

			await translateEntity(input);
			showTranslateModal = false;
			await fetchTranslations();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to translate entity');
			console.error('Error translating entity:', err);
		}
	}

	function filteredTranslations() {
		if (!searchQuery.trim()) return translations;
		const query = searchQuery.toLowerCase();
		return translations.filter(
			(t) =>
				t.entityType.toLowerCase().includes(query) ||
				t.language.toLowerCase().includes(query) ||
				t.status.toLowerCase().includes(query)
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
			<h1>Translations Management</h1>
			<p>Manage translations for entities</p>
		</div>
		<div>
			<button onclick={openRequestModal}>Request Translation</button>
			<button onclick={openTranslateModal}>Translate Entity</button>
		</div>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search translations..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredTranslations().length === 0}
		<div class="empty">No translations found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Entity Type</th>
					<th>Entity ID</th>
					<th>Language</th>
					<th>Status</th>
					<th>Fields Count</th>
					<th>Created</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredTranslations() as translation}
					<tr>
						<td>{translation.entityType}</td>
						<td>
							<code>{translation.entityId.substring(0, 8)}...</code>
						</td>
						<td>{translation.language}</td>
						<td>
							<span class="status status-{translation.status}">{translation.status}</span>
						</td>
						<td>{Object.keys(translation.fields || {}).length}</td>
						<td>{formatDate(translation.createdAt)}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- Request Translation Modal -->
{#if showRequestModal}
	<div class="modal-overlay" onclick={() => (showRequestModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Request Translation</h2>
			<div class="form">
				<div class="form-group">
					<label>Entity Type *</label>
					<select bind:value={requestEntityType}>
						{#each entityTypes as et}
							<option value={et}>{et}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Entity ID *</label>
					<input type="text" bind:value={requestEntityId} />
				</div>
				<div class="form-group">
					<label>Target Languages *</label>
					<div class="language-checkboxes">
						{#each languages as lang}
							<label>
								<input
									type="checkbox"
									checked={requestTargetLanguages.includes(lang)}
									onchange={() => toggleLanguage(lang)}
								/>
								{lang}
							</label>
						{/each}
					</div>
				</div>
				<div class="form-actions">
					<button onclick={handleRequest}>Request</button>
					<button onclick={() => (showRequestModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Translate Entity Modal -->
{#if showTranslateModal}
	<div class="modal-overlay" onclick={() => (showTranslateModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Translate Entity</h2>
			<div class="form">
				<div class="form-group">
					<label>Entity Type *</label>
					<select bind:value={translateEntityType}>
						{#each entityTypes as et}
							<option value={et}>{et}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Entity ID *</label>
					<input type="text" bind:value={translateEntityId} />
				</div>
				<div class="form-group">
					<label>Language *</label>
					<select bind:value={translateLanguage}>
						{#each languages as lang}
							<option value={lang}>{lang}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Fields (JSON) *</label>
					<textarea
						bind:value={translateFields}
						rows="8"
						placeholder={translateFieldsPlaceholder}
					></textarea>
				</div>
				<div class="form-actions">
					<button onclick={handleTranslate}>Translate</button>
					<button onclick={() => (showTranslateModal = false)}>Cancel</button>
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

	.header > div {
		display: flex;
		gap: 0.5rem;
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

	code {
		background: #f5f5f5;
		padding: 0.125rem 0.375rem;
		border-radius: 3px;
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

	.status-processing {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-completed {
		background: #d4edda;
		color: #155724;
	}

	.status-failed {
		background: #f8d7da;
		color: #721c24;
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

	.modal-large {
		max-width: 800px;
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

	.language-checkboxes {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 0.5rem;
	}

	.language-checkboxes label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: normal;
		cursor: pointer;
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

