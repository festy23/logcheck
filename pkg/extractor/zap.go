package extractor

import (
	"go/ast"
	"strings"

	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// zapLoggerMethods — методы логирования *zap.Logger.
var zapLoggerMethods = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"DPanic": true, "Fatal": true, "Panic": true,
}

// zapSugarBaseMethods — базовые методы логирования *zap.SugaredLogger
// (без суффикса — аргументы конкатенируются, без KV).
var zapSugarBaseMethods = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"DPanic": true, "Fatal": true, "Panic": true,
}

// zapIgnoreMethods — методы, которые не создают сообщений логирования.
var zapIgnoreMethods = map[string]bool{
	"With": true, "Named": true, "Sugar": true, "Desugar": true,
	"Sync": true, "Core": true, "Check": true, "WithOptions": true,
}

// extractZap извлекает LogCall из выражения вызова go.uber.org/zap.
// typeName различает "Logger" и "SugaredLogger".
// Возвращает nil, если метод не является известным методом логирования.
func extractZap(call *ast.CallExpr, method string, typeName string, pass *analysis.Pass) *model.LogCall {
	if zapIgnoreMethods[method] {
		return nil
	}

	switch typeName {
	case "Logger":
		return extractZapLogger(call, method, pass)
	case "SugaredLogger":
		return extractZapSugar(call, method, pass)
	default:
		return nil
	}
}

// extractZapLogger обрабатывает методы *zap.Logger.
func extractZapLogger(call *ast.CallExpr, method string, pass *analysis.Pass) *model.LogCall {
	if !zapLoggerMethods[method] {
		return nil
	}
	if len(call.Args) == 0 {
		return nil
	}

	literal, parts, hasLiteral, msgPos := resolveMessage(call.Args[0], pass)

	lc := &model.LogCall{
		Pos:         call.Pos(),
		Logger:      model.LoggerZap,
		Method:      method,
		MsgLiteral:  literal,
		HasLiteral:  hasLiteral,
		MsgPos:      msgPos,
		ConcatParts: parts,
	}

	// Извлекаем ключи из конструкторов zap.Field (аргументы после msg).
	if len(call.Args) > 1 {
		lc.KeyValues = extractZapFieldKeys(call.Args[1:], pass)
	}

	return lc
}

// extractZapSugar обрабатывает методы *zap.SugaredLogger.
func extractZapSugar(call *ast.CallExpr, method string, pass *analysis.Pass) *model.LogCall {
	if len(call.Args) == 0 {
		return nil
	}

	// Определяем категорию метода по суффиксу.
	isW := strings.HasSuffix(method, "w")
	isF := strings.HasSuffix(method, "f")

	// Проверяем, что это известный шаблон метода.
	baseName := method
	if isW || isF {
		baseName = method[:len(method)-1]
	}
	if !zapSugarBaseMethods[baseName] {
		return nil
	}

	literal, parts, hasLiteral, msgPos := resolveMessage(call.Args[0], pass)

	lc := &model.LogCall{
		Pos:         call.Pos(),
		Logger:      model.LoggerZapSugar,
		Method:      method,
		MsgLiteral:  literal,
		HasLiteral:  hasLiteral,
		MsgPos:      msgPos,
		ConcatParts: parts,
	}

	// Методы *w используют чередующиеся пары ключ-значение после сообщения.
	if isW && len(call.Args) > 1 {
		lc.KeyValues = extractAlternatingKV(call.Args, 1, pass)
	}

	return lc
}
