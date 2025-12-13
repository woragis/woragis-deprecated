<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import ConversationItem from '../_components/ConversationItem.svelte';
	import type { Conversation } from '$lib/api/chats';
	
	let {
		conversations = [],
		onOpen
	}: {
		conversations?: Conversation[];
		onOpen?: (id: string) => void;
	} = $props();
</script>

{#if conversations.length > 0}
	<div class="list">
		{#each conversations as conv}
			<ConversationItem conversation={conv} {onOpen} />
		{/each}
	</div>
{:else}
	<div class="empty-state">
		<p>No conversations yet. Start a chat to begin.</p>
	</div>
{/if}

<style>
	.list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.empty-state {
		padding: var(--spacing-lg);
		text-align: center;
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
	}
</style>
