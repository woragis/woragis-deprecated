# Woragis Auth Backend Overview

This document highlights the backend auth workflows implemented in November 2025, focusing on email confirmations, session management, and MFA support.

## User Lifecycle

- **Registration** – `POST /api/auth/register`
  - Creates a user record and dispatches a confirmation email.
  - Response returns `status: confirmation_required`; access tokens are no longer issued automatically on registration.
- **Email Confirmation** – `POST /api/auth/confirm`
  - Accepts a confirmation token and marks `email_confirmed_at`.
  - `POST /api/auth/confirm/resend` can resend tokens for unconfirmed users.
- **Login** – `POST /api/auth/login`
  - Verifies credentials, enforces confirmation & MFA policies, and issues:
    - Access token (JWT)
    - Refresh token (per-device session, stored hashed)
    - Session identifier for managing active sessions.
- **Session Refresh** – `POST /api/auth/refresh`
  - Rotates refresh tokens and extends session expiration.
- **Logout & Session Management**
  - `POST /api/auth/logout` revokes a single session.
  - `GET /api/auth/sessions` lists active sessions.
  - `POST /api/auth/sessions/revoke` revokes all other sessions (optional keep).

## Email Templates & SMTP

- Transactional emails render via `app/internal/services/email/templates`.
  - Templates: confirmation, password reset, welcome, session alert.
  - Renderer falls back to `en` locale if user preference is missing.
- SMTP senders now create multipart/alternative emails (text + HTML).
- Redis email worker accepts structured payloads with `text_message` and `html_message`.

## Multi-Factor Authentication (MFA)

- TOTP enrolment supported through:
  - `POST /api/auth/mfa/enable`
  - `POST /api/auth/mfa/verify`
  - `POST /api/auth/mfa/disable`
- Backup codes are generated when enabling MFA.
- Login requires `mfa_code` when `mfa_enabled` is true.

## Domain Models

- `User` entity extended with:
  - `email_confirmed_at`, `last_login_at`, password & MFA metadata.
- New aggregates:
  - `Session` (per refresh token/device)
  - `Device`
  - `MFAToken`
  - `AuditLog`
  - `OAuthAccount`
  - `EmailToken`
- Repository now supports CRUD for sessions, devices, MFA tokens, audit logs, OAuth accounts, and email tokens.

## Configuration Notes

Existing env vars still apply:

- `APP_PUBLIC_URL` (used for email links)
- SMTP variables (`SMTP_HOST`, `SMTP_FROM`, etc.)

When enabling MFA in production, ensure time synchronisation across servers.

## Next Steps / Frontend Alignment

- Update frontend flows to handle registration confirmation, session management, and MFA prompts.
- Provide user settings UI for device/session review and MFA toggles.
- Instrument new audit log endpoints once exposed.


