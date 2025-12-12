<script lang="ts">
	import { language, languages, type Language } from '$lib/i18n';
	import { Globe, Check } from 'lucide-svelte';

	let showDropdown = $state(false);

	function setLanguage(lang: Language) {
		language.set(lang);
		showDropdown = false;
	}

	function toggleDropdown() {
		showDropdown = !showDropdown;
	}

	// Close dropdown when clicking outside
	import { onMount } from 'svelte';

	onMount(() => {
		function handleClickOutside(event: MouseEvent) {
			const target = event.target as HTMLElement;
			if (!target.closest('.language-switcher')) {
				showDropdown = false;
			}
		}

		window.addEventListener('click', handleClickOutside);
		
		return () => {
			window.removeEventListener('click', handleClickOutside);
		};
	});
</script>

<div class="language-switcher relative">
	<button
		onclick={toggleDropdown}
		class="flex items-center gap-2 px-3 py-2 bg-gray-800/50 hover:bg-gray-700/50 rounded-lg border border-gray-700 transition-colors"
		aria-label="Change language"
	>
		<Globe class="w-5 h-5 text-gray-300" />
		<span class="text-sm text-gray-300 hidden sm:inline">
			{#if $language === 'en'}🇺🇸
			{:else if $language === 'pt-BR'}🇧🇷
			{:else if $language === 'fr'}🇫🇷
			{:else if $language === 'es'}🇪🇸
			{:else if $language === 'de'}🇩🇪
			{:else if $language === 'ru'}🇷🇺
			{:else if $language === 'ja'}🇯🇵
			{:else if $language === 'ko'}🇰🇷
			{:else if $language === 'zh-CN'}🇨🇳
			{:else if $language === 'el'}🇬🇷
			{:else if $language === 'la'}🏛️
			{/if}
		</span>
	</button>

	{#if showDropdown}
		<div
			class="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-lg shadow-xl z-50 overflow-hidden"
		>
			{#each languages as lang}
				<button
					onclick={() => setLanguage(lang.code)}
					class="w-full flex items-center justify-between px-4 py-3 hover:bg-gray-700/50 transition-colors text-left {lang.code === $language
						? 'bg-gray-700/30'
						: ''}"
				>
					<div class="flex items-center gap-3">
						<span class="text-xl">{lang.flag}</span>
						<span class="text-white text-sm">{lang.name}</span>
					</div>
					{#if lang.code === $language}
						<Check class="w-4 h-4 text-blue-400" />
					{/if}
				</button>
			{/each}
		</div>
	{/if}
</div>

