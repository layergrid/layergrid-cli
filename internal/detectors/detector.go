package detectors

import "github.com/layergrid/layergrid-cli/internal/model"

type Detector interface {
	Name() string
	Framework() model.Framework
	Detect(root string, s *model.Stack) error
}
