<script lang="ts">
	import { onMount } from 'svelte';
	import {
		Search,
		Edit,
		Trash2,
		Check,
		X,
		Users,
		AlertCircle,
		ChevronLeft,
		ChevronRight,
		MoreVertical,
		Shield,
		Mail,
		Phone,
		Globe
	} from 'lucide-svelte';
	import {
		listUsers,
		updateUser,
		bulkUpdateUsers,
		getUserAuditLogs,
		type AdminUser,
		type AdminUserListResponse
	} from '$lib/api/admin';
	import { toastSuccess, toastError } from '$lib/utils/toast';

	let users: AdminUser[] = $state([]);
	let total = $state(0);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');
	let currentPage = $state(1);
	let limit = $state(20);
	let selectedUsers = $state<Set<string>>(new Set());
	let showBulkActions = $state(false);
	let showEditModal = $state(false);
	let showAuditModal = $state(false);
	let editingUser: AdminUser | null = $state(null);
	let viewingAuditLogsUserId: string | null = $state(null);
	let auditLogs = $state<any[]>([]);

	// Edit form state
	let editRole = $state('');
	let editEmail = $state('');
	let editPhoneNumber = $state('');
	let editPreferredLocale = $state('');
	let editConfirmEmail = $state(false);
	let editDisableMFA = $state(false);

	// Bulk actions
	let bulkSetRole = $state('');
	let bulkConfirmEmail = $state(false);
	let bulkDisableMFA = $state(false);

	onMount(async () => {
		await fetchUsers();
	});

	async function fetchUsers() {
		loading = true;
		error = null;
		try {
			const response: AdminUserListResponse = await listUsers({
				limit,
				offset: (currentPage - 1) * limit,
				search: searchQuery || undefined
			});
			users = response.users;
			total = response.total;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch users';
			toastError(error);
			console.error('Error fetching users:', err);
		} finally {
			loading = false;
		}
	}

	function handleSearch() {
		currentPage = 1;
		fetchUsers();
	}

	function handlePageChange(newPage: number) {
		currentPage = newPage;
		fetchUsers();
	}

	function toggleUserSelection(userId: string) {
		if (selectedUsers.has(userId)) {
			selectedUsers.delete(userId);
		} else {
			selectedUsers.add(userId);
		}
		selectedUsers = new Set(selectedUsers);
		showBulkActions = selectedUsers.size > 0;
	}

	function toggleSelectAll() {
		if (selectedUsers.size === users.length) {
			selectedUsers = new Set();
		} else {
			selectedUsers = new Set(users.map((u) => u.id));
		}
		showBulkActions = selectedUsers.size > 0;
	}

	function startEdit(user: AdminUser) {
		editingUser = user;
		editRole = user.role;
		editEmail = user.email;
		editPhoneNumber = user.phone_number || '';
		editPreferredLocale = user.preferred_locale || 'en';
		editConfirmEmail = !!user.email_confirmed_at;
		editDisableMFA = false;
		showEditModal = true;
	}

	function cancelEdit() {
		editingUser = null;
		showEditModal = false;
	}

	async function handleUpdate() {
		if (!editingUser) return;

		try {
			const payload: any = {};
			if (editRole !== editingUser.role) {
				payload.set_role = editRole;
			}
			if (editEmail !== editingUser.email) {
				payload.set_email = editEmail;
			}
			if (editPhoneNumber !== (editingUser.phone_number || '')) {
				payload.set_phone_number = editPhoneNumber || null;
			}
			if (editPreferredLocale !== (editingUser.preferred_locale || 'en')) {
				payload.set_preferred_locale = editPreferredLocale;
			}
			if (editConfirmEmail && !editingUser.email_confirmed_at) {
				payload.confirm_email = true;
			}
			if (editDisableMFA && editingUser.mfa_enabled) {
				payload.disable_mfa = true;
			}

			if (Object.keys(payload).length === 0) {
				toastError('No changes to save');
				return;
			}

			await updateUser(editingUser.id, payload);
			toastSuccess('User updated successfully');
			cancelEdit();
			await fetchUsers();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to update user';
			toastError(errorMsg);
			console.error('Error updating user:', err);
		}
	}

	async function handleBulkUpdate() {
		if (selectedUsers.size === 0) return;

		try {
			const payload: any = {
				user_ids: Array.from(selectedUsers)
			};

			if (bulkSetRole) {
				payload.set_role = bulkSetRole;
			}
			if (bulkConfirmEmail) {
				payload.confirm_email = true;
			}
			if (bulkDisableMFA) {
				payload.disable_mfa = true;
			}

			if (Object.keys(payload).length === 1) {
				toastError('Please select at least one action');
				return;
			}

			await bulkUpdateUsers(payload);
			toastSuccess(`Updated ${selectedUsers.size} user(s) successfully`);
			selectedUsers = new Set();
			showBulkActions = false;
			bulkSetRole = '';
			bulkConfirmEmail = false;
			bulkDisableMFA = false;
			await fetchUsers();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to update users';
			toastError(errorMsg);
			console.error('Error bulk updating users:', err);
		}
	}

	async function viewAuditLogs(userId: string) {
		viewingAuditLogsUserId = userId;
		showAuditModal = true;
		auditLogs = [];
		try {
			auditLogs = await getUserAuditLogs(userId, 50);
		} catch (err) {
			toastError('Failed to load audit logs');
			console.error('Error loading audit logs:', err);
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	function formatDateShort(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}

	const totalPages = Math.ceil(total / limit);
</script>

<div class="page-container">
	<!-- Header -->
	<div class="page-header">
		<div>
			<h1 class="page-title">User Management</h1>
			<p class="page-description">Manage users, roles, and permissions</p>
		</div>
		<div class="header-actions">
			<div class="search-box">
				<Search class="search-icon" />
				<input
					type="text"
					placeholder="Search users..."
					bind:value={searchQuery}
					onkeydown={(e) => {
						if (e.key === 'Enter') handleSearch();
					}}
					class="search-input"
				/>
				<button type="button" class="btn btn-primary btn-sm" onclick={handleSearch}>
					Search
				</button>
			</div>
		</div>
	</div>

	<!-- Bulk Actions Panel -->
	{#if showBulkActions}
		<div class="bulk-actions-panel">
			<div class="bulk-actions-info">
				<span>{selectedUsers.size} user(s) selected</span>
			</div>
			<div class="bulk-actions-controls">
				<select bind:value={bulkSetRole} class="input input-sm">
					<option value="">Set Role...</option>
					<option value="user">User</option>
					<option value="admin">Admin</option>
				</select>
				<label class="checkbox-label">
					<input type="checkbox" bind:checked={bulkConfirmEmail} />
					Confirm Email
				</label>
				<label class="checkbox-label">
					<input type="checkbox" bind:checked={bulkDisableMFA} />
					Disable MFA
				</label>
				<button type="button" class="btn btn-primary btn-sm" onclick={handleBulkUpdate}>
					Apply
				</button>
				<button
					type="button"
					class="btn btn-secondary btn-sm"
					onclick={() => {
						selectedUsers = new Set();
						showBulkActions = false;
					}}
				>
					Cancel
				</button>
			</div>
		</div>
	{/if}

	<!-- Error State -->
	{#if error}
		<div class="alert alert-error">
			<AlertCircle class="icon" />
			<p>{error}</p>
		</div>
	{/if}

	<!-- Loading State -->
	{#if loading}
		<div class="loading-container">
			<div class="spinner"></div>
		</div>
	{:else if users.length === 0}
		<!-- Empty State -->
		<div class="empty-state">
			<Users class="empty-icon" />
			<p class="empty-title">No users found</p>
			<p class="empty-description">
				{searchQuery ? 'Try adjusting your search query' : 'No users in the system'}
			</p>
		</div>
	{:else}
		<!-- Users Table -->
		<div class="table-container">
			<table class="table">
				<thead>
					<tr>
						<th class="checkbox-col">
							<input
								type="checkbox"
								checked={selectedUsers.size === users.length && users.length > 0}
								onchange={toggleSelectAll}
							/>
						</th>
						<th>Email</th>
						<th>Role</th>
						<th>Status</th>
						<th>MFA</th>
						<th>Created</th>
						<th>Last Login</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each users as user}
						<tr>
							<td class="checkbox-col">
								<input
									type="checkbox"
									checked={selectedUsers.has(user.id)}
									onchange={() => toggleUserSelection(user.id)}
								/>
							</td>
							<td>
								<div class="user-email">
									<Mail class="icon-sm" />
									<span class="font-medium">{user.email}</span>
								</div>
							</td>
							<td>
								<span class={`role-badge role-${user.role}`}>
									<Shield class="icon-xs" />
									{user.role}
								</span>
							</td>
							<td>
								{#if user.email_confirmed_at}
									<span class="status-badge status-active">
										<Check class="icon-xs" />
										Confirmed
									</span>
								{:else}
									<span class="status-badge status-pending">
										<X class="icon-xs" />
										Unconfirmed
									</span>
								{/if}
							</td>
							<td>
								{#if user.mfa_enabled}
									<span class="status-badge status-active">Enabled</span>
								{:else}
									<span class="status-badge status-inactive">Disabled</span>
								{/if}
							</td>
							<td class="text-muted">{formatDateShort(user.created_at)}</td>
							<td class="text-muted">{formatDateShort(user.last_login_at)}</td>
							<td class="text-right">
								<div class="actions">
									<button
										type="button"
										class="btn btn-sm btn-primary"
										onclick={() => startEdit(user)}
										title="Edit user"
									>
										<Edit class="icon-sm" />
									</button>
									<button
										type="button"
										class="btn btn-sm btn-secondary"
										onclick={() => viewAuditLogs(user.id)}
										title="View audit logs"
									>
										<MoreVertical class="icon-sm" />
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="pagination">
				<button
					type="button"
					class="btn btn-secondary btn-sm"
					disabled={currentPage === 1}
					onclick={() => handlePageChange(currentPage - 1)}
				>
					<ChevronLeft class="icon-sm" />
					Previous
				</button>
				<span class="pagination-info">
					Page {currentPage} of {totalPages} ({total} total)
				</span>
				<button
					type="button"
					class="btn btn-secondary btn-sm"
					disabled={currentPage === totalPages}
					onclick={() => handlePageChange(currentPage + 1)}
				>
					Next
					<ChevronRight class="icon-sm" />
				</button>
			</div>
		{/if}
	{/if}
</div>

<!-- Edit User Modal -->
{#if showEditModal && editingUser}
	<div class="modal-overlay" onclick={cancelEdit}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Edit User</h2>
			<div class="modal-content">
				<div class="form-group">
					<label class="form-label">Email</label>
					<input type="email" bind:value={editEmail} class="input" />
				</div>
				<div class="form-group">
					<label class="form-label">Role</label>
					<select bind:value={editRole} class="input">
						<option value="user">User</option>
						<option value="admin">Admin</option>
					</select>
				</div>
				<div class="form-group">
					<label class="form-label">Phone Number</label>
					<input type="tel" bind:value={editPhoneNumber} class="input" placeholder="+1234567890" />
				</div>
				<div class="form-group">
					<label class="form-label">Preferred Locale</label>
					<input type="text" bind:value={editPreferredLocale} class="input" placeholder="en" />
				</div>
				<div class="form-group">
					<label class="checkbox-label">
						<input type="checkbox" bind:checked={editConfirmEmail} />
						Confirm Email Address
					</label>
				</div>
				<div class="form-group">
					<label class="checkbox-label">
						<input type="checkbox" bind:checked={editDisableMFA} />
						Disable MFA
					</label>
				</div>
				<div class="modal-actions">
					<button type="button" class="btn btn-primary" onclick={handleUpdate}>Save</button>
					<button type="button" class="btn btn-secondary" onclick={cancelEdit}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Audit Logs Modal -->
{#if showAuditModal && viewingAuditLogsUserId}
	<div class="modal-overlay" onclick={() => (showAuditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2 class="modal-title">Audit Logs</h2>
			<div class="modal-content">
				{#if auditLogs.length === 0}
					<p class="text-muted">No audit logs found</p>
				{:else}
					<div class="audit-logs-list">
						{#each auditLogs as log}
							<div class="audit-log-item">
								<div class="audit-log-header">
									<span class="audit-log-action">{log.action}</span>
									<span class="audit-log-time">{formatDate(log.created_at)}</span>
								</div>
								<div class="audit-log-details">
									{#if log.resource_type}
										<span class="audit-log-resource">{log.resource_type}</span>
									{/if}
									{#if log.ip_address}
										<span class="audit-log-ip">IP: {log.ip_address}</span>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				{/if}
				<div class="modal-actions">
					<button
						type="button"
						class="btn btn-secondary"
						onclick={() => (showAuditModal = false)}
					>
						Close
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.page-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.page-title {
		font-size: 1.875rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 0.5rem;
	}

	.page-description {
		color: rgba(148, 163, 184, 0.9);
		font-size: 0.9rem;
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.search-box {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background: rgba(15, 15, 15, 0.6);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.5rem;
		padding: 0.5rem 0.75rem;
	}

	.search-icon {
		width: 1rem;
		height: 1rem;
		color: rgba(148, 163, 184, 0.8);
	}

	.search-input {
		border: none;
		background: transparent;
		color: #f8fafc;
		font-size: 0.875rem;
		width: 200px;
		outline: none;
	}

	.bulk-actions-panel {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: 0.5rem;
		flex-wrap: wrap;
	}

	.bulk-actions-info {
		color: #d4d4d4;
		font-weight: 500;
	}

	.bulk-actions-controls {
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #cbd5e1;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		border: 1px solid;
		transition: all 120ms ease;
		cursor: pointer;
	}

	.btn-primary {
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.12);
		color: #d4d4d4;
	}

	.btn-primary:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.12);
		border-color: rgba(255, 255, 255, 0.2);
	}

	.btn-secondary {
		background: rgba(71, 85, 105, 0.15);
		border-color: rgba(255, 255, 255, 0.08);
		color: #cbd5e1;
	}

	.btn-secondary:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.08);
		border-color: rgba(255, 255, 255, 0.12);
	}

	.btn-sm {
		padding: 0.375rem 0.75rem;
		font-size: 0.8rem;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.icon {
		width: 1rem;
		height: 1rem;
	}

	.icon-sm {
		width: 0.875rem;
		height: 0.875rem;
	}

	.icon-xs {
		width: 0.75rem;
		height: 0.75rem;
	}

	.loading-container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 4rem 0;
	}

	.spinner {
		width: 3rem;
		height: 3rem;
		border: 2px solid rgba(255, 255, 255, 0.06);
		border-top-color: #737373;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
	}

	.empty-icon {
		width: 4rem;
		height: 4rem;
		color: rgba(255, 255, 255, 0.12);
		margin: 0 auto 1rem;
	}

	.empty-title {
		font-size: 1.125rem;
		font-weight: 600;
		color: rgba(203, 213, 225, 0.9);
		margin-bottom: 0.5rem;
	}

	.empty-description {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.table-container {
		background: rgba(15, 15, 15, 0.4);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.75rem;
		overflow: hidden;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
	}

	.table thead {
		background: rgba(15, 15, 15, 0.6);
	}

	.table th {
		padding: 1rem 1.5rem;
		text-align: left;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: rgba(148, 163, 184, 0.9);
	}

	.table td {
		padding: 1rem 1.5rem;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
	}

	.table tbody tr:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.checkbox-col {
		width: 3rem;
		text-align: center;
	}

	.text-right {
		text-align: right;
	}

	.text-muted {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.user-email {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.font-medium {
		font-weight: 500;
	}

	.role-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.role-user {
		background: rgba(255, 255, 255, 0.05);
		color: #cbd5e1;
	}

	.role-admin {
		background: rgba(239, 68, 68, 0.2);
		color: #fca5a5;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-active {
		background: rgba(34, 197, 94, 0.2);
		color: #86efac;
	}

	.status-pending {
		background: rgba(251, 191, 36, 0.2);
		color: #fde047;
	}

	.status-inactive {
		background: rgba(255, 255, 255, 0.05);
		color: #cbd5e1;
	}

	.actions {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		gap: 0.5rem;
	}

	.pagination {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem;
	}

	.pagination-info {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.875rem;
	}

	.alert {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 0.5rem;
		border: 1px solid;
	}

	.alert-error {
		background: rgba(239, 68, 68, 0.1);
		border-color: rgba(239, 68, 68, 0.3);
		color: #fca5a5;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.75);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
		padding: 1rem;
	}

	.modal {
		background: rgba(15, 15, 15, 0.98);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.75rem;
		padding: 1.5rem;
		width: 100%;
		max-width: 28rem;
		box-shadow: 0 20px 45px rgba(0, 0, 0, 0.8);
	}

	.modal-large {
		max-width: 42rem;
	}

	.modal-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #f8fafc;
		margin-bottom: 1rem;
	}

	.modal-content {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: rgba(203, 213, 225, 0.9);
	}

	.input {
		width: 100%;
		padding: 0.5rem 0.75rem;
		background: rgba(15, 15, 15, 0.6);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 0.5rem;
		color: #f8fafc;
		font-size: 0.875rem;
	}

	.input:focus {
		outline: none;
		border-color: rgba(255, 255, 255, 0.2);
		box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.05);
	}

	.input-sm {
		padding: 0.375rem 0.5rem;
		font-size: 0.8rem;
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}

	.audit-logs-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		max-height: 400px;
		overflow-y: auto;
	}

	.audit-log-item {
		padding: 0.75rem;
		background: rgba(15, 15, 15, 0.4);
		border: 1px solid rgba(255, 255, 255, 0.06);
		border-radius: 0.5rem;
	}

	.audit-log-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.audit-log-action {
		font-weight: 600;
		color: #f8fafc;
		font-size: 0.875rem;
	}

	.audit-log-time {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.75rem;
	}

	.audit-log-details {
		display: flex;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.audit-log-resource,
	.audit-log-ip {
		color: rgba(148, 163, 184, 0.8);
		font-size: 0.75rem;
	}
</style>

