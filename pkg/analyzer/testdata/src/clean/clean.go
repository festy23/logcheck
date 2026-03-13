package clean

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/zap"
)

// validSlog — различные методы slog с чистыми сообщениями, ожидается 0 диагностик.
func validSlog() {
	slog.Info("hello world")
	slog.Debug("processing request")
	slog.Warn("connection slow")
	slog.Error("request failed")
	slog.Info("3 retries left")
	slog.Info("")
	slog.Info("_internal thing")
	slog.Info("[debug] something")
	slog.Info("file: /tmp/data.json")
	slog.Info("100% complete")
	slog.Info("key=value pairs")
	slog.Info("hello world!")
	slog.Info("user@example.com logged in")
	slog.Info("reset password")
	slog.Info("token count exceeded")
}

func validSlogKV() {
	slog.Info("user logged in", "user_id", "u123")
	slog.Info("request completed", "method", "GET", "status", 200)
	slog.Info("item processed", slog.String("item_id", "abc"))
}

func validSlogContext() {
	ctx := context.Background()
	slog.InfoContext(ctx, "request started")
	slog.WarnContext(ctx, "connection slow")
	slog.ErrorContext(ctx, "request failed")
	slog.DebugContext(ctx, "processing item")
}

func validSlogLog() {
	ctx := context.Background()
	slog.Log(ctx, slog.LevelInfo, "custom log entry")
	slog.LogAttrs(ctx, slog.LevelInfo, "log with attrs", slog.String("key", "val"))
}

func validVariable() {
	msg := getMessage()
	slog.Info(msg)
}

func validRawString() {
	slog.Info(`raw string message`)
}

func validFmtSprintf() {
	slog.Info(fmt.Sprintf("request %d completed", 42))
}

func validConcat() {
	slog.Info("hello " + "world")
}

func validConst() {
	const msg = "server started"
	slog.Info(msg)
}

func validZapLogger() {
	logger := zap.NewNop()
	logger.Info("hello world")
	logger.Debug("processing request")
	logger.Warn("connection slow")
	logger.Error("request failed")
	logger.Info("request completed", zap.String("method", "GET"), zap.Int("status", 200))
}

func validZapSugar() {
	sugar := zap.NewNop().Sugar()
	sugar.Info("hello world")
	sugar.Infow("request completed", "method", "GET", "status", 200)
	sugar.Infof("request %d completed", 42)
	sugar.Debug("processing request")
	sugar.Debugw("processing item", "item_id", "abc")
	sugar.Warn("connection slow")
	sugar.Warnw("connection issue", "retry", 3)
	sugar.Error("request failed")
	sugar.Errorw("request failed", "code", 500)
}

func validNonLogging() {
	logger := zap.NewNop()
	logger.With(zap.String("module", "auth"))
	logger.Named("mylogger")
	_ = logger.Sugar()
	_ = logger.Sync()

	sugar := logger.Sugar()
	sugar.With("key", "val")
	_ = sugar.Desugar()
}

func getMessage() string { return "dynamic" }
