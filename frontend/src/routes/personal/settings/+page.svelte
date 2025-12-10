<script lang="ts">
	import { onMount } from 'svelte';
	import { getUserProfile, upsertUserProfile, type UserProfile } from '$lib/api/userprofiles';
	import { getUserPreferences, updateUserPreferences, type UserPreferences } from '$lib/api/userpreferences';
	import { toastError, toastSuccess } from '$lib/utils/toast';

	let profile: UserProfile | null = $state(null);
	let aboutMe = $state('');
	let loading = $state(true);
	let saving = $state(false);
	let error: string | null = $state(null);

	let preferences: UserPreferences | null = $state(null);
	let loadingPreferences = $state(true);
	let savingPreferences = $state(false);
	let errorPreferences: string | null = $state(null);
	let defaultLanguage = $state('en');
	let defaultCurrency = $state('USD');

	// Common language options (ISO 639-1)
	const languages = [
		{ code: 'en', name: 'English' },
		{ code: 'pt', name: 'Portuguese' },
		{ code: 'es', name: 'Spanish' },
		{ code: 'fr', name: 'French' },
		{ code: 'de', name: 'German' },
		{ code: 'it', name: 'Italian' },
		{ code: 'ja', name: 'Japanese' },
		{ code: 'zh', name: 'Chinese' },
		{ code: 'ko', name: 'Korean' },
		{ code: 'ru', name: 'Russian' },
		{ code: 'ar', name: 'Arabic' },
		{ code: 'hi', name: 'Hindi' },
		{ code: 'nl', name: 'Dutch' },
		{ code: 'pl', name: 'Polish' },
		{ code: 'tr', name: 'Turkish' }
	];

	// Common currency options (ISO 4217)
	const currencies = [
		{ code: 'USD', name: 'US Dollar (USD)' },
		{ code: 'EUR', name: 'Euro (EUR)' },
		{ code: 'GBP', name: 'British Pound (GBP)' },
		{ code: 'JPY', name: 'Japanese Yen (JPY)' },
		{ code: 'CNY', name: 'Chinese Yuan (CNY)' },
		{ code: 'BRL', name: 'Brazilian Real (BRL)' },
		{ code: 'INR', name: 'Indian Rupee (INR)' },
		{ code: 'CAD', name: 'Canadian Dollar (CAD)' },
		{ code: 'AUD', name: 'Australian Dollar (AUD)' },
		{ code: 'CHF', name: 'Swiss Franc (CHF)' },
		{ code: 'MXN', name: 'Mexican Peso (MXN)' },
		{ code: 'RUB', name: 'Russian Ruble (RUB)' },
		{ code: 'KRW', name: 'South Korean Won (KRW)' },
		{ code: 'SGD', name: 'Singapore Dollar (SGD)' },
		{ code: 'HKD', name: 'Hong Kong Dollar (HKD)' },
		{ code: 'NZD', name: 'New Zealand Dollar (NZD)' },
		{ code: 'SEK', name: 'Swedish Krona (SEK)' },
		{ code: 'NOK', name: 'Norwegian Krone (NOK)' },
		{ code: 'DKK', name: 'Danish Krone (DKK)' },
		{ code: 'PLN', name: 'Polish Złoty (PLN)' },
		{ code: 'TRY', name: 'Turkish Lira (TRY)' },
		{ code: 'ZAR', name: 'South African Rand (ZAR)' },
		{ code: 'AED', name: 'UAE Dirham (AED)' },
		{ code: 'SAR', name: 'Saudi Riyal (SAR)' },
		{ code: 'THB', name: 'Thai Baht (THB)' },
		{ code: 'MYR', name: 'Malaysian Ringgit (MYR)' },
		{ code: 'IDR', name: 'Indonesian Rupiah (IDR)' },
		{ code: 'PHP', name: 'Philippine Peso (PHP)' },
		{ code: 'VND', name: 'Vietnamese Dong (VND)' },
		{ code: 'ILS', name: 'Israeli Shekel (ILS)' }
	];

	onMount(async () => {
		await Promise.all([loadProfile(), loadPreferences()]);
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

	async function loadPreferences() {
		loadingPreferences = true;
		errorPreferences = null;
		try {
			preferences = await getUserPreferences();
			defaultLanguage = preferences.defaultLanguage || 'en';
			defaultCurrency = preferences.defaultCurrency || 'USD';
		} catch (err) {
			errorPreferences = err instanceof Error ? err.message : 'Failed to load preferences';
			console.error('Error loading preferences:', err);
			// Don't show toast for initial load failure - preferences might not exist yet
		} finally {
			loadingPreferences = false;
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

	async function handleSavePreferences() {
		if (savingPreferences) return;

		savingPreferences = true;
		errorPreferences = null;
		try {
			preferences = await updateUserPreferences({
				defaultLanguage,
				defaultCurrency
			});
			toastSuccess('Preferences saved successfully');
		} catch (err) {
			errorPreferences = err instanceof Error ? err.message : 'Failed to save preferences';
			console.error('Error saving preferences:', err);
			toastError('Failed to save preferences');
		} finally {
			savingPreferences = false;
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

		<!-- Job Application Preferences Section -->
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mt-6">
			<h2 class="text-xl font-semibold mb-4 text-gray-900 dark:text-white">Job Application Preferences</h2>
			<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
				Set your default language and currency for job applications. These will be automatically applied when
				creating new job applications if not specified.
			</p>

			{#if loadingPreferences}
				<div class="text-center py-8">
					<div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
					<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">Loading preferences...</p>
				</div>
			{:else}
				<div class="space-y-4">
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label
								for="default-language"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
							>
								Default Language
							</label>
							<select
								id="default-language"
								bind:value={defaultLanguage}
								class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
							>
								{#each languages as lang}
									<option value={lang.code}>{lang.name} ({lang.code})</option>
								{/each}
							</select>
							<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
								ISO 639-1 language code (2 characters)
							</p>
						</div>

						<div>
							<label
								for="default-currency"
								class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
							>
								Default Currency
							</label>
							<select
								id="default-currency"
								bind:value={defaultCurrency}
								class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
							>
								{#each currencies as curr}
									<option value={curr.code}>{curr.name}</option>
								{/each}
							</select>
							<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
								ISO 4217 currency code (3 characters)
							</p>
						</div>
					</div>

					{#if errorPreferences}
						<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
							<p class="text-sm text-red-800 dark:text-red-200">{errorPreferences}</p>
						</div>
					{/if}

					<div class="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
						<button
							onclick={handleSavePreferences}
							disabled={savingPreferences}
							class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
						>
							{#if savingPreferences}
								<div class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
								Saving...
							{:else}
								Save Preferences
							{/if}
						</button>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>

