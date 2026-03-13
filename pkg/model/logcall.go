package model

import "go/token"

// LoggerType определяет, какой фреймворк логирования произвёл вызов.
type LoggerType int

const (
	LoggerSlog     LoggerType = iota // log/slog
	LoggerZap                        // *zap.Logger
	LoggerZapSugar                   // *zap.SugaredLogger
)

// ConcatPart представляет один сегмент выражения конкатенации строк.
type ConcatPart struct {
	Value     string    // разрешённое строковое содержимое (пустое, если не литерал)
	Pos       token.Pos // позиция этой части в исходном коде
	IsLiteral bool      // true, если часть является разрешённым строковым литералом или константой
}

// KeyValue представляет одну пару ключ-значение из аргументов структурированного логирования.
type KeyValue struct {
	Key    string    // строка ключа (например, "user_id")
	KeyPos token.Pos // позиция ключа в исходном коде
}

// LogCall представляет один обнаруженный вызов логирования с извлечёнными метаданными.
type LogCall struct {
	Pos         token.Pos    // позиция всего выражения вызова
	Logger      LoggerType   // фреймворк логирования
	Method      string       // имя метода (например, "Info", "Warnf", "Errorw")
	MsgLiteral  string       // разрешённое содержимое строки сообщения (пустое, если не литерал)
	HasLiteral  bool         // true, если аргумент сообщения является разрешимым строковым литералом
	MsgPos      token.Pos    // позиция аргумента сообщения
	ConcatParts []ConcatPart // части, если сообщение является выражением конкатенации
	KeyValues   []KeyValue   // извлечённые пары ключ-значение из структурированных аргументов
}
