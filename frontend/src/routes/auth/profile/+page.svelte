<script lang="ts">
	import { onMount } from 'svelte';
	import { getCurrentUser, updateProfile } from '$lib/api/auth';
	import { authStore } from '$lib';
	import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';

	let loading = false;
	let saving = false;
	let error = '';
	let phoneNumber = '';
	let preferredLocale = '';

	const loadProfile = async () => {
		loading = true;
		error = '';
		try {
			const user = await getCurrentUser();
			phoneNumber = user.phone_number || '';
			preferredLocale = user.preferred_locale || '';
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to load profile.');
			toastError(error);
		} finally {
			loading = false;
		}
	};

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		saving = true;
		error = '';

		try {
			const updated = await updateProfile({
				phone_number: phoneNumber.trim() || undefined,
				preferred_locale: preferredLocale.trim() || undefined
			});

			authStore.updateUser({
				phone_number: updated.phone_number,
				preferred_locale: updated.preferred_locale
			});

			toastSuccess('Profile updated successfully.');
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Unable to update profile.');
			toastError(error);
		} finally {
			saving = false;
		}
	};

	onMount(() => {
		loadProfile();
	});
</script>

<section class="space-y-6">
	<div class="space-y-2">
		<h2 class="text-xl font-semibold text-slate-100">Profile Settings</h2>
		<p class="text-sm text-slate-400">
			Manage your profile information, including phone number for WhatsApp notifications.
		</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">{error}</p>
	{/if}

	{#if loading}
		<div class="card">
			<div class="flex items-center justify-center py-8">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
			</div>
		</div>
	{:else}
		<form on:submit={handleSubmit} class="space-y-6">
			<div class="card space-y-6">
				<div>
					<label for="phone_number" class="block text-sm font-medium text-slate-300 mb-2">
						Phone Number
					</label>
					<input
						id="phone_number"
						type="tel"
						bind:value={phoneNumber}
						placeholder="+1234567890"
						class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
					<p class="mt-1 text-xs text-slate-400">
						Include country code (e.g., +1 for US, +55 for Brazil). Used for WhatsApp notifications.
					</p>
				</div>

				<div>
					<label for="preferred_locale" class="block text-sm font-medium text-slate-300 mb-2">
						Preferred Locale
					</label>
					<input
						id="preferred_locale"
						type="text"
						bind:value={preferredLocale}
						placeholder="en"
						maxlength="10"
						class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
					<p class="mt-1 text-xs text-slate-400">
						Language code (e.g., en, pt-BR, es). Used for system preferences.
					</p>
				</div>

				<div class="flex justify-end gap-3 pt-4 border-t border-slate-700">
					<button
						type="submit"
						disabled={saving}
						class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors font-medium"
					>
						{saving ? 'Saving...' : 'Save Changes'}
					</button>
				</div>
			</div>
		</form>
	{/if}
</section>

