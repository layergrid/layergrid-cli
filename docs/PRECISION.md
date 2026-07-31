# Precision Report

Run date: July 31, 2026.

Corpus:

- `langchain-ai/langchain`
- `run-llama/llama_index`
- `crewAIInc/crewAI`
- `microsoft/autogen`
- `browser-use/browser-use`
- `openai/openai-agents-python`
- `anthropics/anthropic-cookbook`
- `modelcontextprotocol/servers`
- `smol-ai/developer`
- `mem0ai/mem0`

Reproduce:

```sh
make build
./hack/scan-corpus.sh
for f in .corpus-results/*.json; do
  name=$(basename "$f" .json)
  printf "\n== %s ==\n" "$name"
  jq -r '.findings[]? | [.ruleId,.confidence,.scoreImpact] | @tsv' "$f" | sort | uniq -c | sort -nr
done
```

## Summary

The first corpus pass exposed noisy detections from fake docs credentials, arbitrary `os.environ` reads, file-wide OpenAI SDK tool inference, and CrewAI cascade inference. Those are now regression-covered or downgraded:

- Fake `sk-...`, short fake `sk-proj-*`, and placeholder `ghp_XXXX...` values no longer fire as high-confidence hardcoded keys.
- Generic `os.environ` reads no longer emit `LG-CREDENTIAL-ENV-IN-CONTEXT-01`; that rule now requires detector evidence from a tool-like block.
- OpenAI SDK code execution tools are tied to the local `assistants.create` / `responses.create` block instead of the whole file.
- CrewAI cascade inference no longer connects every prior agent in a file to every `Crew(` call.
- Findings in tests, examples, docs, notebooks, and cookbook paths emit `confidence: low` and `scoreImpact: 0`.

Regression fixtures added:

- `testdata/regression/LG-CREDENTIAL-KEY-HARDCODED-01/fake-doc-key`
- `testdata/regression/LG-TOOL-CODE-EXEC-01/openai-source-support`

## Counts

`TP` means the risky capability or composition is genuinely present. `FP` means a high-confidence finding was wrong after tuning. `Ambiguous` means the finding is real capability inventory in framework/library source, but not necessarily a deployable app vulnerability. Low-confidence zero-score findings are tracked separately.

| Rule | TP | FP | Ambiguous | Low-confidence zero-score | Notes |
|---|---:|---:|---:|---:|---|
| LG-AGENT-NO-GUARDRAIL-01 | 0 | 0 | 8 | 117 | Ambiguous cases are framework source paths with external/code tools but no app-level guardrail context. |
| LG-AUTOGEN-LOCAL-EXEC-01 | 11 | 0 | 0 | 56 | High-confidence findings are actual `LocalCommandLineCodeExecutor` construction or exported local executor capabilities. |
| LG-CREDENTIAL-ENV-IN-CONTEXT-01 | 1 | 0 | 0 | 0 | One CrewAI tool reads an API token inside a tool implementation. |
| LG-CREDENTIAL-KEY-HARDCODED-01 | 0 | 0 | 0 | 3 | Remaining detections are placeholders/examples and low-confidence. |
| LG-MCP-PUBLISHER-UNKNOWN-01 | 1 | 0 | 0 | 0 | `mem0` MCP config declares a remote server without a publisher signal. |
| LG-MEMORY-UNBOUNDED-01 | 3 | 0 | 0 | 0 | Persistent/checkpointer memory detected without retention metadata. |
| LG-TOOL-CODE-EXEC-01 | 1 | 0 | 0 | 0 | AutoGen code-executor agent capability. |
| LG-TOOL-EXFIL-EMAIL-01 | 0 | 0 | 0 | 11 | All observed cases are in tests/examples and score zero. |
| LG-TOOL-SHELL-EXEC-01 | 5 | 0 | 0 | 59 | High-confidence findings are agents wired to local shell/code execution tools. |
| All other seed rules | 0 | 0 | 0 | 0 | No hits in this corpus. Covered by engine tests and targeted fixtures. |

## Precision Budget

Across the 10-repo corpus, there were no known high-confidence false positives after tuning. Ambiguous findings remain below 1 per 1,000 scanned LoC per rule:

- `LG-AGENT-NO-GUARDRAIL-01`: 8 ambiguous findings across 2,865,948 scanned LoC, about 0.003 per 1,000 LoC.

Low-confidence findings do not deduct from score and are meant to guide manual review in framework/example-heavy repositories.
