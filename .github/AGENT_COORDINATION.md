# GitHub Agent Coordination Protocol

本文件定義 coding agents 在 GitHub issue / PR 上協作時的共用互斥與進度紀錄規則。GitHub 是跨 Chat、排程與 coding agent 的協作 source of truth；conversation history 不能取代 GitHub 上可驗證的狀態。

## 適用範圍

所有會認領 GitHub issue、建立 branch / PR、修改程式碼或處理 review 的 agent 都適用，包括手動 Chat session、排程 Autopilot、Codex 或其他 coding agent。

Repository 的 `AGENTS.md`、local `AGENTS.md`、active roadmap、issue acceptance criteria / non-goals 與 repo-local skills 仍具有原本的 authority。本文件只定義「誰正在做、如何避免撞工、如何留下可接手狀態」。

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

Commit、CI 完成、review disposition、PR merge 等重要事件都應更新 heartbeat / progress。

## Claim lease / stale takeover

預設 claim lease 為 **8 小時**。

超過 8 小時沒有 coordination heartbeat 時，不得直接搶工作。接手前必須再確認：

- linked PR 沒有更新；
- branch 沒有新 commit；
- CI / review 沒有新活動；
- issue 沒有較新的進度訊息。

只有以上都沒有活動時，新的 agent 才能留下 takeover comment，將舊 claim 視為 stale，並記錄接手理由與新的 session identity。

## Existing PR wins

若 issue 已有 active implementation PR，優先推進該 PR 的 CI / review / merge loop，而不是另開新 branch 或新 PR。除非 repository contract 或 issue 明確允許 stacking，否則每個 repository 優先維持一個 active implementation PR。

## Work selection priority

Autopilot 選工作時使用以下優先順序，仍須服從 repo roadmap / dependency / issue priority：

1. broken / blocked active PR；
2. actionable review feedback；
3. failing canonical CI；
4. live claimed unfinished issue；
5. roadmap blocker / dependency-unblocking issue；
6. 其他沒有 live claim 的 eligible issue。

不得因為 issue number 較小、建立時間較早或單純標成高 priority，就跳過正在收尾的 active PR。

## Release claim

PR merge 或 issue 完成後，留下最後狀態：

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

若因真正的人類-only blocker 停止，使用 `status: blocked`，清楚寫出需要的 secret、permission、irreversible product decision 或其他無法由 agent 解決的輸入；不要假裝完成，也不要無限持有 claim。

## Rules for manual Chat sessions

手動 Chat 在開始 coding 前也必須讀取最新 coordination comment。看到 `side-project-autopilot` 的 live claim 時應避免撞工；反之 Autopilot 看到 `manual-chat` live claim 時必須跳過該 issue，除非該 claim 已依 stale takeover 規則失效。
