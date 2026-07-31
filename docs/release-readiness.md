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
- [ ] Validate GoReleaser output with a dry run on Linux and macOS
- [ ] Validate SARIF upload in a private GitHub code-scanning repo
- [ ] Benchmark against a real 10k+ file project and record timings
- [ ] Pixel-review human output against launch-demo copy
- [ ] Verify Homebrew tap publishing in a private dry run
- [ ] Verify the curl installer against a signed pre-release artifact
