<script lang="ts">
	import { PUBLIC_API_BASE_URL } from '$env/static/public';
</script>

<section class="space-y-4">
	<h2 class="text-2xl font-semibold">Welcome to Woragis</h2>
	<p class="max-w-3xl text-sm leading-relaxed text-slate-300">
		This console stitches together the new Woragis backend services. Use the navigation above to
		manage finance transactions, automate reports, and orchestrate project workflows from a single
		interface. All API calls are proxied through the Axios client configured with your <code
			>PUBLIC_API_BASE_URL</code
		> environment variable.
	</p>
	<div class="rounded-lg border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-400">
		<strong class="block text-sm text-slate-200">Connected backend</strong>
		<span>{PUBLIC_API_BASE_URL ?? 'http://localhost:8080'} </span>
		<p class="mt-2">
			Define <code>PUBLIC_API_BASE_URL</code> and <code>PUBLIC_DEFAULT_USER_ID</code> in a local
			<code>.env</code>
			file (see
			<code>env.sample</code>) to point the client to your runtime environment.
		</p>
	</div>
</section>
