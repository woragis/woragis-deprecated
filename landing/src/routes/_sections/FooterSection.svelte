<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiGithub, SiInstagram, SiLinkedin } from 'svelte-icons-pack/si';
	import { contact } from '$lib/constants';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<footer class="container mx-auto px-6 py-12 border-t border-gray-700/50 bg-gradient-to-b from-gray-900/50 to-gray-900">
	<div class="max-w-4xl mx-auto text-center">
		<div class="flex justify-center gap-6 mb-6">
			<a
				href={contact.github}
				target="_blank"
				rel="noopener noreferrer"
				class="text-gray-400 hover:text-white transition-all duration-200 hover:scale-110"
				aria-label="GitHub"
			>
				<Icon src={SiGithub} size="1.5rem" />
			</a>
			<a
				href={contact.instagram}
				target="_blank"
				rel="noopener noreferrer"
				class="text-gray-400 hover:text-white transition-all duration-200 hover:scale-110"
				aria-label="Instagram"
			>
				<Icon src={SiInstagram} size="1.5rem" />
			</a>
			{#if contact.linkedin}
				<a
					href={contact.linkedin}
					target="_blank"
					rel="noopener noreferrer"
					class="text-gray-400 hover:text-white transition-all duration-200 hover:scale-110"
					aria-label="LinkedIn"
				>
					<Icon src={SiLinkedin} size="1.5rem" />
				</a>
			{/if}
		</div>
		<p class="mb-4 text-gray-300">
			{t('footer.tagline')}
			<span class="text-gray-500 mx-2">|</span> {t('footer.backendDeveloper')}
			<span class="text-gray-500 mx-2">|</span> {t('footer.devopsEnthusiast')}
		</p>
		<p class="text-sm text-gray-400">© 2024</p>
	</div>
</footer>

