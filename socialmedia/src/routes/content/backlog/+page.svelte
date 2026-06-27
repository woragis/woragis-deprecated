<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getContentBacklog, type ContentPost } from '$lib/api/content';
	import PageHeader from '../posts/_sections/PageHeader.svelte';
	import ContentPostsTable from '../posts/_sections/ContentPostsTable.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';

	let contentBacklog: ContentPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchBacklog();
	});

	async function fetchBacklog() {
		loading = true;
		error = null;
		try {
			contentBacklog = await getContentBacklog();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load content backlog';
			console.error('Error fetching content backlog:', err);
		} finally {
			loading = false;
		}
	}

	function handleView(post: ContentPost) {
		goto(`/content/posts/${post.id}`);
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="page-header">
		<div class="header-content">
			<h1 class="page-title">Content Backlog</h1>
			<p class="page-subtitle">Posts ready for repurposing to social media platforms</p>
		</div>
	</div>
	
	{#if error}
		<ErrorState message={error} onRetry={fetchBacklog} />
	{:else if loading}
		<LoadingState message="Loading content backlog..." />
	{:else}
		<ContentPostsTable contentPosts={contentBacklog} onView={handleView} />
	{/if}
</div>

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}
	
	.header-content {
		flex: 1;
	}
	
	.page-title {
		margin: 0 0 var(--spacing-xs) 0;
		font-size: var(--font-size-2xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}
	
	.page-subtitle {
		margin: 0;
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
	}
</style>
