#!/usr/bin/env bash
set -euo pipefail

mode="${1:?mode (pr|push) is required}"
base_sha="${2:?base SHA is required}"
head_sha="${3:?head SHA is required}"

case "$mode" in
  pr)
    diff_base="$(git merge-base "$base_sha" "$head_sha")"
    ;;
  push)
    if [[ "$base_sha" =~ ^0{40}$ ]]; then
      diff_base="$(git hash-object -t tree /dev/null)"
    else
      diff_base="$base_sha"
    fi
    ;;
  *)
    echo "unsupported mode: $mode" >&2
    exit 2
    ;;
esac

relevant='^(apps/api/|Makefile$|\.github/workflows/ci\.yml$|scripts/ci/)'

diff_file="$(mktemp)"
trap 'rm -f "$diff_file"' EXIT

# NUL delimiters preserve every legal Git pathname byte without C-style quoting.
# Rename detection is required because moving a relevant file out of its subtree
# changes that subtree even when the destination is non-relevant. Copy attribution
# is intentionally not used: copying an unchanged backend blob elsewhere does not
# change the backend tree, and Git's chosen copy source is heuristic.
# Writing to a temporary file also lets git diff failures propagate before parsing.
git diff --name-status -z -M "$diff_base" "$head_sha" > "$diff_file"

while IFS= read -r -d '' status; do
  IFS= read -r -d '' path1
  path2=""

  if [[ "$status" == R* ]]; then
    IFS= read -r -d '' path2
  fi

  if [[ "$path1" =~ $relevant ]] || [[ "$path2" =~ $relevant ]]; then
    echo "true"
    exit 0
  fi
done < "$diff_file"

echo "false"
