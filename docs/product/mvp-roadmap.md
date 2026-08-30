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

目前 **Active Slice：MVP 1 — 使用者提供一則 Meaningful Activity**。

建議順序：

1. **MVP 0 — 建立 Identity 與 Close-friend Relationship**
2. **MVP 1 — 使用者提供一則 Meaningful Activity** ← Active
3. **MVP 2 — Activity 成為 Derived Social Context**
4. **MVP 3 — Privacy 決定 Specific Friend 可以知道什麼**
5. **MVP 4 — Friend 收到有用的 Friend Pulse**
6. **MVP 5 — Context 幫助開始真實 Conversation**
7. **MVP 6 — 第一個 External Integration 自動貢獻 Activity**

這個順序代表目前的最小驗證路徑，不是永久 roadmap。

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

- 作為使用者，看到朋友的一則 context 後，我可以快速得到一個自然、不尷尬、與該 context 有關的 conversation starter。
- 作為使用者，我可以用「聊聊這個」之類的 CTA 把 context 轉成真實互動意圖。

## Use Cases / Interactions

### Generate Conversation Starter

針對 viewer 已有權看到的 Context Projection 產生自然且相關的 starter。若無法產生具體、可使用的 starter，不能只回傳空字串或 generic placeholder 來視為成功。

### Start Conversation Intent

使用者表達「我要拿這個 context 去跟朋友聊」的 product intent。

### Record Lightweight Outcome

若不增加過多摩擦，可以記錄使用者是否點擊 / 採用 conversation CTA，作為產品驗證訊號。

Kin 若無法觀測外部聊天是否真的開始，不得把 intent 記錄成 confirmed conversation。

## Domains Involved

- Conversation / Interaction
- Privacy & Sharing
- Social Context

AI 如用於 starter generation，仍是 outer adapter。

## Expected Domain Responsibilities

### Conversation / Interaction

負責：

- conversation-starting intent
- context-based starter concept
- MVP 所需的 lightweight interaction outcome

不負責：

- friend visibility
- private chat history
- generic messaging system

## Candidate Commands

- `ExpressConversationIntentFromContext`
- `RecordConversationIntentOutcome`

## Candidate Queries

- `GetConversationStarter`

## Candidate Domain Events

- `ConversationIntentExpressedFromContext`

若未來真的能確認外部 conversation 已開始，再另外定義 confirmed outcome event；本 MVP 不過度宣稱。

## Acceptance Criteria

- [ ] Starter 只能使用 viewer 已授權可見的 Context Projection。
- [ ] Starter 不可藉由 prompt 還原更敏感的 raw Social Context / Activity。
- [ ] 對可用的 Context Projection，系統能回傳一個非空白、非 generic placeholder、且與該 projection 明確相關的 conversation starter。
- [ ] 使用者可以從 Pulse/context 明確進入 conversation CTA。
- [ ] 至少能收集一個 lightweight signal 判斷 context 是否促成 interaction intent。
- [ ] Product validation 必須能區分「starter 有生成」與「starter 對使用者真的有幫助 / 被採用」。
- [ ] 未能觀測外部聊天開始時，不得將 CTA click / intent 計為 confirmed conversation started。
- [ ] 不要求 Kin 自己成為 messaging app。

## Non-goals

- 完整 chat product
- 私訊紀錄 ingestion
- AI friend
- relationship coaching
- conversation transcript analysis
- 長期 Social Memory

## Dependencies

- MVP 4 的 Friend Pulse / Context Projection。

## Slice Completion Signal

當團隊可以驗證「projection-based starter 是否真的降低開口摩擦」，並能量測 conversation intent，而不把 intent 誤當 confirmed conversation 時，核心產品 loop 才具備最小驗證能力。

---

# MVP 6 — 第一個 External Integration 自動貢獻 Activity

## Goal

在核心 friendship-context loop 已能驗證後，選擇一個低風險、可明確授權的 provider，自動產生 Activity，降低 manual contribution friction。

## Validation Hypothesis

若核心 loop 有價值，下一個主要 friction 會是「使用者必須記得主動提供 Activity」。一個適合的 external integration 可以測試 ambient / low-friction ingestion 是否提升 Kin 的持續價值。

## Primary Actors

- User
- External Activity Provider

## User Stories

- 作為使用者，我可以明確連接一個 provider。
- 作為使用者，我知道 Kin 會讀取哪一類 signal。
- 作為使用者，我可以停止連接，未來不再自動取得新 Activity。

## Use Cases / Interactions

### Connect Provider

使用者明確授權一個 provider connection。

### Sync Authorized Activity

Integration adapter 取得 provider payload，在 boundary normalization / validation 後，呼叫既有 Activity contribution application flow。

同一 external item 必須具有可穩定識別的 provenance / provider item identity，使 overlapping pages、retry 或 timeout recovery 不會重複產生相同 Kin Activity。

### Disconnect Provider

停止後續 sync；是否刪除既有 derived data 應依獨立 privacy / retention policy 決定，不在此 slice 偷渡假設。

## Domains Involved

- Integration
- Activity
- Identity

後續既有 Social Context / Privacy / Pulse flow 應重用，不為 provider 建另一套 domain path。

## Expected Domain Responsibilities

### Integration

負責：

- provider connection state
- external authorization lifecycle
- checkpoint / sync concerns
- stable external provenance identity
- repeated-sync idempotency / deduplication
- provider error translation
- provider-specific payload normalization boundary

### Activity

只接收 provider-independent Kin Activity contract。

## Candidate Commands

- `ConnectActivityProvider`
- `SyncProviderActivities`
- `DisconnectActivityProvider`

## Candidate Queries

- `GetProviderConnection`
- `ListConnectedProviders`
- `GetSyncStatus`

## Candidate Domain Events

- `ProviderConnected`
- `ProviderDisconnected`
- `ActivityIngested`

## Acceptance Criteria

- [ ] 使用者可以明確 connect / disconnect 第一個 provider。
- [ ] 系統只取得使用者授權範圍內的 signal。
- [ ] Provider-specific DTO 不會穿透進 Activity / Social Context domain model。
- [ ] Provider error 會在 Integration boundary 被 translate。
- [ ] 相同 provider item 在 overlapping page、重試或 timeout recovery 後不會重複建立 Kin Activity。
- [ ] Dedup / idempotency 依 stable provenance / provider item identity 判斷，而不是依脆弱的顯示文字比較。
- [ ] 自動 ingestion 會重用既有 Activity → Social Context → Privacy → Pulse flow。
- [ ] Disconnect 後不再取得新的 provider Activity。
- [ ] Integration 不會改變既有 Privacy rules。

## Non-goals

- 一次支援多個 providers
- generic plugin marketplace
- browser-wide surveillance
- full ChatGPT conversation history ingestion
- background sync optimization for every platform
- provider-specific recommendation engine

## Dependencies

- MVP 1 至 MVP 5 已經證明核心 loop 至少具有初步產品價值。

## Slice Completion Signal

當 automatic activity ingestion 能降低 contribution friction、重試不產生重複 Activity，同時不破壞 privacy / domain boundaries，MVP 才算完成第一輪 ambient-loop 驗證。

---

# Implementation Authorization Rules

## 下一個 Issue 應如何選

AI coding agent 不應直接把整份 roadmap 當成一張 implementation task。

每次只能從 **目前 Active Slice** 中挑一個最小 coherent vertical change，建立獨立 GitHub Issue。

合法順序：

`Root AGENTS.md → Active GitHub Issue → applicable local AGENTS.md → Current MVP Slice → Product Scope → relevant Skills → Implementation`

如果一個 proposed Issue 需要尚未 Active 或尚未完成 slice 才存在的 domain capability，該 Issue 應被視為 premature。

## Slice 可以被拆成多張 Issue

一個 slice 不必等於一張 Issue。

例如 MVP 3 可以拆成：

1. 定義最小 disclosure policy domain behavior。
2. 產生 relationship-specific Context Projection。
3. 建立 friend-visible query model。
4. 補 revocation flow。

但每張 Issue 都必須留下可驗證的 vertical progress，而不是只建立未被 use case 使用的 infrastructure。

## Active Slice 的變更規則

- 初始 Active Slice 是 MVP 0。
- 完成某張 Issue 不會自動切換 Active Slice。
- 當前 slice 全部 Acceptance Criteria 與 Slice Completion Signal 被驗證後，建立 roadmap PR 把 Active 標記移到下一 slice。
- 該 roadmap PR 本身也必須通過 repository review / merge gate。
- Agent 不得僅依 Issue 編號、unchecked checklist 或 dependency prose 自行推斷 Active Slice 已前進。

## 不應建立的 Issue 類型

除非有 active use case 明確需要，避免：

- 「先把所有 DB tables 建好」
- 「先完成整套 REST API」
- 「先建立 generic event bus」
- 「先導入 Kafka / NATS」
- 「先拆成 microservices」
- 「先建立 vector database」
- 「先完成完整 provider abstraction」
- 「先做完整 relationship hierarchy」

Infrastructure 必須服務目前 user interaction，而不是為假想未來預建。

---

# MVP Success Signals

MVP 最終不是以 feature count 驗收，而是要開始能回答：

- 使用者是否願意讓 meaningful activity 進入 Kin？
- Derived context 是否比 raw activity 更有價值？
- 使用者是否信任 relationship-aware disclosure？
- Friend Pulse 是否真的降低理解朋友近況的成本？
- Context 是否提升 conversation-start intent？
- Automatic ingestion 是否提升持續使用價值，而不是只增加資料量？

長期優先 metrics 仍應偏向：

- Conversations Started（只在真的能確認時使用）
- Conversation Intents Expressed
- Friendships Maintained
- Dormant Friendships Reactivated
- 使用者是否覺得更了解 close friends 最近在意什麼

DAU、session time、feed impression 不應成為 Kin MVP 的主要成功定義。

---

# 明確不在本 MVP Roadmap 的 Future Scope

以下能力仍保留在 Product Scope，但沒有被本 MVP 自動授權：

- 完整 Relationship Level hierarchy
- Friendship Drift Detection
- Weekly Friendship Digest
- Shared Rabbit Hole
- Social Memory
- AI Friendship Concierge
- Friend-aware AI Q&A
- rich notification strategy
- widgets / app intents / ambient surfaces
- 多 provider integrations
- public / acquaintance-scale social graph

若未來要實作，必須有新的 validation hypothesis、active GitHub Issue 與對應 architecture / privacy review。

---

## 文件閱讀順序

任何 coding agent 在開始工作前，必須先遵守 root agent contract，而不是先由本文件決定 scope。

閱讀與決策順序：

1. Root `AGENTS.md`
2. Current active GitHub Issue（含 Acceptance Criteria / Non-goals）
3. Nearest applicable local `AGENTS.md`（若存在）
4. 本文件中的 Current Active Slice
5. `docs/product/product-scope.md`
6. Relevant repo-local skills
7. Implementation

若內容衝突，以 root `AGENTS.md` 定義的 source-of-truth precedence 為準；本 roadmap 不可自行覆蓋 repository agent contract。