<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import SidebarNav from '$lib/components/SidebarNav.svelte';
	import { authNav, primaryNav } from '$lib/navigation';
	import { authStore } from '$lib';
	import type { AuthUser } from '$lib';
	import { initAuthLifecycle, performLogout, teardownAuthLifecycle } from '$lib/auth/lifecycle';
	import { queryClient } from '@clients/queryClient';
	import ToastProvider from '$lib/providers/ToastProvider.svelte';
	import '../app.css';

	let sidebarOpen = false;
	let mediaQuery: MediaQueryList | null = null;
	let isAuthenticated = false;
	let currentUser: AuthUser | null = null;

	const closeSidebar = () => {
		sidebarOpen = false;
	};

	const handleLogout = async () => {
		await performLogout();
	};

	onMount(() => {
		initAuthLifecycle();
		const unsubscribe = authStore.subscribe((state) => {
			isAuthenticated = state.isAuthenticated;
			currentUser = state.user;
		});

		mediaQuery = window.matchMedia('(min-width: 640px)');
		const handleChange = (event: MediaQueryListEvent) => {
			if (event.matches) {
				closeSidebar();
			}
		};

		mediaQuery.addEventListener('change', handleChange);

		return () => {
			mediaQuery?.removeEventListener('change', handleChange);
			unsubscribe();
			teardownAuthLifecycle();
		};
	});
</script>

<svelte:head>
	<title>Woragis Console</title>
</svelte:head>

<QueryClientProvider client={queryClient}>
	<ToastProvider />
	<div class="app-shell">
		<SidebarNav
			title="Woragis Console"
			primary={primaryNav}
			secondaryLabel="Auth"
			secondary={authNav}
			currentPath={$page.url.pathname}
			open={sidebarOpen}
			user={currentUser}
			isAuthenticated={isAuthenticated}
			on:close={closeSidebar}
			on:logout={handleLogout}
		/>

		<div class="app-shell__main">
			<header class="topbar">
				<button
					class="topbar__trigger"
					type="button"
					aria-label="Open navigation"
					aria-expanded={sidebarOpen}
					on:click={() => (sidebarOpen = true)}
				>
					<span aria-hidden="true">☰</span>
				</button>
				<h1 class="topbar__title">Woragis Console</h1>
				{#if isAuthenticated}
					<button class="logout" type="button" on:click={handleLogout}>
						Sign out
					</button>
				{/if}
			</header>

			<main class="content">
				<slot />
			</main>
		</div>
	</div>
</QueryClientProvider>

<style>
	.app-shell {
		display: flex;
		min-height: 100vh;
		background: #020617;
		color: #f8fafc;
	}

	.app-shell__main {
		display: flex;
		flex: 1 1 auto;
		flex-direction: column;
		min-height: 100vh;
		background: radial-gradient(circle at top left, rgba(37, 99, 235, 0.15), transparent 55%),
			radial-gradient(circle at bottom right, rgba(7, 89, 133, 0.2), transparent 60%),
			rgba(2, 6, 23, 0.92);
		backdrop-filter: blur(6px);
	}

	.topbar {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.5rem;
		border-bottom: 1px solid rgba(15, 118, 110, 0.2);
		background: rgba(15, 15, 15, 0.4);
		backdrop-filter: blur(12px);
		box-shadow: 0 12px 32px rgba(0, 0, 0, 0.8);
	}

	.topbar__trigger {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 0.75rem;
		border: 1px solid rgba(51, 65, 85, 0.7);
		background: rgba(15, 23, 42, 0.7);
		color: #e2e8f0;
		font-size: 1.25rem;
		transition: background 150ms ease, transform 150ms ease, border-color 150ms ease;
	}

	.topbar__trigger:hover,
	.topbar__trigger:focus-visible {
		background: rgba(51, 65, 85, 0.7);
		border-color: rgba(94, 234, 212, 0.4);
		transform: translateY(-1px);
		outline: none;
	}

	.topbar__title {
		font-size: 1rem;
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;
	}

	.logout {
		margin-left: auto;
		padding: 0.35rem 0.75rem;
		border-radius: 0.5rem;
		border: 1px solid rgba(94, 234, 212, 0.4);
		background: rgba(15, 23, 42, 0.7);
		color: #f1f5f9;
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		font-weight: 600;
		transition: background 120ms ease, border-color 120ms ease, color 120ms ease;
	}

	.logout:hover,
	.logout:focus-visible {
		background: rgba(94, 234, 212, 0.1);
		border-color: rgba(94, 234, 212, 0.6);
		color: #f8fafc;
		outline: none;
	}

	.content {
		display: flex;
		flex-direction: column;
		width: 100%;
		max-width: 72rem;
		margin: 0 auto;
		flex: 1 1 auto;
		padding: 2.5rem 1.5rem 3.5rem;
		gap: 1.75rem;
	}

	@media (min-width: 640px) {
		.topbar {
			display: none;
		}

		.app-shell__main {
			padding-left: 0;
		}
	}

	@media (max-width: 639px) {
		.app-shell {
			flex-direction: column;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.topbar__trigger {
			transition: none;
		}
	}
</style>
