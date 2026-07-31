#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export GOCACHE="${GOCACHE:-$root/.gocache}"
export GOMODCACHE="${GOMODCACHE:-$root/.gomodcache}"

goreleaser check
goreleaser release --snapshot --clean --skip=publish,sign

latest_archive="$(find "$root/dist" -name 'layergrid_*_darwin_arm64.tar.gz' | sort | tail -n 1)"
if [[ -z "$latest_archive" ]]; then
  echo "missing darwin arm64 archive" >&2
  exit 1
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT
tar -xzf "$latest_archive" -C "$tmpdir"
"$tmpdir/layergrid" version
