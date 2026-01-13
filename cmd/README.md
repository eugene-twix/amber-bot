# cmd/

Точки входа приложения.

## Структура

```
cmd/
├── api/main.go   # HTTP API сервер (Mini App backend)
└── bot/main.go   # Telegram бот
```

## Запуск

```bash
# API сервер (порт 8080)
go run cmd/api/main.go

# Telegram бот
go run cmd/bot/main.go
```

## Примечания

- API сервер автоматически применяет миграции при старте
- Бот и API — независимые процессы, можно запускать отдельно
- Для Mini App достаточно только API сервера
