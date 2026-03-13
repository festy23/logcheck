---
name: go-analysis
description: >
  Go static analysis framework (go/analysis) patterns and best practices.
  Use this skill whenever working on analyzer code, creating or modifying
  the analysis.Analyzer, working with ast.Node types, using pass.TypesInfo
  for type resolution, or debugging false positives/missed detections.
  Also use when you see imports of go/ast, go/types, or golang.org/x/tools.
---

# go/analysis Framework

## Analyzer structure

```go
var Analyzer = &analysis.Analyzer{
    Name:     "logcheck",
    Doc:      "checks log messages for common issues",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
}
```

## analysis.Pass — available in Run()

- `pass.Files` — `[]*ast.File` for the current package
- `pass.TypesInfo` — `*types.Info` (Uses, Defs, Types, Selections)
- `pass.Pkg` — `*types.Package`
- `pass.Report(analysis.Diagnostic{})` — report with optional SuggestedFixes
- `pass.Reportf(pos, format, args...)` — simple report
- `pass.ResultOf[dep]` — results of required analyzers

## Inspector pattern (prefer over ast.Inspect)

```go
insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
insp.Preorder(nodeFilter, func(n ast.Node) {
    call := n.(*ast.CallExpr)
    // ...
})
```

## SuggestedFixes for autofix

```go
pass.Report(analysis.Diagnostic{
    Pos:      pos,
    Message:  "logcheck: lowercase: message must start with lowercase",
    Category: "lowercase",
    SuggestedFixes: []analysis.SuggestedFix{{
        Message: "lowercase first character",
        TextEdits: []analysis.TextEdit{{
            Pos: editStart, End: editEnd,
            NewText: []byte(replacement),
        }},
    }},
})
```

For AST node patterns: see `references/ast-patterns.md`
For TypesInfo usage: see `references/type-checking.md`
