# Kin Coding Agent Contract

This file is the canonical contract for coding agents working in this repository.

## Coding Agent Contract summary

### Role

Repository coding agent，負責依目前授權的 issue scope 實作、驗證並提供可追溯 evidence。

### Goal

在 current issue 的 Acceptance Criteria / non-goals 與 active roadmap slice 範圍內，交付最小、coherent、可驗證的 change。

### Can

- 檢查 repository、current issue、相關 product/architecture docs、CI 與 review state。
- 在授權 scope 內修改檔案、執行 tests/checks、建立 branch/commit/PR。
- 對 adjacent blocker 提出 follow-up issue 或建議，但不得 silently 擴大目前 scope。

### Cannot

- 只因 long-term product docs 提及某能力就提前實作 future scope。
- 繞過本文件定義的 source-of-truth precedence、architecture、testing、coordination 或 required review/merge gates。
- 將未執行或失敗的 verification 宣稱為通過。
- 未經明確 operational contract 授權執行 destructive 或 production-impacting action。

### Escalate When

- current issue、roadmap 或 product scope 存在會影響實作方向的 material conflict。
- domain terminology / ownership 的歧義足以影響 modeling。
- work 需要 current contracts 尚未授權的新 architectural decision。
- required validation 無法完成，或 production/destructive/external side effect 缺乏明確 authority。

### Done When

- Acceptance Criteria 已滿足，non-goals 與 active roadmap boundary 均被遵守。
- applicable tests/checks 已實際執行，或明確記錄 blocked/skipped 原因。
- architecture/scope self-review 已完成，required review feedback 已 disposition。
- completion evidence 已記錄於 GitHub；若 required review/CI 尚未完成，不得宣告 merge-ready 或 done。

此 summary 只固定 authority 與 completion boundary；下面的 source-of-truth、architecture、testing、workflow 與 skills 仍是詳細 canonical rules。Local `AGENTS.md` 只補充 subtree-specific constraints，不重複這份 global contract。

## Source-of-truth precedence

When instructions conflict, follow this order:

1. The current GitHub issue and its explicit acceptance criteria / non-goals.
2. `docs/product/mvp-roadmap.md` for the currently authorized product slice.
3. `docs/product/product-scope.md` for long-term product/domain context.
4. This root `AGENTS.md` and any applicable local `AGENTS.md` files.
5. Architecture contracts and repo-local skills under `.agents/skills/`.
6. Existing implementation conventions.

Do not implement future scope merely because it appears in product documentation.

## Hierarchical agent contracts

Kin uses hierarchical agent contracts so the monorepo can scale across backend, mobile, web, workers, and future applications without bloating the root policy.

- The root `AGENTS.md` owns repository-wide policy.
- A local `AGENTS.md` applies to its directory subtree and owns context-specific constraints for that area.
- Local contracts inherit all root rules. They may make root rules more concrete or stricter, but must never weaken or contradict them.
- Repo-local skills under `.agents/skills/` own reusable procedures, checklists, and technology-specific how-to guidance.
- Agents working inside a subtree must read the root contract first, then the nearest applicable local contract, then only the skills relevant to the task.
- Keep policy in `AGENTS.md`; keep repeatable procedures in skills. Avoid duplicating the same rule across layers.

Principle:

> Root AGENTS.md owns policy; local AGENTS.md owns context-specific constraints; skills own procedures.

This structure should remain extensible to additional areas such as `apps/web`, `apps/admin`, or worker applications.

## Required implementation order

Feature work must proceed from the inside out:

1. Discover or confirm ubiquitous language and domain responsibility.
2. Define the user interaction / use case, including success and failure paths.
3. Define or refine domain entities, value objects, aggregates, invariants, and domain events.
4. Write or update domain tests.
5. Define application commands / queries and required ports.
6. Write application tests using fakes or in-memory ports.
7. Implement adapters and infrastructure only after the inner contracts are clear.
8. Add adapter / integration tests.
9. Add delivery-layer code such as HTTP handlers, workers, or mobile integration last.

Do not start a feature by designing database tables, ORM models, HTTP endpoints, or provider SDK usage.

## Architecture rules

Kin uses a Go modular monolith with DDD, Clean Architecture, Hexagonal Architecture, CQRS, and Domain Events.

- Dependencies point inward. Domain code must not depend on infrastructure, frameworks, persistence, HTTP, queues, or external providers.
- Infrastructure is a plugin to the application/domain, not the foundation of the design.
- Domain models, persistence models, and API DTOs are separate concepts and must not be conflated for convenience.
- Cross-domain behavior must use explicit application orchestration or contracts; avoid arbitrary imports between modules.
- Prefer high cohesion inside a domain/module and low coupling between modules.
- Keep bounded-context boundaries explicit and treat early boundaries as hypotheses that may evolve through domain discovery.

See `.agents/skills/architecture/SKILL.md` for the detailed checklist.

## CQRS rules

Commands and queries are conceptually separate.

Write side:
- Commands express intent to change state.
- Domain invariants are enforced before state changes are persisted.
- Aggregates exist to protect consistency boundaries, not to make UI queries convenient.
- Meaningful state changes may emit domain events.

Read side:
- Queries must not mutate domain state.
- Read models may be denormalized or projection-oriented when useful.
- Queries may return purpose-built DTOs without reconstructing full write-side aggregates.
- Read-side convenience must not distort the write-side domain model.

Event Sourcing is not a default requirement.

## AI / external-provider rules

LLMs and external providers are outer adapters.

- Domain code must not know about OpenAI, model names, prompts, JSON response shapes, Spotify, YouTube, or other providers.
- Provider output must be normalized and validated before entering domain/application logic.
- Provider-specific failures must be translated at the adapter boundary.
- Prefer ports that describe required business capabilities rather than vendor APIs.

## Testing contract

Tests must mirror architecture boundaries:

- Domain: fast deterministic unit tests for invariants, value objects, state transitions, and domain events.
- Application command/use case: tests with fakes/in-memory ports; no real database or external network.
- Query/read model: filtering, ordering, privacy projection, pagination, and returned shape.
- Adapters: integration or contract tests against concrete persistence/provider boundaries.
- End-to-end: only critical vertical flows; do not use E2E tests to compensate for missing lower-level tests.

See `.agents/skills/testing/SKILL.md` for detailed rules.

## Issue-driven workflow

Before editing code, every agent must:

1. Read this file.
2. Read the current issue completely, including acceptance criteria and non-goals, and confirm it is eligible under the active roadmap slice.
3. Read `.github/AGENT_COORDINATION.md` and inspect current claims / active PRs.
4. Read the nearest applicable local `AGENTS.md` for the files being changed.
5. Read only the relevant product / MVP / architecture documents and skills.
6. Identify the affected domain(s), interaction(s), command(s), query(s), and boundaries.
7. Complete the coordination preflight / claim arbitration before creating a branch, PR, commit, or changing code for the issue.
8. State assumptions explicitly in the PR when the issue leaves material ambiguity.

The coordination protocol governs ownership/hand-off only; it does not override the source-of-truth precedence, roadmap authorization, issue scope, architecture, testing, or review requirements in this contract.

During implementation:

- Keep scope limited to the issue.
- Implement from inner layers outward.
- Prefer the smallest coherent vertical change.
- Do not introduce infrastructure “for later” without an active use case.
- Do not introduce microservices, distributed CQRS, or Event Sourcing without a dedicated architectural decision.

Before completion:

- Run the smallest relevant tests first, then all repository-required checks.
- Review the diff for architecture leakage, accidental coupling, and scope creep.
- Verify every acceptance criterion explicitly.
- Record test/validation evidence in the PR.
- Process external PR review feedback, including CodeRabbit when enabled for the repository.
- When CodeRabbit is enabled, its review for the current non-draft PR and current diff must complete before merge.
- Every actionable review comment must be fixed or explicitly declined with a technical reason before merge.
- If a review comment has been read but is still being worked on, acknowledge it with an 👀 reaction where supported.
- CI passing alone is not sufficient for merge readiness when the required external review has not completed or actionable review feedback remains unresolved.

See `.agents/skills/workflow/SKILL.md` for the detailed issue-to-PR workflow.

## Skills

- Architecture: `.agents/skills/architecture/SKILL.md`
- Testing: `.agents/skills/testing/SKILL.md`
- Workflow: `.agents/skills/workflow/SKILL.md`

Keep detailed procedures in skills. Keep this file concise, normative, and stable.
