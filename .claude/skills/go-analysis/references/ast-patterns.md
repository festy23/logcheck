# AST Patterns for Log Call Detection

## Package function: slog.Info("msg")

```
CallExpr.Fun = SelectorExpr {
  X:   Ident{Name: "slog"}     ← TypesInfo.Uses → *types.PkgName
  Sel: Ident{Name: "Info"}
}
```

## Method on instance: logger.Info("msg")

```
CallExpr.Fun = SelectorExpr {
  X:   Ident{Name: "logger"}   ← TypesInfo.Uses → *types.Var
  Sel: Ident{Name: "Info"}
}
```

Use TypesInfo.Selections[sel] to get receiver type for methods.

## Resolving package identity

```go
obj := pass.TypesInfo.Uses[ident]
switch o := obj.(type) {
case *types.PkgName:
    pkgPath := o.Imported().Path() // "log/slog"
case *types.Var:
    // method on variable — resolve receiver type
}
```

NEVER check `ident.Name == "slog"` — breaks with import aliases.

## Resolving method receiver type

```go
selection := pass.TypesInfo.Selections[sel]
if selection != nil {
    recv := selection.Recv()
    if ptr, ok := recv.(*types.Pointer); ok {
        recv = ptr.Elem()
    }
    if named, ok := recv.(*types.Named); ok {
        pkgPath := named.Obj().Pkg().Path() // "go.uber.org/zap"
        typeName := named.Obj().Name()      // "Logger" or "SugaredLogger"
    }
}
```

Note: Selections is populated only for method/field access, NOT for
package-level calls (slog.Info). Those use TypesInfo.Uses directly.

## String literal extraction

```go
func unquoteStringLit(lit *ast.BasicLit) (string, bool) {
    if lit.Kind != token.STRING { return "", false }
    if strings.HasPrefix(lit.Value, "`") {
        return strings.Trim(lit.Value, "`"), true // raw string
    }
    s, err := strconv.Unquote(lit.Value)
    return s, err == nil
}
```

## Concatenation: "hello " + name

Left-associative BinaryExpr with Op == token.ADD.
Walk recursively: left side may be another BinaryExpr.

## Constant resolution

```go
obj := pass.TypesInfo.Uses[ident]
if c, ok := obj.(*types.Const); ok && c.Val().Kind() == constant.String {
    val := constant.StringVal(c.Val())
}
```
