<script lang="ts">
	import type { Client } from '$lib/api/clients';
	import { sendMessage } from '$lib/api/whatsapp';
	import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';

	export let client: Client | null;
	export let mode: 'manual' | 'template' | 'instructions' | 'report' = 'manual';
	export let onClose: () => void;
	export let onSent: () => void;

	let message = '';
	let template = '';
	let instructions = '';
	let loading = false;
	let error = '';

	$: clientContext = (() => {
		if (!client) return '';
		const parts: string[] = [];
		if (client.name) parts.push(`Name: ${client.name}`);
		if (client.company) parts.push(`Company: ${client.company}`);
		if (client.email) parts.push(`Email: ${client.email}`);
		if (client.notes) parts.push(`Notes: ${client.notes}`);
		return parts.join('\n');
	})();

	const handleSubmit = async () => {
		if (!client) return;

		if (mode === 'manual') {
			if (!message.trim()) {
				error = 'Please enter a message';
				return;
			}
		} else if (mode === 'template') {
			if (!template.trim()) {
				error = 'Please provide a template';
				return;
			}
		} else if (mode === 'instructions') {
			if (!instructions.trim()) {
				error = 'Please provide instructions for AI generation';
				return;
			}
		} else if (mode === 'report') {
			error = 'Report sending is not yet implemented';
			toastError(error);
			return;
		}

		loading = true;
		error = '';

		try {
			await sendMessage({
				client_id: client.id,
				message: mode === 'manual' ? message.trim() : undefined,
				use_ai: mode === 'template' || mode === 'instructions',
				template: mode === 'template' ? template.trim() : undefined,
				instructions: mode === 'instructions' ? instructions.trim() : undefined,
				client_context: clientContext || undefined
			});

			toastSuccess('Message sent successfully!');
			onSent();
			onClose();
		} catch (err: unknown) {
			error = getApiErrorMessage(err, 'Failed to send message');
			toastError(error);
		} finally {
			loading = false;
		}
	};

	const handleClose = () => {
		message = '';
		template = '';
		instructions = '';
		error = '';
		onClose();
	};

	$: modalTitle = (() => {
		switch (mode) {
			case 'manual':
				return 'Send Message';
			case 'template':
				return 'Send from Template';
			case 'instructions':
				return 'Generate with AI';
			case 'report':
				return 'Send Report';
			default:
				return 'Send Message';
		}
	})();
</script>

{#if client}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
		on:click={handleClose}
		role="button"
		tabindex="0"
		on:keydown={(e) => e.key === 'Escape' && handleClose()}
	>
		<div
			class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
			on:click|stopPropagation
		>
			<div class="p-6 border-b border-gray-200 dark:border-gray-700">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="text-2xl font-semibold text-gray-900 dark:text-white">
							{modalTitle}
						</h2>
						<p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
							To: {client.name} ({client.phoneNumber})
						</p>
					</div>
					<button
						on:click={handleClose}
						class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							></path>
						</svg>
					</button>
				</div>
			</div>

			<div class="p-6 space-y-6">
				{#if error}
					<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
						<p class="text-red-800 dark:text-red-200 text-sm">{error}</p>
					</div>
				{/if}

				{#if mode === 'manual'}
					<div>
						<label for="message" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
							Message <span class="text-red-500">*</span>
						</label>
						<textarea
							id="message"
							bind:value={message}
							rows="8"
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
							placeholder="Type your message here..."
						></textarea>
						<p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
							Write your message directly. It will be sent as-is to {client?.name}.
						</p>
					</div>
				{:else if mode === 'template'}
					<div>
						<label
							for="template"
							class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
						>
							Template <span class="text-red-500">*</span>
						</label>
						<textarea
							id="template"
							bind:value={template}
							rows="8"
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
							placeholder="Example: Hi {name}, I wanted to follow up on our recent conversation. Let me know if you have any questions..."
						></textarea>
						<p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
							Provide a template structure. The AI will generate a message following this template format.
						</p>
					</div>
				{:else if mode === 'instructions'}
					<div>
						<label
							for="instructions"
							class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
						>
							Instructions <span class="text-red-500">*</span>
						</label>
						<textarea
							id="instructions"
							bind:value={instructions}
							rows="8"
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="Example: Write a friendly follow-up message asking about their project status. Be professional but warm..."
						></textarea>
						<p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
							Describe what kind of message you want the AI to generate. The AI will use client information (name, company, etc.) to personalize the message.
						</p>
					</div>
				{:else if mode === 'report'}
					<div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
						<p class="text-yellow-800 dark:text-yellow-200 text-sm">
							Report sending functionality is coming soon. This will allow you to send automated reports to clients.
						</p>
					</div>
				{/if}

				<div class="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
					<button
						type="button"
						on:click={handleClose}
						disabled={loading}
						class="px-4 py-2 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-900 dark:text-white rounded-lg transition-colors disabled:opacity-50"
					>
						Cancel
					</button>
					<button
						type="button"
						on:click={handleSubmit}
						disabled={loading || mode === 'report'}
						class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{loading ? 'Sending...' : mode === 'report' ? 'Coming Soon' : 'Send Message'}
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

