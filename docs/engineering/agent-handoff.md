# Agent Handoff Contract

## Purpose

Kin agent work must be transferable without relying on a previous agent's conversation history. This document defines the canonical, human-readable and machine-parseable handoff contract for implementation, review, and verification work.

The contract records evidence and unresolved state; it does not replace the GitHub issue, active MVP roadmap, `AGENTS.md`, `.github/AGENT_COORDINATION.md`, CI, or review requirements. Those sources remain authoritative for scope, ownership, and merge gates.

Use YAML in a fenced code block when publishing a handoff in a GitHub issue or PR. Additional prose may explain context, but the structured artifact is the canonical handoff payload.

## Common semantics

Every artifact uses explicit status values instead of ambiguous statements such as "done" or "looks good".

- `pass`: observed evidence satisfies the stated check.
- `fail`: observed evidence does not satisfy the stated check.
- `blocked`: the work/check cannot continue until a named dependency or human/tooling condition changes.
- `skipped`: an applicable check was intentionally not executed; a reason is required.
- `not_applicable`: the check does not apply to this change; a reason is required.
- `unknown`: evidence is not available; never treat this as pass.

Evidence should be stable and reviewable: GitHub issue/PR references, exact CI/check names, commands actually executed, review disposition, or other repository artifacts. Do not claim a command passed if it was not run.

`applicable_contracts` should identify the root and local `AGENTS.md` files that governed the affected subtree. It references contracts; it does not copy their policy text.

## `ImplementationResult`

Use after an implementation slice reaches a handoff point: ready for review, partial, or blocked.

```yaml
kind: ImplementationResult
status: ready | partial | blocked
issue: "#123"
authorized_slice: "MVP 3 — Privacy"
affected:
  apps_or_subtrees:
    - "apps/api"
  applicable_contracts:
    - "AGENTS.md"
    - "apps/api/AGENTS.md"
  domains:
    - "Privacy"
  interactions:
    - "Authorize pending delivery"
  commands: []
  queries: []
  events: []
branch: "feat/example"
pr: "#124"
commits:
  - "<commit-or-range>"
files_changed:
  - "apps/api/..."
decisions:
  - "<decision and why>"
assumptions:
  - "<assumption that reviewers must verify>"
verification:
  executed:
    - check: "<command or canonical CI check>"
      result: pass | fail
      evidence: "<observable result/reference>"
  skipped:
    - check: "<applicable check>"
      result: skipped
      reason: "<why>"
acceptance_criteria:
  completed: 0
  total: 0
  items:
    - criterion: "<criterion>"
      status: pass | fail | blocked | unknown
      evidence: "<reference>"
non_goals:
  respected: true
  notes: []
architecture_domain_impact:
  - "<boundary/invariant/read-model impact or none>"
known_risks: []
unresolved_items: []
next_action: "<single concrete next action>"
```

Rules:

1. `status: ready` requires all implementation-owned AC to have evidence and no known implementation blocker.
2. `status: partial` is valid when useful work exists but AC remain incomplete; list the incomplete items explicitly.
3. `status: blocked` requires the blocker and the condition that would unblock it.
4. `files_changed` and `affected.apps_or_subtrees` must be consistent with the PR diff.
5. `verification.executed` records observations only. An unavailable check belongs in `skipped` or an unresolved blocker, never as pass.

## `ReviewResult`

Use when reviewing a current PR diff. It distinguishes CI state from review readiness.

```yaml
kind: ReviewResult
status: ready | changes_required | blocked
issue: "#123"
pr: "#124"
reviewed_head: "<exact-head-sha>"
applicable_contracts:
  - "AGENTS.md"
  - "apps/api/AGENTS.md"
findings:
  blocking: []
  non_blocking: []
  accepted_fixed:
    - finding: "<finding>"
      evidence: "<fix/reference>"
  declined:
    - finding: "<finding>"
      reason: "<technical reason grounded in canonical scope/architecture>"
      evidence: "<reference>"
acceptance_criteria:
  status: pass | fail | blocked | unknown
  notes: []
non_goals:
  status: pass | fail | unknown
  notes: []
architecture_boundary:
  status: pass | fail | unknown
  notes: []
gates:
  ci:
    status: pass | fail | pending | unknown
    evidence: "<current-head check reference>"
  coderabbit:
    status: pass | findings | pending | not_applicable | unknown
    evidence: "<current-diff review reference/reason>"
  human_review:
    status: pass | findings | pending | not_required | unknown
    evidence: "<review reference/reason>"
unresolved_risks: []
recommendation: merge | changes_required | blocked
```

Rules:

- `recommendation: merge` requires the repository's actual merge gates to be satisfied. CI green alone is insufficient when current-diff CodeRabbit or required human review is pending.
- Every substantive finding must be fixed or explicitly declined with technical evidence.
- `reviewed_head` must match the diff being recommended. A later commit makes the result stale until the affected review/verification is refreshed.

## `VerifierResult`

Use for an independent verification pass or for a worker handing off concrete test/check evidence.

```yaml
kind: VerifierResult
status: pass | fail | blocked
issue: "#123"
pr: "#124"
verified_head: "<exact-head-sha>"
checks:
  executed:
    - check: "<command/check>"
      result: pass | fail
      observed: "<what was observed>"
      evidence: "<reference>"
  skipped:
    - check: "<check>"
      result: skipped | not_applicable
      reason: "<why>"
failures: []
infra_tooling_limitations: []
conclusion: pass | fail | blocked
```

Rules:

- `conclusion: pass` means all verification required for this verifier's declared scope was actually executed and passed.
- A required check that cannot run normally makes the result `blocked` unless the repository contract explicitly permits a skip.
- A failed check is `fail`; do not relabel it as tooling noise without evidence.

## Compact example — API privacy task

```yaml
kind: ImplementationResult
status: ready
issue: "#64"
authorized_slice: "MVP 3 — Privacy"
affected:
  apps_or_subtrees: ["apps/api"]
  applicable_contracts: ["AGENTS.md", "apps/api/AGENTS.md"]
  domains: ["Privacy", "Friendship"]
  interactions: ["Authorize pending delivery"]
  commands: []
  queries: []
  events: []
branch: "feat/revision-bound-pending-intent"
pr: "#65"
commits: ["<current-head>"]
files_changed: ["apps/api/src/...", "apps/api/tests/..."]
decisions: ["Re-read privacy and relationship revisions inside the atomic authorization guard."]
assumptions: []
verification:
  executed:
    - check: "canonical CI"
      result: pass
      evidence: "PR #65 current-head checks"
  skipped: []
acceptance_criteria:
  completed: 12
  total: 12
  items:
    - criterion: "Privacy/relation race cannot dispatch stale authorization"
      status: pass
      evidence: "deterministic race coverage + canonical CI"
non_goals:
  respected: true
  notes: ["No notification transport or MVP 4 behavior added."]
architecture_domain_impact: ["Atomic authorization remains an application/domain boundary; transport is unchanged."]
known_risks: []
unresolved_items: []
next_action: "Obtain current-diff review and disposition findings."
```

## Compact example — mobile task with unavailable check

```yaml
kind: VerifierResult
status: blocked
issue: "#88"
pr: "#89"
verified_head: "<current-head>"
checks:
  executed:
    - check: "unit tests"
      result: pass
      observed: "Mobile unit suite completed successfully."
      evidence: "CI unit-test job"
  skipped:
    - check: "iOS device smoke"
      result: skipped
      reason: "No iOS simulator/device runner is available in the current CI environment."
failures: []
infra_tooling_limitations:
  - "iOS simulator/device runner unavailable"
conclusion: blocked
```

The example intentionally remains blocked: an unavailable required check is not equivalent to a pass.

## Handoff lifecycle

1. Implementation publishes or updates `ImplementationResult` when handing work to review/verification.
2. Reviewer evaluates the current diff and publishes `ReviewResult`; findings are dispositioned explicitly.
3. Verifier publishes `VerifierResult` for the checks it actually observed.
4. Merge readiness is computed from current repository contracts and evidence, not from any artifact's self-declared status alone.
5. If the PR head changes materially, stale review/verification artifacts must be refreshed for the new head where required.

Keep artifacts concise. Link to canonical evidence rather than duplicating logs, policy, or large diffs.