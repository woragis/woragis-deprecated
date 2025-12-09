<script lang="ts">
	import { onMount } from 'svelte';
	import { getUserProfile, upsertUserProfile, type UserProfile } from '$lib/api/userprofiles';
	import { toastError, toastSuccess } from '$lib/utils/toast';

	let profile: UserProfile | null = $state(null);
	let aboutMe = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let error: string | null = $state(null);

	onMount(async () => {
		await loadProfile();
	});

	async function loadProfile() {
		loading = true;
		error = null;
		try {
			profile = await getUserProfile();
			aboutMe = profile.aboutMe || '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load profile';
			console.error('Error loading profile:', err);
			toastError('Failed to load profile');
		} finally {
			loading = false;
		}
	}

	async function handleSave() {
		if (saving) return;

		saving = true;
		error = null;
		try {
			profile = await upsertUserProfile({ aboutMe });
			toastSuccess('Profile saved successfully');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save profile';
			console.error('Error saving profile:', err);
			toastError('Failed to save profile');
		} finally {
			saving = false;
		}
	}
</script>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-6">
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Settings</h1>
		<p class="text-gray-600 dark:text-gray-400 mt-1">Manage your knowledge base and profile information</p>
	</div>

	{#if loading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600 dark:text-gray-400">Loading profile...</p>
		</div>
	{:else if error && !profile}
		<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
			<p class="text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-4 text-gray-900 dark:text-white">Knowledge Base</h2>
			<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
				This information will be used by AI assistants to provide better context about you when discussing
				job applications, projects, and other topics. Include information about your languages, hobbies,
				interests, background, and anything else relevant.
			</p>

			<div class="space-y-4">
				<div>
					<label
						for="about-me"
						class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
					>
						About Me
					</label>
					<textarea
						id="about-me"
						bind:value={aboutMe}
						rows="20"
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white font-mono text-sm"
						placeholder="Write about yourself here. You can use HTML formatting.&#10;&#10;Example:&#10;&lt;p&gt;I am a software engineer with expertise in...&lt;/p&gt;&#10;&lt;h3&gt;Languages&lt;/h3&gt;&#10;&lt;ul&gt;&#10;  &lt;li&gt;English (Native)&lt;/li&gt;&#10;  &lt;li&gt;Portuguese (Fluent)&lt;/li&gt;&#10;&lt;/ul&gt;&#10;&lt;h3&gt;Hobbies&lt;/h3&gt;&#10;&lt;ul&gt;&#10;  &lt;li&gt;Reading&lt;/li&gt;&#10;  &lt;li&gt;Hiking&lt;/li&gt;&#10;&lt;/ul&gt;"
					></textarea>
					<p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
						You can use HTML tags for formatting (e.g., &lt;p&gt;, &lt;h3&gt;, &lt;ul&gt;, &lt;li&gt;)
					</p>
				</div>

				{#if error}
					<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
						<p class="text-sm text-red-800 dark:text-red-200">{error}</p>
					</div>
				{/if}

				<div class="flex justify-end gap-3">
					<button
						onclick={handleSave}
						disabled={saving}
						class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if saving}
							<div class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
							Saving...
						{:else}
							Save
						{/if}
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

