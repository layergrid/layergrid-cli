# Contributing

Thanks for helping harden LayerGrid during the private release-candidate phase.

Good bug reports include:

- LayerGrid version: `layergrid version`
- Command run
- Output format used
- Repo/framework scanned
- Whether the finding is true positive, false positive, or ambiguous
- A small redacted fixture if possible

Before sending a patch:

```sh
go test ./...
golangci-lint run
./scripts/principles-check.sh
```

For detector changes, add a regression fixture under `testdata/regression/<rule-id>/<case-name>/`.

For rule changes, update the YAML in `rules/` and regenerate docs:

```sh
make docs
```
