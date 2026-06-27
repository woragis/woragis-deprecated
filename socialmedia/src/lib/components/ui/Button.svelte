<script lang="ts">
	let {
		variant = 'primary' as 'primary' | 'secondary' | 'danger' | 'success' | 'ghost',
		size = 'md' as 'sm' | 'md' | 'lg',
		disabled = false,
		type = 'button' as 'button' | 'submit' | 'reset',
		id,
		onclick
	}: {
		variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
		id?: string;
		onclick?: (event: MouseEvent) => void;
	} = $props();
	
	let className = $derived.by(() => {
		const base = 'btn';
		const variants = {
			primary: 'btn-primary',
			secondary: 'btn-secondary',
			danger: 'btn-danger',
			success: 'btn-success',
			ghost: 'btn-ghost'
		};
		const sizes = {
			sm: 'btn-sm',
			md: 'btn-md',
			lg: 'btn-lg'
		};
		return `${base} ${variants[variant]} ${sizes[size]} ${disabled ? 'btn-disabled' : ''}`;
	});
</script>

<button {type} {disabled} class={className} {id} {onclick}>
	<slot />
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		border: none;
		border-radius: var(--radius-md);
		font-weight: var(--font-weight-medium);
		cursor: pointer;
		transition: background-color var(--transition-base), color var(--transition-base);
		font-family: inherit;
	}
	
	.btn:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}
	
	.btn-disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	
	.btn-sm {
		padding: var(--spacing-xs) var(--spacing-sm);
		font-size: var(--font-size-sm);
	}
	
	.btn-md {
		padding: var(--spacing-sm) var(--spacing-md);
		font-size: var(--font-size-sm);
	}
	
	.btn-lg {
		padding: var(--spacing-md) var(--spacing-lg);
		font-size: var(--font-size-base);
	}
	
	.btn-primary {
		background-color: var(--color-primary);
		color: white;
	}
	
	.btn-primary:hover:not(.btn-disabled) {
		background-color: var(--color-primary-hover);
	}
	
	.btn-secondary {
		background-color: var(--color-bg-tertiary);
		color: var(--color-text-primary);
	}
	
	.btn-secondary:hover:not(.btn-disabled) {
		background-color: var(--color-bg-hover);
	}
	
	.btn-danger {
		background-color: var(--color-danger);
		color: white;
	}
	
	.btn-danger:hover:not(.btn-disabled) {
		background-color: var(--color-danger-hover);
	}
	
	.btn-success {
		background-color: var(--color-success);
		color: white;
	}
	
	.btn-success:hover:not(.btn-disabled) {
		background-color: var(--color-success-hover);
	}
	
	.btn-ghost {
		background-color: transparent;
		color: var(--color-text-secondary);
	}
	
	.btn-ghost:hover:not(.btn-disabled) {
		background-color: var(--color-bg-hover);
	}
</style>
