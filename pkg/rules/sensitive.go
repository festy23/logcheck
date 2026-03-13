package rules

import (
	"strings"

	"github.com/festy23/loglinter/pkg/model"
	"golang.org/x/tools/go/analysis"
)

type sensitiveRule struct {
	extraTerms    []string // дополнительные нормализованные термины чувствительных ключей
	extraKeywords []string // дополнительные ключевые слова для сообщений
}

func (r *sensitiveRule) Name() string        { return "sensitive" }
func (r *sensitiveRule) Description() string { return "log messages must not contain sensitive data" }

// sensitiveTerms — нормализованные (строчные буквы, без подчёркиваний/дефисов) шаблоны чувствительных ключей.
var sensitiveTerms = []string{
	"password", "passwd", "pwd",
	"secret",
	"token", "accesstoken", "refreshtoken", "authtoken",
	"apikey",
	"privatekey",
	"credential", "credentials",
	"authorization",
	"ssn",
	"creditcard", "cardnumber",
}

// sensitiveMessageKeywords — ключевые слова для поиска в тексте сообщения, за которыми следует = или :.
var sensitiveMessageKeywords = []string{
	"password", "passwd", "pwd",
	"secret",
	"token",
	"apikey", "api_key",
	"privatekey", "private_key",
	"credential",
	"authorization",
	"ssn",
	"creditcard", "credit_card",
}

// addPatterns добавляет дополнительные шаблоны чувствительных данных (как термины ключей, так и ключевые слова сообщений).
func (r *sensitiveRule) addPatterns(patterns []string) {
	for _, p := range patterns {
		r.extraTerms = append(r.extraTerms, normalizeKey(p))
		r.extraKeywords = append(r.extraKeywords, strings.ToLower(p))
	}
}

func (r *sensitiveRule) Check(call *model.LogCall, pass *analysis.Pass) {
	// Проверяем ключи структурированного логирования.
	for _, kv := range call.KeyValues {
		if r.isSensitiveKey(kv.Key) {
			pass.Report(analysis.Diagnostic{
				Pos:      kv.KeyPos,
				Message:  "logcheck: sensitive: key \"" + kv.Key + "\" may contain sensitive data",
				Category: "sensitive",
			})
		}
	}

	// Проверяем литеральные части сообщения на встроенные шаблоны учётных данных.
	for _, part := range call.ConcatParts {
		if part.IsLiteral && r.hasSensitivePattern(part.Value) {
			pass.Report(analysis.Diagnostic{
				Pos:      call.MsgPos,
				Message:  "logcheck: sensitive: message may contain embedded credentials",
				Category: "sensitive",
			})
			return // одна диагностика на сообщение достаточно
		}
	}
}

// normalizeKey приводит к нижнему регистру и удаляет подчёркивания и дефисы.
func normalizeKey(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// isSensitiveKey проверяет, совпадает ли имя ключа структурированного логирования с чувствительным термином.
func (r *sensitiveRule) isSensitiveKey(key string) bool {
	normalized := normalizeKey(key)
	for _, term := range sensitiveTerms {
		if strings.Contains(normalized, term) {
			return true
		}
	}
	for _, term := range r.extraTerms {
		if strings.Contains(normalized, term) {
			return true
		}
	}
	return false
}

// hasSensitivePattern проверяет, содержит ли строка чувствительное ключевое слово, за которым следует = или :.
func (r *sensitiveRule) hasSensitivePattern(s string) bool {
	lower := strings.ToLower(s)
	for _, kw := range sensitiveMessageKeywords {
		if matchKeywordPattern(lower, kw) {
			return true
		}
	}
	for _, kw := range r.extraKeywords {
		if matchKeywordPattern(lower, kw) {
			return true
		}
	}
	return false
}

// matchKeywordPattern проверяет, содержит ли lower ключевое слово kw, за которым следует = или :.
func matchKeywordPattern(lower, kw string) bool {
	idx := strings.Index(lower, kw)
	if idx < 0 {
		return false
	}
	rest := lower[idx+len(kw):]
	rest = strings.TrimLeft(rest, " ")
	return len(rest) > 0 && (rest[0] == '=' || rest[0] == ':')
}
