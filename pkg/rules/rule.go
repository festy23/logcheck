package rules

import (
	"github.com/festy23/logcheck/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// Rule определяет одну проверку линтера для сообщений логирования.
type Rule interface {
	Name() string
	Description() string
	Check(call *model.LogCall, pass *analysis.Pass)
}
