# LayerGrid CLI Release Readiness

This checklist tracks the remaining work before tagging `v0.1.0`.

- [x] Offline scan path with embedded rules
- [x] Deterministic finding IDs
- [x] Seed rule YAML catalog
- [x] Engine test proving every seed rule can fire
- [x] LangChain, CrewAI, LlamaIndex, AutoGen, Claude SDK, OpenAI SDK, MCP, and generic detectors
- [x] Human, JSON, SARIF, HTML, and Markdown output
- [x] Config loading and rule disable/category selection
- [x] CI for tests, lint, principles check, and dogfood scan
- [x] GoReleaser, Cosign signing hooks, SBOM declaration, Dockerfile, and verified install script scaffold
- [x] Validate GoReleaser config with `goreleaser check`
- [x] Validate GoReleaser snapshot output locally on macOS arm64, including cross-platform archives, SBOMs, generated Homebrew cask, archive extraction, and `layergrid version`
- [x] Create private Homebrew tap repo for release validation
- [x] Run and document 10-repo precision corpus
- [x] Add regression fixtures for FP classes found in corpus
- [ ] Validate SARIF upload in GitHub code scanning
- [x] Benchmark harness: 1k-file scan averages about 84ms on Apple M3
- [x] Benchmark against 10 real-world projects and record timings
- [x] Add byte-for-byte human output snapshot test
- [x] Commit screenshot-quality first scan image
- [x] Verify Homebrew cask generation in a private dry run
- [ ] Verify Homebrew tap publishing with a cross-repo release token
- [ ] Verify the curl installer against a signed pre-release artifact
- [x] Public-facing README, principles, rules, contributing, benchmark, precision, and release docs

## Blockers

- GitHub code scanning validation is blocked until Advanced Security/code scanning is enabled for `layergrid/layergrid-cli`.
- Signed release verification requires a real tag or pre-release so GitHub OIDC can issue the Cosign certificate.
- Homebrew tap publishing requires `HOMEBREW_TAP_TOKEN` to be configured in `layergrid/layergrid-cli` repository secrets.
- Curl installer end-to-end verification requires a signed release artifact and GitHub Pages to be enabled.
