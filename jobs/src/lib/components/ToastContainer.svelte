<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Toast from './ui/Toast.svelte';
	import { toastManager, type Toast as ToastType } from '$lib/utils/toast';

	let toasts: ToastType[] = $state([]);
	let unsubscribe: (() => void) | null = null;

	onMount(() => {
		toasts = toastManager.getToasts();
		unsubscribe = toastManager.subscribe((newToasts) => {
			toasts = newToasts;
		});
	});

	onDestroy(() => {
		if (unsubscribe) {
			unsubscribe();
		}
	});
</script>

<div class="toast-container">
	{#each toasts as toast (toast.id)}
		<Toast
			message={toast.message}
			type={toast.type}
			duration={toast.duration}
			onClose={() => toastManager.remove(toast.id)}
		/>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		top: var(--spacing-md);
		right: var(--spacing-md);
		z-index: 9999;
		display: flex;
		flex-direction: column;
		gap: var(--spacing-sm);
		pointer-events: none;
	}

	.toast-container :global(.toast) {
		pointer-events: auto;
	}
</style>

