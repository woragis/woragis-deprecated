<script lang="ts">
	import { writable } from 'svelte/store';
	import type {
		DocumentationSection as DocSection,
		DocumentationSectionType,
		DocumentationVisibility,
		ProjectTechnology,
		ProjectFileStructure,
		ProjectArchitectureDiagram,
		TechnologyCategory,
		ArchitectureDiagramType,
		UUID
	} from '$lib/api/types';

	interface Props {
		projectId: UUID;
		documentation: {
			id: UUID;
			visibility: DocumentationVisibility;
			version: number;
		} | null;
		sections: DocSection[];
		technologies: ProjectTechnology[];
		fileStructures: ProjectFileStructure[];
		diagrams: ProjectArchitectureDiagram[];
		isLoading?: boolean;
		onCreateDocumentation?: () => Promise<void>;
		onUpdateVisibility?: (visibility: DocumentationVisibility) => Promise<void>;
		onCreateSection?: (payload: {
			type: DocumentationSectionType;
			title: string;
			content: string;
			position?: number;
		}) => Promise<void>;
		onUpdateSection?: (
			sectionId: UUID,
			payload: { title?: string; content?: string; position?: number }
		) => Promise<void>;
		onDeleteSection?: (sectionId: UUID) => Promise<void>;
		onReorderSections?: (sectionOrder: UUID[]) => Promise<void>;
		onCreateTechnology?: (payload: {
			name: string;
			version: string;
			category: TechnologyCategory;
			purpose: string;
			link?: string;
		}) => Promise<void>;
		onUpdateTechnology?: (
			techId: UUID,
			payload: {
				name?: string;
				version?: string;
				category?: TechnologyCategory;
				purpose?: string;
				link?: string;
			}
		) => Promise<void>;
		onDeleteTechnology?: (techId: UUID) => Promise<void>;
		onCreateFileStructure?: (payload: {
			path: string;
			name: string;
			is_directory: boolean;
			parent_id?: UUID;
			language?: string;
			line_count?: number;
			purpose?: string;
			position?: number;
		}) => Promise<void>;
		onUpdateFileStructure?: (
			fileStructureId: UUID,
			payload: {
				purpose?: string;
				line_count?: number;
				language?: string;
				position?: number;
			}
		) => Promise<void>;
		onDeleteFileStructure?: (fileStructureId: UUID) => Promise<void>;
		onCreateDiagram?: (payload: {
			type: ArchitectureDiagramType;
			title: string;
			description: string;
			content: string;
			format?: string;
			image_url?: string;
		}) => Promise<void>;
		onUpdateDiagram?: (
			diagramId: UUID,
			payload: {
				title?: string;
				description?: string;
				content?: string;
				image_url?: string;
			}
		) => Promise<void>;
		onDeleteDiagram?: (diagramId: UUID) => Promise<void>;
	}

	let {
		documentation,
		sections = [],
		technologies = [],
		fileStructures = [],
		diagrams = [],
		isLoading = false,
		onCreateDocumentation,
		onUpdateVisibility,
		onCreateSection,
		onUpdateSection,
		onDeleteSection,
		onReorderSections,
		onCreateTechnology,
		onUpdateTechnology,
		onDeleteTechnology,
		onCreateFileStructure,
		onUpdateFileStructure,
		onDeleteFileStructure,
		onCreateDiagram,
		onUpdateDiagram,
		onDeleteDiagram
	}: Props = $props();

	const activeTab = writable<'sections' | 'technologies' | 'file-structure' | 'diagrams'>('sections');
	const showCreateSectionForm = writable(false);
	const showCreateTechnologyForm = writable(false);
	const showCreateDiagramForm = writable(false);

	const sectionTypeOptions: { value: DocumentationSectionType; label: string }[] = [
		{ value: 'overview', label: 'Overview' },
		{ value: 'architecture', label: 'Architecture' },
		{ value: 'tech_stack', label: 'Tech Stack' },
		{ value: 'file_structure', label: 'File Structure' },
		{ value: 'api_documentation', label: 'API Documentation' },
		{ value: 'deployment', label: 'Deployment' },
		{ value: 'contributing', label: 'Contributing' },
		{ value: 'custom', label: 'Custom' }
	];

	const visibilityOptions: { value: DocumentationVisibility; label: string }[] = [
		{ value: 'public', label: 'Public' },
		{ value: 'authenticated', label: 'Authenticated Users' },
		{ value: 'collaborators', label: 'Collaborators Only' }
	];

	const technologyCategoryOptions: { value: TechnologyCategory; label: string }[] = [
		{ value: 'backend', label: 'Backend' },
		{ value: 'database', label: 'Database' },
		{ value: 'frontend', label: 'Frontend' },
		{ value: 'infrastructure', label: 'Infrastructure' },
		{ value: 'monitoring', label: 'Monitoring' },
		{ value: 'devops', label: 'DevOps' },
		{ value: 'testing', label: 'Testing' },
		{ value: 'other', label: 'Other' }
	];

	const diagramTypeOptions: { value: ArchitectureDiagramType; label: string }[] = [
		{ value: 'dependency', label: 'Dependency Graph' },
		{ value: 'component', label: 'Component Diagram' },
		{ value: 'data_flow', label: 'Data Flow' },
		{ value: 'infrastructure', label: 'Infrastructure' },
		{ value: 'custom', label: 'Custom' }
	];

	const newSectionForm = writable({
		type: 'overview' as DocumentationSectionType,
		title: '',
		content: ''
	});

	const newTechnologyForm = writable({
		name: '',
		version: '',
		category: 'backend' as TechnologyCategory,
		purpose: '',
		link: ''
	});

	const newDiagramForm = writable({
		type: 'component' as ArchitectureDiagramType,
		title: '',
		description: '',
		content: '',
		format: 'mermaid',
		image_url: ''
	});

	const handleCreateSection = async () => {
		const form = $newSectionForm;
		if (!form.title.trim() || !onCreateSection) return;

		await onCreateSection({
			type: form.type,
			title: form.title,
			content: form.content,
			position: sections.length
		});

		newSectionForm.set({
			type: 'overview',
			title: '',
			content: ''
		});
		showCreateSectionForm.set(false);
	};

	const handleCreateTechnology = async () => {
		const form = $newTechnologyForm;
		if (!form.name.trim() || !onCreateTechnology) return;

		await onCreateTechnology({
			name: form.name,
			version: form.version,
			category: form.category,
			purpose: form.purpose,
			link: form.link || undefined
		});

		newTechnologyForm.set({
			name: '',
			version: '',
			category: 'backend',
			purpose: '',
			link: ''
		});
		showCreateTechnologyForm.set(false);
	};

	const handleCreateDiagram = async () => {
		const form = $newDiagramForm;
		if (!form.title.trim() || !onCreateDiagram) return;

		await onCreateDiagram({
			type: form.type,
			title: form.title,
			description: form.description,
			content: form.content,
			format: form.format,
			image_url: form.image_url || undefined
		});

		newDiagramForm.set({
			type: 'component',
			title: '',
			description: '',
			content: '',
			format: 'mermaid',
			image_url: ''
		});
		showCreateDiagramForm.set(false);
	};
</script>

<div class="rounded-lg border border-slate-800 bg-slate-950/60">
	<div class="border-b border-slate-800 p-4">
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-semibold text-slate-100">Documentation</h2>
			{#if !documentation && onCreateDocumentation}
				<button
					on:click={async () => {
						if (onCreateDocumentation) await onCreateDocumentation();
					}}
					class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
				>
					Create Documentation
				</button>
			{:else if documentation && onUpdateVisibility}
				{@const visibility = documentation.visibility}
				<select
					value={visibility}
					on:change={async (e) => {
						const newVisibility = (e.target as HTMLSelectElement).value as DocumentationVisibility;
						if (onUpdateVisibility) {
							await onUpdateVisibility(newVisibility);
						}
					}}
					class="rounded border border-slate-700 bg-slate-900 px-3 py-1.5 text-sm text-slate-200"
				>
					{#each visibilityOptions as option}
						<option value={option.value} selected={visibility === option.value}>
							{option.label}
						</option>
					{/each}
				</select>
			{/if}
		</div>
		{#if documentation}
			<p class="mt-1 text-xs text-slate-400">Version {documentation.version}</p>
		{/if}
	</div>

	<div class="border-b border-slate-800">
		<div class="flex space-x-1 p-2">
			<button
				on:click={() => activeTab.set('sections')}
				class="rounded px-4 py-2 text-sm font-medium transition-colors {$activeTab === 'sections'
					? 'bg-indigo-600 text-white'
					: 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
			>
				Sections
			</button>
			<button
				on:click={() => activeTab.set('technologies')}
				class="rounded px-4 py-2 text-sm font-medium transition-colors {$activeTab === 'technologies'
					? 'bg-indigo-600 text-white'
					: 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
			>
				Tech Stack
			</button>
			<button
				on:click={() => activeTab.set('file-structure')}
				class="rounded px-4 py-2 text-sm font-medium transition-colors {$activeTab === 'file-structure'
					? 'bg-indigo-600 text-white'
					: 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
			>
				File Structure
			</button>
			<button
				on:click={() => activeTab.set('diagrams')}
				class="rounded px-4 py-2 text-sm font-medium transition-colors {$activeTab === 'diagrams'
					? 'bg-indigo-600 text-white'
					: 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
			>
				Diagrams
			</button>
		</div>
	</div>

	<div class="p-4">
		{#if !documentation}
			<div class="py-8 text-center text-sm text-slate-400">
				Documentation not initialized. Click "Create Documentation" to get started.
			</div>
		{:else if isLoading}
			<div class="py-8 text-center text-sm text-slate-400">Loading documentation...</div>
		{:else if $activeTab === 'sections'}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h3 class="text-base font-medium text-slate-200">Documentation Sections</h3>
					{#if onCreateSection}
						<button
							on:click={() => showCreateSectionForm.set(!$showCreateSectionForm)}
							class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
						>
							{$showCreateSectionForm ? 'Cancel' : 'Add Section'}
						</button>
					{/if}
				</div>

				{#if $showCreateSectionForm && onCreateSection}
					<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
						<div class="space-y-3">
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Type</label>
								<select
									bind:value={$newSectionForm.type}
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200"
								>
									{#each sectionTypeOptions as option}
										<option value={option.value}>{option.label}</option>
									{/each}
								</select>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Title</label>
								<input
									type="text"
									bind:value={$newSectionForm.title}
									placeholder="Section title"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Content (Markdown)</label>
								<textarea
									bind:value={$newSectionForm.content}
									placeholder="Section content in Markdown..."
									rows="6"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								></textarea>
							</div>
							<button
								on:click={handleCreateSection}
								class="w-full rounded bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
							>
								Create Section
							</button>
						</div>
					</div>
				{/if}

				{#if sections.length === 0}
					<div class="py-8 text-center text-sm text-slate-400">
						No documentation sections yet. Create one to get started.
					</div>
				{:else}
					<div class="space-y-3">
						{#each sections as section (section.id)}
							<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-2">
											<span class="text-xs font-medium text-indigo-400">
												{sectionTypeOptions.find((o) => o.value === section.type)?.label ?? section.type}
											</span>
											<h4 class="text-base font-medium text-slate-200">{section.title}</h4>
										</div>
										{#if section.content}
											<div class="mt-2 text-sm text-slate-400 line-clamp-3">
												{section.content.substring(0, 200)}
												{#if section.content.length > 200}...{/if}
											</div>
										{/if}
									</div>
									{#if onDeleteSection}
										<button
											on:click={async () => {
												if (onDeleteSection) await onDeleteSection(section.id);
											}}
											class="ml-2 text-red-400 hover:text-red-300"
										>
											Delete
										</button>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else if $activeTab === 'technologies'}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h3 class="text-base font-medium text-slate-200">Technology Stack</h3>
					{#if onCreateTechnology}
						<button
							on:click={() => showCreateTechnologyForm.set(!$showCreateTechnologyForm)}
							class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
						>
							{$showCreateTechnologyForm ? 'Cancel' : 'Add Technology'}
						</button>
					{/if}
				</div>

				{#if $showCreateTechnologyForm && onCreateTechnology}
					<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
						<div class="space-y-3">
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Name</label>
								<input
									type="text"
									bind:value={$newTechnologyForm.name}
									placeholder="e.g., Go, React, PostgreSQL"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								/>
							</div>
							<div class="grid grid-cols-2 gap-3">
								<div>
									<label class="mb-1 block text-sm font-medium text-slate-300">Version</label>
									<input
										type="text"
										bind:value={$newTechnologyForm.version}
										placeholder="1.0.0"
										class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
									/>
								</div>
								<div>
									<label class="mb-1 block text-sm font-medium text-slate-300">Category</label>
									<select
										bind:value={$newTechnologyForm.category}
										class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200"
									>
										{#each technologyCategoryOptions as option}
											<option value={option.value}>{option.label}</option>
										{/each}
									</select>
								</div>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Purpose</label>
								<textarea
									bind:value={$newTechnologyForm.purpose}
									placeholder="Why this technology was chosen..."
									rows="3"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								></textarea>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Link (optional)</label>
								<input
									type="url"
									bind:value={$newTechnologyForm.link}
									placeholder="https://..."
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								/>
							</div>
							<button
								on:click={handleCreateTechnology}
								class="w-full rounded bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
							>
								Add Technology
							</button>
						</div>
					</div>
				{/if}

				{#if technologies.length === 0}
					<div class="py-8 text-center text-sm text-slate-400">
						No technologies added yet. Add technologies to showcase your tech stack.
					</div>
				{:else}
					<div class="space-y-3">
						{#each technologies as tech (tech.id)}
							<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-2">
											<h4 class="text-base font-medium text-slate-200">{tech.name}</h4>
											{#if tech.version}
												<span class="text-xs text-slate-400">v{tech.version}</span>
											{/if}
											<span class="rounded bg-slate-800 px-2 py-0.5 text-xs text-slate-300">
												{technologyCategoryOptions.find((o) => o.value === tech.category)?.label ?? tech.category}
											</span>
										</div>
										{#if tech.purpose}
											<p class="mt-1 text-sm text-slate-400">{tech.purpose}</p>
										{/if}
										{#if tech.link}
											<a
												href={tech.link}
												target="_blank"
												rel="noopener noreferrer"
												class="mt-1 text-xs text-indigo-400 hover:text-indigo-300"
											>
												{tech.link}
											</a>
										{/if}
									</div>
									{#if onDeleteTechnology}
										<button
											on:click={async () => {
												if (onDeleteTechnology) await onDeleteTechnology(tech.id);
											}}
											class="ml-2 text-red-400 hover:text-red-300"
										>
											Delete
										</button>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else if $activeTab === 'file-structure'}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h3 class="text-base font-medium text-slate-200">File Structure</h3>
					{#if onCreateFileStructure}
						<button
							on:click={() => {
								/* TODO: Implement file structure creation UI */
							}}
							class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
						>
							Add File/Folder
						</button>
					{/if}
				</div>

				{#if fileStructures.length === 0}
					<div class="py-8 text-center text-sm text-slate-400">
						No file structure defined yet. Add files and folders to document your project structure.
					</div>
				{:else}
					<div class="space-y-1">
						{#each fileStructures as fs (fs.id)}
							<div class="flex items-center gap-2 rounded px-2 py-1 hover:bg-slate-800/50">
								<span class="text-sm text-slate-400">
									{fs.is_directory ? '📁' : '📄'}
								</span>
								<span class="text-sm text-slate-200">{fs.path}</span>
								{#if fs.language}
									<span class="text-xs text-slate-500">{fs.language}</span>
								{/if}
								{#if fs.line_count > 0}
									<span class="text-xs text-slate-500">{fs.line_count} lines</span>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else if $activeTab === 'diagrams'}
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<h3 class="text-base font-medium text-slate-200">Architecture Diagrams</h3>
					{#if onCreateDiagram}
						<button
							on:click={() => showCreateDiagramForm.set(!$showCreateDiagramForm)}
							class="rounded bg-indigo-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-indigo-700"
						>
							{$showCreateDiagramForm ? 'Cancel' : 'Add Diagram'}
						</button>
					{/if}
				</div>

				{#if $showCreateDiagramForm && onCreateDiagram}
					<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
						<div class="space-y-3">
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Type</label>
								<select
									bind:value={$newDiagramForm.type}
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200"
								>
									{#each diagramTypeOptions as option}
										<option value={option.value}>{option.label}</option>
									{/each}
								</select>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Title</label>
								<input
									type="text"
									bind:value={$newDiagramForm.title}
									placeholder="Diagram title"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Description</label>
								<input
									type="text"
									bind:value={$newDiagramForm.description}
									placeholder="Brief description"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-slate-300">Content (Mermaid/PlantUML)</label>
								<textarea
									bind:value={$newDiagramForm.content}
									placeholder="graph TD&#10;  A[Start] --> B[Process]"
									rows="8"
									class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 font-mono text-sm text-slate-200 placeholder:text-slate-500"
								></textarea>
							</div>
							<div class="grid grid-cols-2 gap-3">
								<div>
									<label class="mb-1 block text-sm font-medium text-slate-300">Format</label>
									<select
										bind:value={$newDiagramForm.format}
										class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200"
									>
										<option value="mermaid">Mermaid</option>
										<option value="plantuml">PlantUML</option>
										<option value="json">JSON</option>
									</select>
								</div>
								<div>
									<label class="mb-1 block text-sm font-medium text-slate-300">Image URL (optional)</label>
									<input
										type="url"
										bind:value={$newDiagramForm.image_url}
										placeholder="https://..."
										class="w-full rounded border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500"
									/>
								</div>
							</div>
							<button
								on:click={handleCreateDiagram}
								class="w-full rounded bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700"
							>
								Create Diagram
							</button>
						</div>
					</div>
				{/if}

				{#if diagrams.length === 0}
					<div class="py-8 text-center text-sm text-slate-400">
						No architecture diagrams yet. Create diagrams to visualize your project architecture.
					</div>
				{:else}
					<div class="space-y-3">
						{#each diagrams as diagram (diagram.id)}
							<div class="rounded border border-slate-700 bg-slate-900/50 p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="flex items-center gap-2">
											<span class="text-xs font-medium text-indigo-400">
												{diagramTypeOptions.find((o) => o.value === diagram.type)?.label ?? diagram.type}
											</span>
											<h4 class="text-base font-medium text-slate-200">{diagram.title}</h4>
										</div>
										{#if diagram.description}
											<p class="mt-1 text-sm text-slate-400">{diagram.description}</p>
										{/if}
										{#if diagram.content}
											<pre class="mt-2 overflow-x-auto rounded bg-slate-950 p-2 text-xs text-slate-300">
{diagram.content.substring(0, 200)}
{#if diagram.content.length > 200}...{/if}
</pre>
										{/if}
									</div>
									{#if onDeleteDiagram}
										<button
											on:click={async () => {
												if (onDeleteDiagram) await onDeleteDiagram(diagram.id);
											}}
											class="ml-2 text-red-400 hover:text-red-300"
										>
											Delete
										</button>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>

