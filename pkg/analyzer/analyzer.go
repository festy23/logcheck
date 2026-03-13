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
	flagConfig            string
)

func init() {
	Analyzer.Flags.StringVar(&flagDisable, "disable", "", "comma-separated list of rules to disable (e.g. lowercase,english)")
	Analyzer.Flags.StringVar(&flagSensitivePatterns, "sensitive-patterns", "", "comma-separated list of additional sensitive key patterns")
	Analyzer.Flags.StringVar(&flagConfig, "config", "", "path to JSON config file (e.g. .logcheck.json)")
}

func run(pass *analysis.Pass) (any, error) {
	calls := extractor.Extract(pass)
	reg := rules.NewRegistry()

	cfg := loadConfig(flagConfig)

	// Флаги имеют приоритет над конфигурационным файлом.
	disable := flagDisable
	if disable == "" && len(cfg.Disable) > 0 {
		disable = strings.Join(cfg.Disable, ",")
	}

	sensitivePatterns := flagSensitivePatterns
	if sensitivePatterns == "" && len(cfg.SensitivePatterns) > 0 {
		sensitivePatterns = strings.Join(cfg.SensitivePatterns, ",")
	}

	if disable != "" {
		for name := range strings.SplitSeq(disable, ",") {
			reg.Disable(strings.TrimSpace(name))
		}
	}

	if sensitivePatterns != "" {
		var patterns []string
		for p := range strings.SplitSeq(sensitivePatterns, ",") {
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
