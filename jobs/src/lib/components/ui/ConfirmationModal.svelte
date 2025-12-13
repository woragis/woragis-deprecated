<script lang="ts">
	import Modal from './Modal.svelte';
	import Button from './Button.svelte';

	let {
		open = $bindable(false),
		title = 'Confirm Action',
		message = 'Are you sure you want to proceed?',
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		variant = 'danger',
		onConfirm,
		onCancel
	}: {
		open?: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		variant?: 'danger' | 'warning' | 'primary';
		onConfirm?: () => void | Promise<void>;
		onCancel?: () => void;
	} = $props();

	let confirming = $state(false);

	async function handleConfirm() {
		confirming = true;
		try {
			if (onConfirm) {
				await onConfirm();
			}
			open = false;
		} catch (error) {
			console.error('Error in confirmation action:', error);
		} finally {
			confirming = false;
		}
	}

	function handleCancel() {
		if (onCancel) {
			onCancel();
		}
		open = false;
	}
</script>

<Modal bind:open size="sm" {title}>
	<div class="modal-content">
		<p class="message">{message}</p>
		<div class="actions">
			<Button onclick={handleCancel} variant="secondary" disabled={confirming}>
				{cancelText}
			</Button>
			<Button onclick={handleConfirm} variant={variant} disabled={confirming}>
				{confirming ? 'Processing...' : confirmText}
			</Button>
		</div>
	</div>
</Modal>

<style>
	.modal-content {
		padding: var(--spacing-lg);
	}

	.message {
		margin-bottom: var(--spacing-lg);
		color: var(--color-text-primary);
		font-size: var(--font-size-base);
		line-height: 1.6;
	}

	.actions {
		display: flex;
		gap: var(--spacing-sm);
		justify-content: flex-end;
	}
</style>
