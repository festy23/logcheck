package rules

import (
	"strings"
	"unicode"

	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

type specialcharsRule struct{}

func (r *specialcharsRule) Name() string        { return "specialchars" }
func (r *specialcharsRule) Description() string { return "log messages must not contain special characters or emoji" }

func (r *specialcharsRule) Check(call *model.LogCall, pass *analysis.Pass) {
	if !call.HasLiteral {
		return
	}

	for _, ch := range call.MsgLiteral {
		if isSpecialChar(ch) {
			pass.Report(analysis.Diagnostic{
				Pos:      call.MsgPos,
				Message:  "logcheck: specialchars: message contains special characters",
				Category: "specialchars",
			})
			return
		}
	}

	if hasDecorativePunctuation(call.MsgLiteral) {
		pass.Report(analysis.Diagnostic{
			Pos:      call.MsgPos,
			Message:  "logcheck: specialchars: message contains special characters",
			Category: "specialchars",
		})
	}
}

// hasDecorativePunctuation проверяет наличие повторяющейся ASCII-пунктуации,
// не несущей информации: "!!!", "...", "???" и т.д.
func hasDecorativePunctuation(s string) bool {
	return strings.Contains(s, "!!") || strings.Contains(s, "??") || strings.Contains(s, "...")
}

// isSpecialChar возвращает true, если руна является управляющим символом,
// эмодзи или не-ASCII символом, которому не место в сообщениях логирования.
func isSpecialChar(r rune) bool {
	if unicode.IsControl(r) {
		return true
	}
	// Не-ASCII символы, не являющиеся буквами или цифрами (эмодзи, символы и т.д.).
	// Не-ASCII буквы обрабатываются правилом english.
	if r > 127 && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
		return true
	}
	return false
}
