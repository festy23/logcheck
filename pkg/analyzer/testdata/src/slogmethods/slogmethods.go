package slogmethods

import (
	"context"
	"fmt"
	"log/slog"
)

// --- дополнительные методы slog ---

func slogDebugWarn() {
	slog.Debug("processing item")
	slog.Debug("Processing item") // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.Warn("connection slow")
	slog.Warn("Connection slow") // want `logcheck: lowercase: message must start with a lowercase letter`
}

func slogContextMethods() {
	ctx := context.Background()
	slog.InfoContext(ctx, "request started")
	slog.InfoContext(ctx, "Request started") // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.WarnContext(ctx, "something happened")
	slog.WarnContext(ctx, "Something happened") // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.ErrorContext(ctx, "request failed")
	slog.DebugContext(ctx, "processing item")
}

func slogLog() {
	ctx := context.Background()
	slog.Log(ctx, slog.LevelInfo, "custom log entry")
	slog.Log(ctx, slog.LevelInfo, "Custom log entry") // want `logcheck: lowercase: message must start with a lowercase letter`
}

// --- сырые строки (обратные кавычки) ---

func rawStrings() {
	slog.Info(`hello raw string`)
	slog.Info(`Hello raw string`) // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.Info(`ошибка в запросе`) // want `logcheck: english: message contains non-English characters`
	slog.Info(`launch 🚀`)       // want `logcheck: specialchars: message contains special characters`
	slog.Info(`password=secret`)  // want `logcheck: sensitive: message may contain embedded credentials`
}

// --- fmt.Sprintf ---
func fmtSprintf() {
	slog.Info(fmt.Sprintf("request %d completed", 42))
	slog.Info(fmt.Sprintf("Request %d failed", 42))    // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.Info(fmt.Sprintf("ошибка %d", 42))             // want `logcheck: english: message contains non-English characters`
	slog.Info(fmt.Sprintf("done 🚀 %d", 42))            // want `logcheck: specialchars: message contains special characters`
	slog.Info(fmt.Sprintf("password=%s", "hunter2"))    // want `logcheck: sensitive: message may contain embedded credentials`
}

// --- константы ---

func constMessages() {
	const good = "server started"
	slog.Info(good)

	const bad = "Server started"
	slog.Info(bad) // want `logcheck: lowercase: message must start with a lowercase letter`

	const cyrillic = "ошибка"
	slog.Info(cyrillic) // want `logcheck: english: message contains non-English characters`
}

// --- конкатенация ---

func concatMessages() {
	slog.Info("hello " + "world")
	slog.Info("Hello " + "world") // want `logcheck: lowercase: message must start with a lowercase letter`
	slog.Info("hello " + "мир")   // want `logcheck: english: message contains non-English characters`
}

// --- переменная (динамическая) — должна быть пропущена ---

func dynamicMessages() {
	msg := getMessage()
	slog.Info(msg)
	slog.Error(msg)
	slog.Debug(msg)
	slog.Warn(msg)
}

// --- пустая строка ---

func emptyString() {
	slog.Info("")
	slog.Debug("")
}

func getMessage() string { return "dynamic" }
