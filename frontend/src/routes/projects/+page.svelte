<script lang="ts">
	import { onMount } from 'svelte';
	import {
		addMilestone,
		bulkUpdateMilestones,
		createDependency,
		createKanbanCard,
		createKanbanColumn,
		createProject,
		deleteDependency,
		deleteKanbanCard,
		deleteKanbanColumn,
		duplicateProject,
		getKanbanBoard,
		listDependencies,
		listMilestones,
		listProjects,
		moveKanbanCard,
		updateProjectMetrics,
		updateProjectStatus
	} from '$lib/api/projects';
	import { authStore } from '$lib';
	import type {
		KanbanBoard,
		KanbanCard,
		Milestone,
		Project,
		ProjectDependency,
		ProjectStatus,
		UUID
	} from '$lib/api/types';

	const statusOptions: ProjectStatus[] = [
		'idea',
		'planning',
		'executing',
		'monitoring',
		'completed'
	];

	let isAuthenticated = false;
	let projects: Project[] = [];
	let activeProject: Project | null = null;
	let board: KanbanBoard | null = null;
	let milestones: Milestone[] = [];
	let dependencies: ProjectDependency[] = [];
	let loading = false;
	let error: string | null = null;

	const projectForm = {
		name: '',
		description: '',
		status: 'planning' as ProjectStatus,
		healthScore: 60
	};

	const columnForm = {
		name: '',
		position: undefined as number | undefined,
		wipLimit: undefined as number | undefined
	};
	const cardForm = {
		columnId: '' as UUID | '',
		title: '',
		description: '',
		dueDate: '',
		milestoneId: '' as UUID | '',
		position: undefined as number | undefined
	};

	const milestoneForm = {
		title: '',
		description: '',
		dueDate: new Date().toISOString().slice(0, 10)
	};
	const dependencyForm = {
		dependsOnProjectId: '' as UUID | '',
		type: 'relates' as 'blocks' | 'relates' | 'supports'
	};
	const duplicateForm = {
		name: '',
		description: '',
		status: 'planning' as ProjectStatus,
		copyBoard: true,
		copyMilestones: true,
		copyDependencies: false
	};

	onMount(() => {
		const unsubscribe = authStore.subscribe(async (state) => {
			isAuthenticated = state.isAuthenticated;

			if (!state.isAuthenticated) {
				projects = [];
				activeProject = null;
				board = null;
				milestones = [];
				dependencies = [];
				return;
			}

			await loadProjects();
		});

		return () => {
			unsubscribe();
		};
	});

	async function loadProjects() {
		if (!isAuthenticated) {
			return;
		}
		loading = true;
		try {
			projects = await listProjects();
			if (projects.length && !activeProject) {
				await selectProject(projects[0]);
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to load projects';
		} finally {
			loading = false;
		}
	}

	async function selectProject(project: Project) {
		activeProject = project;
		await Promise.all([
			loadBoard(project.id),
			loadMilestones(project.id),
			loadDependenciesList(project.id)
		]);
	}

	async function loadBoard(projectId: UUID) {
		if (!isAuthenticated) return;
		loading = true;
		try {
			board = await getKanbanBoard(projectId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to load kanban board';
		} finally {
			loading = false;
		}
	}

	async function loadMilestones(projectId: UUID) {
		if (!isAuthenticated) return;
		try {
			milestones = await listMilestones(projectId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to load milestones';
		}
	}

	async function loadDependenciesList(projectId: UUID) {
		if (!isAuthenticated) return;
		try {
			dependencies = await listDependencies(projectId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to load dependencies';
		}
	}

	async function handleCreateProject() {
		if (!isAuthenticated) {
			error = 'You must be signed in to create projects.';
			return;
		}
		if (!projectForm.name.trim()) {
			error = 'Project name is required';
			return;
		}
		loading = true;
		try {
			const created = await createProject({
				name: projectForm.name,
				description: projectForm.description,
				status: projectForm.status,
				healthScore: projectForm.healthScore
			});
			projectForm.name = '';
			projectForm.description = '';
			await loadProjects();
			await selectProject(created);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to create project';
		} finally {
			loading = false;
		}
	}

	async function handleCreateColumn() {
		if (!isAuthenticated || !activeProject || !columnForm.name.trim()) return;
		loading = true;
		try {
			board = await createKanbanColumn(activeProject.id, {
				name: columnForm.name,
				position: columnForm.position,
				wipLimit: columnForm.wipLimit
			});
			columnForm.name = '';
			columnForm.position = undefined;
			columnForm.wipLimit = undefined;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to create column';
		} finally {
			loading = false;
		}
	}

	async function handleCreateCard() {
		if (!isAuthenticated || !activeProject || !cardForm.columnId || !cardForm.title.trim()) return;
		loading = true;
		try {
			board = await createKanbanCard(activeProject.id, {
				columnId: cardForm.columnId,
				title: cardForm.title,
				description: cardForm.description,
				dueDate: cardForm.dueDate ? new Date(cardForm.dueDate).toISOString() : undefined,
				milestoneId: cardForm.milestoneId || undefined,
				position: cardForm.position
			});
			cardForm.title = '';
			cardForm.description = '';
			cardForm.dueDate = '';
			cardForm.milestoneId = '';
			cardForm.position = undefined;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to create card';
		} finally {
			loading = false;
		}
	}

	async function handleMoveCard(card: KanbanCard, targetColumnId: UUID) {
		if (!isAuthenticated || !activeProject) return;
		loading = true;
		try {
			board = await moveKanbanCard(activeProject.id, card.id, targetColumnId, 0);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to move card';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteCard(card: KanbanCard) {
		if (!isAuthenticated || !activeProject) return;
		loading = true;
		try {
			board = await deleteKanbanCard(activeProject.id, card.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to delete card';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteColumn(columnId: UUID) {
		if (!isAuthenticated || !activeProject) return;
		loading = true;
		try {
			board = await deleteKanbanColumn(activeProject.id, columnId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to delete column';
		} finally {
			loading = false;
		}
	}

	async function handleAddMilestone() {
		if (!isAuthenticated || !activeProject || !milestoneForm.title.trim()) return;
		loading = true;
		try {
			await addMilestone(activeProject.id, {
				title: milestoneForm.title,
				description: milestoneForm.description,
				dueDate: new Date(milestoneForm.dueDate).toISOString()
			});
			milestoneForm.title = '';
			milestoneForm.description = '';
			await loadMilestones(activeProject.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to add milestone';
		} finally {
			loading = false;
		}
	}

	async function handleToggleMilestone(milestone: Milestone) {
		if (!isAuthenticated) return;
		loading = true;
		try {
			await bulkUpdateMilestones(milestone.project_id, [
				{ milestoneId: milestone.id, completed: !milestone.completed }
			]);
			await loadMilestones(milestone.project_id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to update milestone';
		} finally {
			loading = false;
		}
	}

	async function handleCreateDependency() {
		if (!isAuthenticated || !activeProject || !dependencyForm.dependsOnProjectId) return;
		loading = true;
		try {
			await createDependency(activeProject.id, {
				dependsOnProjectId: dependencyForm.dependsOnProjectId,
				type: dependencyForm.type
			});
			dependencyForm.dependsOnProjectId = '';
			dependencyForm.type = 'relates';
			await loadDependenciesList(activeProject.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to create dependency';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteDependency(dependencyId: UUID) {
		if (!isAuthenticated || !activeProject) return;
		try {
			await deleteDependency(activeProject.id, dependencyId);
			await loadDependenciesList(activeProject.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to delete dependency';
		}
	}

	async function handleDuplicateProject(templateId: UUID) {
		if (!isAuthenticated) return;
		if (!templateId || !duplicateForm.name.trim()) {
			error = 'Provide a name for the duplicate project';
			return;
		}
		loading = true;
		try {
			const duplicate = await duplicateProject(templateId, {
				name: duplicateForm.name,
				description: duplicateForm.description,
				status: duplicateForm.status,
				copyBoard: duplicateForm.copyBoard,
				copyMilestones: duplicateForm.copyMilestones,
				copyDependencies: duplicateForm.copyDependencies
			});
			duplicateForm.name = '';
			duplicateForm.description = '';
			await loadProjects();
			await selectProject(duplicate);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unable to duplicate project';
		} finally {
			loading = false;
		}
	}

	function otherProjects(projectId?: UUID) {
		return projects.filter((project) => project.id !== projectId);
	}
</script>

<section class="space-y-6">
	<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
		<h2 class="text-lg font-semibold text-slate-100">Projects Workspace</h2>
		<div class="mt-4 flex items-center justify-end gap-3">
			<button
				class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white disabled:opacity-50"
				on:click={loadProjects}
				disabled={!isAuthenticated}
			>
				Refresh Projects
			</button>
		</div>
	</div>

	{#if error}
		<p class="rounded border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-200">
			{error}
		</p>
	{/if}

	{#if !isAuthenticated}
		<div
			class="rounded border border-slate-800 bg-slate-900/60 p-6 text-center text-sm text-slate-300"
		>
			<p class="mb-2 text-base font-semibold text-slate-100">Sign in to manage your projects</p>
			<p>You need to be authenticated before you can see or edit your project workspace.</p>
		</div>
	{:else}
		<section class="grid gap-6 lg:grid-cols-[1.05fr_2fr]">
			<div class="space-y-6">
				<div class="space-y-3 rounded border border-slate-800 bg-slate-900/60 p-4">
					<h3 class="text-sm font-semibold text-slate-100">Create Project</h3>
					<label class="flex flex-col gap-1 text-xs text-slate-300">
						Name
						<input
							class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							bind:value={projectForm.name}
						/>
					</label>
					<label class="flex flex-col gap-1 text-xs text-slate-300">
						Description
						<textarea
							class="min-h-[80px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							bind:value={projectForm.description}
						></textarea>
					</label>
					<label class="flex flex-col gap-1 text-xs text-slate-300">
						Status
						<select
							class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							bind:value={projectForm.status}
						>
							{#each statusOptions as status}
								<option value={status}>{status}</option>
							{/each}
						</select>
					</label>
					<label class="flex flex-col gap-1 text-xs text-slate-300">
						Health Score
						<input
							class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
							type="number"
							min="0"
							max="100"
							bind:value={projectForm.healthScore}
						/>
					</label>
					<button
						class="w-full rounded bg-emerald-500 px-3 py-2 text-xs font-semibold text-slate-900"
						on:click={handleCreateProject}
					>
						Create Project
					</button>
				</div>

				{#if activeProject}
					<div class="space-y-3 rounded border border-slate-800 bg-slate-900/60 p-4">
						<h3 class="text-sm font-semibold text-slate-100">Duplicate Active Project</h3>
						<label class="flex flex-col gap-1 text-xs text-slate-300">
							Name
							<input
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={duplicateForm.name}
								placeholder="Copy of ..."
							/>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-300">
							Description
							<textarea
								class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={duplicateForm.description}
							></textarea>
						</label>
						<label class="flex flex-col gap-1 text-xs text-slate-300">
							Status
							<select
								class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
								bind:value={duplicateForm.status}
							>
								{#each statusOptions as status}
									<option value={status}>{status}</option>
								{/each}
							</select>
						</label>
						<div class="flex flex-col gap-2 text-xs text-slate-300">
							<label class="flex items-center gap-2">
								<input type="checkbox" bind:checked={duplicateForm.copyBoard} /> Copy board
							</label>
							<label class="flex items-center gap-2">
								<input type="checkbox" bind:checked={duplicateForm.copyMilestones} /> Copy milestones
							</label>
							<label class="flex items-center gap-2">
								<input type="checkbox" bind:checked={duplicateForm.copyDependencies} /> Copy dependencies
							</label>
						</div>
						<button
							class="w-full rounded bg-sky-500 px-3 py-2 text-xs font-semibold text-white"
							on:click={() => handleDuplicateProject(activeProject?.id ?? '')}
						>
							Duplicate {activeProject?.name}
						</button>
					</div>
				{/if}
			</div>

			<div class="space-y-6">
				<div class="rounded border border-slate-800 bg-slate-900/60 p-4">
					<h3 class="text-sm font-semibold text-slate-100">Projects ({projects.length})</h3>
					<div class="mt-3 overflow-x-auto text-xs text-slate-200">
						<table class="min-w-full border-separate border-spacing-y-2">
							<thead class="text-[11px] tracking-wide text-slate-400 uppercase">
								<tr>
									<th class="text-left">Name</th>
									<th class="text-left">Status</th>
									<th class="text-left">Health</th>
									<th class="text-left">MRR</th>
									<th></th>
								</tr>
							</thead>
							<tbody>
								{#each projects as project}
									<tr class="rounded border border-slate-800 bg-slate-950/40">
										<td class="px-3 py-2 font-semibold">{project.name}</td>
										<td class="px-3 py-2 text-slate-300">{project.status}</td>
										<td class="px-3 py-2">{project.health_score}</td>
										<td class="px-3 py-2">{project.mrr.toFixed(2)}</td>
										<td class="px-3 py-2 text-right">
											<button
												class="rounded bg-slate-800 px-3 py-1 text-xs"
												on:click={() => selectProject(project)}
											>
												View
											</button>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>

				{#if activeProject}
					<div class="space-y-6">
						<section
							class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300"
						>
							<header class="flex flex-wrap items-center justify-between gap-2">
								<div>
									<h3 class="text-sm font-semibold text-slate-100">{activeProject.name}</h3>
									<p class="text-[11px] text-slate-400">Status: {activeProject.status}</p>
								</div>
								<div class="flex gap-2">
									<select
										bind:value={activeProject.status}
										class="rounded border border-slate-700 bg-slate-950 px-3 py-2 text-slate-100"
									>
										{#each statusOptions as status}
											<option value={status}>{status}</option>
										{/each}
									</select>
									<button
										class="rounded bg-indigo-500 px-3 py-2 text-xs font-semibold text-white"
										on:click={async () => {
											if (!activeProject) return;
											await updateProjectStatus(activeProject.id, activeProject.status);
											await loadProjects();
										}}
									>
										Save Status
									</button>
								</div>
							</header>
							<div class="mt-4 grid gap-3 md:grid-cols-4">
								<label class="flex flex-col gap-1"
									>MRR<input
										class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
										type="number"
										bind:value={activeProject.mrr}
									/></label
								>
								<label class="flex flex-col gap-1"
									>CAC<input
										class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
										type="number"
										bind:value={activeProject.cac}
									/></label
								>
								<label class="flex flex-col gap-1"
									>LTV<input
										class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
										type="number"
										bind:value={activeProject.ltv}
									/></label
								>
								<label class="flex flex-col gap-1"
									>Churn<input
										class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
										type="number"
										bind:value={activeProject.churn_rate}
									/></label
								>
							</div>
							<button
								class="mt-3 rounded bg-emerald-500 px-3 py-2 text-xs font-semibold text-slate-900"
								on:click={async () => {
									if (!activeProject) return;
									await updateProjectMetrics(activeProject.id, {
										healthScore: activeProject.health_score,
										mrr: activeProject.mrr,
										cac: activeProject.cac,
										ltv: activeProject.ltv,
										churnRate: activeProject.churn_rate
									});
									await loadProjects();
								}}
							>
								Save Metrics
							</button>
						</section>

						{#if board}
							<section class="space-y-4">
								<h3 class="text-sm font-semibold text-slate-100">Kanban Board</h3>
								<div class="overflow-x-auto">
									<div class="flex min-w-full gap-4">
										{#each board.columns as column}
											<div
												class="w-72 flex-shrink-0 space-y-3 rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-200"
											>
												<header class="flex items-center justify-between">
													<strong>{column.column.name}</strong>
													<button
														class="rounded bg-slate-800 px-2 py-1"
														on:click={() => handleDeleteColumn(column.column.id)}
													>
														Delete
													</button>
												</header>
												<p class="text-[11px] text-slate-400">
													WIP limit: {column.column.wip_limit || '∞'}
												</p>
												<ul class="space-y-2">
													{#each column.cards as card}
														<li class="rounded border border-slate-800 bg-slate-950/60 p-2">
															<p class="font-semibold text-slate-100">{card.title}</p>
															<p class="text-[11px] text-slate-400">{card.description}</p>
															<div class="mt-2 flex flex-wrap gap-2">
																<select
																	class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
																	on:change={(event) =>
																		handleMoveCard(
																			card,
																			(event.target as HTMLSelectElement).value as UUID
																		)}
																>
																	<option value="">Move to…</option>
																	{#each board.columns as target}
																		{#if target.column.id !== column.column.id}
																			<option value={target.column.id}>{target.column.name}</option>
																		{/if}
																	{/each}
																</select>
																<button
																	class="rounded bg-rose-500/80 px-2 py-1 text-white"
																	on:click={() => handleDeleteCard(card)}
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
									<div
										class="rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-300"
									>
										<h4 class="text-sm font-semibold text-slate-100">Add Column</h4>
										<label class="mt-2 flex flex-col gap-1">
											Name
											<input
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
												bind:value={columnForm.name}
											/>
										</label>
										<label class="flex flex-col gap-1">
											Position
											<input
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
												type="number"
												bind:value={columnForm.position}
											/>
										</label>
										<label class="flex flex-col gap-1">
											WIP Limit
											<input
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
												type="number"
												bind:value={columnForm.wipLimit}
											/>
										</label>
										<button
											class="mt-2 w-full rounded bg-emerald-500 px-2 py-2 text-xs font-semibold text-slate-900"
											on:click={handleCreateColumn}
										>
											Add Column
										</button>
									</div>

									<div
										class="rounded border border-slate-800 bg-slate-900/60 p-3 text-xs text-slate-300"
									>
										<h4 class="text-sm font-semibold text-slate-100">Add Card</h4>
										<label class="mt-2 flex flex-col gap-1">
											Column
											<select
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
												bind:value={cardForm.columnId}
											>
												<option value="">Choose column…</option>
												{#each board.columns as col}
													<option value={col.column.id}>{col.column.name}</option>
												{/each}
											</select>
										</label>
										<label class="flex flex-col gap-1">
											Title
											<input
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
												bind:value={cardForm.title}
											/>
										</label>
										<label class="flex flex-col gap-1">
											Description
											<textarea
												class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
												bind:value={cardForm.description}
											></textarea>
										</label>
										<label class="flex flex-col gap-1">
											Due Date
											<input
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
												type="date"
												bind:value={cardForm.dueDate}
											/>
										</label>
										<label class="flex flex-col gap-1">
											Milestone
											<select
												class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
												bind:value={cardForm.milestoneId}
											>
												<option value="">None</option>
												{#each milestones as milestone}
													<option value={milestone.id}>{milestone.title}</option>
												{/each}
											</select>
										</label>
										<button
											class="mt-2 w-full rounded bg-sky-500 px-2 py-2 text-xs font-semibold text-white"
											on:click={handleCreateCard}
										>
											Add Card
										</button>
									</div>
								</div>
							</section>
						{/if}

						<section class="grid gap-4 md:grid-cols-2">
							<div
								class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300"
							>
								<header class="flex items-center justify-between">
									<h4 class="text-sm font-semibold text-slate-100">
										Milestones ({milestones.length})
									</h4>
								</header>
								<ul class="mt-3 space-y-2">
									{#each milestones as milestone}
										<li class="rounded border border-slate-800 bg-slate-950/60 p-3">
											<div class="flex items-center justify-between">
												<div>
													<p class="font-semibold text-slate-100">{milestone.title}</p>
													<p class="text-[11px] text-slate-400">
														{new Date(milestone.due_date).toLocaleDateString()}
													</p>
												</div>
												<button
													class="rounded bg-slate-800 px-2 py-1"
													on:click={() => handleToggleMilestone(milestone)}
												>
													{milestone.completed ? 'Mark pending' : 'Mark done'}
												</button>
											</div>
											<p class="mt-2 text-[11px] text-slate-300">{milestone.description}</p>
										</li>
									{/each}
								</ul>
								<form class="mt-4 space-y-2" on:submit|preventDefault={handleAddMilestone}>
									<label class="flex flex-col">
										Title
										<input
											class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
											bind:value={milestoneForm.title}
										/>
									</label>
									<label class="flex flex-col">
										Description
										<textarea
											class="min-h-[60px] rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
											bind:value={milestoneForm.description}
										></textarea>
									</label>
									<label class="flex flex-col">
										Due Date
										<input
											class="rounded border border-slate-700 bg-slate-950 px-2 py-1"
											type="date"
											bind:value={milestoneForm.dueDate}
										/>
									</label>
									<button
										class="w-full rounded bg-emerald-500 px-2 py-2 text-xs font-semibold text-slate-900"
										>Add Milestone</button
									>
								</form>
							</div>

							<div
								class="rounded border border-slate-800 bg-slate-900/60 p-4 text-xs text-slate-300"
							>
								<header class="flex items-center justify-between">
									<h4 class="text-sm font-semibold text-slate-100">
										Dependencies ({dependencies.length})
									</h4>
								</header>
								<ul class="mt-3 space-y-2">
									{#each dependencies as dependency}
										<li
											class="flex items-center justify-between rounded border border-slate-800 bg-slate-950/60 px-3 py-2"
										>
											<span>{dependency.type} → {dependency.depends_on_project_id}</span>
											<button
												class="rounded bg-rose-500/70 px-2 py-1 text-white"
												on:click={() => handleDeleteDependency(dependency.id)}
											>
												Remove
											</button>
										</li>
									{/each}
								</ul>
								<form class="mt-4 space-y-2" on:submit|preventDefault={handleCreateDependency}>
									<label class="flex flex-col">
										Depends On
										<select
											class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
											bind:value={dependencyForm.dependsOnProjectId}
										>
											<option value="">Select project</option>
											{#each otherProjects(activeProject?.id) as project}
												<option value={project.id}>{project.name}</option>
											{/each}
										</select>
									</label>
									<label class="flex flex-col">
										Type
										<select
											class="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-slate-100"
											bind:value={dependencyForm.type}
										>
											<option value="blocks">blocks</option>
											<option value="relates">relates</option>
											<option value="supports">supports</option>
										</select>
									</label>
									<button
										class="w-full rounded bg-indigo-500 px-2 py-2 text-xs font-semibold text-white"
										>Add Dependency</button
									>
								</form>
							</div>
						</section>
					</div>
				{/if}
			</div>
		</section>
	{/if}
</section>
