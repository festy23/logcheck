# log/slog API Reference

## Package-level functions (also available as methods on *slog.Logger)

| Function                         | Msg Index | Notes                           |
|----------------------------------|-----------|----------------------------------|
| slog.Debug(msg, args...)         | 0         |                                  |
| slog.Info(msg, args...)          | 0         |                                  |
| slog.Warn(msg, args...)          | 0         |                                  |
| slog.Error(msg, args...)         | 0         |                                  |
| slog.DebugContext(ctx, msg, ...) | 1         | first arg is context.Context     |
| slog.InfoContext(ctx, msg, ...)  | 1         | first arg is context.Context     |
| slog.WarnContext(ctx, msg, ...)  | 1         | first arg is context.Context     |
| slog.ErrorContext(ctx, msg, ...) | 1         | first arg is context.Context     |
| slog.Log(ctx, level, msg, ...)   | 2         | ctx + slog.Level before msg      |
| slog.LogAttrs(ctx, lvl, msg, ..) | 2         | ctx + slog.Level before msg      |

## Key-Value extraction

After msg index, slog uses alternating key-value pairs:
```go
slog.Info("msg", "key1", val1, "key2", val2)
```
Keys are string literals, values are any type.
Exception: `slog.Attr` values can also appear (not string-keyed).

## Methods to IGNORE (no message argument)

- `slog.With(args...)` — returns new Logger, no message
- `slog.Default()` — returns default Logger
- `slog.SetDefault(l)` — sets default Logger
- `slog.NewTextHandler(...)`, `slog.NewJSONHandler(...)` — constructors
- `slog.New(handler)` — constructor
- `slog.Group(key, args...)` — creates attribute group
