# config/

Загрузка конфигурации из переменных окружения.

## Файлы

- `config.go` — структура Config и функция Load()
- `config_test.go` — тесты

## Структура Config

```go
type Config struct {
    TelegramToken string   // TELEGRAM_TOKEN
    DatabaseURL   string   // DATABASE_URL
    RedisURL      string   // REDIS_URL
    AdminIDs      []int64  // ADMIN_IDS (через запятую)
}
```

## Использование

```go
cfg, err := config.Load()
```

## Переменные окружения

| Переменная | Обязательная | Описание |
|------------|--------------|----------|
| `TELEGRAM_TOKEN` | да | Токен от @BotFather |
| `DATABASE_URL` | да | PostgreSQL DSN |
| `REDIS_URL` | да | Redis/Dragonfly URL |
| `ADMIN_IDS` | нет | ID админов через запятую |
