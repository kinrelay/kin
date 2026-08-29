# Kin 產品範圍與 Domain Map

## 文件目的

這份文件以「完整，但刻意維持低解析度」的方式定義 Kin 的產品範圍與主要 Domain。

它的目的，是在選擇 database schema、HTTP API、provider SDK、queue、deployment topology 或其他 infrastructure 細節之前，先建立一份穩定的產品與 domain 概念地圖，讓人類與 AI coding agents 都能使用一致的語彙理解系統責任與邊界。

本文件中的部分 domain / bounded-context 邊界只是目前的 hypothesis。若未來從真實使用者行為、domain discovery 或實作證據中發現更高 cohesion 的切分方式，應允許調整，而不是把早期文件當成不可變的服務邊界。

---

## 產品使命

Kin 希望幫助真實世界中的朋友維持關係，透過低摩擦、經授權、具隱私邊界的方式，保留彼此最近生活、興趣與思考方向中的重要 context。

核心觀察很簡單：成年人通常不是不在乎朋友，而是逐漸失去彼此的日常 context。當 context 消失到一定程度，即使還想聯絡，也會因為「不知道對方最近在幹嘛」而增加重新開啟對話的摩擦。

Kin 想維持這種 **context continuity**，但不要求使用者持續發文、整理近況、經營 feed，或記得主動分享。

概念上：

`Authorized Activity Signals → Derived Context Candidates → Permissioned Friendship Context → Better Human Conversation`

Kin 的目標不是讓使用者把更多注意力留在產品裡，而是降低理解朋友近況的成本，增加自然對話發生的機會。

> We don't want your attention. We want to give it back to your friends.

---

## Kin 是什麼

Kin 是一個 **ambient friendship-context system**。

它可以接收使用者明確授權的 digital activity signals，從中辨識較高階、具社交意義的 context，再依照 relationship-aware privacy rules，決定哪些內容值得、而且允許被特定朋友知道。

有價值的 context 可能包括：

- 某位朋友最近開始持續關注一個新主題。
- 某位朋友反覆研究同一個問題，代表這件事近期可能很重要。
- 兩位朋友正在獨立探索相近主題，形成可自然開啟對話的共同交集。
- 某位朋友的興趣出現明顯變化，而不是只有單次行為。
- 一段原本活躍的關係安靜了一段時間，而現在剛好出現值得重新聯絡的 context。

系統應優先分享 **derived meaning**，而不是 raw activity log。

例如：

- Raw signal：最近按讚了五支 AI agents 影片。
- Derived context：最近持續在研究 AI agents 如何用於產品工作流程。

---

## Kin 不是什麼

Kin 不應被設計成：

- 高流量的公開 social network。
- 以 engagement、scrolling time 或 impression 最大化為核心的 feed。
- 把每一個 digital action 直接重新發布的系統。
- surveillance product。
- 取代真實人際互動的 AI friend。
- generic life-logging database。
- 需要使用者一直記得發文，產品才有價值的工具。

Kin 的存在目的是幫助真實朋友更容易理解彼此與重新開始對話，而不是增加另一個需要被經營的數位人格。

---

## 核心產品原則

### 1. Context continuity 優先於 content production

Kin 應幫助朋友維持對彼此的基本認知，不要求持續手動產生 social content。

### 2. Derived context 優先於 raw activity

Raw behavior 是 private input。可分享的 output 應通常是更高階的語意，而不是原始事件。

建議流程：

`Raw Behavior → Interest Detection → Significant Change → Friendship Context`

而不是：

`Raw Behavior → Social Feed`

### 3. Permissioned by default

使用者必須能控制某類 signal 是否可參與 social-context derivation，以及最後能分享給哪些 relationship。

### 4. Sensitive by default

當是否適合分享存在不確定性時，系統應優先隱藏、降低細節或不分享，而不是過度揭露。

### 5. Relationship-aware disclosure

同一份 underlying context 對不同 relationship 可以有不同 detail level。

例如：

- Friend：最近在研究 AI 產品。
- Close Friend：最近在研究給非技術創業者使用的 AI 產品。
- Inner Circle：最近在思考一個專門協助維持朋友關係的 AI 產品。

### 6. Explainability

當 Kin 決定顯示一則 context 時，應能在合理抽象層級說明：

- 為什麼這件事被判斷為有意義。
- 為什麼現在值得顯示。
- 為什麼這位朋友有權看到這個 detail level。

### 7. Human conversation 才是 outcome

Kin 應優先優化：

- 是否產生有意義的 conversation。
- 是否維持 friendship continuity。
- 是否幫助長時間未互動的朋友重新連結。

而不是單純增加 app 使用時間。

---

## 成功衡量方向

主要 outcome metrics 應反映 relationship value，而不是 attention capture。

候選指標：

- Conversations Started
- Friendships Maintained
- Dormant Friendships Reactivated
- 使用者是否更了解 close friends 最近的生活與興趣
- 有多少 context 真正帶來 follow-up conversation
- Shared-interest discovery 是否促成互動

DAU、session time、feed impressions 可以作為 operation metrics，但不應成為產品核心 success definition。

---

## 主要 Actors

### User

擁有 Kin identity 的使用者，可以：

- 授權或主動提供 activity signals。
- 管理 friendship 與 privacy settings。
- 接收朋友的 permissioned context。
- 從 context 開啟 conversation。

### Friend

與 User 建立有效 relationship 的另一位 Kin user。

### Relationship Participant

從「某一段 friendship」角度來看的使用者。

同一個 User 在不同 relationship 中，可以擁有不同 closeness、privacy policy 與 context projection。

### External Activity Provider

能提供經授權 signals 的外部系統，例如 music、video、reading、browser、AI assistant 或其他 digital services。

Provider 是外部 actor，不是 Kin domain model 的一部分，也不能反過來定義 Kin 的 domain architecture。

### AI Context Generator

用於 classify、summarize、normalize 或 derive context 的外部能力。

AI 是 adapter capability，不是 domain truth 的來源。

---

## High-level Capability Map

### Identity 與 relationship setup

- 建立 Kin identity。
- 發現或邀請另一位使用者。
- 建立 friendship。
- 表達 relationship closeness / trust。
- 管理 relationship lifecycle。

### Activity contribution

- 接收明確授權的 digital activity。
- 當 integration 不可用時，接受低摩擦的 manual contribution。
- 將不同 provider 的活動正規化成 Kin 可理解的概念。
- 區分 raw signal、normalized activity 與 socially meaningful context。

### Social-context derivation

- 辨識 recurring interests。
- 偵測 Significant Change。
- 從多個 signals 推導 Context Candidate。
- 抑制低訊號、重複或過度細碎的資訊。
- 保留適當的 provenance / explanation metadata。

### Privacy & Sharing

- 判斷 Context Candidate 是否可被分享。
- 根據 relationship 投影不同 detail level。
- 套用 user preferences 與 sensitive-data policy。
- 支援 opt-out、revocation 與 suppression。

### Relevance & Prioritization

- 判斷一則已允許分享的 context 是否值得現在顯示。
- 針對特定 relationship ranking。
- 抑制 stale / repetitive context。
- 偵測 serendipitous overlap。

### Friend Pulse

- 濃縮「目前最值得知道的朋友近況」。
- 每次只顯示少量高訊號 context。
- 避免退化成高流量 chronological feed。

### Conversation Support

- 產生 context-aware conversation starter。
- 幫助使用者問出自然且相關的問題。
- 解釋為什麼某個主題值得聊。
- 在合適範圍內追蹤 context 是否帶來 interaction。

### Integration

- 連接 external providers。
- 維護 provider authorization state。
- ingest / pull provider activity。
- 隔離 provider-specific payload、error 與 retry concerns。

### Notification

- 在有 meaningful context 時通知使用者。
- 避免 excessive interruption。
- 支援 digest、low-frequency delivery 等模式。

### Social Memory

長期而言，Kin 可以協助保留 shared history、recurring interests 與過往重要 context，但不能因此變成對私人對話進行無限制搜尋或監控的系統。

---

## 候選 Domains / Bounded Contexts

以下都是 **候選 domain boundaries**，不是最終 service boundaries。

### Identity

**責任**

代表一位 Kin 使用者作為產品參與者的身份與 account-level 狀態。

**Owns**

- stable user identity
- Kin 必要的 profile-level information
- account lifecycle state
- account / profile-level preferences，例如顯示名稱、語言或 locale 等；不包含 Privacy & Sharing 或 Notification 各自擁有的 policy / delivery preferences

**Does not own**

- friendship state
- provider-specific OAuth internals
- relationship-specific privacy
- sharing / disclosure preferences
- notification cadence / channel preferences
- derived social context

**Classification hypothesis**

Supporting Domain。

---

### Friendship

**責任**

代表兩位 Kin users 之間明確存在的 relationship，以及 friendship-aware product behavior 所需的 relationship state。

**Owns**

- relationship creation / acceptance
- active / inactive lifecycle
- relationship participants
- closeness / trust classification
- 非 privacy policy 本身的 relationship-specific state

**Does not own**

- user identity
- raw activity
- social-context generation
- notification delivery

**重要 invariants**

- Friendship 必須引用有效 participants。
- Relationship transition 必須是明確行為。
- Relationship-specific state 不可無意間變成 user-global state。

**Classification hypothesis**

Core Domain。

---

### Activity

**責任**

代表使用者明確授權 Kin 使用的 signals，例如看過、收藏、按讚、聽過、研究過、閱讀過或主動提供的活動。

**Owns**

- normalized activity concepts
- activity provenance
- 與該 activity 有關的 contribution / authorization intent
- timestamp 與 semantic metadata
- raw / normalized input 與 derived context 的區分

**Does not own**

- provider connection lifecycle
- friend 是否能看到這項 activity
- 最終 social-context wording
- friend-specific relevance ranking

**重要原則**

Activity 不等於 shareable social item。

**Classification hypothesis**

Core 或 Supporting Domain。需等真實 ingestion use cases 出現後再重新判斷。

---

### Social Context

**責任**

把一個或多個 authorized signals 轉換為 socially meaningful、human-readable 的 Context Candidate。

**Owns**

- Context Candidate
- topic / semantic meaning
- significance / change interpretation
- context lifecycle，例如 candidate、approved、suppressed、expired
- 與 source signals 的 abstract provenance

**Does not own**

- provider API payload
- friendship lifecycle
- final friend-specific privacy projection
- delivery channel

**重要 invariants**

- Raw provider output 不能直接成為 domain Context。
- Provider / LLM output 必須先由對應 adapter normalization + validation，轉成 inner contract 後，才可回到 application / domain flow。
- 一個 raw action 不應自動等於一個 context item。

**Classification hypothesis**

Core Domain。

---

### Privacy & Sharing

**責任**

判斷某個 Context 對某一段 specific relationship 是否可揭露，以及可揭露到什麼 detail level。

**Owns**

- disclosure policy
- relationship-aware visibility rules
- user sharing preferences
- sensitive-content policy
- projection / detail level
- revocation / suppression decision

**Does not own**

- friendship lifecycle
- semantic context generation
- 已允許 context 的 relevance ranking
- provider authorization

**重要 invariants**

- Default 採 least-revealing valid representation。
- 沒有 permission 不等於有 permission。
- Derived context 可分享，不代表 raw activity 也能分享。
- Disclosure decision 必須是 relationship-specific。
- 已排程但尚未發送的 disclosure intent，不得把過去的 approval 視為永久有效；若 owner 在 delivery 前撤銷或降低分享權限，pending intent 必須被取消或在 dispatch 前重新驗證。

**Classification hypothesis**

Core Domain。

---

### Relevance

**責任**

從「已經允許該 relationship 看到」的 context 中，判斷什麼內容現在值得 surfaced。

**Owns**

- ranking signals
- novelty / staleness decision
- repetition suppression
- shared-interest opportunity detection
- Friend Pulse / digest prioritization

**Does not own**

- context 是否允許被分享
- semantic context generation
- notification transport

**重要原則**

Relevance 不能繞過 Privacy & Sharing。

**Classification hypothesis**

可能是 Core 或 Supporting Domain，取決於未來 Kin 的 ranking / timing intelligence 是否形成明顯 differentiation。

---

### Conversation / Interaction

**責任**

讓 friendship context 能自然轉化成真實互動。

**Owns**

- conversation-starting intent
- context-based prompt / question
- Kin 明確追蹤的 interaction outcome
- 「聊聊這個」等 product-level interaction state

**Does not own**

- 第三方 private chat history
- generic messaging infrastructure，除非 Kin 未來真的成為 messaging product
- friendship permission
- context derivation

**Classification hypothesis**

屬於核心 product experience，但 MVP 初期可能仍是 Supporting Domain。

---

### Integration

**責任**

管理與 external systems 的連接與資料交換邊界。

**Owns**

- provider connection state
- external authorization lifecycle
- sync / checkpoint state
- provider-specific normalization boundary
- integration error translation / retry concerns

**Does not own**

- normalization 後的 Kin Activity domain behavior
- social-context meaning
- friendship state
- privacy disclosure policy

**Classification hypothesis**

Generic / Supporting Domain。

---

### Notification

**責任**

把已被允許、已被判斷值得 surfaced 的 product intent，在適當時間與 channel 送達使用者。

**Owns**

- notification intent
- delivery timing / preferences
- channel selection
- delivery status

**Does not own**

- context 是否允許分享
- context 是否有 semantic meaning
- friendship state

**重要原則**

Notification 不自行重新解讀 privacy；若 notification 可能延遲發送，application orchestration 必須在 dispatch 前透過 Privacy & Sharing 重新確認 disclosure 仍有效，或在 revocation 發生時取消 pending intent。

**Classification hypothesis**

Generic / Supporting Domain。

---

## Domain Interactions

以下描述的是 conceptual dependency，不代表 synchronous service call，也不代表未來必須拆成 microservices。

### Activity contribution → Social Context

`Integration → Activity → Social Context`

External provider 先透過 Integration boundary 提供 signal，再 normalize 成 Kin Activity，Social Context 才能針對 Kin concepts 判斷意義。

Provider-specific DTO 不得穿透到 Activity / Social Context domain model。

### Context disclosure

`Social Context + Friendship + Privacy & Sharing → Relationship-specific Context Projection`

Context 存在，不代表 friend 就能看到。必須把 relevant Friendship 與 Privacy Policy 一起納入判斷，決定是否可揭露與 detail level。

### Friend Pulse

`Context Projection + Relevance → Friend Pulse`

只有已經通過 privacy projection 的 Context Projection 才能進 ranking。

Relevance 不可 bypass privacy。

### Conversation support

`Friend Pulse / Context Projection → Conversation / Interaction`

Conversation support 只能使用使用者原本就有權看到的 Context Projection，不可透過 prompt 重新暴露更敏感的內容。

### Notification

`Context Projection + Relevance → Notification Intent → Dispatch-time Privacy Revalidation → Notification`

Notification 只負責 delivery，不可自行重新解讀 privacy 或創造新的 social meaning。若 delivery 是延遲的，application orchestration 必須在 dispatch 前重新確認目前的 Privacy & Sharing 決策仍允許該 disclosure；若已被 revocation，pending intent 必須取消或失效。

---

## Initial Ubiquitous Language

### Activity

使用者明確授權的 signal，描述某件他做過、看過、收藏、按讚、聽過、研究過或主動提供的事情。

### Raw Activity

仍保有 provider / source shape、尚未進行 Kin normalization 與 social interpretation 的 activity。

### Normalized Activity

Provider-independent 的 Kin activity representation。

### Context Candidate

由一個或多個 authorized signals 推導出的「可能具有社交意義」資訊。它尚不代表可分享，也不代表值得 surfaced。

### Social Context

已被 validation、具有可理解 social meaning，且可能值得朋友知道的描述。

### Context Projection

某個 Context 經過 relationship-specific privacy / sharing policy 後，得到的可見版本。

### Friendship

兩位 Kin users 之間明確存在的 relationship。

### Relationship Level

表示 closeness / trust 的 relationship classification，可能影響 disclosure policy。

長期候選值：

- Acquaintance
- Friend
- Close Friend
- Inner Circle

MVP 不需要一次實作完整 hierarchy，可以先只驗證 close-friend / inner-circle 類型的關係。

### Friend Pulse

針對某位朋友，少量、經 prioritization、permissioned 的 meaningful context 集合或摘要。

它不是 chronological activity feed。

### Significant Change

相較於單一 raw action，更具資訊價值的 recurring pattern 或明顯興趣變化。

### Shared Rabbit Hole

兩位朋友獨立探索相同或高度相關主題，而這個 overlap 可能形成自然 conversation opportunity。

### Friendship Drift

一段 relationship 長期變得較少互動，或 context continuity 明顯下降的狀態。

這是未來 capability，不代表 MVP 必須實作。

### Social Memory

協助維持 shared history / relationship continuity 的長期 context，同時避免不必要的 raw record disclosure。

---

## Core / Supporting / Generic Domain Hypotheses

以下 classification 都是 provisional。

### Likely Core Domains

- Friendship
- Social Context
- Privacy & Sharing
- 可能是 Relevance
- 可能是 Activity，取決於未來 signal interpretation 是否形成足夠 domain complexity

這些區域最直接體現 Kin 的差異化：

**在真實 friendship 中維持 context，同時精準控制每個 relationship 到底能知道什麼。**

### Likely Supporting Domains

- Identity
- MVP 初期的 Conversation / Interaction
- 若 Activity 最終主要只是 normalization / lifecycle，則 Activity 也可能屬於 Supporting

### Likely Generic Domains

- Integration plumbing
- Notification delivery
- provider connectivity infrastructure

技術上困難，不代表一定是 Core Domain。

---

## Privacy Model 方向

長期 privacy model 應以 relationship-aware 為核心，而不是單一 global public/private switch。

候選 Relationship Levels：

- Acquaintance
- Friend
- Close Friend
- Inner Circle

Taxonomy 尚未定案。

真正重要的 domain capability 是：

**同一個 source context，能對不同 relationship 安全地產生不同 detail level 的 projection。**

MVP 不應在尚未驗證 close-friend context 本身是否有價值前，就先實作完整 relationship hierarchy。

---

## AI / LLM Boundary

AI 是 external capability，不是 domain authority。

概念流程：

`Authorized Signals → Application Use Case → ContextGenerator Port → AI Adapter [provider call + normalization + validation + error translation] → Validated Context Draft → Application / Domain → Context Candidate`

核心規則：

- Domain 不知道 OpenAI、Anthropic、Gemini 或 model names。
- Domain 不知道 prompt template。
- Raw LLM JSON 不可直接進入 application 或 domain object。
- AI Adapter 必須在 adapter boundary 內完成 provider output normalization / validation 與 provider-specific failure translation，再透過 `ContextGenerator` port 回傳 inner contract。
- Port 名稱應描述 business capability，例如 `ContextGenerator`，而不是 vendor API。

Domain truth 最終仍由 Kin 自己的 invariants、policy 與 validated domain objects 決定。

---

## 長期產品能力，不代表 MVP 授權

以下能力合理地屬於 Kin 的 long-term product map，但它們出現在本文件中 **不代表目前可以實作**：

- 完整 Relationship Level hierarchy
- Friendship Drift Detection
- Weekly Friendship Digest
- Shared Rabbit Hole
- Social Memory
- Friend-aware AI Q&A，例如「Jerry 最近在幹嘛？」
- AI Friendship Concierge
- 多 provider 自動 ingestion
- 更進階的 relevance / timing intelligence
- richer conversation outcome tracking
- widgets / app intents / ambient surfaces

實際 implementation authorization 應由：

1. current GitHub Issue
2. current MVP slice
3. repository architecture contracts

共同決定。

Product scope 提供 context，不提供自動實作權限。

---

## 高階不變原則

即使未來 bounded contexts 調整，以下原則應保持穩定：

1. Raw activity private-by-default。
2. Activity 不等於 Social Context。
3. Social Context 不等於 friend-visible Context Projection。
4. Privacy decision 必須先於 Relevance ranking。
5. Relevance 不可繞過 privacy。
6. 延遲 delivery 在 dispatch 前必須確認 disclosure 仍有效；revocation 必須能使 pending intent 失效。
7. AI / provider 永遠是 outer adapter，而不是 domain authority。
8. 一個 future capability 出現在 product map，不代表已獲得 implementation authorization。
9. Infrastructure convenience 不得主導 domain boundary。
10. Kin 最終服務的是 human relationship，不是 engagement metric。

---

## 尚待 Domain Discovery 驗證的問題

以下問題刻意保留，不在 foundation 階段過早定案：

### Activity 與 Integration 的邊界

- provider normalization 應在哪個 boundary 結束？
- 哪些 semantic normalization 屬於 Activity，哪些只是 Integration concern？

### Activity 是否為 Core Domain

- 真實 use cases 是否會出現複雜的 activity lifecycle / invariants？
- 還是 Activity 最終只是 Social Context 上游的 supporting model？

### Social Context 與 Relevance

- Context significance 與 friend-specific relevance 是否應維持分離？
- 是否會因未來使用資料而形成明確獨立 bounded context？

### Friendship 與 Privacy & Sharing

- Relationship Level 是 Friendship 的 state，還是 Privacy Policy 的 input？
- 個別 friendship override 應由哪個 domain 擁有？

### Conversation / Interaction

- Kin 是否只協助「開始 conversation」，還是未來會擁有更多 interaction lifecycle？
- 若真正 conversation 發生在外部 app，Kin 應追蹤到什麼程度？

### Social Memory

- 哪些 long-lived context 值得保存？
- 如何避免 Social Memory 變成過度監控或不必要的 raw-history retention？

這些問題應由後續 user stories、MVP slices 與真實 domain behavior 逐步回答。

---

## 文件與實作的關係

本文件描述的是 **long-term conceptual product system**。

它不定義：

- database schema
- REST / GraphQL endpoint
- persistence technology
- queue / event infrastructure
- provider SDK
- deployment topology
- microservice boundaries

下一層應由 `docs/product/mvp-roadmap.md` 把產品切成可驗證的 vertical slices，再由 active GitHub Issue 決定目前真正被授權的 implementation scope。

建議閱讀順序：

`Product Scope → MVP Slice → User Story → Use Case → Domain Responsibility → Implementation`
