# cache/

Клиент для Dragonfly (Redis-совместимый кеш).

## Файлы

- `dragonfly.go` — обёртка над go-redis клиентом

## Использование

```go
cache, err := cache.New(redisURL)
cache.Set(ctx, key, value, ttl)
cache.Get(ctx, key)
cache.Del(ctx, key)
```

## Назначение

- Хранение состояний FSM (многошаговые диалоги)
- Возможно кеширование данных в будущем

## Зависимости

- `github.com/redis/go-redis/v9`
- `REDIS_URL` — переменная окружения
