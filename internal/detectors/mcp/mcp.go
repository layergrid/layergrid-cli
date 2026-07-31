package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "mcp" }
func (Detector) Framework() model.Framework { return model.FrameworkMCP }

type configFile struct {
	MCPServers map[string]serverConfig `json:"mcpServers"`
}

type serverConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
	URL     string            `json:"url"`
	Scopes  []string          `json:"scopes"`
	Auth    string            `json:"auth"`
}

func (Detector) Detect(root string, s *model.Stack) error {
	names := map[string]bool{
		"mcp.json": true, "mcp_config.json": true, "claude_desktop_config.json": true, "mcp-servers.json": true,
	}
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".venv", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if names[d.Name()] || strings.HasSuffix(filepath.ToSlash(path), "/.cursor/mcp.json") || strings.HasSuffix(filepath.ToSlash(path), "/.vscode/mcp.json") {
			if err := detectConfig(root, path, s); err != nil {
				s.Errors = append(s.Errors, model.ScanError{Detector: "mcp", Message: err.Error(), Location: model.RelativeLocation(root, path, 1)})
			}
		}
		return nil
	})
}

func detectConfig(root, path string, s *model.Stack) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for name, srv := range cfg.MCPServers {
		serverID := model.StableID("mcp", path, name)
		auth := inferAuth(srv)
		server := model.MCPServer{
			ID:         serverID,
			Name:       name,
			Endpoint:   endpoint(srv),
			Transport:  transport(srv),
			Location:   model.RelativeLocation(root, path, 1),
			AuthMode:   auth,
			Scopes:     srv.Scopes,
			IsExternal: isExternal(srv),
			Publisher:  publisher(srv),
			Descriptor: model.Descriptor(name, endpoint(srv), strings.Join(srv.Args, " ")),
		}
		tool := model.Tool{
			ID:          model.StableID("tool", "mcp", path, name),
			Name:        name + " MCP",
			Kind:        model.ToolKindMCP,
			Source:      model.ToolSource{Kind: "mcp", Name: name},
			Location:    server.Location,
			Capability:  inferCapability(name, srv, server),
			Scope:       srv.Scopes,
			MCPServerID: serverID,
			Descriptor:  server.Descriptor,
		}
		server.ToolIDs = []string{tool.ID}
		s.MCPServers = append(s.MCPServers, server)
		s.Tools = append(s.Tools, tool)
	}
	return nil
}

func inferAuth(s serverConfig) model.MCPAuth {
	if strings.EqualFold(s.Auth, "oauth_dcr") || containsFold(s.Scopes, "dynamic_client_registration") {
		return model.MCPAuthOAuthDCR
	}
	if strings.Contains(strings.ToLower(s.Auth), "oauth") {
		return model.MCPAuthOAuth
	}
	for k := range s.Env {
		upper := strings.ToUpper(k)
		if strings.Contains(upper, "TOKEN") || strings.Contains(upper, "KEY") || strings.Contains(upper, "SECRET") {
			return model.MCPAuthPAT
		}
	}
	return model.MCPAuthNone
}

func endpoint(s serverConfig) string {
	if s.URL != "" {
		return s.URL
	}
	if s.Command != "" {
		return strings.TrimSpace(s.Command + " " + strings.Join(s.Args, " "))
	}
	return "stdio"
}

func transport(s serverConfig) string {
	if s.URL != "" {
		if strings.Contains(s.URL, "/sse") {
			return "sse"
		}
		return "http"
	}
	return "stdio"
}

func isExternal(s serverConfig) bool {
	if s.URL == "" {
		return false
	}
	lower := strings.ToLower(s.URL)
	return strings.HasPrefix(lower, "http") &&
		!strings.Contains(lower, "localhost") &&
		!strings.Contains(lower, "127.0.0.1") &&
		!strings.Contains(lower, "10.") &&
		!strings.Contains(lower, "192.168.")
}

func publisher(s serverConfig) string {
	args := strings.Join(s.Args, " ")
	switch {
	case strings.Contains(args, "@modelcontextprotocol/"):
		return "modelcontextprotocol"
	case strings.Contains(args, "@"):
		return strings.Fields(args)[len(strings.Fields(args))-1]
	default:
		return "unknown"
	}
}

func inferCapability(name string, s serverConfig, server model.MCPServer) model.Capability {
	lower := strings.ToLower(name + " " + strings.Join(s.Args, " ") + " " + s.URL + " " + strings.Join(s.Scopes, " "))
	cap := model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
	if server.IsExternal {
		cap.ReadsUntrusted = model.UntrustedMCP
	}
	if strings.Contains(lower, "github") || strings.Contains(lower, "git") {
		cap.ReadsData = model.DataSensitive
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilGit
	}
	if strings.Contains(lower, "slack") || strings.Contains(lower, "discord") || strings.Contains(lower, "teams") {
		cap.ReadsUntrusted = model.UntrustedInbox
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilChat
	}
	if strings.Contains(lower, "gmail") || strings.Contains(lower, "email") {
		cap.ReadsUntrusted = model.UntrustedInbox
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilEmail
	}
	if strings.Contains(lower, "supabase") || strings.Contains(lower, "postgres") || strings.Contains(lower, "db") {
		cap.ReadsData = model.DataSensitive
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilDB
	}
	return cap
}

func containsFold(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(value, needle) {
			return true
		}
	}
	return false
}
