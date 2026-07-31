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
- Homebrew cask generation
- archive extraction
- `layergrid version` from the extracted archive

## Cosign Verification

For a tagged release, verify an artifact with:

```sh
cosign verify-blob \
  --certificate layergrid_VERSION_darwin_arm64.tar.gz.pem \
  --signature layergrid_VERSION_darwin_arm64.tar.gz.sig \
  --certificate-identity-regexp '^https://github.com/layergrid/layergrid-cli' \
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

The tap publish path uses public GitHub release asset URLs. Validate `brew install layergrid/tap/layergrid` only after `layergrid/layergrid-cli` is public.

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

## Known External Blockers Before Public RC

- `brew install layergrid/tap/layergrid` requires public release assets or an authenticated private-release distribution path.
- Curl installer end-to-end validation requires a public release asset URL.
- SARIF ingestion should be validated against `layergrid/layergrid-cli` after it is public.
