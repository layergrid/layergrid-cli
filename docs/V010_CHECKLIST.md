# v0.1.0 Final Checklist

Check these before tagging `v0.1.0`.

- [ ] All open design-partner bug reports resolved or explicitly deferred
- [ ] `go test ./...` passes on main
- [ ] `golangci-lint run` passes
- [ ] `./scripts/principles-check.sh` passes
- [ ] `layergrid scan .` on `layergrid-cli` dogfoods at Grade A or B
- [ ] `docs/PRECISION.md` is current; re-run if any detector changed
- [ ] `docs/BENCHMARKS.md` is current
- [ ] `docs/RULES.md` regenerated
- [ ] No deprecated Trifecta Score domain references remaining in codebase; use `trifecta.report`
- [ ] README screenshot is current; retake if human output changed
- [ ] `CHANGELOG.md` is written for `v0.1.0`
- [ ] RC tag is NOT the final tag; `v0.1.0` must tag from current main
- [ ] `goreleaser check` passes
- [ ] Signed release artifacts are verified post-cut
