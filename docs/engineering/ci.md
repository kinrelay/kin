# CI Foundation

Kin 目前以 GitHub Actions 作為 primary CI execution layer，但 testing semantics 必須保持 provider-neutral。

## Canonical test command

Repository-level canonical command：

```bash
make test
```

目前它會執行 `apps/api` 的 Go tests。GitHub Actions、未來可能導入的 CircleCI，或其他 runner 都應呼叫這個 command，而不是各自維護另一套 test semantics。

## GitHub Actions behavior

目前 CI 在以下情況執行：

- Pull Request
- push 到 `main`

同一 PR / branch 的舊 workflow 會透過 GitHub Actions concurrency 自動取消，降低頻繁 push / AI review loop 的 runner 消耗。

CI workflow 本身維持一個穩定的 `test` check。它會先檢查此次 event 的實際 diff：

- 若變更包含 `apps/api/**`、root `Makefile` 或 CI workflow 本身，執行 `make test`。
- 若 diff 沒有 backend executable/config change，例如純 docs change，保留成功的 CI check，但不執行 Go setup / Go tests。

Pull Request 使用 PR base SHA 到目前 head SHA 的完整 diff；push 到 `main` 使用 event 的 before SHA 到目前 SHA。這避免只看最後一個 commit 而漏掉同一 PR / push 中較早的 executable change。

## TDD relationship

CI 驗證的是 committed state 是否為 Green。Feature development 的 Red → Green → Refactor 過程仍由 TDD / workflow contract 負責，並應在 PR evidence 中誠實記錄。

## Future CircleCI portability

如果未來因 GitHub Actions quota / cost 需要導入 CircleCI overflow，CircleCI 應直接重用 `make test`。不要把 business/test semantics 複製進 CircleCI config。

Quota threshold、`CI_PROVIDER` switching、billing automation 或獨立 CI Router 不屬於目前 Phase 1 foundation。
