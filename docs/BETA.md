# Design Partner Beta Guide

LayerGrid is an offline static scanner for AI agent stacks. It looks for risky composition across agents, tools, MCP servers, memory, untrusted input, and external communication channels.

This guide assumes you are comfortable reading findings in application or platform code and know the agent frameworks in your environment.

## Install

Use one path.

### Homebrew

```sh
brew install layergrid/tap/layergrid
```

### curl

macOS and Linux:

```sh
curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
```

### Go

```sh
go install github.com/layergrid/layergrid-cli/cmd/layergrid@latest
```

Confirm the install:

```sh
layergrid version
```

## First Scan

Run LayerGrid at the root of a codebase that contains agent code or MCP configuration:

```sh
layergrid scan .
```

For CI-style behavior:

```sh
layergrid scan . --fail-on high
```

For machine-readable output:

```sh
layergrid scan . --format json --output layergrid.json
layergrid scan . --format sarif --output layergrid.sarif
```

## Read The Output

The human report starts with discovered inventory:

- agents
- tools
- MCP servers
- datasources

The score is a 0-100 Trifecta Score. Higher is better. Grade A or B is generally clean enough for this beta; Grade C or lower means LayerGrid found one or more risky paths worth reviewing.

A finding path shows the composition LayerGrid inferred, for example:

```text
support-agent -> Slack MCP -> Secrets store
```

Review the rule ID, severity, location, evidence, and suggested fix. Use `layergrid explain <rule-id>` for rule details.

## Useful Feedback

If a finding looks wrong, file a false-positive issue:

https://github.com/layergrid/layergrid-cli/issues/new?template=false_positive.md

Include:

- rule ID
- minimal code or config that triggered it
- why the finding is not a real risk
- framework and version
- `layergrid version` output

If LayerGrid misses a real risk, file a bug report and label it `false-negative` if you can.

## Coming Soon

These are not expected in `v0.1.0-rc.1`:

- `layergrid share`
- cloud dashboard
- GitHub Action

They are coming soon, not broken.

Public Trifecta Score reports and leaderboard will live at [trifecta.report](https://trifecta.report) (coming soon). The goal is to make shareable scan summaries and score history available without changing the local-first scanner.
