package rules

import "embed"

//go:embed trifecta/*.yaml mcp/*.yaml permissions/*.yaml memory/*.yaml
var FS embed.FS
