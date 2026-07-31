# Engineering Principles

LayerGrid is built to feel like security tooling engineers already trust: deterministic, offline, fast, inspectable, and boring in the best way.

## Deterministic

The same input, LayerGrid version, and rule library should produce the same findings. The scanner does not call an LLM, download rules, or depend on network lookups during scanning.

## Offline

Rules are embedded in the binary. A user can install LayerGrid once and run it in an air-gapped environment.

## Precise

Findings need a rule ID, a location, a rationale, and a concrete fix. Low-confidence matches are marked as low confidence and do not affect score.

## Panic-Free

Malformed files should produce scan errors, not stack traces.

## Fast

The scanner is designed for CI and quick local feedback. Real-world scans in the v0.1 corpus complete well under the 10-second cold-start budget.

## Single Binary

The core CLI is Go and ships as a static binary. No Python, Node, Ruby, Docker, or cloud account is required to scan.

## Stable Output

JSON includes `schemaVersion` and `rubricVersion`. Finding IDs are deterministic so suppressions and score trends can work across runs.

## Human-Readable

The default report should be screenshotable and understandable without opening a dashboard.

## Inspectable Rules

Rules are YAML files in this repository. `layergrid list-rules` and `layergrid explain <rule-id>` expose what the scanner checks and why.

## No Telemetry by Default

LayerGrid does not phone home during scan.
