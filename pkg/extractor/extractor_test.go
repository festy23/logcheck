package extractor_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/festy23/loglinter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func analyzerTestdata() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "analyzer", "testdata")
}

// TestExtract проверяет экстрактор, запуская полный анализатор на
// всех testdata-пакетах. Покрытие относится к pkg/extractor.
func TestExtract(t *testing.T) {
	analysistest.Run(t, analyzerTestdata(), analyzer.Analyzer,
		"basic", "clean", "zapbasic", "slogmethods", "alias",
		"lowercase", "english", "specialchars", "sensitive",
	)
}
