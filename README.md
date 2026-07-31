# LayerGrid

LayerGrid is an offline static scanner for AI agent stacks. It discovers agents, tools, MCP servers, and capability composition risks, then reports a deterministic Trifecta Score.

## Status

Private v0.1 build. The scanner is local-only and has no cloud dependency or telemetry.

## Quick Start

```sh
go run ./cmd/layergrid scan .
go run ./cmd/layergrid list-rules
go run ./cmd/layergrid explain LG-LETHAL-TRIFECTA-01
```

Install from source:

```sh
go install github.com/layergrid/layergrid-cli/cmd/layergrid@latest
```

## Engineering Principles

- Deterministic scans: no LLM calls and no network calls in the scan path.
- Offline by default: built-in rules are embedded in the binary.
- Precision over recall: findings need a rule ID, rationale, location, and fix.
- Stable output: JSON includes `schemaVersion` and `rubricVersion`.
- Human-readable by default: terminal output is meant to be understandable in one glance.

Release tracking lives in [docs/release-readiness.md](docs/release-readiness.md).

## License

Apache-2.0.
