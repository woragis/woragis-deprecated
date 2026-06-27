# ADR-007: Docker Compose: Multi-Service Architecture

## Context
We need to run 8+ services locally for development. Need orchestration for dependencies and networking.

## Decision
Use Docker Compose with single `docker-compose.yml` in backend root, all services and workers defined.

## Rationale

### Docker Compose
- **Orchestration**: Manage all services together
- **Networking**: Automatic service discovery
- **Dependencies**: Start order management
- **Development**: Easy local setup

### Single Compose File
- **Simplicity**: One file to manage
- **Visibility**: See all services at once
- **Consistency**: Same network, same config

## Consequences

### Pros
- Easy local development
- Consistent environment
- Service discovery
- Dependency management

### Cons
- Large compose file
- All services start together
- Resource usage

## Status
Accepted - 2024-01-15

## Alternatives Considered
- Multiple compose files: More complex
- Kubernetes: Overkill for local dev
- Manual setup: Too error-prone
