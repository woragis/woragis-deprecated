<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { ArrowLeft, Calendar, FileText, Globe, Tag as TagIcon, Folder, Code } from 'lucide-svelte';
	import { getPost, getPostCategories, getPostTags, getPostSkills, type Post } from '$lib/api/landing';
	import { listCategories, listTags, listSkills, type Category, type Tag, type Skill } from '$lib/api/landing';
	import PageHero from '$lib/components/PageHero.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';

	let post: Post | null = $state(null);
	let categories: Category[] = $state([]);
	let tags: Tag[] = $state([]);
	let skills: Skill[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	onMount(async () => {
		const postId = $page.params.id;
		if (!postId) {
			error = 'Post ID is required';
			loading = false;
			return;
		}

		try {
			post = await getPost(postId);
			// Load relationships
			const [categoryIds, tagIds, skillIds] = await Promise.all([
				getPostCategories(postId).catch(() => [] as Category[]),
				getPostTags(postId).catch(() => [] as Tag[]),
				getPostSkills(postId).catch(() => [] as string[])
			]);

			// Get full details for relationships
			const [allCategories, allTags, allSkills] = await Promise.all([
				listCategories().catch(() => [] as Category[]),
				listTags().catch(() => [] as Tag[]),
				listSkills().catch(() => [] as Skill[])
			]);

			categories = categoryIds;
			tags = tagIds;
			// getPostSkills returns string[], so we need to map to Skill objects
			skills = allSkills.filter((s) => skillIds.includes(String(s.id)));
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load post';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<PageHero
		title={post?.title || 'Post Details'}
		description={post?.excerpt}
		gradientFrom="from-blue-950/30"
		gradientVia="via-indigo-950/30"
		gradientTo="to-blue-950/30"
	>
		<button
			slot="actions"
			class="flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-blue-500/50 hover:bg-slate-800/80"
			onclick={() => goto('/landing/posts')}
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Posts
		</button>
	</PageHero>

	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if loading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-blue-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading post...</p>
			</div>
		{:else if error || !post}
			<EmptyState
				title={error || 'Post not found'}
				description="The post you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-blue-700"
					onclick={() => goto('/landing/posts')}
				>
					Return to Posts
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Status"
						value={post.status}
						accentColor="blue"
					/>
					<StatCard
						label="Created"
						value={formatDate(post.created_at)}
						accentColor="purple"
					/>
					<StatCard
						label="Featured"
						value={post.featured ? 'Yes' : 'No'}
						accentColor={post.featured ? 'emerald' : 'slate'}
					/>
				</div>

				<!-- Content Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Content</h2>
					<div class="prose prose-invert max-w-none">
						<div class="text-slate-300 whitespace-pre-wrap">{post.content}</div>
					</div>
				</div>

				<!-- Metadata Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Metadata</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#if post.slug}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Slug</p>
								<p class="text-sm font-medium text-slate-200">/{post.slug}</p>
							</div>
						{/if}
						{#if post.published_at}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Published</p>
								<p class="text-sm font-medium text-slate-200">{formatDate(post.published_at)}</p>
							</div>
						{/if}
						{#if post.meta_title}
							<div>
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Meta Title</p>
								<p class="text-sm font-medium text-slate-200">{post.meta_title}</p>
							</div>
						{/if}
						{#if post.meta_description}
							<div class="md:col-span-2">
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Meta Description</p>
								<p class="text-sm font-medium text-slate-200">{post.meta_description}</p>
							</div>
						{/if}
						{#if post.featured_image}
							<div class="md:col-span-2">
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Featured Image</p>
								<a
									href={post.featured_image}
									target="_blank"
									class="text-sm font-medium text-blue-400 hover:text-blue-300"
								>
									{post.featured_image}
								</a>
							</div>
						{/if}
					</div>
				</div>

				<!-- Relationships -->
				{#if categories.length > 0 || tags.length > 0 || skills.length > 0}
					<div
						class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
					>
						<h2 class="text-lg font-semibold text-slate-100 mb-4">Relationships</h2>
						<div class="space-y-4">
							{#if categories.length > 0}
								<div>
									<div class="flex items-center gap-2 mb-2">
										<Folder class="h-4 w-4 text-slate-400" />
										<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Categories</p>
									</div>
									<div class="flex flex-wrap gap-2">
										{#each categories as category}
											<span
												class="inline-flex items-center rounded-full bg-blue-500/20 px-3 py-1 text-xs font-medium text-blue-300"
											>
												{category.name}
											</span>
										{/each}
									</div>
								</div>
							{/if}
							{#if tags.length > 0}
								<div>
									<div class="flex items-center gap-2 mb-2">
										<TagIcon class="h-4 w-4 text-slate-400" />
										<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Tags</p>
									</div>
									<div class="flex flex-wrap gap-2">
										{#each tags as tag}
											<span
												class="inline-flex items-center rounded-full bg-purple-500/20 px-3 py-1 text-xs font-medium text-purple-300"
											>
												{tag.name}
											</span>
										{/each}
									</div>
								</div>
							{/if}
							{#if skills.length > 0}
								<div>
									<div class="flex items-center gap-2 mb-2">
										<Code class="h-4 w-4 text-slate-400" />
										<p class="text-xs font-medium uppercase tracking-wider text-slate-400">Skills</p>
									</div>
									<div class="flex flex-wrap gap-2">
										{#each skills as skill}
											<span
												class="inline-flex items-center rounded-full bg-emerald-500/20 px-3 py-1 text-xs font-medium text-emerald-300"
											>
												{skill.name}
											</span>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>

