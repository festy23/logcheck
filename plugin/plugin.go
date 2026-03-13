//go:build plugin

package main

import (
	"github.com/festy23/logcheck/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("logcheck", func(conf any) (register.LinterPlugin, error) {
		if m, ok := conf.(map[string]any); ok {
			if v, ok := m["disable"].(string); ok {
				_ = analyzer.Analyzer.Flags.Set("disable", v)
			}
			if v, ok := m["sensitive-patterns"].(string); ok {
				_ = analyzer.Analyzer.Flags.Set("sensitive-patterns", v)
			}
			if v, ok := m["config"].(string); ok {
				_ = analyzer.Analyzer.Flags.Set("config", v)
			}
		}
		return &plugin{}, nil
	})
}

type plugin struct{}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
