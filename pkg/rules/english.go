package rules

import (
	"unicode"

	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
)

type englishRule struct{}

func (r *englishRule) Name() string        { return "english" }
func (r *englishRule) Description() string { return "log messages must contain only English (ASCII) text" }

func (r *englishRule) Check(call *model.LogCall, pass *analysis.Pass) {
	if !call.HasLiteral {
		return
	}

	for _, ch := range call.MsgLiteral {
		if unicode.IsLetter(ch) && ch > 127 {
			pass.Report(analysis.Diagnostic{
				Pos:      call.MsgPos,
				Message:  "logcheck: english: message contains non-English characters",
				Category: "english",
			})
			return
		}
	}
}
