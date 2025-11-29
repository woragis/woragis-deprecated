<script lang="ts">
	import { Wallet, Lightbulb } from 'lucide-svelte';
	
	const navItems = [
		{ href: '/personal/finances', label: 'Finances', icon: Wallet },
		{ href: '/personal/ideas', label: 'Ideas', icon: Lightbulb }
	];
</script>

<div class="personal-layout">
	<nav class="personal-nav">
		<div class="nav-header">
			<h2 class="nav-title">Personal</h2>
			<p class="nav-description">Manage personal data</p>
		</div>
		<ul class="nav-list">
			{#each navItems as item}
				<li>
					<a href={item.href} class="nav-link">
						<svelte:component this={item.icon} class="nav-icon" />
						<span>{item.label}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>
	<main class="personal-content">
		<slot />
	</main>
</div>

<style>
	.personal-layout {
		display: flex;
		gap: 2rem;
		min-height: calc(100vh - 4rem);
	}

	.personal-nav {
		width: 200px;
		flex-shrink: 0;
		background: rgba(15, 23, 42, 0.6);
		border: 1px solid rgba(71, 85, 105, 0.4);
		border-radius: 0.75rem;
		padding: 1.5rem;
		height: fit-content;
		position: sticky;
		top: 2rem;
	}

	.nav-header {
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid rgba(71, 85, 105, 0.3);
	}

	.nav-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: #f8fafc;
		margin-bottom: 0.25rem;
	}

	.nav-description {
		font-size: 0.875rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.nav-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.nav-link {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		color: rgba(203, 213, 225, 0.9);
		text-decoration: none;
		transition: all 150ms ease;
		font-size: 0.875rem;
	}

	.nav-link:hover {
		background: rgba(51, 65, 85, 0.3);
		color: #f8fafc;
	}

	.nav-link:global(.active) {
		background: rgba(59, 130, 246, 0.15);
		color: #93c5fd;
		border-left: 3px solid #3b82f6;
	}

	.nav-icon {
		width: 1.125rem;
		height: 1.125rem;
	}

	.personal-content {
		flex: 1;
		min-width: 0;
	}
</style>

