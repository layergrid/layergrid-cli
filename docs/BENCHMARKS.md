# Benchmarks

Benchmarks were run on July 31, 2026 on macOS arm64, Apple M3, using the local `bin/layergrid` binary.

Reproduce:

```sh
make build
./hack/bench-realworld.sh
```

P5 budget:

- Cold scan target: under 10 seconds
- Warm scan target: under 3 seconds

All measured scans met both budgets.

| Repo | Scanned LoC | Cold seconds | Warm seconds | Peak RSS | Findings |
|---|---:|---:|---:|---:|---:|
| langchain-ai/langchain | 401,209 | 0.84 | 1.01 | 17.2 MiB | 7 |
| run-llama/llama_index | 942,761 | 1.39 | 1.42 | 17.2 MiB | 0 |
| crewAIInc/crewAI | 644,189 | 1.00 | 1.11 | 23.8 MiB | 69 |
| microsoft/autogen | 120,437 | 0.27 | 0.27 | 16.9 MiB | 197 |
| browser-use/browser-use | 126,339 | 0.17 | 0.17 | 14.5 MiB | 1 |
| openai/openai-agents-python | 314,488 | 0.47 | 0.43 | 14.4 MiB | 0 |
| anthropics/anthropic-cookbook | 170,073 | 0.06 | 0.06 | 13.9 MiB | 0 |
| modelcontextprotocol/servers | 7,546 | 0.01 | 0.01 | 11.0 MiB | 0 |
| smol-ai/developer | 1,268 | 0.01 | 0.01 | 9.8 MiB | 0 |
| mem0ai/mem0 | 137,638 | 0.18 | 0.18 | 13.8 MiB | 1 |

Notes:

- The benchmark script clones repos into `/tmp/layergrid-corpus` unless `LAYERGRID_CORPUS_DIR` is set.
- Peak RSS is captured with `/usr/bin/time -l`; on macOS this value is bytes.
- The `crewAI` and `autogen` finding totals include low-confidence, zero-score findings from tests/examples and benchmark templates. These are intentionally preserved as advisories for reviewer inspection but do not affect score.
