# GitHub Agent Coordination Protocol

本文件定義 coding agents 在 GitHub issue / PR 上協作時的共用互斥與進度紀錄規則。GitHub 是跨 Chat、排程與 coding agent 的協作 source of truth；conversation history 不能取代 GitHub 上可驗證的狀態。

## 適用範圍

所有會認領 GitHub issue、建立 branch / PR、修改程式碼或處理 review 的 agent 都適用，包括手動 Chat session、排程 Autopilot、Codex 或其他 coding agent。

Repository 的 `AGENTS.md`、local `AGENTS.md`、active roadmap、issue acceptance criteria / non-goals 與 repo-local skills 仍具有原本的 authority。本文件只定義「誰正在做、如何避免撞工、如何留下可接手狀態」。Root `AGENTS.md` 與 workflow skill 必須明確引用本文件，避免 agent 只依 canonical onboarding 卻漏掉 coordination preflight。

## Canonical coordination state

Issue comment 是 claim / progress 的 canonical coordination record。Assignee 與 label 可輔助顯示，但不能單獨代表 agent ownership，因為多個 agent 可能共用同一個 GitHub identity。

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

`agent` 表示 worker 類型；`session` 用來區分同一類 worker 的不同執行實例。

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
2. 找出仍有效且互相衝突的 active claims；
3. 以 GitHub comment `created_at` 較早者為唯一 winner；若 timestamp 相同，以完整 `session` 字串 lexicographic ascending 作 deterministic tie-break；
4. winner 再確認自己仍是唯一有效 owner 後，才能開始 mutation；
5. loser 必須立即留下 `status: stale`、`claim: released` comment，記錄 winner session，且不得建立 branch、commit 或 PR。

若 preflight 時已存在 linked active implementation PR，`Existing PR wins` 優先，新的 claimant 必須退出。若平台無法取得足以做 deterministic arbitration 的 comment metadata，fail closed：不得開始 implementation。

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

預設 claim lease 為 **8 小時**。

超過 8 小時沒有 coordination heartbeat 時，不得直接搶工作。接手前必須再確認：

- linked PR 沒有更新；
- branch 沒有新 commit；
- CI / review 沒有新活動；
- issue 沒有較新的進度訊息。

只有以上都沒有活動時，新的 agent 才能留下 takeover comment，將舊 claim 視為 stale，並記錄接手理由與新的 session identity；takeover claimant 仍須完成 post-claim ownership arbitration。

## Existing PR wins

若 issue 已有 active implementation PR，優先推進該 PR 的 CI / review / merge loop，而不是另開新 branch 或新 PR。除非 repository contract 或 issue 明確允許 stacking，否則每個 repository 優先維持一個 active implementation PR。

## Work selection priority

Autopilot 選工作時使用以下優先順序，仍須服從 repo roadmap / dependency / issue priority：

1. **可由 agent 推進**的 broken active implementation PR；
2. active implementation PR 上的 actionable review feedback；
3. active implementation PR 上的 failing canonical CI；
4. 其他 active implementation PR，包括 CI/review 已乾淨、等待 merge 的 merge-ready PR；
5. 由**目前 agent/session 自己持有 live claim**的 unfinished issue；
6. roadmap blocker / dependency-unblocking issue；
7. 其他沒有 live claim 的 eligible issue。

由其他 agent/session 持有 live claim 的 issue 必須 skip，不屬於第 5 項。需要 secret、permission、irreversible product decision 或其他 human-only input 的 `blocked` 工作必須使用 `claim: inactive`，不應在每輪被當成最高優先工作；留下 blocker evidence 後繼續其他 eligible work。

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

手動 Chat 在開始 coding 前也必須讀取最新 coordination comment，並遵守相同 post-claim arbitration。看到 `side-project-autopilot` 的 live claim 時應避免撞工；反之 Autopilot 看到 `manual-chat` live claim 時必須跳過該 issue，除非該 claim 已依 stale takeover 規則失效。
