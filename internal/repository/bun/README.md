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
db := bun.NewDB(sqldb, pgdialect.New())
userRepo := bunrepo.NewUserRepo(db)
```

## Зависимости

- `github.com/uptrace/bun` — ORM
- `github.com/uptrace/bun/dialect/pgdialect` — PostgreSQL диалект
