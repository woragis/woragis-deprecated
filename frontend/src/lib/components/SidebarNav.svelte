<script lang="ts">
	import type { NavItem } from '$lib/navigation';
	import type { AuthUser } from '$lib';
	import { createEventDispatcher, onMount } from 'svelte';

	export let title = 'Woragis Console';
	export let primary: NavItem[] = [];
	export let secondaryLabel = 'Auth';
	export let secondary: NavItem[] = [];
	export let currentPath = '';
	export let open = false;
	export let user: AuthUser | null = null;
	export let isAuthenticated = false;

	const dispatch = createEventDispatcher<{ close: void; logout: void }>();

	const itemMatches = (item: NavItem) => (item.match ? item.match(currentPath) : currentPath === item.href);

	const handleNavigate = () => {
		dispatch('close');
	};

	const getDisplayName = (value: AuthUser | null) => {
		if (!value) {
			return 'Guest';
		}

		const preferred = value.display_name?.trim();
		if (preferred) {
			return preferred;
		}

		const emailLocalPart = value.email?.split('@')?.[0];
		if (emailLocalPart) {
			return emailLocalPart;
		}

		return 'Account';
	};

	let displayName = '';
	let initials = '';

	$: displayName = getDisplayName(user);
	$: initials =
		displayName
			.split(/\s+/)
			.filter(Boolean)
			.slice(0, 2)
			.map((part) => part[0]?.toUpperCase() ?? '')
			.join('') || 'U';

	let previouslyOpen = open;

	onMount(() => {
		previouslyOpen = open;
	});
</script>

<div class="SidebarNav">
	<div
		class="overlay"
		class:overlay--visible={open}
		tabindex="-1"
		role="button"
		aria-hidden={!open}
		on:click={() => dispatch('close')}
	>
		<span class="sr-only">Close navigation</span>
	</div>

	<nav class="panel" class:panel--open={open} aria-label="Main navigation">
		<div class="panel__header">
			<h1 class="panel__title">{title}</h1>
		</div>

		<div class="panel__content">
			{#if isAuthenticated && user}
				<section class="panel__section panel__section--profile" aria-label="Profile">
					<div class="profile-card">
						<div class="profile-card__avatar" aria-hidden="true">
							{initials}
						</div>
						<div class="profile-card__meta">
							<p class="profile-card__name">{displayName}</p>
							<p class="profile-card__email">{user.email}</p>
						</div>
					</div>
					<button class="profile-card__logout" type="button" on:click={() => dispatch('logout')}>
						Log out
					</button>
				</section>
			{/if}

			<section class="panel__section" aria-label="Primary navigation">
				{#each primary as item (item.href)}
					<a
						class="nav-link"
						class:nav-link--active={itemMatches(item)}
						aria-current={itemMatches(item) ? 'page' : undefined}
						href={item.href}
						on:click={handleNavigate}
					>
						{item.label}
					</a>
				{/each}
			</section>

			{#if secondary.length}
				<section class="panel__section panel__section--secondary" aria-label="Secondary navigation">
					<div class="panel__section-label">{secondaryLabel}</div>
					{#each secondary as item (item.href)}
						<a
							class="nav-link"
							class:nav-link--active={itemMatches(item)}
							aria-current={itemMatches(item) ? 'page' : undefined}
							href={item.href}
							on:click={handleNavigate}
						>
							{item.label}
						</a>
					{/each}
				</section>
			{/if}
		</div>
	</nav>
</div>

<style>
	.SidebarNav {
		position: relative;
	}

	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(2, 6, 23, 0.7);
		opacity: 0;
		pointer-events: none;
		transition: opacity 150ms ease-in-out;
	}

	.overlay--visible {
		opacity: 1;
		pointer-events: auto;
	}

	.panel {
		position: fixed;
		inset: 0 auto 0 0;
		transform: translateX(-100%);
		display: flex;
		flex-direction: column;
		width: min(80vw, 20rem);
		max-width: 20rem;
		background: rgba(15, 23, 42, 0.95);
		border-right: 1px solid rgba(71, 85, 105, 0.4);
		box-shadow: 0 20px 45px rgba(15, 23, 42, 0.55);
		transition: transform 180ms ease-in-out;
		backdrop-filter: blur(16px);
		z-index: 40;
	}

	.panel--open {
		transform: translateX(0);
	}

	.panel__header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.75rem 1.5rem 1.25rem;
	}

	.panel__title {
		font-size: 0.95rem;
		line-height: 1.4;
		font-weight: 600;
		letter-spacing: 0.04em;
		color: #e2e8f0;
		text-transform: uppercase;
	}

	.panel__content {
		display: flex;
		flex-direction: column;
		padding: 0 1rem 2rem;
		overflow-y: auto;
	}

	.panel__section {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
		margin-top: 1.5rem;
	}

	.panel__section:first-child {
		margin-top: 0.5rem;
	}

	.panel__section--secondary {
		padding-top: 1.25rem;
		border-top: 1px solid rgba(71, 85, 105, 0.4);
	}

	.panel__section--profile {
		margin-top: 0.25rem;
		padding: 1.25rem 1rem 1rem;
		background: rgba(15, 23, 42, 0.7);
		border: 1px solid rgba(51, 65, 85, 0.6);
		border-radius: 1rem;
		gap: 1rem;
	}

	.panel__section-label {
		font-size: 0.65rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: rgba(148, 163, 184, 0.9);
		margin-bottom: 1rem;
	}

	.nav-link {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.55rem 0.75rem;
		border-radius: 0.5rem;
		font-size: 0.85rem;
		font-weight: 500;
		color: rgba(203, 213, 225, 0.9);
		transition: background-color 120ms ease-in-out, color 120ms ease-in-out, transform 120ms;
		text-decoration: none;
	}

	.nav-link:hover,
	.nav-link:focus-visible {
		color: #f8fafc;
		background-color: rgba(51, 65, 85, 0.35);
		transform: translateX(2px);
		outline: none;
	}

	.nav-link--active {
		color: #38bdf8;
		background: rgba(15, 118, 110, 0.18);
	}

	.profile-card {
		display: flex;
		align-items: center;
		gap: 0.85rem;
	}

	.profile-card__avatar {
		width: 3rem;
		height: 3rem;
		border-radius: 999px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		background: radial-gradient(circle at 35% 30%, rgba(94, 234, 212, 0.35), rgba(59, 130, 246, 0.25));
		color: #f8fafc;
		font-weight: 600;
		font-size: 1rem;
		letter-spacing: 0.02em;
		text-transform: uppercase;
		border: 1px solid rgba(148, 163, 184, 0.35);
		box-shadow: inset 0 0 0 1px rgba(15, 23, 42, 0.4), 0 8px 18px rgba(2, 6, 23, 0.6);
	}

	.profile-card__meta {
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
	}

	.profile-card__name {
		font-size: 0.95rem;
		font-weight: 600;
		color: #f8fafc;
	}

	.profile-card__email {
		font-size: 0.8rem;
		color: rgba(148, 163, 184, 0.85);
	}

	.profile-card__logout {
		margin-top: 0.25rem;
		padding: 0.45rem 0.75rem;
		border-radius: 0.6rem;
		border: 1px solid rgba(239, 68, 68, 0.4);
		background: rgba(239, 68, 68, 0.08);
		color: #fecaca;
		font-size: 0.78rem;
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		transition: background 120ms ease, border-color 120ms ease, color 120ms ease, transform 120ms ease;
	}

	.profile-card__logout:hover,
	.profile-card__logout:focus-visible {
		background: rgba(239, 68, 68, 0.18);
		border-color: rgba(239, 68, 68, 0.6);
		color: #fee2e2;
		transform: translateY(-1px);
		outline: none;
	}

	@media (min-width: 640px) {
		.SidebarNav {
			position: static;
			display: flex;
			width: 20rem;
			max-width: 20rem;
			flex-shrink: 0;
		}

		.overlay {
			display: none;
		}

		.panel {
			position: sticky;
			top: 0;
			height: 100vh;
			transform: translateX(0);
			border-right: 1px solid rgba(71, 85, 105, 0.6);
		}

		.panel__header {
			padding-top: 2.25rem;
		}

		.panel__content {
			padding-bottom: 3rem;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.overlay,
		.panel,
		.nav-link {
			transition: none;
		}
	}
</style>

