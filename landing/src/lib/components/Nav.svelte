<script lang="ts">
	import LanguageSwitcher from './LanguageSwitcher.svelte';
	import { Menu, X } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let mobileMenuOpen = $state(false);
	let scrolled = $state(false);

	const navLinks = [
		{ href: '#about-me', label: 'About' },
		{ href: '#experience', label: 'Experience' },
		{ href: '#skills', label: 'Skills' },
		{ href: '#projects', label: 'Projects' },
		{ href: '#languages', label: 'Languages' },
		{ href: '#contact', label: 'Contact' }
	];

	function handleScroll() {
		scrolled = window.scrollY > 20;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	function handleNavClick(href: string) {
		closeMobileMenu();
		// Smooth scroll to section
		const element = document.querySelector(href);
		if (element) {
			element.scrollIntoView({ behavior: 'smooth', block: 'start' });
		}
	}

	onMount(() => {
		window.addEventListener('scroll', handleScroll);
		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	});
</script>

<nav
	class="fixed top-0 left-0 right-0 z-50 transition-all duration-300 
		 bg-gray-900/95 backdrop-blur-xl border-b border-gray-800/50 shadow-lg shadow-black/20"
>
	<div class="container mx-auto px-6 py-4">
		<div class="flex items-center justify-between">
			<!-- Logo -->
			<a
				href="/"
				class="text-xl md:text-2xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent hover:from-blue-300 hover:via-purple-300 hover:to-pink-300 transition-all duration-300 relative group"
			>
				<span class="relative z-10">Woragis</span>
				<span
					class="absolute inset-0 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent opacity-0 group-hover:opacity-100 blur-sm transition-opacity duration-300"
				>
					Woragis
				</span>
			</a>

			<!-- Desktop Navigation -->
			<div class="hidden md:flex items-center gap-1">
				{#each navLinks as link}
					<a
						href={link.href}
						onclick={(e) => {
							e.preventDefault();
							handleNavClick(link.href);
						}}
						class="px-4 py-2 text-sm font-medium text-gray-300 hover:text-white rounded-lg hover:bg-gray-800/50 transition-all duration-200 relative group"
					>
						{link.label}
						<span
							class="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300"
						></span>
					</a>
				{/each}
			</div>

			<!-- Right side actions -->
			<div class="flex items-center gap-4">
				<LanguageSwitcher />
				
				<!-- Mobile menu button -->
				<button
					onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
					class="md:hidden p-2 text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg transition-colors duration-200"
					aria-label="Toggle menu"
				>
					{#if mobileMenuOpen}
						<X class="w-6 h-6" />
					{:else}
						<Menu class="w-6 h-6" />
					{/if}
				</button>
			</div>
		</div>

		<!-- Mobile Navigation -->
		{#if mobileMenuOpen}
			<div
				class="md:hidden mt-4 pb-4 border-t border-gray-800/50 pt-4 animate-slideDown"
			>
				<div class="flex flex-col gap-2">
					{#each navLinks as link}
						<a
							href={link.href}
							onclick={(e) => {
								e.preventDefault();
								handleNavClick(link.href);
							}}
							class="px-4 py-3 text-base font-medium text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg transition-all duration-200 flex items-center gap-3 group"
						>
							<span
								class="w-1.5 h-1.5 rounded-full bg-gradient-to-r from-blue-400 to-purple-400 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
							></span>
							{link.label}
						</a>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</nav>

<style>
	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.animate-slideDown) {
		animation: slideDown 0.3s ease-out;
	}
</style>

