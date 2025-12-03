<script lang="ts">
	import { Icon } from 'svelte-icons-pack';
	import { SiGithub, SiInstagram, SiLinkedin, SiWhatsapp } from 'svelte-icons-pack/si';
	import { Mail, Phone, MapPin } from 'lucide-svelte';
	import { contact } from '$lib/constants';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();

	// Generate WhatsApp URL with pre-filled message
	function getWhatsAppUrl(): string {
		if (!contact.whatsapp) return '#';
		const phoneNumber = contact.whatsapp.replace(/\D/g, '');
		const message = encodeURIComponent(t('contact.whatsappMessage'));
		return `https://wa.me/${phoneNumber}?text=${message}`;
	}
</script>

<section id="contact" class="container mx-auto px-6 py-20 relative overflow-hidden">
	<!-- Background decorative elements -->
	<div
		class="absolute inset-0 opacity-10 pointer-events-none"
		style="background-image: radial-gradient(circle at 20% 50%, rgba(59, 130, 246, 0.5) 0%, transparent 50%), radial-gradient(circle at 80% 80%, rgba(147, 51, 234, 0.5) 0%, transparent 50%);"
	></div>

	<div class="max-w-6xl mx-auto relative z-10">
		<!-- Header -->
		<div class="text-center mb-16">
			<div
				class="inline-block mb-4 px-4 py-2 bg-gradient-to-r from-blue-600/20 to-purple-600/20 border border-blue-500/30 rounded-full text-blue-300 text-sm font-medium"
			>
				{t('contact.title')}
			</div>
			<h2 class="text-4xl md:text-5xl font-bold mb-6 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('contact.title')}
			</h2>
			<p class="text-xl text-gray-300 max-w-2xl mx-auto leading-relaxed">
				{t('contact.subtitle')}
			</p>
		</div>

		<!-- Contact Cards Grid -->
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
			<!-- Email Card -->
			<a
				href="mailto:{contact.email}"
				class="group relative bg-gradient-to-br from-blue-900/40 via-blue-800/30 to-blue-900/40 backdrop-blur-sm rounded-2xl p-6 border border-blue-700/30 hover:border-blue-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-blue-500/30 hover:-translate-y-1 overflow-hidden"
			>
				<!-- Animated background gradient -->
				<div
					class="absolute inset-0 bg-gradient-to-br from-blue-600/0 via-blue-600/0 to-blue-600/0 group-hover:from-blue-600/10 group-hover:via-blue-600/5 group-hover:to-blue-600/10 transition-all duration-500"
				></div>

				<div class="relative z-10">
					<div
						class="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-blue-500/30"
					>
						<Mail class="w-8 h-8 text-white" />
					</div>
					<h3 class="text-xl font-bold text-white mb-2 group-hover:text-blue-300 transition-colors">
						{t('contact.email')}
					</h3>
					<p class="text-gray-300 text-sm break-all group-hover:text-white transition-colors">
						{contact.email}
					</p>
					<p class="text-blue-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
						{t('contact.clickToEmail')}
					</p>
				</div>
			</a>

			<!-- Phone Card -->
			{#if contact.phone}
				<a
					href="tel:{contact.phone}"
					class="group relative bg-gradient-to-br from-green-900/40 via-green-800/30 to-green-900/40 backdrop-blur-sm rounded-2xl p-6 border border-green-700/30 hover:border-green-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-green-500/30 hover:-translate-y-1 overflow-hidden"
				>
					<div
						class="absolute inset-0 bg-gradient-to-br from-green-600/0 via-green-600/0 to-green-600/0 group-hover:from-green-600/10 group-hover:via-green-600/5 group-hover:to-green-600/10 transition-all duration-500"
					></div>

					<div class="relative z-10">
						<div
							class="w-16 h-16 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-green-500/30"
						>
							<Phone class="w-8 h-8 text-white" />
						</div>
						<h3 class="text-xl font-bold text-white mb-2 group-hover:text-green-300 transition-colors">
							{t('contact.phone')}
						</h3>
						<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
							{contact.phone}
						</p>
						<p class="text-green-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
							{t('contact.clickToCall')}
						</p>
					</div>
				</a>
			{/if}

			<!-- Location Card -->
			{#if contact.location}
				<div
					class="group relative bg-gradient-to-br from-purple-900/40 via-purple-800/30 to-purple-900/40 backdrop-blur-sm rounded-2xl p-6 border border-purple-700/30 overflow-hidden"
				>
					<div class="relative z-10">
						<div
							class="w-16 h-16 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-purple-500/30"
						>
							<MapPin class="w-8 h-8 text-white" />
						</div>
						<h3 class="text-xl font-bold text-white mb-2">{t('contact.location')}</h3>
						<p class="text-gray-300 text-sm">{contact.location}</p>
					</div>
				</div>
			{/if}

			<!-- LinkedIn Card -->
			{#if contact.linkedin}
				<a
					href={contact.linkedin}
					target="_blank"
					rel="noopener noreferrer"
					class="group relative bg-gradient-to-br from-indigo-900/40 via-indigo-800/30 to-indigo-900/40 backdrop-blur-sm rounded-2xl p-6 border border-indigo-700/30 hover:border-indigo-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-indigo-500/30 hover:-translate-y-1 overflow-hidden"
				>
					<div
						class="absolute inset-0 bg-gradient-to-br from-indigo-600/0 via-indigo-600/0 to-indigo-600/0 group-hover:from-indigo-600/10 group-hover:via-indigo-600/5 group-hover:to-indigo-600/10 transition-all duration-500"
					></div>

					<div class="relative z-10">
						<div
							class="w-16 h-16 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-indigo-500/30"
						>
							<Icon src={SiLinkedin} size="2rem" color="white" />
						</div>
						<h3 class="text-xl font-bold text-white mb-2 group-hover:text-indigo-300 transition-colors">
							{t('contact.linkedin')}
						</h3>
						<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
							{t('contact.subtitle')}
						</p>
						<p class="text-indigo-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
							View profile →
						</p>
					</div>
				</a>
			{/if}

			<!-- WhatsApp Card -->
			{#if contact.whatsapp}
				<a
					href={getWhatsAppUrl()}
					target="_blank"
					rel="noopener noreferrer"
					class="group relative bg-gradient-to-br from-emerald-900/40 via-emerald-800/30 to-emerald-900/40 backdrop-blur-sm rounded-2xl p-6 border border-emerald-700/30 hover:border-emerald-500/60 transition-all duration-300 hover:shadow-2xl hover:shadow-emerald-500/30 hover:-translate-y-1 overflow-hidden"
				>
					<div
						class="absolute inset-0 bg-gradient-to-br from-emerald-600/0 via-emerald-600/0 to-emerald-600/0 group-hover:from-emerald-600/10 group-hover:via-emerald-600/5 group-hover:to-emerald-600/10 transition-all duration-500"
					></div>

					<div class="relative z-10">
						<div
							class="w-16 h-16 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-lg shadow-emerald-500/30"
						>
							<Icon src={SiWhatsapp} size="2rem" color="white" />
						</div>
						<h3 class="text-xl font-bold text-white mb-2 group-hover:text-emerald-300 transition-colors">
							{t('contact.whatsapp')}
						</h3>
						<p class="text-gray-300 text-sm group-hover:text-white transition-colors">
							{t('contact.subtitle')}
						</p>
						<p class="text-emerald-400 text-xs mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
							Start conversation →
						</p>
					</div>
				</a>
			{/if}
		</div>

		<!-- Social Links Section -->
		<div
			class="relative bg-gradient-to-br from-gray-800/60 via-gray-800/40 to-gray-800/60 backdrop-blur-sm rounded-3xl p-8 md:p-12 border border-gray-700/50 overflow-hidden"
		>
			<!-- Decorative gradient overlay -->
			<div
				class="absolute inset-0 bg-gradient-to-r from-blue-600/5 via-purple-600/5 to-pink-600/5 opacity-50"
			></div>

			<div class="relative z-10">
				<div class="text-center mb-8">
					<h3 class="text-2xl font-bold text-white mb-2">Follow My Journey</h3>
					<p class="text-gray-400">Connect with me on social media</p>
				</div>

				<div class="flex flex-wrap justify-center gap-8">
					<a
						href={contact.github}
						target="_blank"
						rel="noopener noreferrer"
						class="group relative flex flex-col items-center gap-3 p-6 bg-gray-700/30 hover:bg-gray-700/50 rounded-xl border border-gray-600/50 hover:border-gray-500 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-gray-500/20"
						aria-label="GitHub"
					>
						<div
							class="w-14 h-14 bg-gradient-to-br from-gray-700 to-gray-800 rounded-xl flex items-center justify-center group-hover:from-gray-600 group-hover:to-gray-700 transition-all duration-300 shadow-lg"
						>
							<Icon src={SiGithub} size="2rem" color="white" />
						</div>
						<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
							>{t('contact.github')}</span
						>
					</a>

					<a
						href={contact.instagram}
						target="_blank"
						rel="noopener noreferrer"
						class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-pink-900/30 to-purple-900/30 hover:from-pink-800/40 hover:to-purple-800/40 rounded-xl border border-pink-700/30 hover:border-pink-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-pink-500/20"
						aria-label="Instagram"
					>
						<div
							class="w-14 h-14 bg-gradient-to-br from-pink-500 to-purple-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-pink-500/30"
						>
							<Icon src={SiInstagram} size="2rem" color="white" />
						</div>
						<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
							>{t('contact.instagram')}</span
						>
					</a>

					{#if contact.linkedin}
						<a
							href={contact.linkedin}
							target="_blank"
							rel="noopener noreferrer"
							class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-blue-900/30 to-indigo-900/30 hover:from-blue-800/40 hover:to-indigo-800/40 rounded-xl border border-blue-700/30 hover:border-blue-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-blue-500/20"
							aria-label="LinkedIn"
						>
							<div
								class="w-14 h-14 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-blue-500/30"
							>
								<Icon src={SiLinkedin} size="2rem" color="white" />
							</div>
							<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
								>{t('contact.linkedin')}</span
							>
						</a>
					{/if}

					{#if contact.whatsapp}
						<a
							href={getWhatsAppUrl()}
							target="_blank"
							rel="noopener noreferrer"
							class="group relative flex flex-col items-center gap-3 p-6 bg-gradient-to-br from-green-900/30 to-emerald-900/30 hover:from-green-800/40 hover:to-emerald-800/40 rounded-xl border border-green-700/30 hover:border-green-500/50 transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-green-500/20"
							aria-label="WhatsApp"
						>
							<div
								class="w-14 h-14 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-green-500/30"
							>
								<Icon src={SiWhatsapp} size="2rem" color="white" />
							</div>
							<span class="text-gray-300 text-sm font-medium group-hover:text-white transition-colors"
								>{t('contact.whatsapp')}</span
							>
						</a>
					{/if}
				</div>
			</div>
		</div>
	</div>
</section>

