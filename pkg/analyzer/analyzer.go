package analyzer

import (
	"github.com/festy23/loglinter/pkg/extractor"
	"github.com/festy23/loglinter/pkg/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Analyzer — точка входа analysis.Analyzer для logcheck.
var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "checks log messages for common issues",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	calls := extractor.Extract(pass)
	reg := rules.NewRegistry()
	for i := range calls {
		reg.RunAll(&calls[i], pass)
	}
	return nil, nil
}
