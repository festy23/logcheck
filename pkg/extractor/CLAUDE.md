# Extractor Package

Detects log calls in Go source and extracts message content.

## Files

- `extractor.go` — orchestrator: AST traversal, package routing
- `slog.go` — log/slog specifics: method→index mapping, KV extraction
- `zap.go` — go.uber.org/zap specifics: Logger vs SugaredLogger, field extraction
- `message.go` — recursive message expression unwinding (literals, concat, Sprintf, const)

## Critical rules

- ALWAYS resolve packages via `pass.TypesInfo.Uses` → `*types.PkgName` → `.Imported().Path()`
- Handle import aliases: `import s "log/slog"` must work
- Distinguish `*zap.Logger` from `*zap.SugaredLogger` via receiver type
- `strconv.Unquote` for quoted strings, manual backtick trim for raw strings
- Skip files with `// Code generated` comment
- Changes to one logger MUST NOT break the other — run full test suite
