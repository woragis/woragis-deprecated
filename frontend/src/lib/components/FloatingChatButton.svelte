<script lang="ts">
	import { MessageSquare, X } from 'lucide-svelte';
	import { goto } from '$app/navigation';

	interface Props {
		jobApplicationId?: string;
		jobTitle?: string;
		companyName?: string;
	}

	let { jobApplicationId, jobTitle, companyName }: Props = $props();

	let isOpen = $state(false);
	let selectedJobApplicationId = $state<string | undefined>(jobApplicationId);

	function handleOpen() {
		if (jobApplicationId) {
			// Navigate to chat with this job application
			goto(`/chats?jobApplicationId=${jobApplicationId}`);
		} else {
			// Open chat selection modal or navigate to chats page
			isOpen = true;
		}
	}

	function handleSelectJobApplication(id: string) {
		selectedJobApplicationId = id;
		goto(`/chats?jobApplicationId=${id}`);
		isOpen = false;
	}
</script>

<!-- Floating Button -->
<button
	onclick={handleOpen}
	class="fixed bottom-6 right-6 w-14 h-14 bg-blue-600 hover:bg-blue-700 text-white rounded-full shadow-lg flex items-center justify-center z-50 transition-all hover:scale-110"
	title="Open Chat"
>
	<MessageSquare class="w-6 h-6" />
</button>

<!-- Selection Modal (if no jobApplicationId provided) -->
{#if isOpen && !jobApplicationId}
	<div
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
		onclick={() => (isOpen = false)}
	>
		<div
			class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white">Select Job Application</h3>
				<button
					onclick={() => (isOpen = false)}
					class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
				>
					<X class="w-5 h-5" />
				</button>
			</div>
			<p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
				Choose a job application to chat about, or navigate to the chats page to start a general
				conversation.
			</p>
			<div class="flex gap-3">
				<button
					onclick={() => {
						goto('/chats');
						isOpen = false;
					}}
					class="flex-1 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-md hover:bg-gray-300 dark:hover:bg-gray-600"
				>
					Go to Chats
				</button>
			</div>
		</div>
	</div>
{/if}

