# CI Foundation

目前以 GitHub Actions 作為 primary CI execution layer，但 testing semantics 必須保持 provider-neutral。

## Canonical test command

Repository-level canonical command：

```bash
make test
```

目前它會先跑 CI diff detector regression tests，再執行 `apps/api` 的 Go tests。GitHub Actions、未來可能導入的 CircleCI，或其他 runner 都應呼叫這個 command，而不是各自維護另一套 test semantics。

## GitHub Actions behavior

目前 CI 在以下情況執行：

- Pull Request
- push 到 `main`

為避免同一 PR branch 同時跑 push CI 與 PR CI，目前 feature branch push 不另外觸發 workflow。

同一 PR / branch 有新 commit 時，舊的 in-progress workflow 會由 GitHub Actions concurrency 自動取消。

## Diff-based optimization

workflow 會先判斷此次 event 的完整 diff：

- Pull Request：使用 PR base/head 的 merge base → current head SHA
- push `main`：event before SHA → current SHA

changed-file detection 會保留 rename 的 source 與 destination path；因此 backend 檔案即使被搬出 `apps/api/**`，仍會被視為 backend-relevant change。

只有 diff 包含 backend executable/config concern 時才 setup Go 並執行 `make test`。目前包括：

- `apps/api/**`
- `Makefile`
- `.github/workflows/ci.yml`
- `scripts/ci/**`

純文件或與 backend 無關的變更仍會留下穩定的 `test` job/check，但不消耗 Go setup / test workload。

這個判斷以 event 的完整 diff 為準，不以 branch 歷史或單一最後 commit 判斷。

## TDD relationship

CI 驗證的是 committed state 是否為 Green。Feature-level Red → Green → Refactor evidence 由 development workflow / TDD contract 負責，不應把 intentionally failing Red commit 當成可合併狀態。

## Future provider portability

未來若加入 CircleCI，它應重用相同 repository-level commands（目前是 `make test`），而不是在 CircleCI config 裡重新定義另一套 testing semantics。

完整 quota router、billing polling、`CI_PROVIDER` 自動切換與 expensive-job routing 不屬於目前 foundation scope。
