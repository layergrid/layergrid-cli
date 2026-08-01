# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to semantic versioning.

## [0.1.0] - 2026-08-01

### Added

- Offline `layergrid` CLI for static scanning of AI agent stacks.
- Agent stack discovery for agents, tools, MCP servers, datasources, memory, and risky capability paths.
- Detector coverage for LangChain, CrewAI, LlamaIndex, AutoGen, MCP, Claude SDK, OpenAI SDK, and generic Python/agent patterns.
- Twenty embedded seed rules:
  - `LG-MCP-DCR-01`
  - `LG-MCP-EXTERNAL-WRITE-01`
  - `LG-MCP-NOAUTH-01`
  - `LG-MCP-OVERSCOPE-01`
  - `LG-MCP-PUBLISHER-UNKNOWN-01`
  - `LG-MEMORY-CROSS-USER-01`
  - `LG-MEMORY-UNBOUNDED-01`
  - `LG-AGENT-NO-GUARDRAIL-01`
  - `LG-AUTOGEN-LOCAL-EXEC-01`
  - `LG-CREDENTIAL-ENV-IN-CONTEXT-01`
  - `LG-CREDENTIAL-KEY-HARDCODED-01`
  - `LG-RAG-UNTRUSTED-01`
  - `LG-TOOL-CODE-EXEC-01`
  - `LG-TOOL-EXFIL-CHAT-01`
  - `LG-TOOL-EXFIL-DBWRITE-01`
  - `LG-TOOL-EXFIL-EMAIL-01`
  - `LG-TOOL-SHELL-EXEC-01`
  - `LG-LETHAL-TRIFECTA-01`
  - `LG-LETHAL-TRIFECTA-02`
  - `LG-LETHAL-TRIFECTA-03`
- CLI commands:
  - `layergrid scan`
  - `layergrid list-rules`
  - `layergrid explain`
  - `layergrid init`
  - `layergrid version`
- Human, JSON, SARIF, HTML, and Markdown report formats.
- Config support for excludes, disabled rules, category selection, and fail thresholds.
- Deterministic finding IDs suitable for suppressions, SARIF, and score trend analysis.
- Trifecta Score rubric v1.0 with grades, severity counts, and low-confidence zero-score findings.
- JSON schema version `1.0.0` and rubric version `1.0` in machine-readable output.
- Golden tests for human output and SARIF/JSON report shape.
- Regression fixtures for false-positive classes found during the 10-repo corpus pass.
- Release infrastructure with GoReleaser, Cosign keyless signing, SBOM generation, checksums, and GitHub Releases.
- Homebrew formula publishing to `layergrid/homebrew-tap`.
- Signature-verifying curl installer served from GitHub Pages.
- GitHub Code Scanning SARIF validation.
- Public documentation for rules, precision, benchmarks, principles, release process, beta onboarding, and contribution workflow.

### Notes

- LayerGrid scans locally and does not call an LLM or require a cloud account during scanning.
- `trifecta.report` is reserved for public Trifecta Score reports and leaderboard work; the site is coming soon.
