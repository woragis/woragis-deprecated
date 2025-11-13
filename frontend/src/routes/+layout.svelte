<script lang="ts">
	import { page } from '$app/stores';

	const linkClass = ({ match }: { match: boolean }) =>
		`flex items-center gap-3 rounded px-3 py-2 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 ${
			match
				? 'bg-slate-800/80 text-primary'
				: 'text-slate-300 hover:bg-slate-800/60 hover:text-slate-100'
		}`;

	const mobileLinkClass = ({ match }: { match: boolean }) =>
		`rounded px-2 py-1 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 ${
			match ? 'bg-slate-800/80 text-primary' : 'text-slate-300 hover:text-slate-100'
		}`;

	const navLinks = [
		{ href: '/', label: 'Home', match: (pathname: string) => pathname === '/' },
		{
			href: '/finances',
			label: 'Finances',
			match: (pathname: string) => pathname.startsWith('/finances')
		},
		{
			href: '/projects',
			label: 'Projects',
			match: (pathname: string) => pathname.startsWith('/projects')
		},
		{
			href: '/monitoring',
			label: 'Monitoring',
			match: (pathname: string) => pathname.startsWith('/monitoring')
		}
	];

	const authLinks = [
		{
			href: '/auth/login',
			label: 'Sign in',
			match: (pathname: string) => pathname.startsWith('/auth/login')
		},
		{
			href: '/auth/register',
			label: 'Register',
			match: (pathname: string) => pathname.startsWith('/auth/register')
		}
	];
</script>

<svelte:head>
	<title>Woragis Console</title>
</svelte:head>

<div class="flex min-h-screen bg-slate-950 text-slate-100">
	<aside
		class="hidden w-64 flex-shrink-0 flex-col border-r border-slate-800 bg-slate-900/70 px-6 py-8 sm:flex"
	>
		<h1 class="text-lg font-semibold tracking-wide">Woragis Console</h1>
		<nav class="mt-8 flex flex-col gap-1 text-sm font-medium">
			{#each navLinks as item (item.href)}
				<a
					class={linkClass({ match: item.match($page.url.pathname) })}
					aria-current={item.match($page.url.pathname) ? 'page' : undefined}
					href={item.href}
				>
					{item.label}
				</a>
			{/each}
			<div
				class="mt-6 border-t border-slate-800/60 pt-4 text-xs tracking-wide text-slate-500 uppercase"
			>
				Auth
			</div>
			{#each authLinks as item (item.href)}
				<a
					class={linkClass({ match: item.match($page.url.pathname) })}
					aria-current={item.match($page.url.pathname) ? 'page' : undefined}
					href={item.href}
				>
					{item.label}
				</a>
			{/each}
		</nav>
	</aside>

	<div class="flex flex-1 flex-col">
		<header class="border-b border-slate-800 bg-slate-900/70 px-6 py-4 sm:hidden">
			<div class="flex flex-col gap-2">
				<h1 class="text-lg font-semibold tracking-wide">Woragis Console</h1>
				<nav class="flex items-center gap-3 text-sm font-medium">
					{#each navLinks as item (item.href + '-mobile')}
						<a
							class={mobileLinkClass({ match: item.match($page.url.pathname) })}
							aria-current={item.match($page.url.pathname) ? 'page' : undefined}
							href={item.href}
						>
							{item.label}
						</a>
					{/each}
					{#each authLinks as item (item.href + '-mobile')}
						<a
							class={mobileLinkClass({ match: item.match($page.url.pathname) })}
							aria-current={item.match($page.url.pathname) ? 'page' : undefined}
							href={item.href}
						>
							{item.label}
						</a>
					{/each}
				</nav>
			</div>
		</header>

		<main class="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-6 py-8">
			<slot />
		</main>
	</div>
</div>

<style>
	a {
		color: inherit;
		text-decoration: none;
	}
</style>
