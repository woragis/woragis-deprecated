<script lang="ts">
	export let isLoading = false;
	export let onRefresh: () => void;
</script>

<div class="flex items-center justify-between">
	<div>
		<h1 class="text-2xl font-bold text-slate-100">Schedules</h1>
		<p class="mt-1 text-sm text-slate-400">
			Manage automated report schedules across all your reports
		</p>
	</div>
	<button
		class="rounded-lg border border-slate-700/60 px-4 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:bg-slate-800/60 disabled:opacity-50"
		type="button"
		on:click={onRefresh}
		disabled={isLoading}
	>
		{isLoading ? 'Refreshing...' : 'Refresh'}
	</button>
</div>

