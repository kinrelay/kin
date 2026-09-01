# GitHub Agent Coordination Protocol

本文件定義 coding agents 在 GitHub issue / PR 上協作時的互斥、ownership 與進度紀錄規則。GitHub 是跨 Chat、排程與 coding agent 的 coordination source of truth；conversation history 不能取代 GitHub 上可驗證的狀態。

## 適用範圍與 authority

所有會認領 GitHub issue、建立 branch / PR、修改程式碼或處理 review 的 agent 都適用，包括 manual Chat、排程 Autopilot、Codex 與其他 coding agent。

Repository 的 root/local `AGENTS.md`、active roadmap、issue acceptance criteria / non-goals 與 repo-local skills 維持原本 authority。本文件只定義 ownership / handoff procedure；不得擴張 issue scope或繞過 roadmap、architecture、testing、review gate。Root `AGENTS.md` 與 workflow skill 必須明確引用本文件。

## Canonical coordination state

Issue comment 是 issue-backed work 的 claim / progress canonical record。Assignee / label 只能輔助顯示，不能單獨代表 agent ownership。

### Trusted coordination comments

只有可信任 GitHub actor 建立的 `<!-- agent-coordination:v1 -->` comment 才能改變 ownership。可信任 actor 必須能驗證為 repository owner、具有 `write` / `maintain` / `admin` 權限，或 repository 明確 allowlist 的 GitHub App / bot identity。`author_association=MEMBER` / `COLLABORATOR` 本身不足以建立 trust。

無法驗證 author trust 時，該 comment 只能視為 discussion，不得成為 claim、heartbeat、release、takeover 或 revocation record；新 implementation ownership 必須 fail closed。

### Worker identity、latest state 與 lease anchor

Canonical worker identity 是完整 `(agent, session)`。Session 應為可辨識且不與其他 active worker 混淆的 stable run/session identifier。

對每個 trusted `(agent, session)`，依 GitHub server `created_at`、必要時 comment id，選出最新一筆 coordination comment作 authoritative current state。舊 `claim: active` 不得覆蓋較新的 `claim: released` / `inactive` / `done` / `blocked` / `stale`。

Ownership acquisition metadata 與 latest state 分離。每段 continuous active lease 以第一次 trusted `claim: active` comment 作 immutable **server-side lease anchor**：

- anchor ordering 只使用 GitHub server `created_at`；timestamp 相同再比較 comment id；
- claimant 自填的 `claimed_at`、`heartbeat` 或本機時間不得決定 ownership；
- 後續 active heartbeat / progress 必須重複完整 state，並引用同一 `lease_anchor: <claim-comment-id-or-url>`；
- sparse heartbeat 不得取代完整 current state，也不得改變 lease acquisition order。

## Revoked anchor set

Arbitration 不只歸約每個 worker 的 latest state，也必須先建立 repository issue scope 內的 **revoked anchor set**。

任何 trusted coordination comment 若包含：

```text
takeover_of_anchor: <old-anchor-id-or-url>
```

即為該 old anchor 的永久 revocation tombstone。只要 takeover 已完成並留下此 metadata，所有 arbiter 在處理任何 worker state 前，都必須先把 `takeover_of_anchor` 加入 revoked anchor set；之後凡引用該 anchor 的 heartbeat、progress、release 或其他 late state 一律忽略，不得復活 ownership，即使它的 GitHub `created_at` 較新。

Revocation 是跨 worker reduction：不能只看舊 owner 自己的 latest comment。Trusted takeover tombstone 一旦成立，在同一 issue ownership history 中永久有效；除非該 takeover comment 本身被 GitHub 刪除/失去 trust，否則後續 arbitration 都必須保留此 evidence。

只有 latest trusted comment 同時為 `claim: active`、引用未被 revoked 的有效 anchor，且 lease 未過期，worker 才算 live claim。Release / inactive / stale / done / blocked、lease expiration，以及完成 takeover 都是不可逆 lease boundary。原 session 若之後重新工作，必須發布新 claim、取得新 anchor並重新 arbitration。

## Preflight

開始 implementation ownership 前，agent 必須先：

1. 完整讀取 issue AC/non-goals並確認 active roadmap eligibility；
2. 檢查最新 trusted coordination comments並建立 revoked anchor set；
3. 檢查 live claims；
4. 檢查 linked / active implementation PR；
5. 檢查 branch commits、CI、review與可歸因其他 session 的 activity；
6. 確認 dependency與 human-only blocker 狀態。

另一個 agent/session 有有效 live claim 時，不得平行實作相同 issue，也不得另開 duplicate PR。

## Agent identity

至少區分：

- `manual-chat`
- `side-project-autopilot`
- 其他穩定名稱，例如 `codex`

Ownership 一律比較完整 `(agent, session)`，不能只比較 session。

## Claim format

第一次認領 issue：

```text
<!-- agent-coordination:v1 -->
agent: side-project-autopilot
session: <stable-session-or-run-id>
status: claimed
claim: active
branch: <branch-or-none>
pr: <pr-number-or-none>
head: <sha-or-none>
heartbeat: <ISO-8601 UTC timestamp; informational only>
next: <next concrete action>
```

發布後重新抓取該 comment 的 GitHub server `created_at` 與 comment id/url，固定為 lease anchor。後續 active progress：

```text
<!-- agent-coordination:v1 -->
agent: <agent>
session: <session>
status: review
claim: active
lease_anchor: <first-claim-comment-id-or-url>
claimed_at: <anchor GitHub created_at; informational copy only>
branch: <branch>
pr: <number>
head: <sha>
validation: <canonical CI/test evidence>
heartbeat: <ISO-8601 UTC timestamp>
next: <next concrete action>
```

## Post-claim ownership arbitration

發布 `claim: active` 不代表立即取得 ownership。在建立 branch、PR、commit 或修改程式碼前：

1. 重新抓最新 trusted coordination comments與剛發布 claim 的 server metadata；
2. 先建立完整 revoked anchor set；
3. 對每個 `(agent, session)` 做 latest-state reduction；
4. 丟棄引用 revoked anchor、released/inactive/stale/done/blocked、或已過期的 state；
5. 對剩餘衝突 live claims解析 immutable server-side lease anchor；
6. anchor `created_at` 較早者為唯一 winner；相同時比較 comment id；無法排序時 fail closed；
7. winner 再確認自己仍是唯一有效 owner後才能 mutation；
8. loser 留 `status: stale`、`claim: released`並停止。

舊資料若沒有 `lease_anchor`，只能回溯該 continuous active lease 第一筆 trusted `claim: active` comment 的 server `created_at` + id 作 migration fallback。metadata 不足時 fail closed。

等待 dependency 使用 `status: waiting-dependency`、`claim: inactive`。

## Progress / heartbeat

Material transition 至少使用：

- `claimed`
- `red`
- `implementing`
- `ci`
- `review`
- `blocked`
- `merge-ready`
- `done`
- `waiting-dependency`
- `stale`

`red` 是 Kin strict TDD 專用狀態：只有 issue / testing contract 要求 Red → Green，且 regression test 已因正確原因客觀失敗時才能使用，不是 generic error state。

所有 `claim: active` heartbeat / progress 必須是完整 state snapshot並引用同一有效 anchor。重要 commit、CI、review disposition、merge 等事件應更新；即使沒有 material transition，也至少每 **2 小時**更新一次。Heartbeat 不改 lease acquisition order，也不能復活 expired/revoked anchor。

## Claim lease / stale takeover

預設 lease **8 小時**。只有該 trusted `(agent, session)` 在 lease 尚未過期前發布、引用有效且未 revoked anchor 的最新 coordination heartbeat可續期。CI、review bot、其他 reviewer/session activity都不能延長 owner lease。

Lease 超過 8 小時即 terminal expired；舊 anchor 後續 heartbeat不得復活。Takeover preflight仍須檢查 linked PR、branch commits、review replies與 issue discussion，辨識是否有可歸因原 owner、但漏寫 heartbeat 的近期 mutation。若存在，可暫緩最多一個 2 小時 heartbeat cadence；純 CI/bot/他人 activity只是 context，不是 veto。

只有 recent-activity hold 結束、inactivity preflight 通過並確認 safe 後，新 agent 才可發布 takeover claim並取得新 anchor。完成 takeover 時必須再留下 machine-readable revocation evidence：

```text
<!-- agent-coordination:v1 -->
agent: <new-agent>
session: <new-session>
status: implementing
claim: active
lease_anchor: <new-anchor-id-or-url>
takeover_of_anchor: <old-anchor-id-or-url>
takeover_from: <old-agent>/<old-session>
takeover_reason: lease-expired-and-preflight-safe
branch: <existing-or-new-allowed-branch>
pr: <existing-pr-or-none>
heartbeat: <ISO-8601 UTC timestamp>
next: <next concrete action>
```

此 `takeover_of_anchor` comment 是永久 tombstone；任何 future arbiter 必須先套用 revocation，再做 latest-state reduction。尚在 hold 期間不得 adoption existing PR 或推送既有 branch。

## Existing PR wins / adoption

Issue 已有 active implementation PR 時優先沿用 existing PR，不另開 duplicate，但不能無條件接管：

- 有有效 live owner：新 claimant退出；
- 有尚未清除的 human-only `status: blocked`：不得選擇/adoption；
- author / trusted owner 明確 handoff、release或 inactive：eligible agent可 claim existing PR並完成 arbitration；
- owner expired/stale：必須先完成 stale-takeover preflight與 hold，完成 takeover時以 `takeover_of_anchor` revoke舊 lease，才能 adoption；
- 沒有 coordination metadata：不得視為 ownerless；先把 PR author / branch owner視為可能 active，只有明確 handoff/consent，或超過 8 小時且完整 inactivity preflight通過後才可 adoption；
- branch不可寫或 ownership無法安全驗證：記錄 blocker/release，不得用 duplicate PR規避。

除非 repo contract / issue 明確允許 stacking，否則每 repo優先一個 active implementation PR。

## Work selection priority

Autopilot 仍須服從 roadmap、dependencies與 issue priority，依序：

1. 可由 agent 推進的 broken active implementation PR；
2. active implementation PR actionable review；
3. active implementation PR failing canonical CI；
4. 其他未被 human-only blocker 阻擋的 active implementation PR，包括 merge-ready；
5. 目前 agent/session自己持有 live claim 的 unfinished issue；
6. roadmap / dependency blocker；
7. 其他沒有 live claim、沒有未清除 blocked state且合法 eligible issue。

其他 session live claim 必須 skip。Human-only blocker使用 `status: blocked`、`claim: inactive`，在清除前同時排除 active-PR queue、adoption與 new-issue eligibility。

## Release claim

下列情況 terminate / release current lease：PR merged、issue closed、PR closed without merge、implementation abandoned、PR/branch superseded、lease expired/completed takeover，或進入 human-only blocked。

一般完成：

```text
<!-- agent-coordination:v1 -->
agent: <agent>
session: <session>
status: done
claim: released
lease_anchor: <current-anchor-if-applicable>
pr: <number-or-none>
head: <merged-head-or-merge-sha>
heartbeat: <ISO-8601 UTC timestamp>
next: <next eligible issue or none>
```

Takeover completion 不應由舊 owner自行 release來建立 revocation；新的 trusted owner必須用 `takeover_of_anchor` tombstone明確撤銷舊 anchor。

若 PR closed/abandoned/superseded但 issue未完成，使用 `claim: released` 並寫 `next: issue remains eligible for reassignment`。Expired/taken-over anchor不得再次使用。

Human-only blocker使用 `status: blocked`、`claim: inactive`並記錄所需輸入；不得假裝完成或無限持有 claim。

## Rules for manual Chat sessions

Manual Chat 在 coding 前也必須讀取最新 coordination history、建立 revoked anchor set，並遵守相同 trusted-author validation、server-side anchor、state reduction、post-claim arbitration與 stale takeover。看到 `side-project-autopilot` live claim應避免撞工；Autopilot看到 `manual-chat` live claim同樣必須跳過，除非依規則確定失效。
