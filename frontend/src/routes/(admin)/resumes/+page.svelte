<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listResumes,
		createResume,
		uploadResume,
		updateResume,
		deleteResume,
		markAsMain,
		markAsFeatured,
		unmarkAsMain,
		unmarkAsFeatured,
		getResumeDownloadUrl,
		getResumePreviewUrl,
		type Resume,
		type CreateResumeInput,
		type UpdateResumeInput
	} from '$lib/api/resumes';
	import { Download, Eye, Star, Trash2, Edit2, Check, X, Plus, Upload } from 'lucide-svelte';
	import { toastError, toastSuccess } from '$lib/utils/toast';
	// User ID for resume downloads - should match your backend user ID
	// You can move this to a constants file or get from environment
	const userId = import.meta.env.PUBLIC_DEFAULT_USER_ID || '6ad0d828-f605-45fc-a545-3441e17a015c';

	let resumes: Resume[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showEditModal = $state(false);
	let showCreateModal = $state(false);
	let editingResume: Resume | null = $state(null);
	let formTitle = $state('');
	let createFormTitle = $state('');
	let createFormFile: File | null = $state(null);
	let createFormFilePath = $state('');
	let createFormFileName = $state('');
	let createFormFileSize = $state(0);
	let uploading = $state(false);

	onMount(async () => {
		await fetchResumes();
	});

	async function fetchResumes() {
		loading = true;
		error = null;
		try {
			resumes = await listResumes();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load resumes';
			console.error('Error fetching resumes:', err);
			toastError('Failed to load resumes');
		} finally {
			loading = false;
		}
	}

	function openEditModal(resume: Resume) {
		editingResume = resume;
		formTitle = resume.title;
		showEditModal = true;
	}

	function closeEditModal() {
		showEditModal = false;
		editingResume = null;
		formTitle = '';
	}

	function openCreateModal() {
		showCreateModal = true;
		createFormTitle = '';
		createFormFile = null;
		createFormFilePath = '';
		createFormFileName = '';
		createFormFileSize = 0;
	}

	function closeCreateModal() {
		showCreateModal = false;
		createFormTitle = '';
		createFormFile = null;
		createFormFilePath = '';
		createFormFileName = '';
		createFormFileSize = 0;
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
			toastError('Please enter a title');
			return;
		}

		uploading = true;
		try {
			// If a file is selected, use the upload endpoint
			if (createFormFile) {
				await uploadResume(createFormFile, createFormTitle.trim());
				toastSuccess('Resume uploaded successfully');
				closeCreateModal();
				await fetchResumes();
				return;
			}

			// Otherwise, use the manual create endpoint (for files already on server)
			if (!createFormFilePath) {
				toastError('Please select a file or provide a file path');
				return;
			}

			const input: CreateResumeInput = {
				title: createFormTitle.trim(),
				filePath: createFormFilePath,
				fileName: createFormFileName || 'resume.pdf',
				fileSize: createFormFileSize || 0
			};

			await createResume(input);
			toastSuccess('Resume created successfully');
			closeCreateModal();
			await fetchResumes();
		} catch (err) {
			console.error('Error creating resume:', err);
			toastError(err instanceof Error ? err.message : 'Failed to create resume');
		} finally {
			uploading = false;
		}
	}

	async function handleUpdate() {
		if (!editingResume) return;

		try {
			const input: UpdateResumeInput = { title: formTitle };
			await updateResume(editingResume.id, input);
			toastSuccess('Resume updated successfully');
			closeEditModal();
			await fetchResumes();
		} catch (err) {
			console.error('Error updating resume:', err);
			toastError('Failed to update resume');
		}
	}

	async function handleDelete(resume: Resume) {
		if (!confirm(`Are you sure you want to delete "${resume.title}"?`)) {
			return;
		}

		try {
			await deleteResume(resume.id);
			toastSuccess('Resume deleted successfully');
			await fetchResumes();
		} catch (err) {
			console.error('Error deleting resume:', err);
			toastError('Failed to delete resume');
		}
	}

	async function handleMarkAsMain(resume: Resume) {
		try {
			await markAsMain(resume.id);
			toastSuccess('Resume marked as main');
			await fetchResumes();
		} catch (err) {
			console.error('Error marking resume as main:', err);
			toastError('Failed to mark resume as main');
		}
	}

	async function handleUnmarkAsMain(resume: Resume) {
		try {
			await unmarkAsMain(resume.id);
			toastSuccess('Resume unmarked as main');
			await fetchResumes();
		} catch (err) {
			console.error('Error unmarking resume as main:', err);
			toastError('Failed to unmark resume as main');
		}
	}

	async function handleMarkAsFeatured(resume: Resume) {
		try {
			await markAsFeatured(resume.id);
			toastSuccess('Resume marked as featured');
			await fetchResumes();
		} catch (err) {
			console.error('Error marking resume as featured:', err);
			toastError('Failed to mark resume as featured');
		}
	}

	async function handleUnmarkAsFeatured(resume: Resume) {
		try {
			await unmarkAsFeatured(resume.id);
			toastSuccess('Resume unmarked as featured');
			await fetchResumes();
		} catch (err) {
			console.error('Error unmarking resume as featured:', err);
			toastError('Failed to unmark resume as featured');
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

	function handleDownload(resume: Resume, language: string = 'en') {
		const url = getResumeDownloadUrl(userId, language);
		window.open(url, '_blank');
	}

	function handlePreview(resume: Resume, language: string = 'en') {
		const url = getResumePreviewUrl(userId, language);
		window.open(url, '_blank');
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Resumes</h1>
			<p class="text-gray-600 dark:text-gray-400 mt-1">Manage your resume files</p>
		</div>
		<button
			onclick={openCreateModal}
			class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
		>
			<Plus class="w-4 h-4" />
			Create Resume
		</button>
	</div>

	{#if loading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600 dark:text-gray-400">Loading resumes...</p>
		</div>
	{:else if error}
		<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
			<p class="text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if resumes.length === 0}
		<div class="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-lg">
			<p class="text-gray-600 dark:text-gray-400">No resumes found</p>
			<p class="text-sm text-gray-500 dark:text-gray-500 mt-2">
				Resumes are created when you generate them using the resume-worker
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
								Title
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								File
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Status
							</th>
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Created
							</th>
							<th
								class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
							>
								Actions
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
												Main
											</span>
										{/if}
										{#if resume.isFeatured}
											<span
												class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
											>
												<Star class="w-3 h-3 mr-1" />
												Featured
											</span>
										{/if}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
									{formatDate(resume.createdAt)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex items-center justify-end gap-2">
										<button
											onclick={() => handlePreview(resume, 'en')}
											class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
											title="Preview (EN)"
										>
											<Eye class="w-4 h-4" />
										</button>
										<button
											onclick={() => handleDownload(resume, 'en')}
											class="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300"
											title="Download (EN)"
										>
											<Download class="w-4 h-4" />
										</button>
										{#if !resume.isMain}
											<button
												onclick={() => handleMarkAsMain(resume)}
												class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
												title="Mark as Main"
											>
												<Check class="w-4 h-4" />
											</button>
										{:else}
											<button
												onclick={() => handleUnmarkAsMain(resume)}
												class="text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
												title="Unmark as Main"
											>
												<X class="w-4 h-4" />
											</button>
										{/if}
										{#if !resume.isFeatured}
											<button
												onclick={() => handleMarkAsFeatured(resume)}
												class="text-yellow-600 hover:text-yellow-900 dark:text-yellow-400 dark:hover:text-yellow-300"
												title="Mark as Featured"
											>
												<Star class="w-4 h-4" />
											</button>
										{:else}
											<button
												onclick={() => handleUnmarkAsFeatured(resume)}
												class="text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
												title="Unmark as Featured"
											>
												<Star class="w-4 h-4 fill-current" />
											</button>
										{/if}
										<button
											onclick={() => openEditModal(resume)}
											class="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-200"
											title="Edit"
										>
											<Edit2 class="w-4 h-4" />
										</button>
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

<!-- Edit Modal -->
{#if showEditModal && editingResume}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
		onclick={closeEditModal}
	>
		<div
			class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Edit Resume</h2>
			<div class="space-y-4">
				<div>
					<label
						for="title"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						Title
					</label>
					<input
						id="title"
						type="text"
						bind:value={formTitle}
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
						placeholder="Resume title"
					/>
				</div>
			</div>
			<div class="mt-6 flex justify-end gap-3">
				<button
					onclick={closeEditModal}
					class="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600"
				>
					Cancel
				</button>
				<button
					onclick={handleUpdate}
					class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
				>
					Save
				</button>
			</div>
		</div>
	</div>
{/if}

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
			<h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">Create Resume</h2>
			<div class="space-y-4">
				<div>
					<label
						for="create-title"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						Title <span class="text-red-500">*</span>
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
						Resume File (PDF) <span class="text-red-500">*</span>
					</label>
					<div class="mt-1 flex items-center gap-4">
						<label
							for="create-file-input"
							class="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-200 dark:hover:bg-gray-600 cursor-pointer transition-colors"
						>
							<Upload class="w-4 h-4" />
							<span>Choose File</span>
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
						Or manually enter file details below
					</p>
				</div>
				<div>
					<label
						for="create-file-path"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
					>
						File Path
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
						File Name
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
						File Size (bytes)
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
			</div>
			<div class="mt-6 flex justify-end gap-3">
				<button
					onclick={closeCreateModal}
					disabled={uploading}
					class="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-700 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50"
				>
					Cancel
				</button>
				<button
					onclick={handleCreate}
					disabled={uploading}
					class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 flex items-center gap-2"
				>
					{#if uploading}
						<div class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
						Creating...
					{:else}
						Create
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

