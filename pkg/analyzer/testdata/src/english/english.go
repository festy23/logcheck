package english

import "log/slog"

func good() {
	slog.Info("connection established")
	slog.Info("3 retries left")
	slog.Info("")
	slog.Info("request took 42ms")
	slog.Info("file: /tmp/data.json")
	slog.Info("100% complete")
}

func bad() {
	slog.Info("соединение установлено")    // want `logcheck: english: message contains non-English characters`
	slog.Info("接続が確立されました")       // want `logcheck: english: message contains non-English characters`
	slog.Info("café is ready")             // want `logcheck: english: message contains non-English characters`
	slog.Info("connexion établie")         // want `logcheck: english: message contains non-English characters`
	slog.Error("اتصال برقرار شد")         // want `logcheck: english: message contains non-English characters`
}

func edgeCases() {
	msg := getMessage()
	slog.Info(msg)

	const cyrillic = "ошибка сервера"
	slog.Info(cyrillic) // want `logcheck: english: message contains non-English characters`

	slog.Info("hello " + "мир")   // want `logcheck: english: message contains non-English characters`
	slog.Info("hello " + "world")
}

func getMessage() string { return "dynamic" }
