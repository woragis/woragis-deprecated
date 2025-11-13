# Woragis Backend TODO

## Domains Requiring Deeper CRUD / Advanced Workflows
- **Auth**: Implement session/device management, multi-factor tokens, audit trails, OAuth provider links, bulk user admin actions.
- **Finances**: Extend reporting (cash-flow projections, tagging), multi-currency normalization, recurring schedule templates, bulk upload from CSV/OFX.
- **Languages**: Add spaced-repetition scheduling, bulk vocabulary import/export, proficiency analytics, AI-generated practice sets.
- **Projects**: Implement kanban workflows, dependency graphs, bulk milestone updates, templated project duplication.
- **Chats**: Conversation search, bulk archive/delete, shared transcripts, agent assignment history, streaming responses with websockets.
- **Ideas**: Version history, bulk node operations, collaborative editing, advanced filters for relationships.
- **Scheduler**: Complex recurrence rules (cron/rrule), bulk activation/deactivation, execution history dashboard, alerting workflows.
- **Reports**: Custom report builder, scheduling templates, multi-channel delivery, bulk regeneration and export management.
- **Monitoring**: Persist structured events in production DB, alert threshold management, Grafana dashboard provisioning.

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
