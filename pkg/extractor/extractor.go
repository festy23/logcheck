package extractor

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Extract обходит AST указанного пакета и возвращает все обнаруженные
// вызовы логирования из log/slog и go.uber.org/zap.
func Extract(pass *analysis.Pass) []model.LogCall {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Собираем множество сгенерированных файлов для пропуска.
	generated := generatedFiles(pass)

	var calls []model.LogCall

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Обрабатываем только выражения-селекторы (pkg.Func или obj.Method).
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// Пропускаем сгенерированные файлы.
		pos := pass.Fset.Position(call.Pos())
		if generated[pos.Filename] {
			return
		}

		method := sel.Sel.Name

		// Сначала пробуем функцию уровня пакета (например, slog.Info).
		if ident, ok := sel.X.(*ast.Ident); ok {
			obj := pass.TypesInfo.Uses[ident]
			if pkg, ok := obj.(*types.PkgName); ok {
				pkgPath := pkg.Imported().Path()
				lc := extractByPkgPath(call, method, pkgPath, "", pass)
				if lc != nil {
					calls = append(calls, *lc)
				}
				return
			}
		}

		// Пробуем метод на экземпляре (например, logger.Info).
		selection := pass.TypesInfo.Selections[sel]
		if selection == nil {
			return
		}
		named := namedType(selection.Recv())
		if named == nil {
			return
		}
		obj := named.Obj()
		if obj.Pkg() == nil {
			return
		}
		lc := extractByPkgPath(call, method, obj.Pkg().Path(), obj.Name(), pass)
		if lc != nil {
			calls = append(calls, *lc)
		}
	})

	return calls
}

// extractByPkgPath маршрутизирует к соответствующему экстрактору логгера.
// typeName пуст для вызовов уровня пакета и задан для вызовов методов.
func extractByPkgPath(call *ast.CallExpr, method, pkgPath, typeName string, pass *analysis.Pass) *model.LogCall {
	switch pkgPath {
	case "log/slog":
		return extractSlog(call, method, pass)
	case "go.uber.org/zap":
		if typeName == "" {
			// Функция уровня пакета zap — не вызов логирования, который мы отслеживаем.
			return nil
		}
		return extractZap(call, method, typeName, pass)
	default:
		return nil
	}
}

// generatedFiles возвращает множество имён файлов, содержащих комментарий
// "Code generated", которые следует пропускать.
func generatedFiles(pass *analysis.Pass) map[string]bool {
	gen := make(map[string]bool)
	for _, f := range pass.Files {
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				if strings.Contains(c.Text, "Code generated") {
					fname := pass.Fset.Position(f.Pos()).Filename
					gen[fname] = true
				}
			}
		}
	}
	return gen
}
