<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { listContentPosts, type ContentPost } from '$lib/api/content';
	import PageHeader from './_sections/PageHeader.svelte';
	import ContentPostsTable from './_sections/ContentPostsTable.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';

	let contentPosts: ContentPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		await fetchContentPosts();
	});

	async function fetchContentPosts() {
		loading = true;
		error = null;
		try {
			contentPosts = await listContentPosts();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load content posts';
			console.error('Error fetching content posts:', err);
		} finally {
			loading = false;
		}
	}

	function handleView(post: ContentPost) {
		goto(`/content/posts/${post.id}`);
	}

	function handleCreate() {
		goto('/content/posts/new');
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<PageHeader onCreateClick={handleCreate} />
	
	{#if error}
		<ErrorState message={error} onRetry={fetchContentPosts} />
	{:else if loading}
		<LoadingState message="Loading content posts..." />
	{:else}
		<ContentPostsTable {contentPosts} onView={handleView} />
	{/if}
</div>
