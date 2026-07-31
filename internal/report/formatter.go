package report

import "github.com/layergrid/layergrid/internal/scan"

type Formatter interface {
	Format(scan.Result) ([]byte, error)
}

func New(format string, noColor bool) Formatter {
	switch format {
	case "json":
		return JSON{}
	case "markdown":
		return Markdown{}
	case "html":
		return HTML{}
	case "sarif":
		return SARIF{}
	default:
		return Human{NoColor: noColor}
	}
}
