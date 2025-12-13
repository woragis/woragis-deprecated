<script lang="ts">
	let {
		label = '',
		error = '',
		required = false,
		disabled = false,
		rows = 4,
		value = $bindable(''),
		placeholder = '',
		name = ''
	}: {
		label?: string;
		error?: string;
		required?: boolean;
		disabled?: boolean;
		rows?: number;
		value?: string;
		placeholder?: string;
		name?: string;
	} = $props();
	
	let textareaId = `textarea-${Math.random().toString(36).substr(2, 9)}`;
</script>

<div class="form-group">
	{#if label}
		<label for={textareaId} class="form-label">
			{label} {#if required}<span class="required">*</span>{/if}
		</label>
	{/if}
	<textarea
		id={textareaId}
		class="form-textarea"
		class:error={error}
		{disabled}
		{required}
		{rows}
		{placeholder}
		{name}
		bind:value
	></textarea>
	{#if error}
		<span class="form-error">{error}</span>
	{/if}
</div>

<style>
	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}
	
	.form-label {
		font-size: var(--font-size-sm);
		font-weight: var(--font-weight-medium);
		color: var(--color-text-primary);
	}
	
	.required {
		color: var(--color-danger);
	}
	
	.form-textarea {
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		background-color: var(--color-bg-primary);
		transition: border-color var(--transition-base);
		width: 100%;
		font-family: inherit;
		resize: vertical;
	}
	
	.form-textarea:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.form-textarea:disabled {
		background-color: var(--color-bg-tertiary);
		cursor: not-allowed;
		opacity: 0.6;
	}
	
	.form-textarea.error {
		border-color: var(--color-danger);
	}
	
	.form-error {
		font-size: var(--font-size-xs);
		color: var(--color-danger);
	}
</style>
