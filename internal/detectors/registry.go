package detectors

import (
	"github.com/layergrid/layergrid-cli/internal/detectors/autogen"
	"github.com/layergrid/layergrid-cli/internal/detectors/claude_sdk"
	"github.com/layergrid/layergrid-cli/internal/detectors/crewai"
	"github.com/layergrid/layergrid-cli/internal/detectors/generic"
	"github.com/layergrid/layergrid-cli/internal/detectors/langchain"
	"github.com/layergrid/layergrid-cli/internal/detectors/llamaindex"
	"github.com/layergrid/layergrid-cli/internal/detectors/mcp"
	"github.com/layergrid/layergrid-cli/internal/detectors/openai_sdk"
)

func Registry() []Detector {
	return []Detector{
		langchain.New(),
		crewai.New(),
		llamaindex.New(),
		autogen.New(),
		claude_sdk.New(),
		openai_sdk.New(),
		mcp.New(),
		generic.New(),
	}
}
