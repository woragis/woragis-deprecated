<script lang="ts">
	import type { ReportFilters } from '../reports.logic';

	export let filters: ReportFilters;
	export let onChange: (partial: Partial<ReportFilters>) => void;
	export let onApply: () => void;
</script>

<div class="flex flex-col gap-3 rounded-xl border border-slate-800/60 bg-slate-900/40 p-4">
	<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
		Search
		<input
			class="mt-1 w-full rounded-lg border border-slate-700/60 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
			placeholder="Search by name or description"
			value={filters.search}
			on:input={(event) => onChange({ search: (event.target as HTMLInputElement).value })}
			on:keydown={(event) => {
				if (event.key === 'Enter') onApply();
			}}
		/>
	</label>
	<div class="grid grid-cols-1 gap-2 text-sm text-slate-300">
		<label class="flex items-center justify-between gap-2 rounded-lg border border-slate-800/60 bg-slate-900/60 px-3 py-2">
			<span>Include archived</span>
			<input
				type="checkbox"
				checked={filters.includeArchived}
				on:change={(event) => onChange({ includeArchived: (event.target as HTMLInputElement).checked })}
			/>
		</label>
		<label class="flex items-center justify-between gap-2 rounded-lg border border-slate-800/60 bg-slate-900/60 px-3 py-2">
			<span>Favorites only</span>
			<input
				type="checkbox"
				checked={filters.favoritesOnly}
				on:change={(event) => onChange({ favoritesOnly: (event.target as HTMLInputElement).checked })}
			/>
		</label>
	</div>
	<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
		Channel
		<input
			class="mt-1 w-full rounded-lg border border-slate-700/60 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
			placeholder="Filter by delivery channel"
			value={filters.channel}
			on:input={(event) => onChange({ channel: (event.target as HTMLInputElement).value })}
			on:keydown={(event) => {
				if (event.key === 'Enter') onApply();
			}}
		/>
	</label>
	<button
		class="rounded-lg bg-slate-800/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-700/70"
		type="button"
		on:click={onApply}
	>
		Apply filters
	</button>
</div>

