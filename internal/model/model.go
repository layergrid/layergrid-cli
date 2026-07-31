package model

import "time"

type Stack struct {
	Root        string       `json:"root"`
	ScanID      string       `json:"scanId"`
	ScannedAt   time.Time    `json:"scannedAt"`
	Agents      []Agent      `json:"agents"`
	Tools       []Tool       `json:"tools"`
	MCPServers  []MCPServer  `json:"mcpServers"`
	Datasources []Datasource `json:"datasources"`
	Models      []Model      `json:"models"`
	Errors      []ScanError  `json:"errors"`
}

type Agent struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Framework  Framework         `json:"framework"`
	Location   Location          `json:"location"`
	Tools      []ToolRef         `json:"tools"`
	Model      ModelRef          `json:"model,omitempty"`
	Memory     MemoryConfig      `json:"memory"`
	Guardrails []GuardrailRef    `json:"guardrails,omitempty"`
	SubAgents  []AgentRef        `json:"subAgents,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Framework string

const (
	FrameworkLangChain  Framework = "langchain"
	FrameworkLangGraph  Framework = "langgraph"
	FrameworkLlamaIndex Framework = "llamaindex"
	FrameworkCrewAI     Framework = "crewai"
	FrameworkAutoGen    Framework = "autogen"
	FrameworkClaudeSDK  Framework = "claude_sdk"
	FrameworkOpenAISDK  Framework = "openai_sdk"
	FrameworkCustom     Framework = "custom"
	FrameworkMCP        Framework = "mcp"
)

type Tool struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Kind        ToolKind   `json:"kind"`
	Source      ToolSource `json:"source"`
	Location    Location   `json:"location"`
	Capability  Capability `json:"capability"`
	Scope       []string   `json:"scope,omitempty"`
	MCPServerID string     `json:"mcpServerId,omitempty"`
	Description string     `json:"description,omitempty"`
	Descriptor  string     `json:"descriptor"`
}

type ToolKind string

const (
	ToolKindFunction ToolKind = "function"
	ToolKindMCP      ToolKind = "mcp"
	ToolKindHTTP     ToolKind = "http"
	ToolKindShell    ToolKind = "shell"
	ToolKindCode     ToolKind = "code_exec"
	ToolKindDB       ToolKind = "db"
)

type ToolSource struct {
	Kind string `json:"kind"`
	Name string `json:"name,omitempty"`
}

type Capability struct {
	ReadsData      DataAccess     `json:"readsData"`
	ReadsUntrusted UntrustedInput `json:"readsUntrusted"`
	Writes         WriteScope     `json:"writes"`
	Exfil          ExfilChannel   `json:"exfil"`
	Irreversible   bool           `json:"irreversible"`
	NetworkEgress  []string       `json:"networkEgress,omitempty"`
}

type DataAccess string

const (
	DataNone       DataAccess = "none"
	DataPublic     DataAccess = "public"
	DataInternal   DataAccess = "internal"
	DataSensitive  DataAccess = "sensitive"
	DataRestricted DataAccess = "restricted"
)

type UntrustedInput string

const (
	UntrustedNone  UntrustedInput = "none"
	UntrustedTool  UntrustedInput = "tool_output"
	UntrustedRAG   UntrustedInput = "rag"
	UntrustedInbox UntrustedInput = "inbox"
	UntrustedWeb   UntrustedInput = "web"
	UntrustedMCP   UntrustedInput = "mcp_external"
)

type WriteScope string

const (
	WriteNone     WriteScope = "none"
	WriteLocal    WriteScope = "local"
	WriteInternal WriteScope = "internal"
	WriteExternal WriteScope = "external"
)

type ExfilChannel string

const (
	ExfilNone  ExfilChannel = "none"
	ExfilEmail ExfilChannel = "email"
	ExfilChat  ExfilChannel = "chat"
	ExfilHTTP  ExfilChannel = "http"
	ExfilDB    ExfilChannel = "db_write"
	ExfilGit   ExfilChannel = "git"
	ExfilShell ExfilChannel = "shell"
)

type MCPServer struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Endpoint     string   `json:"endpoint"`
	Transport    string   `json:"transport"`
	Location     Location `json:"location"`
	ToolIDs      []string `json:"toolIds,omitempty"`
	AuthMode     MCPAuth  `json:"authMode"`
	Scopes       []string `json:"scopes,omitempty"`
	IsExternal   bool     `json:"isExternal"`
	Publisher    string   `json:"publisher,omitempty"`
	ManifestHash string   `json:"manifestHash,omitempty"`
	Descriptor   string   `json:"descriptor"`
}

type MCPAuth string

const (
	MCPAuthNone     MCPAuth = "none"
	MCPAuthPAT      MCPAuth = "pat"
	MCPAuthOAuth    MCPAuth = "oauth"
	MCPAuthOAuthDCR MCPAuth = "oauth_dcr"
)

type Datasource struct {
	ID          string     `json:"id"`
	Kind        string     `json:"kind"`
	Sensitivity DataAccess `json:"sensitivity"`
	Trust       string     `json:"trust"`
	Location    Location   `json:"location"`
	Description string     `json:"description,omitempty"`
}

type Model struct {
	ID       string   `json:"id"`
	Provider string   `json:"provider"`
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint,omitempty"`
	Location Location `json:"location"`
}

type Location struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Col  int    `json:"col,omitempty"`
}

type ScanError struct {
	Detector string   `json:"detector"`
	Message  string   `json:"message"`
	Location Location `json:"location"`
}

type MemoryConfig struct {
	Persistent        bool   `json:"persistent"`
	Backend           string `json:"backend,omitempty"`
	RetentionPolicy   string `json:"retentionPolicy,omitempty"`
	SharedAcrossUsers bool   `json:"sharedAcrossUsers,omitempty"`
}

type ToolRef string
type AgentRef string
type ModelRef string
type GuardrailRef string
