<script lang="ts">
	import type { Idea } from '$lib/api/types';
	import { Handle, Position } from '@xyflow/svelte';

	export let data: {
		idea: Idea;
		onSelect?: (idea: Idea) => void;
	};
	export let selected = false;

	const idea = data.idea;

	const borderColor = idea.color ?? '#2563eb';
	const ownerLabel = idea.user_id ? `${idea.user_id.slice(0, 6)}…` : 'Unknown';
</script>

<div
	class="idea-node group relative flex min-w-[160px] max-w-[280px] cursor-grab flex-col rounded-xl border-2 bg-slate-900/90 p-3 shadow-lg transition-shadow hover:shadow-primary/30"
	class:selected={selected}
	style={`--border-color:${borderColor}; border-color:${borderColor};`}
>
	<Handle class="handle handle--left" position={Position.Left} type="target" />
	<Handle class="handle handle--right" position={Position.Right} type="source" />
	<div class="idea-node__header flex items-center justify-between gap-2">
		<h3 class="text-sm font-semibold leading-tight text-slate-100">
			{idea.title}
		</h3>
		<span class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-xs text-slate-200">
			v{idea.version}
		</span>
	</div>
	{#if idea.description}
		<p class="mt-2 line-clamp-3 text-xs leading-relaxed text-slate-300">{idea.description}</p>
	{:else}
		<p class="mt-2 text-xs italic text-slate-500">Click to add details…</p>
	{/if}
	<div class="mt-3 flex items-center justify-between text-[10px] uppercase tracking-wide text-slate-500">
		<span>{new Date(idea.updated_at ?? idea.created_at).toLocaleDateString()}</span>
		<span>{ownerLabel}</span>
	</div>
</div>

<style>
	.idea-node.selected {
		box-shadow:
			0 0 0 2px color-mix(in srgb, var(--border-color) 60%, white 40%),
			0 10px 30px -12px rgba(15, 118, 255, 0.4);
	}

	.idea-node:hover {
		cursor: grab;
	}

	.idea-node :global(.handle) {
		width: 10px;
		height: 10px;
		border-radius: 9999px;
		background-color: rgba(255, 255, 255, 0.8);
		border: 1px solid rgba(15, 118, 255, 0.5);
	}

	.idea-node :global(.handle--left) {
		left: -5px;
	}

	.idea-node :global(.handle--right) {
		right: -5px;
	}
</style>

