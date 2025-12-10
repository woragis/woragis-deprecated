<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import type { Conversation } from '$lib/api/chats';
	
	let {
		conversation,
		onOpen
	}: {
		conversation: Conversation;
		onOpen?: (id: string) => void;
	} = $props();
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="list-item">
	<div class="list-item-header">
		<span class="conversation-title">{conversation.title}</span>
		<span class="date">{formatDate(conversation.updatedAt)}</span>
	</div>
	{#if conversation.description}
		<p class="list-item-text">{conversation.description}</p>
	{/if}
	<div class="list-item-actions">
		{#if onOpen}
			<Button variant="primary" size="sm" onclick={() => onOpen(conversation.id)}>Open Chat</Button>
		{/if}
	</div>
</div>

<style>
	.list-item {
		background-color: var(--color-bg-primary);
		padding: var(--spacing-md);
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
	}
	
	.list-item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-sm);
	}
	
	.conversation-title {
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}
	
	.date {
		font-size: var(--font-size-xs);
		color: var(--color-text-tertiary);
	}
	
	.list-item-text {
		margin: var(--spacing-xs) 0;
		font-size: var(--font-size-sm);
		color: var(--color-text-secondary);
		line-height: 1.5;
	}
	
	.list-item-actions {
		display: flex;
		gap: var(--spacing-sm);
		margin-top: var(--spacing-sm);
	}
</style>
