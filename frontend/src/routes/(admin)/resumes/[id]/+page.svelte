<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getResume,
		updateResume,
		deleteResume,
		markAsMain,
		markAsFeatured,
		unmarkAsMain,
		unmarkAsFeatured,
		getResumeDownloadUrl,
		getResumePreviewUrl,
		type Resume,
		type UpdateResumeInput
	} from '$lib/api/resumes';
	import { Download, Eye, Star, Edit2, Check, X } from 'lucide-svelte';
	import { toastError, toastSuccess } from '$lib/utils/toast';

	const userId = import.meta.env.PUBLIC_DEFAULT_USER_ID || '6ad0d828-f605-45fc-a545-3441e17a015c';

	let resume: Resume | null = $state(null);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showEditModal = $state(false);
	let formTitle = $state('');

	const resumeId = $derived($page.params.id);

	onMount(async () => {
		if (resumeId) {
			await loadResume();
		}
	});

	async function loadResume() {
		if (!resumeId) return;
		loading = true;
		error = null;
		try {
			const loadedResume = await getResume(resumeId);
			if (!loadedResume) {
				error = 'Resume not found';
				return;
			}
			resume = loadedResume;
			formTitle = resume.title;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load resume';
			console.error('Error loading resume:', err);
		} finally {
			loading = false;
		}
	}

	function openEditModal() {
		if (!resume) return;
		formTitle = resume.title;
		showEditModal = true;
	}

	async function handleUpdate() {
		if (!resume) return;

		try {
			const input: UpdateResumeInput = { title: formTitle };
			await updateResume(resume.id, input);
			toastSuccess('Resume updated successfully');
			showEditModal = false;
			await loadResume();
		} catch (err) {
			console.error('Error updating resume:', err);
			toastError('Failed to update resume');
		}
	}

	async function handleDelete() {
		if (!resume || !confirm(`Are you sure you want to delete "${resume.title}"?`)) {
			return;
		}

		try {
			await deleteResume(resume.id);
			toastSuccess('Resume deleted successfully');
			await goto('/resumes');
		} catch (err) {
			console.error('Error deleting resume:', err);
			toastError('Failed to delete resume');
		}
	}

	async function handleMarkAsMain() {
		if (!resume) return;
		try {
			await markAsMain(resume.id);
			toastSuccess('Resume marked as main');
			await loadResume();
		} catch (err) {
			console.error('Error marking resume as main:', err);
			toastError('Failed to mark resume as main');
		}
	}

	async function handleUnmarkAsMain() {
		if (!resume) return;
		try {
			await unmarkAsMain(resume.id);
			toastSuccess('Resume unmarked as main');
			await loadResume();
		} catch (err) {
			console.error('Error unmarking resume as main:', err);
			toastError('Failed to unmark resume as main');
		}
	}

	async function handleMarkAsFeatured() {
		if (!resume) return;
		try {
			await markAsFeatured(resume.id);
			toastSuccess('Resume marked as featured');
			await loadResume();
		} catch (err) {
			console.error('Error marking resume as featured:', err);
			toastError('Failed to mark resume as featured');
		}
	}

	async function handleUnmarkAsFeatured() {
		if (!resume) return;
		try {
			await unmarkAsFeatured(resume.id);
			toastSuccess('Resume unmarked as featured');
			await loadResume();
		} catch (err) {
			console.error('Error unmarking resume as featured:', err);
			toastError('Failed to unmark resume as featured');
		}
	}

	function handleDownload(language: string = 'en') {
		const url = getResumeDownloadUrl(userId, language);
		window.open(url, '_blank');
	}

	function handlePreview(language: string = 'en') {
		const url = getResumePreviewUrl(userId, language);
		window.open(url, '_blank');
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
	<div class="mb-6 flex items-center justify-between">
		<a href="/resumes" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
			← Back to Resumes
		</a>
		<div class="flex gap-2">
			{#if resume}
				<button
					onclick={openEditModal}
					class="flex items-center gap-2 px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors"
				>
					<Edit2 class="w-4 h-4" />
					Edit
				</button>
				<button
					onclick={handleDelete}
					class="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors"
				>
					Delete
				</button>
			{/if}
		</div>
	</div>

	{#if loading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600 dark:text-gray-400">Loading resume...</p>
		</div>
	{:else if error}
		<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
			<p class="text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if resume}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
			<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
				<h1 class="text-2xl font-bold text-gray-900 dark:text-white">{resume.title}</h1>
			</div>

			<div class="px-6 py-4 space-y-6">
				<div>
					<h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">File Information</h2>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">File Name</label>
							<p class="mt-1 text-sm text-gray-900 dark:text-white">{resume.fileName}</p>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">File Size</label>
							<p class="mt-1 text-sm text-gray-900 dark:text-white">{formatFileSize(resume.fileSize)}</p>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">File Path</label>
							<p class="mt-1 text-sm text-gray-900 dark:text-white break-all">{resume.filePath}</p>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Created</label>
							<p class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(resume.createdAt)}</p>
						</div>
					</div>
				</div>

				<div>
					<h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Status</h2>
					<div class="flex items-center gap-4">
						{#if resume.isMain}
							<span
								class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
							>
								<Check class="w-4 h-4 mr-1" />
								Main Resume
							</span>
						{:else}
							<button
								onclick={handleMarkAsMain}
								class="inline-flex items-center px-3 py-1 rounded-md text-sm font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 border border-blue-300 dark:border-blue-700 hover:bg-blue-50 dark:hover:bg-blue-900/20"
							>
								<Check class="w-4 h-4 mr-1" />
								Mark as Main
							</button>
						{/if}
						{#if resume.isFeatured}
							<span
								class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
							>
								<Star class="w-4 h-4 mr-1 fill-current" />
								Featured
							</span>
						{:else}
							<button
								onclick={handleMarkAsFeatured}
								class="inline-flex items-center px-3 py-1 rounded-md text-sm font-medium text-yellow-600 hover:text-yellow-800 dark:text-yellow-400 dark:hover:text-yellow-300 border border-yellow-300 dark:border-yellow-700 hover:bg-yellow-50 dark:hover:bg-yellow-900/20"
							>
								<Star class="w-4 h-4 mr-1" />
								Mark as Featured
							</button>
						{/if}
					</div>
					{#if resume.isMain}
						<button
							onclick={handleUnmarkAsMain}
							class="mt-2 text-sm text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
						>
							Unmark as Main
						</button>
					{/if}
					{#if resume.isFeatured}
						<button
							onclick={handleUnmarkAsFeatured}
							class="mt-2 ml-4 text-sm text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
						>
							Unmark as Featured
						</button>
					{/if}
				</div>

				<div>
					<h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Actions</h2>
					<div class="flex gap-4">
						<button
							onclick={() => handlePreview('en')}
							class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
						>
							<Eye class="w-4 h-4" />
							Preview (EN)
						</button>
						<button
							onclick={() => handleDownload('en')}
							class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
						>
							<Download class="w-4 h-4" />
							Download (EN)
						</button>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Edit Modal -->
{#if showEditModal && resume}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
		onclick={() => (showEditModal = false)}
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
					onclick={() => (showEditModal = false)}
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

