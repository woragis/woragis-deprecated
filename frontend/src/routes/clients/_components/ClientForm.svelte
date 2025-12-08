<script lang="ts">
	import type { ClientFormState } from '../clients.logic';
	import type { Client } from '$lib/api/clients';

	export let formState: ClientFormState;
	export let editingClient: Client | null = null;
	export let onFieldChange: (field: keyof ClientFormState, value: string) => void;
	export let onSubmit: () => void;
	export let onCancel: (() => void) | undefined = undefined;
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
	<h2 class="text-xl font-semibold mb-4 text-gray-900 dark:text-white">
		{editingClient ? 'Edit Client' : 'Add New Client'}
	</h2>

	<form
		on:submit|preventDefault={onSubmit}
		class="space-y-4"
	>
		<div>
			<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Name <span class="text-red-500">*</span>
			</label>
			<input
				id="name"
				type="text"
				value={formState.name}
				on:input={(e) => onFieldChange('name', e.currentTarget.value)}
				required
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="John Doe"
			/>
		</div>

		<div>
			<label for="phone_number" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Phone Number <span class="text-red-500">*</span>
			</label>
			<input
				id="phone_number"
				type="tel"
				value={formState.phone_number}
				on:input={(e) => onFieldChange('phone_number', e.currentTarget.value)}
				required
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="+1234567890"
			/>
			<p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
				Include country code (e.g., +1 for US, +55 for Brazil)
			</p>
		</div>

		<div>
			<label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Email
			</label>
			<input
				id="email"
				type="email"
				value={formState.email}
				on:input={(e) => onFieldChange('email', e.currentTarget.value)}
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="john@example.com"
			/>
		</div>

		<div>
			<label for="company" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Company
			</label>
			<input
				id="company"
				type="text"
				value={formState.company}
				on:input={(e) => onFieldChange('company', e.currentTarget.value)}
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Acme Inc."
			/>
		</div>

		<div>
			<label for="notes" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
				Notes
			</label>
			<textarea
				id="notes"
				value={formState.notes}
				on:input={(e) => onFieldChange('notes', e.currentTarget.value)}
				rows="3"
				class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				placeholder="Additional notes about this client..."
			></textarea>
		</div>

		<div class="flex gap-2">
			<button
				type="submit"
				class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium"
			>
				{editingClient ? 'Update Client' : 'Create Client'}
			</button>
			{#if onCancel}
				<button
					type="button"
					on:click={onCancel}
					class="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg transition-colors"
				>
					Cancel
				</button>
			{/if}
		</div>
	</form>
</div>

