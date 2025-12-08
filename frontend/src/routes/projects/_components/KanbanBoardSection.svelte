<script lang="ts">
	import type { KanbanBoard, KanbanCard, Milestone, UUID } from '$lib/api/types';
	import type { CardFormState, ColumnFormState } from '../[slug]/project-detail.logic';

	export let board: KanbanBoard | null = null;
	export let columnForm: ColumnFormState;
	export let cardForm: CardFormState;
	export let milestones: Milestone[] = [];

	export let onDeleteColumn: (columnId: UUID) => void;
	export let onMoveCard: (card: KanbanCard, columnId: UUID) => void;
	export let onDeleteCard: (card: KanbanCard) => void;
	export let onColumnFieldChange: <K extends keyof ColumnFormState>(
		field: K,
		value: ColumnFormState[K]
	) => void;
	export let onCardFieldChange: <K extends keyof CardFormState>(field: K, value: CardFormState[K]) => void;
	export let onAddColumn: () => void;
	export let onAddCard: () => void;
</script>

{#if board}
	<section class="space-y-4">
		<h3 class="text-sm font-semibold text-slate-100">Kanban Board</h3>
		<div class="overflow-x-auto">
			<div class="flex min-w-full gap-4">
				{#each board.columns as column (column.column.id)}
					<div class="w-72 flex-shrink-0 space-y-3 rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-200">
						<header class="flex items-center justify-between">
							<strong>{column.column.name}</strong>
							<button class="rounded bg-slate-800 px-2 py-1" on:click={() => onDeleteColumn(column.column.id)}>
								Delete
							</button>
						</header>
						<p class="text-[11px] text-slate-400">WIP limit: {column.column.wip_limit || '∞'}</p>
						<ul class="space-y-2">
							{#each column.cards as card (card.id)}
								<li class="rounded border border-slate-800 bg-slate-950/60 p-2">
									<p class="font-semibold text-slate-100">{card.title}</p>
									<p class="text-[11px] text-slate-400">{card.description}</p>
									<div class="mt-2 flex flex-wrap gap-2">
										<select
											class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
											on:change={(event) =>
												onMoveCard(card, (event.target as HTMLSelectElement).value as UUID)}
										>
											<option value="">Move to…</option>
											{#each board.columns as target (target.column.id)}
												{#if target.column.id !== column.column.id}
													<option value={target.column.id}>{target.column.name}</option>
												{/if}
											{/each}
										</select>
										<button
											class="rounded bg-rose-500/80 px-2 py-1 text-white"
											on:click={() => onDeleteCard(card)}
										>
											Delete
										</button>
									</div>
								</li>
							{/each}
						</ul>
					</div>
				{/each}
			</div>
		</div>

		<div class="grid gap-4 md:grid-cols-2">
			<div class="rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-300">
				<h4 class="text-sm font-semibold text-slate-100">Add Column</h4>
				<label class="mt-2 flex flex-col gap-1">
					Name
					<input
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
						value={columnForm.name}
						on:input={(event) => onColumnFieldChange('name', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<label class="flex flex-col gap-1">
					Position
					<input
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
						type="number"
						value={columnForm.position ?? ''}
						on:input={(event) =>
							onColumnFieldChange(
								'position',
								(event.target as HTMLInputElement).value
									? Number((event.target as HTMLInputElement).value)
									: undefined
							)}
					/>
				</label>
				<label class="flex flex-col gap-1">
					WIP Limit
					<input
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
						type="number"
						value={columnForm.wipLimit ?? ''}
						on:input={(event) =>
							onColumnFieldChange(
								'wipLimit',
								(event.target as HTMLInputElement).value
									? Number((event.target as HTMLInputElement).value)
									: undefined
							)}
					/>
				</label>
				<button
					class="mt-2 w-full rounded bg-emerald-500 px-2 py-2 text-xs font-semibold text-slate-900"
					on:click={onAddColumn}
				>
					Add Column
				</button>
			</div>

			<div class="rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-300">
				<h4 class="text-sm font-semibold text-slate-100">Add Card</h4>
				<label class="mt-2 flex flex-col gap-1">
					Column
					<select
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
						value={cardForm.columnId}
						on:change={(event) =>
							onCardFieldChange('columnId', (event.target as HTMLSelectElement).value as UUID | '')}
					>
						<option value="">Choose column…</option>
						{#each board.columns as col (col.column.id)}
							<option value={col.column.id}>{col.column.name}</option>
						{/each}
					</select>
				</label>
				<label class="flex flex-col gap-1">
					Title
					<input
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
						value={cardForm.title}
						on:input={(event) => onCardFieldChange('title', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<label class="flex flex-col gap-1">
					Description
					<textarea
						class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
						value={cardForm.description}
						on:input={(event) =>
							onCardFieldChange('description', (event.target as HTMLTextAreaElement).value)}
					></textarea>
				</label>
				<label class="flex flex-col gap-1">
					Due Date
					<input
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
						type="date"
						value={cardForm.dueDate}
						on:input={(event) => onCardFieldChange('dueDate', (event.target as HTMLInputElement).value)}
					/>
				</label>
				<label class="flex flex-col gap-1">
					Milestone
					<select
						class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
						value={cardForm.milestoneId}
						on:change={(event) =>
							onCardFieldChange('milestoneId', (event.target as HTMLSelectElement).value as UUID | '')}
					>
						<option value="">None</option>
						{#each milestones as milestone (milestone.id)}
							<option value={milestone.id}>{milestone.title}</option>
						{/each}
					</select>
				</label>
				<button
					class="mt-2 w-full rounded bg-sky-500 px-2 py-2 text-xs font-semibold text-white"
					on:click={onAddCard}
				>
					Add Card
				</button>
			</div>
		</div>
	</section>
{/if}

