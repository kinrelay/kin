# Kin

Kin is an AI-assisted friendship context layer designed to help real friends stay connected without requiring them to constantly post, browse, or maintain another attention-driven social feed.

Kin turns explicitly authorized digital activity into lightweight, permissioned social context that can help close friends understand what matters to each other and restart meaningful conversations naturally.

## Repository purpose

This repository contains the Kin product implementation and its shared engineering contracts. It is organized as a monorepo so backend, mobile, and future applications can evolve under the same product and architecture rules while keeping app-local constraints explicit.

Current planned application areas:

- `apps/api` — Go backend
- `apps/mobile` — Expo / React Native mobile application

Additional apps or workers may be added later when justified by real product needs.

## Engineering direction

The current architecture direction is:

- Go modular monolith for the backend
- Expo + React Native + TypeScript for mobile
- Domain-Driven Design where meaningful domain behavior exists
- Clean Architecture
- Hexagonal Architecture / Ports and Adapters
- CQRS-style command/query separation
- Domain Events
- Infrastructure and external providers treated as outer adapters

The core development rule is:

> Start from product interaction and domain behavior. Design persistence, APIs, frameworks, and providers only after the inner contracts are clear.

Do not start feature design from database tables, ORM models, HTTP endpoints, SDKs, or UI components.

## Agent contracts

All coding agents must read [`AGENTS.md`](./AGENTS.md) before making changes.

Kin uses hierarchical agent contracts:

- Root `AGENTS.md` defines repository-wide policy.
- App-local `AGENTS.md` files add constraints for their directory subtree.
- `.agents/skills/` contains reusable procedures and detailed implementation guidance.

Current local contracts:

- [`apps/api/AGENTS.md`](./apps/api/AGENTS.md)
- [`apps/mobile/AGENTS.md`](./apps/mobile/AGENTS.md)

Current shared skills:

- [Architecture](./.agents/skills/architecture/SKILL.md)
- [Testing](./.agents/skills/testing/SKILL.md)
- [Workflow](./.agents/skills/workflow/SKILL.md)

## Documentation map

The intended documentation hierarchy is:

1. `docs/product/product-scope.md` — long-term product and domain map
2. `docs/product/mvp-roadmap.md` — current validation sequence and vertical MVP slices
3. Architecture contracts and repo-local skills — implementation constraints and procedures
4. GitHub Issues — currently authorized unit of work

The product scope may describe future capabilities. Their presence there does **not** authorize implementation. The active issue and current MVP slice define what should be built now.

The product scope and MVP roadmap documents are part of the current foundation work and may not exist on `main` yet.

## Development workflow

Feature work should generally flow as:

```text
Product / user interaction
        ↓
Use case and domain responsibility
        ↓
Domain rules and tests
        ↓
Application commands / queries and ports
        ↓
Adapters and infrastructure
        ↓
Delivery layer
        ↓
PR validation and review
```

Changes should be issue-driven and delivered through branch → pull request → review → merge unless an exceptional repository-bootstrap situation requires otherwise.

## Project status

Kin is currently in the **foundation / product-definition phase**.

The immediate goal is to establish a durable product map, MVP sequence, architecture contracts, and engineering workflow before substantial feature implementation begins.
