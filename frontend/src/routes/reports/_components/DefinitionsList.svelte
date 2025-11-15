<script lang="ts">
	import type { ReportDefinition, UUID } from '$lib/api/types';

	export let definitions: ReportDefinition[] = [];
	export let isLoading = false;
	export let selectedDefinitionId: UUID | null = null;
	export let selectedIds: Set<UUID>;
	export let allSelected = false;

	export let onSelectDefinition: (definitionId: UUID) => void;
	export let onToggleSelection: (definitionId: UUID) => void;
	export let onToggleSelectAll: () => void;
	export let onToggleFavorite: (definitionId: UUID, isFavorite: boolean) => void;
</script>

<div class="flex flex-col gap-2">
	<div class="flex items-center justify-between text-xs uppercase tracking-wide text-slate-400">
		<span>Definitions</span>
		<button
			class="rounded-md border border-slate-800/60 px-2 py-1 text-[11px] text-slate-300 transition hover:border-slate-600 hover:text-slate-100"
			on:click={onToggleSelectAll}
			type="button"
		>
			{allSelected ? 'Clear all' : 'Select all'}
		</button>
	</div>
	<div class="space-y-2 overflow-auto rounded-xl border border-slate-800/60 bg-slate-900/50 p-2 max-h-[540px]">
		{#if isLoading}
			<div class="flex items-center justify-center py-10 text-sm text-slate-400">
				Loading definitions…
			</div>
		{:else if definitions.length === 0}
			<div class="rounded-lg border border-slate-800/60 bg-slate-900/70 px-3 py-6 text-center text-sm text-slate-400">
				No definitions found. Try adjusting your filters.
			</div>
		{:else}
			{#each definitions as def}
				<button
					type="button"
					class={`w-full rounded-lg border px-3 py-3 text-left transition ${
						def.id === selectedDefinitionId
							? 'border-primary/60 bg-primary/10 text-slate-50'
							: 'border-slate-800/60 bg-slate-900/60 text-slate-200 hover:border-primary/40 hover:bg-slate-800/60'
					}`}
					on:click={() => onSelectDefinition(def.id)}
				>
					<div class="flex items-start justify-between gap-2">
						<div class="flex items-center gap-2">
							<input
								type="checkbox"
								checked={selectedIds.has(def.id)}
								on:click={(event) => {
									event.stopPropagation();
									onToggleSelection(def.id);
								}}
							/>
							<h3 class="text-sm font-semibold">{def.name}</h3>
						</div>
						<span
							role="button"
							tabindex="0"
							aria-pressed={def.is_favorite}
							class={`cursor-pointer text-xs ${
								def.is_favorite ? 'text-amber-300' : 'text-slate-500 hover:text-slate-200'
							}`}
							on:click={(event) => {
								event.stopPropagation();
								onToggleFavorite(def.id, def.is_favorite);
							}}
							on:keydown={(event) => {
								if (event.key === 'Enter' || event.key === ' ') {
									event.preventDefault();
									event.stopPropagation();
									onToggleFavorite(def.id, def.is_favorite);
								}
							}}
						>
							{def.is_favorite ? '★' : '☆'}
						</span>
					</div>
					<p class="mt-1 line-clamp-2 text-xs text-slate-400">{def.description}</p>
					<div class="mt-2 flex items-center justify-between text-[11px] uppercase tracking-wide text-slate-500">
						<span>{new Date(def.updated_at).toLocaleDateString()}</span>
						{#if def.archived_at}
							<span class="rounded-full border border-red-500/40 px-2 py-[1px] text-[10px] text-red-200">
								Archived
							</span>
						{/if}
					</div>
				</button>
			{/each}
		{/if}
	</div>
</div>

