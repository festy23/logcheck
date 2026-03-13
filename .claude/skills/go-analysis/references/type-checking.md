# TypesInfo Deep Dive

## pass.TypesInfo.Uses

Maps `*ast.Ident` → `types.Object` for every identifier USE (not definition).

Key Object types:
- `*types.PkgName` — import reference (slog, zap)
- `*types.Func` — function/method reference
- `*types.Var` — variable reference
- `*types.Const` — constant reference
- `*types.TypeName` — type reference

## pass.TypesInfo.Selections

Maps `*ast.SelectorExpr` → `*types.Selection` for method/field access.
Only populated for EXPRESSIONS (method calls on values), NOT for
qualified identifiers (package-level calls like slog.Info).

```go
sel := call.Fun.(*ast.SelectorExpr)
selection := pass.TypesInfo.Selections[sel]
if selection != nil {
    // Method call on an instance — get receiver type
    recvType := selection.Recv()
    named := namedType(recvType)
    if named != nil {
        pkg := named.Obj().Pkg()
        if pkg != nil {
            pkgPath := pkg.Path()          // "log/slog" or "go.uber.org/zap"
            typeName := named.Obj().Name() // "Logger", "SugaredLogger"
        }
    }
}
```

Helper to strip pointer and get Named type:
```go
func namedType(t types.Type) *types.Named {
    t = types.Unalias(t)
    if ptr, ok := t.(*types.Pointer); ok {
        t = ptr.Elem()
    }
    named, _ := t.(*types.Named)
    return named
}
```

`selection.Kind()`:
- `types.MethodVal` — method call (`logger.Info`)
- `types.MethodExpr` — method expression (`(*Logger).Info`)
- `types.FieldVal` — field access (`config.Port`)

## pass.TypesInfo.Types

Maps `ast.Expr` → `types.TypeAndValue`. Use for:
- Getting type of any expression: `pass.TypesInfo.Types[expr].Type`
- Checking if expression is a constant: `tv.Value != nil`

## Resolving constants

```go
ident := expr.(*ast.Ident)
obj := pass.TypesInfo.Uses[ident]
if c, ok := obj.(*types.Const); ok {
    val := constant.StringVal(c.Val()) // extract string value
}
```

## Package path comparison

ALWAYS compare by full import path:
- `"log/slog"` not `"slog"`
- `"go.uber.org/zap"` not `"zap"`
- `"fmt"` (happens to match, but use the full path for consistency)

## Distinguishing Logger vs SugaredLogger (zap)

Both live in `"go.uber.org/zap"`. After getting the Named type:
- `named.Obj().Name() == "Logger"` → structured API, zap.Field args
- `named.Obj().Name() == "SugaredLogger"` → printf/kv-style args

## Chained calls: logger.With(...).Info("msg")

With() returns same type (*Logger or *SugaredLogger).
TypesInfo.Selections resolves .Info() receiver through the chain automatically.
No special handling needed — the type checker does the work.
