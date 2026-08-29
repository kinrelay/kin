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
mkdir -p apps/api docs
echo 'package api' > apps/api/foo.go
git add .
git commit -qm 'base'
base="$(git rev-parse HEAD)"

git checkout -qb feature
git mv apps/api/foo.go docs/foo.go
git commit -qm 'move backend file out'
head="$(git rev-parse HEAD)"

result="$(bash "$detector" pr "$base" "$head")"
[[ "$result" == "true" ]] || { echo "expected backend rename to be relevant" >&2; exit 1; }

git checkout -q -b main "$base"
echo 'package api' > apps/api/base-only.go
git add apps/api/base-only.go
git commit -qm 'base branch backend change'
main_head="$(git rev-parse HEAD)"

git checkout -q feature
echo 'docs only' > docs/readme.md
git add docs/readme.md
git commit -qm 'feature docs change'
feature_head="$(git rev-parse HEAD)"

result="$(bash "$detector" pr "$main_head" "$feature_head")"
[[ "$result" == "false" ]] || { echo "expected unrelated base-branch backend change to be ignored" >&2; exit 1; }

echo 'CI diff detector regression tests passed.'
