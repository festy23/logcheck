package rules

import (
	"go/token"
	"os"
	"strconv"
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

	hasSpecial := false
	for _, ch := range call.MsgLiteral {
		if isSpecialChar(ch) {
			hasSpecial = true
			break
		}
	}

	if !hasSpecial && !hasDecorativePunctuation(call.MsgLiteral) {
		return
	}

	diag := analysis.Diagnostic{
		Pos:      call.MsgPos,
		Message:  "logcheck: specialchars: message contains special characters",
		Category: "specialchars",
	}

	if len(call.ConcatParts) == 1 && call.ConcatParts[0].IsLiteral {
		if fix, ok := buildSpecialCharsFix(pass, call.ConcatParts[0].Pos, call.MsgLiteral); ok {
			diag.SuggestedFixes = []analysis.SuggestedFix{fix}
		}
	}

	pass.Report(diag)
}

// buildSpecialCharsFix создаёт SuggestedFix, заменяющий строковый литерал на очищенную версию.
func buildSpecialCharsFix(pass *analysis.Pass, pos token.Pos, msg string) (analysis.SuggestedFix, bool) {
	f := pass.Fset.File(pos)
	if f == nil {
		return analysis.SuggestedFix{}, false
	}
	offset := f.Offset(pos)
	src, err := os.ReadFile(f.Name())
	if err != nil || offset >= len(src) {
		return analysis.SuggestedFix{}, false
	}
	quote := src[offset]
	if quote != '"' && quote != '`' {
		return analysis.SuggestedFix{}, false
	}
	endOffset := findLiteralEnd(src, offset)
	if endOffset < 0 {
		return analysis.SuggestedFix{}, false
	}

	cleaned := cleanMessage(msg)
	var newLit string
	if quote == '`' {
		newLit = "`" + cleaned + "`"
	} else {
		newLit = strconv.Quote(cleaned)
	}

	return analysis.SuggestedFix{
		Message: "remove special characters",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     pos,
				End:     pos + token.Pos(endOffset-offset+1),
				NewText: []byte(newLit),
			},
		},
	}, true
}

// findLiteralEnd находит смещение закрывающей кавычки строкового литерала.
func findLiteralEnd(src []byte, start int) int {
	if start >= len(src) {
		return -1
	}
	quote := src[start]
	if quote == '`' {
		for i := start + 1; i < len(src); i++ {
			if src[i] == '`' {
				return i
			}
		}
		return -1
	}
	if quote == '"' {
		for i := start + 1; i < len(src); i++ {
			if src[i] == '\\' {
				i++
				continue
			}
			if src[i] == '"' {
				return i
			}
		}
		return -1
	}
	return -1
}

// cleanMessage удаляет спецсимволы, эмодзи и декоративную пунктуацию из сообщения.
func cleanMessage(s string) string {
	var buf strings.Builder
	for _, ch := range s {
		if isSpecialChar(ch) {
			if unicode.IsControl(ch) {
				buf.WriteRune(' ')
			}
			continue
		}
		buf.WriteRune(ch)
	}
	result := buf.String()
	result = removeDecorativePunctuation(result)
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	return strings.TrimSpace(result)
}

// removeDecorativePunctuation удаляет повторяющуюся пунктуацию: !! → пусто, ?? → пусто, ... → пусто.
func removeDecorativePunctuation(s string) string {
	runes := []rune(s)
	var result strings.Builder
	i := 0
	for i < len(runes) {
		ch := runes[i]
		if ch == '!' || ch == '?' {
			j := i + 1
			for j < len(runes) && runes[j] == ch {
				j++
			}
			if j-i >= 2 {
				i = j
				continue
			}
		} else if ch == '.' {
			j := i + 1
			for j < len(runes) && runes[j] == '.' {
				j++
			}
			if j-i >= 3 {
				i = j
				continue
			}
			for k := i; k < j; k++ {
				result.WriteRune('.')
			}
			i = j
			continue
		}
		result.WriteRune(ch)
		i++
	}
	return result.String()
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
