package extractor

import (
	"go/ast"
	"go/types"

	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// extractAlternatingKV извлекает пары ключ-значение из чередующихся списков
// аргументов string/any, начиная с указанного индекса. Используется slog и
// методами *w zap SugaredLogger.
func extractAlternatingKV(args []ast.Expr, startIdx int, pass *analysis.Pass) []model.KeyValue {
	var kvs []model.KeyValue
	for i := startIdx; i < len(args); i += 2 {
		arg := args[i]
		// Пропускаем аргументы типа slog.Attr — они не являются парами со строковым ключом.
		if isSlogAttr(arg, pass) {
			// slog.Attr занимает один слот, а не два.
			i--
			continue
		}
		lit, ok := arg.(*ast.BasicLit)
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

// isSlogAttr проверяет, имеет ли выражение тип slog.Attr.
func isSlogAttr(expr ast.Expr, pass *analysis.Pass) bool {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok {
		return false
	}
	named := namedType(tv.Type)
	if named == nil {
		return false
	}
	obj := named.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == "log/slog" && obj.Name() == "Attr"
}

// namedType убирает алиасы и указатели, возвращая базовый *types.Named,
// или nil, если тип не является именованным.
func namedType(t types.Type) *types.Named {
	t = types.Unalias(t)
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, _ := t.(*types.Named)
	return named
}
