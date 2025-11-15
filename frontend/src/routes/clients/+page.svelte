<script lang="ts">
	import { createClientsLogic } from './clients.logic';
	import ClientsList from './_components/ClientsList.svelte';
	import ClientForm from './_components/ClientForm.svelte';
	import SendMessageModal from './_components/SendMessageModal.svelte';
	import AuthNotice from '../projects/_components/AuthNotice.svelte';
	import ErrorBanner from '../projects/_components/ErrorBanner.svelte';
	import type { Client } from '$lib/api/clients';

	const {
		isAuthenticated,
		error,
		loading,
		clients,
		clientForm,
		editingClient,
		showArchived,
		loadClients,
		handleCreateClient,
		handleUpdateClient,
		handleDeleteClient,
		handleToggleArchived,
		startEdit,
		cancelEdit,
		updateClientFormField
	} = createClientsLogic();

	let selectedClientForMessage: Client | null = null;
	let messageMode: 'manual' | 'template' | 'instructions' | 'report' = 'manual';

	const handleSendMessage = (client: Client, mode: 'manual' | 'template' | 'instructions' | 'report') => {
		selectedClientForMessage = client;
		messageMode = mode;
	};

	const handleMessageSent = () => {
		selectedClientForMessage = null;
		messageMode = 'manual';
		// Optionally reload clients or show success
	};

	const handleCloseModal = () => {
		selectedClientForMessage = null;
		messageMode = 'manual';
	};
</script>

<section class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900 dark:text-white">Clients</h1>
			<p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
				Manage your clients and contacts for WhatsApp messaging
			</p>
		</div>
		<div class="flex items-center gap-4">
			<label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
				<input
					type="checkbox"
					bind:checked={$showArchived}
					on:change={loadClients}
					class="rounded border-gray-300 dark:border-gray-600"
				/>
				<span>Show archived</span>
			</label>
			<button
				on:click={loadClients}
				disabled={$loading}
				class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg transition-colors disabled:opacity-50"
			>
				Refresh
			</button>
		</div>
	</div>

	{#if $error}
		<ErrorBanner message={$error} />
	{/if}

	{#if !$isAuthenticated}
		<AuthNotice />
	{:else}
		<div class="grid gap-6 lg:grid-cols-[1.05fr_2fr]">
			<div class="space-y-6">
				<ClientForm
					formState={$clientForm}
					editingClient={$editingClient}
					onFieldChange={updateClientFormField}
					onSubmit={$editingClient ? () => handleUpdateClient($editingClient) : handleCreateClient}
					onCancel={$editingClient ? cancelEdit : undefined}
				/>
			</div>

			<div class="space-y-6">
				<ClientsList
					clients={$clients}
					loading={$loading}
					onEdit={startEdit}
					onDelete={handleDeleteClient}
					onToggleArchived={handleToggleArchived}
					onSendMessage={handleSendMessage}
				/>
			</div>
		</div>
	{/if}

	{#if selectedClientForMessage}
		<SendMessageModal
			client={selectedClientForMessage}
			mode={messageMode}
			onClose={handleCloseModal}
			onSent={handleMessageSent}
		/>
	{/if}
</section>

