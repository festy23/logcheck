---
name: analysistest
description: >
  How to write tests for Go analyzers using the analysistest framework.
  Use this skill when creating or modifying test files in testdata/,
  writing analyzer_test.go, debugging test failures with // want
  comments, or working with SuggestedFixes golden files. Also use
  when you see analysistest import or "// want" in conversation.
---

# Testing with analysistest

## Integration test entry point

```go
// pkg/analyzer/analyzer_test.go
func TestAll(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, analyzer.Analyzer,
        "lowercase", "english", "specialchars", "sensitive",
    )
}
```

`TestData()` returns path to `testdata/` relative to the test file.
Last args are package names inside `testdata/src/`.

## // want comments

```go
slog.Info("Starting server") // want `logcheck: lowercase: message must start with lowercase`
slog.Info("starting server") // no comment = assert ZERO diagnostics here
```

- String after `// want` is a REGEX matched against diagnostic message
- Multiple diagnostics on one line: `// want "msg1" "msg2"`
- Lines WITHOUT `// want` MUST produce zero diagnostics — this catches false positives
- Use backtick quotes to avoid escaping: `` // want `logcheck: lowercase:` ``

## Testdata file structure

```go
package lowercase

import "log/slog"

// SHOULD trigger
func bad() {
    slog.Info("Starting server") // want `logcheck: lowercase:`
}

// Should NOT trigger
func good() {
    slog.Info("starting server")
    slog.Info("123 items processed")
    slog.Info("")
}
```

CRITICAL: testdata files must be valid, compilable Go code.
Missing imports = compilation error = test failure.

## SuggestedFixes testing

Use `analysistest.RunWithSuggestedFixes` instead of `Run`.
Create golden file alongside test file:

```
testdata/src/lowercase/
├── lowercase.go           # test input with // want
└── lowercase.go.golden    # expected output after autofix
```

## Common pitfalls

- Regex special chars in `// want` must be escaped: `.` matches anything
- For zap tests: testdata may need its own go.mod to resolve zap imports
- analysistest runs the type checker — type errors cause test failure
- Always include BOTH positive and negative test cases
