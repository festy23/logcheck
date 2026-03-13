package analyzer

import (
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
	return nil, nil
}
