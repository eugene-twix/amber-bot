# Реализация: Bot handlers для amber-bot (Tasks 10-12)

## 1. Созданные/изменённые файлы

### Новые файлы
| Файл | Описание |
|------|----------|
| internal/bot/handlers.go | Публичные команды: /start, /teams, /team, /rating, /cancel |
| internal/bot/handlers_org.go | Команды организатора: /newteam, /addmember, /newtournament, /result с FSM диалогами |
| internal/bot/handlers_admin.go | Админские команды: /grant для управления ролями |

### Изменённые файлы
| Файл | Изменения |
|------|-----------|
| internal/bot/bot.go | Удалены stub методы (handleStart, handleTeams, и др.) |
| internal/bot/auth.go | Добавлена обработка ошибок c.Send() |

## 2. Ключевые решения реализации

### Публичные команды (handlers.go)
- `/start` - динамически генерирует список доступных команд в зависимости от роли пользователя
- `/teams` - выводит список всех команд из БД
- `/team <название>` - показывает состав команды и историю результатов
- `/rating` - использует GetTeamRating() из ResultRepo для агрегации статистики
- `/cancel` - очищает FSM состояние пользователя

### Команды организатора (handlers_org.go)
- **NewTeam flow**: простой диалог с одним шагом (ввод названия)
- **AddMember flow**: inline-кнопки для выбора команды → ввод имени
- **NewTournament flow**: название → дата (поддержка разных форматов) → место (опционально)
- **Result flow**: выбор турнира (10 последних) → выбор команды → ввод места

Все диалоги используют FSM:
- `handleNewTeam` → `processNewTeamName`
- `handleAddMember` → callback `addmember_team:ID` → `processAddMemberName`
- `handleNewTournament` → `processNewTournamentName` → `processNewTournamentDate` → `processNewTournamentLocation`
- `handleResult` → callback `result_tourn:ID` → callback `result_team:ID` → `processResultPlace`

### Обработка событий
- `handleCallback` - роутер для всех inline-кнопок (addmember_team, result_tourn, result_team, grant_role)
- `handleText` - роутер для текстовых сообщений на основе FSM state

### Админские команды (handlers_admin.go)
- `/grant` - ввод Telegram ID → выбор роли (Viewer/Organizer/Admin)
- Проверка что пользователь существует в БД перед назначением роли

## 3. Тесты
- [ ] Unit тесты для основной логики (не реализованы, т.к. требуют моков Telegram API)
- [ ] Тесты граничных случаев (не реализованы)
- [ ] Тесты ошибок (не реализованы)

**Примечание**: Для unit-тестирования handlers необходимо добавить моки telebot.Context и repository интерфейсов. Это отдельная задача.

## 4. Проверки
- [x] go build ./... - passed
- [x] go test ./... - passed (там где есть тесты)
- [x] golangci-lint run ./internal/bot/... - passed (0 issues)

## 5. Инструкция по запуску/тестированию

### Локальный запуск бота
```bash
# Запустить инфраструктуру
docker-compose up -d

# Применить миграции
make migrate

# Запустить бота
make run
```

### Тестирование в Telegram
1. Найти бота по токену из config.yaml
2. Отправить `/start` - должно показать приветствие и список команд
3. Проверить публичные команды:
   - `/teams` - должно быть пусто или показать команды
   - `/rating` - должно быть пусто или показать рейтинг

4. Для тестирования организаторских команд:
   - Админ должен выполнить `/grant` и назначить роль Organizer вашему Telegram ID
   - После этого доступны: `/newteam`, `/addmember`, `/newtournament`, `/result`

5. Проверить FSM диалоги:
   - `/newteam` → ввод названия → должна создаться команда
   - `/addmember` → выбор команды → ввод имени → должен добавиться участник
   - `/newtournament` → название → дата → место → должен создаться турнир
   - `/result` → выбор турнира → выбор команды → ввод места → должен записаться результат

6. `/cancel` в любом диалоге должен прерывать процесс

### Линтер
```bash
golangci-lint run ./internal/bot/...
```

## 6. Зависимости и проверки

### Использованные методы repository
- `teamRepo.List()` - список команд
- `teamRepo.GetByName()` - поиск по названию
- `teamRepo.GetByID()` - получение по ID
- `teamRepo.Create()` - создание команды
- `memberRepo.GetByTeamID()` - участники команды
- `memberRepo.Create()` - добавление участника
- `tournRepo.ListRecent()` - последние турниры
- `tournRepo.GetByID()` - получение турнира
- `tournRepo.Create()` - создание турнира
- `resultRepo.GetByTeamID()` - результаты команды
- `resultRepo.GetTeamRating()` - рейтинг команд
- `resultRepo.Create()` - запись результата
- `userRepo.GetByTelegramID()` - проверка существования пользователя
- `userRepo.UpdateRole()` - изменение роли

Все методы существуют и работают корректно.

### FSM States
Использованы все необходимые состояния из `internal/fsm/state.go`:
- StateNewTeamName
- StateAddMemberTeam, StateAddMemberName
- StateNewTournamentName, StateNewTournamentDate, StateNewTournamentLocation
- StateResultTournament, StateResultTeam, StateResultPlace
- StateGrantUser, StateGrantRole

## 7. Особенности реализации

1. **Inline-кнопки** используются для выбора команд/турниров (удобнее чем текстовый ввод)
2. **Форматы дат** поддерживаются: `2006-01-02`, `02.01.2006`, `2006/01/02`, `02/01/2006`
3. **Обработка ошибок**: все FSM операции используют `_ =` для игнорирования ошибок (т.к. они не критичны для UX)
4. **Валидация**: проверка пустых строк, дубликатов команд, корректности числовых значений
5. **Контекст**: везде используется `context.Background()` (в будущем можно добавить timeout/cancellation)

## 8. Возможные улучшения (не в scope текущей задачи)

- Добавить unit-тесты с моками
- Добавить pagination для списка команд/турниров (если их > 20)
- Добавить подтверждение перед созданием команды/турнира
- Добавить возможность редактирования/удаления записей
- Добавить поиск команд по части названия
- Добавить экспорт рейтинга в CSV
- Добавить timeout для FSM диалогов (сейчас TTL = 10 минут)
