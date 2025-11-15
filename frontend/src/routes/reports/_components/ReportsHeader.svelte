<script lang="ts">
	export let isLoading = false;
	export let onRefresh: () => void;
	export let onCreate: () => void;
</script>

<header class="flex flex-wrap items-center justify-between gap-4">
	<div>
		<h1 class="text-2xl font-semibold text-slate-100">Reports</h1>
		<p class="text-sm text-slate-400">
			Manage report definitions, schedules, deliveries, and run history.
		</p>
	</div>
	<div class="flex items-center gap-3">
		<button
			class="rounded-lg border border-slate-700/70 px-3 py-2 text-sm text-slate-200 transition hover:border-slate-500 hover:text-slate-100 disabled:opacity-70"
			type="button"
			on:click={onRefresh}
			disabled={isLoading}
		>
			{isLoading ? 'Refreshing…' : 'Refresh'}
		</button>
		<button
			class="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-primary/80"
			type="button"
			on:click={onCreate}
		>
			New Report
		</button>
	</div>
</header>

