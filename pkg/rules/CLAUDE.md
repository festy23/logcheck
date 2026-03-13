# Rules Package

Every rule implements the Rule interface from `rule.go`:

```go
type Rule interface {
    Name() string
    Description() string
    Check(call *model.LogCall, pass *analysis.Pass)
}
```

Each rule is ONE struct in ONE file: `<name>.go` with matching `<name>_test.go`.

## Conventions

- Rules receive `*model.LogCall` with pre-extracted message data
- DO NOT re-parse AST inside rules — use LogCall fields
- Use `call.MsgLiteral` for string content checks (lowercase, english, specialchars)
- Use `call.ConcatParts` and `call.KeyValues` for data flow checks (sensitive)
- ALWAYS skip when `call.HasLiteral == false`
- Report via `pass.Report(analysis.Diagnostic{...})`:
  - Pos: `call.MsgPos`
  - Message: `"logcheck: <rule-name>: <specific problem>"`
  - Category: `"<rule-name>"`
- For autofix: add SuggestedFixes to the Diagnostic
- Register new rules in `registry.go` defaultRules slice
