# repository/bun/

Реализация репозиториев на Bun ORM.

## Файлы

| Файл | Репозиторий | Описание |
|------|-------------|----------|
| `db.go` | — | Инициализация подключения к PostgreSQL |
| `user.go` | `UserRepo` | CRUD пользователей |
| `team.go` | `TeamRepo` | CRUD команд |
| `member.go` | `MemberRepo` | CRUD участников команд |
| `tournament.go` | `TournamentRepo` | CRUD турниров |
| `result.go` | `ResultRepo` | CRUD результатов + рейтинг |

## Использование

```go
db := bunrepo.NewDB(databaseURL, debug)
defer db.Close()

userRepo := bunrepo.NewUserRepo(db)
teamRepo := bunrepo.NewTeamRepo(db)
// ...
```

## Soft Delete

Все методы Get/List фильтруют `WHERE deleted_at IS NULL`:
- Удалённые записи не возвращаются в запросах
- Delete() ставит `deleted_at = NOW()` вместо физического удаления

## Зависимости

- `github.com/uptrace/bun` — ORM
- `github.com/uptrace/bun/dialect/pgdialect` — PostgreSQL диалект
- `github.com/uptrace/bun/driver/pgdriver` — драйвер
