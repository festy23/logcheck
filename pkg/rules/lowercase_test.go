package rules

import (
	"testing"

	"github.com/festy23/logcheck/pkg/model"
)

func TestLowercaseRule_Name(t *testing.T) {
	r := &lowercaseRule{}
	if r.Name() != "lowercase" {
		t.Errorf("got %q, want %q", r.Name(), "lowercase")
	}
}

func TestLowercaseRule_Description(t *testing.T) {
	r := &lowercaseRule{}
	if r.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestLowercaseRule_SkipNonLiteral(t *testing.T) {
	r := &lowercaseRule{}
	call := &model.LogCall{
		HasLiteral: false,
		MsgLiteral: "Hello",
	}
	// pass равен nil — Check должен вернуться до его использования.
	r.Check(call, nil)
}

func TestLowercaseRule_SkipEmptyString(t *testing.T) {
	r := &lowercaseRule{}
	call := &model.LogCall{
		HasLiteral: true,
		MsgLiteral: "",
	}
	// Пустая строка: DecodeRuneInString возвращает RuneError, size 0 → ранний возврат.
	r.Check(call, nil)
}

func TestLowercaseRule_SkipNonLetter(t *testing.T) {
	r := &lowercaseRule{}
	tests := []string{
		"3 retries left",
		"_internal thing",
		"[debug] something",
		"123",
	}
	for _, msg := range tests {
		call := &model.LogCall{
			HasLiteral: true,
			MsgLiteral: msg,
		}
		// Первый символ не заглавная буква → нет отчёта → pass может быть nil.
		r.Check(call, nil)
	}
}
