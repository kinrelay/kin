# Testing Skill

Use this skill when designing, implementing, or reviewing tests.

## Objective

Make the test architecture mirror the software architecture. Tests should provide fast feedback at inner layers and use slower integration/E2E coverage only where concrete boundaries require it.

## Domain tests

Characteristics:

- Fast and deterministic.
- No database, network, filesystem, queue, clock, or provider dependency unless represented by an explicit controllable abstraction.
- Focus on Value Object validation, invariants, state transitions, policies, and Domain Events.

Examples of what belongs here:

- whether a relationship level permits a visibility level
- whether a state transition is valid
- whether an invalid domain value can be constructed
- whether a meaningful state change emits the expected event

Prefer table-driven tests where they improve clarity.

## Application command / use-case tests

Use fakes or in-memory implementations of ports.

Verify:

- orchestration order where behavior depends on it
- correct domain behavior is invoked
- persistence ports receive the intended state
- domain events are published when appropriate
- failure paths are handled without leaking provider/storage details

Do not use a real database merely to test application orchestration.

## Query / read-model tests

Read-side tests should focus on the consumer-facing query contract:

- filtering
- ordering
- pagination
- visibility/privacy projection
- aggregation/projection behavior
- returned DTO shape and semantics

Queries must remain read-only.

## Adapter integration tests

Use these for concrete boundaries such as:

- Postgres repositories/read stores
- migrations and persistence mapping
- queue adapters
- HTTP/provider adapters
- serialization/deserialization
- OAuth/provider contract behavior

Test the adapter contract, not business rules already covered by domain/application tests.

For data adapters, prefer isolated disposable test infrastructure such as Testcontainers when introduced by an actual implementation need.

## External integration contract tests

For providers such as Spotify, YouTube, or LLM APIs:

- isolate provider-specific DTOs in the adapter
- test normalization into inner contracts
- test representative error translation
- avoid relying on live external APIs in the default test suite
- use fixtures/fakes for routine CI; add opt-in live smoke tests only when useful and safe

Never require production credentials for ordinary tests.

## End-to-end tests

Use E2E tests sparingly for critical vertical flows.

Candidate critical flow shape:

`activity contribution -> context derivation -> privacy decision -> friend-visible pulse`

E2E tests should validate composition, not replace missing domain/application tests.

## Test order during development

Run tests from fastest/narrowest to broadest:

1. affected domain tests
2. affected application/query tests
3. affected adapter integration tests
4. repository-wide unit/integration checks
5. E2E/smoke checks when required by the issue or changed boundary

## Test design rules

- Test behavior, not private implementation details.
- Keep fixtures small and semantically named.
- Prefer builders/helpers that express domain language over giant generic fixtures.
- Do not mock concrete implementation internals when a port-level fake is the architectural seam.
- A bug fix should include a regression test at the lowest layer capable of reproducing the bug.
- A documentation-only change does not require runtime tests unless the diff also changes executable/configuration behavior.

## Completion evidence

A PR should state:

- which test layers were changed
- exact commands run
- whether any required test was skipped and why
- any live-provider/manual validation performed

Do not claim a test passed unless it was actually run.
