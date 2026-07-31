package detectors

import (
	"github.com/layergrid/layergrid/internal/detectors/generic"
	"github.com/layergrid/layergrid/internal/detectors/langchain"
	"github.com/layergrid/layergrid/internal/detectors/mcp"
)

func Registry() []Detector {
	return []Detector{
		langchain.New(),
		mcp.New(),
		generic.New(),
	}
}
