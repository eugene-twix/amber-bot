# bot/

Telegram бот: инициализация, handlers, middleware.

## Файлы

| Файл | Описание |
|------|----------|
| `bot.go` | Структура Bot, инициализация, регистрация handlers |
| `auth.go` | Middleware авторизации, проверка ролей |
| `handlers.go` | Публичные команды (start, teams, rating, cancel) |
| `handlers_org.go` | Команды организаторов (newteam, addmember, newtournament, result) |
| `handlers_admin.go` | Админские команды (grant) |

## Команды

### Публичные
- `/start` — приветствие
- `/teams` — список команд
- `/team <name>` — статистика команды
- `/rating` — рейтинг команд
- `/cancel` — отмена FSM

### Организаторы
- `/newteam` — создать команду
- `/addmember` — добавить участника
- `/newtournament` — создать турнир
- `/result` — записать результат

### Админы
- `/grant` — назначить роль

## Архитектура

```
Bot struct {
    tg         *tele.Bot       // telebot
    cfg        *config.Config
    db         *bun.DB
    cache      *cache.Cache
    fsm        *fsm.Manager
    *Repo      // репозитории
}
```

## Middleware

`authMiddleware` — создаёт/получает User и кладёт в context.
