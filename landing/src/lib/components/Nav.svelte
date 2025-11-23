<script lang="ts">
	import LanguageSwitcher from './LanguageSwitcher.svelte';
</script>

<nav class="fixed top-0 left-0 right-0 z-50 bg-gray-900/80 backdrop-blur-md border-b border-gray-800">
	<div class="container mx-auto px-6 py-4">
		<div class="flex items-center justify-between">
			<a href="/" class="text-xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				Woragis
			</a>
			<div class="flex items-center gap-4">
				<LanguageSwitcher />
			</div>
		</div>
	</div>
</nav>

