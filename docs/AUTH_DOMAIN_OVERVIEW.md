# Woragis Authentication Domain Overview

This document outlines the authentication capabilities that now exist across the Woragis backend and frontend, and explains how the two layers interact. It is intended as the starting point for anyone touching auth code or trying to extend the domain.

---

## 1. Backend Features (Go / Fiber)

### 1.1 Core REST Endpoints

| Route | Method | Purpose |
|-------|--------|---------|
| `/api/auth/register` | `POST` | Registers a user, enforces password rules, schedules confirmation email |
| `/api/auth/confirm` | `POST` | Confirms email ownership with a token |
| `/api/auth/confirm/resend` | `POST` | Sends a fresh confirmation email if the user is still unverified |
| `/api/auth/login` | `POST` | Validates credentials, MFA, device metadata, issues access/refresh tokens |
| `/api/auth/refresh` | `POST` | Rotates refresh tokens and issues a new JWT |
| `/api/auth/logout` | `POST` | Revokes a single session by ID |
| `/api/auth/sessions` | `GET` | Lists active sessions (requires auth middleware) |
| `/api/auth/sessions/revoke` | `POST` | Revokes all sessions except an optional “keep” ID |
| `/api/auth/password/reset/request` | `POST` | Generates a one-time reset token and emails it |
| `/api/auth/password/reset/confirm` | `POST` | Applies a new password using the reset token |
| `/api/auth/mfa/enable` | `POST` | Generates a TOTP secret and backup codes for a user |
| `/api/auth/mfa/verify` | `POST` | Confirms MFA enrolment with a valid TOTP code |
| `/api/auth/mfa/disable` | `POST` | Removes TOTP information from the account |
| `/api/auth/oauth/providers` | `GET` | Exposes the configured OAuth providers (Google/GitHub/Microsoft) |
| `/api/auth/oauth/start` | `POST` | Produces an authorization URL + state token for an OAuth provider |
| `/api/auth/oauth/callback/:provider` | `GET` | Handles provider redirect, finalizes login or linking |
| `/api/auth/oauth/accounts` | `GET` | Lists linked OAuth providers for the authenticated user |
| `/api/auth/oauth/accounts/:provider` | `DELETE` | Unlinks an OAuth provider |

### 1.2 Domain Highlights

- **Session Management**: access token (JWT) + refresh token pair with per-device session records. Device metadata (fingerprint, friendly name, user agent, IP) is stored.
- **Email Confirmation**: registration is always staged; email must be confirmed before login works. Tokens are issued via Redis-backed store.
- **Password Reset**: one-time tokens with TTL issued via email; tokens are invalidated upon use.
- **OAuth 2.0**:
  - Supports Google, GitHub, Microsoft via `golang.org/x/oauth2`.
  - Handles login (create account if needed), linking, token refresh metadata, audit logging.
  - State tokens are stored in-memory with TTL (default 10 minutes).
- **MFA (TOTP)**: uses `github.com/pquerna/otp/totp` to generate secrets, backup codes, and enforce verification during login.
- **Audit Logging**: key actions (login, OAuth linking, MFA enable, etc.) are stored via the existing audit repository.
- **Configuration**:
  - `.env` variables for OAuth clients, redirect URLs, allowed scopes.
  - `APP_PUBLIC_URL` drives callback defaults and templated emails.
  - `AUTH_JWT_SECRET` / `AUTH_JWT_TTL` for token issuance.
- **Middlewares**: JWT middleware attaches `auth.user_id` / `auth.user_email` to the Fiber context, enabling the protected routes.

### 1.3 Important Packages

- `app/internal/domains/auth/service.go`: core business rules.
- `service_oauth.go`: OAuth-specific orchestration (state, token exchange, user provisioning).
- `handler.go`: HTTP binding + normalization of responses/errors.
- `repository.go`: persistence (users, sessions, devices, MFA tokens, OAuth accounts).
- `token_store.go`: Redis-backed token storage for password resets.
- `jwt_manager.go`: JWT signing/verification.
- `routes.go`: public + protected route registration.
- `pkg/config/oauth.go`: environment-driven provider configuration.

---

## 2. Frontend Features (Svelte + Axios)

### 2.1 Visible Routes

| Route | Purpose |
|-------|---------|
| `/auth/login` | Email/password login, OAuth providers, forgot password links |
| `/auth/register` | Registration form; routes to confirmation if required |
| `/auth/confirm` | Token submission + resend confirmation email |
| `/auth/forgot` | Password reset request |
| `/auth/reset` | Password reset confirmation (token + new password) |
| `/auth/mfa` | Enrol/disable MFA, show secrets and backup codes |
| `/auth/sessions` | Display active sessions, revoke others/devices |
| `/auth/connections` | Manage linked OAuth providers |

### 2.2 Stores & Helpers

- `authStore` (`src/lib/stores/auth.ts`):
  - Persists `{user, token, refreshToken, sessionId}` to cookies.
  - Normalizes camel_case/snake_case payload mismatches.
  - Provides `setSession`, `updateTokens`, `updateUser`, `clear`.
- `auth/lifecycle.ts`:
  - Refreshes access tokens (default every 9 minutes) using refresh token.
  - Clears state on 401 responses.
- `device.ts`:
  - Generates and persists a browser fingerprint.
  - Returns vendor/platform-based device name and user agent string.
- `api/auth.ts`:
  - Axios wrappers for every auth endpoint (login, refresh, logout, sessions, MFA, OAuth, confirmation, password reset).
- `api/client.ts`:
  - Adds bearer tokens to requests, clears auth store on 401.

### 2.3 UX Behaviors

- **Login**: collects device metadata, drives OAuth winners via popup + `postMessage`.
- **Registration**: detects `confirmation_required` response and redirects to `/auth/confirm` while caching the email in `localStorage`.
- **Lifecycle**: top-level layout triggers auto-refresh and exposes a sign-out button.
- **Sessions View**: surfaces server data, supports revoking others or individual sessions.
- **MFA UI**: displays secrets, provisioning URIs, backup codes; verifies code entry before marking MFA enabled.
- **Password Reset**: two-step flow that matches backend semantics (request + confirm).
- **Confirm Email**: handles inbound tokens via query params and provides resend forms.
- **OAuth Connections**: friendly linking/unlinking with provider-specific error feedback.

---

## 3. Backend <-> Frontend Integration Notes

### 3.1 Token Handling

- **Access Tokens (JWT)**: stored in cookies via `authStore`; included in `Authorization: Bearer` headers for non-bypassed routes.
- **Refresh Tokens**: stored in cookies; auto-rotated by `auth/lifecycle`, ensuring the client always uses the latest refresh token.
- **Session IDs**: retained so the frontend can call `/auth/logout`, `/auth/sessions/revoke`, and highlight “This device” in the sessions table.

### 3.2 Device Metadata

- Login/OAuth flows provide `device_fingerprint`, `device_name`, and `user_agent` — enabling backend session insights.
- The fingerprint (UUID) is persisted in `localStorage` to remain stable across refreshes.

### 3.3 OAuth Popups

- Frontend uses `window.open()` for provider authorization.
- Backend callback returns an HTML page that `postMessage`s the outcome (`oauth:result`) with login payload or error.
- Frontend listens for messages from `API_BASE_URL` origin and either stores credentials or surfaces errors.

### 3.4 MFA & Confirmation Coupling

- Backend blocks login until email is confirmed and MFA codes (if required) are verified.
- Frontend surfaces appropriate UI combinations:
  - Registration -> Confirmation route -> Login.
  - MFA prompts during login (provided by backend response codes; future enhancement could inspect error codes to open TOTP prompt).

### 3.5 Password Reset

- Backend email contains reset token; user lands on `/auth/reset?token=...`.
- Frontend simply posts token/password pairs back to `/auth/password/reset/confirm`.

---

## 4. Environment & Configuration Checklist

| Variable | Description |
|----------|-------------|
| `APP_PUBLIC_URL` | Used for OAuth redirect defaults, email links, frontend configuration |
| `AUTH_JWT_SECRET`, `AUTH_JWT_TTL` | JWT signing secret and TTL |
| `GOOGLE_OAUTH_CLIENT_ID` / `_SECRET` | Enable Google OAuth |
| `GOOGLE_OAUTH_REDIRECT_URL` | Optional override; defaults to `${APP_PUBLIC_URL}/api/auth/oauth/callback/google` |
| `GOOGLE_OAUTH_SCOPES` | Optional; defaults to `openid email profile` |
| `GITHUB_OAUTH_CLIENT_ID` / `_SECRET` | Enable GitHub OAuth |
| `MICROSOFT_OAUTH_CLIENT_ID` / `_SECRET` | Enable Microsoft OAuth |
| SMTP variables | Required for email confirmation + password reset flows |
| Redis URL | Required by `RedisTokenStore` for password reset tokens |

On the frontend, define `PUBLIC_API_BASE_URL` in `frontend/.env` so that Axios hits the correct API host.

---

## 5. Testing Checklist

1. **Registration**: ensure `/auth/register` returns `confirmation_required`, email arrives, `/auth/confirm` accepts token.
2. **Login**: verify login yields session + cookies; check session view reflects device metadata.
3. **Token Refresh**: wait for auto-refresh or call `refreshSessionNow()` in console; confirm new access token.
4. **Logout**: click “Sign out”; ensure cookies cleared and server session revoked.
5. **Password Reset**: request + confirm flows should update password and log audit entries.
6. **MFA**: run through enable → verify → login (with code) → disable.
7. **OAuth**: start link/login flows, revoke provider, re-link. Validate audit logs and session events.
8. **Sessions**: open another browser, log in, then revoke from primary device to check enforcement.

---

## 6. Future Enhancements

- Enforce MFA prompts in the login UI based on backend error codes (`ErrCodeMFARequired`).
- Add QR code generation on the frontend (currently we surface raw provisioning URI/secret).
- Expose audit logs to the frontend for better transparency.
- Implement refresh token reuse detection (refresh session invalidation in case of reuse).
- Provide scaffolding for WebAuthn or additional OAuth providers using the same service hooks.

---

**Maintainers:** @auth-domain contributors  
**Last updated:** November 2025


