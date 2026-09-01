# GitHub Agent Coordination Protocol

本文件定義 coding agents 在 GitHub issue / PR 上協作時的共用互斥與進度紀錄規則。GitHub 是跨 Chat、排程與 coding agent 的協作 source of truth；conversation history 不能取代 GitHub 上可驗證的狀態。

## 適用範圍

所有會認領 GitHub issue、建立 branch / PR、修改程式碼或處理 review 的 agent 都適用，包括手動 Chat session、排程 Autopilot、Codex 或其他 coding agent。

Repository 的 `AGENTS.md`、local `AGENTS.md`、active roadmap、issue acceptance criteria / non-goals 與 repo-local skills 仍具有原本的 authority。本文件只定義「誰正在做、如何避免撞工、如何留下可接手狀態」。Root `AGENTS.md` 與 workflow skill 必須明確引用本文件，避免 agent 只依 canonical onboarding 卻漏掉 coordination preflight。

## Canonical coordination state

Issue comment 是 claim / progress 的 canonical coordination record。Assignee 與 label 可輔助顯示，但不能單獨代表 agent ownership，因為多個 agent 可能共用同一個 GitHub identity。

### Trusted coordination comments

只有**可信任 GitHub actor** 建立的 `<!-- agent-coordination:v1 -->` comment 才能影響 ownership。可信任 actor 必須至少符合其中一項：

- GitHub 將 comment author 標示為 repository `OWNER`、`MEMBER` 或 `COLLABORATOR`；
- repository permission evidence 可驗證該 actor 具有 write / maintain / admin 權限；
- actor 是 repository 明確核准、可驗證的 GitHub App / bot identity。

無法驗證 author trust 時，該 comment 只能視為一般 discussion，不得成為 live claim、winner、heartbeat 或 release record。若目前工具無法取得足以驗證 claim author 的 metadata，對新 implementation ownership 必須 fail closed，不得因未驗證 comment 猜測 owner。

### Current state reduction

Canonical worker identity 是 `(agent, session)`，不是只有 `session`。對每個 trusted `(agent, session)`，先依 GitHub `created_at`（必要時再依 comment id）選出**最新一筆 coordination comment**作為該 worker 的 authoritative current state；更早的 `claim: active` 不得覆蓋較新的 `claim: released` / `inactive` / `done` / `blocked`。

只有 latest trusted comment 同時滿足 `claim: active` 且 lease 未過期的 worker，才算 live claim。

每次開始工作前，agent 必須重新檢查：

1. issue 是否已有 live claim；
2. 是否已有 linked / active implementation PR；
3. branch 是否已有近期 commit；
4. CI / review 是否顯示另一個 session 正在推進；
5. dependency / roadmap 是否允許目前 issue 開始。

若另一個 agent/session 有 live claim，不得平行實作相同 issue，也不得另開重複 PR。

## Agent identity

至少區分：

- `manual-chat`：使用者直接開啟的互動式 Chat session；
- `side-project-autopilot`：跨 side-project 的排程 agent；
- 其他 agent 使用穩定、可辨識的名稱，例如 `codex`。

`agent` 表示 worker 類型；`session` 用來區分同一類 worker 的不同執行實例。Session id 應盡量全域唯一；ownership 比較一律使用完整 `(agent, session)` identity。

## Claim format

認領 issue 時新增 comment：

```text
<!-- agent-coordination:v1 -->
agent: side-project-autopilot
session: <stable-session-or-run-id>
status: claimed
claim: active
branch: <branch-or-none>
pr: <pr-number-or-none>
head: <sha-or-none>
heartbeat: <ISO-8601 UTC timestamp>
next: <next concrete action>
```

### Post-claim ownership arbitration

發布 `claim: active` **不代表立即取得 ownership**。為避免兩個 session 同時 preflight、同時 claim 的 race condition，在建立 branch、PR、commit 或修改程式碼前，claimant 必須：

1. 發布 claim 後重新抓取 issue 最新 coordination comments；
2. 先依「Trusted coordination comments」與「Current state reduction」規則歸約每個 worker 的 current state；
3. 找出仍有效且互相衝突的 live claims；
4. 以 winning claim comment 的 GitHub `created_at` 較早者為唯一 winner；若 timestamp 相同，以完整 `(agent, session)` 字串 lexicographic ascending 作 deterministic tie-break；
5. winner 再確認自己仍是唯一有效 owner 後，才能開始 mutation；
6. loser 必須立即留下 `status: stale`、`claim: released` comment，記錄 winner identity，且不得建立 branch、commit 或 PR。

若平台無法取得足以完成 trusted-author validation、state reduction 或 deterministic arbitration 的 metadata，fail closed：不得開始新的 implementation。

若 issue 只是等待 dependency，不應宣告 active implementation claim；使用：

```text
<!-- agent-coordination:v1 -->
agent: side-project-autopilot
session: <id>
status: waiting-dependency
claim: inactive
heartbeat: <ISO-8601 UTC timestamp>
next: <dependency that must complete>
```

## Progress / heartbeat

每個 material state transition 都必須回寫 issue comment，不需要為每個小動作洗版。至少涵蓋：

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

建議格式：

```text
<!-- agent-coordination:v1 -->
agent: <agent>
session: <session>
status: review
claim: active
branch: <branch>
pr: <number>
head: <sha>
validation: <canonical CI/test evidence>
heartbeat: <ISO-8601 UTC timestamp>
next: <next concrete action>
```

Commit、CI 完成、review disposition、PR merge 等重要事件都應更新 heartbeat / progress。除此之外，任何 `claim: active` 的 session 即使狀態沒有改變，也必須至少每 **2 小時**更新一次 heartbeat-only progress comment；長時間本地修改、等待外部 reviewer 或長任務都不例外。這個週期刻意短於 8 小時 lease，避免仍在工作的 session 因沒有 material transition 被誤判 stale。

## Claim lease / stale takeover

預設 claim lease 為 **8 小時**。Lease 只由該 trusted `(agent, session)` 最新 coordination heartbeat 續期；排程 CI、review bot、其他 reviewer、其他 session 的 branch/PR/issue 活動都**不能**延長 owner lease。

當最新 trusted owner heartbeat 超過 8 小時時，新的 agent 才能進入 takeover preflight。接手前仍應檢查 linked PR、branch commits、review replies 與 issue discussion，以辨識是否存在**可歸因於原 owner identity**、但漏寫 coordination heartbeat 的近期 mutation；若存在，可暫緩最多一個 heartbeat cadence（2 小時）並要求 coordination state 修復。純 CI、bot review 或其他人的活動只提供 context，不是 takeover veto，也不得讓 abandoned claim 永久存活。

確認 lease stale 後，新的 agent 留下 takeover claim，記錄舊 owner、過期 heartbeat 與 evidence，並再次完成 post-claim ownership arbitration。

## Existing PR wins / adoption

若 issue 已有 active implementation PR，原則是**沿用該 PR，而不是另開重複 PR**。

- 若該 PR 仍有有效 live owner claim，新 claimant 必須退出。
- 若該 PR 的 owner 已 `released`、inactive、stale，或 PR 沒有有效 owner，eligible agent 可以用 `status: claimed`、`claim: active`、`pr: <existing PR>`、`next: adopt existing PR` 宣告 adoption，完成 trusted post-claim arbitration 後接手該既有 PR / 可寫 branch。
- 若 branch 不可寫或 ownership 無法安全驗證，記錄 blocker / release claim；不得用另開 duplicate implementation PR 規避。

除非 repository contract 或 issue 明確允許 stacking，否則每個 repository 優先維持一個 active implementation PR。

## Work selection priority

Autopilot 選工作時使用以下優先順序，仍須服從 repo roadmap / dependency / issue priority：

1. **可由 agent 推進**的 broken active implementation PR；
2. active implementation PR 上的 actionable review feedback；
3. active implementation PR 上的 failing canonical CI；
4. 其他**未被 human-only blocker 阻擋**的 active implementation PR，包括 CI/review 已乾淨、等待 merge 的 merge-ready PR；
5. 由**目前 agent/session 自己持有 live claim**的 unfinished issue；
6. roadmap blocker / dependency-unblocking issue；
7. 其他沒有 live claim 的 eligible issue。

由其他 agent/session 持有 live claim 的 issue 必須 skip，不屬於第 5 項。需要 secret、permission、irreversible product decision 或其他 human-only input 的 `blocked` 工作必須使用 `claim: inactive`，並從 active-PR work queue 排除直到 blocker 被清除；留下 blocker evidence 後繼續其他 eligible work。

不得因為 issue number 較小、建立時間較早或單純標成高 priority，就跳過正在收尾的 active implementation PR。

## Release claim

以下任一情況都必須 release claim：

- PR merged；
- issue completed / closed；
- linked PR closed without merge；
- implementation abandoned；
- PR/branch 被 supersede；
- 進入需要 human-only input 的 blocked 狀態。

完成時留下：

```text
<!-- agent-coordination:v1 -->
agent: <agent>
session: <session>
status: done
claim: released
pr: <number-or-none>
head: <merged-head-or-merge-sha>
heartbeat: <ISO-8601 UTC timestamp>
next: <next eligible issue or none>
```

若 linked PR 關閉 / abandoned / superseded 但 issue 尚未完成，使用 `claim: released` 並明確寫 `next: issue remains eligible for reassignment`，不得把舊 claim 留成 live。

若因真正的人類-only blocker 停止，使用 `status: blocked`、`claim: inactive`，清楚寫出需要的 secret、permission、irreversible product decision 或其他無法由 agent 解決的輸入；不要假裝完成，也不要無限持有 claim。

## Rules for manual Chat sessions

手動 Chat 在開始 coding 前也必須讀取最新 coordination comment，並遵守相同 trusted-author validation、state reduction 與 post-claim arbitration。看到 `side-project-autopilot` 的 live claim 時應避免撞工；反之 Autopilot 看到 `manual-chat` live claim 時必須跳過該 issue，除非該 claim 已依 stale takeover 規則失效。
