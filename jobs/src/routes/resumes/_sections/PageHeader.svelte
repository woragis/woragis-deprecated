<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import { Upload, FileText } from 'lucide-svelte';
	
	let {
		onUploadClick,
		onGenerateClick
	}: {
		onUploadClick?: () => void;
		onGenerateClick?: () => void;
	} = $props();
</script>

<div class="page-header">
	<div class="header-content">
		<h1 class="page-title">Curriculum Vitae</h1>
		<p class="page-subtitle">Manage and organize your CVs and resumes</p>
	</div>
	<div class="header-actions">
		{#if onGenerateClick}
			<Button onclick={onGenerateClick} variant="secondary">
				<FileText size={16} />
				Generate CV
			</Button>
		{/if}
		{#if onUploadClick}
			<Button onclick={onUploadClick} variant="primary">
				<Upload size={16} />
				Upload CV
			</Button>
		{/if}
	</div>
</div>

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}
	
	.header-content {
		flex: 1;
	}
	
	.page-title {
		margin: 0 0 var(--spacing-xs) 0;
		font-size: var(--font-size-2xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}
	
	.page-subtitle {
		margin: 0;
		color: var(--color-text-secondary);
		font-size: var(--font-size-sm);
	}

	.header-actions {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}
</style>
