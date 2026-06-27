<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		listSocialMediaPosts,
		deleteSocialMediaPost,
		type SocialMediaPost,
		type Platform,
		type PostStatus
	} from '$lib/api/socialmediaposts';
	import PageHeader from './_sections/PageHeader.svelte';
	import SearchBar from './_sections/SearchBar.svelte';
	import PostFilters from './_sections/PostFilters.svelte';
	import PostsTable from './_sections/PostsTable.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';

	let posts: SocialMediaPost[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');
	let platform: Platform | undefined = $state(undefined);
	let status: PostStatus | undefined = $state(undefined);

	onMount(async () => {
		await fetchPosts();
	});

	async function fetchPosts() {
		loading = true;
		error = null;
		try {
			const filters: { platform?: Platform; status?: PostStatus } = {};
			if (platform) filters.platform = platform;
			if (status) filters.status = status;
			posts = await listSocialMediaPosts(filters);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load posts';
			console.error('Error fetching posts:', err);
		} finally {
			loading = false;
		}
	}

	function handleDelete(id: string) {
		if (confirm('Are you sure you want to delete this post?')) {
			deleteSocialMediaPost(id)
				.then(() => {
					posts = posts.filter((p) => p.id !== id);
				})
				.catch((err) => {
					console.error('Error deleting post:', err);
					alert('Failed to delete post');
				});
		}
	}

	function handleView(post: SocialMediaPost) {
		goto(`/posts/${post.id}`);
	}

	function handleCreate() {
		goto('/posts/new');
	}

	const filteredPosts = $derived(() => {
		if (!searchQuery) return posts;
		const query = searchQuery.toLowerCase();
		return posts.filter(
			(post) =>
				post.title.toLowerCase().includes(query) || post.content.toLowerCase().includes(query)
		);
	});
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<PageHeader onCreateClick={handleCreate} />
	<SearchBar bind:searchQuery />
	<PostFilters bind:platform bind:status />
	
	{#if error}
		<ErrorState message={error} onRetry={fetchPosts} />
	{:else if loading}
		<LoadingState message="Loading posts..." />
	{:else}
		<PostsTable posts={filteredPosts()} onDelete={handleDelete} onView={handleView} />
	{/if}
</div>
