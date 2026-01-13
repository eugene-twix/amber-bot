# domain/

Доменные сущности приложения.

## Сущности

| Файл | Структура | Описание |
|------|-----------|----------|
| `user.go` | `User` | Пользователь Telegram с ролью (viewer/organizer/admin) |
| `team.go` | `Team` | Команда квиза |
| `member.go` | `Member` | Участник команды |
| `tournament.go` | `Tournament` | Турнир (название, дата, место) |
| `result.go` | `Result` | Результат команды на турнире (место) |

## Роли пользователей

```go
RoleViewer    // просмотр
RoleOrganizer // + создание команд, турниров, результатов
RoleAdmin     // + назначение ролей
```

## Связи

```
User (1) ──── (*) управляет командами
Team (1) ──── (*) Member
Team (1) ──── (*) Result
Tournament (1) ──── (*) Result
```
