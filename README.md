# logcheck

Go-линтер для проверки сообщений в вызовах `log/slog` и `go.uber.org/zap`. Проверяет стиль, язык, спецсимволы и утечку чувствительных данных.

## Правила

| Правило | Описание | ❌ Плохо | ✅ Хорошо |
|---------|----------|----------|-----------|
| **lowercase** | Сообщение должно начинаться со строчной буквы | `slog.Info("User logged in")` | `slog.Info("user logged in")` |
| **english** | Только английский текст (ASCII) | `slog.Info("пользователь вошёл")` | `slog.Info("user logged in")` |
| **specialchars** | Без спецсимволов и эмодзи | `slog.Info("done! 🎉")` | `slog.Info("done")` |
| **sensitive** | Нет чувствительных данных в ключах и сообщениях | `slog.Info("auth", "password", pw)` | `slog.Info("auth", "user_id", id)` |

Правила **lowercase** и **specialchars** поддерживают автоматическое исправление (`SuggestedFix`).

## Установка

### Standalone

```bash
go install github.com/festy23/logcheck/cmd/logcheck@latest
```

### golangci-lint (module plugin)

Добавьте в `.custom-gcl.yml`:

```yaml
version: v2.10.1
plugins:
  - module: 'github.com/festy23/logcheck'
    import: 'github.com/festy23/logcheck/plugin'
```

Соберите кастомный golangci-lint:

```bash
golangci-lint custom
```

## Конфигурация

### Конфигурационный файл (JSON)

Создайте `.logcheck.json` в корне проекта:

```json
{
  "disable": ["english"],
  "sensitive_patterns": ["session_key", "internal_id"]
}
```

Передайте путь через флаг:

```bash
logcheck -logcheck.config=.logcheck.json ./...
```

### Флаги анализатора

| Флаг | Описание | Пример |
|------|----------|--------|
| `-logcheck.config` | Путь к JSON-файлу конфигурации | `-logcheck.config=.logcheck.json` |
| `-logcheck.disable` | Отключить правила (через запятую) | `-logcheck.disable=english,specialchars` |
| `-logcheck.sensitive-patterns` | Дополнительные шаблоны чувствительных данных | `-logcheck.sensitive-patterns=session_key,internal_id` |

Флаги имеют приоритет над конфигурационным файлом.

### Пример `.golangci.yml`

```yaml
linters:
  enable:
    - logcheck

linters-settings:
  custom:
    logcheck:
      type: module
      description: Checks log messages for common issues
      settings:
        config: ".logcheck.json"
        disable: "english"
        sensitive-patterns: "session_key,internal_id"
```

## Запуск

### Standalone

```bash
# Проверить текущий пакет
logcheck ./...

# С отключением правила
logcheck -logcheck.disable=english ./...

# С дополнительными паттернами
logcheck -logcheck.sensitive-patterns=session_key ./...
```

### Пример вывода

```
main.go:12:12: logcheck: lowercase: message must start with a lowercase letter
main.go:15:12: logcheck: english: message contains non-English characters
main.go:19:12: logcheck: specialchars: message contains special characters
main.go:24:12: logcheck: sensitive: message may contain embedded credentials
main.go:27:12: logcheck: sensitive: key "password" may contain sensitive data
```

### Через go vet

```bash
go vet -vettool=$(which logcheck) ./...
```

### Через golangci-lint

```bash
./custom-gcl run --enable logcheck ./...
```

## Разработка

### Сборка и тесты

```bash
# Собрать бинарник
go build ./cmd/logcheck/

# Запустить все тесты
go test ./... -count=1

# Тесты с race detector
go test ./... -count=1 -race

# go vet
go vet ./...

# Интеграционные тесты (analysistest)
go test ./pkg/analyzer/ -run TestAll -v
```

### Git-хуки (Lefthook)

Проект использует [Lefthook](https://github.com/evilmartians/lefthook) для pre-commit проверок.

Установка:
```bash
go install github.com/evilmartians/lefthook@latest
go install golang.org/x/tools/cmd/goimports@latest
lefthook install
```

Перед каждым коммитом автоматически запускаются:

**Форматтеры** (последовательно):
1. `go fix` — модернизация кода (Go 1.26)
2. `gofmt` — форматирование
3. `goimports` — порядок импортов
4. `go mod tidy` — проверка зависимостей

**Проверки** (параллельно):
- `go vet` — статический анализ
- `go test` — юнит-тесты
- `go build` — компиляция

### Структура проекта

```
pkg/
├── model/       — LogCall, KeyValue, ConcatPart
├── extractor/   — обнаружение вызовов логгера, извлечение сообщений
├── rules/       — правила проверки + реестр
└── analyzer/    — точка входа analysis.Analyzer
cmd/logcheck/    — standalone CLI (singlechecker)
plugin/          — golangci-lint module plugin
```

**Pipeline:** extractor (находит вызовы) → rules (проверяет каждый) → analyzer (связывает всё вместе)

### Добавление нового правила

1. Создайте `pkg/rules/<name>.go` с реализацией интерфейса `Rule`:

```go
type Rule interface {
    Name() string
    Description() string
    Check(call *model.LogCall, pass *analysis.Pass)
}
```

2. Зарегистрируйте правило в `pkg/rules/registry.go` → `defaultRules()`
3. Создайте `pkg/rules/<name>_test.go` с unit-тестами
4. Создайте testdata-пакет `pkg/analyzer/testdata/src/<name>/` с `// want` комментариями
5. Добавьте пакет в `pkg/analyzer/analyzer_test.go` → `TestAll`
6. Запустите: `go vet ./... && go test ./... -count=1`
