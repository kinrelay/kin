# API Agent Contract

This contract applies to work under `apps/api/` and inherits the repository root `AGENTS.md`.

## Scope

`apps/api` is the Go backend application. Its local contract makes the repository-wide architecture rules concrete for backend work without choosing implementation libraries prematurely.

## Required boundaries

- Keep domain, application, adapter/infrastructure, and delivery concerns distinct.
- Commands and queries remain conceptually separated.
- Domain behavior and invariants must not depend on persistence, HTTP, queues, serialization, or provider SDKs.
- Application use cases orchestrate domain behavior through explicit ports.
- Adapters implement ports for persistence, external providers, queues, or other infrastructure.
- Delivery code such as HTTP handlers or workers depends inward and must not become the home for business policy.
- Persistence models, API DTOs, and domain models are separate concepts.

## Design order

For backend features, start from domain responsibility and use cases, then define ports, then implement infrastructure and delivery adapters.

Do not begin feature design from database tables, ORM structures, route definitions, request/response payloads, or vendor SDKs.

## CQRS

- Write-side behavior protects invariants and consistency boundaries.
- Read-side behavior may use purpose-built projections or DTOs.
- Query convenience must not distort write-side aggregates.

## Extensibility

Keep this file minimal. Go-specific implementation procedures, package conventions, error handling, tooling, persistence patterns, and testing commands should move to a dedicated backend skill when the Go application is bootstrapped.

Do not introduce persistence libraries, router frameworks, queue technology, or API schemas through this contract alone.
