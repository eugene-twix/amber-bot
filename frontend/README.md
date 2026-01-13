# frontend/

Telegram Mini App для Amber Bot.

## Технологии

- **React 19** + TypeScript
- **Vite** — сборка
- **TanStack Query** — кеширование данных
- **shadcn/ui** — UI компоненты
- **Tailwind CSS** — стили

## Структура

```
frontend/
├── src/
│   ├── components/     # UI компоненты
│   │   ├── ui/         # shadcn/ui компоненты
│   │   └── TabBar.tsx  # Нижняя навигация
│   ├── hooks/          # React hooks (useApi)
│   ├── pages/          # Страницы приложения
│   ├── services/       # API клиент
│   ├── lib/            # Утилиты
│   ├── App.tsx         # Роутинг
│   └── main.tsx        # Entry point
├── index.html          # HTML + Telegram Web App SDK
└── vite.config.ts      # Конфигурация Vite
```

## Страницы

| Страница | Путь | Описание |
|----------|------|----------|
| TeamsPage | `/` | Список команд с поиском |
| TeamDetailPage | `/teams/:id` | Детали команды + участники + результаты |
| RatingPage | `/rating` | Рейтинг команд с сортировкой |
| TournamentsPage | `/tournaments` | Список турниров |
| TournamentDetailPage | `/tournaments/:id` | Детали турнира + результаты |
| ManagePage | `/manage` | Управление (только organizer/admin) |

## Запуск

```bash
cd frontend

# Установка зависимостей
npm install

# Dev сервер (порт 5173)
npm run dev

# Сборка
npm run build
```

## API

Фронтенд обращается к API через `/api/v1/`:
- Публичные endpoints: `/api/v1/public/*`
- Приватные endpoints: `/api/v1/private/*` (organizer/admin)

Авторизация через Telegram Web App `initData`.

## Dev режим

При `DEV_MODE=true` на бэкенде авторизация отключена, используется `DEV_USER_ID`.
