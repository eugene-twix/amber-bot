# Amber Bot

Telegram Mini App + бот для учёта команд, участников и результатов квизов и настольных игр.

## Стек технологий

- **Backend:** Go, Gin, Bun ORM
- **Frontend:** React, TypeScript, Vite, shadcn/ui, TanStack Query
- **База данных:** PostgreSQL
- **Кеш:** Dragonfly (Redis-совместимый)
- **Развертывание:** Docker

---

## Возможности

### Mini App (основной интерфейс)

- **Команды**: Просмотр списка команд, поиск, детали команды с участниками и результатами
- **Рейтинг**: Таблица рейтинга команд с сортировкой по победам и среднему месту
- **Турниры**: Список турниров и их результаты
- **Управление** (organizer/admin): Создание команд, турниров, участников, запись результатов

### Роли пользователей

| Роль | Права |
|------|-------|
| `viewer` | Просмотр команд, рейтинга, турниров |
| `organizer` | + создание команд, турниров, запись результатов |
| `admin` | + управление ролями пользователей |

---

## Структура репозитория

```
.
├── cmd/
│   ├── api/main.go      # HTTP API сервер
│   └── bot/main.go      # Telegram бот
├── frontend/            # React Mini App
│   ├── src/
│   │   ├── components/  # UI компоненты
│   │   ├── pages/       # Страницы
│   │   ├── hooks/       # React hooks
│   │   └── services/    # API клиент
│   └── ...
├── internal/
│   ├── api/             # REST API (Gin)
│   ├── bot/             # Telegram бот (telebot)
│   ├── cache/           # Redis клиент
│   ├── config/          # Конфигурация
│   ├── domain/          # Доменные сущности
│   ├── fsm/             # FSM для диалогов бота
│   ├── migrations/      # SQL миграции
│   └── repository/      # Слой данных (Bun ORM)
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Установка и запуск

### 1. Предварительные требования

- `git`, `go` (1.21+), `node` (20+), `docker`, `docker-compose`
- `age` (для шифрования секретов)

### 2. Клонирование

```bash
git clone <repository-url>
cd amber-bot
```

### 3. Запуск инфраструктуры

```bash
make up  # PostgreSQL + Dragonfly
```

### 4. Настройка секретов

```bash
# Сгенерировать age ключ
make secrets-keygen

# Создать .secret.env
cp .secret.env.example .secret.env
# Заполнить TELEGRAM_TOKEN, DATABASE_URL, ADMIN_IDS

# Зашифровать
make secrets-encrypt
```

### 5. Конфигурация (.env)

```env
DATABASE_URL=postgres://amber:amber@localhost:5432/amber_bot?sslmode=disable
REDIS_URL=redis://localhost:6379/0
ADMIN_IDS=123456789

# Для Mini App
API_PORT=8080
MINI_APP_URL=https://your-domain.com

# Dev режим (без Telegram auth)
DEV_MODE=true
DEV_USER_ID=123456789
```

### 6. Запуск

```bash
# Backend (API сервер) — миграции применяются автоматически
go run cmd/api/main.go

# Frontend (dev)
cd frontend && npm install && npm run dev

# Telegram бот (опционально)
go run cmd/bot/main.go
```

Или через Makefile:

```bash
make run-api      # API сервер
make run-bot      # Telegram бот
make run-frontend # Frontend dev server
```

---

## API

### Публичные endpoints (`/api/v1/public/*`)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/me` | Текущий пользователь |
| GET | `/teams` | Список команд |
| GET | `/teams/:id` | Детали команды |
| GET | `/teams/:id/members` | Участники команды |
| GET | `/teams/:id/results` | Результаты команды |
| GET | `/tournaments` | Список турниров |
| GET | `/tournaments/:id` | Детали турнира |
| GET | `/tournaments/:id/results` | Результаты турнира |
| GET | `/rating` | Рейтинг команд |

### Приватные endpoints (`/api/v1/private/*`)

Требуют роль `organizer` или `admin`.

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/teams` | Создать команду |
| PATCH | `/teams/:id` | Обновить команду |
| DELETE | `/teams/:id` | Удалить команду |
| POST | `/teams/:id/members` | Добавить участника |
| POST | `/tournaments` | Создать турнир |
| POST | `/tournaments/:id/results` | Записать результат |
| ... | ... | ... |

---

## Разработка

```bash
make up           # Запустить Docker (DB + Cache)
make down         # Остановить Docker
make run-api      # Запустить API сервер
make run-bot      # Запустить бота
make test         # Запустить тесты
make build        # Собрать бинарники
```

### Секреты

```bash
make secrets-keygen   # Сгенерировать age ключ
make secrets-encrypt  # Зашифровать .secret.env
make secrets-decrypt  # Расшифровать .secret.enc.env
```

---

## Лицензия

Проект не имеет определенной лицензии. Все права защищены.
