package extractor

import (
	"go/ast"
	"go/types"

	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// extractZapFieldKeys извлекает ключи из вызовов конструкторов zap.Field
// (например, zap.String("key", val), zap.Int("count", 42)).
func extractZapFieldKeys(args []ast.Expr, pass *analysis.Pass) []model.KeyValue {
	var kvs []model.KeyValue
	for _, arg := range args {
		call, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}
		if !isZapFieldConstructor(call, pass) {
			continue
		}
		if len(call.Args) == 0 {
			continue
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			continue
		}
		key, ok := unquoteStringLit(lit)
		if !ok {
			continue
		}
		kvs = append(kvs, model.KeyValue{Key: key, KeyPos: lit.Pos()})
	}
	return kvs
}

// isZapFieldConstructor проверяет, является ли выражение вызовом конструктора
// поля zap (zap.String, zap.Int и т.д.) из пакета "go.uber.org/zap".
func isZapFieldConstructor(call *ast.CallExpr, pass *analysis.Pass) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	obj := pass.TypesInfo.Uses[ident]
	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}
	return pkg.Imported().Path() == "go.uber.org/zap"
}
