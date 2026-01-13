# repository/

Слой доступа к данным.

## Файлы

- `interfaces.go` — интерфейсы репозиториев
- `bun/` — реализация на Bun ORM

## Интерфейсы

| Интерфейс | Методы |
|-----------|--------|
| `UserRepository` | GetOrCreate, GetByTelegramID, UpdateRole |
| `TeamRepository` | Create, GetByID, GetByName, List, Delete |
| `MemberRepository` | Create, GetByTeamID, Delete |
| `TournamentRepository` | Create, GetByID, List, ListRecent |
| `ResultRepository` | Create, GetByTeamID, GetByTournamentID, GetTeamRating |

## Типы

```go
type TeamRating struct {
    TeamID, TeamName string
    Wins, TotalGames int
    AvgPlace         float64
}
```
