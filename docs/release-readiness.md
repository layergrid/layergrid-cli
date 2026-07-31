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
- [x] Validate GoReleaser snapshot output locally on macOS arm64, including cross-platform archives, SBOMs, generated Homebrew formula, archive extraction, and `layergrid version`
- [x] Create private Homebrew tap repo for release validation
- [x] Run and document 10-repo precision corpus
- [x] Add regression fixtures for FP classes found in corpus
- [x] Validate SARIF upload in GitHub code scanning against `layergrid/layergrid-cli`
- [x] Benchmark harness: 1k-file scan averages about 84ms on Apple M3
- [x] Benchmark against 10 real-world projects and record timings
- [x] Add byte-for-byte human output snapshot test
- [x] Commit screenshot-quality first scan image
- [x] Verify Homebrew formula generation in a private dry run
- [x] Verify Homebrew tap publishing with a cross-repo release token
- [x] Commit curl installer script
- [x] Public-facing README, principles, rules, contributing, benchmark, precision, and release docs
- [x] Verify Cosign signatures on `checksums.txt` and public release archives
- [x] Verify Homebrew install from public release assets
- [x] Verify curl installer against public release assets on macOS and Ubuntu
- [x] Cut public `v0.1.0-rc.1`

## Validation Notes

- `v0.1.0-alpha.0` was used as a throwaway signing/install dry run and deleted after validation.
- `v0.1.0-rc.1` is the current public release candidate.
- Dogfood self-scan passes with Grade A and zero findings.
- `layergrid/layergrid-cli` code scanning accepted the LayerGrid SARIF self-scan.
