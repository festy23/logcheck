package zap

// Field — поле структурированного логирования.
type Field struct{}

// Logger — быстрый структурированный логгер.
type Logger struct{}

// NewNop возвращает Logger без операций.
func NewNop() *Logger { return &Logger{} }

func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) With(fields ...Field) *Logger       { return l }
func (l *Logger) Named(name string) *Logger          { return l }
func (l *Logger) Sugar() *SugaredLogger              { return &SugaredLogger{} }
func (l *Logger) Sync() error                        { return nil }

// SugaredLogger оборачивает Logger в упрощённый API.
type SugaredLogger struct{}

func (s *SugaredLogger) Info(args ...any)                         {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...any)   {}
func (s *SugaredLogger) Infof(template string, args ...any)       {}
func (s *SugaredLogger) Error(args ...any)                        {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Errorf(template string, args ...any)      {}
func (s *SugaredLogger) Debug(args ...any)                        {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Debugf(template string, args ...any)      {}
func (s *SugaredLogger) Warn(args ...any)                         {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...any)   {}
func (s *SugaredLogger) Warnf(template string, args ...any)       {}
func (s *SugaredLogger) DPanic(args ...any)                       {}
func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...any) {}
func (s *SugaredLogger) DPanicf(template string, args ...any)     {}
func (s *SugaredLogger) Fatal(args ...any)                        {}
func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Fatalf(template string, args ...any)      {}
func (s *SugaredLogger) Panic(args ...any)                        {}
func (s *SugaredLogger) Panicw(msg string, keysAndValues ...any)  {}
func (s *SugaredLogger) Panicf(template string, args ...any)      {}
func (s *SugaredLogger) With(args ...any) *SugaredLogger          { return s }
func (s *SugaredLogger) Desugar() *Logger                         { return &Logger{} }

// Конструкторы полей.
func String(key string, val string) Field   { return Field{} }
func Int(key string, val int) Field         { return Field{} }
func Bool(key string, val bool) Field       { return Field{} }
func Any(key string, val any) Field         { return Field{} }
func Float64(key string, val float64) Field { return Field{} }
