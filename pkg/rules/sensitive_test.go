package rules

import (
	"testing"

	"github.com/festy23/loglinter/pkg/model"
)

func TestSensitiveRule_Name(t *testing.T) {
	r := &sensitiveRule{}
	if r.Name() != "sensitive" {
		t.Errorf("got %q, want %q", r.Name(), "sensitive")
	}
}

func TestSensitiveRule_Description(t *testing.T) {
	r := &sensitiveRule{}
	if r.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestSensitiveRule_SkipEmpty(t *testing.T) {
	r := &sensitiveRule{}
	call := &model.LogCall{
		HasLiteral: false,
	}
	// pass равен nil — Check должен вернуться до его использования.
	r.Check(call, nil)
}

func TestSensitiveRule_KeyDetection(t *testing.T) {
	r := &sensitiveRule{}
	tests := []struct {
		key  string
		want bool
	}{
		{"password", true},
		{"Password", true},
		{"user_password", true},
		{"api_key", true},
		{"API-KEY", true},
		{"auth_token", true},
		{"access_token", true},
		{"ssn", true},
		{"credit_card", true},
		{"private-key", true},
		{"authorization", true},
		{"credentials", true},
		{"refresh_token", true},
		{"user_id", false},
		{"name", false},
		{"count", false},
		{"request_id", false},
		{"method", false},
		{"status", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := r.isSensitiveKey(tt.key); got != tt.want {
				t.Errorf("isSensitiveKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestSensitiveRule_PatternDetection(t *testing.T) {
	r := &sensitiveRule{}
	tests := []struct {
		msg  string
		want bool
	}{
		{"password=abc123", true},
		{"token: xyz", true},
		{"secret =hidden", true},
		{"api_key=foo", true},
		{"ssn:123-45-6789", true},
		{"password =abc", true},
		{"private_key=xxx", true},
		{"credential:abc", true},
		{"credit_card=4111", true},
		{"reset password", false},
		{"hello world", false},
		{"connection established", false},
		{"token count exceeded", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			if got := r.hasSensitivePattern(tt.msg); got != tt.want {
				t.Errorf("hasSensitivePattern(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestSensitiveRule_ExtraPatterns(t *testing.T) {
	r := &sensitiveRule{}
	r.addPatterns([]string{"custom_secret", "mytoken"})

	if !r.isSensitiveKey("custom_secret") {
		t.Error("extra pattern custom_secret not detected as key")
	}
	if !r.isSensitiveKey("MY-TOKEN") {
		t.Error("extra pattern mytoken not detected as key (normalized)")
	}
	if r.isSensitiveKey("username") {
		t.Error("username should not be sensitive")
	}

	if !r.hasSensitivePattern("custom_secret=abc") {
		t.Error("extra keyword custom_secret not detected in message")
	}
	if !r.hasSensitivePattern("mytoken: xyz") {
		t.Error("extra keyword mytoken not detected in message")
	}
}

func TestSensitiveRule_AddPatternsViaRegistry(t *testing.T) {
	reg := NewRegistry()
	reg.AddSensitivePatterns([]string{"internal_id"})

	// Проверяем через прямое приведение типа для тестирования.
	for _, rule := range reg.Rules() {
		if sr, ok := rule.(*sensitiveRule); ok {
			if !sr.isSensitiveKey("internal_id") {
				t.Error("AddSensitivePatterns did not add pattern")
			}
			return
		}
	}
	t.Error("sensitive rule not found in registry")
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"api_key", "apikey"},
		{"API-KEY", "apikey"},
		{"auth_token", "authtoken"},
		{"Password", "password"},
		{"user-name", "username"},
		{"simple", "simple"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeKey(tt.input); got != tt.want {
				t.Errorf("normalizeKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
