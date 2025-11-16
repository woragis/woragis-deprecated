<script lang="ts">
	import type { IdeaNode } from '$lib/api/types';
	import { Handle, Position } from '@xyflow/svelte';

	export let data: {
		node: IdeaNode;
	};
	export let selected = false;

	const node = data.node;

	const borderColor = node.color ?? '#2563eb';
</script>

<div
	class="idea-node group relative flex cursor-grab flex-col rounded-xl border-2 bg-slate-900/90 p-3 shadow-lg transition-shadow hover:shadow-primary/30"
	class:selected={selected}
	style={`--border-color:${borderColor}; border-color:${borderColor};`}
>
	<!-- 4-directional connection handles -->
	<Handle class="handle handle--top" position={Position.Top} type="source" id="top" />
	<Handle class="handle handle--bottom" position={Position.Bottom} type="source" id="bottom" />
	<Handle class="handle handle--left" position={Position.Left} type="source" id="left" />
	<Handle class="handle handle--right" position={Position.Right} type="source" id="right" />
	
	<Handle class="handle handle--top" position={Position.Top} type="target" id="top-target" />
	<Handle class="handle handle--bottom" position={Position.Bottom} type="target" id="bottom-target" />
	<Handle class="handle handle--left" position={Position.Left} type="target" id="left-target" />
	<Handle class="handle handle--right" position={Position.Right} type="target" id="right-target" />

	<div class="idea-node__header flex items-center justify-between gap-2">
		<h3 class="text-sm font-semibold leading-tight text-slate-100">
			{node.title}
		</h3>
		<span class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-xs text-slate-200">
			{node.type || 'default'}
		</span>
	</div>
	{#if node.description}
		<p class="mt-2 line-clamp-3 text-xs leading-relaxed text-slate-300">{node.description}</p>
	{:else}
		<p class="mt-2 text-xs italic text-slate-500">Click to add details…</p>
	{/if}
	<div class="mt-3 flex items-center justify-between text-[10px] uppercase tracking-wide text-slate-500">
		<span>{new Date(node.updated_at ?? node.created_at).toLocaleDateString()}</span>
		<span class="rounded px-1.5 py-0.5 bg-slate-800/50" style={`background-color: ${borderColor}20; color: ${borderColor}`}>
			v{node.version}
		</span>
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
		opacity: 0;
		transition: opacity 0.2s;
	}

	.idea-node:hover :global(.handle),
	.idea-node.selected :global(.handle) {
		opacity: 1;
	}

	.idea-node :global(.handle--top) {
		top: -5px;
		left: 50%;
		transform: translateX(-50%);
	}

	.idea-node :global(.handle--bottom) {
		bottom: -5px;
		left: 50%;
		transform: translateX(-50%);
	}

	.idea-node :global(.handle--left) {
		left: -5px;
		top: 50%;
		transform: translateY(-50%);
	}

	.idea-node :global(.handle--right) {
		right: -5px;
		top: 50%;
		transform: translateY(-50%);
	}
</style>
