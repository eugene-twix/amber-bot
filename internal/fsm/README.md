# fsm/

Finite State Machine для многошаговых диалогов.

## Файлы

- `state.go` — определение состояний и типа Data
- `manager.go` — Manager для работы с состояниями (Redis)

## Состояния

| Flow | Состояния |
|------|-----------|
| NewTeam | `new_team:name` |
| AddMember | `add_member:team` → `add_member:name` |
| NewTournament | `new_tournament:name` → `date` → `location` |
| Result | `result:tournament` → `team` → `place` |
| Grant | `grant:user` → `role` |

## API Manager

```go
m := fsm.NewManager(cache)
m.SetState(ctx, userID, state, data)
state, data, _ := m.GetState(ctx, userID)
m.Clear(ctx, userID)
```

## Data helpers

```go
data.GetInt64("team_id")
data.GetString("name")
```

## Хранение

Состояния хранятся в Redis с TTL (обычно 5-10 минут).
