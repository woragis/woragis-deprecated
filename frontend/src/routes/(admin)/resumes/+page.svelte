<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listResumes,
		listResumesByTags,
		listResumeTags,
		createResume,
		uploadResume,
		deleteResume,
		type Resume,
		type CreateResumeInput
	} from '$lib/api/resumes';
	import TagInput from '$lib/components/ui/TagInput.svelte';
	import { Eye, Trash2, Plus, Upload, Star } from 'lucide-svelte';
	import { toastError, toastSuccess } from '$lib/utils/toast';
	import { locale, t } from '$lib/i18n';
	// User ID for resume downloads - should match your backend user ID
	// You can move this to a constants file or get from environment
	const userId = import.meta.env.PUBLIC_DEFAULT_USER_ID || '6ad0d828-f605-45fc-a545-3441e17a015c';

	let resumes: Resume[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let createFormTitle = $state('');
	let createFormFile: File | null = $state(null);
	let createFormFilePath = $state('');
	let createFormFileName = $state('');
	let createFormFileSize = $state(0);
	let createFormTags = $state<string[]>([]);
	let uploading = $state(false);
	let filterTags = $state<string[]>([]);
	let availableTags = $state<string[]>([]);

	onMount(async () => {
		await Promise.all([fetchResumes(), fetchAvailableTags()]);
	});

	async function fetchResumes() {
		loading = true;
		error = null;
		try {
			if (filterTags.length > 0) {
				resumes = await listResumesByTags(filterTags);
			} else {
				resumes = await listResumes();
			}
		} catch (err) {
			error = err instanceof Error ? err.message : $t('resumes.error');
			console.error('Error fetching resumes:', err);
			toastError($t('resumes.loadError'));
		} finally {
			loading = false;
		}
	}

	async function fetchAvailableTags() {
		try {
			availableTags = await listResumeTags();
		} catch (err) {
			console.error('Error fetching tags:', err);
		}
	}

	function handleFilterTagsChange(newTags: string[]) {
		filterTags = newTags;
		fetchResumes();
	}


	function openCreateModal() {
		showCreateModal = true;
		createFormTitle = '';
		createFormFile = null;
		createFormFilePath = '';
		createFormFileName = '';
		createFormFileSize = 0;
		createFormTags = [];
		if (availableTags.length === 0) {
			fetchAvailableTags();
		}
	}

	function closeCreateModal() {
		showCreateModal = false;
		createFormTitle = '';
		createFormFile = null;
		createFormFilePath = '';
		createFormFileName = '';
		createFormFileSize = 0;
		createFormTags = [];
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			createFormFile = file;
			createFormFileName = file.name;
			createFormFileSize = file.size;
			// Generate a file path - in a real scenario, this would be set by the backend after upload
			// For now, we'll use a relative path that the backend can handle
			createFormFilePath = `uploads/${Date.now()}_${file.name}`;
		}
	}

	async function handleCreate() {
		if (!createFormTitle.trim()) {
			toastError($t('resumes.modal.title') + ' ' + $t('resumes.modal.required'));
			return;
		}

		uploading = true;
		try {
			// If a file is selected, use the upload endpoint
			if (createFormFile) {
				await uploadResume(createFormFile, createFormTitle.trim());
				toastSuccess($t('resumes.uploadSuccess'));
				closeCreateModal();
				await fetchResumes();
				return;
			}

			// Otherwise, use the manual create endpoint (for files already on server)
			if (!createFormFilePath) {
				toastError($t('resumes.createError'));
				return;
			}

			const input: CreateResumeInput = {
				title: createFormTitle.trim(),
				filePath: createFormFilePath,
				fileName: createFormFileName || 'resume.pdf',
				fileSize: createFormFileSize || 0,
				tags: createFormTags.length > 0 ? createFormTags : undefined
			};

			await createResume(input);
			toastSuccess($t('resumes.createSuccess'));
			closeCreateModal();
			await Promise.all([fetchResumes(), fetchAvailableTags()]);
		} catch (err) {
			console.error('Error creating resume:', err);
			toastError(err instanceof Error ? err.message : $t('resumes.createError'));
		} finally {
			uploading = false;
		}
	}


	async function handleDelete(resume: Resume) {
		if (!confirm(`${$t('resumes.deleteConfirm')} "${resume.title}"?`)) {
			return;
		}

		try {
			await deleteResume(resume.id);
			toastSuccess($t('resumes.deleteSuccess'));
			await fetchResumes();
		} catch (err) {
			console.error('Error deleting resume:', err);
			toastError($t('resumes.deleteError'));
		}
	}


	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

</script>

<div class="container mx-auto px-4 py-8">
	<div class="mb-6">
		<div class="flex items-center justify-between mb-4">
			<div>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white">{$t('resumes.title')}</h1>
				<p class="text-gray-600 dark:text-gray-400 mt-1">{$t('resumes.subtitle')}</p>
			</div>
			<button
				onclick={openCreateModal}
				class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
			>
				<Plus class="w-4 h-4" />
				{$t('resumes.createButton')}
			</button>
		</div>
		<div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
			<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
				Filter by Tags
			</label>
			<TagInput
				bind:tags={filterTags}
				{availableTags}
				onFetchTags={listResumeTags}
				placeholder="Filter by tags..."
				maxTags={10}
			/>
		</div>
	</div>

	{#if loading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600 dark:text-gray-400">{$t('resumes.loading')}</p>
		</div>
	{:else if error}
		<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
			<p class="text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if resumes.length === 0}
		<div class="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg">
			<p class="text-gray-600 dark:text-gray-400">{$t('resumes.empty')}</p>
			<p class="text-sm text-gray-500 dark:text-gray-500 mt-2">
				{$t('resumes.emptySubtext')}
			</p>
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-900">
						<tr>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								{$t('resumes.table.title')}
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								{$t('resumes.table.file')}
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								{$t('resumes.table.status')}
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Tags
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Applications
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Interview Rate
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Offer Rate
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								{$t('resumes.table.created')}
							</th>
							<th
								class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								{$t('resumes.table.actions')}
							</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each resumes as resume}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div>
											<div class="text-sm font-medium text-gray-900 dark:text-white">
												{resume.title}
											</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-gray-900 dark:text-gray-300">{resume.fileName}</div>
									<div class="text-xs text-gray-500 dark:text-gray-400">
										{formatFileSize(resume.fileSize)}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center gap-2">
										{#if resume.isMain}
											<span
												class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
											>
												{$t('resumes.table.main')}
											</span>
										{/if}
										{#if resume.isFeatured}
											<span
												class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
											>
												<Star class="w-3 h-3 mr-1" />
												{$t('resumes.table.featured')}
											</span>
										{/if}
									</div>
								</td>
								<td class="px-6 py-4">
									<div class="flex flex-wrap gap-1">
										{#if resume.tags && resume.tags.length > 0}
											{#each resume.tags as tag}
												<span
													class="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
												>
													{tag}
												</span>
											{/each}
										{:else}
											<span class="text-xs text-gray-400 dark:text-gray-500">—</span>
										{/if}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
									{resume.applicationsUsed || 0}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div
										class="text-sm font-medium text-gray-900 dark:text-gray-300 cursor-help"
										title="Interviews: {Math.round((resume.interviewRate || 0) / 100 * (resume.applicationsUsed || 0))} / {resume.applicationsUsed || 0}"
									>
										{(resume.interviewRate || 0).toFixed(1)}%
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div
										class="text-sm font-medium text-gray-900 dark:text-gray-300 cursor-help"
										title="Offers: {Math.round((resume.offerRate || 0) / 100 * (resume.applicationsUsed || 0))} / {resume.applicationsUsed || 0}"
									>
										{(resume.offerRate || 0).toFixed(1)}%
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
									{formatDate(resume.createdAt)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex items-center justify-end gap-2">
										<a
											href="/resumes/{resume.id}"
											class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
											title="View Details"
										>
											<Eye class="w-4 h-4" />
										</a>
										<button
											onclick={() => handleDelete(resume)}
											class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
											title="Delete"
										>
											<Trash2 class="w-4 h-4" />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>


<!-- Create Modal -->
{#if showCreateModal}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
		onclick={closeCreateModal}
	>
		<div
			class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">{$t('resumes.modal.createTitle')}</h2>
			<div class="space-y-4">
				<div>
					<label
						for="create-title"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						{$t('resumes.modal.title')} <span class="text-red-500">{$t('resumes.modal.required')}</span>
					</label>
					<input
						id="create-title"
						type="text"
						bind:value={createFormTitle}
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						placeholder="Resume title"
					/>
				</div>
				<div>
					<label
						for="create-file"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						{$t('resumes.modal.file')} <span class="text-red-500">{$t('resumes.modal.required')}</span>
					</label>
					<div class="mt-1 flex items-center gap-4">
						<label
							for="create-file-input"
							class="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-200 dark:hover:bg-gray-600 cursor-pointer transition-colors"
						>
							<Upload class="w-4 h-4" />
							<span>{$t('resumes.modal.chooseFile')}</span>
						</label>
						<input
							id="create-file-input"
							type="file"
							accept=".pdf,application/pdf"
							onchange={handleFileSelect}
							class="hidden"
						/>
						{#if createFormFile}
							<span class="text-sm text-gray-600 dark:text-gray-400">
								{createFormFile.name} ({formatFileSize(createFormFile.size)})
							</span>
						{/if}
					</div>
					<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
						{$t('resumes.modal.manualEntry')}
					</p>
				</div>
				<div>
					<label
						for="create-file-path"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						{$t('resumes.modal.filePath')}
					</label>
					<input
						id="create-file-path"
						type="text"
						bind:value={createFormFilePath}
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						placeholder="e.g., resumes/resume_2024.pdf"
					/>
				</div>
				<div>
					<label
						for="create-file-name"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						{$t('resumes.modal.fileName')}
					</label>
					<input
						id="create-file-name"
						type="text"
						bind:value={createFormFileName}
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						placeholder="e.g., resume_2024.pdf"
					/>
				</div>
				<div>
					<label
						for="create-file-size"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						{$t('resumes.modal.fileSize')}
					</label>
					<input
						id="create-file-size"
						type="number"
						bind:value={createFormFileSize}
						min="0"
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						placeholder="0"
					/>
				</div>
				<div>
					<label
						for="create-tags"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						Tags
					</label>
					<TagInput
						bind:tags={createFormTags}
						{availableTags}
						onFetchTags={listResumeTags}
						placeholder="Add tags (e.g., golang, python, aws)..."
						maxTags={10}
					/>
					<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
						Add technology tags like programming languages, frameworks, tools, etc.
					</p>
				</div>
			</div>
			<div class="mt-6 flex justify-end gap-3">
				<button
					onclick={closeCreateModal}
					disabled={uploading}
					class="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50"
				>
					{$t('resumes.modal.cancel')}
				</button>
				<button
					onclick={handleCreate}
					disabled={uploading}
					class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 flex items-center gap-2"
				>
					{#if uploading}
						<div class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
						{$t('resumes.modal.creating')}
					{:else}
						{$t('resumes.modal.create')}
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

