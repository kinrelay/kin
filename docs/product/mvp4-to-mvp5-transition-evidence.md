# MVP 4 → MVP 5 Transition Evidence

## Purpose

這份文件補充 `docs/product/mvp-roadmap.md` 的 MVP 4 → MVP 5 transition evidence，保留這次 closure 的完整可追溯歷史。

## Why a fresh transition was required

先前 PR #78 曾嘗試把 Active Slice 從 MVP 4 切到 MVP 5，但 review / Issue #74 的 reconciliation 發現三個仍阻擋 MVP 4 completion 的實作 gap，因此 #78 被關閉且未合併：

1. 缺少 non-test runnable Friend Pulse delivery wiring。
2. production path 尚未客觀保證 expired / repetition-suppressed context 不進入 Pulse。
3. privacy-first ordering 尚未成立；Relevance 仍可能在 projection 前影響排序。

因此不得把 #78 的舊 review 或舊 transition 判斷當作 closure evidence。

## Closure evidence

Issue #79 / PR #80 專門補齊上述 closure gaps：

- Friend Pulse 先對候選做 relationship-specific privacy projection，再對已可見 projection ranking / bounding，避免 hidden context 影響結果。
- production eligibility 明確排除 expired 與 repetition-suppressed candidates，並以 injectable clock 提供 deterministic boundary evidence。
- composition / delivery path 有非單純 domain-unit-test 的 executable wiring / integration coverage。
- PR #80 exact-current-head CI 與 review gate 通過後 squash merge，merge commit：`916d40146bef478839bbca6264c82fa4b4a806c5`。
- merge 後 `main` GitHub Actions CI #537 成功，提供 post-merge repository evidence。

## MVP 4 reconciliation

依 current roadmap 的六項 Acceptance Criteria：

1. active friend 可取得 Friend Pulse：已由 merged Friend Pulse query / composition flow 覆蓋。
2. Pulse 只包含 viewer 當下可見的 Context Projection：privacy-first projection 已在 production use case enforcement。
3. Pulse 維持 bounded 1–3 high-signal items：既有 `maxPulseItems` behavior 與 tests 保留。
4. revoked / expired / suppressed context 不出現在 Pulse：revocation 由既有 projection contract保護；#79 補齊 expiry / repetition suppression production eligibility。
5. Relevance 不讀取或重新揭露 projection 前敏感內容：#79 將 ranking 移到 privacy projection 之後。
6. 不依賴 chronological feed：Friend Pulse 仍是 bounded purpose-built query / read model。

Slice Completion Signal「使用者可以在極短時間內理解朋友最近最值得知道什麼」目前由 bounded 1–3 item Friend Pulse interaction 作 MVP 可執行 proxy；真正的 conversation-value 驗證留給下一個 Active Slice MVP 5，不把尚未觀測的 conversation outcome偽造成 MVP 4 evidence。

## Transition decision

MVP 4 的 implementation / safety / delivery blockers 已由 #79 / PR #80 補齊，且 post-merge main CI 綠；因此允許把 **MVP 5 — Context 幫助開始真實 Conversation** 設為唯一 Active Slice。

這個 transition 不授權 MVP 6 integration，也不授權 future-scope capability；後續 implementation issue 必須只從 MVP 5 的 Acceptance Criteria / Slice Completion Signal 選擇最小 coherent vertical slice。
