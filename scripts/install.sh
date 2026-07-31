#!/usr/bin/env bash
set -euo pipefail

repo="layergrid/layergrid-cli"
version="${LAYERGRID_VERSION:-latest}"
bindir="${LAYERGRID_INSTALL_DIR:-/usr/local/bin}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

if [[ "$version" == "latest" ]]; then
  version="$(curl -fsSL "https://api.github.com/repos/${repo}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)"
fi
if [[ -z "$version" ]]; then
  echo "could not resolve LayerGrid version" >&2
  exit 1
fi

base="https://github.com/${repo}/releases/download/${version}"
archive="layergrid_${version#v}_${os}_${arch}.tar.gz"
if [[ "$os" == "windows" ]]; then
  archive="layergrid_${version#v}_${os}_${arch}.zip"
fi

curl -fsSLo "$tmpdir/$archive" "$base/$archive"
curl -fsSLo "$tmpdir/checksums.txt" "$base/checksums.txt"
curl -fsSLo "$tmpdir/checksums.txt.sig" "$base/checksums.txt.sig"
curl -fsSLo "$tmpdir/checksums.txt.pem" "$base/checksums.txt.pem"

if ! command -v cosign >/dev/null 2>&1; then
  echo "cosign is required to verify LayerGrid releases" >&2
  exit 1
fi

cosign verify-blob \
  --certificate "$tmpdir/checksums.txt.pem" \
  --signature "$tmpdir/checksums.txt.sig" \
  --certificate-identity-regexp "https://github.com/layergrid/layergrid-cli/.github/workflows/release.yaml@refs/tags/.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$tmpdir/checksums.txt"

cd "$tmpdir"
grep "  $archive$" checksums.txt | shasum -a 256 -c -

case "$archive" in
  *.zip) unzip -q "$archive" ;;
  *.tar.gz) tar -xzf "$archive" ;;
esac

install -m 0755 layergrid "$bindir/layergrid"
echo "installed layergrid to $bindir/layergrid"
