package rules

import (
	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
)

// Registry хранит набор включённых правил.
type Registry struct {
	rules []Rule
}

// NewRegistry создаёт реестр со всеми правилами по умолчанию.
func NewRegistry() *Registry {
	return &Registry{
		rules: defaultRules(),
	}
}

// defaultRules возвращает встроенные правила.
func defaultRules() []Rule {
	return []Rule{
		&lowercaseRule{},
		&englishRule{},
		&specialcharsRule{},
	}
}

// RunAll выполняет все зарегистрированные правила для данного вызова.
func (r *Registry) RunAll(call *model.LogCall, pass *analysis.Pass) {
	for _, rule := range r.rules {
		rule.Check(call, pass)
	}
}

// Disable удаляет правило по имени. Возвращает true, если найдено и удалено.
func (r *Registry) Disable(name string) bool {
	for i, rule := range r.rules {
		if rule.Name() == name {
			r.rules = append(r.rules[:i], r.rules[i+1:]...)
			return true
		}
	}
	return false
}

// Rules возвращает копию текущего списка правил.
func (r *Registry) Rules() []Rule {
	return append([]Rule(nil), r.rules...)
}
