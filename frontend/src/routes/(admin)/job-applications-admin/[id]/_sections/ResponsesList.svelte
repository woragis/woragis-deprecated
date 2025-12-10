<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import ResponseItem from '../_components/ResponseItem.svelte';
	import type { JobApplicationResponse } from '$lib/api/jobapplicationresponses';
	
	let {
		responses = [],
		onEdit,
		onDelete
	}: {
		responses?: JobApplicationResponse[];
		onEdit?: (response: JobApplicationResponse) => void;
		onDelete?: (id: string) => void;
	} = $props();
</script>

<Card title="Responses ({responses.length})">
	{#if responses.length === 0}
		<EmptyState title="No responses yet" description="Add a response when you hear back from the company" />
	{:else}
		<div class="list">
			{#each responses as response}
				<ResponseItem {response} {onEdit} {onDelete} />
			{/each}
		</div>
	{/if}
</Card>

<style>
	.list {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}
</style>
