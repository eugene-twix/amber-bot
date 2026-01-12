// internal/bot/handlers_org.go
package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/fsm"
	tele "gopkg.in/telebot.v3"
)

// /newteam - создать команду
func (b *Bot) handleNewTeam(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	_ = b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTeamName, fsm.Data{})
	return c.Send("Введите название команды:")
}

func (b *Bot) processNewTeamName(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	teamName := strings.TrimSpace(c.Text())

	if teamName == "" {
		return c.Send("Название не может быть пустым. Введите название команды:")
	}

	// Check if team already exists
	_, err := b.teamRepo.GetByName(ctx, teamName)
	if err == nil {
		return c.Send("Команда с таким названием уже существует. Введите другое название:")
	}

	// Create team
	team := &domain.Team{
		Name:      teamName,
		CreatedBy: c.Sender().ID,
	}

	if err := b.teamRepo.Create(ctx, team); err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка при создании команды")
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send(fmt.Sprintf("Команда '%s' создана!", teamName))
}

// /addmember - добавить участника
func (b *Bot) handleAddMember(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	teams, err := b.teamRepo.List(ctx)
	if err != nil || len(teams) == 0 {
		return c.Send("Сначала создайте команду через /newteam")
	}

	_ = b.fsm.Set(ctx, c.Sender().ID, fsm.StateAddMemberTeam, fsm.Data{})

	// Create inline keyboard with teams
	var buttons [][]tele.InlineButton
	for _, team := range teams {
		buttons = append(buttons, []tele.InlineButton{
			{Text: team.Name, Data: fmt.Sprintf("addmember_team:%d", team.ID)},
		})
	}

	return c.Send("Выберите команду:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}

func (b *Bot) processAddMemberName(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	memberName := strings.TrimSpace(c.Text())

	if memberName == "" {
		return c.Send("Имя не может быть пустым. Введите имя участника:")
	}

	teamID := state.Data.GetInt64("team_id")
	if teamID == 0 {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка: команда не выбрана")
	}

	member := &domain.Member{
		Name:   memberName,
		TeamID: teamID,
	}

	if err := b.memberRepo.Create(ctx, member); err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка при добавлении участника")
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send(fmt.Sprintf("Участник '%s' добавлен!", memberName))
}

// /newtournament - создать турнир
func (b *Bot) handleNewTournament(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	_ = b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTournamentName, fsm.Data{})
	return c.Send("Введите название турнира:")
}

func (b *Bot) processNewTournamentName(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	tournamentName := strings.TrimSpace(c.Text())

	if tournamentName == "" {
		return c.Send("Название не может быть пустым. Введите название турнира:")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateNewTournamentDate, "name", tournamentName)
	return c.Send("Введите дату турнира (например: 2026-01-15 или 15.01.2026):")
}

func (b *Bot) processNewTournamentDate(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	dateStr := strings.TrimSpace(c.Text())

	// Try to parse date in different formats
	var date time.Time
	var err error

	formats := []string{
		"2006-01-02",
		"02.01.2006",
		"2006/01/02",
		"02/01/2006",
	}

	for _, format := range formats {
		date, err = time.Parse(format, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return c.Send("Неверный формат даты. Используйте формат: 2026-01-15 или 15.01.2026")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateNewTournamentLocation, "date", date.Format("2006-01-02"))
	return c.Send("Введите место проведения (или отправьте '-' если не указываете):")
}

func (b *Bot) processNewTournamentLocation(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	location := strings.TrimSpace(c.Text())

	if location == "-" {
		location = ""
	}

	tournamentName := state.Data.GetString("name")
	dateStr := state.Data.GetString("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка: неверная дата")
	}

	tournament := &domain.Tournament{
		Name:      tournamentName,
		Date:      date,
		Location:  location,
		CreatedBy: c.Sender().ID,
	}

	if err := b.tournRepo.Create(ctx, tournament); err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка при создании турнира")
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send(fmt.Sprintf("Турнир '%s' создан!", tournamentName))
}

// /result - записать результат
func (b *Bot) handleResult(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	tournaments, err := b.tournRepo.ListRecent(ctx, 10)
	if err != nil || len(tournaments) == 0 {
		return c.Send("Сначала создайте турнир через /newtournament")
	}

	_ = b.fsm.Set(ctx, c.Sender().ID, fsm.StateResultTournament, fsm.Data{})

	// Create inline keyboard with tournaments
	var buttons [][]tele.InlineButton
	for _, t := range tournaments {
		dateStr := t.Date.Format("02.01.2006")
		buttonText := fmt.Sprintf("%s (%s)", t.Name, dateStr)
		buttons = append(buttons, []tele.InlineButton{
			{Text: buttonText, Data: fmt.Sprintf("result_tourn:%d", t.ID)},
		})
	}

	return c.Send("Выберите турнир:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}

func (b *Bot) processResultTeam(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	teams, err := b.teamRepo.List(ctx)
	if err != nil || len(teams) == 0 {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Нет доступных команд")
	}

	// Create inline keyboard with teams
	var buttons [][]tele.InlineButton
	for _, team := range teams {
		buttons = append(buttons, []tele.InlineButton{
			{Text: team.Name, Data: fmt.Sprintf("result_team:%d", team.ID)},
		})
	}

	return c.Send("Выберите команду:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}

func (b *Bot) processResultPlace(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	placeStr := strings.TrimSpace(c.Text())

	place, err := strconv.Atoi(placeStr)
	if err != nil || place < 1 {
		return c.Send("Введите корректное место (число от 1 и выше):")
	}

	tournamentID := state.Data.GetInt64("tournament_id")
	teamID := state.Data.GetInt64("team_id")

	if tournamentID == 0 || teamID == 0 {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка: турнир или команда не выбраны")
	}

	result := &domain.Result{
		TeamID:       teamID,
		TournamentID: tournamentID,
		Place:        place,
		RecordedBy:   c.Sender().ID,
	}

	if err := b.resultRepo.Create(ctx, result); err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка при сохранении результата")
	}

	// Get team and tournament names for confirmation
	team, _ := b.teamRepo.GetByID(ctx, teamID)
	tournament, _ := b.tournRepo.GetByID(ctx, tournamentID)

	_ = b.fsm.Clear(ctx, c.Sender().ID)

	msg := fmt.Sprintf("Результат записан: %s заняла %d место", team.Name, place)
	if tournament != nil {
		msg = fmt.Sprintf("Результат записан: %s заняла %d место в турнире '%s'",
			team.Name, place, tournament.Name)
	}

	return c.Send(msg)
}

// handleCallback - обработка inline кнопок
func (b *Bot) handleCallback(c tele.Context) error {
	data := c.Callback().Data

	// Acknowledge callback
	if err := c.Respond(); err != nil {
		return err
	}

	// Parse callback data
	parts := strings.SplitN(data, ":", 2)
	if len(parts) < 2 {
		return c.Send("Ошибка: неверный формат данных")
	}

	action := parts[0]
	payload := parts[1]

	switch action {
	case "addmember_team":
		return b.handleAddMemberTeamCallback(c, payload)
	case "result_tourn":
		return b.handleResultTournamentCallback(c, payload)
	case "result_team":
		return b.handleResultTeamCallback(c, payload)
	case "grant_role":
		return b.handleGrantRoleCallback(c, payload)
	default:
		return c.Send("Неизвестное действие")
	}
}

func (b *Bot) handleAddMemberTeamCallback(c tele.Context, payload string) error {
	ctx := context.Background()
	teamID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID команды")
	}

	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return c.Send("Ошибка: команда не найдена")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateAddMemberName, "team_id", teamID)
	return c.Send(fmt.Sprintf("Команда '%s' выбрана. Введите имя участника:", team.Name))
}

func (b *Bot) handleResultTournamentCallback(c tele.Context, payload string) error {
	ctx := context.Background()
	tournamentID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID турнира")
	}

	_, err = b.tournRepo.GetByID(ctx, tournamentID)
	if err != nil {
		return c.Send("Ошибка: турнир не найден")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateResultTeam, "tournament_id", tournamentID)

	// Show teams selection
	return b.processResultTeam(c, nil)
}

func (b *Bot) handleResultTeamCallback(c tele.Context, payload string) error {
	ctx := context.Background()
	teamID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID команды")
	}

	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return c.Send("Ошибка: команда не найдена")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateResultPlace, "team_id", teamID)

	return c.Send(fmt.Sprintf("Команда '%s' выбрана. Введите место (число):", team.Name))
}

func (b *Bot) handleGrantRoleCallback(c tele.Context, payload string) error {
	ctx := context.Background()

	// Parse payload: userID:role
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		return c.Send("Ошибка: неверный формат данных")
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID пользователя")
	}

	role := domain.Role(parts[1])

	// Update role
	if err := b.userRepo.UpdateRole(ctx, userID, role); err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка при изменении роли")
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send(fmt.Sprintf("Роль '%s' назначена пользователю %d", role, userID))
}

// handleText - обработка текстовых сообщений (для FSM)
func (b *Bot) handleText(c tele.Context) error {
	ctx := context.Background()
	state, err := b.fsm.Get(ctx, c.Sender().ID)
	if err != nil || state.State == fsm.StateNone {
		return c.Send("Используйте команды для работы с ботом. Наберите /start для списка команд.")
	}

	// Route to appropriate handler based on current state
	switch state.State {
	case fsm.StateNewTeamName:
		return b.processNewTeamName(c, state)
	case fsm.StateAddMemberName:
		return b.processAddMemberName(c, state)
	case fsm.StateNewTournamentName:
		return b.processNewTournamentName(c, state)
	case fsm.StateNewTournamentDate:
		return b.processNewTournamentDate(c, state)
	case fsm.StateNewTournamentLocation:
		return b.processNewTournamentLocation(c, state)
	case fsm.StateResultPlace:
		return b.processResultPlace(c, state)
	case fsm.StateGrantUser:
		return b.processGrantUser(c, state)
	default:
		return c.Send("Неизвестное состояние. Используйте /cancel для отмены.")
	}
}
