package rules

import (
	"testing"

	"github.com/festy23/logcheck/pkg/model"
)

func TestSpecialcharsRule_Name(t *testing.T) {
	r := &specialcharsRule{}
	if r.Name() != "specialchars" {
		t.Errorf("got %q, want %q", r.Name(), "specialchars")
	}
}

func TestSpecialcharsRule_Description(t *testing.T) {
	r := &specialcharsRule{}
	if r.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestSpecialcharsRule_SkipNonLiteral(t *testing.T) {
	r := &specialcharsRule{}
	call := &model.LogCall{
		HasLiteral: false,
		MsgLiteral: "hello 🚀",
	}
	// pass равен nil — Check должен вернуться до его использования.
	r.Check(call, nil)
}

func TestSpecialcharsRule_AcceptClean(t *testing.T) {
	r := &specialcharsRule{}
	tests := []string{
		"hello world",
		"file: /tmp/data.json",
		"100% complete",
		"key=value",
		"hello world!",
		"user@example.com",
		"",
		"single dot.",
		"two dots..",
	}
	for _, msg := range tests {
		call := &model.LogCall{
			HasLiteral: true,
			MsgLiteral: msg,
		}
		// Чистые сообщения — нет отчёта → pass может быть nil.
		r.Check(call, nil)
	}
}

func TestHasDecorativePunctuation(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"hello world", false},
		{"hello world!", false},
		{"connection failed!!!", true},
		{"what!!", true},
		{"something went wrong...", true},
		{"really???", true},
		{"file.txt", false},
		{"two dots..", false},
		{"a.b.c", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := hasDecorativePunctuation(tt.s); got != tt.want {
				t.Errorf("hasDecorativePunctuation(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsSpecialChar(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'a', false},
		{'Z', false},
		{'0', false},
		{'!', false},
		{' ', false},
		{'.', false},
		{'\t', true},
		{'\n', true},
		{'\r', true},
		{'\x00', true},
		{'🚀', true},
		{'✅', true},
		{'❌', true},
		{'★', true},
		{'©', true},
		{'→', true},
		{'ä', false},  // не-ASCII буква — обрабатывается правилом english
		{'ü', false},  // не-ASCII буква — обрабатывается правилом english
		{'я', false},  // не-ASCII кириллическая буква
		{'5', false},
	}
	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			if got := isSpecialChar(tt.r); got != tt.want {
				t.Errorf("isSpecialChar(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}
