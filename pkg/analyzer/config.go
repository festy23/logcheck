package analyzer

import (
	"encoding/json"
	"os"
)

// Config представляет конфигурацию logcheck из JSON-файла.
type Config struct {
	Disable           []string `json:"disable"`
	SensitivePatterns []string `json:"sensitive_patterns"`
}

// loadConfig читает конфигурацию из JSON-файла.
// Возвращает пустую конфигурацию, если путь пуст или файл не найден.
func loadConfig(path string) Config {
	if path == "" {
		return Config{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}
	}
	return cfg
}
