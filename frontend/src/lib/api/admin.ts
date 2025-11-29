import { apiClient } from '@clients/apiClient';

export interface AdminUser {
	id: string;
	email: string;
	created_at: string;
	updated_at: string;
	email_confirmed_at?: string;
	last_login_at?: string;
	role: string;
	mfa_enabled: boolean;
	preferred_locale: string;
	phone_number?: string;
}

export interface AdminUserListResponse {
	users: AdminUser[];
	total: number;
	limit: number;
	offset: number;
}

export interface AdminUserListParams {
	limit?: number;
	offset?: number;
	search?: string;
}

export interface AdminUpdateUserPayload {
	set_role?: string;
	set_email?: string;
	confirm_email?: boolean;
	disable_mfa?: boolean;
	set_phone_number?: string;
	set_preferred_locale?: string;
}

export interface AdminBulkUpdateUsersPayload {
	user_ids: string[];
	set_role?: string;
	confirm_email?: boolean;
	disable_mfa?: boolean;
}

export interface AuditLog {
	id: string;
	user_id: string;
	action: string;
	resource_type: string;
	resource_id?: string;
	ip_address?: string;
	user_agent?: string;
	metadata?: Record<string, any>;
	created_at: string;
}

export interface AuditLogsResponse {
	audit_logs: AuditLog[];
}

// List users with pagination and search
export async function listUsers(params?: AdminUserListParams): Promise<AdminUserListResponse> {
	const queryParams = new URLSearchParams();
	if (params?.limit) queryParams.append('limit', params.limit.toString());
	if (params?.offset) queryParams.append('offset', params.offset.toString());
	if (params?.search) queryParams.append('search', params.search);

	const queryString = queryParams.toString();
	const url = `/admin/users${queryString ? `?${queryString}` : ''}`;

	const response = await apiClient.get<{ data: AdminUserListResponse }>(url);
	return response.data.data;
}

// Get user by ID
export async function getUser(id: string): Promise<AdminUser> {
	const response = await apiClient.get<{ data: AdminUser }>(`/admin/users/${id}`);
	return response.data.data;
}

// Update user
export async function updateUser(id: string, payload: AdminUpdateUserPayload): Promise<AdminUser> {
	const response = await apiClient.patch<{ data: AdminUser }>(`/admin/users/${id}`, payload);
	return response.data.data;
}

// Bulk update users
export async function bulkUpdateUsers(payload: AdminBulkUpdateUsersPayload): Promise<{ status: string; count: number }> {
	const response = await apiClient.post<{ data: { status: string; count: number } }>(
		'/admin/users/bulk-update',
		payload
	);
	return response.data.data;
}

// Get user audit logs
export async function getUserAuditLogs(
	userId: string,
	limit?: number
): Promise<AuditLog[]> {
	const queryParams = new URLSearchParams();
	if (limit) queryParams.append('limit', limit.toString());

	const queryString = queryParams.toString();
	const url = `/admin/users/${userId}/audit-logs${queryString ? `?${queryString}` : ''}`;

	const response = await apiClient.get<{ data: AuditLogsResponse }>(url);
	return response.data.data.audit_logs;
}

