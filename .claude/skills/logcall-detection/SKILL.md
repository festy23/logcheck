---
name: logcall-detection
description: >
  How to detect and extract log calls from slog and zap in this project.
  Use this skill when working on pkg/extractor/, adding support for new
  logger methods, fixing missed log call detection, working with TypesInfo
  resolution, or understanding the LogCall model. Also use when modifying
  how message arguments are extracted from different logging APIs.
---

# Log Call Detection

## Pipeline

1. Inspector finds all `*ast.CallExpr`
2. Classify: package function vs method on instance
3. Resolve package path via TypesInfo
4. Delegate to slog.go or zap.go based on path
5. Extract message argument by index (varies per method)
6. Recursively unwrap message expression (literal, concat, Sprintf, const)
7. Extract key-value pairs for structured logging args
8. Return `[]model.LogCall`

## What we skip

- Non-literal message args: set `HasLiteral = false`, rules skip these
- Methods not in our known list: ignore silently
- Generated files with `// Code generated` comment: skip entirely
- Chained calls like `.With().Info()`: TypesInfo resolves receiver type through chain, no special handling needed

## Adding a new logger

1. Create `pkg/extractor/<loggername>.go`
2. Map all methods → message argument index
3. Map key-value extraction pattern
4. Add package path to the router in `extractor.go`
5. Add testdata covering all method variants
6. Run full test suite — MUST NOT break existing slog/zap tests

For slog method details: see `references/slog-api.md`
For zap method details: see `references/zap-api.md`
