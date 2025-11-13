<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiClient } from '$lib/api/client';
	import { authStore, type AuthUser } from '$lib';

	type LoginPayload = AuthUser & { token?: string | null };

	interface LoginResponse {
		data: LoginPayload | null;
		success: boolean;
	}

	let email = '';
	let password = '';
	let loading = false;
	let error = '';

	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		error = '';
		loading = true;

		try {
			const response = await apiClient.post<LoginResponse>('/auth/login', { email, password });

			const payload = response.data?.data;
			if (payload?.token) {
				const { token, ...user } = payload;
				authStore.setSession(user, token);
				goto('/');
				return;
			}

			error = 'Unable to sign in. Please try again.';
		} catch (err) {
			error = 'Unable to sign in. Check your credentials and try again.';
			console.error(err);
		} finally {
			loading = false;
		}
	};
</script>

<form class="space-y-6" on:submit={handleSubmit}>
	<div class="space-y-2">
		<h2 class="text-lg font-semibold">Sign in</h2>
		<p class="text-sm text-slate-400">Enter your email and password to continue.</p>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
			{error}
		</p>
	{/if}

	<div class="space-y-4">
		<label class="flex flex-col gap-1 text-sm text-slate-200">
			Email
			<input
				class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
				type="email"
				placeholder="you@example.com"
				bind:value={email}
				required
			/>
		</label>
		<label class="flex flex-col gap-1 text-sm text-slate-200">
			Password
			<input
				class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
				type="password"
				placeholder="********"
				bind:value={password}
				required
			/>
		</label>
	</div>

	<button
		class="w-full rounded bg-emerald-500 px-3 py-2 text-sm font-semibold text-slate-900 transition-opacity disabled:opacity-50"
		type="submit"
		disabled={loading}
	>
		{#if loading}
			Signing in...
		{:else}
			Sign in
		{/if}
	</button>

	<p class="text-center text-xs text-slate-400">
		Don't have an account?
		<a class="text-primary hover:underline" href="/auth/register">Create one</a>
	</p>
</form>
