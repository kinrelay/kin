#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
detector="$repo_root/scripts/ci/detect-backend-changes.sh"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

cd "$tmp"
git init -q
git config user.email ci-test@example.com
git config user.name ci-test
git config commit.gpgsign false
git config core.hooksPath "$tmp/.empty-hooks"
mkdir -p "$tmp/.empty-hooks" apps/api docs
echo 'package api' > apps/api/foo.go
git add .
git commit -qm 'base'
base="$(git rev-parse HEAD)"

# Case 1: moving a backend file outside apps/api must still be relevant.
git checkout -qb rename-feature "$base"
git mv apps/api/foo.go docs/foo.go
git commit -qm 'move backend file out'
rename_head="$(git rev-parse HEAD)"
result="$(bash "$detector" pr "$base" "$rename_head")"
[[ "$result" == "true" ]] || { echo "expected backend rename to be relevant" >&2; exit 1; }

# Case 2: backend-only changes on the base branch after the split must not make
# a docs-only PR backend-relevant.
git checkout -q -b docs-feature "$base"
echo 'docs only' > docs/readme.md
git add docs/readme.md
git commit -qm 'feature docs change'
docs_head="$(git rev-parse HEAD)"

git checkout -q -b base-updated "$base"
echo 'package api' > apps/api/base-only.go
git add apps/api/base-only.go
git commit -qm 'base branch backend change'
base_updated_head="$(git rev-parse HEAD)"
result="$(bash "$detector" pr "$base_updated_head" "$docs_head")"
[[ "$result" == "false" ]] || { echo "expected unrelated base-branch backend change to be ignored" >&2; exit 1; }

# Case 3: non-ASCII backend paths must remain backend-relevant.
git checkout -q -b utf8-feature "$base"
echo 'package api' > 'apps/api/測試.go'
git add 'apps/api/測試.go'
git commit -qm 'add utf8 backend file'
utf8_head="$(git rev-parse HEAD)"
result="$(bash "$detector" pr "$base" "$utf8_head")"
[[ "$result" == "true" ]] || { echo "expected UTF-8 backend path to be relevant" >&2; exit 1; }

# Case 4: quoted/special-character backend paths must not be hidden by Git path quoting.
git checkout -q -b quoted-feature "$base"
printf 'package api\n' > 'apps/api/a"b.go'
git add 'apps/api/a"b.go'
git commit -qm 'add quoted backend file'
quoted_head="$(git rev-parse HEAD)"
result="$(bash "$detector" pr "$base" "$quoted_head")"
[[ "$result" == "true" ]] || { echo "expected quoted backend path to be relevant" >&2; exit 1; }

# Case 5: an unavailable revision must fail closed, not report false.
set +e
bash "$detector" push "$base" deadbeefdeadbeefdeadbeefdeadbeefdeadbeef >/dev/null 2>&1
status=$?
set -e
[[ "$status" -ne 0 ]] || { echo "expected unavailable SHA to fail detection" >&2; exit 1; }

# Case 6: branch-creation pushes use an all-zero before SHA and must compare
# against the empty tree rather than fail before test selection.
zero_sha="0000000000000000000000000000000000000000"
result="$(bash "$detector" push "$zero_sha" "$base")"
[[ "$result" == "true" ]] || { echo "expected branch-creation push to inspect full tree" >&2; exit 1; }

echo 'CI diff detector regression tests passed.'
