<script lang="ts">
	export let searchValue = '';
	export let includeArchived = false;
	export let bulkStatus = '';
	export let selectionCount = 0;
	export let onSearchInput: (value: string) => void;
	export let onIncludeArchivedChange: (value: boolean) => void;
	export let onSearch: () => void;
	export let onBulkArchive: () => void;
	export let onBulkRestore: () => void;
	export let onBulkDelete: () => void;
	export let onClearSelection: () => void;

	const handleKeydown = (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			event.preventDefault();
			onSearch();
		}
	};
</script>

<section class="rounded-2xl border border-slate-800/80 bg-slate-950/60 p-5 shadow-inner">
	<div class="flex flex-wrap items-center gap-3">
		<input
			class="w-full flex-1 rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 sm:w-64"
			placeholder="Search conversations…"
			value={searchValue}
			on:input={(event) => onSearchInput((event.target as HTMLInputElement).value)}
			on:keydown={handleKeydown}
		/>
		<label class="flex items-center gap-2 text-xs text-slate-300">
			<input
				type="checkbox"
				checked={includeArchived}
				on:change={(event) => onIncludeArchivedChange((event.target as HTMLInputElement).checked)}
			/>
			Include archived
		</label>
		<button
			class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/90"
			type="button"
			on:click={onSearch}
		>
			Search
		</button>
		{#if bulkStatus}
			<span class="text-xs text-slate-400">{bulkStatus}</span>
		{/if}
	</div>

	{#if selectionCount > 0}
		<div class="mt-4 flex flex-wrap items-center gap-2 text-xs text-slate-300">
			<span>{selectionCount} selected</span>
			<button class="rounded border border-slate-700 px-2 py-1" type="button" on:click={onBulkArchive}>
				Archive
			</button>
			<button class="rounded border border-slate-700 px-2 py-1" type="button" on:click={onBulkRestore}>
				Restore
			</button>
			<button
				class="rounded border border-red-500/40 px-2 py-1 text-red-200"
				type="button"
				on:click={onBulkDelete}
			>
				Delete
			</button>
			<button class="rounded border border-slate-700 px-2 py-1" type="button" on:click={onClearSelection}>
				Clear
			</button>
		</div>
	{/if}
</section>

