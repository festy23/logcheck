package specialchars

import "log/slog"

func good() {
	slog.Info("connection established")
	slog.Info("3 retries left")
	slog.Info("")
	slog.Info("file: /tmp/data.json")
	slog.Info("100% complete")
	slog.Info("[debug] something")
	slog.Info("key=value pairs")
	slog.Info("hello world!")
	slog.Info("user@example.com logged in")
	slog.Info("single dot.")
	slog.Info("two dots..")
}

func bad() {
	slog.Info("server started 🚀")        // want `logcheck: specialchars: message contains special characters`
	slog.Error("build failed ❌")          // want `logcheck: specialchars: message contains special characters`
	slog.Info("check passed ✅")           // want `logcheck: specialchars: message contains special characters`
	slog.Info("arrow → here")              // want `logcheck: specialchars: message contains special characters`
	slog.Info("star ★ rating")             // want `logcheck: specialchars: message contains special characters`
	slog.Info("copyright ©")               // want `logcheck: specialchars: message contains special characters`
	slog.Info("hello\tworld")              // want `logcheck: specialchars: message contains special characters`
	slog.Info("line1\nline2")              // want `logcheck: specialchars: message contains special characters`
	slog.Error("connection failed!!!")     // want `logcheck: specialchars: message contains special characters`
	slog.Warn("something went wrong...")   // want `logcheck: specialchars: message contains special characters`
	slog.Info("really???")                 // want `logcheck: specialchars: message contains special characters`
}

func edgeCases() {
	msg := getMessage()
	slog.Info(msg)

	const rocket = "launch 🚀"
	slog.Info(rocket) // want `logcheck: specialchars: message contains special characters`

	slog.Info("hello " + "🌍")  // want `logcheck: specialchars: message contains special characters`
	slog.Info("hello " + "world")
}

func getMessage() string { return "dynamic" }
