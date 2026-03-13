package analyzer_test

import (
	"testing"

	"github.com/festy23/logcheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer.Analyzer,
		"basic", "lowercase", "english", "specialchars", "sensitive",
	)
	analysistest.Run(t, testdata, analyzer.Analyzer,
		"clean", "zapbasic", "slogmethods", "alias",
	)
}
