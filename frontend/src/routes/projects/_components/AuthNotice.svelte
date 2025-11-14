<script lang="ts">
	export let title = 'Sign in to manage your projects';
	export let message =
		'You need to be authenticated before you can see or edit your project workspace.';
</script>

<div class="rounded border border-slate-800 bg-slate-900/60 p-6 text-center text-sm text-slate-300">
	<p class="mb-2 text-base font-semibold text-slate-100">{title}</p>
	<p>{message}</p>
</div>

