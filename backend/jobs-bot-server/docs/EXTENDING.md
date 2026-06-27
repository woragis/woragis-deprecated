# Extending & Customizing

## Adding New Platforms

- Create a new file in `app/platforms/` inheriting from `JobPlatform`
- Implement `login`, `collect_jobs`, `apply_to_job`
- Add to allowed platforms in API validation

## Adding Event Types

- Update `emit_ws_event` and event docs

## Feature Flags & Config

- Add new config to `Config` class in `app/config.py`
- Document in README and `docs/`
