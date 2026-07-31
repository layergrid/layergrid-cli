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
  --certificate layergrid_0.1.0-rc.1_darwin_arm64.tar.gz.pem \
  --signature layergrid_0.1.0-rc.1_darwin_arm64.tar.gz.sig \
  --certificate-identity-regexp '^https://github.com/layergrid/layergrid-cli' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  layergrid_0.1.0-rc.1_darwin_arm64.tar.gz
```

The installer verifies `checksums.txt` with the same issuer and repository identity before installing a binary.

## Homebrew

Design partner install command after `v0.1.0-rc.1` is published:

```sh
brew install layergrid/tap/layergrid
```

Cross-repo publishing uses `HOMEBREW_TAP_TOKEN`, scoped to `contents:write` on `layergrid/homebrew-tap`.

For `v0.1.0-rc.1`, GoReleaser published `Casks/layergrid.rb` to the private tap. A local `brew install layergrid/tap/layergrid` reached the release asset download but GitHub returned 404 because the CLI repo and release assets are still private. This is expected while the repos remain private.

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
layergrid scan ./testdata/langchain-trifecta --format sarif --output layergrid.sarif
gzip -c layergrid.sarif | base64 | tr -d '\n' > layergrid.sarif.b64
gh api repos/layergrid/sarif-test/code-scanning/sarifs --method POST --input payload.json
```

For `v0.1.0-rc.1`, SARIF upload to public repo `layergrid/sarif-test` completed and produced four open LayerGrid code scanning alerts:

- `LG-LETHAL-TRIFECTA-01`
- `LG-RAG-UNTRUSTED-01`
- `LG-TOOL-EXFIL-CHAT-01`
- `LG-AGENT-NO-GUARDRAIL-01`

## Known External Blockers Before RC

- `brew install layergrid/tap/layergrid` requires public release assets or an authenticated private-release distribution path.
- Curl installer end-to-end public URL validation is descoped from this RC batch.
