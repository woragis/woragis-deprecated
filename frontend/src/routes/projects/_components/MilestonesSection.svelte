<script lang="ts">
	import type { Milestone } from '$lib/api/types';
	import type { MilestoneFormState } from '../[slug]/project-detail.logic';

	export let milestones: Milestone[] = [];
	export let formState: MilestoneFormState;
	export let onFieldChange: <K extends keyof MilestoneFormState>(
		field: K,
		value: MilestoneFormState[K]
	) => void;
	export let onAdd: () => void;
	export let onToggleMilestone: (milestone: Milestone) => void;
</script>

<div class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300">
	<header class="flex items-center justify-between">
		<h4 class="text-sm font-semibold text-slate-100">Milestones ({milestones.length})</h4>
	</header>
	<ul class="mt-3 space-y-2">
		{#each milestones as milestone (milestone.id)}
			<li class="rounded border border-slate-800 bg-slate-950/60 p-3">
				<div class="flex items-center justify-between">
					<div>
						<p class="font-semibold text-slate-100">{milestone.title}</p>
						<p class="text-[11px] text-slate-400">
							{new Date(milestone.due_date).toLocaleDateString()}
						</p>
					</div>
					<button class="rounded bg-slate-800 px-2 py-1" on:click={() => onToggleMilestone(milestone)}>
						{milestone.completed ? 'Mark pending' : 'Mark done'}
					</button>
				</div>
				<p class="mt-2 text-[11px] text-slate-300">{milestone.description}</p>
			</li>
		{/each}
	</ul>
	<form class="mt-4 space-y-2" on:submit|preventDefault={onAdd}>
		<label class="flex flex-col">
			Title
			<input
				class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
				value={formState.title}
				on:input={(event) => onFieldChange('title', (event.target as HTMLInputElement).value)}
			/>
		</label>
		<label class="flex flex-col">
			Description
			<textarea
				class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
				value={formState.description}
				on:input={(event) => onFieldChange('description', (event.target as HTMLTextAreaElement).value)}
			></textarea>
		</label>
		<label class="flex flex-col">
			Due Date
			<input
				class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
				type="date"
				value={formState.dueDate}
				on:input={(event) => onFieldChange('dueDate', (event.target as HTMLInputElement).value)}
			/>
		</label>
		<button class="w-full rounded bg-emerald-500 px-2 py-2 text-xs font-semibold text-slate-900">
			Add Milestone
		</button>
	</form>
</div>

