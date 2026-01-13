# bot/

Telegram бот: инициализация, handlers, middleware.

## Файлы

| Файл | Описание |
|------|----------|
| `bot.go` | Структура Bot, инициализация, регистрация handlers |
| `auth.go` | Middleware авторизации, создание/получение User |
| `handlers.go` | Публичные команды (/start, teams, rating, cancel) |
| `handlers_org.go` | Команды организаторов (newteam, addmember, newtournament, result) |
| `handlers_admin.go` | Админские команды (grant) |
| `keyboard.go` | Reply/Inline клавиатуры, пагинация |

## Интерфейс

Бот использует **Mini App** как основной интерфейс. При `/start` показывается inline-кнопка для открытия Mini App.

### Команда /start

- Убирает Reply Keyboard (если была)
- Показывает WebApp кнопку (если `MINI_APP_URL` настроен)
- Приветствует пользователя с указанием роли

### FSM диалоги (устаревшие, для обратной совместимости)

Многошаговые диалоги через Reply Keyboard:
- Создание команды
- Добавление участника
- Создание турнира
- Запись результата
- Назначение роли (admin)

## Архитектура

```go
Bot struct {
    tg         *tele.Bot       // telebot
    cfg        *config.Config
    db         *bun.DB
    cache      *cache.Cache
    fsm        *fsm.Manager
    miniAppURL string          // URL Mini App
    // репозитории...
}
```

## Middleware

`authMiddleware` — создаёт/получает User и кладёт в context.
