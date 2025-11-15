<script lang="ts">
	import type { ChatConversation, UUID } from '$lib/api/types';

	export let conversations: ChatConversation[] = [];
	export let isLoading = false;
	export let isFetching = false;
	export let selectedConversationId: UUID | null = null;
	export let selectedIds: Set<UUID> = new Set();
	export let onSelect: (conversation: ChatConversation) => void;
	export let onToggleSelection: (id: UUID) => void;
</script>

<section class="flex h-[650px] flex-col rounded-2xl border border-slate-800/80 bg-slate-950/60">
	<header class="flex items-center justify-between border-b border-slate-800/80 p-4 text-xs text-slate-400">
		<h2 class="font-semibold text-slate-200">Conversations</h2>
		{#if isFetching}
			<span>Refreshing…</span>
		{/if}
	</header>
	<div class="flex-1 overflow-y-auto">
		{#if isLoading}
			<div class="flex h-full items-center justify-center text-sm text-slate-400">
				Loading conversations…
			</div>
		{:else if conversations.length === 0}
			<div class="flex h-full items-center justify-center px-4 text-center text-sm text-slate-400">
				No conversations found. Adjust your filters or start a new chat from another module.
			</div>
		{:else}
			<ul class="divide-y divide-slate-800/80">
				{#each conversations as conversation (conversation.id)}
					<li>
						<button
							type="button"
							class={`flex w-full cursor-pointer items-center justify-between gap-3 px-4 py-3 text-left transition hover:bg-slate-900/80 ${
								selectedConversationId === conversation.id ? 'bg-slate-900/80' : ''
							}`}
							on:click={() => onSelect(conversation)}
						>
							<div class="flex items-center gap-2">
								<input
									type="checkbox"
									checked={selectedIds.has(conversation.id)}
									on:click|stopPropagation={() => onToggleSelection(conversation.id)}
								/>
								<div>
									<p class="text-sm font-semibold text-slate-100">{conversation.title}</p>
									<p class="text-xs text-slate-400 line-clamp-1">
										{conversation.description || 'No description'}
									</p>
								</div>
							</div>
							<div class="flex flex-col items-end gap-1 text-[10px] text-slate-400">
								<span>{new Date(conversation.updated_at).toLocaleString()}</span>
								{#if conversation.archived_at}
									<span class="rounded-full border border-amber-500/40 bg-amber-500/10 px-2 py-0.5 text-amber-200">
										Archived
									</span>
								{/if}
							</div>
						</button>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</section>

