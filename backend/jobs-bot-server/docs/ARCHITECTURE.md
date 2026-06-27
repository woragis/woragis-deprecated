# System Architecture

## Overview

This system is an autonomous, observable, and robust job application automation service. It separates LLM-powered decision-making from deterministic browser automation, and is designed for reliability, traceability, and maintainability.

## Core Design Principles

- LLMs only for language/decision tasks
- All browser actions are deterministic and blocking
- System is observable and restartable
- One run at a time (single worker)

## High-Level Diagram

```
[API] <-> [Orchestrator] <-> [Worker] <-> [Platform] <-> [Browser]
   |           |                |             |
 [WebSocket] [DB/Logs]      [LLM]         [Playwright]
```

## Folder Structure

- `app/` — Main application code
- `docs/` — Documentation
- `logs/` — Persistent logs
- `profiles/` — User profiles

## Key Modules

- API (FastAPI)
- Orchestrator (run lifecycle, scheduling)
- Worker (job processing)
- Platforms (Greenhouse, Lever, etc.)
- Browser (Playwright session, humanization)
- LLM (fit check, text generation)
- Persistence (DB, models)
- Utils (logging, errors, time)
