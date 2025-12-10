<script lang="ts">
	import Tabs from '$lib/components/ui/Tabs.svelte';
	import { MessageSquare } from 'lucide-svelte';
	
	let {
		activeTab = $bindable('details' as 'details' | 'chat'),
		onTabChange
	}: {
		activeTab?: 'details' | 'chat';
		onTabChange?: (tab: 'details' | 'chat') => void;
	} = $props();
	
	const tabs = [
		{ id: 'details', label: 'Details' },
		{ id: 'chat', label: 'Chat', icon: MessageSquare }
	];
	
	function handleTabChange(tabId: string) {
		activeTab = tabId as 'details' | 'chat';
		onTabChange?.(activeTab);
	}
</script>

<Tabs {tabs} bind:activeTab onTabChange={handleTabChange} />
