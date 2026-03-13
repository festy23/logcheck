package rules

import (
	"go/token"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
)

type lowercaseRule struct{}

func (r *lowercaseRule) Name() string        { return "lowercase" }
func (r *lowercaseRule) Description() string { return "log messages must start with a lowercase letter" }

func (r *lowercaseRule) Check(call *model.LogCall, pass *analysis.Pass) {
	if !call.HasLiteral {
		return
	}

	first, size := utf8.DecodeRuneInString(call.MsgLiteral)
	if first == utf8.RuneError || size == 0 {
		return
	}
	if !unicode.IsUpper(first) {
		return
	}

	diag := analysis.Diagnostic{
		Pos:      call.MsgPos,
		Message:  "logcheck: lowercase: message must start with a lowercase letter",
		Category: "lowercase",
	}

	// Предлагаем автоисправление только для одиночных строковых литералов (не констант и не конкатенаций).
	if len(call.ConcatParts) == 1 && call.ConcatParts[0].IsLiteral {
		if fix, ok := buildFix(pass, call.ConcatParts[0].Pos, first); ok {
			diag.SuggestedFixes = []analysis.SuggestedFix{fix}
		}
	}

	pass.Report(diag)
}

// buildFix создаёт SuggestedFix, заменяющий первую букву строкового литерала на строчную.
// Проверяет, что байт исходного кода в позиции pos является символом кавычки (отличая BasicLit от Ident).
func buildFix(pass *analysis.Pass, pos token.Pos, upper rune) (analysis.SuggestedFix, bool) {
	f := pass.Fset.File(pos)
	if f == nil {
		return analysis.SuggestedFix{}, false
	}
	offset := f.Offset(pos)
	src, err := os.ReadFile(f.Name())
	if err != nil || offset >= len(src) {
		return analysis.SuggestedFix{}, false
	}
	if src[offset] != '"' && src[offset] != '`' {
		return analysis.SuggestedFix{}, false
	}

	lower := unicode.ToLower(upper)
	upperBytes := make([]byte, utf8.UTFMax)
	upperSize := utf8.EncodeRune(upperBytes, upper)
	lowerBytes := make([]byte, utf8.UTFMax)
	lowerSize := utf8.EncodeRune(lowerBytes, lower)

	return analysis.SuggestedFix{
		Message: "lowercase the first letter",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     pos + 1,                       // пропускаем открывающую кавычку
				End:     pos + 1 + token.Pos(upperSize), // охватываем заглавную руну
				NewText: lowerBytes[:lowerSize],
			},
		},
	}, true
}
