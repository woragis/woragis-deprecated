<script lang="ts">
	import SkillsTimeline from '../_components/SkillsTimeline.svelte';

	interface Props {
		translations: any;
	}

	let { translations: t }: Props = $props();
</script>

<section id="skills-timeline" class="container mx-auto px-6 py-20">
	<div class="max-w-7xl mx-auto">
		<div class="text-center mb-12">
			<h2 class="text-4xl font-bold mb-4 bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
				{t('skills.timeline.title')}
			</h2>
			<p class="text-gray-400 text-lg max-w-2xl mx-auto">
				{t('skills.timeline.subtitle')}
			</p>
		</div>
		<SkillsTimeline />
	</div>
</section>

