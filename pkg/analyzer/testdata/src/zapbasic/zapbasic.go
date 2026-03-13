package zapbasic

import "go.uber.org/zap"

// --- правило lowercase ---

func zapLoggerLowercase() {
	logger := zap.NewNop()
	logger.Info("Hello world")      // want `logcheck: lowercase: message must start with a lowercase letter`
	logger.Error("Something broke") // want `logcheck: lowercase: message must start with a lowercase letter`
	logger.Debug("Debug message")   // want `logcheck: lowercase: message must start with a lowercase letter`
	logger.Warn("Warning here")     // want `logcheck: lowercase: message must start with a lowercase letter`
	logger.Info("hello world")
}

func zapSugarLowercase() {
	sugar := zap.NewNop().Sugar()
	sugar.Info("Hello world")          // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Infow("Hello structured")    // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Infof("Hello %s", "world")   // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Error("Error happened")      // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Errorw("Error structured")   // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Debug("Debug started")       // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Warn("Warning issued")       // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Info("hello world")
}

// --- правило english ---

func zapLoggerEnglish() {
	logger := zap.NewNop()
	logger.Info("ошибка сервера") // want `logcheck: english: message contains non-English characters`
	logger.Info("hello world")
}

func zapSugarEnglish() {
	sugar := zap.NewNop().Sugar()
	sugar.Info("ошибка")          // want `logcheck: english: message contains non-English characters`
	sugar.Infow("ошибка запроса") // want `logcheck: english: message contains non-English characters`
	sugar.Info("hello world")
}

// --- правило specialchars ---

func zapLoggerSpecialchars() {
	logger := zap.NewNop()
	logger.Info("server started 🚀") // want `logcheck: specialchars: message contains special characters`
	logger.Info("hello world")
}

func zapSugarSpecialchars() {
	sugar := zap.NewNop().Sugar()
	sugar.Info("launch 🚀")          // want `logcheck: specialchars: message contains special characters`
	sugar.Infow("check passed ✅")   // want `logcheck: specialchars: message contains special characters`
	sugar.Info("hello world")
}

// --- правило sensitive ---

func zapLoggerSensitive() {
	logger := zap.NewNop()
	logger.Info("auth", zap.String("password", "hunter2"))   // want `logcheck: sensitive: key "password" may contain sensitive data`
	logger.Info("config", zap.String("api_key", "ak-xxx"))   // want `logcheck: sensitive: key "api_key" may contain sensitive data`
	logger.Info("tls", zap.String("private_key", "---BEGIN")) // want `logcheck: sensitive: key "private_key" may contain sensitive data`
	logger.Info("password=hunter2")                           // want `logcheck: sensitive: message may contain embedded credentials`
	logger.Info("request", zap.String("method", "GET"))
}

func zapSugarSensitive() {
	sugar := zap.NewNop().Sugar()
	sugar.Infow("auth", "password", "hunter2")    // want `logcheck: sensitive: key "password" may contain sensitive data`
	sugar.Infow("config", "api_key", "ak-xxx")    // want `logcheck: sensitive: key "api_key" may contain sensitive data`
	sugar.Infow("session", "auth_token", "tok123") // want `logcheck: sensitive: key "auth_token" may contain sensitive data`
	sugar.Info("password=hunter2")                  // want `logcheck: sensitive: message may contain embedded credentials`
	sugar.Infow("request", "method", "GET")
}

// --- крайние случаи ---

func zapEdgeCases() {
	logger := zap.NewNop()
	msg := getMessage()
	logger.Info(msg) // динамическое — пропускается

	sugar := zap.NewNop().Sugar()
	sugar.Info(msg) // динамическое — пропускается
}

func getMessage() string { return "dynamic" }
