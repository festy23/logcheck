package rules

import (
	"testing"

	"github.com/festy23/loglinter/pkg/model"
)

func TestEnglishRule_Name(t *testing.T) {
	r := &englishRule{}
	if r.Name() != "english" {
		t.Errorf("got %q, want %q", r.Name(), "english")
	}
}

func TestEnglishRule_Description(t *testing.T) {
	r := &englishRule{}
	if r.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestEnglishRule_SkipNonLiteral(t *testing.T) {
	r := &englishRule{}
	call := &model.LogCall{
		HasLiteral: false,
		MsgLiteral: "café",
	}
	// pass равен nil — Check должен вернуться до его использования.
	r.Check(call, nil)
}

func TestEnglishRule_AcceptASCII(t *testing.T) {
	r := &englishRule{}
	tests := []string{
		"hello world",
		"request took 42ms",
		"file: /tmp/data.json",
		"100% complete",
		"",
		"key=value",
	}
	for _, msg := range tests {
		call := &model.LogCall{
			HasLiteral: true,
			MsgLiteral: msg,
		}
		// Все символы ASCII — нет отчёта → pass может быть nil.
		r.Check(call, nil)
	}
}
