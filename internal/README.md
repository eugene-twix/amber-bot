# internal/

Внутренние пакеты приложения (не экспортируются).

## Структура

```
internal/
├── api/        # REST API (Gin): handlers, middleware
├── bot/        # Telegram бот: handlers, middleware
├── cache/      # Dragonfly (Redis) клиент
├── config/     # Загрузка конфигурации из ENV
├── domain/     # Доменные сущности (User, Team, etc.)
├── fsm/        # FSM для многошаговых диалогов бота
├── migrations/ # SQL миграции (применяются автоматически)
└── repository/ # Интерфейсы и реализации репозиториев
```

## Поток данных

### Mini App (API)

```
Frontend (React)
    ↓
api/           ← REST handlers, auth middleware
    ↓
repository/    ← Bun ORM
    ↓
PostgreSQL
```

### Telegram Bot

```
Telegram Message
    ↓
bot/           ← handlers, middleware авторизации
    ↓
fsm/           ← состояния диалогов (Redis)
    ↓
repository/    ← Bun ORM
    ↓
PostgreSQL
```
