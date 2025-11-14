# Woragis Chats Backend Enhancements

The chats domain now supports collaboration workflows around search, bulk management, transcripts, assignments, and realtime streaming. This document summarises the server-side features added in November 2025.

## Conversation Search

- `GET /api/chats/conversations/search`
  - Query parameters:
    - `q` – free text search across titles, descriptions, and message bodies (case-insensitive).
    - `include_archived` (optional) – include archived conversations when set to `true`.
    - `limit` (optional) – cap results (defaults to 20).
  - The authenticated user is inferred from the JWT middleware; no explicit `user_id` is required.
- Backend performs SQL filtering with safe `LOWER(...) LIKE` clauses and message sub-selects, supporting both Postgres and SQLite dev environments.

## Bulk Archive/Delete/Restore

- New POST endpoints on `/api/chats/conversations/archive|delete|restore`.
- Payload shape:
  ```json
  {
    "conversation_ids": ["<uuid>", "..."]
  }
  ```
- Soft-delete behaviour:
  - `archived_at` and `deleted_at` timestamps captured on the conversation record.
  - Restoring conversations clears both fields.
  - `ListConversations` excludes deleted threads by default.

## Shared Transcripts

- `POST /api/chats/conversations/:id/transcripts`
  - Generates a signed `share_code` and stores a JSON snapshot of messages.
  - Optional `expire_after` duration (Go duration syntax, e.g. `72h`); defaults to 7 days.
- `GET /api/chats/conversations/:id/transcripts` (owner scoped) lists prior exports.
- `GET /api/chats/transcripts/:code` returns the shared transcript payload for external viewers.
- Conversations store the latest share code in `shared_transcript` for quick access.

## Assignment History

- `POST /api/chats/conversations/:id/assign`
- `POST /api/chats/conversations/:id/unassign`
- `GET /api/chats/conversations/:id/assignments`
- Assignment metadata persists in `conversation_assignments` with audit timestamps and optional notes.
- Assigning closes previous assignments automatically to maintain a clean history.

## Realtime Streaming

- WebSocket endpoint: `GET /api/chats/conversations/:id/stream`
- Streaming hub broadcasts `messageResponse` JSON whenever a message or AI reply is appended.
- Requires authentication middleware (same JWT guard as REST endpoints).
- Clients should connect with standard websocket protocol and listen for JSON payloads representing appended messages.

## Data Model Additions

- `Conversation` fields: `AssignedAgentID`, `SharedTranscript`, `ArchivedAt`, `DeletedAt`, `LastAssignedAt`.
- New tables:
  - `conversation_transcripts`
  - `conversation_assignments`
- Auto-migration updated in `app/cmd/server/main.go`.

## Configuration / Infrastructure Notes

- No new environment variables required.
- Streaming uses the existing Fiber/WebSocket stack—ensure load balancers permit websocket upgrades on `/api/chats/conversations/*/stream`.
- Transcripts are stored as JSON strings inside the database; consider rotating them to blob storage if the payload grows beyond acceptable size.

## Frontend & QA Follow-up

- Update frontend flows to integrate search filters, archive/delete controls, transcript downloads, assignment timelines, and websocket listeners.
- Expand integration tests to cover bulk operations, transcript sharing, and MFA-protected websocket access once UI endpoints are wired up.

