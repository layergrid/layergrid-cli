# LayerGrid Rules

Generated from the embedded YAML rule library.

| Rule | Severity | Category | Description | References |
|---|---|---|---|---|
| `LG-MCP-DCR-01` | medium | mcp | OAuth Dynamic Client Registration can widen trust if not pinned and reviewed. | https://modelcontextprotocol.io |
| `LG-MCP-EXTERNAL-WRITE-01` | high | mcp | An agent appears to combine sensitive-data reads with an external MCP write channel. | https://modelcontextprotocol.io |
| `LG-MCP-NOAUTH-01` | high | mcp | An externally reachable MCP server has no detected authentication mode. | https://modelcontextprotocol.io |
| `LG-MCP-OVERSCOPE-01` | high | mcp | An MCP server grants wildcard scope, making least-privilege review impossible. | https://modelcontextprotocol.io |
| `LG-MCP-PUBLISHER-UNKNOWN-01` | medium | mcp | MCP server publisher could not be identified from the local config. | https://modelcontextprotocol.io |
| `LG-MEMORY-CROSS-USER-01` | high | memory | Cross-user shared memory can allow one user to poison another user's agent context. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-MEMORY-UNBOUNDED-01` | medium | memory | Persistent memory without retention controls can preserve prompt-injection payloads and sensitive content. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-AGENT-NO-GUARDRAIL-01` | medium | permissions | Agent has external communication tools and no detected guardrail middleware. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-AUTOGEN-LOCAL-EXEC-01` | critical | permissions | AutoGen local command execution allows broad system access from agent flows. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-CREDENTIAL-ENV-IN-CONTEXT-01` | medium | permissions | Inline environment credential reads can leak secrets into model context. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-CREDENTIAL-KEY-HARDCODED-01` | high | permissions | A credential-like literal appears in agent code. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-RAG-UNTRUSTED-01` | high | permissions | User-uploaded RAG content can carry prompt-injection payloads into an agent. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-TOOL-CODE-EXEC-01` | high | permissions | Agent code execution increases prompt-injection impact and data exfiltration risk. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-TOOL-EXFIL-CHAT-01` | high | permissions | An agent can read sensitive data and post externally to chat. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-TOOL-EXFIL-DBWRITE-01` | high | permissions | An agent can read sensitive data and write to an external database. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-TOOL-EXFIL-EMAIL-01` | high | permissions | An agent can read untrusted inbox content and send email externally. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-TOOL-SHELL-EXEC-01` | critical | permissions | Shell execution gives an agent a broad external communication and system mutation surface. | https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/ |
| `LG-LETHAL-TRIFECTA-01` | critical | trifecta | A single agent combines sensitive data access, untrusted input, and an external communication channel. | https://simonwillison.net/2025/Jun/16/the-lethal-trifecta/ |
| `LG-LETHAL-TRIFECTA-02` | critical | trifecta | Multiple agents appear connected by shared memory and collectively satisfy all three lethal-trifecta legs. | https://simonwillison.net/2025/Jun/16/the-lethal-trifecta/ |
| `LG-LETHAL-TRIFECTA-03` | high | trifecta | A parent/sub-agent handoff path appears to combine lethal-trifecta capabilities. | https://simonwillison.net/2025/Jun/16/the-lethal-trifecta/ |
