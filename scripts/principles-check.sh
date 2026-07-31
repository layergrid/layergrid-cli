#!/usr/bin/env bash
set -euo pipefail

if grep -R "http.Get\|http.Post\|panic(" internal/scan internal/detectors internal/trifecta 2>/dev/null; then
  echo "principles-check failed: scanner code contains network calls or panic"
  exit 1
fi

find rules -name '*.yaml' -type f | sort | while read -r rule; do
  grep -q '^id:' "$rule"
  grep -q '^fix:' "$rule"
  grep -q '^score_impact:' "$rule"
done
