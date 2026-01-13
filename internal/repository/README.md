# repository/

Слой доступа к данным.

## Файлы

- `interfaces.go` — интерфейсы репозиториев
- `bun/` — реализация на Bun ORM

## Интерфейсы

| Интерфейс | Методы |
|-----------|--------|
| `UserRepository` | GetOrCreate, GetByTelegramID, UpdateRole, List |
| `TeamRepository` | Create, GetByID, GetByName, List, Update, Delete |
| `MemberRepository` | Create, GetByID, GetByTeamID, Update, Delete |
| `TournamentRepository` | Create, GetByID, List, ListRecent, Update, Delete |
| `ResultRepository` | Create, GetByID, GetByTeamID, GetByTournamentID, GetTeamRating, Update, Delete, DeleteWithShift |

## Типы

```go
type TeamRating struct {
    TeamID     int64
    TeamName   string
    Wins       int     // количество 1-х мест
    TotalGames int
    AvgPlace   float64
}
```

## Soft Delete

Все репозитории поддерживают soft delete:
- Delete() устанавливает `deleted_at` вместо физического удаления
- Все запросы фильтруют `WHERE deleted_at IS NULL`
- GetTeamRating() фильтрует удалённые записи в JOIN

## DeleteWithShift

`ResultRepository.DeleteWithShift()` — удаляет результат и сдвигает места:
- Если удаляем 3-е место из [1,2,3,4,5], получаем [1,2,3,4]
- Выполняется в транзакции
