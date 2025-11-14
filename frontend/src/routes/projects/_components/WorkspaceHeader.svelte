<script lang="ts">
	export let isAuthenticated = false;
	export let onRefresh: () => void;
</script>

<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
	<h2 class="text-lg font-semibold text-slate-100">Projects Workspace</h2>
	<div class="mt-4 flex items-center justify-end gap-3">
		<button
			class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
			on:click={onRefresh}
			disabled={!isAuthenticated}
		>
			Refresh Projects
		</button>
	</div>
</div>

