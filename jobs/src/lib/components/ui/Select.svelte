<script lang="ts">
	let {
		label = '',
		error = '',
		required = false,
		disabled = false,
		value = $bindable(''),
		name = ''
	}: {
		label?: string;
		error?: string;
		required?: boolean;
		disabled?: boolean;
		value?: string;
		name?: string;
	} = $props();
	
	let selectId = `select-${Math.random().toString(36).substr(2, 9)}`;
</script>

<div class="form-group">
	{#if label}
		<label for={selectId} class="form-label">
			{label} {#if required}<span class="required">*</span>{/if}
		</label>
	{/if}
	<select
		id={selectId}
		class="form-select"
		class:error={error}
		{disabled}
		{required}
		{name}
		bind:value
	>
		<slot />
	</select>
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
	
	.form-select {
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		background-color: var(--color-bg-primary);
		transition: border-color var(--transition-base);
		width: 100%;
		cursor: pointer;
	}
	
	.form-select:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.form-select:disabled {
		background-color: var(--color-bg-tertiary);
		cursor: not-allowed;
		opacity: 0.6;
	}
	
	.form-select.error {
		border-color: var(--color-danger);
	}
	
	.form-error {
		font-size: var(--font-size-xs);
		color: var(--color-danger);
	}
</style>
