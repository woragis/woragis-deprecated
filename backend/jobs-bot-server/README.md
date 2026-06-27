# Autonomous Job Application System

## Overview

This project is a reliable, observable, and scalable automation system for applying to online job postings. It is designed as a long-running service, not a script, and separates LLM-powered decision-making from deterministic browser automation.

## Features

- Automatically applies to jobs on supported platforms
- Manual and scheduled activation
- Daily application limits
- Real-time progress via WebSocket
- Persistent logs and observable state
- Safe, debuggable, and restartable
- Deterministic browser automation (Playwright)
- LLM used only for language tasks (fit check, cover letter, answers)

## Tech Stack

- Python
- FastAPI (HTTP API)
- Playwright (browser automation)
- OpenAI-compatible LLM
- APScheduler (scheduling)
- SQLite (initial storage)

## File Structure

```
app/
    main.py
    config.py
    api/
        routes/
        runs.py
        health.py
        websocket.py
    orchestrator/
        runner.py
        scheduler.py
        queue.py
    workers/
        worker.py
        states.py
    platforms/
        base.py
        greenhouse.py
        lever.py
        linkedin.py
        gupy.py
    browser/
        session.py
        humanize.py
    llm/
        fit_check.py
        text_gen.py
    persistence/
        db.py
        models.py
    utils/
        logging.py
        time.py
profiles/
    default_profile.json
logs/
    runs.log
requirements.txt
```

## Usage

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   playwright install
   ```
2. Run the API server:
   ```bash
   uvicorn job_applier.app.main:app --reload
   ```
3. Use the HTTP API to start job application runs.

## System Principles

- LLMs do NOT control UI
- All browser actions are deterministic and blocking
- LLMs only for language/decision tasks
- System is observable and restartable
- One run at a time (single worker)

See the system prompt for full design details.
