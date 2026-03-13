package analyzer

import (
	"strings"

	"github.com/festy23/logcheck/pkg/extractor"
	"github.com/festy23/logcheck/pkg/rules"
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

var (
	flagDisable           string
	flagSensitivePatterns string
)

func init() {
	Analyzer.Flags.StringVar(&flagDisable, "disable", "", "comma-separated list of rules to disable (e.g. lowercase,english)")
	Analyzer.Flags.StringVar(&flagSensitivePatterns, "sensitive-patterns", "", "comma-separated list of additional sensitive key patterns")
}

func run(pass *analysis.Pass) (any, error) {
	calls := extractor.Extract(pass)
	reg := rules.NewRegistry()

	if flagDisable != "" {
		for _, name := range strings.Split(flagDisable, ",") {
			reg.Disable(strings.TrimSpace(name))
		}
	}

	if flagSensitivePatterns != "" {
		var patterns []string
		for _, p := range strings.Split(flagSensitivePatterns, ",") {
			if s := strings.TrimSpace(p); s != "" {
				patterns = append(patterns, s)
			}
		}
		reg.AddSensitivePatterns(patterns)
	}

	for i := range calls {
		reg.RunAll(&calls[i], pass)
	}
	return nil, nil
}
