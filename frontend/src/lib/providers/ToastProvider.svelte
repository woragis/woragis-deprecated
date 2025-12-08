<script lang="ts">
	import { Toaster } from 'svelte-sonner';

	export let position: 'top-center' | 'top-right' | 'top-left' | 'bottom-center' | 'bottom-right' | 'bottom-left' =
		'top-center';
	export let theme: 'dark' | 'light' | 'system' = 'dark';
	export let duration = 4500;
</script>

<Toaster
	position={position}
	theme={theme}
	richColors
	closeButton
	toastOptions={{
		duration,
		style: 'background: rgba(15, 23, 42, 0.95); color: #f8fafc; border: 1px solid rgba(148, 163, 184, 0.35); box-shadow: 0 15px 35px rgba(2, 6, 23, 0.7); backdrop-filter: blur(12px);'
	}}
/>

