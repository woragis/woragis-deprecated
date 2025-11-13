# Woragis Backend TODO

## Domains Requiring Deeper CRUD / Advanced Workflows
- **Auth**
  - [x] Implement session/device management, multi-factor tokens, audit trails, OAuth provider links, bulk user admin actions.
  - [x] Build email confirmation workflow (token issuance, expiry handling, resend limits, transactional templates).
  - [x] Integrate SMTP provider or notification service for confirmation, password reset, and magic link emails.
  - [x] Create HTML/email templates for account confirmation, password reset, welcome, and session notifications with responsive design.
  - [x] Implement template rendering pipeline (layout partials, localization support, preview tooling).
  - [x] Harden domain services for `RegisterUser`, `ConfirmEmail`, `Login`, `Logout`, `RefreshSession`, `RequestPasswordReset`, `CompletePasswordReset`.
  - [x] Ensure domain validation rules (password policy, unique email, unconfirmed account login handling) emit typed errors and domain events.
  - [x] Document auth use cases with sequence diagrams and acceptance criteria for QA.
  - [ ] Align frontend flows with new auth APIs (confirmation, session management, MFA) and update QA scripts.
- **Finances**
  - [x] Extend reporting (cash-flow projections, tagging).
  - [x] Multi-currency normalization pipeline.
  - [x] Recurring schedule templates for transactions.
  - [x] Bulk upload from CSV/OFX.
  - [ ] Additional forecasting analytics (scenario planning, variance alerts).
- **Languages**
  - [ ] Add spaced-repetition scheduling, bulk vocabulary import/export, proficiency analytics, AI-generated practice sets.
- **Projects**
  - [x] Implement kanban workflows.
  - [x] Build dependency graphs.
  - [x] Bulk milestone updates.
  - [x] Templated project duplication.
- **Chats**
  - [ ] Conversation search, bulk archive/delete, shared transcripts, agent assignment history, streaming responses with websockets.
- **Ideas**
  - [ ] Version history, bulk node operations, collaborative editing, advanced filters for relationships.
- **Scheduler**
  - [ ] Complex recurrence rules (cron/rrule), bulk activation/deactivation, execution history dashboard, alerting workflows.
- **Reports**
  - [ ] Custom report builder, scheduling templates, multi-channel delivery, bulk regeneration and export management.
- **Monitoring**
  - [ ] Persist structured events in production DB, alert threshold management, Grafana dashboard provisioning.

## Observability / Platform Enhancements
- Add **WebSocket/Server-Sent Events** pipeline for live metrics and chat streaming.
- Evaluate **gRPC** endpoint layer for Flutter clients; define protobuf schema and gateway.
- Introduce centralized tracing (OpenTelemetry) and log aggregation for cross-service diagnostics.

## Infrastructure Follow-ups
- Dedicated monitoring database connection string for production.
- CI pipeline to build/push Docker images and validate Grafana/Prometheus configs.
- Secrets management for SMTP, AI providers, and future gRPC credentials.
- Comprehensive automated tests (unit/integration) covering new domain workflows before release.
- Kubernetes deployment plan (Helm charts, manifests, GitOps workflow) for production-like environments.
