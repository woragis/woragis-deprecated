<script lang="ts">
	let {
		open = $bindable(false),
		size = 'md' as 'sm' | 'md' | 'lg',
		title = ''
	}: {
		open?: boolean;
		size?: 'sm' | 'md' | 'lg';
		title?: string;
	} = $props();
	
	let className = $derived.by(() => {
		const sizes = {
			sm: 'modal-sm',
			md: 'modal-md',
			lg: 'modal-lg'
		};
		return `modal ${sizes[size]}`;
	});
	
	function handleOverlayClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			open = false;
		}
	}
	
	function handleKeydown(e: KeyboardEvent) {
		if (open && e.key === 'Escape') {
			open = false;
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<div class="modal-overlay" onclick={handleOverlayClick} role="dialog" aria-modal="true" aria-labelledby={title ? 'modal-title' : undefined}>
		<div class={className}>
			{#if title}
				<h2 id="modal-title" class="modal-title">{title}</h2>
			{/if}
			<slot />
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		inset: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: var(--z-modal);
		padding: var(--spacing-md);
	}
	
	.modal {
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		width: 100%;
		max-height: 90vh;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		border: 1px solid var(--color-border);
	}
	
	.modal-sm {
		max-width: 400px;
	}
	
	.modal-md {
		max-width: 600px;
	}
	
	.modal-lg {
		max-width: 900px;
	}
	
	.modal-title {
		margin: 0 0 var(--spacing-md) 0;
		font-size: var(--font-size-xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
		padding: var(--spacing-lg);
		padding-bottom: 0;
	}
	
	:global(.modal-content) {
		padding: var(--spacing-lg);
	}
</style>
