# Contributing

LayerGrid is built around deterministic, offline static analysis.

Before opening a PR:

- Run `go test ./...`.
- Add or update fixtures for any new rule.
- Include a "Principles Upheld" section in the PR description.
- Do not add network calls to scanner code.
- Do not add panic paths outside the CLI boundary.
