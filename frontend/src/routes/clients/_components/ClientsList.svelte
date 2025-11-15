<script lang="ts">
	import type { Client } from '$lib/api/clients';

	export let clients: Client[];
	export let loading: boolean;
	export let onEdit: (client: Client) => void;
	export let onDelete: (id: string) => void;
	export let onToggleArchived: (id: string, archived: boolean) => void;
	export let onSendMessage: (client: Client, mode: 'manual' | 'template' | 'instructions' | 'report') => void;
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg">
	<div class="p-6 border-b border-gray-200 dark:border-gray-700">
		<h2 class="text-xl font-semibold text-gray-900 dark:text-white">Clients List</h2>
		<p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
			{clients.length} {clients.length === 1 ? 'client' : 'clients'}
		</p>
	</div>

	{#if loading && clients.length === 0}
		<div class="p-6 text-center text-gray-600 dark:text-gray-400">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
			<p>Loading clients...</p>
		</div>
	{:else if clients.length === 0}
		<div class="p-6 text-center text-gray-600 dark:text-gray-400">
			<p>No clients found. Create your first client to get started.</p>
		</div>
	{:else}
		<div class="divide-y divide-gray-200 dark:divide-gray-700">
			{#each clients as client (client.id)}
				<div
					class="p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors {client.is_archived
						? 'opacity-60'
						: ''}"
				>
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<div class="flex items-center gap-2">
								<h3 class="font-semibold text-gray-900 dark:text-white">{client.name}</h3>
								{#if client.is_archived}
									<span
										class="px-2 py-0.5 text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded"
									>
										Archived
									</span>
								{/if}
							</div>
							<div class="mt-1 space-y-1 text-sm text-gray-600 dark:text-gray-400">
								{#if client.phone_number}
									<p>
										<span class="font-medium">Phone:</span> {client.phone_number}
									</p>
								{/if}
								{#if client.email}
									<p>
										<span class="font-medium">Email:</span> {client.email}
									</p>
								{/if}
								{#if client.company}
									<p>
										<span class="font-medium">Company:</span> {client.company}
									</p>
								{/if}
								{#if client.notes}
									<p class="text-xs text-gray-500 dark:text-gray-500 line-clamp-2">
										{client.notes}
									</p>
								{/if}
							</div>
						</div>
						<div class="flex items-center gap-2 ml-4 flex-wrap">
							{#if client.phone_number}
								<div class="flex items-center gap-1 border-r border-gray-300 dark:border-gray-600 pr-2 mr-1">
									<button
										on:click={() => onSendMessage(client, 'manual')}
										class="px-3 py-1.5 text-sm bg-green-100 hover:bg-green-200 dark:bg-green-900/30 dark:hover:bg-green-900/50 text-green-700 dark:text-green-300 rounded transition-colors font-medium"
										title="Send a message you write"
									>
										💬 Send Message
									</button>
									<button
										on:click={() => onSendMessage(client, 'template')}
										disabled
										class="px-3 py-1.5 text-sm bg-purple-100 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400 rounded transition-colors opacity-50 cursor-not-allowed"
										title="Send from template (coming soon)"
									>
										📝 Template
									</button>
									<button
										on:click={() => onSendMessage(client, 'instructions')}
										disabled
										class="px-3 py-1.5 text-sm bg-blue-100 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded transition-colors opacity-50 cursor-not-allowed"
										title="Send from AI instructions (coming soon)"
									>
										🤖 AI Generate
									</button>
									<button
										on:click={() => onSendMessage(client, 'report')}
										disabled
										class="px-3 py-1.5 text-sm bg-orange-100 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400 rounded transition-colors opacity-50 cursor-not-allowed"
										title="Send report (coming soon)"
									>
										📊 Report
									</button>
								</div>
							{/if}
							<button
								on:click={() => onEdit(client)}
								class="px-3 py-1.5 text-sm bg-blue-100 hover:bg-blue-200 dark:bg-blue-900/30 dark:hover:bg-blue-900/50 text-blue-700 dark:text-blue-300 rounded transition-colors"
							>
								Edit
							</button>
							<button
								on:click={() => onToggleArchived(client.id, !client.is_archived)}
								class="px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded transition-colors"
							>
								{client.is_archived ? 'Restore' : 'Archive'}
							</button>
							<button
								on:click={() => onDelete(client.id)}
								class="px-3 py-1.5 text-sm bg-red-100 hover:bg-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 text-red-700 dark:text-red-300 rounded transition-colors"
							>
								Delete
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

