<script lang="ts">
	import CertificationsSectionComponent from '../_components/CertificationsSection.svelte';
	import type { Certification } from '$lib/types/certification';

	interface Props {
		certifications: Certification[];
		loading: boolean;
		translations: any;
	}

	let { certifications, loading, translations: t }: Props = $props();
</script>

<section id="certifications" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-green-400 to-emerald-400 bg-clip-text text-transparent">
				{t('certifications.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('certifications.subtitle')}
			</p>
		</div>
		<CertificationsSectionComponent certifications={certifications} loading={loading} />
	</div>
</section>

