<script lang="ts">
	import { onMount } from 'svelte';
	import { getContextPreview, type ContextPreview } from '$lib/api/chats';
	import { toastError } from '$lib/utils/toast';
	import { Info, ChevronDown, ChevronUp } from 'lucide-svelte';

	interface Props {
		conversationId: string;
	}

	let { conversationId }: Props = $props();

	let contextPreview: ContextPreview | null = $state(null);
	let loading = $state(false);
	let showContext = $state(false);
	let showOptions = $state(false);

	onMount(async () => {
		await loadContext();
	});

	async function loadContext() {
		if (!conversationId) return;

		loading = true;
		try {
			contextPreview = await getContextPreview(conversationId);
		} catch (err) {
			console.error('Error loading context preview:', err);
			toastError('Failed to load context preview');
		} finally {
			loading = false;
		}
	}

	function getEnabledOptionsCount(): number {
		if (!contextPreview) return 0;
		const opts = contextPreview.options;
		return [
			opts.includeJobApplication,
			opts.includeResume,
			opts.includeUserProfile,
			opts.includeProjects,
			opts.includeCaseStudies,
			opts.includeTechnicalWritings,
			opts.includePosts,
			opts.includeProblemSolutions,
			opts.includeSkills,
			opts.includeExperiences
		].filter(Boolean).length;
	}
</script>

<div class="context-viewer">
	<button
		class="context-toggle"
		onclick={() => (showContext = !showContext)}
		disabled={loading || !contextPreview}
	>
		<Info size={16} />
		<span>Chat Context</span>
		{#if contextPreview}
			<span class="badge">{getEnabledOptionsCount()} enabled</span>
		{/if}
		{#if showContext}
			<ChevronUp size={16} />
		{:else}
			<ChevronDown size={16} />
		{/if}
	</button>

	{#if showContext}
		<div class="context-content">
			{#if loading}
				<p>Loading context...</p>
			{:else if !contextPreview}
				<p>No context available</p>
			{:else}
				<div class="context-options">
					<button
						class="options-toggle"
						onclick={() => (showOptions = !showOptions)}
					>
						<span>Context Options</span>
						{#if showOptions}
							<ChevronUp size={14} />
						{:else}
							<ChevronDown size={14} />
						{/if}
					</button>

					{#if showOptions}
						<div class="options-list">
							<div class="option-item">
								<span class="option-label">Job Application:</span>
								<span class="option-status {contextPreview.options.includeJobApplication ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeJobApplication ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Resume:</span>
								<span class="option-status {contextPreview.options.includeResume ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeResume ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">User Profile:</span>
								<span class="option-status {contextPreview.options.includeUserProfile ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeUserProfile ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Projects:</span>
								<span class="option-status {contextPreview.options.includeProjects ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeProjects ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Case Studies:</span>
								<span class="option-status {contextPreview.options.includeCaseStudies ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeCaseStudies ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Technical Writings:</span>
								<span class="option-status {contextPreview.options.includeTechnicalWritings ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeTechnicalWritings ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Posts:</span>
								<span class="option-status {contextPreview.options.includePosts ? 'enabled' : 'disabled'}">
									{contextPreview.options.includePosts ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Problem Solutions:</span>
								<span class="option-status {contextPreview.options.includeProblemSolutions ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeProblemSolutions ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Skills:</span>
								<span class="option-status {contextPreview.options.includeSkills ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeSkills ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
							<div class="option-item">
								<span class="option-label">Experiences:</span>
								<span class="option-status {contextPreview.options.includeExperiences ? 'enabled' : 'disabled'}">
									{contextPreview.options.includeExperiences ? '✓ Enabled' : '✗ Disabled'}
								</span>
							</div>
						</div>
					{/if}
				</div>

				{#if contextPreview.message}
					<div class="context-message">
						<p>{contextPreview.message}</p>
					</div>
				{/if}

				{#if contextPreview.context}
					<div class="context-text">
						<h4>Context Preview:</h4>
						<pre>{contextPreview.context}</pre>
					</div>
				{:else}
					<p class="no-context">No context available for this conversation.</p>
				{/if}
			{/if}
		</div>
	{/if}
</div>

<style>
	.context-viewer {
		border: 1px solid #e2e8f0;
		border-radius: 8px;
		background-color: #f8fafc;
		margin-bottom: 16px;
	}

	.context-toggle {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 16px;
		background: transparent;
		border: none;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		color: #334155;
		transition: background-color 0.2s;
	}

	.context-toggle:hover:not(:disabled) {
		background-color: #e2e8f0;
	}

	.context-toggle:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.badge {
		background-color: #3b82f6;
		color: white;
		padding: 2px 8px;
		border-radius: 12px;
		font-size: 12px;
		font-weight: 600;
	}

	.context-content {
		padding: 16px;
		border-top: 1px solid #e2e8f0;
		background-color: white;
	}

	.context-options {
		margin-bottom: 16px;
	}

	.options-toggle {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 8px 12px;
		background-color: #f1f5f9;
		border: 1px solid #cbd5e1;
		border-radius: 4px;
		cursor: pointer;
		font-size: 13px;
		font-weight: 500;
		color: #475569;
	}

	.options-toggle:hover {
		background-color: #e2e8f0;
	}

	.options-list {
		margin-top: 8px;
		padding: 12px;
		background-color: #f8fafc;
		border-radius: 4px;
	}

	.option-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 6px 0;
		font-size: 13px;
	}

	.option-label {
		color: #64748b;
		font-weight: 500;
	}

	.option-status {
		font-weight: 600;
		font-size: 12px;
	}

	.option-status.enabled {
		color: #10b981;
	}

	.option-status.disabled {
		color: #ef4444;
	}

	.context-message {
		padding: 12px;
		background-color: #fef3c7;
		border-left: 4px solid #f59e0b;
		border-radius: 4px;
		margin-bottom: 16px;
	}

	.context-message p {
		margin: 0;
		font-size: 13px;
		color: #92400e;
	}

	.context-text {
		margin-top: 16px;
	}

	.context-text h4 {
		margin: 0 0 8px 0;
		font-size: 14px;
		font-weight: 600;
		color: #334155;
	}

	.context-text pre {
		background-color: #1e293b;
		color: #e2e8f0;
		padding: 16px;
		border-radius: 4px;
		overflow-x: auto;
		font-size: 12px;
		line-height: 1.5;
		white-space: pre-wrap;
		word-wrap: break-word;
		max-height: 400px;
		overflow-y: auto;
	}

	.no-context {
		color: #64748b;
		font-size: 13px;
		font-style: italic;
		text-align: center;
		padding: 24px;
	}
</style>

