# LayerGrid Baselines

`layergrid baseline` captures the trust-boundary shape of an agent stack so later scans can detect drift.

## Save a Baseline

```sh
layergrid baseline save
```

By default this writes `.layergrid/baseline.json`. Commit that file with the application so CI can compare future changes against it.

Useful flags:

```sh
layergrid baseline save ./services/agent --output .layergrid/baseline.json
layergrid baseline save --config .layergrid.yaml --quiet
```

The baseline stores:

- MCP servers
- tools
- agents
- deterministic descriptors based on name, description, and sorted scopes

It does not store source code.

## Compare a Baseline

```sh
layergrid baseline compare
```

By default this loads `.layergrid/baseline.json`, scans the current directory, and reports:

- new MCP servers
- MCP server descriptor drift
- tool descriptor drift
- new tools added to existing agents
- widened permission scopes

Useful flags:

```sh
layergrid baseline compare --baseline .layergrid/baseline.json
layergrid baseline compare --format json
layergrid baseline compare --fail-on scope-widening
```

`--fail-on` accepts:

- `scope-widening`
- `tool-added`
- `descriptor-drift`
- `any`
- `never`

Exit codes:

- `0`: no drift, or drift below the selected threshold
- `1`: drift meets the selected failure threshold
- `2`: baseline not found, scan failed, or output could not be generated

## CI Pattern

```yaml
- name: Compare LayerGrid baseline
  run: |
    layergrid baseline compare --fail-on any
```

Use `--fail-on scope-widening` if the team wants descriptor changes to appear in logs without blocking merges.
