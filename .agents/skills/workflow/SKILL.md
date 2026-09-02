# Workflow Skill

Use this skill for issue-to-PR execution.

## Objective

Keep agent work scoped, architecture-aligned, reviewable, coordinated, and evidence-based.

## Before implementation

1. Read root `AGENTS.md`.
2. Read `.github/AGENT_COORDINATION.md`.
3. Read the current issue completely.
4. Identify acceptance criteria and explicit non-goals.
5. Read the minimum relevant product/MVP/architecture documents.
6. Identify affected domain(s), actor(s), interaction(s), command(s), query(s), events, and external boundaries.
7. Confirm whether the requested change is domain behavior, application orchestration, read-model work, adapter work, delivery work, or documentation.
8. Inspect current issue/PR/branch/CI/review state and complete the coordination protocol's preflight and deterministic post-claim arbitration before creating a branch, PR, commit, or changing code.
9. Do not create infrastructure or abstractions for hypothetical future work.

If another live agent/session owns the issue, stop work on that issue and select another eligible task. If issue comments cannot be read well enough to establish unique ownership, fail closed rather than creating duplicate work.

If product or domain terminology is unclear, prefer updating/discussing the language before encoding ambiguous concepts into persistence or APIs.

## Implementation order

For feature changes, work inside-out:

1. Domain concepts and invariants.
2. Domain tests.
3. Application command/query/use-case contract.
4. Application tests with fakes/in-memory ports.
5. Ports required by the use case.
6. Concrete adapters/infrastructure.
7. Adapter/integration tests.
8. Delivery layer.
9. Critical vertical/E2E validation when needed.

A change may skip irrelevant layers, but must not invert the dependency direction.

## Scope discipline

- The GitHub issue defines the authorized unit of work.
- Product scope provides context, not permission to implement all future capabilities.
- MVP roadmap determines sequencing when available.
- Avoid opportunistic refactors unrelated to the issue.
- If an adjacent problem materially blocks the issue, document it and prefer a follow-up issue rather than silently expanding scope.

## Branch and commit expectations

After the repository has an initial commit:

- Work on a dedicated branch.
- Keep commits coherent and reviewable.
- Prefer conventional-style messages where practical.
- Do not mix unrelated feature work into the same branch.

## Validation

Run the smallest relevant checks first, then broader required checks.

For code changes, record exact commands and outcomes.
For documentation-only changes, inspect links, paths, terminology, and consistency; runtime tests are not required unless executable/configuration behavior changed.

Never state that a test or check passed unless it was actually executed.

## Self-review before PR

Review the final diff and answer:

- Does every changed file belong to the issue?
- Did any infrastructure detail leak into domain/application contracts?
- Did a query mutate state?
- Did read-side needs distort the write model?
- Are provider-specific DTOs/types isolated in adapters?
- Are new abstractions justified by an active use case?
- Are tests placed at the lowest useful architecture layer?
- Are all acceptance criteria satisfied?
- Are non-goals still respected?

## Pull request contract

PR descriptions should include:

- Summary
- Why this change exists
- Architecture/domain impact
- Scope / non-goals when useful
- Validation evidence with exact commands
- Known follow-ups or intentionally deferred work
- Issue linkage (`Closes #...` when completion is intended)

Do not mark work complete merely because code compiles or an agent reports “done”. Completion means acceptance criteria are demonstrably satisfied.

## Review handling

Every substantive human or automated review comment must be processed explicitly.

For each actionable comment:

1. Read and assess the suggestion against the issue scope and repository contracts.
2. If work is still in progress after reading it, add an 👀 reaction where supported so reviewers can see it has been acknowledged.
3. Either:
   - fix the issue and reply with what changed, or
   - decline the suggestion and reply with a concrete technical reason.
4. If the resulting diff changes executable behavior or relevant configuration, rerun the smallest affected validation checks.
5. Resolve the review thread when the platform supports it and the disposition is complete.

Do not silently ignore review feedback.

### CodeRabbit

When CodeRabbit is enabled for the repository:

- Treat its actionable findings as external review input, not as architectural authority.
- Repository contracts (`AGENTS.md`, local contracts, issue acceptance criteria, and relevant skills) remain the source of truth.
- Automatic review is required on non-draft PRs.
- Do not request automatic review for draft PRs unless a task explicitly requires it.
- The CodeRabbit review for the current non-draft PR and current diff must complete before merge.
- Review feedback must be handled with the same fix-or-explicit-decline policy as human feedback.
- A green CI result does not make the PR merge-ready while the required CodeRabbit review is incomplete or actionable CodeRabbit feedback remains unresolved.

## Merge readiness

Before merge, verify all of the following:

- Acceptance criteria are satisfied.
- Required CI/checks have passed for the current diff.
- The final diff has been self-reviewed for scope and architecture leakage.
- When CodeRabbit is enabled, its review has completed for the current non-draft PR and current diff.
- All actionable human and automated review feedback is fixed or explicitly declined.
- Review-driven changes have been revalidated where relevant.
- No known unresolved blocker remains hidden in a review thread or PR conversation.
