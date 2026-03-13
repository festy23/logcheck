# logcheck — Go linter for log message validation

Linter for golangci-lint. Checks log/slog and go.uber.org/zap messages.
Rules: lowercase start, English only, no special chars/emoji, no sensitive data.

## Build & Test

- `go test ./...` — all tests
- `go test ./pkg/analyzer/ -run TestAll -v` — analysistest integration
- `go build ./cmd/logcheck/` — standalone binary
- `go vet ./...` — must pass before any commit

## Architecture

Pipeline: **extractor** (finds log calls via TypesInfo) → **rules** (checks each call) → **analyzer** (wires together).

- `pkg/model/` — LogCall, KeyValue, ConcatPart structs
- `pkg/extractor/` — detects log calls, extracts messages. Files: extractor.go (orchestrator), slog.go, zap.go, message.go
- `pkg/rules/` — each rule = one file implementing Rule interface. rule.go has interface, registry.go manages enabled rules
- `pkg/analyzer/` — analysis.Analyzer entry point + config
- `cmd/logcheck/` — singlechecker entry point
- `plugin/` — golangci-lint module plugin

## IMPORTANT: mistakes to avoid

- ALWAYS resolve packages via `pass.TypesInfo.Uses` → `*types.PkgName` → `.Imported().Path()`. NEVER match by identifier name like "slog" — import aliases break this
- `strconv.Unquote` does NOT handle backtick raw strings — check prefix and trim manually
- For `slog.InfoContext(ctx, msg)` the message is arg index 1, not 0. For `slog.Log(ctx, level, msg)` it's index 2
- Distinguish `*zap.Logger` from `*zap.SugaredLogger` via receiver type — they have different APIs
- Diagnostic message format: `logcheck: <rule-name>: <description>`
- Each rule MUST skip analysis when `call.HasLiteral == false` — can't check non-literal messages
- No external deps beyond `golang.org/x/tools`
- testdata/ files must be valid compilable Go code

## After every change

Run at minimum: `go vet ./... && go test ./... -count=1`
For thorough check: delegate to the **verifier** agent.
