package rules_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/festy23/logcheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func analyzerTestdata() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "analyzer", "testdata")
}

// TestRulesIntegration проверяет все правила, запуская полный анализатор
// на всех testdata-пакетах. Покрытие относится к pkg/rules.
func TestRulesIntegration(t *testing.T) {
	analysistest.Run(t, analyzerTestdata(), analyzer.Analyzer,
		"basic", "clean", "zapbasic", "slogmethods", "alias",
		"lowercase", "english", "specialchars", "sensitive",
	)
}
