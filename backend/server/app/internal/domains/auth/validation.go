package auth

import (
	"fmt"
	"strings"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateRegisterPayload validates registration payload
func ValidateRegisterPayload(payload *registerPayload) error {
	// Validate email
	if err := validation.ValidateEmail(payload.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}

	// Validate password (min 8 chars, max 128 chars)
	if err := validation.ValidateString(payload.Password, 8, 128, "password"); err != nil {
		return fmt.Errorf("password: %w", err)
	}

	// Validate display name (optional, but if provided, validate)
	if payload.DisplayName != "" {
		if err := validation.ValidateString(payload.DisplayName, 1, 100, "display_name"); err != nil {
			return fmt.Errorf("display_name: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.DisplayName); err != nil {
			return fmt.Errorf("display_name: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.DisplayName); err != nil {
			return fmt.Errorf("display_name: %w", err)
		}
	}

	// Validate locale (optional, but if provided, validate format)
	if payload.Locale != "" {
		if err := validation.ValidateString(payload.Locale, 2, 10, "locale"); err != nil {
			return fmt.Errorf("locale: %w", err)
		}
	}

	return nil
}

// ValidateLoginPayload validates login payload
func ValidateLoginPayload(payload *loginPayload) error {
	// Validate email
	if err := validation.ValidateEmail(payload.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}

	// Validate password (required, but no length check for login - service handles it)
	if payload.Password == "" {
		return fmt.Errorf("password is required")
	}

	// Validate device fingerprint (optional, but if provided, validate length)
	if payload.DeviceFingerprint != "" {
		if err := validation.ValidateString(payload.DeviceFingerprint, 1, 255, "device_fingerprint"); err != nil {
			return fmt.Errorf("device_fingerprint: %w", err)
		}
	}

	// Validate device name (optional, but if provided, validate)
	if payload.DeviceName != "" {
		if err := validation.ValidateString(payload.DeviceName, 1, 100, "device_name"); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.DeviceName); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.DeviceName); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
	}

	// Validate MFA code (optional, but if provided, validate format - typically 6 digits)
	if payload.MFACode != "" {
		if err := validation.ValidateString(payload.MFACode, 6, 6, "mfa_code"); err != nil {
			return fmt.Errorf("mfa_code: %w", err)
		}
	}

	return nil
}

// ValidateResendPayload validates resend confirmation payload
func ValidateResendPayload(payload *resendPayload) error {
	if err := validation.ValidateEmail(payload.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}
	return nil
}

// ValidateRequestResetPayload validates password reset request payload
func ValidateRequestResetPayload(payload *requestResetPayload) error {
	if err := validation.ValidateEmail(payload.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}
	return nil
}

// ValidateResetConfirmPayload validates password reset confirmation payload
func ValidateResetConfirmPayload(payload *resetConfirmPayload) error {
	// Validate token (required)
	if payload.Token == "" {
		return fmt.Errorf("token is required")
	}

	// Validate password (min 8 chars, max 128 chars)
	if err := validation.ValidateString(payload.Password, 8, 128, "password"); err != nil {
		return fmt.Errorf("password: %w", err)
	}

	return nil
}

// ValidateRefreshPayload validates refresh token payload
func ValidateRefreshPayload(payload *refreshPayload) error {
	// Validate refresh token (required)
	if payload.RefreshToken == "" {
		return fmt.Errorf("refresh_token is required")
	}

	return nil
}

// ValidateLogoutPayload validates logout payload
func ValidateLogoutPayload(payload *logoutPayload) error {
	// Session ID is optional (if not provided, logout all sessions)
	if payload.SessionID != "" {
		if err := validation.ValidateUUID(payload.SessionID); err != nil {
			return fmt.Errorf("session_id: %w", err)
		}
	}
	return nil
}

// ValidateRevokeSessionsPayload validates revoke sessions payload
func ValidateRevokeSessionsPayload(payload *revokeSessionsPayload) error {
	// Keep session ID is optional, but if provided, validate UUID
	if payload.KeepSessionID != "" {
		if err := validation.ValidateUUID(payload.KeepSessionID); err != nil {
			return fmt.Errorf("keep_session_id: %w", err)
		}
	}
	return nil
}

// ValidateEnableMFAPayload validates enable MFA payload
func ValidateEnableMFAPayload(payload *enableMFAPayload) error {
	// Validate issuer (optional, but if provided, validate)
	if payload.Issuer != "" {
		if err := validation.ValidateString(payload.Issuer, 1, 100, "issuer"); err != nil {
			return fmt.Errorf("issuer: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Issuer); err != nil {
			return fmt.Errorf("issuer: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Issuer); err != nil {
			return fmt.Errorf("issuer: %w", err)
		}
	}

	// Validate label (optional, but if provided, validate)
	if payload.Label != "" {
		if err := validation.ValidateString(payload.Label, 1, 100, "label"); err != nil {
			return fmt.Errorf("label: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Label); err != nil {
			return fmt.Errorf("label: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Label); err != nil {
			return fmt.Errorf("label: %w", err)
		}
	}

	// Validate code (optional, but if provided, validate format - typically 6 digits)
	if payload.Code != "" {
		if err := validation.ValidateString(payload.Code, 6, 6, "code"); err != nil {
			return fmt.Errorf("code: %w", err)
		}
	}

	return nil
}

// ValidateVerifyMFAPayload validates verify MFA payload
func ValidateVerifyMFAPayload(payload *verifyMFAPayload) error {
	// Validate code (required, 6 digits)
	if err := validation.ValidateString(payload.Code, 6, 6, "code"); err != nil {
		return fmt.Errorf("code: %w", err)
	}
	return nil
}

// ValidateOAuthStartPayload validates OAuth start payload
func ValidateOAuthStartPayload(payload *oauthStartPayload) error {
	// Validate provider (required)
	if payload.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if err := validation.ValidateString(payload.Provider, 1, 50, "provider"); err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	// Validate mode (optional, but if provided, validate)
	if payload.Mode != "" {
		if err := validation.ValidateString(payload.Mode, 1, 50, "mode"); err != nil {
			return fmt.Errorf("mode: %w", err)
		}
	}

	// Validate redirect origin (optional, but if provided, validate URL)
	if payload.RedirectOrigin != "" {
		if err := validation.ValidateURL(payload.RedirectOrigin); err != nil {
			return fmt.Errorf("redirect_origin: %w", err)
		}
	}

	// Validate device fingerprint (optional)
	if payload.DeviceFingerprint != "" {
		if err := validation.ValidateString(payload.DeviceFingerprint, 1, 255, "device_fingerprint"); err != nil {
			return fmt.Errorf("device_fingerprint: %w", err)
		}
	}

	// Validate device name (optional)
	if payload.DeviceName != "" {
		if err := validation.ValidateString(payload.DeviceName, 1, 100, "device_name"); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.DeviceName); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.DeviceName); err != nil {
			return fmt.Errorf("device_name: %w", err)
		}
	}

	return nil
}

// ValidateUpdateProfilePayload validates update profile payload
func ValidateUpdateProfilePayload(payload *updateProfilePayload) error {
	// Validate phone number (optional, but if provided, validate format)
	if payload.PhoneNumber != nil && *payload.PhoneNumber != "" {
		phone := *payload.PhoneNumber
		// Basic phone validation (can be enhanced with regex)
		if len(phone) > 20 {
			return fmt.Errorf("phone_number: too long (maximum 20 characters)")
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(phone); err != nil {
			return fmt.Errorf("phone_number: %w", err)
		}
		if err := validation.ValidateNoXSS(phone); err != nil {
			return fmt.Errorf("phone_number: %w", err)
		}
	}

	// Validate preferred locale (optional, but if provided, validate format)
	if payload.PreferredLocale != nil && *payload.PreferredLocale != "" {
		locale := *payload.PreferredLocale
		if err := validation.ValidateString(locale, 2, 10, "preferred_locale"); err != nil {
			return fmt.Errorf("preferred_locale: %w", err)
		}
	}

	return nil
}

// ValidateAdminUpdateUserPayload validates admin update user payload
func ValidateAdminUpdateUserPayload(payload *adminUpdateUserPayload) error {
	// Validate role (optional, but if provided, validate)
	if payload.SetRole != nil && *payload.SetRole != "" {
		role := *payload.SetRole
		validRoles := []string{"user", "admin", "moderator"}
		isValid := false
		for _, validRole := range validRoles {
			if strings.ToLower(role) == validRole {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("set_role: must be one of: user, admin, moderator")
		}
	}

	// Validate email (optional, but if provided, validate format)
	if payload.SetEmail != nil && *payload.SetEmail != "" {
		if err := validation.ValidateEmail(*payload.SetEmail); err != nil {
			return fmt.Errorf("set_email: %w", err)
		}
	}

	// Validate phone number (optional, but if provided, validate)
	if payload.SetPhoneNumber != nil && *payload.SetPhoneNumber != "" {
		phone := *payload.SetPhoneNumber
		if len(phone) > 20 {
			return fmt.Errorf("set_phone_number: too long (maximum 20 characters)")
		}
		if err := validation.ValidateNoSQLInjection(phone); err != nil {
			return fmt.Errorf("set_phone_number: %w", err)
		}
		if err := validation.ValidateNoXSS(phone); err != nil {
			return fmt.Errorf("set_phone_number: %w", err)
		}
	}

	// Validate preferred locale (optional, but if provided, validate)
	if payload.SetPreferredLocale != nil && *payload.SetPreferredLocale != "" {
		locale := *payload.SetPreferredLocale
		if err := validation.ValidateString(locale, 2, 10, "set_preferred_locale"); err != nil {
			return fmt.Errorf("set_preferred_locale: %w", err)
		}
	}

	return nil
}

// ValidateAdminBulkUpdateUsersPayload validates admin bulk update users payload
func ValidateAdminBulkUpdateUsersPayload(payload *adminBulkUpdateUsersPayload) error {
	// Validate user IDs (required, at least one)
	if len(payload.UserIDs) == 0 {
		return fmt.Errorf("user_ids: at least one user ID is required")
	}
	if len(payload.UserIDs) > 100 {
		return fmt.Errorf("user_ids: too many user IDs (maximum 100)")
	}

	// Validate each user ID
	for i, idStr := range payload.UserIDs {
		if err := validation.ValidateUUID(idStr); err != nil {
			return fmt.Errorf("user_ids[%d]: %w", i, err)
		}
	}

	// Validate role (optional, but if provided, validate)
	if payload.SetRole != nil && *payload.SetRole != "" {
		role := *payload.SetRole
		validRoles := []string{"user", "admin", "moderator"}
		isValid := false
		for _, validRole := range validRoles {
			if strings.ToLower(role) == validRole {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("set_role: must be one of: user, admin, moderator")
		}
	}

	return nil
}

// ValidateListUsersQueryParams validates query parameters for ListUsers
func ValidateListUsersQueryParams(limit, offset int, search string) error {
	// Validate limit
	if limit < 1 {
		return fmt.Errorf("limit: must be at least 1")
	}
	if limit > 100 {
		return fmt.Errorf("limit: must be at most 100")
	}

	// Validate offset
	if offset < 0 {
		return fmt.Errorf("offset: must be at least 0")
	}

	// Validate search (optional, but if provided, validate length)
	if search != "" {
		if err := validation.ValidateString(search, 1, 200, "search"); err != nil {
			return fmt.Errorf("search: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(search); err != nil {
			return fmt.Errorf("search: %w", err)
		}
		if err := validation.ValidateNoXSS(search); err != nil {
			return fmt.Errorf("search: %w", err)
		}
	}

	return nil
}

// ValidateGetUserAuditLogsQueryParams validates query parameters for GetUserAuditLogs
func ValidateGetUserAuditLogsQueryParams(limit, offset int) error {
	// Validate limit
	if limit < 1 {
		return fmt.Errorf("limit: must be at least 1")
	}
	if limit > 200 {
		return fmt.Errorf("limit: must be at most 200")
	}

	// Validate offset
	if offset < 0 {
		return fmt.Errorf("offset: must be at least 0")
	}

	return nil
}
