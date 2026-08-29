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
    diff_base="$base_sha"
    ;;
  *)
    echo "unsupported mode: $mode" >&2
    exit 2
    ;;
esac

relevant='^(apps/api/|Makefile$|\.github/workflows/ci\.yml$|scripts/ci/)'

# Capture the diff first so an unavailable revision fails the script instead of
# being mistaken for an empty/non-backend diff. Disable path quoting so UTF-8
# filenames retain their real prefixes for classification.
diff_output="$(git -c core.quotePath=false diff --name-status -M "$diff_base" "$head_sha")"

while IFS=$'\t' read -r status path1 path2; do
  [[ -z "${status:-}" ]] && continue

  if [[ "${path1:-}" =~ $relevant ]] || [[ "${path2:-}" =~ $relevant ]]; then
    echo "true"
    exit 0
  fi
done <<< "$diff_output"

echo "false"
