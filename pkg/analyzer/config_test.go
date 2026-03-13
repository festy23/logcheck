package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		cfg := loadConfig("")
		if len(cfg.Disable) != 0 || len(cfg.SensitivePatterns) != 0 {
			t.Error("expected empty config for empty path")
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		cfg := loadConfig("/nonexistent/path.json")
		if len(cfg.Disable) != 0 || len(cfg.SensitivePatterns) != 0 {
			t.Error("expected empty config for nonexistent file")
		}
	})

	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".logcheck.json")
		data := `{"disable": ["english", "specialchars"], "sensitive_patterns": ["session_key"]}`
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg := loadConfig(path)
		if len(cfg.Disable) != 2 || cfg.Disable[0] != "english" || cfg.Disable[1] != "specialchars" {
			t.Errorf("Disable = %v, want [english specialchars]", cfg.Disable)
		}
		if len(cfg.SensitivePatterns) != 1 || cfg.SensitivePatterns[0] != "session_key" {
			t.Errorf("SensitivePatterns = %v, want [session_key]", cfg.SensitivePatterns)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.json")
		if err := os.WriteFile(path, []byte("{invalid}"), 0o644); err != nil {
			t.Fatal(err)
		}
		cfg := loadConfig(path)
		if len(cfg.Disable) != 0 {
			t.Error("expected empty config for invalid JSON")
		}
	})
}
