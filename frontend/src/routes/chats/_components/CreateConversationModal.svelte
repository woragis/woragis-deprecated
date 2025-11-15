<script lang="ts">
	import type { Idea } from '$lib/api/types';
	import type { CreateConversationForm } from '../chats.logic';

	export let open = false;
	export let form: CreateConversationForm;
	export let createError = '';
	export let ideaQueryValue = '';
	export let projectQueryValue = '';
	export let ideaDropdownOpen = false;
	export let projectDropdownOpen = false;
	export let filteredProjects: Project[] = [];
	export let filteredIdeas: Idea[] = [];
	export let isSubmitting = false;
	export let onClose: () => void;
	export let onSubmit: (event: SubmitEvent) => void;
	export let onFieldChange: <K extends keyof CreateConversationForm>(
		field: K,
		value: CreateConversationForm[K]
	) => void;
	export let onIdeaQueryChange: (value: string) => void;
	export let onProjectQueryChange: (value: string) => void;
	export let onIdeaDropdownChange: (value: boolean) => void;
	export let onProjectDropdownChange: (value: boolean) => void;
	export let onSelectIdea: (idea: Idea) => void;
	export let onClearIdea: () => void;
	export let onSelectProject: (project: Project) => void;
	export let onClearProject: () => void;
</script>

{#if open}
	<div class="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/80 backdrop-blur">
		<form
			class="w-full max-w-lg space-y-4 rounded-2xl border border-slate-800/80 bg-slate-900/90 p-6 shadow-2xl"
			on:submit|preventDefault={onSubmit}
		>
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-lg font-semibold text-slate-100">Create conversation</h2>
					<p class="text-xs text-slate-400">
						Start a new chat thread and optionally link it to an idea or project.
					</p>
				</div>
				<button
					class="rounded-full border border-slate-700/70 px-3 py-1 text-xs text-slate-400 transition hover:border-slate-500 hover:text-slate-200"
					type="button"
					on:click={onClose}
				>
					Close
				</button>
			</div>

			<label class="flex flex-col gap-2 text-sm text-slate-200">
				<span class="text-xs uppercase tracking-wide text-slate-400">Title</span>
				<input
					class="rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					placeholder="Conversation subject"
					value={form.title}
					on:input={(event) => onFieldChange('title', (event.target as HTMLInputElement).value)}
					required
				/>
			</label>

			<label class="flex flex-col gap-2 text-sm text-slate-200">
				<span class="text-xs uppercase tracking-wide text-slate-400">Description</span>
				<textarea
					class="h-24 rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
					placeholder="Optional context"
					value={form.description}
					on:input={(event) => onFieldChange('description', (event.target as HTMLTextAreaElement).value)}
				></textarea>
			</label>

			<div class="grid gap-4 sm:grid-cols-2">
				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Idea</span>
					<div class="relative">
						<input
							class="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Search ideas…"
							value={ideaQueryValue}
							on:input={(event) => {
								onIdeaQueryChange((event.target as HTMLInputElement).value);
								onClearIdea();
								onIdeaDropdownChange(true);
							}}
							on:focus={() => onIdeaDropdownChange(true)}
							on:blur={() => setTimeout(() => onIdeaDropdownChange(false), 120)}
						/>
						{#if form.ideaId}
							<p class="mt-1 text-[11px] text-slate-400">Selected ID: {form.ideaId}</p>
						{/if}
						{#if ideaDropdownOpen && filteredIdeas.length > 0}
							<ul class="absolute z-10 mt-2 max-h-48 w-full overflow-y-auto rounded-lg border border-slate-800/80 bg-slate-900 shadow-lg">
								{#each filteredIdeas as idea (idea.id)}
									<li>
										<button
											type="button"
											class="flex w-full flex-col items-start gap-1 px-3 py-2 text-left text-xs text-slate-100 hover:bg-slate-800/60"
											on:click={() => {
												onSelectIdea(idea);
												onIdeaDropdownChange(false);
											}}
										>
											<span class="font-semibold">{idea.title}</span>
											<span class="text-[10px] text-slate-400">#{idea.id}</span>
										</button>
									</li>
								{/each}
							</ul>
						{/if}
					</div>
					{#if form.ideaId}
						<button
							class="text-left text-xs text-primary transition hover:text-primary/80"
							type="button"
							on:click={onClearIdea}
						>
							Clear selected idea
						</button>
					{/if}
				</label>

				<label class="flex flex-col gap-2 text-sm text-slate-200">
					<span class="text-xs uppercase tracking-wide text-slate-400">Project</span>
					<div class="relative">
						<input
							class="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Search projects…"
							value={projectQueryValue}
							on:input={(event) => {
								onProjectQueryChange((event.target as HTMLInputElement).value);
								onClearProject();
								onProjectDropdownChange(true);
							}}
							on:focus={() => onProjectDropdownChange(true)}
							on:blur={() => setTimeout(() => onProjectDropdownChange(false), 120)}
						/>
						{#if form.projectId}
							<p class="mt-1 text-[11px] text-slate-400">Selected ID: {form.projectId}</p>
						{/if}
						{#if projectDropdownOpen && filteredProjects.length > 0}
							<ul class="absolute z-10 mt-2 max-h-48 w-full overflow-y-auto rounded-lg border border-slate-800/80 bg-slate-900 shadow-lg">
								{#each filteredProjects as project (project.id)}
									<li>
										<button
											type="button"
											class="flex w-full flex-col items-start gap-1 px-3 py-2 text-left text-xs text-slate-100 hover:bg-slate-800/60"
											on:click={() => {
												onSelectProject(project);
												onProjectDropdownChange(false);
											}}
										>
											<span class="font-semibold">{project.name}</span>
											<span class="text-[10px] text-slate-400">#{project.id}</span>
										</button>
									</li>
								{/each}
							</ul>
						{/if}
					</div>
					{#if form.projectId}
						<button
							class="text-left text-xs text-primary transition hover:text-primary/80"
							type="button"
							on:click={onClearProject}
						>
							Clear selected project
						</button>
					{/if}
				</label>
			</div>

			{#if createError}
				<p class="text-xs text-red-400">{createError}</p>
			{/if}

			<div class="flex items-center justify-end gap-2">
				<button
					class="rounded-lg border border-slate-700 px-4 py-2 text-sm text-slate-300 transition hover:border-slate-500 hover:text-slate-100"
					type="button"
					on:click={onClose}
				>
					Cancel
				</button>
				<button
					class="rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white transition hover:bg-primary/90 disabled:opacity-60"
					type="submit"
					disabled={isSubmitting}
				>
					{isSubmitting ? 'Creating…' : 'Create conversation'}
				</button>
			</div>
		</form>
	</div>
{/if}

