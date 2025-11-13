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
			<a
				class={linkClass({ match: $page.url.pathname === '/' })}
				aria-current={$page.url.pathname === '/' ? 'page' : undefined}
				href="/"
			>
				Home
			</a>
			<a
				class={linkClass({ match: $page.url.pathname.startsWith('/finances') })}
				aria-current={$page.url.pathname.startsWith('/finances') ? 'page' : undefined}
				href="/finances"
			>
				Finances
			</a>
			<a
				class={linkClass({ match: $page.url.pathname.startsWith('/projects') })}
				aria-current={$page.url.pathname.startsWith('/projects') ? 'page' : undefined}
				href="/projects"
			>
				Projects
			</a>
			<div
				class="mt-6 border-t border-slate-800/60 pt-4 text-xs tracking-wide text-slate-500 uppercase"
			>
				Auth
			</div>
			<a
				class={linkClass({ match: $page.url.pathname.startsWith('/auth/login') })}
				aria-current={$page.url.pathname.startsWith('/auth/login') ? 'page' : undefined}
				href="/auth/login"
			>
				Sign in
			</a>
			<a
				class={linkClass({ match: $page.url.pathname.startsWith('/auth/register') })}
				aria-current={$page.url.pathname.startsWith('/auth/register') ? 'page' : undefined}
				href="/auth/register"
			>
				Register
			</a>
		</nav>
	</aside>

	<div class="flex flex-1 flex-col">
		<header class="border-b border-slate-800 bg-slate-900/70 px-6 py-4 sm:hidden">
			<div class="flex flex-col gap-2">
				<h1 class="text-lg font-semibold tracking-wide">Woragis Console</h1>
				<nav class="flex items-center gap-3 text-sm font-medium">
					<a
						class={mobileLinkClass({ match: $page.url.pathname === '/' })}
						aria-current={$page.url.pathname === '/' ? 'page' : undefined}
						href="/"
					>
						Home
					</a>
					<a
						class={mobileLinkClass({ match: $page.url.pathname.startsWith('/finances') })}
						aria-current={$page.url.pathname.startsWith('/finances') ? 'page' : undefined}
						href="/finances"
					>
						Finances
					</a>
					<a
						class={mobileLinkClass({ match: $page.url.pathname.startsWith('/projects') })}
						aria-current={$page.url.pathname.startsWith('/projects') ? 'page' : undefined}
						href="/projects"
					>
						Projects
					</a>
					<a
						class={mobileLinkClass({ match: $page.url.pathname.startsWith('/auth/login') })}
						aria-current={$page.url.pathname.startsWith('/auth/login') ? 'page' : undefined}
						href="/auth/login"
					>
						Sign in
					</a>
					<a
						class={mobileLinkClass({ match: $page.url.pathname.startsWith('/auth/register') })}
						aria-current={$page.url.pathname.startsWith('/auth/register') ? 'page' : undefined}
						href="/auth/register"
					>
						Register
					</a>
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
