package sensitive

import "log/slog"

func good() {
	slog.Info("connection established")
	slog.Info("user logged in", "user_id", "u123")
	slog.Info("request completed", "method", "GET", "status", 200)
	slog.Info("reset password")
	slog.Info("token count exceeded")
	slog.Info("")
}

func bad() {
	pw := "hunter2"
	slog.Info("login attempt", "password", pw)    // want `logcheck: sensitive: key "password" may contain sensitive data`
	slog.Info("auth", "secret", "s3cr3t")          // want `logcheck: sensitive: key "secret" may contain sensitive data`
	slog.Info("session", "auth_token", "tok123")   // want `logcheck: sensitive: key "auth_token" may contain sensitive data`
	slog.Info("config", "api_key", "ak-xxx")       // want `logcheck: sensitive: key "api_key" may contain sensitive data`
	slog.Info("user", "ssn", "123-45-6789")        // want `logcheck: sensitive: key "ssn" may contain sensitive data`
	slog.Info("payment", "credit_card", "4111xxx") // want `logcheck: sensitive: key "credit_card" may contain sensitive data`
	slog.Info("tls", "private_key", "-----BEGIN")  // want `logcheck: sensitive: key "private_key" may contain sensitive data`
	slog.Info("auth", "authorization", "Bearer x") // want `logcheck: sensitive: key "authorization" may contain sensitive data`
}

func badMessage() {
	slog.Info("password=hunter2")  // want `logcheck: sensitive: message may contain embedded credentials`
	slog.Info("token: abc123")     // want `logcheck: sensitive: message may contain embedded credentials`
	slog.Info("secret =s3cr3t")    // want `logcheck: sensitive: message may contain embedded credentials`
	slog.Info("api_key=ak-xxx")    // want `logcheck: sensitive: message may contain embedded credentials`
}

func edgeCases() {
	msg := getMessage()
	slog.Info(msg)

	slog.Info("password=" + getMessage()) // want `logcheck: sensitive: message may contain embedded credentials`

	slog.Info("user " + "logged in")

	slog.Info("msg", "Password", "x")  // want `logcheck: sensitive: key "Password" may contain sensitive data`
	slog.Info("msg", "API-KEY", "x")   // want `logcheck: sensitive: key "API-KEY" may contain sensitive data`
	slog.Info("msg", "user_name", "x")
}

func getMessage() string { return "dynamic" }
