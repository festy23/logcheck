Add a new lint rule named $ARGUMENTS to the logcheck linter.

## Deliverables

1. **Rule**: `pkg/rules/$ARGUMENTS.go` — implements Rule interface from rule.go
2. **Unit tests**: `pkg/rules/$ARGUMENTS_test.go` — test logic with constructed LogCall structs
3. **Testdata**: `testdata/src/$ARGUMENTS/$ARGUMENTS.go` — both slog and zap, lines with `// want` and clean lines
4. **Registration**: add to defaultRules in `pkg/rules/registry.go`
5. **Autofix** (if possible): SuggestedFixes in Diagnostic + `.golden` file

## Workflow

Start in plan mode. Think through:
- What exactly does this rule check?
- What are the edge cases and potential false positives?
- Can it be autofixed safely?

Read existing rules in `pkg/rules/` as reference before writing.
Read `pkg/rules/rule.go` for the interface and `pkg/model/logcall.go` for LogCall fields.

After implementation, delegate to the **verifier** agent.
If verifier finds issues, fix and re-verify.
