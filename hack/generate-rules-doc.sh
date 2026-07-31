#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
out="$root/docs/RULES.md"

{
  echo "# LayerGrid Rules"
  echo
  echo "Generated from the embedded YAML rule library."
  echo
  echo "| Rule | Severity | Category | Description | References |"
  echo "|---|---|---|---|---|"
  find "$root/rules" -name '*.yaml' -type f | sort | while read -r rule; do
    id="$(awk -F': ' '/^id:/ {print $2; exit}' "$rule")"
    severity="$(awk -F': ' '/^severity:/ {print $2; exit}' "$rule")"
    category="$(awk -F': ' '/^category:/ {print $2; exit}' "$rule")"
    description="$(awk -F': ' '/^description:/ {print $2; exit}' "$rule")"
    ref="$(awk '/^  - / {print substr($0,5); exit}' "$rule")"
    printf '| `%s` | %s | %s | %s | %s |\n' "$id" "$severity" "$category" "$description" "$ref"
  done
} > "$out"

echo "wrote $out"
