<script lang="ts">
	import type { UUID } from '$lib/api/types';

	export let selectedCount = 0;
	export let queueMetadataText = '';
	export let onArchive: () => void;
	export let onRestore: () => void;
	export let onDelete: () => void;
	export let onQueueRuns: () => void;
	export let onMetadataChange: (value: string) => void;

	const metadataPlaceholder = '{"key":"value"}';
</script>

<div class="flex flex-col gap-2 rounded-xl border border-slate-800/70 bg-slate-900/60 p-4">
	<h4 class="text-sm font-semibold text-slate-100">Bulk actions</h4>
	<p class="text-xs text-slate-400">{selectedCount} selected</p>
	<div class="grid gap-2 text-sm text-slate-200">
		<button
			type="button"
			class="rounded-lg border border-slate-700/60 px-3 py-2 text-left transition hover:border-slate-500 hover:bg-slate-800/60"
			on:click={onArchive}
		>
			Archive selected
		</button>
		<button
			type="button"
			class="rounded-lg border border-slate-700/60 px-3 py-2 text-left transition hover:border-slate-500 hover:bg-slate-800/60"
			on:click={onRestore}
		>
			Restore selected
		</button>
		<button
			type="button"
			class="rounded-lg border border-red-500/60 px-3 py-2 text-left text-red-200 transition hover:border-red-400 hover:bg-red-500/10"
			on:click={onDelete}
		>
			Delete selected
		</button>
		<div class="mt-2 rounded-lg border border-slate-800/60 bg-slate-900/60 p-3">
			<label class="text-xs uppercase tracking-wide text-slate-400">
				Run metadata (JSON)
				<textarea
					class="mt-2 h-20 w-full rounded-lg border border-slate-700/60 bg-slate-950/80 px-3 py-2 text-xs text-slate-200 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					value={queueMetadataText}
					on:input={(event) => onMetadataChange((event.target as HTMLTextAreaElement).value)}
					placeholder={metadataPlaceholder}
				/>
			</label>
			<button
				type="button"
				class="mt-2 w-full rounded-lg bg-primary/80 px-3 py-2 text-sm font-medium text-white transition hover:bg-primary/70"
				on:click={onQueueRuns}
			>
				Queue runs
			</button>
		</div>
	</div>
</div>

