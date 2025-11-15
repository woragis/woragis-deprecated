<script lang="ts">
	import type { DefinitionFormState } from '../reports.logic';

	export let open = false;
	export let mode: 'create' | 'edit' = 'create';
	export let form: DefinitionFormState;
	export let onFieldChange: <K extends keyof DefinitionFormState>(
		field: K,
		value: DefinitionFormState[K]
	) => void;
	export let onClose: () => void;
	export let onSubmit: () => void;
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
		<div class="w-full max-w-xl rounded-2xl border border-slate-800 bg-slate-950 p-6 shadow-2xl">
			<header class="mb-4 flex items-center justify-between">
				<h2 class="text-lg font-semibold text-slate-100">
					{mode === 'create' ? 'New Report Definition' : 'Edit Report Definition'}
				</h2>
				<button
					class="rounded border border-slate-700 px-2 py-1 text-sm text-slate-300"
					on:click={onClose}
				>
					Close
				</button>
			</header>
			<div class="space-y-3">
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Name
					<input
						class="mt-1 w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						value={form.name}
						on:input={(event) => onFieldChange('name', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Description
					<textarea
						class="mt-1 min-h-[80px] w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
						value={form.description}
						on:input={(event) =>
							onFieldChange('description', (event.target as HTMLTextAreaElement).value)}
					></textarea>
				</label>
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Sections (JSON)
					<textarea
						class="mt-1 min-h-[120px] w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 font-mono"
						value={form.sectionsText}
						on:input={(event) =>
							onFieldChange('sectionsText', (event.target as HTMLTextAreaElement).value)}
					></textarea>
				</label>
				<label class="text-xs font-medium uppercase tracking-wide text-slate-400">
					Filters (JSON)
					<textarea
						class="mt-1 min-h-[120px] w-full rounded-lg border border-slate-700 bg-slate-900/80 px-3 py-2 text-xs text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 font-mono"
						value={form.filtersText}
						on:input={(event) =>
							onFieldChange('filtersText', (event.target as HTMLTextAreaElement).value)}
					></textarea>
				</label>
				<label class="flex items-center gap-2 text-sm text-slate-300">
					<input
						type="checkbox"
						checked={form.favorite}
						on:change={(event) => onFieldChange('favorite', (event.target as HTMLInputElement).checked)}
					/>
					Mark as favorite
				</label>
			</div>
			<div class="mt-4 flex justify-end gap-3">
				<button
					class="rounded-lg border border-slate-700 px-4 py-2 text-sm text-slate-200 hover:border-slate-500"
					type="button"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white hover:bg-primary/80"
					type="button"
					on:click={onSubmit}
				>
					{mode === 'create' ? 'Create' : 'Update'}
				</button>
			</div>
		</div>
	</div>
{/if}

