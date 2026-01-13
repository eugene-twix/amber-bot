// internal/bot/handlers_org.go
package bot

import (
	"context"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/fsm"
	tele "gopkg.in/telebot.v3"
)

// Validation constants
const (
	maxTeamNameLen       = 100
	maxMemberNameLen     = 100
	maxTournamentNameLen = 100
	maxLocationLen       = 200
	maxPlace             = 1000
	minDateYearsAgo      = 1
	maxDateYearsAhead    = 5
)

// verifyState re-reads FSM state and verifies it matches expected state (race condition protection)
func (b *Bot) verifyState(ctx context.Context, userID int64, expected fsm.State) (*fsm.UserState, error) {
	state, err := b.fsm.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get FSM state: %w", err)
	}
	if state.State != expected {
		return nil, fmt.Errorf("state mismatch: expected %s, got %s", expected, state.State)
	}
	return state, nil
}

// /newteam - создать команду
func (b *Bot) handleNewTeam(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTeamName, fsm.Data{}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}
	return c.Send("Введите название команды:", CancelMenu())
}

func (b *Bot) processNewTeamName(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	if _, err := b.verifyState(ctx, c.Sender().ID, fsm.StateNewTeamName); err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	teamName := strings.TrimSpace(c.Text())

	if teamName == "" {
		return c.Send("Название не может быть пустым. Введите название команды:", CancelMenu())
	}

	// Validate length
	if len(teamName) > maxTeamNameLen {
		return c.Send(fmt.Sprintf("Название слишком длинное (макс %d символов). Введите другое название:", maxTeamNameLen), CancelMenu())
	}

	// Check if team already exists
	_, err := b.teamRepo.GetByName(ctx, teamName)
	if err == nil {
		return c.Send("Команда с таким названием уже существует. Введите другое название:", CancelMenu())
	}

	// Create team
	team := &domain.Team{
		Name:      teamName,
		CreatedBy: c.Sender().ID,
	}

	if err := b.teamRepo.Create(ctx, team); err != nil {
		log.Printf("ERROR: failed to create team: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка при создании команды", MainMenu(user.Role))
	}

	// Сохраняем team_id и спрашиваем о добавлении участников
	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTeamAddMembers, fsm.Data{"team_id": team.ID, "team_name": team.Name}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send(fmt.Sprintf("✅ Команда '%s' создана!", teamName), MainMenu(user.Role))
	}

	buttons := [][]tele.InlineButton{
		{
			{Text: "✅ Да", Data: fmt.Sprintf("newteam_addmembers:%d:yes", team.ID)},
			{Text: "❌ Нет", Data: fmt.Sprintf("newteam_addmembers:%d:no", team.ID)},
		},
	}

	return c.Send(fmt.Sprintf("✅ Команда '%s' создана!\n\nДобавить участников?", teamName), &tele.ReplyMarkup{InlineKeyboard: buttons})
}

// handleNewTeamAddMembersCallback - ответ на "Добавить участников?"
func (b *Bot) handleNewTeamAddMembersCallback(c tele.Context, payload string) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()

	// payload format: "teamID:yes" or "teamID:no"
	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		return c.Send("Ошибка формата данных")
	}

	teamID, _ := strconv.ParseInt(parts[0], 10, 64)
	answer := parts[1]

	if answer == "no" {
		// Спрашиваем о записи результата
		buttons := [][]tele.InlineButton{
			{
				{Text: "✅ Да", Data: fmt.Sprintf("newteam_result:%d:yes", teamID)},
				{Text: "❌ Нет", Data: fmt.Sprintf("newteam_result:%d:no", teamID)},
			},
		}
		return c.Edit("Записать результат турнира для этой команды?", &tele.ReplyMarkup{InlineKeyboard: buttons})
	}

	// Переходим к добавлению участника
	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		user := b.getUser(c)
		return c.Edit("Команда не найдена", MainMenu(user.Role))
	}

	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTeamMoreMember, fsm.Data{"team_id": teamID, "team_name": team.Name}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса")
	}

	return c.Edit(fmt.Sprintf("Введите имя участника для команды '%s':", team.Name))
}

// processNewTeamMemberName - обработка ввода имени участника в цепочке создания команды
func (b *Bot) processNewTeamMemberName(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	state, err := b.verifyState(ctx, c.Sender().ID, fsm.StateNewTeamMoreMember)
	if err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	memberName := strings.TrimSpace(c.Text())

	if memberName == "" {
		return c.Send("Имя не может быть пустым. Введите имя участника:", CancelMenu())
	}

	if len(memberName) > maxMemberNameLen {
		return c.Send(fmt.Sprintf("Имя слишком длинное (макс %d символов):", maxMemberNameLen), CancelMenu())
	}

	teamID := state.Data.GetInt64("team_id")
	teamName := state.Data.GetString("team_name")

	member := &domain.Member{
		Name:   memberName,
		TeamID: teamID,
	}

	if err := b.memberRepo.Create(ctx, member); err != nil {
		log.Printf("ERROR: failed to create member: %v", err)
		return c.Send("Ошибка при добавлении участника", CancelMenu())
	}

	// Спрашиваем "Ещё участника?"
	buttons := [][]tele.InlineButton{
		{
			{Text: "➕ Ещё", Data: fmt.Sprintf("newteam_more:%d:yes", teamID)},
			{Text: "✅ Готово", Data: fmt.Sprintf("newteam_more:%d:no", teamID)},
		},
	}

	return c.Send(fmt.Sprintf("✅ Участник '%s' добавлен в команду '%s'!\n\nДобавить ещё?", memberName, teamName), &tele.ReplyMarkup{InlineKeyboard: buttons})
}

// handleNewTeamMoreCallback - ответ на "Ещё участника?"
func (b *Bot) handleNewTeamMoreCallback(c tele.Context, payload string) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()

	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		return c.Send("Ошибка формата данных")
	}

	teamID, _ := strconv.ParseInt(parts[0], 10, 64)
	answer := parts[1]

	if answer == "no" {
		_ = b.fsm.Clear(ctx, c.Sender().ID)

		// Спрашиваем о записи результата
		buttons := [][]tele.InlineButton{
			{
				{Text: "✅ Да", Data: fmt.Sprintf("newteam_result:%d:yes", teamID)},
				{Text: "❌ Нет", Data: fmt.Sprintf("newteam_result:%d:no", teamID)},
			},
		}
		return c.Edit("Записать результат турнира для этой команды?", &tele.ReplyMarkup{InlineKeyboard: buttons})
	}

	// Продолжаем добавлять участников
	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		user := b.getUser(c)
		return c.Edit("Команда не найдена", MainMenu(user.Role))
	}

	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTeamMoreMember, fsm.Data{"team_id": teamID, "team_name": team.Name}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса")
	}

	return c.Edit(fmt.Sprintf("Введите имя участника для команды '%s':", team.Name))
}

// handleNewTeamResultCallback - ответ на "Записать результат?"
func (b *Bot) handleNewTeamResultCallback(c tele.Context, payload string) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()

	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		return c.Send("Ошибка формата данных")
	}

	teamID, _ := strconv.ParseInt(parts[0], 10, 64)
	answer := parts[1]

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	user := b.getUser(c)

	if answer == "no" {
		_ = c.Edit("Готово!")
		return c.Send("Используйте кнопки меню для дальнейших действий.", MainMenu(user.Role))
	}

	// Показываем список турниров для выбора
	tournaments, err := b.tournRepo.ListRecent(ctx, 10)
	if err != nil || len(tournaments) == 0 {
		return c.Edit("Нет турниров. Сначала создайте турнир.", MainMenu(user.Role))
	}

	// Сохраняем team_id для записи результата
	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateResultTournament, fsm.Data{"team_id": teamID}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Edit("Ошибка сервиса", MainMenu(user.Role))
	}

	var buttons [][]tele.InlineButton
	for _, t := range tournaments {
		buttons = append(buttons, []tele.InlineButton{
			{Text: fmt.Sprintf("%s (%s)", t.Name, t.Date.Format("02.01.2006")), Data: fmt.Sprintf("result_tourn:%d", t.ID)},
		})
	}

	return c.Edit("Выберите турнир:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}

// /addmember - добавить участника
func (b *Bot) handleAddMember(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	teams, err := b.teamRepo.List(ctx)
	if err != nil {
		log.Printf("ERROR: failed to list teams: %v", err)
		return c.Send("Ошибка получения списка команд")
	}
	if len(teams) == 0 {
		return c.Send("Сначала создайте команду через кнопку «➕ Команда»")
	}

	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateAddMemberTeam, fsm.Data{}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}

	// Create inline keyboard with teams
	var buttons [][]tele.InlineButton
	for _, team := range teams {
		buttons = append(buttons, []tele.InlineButton{
			{Text: team.Name, Data: fmt.Sprintf("addmember_team:%d", team.ID)},
		})
	}

	return c.Send("Выберите команду:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}

func (b *Bot) processAddMemberName(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	state, err := b.verifyState(ctx, c.Sender().ID, fsm.StateAddMemberName)
	if err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	memberName := strings.TrimSpace(c.Text())

	if memberName == "" {
		return c.Send("Имя не может быть пустым. Введите имя участника:", CancelMenu())
	}

	// Validate length
	if len(memberName) > maxMemberNameLen {
		return c.Send(fmt.Sprintf("Имя слишком длинное (макс %d символов). Введите другое имя:", maxMemberNameLen), CancelMenu())
	}

	teamID := state.Data.GetInt64("team_id")
	if teamID == 0 {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка: команда не выбрана", MainMenu(user.Role))
	}

	member := &domain.Member{
		Name:   memberName,
		TeamID: teamID,
	}

	if err := b.memberRepo.Create(ctx, member); err != nil {
		log.Printf("ERROR: failed to create member: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка при добавлении участника", MainMenu(user.Role))
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	user := b.getUser(c)
	return c.Send(fmt.Sprintf("✅ Участник '%s' добавлен!", memberName), MainMenu(user.Role))
}

// handleNewTournament - создать турнир
func (b *Bot) handleNewTournament(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateNewTournamentName, fsm.Data{}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}
	return c.Send("Введите название турнира:", CancelMenu())
}

func (b *Bot) processNewTournamentName(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	if _, err := b.verifyState(ctx, c.Sender().ID, fsm.StateNewTournamentName); err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	tournamentName := strings.TrimSpace(c.Text())

	if tournamentName == "" {
		return c.Send("Название не может быть пустым. Введите название турнира:", CancelMenu())
	}

	// Validate length
	if len(tournamentName) > maxTournamentNameLen {
		return c.Send(fmt.Sprintf("Название слишком длинное (макс %d символов). Введите другое название:", maxTournamentNameLen), CancelMenu())
	}

	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateNewTournamentDate, "name", tournamentName); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}
	return c.Send("Введите дату турнира (например: 2026-01-15 или 15.01.2026):", CancelMenu())
}

func (b *Bot) processNewTournamentDate(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	if _, err := b.verifyState(ctx, c.Sender().ID, fsm.StateNewTournamentDate); err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

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
		return c.Send("Неверный формат даты. Используйте формат: 2026-01-15 или 15.01.2026", CancelMenu())
	}

	// Validate date range
	now := time.Now()
	minDate := now.AddDate(-minDateYearsAgo, 0, 0)
	maxDate := now.AddDate(maxDateYearsAhead, 0, 0)
	if date.Before(minDate) || date.After(maxDate) {
		return c.Send(fmt.Sprintf("Дата должна быть от %d года назад до %d лет вперёд",
			minDateYearsAgo, maxDateYearsAhead), CancelMenu())
	}

	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateNewTournamentLocation, "date", date.Format("2006-01-02")); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}
	return c.Send("Введите место проведения (или отправьте '-' если не указываете):", CancelMenu())
}

func (b *Bot) processNewTournamentLocation(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	state, err := b.verifyState(ctx, c.Sender().ID, fsm.StateNewTournamentLocation)
	if err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	location := strings.TrimSpace(c.Text())

	if location == "-" {
		location = ""
	}

	// Validate length
	if len(location) > maxLocationLen {
		return c.Send(fmt.Sprintf("Место проведения слишком длинное (макс %d символов). Введите другое:", maxLocationLen), CancelMenu())
	}

	tournamentName := state.Data.GetString("name")
	dateStr := state.Data.GetString("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка: неверная дата", MainMenu(user.Role))
	}

	tournament := &domain.Tournament{
		Name:      tournamentName,
		Date:      date,
		Location:  location,
		CreatedBy: c.Sender().ID,
	}

	if err := b.tournRepo.Create(ctx, tournament); err != nil {
		log.Printf("ERROR: failed to create tournament: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка при создании турнира", MainMenu(user.Role))
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	user := b.getUser(c)
	return c.Send(fmt.Sprintf("✅ Турнир '%s' создан!", tournamentName), MainMenu(user.Role))
}

// handleResult - записать результат
func (b *Bot) handleResult(c tele.Context) error {
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	tournaments, err := b.tournRepo.ListRecent(ctx, 10)
	if err != nil {
		log.Printf("ERROR: failed to list tournaments: %v", err)
		return c.Send("Ошибка получения списка турниров")
	}
	if len(tournaments) == 0 {
		return c.Send("Сначала создайте турнир через кнопку «🎯 Турнир»")
	}

	if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateResultTournament, fsm.Data{}); err != nil {
		log.Printf("ERROR: failed to set FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}

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
	if err != nil {
		log.Printf("ERROR: failed to list teams: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		return c.Send("Ошибка получения списка команд")
	}
	if len(teams) == 0 {
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

func (b *Bot) processResultPlace(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	state, err := b.verifyState(ctx, c.Sender().ID, fsm.StateResultPlace)
	if err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	placeStr := strings.TrimSpace(c.Text())

	place, parseErr := strconv.Atoi(placeStr)
	if parseErr != nil || place < 1 || place > maxPlace {
		return c.Send(fmt.Sprintf("Введите корректное место (число от 1 до %d):", maxPlace), CancelMenu())
	}

	tournamentID := state.Data.GetInt64("tournament_id")
	teamID := state.Data.GetInt64("team_id")

	if tournamentID == 0 || teamID == 0 {
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка: турнир или команда не выбраны", MainMenu(user.Role))
	}

	result := &domain.Result{
		TeamID:       teamID,
		TournamentID: tournamentID,
		Place:        place,
		RecordedBy:   c.Sender().ID,
	}

	if err := b.resultRepo.Create(ctx, result); err != nil {
		log.Printf("ERROR: failed to create result: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка при сохранении результата", MainMenu(user.Role))
	}

	// Get team and tournament names for confirmation
	team, teamErr := b.teamRepo.GetByID(ctx, teamID)
	if teamErr != nil {
		log.Printf("ERROR: failed to get team by ID: %v", teamErr)
	}
	tournament, tournamentErr := b.tournRepo.GetByID(ctx, tournamentID)
	if tournamentErr != nil {
		log.Printf("ERROR: failed to get tournament by ID: %v", tournamentErr)
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)
	user := b.getUser(c)

	msg := "✅ Результат записан."
	if team != nil && tournament != nil {
		msg = fmt.Sprintf("✅ Результат записан: %s заняла %d место в турнире '%s'",
			team.Name, place, tournament.Name)
	} else if team != nil {
		msg = fmt.Sprintf("✅ Результат записан: %s заняла %d место", team.Name, place)
	}

	return c.Send(msg, MainMenu(user.Role))
}

// handleCallback - обработка inline кнопок
func (b *Bot) handleCallback(c tele.Context) error {
	data := c.Callback().Data

	// Acknowledge callback
	if err := c.Respond(); err != nil {
		return err
	}

	// Handle noop (page number buttons)
	if data == "noop" {
		return nil
	}

	// Parse callback data
	parts := strings.SplitN(data, ":", 2)
	if len(parts) < 2 {
		return c.Send("Ошибка: неверный формат данных")
	}

	action := parts[0]
	payload := parts[1]

	switch action {
	case "team_page":
		page, _ := strconv.Atoi(payload)
		return b.showTeamsPage(c, page, true)
	case "rating_page":
		page, _ := strconv.Atoi(payload)
		return b.showRatingPage(c, page, true)
	case "team_info":
		return b.handleTeamInfoCallback(c, payload)
	case "newteam_addmembers":
		return b.handleNewTeamAddMembersCallback(c, payload)
	case "newteam_more":
		return b.handleNewTeamMoreCallback(c, payload)
	case "newteam_result":
		return b.handleNewTeamResultCallback(c, payload)
	case "addmember_team":
		return b.handleAddMemberTeamCallback(c, payload)
	case "result_tourn":
		return b.handleResultTournamentCallback(c, payload)
	case "result_team":
		return b.handleResultTeamCallback(c, payload)
	case "grant_page":
		page, _ := strconv.Atoi(payload)
		return b.showGrantUsersPage(c, page, true)
	case "grant_user":
		return b.handleGrantUserCallback(c, payload)
	case "grant_role":
		return b.handleGrantRoleCallback(c, payload)
	default:
		return c.Send("Неизвестное действие")
	}
}

func (b *Bot) handleTeamInfoCallback(c tele.Context, payload string) error {
	ctx := context.Background()
	teamID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID команды")
	}

	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		log.Printf("ERROR: failed to get team by ID: %v", err)
		return c.Send("Команда не найдена")
	}

	members, _ := b.memberRepo.GetByTeamID(ctx, team.ID)
	results, _ := b.resultRepo.GetByTeamID(ctx, team.ID)

	var sb strings.Builder

	// Заголовок команды
	sb.WriteString(fmt.Sprintf("<b>📋 %s</b>\n", html.EscapeString(team.Name)))
	sb.WriteString(Separator + "\n\n")

	// Участники
	sb.WriteString(fmt.Sprintf("<b>👥 Участники (%d)</b>\n", len(members)))
	if len(members) == 0 {
		sb.WriteString("  <i>нет участников</i>\n")
	} else {
		for i, m := range members {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, html.EscapeString(m.Name)))
		}
	}

	// Результаты
	sb.WriteString(fmt.Sprintf("\n<b>🏆 Результаты (%d)</b>\n", len(results)))
	if len(results) == 0 {
		sb.WriteString("  <i>нет результатов</i>\n")
	} else {
		for _, r := range results {
			tournament, _ := b.tournRepo.GetByID(ctx, r.TournamentID)
			if tournament != nil {
				medal := "   "
				switch r.Place {
				case 1:
					medal = "🥇"
				case 2:
					medal = "🥈"
				case 3:
					medal = "🥉"
				}
				sb.WriteString(fmt.Sprintf("  %s %s — <code>%d место</code>\n", medal, html.EscapeString(tournament.Name), r.Place))
			}
		}
	}

	// Статистика
	if len(results) > 0 {
		wins := 0
		totalPlace := 0
		for _, r := range results {
			if r.Place == 1 {
				wins++
			}
			totalPlace += r.Place
		}
		avgPlace := float64(totalPlace) / float64(len(results))

		sb.WriteString(fmt.Sprintf("\n<b>📊 Статистика</b>\n"))
		sb.WriteString(fmt.Sprintf("  Игр: <code>%d</code> | Побед: <code>%d</code> | Ср: <code>%.1f</code>\n", len(results), wins, avgPlace))
	}

	return c.Send(sb.String(), tele.ModeHTML)
}

func (b *Bot) handleAddMemberTeamCallback(c tele.Context, payload string) error {
	// Check permissions
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	teamID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID команды")
	}

	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		log.Printf("ERROR: failed to get team by ID: %v", err)
		return c.Send("Ошибка: команда не найдена")
	}

	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateAddMemberName, "team_id", teamID); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}
	return c.Send(fmt.Sprintf("Команда '%s' выбрана. Введите имя участника:", team.Name), CancelMenu())
}

func (b *Bot) handleResultTournamentCallback(c tele.Context, payload string) error {
	// Check permissions
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	tournamentID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID турнира")
	}

	_, err = b.tournRepo.GetByID(ctx, tournamentID)
	if err != nil {
		log.Printf("ERROR: failed to get tournament by ID: %v", err)
		return c.Send("Ошибка: турнир не найден")
	}

	// Проверяем, есть ли уже team_id (из flow создания команды)
	state, _ := b.fsm.Get(ctx, c.Sender().ID)
	if state != nil && state.Data.GetInt64("team_id") != 0 {
		// team_id уже есть — сразу к вводу места
		if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateResultPlace, "tournament_id", tournamentID); err != nil {
			log.Printf("ERROR: failed to update FSM state: %v", err)
			return c.Send("Ошибка сервиса. Попробуйте позже.")
		}
		return c.Edit("Введите место (число):")
	}

	// team_id нет — показываем выбор команды
	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateResultTeam, "tournament_id", tournamentID); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}

	return b.processResultTeam(c, nil)
}

func (b *Bot) handleResultTeamCallback(c tele.Context, payload string) error {
	// Check permissions
	if !b.requireOrganizer(c) {
		return nil
	}

	ctx := context.Background()
	teamID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send("Ошибка: неверный ID команды")
	}

	team, err := b.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		log.Printf("ERROR: failed to get team by ID: %v", err)
		return c.Send("Ошибка: команда не найдена")
	}

	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateResultPlace, "team_id", teamID); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}

	return c.Send(fmt.Sprintf("Команда '%s' выбрана. Введите место (число):", team.Name), CancelMenu())
}

func (b *Bot) handleGrantRoleCallback(c tele.Context, payload string) error {
	// Check permissions - only admin can grant roles
	if !b.requireAdmin(c) {
		return nil
	}

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
		log.Printf("ERROR: failed to update user role: %v", err)
		_ = b.fsm.Clear(ctx, c.Sender().ID)
		user := b.getUser(c)
		return c.Send("Ошибка при изменении роли", MainMenu(user.Role))
	}

	_ = b.fsm.Clear(ctx, c.Sender().ID)

	// Edit the message to remove buttons and show result
	if err := c.Edit(fmt.Sprintf("✅ Роль '%s' назначена пользователю %d", role, userID)); err != nil {
		log.Printf("WARN: failed to edit message: %v", err)
	}

	user := b.getUser(c)
	return c.Send("Готово!", MainMenu(user.Role))
}
