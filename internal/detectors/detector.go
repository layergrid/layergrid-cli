package detectors

import (
	"github.com/layergrid/layergrid-cli/internal/detectors/detectopts"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector interface {
	Name() string
	Framework() model.Framework
	Detect(root string, s *model.Stack, opts detectopts.Options) error
}
