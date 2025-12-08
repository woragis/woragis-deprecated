import { onDestroy, onMount } from 'svelte';
import { derived, get, writable } from 'svelte/store';
import { useQueryClient } from '@tanstack/svelte-query';

import { authStore } from '$lib';
import type { Client, CreateClientInput, UpdateClientInput } from '$lib/api/clients';
import { getApiErrorMessage, toastError, toastSuccess } from '$lib/utils/toast';
import {
	createClient,
	deleteClient,
	fetchClients,
	toggleClientArchived,
	updateClient
} from '$lib/api/clients';

export interface ClientFormState {
	name: string;
	email: string;
	phone_number: string;
	company: string;
	notes: string;
}

const defaultClientForm = (): ClientFormState => ({
	name: '',
	email: '',
	phone_number: '',
	company: '',
	notes: ''
});

interface AuthState {
	isAuthenticated: boolean;
	userId: string | null;
}

export function createClientsLogic() {
	const queryClient = useQueryClient();

	const authStateStore = writable<AuthState>({ isAuthenticated: false, userId: null });
	const isAuthenticated = derived(authStateStore, ($auth) => $auth.isAuthenticated);
	const error = writable<string | null>(null);
	const loading = writable(false);

	const clients = writable<Client[]>([]);
	const clientForm = writable<ClientFormState>(defaultClientForm());
	const editingClient = writable<Client | null>(null);
	const showArchived = writable(false);

	const resetStateForUnauthenticated = () => {
		clients.set([]);
		clientForm.set(defaultClientForm());
		error.set(null);
		editingClient.set(null);
	};

	const loadClients = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			const message = 'You must be signed in to load clients.';
			error.set(message);
			toastError(message);
			return;
		}

		loading.set(true);
		error.set(null);
		try {
			const data = await fetchClients(get(showArchived));
			clients.set(data);
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to load clients.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	};

	const handleCreateClient = async () => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			toastError('You must be signed in to create clients.');
			return;
		}

		const form = get(clientForm);
		if (!form.name.trim() || !form.phone_number.trim()) {
			toastError('Name and phone number are required.');
			return;
		}

		loading.set(true);
		error.set(null);
		try {
			const input: CreateClientInput = {
				name: form.name.trim(),
				email: form.email.trim() || undefined,
				phoneNumber: form.phone_number.trim(),
				company: form.company.trim() || undefined,
				notes: form.notes.trim() || undefined
			};

			await createClient(input);
			clientForm.set(defaultClientForm());
			toastSuccess('Client created successfully.');
			await loadClients();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to create client.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	};

	const handleUpdateClient = async (client: Client) => {
		const auth = get(authStateStore);
		if (!auth.isAuthenticated) {
			toastError('You must be signed in to update clients.');
			return;
		}

		loading.set(true);
		error.set(null);
		try {
			const form = get(clientForm);
			const input: UpdateClientInput = {
				name: form.name.trim() || undefined,
				email: form.email.trim() || undefined,
				phoneNumber: form.phone_number.trim() || undefined,
				company: form.company.trim() || undefined,
				notes: form.notes.trim() || undefined
			};

			await updateClient(client.id, input);
			editingClient.set(null);
			clientForm.set(defaultClientForm());
			toastSuccess('Client updated successfully.');
			await loadClients();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update client.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	};

	const handleDeleteClient = async (id: string) => {
		if (!confirm('Are you sure you want to delete this client? This action cannot be undone.')) {
			return;
		}

		loading.set(true);
		error.set(null);
		try {
			await deleteClient(id);
			toastSuccess('Client deleted successfully.');
			await loadClients();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to delete client.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	};

	const handleToggleArchived = async (id: string, archived: boolean) => {
		loading.set(true);
		error.set(null);
		try {
			await toggleClientArchived(id, archived);
			toastSuccess(archived ? 'Client archived.' : 'Client restored.');
			await loadClients();
		} catch (err) {
			const message = getApiErrorMessage(err, 'Unable to update client.');
			error.set(message);
			toastError(message);
		} finally {
			loading.set(false);
		}
	};

	const startEdit = (client: Client) => {
		editingClient.set(client);
		clientForm.set({
			name: client.name,
			email: client.email || '',
			phone_number: client.phoneNumber,
			company: client.company || '',
			notes: client.notes || ''
		});
	};

	const cancelEdit = () => {
		editingClient.set(null);
		clientForm.set(defaultClientForm());
	};

	const updateClientFormField = <K extends keyof ClientFormState>(
		field: K,
		value: ClientFormState[K]
	) => {
		clientForm.update((form) => ({ ...form, [field]: value }));
	};

	const authUnsubscribe = authStore.subscribe((state) => {
		authStateStore.set({
			isAuthenticated: state.isAuthenticated,
			userId: state.user?.id ?? null
		});

		if (!state.isAuthenticated) {
			resetStateForUnauthenticated();
		} else {
			loadClients();
		}
	});

	onMount(() => {
		const auth = get(authStore);
		if (auth.isAuthenticated) {
			loadClients();
		}
	});

	onDestroy(() => {
		authUnsubscribe();
	});

	return {
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
	};
}

