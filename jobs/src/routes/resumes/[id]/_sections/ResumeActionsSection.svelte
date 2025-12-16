<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import { Star, StarOff, Award, AwardOff, RefreshCw, Trash2 } from 'lucide-svelte';
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resume,
		onMarkAsMain,
		onUnmarkAsMain,
		onMarkAsFeatured,
		onUnmarkAsFeatured,
		onRecalculateMetrics,
		onDelete
	}: {
		resume: Resume;
		onMarkAsMain?: () => void;
		onUnmarkAsMain?: () => void;
		onMarkAsFeatured?: () => void;
		onUnmarkAsFeatured?: () => void;
		onRecalculateMetrics?: () => void;
		onDelete?: () => void;
	} = $props();
</script>

<div class="section">
	<h2 class="section-title">Actions</h2>
	<div class="actions-list">
		{#if resume.isMain}
			<Button onclick={onUnmarkAsMain} variant="secondary" class="action-button">
				<StarOff size={16} />
				Unmark as Main
			</Button>
		{:else}
			<Button onclick={onMarkAsMain} variant="secondary" class="action-button">
				<Star size={16} />
				Mark as Main
			</Button>
		{/if}
		
		{#if resume.isFeatured}
			<Button onclick={onUnmarkAsFeatured} variant="secondary" class="action-button">
				<AwardOff size={16} />
				Unmark as Featured
			</Button>
		{:else}
			<Button onclick={onMarkAsFeatured} variant="secondary" class="action-button">
				<Award size={16} />
				Mark as Featured
			</Button>
		{/if}
		
		<Button onclick={onRecalculateMetrics} variant="secondary" class="action-button">
			<RefreshCw size={16} />
			Recalculate Metrics
		</Button>
		
		<Button onclick={onDelete} variant="danger" class="action-button">
			<Trash2 size={16} />
			Delete Resume
		</Button>
	</div>
</div>

<style>
	.actions-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.action-button {
		width: 100%;
		justify-content: flex-start;
	}
</style>
