<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import StageItem from '../_components/StageItem.svelte';
	import type { InterviewStage } from '$lib/api/jobapplicationstages';
	
	let {
		stages = [],
		onEdit,
		onComplete,
		onDelete
	}: {
		stages?: InterviewStage[];
		onEdit?: (stage: InterviewStage) => void;
		onComplete?: (stage: InterviewStage) => void;
		onDelete?: (id: string) => void;
	} = $props();
</script>

<Card title="Interview Stages ({stages.length})">
	{#if stages.length === 0}
		<EmptyState title="No interview stages yet" description="Add interview stages as you progress through the process" />
	{:else}
		<div class="list">
			{#each stages as stage}
				<StageItem {stage} {onEdit} {onComplete} {onDelete} />
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
