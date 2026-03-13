package extractor

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// resolveMessage рекурсивно разрешает строковое выражение в его литеральное
// содержимое и составные части. Возвращает объединённую литеральную строку,
// отдельные части конкатенации, признак разрешения в литерал и позицию
// выражения сообщения.
func resolveMessage(expr ast.Expr, pass *analysis.Pass) (literal string, parts []model.ConcatPart, hasLiteral bool, msgPos token.Pos) {
	msgPos = expr.Pos()
	parts = collectParts(expr, pass)
	if len(parts) == 0 {
		return "", nil, false, msgPos
	}
	// Проверяем, что все части являются литералами.
	for _, p := range parts {
		if !p.IsLiteral {
			return "", parts, false, msgPos
		}
	}
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(p.Value)
	}
	return sb.String(), parts, true, msgPos
}

// collectParts рекурсивно обходит выражение и возвращает его строковые части.
func collectParts(expr ast.Expr, pass *analysis.Pass) []model.ConcatPart {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return nil
		}
		s, ok := unquoteStringLit(e)
		if !ok {
			return nil
		}
		return []model.ConcatPart{{Value: s, Pos: e.Pos(), IsLiteral: true}}

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return nil
		}
		left := collectParts(e.X, pass)
		right := collectParts(e.Y, pass)
		if left == nil && right == nil {
			return nil
		}
		// Если одна из сторон nil, всё равно отслеживаем её как нелитеральную часть.
		if left == nil {
			left = []model.ConcatPart{{Pos: e.X.Pos(), IsLiteral: false}}
		}
		if right == nil {
			right = []model.ConcatPart{{Pos: e.Y.Pos(), IsLiteral: false}}
		}
		return append(left, right...)

	case *ast.Ident:
		obj := pass.TypesInfo.Uses[e]
		if c, ok := obj.(*types.Const); ok && c.Val().Kind() == constant.String {
			val := constant.StringVal(c.Val())
			return []model.ConcatPart{{Value: val, Pos: e.Pos(), IsLiteral: true}}
		}
		return []model.ConcatPart{{Pos: e.Pos(), IsLiteral: false}}

	case *ast.CallExpr:
		// Обработка fmt.Sprintf — извлекаем строку формата как частичную информацию.
		if isFmtSprintf(e, pass) && len(e.Args) > 0 {
			return collectParts(e.Args[0], pass)
		}
		return nil

	case *ast.ParenExpr:
		return collectParts(e.X, pass)

	default:
		return nil
	}
}

// isFmtSprintf проверяет, является ли выражение вызовом fmt.Sprintf.
func isFmtSprintf(call *ast.CallExpr, pass *analysis.Pass) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Sprintf" {
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
	return pkg.Imported().Path() == "fmt"
}

// unquoteStringLit извлекает строковое значение из строкового литерала.
// Обрабатывает как интерпретируемые ("..."), так и сырые (`...`) строки.
func unquoteStringLit(lit *ast.BasicLit) (string, bool) {
	if lit.Kind != token.STRING {
		return "", false
	}
	if strings.HasPrefix(lit.Value, "`") {
		return strings.Trim(lit.Value, "`"), true
	}
	s, err := strconv.Unquote(lit.Value)
	return s, err == nil
}
