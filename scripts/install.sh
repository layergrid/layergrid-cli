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

cosign_bin="$(command -v cosign || true)"
if [[ -z "$cosign_bin" ]]; then
  cosign_version="${COSIGN_VERSION:-v2.6.0}"
  cosign_name="cosign-${os}-${arch}"
  cosign_url="https://github.com/sigstore/cosign/releases/download/${cosign_version}/${cosign_name}"
  cosign_bin="$tmpdir/cosign"
  curl -fsSLo "$cosign_bin" "$cosign_url"
  chmod +x "$cosign_bin"
fi

"$cosign_bin" verify-blob \
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

if ! install -m 0755 layergrid "$bindir/layergrid" 2>/dev/null; then
  bindir="$HOME/.local/bin"
  mkdir -p "$bindir"
  install -m 0755 layergrid "$bindir/layergrid"
fi

installed_version="$("$bindir/layergrid" version | head -n 1)"
echo "LayerGrid installed: $installed_version"
echo "Next: $bindir/layergrid scan ."
