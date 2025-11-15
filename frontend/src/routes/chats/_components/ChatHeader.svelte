<script lang="ts">
	export let isRefreshing = false;
	export let onRefresh: () => void;
	export let onCreate: () => void;
</script>

<header class="flex flex-wrap items-center justify-between gap-4">
	<div>
		<h1 class="text-2xl font-semibold text-slate-100">Chats</h1>
		<p class="text-sm text-slate-400">
			Browse, search, and manage recent assistant conversations with live streaming updates.
		</p>
	</div>
	<div class="flex items-center gap-3">
		<button
			class="rounded-lg border border-primary/60 px-4 py-2 text-sm font-medium text-primary transition hover:bg-primary/10"
			type="button"
			on:click={onCreate}
		>
			New conversation
		</button>
		<button
			class="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100 disabled:opacity-60"
			type="button"
			on:click={onRefresh}
			disabled={isRefreshing}
		>
			{isRefreshing ? 'Refreshing…' : 'Refresh'}
		</button>
	</div>
</header>

