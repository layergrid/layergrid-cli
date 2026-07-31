#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
scratch="${LAYERGRID_CORPUS_DIR:-/tmp/layergrid-corpus}"
results="$root/.corpus-results"
mkdir -p "$scratch" "$results"

export GOCACHE="${GOCACHE:-$root/.gocache}"

if [[ ! -x "$root/bin/layergrid" ]]; then
  (cd "$root" && go build -o bin/layergrid ./cmd/layergrid)
fi

while read -r repo; do
  [[ -z "$repo" || "$repo" =~ ^# ]] && continue
  name="${repo##*/}"
  dest="$scratch/$name"
  if [[ -d "$dest/.git" ]]; then
    git -C "$dest" fetch --depth=1 origin >/dev/null 2>&1 || true
  else
    git clone --depth=1 "https://github.com/$repo.git" "$dest"
  fi
  out="$results/${name}.json"
  "$root/bin/layergrid" scan "$dest" --format json --no-color > "$out"
  echo "$repo -> $out"
done < "$root/hack/corpus-repos.txt"
