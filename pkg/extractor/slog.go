package extractor

import (
	"go/ast"

	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// slogMsgIndex сопоставляет имена методов slog с индексом аргумента сообщения.
var slogMsgIndex = map[string]int{
	"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
	"DebugContext": 1, "InfoContext": 1, "WarnContext": 1, "ErrorContext": 1,
	"Log": 2, "LogAttrs": 2,
}

// extractSlog извлекает LogCall из выражения вызова log/slog.
// Возвращает nil, если метод не является известным методом логирования.
func extractSlog(call *ast.CallExpr, method string, pass *analysis.Pass) *model.LogCall {
	msgIdx, ok := slogMsgIndex[method]
	if !ok {
		return nil
	}
	if msgIdx >= len(call.Args) {
		return nil
	}

	literal, parts, hasLiteral, msgPos := resolveMessage(call.Args[msgIdx], pass)

	lc := &model.LogCall{
		Pos:         call.Pos(),
		Logger:      model.LoggerSlog,
		Method:      method,
		MsgLiteral:  literal,
		HasLiteral:  hasLiteral,
		MsgPos:      msgPos,
		ConcatParts: parts,
	}

	// Извлекаем пары ключ-значение из аргументов после сообщения.
	kvStart := msgIdx + 1
	if kvStart < len(call.Args) {
		lc.KeyValues = extractAlternatingKV(call.Args, kvStart, pass)
	}

	return lc
}
