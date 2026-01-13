# config/

Загрузка конфигурации из переменных окружения.

## Файлы

- `config.go` — структура Config и функция Load()

## Структура Config

```go
type Config struct {
    TelegramToken string   // TELEGRAM_TOKEN
    DatabaseURL   string   // DATABASE_URL
    RedisURL      string   // REDIS_URL
    AdminIDs      []int64  // ADMIN_IDS (через запятую)

    // API сервер
    APIPort      int    // API_PORT (default: 8080)
    FrontendPath string // FRONTEND_PATH (default: ./frontend/dist)

    // Mini App
    MiniAppURL string // MINI_APP_URL

    // Dev режим
    DevMode   bool  // DEV_MODE (default: false)
    DevUserID int64 // DEV_USER_ID (default: 123456789)
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
| `API_PORT` | нет | Порт API сервера (default: 8080) |
| `FRONTEND_PATH` | нет | Путь к frontend/dist |
| `MINI_APP_URL` | нет | URL Mini App для кнопки в боте |
| `DEV_MODE` | нет | Режим разработки (без Telegram auth) |
| `DEV_USER_ID` | нет | User ID для dev режима |
