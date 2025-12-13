<script lang="ts">
	let {
		message = '',
		type = 'success',
		duration = 3000,
		onClose
	}: {
		message?: string;
		type?: 'success' | 'error' | 'info' | 'warning';
		onClose?: () => void;
		duration?: number;
	} = $props();

	let visible = $state(true);

	$effect(() => {
		if (duration > 0) {
			const timer = setTimeout(() => {
				visible = false;
				setTimeout(() => onClose?.(), 300); // Wait for fade out
			}, duration);
			return () => clearTimeout(timer);
		}
	});

	function handleClose() {
		visible = false;
		setTimeout(() => onClose?.(), 300);
	}
</script>

{#if visible}
	<div class="toast toast-{type}" role="alert">
		<span class="toast-message">{message}</span>
		<button class="toast-close" onclick={handleClose} aria-label="Close">×</button>
	</div>
{/if}

<style>
	.toast {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--spacing-md);
		padding: var(--spacing-md);
		border-radius: var(--radius-md);
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		min-width: 300px;
		max-width: 500px;
		animation: slideIn 0.3s ease-out;
		opacity: 1;
		transition: opacity 0.3s ease-out;
	}

	.toast:not(.visible) {
		opacity: 0;
	}

	.toast-success {
		background-color: var(--color-success, #10b981);
		color: white;
	}

	.toast-error {
		background-color: var(--color-danger, #ef4444);
		color: white;
	}

	.toast-info {
		background-color: var(--color-primary, #3b82f6);
		color: white;
	}

	.toast-warning {
		background-color: var(--color-warning, #f59e0b);
		color: white;
	}

	.toast-message {
		flex: 1;
		font-size: var(--font-size-sm);
	}

	.toast-close {
		background: none;
		border: none;
		color: inherit;
		font-size: 1.5rem;
		cursor: pointer;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		opacity: 0.8;
		transition: opacity 0.2s;
	}

	.toast-close:hover {
		opacity: 1;
	}

	@keyframes slideIn {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}
</style>

