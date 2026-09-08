# Kin MVP Roadmap

## 文件目的

這份文件把 Kin 的 MVP 定義成一系列有順序的 **vertical product slices**。

每個 slice 都必須能回答一個具體的使用者價值問題，而不是單純完成某個技術層。實作順序應從 user interaction 與 domain behavior 出發，再逐步落到 application、ports、adapters 與 delivery。

核心驗證問題：

> 如果 Kin 能替我們維持關於最親近朋友的輕量、經授權 context，是否能幫助真實友情保持活性，並降低重新開啟對話的摩擦？

本文件只定義 MVP 驗證順序與目前允許進入 MVP 的 product scope。它不定義 database schema、HTTP API、provider SDK、queue、deployment topology，也不代表 Product Scope 中所有 long-term capabilities 都已獲得 implementation authorization。

---

## MVP slicing 原則

### 1. 每個 slice 必須驗證 user outcome

不以「建立 database」、「完成 backend」、「完成 frontend」作為 slice。

一個合法 slice 應能用真實 interaction 描述，例如：

- 兩個人可以建立一段 close-friend relationship。
- 一位使用者可以提供一則有意義的 Activity。
- 一位朋友可以安全地看到一則 relationship-specific Context Projection。

### 2. Domain / interaction first

每個 slice 的設計順序：

`User Outcome → Use Case → Domain Responsibility → Command / Query → Ports → Adapter / Delivery`

不要反過來從 table、route、SDK 或 UI component 推導 use case。

### 3. Slice 只授權當下需要的能力

Product Scope 是 long-term conceptual map，不是 implementation backlog。

若某個 future capability 沒有被目前 slice 明確需要，例如：

- 完整 Relationship Level hierarchy
- Friendship Drift Detection
- Social Memory
- Weekly Friendship Digest
- Shared Rabbit Hole
- AI Friendship Concierge

就不能只因為它存在於 Product Scope 而順手實作。

### 4. Privacy 先於 Relevance

任何 friend-visible output 都必須先產生 relationship-specific `Context Projection`，再做 ranking、Friend Pulse、Conversation Support 或 Notification。

`Social Context ≠ Context Projection`

Relevance 不得直接處理未經 Privacy & Sharing 投影的 Social Context。

### 5. AI / provider 永遠是 outer adapter

MVP 可以使用 AI 或 external provider，但它們不能成為 domain authority。

Provider / LLM output 必須先在 adapter boundary normalization、validation、error translation，再轉成 Kin inner contract。

### 6. Manual / explicit flow 先於 automation

第一輪 MVP 應優先驗證核心 friendship-context loop，而不是先投入 external integration complexity。

因此 automatic provider ingestion 放在核心 manual/explicit flow 被證明有產品價值之後。

---

## MVP 全貌與 Active Slice

目前 **Active Slice：MVP 5 — Context 幫助開始真實 Conversation**。

建議順序：

1. **MVP 0 — 建立 Identity 與 Close-friend Relationship**
2. **MVP 1 — 使用者提供一則 Meaningful Activity**
3. **MVP 2 — Activity 成為 Derived Social Context**
4. **MVP 3 — Privacy 決定 Specific Friend 可以知道什麼**
5. **MVP 4 — Friend 收到有用的 Friend Pulse**
6. **MVP 5 — Context 幫助開始真實 Conversation** ← Active
7. **MVP 6 — 第一個 External Integration 自動貢獻 Activity**

這個順序代表目前的最小驗證路徑，不是永久 roadmap。

### MVP 2 → MVP 3 transition evidence

MVP 2 的 completion signal 已由 #33、#34、#35 與 #48 覆蓋，並經 #52、#55 的後續驗證與修正完成收斂：authorized Activity 會先經 significance / suppression，再產生並驗證可追溯 provenance、非 raw replay 的 private Social Context；owner 可透過 purpose-built read model 查看結果，且既有實作沒有提前建立 friend-visible disclosure。基於這些 evidence 與 MVP 2 Slice Completion Signal，MVP 3 已完成其 Active Slice transition。

### MVP 3 → MVP 4 transition evidence

MVP 3 的 completion signal 已由 #60 / PR #61、#62 / PR #63 與 #64 / PR #65 覆蓋：Context Owner 可建立、修改、降低或撤銷 disclosure；authenticated active-friend read boundary 只回傳 relationship-specific `Context Projection`；pending delivery 會綁定 privacy / relationship revision，並在 dispatch-time 重新授權、重新投影或取消，避免 stale / over-detailed payload 進入可送出狀態。#70 已逐項對照 MVP 3 全部 Acceptance Criteria 與 Slice Completion Signal，未發現阻擋核心 hypothesis 的 implementation gap，因此 MVP 4 現在成為唯一 Active Slice。

### MVP 4 → MVP 5 transition evidence

MVP 4 的 completion signal 已由 #72 / PR #73 與 #74 的 closure reconciliation 覆蓋：authenticated viewer 可以取得 active friend 的 permissioned Friend Pulse；read side 只對 relationship-specific `Context Projection` 做 deterministic prioritization，限制為 1–3 個高訊號 item，排除 revoked / expired / suppressed context，且不讀取 privacy projection 前的敏感內容，也不依賴 chronological feed。#74 已逐項對照 MVP 4 的 6 項 Acceptance Criteria 與 Slice Completion Signal，並確認目前沒有其他 open product gap 阻擋「使用者能在極短時間理解朋友最近最值得知道什麼」的核心 hypothesis。因此 MVP 5 現在成為唯一 Active Slice；MVP 6 仍未授權。

### Active Slice 如何前進

- 只有目前 Active Slice 可以直接產生 implementation issues。
- 一個 slice 可以拆成多張 Issue；完成其中一張不代表 slice 自動完成。
- 只有當該 slice 的 Acceptance Criteria 與 Slice Completion Signal 都已透過實作與驗證滿足，且沒有仍阻擋核心 hypothesis 的 open issue，才可把下一個 slice 標成 Active。
- 進入下一個 slice 必須用獨立 docs / roadmap PR 明確更新此標記；不得由 agent 自行推測。
- 後續 slice 即使已寫在本文件中，在被標成 Active 前仍屬未授權 implementation scope。

---

# MVP 0 — 建立 Identity 與 Close-friend Relationship

## Goal

讓兩位 Kin 使用者可以形成一段明確、雙方同意、可被後續 privacy 與 context flow 引用的 close-friend relationship。

## Validation Hypothesis

如果 Kin 的核心價值建立在「特定朋友之間的 context continuity」，那系統首先必須能表達一段真實、雙方明確參與的 relationship，而不是 generic follower graph。

## Primary Actors

- User
- Friend / Invitee

## User Stories

- 作為 Kin 使用者，我可以建立基本 identity，讓朋友能辨識我是誰。
- 作為 Kin 使用者，我可以邀請另一位使用者成為 close friend。
- 作為受邀者，我可以明確接受邀請；不是受邀者的人不能替我接受。
- 作為 relationship participant，我可以知道目前這段 friendship 是否已成立。

## Use Cases / Interactions

### 建立 Identity

Actor 建立可參與 Kin relationship 的最小 identity。

### 發起 Friendship

一位 User 對另一位 distinct User 表達建立 relationship 的 intent。

### 接受 Friendship

只有該 invitation 指定的 invitee 可以接受；inviter 本人或不相關第三人嘗試接受必須失敗。接受後 relationship 才成為 active。

### 查看 Relationship State

參與者可以知道某位使用者是否已經是自己的 active friend。

## Domains Involved

- Identity
- Friendship

## Expected Domain Responsibilities

### Identity

負責：

- stable Kin user identity
- MVP 必要的最小 profile state

不負責：

- friendship lifecycle
- sharing policy

### Friendship

負責：

- relationship participants
- invitation / acceptance lifecycle
- invitation target identity
- active relationship state
- MVP 所需的 close-friend semantic

## Candidate Commands

- `CreateIdentity`
- `InviteFriend`
- `AcceptFriendship`

## Candidate Queries

- `GetMyIdentity`
- `GetFriendship`
- `ListMyFriends`

## Candidate Domain Events

- `IdentityCreated`
- `FriendshipInvited`
- `FriendshipCreated`

## Acceptance Criteria

- [ ] 兩位不同使用者可擁有可識別的 Kin identity。
- [ ] 一位使用者可對另一位 distinct 使用者發起 friendship intent。
- [ ] 未接受前不可視為 active friendship。
- [ ] 只有 invitation 指定的 invitee 可以執行 acceptance。
- [ ] inviter 自行接受自己的邀請必須失敗。
- [ ] 非 invitation participant 的第三人接受邀請必須失敗。
- [ ] 接受後雙方都能查詢到一致的 active relationship state。
- [ ] MVP 不依賴完整 social graph 或公開 follower model。

## Non-goals

- 完整 Relationship Level hierarchy
- contacts import
- friend recommendation
- block/report 系統
- relationship strength scoring
- public profiles

## Dependencies

無。

## Slice Completion Signal

當系統可以可靠表達「A 邀請 B，且只有 B 明確接受後，A 與 B 才成為 active close-friend relationship」時，MVP 0 才具備進入下一 slice 的條件。

---

# MVP 1 — 使用者提供一則 Meaningful Activity

## Goal

讓一位使用者能以低摩擦、明確授權的方式，提供一則 Kin 可以理解的 Activity。

## Validation Hypothesis

如果 friendship context 能帶來價值，Kin 必須先證明使用者願意讓某些 digital/life signals 進入系統，而且不需要先完成 external integrations。

## Primary Actors

- User

## User Stories

- 作為使用者，我可以提供一則最近在做、看、研究、收藏或喜歡的事情。
- 作為使用者，我知道這個 Activity 目前只是 private input，不會自動分享給朋友。

## Use Cases / Interactions

### Contribute Activity

使用者主動提供一則 Activity，包括足以讓 Kin 理解基本 meaning 的資訊。

### View My Activities

使用者可以確認自己曾提供哪些 Activity。

## Domains Involved

- Identity
- Activity

## Expected Domain Responsibilities

### Activity

負責：

- authorized contribution intent
- normalized Kin Activity
- provenance / timestamp
- Activity lifecycle 的最小必要狀態

不負責：

- friend visibility
- Context wording
- relevance ranking

## Candidate Commands

- `ContributeActivity`

## Candidate Queries

- `ListMyActivities`
- `GetActivity`

## Candidate Domain Events

- `ActivityContributed`

## Acceptance Criteria

- [ ] 使用者可以提供至少一種 MVP-defined Activity。
- [ ] Activity 與 contributing user 有明確 ownership。
- [ ] Activity 預設 private，不會因建立 friendship 就自動變成 friend-visible item。
- [ ] 系統能區分 raw contribution 與 Kin normalized Activity concept。
- [ ] 使用者能查看自己已提供的 Activity。

## Non-goals

- Spotify / YouTube / ChatGPT 自動同步
- browser extension
- Share Extension
- bulk import
- relevance ranking
- Social Context generation

## Dependencies

- MVP 0 的 Identity。

Friendship 可以存在，但不是 contribution 的必要條件。

## Slice Completion Signal

當使用者可以主動提供一則 private Activity，且系統沒有把 Activity 誤當 social post 時，即具備進入 MVP 2 的條件。

---

# MVP 2 — Activity 成為 Derived Social Context

## Goal

把一則或多則 authorized Activity 轉換成較高階、具社交意義的 Social Context，而不是一對一改寫成另一種 activity feed。

## Validation Hypothesis

如果 Kin 只是重新發布或輕量 paraphrase 每一筆 Activity，它仍會退化成 activity feed。差異化必須來自「辨識 significance、抑制低訊號，並從 signals 推導 meaning」。

## Primary Actors

- User
- System（透過 application orchestration）

## User Stories

- 作為使用者，我希望 Kin 理解我最近在意的主題，而不是逐筆重播我的行為。
- 作為使用者，我希望低訊號、重複或過度細碎的 Activity 不會自動變成 social content。

## Use Cases / Interactions

### Evaluate Activity Significance

Application / domain 先判斷一則或多則 Activity 是否具有足以形成 social meaning 的 signal。

低訊號、純重複或僅能逐字改寫的輸入應被 suppression / no-op，而不是強制產生 Context Candidate。

### Derive Context Candidate

只有通過最小 significance rule 的 authorized Activity，才可透過 domain rules 與必要的 `ContextGenerator` port 產生 `Context Candidate`。

### Validate / Promote Social Context

`Context Candidate` 必須經 Kin 自己的 validation rules 後，才可 promotion 成 `Social Context`。Validation 至少確認：

- 有可理解且非逐字重播的 social meaning。
- provenance 仍可追溯到 authorized Activity。
- 不包含 provider-specific raw payload。
- 若來自 LLM，已在 adapter boundary 完成 normalization / validation。

無法通過 validation 的 candidate 不得成為 Social Context。

### Review My Derived Context

Owner 可以查看 derived Social Context，以驗證 wording 與 meaning 是否合理。

## Domains Involved

- Activity
- Social Context

AI 如被使用，只是 outer adapter。

## Expected Domain Responsibilities

### Activity

提供 normalized signals 與 provenance。

### Social Context

負責：

- significance interpretation
- Context Candidate
- candidate validation / promotion
- Social Context lifecycle
- semantic meaning
- suppression
- abstract provenance

## Candidate Commands

- `GenerateContextFromActivity`
- `ValidateContextCandidate`
- `SuppressContext`

## Candidate Queries

- `ListMyContextCandidates`
- `ListMySocialContexts`
- `GetSocialContext`

## Candidate Domain Events

- `ContextCandidateGenerated`
- `SocialContextValidated`
- `ContextSuppressed`

## Acceptance Criteria

- [ ] 一則或多則 authorized Activity 可以在具有足夠 significance 時產生 derived Context Candidate。
- [ ] 低訊號、重複、純逐筆 paraphrase 或過度細碎的 Activity 不會被強制轉成 Context Candidate。
- [ ] Raw Activity 不會被直接當成 friend-visible context。
- [ ] Context Candidate 必須經 validation 才能成為 Social Context。
- [ ] 無法通過 validation 的 candidate 不得進入 Social Context state。
- [ ] Context wording 不包含 provider-specific payload shape。
- [ ] 若使用 LLM，provider output 必須在 adapter boundary 完成 normalization / validation。
- [ ] 無法通過 adapter validation 的 AI/provider output 不得進入 application/domain flow。
- [ ] 使用者可以看到目前為自己產生的 Social Context，以進行產品驗證。

## Non-goals

- friend-specific disclosure
- Friend Pulse
- notification
- complex recurring-interest engine
- multi-model orchestration
- semantic vector infrastructure
- automatic provider ingestion

## Dependencies

- MVP 1 的 Activity contribution。

## Slice Completion Signal

當 Kin 能把 private Activity 經 significance 判斷與 candidate validation，轉成較高階、可理解、但尚未對朋友揭露的 `Social Context` 時，MVP 2 才具備進入下一 slice 的條件。

---

# MVP 3 — Privacy 決定 Specific Friend 可以知道什麼

## Goal

讓同一份 Social Context 依 relationship 與 sharing policy，安全地產生 specific friend 可看的 `Context Projection`。

## Validation Hypothesis

使用者只有在相信 Kin 不會過度揭露時，才可能願意長期讓 meaningful signals 進入系統。Privacy 必須是 core product behavior，不是最後補上的 filter。

## Primary Actors

- Context Owner
- Friend Viewer

## User Stories

- 作為 Context Owner，我可以決定某個 context 是否允許朋友知道。
- 作為 Context Owner，只有我可以建立、修改或撤銷我的 disclosure decision。
- 作為 Friend，我只能看到這段 relationship 被允許看到的 detail level。
- 作為 Context Owner，我降低或撤銷分享後，未來 surface / notification 不應繼續洩漏舊的或過度詳細的資料。

## Use Cases / Interactions

### Define Sharing Decision

只有 authenticated Context Owner 可以對自己的 context / category / relationship 建立或修改 disclosure decision。

### Project Context For Friend

Privacy & Sharing 根據 Social Context、Friendship 與 policy 產生 relationship-specific Context Projection。

Friend-facing read side 必須從 authenticated caller identity 判斷 viewer，並確認 caller 是該 **active relationship** 的 participant；不得只信任 request 傳入的 viewer identity 或 relationship identifier。

### Revoke / Downgrade Disclosure

只有 authenticated Context Owner 可以撤銷既有 disclosure 或降低 sharing detail。撤銷或降級後，後續 read / surface 必須反映最新 policy。

若已有 pending notification intent，application orchestration 必須在實際 dispatch 前重新取得最新 Privacy & Sharing decision，並**重新計算當下 relationship-specific `Context Projection` / detail level**，不能只檢查 disclosure 是否仍為 valid。

- 每個 delivery intent 必須記錄產生目前 projection 時所依據的 privacy policy revision 與 relationship revision。
- Dispatch authorization 必須在同一個 atomic boundary 內完成「確認 revision 仍為最新」與「把 delivery intent 轉成可送出的 committed / dispatchable state」；不能先驗證、離開 atomic boundary，再直接送出舊 payload。
- 若 atomic boundary 內發現 privacy 或 relationship revision 已改變，舊 authorization 立即失效，必須重新取得最新 decision 並重新投影，或取消 delivery。
- 若 disclosure 在排程後、dispatch 前已 revoked / invalid，pending intent 必須取消或標記 invalid，不得送出。
- 若 disclosure 仍有效但 detail level 已降低，舊 intent 若包含超出最新 policy 的 projection，必須以最新較低 detail 的 projection 更新 delivery intent，或在無法安全更新時取消。
- Notification 本身不解讀 privacy policy，也不得沿用舊的 friend-visible payload；它只接受已在 atomic dispatch authorization boundary 中確認 revision 並重新授權 / 重新投影後的 delivery intent。

這個 contract 不指定 database transaction、lock、CAS 或 queue implementation；實作可以依當下 architecture 選擇機制，但不得讓「reauthorize → policy/relationship 變更 → send 舊 payload」的 race 成立。

## Domains Involved

- Friendship
- Social Context
- Privacy & Sharing

## Expected Domain Responsibilities

### Privacy & Sharing

負責：

- disclosure decision
- ownership authorization for policy mutation
- least-revealing valid projection
- relationship-specific visibility
- revocation / detail downgrade
- privacy policy revision
- dispatch-time projection re-evaluation contract

### Friendship

提供 relationship state / closeness input 與 relationship revision，但不直接決定 disclosure policy。

### Social Context

提供可被投影的 semantic context，不決定 specific friend 的 visibility。

## Candidate Commands

- `SetContextSharingPolicy`
- `RevokeContextSharing`

## Candidate Queries

- `GetContextProjectionForFriend`
- `ListVisibleContextsForFriend`

## Candidate Domain Events

- `ContextSharingPolicyChanged`
- `ContextDisclosureRevoked`

## Acceptance Criteria

- [ ] Friend 無法看到沒有明確 disclosure permission 的 Social Context。
- [ ] Friend-facing read side 使用 `Context Projection`，不是 raw Social Context。
- [ ] Friend-facing read side 必須驗證 authenticated caller identity，且 caller 必須是目標 active relationship 的 participant；不得信任 caller 可自行指定的 viewer / relationship id 來授權讀取。
- [ ] 同一份 Social Context 可以對不同 relationship 產生不同結果，至少支援「可見 / 不可見」或一個最小 detail-level variation。
- [ ] 只有 authenticated Context Owner 可以建立、修改、降低或撤銷該 context 的 sharing decision。
- [ ] Friend Viewer 嘗試替自己增加 visibility 必須失敗。
- [ ] 非 owner 的第三人修改 disclosure policy 必須失敗。
- [ ] Revocation 或 detail downgrade 會影響後續 query / surface。
- [ ] Pending notification intent 在 dispatch 前必須重新取得最新 privacy decision，並重新計算 relationship-specific `Context Projection` / detail level；不得直接沿用排程時保存的舊 projection。
- [ ] Delivery intent 必須綁定產生 projection 時的 privacy policy revision 與 relationship revision。
- [ ] Revision check 與 delivery intent 進入 committed / dispatchable state 的 transition 必須位於同一 atomic boundary；若 revision 不符，舊 payload 不得成為可送出狀態。
- [ ] 若 disclosure 已 revoked / invalid，pending intent 必須取消或失效且不可送出。
- [ ] 若 disclosure 仍有效但 detail level 已降低，pending intent 必須更新成最新較低 detail 的 projection，或在無法安全更新時取消；不得送出舊的較詳細 context。
- [ ] Acceptance test 必須覆蓋 race：在 reauthorization / projection 完成後、delivery commit 前插入 revoke 或 detail downgrade，驗證 revision mismatch 會觸發重新投影或取消，且舊 payload 不會被送出。
- [ ] Acceptance test 必須覆蓋 relationship revision race：若 relationship 在 reauthorization 與 delivery commit 之間失效或改變，舊 delivery intent 不得被送出。
- [ ] 沒有 permission 必須解讀為不可揭露。
- [ ] Privacy evaluation 必須發生在 Relevance / Friend Pulse 之前。

## Non-goals

- 完整 Acquaintance / Friend / Close Friend / Inner Circle hierarchy
- rule-builder UI
- machine-learned privacy policy
- organization / group sharing
- public sharing
- 指定 transaction / lock / CAS / queue 技術

## Dependencies

- MVP 0 的 Friendship
- MVP 2 的 Social Context

## Slice Completion Signal

當系統能可靠回答「對這位 specific friend，這份 context 現在到底能不能看、能看到多少」，只有 active relationship participant 能取得 projection、只有 Context Owner 能控制 disclosure，且任何尚未 dispatch 的 disclosure 都會綁定 privacy / relationship revision，並在同一 atomic boundary 內完成 revision validation 與 delivery state transition；任何中途 revocation、detail downgrade 或 relationship change 都會讓舊 authorization 失效並觸發重新投影或取消，確保過期或過度詳細的 payload 永遠不會進入可送出狀態時，MVP 3 才具備進入下一 slice 的條件。

---

# MVP 4 — Friend 收到有用的 Friend Pulse

## Goal

讓使用者可以看到某位朋友目前最值得知道的少量、permissioned context，而不是 activity feed。

## Validation Hypothesis

Kin 的核心使用體驗應降低「我不知道朋友最近在幹嘛」的成本。若只提供大量 context list，產品仍可能製造另一個需要刷的 feed。

## Primary Actors

- Friend Viewer

## User Stories

- 作為使用者，我可以快速知道一位 close friend 最近最值得知道的 1–3 件事。
- 作為使用者，我不需要閱讀對方完整 activity history。

## Use Cases / Interactions

### Get Friend Pulse

Read-side 先取得 viewer 有權看到的 Context Projections，再由 Relevance 做最小必要 prioritization，形成 Friend Pulse。

### Explain Pulse Item

若產品需要，可以提供簡短「為什麼現在顯示」的 explanation，但不得暴露 viewer 無權看到的 raw evidence。

## Domains Involved

- Privacy & Sharing
- Relevance
- Social Context
- Friendship

## Expected Domain Responsibilities

### Privacy & Sharing

先提供合法 Context Projection。

### Relevance

只對已可見的 projection 做：

- prioritization
- staleness suppression
- repetition suppression

### Friend Pulse

在 MVP 可先視為 application/read-model concept，不急著宣告獨立 bounded context。

## Candidate Commands

本 slice 可能不需要新的 write command。

## Candidate Queries

- `GetFriendPulse`

## Candidate Domain Events

通常不需要因為單純 query 產生 domain event。

## Acceptance Criteria

- [ ] 使用者可以取得某位 active friend 的 Friend Pulse。
- [ ] Pulse 只包含該 viewer 目前有權看到的 Context Projection。
- [ ] Pulse 數量刻意保持少量，例如 1–3 個高訊號 item。
- [ ] Revoked / expired / suppressed context 不會出現在 Pulse。
- [ ] Relevance 不會讀取或重新揭露 privacy projection 前的敏感內容。
- [ ] 產品不需要 chronological feed 才能完成此 interaction。

## Non-goals

- infinite feed
- engagement ranking
- ads
- push notification
- Weekly Digest
- Shared Rabbit Hole
- Friendship Drift

## Dependencies

- MVP 3 的 Context Projection。

## Slice Completion Signal

當使用者可以在極短時間內理解「朋友最近最值得知道什麼」時，MVP 4 才具備進入下一 slice 的條件。

---

# MVP 5 — Context 幫助開始真實 Conversation

## Goal

驗證 Friend Pulse / Context Projection 是否真的能降低重新開口的摩擦，而不只是提供資訊。

## Validation Hypothesis

Kin 真正的 outcome 不是「看過 context」，而是讓真實 relationship 更容易產生 conversation intent，並最終促成真實對話。

## Primary Actors

- Friend Viewer

## User Stories

- 作為使用者，我看到朋友最近的 context 後，可以自然地找到一個開口方式。
- 作為使用者，我不希望 Kin 代替我聊天，而是幫我降低 conversation startup friction。

## Use Cases / Interactions

### Get Conversation Support

使用者從某個 permissioned Context Projection / Pulse Item 請求 conversation support。

第一版可以是 deterministic template 或非常薄的 `ConversationSupportGenerator` port；不需要先建立 agentic conversation system。

### Start Real Conversation

產品至少要能表達「使用者因這個 context 產生了 conversation intent」。

MVP 不一定需要內建 chat；可以是 copy suggestion、open external app、mark intent，或其他最小 interaction。

## Domains Involved

- Social Context
- Privacy & Sharing
- Relevance
- Conversation Support

## Expected Domain Responsibilities

### Conversation Support

負責：

- conversation starter intent
- context-aware but permission-bounded support
- 不越界補充未揭露資訊

## Candidate Commands

- `RequestConversationSupport`
- `RecordConversationIntent`

## Candidate Queries

- `GetConversationSupport`

## Candidate Domain Events

- `ConversationSupportRequested`
- `ConversationIntentRecorded`

## Acceptance Criteria

- [ ] Conversation support 只能基於 viewer 已有權看到的 Context Projection。
- [ ] Support 不會引入 projection 中不存在的敏感 detail。
- [ ] 使用者可以取得至少一個可自然開口的 suggestion / prompt。
- [ ] 產品可以記錄最小 conversation intent signal，例如 `opened` / `copied` / `started`。
- [ ] MVP 不要求 Kin 代替使用者進行 autonomous conversation。

## Non-goals

- autonomous messaging
- full in-app chat
- long-running conversation agent
- sentiment coaching
- relationship therapy
- automatic follow-up

## Dependencies

- MVP 4 的 Friend Pulse / Context Projection。

## Slice Completion Signal

當至少一位使用者能從 permissioned context 取得自然的 conversation support，並產生可觀察的 conversation intent signal 時，MVP 5 才具備進入下一 slice 的條件。

---

# MVP 6 — 第一個 External Integration 自動貢獻 Activity

## Goal

在核心 manual loop 已證明有價值後，驗證 external provider 是否能降低 Activity contribution friction，而不破壞 consent / privacy。

## Validation Hypothesis

若 Kin 需要長期運作，完全依賴 manual contribution 可能 friction 太高；但 automation 只有在核心 context → relationship loop 已有價值後才值得投入。

## Primary Actors

- User
- External Provider

## User Stories

- 作為使用者，我可以明確連接一個 provider。
- 作為使用者，我知道哪些 provider data 會進入 Kin。
- 作為使用者，我可以停止同步。

## Use Cases / Interactions

### Connect Provider

使用者明確授權一個 provider adapter。

### Import Provider Activity

Adapter 把 provider-specific payload normalization 成 Kin Activity candidate。

### Disconnect Provider

停止後續 ingestion；是否刪除歷史資料由獨立 policy 決定。

## Domains Involved

- Activity
- Integration

## Expected Domain Responsibilities

### Integration

負責：

- provider connection state
- consent boundary
- sync lifecycle

### Activity

仍只接受 normalized contribution，不知道 provider SDK。

## Candidate Commands

- `ConnectActivityProvider`
- `DisconnectActivityProvider`
- `ImportProviderActivity`

## Candidate Queries

- `ListConnectedProviders`

## Candidate Domain Events

- `ActivityProviderConnected`
- `ActivityProviderDisconnected`

## Acceptance Criteria

- [ ] 使用者必須明確 opt in 才會啟用 provider ingestion。
- [ ] Provider payload 在 adapter boundary normalization。
- [ ] Provider adapter failure 不會污染 domain state。
- [ ] 使用者可以停止後續 ingestion。
- [ ] External Activity 仍遵守既有 private-by-default 與 downstream privacy flow。

## Non-goals

- 多 provider 同時整合
- background sync optimization
- recommendation engine
- data resale
- provider-specific logic 進 domain

## Dependencies

- MVP 1 Activity
- 核心 context → privacy → pulse → conversation loop 已有初步產品證據。

## Slice Completion Signal

當一個 external provider 能在明確 consent 下穩定產生 normalized private Activity，且後續仍完整經過 Kin 既有 context / privacy pipeline 時，MVP 6 才具備完成條件。

---

## MVP 之外暫不實作

除非 roadmap 明確更新，以下能力不應進入目前 implementation scope：

- Relationship Level hierarchy
- Friendship Drift Detection
- Social Memory
- Weekly Friendship Digest
- Shared Rabbit Hole
- AI Friendship Concierge
- autonomous messaging
- multi-agent social orchestration
- public social graph
- monetization / ads
