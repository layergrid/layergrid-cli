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
  --certificate layergrid_0.1.0-alpha.0_darwin_arm64.tar.gz.pem \
  --signature layergrid_0.1.0-alpha.0_darwin_arm64.tar.gz.sig \
  --certificate-identity-regexp '^https://github.com/layergrid/layergrid-cli' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  layergrid_0.1.0-alpha.0_darwin_arm64.tar.gz
```

The installer verifies `checksums.txt` with the same issuer and repository identity before installing a binary.

## Homebrew

Design partner install command after `v0.1.0-rc.1` is published:

```sh
brew install layergrid/tap/layergrid
```

Cross-repo publishing uses `HOMEBREW_TAP_TOKEN`, scoped to `contents:write` on `layergrid/homebrew-tap`.

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
gh repo create layergrid/sarif-test --public
layergrid scan ./small-agent-fixture --format sarif --output layergrid.sarif
gh code-scanning upload --repo layergrid/sarif-test --sarif layergrid.sarif
```

Current blocker: GitHub code scanning API returned `Advanced Security must be enabled` for the private CLI repo. The validation should be run against a temporary public repo where code scanning is available, then the test repo should be deleted.

## Known External Blockers Before RC

- Configure `HOMEBREW_TAP_TOKEN` in `layergrid/layergrid-cli` repo secrets.
- Enable GitHub Pages for the repo and publish `scripts/install.sh`.
- Run one real tagged alpha release to verify Cosign OIDC certificates and signed artifacts end-to-end.
- Run SARIF upload against a temporary public repo.
