# Project Engineering Profile

This file classifies the project against the shared Personal Engineering Guidelines. It is a routing aid, not a second copy of repository policy.

Shared baseline: https://github.com/koshuang/personal-engineering-guidelines

## Current maturity

**MVP / roadmap-gated product development.**

Kin is still validating product/domain hypotheses. Engineering should preserve strong domain boundaries and testability without prematurely paying for production-scale infrastructure or future-scope complexity.

## Primary outcome

Advance the currently authorized MVP slice with the smallest coherent, behavior-correct implementation and enough evidence to decide whether the product hypothesis is working.

## Required capabilities

- explicit active roadmap slice and issue acceptance criteria;
- domain/use-case-first modeling with provider-neutral boundaries;
- deterministic domain and application tests;
- adapter/integration tests only where concrete boundary fidelity matters;
- minimal end-to-end evidence for critical vertical flows;
- issue/PR coordination and current-head review/CI evidence;
- deliberate scope control so future product ideas do not leak into current MVP work.

## Canonical sources

Source-of-truth precedence:

1. Current GitHub issue and explicit acceptance criteria / non-goals.
2. `docs/product/mvp-roadmap.md` for the authorized slice.
3. `docs/product/product-scope.md` for longer-term context.
4. Root/scoped `AGENTS.md` and relevant repo-local skills.
5. Current code, tests, CI, PR/review state, and acceptance evidence.

The repository's source-of-truth precedence remains authoritative. This profile does not override it.

## Deliberate profile choices

- Go modular monolith, DDD/Clean/Hexagonal/CQRS/domain-event patterns are current architectural choices justified by the product model; they are not cross-project mandates.
- Infrastructure should be introduced only for an active use case; do not pre-build microservices, distributed CQRS, or event sourcing.
- Strong inner-layer tests are preferred over broad E2E suites.
- MVP speed may reduce ceremony, but not by weakening accepted behavior, roadmap authority, or required evidence.

## Revisit when

Reassess when an MVP slice moves into bounded real-world pilot, when production operations become material, or when observed domain boundaries repeatedly differ from the current model.
