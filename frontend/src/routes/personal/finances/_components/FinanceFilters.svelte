<script lang="ts">
	export let search = '';
	export let includeArchived = false;
	export let isAuthenticated = false;
	export let onSearchChange: (value: string) => void;
	export let onIncludeArchivedChange: (value: boolean) => void;
	export let onRefresh: () => void;
</script>

<div class="space-y-4">
	<h3 class="text-lg font-semibold text-slate-100">Filters</h3>
	<div class="grid gap-4 md:grid-cols-2">
		<label class="flex flex-col gap-1.5 text-sm text-slate-300">
			<span class="text-xs font-medium uppercase tracking-wider text-slate-400">Search</span>
			<input
				class="rounded-lg border border-slate-700 bg-slate-950/50 px-4 py-2.5 text-slate-100 transition-all focus:border-emerald-500/50 focus:bg-slate-950 focus:outline-none focus:ring-2 focus:ring-emerald-500/20"
				type="search"
				value={search}
				on:input={(event) => onSearchChange((event.target as HTMLInputElement).value)}
				placeholder="category, description, tags"
			/>
		</label>
	</div>
	<div class="flex flex-wrap items-center gap-4 text-sm">
		<label class="flex cursor-pointer items-center gap-2.5 text-slate-300">
			<input
				type="checkbox"
				class="h-4 w-4 rounded border-slate-700 bg-slate-950 text-emerald-600 focus:ring-2 focus:ring-emerald-500/20"
				checked={includeArchived}
				on:change={(event) => onIncludeArchivedChange((event.target as HTMLInputElement).checked)}
			/>
			<span>Include archived</span>
		</label>
		<button
			class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2.5 text-sm font-medium text-slate-200 backdrop-blur-sm transition-all hover:border-emerald-500/50 hover:bg-slate-800/80 disabled:opacity-50"
			on:click={onRefresh}
			disabled={!isAuthenticated}
		>
			<svg class="mr-2 inline h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
				></path>
			</svg>
			Refresh
		</button>
	</div>
</div>

