package alias

import (
	s "log/slog"

	z "go.uber.org/zap"
)

func aliasedSlog() {
	s.Info("hello world")
	s.Info("Hello world") // want `logcheck: lowercase: message must start with a lowercase letter`
	s.Error("request failed")
	s.Debug("processing")
}

func aliasedZap() {
	logger := z.NewNop()
	logger.Info("hello world")
	logger.Info("Hello world") // want `logcheck: lowercase: message must start with a lowercase letter`
}

func aliasedZapSugar() {
	sugar := z.NewNop().Sugar()
	sugar.Info("hello world")
	sugar.Info("Hello world") // want `logcheck: lowercase: message must start with a lowercase letter`
	sugar.Infow("auth", "password", "hunter2") // want `logcheck: sensitive: key "password" may contain sensitive data`
}
