<script lang="ts">
	export let search = '';
	export let includeArchived = false;
	export let isAuthenticated = false;
	export let onSearchChange: (value: string) => void;
	export let onIncludeArchivedChange: (value: boolean) => void;
	export let onRefresh: () => void;
</script>

<div class="flex flex-col gap-4 rounded-lg border border-slate-800 bg-slate-900/60 p-4">
	<h2 class="text-lg font-semibold text-slate-100">Finance Workspace</h2>
	<div class="grid gap-3 md:grid-cols-2">
		<label class="flex flex-col gap-1 text-xs font-medium tracking-wide text-slate-400 uppercase">
			Search
			<input
				class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
				type="search"
				value={search}
				on:input={(event) => onSearchChange((event.target as HTMLInputElement).value)}
				placeholder="category, description, tags"
			/>
		</label>
	</div>
	<div class="flex flex-wrap items-center gap-3 text-xs">
		<label class="flex items-center gap-2">
			<input
				type="checkbox"
				checked={includeArchived}
				on:change={(event) => onIncludeArchivedChange((event.target as HTMLInputElement).checked)}
			/>
			<span>Include archived</span>
		</label>
		<button
			class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
			on:click={onRefresh}
			disabled={!isAuthenticated}
		>
			Refresh
		</button>
	</div>
</div>

