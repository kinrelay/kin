# GitHub Agent Coordination Protocol

本文件定義 coding agents 在 GitHub issue / PR 上協作時的共用互斥與進度紀錄規則。GitHub 是跨 Chat、排程與 coding agent 的協作 source of truth；conversation history 不能取代 GitHub 上可驗證的狀態。

## 適用範圍

所有會認領 GitHub issue、建立 branch / PR、修改程式碼或處理 review 的 agent 都適用，包括手動 Chat session、排程 Autopilot、Codex 或其他 coding agent。

Repository 的 `AGENTS.md`、local `AGENTS.md`、active roadmap、issue acceptance criteria / non-goals 與 repo-local skills 仍具有原本的 authority。本文件只定義「誰正在做、如何避免撞工、如何留下可接手狀態」。Root `AGENTS.md` 與 workflow skill 必須明確引用本文件，避免 agent 只依 canonical onboarding 卻漏掉 coordination preflight。

## Canonical coordination state

Issue comment 是 claim / progress 的 canonical coordination record。Assignee 與 label 可輔助顯示，但不能單獨代表 agent ownership，因為多個 agent 可能共用同一個 GitHub identity。

### Trusted coordination comments

只有**可信任 GitHub actor** 建立的 `<!-- agent-coordination:v1 -->` comment 才能影響 ownership。可信任 actor 必須能由 repository 權限或明確 allowlist 驗證：repository owner，或具有 `write` / `maintain` / `admin` 權限的 actor，或 repository 明確核准、可驗證的 GitHub App / bot identity。`author_association=MEMBER` / `COLLABORATOR` 等關聯資訊本身不足以證明目前具有寫入權，不得單獨建立 trust。

無法驗證 author trust 時，該 comment 只能視為一般 discussion，不得成為 live claim、winner、heartbeat 或 release record。若目前工具無法取得足以驗證 claim author 的 metadata，對新 implementation ownership 必須 fail closed，不得因未驗證 comment 猜測 owner。

### Current state reduction and lease anchor

Canonical worker identity 是 `(agent, session)`，不是只有 `session`。對每個 trusted `(agent, session)`，先依 GitHub server `created_at`（必要時再依 comment id）選出**最新一筆 coordination comment**作為該 worker 的 authoritative current state；更早的 `claim: active` 不得覆蓋較新的 `claim: released` / `inactive` / `done` / `blocked`。

Ownership acquisition metadata 必須與 latest state 分離。每段 continuous active lease 以「第一次成功發布且通過 trusted validation 的 `claim: active` comment」作 immutable **server-side lease anchor**。Anchor 的排序來源只能是 GitHub server `created_at`，timestamp 相同再比較 GitHub comment id；不得使用 claimant 自填 timestamp、本機時間或最新 heartbeat comment 的時間決定 ownership。後續 heartbeat / progress 必須重複完整 active state，並引用同一 `lease_anchor: <first-claim-comment-id-or-url>`；不得以 sparse heartbeat 讓 `claim: active` 消失，也不得改變 lease acquisition order。

只有 latest trusted comment 同時滿足 `claim: active`、引用目前有效 anchor，且 lease 未過期的 worker，才算 live claim。Release / inactive / stale / done / blocked、lease expiration，或完成 takeover 都是不可逆 lease boundary。舊 anchor 在 boundary 後永久失效；遲到 heartbeat 即使引用舊 anchor也不能復活 ownership。原 session 若之後重新工作，必須發布新 claim、取得新 server-side anchor、重新 arbitration。

每次開始工作前，agent 必須重新檢查：

1. issue 是否已有 live claim；
2. 是否已有 linked / active implementation PR；
3. branch 是否已有近期 commit；
4. CI / review 是否顯示另一個 session 正在推進；
5. dependency / roadmap 是否允許目前 issue 開始；
6. 是否存在尚未清除的 `status: blocked` human-only blocker。

若另一個 agent/session 有 live claim，不得平行實作相同 issue，也不得另開重複 PR。

## Agent identity

至少區分：

- `manual-chat`：使用者直接開啟的互動式 Chat session；
- `side-project-autopilot`：跨 side-project 的排程 agent；
- 其他 agent 使用穩定、可辨識的名稱，例如 `codex`。

`agent` 表示 worker 類型；`session` 用來區分同一類 worker 的不同執行實例。Session id 應盡量全域唯一；ownership 比較一律使用完整 `(agent, session)` identity。

## Claim format

第一次認領 issue 時新增 comment：

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

發布後必須重新抓取該 comment 的 GitHub server `created_at` 與 comment id/url，將它固定為這段 lease 的 server-side anchor。後續 active progress 必須帶完整 state，例如：

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

### Post-claim ownership arbitration

發布 `claim: active` **不代表立即取得 ownership**。為避免兩個 session 同時 preflight、同時 claim 的 race condition，在建立 branch、PR、commit 或修改程式碼前，claimant 必須：

1. 發布 claim 後重新抓取最新 coordination comments，以及剛發布 comment 的 GitHub server metadata；
2. 先依「Trusted coordination comments」與「Current state reduction and lease anchor」歸約每個 worker；
3. 找出仍有效且互相衝突的 live claims，解析各自的有效 server-side lease anchor；
4. anchor 的 GitHub `created_at` 較早者為唯一 winner；timestamp 相同時比較 GitHub comment id；完全無法排序時 fail closed。不得用 claimant 自填 `claimed_at`、heartbeat 或本機時間決定 winner；
5. winner 再確認自己仍是唯一有效 owner 後，才能開始 mutation；
6. loser 必須立即留下 `status: stale`、`claim: released` comment，記錄 winner identity，且不得建立 branch、commit 或 PR。

舊資料若沒有 `lease_anchor`，只能用該 session 目前 continuous active lease 的第一筆 trusted `claim: active` comment 的 GitHub server `created_at` + comment id 作 migration fallback；不得用最新 heartbeat 時間。若平台無法取得足以完成 trusted-author validation、state reduction 或 deterministic arbitration 的 metadata，fail closed：不得開始新的 implementation。

若 issue 只是等待 dependency，不應宣告 active implementation claim；使用 `status: waiting-dependency` 與 `claim: inactive`。

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

`red` 是 Kin strict TDD 專用狀態：只有在 issue / repo testing contract 要求 Red → Green，且 regression test 已在正確原因上證明失敗時才能使用。它不是 generic「有錯」狀態，也不能在沒有客觀 failing test evidence 時宣告。

任何 `claim: active` progress / heartbeat 都必須重複完整 claim state並引用同一個有效 `lease_anchor`。Commit、CI 完成、review disposition、PR merge 等重要事件應更新 heartbeat / progress；除此之外，active session 即使狀態沒有改變，也必須至少每 **2 小時**更新一次 heartbeat-only progress comment。Heartbeat 只代表最近活動，不改 lease acquisition order，也不能延長已過期或已被 takeover 的舊 lease。

## Claim lease / stale takeover

預設 claim lease 為 **8 小時**。Lease 只由該 trusted `(agent, session)` 在 lease 尚未過期前發布、且引用有效 `lease_anchor` 的最新 coordination heartbeat 續期；排程 CI、review bot、其他 reviewer、其他 session 的 branch/PR/issue 活動都不能延長 owner lease。

一旦有效 heartbeat 超過 8 小時，該 lease 進入 terminal expired state；之後引用舊 anchor 的遲到 heartbeat不得復活。新的 agent 進入 takeover preflight 時，仍需檢查 linked PR、branch commits、review replies 與 issue discussion，以辨識是否存在**可歸因於原 owner identity**、但漏寫 coordination heartbeat 的近期 mutation；若存在，可暫緩最多一個 heartbeat cadence（2 小時）並要求 coordination state 修復。純 CI、bot review 或其他人的活動只提供 context，不是 takeover veto。

只有 recent-activity hold 已結束、inactivity preflight 通過並確認 takeover safe 後，新的 agent 才能發布 takeover claim、取得新的 server-side lease anchor、記錄舊 owner / expired evidence，並重新 arbitration。尚在 hold 期間不得 adoption existing PR 或推送既有 branch。

## Existing PR wins / adoption

若 issue 已有 active implementation PR，原則是**沿用該 PR，而不是另開重複 PR**；但既有 PR 不代表可無條件接管。

- 若該 PR 仍有有效 live owner claim，新 claimant 必須退出。
- 若 current state 是尚未清除的 human-only `status: blocked`，不得選擇或 adoption 該 PR，直到 blocker 明確清除。
- 若 PR author / trusted owner 已留下明確 handoff、release 或 inactive 記錄，eligible agent 可對 existing PR 發 claim，取得新 server-side anchor並完成 arbitration 後接手可寫 branch。
- 若 owner lease expired / stale，必須先完成上一節 stale-takeover inactivity preflight與 recent-activity hold；只有確認 safe 後才能 adoption。
- 若 PR 沒有 coordination metadata，**不得**直接視為 ownerless。先把 PR author / branch owner 視為可能 active，檢查 author commits、review replies、PR/issue comments、CI-triggering pushes 等可歸因活動；只有明確 handoff/consent，或至少超過 8 小時且完整 inactivity preflight 通過後，才可 adoption。
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
7. 其他沒有 live claim、沒有尚未清除 `status: blocked`、且依 active roadmap / dependencies 合法的 eligible issue。

由其他 agent/session 持有 live claim 的 issue 必須 skip，不屬於第 5 項。需要 secret、permission、irreversible product decision 或其他 human-only input 的 `blocked` 工作必須使用 `claim: inactive`；在 blocker 清除前，它同時從 active-PR queue、existing-PR adoption 與 new-issue eligibility 排除。留下 blocker evidence 後繼續其他 eligible work。

不得因為 issue number 較小、建立時間較早或單純標成高 priority，就跳過正在收尾的 active implementation PR。

## Release claim

以下任一情況都必須 terminate / release current lease：

- PR merged；
- issue completed / closed；
- linked PR closed without merge；
- implementation abandoned；
- PR/branch 被 supersede；
- lease expired / completed takeover；
- 進入需要 human-only input 的 blocked 狀態。

完成時留下：

```text
<!-- agent-coordination:v1 -->
agent: <agent>
session: <session>
status: done
claim: released
lease_anchor: <current anchor when applicable>
pr: <number-or-none>
head: <merged-head-or-merge-sha>
heartbeat: <ISO-8601 UTC timestamp>
next: <next eligible issue or none>
```

若 linked PR 關閉 / abandoned / superseded 但 issue 尚未完成，使用 `claim: released` 並明確寫 `next: issue remains eligible for reassignment`，不得把舊 claim 留成 live。Expired/taken-over lease 的舊 anchor 不得再次使用；同一 session 後續若重回工作，必須取得新 anchor。

若因真正的人類-only blocker 停止，使用 `status: blocked`、`claim: inactive`，清楚寫出需要的 secret、permission、irreversible product decision 或其他無法由 agent 解決的輸入；不要假裝完成，也不要無限持有 claim。

## Rules for manual Chat sessions

手動 Chat 在開始 coding 前也必須讀取最新 coordination comment，並遵守相同 trusted-author validation、server-side lease anchor、state reduction 與 post-claim arbitration。看到 `side-project-autopilot` 的 live claim 時應避免撞工；反之 Autopilot 看到 `manual-chat` live claim 時必須跳過該 issue，除非該 claim 已依 stale takeover 規則失效。
