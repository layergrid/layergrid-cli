package graph

import "github.com/layergrid/layergrid/internal/model"

type NodeKind string

const (
	NodeAgent NodeKind = "agent"
	NodeTool  NodeKind = "tool"
	NodeMCP   NodeKind = "mcp"
)

type Node struct {
	Kind NodeKind `json:"kind"`
	ID   string   `json:"id"`
	Name string   `json:"name"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type Graph struct {
	Nodes map[string]Node `json:"nodes"`
	Edges []Edge          `json:"edges"`
}

func Build(s *model.Stack) *Graph {
	g := &Graph{Nodes: map[string]Node{}}
	for _, agent := range s.Agents {
		g.Nodes[agent.ID] = Node{Kind: NodeAgent, ID: agent.ID, Name: agent.Name}
		for _, toolID := range agent.Tools {
			g.Edges = append(g.Edges, Edge{From: agent.ID, To: string(toolID), Type: "invokes"})
		}
		for _, child := range agent.SubAgents {
			g.Edges = append(g.Edges, Edge{From: agent.ID, To: string(child), Type: "cascades_to"})
		}
	}
	for _, tool := range s.Tools {
		g.Nodes[tool.ID] = Node{Kind: NodeTool, ID: tool.ID, Name: tool.Name}
		if tool.MCPServerID != "" {
			g.Edges = append(g.Edges, Edge{From: tool.ID, To: tool.MCPServerID, Type: "served_by"})
		}
	}
	for _, server := range s.MCPServers {
		g.Nodes[server.ID] = Node{Kind: NodeMCP, ID: server.ID, Name: server.Name}
	}
	return g
}
