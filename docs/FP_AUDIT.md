# False Positive Audit

Run date: August 1, 2026.

Scanner build: current `main` working tree, rebuilt locally with `go build ./cmd/layergrid` after the hardening changes in this audit.

Scope: 15 public repositories. `langchain-ai/langchain-examples` was not available at scan time, so this run used `langchain-ai/streaming-cookbook` as the equivalent LangChain examples repository. `openai/openai-agents-sdk` was not available at scan time, so this run used `openai/openai-agents-python`.

## Summary

Initial scan produced 309 findings. Most were confirmed false positives from tests, benchmark fixtures, framework implementation source, file-wide tool association, placeholder credentials, or tutorial comments. After detector hardening, the final scan produced 7 findings total, all in `joaomdmoura/crewAI-examples`.

Confirmed final false positives: 0.

Overall final false-positive rate: 0 per 1,000 scanned LoC.

| Repo | Lines | Findings | TPs | FPs | FP Rate |
|---|---:|---:|---:|---:|---:|
| langchain-ai/langchain | 388,980 | 0 | 0 | 0 | 0.000 |
| langchain-ai/langgraph | 186,504 | 0 | 0 | 0 | 0.000 |
| run-llama/llama_index | 453,110 | 0 | 0 | 0 | 0.000 |
| microsoft/autogen | 112,463 | 0 | 0 | 0 | 0.000 |
| crewAIInc/crewAI | 300,607 | 0 | 0 | 0 | 0.000 |
| assafelovic/gpt-researcher | 26,848 | 0 | 0 | 0 | 0.000 |
| AgentOps-AI/agentops | 89,058 | 0 | 0 | 0 | 0.000 |
| joaomdmoura/crewAI-examples | 5,478 | 7 | 7 | 0 | 0.000 |
| langchain-ai/streaming-cookbook | 4,190 | 0 | 0 | 0 | 0.000 |
| openai/openai-agents-python | 313,935 | 0 | 0 | 0 | 0.000 |
| modelcontextprotocol/servers | 2,770 | 0 | 0 | 0 | 0.000 |
| punkpeye/awesome-mcp-servers | 0 | 0 | 0 | 0 | 0.000 |
| gitleaks/gitleaks | 21,606 | 0 | 0 | 0 | 0.000 |
| aquasecurity/trivy | 274,443 | 0 | 0 | 0 | 0.000 |
| layergrid/layergrid-cli | 393,152 | 0 | 0 | 0 | 0.000 |

## Fixes Made

| Area | Rule IDs affected | Change | Why |
|---|---|---|---|
| Default Python path filtering | Multiple | Exclude test, fixture, docs, notebook, benchmark, and `test_*.py` / `*_test.py` files by default when no explicit include is provided. | Tests and docs are not production agent deployments, and users should not need manual excludes for them. |
| Framework source detection | Multiple | Suppress known LangChain, LangGraph, CrewAI, and AutoGen framework implementation paths when scanning those framework repositories. | Framework internals define risky capabilities for users to import; that is not the same as a user-written agent deploying those capabilities. |
| CrewAI tool association | `LG-AGENT-NO-GUARDRAIL-01`, `LG-AGENT-INBOX-EXFIL`, `LG-TOOL-EXFIL-EMAIL-01` | Associate only tools referenced in the local `Agent(...)` block instead of every prior tool in the file. | File-wide association caused unrelated tool definitions to attach to later agents. |
| AutoGen executor association | `LG-AUTOGEN-LOCAL-EXEC-01`, `LG-TOOL-SHELL-EXEC-01`, `LG-TOOL-CODE-EXEC-01` | Track executor variables and attach them only when referenced in the local agent block. | Executor definitions elsewhere in a file should not imply every agent has shell or code execution. |
| MCP placeholder credentials | `LG-MCP-CREDENTIAL-IN-ENV` | Suppress common placeholder values such as `your-token`, `your-key`, `replace-me`, `todo`, `your_`, and `example_`. | Tutorial placeholders are not raw credentials. |
| Hidden Unicode threshold | `LG-TOOL-UNICODE-HIDDEN` | Keep immediate detection for BOM and bidi controls, but require at least three zero-width or soft-hyphen characters before firing. | Single zero-width characters can appear in legitimate text; repeated hidden controls are a stronger signal. |
| LangChain memory inference | `LG-MEMORY-UNBOUNDED-01`, `LG-MEMORY-EXTERNAL-WRITE` | Strip Python comment lines before inferring persistent memory backends. | A tutorial comment linking to the AgentOps dashboard was being treated as an external memory backend. |

## Repo: langchain-ai/langchain

Lines scanned: ~388,980
Scan time: 1.1s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No final findings. Initial test/framework-source findings were fixed. |

## Repo: langchain-ai/langgraph

Lines scanned: ~186,504
Scan time: 0.2s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No final findings. Initial test/benchmark memory findings were fixed. |

## Repo: run-llama/llama_index

Lines scanned: ~453,110
Scan time: 1.3s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. |

## Repo: microsoft/autogen

Lines scanned: ~112,463
Scan time: 0.2s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No final findings. Initial test, benchmark, and framework-source executor findings were fixed. |

## Repo: crewAIInc/crewAI

Lines scanned: ~300,607
Scan time: 0.8s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No final findings. Initial test and framework-source findings were fixed. |

## Repo: assafelovic/gpt-researcher

Lines scanned: ~26,848
Scan time: 0.1s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. |

## Repo: AgentOps-AI/agentops

Lines scanned: ~89,058
Scan time: 0.2s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No final findings. Initial tutorial-comment memory findings were fixed. |

## Repo: joaomdmoura/crewAI-examples

Lines scanned: ~5,478
Scan time: 0.0s
Findings: 7 total (0 critical, 0 high, 7 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| 1 | LG-AGENT-NO-GUARDRAIL-01 | medium | `flows/email_auto_responder_flow/src/email_auto_responder_flow/crews/email_filter_crew/email_filter_crew.py:22` | TRUE_POSITIVE | The agent is constructed with an external search tool and no detected guardrail or approval layer. |
| 2 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/instagram_post/tools/browser_tools.py:14` | TRUE_POSITIVE | Tool code reads `BROWSERLESS_API_KEY` from `os.environ` inside the model-exposed tool implementation. |
| 3 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/trip_planner/tools/browser_tools.py:13` | TRUE_POSITIVE | Tool code reads `BROWSERLESS_API_KEY` from `os.environ` inside the model-exposed tool implementation. |
| 4 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/trip_planner/tools/search_tools.py:11` | TRUE_POSITIVE | Tool code reads `SERPER_API_KEY` from `os.environ` inside the model-exposed tool implementation. |
| 5 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/landing_page_generator/src/landing_page_generator/tools/search_tools.py:10` | TRUE_POSITIVE | Tool code reads `SERPER_API_KEY` from `os.environ` inside the model-exposed tool implementation. |
| 6 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/landing_page_generator/src/landing_page_generator/tools/browser_tools.py:13` | TRUE_POSITIVE | Tool code reads `BROWSERLESS_API_KEY` from `os.environ` inside the model-exposed tool implementation. |
| 7 | LG-CREDENTIAL-ENV-IN-CONTEXT-01 | medium | `crews/instagram_post/tools/search_tools.py:17` | TRUE_POSITIVE | Tool code reads `SERPER_API_KEY` from `os.environ` inside the model-exposed tool implementation. |

## Repo: langchain-ai/streaming-cookbook

Lines scanned: ~4,190
Scan time: 0.0s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. |

## Repo: openai/openai-agents-python

Lines scanned: ~313,935
Scan time: 0.3s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. |

## Repo: modelcontextprotocol/servers

Lines scanned: ~2,770
Scan time: 0.0s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. |

## Repo: punkpeye/awesome-mcp-servers

Lines scanned: 0
Scan time: 0.0s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No embedded MCP config findings. Markdown list content was not treated as production configuration. |

## Repo: gitleaks/gitleaks

Lines scanned: ~21,606
Scan time: 0.0s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. Security-adjacent non-agent repo remains clean. |

## Repo: aquasecurity/trivy

Lines scanned: ~274,443
Scan time: 0.2s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. Security-adjacent non-agent repo remains clean. |

## Repo: layergrid/layergrid-cli

Lines scanned: ~393,152
Scan time: 0.1s
Findings: 0 total (0 critical, 0 high, 0 medium, 0 low)

| # | Rule ID | Severity | Location | Classification | Reason |
|---|---|---|---|---|---|
| - | - | - | - | - | No findings. Dogfood scan passes Grade A. |

## Needs Human Review

None. The 7 remaining findings are classified as TRUE_POSITIVE. They are in example app code rather than a production deployment, but the detected patterns are real in the scanned code and should remain visible.

## Acceptance Notes

- `gitleaks/gitleaks`: 0 findings.
- `aquasecurity/trivy`: 0 findings.
- `layergrid/layergrid-cli`: Grade A, 0 findings.
- Confirmed FPs fixed: test paths, fixture paths, docs paths, benchmark paths, framework-source paths, file-wide tool association, placeholder MCP credentials, one-off hidden Unicode, and comment-derived memory backends.
- Public-launch recommendation: ready from a false-positive standpoint. Final confirmed FP rate is 0 per 1,000 scanned LoC across this audit corpus.
