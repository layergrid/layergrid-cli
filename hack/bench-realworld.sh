#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
scratch="${LAYERGRID_CORPUS_DIR:-/tmp/layergrid-corpus}"
results="$root/.bench-results.tsv"
mkdir -p "$scratch"

export GOCACHE="${GOCACHE:-$root/.gocache}"

if [[ ! -x "$root/bin/layergrid" ]]; then
  (cd "$root" && go build -o bin/layergrid ./cmd/layergrid)
fi

printf "repo\tloc\tcold_seconds\twarm_seconds\tpeak_rss_bytes\tfindings\n" > "$results"

while read -r repo; do
  [[ -z "$repo" || "$repo" =~ ^# ]] && continue
  name="${repo##*/}"
  dest="$scratch/$name"
  if [[ ! -d "$dest/.git" ]]; then
    git clone --depth=1 "https://github.com/$repo.git" "$dest"
  fi
  loc="$(find "$dest" -type f \( -name '*.py' -o -name '*.json' -o -name '*.yaml' -o -name '*.yml' \) \
    -not -path '*/.git/*' -not -path '*/node_modules/*' -not -path '*/.venv/*' \
    -print0 | xargs -0 wc -l | awk 'END {print $1+0}')"
  cold_stats="$(/usr/bin/time -l sh -c "\"$root/bin/layergrid\" scan \"$dest\" --format json --no-color >/tmp/layergrid-bench-cold.json" 2>&1)"
  warm_stats="$(/usr/bin/time -l sh -c "\"$root/bin/layergrid\" scan \"$dest\" --format json --no-color >/tmp/layergrid-bench-warm.json" 2>&1)"
  cold_seconds="$(printf "%s\n" "$cold_stats" | awk '/real/ {print $1; exit}')"
  warm_seconds="$(printf "%s\n" "$warm_stats" | awk '/real/ {print $1; exit}')"
  peak_rss="$(printf "%s\n" "$warm_stats" | awk '/maximum resident set size/ {print $1; exit}')"
  findings="$(grep -o '"ruleId"' /tmp/layergrid-bench-warm.json 2>/dev/null | wc -l | tr -d ' ' || true)"
  printf "%s\t%s\t%s\t%s\t%s\t%s\n" "$repo" "$loc" "${cold_seconds:-n/a}" "${warm_seconds:-n/a}" "${peak_rss:-n/a}" "$findings" | tee -a "$results"
done < "$root/hack/corpus-repos.txt"

echo "wrote $results"
