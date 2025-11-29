<script lang="ts">
	import type { Project } from '$lib/api/types';
	import StatCard from '$lib/components/StatCard.svelte';

	export let projects: Project[] = [];
</script>

<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
	<StatCard label="Total Projects" value={projects.length} accentColor="indigo" />
	<StatCard
		label="Active Projects"
		value={projects.filter((p) => p.status === 'active').length}
		accentColor="emerald"
	/>
	<StatCard
		label="Total MRR"
		value={`$${projects.reduce((sum, p) => sum + p.mrr, 0).toFixed(2)}`}
		accentColor="purple"
	/>
</div>

