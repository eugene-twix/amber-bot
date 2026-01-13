# internal/

Внутренние пакеты приложения (не экспортируются).

## Структура

```
internal/
├── bot/        # Telegram бот: handlers, middleware
├── cache/      # Dragonfly (Redis) клиент
├── config/     # Загрузка конфигурации из ENV
├── domain/     # Доменные сущности (User, Team, etc.)
├── fsm/        # FSM для многошаговых диалогов
└── repository/ # Интерфейсы и реализации репозиториев
```

## Поток данных

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
