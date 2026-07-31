# Release Process

LayerGrid CLI releases are built with GoReleaser and signed with Sigstore Cosign keyless signing.

## Local Validation

```sh
make validate-release
```

This runs:

- `goreleaser check`
- cross-platform snapshot builds
- archive generation
- SBOM generation
- checksum generation
- Homebrew formula generation
- archive extraction
- `layergrid version` from the extracted archive

## Cosign Verification

For a tagged release, verify an artifact with:

```sh
cosign verify-blob \
  --certificate layergrid_VERSION_darwin_arm64.tar.gz.pem \
  --signature layergrid_VERSION_darwin_arm64.tar.gz.sig \
  --certificate-identity-regexp '^https://github.com/layergrid/layergrid-cli/.github/workflows/release.yaml@refs/tags/VERSION_TAG$' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  layergrid_VERSION_darwin_arm64.tar.gz
```

The installer verifies `checksums.txt` with the same issuer and repository identity before installing a binary.

## Homebrew

Design partner install command after a public release is published:

```sh
brew install layergrid/tap/layergrid
```

Cross-repo publishing uses `HOMEBREW_TAP_TOKEN`, scoped to `contents:write` on `layergrid/homebrew-tap`.

The release workflow writes `Formula/layergrid.rb` to `layergrid/homebrew-tap` from the release checksums and public GitHub release asset URLs.

## Curl

Temporary GitHub Pages install command after Pages is enabled:

```sh
curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
```

Fallbacks:

```sh
brew install layergrid/tap/layergrid
go install github.com/layergrid/layergrid-cli/cmd/layergrid@latest
```

## SARIF Validation

Expected validation flow:

```sh
layergrid scan . --format sarif --output layergrid.sarif
gh api repos/layergrid/layergrid-cli/code-scanning/sarifs --method POST --input payload.json
```

## Public RC Validation

- `v0.1.0-rc.1` release workflow completed successfully.
- Release archives, SBOMs, and `checksums.txt` all have Cosign `.sig` and `.pem` assets.
- `cosign verify-blob` passed for `checksums.txt`, Darwin arm64, and Linux amd64 artifacts.
- `brew install --formula layergrid/tap/layergrid` installs `0.1.0-rc.1`.
- `curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash` installs `0.1.0-rc.1` on macOS and Ubuntu.
- LayerGrid SARIF self-scan was accepted by GitHub code scanning for `layergrid/layergrid-cli`.
