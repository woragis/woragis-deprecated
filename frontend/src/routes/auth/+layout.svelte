<script lang="ts">
	export let data;
</script>

<div class="flex min-h-screen items-center justify-center bg-slate-950 text-slate-100">
	<div class="w-full max-w-md rounded-lg border border-slate-800 bg-slate-900/70 p-8 shadow-lg">
		<div class="mb-6 text-center">
			<h1 class="text-xl font-semibold">Woragis Console</h1>
			<p class="mt-1 text-sm text-slate-400">Access your workspace</p>
		</div>
		<slot {data} />
	</div>
</div>
